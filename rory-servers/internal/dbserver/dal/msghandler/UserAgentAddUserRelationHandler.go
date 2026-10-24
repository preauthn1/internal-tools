package msghandler

import (
	"strings"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// UserAgentAddUserRelationHandler USER 加入下线(全民代理)
type UserAgentAddUserRelationHandler struct{}

// IsHook ?
func (rh *UserAgentAddUserRelationHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserAgentAddUserRelationID)
}

// OnMessage ?
func (rh *UserAgentAddUserRelationHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.AddUserRelationMsg{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.AddUserRelationMsg)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddVarcharInput("chvInviteCode", mm.GetInviteCode(), 20)
	ps.AddVarcharOutput("chvMsg", "")

	proc := db.NewProcedure("PrPs_UserAgent_AddUserRelation", ps)
	//mlog.System("[weihao] %v", proc.GetSQLString())
	ret, err := dbutil.DBUtil.ExecProc(proc)

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()

	var retcode int32 = 0
	var chvMsg string
	if err != nil {
		ctx.Error("执行存储过程PrPs_UserAgent_AddUserRelation错误, %v, %s", err, proc.GetSQLString())
		retcode = 0
	} else {
		retcode = int32(ret.GetReturnValue())
		chvMsg = ret.GetOutputParamValue("chvMsg").(string)
		iErr := strings.Index(chvMsg, "错误")
		if iErr > -1 {
			retcode = 0
		}
	}

	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retmsg.SClassID = int32(netproto.HallMsgClassID_UserAgentAddUserRelationRetID)
	data := new(netproto.AddUserRelationRet)
	data.RetCode = proto.Int32(retcode)
	data.ChvMsg = proto.String(chvMsg)
	retmsg.MsgData = data

	return []*msg.Message{retmsg}
}
