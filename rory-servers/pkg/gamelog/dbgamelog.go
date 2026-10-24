package gamelog

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// DBGameLog 基於數據庫的牌局記錄
type DBGameLog struct {
	*EmptyGameLog
	gameid      int32
	serverid    int32
	tableno     int32
	rpcMsgQueue *rpcmsg.Queue
}

// NewDBGameLog 創建基於數據庫的牌局記錄對象
func NewDBGameLog(rpcQueue *rpcmsg.Queue, gameid int32, serverid int32, tableno int32) *DBGameLog {
	dblog := new(DBGameLog)
	dblog.rpcMsgQueue = rpcQueue
	dblog.gameid = gameid
	dblog.serverid = serverid
	dblog.tableno = tableno
	dblog.EmptyGameLog = NewEmptyGameLog()

	return dblog
}

// EndGameLog 遊戲步驟結束
func (l *DBGameLog) EndGameLog() {
	defer util.PanicErrStack()

	if l.ctx == nil {
		return
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

	mlog.Debug("遊戲結束, 玩家人數: %d", len(l.ctx.Users))

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

		mlog.Debug("玩家[%d]: %v", uid, d)
		deallogdata.UserGameDetails = append(deallogdata.UserGameDetails, d)
	}

	deallogdata.UserGameRet = proto.String(urets)
	deallogdata.GameData = proto.String(l.ctx.GetGameLogDesc())

	deallogdata.TaxAmount = proto.Int64(taxamount)

	deallogmsg.MsgData = deallogdata

	l.rpcMsgQueue.PushMessage(netproto.CreateEmptyRPCInfo(l.serverid), false, deallogmsg)
}
