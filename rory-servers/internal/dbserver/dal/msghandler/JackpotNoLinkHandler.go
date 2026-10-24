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

// JackpotNoLinkHandler ?
type JackpotNoLinkHandler struct{}

// IsHook ?
func (fsp *JackpotNoLinkHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && (sclsid == int32(netproto.GameMessageClassID_JackpotNoLinkID))
}

// OnMessage ?
func (fsp *JackpotNoLinkHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	//傳送syspool data
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.JackpotData{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.JackpotData)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", mm.GetGameID())
	ps.AddIntInput("intFlagID", mm.GetFlagID())
	ps.AddIntInput("intTableNo", mm.GetTableNo())
	ps.AddBigIntInput("lngContribution", mm.GetContribution())
	ps.AddBigIntInput("lngTotalContribution", mm.GetTotalContribution())
	ps.AddBigIntInput("lngAwardTotal", mm.GetAwardTotal())
	ps.AddIntInput("intLevelID", mm.GetLevel())
	ps.AddBigIntInput("lngJackpotPool", mm.GetJackpotPool())

	ps.AddVarcharOutput("chvMsg", "")
	proc := db.NewProcedure("PrGs_Poker_JackpotPool_NoLink", ps)
	//mlog.System("-------:",proc.GetSQLString())
	ret, err := dbutil.DBUtil.ExecProc(proc)
	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()
	retcode := int32(ret.GetReturnValue())

	if err != nil {
		ctx.Error("执行存储过程[PrGs_Poker_JackpotPool_NoLink]错误, %v, %s", err, proc.GetSQLString())
		retcode = int32(0)
	}

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_JackpotNoLinkRetID), nil)
	bmsg := new(netproto.JackpotAllData)
	pdata := make([]*netproto.JackpotData, 0)
	//get syspooldata
	if retcode == 1 {
		if ret.GetRetTableCount() > 0 {
			tbData := ret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				detail := new(netproto.JackpotData)
				detail.GameID = proto.Int32(tbData.GetValueByColNameInt32(i, "GameID"))
				detail.FlagID = proto.Int32(tbData.GetValueByColNameInt32(i, "FlagID"))
				detail.TableNo = proto.Int32(tbData.GetValueByColNameInt32(i, "TableNo"))
				detail.Contribution = proto.Int64(tbData.GetValueByColNameInt64(i, "Contribution"))
				detail.TotalContribution = proto.Int64(tbData.GetValueByColNameInt64(i, "TotalContribution"))
				detail.AwardTotal = proto.Int64(tbData.GetValueByColNameInt64(i, "AwardTotal"))
				detail.Level = proto.Int32(tbData.GetValueByColNameInt32(i, "LevelID"))
				detail.JackpotPool = proto.Int64(tbData.GetValueByColNameInt64(i, "JackpotPool"))
				pdata = append(pdata, detail)
			}
		}
	}
	bmsg.AllData = pdata
	retmsg.MsgData = bmsg

	return []*msg.Message{retmsg}
}
