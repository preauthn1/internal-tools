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
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RequestUserInfoHandler ?
type RequestUserInfoHandler struct{}

// IsHook ?
func (h *RequestUserInfoHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_RequestUserHallInfoID)
}

// OnMessage ?
func (h *RequestUserInfoHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.RequestUserHallInfo{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.RequestUserHallInfo)

	userid := rpc.GetUserID()
	areaid := mm.GetSiteID()

	ps := db.NewSQLParameters()
	ps.AddIntOutput("intUserID", userid)
	ps.AddIntInput("intAreaID", areaid)
	ps.AddVarcharOutput("chvCertification", "")
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrPs_User_GetUserInfo", ps)
	mlog.System(proc.GetSQLString())
	//ret, err := DbUtil.ExecProc(proc)
	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserHallInfoID), nil)

	if err != nil || ret.GetRetTableCount() <= 0 {
		mlog.Error("执行存储过程PrPs_User_GetUserInfo出错%v", err)
		logret := &netproto.UserLoginRet{}
		logret.Code = proto.Int32(0)
		logret.Message = proto.String("服务器错误.")
		retmsg.MsgData = logret

	} else {
		logret := dbutil.ParseUserLogonRet(ret)
		logret.Area_ID = proto.Int32(areaid)
		retmsg.MsgData = logret
	}

	return []*msg.Message{retmsg}
}
