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

// SlotJackpotLinkHandler Slot 連接型Jackpot寫入及讀取
type SlotJackpotLinkHandler struct{}

// IsHook ?
func (fsp *SlotJackpotLinkHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && (sclsid == int32(netproto.GameMessageClassID_SlotJackpotLinkID))
}

// OnMessage ?
func (fsp *SlotJackpotLinkHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	//傳送syspool data
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_JackpotData{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_JackpotData)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", mm.GetGameID())
	ps.AddIntInput("intFlagID", mm.GetFlagID())
	ps.AddBigIntInput("lngContribution", mm.GetContribution())
	ps.AddBigIntInput("lngTotalContribution", mm.GetTotalContribution())
	ps.AddBigIntInput("lngAwardTotal", mm.GetAwardTotal())
	ps.AddIntInput("intLevelID", mm.GetLevel())
	ps.AddVarcharOutput("chvMsg", "")
	proc := db.NewProcedure("PrGs_Slot_JackpotPool_Link", ps)
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
		ctx.Error("执行存储过程[PrGs_Slot_JackpotPool_Link]错误, %v, %s", err, proc.GetSQLString())
		retcode = int32(0)
	}

	var retmsg *msg.Message
	retmsg = msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_SlotJackpotLinkRetID), nil)

	//get syspooldata
	detail := new(netproto.Slot_JackpotData)
	if retcode == 1 {
		if ret.GetRetTableCount() > 0 {
			tbData := ret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				detail.GameID = proto.Int32(tbData.GetValueByColNameInt32(i, "GameID"))
				detail.FlagID = proto.Int32(tbData.GetValueByColNameInt32(i, "FlagID"))
				detail.TableNo = proto.Int32(int32(0))
				detail.Contribution = proto.Int64(tbData.GetValueByColNameInt64(i, "Contribution"))
				detail.TotalContribution = proto.Int64(tbData.GetValueByColNameInt64(i, "TotalContribution"))
				detail.AwardTotal = proto.Int64(tbData.GetValueByColNameInt64(i, "AwardTotal"))
				break
			}
		}
	}

	retmsg.MsgData = detail

	return []*msg.Message{retmsg}
}
