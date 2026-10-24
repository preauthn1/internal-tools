package gamelog

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ConsoleGameLog Console 顯示牌局紀錄
type ConsoleGameLog struct {
	*EmptyGameLog
	gameid   int32
	serverid int32
	tableno  int32
}

// NewConsoleGameLog Create ConsoleGameLog
func NewConsoleGameLog(gameid int32, serverid int32, tableno int32) *ConsoleGameLog {
	logger := new(ConsoleGameLog)
	logger.EmptyGameLog = NewEmptyGameLog()
	logger.gameid = gameid
	logger.serverid = serverid
	logger.tableno = tableno

	return logger
}

// EndGameLog 結束
func (l *ConsoleGameLog) EndGameLog() {

	if l.ctx == nil {
		return
	}

	if !conf.AnalysisMode {
		fmt.Println(l.ctx.GetGameLogDesc())
	}

	l.ctx.EndTime = time.Now()

	deallogmsg := new(msg.Message)
	deallogmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	deallogmsg.SClassID = int32(netproto.DBServerClassID_DealLogID)

	deallogdata := new(netproto.DealLog)
	deallogdata.GameBeginTime = proto.String(util.FormatTimeStr(l.ctx.BeginTime))
	deallogdata.GameEndTime = proto.String(util.FormatTimeStr(l.ctx.EndTime))
	deallogdata.GameID = proto.Int32(l.gameid)
	deallogdata.ServerID = proto.Int32(l.serverid)
	deallogdata.TableNo = proto.Int32(l.tableno)
	if l.nextJiangRecord != nil {
		deallogdata.NextJiangRecord = l.nextJiangRecord
	}

	var taxamount int64 = 0
	urets := ""
	for uid, ud := range l.ctx.Users {

		urets += fmt.Sprintf("%d:%d:%d:%d,", uid, ud.Score, ud.UpdateScore, ud.TradeID)

		d := new(netproto.UserGameDetail)
		d.UserID = proto.Int32(uid)
		d.TradeID = proto.Int64(ud.TradeID)
		d.CurrentScore = proto.Int64(ud.Score)
		d.UpdateScore = proto.Int64(ud.UpdateScore)
		d.BetScore = proto.Int64(ud.BetScore)
		d.WinScore = proto.Int64(ud.PayoutScore)
		d.ValidScore = proto.Int64(ud.ValidScore)
		d.TaxAmount = proto.Int64(ud.TaxAmount)
		d.TransactionNo = proto.String(ud.TransactionNo)
		taxamount += ud.TaxAmount

		deallogdata.UserGameDetails = append(deallogdata.UserGameDetails, d)
	}

	deallogdata.UserGameRet = proto.String(urets)
	deallogdata.GameData = proto.String(l.ctx.GetGameLogDesc())

	deallogdata.TaxAmount = proto.Int64(taxamount)

	if !conf.AnalysisMode {
		mlog.Debug("%v", deallogdata)
	}

}
