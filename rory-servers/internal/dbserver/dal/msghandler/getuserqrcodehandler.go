package msghandler

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetUserQRCodeHandler 查詢我的推廣碼
type GetUserQRCodeHandler struct{}

// IsHook ?
func (rh *GetUserQRCodeHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetUserQACodeID)
}

// OnMessage ?
func (rh *GetUserQRCodeHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetUserQRCodeMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetUserQRCodeMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_UserAgent_GetUserQRCode", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserQACodeRetID), nil)
	retdata := new(netproto.GetUserQRCodeRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_UserAgent_GetUserQRCode错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.UserID = proto.Int32(int32(tb.GetValueByColName(0, "UserID").(int64)))
				retdata.QRCode = proto.String(tb.GetValueByColName(0, "QRCode").(string))
				retdata.InviteCode = proto.String(tb.GetValueByColName(0, "InviteCode").(string))
				retdata.UpperUserID = proto.Int32(int32(tb.GetValueByColName(0, "UpperUserID").(int64)))
			} else {
				retdata.UserID = proto.Int32(0)
				retdata.QRCode = proto.String("")
				retdata.InviteCode = proto.String("")
				retdata.UpperUserID = proto.Int32(0)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
