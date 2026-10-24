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

// FishSysPoolHandler Fish 新增/更新系統水池輸贏記錄
type FishSysPoolHandler struct{}

// IsHook ?
func (fsp *FishSysPoolHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_FishSysPoolID)
}

// OnMessage ?
func (fsp *FishSysPoolHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.FishSysPoolData{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.FishSysPoolData)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intFlagID", mm.GetFlagID())
	ps.AddBigIntInput("lngBetMoney", mm.GetSumBetMoney())
	ps.AddBigIntInput("lngBullet", mm.GetSumBullet())
	ps.AddBigIntInput("lngPoolAmount", mm.GetPoolAmount())
	proc := db.NewProcedure("PrGs_SuperTool_Fish_UpdSystemPoolWinLost", ps)
	ret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(ret.GetReturnValue())

	if err != nil {
		ctx.Error("执行存储过程PrGs_SuperTool_Fish_UpdSystemPoolWinLost错误, %v, %s", err, proc.GetSQLString())
		retcode = int32(0)
	}

	var retmsg *msg.Message
	retmsg = msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_FishSysPoolRetID), nil)
	detail := new(netproto.FishSysPoolData)

	if retcode == 1 {
		if ret.GetRetTableCount() > 0 {
			tbData := ret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				detail.FlagID = proto.Int32(int32(tbData.GetValueByColName(i, "FlagID").(int64)))
				detail.SumBetMoney = proto.Int64(tbData.GetValueByColName(i, "BetMoney").(int64))
				detail.SumBullet = proto.Int64(tbData.GetValueByColName(i, "Bullet").(int64))
				detail.SysSumWinLoss = proto.Int64(tbData.GetValueByColName(i, "PlayerTotalWinLost").(int64))
				detail.AfterAdjustSyspoolWinLoss = proto.Int64(tbData.GetValueByColName(i, "PlayerTotalWinLost").(int64))
				break
			}
		}
	}

	retmsg.MsgData = detail

	return []*msg.Message{retmsg}
}
