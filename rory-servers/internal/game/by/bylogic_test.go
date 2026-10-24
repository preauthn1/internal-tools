package by

import (
	"fmt"
	"strconv"
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

//获取房间配置文件
func TestGetRoomTable(t *testing.T) {
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	tgrs.GameLogicManager.GetGameLogic(0)

}

func confInit() {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 100 //40
	conf.TestMode = true
	mylog.InitDatabaseQueue("", 0, 0, 0, mlog.NewDefaultLogger())
}

//發射砲彈, 未打到魚
func TestBoomFish(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	glc.TableManager.Room.FishSysPoolData = new(netproto.FishSysPoolData)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	boomId := int32(0)
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.BY_GameMessageClassID_FrameS2CID) {
			ust := m.MsgData.(*netproto.BY_FrameS2C)
			if ust.FishWave.GetFishWaveID() > int32(0) {
				glc.gameEnd()
			}

			if len(ust.NewFish) > 0 {
				glc.userdata[0].rate = 10
				//補魚
				btmsg := new(msg.Message)
				btmsg.BClassID = int32(netproto.MessageBClassID_Game)
				btmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameC2SID)
				btdata := new(netproto.BY_FrameC2S)
				btdata.FrameID = ust.FrameID
				//用戶發射
				userShoot := new(netproto.BY_ShootC2S)
				pos := new(netproto.BY_Position)
				pos.X = proto.Float32(float32(5))
				pos.Y = proto.Float32(float32(5))
				userShoot.Pos = pos
				userShoot.FishID = ust.NewFish[0].FishID
				btdata.Shoot = userShoot

				//用戶砲彈打到魚
				userBoom := new(netproto.BY_BoomC2S)
				boomId++
				userBoom.BoomID = proto.Int32(boomId)
				userBoom.BoomRate = proto.Int32(1)
				userBoom.FishID = []int32{ust.NewFish[0].GetFishID()}
				userBoom.UserID = proto.Int32(users[1].UserID)
				btdata.Boom = userBoom
				userChangeRate := new(netproto.BY_UserChangeRateC2S)
				userChangeRate.RateAfterChange = proto.Int32(1)
				btdata.UserChangeRate = userChangeRate
				suoDing := new(netproto.BY_SuoDingC2S)
				suoDing.SuoDing = proto.Bool(false)
				btdata.SuoDing = suoDing
				allFish := new(netproto.BY_AllFish)
				btdata.AllFish = allFish
				btdata.UseIce = proto.Bool(false)
				btmsg.MsgData, _ = proto.Marshal(btdata)

				glc.GameMessage(users[1], btmsg)
				mlog.Debug("UpdateScore[%v]", users[1].UpdateScore)
				ud := glc.getuserdatabyuserid(users[1].UserID)
				boomOdds := glc.boommanager[userChangeRate.GetRateAfterChange()-1].boomOdds

				glc.TimerFuncManager.NewTimer(10, func() {

					if boomId < 46 {
						if ud.money != int64(boomId*-1*boomOdds*glc.basescore) {
							//t.Errorf("炮彈金幣money錯誤! money[%v] boomid[%v] boomOdds[%v] basescore[%v]", ud.money, boomId, boomOdds, glc.basescore)
						}
					}

					if ud.allwin != int64(boomId*-1*boomOdds*glc.basescore) {
						t.Error("炮彈金幣allwin錯誤!")
					}

					if ud.winmoney != int64(boomId*-1*boomOdds*glc.basescore) {
						t.Error("炮彈金幣winmoney錯誤!")
					}

					endChannel <- true
				})
			}
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("by game end")
	}
}

//打到魚, 金額與捕魚數據統計
func TestBoomFish6(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	glc.TableManager.Room.FishSysPoolData = new(netproto.FishSysPoolData)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	glc.TableManager.Room.FishSysPoolData = new(netproto.FishSysPoolData)

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(5500)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)

	//打到魚
	fishid := make([]int32, 0)
	fishid = append(fishid, 1)
	bydata := new(byfishdata)
	bydata.fishSeq = 2
	bydata.fishtype = 2
	bydata.fishcount = 1
	glc.fishmanager.fishdata[fishid[0]] = bydata
	glc.fish(fishid, ud, false, nil)

	if ud.money != 1680 {
		t.Errorf("money[%v]錯誤, ", ud.money)
	}
	if ud.allwin != 1680 {
		t.Errorf("allwin[%v]錯誤", ud.allwin)
	}
	if ud.winmoney != 1680 {
		t.Errorf("winmoney[%v]錯誤", ud.winmoney)
	}
	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != 1680 {
		t.Errorf("SysSumWinLoss[%v]錯誤", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

func BenchmarkBoomFish6(t *testing.B) {
	//重置时间点
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		TestBoomFish6(nil)
	}
}

func BenchmarkBoomFish15(t *testing.B) {
	//重置时间点
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		TestBoomFish15(nil)
	}
}

//打到魚, player pool +100
func TestBoomFish15(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(10)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(10)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	playerpool := new(netproto.FishPlayerPoolDetail)
	playerpool.FlagID = proto.Int32(1)
	pdetail := new(netproto.FishPlayerPoolDetailList)
	pdetail.UserID = proto.Int32(users[1].UserID)
	pdetail.HitMoney = proto.Int64(750)
	pdetail.SpendMoney = proto.Int64(1600)
	playerpool.FishRewardList = append(playerpool.FishRewardList, pdetail)
	glc.TableManager.Room.FishPlayerRewardDetail[users[1].UserID] = pdetail

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	//glc.fish(fishid2, ud, false)

	//rate = by.getRate(fishdata.fishtype, ud)
	newRate := glc.getRate(bydata2.fishtype, ud)
	newRate = glc.adjectPlayerPoolRate(newRate, ud)

	if newRate != 2.06 {
		t.Errorf("個人水池player pool +0.1 攻擊率錯誤 rate[%v]", newRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != 10 ||
		glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, player pool +60
func TestBoomFish16(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FishSysPoolData = new(netproto.FishSysPoolData)
	glc.TableManager.Room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(10)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(10)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	playerpool := new(netproto.FishPlayerPoolDetail)
	playerpool.FlagID = proto.Int32(1)
	pdetail := new(netproto.FishPlayerPoolDetailList)
	pdetail.UserID = proto.Int32(users[1].UserID)
	pdetail.HitMoney = proto.Int64(980)
	pdetail.SpendMoney = proto.Int64(1600)
	playerpool.FishRewardList = append(playerpool.FishRewardList, pdetail)
	glc.TableManager.Room.FishPlayerRewardDetail[users[1].UserID] = pdetail

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	//glc.fish(fishid2, ud, false)

	//rate = by.getRate(fishdata.fishtype, ud)
	newRate := glc.getRate(bydata2.fishtype, ud)
	mlog.Debug("newRate[%v]", newRate)
	newRate = glc.adjectPlayerPoolRate(newRate, ud)

	if newRate != 2.02 {
		t.Errorf("個人水池 player pool +0.6 攻擊率錯誤 rate[%v]", newRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != 10 ||
		glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, player pool +30
func TestBoomFish17(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(10)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(10)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	playerpool := new(netproto.FishPlayerPoolDetail)
	playerpool.FlagID = proto.Int32(1)
	pdetail := new(netproto.FishPlayerPoolDetailList)
	pdetail.UserID = proto.Int32(users[1].UserID)
	pdetail.HitMoney = proto.Int64(1150)
	pdetail.SpendMoney = proto.Int64(1600)
	playerpool.FishRewardList = append(playerpool.FishRewardList, pdetail)
	glc.TableManager.Room.FishPlayerRewardDetail[users[1].UserID] = pdetail

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	//glc.fish(fishid2, ud, false)

	//rate = by.getRate(fishdata.fishtype, ud)
	newRate := glc.getRate(bydata2.fishtype, ud)
	newRate = glc.adjectPlayerPoolRate(newRate, ud)

	if newRate != 1.99 {
		t.Errorf("個人水池 +30 攻擊率錯誤 rate[%v]", newRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != 10 ||
		glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, player pool -30
func TestBoomFish18(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(10)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(10)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	playerpool := new(netproto.FishPlayerPoolDetail)
	playerpool.FlagID = proto.Int32(1)
	pdetail := new(netproto.FishPlayerPoolDetailList)
	pdetail.UserID = proto.Int32(users[1].UserID)
	pdetail.HitMoney = proto.Int64(2100)
	pdetail.SpendMoney = proto.Int64(1700)
	playerpool.FishRewardList = append(playerpool.FishRewardList, pdetail)
	glc.TableManager.Room.FishPlayerRewardDetail[users[1].UserID] = pdetail

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	//glc.fish(fishid2, ud, false)

	//rate = by.getRate(fishdata.fishtype, ud)
	newRate := glc.getRate(bydata2.fishtype, ud)
	newRate = glc.adjectPlayerPoolRate(newRate, ud)

	rate, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", newRate), 64)

	if rate != 1.93 {
		t.Errorf("個人水池 player pool -30 攻擊率錯誤 rate[%v]", rate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != 10 ||
		glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, player pool -200
func TestBoomFish20(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(10)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(10)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	playerpool := new(netproto.FishPlayerPoolDetail)
	playerpool.FlagID = proto.Int32(1)
	pdetail := new(netproto.FishPlayerPoolDetailList)
	pdetail.UserID = proto.Int32(users[1].UserID)
	pdetail.HitMoney = proto.Int64(2800)
	pdetail.SpendMoney = proto.Int64(1600)
	playerpool.FishRewardList = append(playerpool.FishRewardList, pdetail)
	glc.TableManager.Room.FishPlayerRewardDetail[users[1].UserID] = pdetail

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	//glc.fish(fishid2, ud, false)

	//rate = by.getRate(fishdata.fishtype, ud)
	newRate := glc.getRate(bydata2.fishtype, ud)
	newRate = glc.adjectPlayerPoolRate(newRate, ud)

	rate, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", newRate), 64)
	if rate != 1.76 {
		t.Errorf("個人水池攻擊率錯誤 rate[%v]", rate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != 10 ||
		glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, SysPool +0.2
func TestBoomFish7(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-15000)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(-15000)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)

	//打到魚
	fishid := make([]int32, 0)
	fishid = append(fishid, 1)
	bydata := new(byfishdata)
	bydata.fishSeq = 2
	bydata.fishtype = 2
	bydata.fishcount = 1
	glc.fishmanager.fishdata[fishid[0]] = bydata
	glc.fish(fishid, ud, false, nil)

	if glc.syspoolAdjustRate != 0.2 {
		t.Errorf("系統調整 SysPool +200 攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

}

//打到魚, SysPool +0.125
func TestBoomFish8(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.TableManager.Room.FishSysPoolData = new(netproto.FishSysPoolData)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-465)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(-465)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(1000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != 0.125 {
		t.Errorf("系統調整 +0.125 攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	//打到魚
	fishid := make([]int32, 0)
	fishid = append(fishid, 1)
	bydata := new(byfishdata)
	bydata.fishSeq = 2
	bydata.fishtype = 2
	bydata.fishcount = 1
	glc.fishmanager.fishdata[fishid[0]] = bydata
	glc.TableManager.Room.FishSysPoolData.SysSumWinLoss = proto.Int64(2000)
	glc.TableManager.Room.AfterAdjustSyspoolWinLoss = proto.Int64(2000)
	glc.fish(fishid, ud, false, nil)

	if glc.syspoolAdjustRate != -0.5 {
		t.Errorf("系統調整 -0.5 攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != 2840 ||
		glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, SysPool +0.08
func TestBoomFish9(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-370)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(-370)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != 0.08 {
		t.Errorf("系統調整攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, SysPool 0.05
func TestBoomFish10(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-290)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(0)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(0)
	syspool.PoolAmount = proto.Int64(0)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != 0.05 {
		t.Errorf("系統調整攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, SysPool +0.02
func TestBoomFish11(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-210)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(-210)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != 0.02 {
		t.Errorf("系統調整 +0.02 攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}

}

//打到魚, SysPool 0
func TestBoomFish12(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-105)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(-105)
	syspool.SumBullet = proto.Int64(1000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != 0 {
		t.Errorf("系統調整攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}
}

//打到魚, SysPool -0.03
func TestBoomFish13(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(80)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(80)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != -0.03 {
		t.Errorf("系統調整 -0.03 攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ syspool[%v] adjust[%v]",
			glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss(), glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss())
	}

}

//打到魚, SysPool -0.06
func TestBoomFish25(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(120)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(120)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != -0.06 {
		t.Errorf("系統調整 -0.06 攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ syspool[%v] adjust[%v]",
			glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss(), glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss())
	}

}

//打到魚, SysPool -0.1
func TestBoomFish26(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(260)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(260)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != -0.1 {
		t.Errorf("系統調整 -0.1 攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ syspool[%v] adjust[%v]",
			glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss(), glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss())
	}

}

//打到魚, SysPool -0.15
func TestBoomFish27(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(350)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(350)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)
	if glc.syspoolAdjustRate != -0.15 {
		t.Errorf("系統調整 -0.15 攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ syspool[%v] adjust[%v]",
			glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss(), glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss())
	}

}

//打到魚, SysPool -0.5
func TestBoomFish14(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(600)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(600)
	syspool.SumBullet = proto.Int64(5000)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)

	//打到魚
	fishid := make([]int32, 0)
	fishid = append(fishid, 1)
	bydata := new(byfishdata)
	bydata.fishSeq = 2
	bydata.fishtype = 2
	bydata.fishcount = 1
	glc.fishmanager.fishdata[fishid[0]] = bydata
	glc.fish(fishid, ud, false, nil)

	if glc.syspoolAdjustRate != -0.5 {
		t.Errorf("系統調整 -0.5攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	}

	if glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss() != glc.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss() {
		t.Errorf("系統輸錢, 系統水池和調整後水池要一樣~ pool[%v]", glc.TableManager.Room.FishSysPoolData.GetSysSumWinLoss())
	}
}

//打到魚, 金額與捕魚數據統計
func TestBoomFish2(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(1000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)

	//打到魚
	fishid := make([]int32, 0)
	fishid = append(fishid, 1)
	bydata := new(byfishdata)
	bydata.fishSeq = 2
	bydata.fishtype = 2
	bydata.fishcount = 1
	glc.fishmanager.fishdata[fishid[0]] = bydata
	glc.fish(fishid, ud, false, nil)

	if ud.money != 1680 {
		t.Errorf("money[%v]錯誤, ", ud.money)
	}
	if ud.allwin != 1680 {
		t.Errorf("allwin[%v]錯誤", ud.allwin)
	}
	if ud.winmoney != 1680 {
		t.Errorf("winmoney[%v]錯誤", ud.winmoney)
	}

}

//打到黃金放分魚, 金額與捕魚數據統計
func TestBoomFish4(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	glc.TableManager.Room.FishPlayerPoolStatus = true
	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(80000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//黃金魚
	goldFish2 := new(netproto.BY_NewFish)
	goldFish2.FishID = proto.Int32(2)
	goldFish2.FishType = proto.Int32(2)
	goldFish2.FishNumber = proto.Int32(2)
	goldFish2.FishCount = proto.Int32(1)

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)

	//黃金魚
	goldFish := new(netproto.BY_NewFish)
	goldFish.FishID = proto.Int32(2)
	goldFish.FishType = proto.Int32(2)
	goldFish.FishNumber = proto.Int32(2)
	goldFish.FishCount = proto.Int32(1)

	//打到魚
	fishid := make([]int32, 0)
	fishid = append(fishid, 1)
	bydata := new(byfishdata)
	bydata.fishSeq = 2
	bydata.fishtype = 2
	bydata.fishcount = 1
	glc.fishmanager.fishdata[fishid[0]] = bydata
	glc.fish(fishid, ud, false, nil)

	if ud.money != 1680 {
		t.Errorf("money[%v]錯誤, ", ud.money)
	}
	if ud.allwin != 1680 {
		t.Errorf("allwin[%v]錯誤", ud.allwin)
	}
	if ud.winmoney != 1680 {
		t.Errorf("winmoney[%v]錯誤", ud.winmoney)
	}
}

//沒打到黃金放分魚, 金額與捕魚數據統計
func TestBoomFish5(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//黃金魚
	goldFish2 := new(netproto.BY_NewFish)
	goldFish2.FishID = proto.Int32(2)
	goldFish2.FishType = proto.Int32(2)
	goldFish2.FishNumber = proto.Int32(2)
	goldFish2.FishCount = proto.Int32(1)

	//打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)

	//黃金魚
	goldFish := new(netproto.BY_NewFish)
	goldFish.FishID = proto.Int32(2)
	goldFish.FishType = proto.Int32(2)
	goldFish.FishNumber = proto.Int32(2)
	goldFish.FishCount = proto.Int32(1)

	//打到魚
	fishid := make([]int32, 0)
	fishid = append(fishid, 1)
	bydata := new(byfishdata)
	bydata.fishSeq = 2
	bydata.fishtype = 2
	bydata.fishcount = 1
	glc.fishmanager.fishdata[fishid[0]] = bydata
	glc.fish(fishid, ud, false, nil)

}

//沒打到魚, 金額與捕魚數據統計
func TestBoomFish3(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(100)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//沒打到魚
	fishid2 := make([]int32, 0)
	fishid2 = append(fishid2, 1)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 1
	bydata2.fishtype = 2
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid2[0]] = bydata2
	glc.fish(fishid2, ud, false, nil)

	//沒打到魚
	fishid := make([]int32, 0)
	fishid = append(fishid, 1)
	bydata := new(byfishdata)
	bydata.fishSeq = 2
	bydata.fishtype = 2
	bydata.fishcount = 1
	glc.fishmanager.fishdata[fishid[0]] = bydata
	glc.fish(fishid, ud, false, nil)

}

//使用冰冻
func TestUserCastIce(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	boomId := int32(0)
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.BY_GameMessageClassID_FrameS2CID) {
			ust := m.MsgData.(*netproto.BY_FrameS2C)
			if ust.FishWave.GetFishWaveID() > int32(0) {
				glc.gameEnd()
			}

			if len(ust.NewFish) > 0 {
				glc.userdata[0].rate = 10
				//補魚
				btmsg := new(msg.Message)
				btmsg.BClassID = int32(netproto.MessageBClassID_Game)
				btmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameC2SID)
				btdata := new(netproto.BY_FrameC2S)
				btdata.FrameID = ust.FrameID
				//使用冰凍
				btdata.UseIce = proto.Bool(true)
				//用戶發射
				userShoot := new(netproto.BY_ShootC2S)
				pos := new(netproto.BY_Position)
				pos.X = proto.Float32(float32(5))
				pos.Y = proto.Float32(float32(5))
				userShoot.Pos = pos
				userShoot.FishID = ust.NewFish[0].FishID
				btdata.Shoot = nil //userShoot
				//
				////用戶砲彈打到魚
				userBoom := new(netproto.BY_BoomC2S)
				boomId++
				userBoom.BoomID = proto.Int32(0)
				userBoom.BoomRate = proto.Int32(1)
				userBoom.FishID = []int32{}
				userBoom.UserID = proto.Int32(users[1].UserID)
				btdata.Boom = userBoom
				userChangeRate := new(netproto.BY_UserChangeRateC2S)
				userChangeRate.RateAfterChange = proto.Int32(1)
				btdata.UserChangeRate = userChangeRate
				suoDing := new(netproto.BY_SuoDingC2S)
				suoDing.SuoDing = proto.Bool(false)
				btdata.SuoDing = suoDing
				allFish := new(netproto.BY_AllFish)
				btdata.AllFish = allFish
				btdata.UseIce = proto.Bool(true)
				btmsg.MsgData, _ = proto.Marshal(btdata)
				//
				glc.GameMessage(users[1], btmsg)

				//ud := glc.getuserdatabyuserid(users[1].UserID)

				glc.TimerFuncManager.NewTimer(10, func() {

					//if ud.money != int64(-glc.icemanager.iceMoney) {
					//	t.Errorf("冰凍金幣money錯誤! money[%v] iceMoney[%v] ", ud.money, -glc.icemanager.iceMoney)
					//}
					//
					//if ud.allwin != int64(-glc.icemanager.iceMoney) {
					//	t.Errorf("冰凍金幣allwin錯誤! allwin[%v]", ud.allwin)
					//}
					//
					//if ud.winmoney != int64(-glc.icemanager.iceMoney) {
					//	t.Errorf("冰凍金幣winmoney錯誤! winmoney[%v]", ud.winmoney)
					//}

					endChannel <- true
				})

			}

		}

	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("by game end")
	}
}

//户破产消息
func TestIsPoChan(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	boomId := int32(0)
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.BY_GameMessageClassID_FrameS2CID) {
			ust := m.MsgData.(*netproto.BY_FrameS2C)
			if ust.FishWave.GetFishWaveID() > int32(0) {
				if len(ust.PoChan) == 0 {
					t.Errorf("破產驗證失敗, 無人破產")
				}
				glc.gameEnd()
			}

			if len(ust.NewFish) > 0 {
				glc.userdata[0].rate = 10
				//補魚
				btmsg := new(msg.Message)
				btmsg.BClassID = int32(netproto.MessageBClassID_Game)
				btmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameC2SID)
				btdata := new(netproto.BY_FrameC2S)
				btdata.FrameID = ust.FrameID
				//用戶發射
				userShoot := new(netproto.BY_ShootC2S)
				pos := new(netproto.BY_Position)
				pos.X = proto.Float32(float32(5))
				pos.Y = proto.Float32(float32(5))
				userShoot.Pos = pos
				userShoot.FishID = ust.NewFish[0].FishID
				btdata.Shoot = userShoot

				//用戶砲彈打到魚
				userBoom := new(netproto.BY_BoomC2S)
				boomId++
				userBoom.BoomID = proto.Int32(boomId)
				userBoom.BoomRate = proto.Int32(5)
				userBoom.FishID = []int32{ust.NewFish[0].GetFishID()}
				//userBoom.UserID = proto.Int32(users[1].UserID)
				btdata.Boom = userBoom
				userChangeRate := new(netproto.BY_UserChangeRateC2S)
				userChangeRate.RateAfterChange = proto.Int32(5)
				btdata.UserChangeRate = userChangeRate
				suoDing := new(netproto.BY_SuoDingC2S)
				suoDing.SuoDing = proto.Bool(false)
				btdata.SuoDing = suoDing
				allFish := new(netproto.BY_AllFish)
				btdata.AllFish = allFish
				btdata.UseIce = proto.Bool(false)
				btmsg.MsgData, _ = proto.Marshal(btdata)

				glc.GameMessage(users[1], btmsg)

				if len(ust.PoChan) > 0 {
					endChannel <- true
				}

				glc.TimerFuncManager.NewTimer(30, func() {
					if len(ust.PoChan) == 0 {
						t.Errorf("破產驗證失敗, 無人破產")
					}
					endChannel <- true
				})

			}
		}
	}

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("by game end")
	}
}

//取得個人放分的攻擊倍率
func TestRateByUserIdSungFen(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(150)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	rate := glc.getRate(1, ud)

	if rate != float32(2.45) {
		t.Errorf("取得個人放分的攻擊倍率錯誤, rate[%v]", rate)
	}
}

//取得個人殺分的攻擊倍率
func TestRateByUserIdShaFen(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(-50000)
	vipList.Rate = proto.Int32(50)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	rate := glc.getRate(1, ud)

	if rate != float32(0.49) {
		t.Errorf("取得個人殺分的攻擊倍率錯誤, rate[%v]", rate)
	}
}

//取得系統放分魚的攻擊倍率
func TestRateBySysSungFen(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(1)
	vipList.Amount = proto.Int64(6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(180)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	newRate := glc.getRate(1, ud)

	rate, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", newRate), 64)
	if rate != 1.764 {
		t.Errorf("取得系統放分的攻擊倍率錯誤, rate[%v]", rate)
	}
}

//取得系統放分全部魚的攻擊倍率
func TestRateBySysSungFenAllFish(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(0)
	vipList.Amount = proto.Int64(6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(180)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	newRate := glc.getRate(1, ud)

	rate, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", newRate), 64)
	if rate != 1.764 {
		t.Errorf("取得系統放分的攻擊倍率錯誤, rate[%v]", rate)
	}
}

//取得系統殺分魚的攻擊倍率
func TestRateBySysShaFen(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(1)
	vipList.Amount = proto.Int64(-6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(20)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	newRate := glc.getRate(1, ud)

	rate, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", newRate), 64)
	if rate != 0.196 {
		t.Errorf("取得系統殺分的攻擊倍率錯誤, rate[%v]", rate)
	}
}

//取得系統殺分全部魚的攻擊倍率
func TestRateBySysShaFenAllFish(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(0)
	vipList.Amount = proto.Int64(-6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(20)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	newRate := glc.getRate(1, ud)

	rate, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", newRate), 64)
	if rate != 0.196 {
		t.Errorf("取得系統殺分的攻擊倍率錯誤, rate[%v]", rate)
	}
}

//刷新抓到魚個人放分數據
func TestUpdateVIPReword(t *testing.T) {
	//testBBBB
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(150)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(true, fish, ud, 300)

	if glc.TableManager.Room.VIPUserRewardList.FishVIPUserList[0].GetFinishAmount() != 200 {
		t.Errorf("刷新抓到魚個人放分錯誤, FinishAmount[%v]", glc.TableManager.Room.VIPUserRewardList.FishVIPUserList[0].GetFinishAmount())
	}
}

//刷新沒抓到魚個人放分數據
func TestUpdateVIPReword2(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(150)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.boomNumber = 5
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(false, fish, ud, 0)

	if glc.TableManager.Room.VIPUserRewardList.FishVIPUserList[0].GetFinishAmount() != -500 {
		t.Errorf("刷新抓到魚個人放分錯誤, FinishAmount[%v]", glc.TableManager.Room.VIPUserRewardList.FishVIPUserList[0].GetFinishAmount())
	}
}

//刷新沒抓到魚個人殺分數據
func TestUpdateVIPReword3(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(-50000)
	vipList.Rate = proto.Int32(30)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(1)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.boomNumber = 5
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(false, fish, ud, 0)

	if glc.TableManager.Room.VIPUserRewardList.FishVIPUserList[0].GetFinishAmount() != -500 {
		t.Errorf("刷新抓到魚個人殺分錯誤, FinishAmount[%v]", glc.TableManager.Room.VIPUserRewardList.FishVIPUserList[0].GetFinishAmount())
	}
}

//刷新抓到魚系統放分數據
func TestUpdateSysReword(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(1)
	vipList.Amount = proto.Int64(6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(180)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.boomNumber = 2
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(true, fish, ud, 300)

	if glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount() != 100 {
		t.Errorf("刷新抓到魚系統放分錯誤, FinishAmount[%v]", glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount())
	}
}

//刷新抓到所有魚系統放分數據
func TestUpdateSysRewordAllFish(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(0)
	vipList.Amount = proto.Int64(6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(180)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.boomNumber = 1
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(true, fish, ud, 300)

	if glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount() != 200 {
		t.Errorf("刷新抓到魚系統放分錯誤, FinishAmount[%v]", glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount())
	}
}

//刷新沒抓到魚系統放分數據
func TestUpdateSysReword2(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(1)
	vipList.Amount = proto.Int64(6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(180)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.boomNumber = 2
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(false, fish, ud, 0)

	if glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount() != -200 {
		t.Errorf("刷新沒抓到魚系統放分錯誤, FinishAmount[%v]", glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount())
	}
}

//刷新沒抓到魚系統殺分數據
func TestUpdateSysReword3(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(1)
	vipList.Amount = proto.Int64(-6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(80)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.boomNumber = 2
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(false, fish, ud, 0)

	if glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount() != -200 {
		t.Errorf("刷新沒抓到魚系統殺分錯誤, FinishAmount[%v]", glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount())
	}
}

//刷新沒抓到All魚系統放分數據
func TestUpdateSysRewordAllFish1(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(0)
	vipList.Amount = proto.Int64(6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(180)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.boomNumber = 2
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(false, fish, ud, 0)

	if glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount() != -200 {
		t.Errorf("刷新沒抓到All魚系統殺分錯誤, FinishAmount[%v]", glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount())
	}
}

//刷新沒抓到All魚系統殺分數據
func TestUpdateSysRewordAllFish2(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	vip := new(netproto.SuperToolFishGetSystemRewardDetail)
	vipList := new(netproto.FishSystemRewardDataList)
	vipList.FishID = proto.Int32(0)
	vipList.Amount = proto.Int64(-6000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Rate = proto.Int32(80)
	vipList.Status = proto.Int32(1)
	vip.FishRewardList = append(vip.FishRewardList, vipList)
	glc.TableManager.Room.SysRewardList = vip

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.boomNumber = 2
	fish := new(byfishdata)
	fish.fishSeq = 10
	fish.fishtype = 1

	glc.updateSysReword(false, fish, ud, 0)

	if glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount() != -200 {
		t.Errorf("刷新沒抓到All魚系統殺分錯誤, FinishAmount[%v]", glc.TableManager.Room.SysRewardList.FishRewardList[0].GetFinishAmount())
	}
}

//打到炸弹螃蟹
func TestBoomFish21(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	ud.user.SumCharge = 10

	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	glc.newTeshuFish()

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-5)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(-5)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(30000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到炸彈蟹
	fishid1 := make([]int32, 0)
	fishid1 = append(fishid1, 1)
	bydata1 := new(byfishdata)
	bydata1.fishSeq = 31
	bydata1.fishtype = 31
	bydata1.fishcount = 1
	glc.fishmanager.fishdata[fishid1[0]] = bydata1
	glc.fish(fishid1, ud, false, nil)
	if ud.diedbosstype != BoomFish {
		t.Errorf("打死了炸彈蟹但diedbosstype不是3.  錯誤資料:diedbosstype[%v]", ud.diedbosstype)
	}
	if ud.money != 0 {
		t.Errorf("打死了炸彈蟹亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:money[%v]", ud.money)
	}
	if ud.allwin != 0 {
		t.Errorf("打死了炸彈蟹亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:allwin[%v]", ud.allwin)
	}
	if ud.winmoney != 0 {
		t.Errorf("打死了炸彈蟹亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:winmoney[%v]", ud.winmoney)
	}
	if ud.teshufishtype != BoomFish {
		t.Errorf("打死了炸彈蟹亂倍魚, 錯誤資料:teshufishtype[%v]", ud.teshufishtype)
	}

	//打到魚
	fishid3 := make([]int32, 0)
	fishid3 = append(fishid3, 1)
	bydata3 := new(byfishdata)
	bydata3.fishSeq = 1
	bydata3.fishtype = 1
	bydata3.fishodd = 1200
	bydata3.fishcount = 1
	glc.fishmanager.fishdata[fishid3[0]] = bydata3
	//打到魚
	fishid3 = append(fishid3, 2)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 2
	bydata2.fishtype = 2
	bydata2.fishodd = 2500
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid3[1]] = bydata2

	ust := new(netproto.BY_FrameC2S)
	ust.AllFish = new(netproto.BY_AllFish)
	ust.AllFish.FishID = fishid3
	glc.boomFishSetAllOdd(ud, ust)
	glc.fish(fishid3, ud, true, nil)

	//if glc.syspoolAdjustRate != -0.1 {
	//	t.Errorf("系統調整攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	//}
	if ud.money != 25900 {
		t.Errorf("打死了炸彈蟹亂倍魚,金額要從其他死掉魚計算. 錯誤資料:money[%v]", ud.money)
	}
	if ud.allwin != 25900 {
		t.Errorf("打死了炸彈蟹亂倍魚,金額要從其他死掉魚計算. 錯誤資料:allwin[%v]", ud.allwin)
	}
	if ud.winmoney != 25900 {
		t.Errorf("打死了炸彈蟹亂倍魚,金額要從其他死掉魚計算. 錯誤資料:winmoney[%v]", ud.winmoney)
	}
	if ud.teshufishtype != BoomFish {
		t.Errorf("打死了炸彈蟹亂倍魚, 錯誤資料:teshufishtype[%v]", ud.teshufishtype)
	}

}

//打到閃電鰻魚
func TestBoomFish22(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	ud.user.SumCharge = 10

	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	glc.newTeshuFish()

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-5)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(-5)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(30000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到閃電鰻魚
	fishid1 := make([]int32, 0)
	fishid1 = append(fishid1, 1)
	bydata1 := new(byfishdata)
	bydata1.fishSeq = 32
	bydata1.fishtype = 32
	bydata1.fishcount = 1
	glc.fishmanager.fishdata[fishid1[0]] = bydata1
	glc.fish(fishid1, ud, false, nil)
	if ud.diedbosstype != LeiLong {
		t.Errorf("打死了閃電鰻魚但diedbosstype不是4.  錯誤資料:diedbosstype[%v]", ud.diedbosstype)
	}
	if ud.money != 0 {
		t.Errorf("打死了閃電鰻魚亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:money[%v]", ud.money)
	}
	if ud.allwin != 0 {
		t.Errorf("打死了閃電鰻魚亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:allwin[%v]", ud.allwin)
	}
	if ud.winmoney != 0 {
		t.Errorf("打死了閃電鰻魚亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:winmoney[%v]", ud.winmoney)
	}
	if ud.teshufishtype != LeiLong {
		t.Errorf("打死了閃電鰻魚亂倍魚, 錯誤資料:teshufishtype[%v]", ud.teshufishtype)
	}

	//打到魚
	fishid3 := make([]int32, 0)
	fishid3 = append(fishid3, 1)
	bydata3 := new(byfishdata)
	bydata3.fishSeq = 1
	bydata3.fishtype = 1
	bydata3.fishodd = 1200
	bydata3.fishcount = 1
	glc.fishmanager.fishdata[fishid3[0]] = bydata3
	//打到魚
	fishid3 = append(fishid3, 2)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 2
	bydata2.fishtype = 2
	bydata2.fishodd = 2500
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid3[1]] = bydata2

	ust := new(netproto.BY_FrameC2S)
	ust.AllFish = new(netproto.BY_AllFish)
	ust.AllFish.FishID = fishid3
	glc.boomFishSetAllOdd(ud, ust)
	glc.fish(fishid3, ud, true, nil)

	//if glc.syspoolAdjustRate != -0.1 {
	//	t.Errorf("系統調整攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	//}
	if ud.money != 25900 {
		t.Errorf("打死了閃電鰻魚亂倍魚,金額要從其他死掉魚計算. 錯誤資料:money[%v]", ud.money)
	}
	if ud.allwin != 25900 {
		t.Errorf("打死了閃電鰻魚亂倍魚,金額要從其他死掉魚計算. 錯誤資料:allwin[%v]", ud.allwin)
	}
	if ud.winmoney != 25900 {
		t.Errorf("打死了閃電鰻魚亂倍魚,金額要從其他死掉魚計算. 錯誤資料:winmoney[%v]", ud.winmoney)
	}
	if ud.teshufishtype != LeiLong {
		t.Errorf("打死了閃電鰻魚亂倍魚, 錯誤資料:teshufishtype[%v]", ud.teshufishtype)
	}

}

//打到磚頭龍蝦
func TestBoomFish23(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	ud.user.SumCharge = 10

	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	glc.newTeshuFish()

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(5)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(5)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(19000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到磚頭龍蝦
	fishid1 := make([]int32, 0)
	fishid1 = append(fishid1, 1)
	bydata1 := new(byfishdata)
	bydata1.fishSeq = 36
	bydata1.fishtype = 36
	bydata1.fishcount = 1
	glc.fishmanager.fishdata[fishid1[0]] = bydata1
	glc.fish(fishid1, ud, false, nil)
	if ud.diedbosstype != Zuantoulongxia {
		t.Errorf("打死了磚頭龍蝦但diedbosstype不是4.  錯誤資料:diedbosstype[%v]", ud.diedbosstype)
	}
	if ud.money != 0 {
		t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:money[%v]", ud.money)
	}
	if ud.allwin != 0 {
		t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:allwin[%v]", ud.allwin)
	}
	if ud.winmoney != 0 {
		t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:winmoney[%v]", ud.winmoney)
	}
	if ud.teshufishtype != Zuantoulongxia {
		t.Errorf("打死了磚頭龍蝦亂倍魚, 錯誤資料:teshufishtype[%v]", ud.teshufishtype)
	}

	//打到魚
	fishid3 := make([]int32, 0)
	fishid3 = append(fishid3, 1)
	bydata3 := new(byfishdata)
	bydata3.fishSeq = 1
	bydata3.fishtype = 1
	bydata3.fishodd = 1200
	bydata3.fishcount = 1
	glc.fishmanager.fishdata[fishid3[0]] = bydata3
	//打到魚
	fishid3 = append(fishid3, 2)
	bydata2 := new(byfishdata)
	bydata2.fishSeq = 2
	bydata2.fishtype = 2
	bydata2.fishodd = 2500
	bydata2.fishcount = 1
	glc.fishmanager.fishdata[fishid3[1]] = bydata2

	ust := new(netproto.BY_FrameC2S)
	ust.AllFish = new(netproto.BY_AllFish)
	ust.AllFish.FishID = fishid3
	glc.boomFishSetAllOdd(ud, ust)
	glc.fish(fishid3, ud, true, nil)

	//if glc.syspoolAdjustRate != -0.1 {
	//	t.Errorf("系統調整攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
	//}
	if ud.money != 25900 {
		t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算. 錯誤資料:money[%v]", ud.money)
	}
	if ud.allwin != 25900 {
		t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算. 錯誤資料:allwin[%v]", ud.allwin)
	}
	if ud.winmoney != 25900 {
		t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算. 錯誤資料:winmoney[%v]", ud.winmoney)
	}
	if ud.teshufishtype != Zuantoulongxia {
		t.Errorf("打死了磚頭龍蝦亂倍魚, 錯誤資料:teshufishtype[%v]", ud.teshufishtype)
	}

}

//打到核彈
//func TestBoomFish25(t *testing.T) {
//	endChannel := make(chan bool)
//	confInit()
//	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
//		return NewByLogic(t, tfManager, cfgs)
//	}
//	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
//		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
//	}
//
//	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
//	robotManager := tgrs.RobotService.GetRobotManager()
//
//	for i := 0; i < 3; i++ {
//		robot := new(robotcomm.Robot)
//		robot.UserInfo = robots[i]
//		robot.RebotLevel = 6
//		robot.UserInfo.UserType = 2
//		robotManager.AddRobot(robot)
//	}
//
//	users := fixtures.GetGameUser()
//	users[1].Score = 1
//	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
//	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
//
//	//安排用戶入座
//	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
//	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
//
//	//glc.UserReady(robots[0])
//	glc.UserReady(users[1])
//
//	ud := glc.getuserdatabyuserid(users[1].UserID)
//	ud.currodd = 7
//	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
//	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
//		UserID: proto.Int32(users[1].UserID) ,
//		Bullet:proto.Int32(600),
//		Score: proto.Int32(100),
//		Status:proto.Int32(1),
//	}
//
//	glc.newfish()
//	glc.newBossFish()
//	glc.newTeshuFish()
//
//	syspool := new(netproto.FishSysPoolData)
//	syspool.FlagID = proto.Int32(1)
//	syspool.SysSumWinLoss = proto.Int64(-5)
//	syspool.SumBullet = proto.Int64(150)
//	syspool.SumBetMoney = proto.Int64(900)
//	glc.TableManager.Room.FishSysPoolData = syspool
//
//	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
//	vipList := new(netproto.FishVIPUserRewardDataList)
//	vipList.UserID = proto.Int32(users[1].UserID)
//	vipList.AssignAmount = proto.Int64(50000)
//	vipList.Rate = proto.Int32(19000)
//	vipList.FinishAmount = proto.Int64(0)
//	vipList.Status = proto.Int32(0)
//	vipList.SumAssignAmount = proto.Int64(0)
//	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
//	glc.TableManager.Room.VIPUserRewardList = vip
//
//	//打到炸彈蟹
//	fishid1 := make([]int32, 0)
//	fishid1 = append(fishid1, 1)
//	bydata1 := new(byfishdata)
//	bydata1.fishSeq = 41
//	bydata1.fishtype = 41
//	bydata1.fishcount = 1
//	glc.fishmanager.fishdata[fishid1[0]] = bydata1
//
//	ust := new(netproto.BY_FrameC2S)
//	ust.AllFish = new(netproto.BY_AllFish)
//	ust.AllFish.FishID = []int32{2, 3}
//	glc.fish(fishid1, ud, false, ust.AllFish.FishID )
//	//if ud.diedbosstype != BoomFish {
//	//	t.Errorf("打死了核彈魚但diedbosstype不是3.  錯誤資料:diedbosstype[%v]", ud.diedbosstype)
//	//}
//
//	if ud.money == 0 {
//		t.Errorf("打死了核彈魚亂倍魚,金額要從其他死掉魚計算,核彈不會有錢. 錯誤資料:money[%v]", ud.money)
//	}
//	if ud.allwin == 0 {
//		t.Errorf("打死了核彈魚亂倍魚,金額要從其他死掉魚計算,核彈不會有錢. 錯誤資料:allwin[%v]", ud.allwin)
//	}
//	if ud.winmoney == 0 {
//		t.Errorf("打死了核彈魚亂倍魚,金額要從其他死掉魚計算,核彈不會有錢. 錯誤資料:winmoney[%v]", ud.winmoney)
//	}
//	//if ud.teshufishtype != BoomFish {
//	//	t.Errorf("打死了核彈魚亂倍魚, 錯誤資料:teshufishtype[%v]", ud.teshufishtype)
//	//}
//
//	//打到魚
//	//fishid3 := make([]int32, 0)
//	//fishid3 = append(fishid3, 1)
//	//bydata3 := new(byfishdata)
//	//bydata3.fishSeq = 1
//	//bydata3.fishtype = 1
//	//bydata3.fishodd = 1200
//	//bydata3.fishcount = 1
//	//glc.fishmanager.fishdata[fishid3[0]] = bydata3
//	////打到魚
//	//fishid3 = append(fishid3, 2)
//	//bydata2 := new(byfishdata)
//	//bydata2.fishSeq = 2
//	//bydata2.fishtype = 2
//	//bydata2.fishodd = 2500
//	//bydata2.fishcount = 1
//	//glc.fishmanager.fishdata[fishid3[1]] = bydata2
//	//
//	//
//	//ust.AllFish = new(netproto.BY_AllFish)
//	//ust.AllFish.FishID = fishid3
//	//glc.boomFishSetAllOdd(ud, ust)
//	//glc.fish(fishid3, ud, false, ust.AllFish.FishID )
//	////glc.fish(fishid3, ud, true, nil)
//	//
//	////if glc.syspoolAdjustRate != -0.1 {
//	////	t.Errorf("系統調整攻擊率錯誤 , err : %v", glc.syspoolAdjustRate)
//	////}
//	//if ud.money != 25900 {
//	//	t.Errorf("打死了核彈魚亂倍魚,金額要從其他死掉魚計算. 錯誤資料:money[%v]", ud.money)
//	//}
//	//if ud.allwin != 25900 {
//	//	t.Errorf("打死了核彈魚亂倍魚,金額要從其他死掉魚計算. 錯誤資料:allwin[%v]", ud.allwin)
//	//}
//	//if ud.winmoney != 25900 {
//	//	t.Errorf("打死了核彈魚亂倍魚,金額要從其他死掉魚計算. 錯誤資料:winmoney[%v]", ud.winmoney)
//	//}
//
//}

//打到鯤
func TestBoomFish24(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7

	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	glc.newTeshuFish()

	syspool := new(netproto.FishSysPoolData)
	syspool.FlagID = proto.Int32(1)
	syspool.SysSumWinLoss = proto.Int64(-5)
	syspool.AfterAdjustSyspoolWinLoss = proto.Int64(-5)
	syspool.SumBullet = proto.Int64(150)
	syspool.SumBetMoney = proto.Int64(900)
	glc.TableManager.Room.FishSysPoolData = syspool

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(200000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	//打到炸彈蟹
	fishid1 := make([]int32, 0)
	fishid1 = append(fishid1, 42)
	bydata1 := new(byfishdata)
	bydata1.fishSeq = 42
	bydata1.fishtype = 42
	bydata1.fishcount = 1
	glc.fishmanager.fishdata[fishid1[0]] = bydata1
	glc.fish(fishid1, ud, false, nil)
	if ud.diedbosstype != Dujiaojing {
		t.Logf("打中鯤但diedbosstype不是3.  錯誤資料:diedbosstype[%v]", ud.diedbosstype)
	}
	if ud.money < 60000 {
		t.Errorf("打中鯤亂倍魚,金額不能低於600000. 錯誤資料:money[%v]", ud.money)
	}
	if ud.allwin < 60000 {
		t.Errorf("打中鯤亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:allwin[%v]", ud.allwin)
	}
	if ud.winmoney < 60000 {
		t.Errorf("打中鯤亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:winmoney[%v]", ud.winmoney)
	}
}

//閃電鰻魚 or 磚頭砲打中魚死亡清單
func TestLeiLongOrZtxAllFishDied(t *testing.T) {
	endChannel := make(chan bool)
	confInit()
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewByLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBYRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()
	users[1].Score = 1
	//robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	//glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//glc.UserReady(robots[0])
	glc.UserReady(users[1])

	ud := glc.getuserdatabyuserid(users[1].UserID)
	ud.currodd = 7
	ud.user.SumCharge = 10

	glc.TableManager.Room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
	glc.TableManager.Room.FisnNewPlayerDetail[users[1].UserID] = &netproto.FishBeginnerOffer{
		UserID: proto.Int32(users[1].UserID),
		Bullet: proto.Int32(600),
		Score:  proto.Int32(100),
		Status: proto.Int32(1),
	}

	vip := new(netproto.SuperToolFishGetVIPUserRewardDetail)
	vipList := new(netproto.FishVIPUserRewardDataList)
	vipList.UserID = proto.Int32(users[1].UserID)
	vipList.AssignAmount = proto.Int64(50000)
	vipList.Rate = proto.Int32(19000)
	vipList.FinishAmount = proto.Int64(0)
	vipList.Status = proto.Int32(0)
	vipList.SumAssignAmount = proto.Int64(0)
	vip.FishVIPUserList = append(vip.FishVIPUserList, vipList)
	glc.TableManager.Room.VIPUserRewardList = vip

	for i := 0; i < 1; i++ {
		//打到磚頭龍蝦
		fishid1 := make([]int32, 0)
		fishid1 = append(fishid1, 1)
		bydata1 := new(byfishdata)
		bydata1.fishSeq = 36
		bydata1.fishtype = 36
		bydata1.fishcount = 1
		glc.fishmanager.fishdata[fishid1[0]] = bydata1
		glc.fish(fishid1, ud, false, nil)
		if ud.diedbosstype != Zuantoulongxia {
			t.Errorf("打死了磚頭龍蝦但diedbosstype不是4.  錯誤資料:diedbosstype[%v]", ud.diedbosstype)
		}
		if ud.money != 0 {
			t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:money[%v]", ud.money)
		}
		if ud.allwin != 0 {
			t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:allwin[%v]", ud.allwin)
		}
		if ud.winmoney != 0 {
			t.Errorf("打死了磚頭龍蝦亂倍魚,金額要從其他死掉魚計算,炸彈蟹不會有錢. 錯誤資料:winmoney[%v]", ud.winmoney)
		}
		if ud.teshufishtype != Zuantoulongxia {
			t.Errorf("打死了磚頭龍蝦亂倍魚, 錯誤資料:teshufishtype[%v]", ud.teshufishtype)
		}

		//打到魚
		fishid3 := make([]int32, 0)
		fishid3 = append(fishid3, 1)
		bydata3 := new(byfishdata)
		bydata3.fishSeq = 1
		bydata3.fishtype = 1
		bydata3.fishodd = 1200
		bydata3.fishcount = 1
		glc.fishmanager.fishdata[fishid3[0]] = bydata3
		//打到魚
		fishid3 = append(fishid3, 2)
		bydata2 := new(byfishdata)
		bydata2.fishSeq = 2
		bydata2.fishtype = 2
		bydata2.fishodd = 2500
		bydata2.fishcount = 1
		glc.fishmanager.fishdata[fishid3[1]] = bydata2
		ust := new(netproto.BY_FrameC2S)
		ust.AllFish = new(netproto.BY_AllFish)
		ust.AllFish.FishID = fishid3
		ud.candiedrate = 0
		glc.leiLongOrZtxAllFishDied(ust.AllFish.FishID, ud)
		checkDone := true
		if ud.candiedrate != 0.2 {
			checkDone = false
			t.Errorf("閃電鰻魚 or 磚頭砲 攻擊率遞增錯誤. 錯誤資料candiedrate[%v]", ud.candiedrate)
		}
		if ud.teshufishfishodd != 15000 {
			checkDone = false
			t.Errorf("閃電鰻魚 or 磚頭砲 特殊魚倍率錯誤. 錯誤資料teshufishfishodd[%v]", ud.teshufishfishodd)
		}

		fishid2 := make([]int32, 0)
		fishid2 = append(fishid2, 1)
		ud.candiedrate = 100
		glc.leiLongOrZtxAllFishDied(fishid2, ud)
		if ud.candiedrate != 0.4 {
			checkDone = false
			t.Errorf("閃電鰻魚 or 磚頭砲 攻擊率遞增錯誤. 錯誤資料candiedrate[%v]", ud.candiedrate)
		}
		if ud.teshufishfishodd != 14890 {
			checkDone = false
			t.Errorf("閃電鰻魚 or 磚頭砲 特殊魚倍率錯誤. 錯誤資料teshufishfishodd[%v]", ud.teshufishfishodd)
		}

		if checkDone {
			break
		}
	}
}

//流程整合測試共用初始化
func IntegrationTestInit(endChannel chan bool, t *testing.T) (*fixtures.TestGameRoomServer, *Bylogic, chan bool, []*gamelogic.GameUser) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*Bylogic)
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
		robotManager.AddRobot(robot)
	}

	glc.TableManager.Room.Seq = proto.Int64(0)

	return tgrs, glc, endChannel, robots
}
