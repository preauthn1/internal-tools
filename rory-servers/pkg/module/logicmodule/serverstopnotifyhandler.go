package logicmodule

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
)

// ServerStatusProvider 服務器狀態
type ServerStatusProvider interface {
	IsServerRuning() bool                           // 服務器是否運行狀態
	GetServerStatus() int32                         // 獲取服務器狀態
	GetStopServerTipMsg() string                    // 獲取停服提示消息
	IsWhiteIP(ip string) bool                       // 是否IP白名單
	SetServerStatus(statusid int32, msg string)     // 設置服務器狀態
	SetMasterSlave(master HostInfo, slave HostInfo) // 設置Master和Slave
}

// ServerStopNotifyHandler 停服通知處理器
type ServerStopNotifyHandler struct {
	serverstopInfo *netproto.ServerStop
	whiteIPList    *netproto.ServerStopWhiteIPList
}

// HostInfo Host Info
type HostInfo struct {
	GameID   int32
	ServerID int32
}

// NewServerStopNotifyHandler 創建停服通知處理器對象
func NewServerStopNotifyHandler() *ServerStopNotifyHandler {
	ssh := new(ServerStopNotifyHandler)
	ssh.serverstopInfo = nil
	ssh.whiteIPList = nil

	return ssh
}

// IsHook 是否處理
func (ssh *ServerStopNotifyHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) &&
		(sclsid == int32(netproto.NotifyServerClassID_ServerStopID) ||
			sclsid == int32(netproto.NotifyServerClassID_ServerStopWhiteIPListID) ||
			sclsid == int32(netproto.NotifyServerClassID_ServerHostAndBackUp))
}

// OnNotifyMessage 響應Notify消息
func (ssh *ServerStopNotifyHandler) OnNotifyMessage(m *msg.Message, ctx *module.MouduleCtx, player *player.Manager) {
	// 停服
	if m.SClassID == int32(netproto.NotifyServerClassID_ServerStopID) {
		mmbyte := m.MsgData.([]byte)

		mm := &netproto.ServerStop{}
		err := proto.Unmarshal(mmbyte, mm)

		if err != nil {
			ctx.Error("解析停服消息錯誤,%v", err)
			return
		}

		ssh.serverstopInfo = mm

		example.invalid("收到停服消息,%v", ssh.serverstopInfo)
	}

	// ip白名單
	if m.SClassID == int32(netproto.NotifyServerClassID_ServerStopWhiteIPListID) {

		mmbyte := m.MsgData.([]byte)

		mm := &netproto.ServerStopWhiteIPList{}
		err := proto.Unmarshal(mmbyte, mm)

		if err != nil {
			ctx.Error("解析停服IP白名單消息錯誤,%v", err)
			return
		}

		ssh.whiteIPList = mm
		example.invalid("收到停服IP白名單消息,%v", ssh.whiteIPList)

	}

}

// IsServerRuning Server 是否執行中
func (ssh *ServerStopNotifyHandler) IsServerRuning() bool {
	if ssh.serverstopInfo == nil {
		return true
	}
	if ssh.serverstopInfo.GetStatus() == conf.ServerStatusRuning || ssh.serverstopInfo.GetStatus() == conf.ServerStatusRecover {
		return true
	}
	return false
}

// GetServerStatus 取得Server狀態
func (ssh *ServerStopNotifyHandler) GetServerStatus() int32 {
	if ssh.serverstopInfo == nil {
		return conf.ServerStatusRuning
	}

	return ssh.serverstopInfo.GetStatus()
}

// GetStopServerTipMsg 取得停服踢人訊息
func (ssh *ServerStopNotifyHandler) GetStopServerTipMsg() string {
	if ssh.serverstopInfo == nil {
		return ""
	}
	return ssh.serverstopInfo.GetMessage()
}

// IsWhiteIP 是否為白名單
func (ssh *ServerStopNotifyHandler) IsWhiteIP(ip string) bool {
	if ssh.whiteIPList == nil {
		return false
	}

	for _, sip := range ssh.whiteIPList.IP {
		if ip == sip {
			return true
		}
	}

	return false
}

// SetServerStatus 設定Server狀態
func (ssh *ServerStopNotifyHandler) SetServerStatus(statusid int32, msg string) {
	if ssh.serverstopInfo == nil {
		ssh.serverstopInfo = new(netproto.ServerStop)
	}

	ssh.serverstopInfo.Status = proto.Int32(statusid)
	ssh.serverstopInfo.Message = proto.String(msg)
}

// SetMasterSlave 設定奴隸
func (ssh *ServerStopNotifyHandler) SetMasterSlave(master HostInfo, slave HostInfo) {
	mlog.System("SetMasterSlave")
}
