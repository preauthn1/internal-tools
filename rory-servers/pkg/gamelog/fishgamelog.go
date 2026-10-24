package gamelog

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// FishGameLog 捕魚遊戲記錄
type FishGameLog struct {
	*EmptyGameLog
	gameid      int32
	serverid    int32
	tableno     int32
	rpcMsgQueue *rpcmsg.Queue
}

// NewFishGameLog 捕魚遊戲記錄
func NewFishGameLog(rpcQueue *rpcmsg.Queue, gameid int32, serverid int32, tableno int32) *FishGameLog {
	flog := new(FishGameLog)
	flog.rpcMsgQueue = rpcQueue
	flog.gameid = gameid
	flog.serverid = serverid
	flog.tableno = tableno
	flog.EmptyGameLog = NewEmptyGameLog()

	return flog
}

// EndGameLog 遊戲步驟結束
func (f *FishGameLog) EndGameLog() {
	defer util.PanicErrStack()

	if f.ctx == nil {
		return
	}

	f.ctx.EndTime = time.Now()

	deallogmsg := new(msg.Message)
	deallogmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	deallogmsg.SClassID = int32(netproto.DBServerClassID_DealLogID)

	deallogdata := new(netproto.DealLog)
	deallogdata.GameBeginTime = proto.String(util.FormatTimeStr(f.ctx.BeginTime))
	deallogdata.GameEndTime = proto.String(util.FormatTimeStr(f.ctx.EndTime))
	deallogdata.GameID = proto.Int32(f.gameid)
	deallogdata.ServerID = proto.Int32(f.serverid)
	deallogdata.TableNo = proto.Int32(f.tableno)
	if f.nextJiangRecord != nil {
		deallogdata.NextJiangRecord = f.nextJiangRecord
	}

	var taxamount int64 = 0
	urets := ""
	for uid, ud := range f.ctx.Users {
		taxamount += ud.TaxAmount
		urets += fmt.Sprintf("%d:%d:%d:%d,", uid, ud.Score, ud.UpdateScore, ud.TradeID)
	}

	deallogdata.UserGameRet = proto.String(urets)
	deallogdata.GameData = proto.String(f.ctx.GetGameLogDesc())
	deallogdata.TaxAmount = proto.Int64(taxamount)

	deallogmsg.MsgData = deallogdata

	f.rpcMsgQueue.PushMessage(netproto.CreateEmptyRPCInfo(f.serverid), false, deallogmsg)

	f.ctx = nil
}
