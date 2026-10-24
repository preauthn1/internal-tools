package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// ContactServiceHandler 聯繫客服
type ContactServiceHandler struct{}

// IsHook ?
func (rh *ContactServiceHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_ContactServiceID)
}

// OnMessage ?
func (rh *ContactServiceHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.ContactService{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.ContactService)
	msgContent := mm.GetMsg()
	questType := mm.GetQuestionType()
	platformID := mm.GetPlatformID()
	areaID := mm.GetAreaID()

	userid := rpcinfo.GetUserID()
	ip := rpcinfo.GetIPAddress()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intPlatformID", platformID)
	ps.AddIntInput("intAreaID", areaID)
	ps.AddIntInput("intUserID", userid)
	ps.AddVarcharInput("chvMsg", msgContent, 512)
	ps.AddIntInput("intQuestionType", questType)

	ps.AddVarcharInput("chvIPAddress", ip, 15)
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrPs_User_ContactService", ps)

	ret, err := dbutil.DBUtil.ExecProc(proc)

	var code int32 = 0
	retmessage := ""
	if err != nil {
		mlog.Error("执行存储过程PrPs_User_ContactService错误, %v", err)
		code = 0
		retmessage = "服务器错误"
	} else {
		if ret.GetReturnValue() == 1 {
			code = 1
		}
		retmessage = ret.GetOutputParamValue("chvErrMsg").(string)
	}

	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retmsg.SClassID = int32(netproto.HallMsgClassID_ContactServiceRetID)
	retmsg.MsgData = dbutil.GetRetMsg(code, retmessage)

	return []*msg.Message{retmsg}
}
