package rpchandler

import (
	"fmt"

	proto "example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// LogonRetHandler ?
type LogonRetHandler struct {
	playerManager    *player.Manager
	onUserLogOutFunc func(int32, string)
	kickUserFunc     func(connid string, msgstr string)
}

// NewLogonRetHandler 創建登錄結果處理器
func NewLogonRetHandler(playerManager *player.Manager, userLogoutFunc func(int32, string), kickUserFunc func(connid string, msgstr string)) *LogonRetHandler {
	reth := new(LogonRetHandler)
	reth.playerManager = playerManager
	reth.onUserLogOutFunc = userLogoutFunc
	reth.kickUserFunc = kickUserFunc

	return reth
}

// IsHook ?
func (h *LogonRetHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_LoginRetID)
}

// OnRPCMessage ?
func (h *LogonRetHandler) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserLoginRet{}
		return mm
	})

	if rmm == nil {
		return
	}

	mm := rmm.(*netproto.UserLoginRet)

	pr := h.playerManager.GetPlayerByConnID(rpcinfo.GetConnID())

	// 登錄成功，把AreaId存放在連線資訊中
	pr.Conn.AreaID = int64(mm.GetArea_ID())

	// 登錄成功，但客戶端連接已經斷開，此時要執行登出操作
	if pr == nil {
		if mm.GetCode() == 1 {
			uid := mm.GetUserID()
			cer := mm.GetCer()

			// 通知用戶登出大廳
			h.onUserLogOutFunc(uid, cer)
		}

		return
	}

	// 登錄失敗,斷開客戶端連接
	if mm.GetCode() == 0 {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("登錄失敗 尚未斷開客戶端連接")).Debug("HallServer壓測時")
		//ctx.Server.CloseConnByConnID(rpcinfo.GetConnID())
	} else {

		// 檢查用戶是否已經在線
		exur := h.playerManager.GetPlayerByUserID(mm.GetUserID())
		if exur != nil {
			exconnid := exur.Conn.ConnID()

			h.playerManager.ClearPlayerInfo(exconnid)

			// 如果不是同一連接,則踢出用戶
			if exconnid != rpcinfo.GetConnID() {
				h.kickUserFunc(exconnid, "您的账号已在别处登录.")
			}
		}

		// 檢查當前連接是否已經有登錄賬號
		if !pr.IsGuest() {
			if pr.UserID != mm.GetUserID() {
				h.onUserLogOutFunc(pr.UserID, pr.Cer)

				h.playerManager.ClearPlayerInfo(rpcinfo.GetConnID())
			}
		}

		h.playerManager.SetPlayerData(rpcinfo.GetConnID(), func(pr *player.Player) {
			// 設置用戶信息
			pr.Cer = mm.GetCer()
			pr.UserID = mm.GetUserID()
			pr.HDType = mm.GetHDType()
			pr.HDCode = mm.GetHDCode()
			pr.UserData = mm.GetUserData()
		})
	}

}
