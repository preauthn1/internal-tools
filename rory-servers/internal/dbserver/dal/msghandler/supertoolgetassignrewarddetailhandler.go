package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// SuperToolGetAssignRewardDetailHandler 寫遊戲記錄詳情
type SuperToolGetAssignRewardDetailHandler struct {
}

// IsHook ?
func (rh *SuperToolGetAssignRewardDetailHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailID)
}

// OnMessage ?
func (rh *SuperToolGetAssignRewardDetailHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolGetAssignRewardDetail{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperToolGetAssignRewardDetail)
	gameid := mm.GetGameID()
	flagid := mm.GetFlagID()
	tableno := mm.GetTableNo()
	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", gameid)
	ps.AddIntInput("intFlagID", flagid)
	ps.AddIntInput("intTableNo", tableno)
	proc := db.NewProcedure("PrGs_SuperTool_GetAssignRewardDetail", ps)
	//mlog.System(proc.GetSQLString())
	ret, err := dbutil.DBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("执行存储过程PrGs_SuperTool_GetAssignRewardDetail错误, %v, %s", err, proc.GetSQLString())
		retcode = 0
	}
	retcode = int32(ret.GetReturnValue())
	var retmsg *msg.Message
	if retcode == 1 {
		logid := ret[0].GetSingleValueInt32("LogID")
		idx := ret[0].GetSingleValueInt32("Idx")
		description := ret[0].GetSingleValueString("Description")
		reward := ret[0].GetSingleValueString("Reward")
		status := ret[0].GetSingleValueInt32("Status")
		retdata := new(netproto.SuperToolAssignRewardDetail)
		retdata.LogID = proto.Int32(logid)
		retdata.Idx = proto.Int32(idx)
		retdata.Description = proto.String(description)
		retdata.Reward = proto.String(reward)
		retdata.Status = proto.Int32(status)
		retmsg = msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailRetID), nil)
		retmsg.MsgData = retdata
		//mlog.System("===============return message")
	} else {
		retmsg = dbutil.GetDBRetMsg(retcode)
	}

	return []*msg.Message{retmsg}

}
