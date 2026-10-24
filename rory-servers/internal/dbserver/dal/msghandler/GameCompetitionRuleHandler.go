package msghandler

import (
	"time"

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

// GameCompetitionRuleHandler 活動比賽規則
type GameCompetitionRuleHandler struct{}

// IsHook ?
func (rh *GameCompetitionRuleHandler) IsHook(bclsid int32, sclsid int32) bool {
	return (bclsid == int32(netproto.MessageBClassID_Hall) && (sclsid == int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID))) ||
		(bclsid == int32(netproto.MessageBClassID_Game) && (sclsid == int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID)))
}

// OnMessage ?
func (rh *GameCompetitionRuleHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetGameCompetitionRuleMsg{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetGameCompetitionRuleMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", mm.GetGameID())
	proc := db.NewProcedure("PrGs_Event_GetContestInfo", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())
	if err != nil {
		ctx.Error("执行存储过程PrGs_Event_GetContestInfo错误, %v, %s", err, proc.GetSQLString())
		retcode = int32(0)
	}

	retdata := new(netproto.GameCompetitionRuleMsg)
	if retcode == 1 {
		if dbret.GetRetTableCount() == 3 {
			//比賽資訊
			tbData1 := dbret[0]
			if len(tbData1.Rows) > 0 {
				retdata.ContestID = proto.Int32(int32(tbData1.GetValueByColName(0, "ContestID").(int64))) //比賽活動ID
				retdata.Charge = proto.Int32(int32(tbData1.GetValueByColName(0, "Charge").(int64)))       //入場費用
				retdata.PlayTime = proto.Int32(int32(tbData1.GetValueByColName(0, "PlayTime").(int64)))   //秒
				retdata.Coin = proto.Int32(int32(tbData1.GetValueByColName(0, "Coin").(int64)))           //固定比賽幣
				retdata.StartTime = proto.String(tbData1.GetValueByColName(0, "StartTime").(time.Time).Format("2006-01-02 15:04:05"))
				retdata.EndTime = proto.String(tbData1.GetValueByColName(0, "EndTime").(time.Time).Format("2006-01-02 15:04:05"))
				retdata.Show = proto.Int32(int32(tbData1.GetValueByColName(0, "Show").(int64)))
				retdata.EnterTime = proto.String(tbData1.GetValueByColName(0, "EnterTime").(time.Time).Format("2006-01-02 15:04:05"))
			} else {
				mlog.System("Gameid[%v] 無比賽資訊", mm.GetGameID())

				retdata.Charge = proto.Int32(0)
				retdata.PlayTime = proto.Int32(0)
				retdata.Coin = proto.Int32(0)
				retdata.StartTime = proto.String("")
				retdata.EndTime = proto.String("")
				retdata.Show = proto.Int32(0)
				retdata.EnterTime = proto.String("")
				retdata.ContestID = proto.Int32(0)
				retdata.Area_ID = proto.Int32(0)
			}

			//充值免費票
			tbData2 := dbret[1]
			for i := 0; i < len(tbData2.Rows); i++ {
				orderTicket := new(netproto.OrderTicketMsg) //儲值遊玩次數
				orderTicket.ID = proto.Int32(int32(tbData2.GetValueByColName(i, "ID").(int64)))
				orderTicket.OrderReach = proto.Int32(int32(tbData2.GetValueByColName(i, "OrderReach").(int64)))
				orderTicket.FreeTicket = proto.Int32(int32(tbData2.GetValueByColName(i, "FreeTicket").(int64)))
				retdata.OrderTicket = append(retdata.OrderTicket, orderTicket)
			}

			//日周獎金
			tbData3 := dbret[2]
			for i := 0; i < len(tbData3.Rows); i++ {
				reward := new(netproto.RewardDetail) //每日獎勵
				hasType := false
				for _, check := range retdata.RewardInfo {
					t := int32(tbData3.GetValueByColName(i, "Type").(int64))
					if check.GetType() == t {
						hasType = true
						break
					}
				}
				if hasType {
					continue
				}

				reward.Type = proto.Int32(int32(tbData3.GetValueByColName(i, "Type").(int64))) //1日 2周
				//if day or week
				if reward.GetType() == 1 {
					reward.DeadLine = proto.String(tbData1.GetValueByColName(0, "Deadline").(time.Time).Format("15:04:05"))
				} else if reward.GetType() == 2 {
					reward.DeadLine = proto.String(tbData1.GetValueByColName(0, "PeriodDeadline").(time.Time).Format("15:04:05"))
				}

				for j := 0; j < len(tbData3.Rows); j++ {
					if reward.GetType() == int32(tbData3.GetValueByColName(j, "Type").(int64)) {
						reward.Bonus = append(reward.Bonus, int32(tbData3.GetValueByColName(j, "Money").(int64)))
					}
				}
				retdata.RewardInfo = append(retdata.RewardInfo, reward)
			}
		} else {
			mlog.System("Gameid[%v] 無比賽資訊", mm.GetGameID())
		}
	}

	bclassid := int32(0)
	if m.BClassID == int32(netproto.MessageBClassID_Hall) {
		bclassid = int32(netproto.MessageBClassID_Hall)
	} else {
		bclassid = int32(netproto.MessageBClassID_Game)
	}

	retmsg := msg.NewMessage(bclassid, int32(netproto.HallMsgClassID_GameCompetioniInfoRuleRetID), nil)

	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
