package fkeb

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func init() {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewFKEBLogic(t, tfManager, cfgs)
	}
	roomserver.SuperToolCreator = func(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) gamelogic.ISuperToolService {
		return gamelogic.NewBrGameSuperToolServiceBase(room, tfm, dbMsgQueue)
	}
}

// TestGameRound
func TestGameRound(t *testing.T) {

	tgrs := fixtures.NewTestGameRoomServer()
	fkeb := tgrs.GameLogicManager.GetGameLogic(0).(*FKEBLogic)
	if testing.Short() {
		fkeb.SetGameLog(new(gamelog.EmptyGameLog))
	}

	fkeb.robotzhuang = 0
	fkeb.syspool = 0
	fkeb.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	fkeb.OnGameStart(func() {
	})
	fkeb.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	fkeb.GameEndFunc = func() {
	}
	fkeb.GameStartFunc()
	users := fixtures.GetGameUser()
	users[0].UserType = 1

	//安排用戶入座
	fkeb.UserEnter(users[0], fkeb.GetTable().TableNo, fkeb.GetTable().GetFreeSeat().SeatNo)

	//下注
	fkeb.userBet(users[0], fixtures.GetUserBet(SHUN, 10))
	fkeb.userBet(users[0], fixtures.GetUserBet(TIAN, 10))
	fkeb.userBet(users[0], fixtures.GetUserBet(DI, 10))

	// 指定開獎 為 順
	fkeb.nextjiangtype = append(fkeb.nextjiangtype, SHUN)

	// 開獎
	fkeb.kaiJiang()

	// 指定牌型
	fkeb.majiangresult[0].majiang1 = int32(8)
	fkeb.majiangresult[0].majiang2 = int32(2)
	fkeb.majiangresult[1].majiang1 = int32(5)
	fkeb.majiangresult[1].majiang2 = int32(6)
	fkeb.majiangresult[2].majiang1 = int32(7)
	fkeb.majiangresult[2].majiang2 = int32(8)
	fkeb.majiangresult[3].majiang1 = int32(1)
	fkeb.majiangresult[3].majiang2 = int32(2)

	// 結算
	fkeb.sendResult()

	// 檢查結果
	// 這裡檢查 2 8 牌型 和 8 2 牌型 是否有錯誤
	if fkeb.getMajiangsType(fkeb.majiangresult[0]) <= fkeb.getMajiangsType(fkeb.majiangresult[1]) {
		t.Errorf("COMPARE CARD ERROR: %v : %v", fkeb.majiangresult[0], fkeb.majiangresult[1])
	}

	// 遊戲結束
	fkeb.gameEnd()

	// 清除資料
	fkeb.clearGameData()
}

//系統上庄, 庄贏, 真人下注輸
func TestGameRound2(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = 1

	//真人上庄
	glc.zhuangID = 0 //users[1].UserID
	betDone := false

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 2))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
			}

			if *ust.State == Wait {
				if len(glc.userdata) <= 0 {
					return
				}
				//KaiJiangActionDone = true

				if glc.syspool != 285 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -300 {
					t.Errorf("userid[%v] 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("userid[%v] 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.zhuangWin != 300 {
					t.Errorf("系統庄結算錯誤, 實際[%v]", glc.zhuangWin)
				}
			}
		}

	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])

	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, ZHUANG)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//系統上庄, 庄輸, 真人下注贏
func TestGameRound3(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = 1

	//真人上庄
	glc.zhuangID = 0 //users[1].UserID
	betDone := false

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 2))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
			}

			if *ust.State == Wait {
				if len(glc.userdata) <= 0 {
					return
				}

				if glc.syspool != -300 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != 285 {
					t.Errorf("userid[%v] 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 585 {
					t.Errorf("userid[%v] 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.zhuangWin != -300 {
					t.Errorf("系統庄結算錯誤, 實際[%v]", glc.zhuangWin)
				}
			}
		}

	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])

	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN, TIAN)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//系統上庄, 庄輸, 真人下注贏
func TestGameRound10(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = 1
	users[1].UserType = 1

	//真人上庄
	glc.zhuangID = 0 //users[1].UserID
	betDone := false

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 3))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
				glc.userBet(users[2], fixtures.GetUserBet(TIAN, 3))
			}

			if *ust.State == Wait {
				if len(glc.userdata) <= 0 {
					return
				}

				if glc.syspool != 95 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != 185 {
					t.Errorf("userid[%v] 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 585 {
					t.Errorf("userid[%v] 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.zhuangWin != 100 {
					t.Errorf("系統庄結算錯誤, 實際[%v]", glc.zhuangWin)
				}

			}
		}

	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])

	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//真人上庄, 庄贏, 真人下注輸
func TestGameRound4(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = 1
	users[1].UserType = 1

	//真人上庄
	glc.zhuangID = users[1].UserID
	betDone := false

	glc.RegisterBeforeGameStartFunc(func() {
		glc.zhuangID = users[1].UserID
	})

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 2))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
			}

			if *ust.State == Wait {
				if len(glc.userdata) <= 0 {
					return
				}

				if glc.syspool != 0 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -300 {
					t.Errorf("userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.zhuangWin != 285 {
					t.Errorf("真人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])

	glc.zhuangID = users[1].UserID

	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, ZHUANG)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//真人上庄, 庄贏, 真人下注輸, 機器人下注輸
func TestGameRound5(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	//真人上庄
	glc.zhuangID = users[1].UserID
	betDone := false

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 2))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
				glc.userBet(users[2], fixtures.GetUserBet(TIAN, 1))
			}

			if *ust.State == Wait {

				if len(glc.userdata) <= 0 {
					return
				}

				if glc.syspool != -100 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -300 {
					t.Errorf("userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != -100 {
					t.Errorf("userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].updatemoney)
				}

				if glc.zhuangWin != 380 {
					t.Errorf("真人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[1].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, ZHUANG)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//機器人上庄, 庄贏, 真人下注輸, 機器人下注輸
func TestGameRound15(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeRebot
	users[2].UserType = gamelogic.UserTypeRebot

	betDone := false

	glc.RegisterBeforeGameStartFunc(func() {
		glc.zhuangID = users[1].UserID
	})

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 2))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
				glc.userBet(users[2], fixtures.GetUserBet(TIAN, 1))
			}

			if *ust.State == Wait {
				if len(glc.userdata) <= 0 {
					return
				}

				if glc.syspool != 400 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -300 {
					t.Errorf("userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != -100 {
					t.Errorf("userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].updatemoney)
				}

				if glc.zhuangWin != 380 {
					t.Errorf("真人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[1].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, ZHUANG)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//真人上庄, 庄贏, 真人下注贏, 機器人下注輸
func TestGameRound6(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	//真人上庄
	glc.zhuangID = users[1].UserID
	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 2))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 7))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != -700 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != 285 {
					t.Errorf("真人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 585 {
					t.Errorf("真人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != -700 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].updatemoney)
				}

				if glc.userdata[1].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].win)
				}

				if glc.zhuangWin != 380 {
					t.Errorf("真人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[1].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN, TIAN)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//機器人人上庄, 庄贏, 真人下注贏, 機器人下注輸
func TestGameRound16(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeRebot
	users[2].UserType = gamelogic.UserTypeRebot

	//真人上庄
	glc.zhuangID = users[1].UserID
	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 2))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 7))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				if len(glc.userdata) <= 0 {
					return
				}

				KaiJiangActionDone = true

				if glc.syspool != 400 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != 285 {
					t.Errorf("真人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 585 {
					t.Errorf("真人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != -700 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].updatemoney)
				}

				if glc.userdata[1].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].win)
				}

				if glc.zhuangWin != 380 {
					t.Errorf("機器人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[1].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN, TIAN)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//真人上庄, 庄輸, 真人下注贏, 機器人下注輸
func TestGameRound7(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	//真人上庄
	glc.zhuangID = users[1].UserID
	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 2))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 1))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 1))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				if len(glc.userdata) <= 0 {
					return
				}

				KaiJiangActionDone = true

				if glc.syspool != -100 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != 285 {
					t.Errorf("真人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 585 {
					t.Errorf("真人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != -100 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].updatemoney)
				}

				if glc.userdata[1].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].win)
				}

				if glc.zhuangWin != -200 {
					t.Errorf("真人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[1].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN, TIAN)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//真人上庄, 庄輸, 真人下注輸, 機器人下注贏
func TestGameRound8(t *testing.T) {

	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	//真人上庄
	glc.zhuangID = users[1].UserID
	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 1))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 3))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 4))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != 400 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -205 {
					t.Errorf("真人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 195 {
					t.Errorf("真人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != 380 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].updatemoney)
				}

				if glc.userdata[1].win != 780 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].win)
				}

				if glc.zhuangWin != -200 {
					t.Errorf("真人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[1].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

//真人上庄, 庄輸, 真人下注贏, 機器人下注輸
func TestGameRound9(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	//真人上庄
	glc.zhuangID = users[1].UserID
	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[0], fixtures.GetUserBet(SHUN, 1))
				glc.userBet(users[0], fixtures.GetUserBet(TIAN, 3))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 2))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != -200 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != 380 {
					t.Errorf("真人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 780 {
					t.Errorf("真人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != -200 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].updatemoney)
				}

				if glc.userdata[1].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].win)
				}

				if glc.zhuangWin != -200 {
					t.Errorf("真人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[1].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN, TIAN)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

// 機器人上庄, 庄輸, 真人下注輸比較多，真人贏比較少, 機器人下注贏
func TestGameRound11(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeRebot
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	//機器人上庄
	glc.zhuangID = users[0].UserID // 指定機器人上莊
	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[1], fixtures.GetUserBet(SHUN, 1)) // 機器人 莊家 // win
				glc.userBet(users[1], fixtures.GetUserBet(TIAN, 3)) // 真人 // lose
				glc.userBet(users[2], fixtures.GetUserBet(DI, 2))   // 機器人 // win
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != 200 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -205 {
					t.Errorf("真人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 195 {
					t.Errorf("真人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}
				mlog.System("%v %v", glc.userdata[0].user.UserID, glc.userdata[0].user.UserType)

				if glc.zhuangWin != 0 {
					t.Errorf("真人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[1].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

// 機器人上庄, 庄輸, 真人機器人通賠
func TestGameRound12(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeRebot
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[1], fixtures.GetUserBet(SHUN, 1))
				glc.userBet(users[1], fixtures.GetUserBet(TIAN, 3))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 2))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != -400 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != 380 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 780 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}
				mlog.System("%v %v", glc.userdata[0].user.UserID, glc.userdata[0].user.UserType)

				if glc.zhuangWin != -600 {
					t.Errorf("機器人莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[0].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, SHUN, TIAN, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

// 機器人上庄, 庄贏, 真人輸, 機器人贏
func TestGameRound13(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeRebot
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[1], fixtures.GetUserBet(SHUN, 1))
				glc.userBet(users[1], fixtures.GetUserBet(TIAN, 3))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 2))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != 200 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -400 {
					t.Errorf("真人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("真人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != 190 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[1].user.UserID, glc.userdata[1].updatemoney)
				}

				if glc.zhuangWin != 190 {
					t.Errorf("機器人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[0].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

// 機器人上庄, 庄贏, 真人輸, 機器人贏
func TestGameRound14(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeRebot
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	betDone := false
	KaiJiangActionDone := false
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[1], fixtures.GetUserBet(SHUN, 1))
				glc.userBet(users[1], fixtures.GetUserBet(TIAN, 3))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 2))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != 200 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -400 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}
				mlog.System("%v %v", glc.userdata[0].user.UserID, glc.userdata[0].user.UserType)

				if glc.zhuangWin != 190 {
					t.Errorf("機器人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[0].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

// 機器人上庄,
func TestGameRound17(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeRebot
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	betDone := false
	KaiJiangActionDone := false

	glc.RegisterBeforeGameStartFunc(func() {
		glc.zhuangID = users[0].UserID
	})

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[1], fixtures.GetUserBet(SHUN, 1))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 8))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != 100 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -100 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}
				mlog.System("%v %v", glc.userdata[0].user.UserID, glc.userdata[0].user.UserType)

				if glc.zhuangWin != -700 {
					t.Errorf("機器人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[0].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

// 真人上庄,
func TestGameRound18(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	betDone := false
	KaiJiangActionDone := false

	glc.RegisterBeforeGameStartFunc(func() {
		glc.zhuangID = users[0].UserID
	})

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[1], fixtures.GetUserBet(SHUN, 1))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 8))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != 800 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -100 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}
				mlog.System("%v %v", glc.userdata[0].user.UserID, glc.userdata[0].user.UserType)

				if glc.zhuangWin != -700 {
					t.Errorf("機器人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[0].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

// 機器人上庄,
func TestGameRound19(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeRebot
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	betDone := false
	KaiJiangActionDone := false

	glc.RegisterBeforeGameStartFunc(func() {
		glc.zhuangID = users[0].UserID
	})

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[1], fixtures.GetUserBet(SHUN, 8))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 1))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != 700 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -800 {
					t.Errorf("真人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}
				mlog.System("%v %v", glc.userdata[0].user.UserID, glc.userdata[0].user.UserType)

				if glc.zhuangWin != 665 {
					t.Errorf("機器人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[0].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

// 真人上庄,
func TestGameRound20(t *testing.T) {
	endChannel := make(chan bool)
	confInit()

	tgrs, glc, endChannel, _ := IntegrationTestInit(endChannel, t)
	glc.robotzhuang = 0
	glc.syspool = 0

	users := fixtures.GetGameUser()
	users[0].UserType = gamelogic.UserTypeNormal
	users[1].UserType = gamelogic.UserTypeNormal
	users[2].UserType = gamelogic.UserTypeRebot

	betDone := false
	KaiJiangActionDone := false

	glc.RegisterBeforeGameStartFunc(func() {
		glc.zhuangID = users[0].UserID
	})

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.FKEB_GameMessageClassID_FKEBStateID) {
			ust := m.MsgData.(*netproto.FKEB_Sate)
			if *ust.State == Bet && !betDone {
				betDone = true
				//下注
				glc.userBet(users[1], fixtures.GetUserBet(SHUN, 8))
				glc.userBet(users[2], fixtures.GetUserBet(DI, 1))
			}

			if *ust.State == Wait && !KaiJiangActionDone {
				KaiJiangActionDone = true

				if glc.syspool != 100 {
					t.Errorf("系統水池錯誤, 實際[%v]", glc.syspool)
				}

				if glc.syslastwin != 0 {
					t.Errorf("上一局收益錯誤, 實際[%v]", glc.syslastwin)
				}

				if glc.userdata[0].updatemoney != -800 {
					t.Errorf("機器人userid[%v]下注 結算updatemoney錯誤錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != 0 {
					t.Errorf("機器人userid[%v]下注 結算win錯誤, 實際[%v]", glc.userdata[0].user.UserID, glc.userdata[0].win)
				}
				mlog.System("%v %v", glc.userdata[0].user.UserID, glc.userdata[0].user.UserType)

				if glc.zhuangWin != 665 {
					t.Errorf("機器人上莊userid[%v]庄結算錯誤, 實際[%v]", glc.zhuangID, glc.zhuangWin)
				}
			}
		}
	}

	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	users[2].UserState.OnUserQueue(users[2], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	glc.zhuangID = users[0].UserID
	// 指定開獎
	glc.nextjiangtype = append(glc.nextjiangtype, DI)

	select {
	case _ = <-endChannel:
		mlog.System("fkeb game end")
	}
}

func confInit() {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 1000 //40
	conf.TestMode = true
	mlog.GetCLogger().Enable(false)
	mylog.InitDatabaseQueue("", 0, 0, 0, mlog.NewDefaultLogger())
}

// 流程整合測試共用初始化
func IntegrationTestInit(endChannel chan bool, t *testing.T) (*fixtures.TestGameRoomServer, *FKEBLogic, chan bool, []*gamelogic.GameUser) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*FKEBLogic)
	if testing.Short() {
		glc.SetGameLog(new(gamelog.EmptyGameLog))
	}
	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.SendFishNotifyFunc = func(m *msg.Message) {

	}

	glc.OnGameStart(func() {
	})

	_, robots := fixtures.GetGameUsers()
	robotManager := tgrs.RobotService.GetRobotManager()
	if robotManager.RobotCount() != 0 {
		t.Errorf("初始化, 機器人數量應為0, RobotCount=>%d", robotManager.RobotCount())
	}
	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		//robotManager.AddRobot(robot)
	}

	glc.TableManager.Room.Seq = proto.Int64(0)

	return tgrs, glc, endChannel, robots
}
