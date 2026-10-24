package hbsl

import (
	"testing"

	"example.invalid/stretchr/testify/assert"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

type TestRoundLimit struct {
	CurrentRound int
	RoundLimit   int
	EndSignal    chan bool
}

func init() {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHbslLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHbslRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}
}

// 測試遊戲流程, 十局
func TestGameRound1(t *testing.T) {
	conf.TestMode = true
	conf.TimerFastMode = 100
	roundLimit := &TestRoundLimit{
		1,
		10,
		make(chan bool),
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HbslLogic)
	//glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	//}
	glc.RegisterAfterGameStartFunc(func() {
		mlog.System("game start ...")
	})

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterAfterGameEndFunc(func() {
		mlog.System("game end ...")
		roundLimit.EndSignal <- true
	})
	robotManager := tgrs.RobotService.GetRobotManager()
	_, robots := fixtures.GetGameUsers()

	for i := 0; i < 2; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	glc.UserEnter(robots[0], 0, 0)
	glc.UserEnter(robots[1], 0, 1)

End:
	for {
		select {
		case <-roundLimit.EndSignal:
			if roundLimit.CurrentRound >= roundLimit.RoundLimit {
				glc.GameStartFunc = func() {

				}
				break End
			}
			roundLimit.CurrentRound++
		}

	}

}

// 測試分紅包, 100次
func TestDivideHB(t *testing.T) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HbslLogic)

	glc.pNowHB = new(HBRecord)
	glc.pNowHB.pHB = new(HB)
	glc.pNowHB.ary_i64SHB = make([]int64, 0)
	glc.pNowHB.pHB.i64Chip = 50
	actual, _ := glc.GetUserWinAndTax(glc.pNowHB.pHB.i64Chip)
	for i := 0; i < 100; i++ {
		//iRnd := rand.Intn(len(glc.ary_i64Chip))
		//glc.pNowHB.pHB.i64Chip = glc.ary_i64Chip[iRnd]

		glc.pNowHB.ary_i64SHB = glc.DivideHB()
		sum := int64(0)
		for _, v := range glc.pNowHB.ary_i64SHB {
			sum += v
		}
		assert.Equal(t, sum, actual, "DivideHB not equal")
		assert.Equal(t, int32(len(glc.pNowHB.ary_i64SHB)), glc.i32HBNum, "DivideHB len not equal")
		t.Log(glc.pNowHB.ary_i64SHB)
	}
}
