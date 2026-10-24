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

// Cd21SysPoolHandler CD21 新增/更新系統水池輸贏記錄
type Cd21SysPoolHandler struct{}

// IsHook ?
func (fsp *Cd21SysPoolHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.CD21_GameMessageClassID_Cd21SysPoolID)
}

// OnMessage ?
func (fsp *Cd21SysPoolHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	//傳送syspool data
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Cd21SysPoolData{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Cd21SysPoolData)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intFlagID", mm.GetFlagID())
	ps.AddBigIntInput("lngBetMoney", mm.GetSumBetMoney())
	ps.AddBigIntInput("lngBetCounts", mm.GetSumBetCount())
	ps.AddBigIntInput("lngPoolAmount", mm.GetPoolAmount())
	proc := db.NewProcedure("PrGs_SuperTool_TwentyOne_UpdSystemPoolWinLost", ps)
	//mlog.System("-------:",proc.GetSQLString())
	ret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(ret.GetReturnValue())

	if err != nil {
		ctx.Error("執行存儲過程[PrGs_SuperTool_TwentyOne_UpdSystemPoolWinLost]錯誤, %v, %s", err, proc.GetSQLString())
		retcode = int32(0)
	}

	var retmsg *msg.Message
	retmsg = msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.CD21_GameMessageClassID_Cd21SysPoolRetID), nil)

	//get syspooldata
	detail := new(netproto.Cd21SysPoolData)
	if retcode == 1 {
		if ret.GetRetTableCount() > 0 {
			tbData := ret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				detail.FlagID = proto.Int32(int32(tbData.GetValueByColName(i, "FlagID").(int64)))
				detail.SumBetMoney = proto.Int64(tbData.GetValueByColName(i, "BetMoney").(int64))
				detail.SumBetCount = proto.Int64(tbData.GetValueByColName(i, "BetCounts").(int64))
				detail.PoolAmount = proto.Int64(tbData.GetValueByColName(i, "PoolAmount").(int64))
				mlog.Debug("[linda] flag[%v] sum[%v]", detail.GetFlagID(), detail.GetPoolAmount())
				break
			}
		}
	}

	retmsg.MsgData = detail

	return []*msg.Message{retmsg}
}
