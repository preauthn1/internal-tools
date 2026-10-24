package network

import (
	"bytes"
	"errors"
	"fmt"
)

// NetServer 網絡服務器
type NetServer interface {
	Start(evtchar *ServerNetEventChan) error                  // 啟動服務器
	Stop()                                                    // 停止服務器
	ServerID() string                                         // 服務器標識
	RangeConn(f func(conn *ConnInfo))                         // 遍歷所有連接
	GetConnByID(connid string) *ConnInfo                      // 獲取連接
	CloseConn(conninfo *ConnInfo)                             // 關閉客戶端連接
	CloseConnByConnID(connid string)                          // 根據連接ID關閉客戶端連接
	SendMessage(conninfo *ConnInfo, data []byte) (int, error) // 發送消息
	SendAllMessage(data []byte)                               // 廣播消息
	ConnNum() int                                             // 當前連接數
	ServerType() int                                          // 服務器類型
	String() string                                           // 服務器描述信息
}

// ServerCfg 服務器配置
type ServerCfg struct {
	Addr       string // 偵聽地址
	MaxConnNum int    // 最大連接數
	ServerID   string // 服務器標識

	OtherCfg interface{}
}

// CheckCfg 檢查服務器參數
func (cfg *ServerCfg) CheckCfg() error {
	if cfg.Addr == "" {
		return errors.New("服務器偵聽地址不能為空")
	}
	if cfg.MaxConnNum <= 0 {
		return errors.New("最大連接數必須為大於0的數字")
	}

	if cfg.ServerID == "" {
		return errors.New("服務器標識不能為空")
	}

	return nil
}

// NewServerCfg 建立服務器配置
func NewServerCfg(addr string, maxconnNum int, serverid string) *ServerCfg {
	cfg := new(ServerCfg)
	cfg.Addr = addr
	cfg.MaxConnNum = maxconnNum
	cfg.ServerID = serverid

	return cfg
}

// NetServerGroup 服務器組
type NetServerGroup struct {
	innerServers map[string]NetServer
	serverID     string
}

// NewNetServerGroup 服務器組
func NewNetServerGroup(sid string) *NetServerGroup {
	svr := new(NetServerGroup)
	svr.serverID = sid
	svr.innerServers = make(map[string]NetServer)

	return svr
}

// AddServer 增加Server
func (svr *NetServerGroup) AddServer(s NetServer) error {
	sid := s.ServerID()

	if _, ok := svr.innerServers[sid]; ok {
		return errors.New("服務器已經存在")
	}

	svr.innerServers[sid] = s

	return nil
}

// Start 啟動服務器
func (svr *NetServerGroup) Start(evtchan *ServerNetEventChan) error {
	errstr := ""

	for _, s := range svr.innerServers {
		err := s.Start(evtchan)

		if err != nil {
			errstr += err.Error()
			errstr += "\r\n"
		}
	}

	if errstr == "" {
		return nil
	}

	return errors.New(errstr)
}

// Stop 停止服務器
func (svr *NetServerGroup) Stop() {
	for _, s := range svr.innerServers {
		s.Stop()
	}
}

// ServerID 服務器標識
func (svr *NetServerGroup) ServerID() string {
	return svr.serverID
}

// RangeConn 遍歷所有連接
func (svr *NetServerGroup) RangeConn(f func(conn *ConnInfo)) {
	for _, s := range svr.innerServers {
		s.RangeConn(f)
	}
}

// GetConnByID 獲取連接
func (svr *NetServerGroup) GetConnByID(connid string) *ConnInfo {
	var ci *ConnInfo

	for _, s := range svr.innerServers {
		ci = s.GetConnByID(connid)

		if ci != nil {
			return ci
		}
	}

	return nil
}

// CloseConn 關閉客戶端連接
func (svr *NetServerGroup) CloseConn(conninfo *ConnInfo) {
	sid := conninfo.ServerID

	s := svr.innerServers[sid]
	if s == nil {
		return
	}

	if conninfo != nil {
		s.CloseConn(conninfo)
	}

	return
}

// CloseConnByConnID 根據連接ID關閉客戶端連接
func (svr *NetServerGroup) CloseConnByConnID(connid string) {
	clt := svr.GetConnByID(connid)

	if clt != nil {
		svr.CloseConn(clt)
	}
}

// SendMessage 發送消息
func (svr *NetServerGroup) SendMessage(conninfo *ConnInfo, data []byte) (int, error) {
	sid := conninfo.ServerID

	s := svr.innerServers[sid]
	if s == nil {
		return 0, errors.New("serverid不存在")
	}

	return s.SendMessage(conninfo, data)
}

// SendAllMessage 廣播消息
func (svr *NetServerGroup) SendAllMessage(data []byte) {
	for _, s := range svr.innerServers {
		s.SendAllMessage(data)
	}
}

// ConnNum 當前連接數
func (svr *NetServerGroup) ConnNum() int {
	num := 0

	for _, s := range svr.innerServers {
		num += s.ConnNum()
	}

	return num
}

// RangeServer 遍歷所有服務器
func (svr *NetServerGroup) RangeServer(f func(s NetServer)) {
	for _, s := range svr.innerServers {
		f(s)
	}
}

// GetServerByID 取得Server
func (svr *NetServerGroup) GetServerByID(id string) NetServer {
	netsvr, ok := svr.innerServers[id]

	if !ok {
		return nil
	}

	return netsvr
}

// ServerType 服務器類型
func (*NetServerGroup) ServerType() int {
	return Mixed
}

// ClearServer 清除服務器
func (svr *NetServerGroup) ClearServer() {
	svr.innerServers = make(map[string]NetServer)
}

// String 服務器描述信息
func (svr *NetServerGroup) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("server count:%d,conn count:%d\r\n", len(svr.innerServers), svr.ConnNum()))

	// 寫服務器組信息
	for _, s := range svr.innerServers {
		buf.WriteString(s.String() + "\r\n")
	}

	return buf.String()
}
