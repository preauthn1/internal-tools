package sgj

import (
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func NewShuiGuoJiLogicByTest() *ShuiGuoJiLogic {
	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)

	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
	rumanager := gamelogic.NewRoomUserManager(new(player.Manager))
	room := gamelogic.NewRoom(servercfg, rumanager)
	table := room.GetTableByNo(0)

	evtdispatcher := event.NewAsyncEventDispatcher(1)
	tfManager := event.NewTimerFuncManager(evtdispatcher)

	glc := new(ShuiGuoJiLogic)

	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(table, tfManager, cfgs)
	glc.getsgjConfig(cfgs)
	glc.gameState = waitState
	glc.chiShui = 0
	glc.zhuangUserID = 0
	glc.syslastwin = 0
	glc.systotalwin = 0
	glc.nextjiangtype = -1
	glc.zhuangList = make([]int32, 0)
	glc.startsyspool = glc.sysPool
	glc.lastResultTypeList = make([]*netproto.SGJ_JiangTypeOdd, 0)
	glc.addCmd()
	glc.betTotal = 0
	glc.RegisterAfterGameEndFunc(func() {
		glc.clearGameData() //清除數據
		glc.GameLog.AddGameStep("庄家%d", glc.zhuangUserID)
		glc.zhuShui()
		glc.broadCastGameState()
	})

	return glc
}

func TestGameRound(t *testing.T) {
	var userData []*sgjUserData
	glc := NewShuiGuoJiLogicByTest()
	glc.robotzhuang = 0
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	//新增玩家資料
	glc.GameStartFunc()
	users := fixtures.GetGameUser()
	for i := 0; i < len(users); i++ {
		userData = append(userData, newUserData(users[i]))
	}

	//安排用戶入座及下注
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//机器人上庄
	//glc.robotzhuang()
	//n := rand.Int31n(glc.maxrobotzhuang)
	//t.Log("機器人隨機人數",n)

	//t.Log("")
	//下注資訊: bettype, count
	glc.userBet(userData[0], fixtures.GetUserBet(1, 100))
	glc.userBet(userData[0], fixtures.GetUserBet(6, 5))
	glc.userBet(userData[0], fixtures.GetUserBet(20, 10))
	glc.userBet(userData[0], fixtures.GetUserBet(16, 10))
	glc.userBet(userData[1], fixtures.GetUserBet(5, 100))
	glc.userBet(userData[1], fixtures.GetUserBet(4, 5))
	glc.userBet(userData[1], fixtures.GetUserBet(17, 10))

	//開獎及發獎
	glc.getKaiJiangType()
	t.Log("开奖结果是:", glc.betResultType)
	t.Log("winline =>", glc.winline)
	if ok, _ := util.InArray(1, glc.betResultType); ok {
		t.Error("開獎類型有問題", 1, glc.betResultType)
	}
	glc.UserLeave(users[0])
	glc.UserLeave(users[1])
	glc.clearGameData()
	return

	//開十次獎來判定結果
	for i := 0; i < 10; i++ {
		glc.GameStartFunc()
		bettype := int32(util.RndInt(1, 24))
		glc.userBet(userData[1], fixtures.GetUserBet(bettype, 3))
		glc.getKaiJiangType()
		t.Log("开奖结果是:", glc.betResultType)
		if ok, _ := util.InArray(bettype, glc.betResultType); ok {
			t.Error("開十次獎有錯")
		}
		glc.UserLeave(users[0])
		glc.clearGameData()
		time.Sleep(100 * time.Microsecond)
	}
}

//系統當莊, 下注龍開龍, 遊戲流程入桌→下注→開獎→結果
//系統當莊, 下注JZMIN開JZMIN, 遊戲流程入桌→下注→開獎→結果
func TestGameRound2(t *testing.T) {
	var userData []*sgjUserData
	sgj := NewShuiGuoJiLogicByTest()
	sgj.robotzhuang = 0
	sgj.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	sgj.OnGameStart(func() {
	})
	sgj.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	sgj.GameEndFunc = func() {
	}
	sgj.GameStartFunc()
	users := fixtures.GetGameUser()
	users[0].UserType = 1

	for i := 0; i < len(users); i++ {
		userData = append(userData, newUserData(users[i]))
	}

	sgj.userdata = userData

	//安排用戶入座
	sgj.UserEnter(users[0], sgj.GetTable().TableNo, sgj.GetTable().GetFreeSeat().SeatNo)
	//sgj.UserEnter(users[1], sgj.GetTable().TableNo, sgj.GetTable().GetFreeSeat().SeatNo)

	//下注
	sgj.userBet(userData[0], fixtures.GetUserBet(JZMIN, 100))

	//sgj.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//sgj.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//sgj.setNextJiangType(0, 1)
	sgj.nextjiangtype = int32(JZMAX)

	sgj.gameResult()
	//sgj.sendResult()
	calWinLose(sgj)

	if sgj.systotalwin != int64(-900) {
		t.Errorf("systotalwin:%v", sgj.systotalwin)
	}

	if sgj.syslastwin != int64(-900) {
		t.Errorf("syslastwin:%v", sgj.syslastwin)
	}

	if sgj.systaxincome != int64(45) {
		t.Errorf("systaxincome:%v", sgj.systaxincome)
	}

	if sgj.userdata[0].updatemoney != int64(855) {
		t.Errorf("updatemoney:%v", sgj.userdata[0].updatemoney)
	}

	if sgj.userdata[0].win != int64(955) {
		t.Errorf("win:%v", sgj.userdata[0].win)
	}

	//if sgj.zhuangWin != int64(-10000) {
	//	t.Errorf("zhuangWin:%v", sgj.zhuangWin)
	//}

	//sgj.gameEnd()
	sgj.clearGameData()
}

//玩家當莊, 下注, 遊戲流程入桌→下注→開獎→結果
func TestGameRound3(t *testing.T) {
	sgj, users := enterGameUser()

	var userData []*sgjUserData
	for i := 0; i < len(users); i++ {
		userData = append(userData, newUserData(users[i]))
	}

	sgj.userdata = userData

	//安排用戶入座
	sgj.UserEnter(users[0], sgj.GetTable().TableNo, sgj.GetTable().GetFreeSeat().SeatNo)
	//sgj.UserEnter(users[1], sgj.GetTable().TableNo, sgj.GetTable().GetFreeSeat().SeatNo)

	//下注
	sgj.userBet(userData[0], fixtures.GetUserBet(JZMAX, 100))

	//sgj.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//sgj.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//sgj.setNextJiangType(0, 1)
	sgj.nextjiangtype = int32(SXMAX)

	sgj.gameResult()
	//sgj.sendResult()
	calWinLose(sgj)

	if sgj.systotalwin != int64(95) {
		t.Errorf("systotalwin:%v", sgj.systotalwin)
	}

	if sgj.syslastwin != int64(95) {
		t.Errorf("syslastwin:%v", sgj.syslastwin)
	}

	if sgj.systaxincome != int64(5) {
		t.Errorf("systaxincome:%v", sgj.systaxincome)
	}

	if sgj.userdata[0].updatemoney != int64(-100) {
		t.Errorf("updatemoney:%v", sgj.userdata[0].updatemoney)
	}

	if sgj.userdata[0].win != int64(0) {
		t.Errorf("win:%v", sgj.userdata[0].win)
	}

	if sgj.zhuangWin != int64(95) {
		t.Errorf("zhuangWin:%v", sgj.zhuangWin)
	}

	//sgj.gameEnd()
	sgj.clearGameData()
}

func enterGameUser() (*ShuiGuoJiLogic, []*gamelogic.GameUser) {
	sgj := NewShuiGuoJiLogicByTest()
	sgj.robotzhuang = 0
	sgj.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	sgj.OnGameStart(func() {
	})
	sgj.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	sgj.GameEndFunc = func() {
	}
	sgj.GameStartFunc()
	users := fixtures.GetGameUserLhd()
	seat := sgj.GetTable().GetFreeSeat()
	seat.SeatUser = users[0]
	seat.SeatUser = users[1]
	//
	//u := new(longhudouuserdata)
	//u.user = users[0]
	//sgj.userdata = append(sgj.userdata, u)
	//u1 := new(longhudouuserdata)
	//u1.user = users[1]
	//sgj.userdata = append(sgj.userdata, u1)

	sgj.zhuangUserID = 1
	users[1].Score = 5000000
	sgj.startmoney = users[1].Score
	//安排用戶入座
	sgj.UserEnter(users[0], sgj.GetTable().TableNo, sgj.GetTable().GetFreeSeat().SeatNo)
	sgj.UserEnter(users[1], sgj.GetTable().TableNo, sgj.GetTable().GetFreeSeat().SeatNo)
	sgj.UserEnter(users[2], sgj.GetTable().TableNo, sgj.GetTable().GetFreeSeat().SeatNo)
	return sgj, users
}

func TestChiShui(t *testing.T) {
	for j := 0; j < 100; j++ {
		c := 0
		for i := 0; i < 20; i++ {
			v := rand.Int31n(100) + 1
			if v <= 2 {
				c++
			}
		}
		if c > 1 {
			mlog.System("中幾次 %d", c)
		}

	}

}
