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

// SuperToolAddAssignRewardHandler 寫遊戲記錄詳情
type SuperToolAddAssignRewardHandler struct{}

// IsHook ?
func (rh *SuperToolAddAssignRewardHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_SuperToolAssignRewardID)
}

// OnMessage ?
func (rh *SuperToolAddAssignRewardHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolAddAssignReward{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperToolAddAssignReward)
	gameid := mm.GetGameID()
	flagid := mm.GetFlagID()
	tableno := mm.GetTableNo()
	result := mm.GetResult()
	operator := mm.GetOperator()
	var logid int32
	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", gameid)
	ps.AddIntInput("intFlagID", flagid)
	ps.AddIntInput("intTableNo", tableno)
	ps.AddVarcharInput("chvResult", result, -1)
	ps.AddVarcharInput("chvOperator", operator, 50)
	ps.AddIntOutput("intOutLogID", logid)
	proc := db.NewProcedure("PrGs_SuperTool_AddAssignReward", ps)
	//mlog.System(proc.GetSQLString())
	ret, err := dbutil.DBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("執行存儲過程PrGs_SuperTool_AddAssignReward錯誤, %v, %s", err, proc.GetSQLString())
		retcode = 0
	} else {
		retcode = int32(ret.GetReturnValue())
		logid := int32(ret.GetOutputParamValue("intOutLogID").(int64))
		if logid > 0 {
			for idx, detail := range mm.GetDetailList() {
				ps := db.NewSQLParameters()
				ps.AddIntInput("intLogID", logid)
				ps.AddIntInput("intIdx", int32(idx))
				ps.AddVarcharInput("chvDescription", detail.GetDescription(), -1)
				ps.AddVarcharInput("chvReward", "", -1)
				ps.AddIntInput("intStatus", 0)
				proc := db.NewProcedure("PrGs_SuperTool_AddAssignRewardDetail", ps)
				mlog.Debug(proc.GetSQLString())
				_, err := dbutil.DBUtil.ExecProc(proc)
				if err != nil {
					retcode = 0
					break
				}
			}
		} else {
			retcode = 0
		}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}

}
