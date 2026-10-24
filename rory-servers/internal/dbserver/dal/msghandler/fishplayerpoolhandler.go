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

// FishPlayerPoolHandler 新增/更新玩家個人水池輸贏記錄
type FishPlayerPoolHandler struct{}

// IsHook ?
func (fpp *FishPlayerPoolHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_FishPlayerPoolID)
}

// OnMessage ?
func (fpp *FishPlayerPoolHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.FishPlayerPoolDetail{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.FishPlayerPoolDetail)
	var retcode int32 = 1
	var ret db.RecordSet
	var err error
	for _, v := range mm.FishRewardList {
		ps := db.NewSQLParameters()
		ps.AddIntInput("intFlagID", mm.GetFlagID())
		ps.AddIntInput("intUserID", v.GetUserID())
		ps.AddBigIntInput("lngHitMoney", v.GetHitMoney())
		ps.AddBigIntInput("lngSpendMoney", v.GetSpendMoney())
		proc := db.NewProcedure("PrGs_SuperTool_Fish_UpdPlayerPoolWinLost", ps)
		ret, err = dbutil.DBUtil.ExecProc(proc)

		if err != nil {
			ctx.Error("执行存储过程PrGs_SuperTool_Fish_UpdSystemReward错误, %v, %s", err, proc.GetSQLString())
			retcode = 0
		}

		retcode = int32(ret.GetReturnValue())
	}

	var retmsg *msg.Message
	retmsg = msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_FishPlayerPoolRetID), nil)
	retdata := new(netproto.FishPlayerPoolDetail)

	if retcode == 1 {
		if ret.GetRetTableCount() > 0 {
			tbData := ret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				detail := new(netproto.FishPlayerPoolDetailList)
				detail.UserID = proto.Int32(int32(tbData.GetValueByColName(i, "UserID").(int64)))
				detail.HitMoney = proto.Int64(tbData.GetValueByColName(i, "HitMoney").(int64))
				detail.SpendMoney = proto.Int64(tbData.GetValueByColName(i, "SpendMoney").(int64))
				retdata.FishRewardList = append(retdata.FishRewardList, detail)
			}
			if len(tbData.Rows) == 0 {
				detail := new(netproto.FishPlayerPoolDetailList)
				detail.UserID = proto.Int32(0)
				detail.HitMoney = proto.Int64(0)
				detail.SpendMoney = proto.Int64(0)
				retdata.FishRewardList = append(retdata.FishRewardList, detail)
			}
		} else {
			detail := new(netproto.FishPlayerPoolDetailList)
			detail.UserID = proto.Int32(0)
			detail.HitMoney = proto.Int64(0)
			detail.SpendMoney = proto.Int64(0)
			retdata.FishRewardList = append(retdata.FishRewardList, detail)
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
