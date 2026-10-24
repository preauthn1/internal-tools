package task

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// UpdateUserNotifyIdxHandler 服務器註冊處理器
type UpdateUserNotifyIdxHandler struct{}

// NewUpdateUserNotifyHandler 創建服務器註冊處理器
func NewUpdateUserNotifyHandler() *UpdateUserNotifyIdxHandler {
	handler := new(UpdateUserNotifyIdxHandler)
	return handler
}

// IsHook ?
func (h *UpdateUserNotifyIdxHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_UpdateUserNotifyIdxH2NID)
}

// OnMessage ?
func (h *UpdateUserNotifyIdxHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	defer util.PanicErrStack()
	mmbyte := m.MsgData.([]byte)

	mm := &netproto.UpdateUserNotifyIdx{}
	err := proto.Unmarshal(mmbyte, mm)

	if err != nil {
		ctx.Error("解析notifyserver收到大廳更新線上用戶錯誤,%v", err)
		ctx.Server.CloseConn(clt)
		return nil
	}

	d := GetDatabase()
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddBigIntInput("intIdx", mm.GetIdx())
	proc := db.NewProcedure("PrNs_UpdateNotifyIdx", ps)
	dbret, err := d.ExecProc(proc)
	if err != nil {
		mlog.Error("執行存儲過程PrNs_UpdateNotifyIdx錯誤, %v", err)
	}
	if dbret.GetReturnValue() != 1 {
		mlog.Error("執行存儲過程PrNs_UpdateNotifyIdx錯誤, %d", dbret.GetReturnValue())
	}
	return nil
}
