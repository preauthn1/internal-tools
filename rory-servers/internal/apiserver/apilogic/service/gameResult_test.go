package service

import (
	"fmt"
	"testing"
	"time"

	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// func TestGameResult_GameResultPerGameLog(t *testing.T) {
// 	httputil.DebugMode(true)
// 	service := &GameResult{}
// 	service.Base = newBase()
// 	userid, cer := GetTestAccessToken()
// 	rpcInfo := &netproto.RPCInfo{}
// 	rpcInfo.UserID = proto.Int(int(userid))
// 	rpcInfo.Cer = proto.String(cer)
// 	rpcInfo.IPAddress = proto.String("0.0.0.0")
// 	params := map[string]interface{}{
// 		"perGameLog": map[string]interface{}{
// 			"serverID":       101000410,
// 			"gameData":       "\r\n2020-01-14 12:22:43 游戏开始.桌子号:0.userid:1000020,seat:0,score:1000000,\r\n12:22:43 玩家進入\r\n12:22:53 玩家 1000020; 底注: 5, 符號倍率: 88, 花費金額: 440\r\n12:22:53 N=2200,AwardAdd=0,jpAmount=20000000|N=6600,AwardAdd=0,jpAmount=7500000|N=13200,AwardAdd=0,jpAmount=375001|N=33000,AwardAdd=0,jpAmount=200003\r\n12:22:53 SymbolGroupID: 9, data: p7, p8, p1, p14, p12|p14, w1, p12, p14, p14|p12, p14, p13, p6, p13\r\n12:22:53 開獎符號: p12:left:3, 500 x [0 + 1(x1)] = 500\r\n12:22:53 即將進入JackpotGame\r\n12:22:53 JackpotGame獎金層級: Mini\r\n12:22:53 更新玩家金幣, 下注金幣: 440, 使用金幣: 440, 贏得獎金: 500\r\n2020-01-14 12:22:53 游戏结束.userid:1000020,winamount:60,\r\n",
// 			"userGameDetail": "15267354:1000000:6000:0,",
// 			"gameBeginTime":  time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339),
// 			"gameEndTime":    time.Now().UTC().Format(time.RFC3339),
// 			"taxAmount":      0,
// 			"gameID":         1010,
// 			"winPrizeID":     "", //超端用
// 		},
// 	}
// 	if ok := service.PostPerGameLog(rpcInfo, params); ok {
// 		t.Log("TestGameResult_GameResultPerGameLog success")
// 	} else {
// 		t.Error("TestGameResult_GameResultPerGameLog fail")
// 	}
// }

func TestGameResult_SettlementV2(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	var uid int = 1015909
	var cer string = "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTkwOSwiTG9naW5OYW1lIjoiZGVtby0xNjEyMDk4Njg0NjEwIiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjEyOTYyNjg5LCJSbmRJRCI6IjljMGUxZWUzLWRhZDktNDY3Yy04OGJlLWU5ZjY3NWExNmY4ZiIsIlNlc3Npb25Ub2tlbiI6ImVhMzUxZWJhLTIyMmMtNGI5Zi04MzU2LTIyMzMwOGFiYWI3MyIsIlNpZ25LZXkiOiIyNmQyNDNhNWQzODE0OTkyNTZiMzYzOGM5OTYyOTAyMyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj+y/j+cZglfyueeAlO6OxOVhT5mlcBNYUjPJHeg60m/udJBYAz9rUFGBfa+t5rSIS5vjhjHWk5sXzOq6rv9cjzo4ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+xwkImZSqPczNB3l7uxmVIS3JDJ/47U7Xl01XIll17nXxiCtuHw4ruATjCazNaKSf1hho6vFbQckz/HoH0l+HgD8TWQRDz9s30tn59nuG2wsyxI4YU9tq+oVz0O3hicsWvSuyoMkAM0l5+7pNqyjuqgkh4ceHqHLvIiRB9yr/eL4wr7M/poYPK+XnTKCqHddlVbU0F/V5XFaVN9FcLc0Em8YV9Eyh6Kaq3+JbASZcdSSnoFYBs89G74KaqjlA4HDb6KvdszgaotF7hEBi8F53nyiAT1DOpVEk+SECSfzEHI794MaLtMAeE+rtw72hrpMsEE+/JGqq2tacREC1f3g5Di"

	var servID int32 = 102100410
	var tradeID int64 = 20191200000095331
	var updateAmt int64 = 1000
	var betAmt int64 = 100
	var ipAddr string = "0.0.0.0"
	var transactionNo string = "9708f8507fc343608100b07195bbdcd9"

	ok, afterAmt := GameResultService.PostSettlementV2(cer, int32(uid), servID, tradeID, updateAmt, ipAddr, transactionNo, betAmt, 0, updateAmt, 0)

	fmt.Println(ok, afterAmt, transactionNo)
}

func TestGameResult_PerGameLogV2(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	var uid int = 1015909
	var cer string = "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTkwOSwiTG9naW5OYW1lIjoiZGVtby0xNjEyMDk4Njg0NjEwIiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjEyOTYyNjg5LCJSbmRJRCI6IjljMGUxZWUzLWRhZDktNDY3Yy04OGJlLWU5ZjY3NWExNmY4ZiIsIlNlc3Npb25Ub2tlbiI6ImVhMzUxZWJhLTIyMmMtNGI5Zi04MzU2LTIyMzMwOGFiYWI3MyIsIlNpZ25LZXkiOiIyNmQyNDNhNWQzODE0OTkyNTZiMzYzOGM5OTYyOTAyMyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj+y/j+cZglfyueeAlO6OxOVhT5mlcBNYUjPJHeg60m/udJBYAz9rUFGBfa+t5rSIS5vjhjHWk5sXzOq6rv9cjzo4ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+xwkImZSqPczNB3l7uxmVIS3JDJ/47U7Xl01XIll17nXxiCtuHw4ruATjCazNaKSf1hho6vFbQckz/HoH0l+HgD8TWQRDz9s30tn59nuG2wsyxI4YU9tq+oVz0O3hicsWvSuyoMkAM0l5+7pNqyjuqgkh4ceHqHLvIiRB9yr/eL4wr7M/poYPK+XnTKCqHddlVbU0F/V5XFaVN9FcLc0Em8YV9Eyh6Kaq3+JbASZcdSSnoFYBs89G74KaqjlA4HDb6KvdszgaotF7hEBi8F53nyiAT1DOpVEk+SECSfzEHI794MaLtMAeE+rtw72hrpMsEE+/JGqq2tacREC1f3g5Di"

	rpcInfo := &netproto.RPCInfo{}
	rpcInfo.UserID = proto.Int(uid)
	rpcInfo.Cer = proto.String(cer)
	rpcInfo.IPAddress = proto.String("0.0.0.0")

	var gameID int32 = 1021
	var servID int32 = 102100410
	var tradeID int64 = 20191200000095331
	var currentAmt int64 = 100081600
	var updateAmt int64 = 1000
	var betAmt int64 = 100
	var tranNo string = "9708f8507fc343608100b07195bbdcd9"

	now := time.Now().UTC()

	dealLog := new(netproto.DealLog)
	dealLog.GameID = proto.Int32(gameID)
	dealLog.ServerID = proto.Int32(servID)
	dealLog.GameData = proto.String("unit test")
	dealLog.TaxAmount = proto.Int64(0)
	dealLog.GameBeginTime = proto.String(now.Format(time.RFC3339))
	dealLog.GameEndTime = proto.String(now.Format(time.RFC3339))

	d1 := new(netproto.UserGameDetail)
	d1.UserID = proto.Int32(int32(uid))
	d1.TradeID = proto.Int64(tradeID)
	d1.CurrentScore = proto.Int64(currentAmt)
	d1.UpdateScore = proto.Int64(updateAmt)
	d1.BetScore = proto.Int64(betAmt)
	d1.ValidScore = proto.Int64(0)
	d1.WinScore = proto.Int64(updateAmt)
	d1.TaxAmount = proto.Int64(0)
	d1.TransactionNo = proto.String(tranNo)

	dealLog.UserGameDetails = append(dealLog.UserGameDetails, d1)

	ok := GameResultService.PostPerGameLogV2(rpcInfo, dealLog)

	fmt.Println(ok)
}
