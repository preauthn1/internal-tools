package msghandler

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GameCompetitionStartHandler 報名參加活動
type GameCompetitionStartHandler struct{}

// IsHook ?
func (fsu *GameCompetitionStartHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_BuyTicketID)
}

// OnMessage ?
func (fsu *GameCompetitionStartHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetBuyTicketMsg{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetBuyTicketMsg)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BuyTicketRetID), nil)
	retdata := new(netproto.BuyTicketMsg)
	retdata.Area_ID = proto.Int32(mm.GetArea_ID())

	contestInfo, err := GetContestInfo(ctx, mm)
	if err != nil {
		retdata.RetCode = proto.Int32(-1)
		retdata.RetMsg = proto.String(err.Error())

		retmsg.MsgData = retdata
		return []*msg.Message{retmsg}
	}

	contestNo := GetContestUserData(ctx, mm)

	contestID := contestInfo.GetValueByColName(0, "ContestID").(int64)

	ps := db.NewSQLParameters()

	ps.AddIntInput("intContestID", int32(contestID))
	ps.AddIntInput("intContestNo", contestNo) //該場次的第幾局(用來判斷要不要扣錢/免費次數)
	ps.AddIntInput("intGameID", mm.GetGameID())
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddVarcharOutput("chvErrMsg", "")
	ps.AddIntOutput("intRtnCode", 0)

	proc := db.NewProcedure("PrGs_Event_StartContest", ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		ctx.Error("執行存儲過程PrGs_Event_StartContest錯誤, %v, %s", err, proc.GetSQLString())
		retdata.RetCode = proto.Int32(-1)
		retdata.RetMsg = proto.String("執行存儲過程錯誤")
		retmsg.MsgData = retdata
		return []*msg.Message{retmsg}
	}

	if dbret.GetRetTableCount() == 0 {
		retdata.RetCode = proto.Int32(-1)
		retdata.RetMsg = proto.String("執行存儲過程錯誤")
		retmsg.MsgData = retdata
		return []*msg.Message{retmsg}
	}

	retdata.RetMsg = proto.String(dbret[0].GetValueByColName(0, "chvErrMsg").(string))
	retdata.RetCode = proto.Int32(int32(dbret[0].GetValueByColName(0, "intRtnCode").(int64)))
	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}

// GetContestUserData 取得玩家比賽場次
func GetContestUserData(ctx *module.MouduleCtx, mm *netproto.GetBuyTicketMsg) int32 {

	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", mm.GetGameID())
	ps.AddIntInput("intUserID", mm.GetUserID())

	proc := db.NewProcedure("PrGs_Event_GetContestUserData", ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		ctx.Error("執行存儲過程PrGs_Event_GetContestUserData錯誤, %v, %s", err, proc.GetSQLString())
		return 0
	}

	var contestNo int32 = 1

	if dbret.GetReturnValue() == 1 && dbret.GetRetTableCount() == 1 && len(dbret[0].Rows) > 0 {
		contestNo = int32(dbret[0].GetValueByColName(0, "ContestNo").(int64)) //比賽場次
		playtime := int32(dbret[0].GetValueByColName(0, "Playtime").(int64))
		if playtime <= 0 {
			contestNo++
		}
	}

	return contestNo
}

// GetContestInfo 取得比賽資訊
func GetContestInfo(ctx *module.MouduleCtx, mm *netproto.GetBuyTicketMsg) (*db.RecordTable, error) {

	sp := db.NewSQLParameters()
	sp.AddIntInput("intGameID", mm.GetGameID())

	proc := db.NewProcedure("PrGs_Event_GetContestInfo", sp)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		ctx.Error("執行存儲過程PrGs_Event_GetContestInfo錯誤, %v, %s", err, proc.GetSQLString())
		return nil, fmt.Errorf("查無比賽資訊, GameID: %d", mm.GetGameID())
	}

	if dbret.GetReturnValue() == 1 && dbret.GetRetTableCount() == 3 && len(dbret[0].Rows) > 0 {
		return dbret[0], nil
	}

	return nil, fmt.Errorf("查無比賽資訊, GameID: %d", mm.GetGameID())
}
