package service

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"example.invalid/spf13/viper"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// Service 向後端取資料邏輯
type Service interface {
	GetUserByID(db *sql.DB, uid int32) (*PlayerInfo, error)
	GetServerAddr(db *sql.DB, gid int32, sid int32) (string, error)
	GetInfo(db *sql.DB, uid, gid, sid int32) (*Info, error)
}

// DDService for DD
type DDService struct{}

// GetServerAddr 查詢 GameServerList 中對應的 server id 聽的位置
func (*DDService) GetServerAddr(db *sql.DB, gid, sid int32) (string, error) {
	row := db.QueryRow("SELECT ServerName, ClientAddr FROM ServerInfo.GameServerList WHERE GameID = ? AND ServerID = ?", gid, sid)

	var sname, saddr string

	if err := row.Scan(&sname, &saddr); err != nil {
		return "", err
	}

	// 拆解port
	strs := strings.Split(saddr, ":")
	if len(strs) != 2 {
		return "", fmt.Errorf("Could not find server port, %s", saddr)
	}
	port := strs[1]

	switch viper.GetString("server.conntype") {
	case "TCP":
		saddr = fmt.Sprintf("%s:%s", viper.GetString("server.dnsprefix"), port)
	case "Websocket":
		saddr = fmt.Sprintf("%s/ws/%s", viper.GetString("server.dnsprefix"), port)
	}

	mlog.Debug("Query ServerAddr %d:%d, Name: %s, Addr: %s", gid, sid, sname, saddr)

	return saddr, nil
}

// GetUserByID get user
func (*DDService) GetUserByID(db *sql.DB, uid int32) (*PlayerInfo, error) {
	row := db.QueryRow("SELECT a.HDType, a.HDCode, b.CerCode FROM Users.UserInfo a INNER JOIN Game.Certification b ON a.UserID = b.UserID WHERE a.UserID = ?", uid)

	var hdt int32
	var hdc, cer string

	if err := row.Scan(&hdt, &hdc, &cer); err != nil {
		return nil, err
	}

	u := new(PlayerInfo)
	u.PlayerID = uid
	// u.UserHDType = hdt
	// u.UserHDCode = hdc
	u.AccessToken = [REDACTED]

	return u, nil
}

// GetInfo get info
func (svc *DDService) GetInfo(db *sql.DB, uid, gid, sid int32) (*Info, error) {
	info := new(Info)
	u, err := svc.GetUserByID(db, uid)
	if err != nil {
		return nil, fmt.Errorf("Get User Error %v", err)
	}
	info.PlayerInfo = *u

	srvAddr, err := svc.GetServerAddr(db, gid, sid)
	if err != nil {
		return nil, fmt.Errorf("Get Server Error %v", err)
	}

	info.ServerInfo.GameID = gid
	info.ServerInfo.ServerID = sid
	info.ServerInfo.ServerConnection = srvAddr
	info.ServerInfo.ConnType = viper.GetString("server.conntype")

	return info, nil
}

// PGSService for target
type PGSService struct{}

// GetServerAddr 查詢 GameServerList 中對應的 server id 聽的位置
func (*PGSService) GetServerAddr(db *sql.DB, gid, sid int32) (string, error) {
	row := db.QueryRow("SELECT ServerName, ClientAddr FROM ServerList WHERE GameID = ? AND ServerID = ?", gid, sid)

	var _n, _addr sql.NullString
	var sname, saddr string

	if err := row.Scan(&_n, &_addr); err != nil {
		return "", err
	}

	if _n.Valid {
		sname = _n.String
	}

	if _addr.Valid {
		saddr = _addr.String
	}

	fmt.Println(_n, _addr, sname, saddr)

	// 拆解port
	strs := strings.Split(saddr, ":")
	if len(strs) != 2 {
		return "", fmt.Errorf("Could not find server port, %s", saddr)
	}
	port := strs[1]

	switch viper.GetString("server.conntype") {
	case "TCP":
		saddr = fmt.Sprintf("%s:%s", viper.GetString("server.dnsprefix"), port)
	case "Websocket":
		saddr = fmt.Sprintf("%s/ws/%s", viper.GetString("server.dnsprefix"), port)
	}

	mlog.Debug("Query ServerAddr %d:%d, Name: %s, Addr: %s", gid, sid, sname, saddr)

	return saddr, nil
}

// GetUserByID get user
func (*PGSService) GetUserByID(db *sql.DB, uid int32) (*PlayerInfo, error) {
	// row := db.QueryRow("SELECT PlayerName FROM Player WHERE PlayerID = ?", uid)

	// var hdt int32
	// var hdc, cer string

	// if err := row.Scan(&hdt, &hdc, &cer); err != nil {
	// 	return nil, err
	// }

	u := new(PlayerInfo)
	u.PlayerID = uid
	// u.UserHDType = hdt
	// u.UserHDCode = hdc
	// u.UserCer = cer

	return u, nil
}

// GetInfo get info
func (svc *PGSService) GetInfo(db *sql.DB, uid, gid, sid int32) (*Info, error) {

	var mid int32
	var uname string
	var ts int64 = time.Now().Unix()
	var signkey string
	var lang string = "EN-US"
	var session string
	var room int32
	var saddr string

	info := new(Info)

	{
		row := db.QueryRow("SELECT MerchantID, PlayerName FROM Player WHERE PlayerID = ?", uid)

		if err := row.Scan(&mid, &uname); err != nil {
			return nil, err
		}
	}
	{
		var _addr sql.NullString
		row := db.QueryRow("SELECT FlagLevel, ClientAddr FROM ServerList WHERE GameID = ? AND ServerID = ?", gid, sid)

		if err := row.Scan(&room, &_addr); err != nil {
			return nil, err
		}

		if _addr.Valid {
			saddr = _addr.String
		}
	}

	session = fmt.Sprintf("%x%d%d%s", ts, mid, gid, uname)
	signkey = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_%d_%d", uname, ts, mid))))

	u := url.URL{
		Scheme: "http",
		Host:   viper.GetString("target.api.public"),
		Path:   "/v1/game/launch",
	}
	q := url.Values{}
	q.Set("merchantID", fmt.Sprintf("%d", mid))
	q.Set("playerName", uname)
	q.Set("playerValidateSession", session)
	q.Set("lang", lang)
	q.Set("gameID", fmt.Sprintf("%d", gid))
	q.Set("ts", fmt.Sprintf("%d", ts))
	q.Set("signKey", signkey)
	q.Set("roomLevelFlagID", fmt.Sprintf("%d", room))

	resp, err := http.Get(u.String() + "?" + q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	mp := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&mp); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Get Access Token Error , %v", mp["message"])
	}

	info.PlayerInfo.PlayerID = uid
	info.PlayerInfo.AccessToken = [REDACTED]"data"].(map[string]interface{})["access_Token"].(string)

	info.ServerInfo.GameID = gid
	info.ServerInfo.ServerID = sid
	info.ServerInfo.ServerConnection = saddr
	info.ServerInfo.ConnType = viper.GetString("server.conntype")

	return info, nil
}
