package dfdc

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
)

// MyRobotLogic abc
type MyRobotLogic struct {
	*robotlogic.LogicBase
}

// NewMyRobotLogic abc
func NewMyRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *MyRobotLogic {
	rblgc := new(MyRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	return rblgc
}

// OnGameMessage 游戏消息
func (lgc *MyRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	// glc := lgc.getGameLogic()
	// 遊戲為下注狀態
	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCStateID) {
		ust := m.MsgData.(*netproto.DFDC_GameStatus)
		//mlog.System("=========== ust state %d", ust.GetState())
		if ust.GetState() == 0 {
			//mlog.System("=========== test %d", robot.UserInfo.UserID)
			lgc.SendGameMessage(robot.UserInfo, lgc.getBetMsg(10, 88))
			return
		}
	}
}

func (lgc *MyRobotLogic) getBetMsg(money int64, goldenSymbol int32) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCGameStartID)
	betdata := new(netproto.DFDC_StartBet)
	betdata.Money = proto.Int64(money)
	betdata.GoldenSymbol = proto.Int32(goldenSymbol)
	btmsg.MsgData, _ = proto.Marshal(betdata)
	return btmsg
}

func (lgc *MyRobotLogic) getGameLogic() *Logic {
	return lgc.GameLogic.(*Logic)
}
