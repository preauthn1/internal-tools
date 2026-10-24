package robot

import (
	"testing"
	"time"

	"example.invalid/ServerTeam/Servers/internal/game/inrm"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

func TestRobotLogic(t *testing.T) {

	cfg.Configs[gamelogic.ConfigFileNameGameServer].Sections["robot"]["gameactionwaitspan"] = "0,0"

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*inrm.GameServer)

	rbt1, p1 := fixtures.NewGameRobot(1000, "阿呆", 100000)
	rbt2, p2 := fixtures.NewGameRobot(2000, "阿瓜", 100000)
	rbt3, p3 := fixtures.NewGameRobot(3000, "阿飛", 100000)
	rbt4, p4 := fixtures.NewGameRobot(4000, "阿機", 100000)
	rbt5, p5 := fixtures.NewGameRobot(5000, "阿掰", 100000)

	rbtmgr := grs.RobotService.GetRobotManager()
	rbtmgr.AddRobot(rbt1)
	rbtmgr.AddRobot(rbt2)
	rbtmgr.AddRobot(rbt3)
	rbtmgr.AddRobot(rbt4)
	rbtmgr.AddRobot(rbt5)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)
	grs.PlayerManager.SetPlayer(p4)
	grs.PlayerManager.SetPlayer(p5)

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	game.UserEnter(rbt1.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(rbt2.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(rbt3.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(rbt4.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(rbt5.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	<-start

	select {
	case <-end:
	case <-time.After(5 * time.Second):
		t.Errorf("timeout")
	}
}

func TestRobotLogic_沒真人離開(t *testing.T) {

	if testing.Short() {
		t.Skip("Skiped")
		return
	}

	conf.CLogLevel = mlog.InfoLevel
	conf.TimerFastMode = 1

	cfg.Configs[gamelogic.ConfigFileNameGameServer].Sections["robot"]["gameactionwaitspan"] = "0,0"

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*inrm.GameServer)

	rbt1, p1 := fixtures.NewGameRobot(1000, "阿呆", 100000)
	rbt2, p2 := fixtures.NewGameRobot(2000, "阿瓜", 100000)
	rbt3, p3 := fixtures.NewGameRobot(3000, "阿飛", 100000)
	rbt4, p4 := fixtures.NewGameRobot(4000, "阿機", 100000)
	u5, p5 := fixtures.NewGameUser(5000, "阿掰", 100000)

	rbtmgr := grs.RobotService.GetRobotManager()
	rbtmgr.AddRobot(rbt1)
	rbtmgr.AddRobot(rbt2)
	rbtmgr.AddRobot(rbt3)
	rbtmgr.AddRobot(rbt4)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)
	grs.PlayerManager.SetPlayer(p4)
	grs.PlayerManager.SetPlayer(p5)

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	game.UserEnter(rbt1.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(rbt2.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(rbt3.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(rbt4.UserInfo, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u5, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	game.UserLeave(u5)

	t.Log(game.State, game.IsGaming())
}
