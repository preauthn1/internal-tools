package hljbp

import (
	"fmt"
	"testing"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
)

type TestRoundLimit struct {
	CurrentRound int
	RoundLimit   int
	EndSignal    chan bool
}

// TestGameRound1 測試遊戲流程, 十局
func TestGameRound1(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	conf.TestMode = true
	conf.TimerFastMode = 100
	roundLimit := &TestRoundLimit{
		1,
		10,
		make(chan bool),
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*Logic)

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

func TestGame_OneRound_TwoUser_OneEnterSlower(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "阿呆", 10000000)
	u2, p2 := fixtures.NewGameUser(1001, "阿瓜", 10000000)

	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)
	go func() {
		done <- <-c1 && <-c2
	}()

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}

		switch userID {
		case u1.UserID:
			switch netproto.GameMessageClassID(m.SClassID) {
			case netproto.GameMessageClassID_GameEndID:
				go func() {
					c1 <- true
				}()
			}
		case u2.UserID:
			switch netproto.GameMessageClassID(m.SClassID) {
			case netproto.GameMessageClassID_GameEndID:
				go func() {
					c2 <- true
				}()
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	go func() { // User2 Enter
		u2.UserState.OnUserQueue(u2, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p2)
		game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		fmt.Println("Done")

		var expUserCnt int = 2

		if len(game.GameLog.Context().Users) != expUserCnt {
			t.Errorf("PerGameLog Context Users Length want %d, but %d", expUserCnt, len(game.GameLog.Context().Users))
		}

	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}
