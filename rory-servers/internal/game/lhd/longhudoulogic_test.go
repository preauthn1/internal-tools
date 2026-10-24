package lhd

import (
	"testing"

	"example.invalid/stretchr/testify/assert"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func init() {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewLongHuDouLogic(t, tfManager, cfgs)
	}
	roomserver.SuperToolCreator = func(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) gamelogic.ISuperToolService {
		return gamelogic.NewBrGameSuperToolServiceBase(room, tfm, dbMsgQueue)
	}
}

func TestCal(t *testing.T) {
	mlog.System("%d", 11*15/10)
	a := make(map[string]string)
	a["1"] = "1"
	a["2"] = "2"

	b := a
	a = make(map[string]string)
	mlog.System("==== a [%v]", a)
	mlog.System("==== b [%v]", b)
}

func ChangeConfig1(lhd *LongHudouLogic) {
	lhd.Configs.Configs[gamelogic.ConfigFileNameGameServer].Sections[gamelogic.ConfigSectionNameGameServer]["taxscale"] = "0"
	lhd.HEandKill = 0
}

func ChangeConfig2(lhd *LongHudouLogic) {
	lhd.Configs.Configs[gamelogic.ConfigFileNameGameServer].Sections[gamelogic.ConfigSectionNameGameServer]["taxscale"] = "5"
	lhd.HEandKill = 1
}

// 以下測試使用taxscale=0 HEandKill=0

//系統當莊, 下注龍開和, 遊戲流程入桌→下注→開獎→結果
func TestGameRound1(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 1
	conf.TestMode = true

	tgrs := fixtures.NewTestGameRoomServer()
	lhd := tgrs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	ChangeConfig1(lhd)

	lhd.robotzhuang = 0
	lhd.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	lhd.OnGameStart(func() {
	})
	lhd.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	lhd.GameEndFunc = func() {
	}
	lhd.GameStartFunc()
	users := fixtures.GetGameUser()
	users[0].UserType = 1

	//安排用戶入座
	lhd.UserEnter(users[0], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)
	//lhd.UserEnter(users[1], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HE)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(100) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(100) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.syslastwin != int64(100) {
		t.Errorf("syslastwin:%v", lhd.syslastwin)
	}

	if lhd.systaxincome != int64(0) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	if lhd.userdata[0].updatemoney != int64(-100) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	if lhd.userdata[0].win != int64(0) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(100) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

// 以下測試使用taxscale=5 HEandKill=1

//系統當莊, 下注龍開龍, 遊戲流程入桌→下注→開獎→結果
func TestGameRound(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 1
	conf.TestMode = true

	tgrs := fixtures.NewTestGameRoomServer()
	lhd := tgrs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	ChangeConfig2(lhd)

	lhd.robotzhuang = 0
	lhd.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	lhd.OnGameStart(func() {
	})
	lhd.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	lhd.GameEndFunc = func() {
	}
	lhd.GameStartFunc()
	users := fixtures.GetGameUser()
	users[0].UserType = 1

	//安排用戶入座
	lhd.UserEnter(users[0], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)
	//lhd.UserEnter(users[1], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(LONG)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(-100) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(-100) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.syslastwin != int64(-100) {
		t.Errorf("syslastwin:%v", lhd.syslastwin)
	}

	if lhd.systaxincome != int64(5) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	if lhd.userdata[0].updatemoney != int64(95) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	if lhd.userdata[0].win != int64(195) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(-100) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//系統當莊, 下注和開和, 遊戲流程入桌→下注→開獎→結果
func TestGameRound7(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 1
	conf.TestMode = true

	tgrs := fixtures.NewTestGameRoomServer()
	lhd := tgrs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	ChangeConfig2(lhd)

	lhd.robotzhuang = 0
	lhd.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	lhd.OnGameStart(func() {
	})
	lhd.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	lhd.GameEndFunc = func() {
	}
	lhd.GameStartFunc()
	users := fixtures.GetGameUser()
	users[0].UserType = 1

	//安排用戶入座
	lhd.UserEnter(users[0], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)
	//lhd.UserEnter(users[1], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HE)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	assert.Equal(t, int64(-800), lhd.SysPool.GetPool())
	assert.Equal(t, int64(-800), lhd.systotalwin)
	assert.Equal(t, int64(-800), lhd.syslastwin)
	assert.Equal(t, int64(40), lhd.systaxincome)
	assert.Equal(t, int64(760), lhd.userdata[0].updatemoney)
	assert.Equal(t, int64(860), lhd.userdata[0].win)
	assert.Equal(t, int64(-800), lhd.zhuangWin)

	lhd.gameEnd()
	lhd.clearGameData()
}

//系統當莊, 下注龍和 開和, 遊戲流程入桌→下注→開獎→結果
func TestGameRound8(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 1
	conf.TestMode = true

	tgrs := fixtures.NewTestGameRoomServer()
	lhd := tgrs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	ChangeConfig2(lhd)

	lhd.robotzhuang = 0
	lhd.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	lhd.OnGameStart(func() {
	})
	lhd.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	lhd.GameEndFunc = func() {
	}
	lhd.GameStartFunc()
	lhd.zhuangID = 0
	users := fixtures.GetGameUser()
	users[0].UserType = 1

	//安排用戶入座
	lhd.UserEnter(users[0], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)
	//lhd.UserEnter(users[1], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HE)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(-800) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(-800) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.syslastwin != int64(-800) {
		t.Errorf("syslastwin:%v", lhd.syslastwin)
	}

	if lhd.systaxincome != int64(40) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	if lhd.userdata[0].updatemoney != int64(760) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	if lhd.userdata[0].win != int64(960) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(-800) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//系統當莊, 下注龍和 開龍, 遊戲流程入桌→下注→開獎→結果
func TestGameRound9(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 1
	conf.TestMode = true

	tgrs := fixtures.NewTestGameRoomServer()
	lhd := tgrs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	ChangeConfig2(lhd)

	lhd.robotzhuang = 0
	lhd.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	lhd.OnGameStart(func() {
	})
	lhd.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	lhd.GameEndFunc = func() {
	}
	lhd.GameStartFunc()
	lhd.zhuangID = 0
	users := fixtures.GetGameUser()
	users[0].UserType = 1

	//安排用戶入座
	lhd.UserEnter(users[0], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)
	//lhd.UserEnter(users[1], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(LONG)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(0) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(0) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.syslastwin != int64(0) {
		t.Errorf("syslastwin:%v", lhd.syslastwin)
	}

	if lhd.systaxincome != int64(0) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	if lhd.userdata[0].updatemoney != int64(0) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	if lhd.userdata[0].win != int64(200) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(0) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//玩家當莊, 下注龍開虎, 遊戲流程入桌→下注→開獎→結果
func TestGameRound2(t *testing.T) {
	lhd, users := enterGameUser()

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HU)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()
	assert.Equal(t, int64(0), lhd.SysPool.GetPool())
	assert.Equal(t, int64(95), lhd.systotalwin)
	assert.Equal(t, int64(5), lhd.systaxincome)
	assert.Equal(t, int64(-100), lhd.userdata[0].updatemoney)
	assert.Equal(t, int64(0), lhd.userdata[0].win)
	assert.Equal(t, int64(95), lhd.zhuangWin)

	lhd.gameEnd()
	lhd.clearGameData()
}

//玩家當莊, 下注虎開虎, 遊戲流程入桌→下注→開獎→結果
func TestGameRound6(t *testing.T) {
	lhd, users := enterGameUser()

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HU)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()
	assert.Equal(t, int64(-100), lhd.systotalwin)
	assert.Equal(t, int64(5), lhd.systaxincome)
	assert.Equal(t, int64(95), lhd.userdata[0].updatemoney)
	assert.Equal(t, int64(195), lhd.userdata[0].win)
	assert.Equal(t, int64(-100), lhd.zhuangWin)

	lhd.gameEnd()
	lhd.clearGameData()
}

//玩家當莊, 下注和開和, 遊戲流程入桌→下注→開獎→結果
func TestGameRound5(t *testing.T) {
	lhd, users := enterGameUser()

	//下注
	//lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HE)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(0) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(-800) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.systaxincome != int64(40) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	if lhd.userdata[0].updatemoney != int64(760) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	if lhd.userdata[0].win != int64(860) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(-800) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//玩家當莊, 下注龍和 開龍, 遊戲流程入桌→下注→開獎→結果
func TestGameRound3(t *testing.T) {
	lhd, users := enterGameUser()

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(LONG)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(0) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(0) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.syslastwin != int64(0) {
		t.Errorf("syslastwin:%v", lhd.syslastwin)
	}

	if lhd.systaxincome != int64(0) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	if lhd.userdata[0].updatemoney != int64(0) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	if lhd.userdata[0].win != int64(200) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(0) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//robot當莊, 下注龍和 開和, 遊戲流程入桌→下注→開獎→結果
func TestGameRound11(t *testing.T) {
	lhd, users := enterGameUser()
	users[1].UserType = int32(gamelogic.UserTypeRebot)
	//lhd.userdata[1].user = users[1]

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HE)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(-800) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(-800) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.systaxincome != int64(0) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	//玩家輸贏
	if lhd.userdata[0].updatemoney != int64(760) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	//莊輸贏
	if lhd.userdata[0].win != int64(960) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(-800) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//robot當莊, 下注龍和 開龍, 遊戲流程入桌→下注→開獎→結果
func TestGameRound12(t *testing.T) {
	lhd, users := enterGameUser()
	users[1].UserType = int32(gamelogic.UserTypeRebot)

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(LONG)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(0) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(0) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.systaxincome != int64(0) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	//玩家輸贏
	if lhd.userdata[0].updatemoney != int64(0) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	//莊輸贏
	if lhd.userdata[0].win != int64(200) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(0) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//robot當莊, 下注龍和 開虎, 遊戲流程入桌→下注→開獎→結果
func TestGameRound13(t *testing.T) {
	lhd, users := enterGameUser()
	users[1].UserType = int32(gamelogic.UserTypeRebot)

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HU)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(190) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(190) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.systaxincome != int64(10) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	//玩家輸贏
	if lhd.userdata[0].updatemoney != int64(-200) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	//莊輸贏
	if lhd.userdata[0].win != int64(0) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(190) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//robot當莊, robot下注龍和 開虎, 遊戲流程入桌→下注→開獎→結果
func TestGameRound16(t *testing.T) {
	lhd, users := enterGameUser()
	users[0].UserType = int32(gamelogic.UserTypeRebot)
	users[1].UserType = int32(gamelogic.UserTypeRebot)

	//下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))
	//lhd.userBet(users[0], fixtures.GetUserBet(HU, 100))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HU)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(0) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(190) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.systaxincome != int64(10) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	//玩家輸贏
	if lhd.userdata[0].updatemoney != int64(-200) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	//莊輸贏
	if lhd.userdata[0].win != int64(0) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(190) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//玩家當莊
//玩家1, 下注龍和 開和, 遊戲流程入桌→下注→開獎→結果
//玩家2, 下注虎和 開和, 遊戲流程入桌→下注→開獎→結果
func TestGameRound14(t *testing.T) {
	lhd, users := enterGameUser()

	//player1 下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))

	//player2 下注
	lhd.userBet(users[2], fixtures.GetUserBet(HU, 50))
	lhd.userBet(users[2], fixtures.GetUserBet(HE, 50))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HE)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(0) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(-1200) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.systaxincome != int64(60) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	//玩家輸贏
	if lhd.userdata[0].updatemoney != int64(760) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	//莊輸贏
	if lhd.userdata[0].win != int64(960) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	if lhd.zhuangWin != int64(-1200) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

//玩家當莊
//玩家1, 下注龍和 開虎, 遊戲流程入桌→下注→開獎→結果
//玩家2, 下注虎和 開虎, 遊戲流程入桌→下注→開獎→結果
func TestGameRound15(t *testing.T) {
	lhd, users := enterGameUser()

	//player1 下注
	lhd.userBet(users[0], fixtures.GetUserBet(LONG, 100))
	lhd.userBet(users[0], fixtures.GetUserBet(HE, 100))

	//player2 下注
	lhd.userBet(users[2], fixtures.GetUserBet(HU, 100))
	lhd.userBet(users[2], fixtures.GetUserBet(HE, 50))

	//lhd.setNextJiangType(0, 1)
	lhd.nextjiangtype = int32(HU)

	lhd.kaiJiang()
	//lhd.sendResult()
	lhd.calWinLose()

	if lhd.SysPool.GetPool() != int64(0) {
		t.Errorf("syspool:%v", lhd.SysPool.GetPool())
	}

	if lhd.systotalwin != int64(143) {
		t.Errorf("systotalwin:%v", lhd.systotalwin)
	}

	if lhd.systaxincome != int64(9) {
		t.Errorf("systaxincome:%v", lhd.systaxincome)
	}

	//玩家1輸贏
	if lhd.userdata[0].updatemoney != int64(-200) {
		t.Errorf("updatemoney:%v", lhd.userdata[0].updatemoney)
	}

	//玩家1輸贏
	if lhd.userdata[0].win != int64(0) {
		t.Errorf("win:%v", lhd.userdata[0].win)
	}

	//玩家2輸贏
	if lhd.userdata[1].updatemoney != int64(48) {
		t.Errorf("updatemoney:%v", lhd.userdata[1].updatemoney)
	}

	//玩家2輸贏
	if lhd.userdata[1].win != int64(198) {
		t.Errorf("win:%v", lhd.userdata[1].win)
	}

	//莊輸贏
	if lhd.zhuangWin != int64(143) {
		t.Errorf("zhuangWin:%v", lhd.zhuangWin)
	}

	lhd.gameEnd()
	lhd.clearGameData()
}

func enterGameUser() (*LongHudouLogic, []*gamelogic.GameUser) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 1
	conf.TestMode = true

	tgrs := fixtures.NewTestGameRoomServer()
	lhd := tgrs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	ChangeConfig2(lhd)

	lhd.robotzhuang = 0
	lhd.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	lhd.OnGameStart(func() {
	})
	lhd.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	lhd.GameEndFunc = func() {
	}
	lhd.GameStartFunc()
	users := fixtures.GetGameUserLhd()
	seat := lhd.GetTable().GetFreeSeat()
	seat.SeatUser = users[0]
	seat.SeatUser = users[1]
	//
	//u := new(longhudouuserdata)
	//u.user = users[0]
	//lhd.userdata = append(lhd.userdata, u)
	//u1 := new(longhudouuserdata)
	//u1.user = users[1]
	//lhd.userdata = append(lhd.userdata, u1)

	lhd.zhuangID = 1
	lhd.startmoney = users[1].Score
	//安排用戶入座
	lhd.UserEnter(users[0], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)
	lhd.UserEnter(users[1], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)
	lhd.UserEnter(users[2], lhd.GetTable().TableNo, lhd.GetTable().GetFreeSeat().SeatNo)
	return lhd, users
}

//
//func NewLongHudouLogicByTest() *LongHudouLogic {
//	_, filename, _, _ := runtime.Caller(1)
//	configDir := path.Dir(filename) + "/../../config"
//	cfgs := new(config.ConfigGroup)
//	cfgs.LoadConfig(configDir, 0)
//	cfg := cfgs.GetConfig(gamelogic.CONFIGFILENAME_GAMESERVER)
//	serverInfo := cfgs.GetConfig(gamelogic.CONFIGFILENAME_SERVERINFO)
//
//	//获取配置信息
//	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
//	rumanager := gamelogic.NewRoomUserManager(new(player.PlayerManager))
//	room := gamelogic.NewRoom(servercfg, rumanager)
//	table := room.GetTableByNo(0)
//
//	evtdispatcher := event.NewAsyncEventDispatcher(1)
//	tfManager := event.NewTimerFuncManager(evtdispatcher)
//
//	glc := new(LongHudouLogic)
//
//	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(table, tfManager, cfgs)
//	glc.getlhdconfig(cfgs)
//	glc.gamestate = Wait
//	glc.zhuangID = 0
//	glc.nextjiangtype = -1
//	glc.syslastwin = 0
//	glc.systotalwin = 0
//	glc.resultcardlist = new(netproto.LHD_ResultCardList)
//	glc.startsyspool = glc.syspool
//
//	//设置桌子对应的逻辑实例
//	for _, t := range room.Tables {
//		t.GameLogic = new(LongHudouLogic)
//	}
//
//	glc.RegisterAfterGameStartFunc(func() {
//		glc.clearGameData() //清除数据
//		glc.fapai()
//		glc.huangZhuang()
//		glc.GameLog.AddGameStep("庄家%d，剩余%d", glc.zhuangID, glc.lastturn)
//		glc.zhuShui()
//		glc.broadCastZhuangInfo()
//		//glc.startBet()
//	})
//
//	return glc
//}

//sendUserResult 给每个玩家发送开奖结果

//TODO 驗證莊贏, 系統水池, 系統總營收, 上一局收益

//zhuShui 注水到下注额

//kaiJiang 开始开奖

//userBet 用户下注处理函数

//getRealUserWinByType 获取系统赢钱数

//TODO
//玩家1, 下注龍和 開虎, 遊戲流程入桌→下注→開獎→結果
//robot, 下注龍和 開虎, 遊戲流程入桌→下注→開獎→結果
