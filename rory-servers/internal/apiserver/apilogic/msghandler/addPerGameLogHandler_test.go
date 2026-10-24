// handler單元測試
package msghandler

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 寫GameLog

func TestAddPerGameLogHandler(t *testing.T) {
	InitTest()

	// 取得 token
	var accessToken = [REDACTED]

	// 驗證 handler
	serverid := int32(101000410)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid)
	rpcInfo.Cer = proto.String(util.ToString(accessToken))
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	//{"perGameLog":{"gameBeginTime":"2020-01-14 14:44:27","gameData":"\r\n2020-01-14 14:44:20 游戏开始.桌子号:1.userid:1000000,seat:0,score:6034760,\r\n14:44:20 玩家進入\r\n14:44:27 玩家 1000000; 底注: 5, 符號倍率: 8, 花費金額: 40\r\n14:44:27 N=200,AwardAdd=0,jpAmount=20000000|N=600,AwardAdd=0,jpAmount=7500000|N=1200,AwardAdd=0,jpAmount=375000|N=3000,AwardAdd=0,jpAmount=200000\r\n14:44:27 SymbolGroupID: 1, data: p3, p9, p5, p10, p3|p9, p9, p6, p4, p3|p2, p1, p6, p1, p4\r\n14:44:27 更新玩家金幣, 下注金幣: 40, 使用金幣: 40, 贏得獎金: 0\r\n2020-01-14 14:44:27 游戏结束.userid:1000000,winamount:-40,\r\n","gameEndTime":"2020-01-14 14:44:27","gameID":25,"serverID":365,"taxAmount":-0.40,"userGameDetail":"1000000:60347.60:-0.40:0,","winPrizeID":""}}
	// 驗證資料
	msgdata := new(netproto.DealLog)
	msgdata.ServerID = proto.Int32(serverid)
	msgdata.GameID = proto.Int32(1010)
	msgdata.TaxAmount = proto.Int64(0)
	msgdata.GameData = proto.String("\r\n2020-01-14 12:22:43 游戏开始.桌子号:0.userid:1000020,seat:0,score:1000000,\r\n12:22:43 玩家進入\r\n12:22:53 玩家 1000020; 底注: 5, 符號倍率: 88, 花費金額: 440\r\n12:22:53 N=2200,AwardAdd=0,jpAmount=20000000|N=6600,AwardAdd=0,jpAmount=7500000|N=13200,AwardAdd=0,jpAmount=375001|N=33000,AwardAdd=0,jpAmount=200003\r\n12:22:53 SymbolGroupID: 9, data: p7, p8, p1, p14, p12|p14, w1, p12, p14, p14|p12, p14, p13, p6, p13\r\n12:22:53 開獎符號: p12:left:3, 500 x [0 + 1(x1)] = 500\r\n12:22:53 即將進入JackpotGame\r\n12:22:53 JackpotGame獎金層級: Mini\r\n12:22:53 更新玩家金幣, 下注金幣: 440, 使用金幣: 440, 贏得獎金: 500\r\n2020-01-14 12:22:53 游戏结束.userid:1000020,winamount:60,\r\n")
	msgdata.UserGameRet = proto.String("1000174:1000000:60:0,")
	msgdata.UserGameDetails = []*netproto.UserGameDetail{
		{
			UserID:       proto.Int32(1000000),
			TradeID:      proto.Int64(0),
			CurrentScore: proto.Int64(1000000),
			UpdateScore:  proto.Int64(60),
			BetScore:     proto.Int64(40),
			WinScore:     proto.Int64(100),
			ValidScore:   proto.Int64(0),
		},
	}

	// msgdata.GameBeginTime = proto.String("2019-05-03T17:30:08+08:00Z")
	// msgdata.GameEndTime = proto.String("2019-05-03T18:30:08+08:00Z")
	msgdata.GameBeginTime = proto.String(time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339))
	msgdata.GameEndTime = proto.String(time.Now().UTC().Format(time.RFC3339))
	// msgdata.GameBeginTime = proto.String(util.FormatTimeStr(time.Now().UTC().Add(-1 * time.Hour)))
	// msgdata.GameEndTime = proto.String(util.FormatTimeStr(time.Now().UTC()))

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData

	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_DealLogID), rpcMessage)
	handler := new(AddPerGameLogHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID) {
		t.Error("OnMessage 錯誤")
	}

}
