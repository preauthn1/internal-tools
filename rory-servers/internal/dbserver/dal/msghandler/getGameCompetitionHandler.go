package msghandler

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GameCompetitionHandler 查詢比賽活動公告
type GameCompetitionHandler struct{}

// IsHook 是否為這隻Handler該處理
func (rh *GameCompetitionHandler) IsHook(bclsid int32, sclsid int32) bool {
	return (bclsid == int32(netproto.MessageBClassID_Hall) || bclsid == int32(netproto.MessageBClassID_GameRoom)) &&
		(sclsid == int32(netproto.HallMsgClassID_GameCompetitionID))
}

// OnMessage 處理訊息
func (rh *GameCompetitionHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetGameCompetitionMsg{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetGameCompetitionMsg)

	ps1 := db.NewSQLParameters()
	ps1.AddIntInput("intGameID", mm.GetGameID())
	proc2 := db.NewProcedure("PrGs_Event_GetContestInfo", ps1)
	dbret2, err := dbutil.DBUtil.ExecProc(proc2)
	retcode2 := int32(dbret2.GetReturnValue())
	if err != nil {
		ctx.Error("执行存储过程PrGs_Event_GetContestInfo错误, %v, %s", err, proc2.GetSQLString())
		retcode2 = int32(0)
	}

	contestID := int32(0)
	maxorderreach := int32(0)
	maxfreeticket := int32(0)
	if retcode2 == 1 {
		if dbret2.GetRetTableCount() == 3 {
			//比賽資訊
			tbData1 := dbret2[0]
			if len(tbData1.Rows) > 0 {
				contestID = int32(tbData1.GetValueByColName(0, "ContestID").(int64)) //比賽活動ID
			}

			//最大充值金額和最大免費票
			tbData2 := dbret2[1]
			if len(tbData2.Rows) > 0 {
				maxorderreach = int32(tbData2.GetValueByColName(len(tbData2.Rows)-1, "OrderReach").(int64)) //最大充值金額
				maxfreeticket = int32(tbData2.GetValueByColName(len(tbData2.Rows)-1, "FreeTicket").(int64)) //最大免費票
			}
		}
	}

	retData := new(netproto.GameCompetitionMsg)

	retData.PiSaiTime = proto.Int32(getUserPiSaiTime(mm))

	if contestID != 0 {
		tbData0 := dbret2[0]
		if len(tbData0.Rows) > 0 {
			retData.Charge = proto.Int32(int32(tbData0.GetValueByColName(0, "Charge").(int64))) //入場費
		}

		ps := db.NewSQLParameters()
		ps.AddIntInput("intContestID", contestID)
		proc := db.NewProcedure("PrGs_Event_GetContestRank", ps)
		dbret, err := dbutil.DBUtil.ExecProc(proc)
		retcode := int32(dbret.GetReturnValue())
		if err != nil {
			ctx.Error("执行存储过程PrGs_Event_GetContestRank错误, %v, %s", err, proc.GetSQLString())
			retcode = int32(0)
		}

		if retcode == 1 {
			if dbret.GetRetTableCount() == 3 {
				//今日
				tbData1 := dbret[0]
				SetAllUserRank(tbData1, retData, "today")
				//yesterday
				tbData2 := dbret[1]
				SetAllUserRank(tbData2, retData, "yesterday")
				//yesterday
				tbData3 := dbret[2]
				SetAllUserRank(tbData3, retData, "week")

			}
		}

	} else {
		mlog.System("無比賽活動資訊 GameID[%v]", mm.GetGameID())
		retmsg := new(msg.Message)
		retmsg.BClassID = m.BClassID //int32(netproto.MessageBClassID_Hall)
		retmsg.SClassID = int32(netproto.HallMsgClassID_GameCompetitionRetID)
		retmsg.MsgData = retData

		return []*msg.Message{retmsg}
	}

	ps3 := db.NewSQLParameters()
	ps3.AddIntInput("intContestID", contestID)
	ps3.AddIntInput("intUserID", mm.GetUserID())

	proc3 := db.NewProcedure("PrGs_Event_GetContestUserInfo", ps3)
	dbret3, err := dbutil.DBUtil.ExecProc(proc3)
	retcode3 := int32(dbret2.GetReturnValue())
	if err != nil {
		ctx.Error("执行存储过程PrGs_Event_GetContestUserInfo错误, %v, %s", err, proc2.GetSQLString())
		retcode3 = int32(0)
	}

	if retcode3 == 1 {
		if dbret3.GetRetTableCount() == 1 {
			tbData1 := dbret3[0]
			if len(tbData1.Rows) == 1 {
				retData.UserID = proto.Int32(mm.GetUserID())
				retData.BestScore = proto.Int32(int32(tbData1.GetValueByColName(0, "HighScore").(int64)))
				retData.FreeTimes = proto.Int32(int32(tbData1.GetValueByColName(0, "FreeTicket").(int64)))
				retData.UserRanking = proto.Int32(int32(tbData1.GetValueByColName(0, "Rank").(int64)))
			}
		}

		s1 := fmt.Sprintf("储值%v可以获得%v次免费", maxorderreach/100, maxfreeticket)
		retData.MsgInfo = []string{"排行每20分钟更新一次", "每日24点统计当日排名并发放奖励", s1}
	}

	retmsg := new(msg.Message)
	retmsg.BClassID = m.BClassID
	retmsg.SClassID = int32(netproto.HallMsgClassID_GameCompetitionRetID)
	retmsg.MsgData = retData

	return []*msg.Message{retmsg}
}

func getUserPiSaiTime(mm *netproto.GetGameCompetitionMsg) int32 {
	ps3 := db.NewSQLParameters()
	ps3.AddIntInput("intGameID", mm.GetGameID())
	ps3.AddIntInput("intUserID", mm.GetUserID())

	proc3 := db.NewProcedure("PrGs_Event_GetContestUserData", ps3)
	dbret3, err := dbutil.DBUtil.ExecProc(proc3)
	retcode3 := int32(dbret3.GetReturnValue())
	if err != nil {
		mlog.Error("执行存储过程PrGs_Event_GetContestUserData错误, %v, %s", err, proc3.GetSQLString())
		retcode3 = int32(0)
	}

	if retcode3 == 1 {
		if dbret3.GetRetTableCount() == 1 {
			tbData1 := dbret3[0]
			if len(tbData1.Rows) == 1 {
				t := int32(tbData1.GetValueByColName(0, "Playtime").(int64))
				return t
			}
		}
	}

	return 0
}

// SetAllUserRank ?
func SetAllUserRank(tbData1 *db.RecordTable, retData *netproto.GameCompetitionMsg, titalName string) {
	leaderboardDetail1 := new(netproto.LeaderboardDetail)
	leaderboardDetail1.TitleName = proto.String(titalName)
	for i := 0; i < len(tbData1.Rows); i++ {
		today := new(netproto.ScoreboardList)
		today.SortID = proto.Int32(int32(tbData1.GetValueByColName(i, "Rank").(int64)))
		today.UserID = proto.String(tbData1.GetValueByColName(i, "UserID").(string))
		today.Score = proto.Int32(int32(tbData1.GetValueByColName(i, "Score").(int64)))
		leaderboardDetail1.LeaderboardList = append(leaderboardDetail1.LeaderboardList, today)
	}
	retData.Leaderboard = append(retData.Leaderboard, leaderboardDetail1)
}
