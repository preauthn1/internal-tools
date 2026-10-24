package ddz

import (
	"fmt"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func TestGameRound_My1(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 20
	//conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*DDZGameLogic)
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		user := glc.GetUserByUserID(userid)
		if user.IsRobot() {
			tgrs.RobotService.SendGameMessage(user, m)
		} else {
			fmt.Println(m.SClassID)
			//分配位子後, 先準備
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_UserEnterID) {
				tgrs.RoomActionController.UserReady(user)
			}

			//游戏结束准备时间消息
			if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_GameEndID) {
				user.UserState.OnUserReady(user, tgrs.Room)
				user.UserState.OnUserQueue(user, tgrs.Room)
			}

		}

	}

	//eventChannel := tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.LogicBase)
	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})
	players := fixtures.GetPlayers()
	tgrs.PlayerManager.SetPlayer(players[0])
	gur := players[0].UserData.(*gamelogic.GameUser)
	gur.UserState.OnUserQueue(gur, tgrs.Room)

	robots := fixtures.GetRobots()

	tgrs.RobotService.GetRobotManager().AddRobot(robots[0])
	tgrs.RobotService.GetRobotManager().AddRobot(robots[1])

	//
	//tgrs.RoomActionController.UserQueue(users[0])

	//users[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	//for {
	//
	//	select {
	//	case _ = <-endChannel:
	//		mlog.NewDefaultLogger().Error("zjh game end")
	//		gur.UserState.OnUserQueue(gur, tgrs.Room)
	//	}
	//	fmt.Println("一局結束")
	//}

}
