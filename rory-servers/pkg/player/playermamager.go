package player

import (
	"errors"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// Manager 用户管理器
type Manager struct {
	users     map[string]*Player
	msgserver *msgserver.MsgServer
}

// NewPlayerManager 創建用戶管理器對象
func NewPlayerManager(svr *msgserver.MsgServer) *Manager {
	pm := new(Manager)
	pm.users = make(map[string]*Player)
	pm.msgserver = svr

	PlayerCmdP = NewPlayerCmd(pm, svr)
	command.AddCommand(PlayerCmdP)

	return pm
}

// SetPlayer 添加用戶
func (pm *Manager) SetPlayer(ur *Player) {

	pm.users[ur.Conn.ConnID()] = ur
}

// RemovePlayer 刪除用戶
func (pm *Manager) RemovePlayer(connid string) {

	delete(pm.users, connid)
}

// GetPlayerByConnID 根據連接獲取用戶
func (pm *Manager) GetPlayerByConnID(connid string) *Player {

	return pm.users[connid]
}

// GetPlayerByUserID 根據用戶ID獲取用戶信息
func (pm *Manager) GetPlayerByUserID(userid int32) *Player {

	for _, u := range pm.users {
		if u.UserID == userid {
			return u
		}
	}

	return nil
}

// LogonPlayer 已登錄用戶
func (pm *Manager) LogonPlayer() []*Player {

	urs := make([]*Player, 0)
	for _, u := range pm.users {
		if !u.IsGuest() {
			urs = append(urs, u)
		}
	}

	return urs
}

// UnLogonPlayer 未登錄用戶
func (pm *Manager) UnLogonPlayer() []*Player {

	urs := make([]*Player, 0)
	for _, u := range pm.users {
		if u.IsGuest() {
			urs = append(urs, u)
		}
	}

	return urs
}

// PlayerCount 用戶數量
func (pm *Manager) PlayerCount() int {

	return len(pm.users)
}

// GuestCount 未登錄用戶數
func (pm *Manager) GuestCount() int {
	return len(pm.UnLogonPlayer())
}

// UserCount 已登錄用戶數
func (pm *Manager) UserCount() int {
	return len(pm.LogonPlayer())
}

// Range 循例
func (pm *Manager) Range(f func(p *Player)) {

	for _, pr := range pm.users {
		f(pr)
	}
}

// SetPlayerData 設定玩家資料
func (pm *Manager) SetPlayerData(connid string, setfunc func(pr *Player)) error {

	pr, ok := pm.users[connid]

	if !ok {
		return errors.New("连接不存在")
	}
	setfunc(pr)

	return nil

}

// SetPlayerOffLine 設置用戶斷線
func (pm *Manager) SetPlayerOffLine(connid string) {

	pr, ok := pm.users[connid]

	if !ok {
		return
	}

	if pr != nil {
		pr.IsOnline = false
	}
}

// SetPlayerRecover 設置用戶斷線重回
func (pm *Manager) SetPlayerRecover(newconnid string, userid int32) {

	var exid string
	var expr *Player

	exid = ""
	expr = nil
	for id, pr := range pm.users {
		if pr.UserID == userid && pr.IsOffLine() {
			exid = id
			expr = pr
			break
		}
	}

	if exid != "" && expr != nil {
		newplayer, ok := pm.users[newconnid]
		if ok {
			newplayer.UserID = expr.UserID
			newplayer.Cer = expr.Cer
			newplayer.HDCode = expr.HDCode
			newplayer.HDType = expr.HDType
			newplayer.UserData = expr.UserData
			newplayer.Conn.ConnTime = expr.Conn.ConnTime

			//fmt.Printf("set newplayer :%v\r\n", newplayer)
		}

		//删除旧的用户信息
		delete(pm.users, exid)

	}
}

// ClearPlayerInfo 清除用戶信息
func (pm *Manager) ClearPlayerInfo(connid string) {

	expr, ok := pm.users[connid]
	if ok {
		expr.UserData = nil
		expr.UserID = GuestUserID
		expr.Cer = GuestCerData
		expr.HDType = 0
		expr.HDCode = ""
	}
}

// SendMessage 發送訊息
func (pm *Manager) SendMessage(clt *network.ConnInfo, msg *msg.Message) {
	pm.msgserver.SendMessage(clt, msg)
}
