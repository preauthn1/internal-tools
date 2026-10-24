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

// GetUserAgentListHandler 查询推廣人數
type GetUserAgentListHandler struct{}

// IsHook ?
func (rh *GetUserAgentListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetUserAgentListID)
}

// OnMessage ?
func (rh *GetUserAgentListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetUserAgentListMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetUserAgentListMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_UserAgent_GetUserAgentList", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentListRetID), nil)
	retdata := new(netproto.GetUserAgentListRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_UserAgent_GetUserAgentList错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.TotalUser = proto.Int32(int32(tb.GetValueByColName(0, "totaluser").(int64)))
				retdata.TempBUser = proto.Int32(int32(tb.GetValueByColName(0, "temp_buser").(int64)))
				retdata.TodayNewUser = proto.Int32(int32(tb.GetValueByColName(0, "todaynewuser").(int64)))
				retdata.TodayLogin = proto.Int32(int32(tb.GetValueByColName(0, "todaylogin").(int64)))
			} else {
				retdata.TotalUser = proto.Int32(0)
				retdata.TempBUser = proto.Int32(0)
				retdata.TodayNewUser = proto.Int32(0)
				retdata.TodayLogin = proto.Int32(0)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
