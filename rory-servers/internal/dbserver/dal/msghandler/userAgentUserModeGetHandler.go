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

// UserAgentUserModeGetHandler USER 取得已開的模式
type UserAgentUserModeGetHandler struct {
}

// IsHook ?
func (rh *UserAgentUserModeGetHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserAgentUserModeGetID)
}

// OnMessage ?
func (rh *UserAgentUserModeGetHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserAgentUserModeGetMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserAgentUserModeGetMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddVarcharOutput("chvMsg", "")

	//開始準備SP指令
	proc := db.NewProcedure("PrPs_UserAgent_UserMode_Get", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentUserModeGetRetID), nil)
	retdata := new(netproto.UserAgentUserModeGetRet)

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_UserAgent_UserMode_Get 错误 %v(MISSING)", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Userid = proto.Int32(int32(tb.GetValueByColName(0, "userid").(int64)))
				retdata.Mode1 = proto.Int32(int32(tb.GetValueByColName(0, "Mode1").(int64)))
				retdata.Mode2 = proto.Int32(int32(tb.GetValueByColName(0, "Mode2").(int64)))
				retdata.Mode3 = proto.Int32(int32(tb.GetValueByColName(0, "Mode3").(int64)))
				retdata.Mode4 = proto.Int32(int32(tb.GetValueByColName(0, "Mode4").(int64)))
				retdata.Mode5 = proto.Int32(int32(tb.GetValueByColName(0, "Mode5").(int64)))
				retdata.Mode6 = proto.Int32(int32(tb.GetValueByColName(0, "Mode6").(int64)))
			} else {
				retdata.Userid = proto.Int32(int32(0))
				retdata.Mode1 = proto.Int32(int32(0))
				retdata.Mode2 = proto.Int32(int32(0))
				retdata.Mode3 = proto.Int32(int32(0))
				retdata.Mode4 = proto.Int32(int32(0))
				retdata.Mode5 = proto.Int32(int32(0))
				retdata.Mode6 = proto.Int32(int32(0))
			}
			tb = dbret[1]
			if len(tb.Rows) > 0 {
				retdata.ChvMsg = proto.String(string(tb.GetValueByColName(0, "chvMsg").(string)))
			} else {
				retdata.ChvMsg = proto.String(string(""))
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
