package robot

import (
	"context"
	"fmt"
	"time"

	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// RealTimeRobot Slot 分析機器人
type RealTimeRobot struct {
	*robotlogic.LogicBase

	ctx    context.Context
	cancel context.CancelFunc

	GameName string

	BeginTime time.Time
	LastTime  time.Time

	TotalTimes int
	LeftTimes  int

	BetMoneys   []int64
	BetRates    []int64
	MainGameIDs []int32

	BetMoney  int64
	BetRate   int64
	SelectIdx int

	GameState   slot.GameState
	BetStatus   int32
	LastSendOpt *netproto.Slot_ChooseGameOption

	PickGame *netproto.TGW_MSG_BonusGame

	JackpotPool           map[slot.JackpotLevel]*netproto.Slot_JackpotInfo
	NetprotoFreeGameInfos []*netproto.Slot_FreeGameInfo

	GoldBall []*netproto.Slot_GoldBallInfo

	MainGameJPTimes  int
	MainGameJPWin    int64
	FreeGameJPTimes  int
	FreeGameJPWin    int64
	LightningJPTimes int
	LightningJPWin   int64

	RandomBet       bool
	RandomSelection bool

	Features uint

	Debug bool
}

// NewRealTimeRobot Create RealTimeRobot Robot
func NewRealTimeRobot(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) (*RealTimeRobot, context.Context, context.CancelFunc) {
	rbt := new(RealTimeRobot)
	rbt.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)

	rbt.ctx, rbt.cancel = context.WithCancel(context.Background())

	return rbt, rbt.ctx, rbt.cancel
}

// OnGameMessage 收到遊戲訊息
func (rbt *RealTimeRobot) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

	// 初始設定
	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_SettingID) {
		rbt.BeginTime = time.Now()
		rbt.LastTime = time.Now()
		rbt.onSlotSetting(robot, m)
		return
	}

	// 遊戲狀態
	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_StatusID) {
		rbt.onGameStatus(robot, m)
		return
	}

	// Jackpot pool
	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_JackpotPoolID) {
		rbt.onJapotPool(robot, m)
		return
	}

	// 下注狀態
	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_BetStatusID) {
		rbt.onBetStatus(robot, m)
	}

	// 遊戲結果
	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_GameResultID) {
		rbt.onGameResult(robot, m)
		return
	}

	// PickGame遊戲
	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_BonusGameID) {
		rbt.onPickGame(robot, m)
		return
	}

	// 送確認玩家選擇封包
	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_ActionPickRetID) {
		rbt.actionPickGame(robot)
		return
	}
}

func (rbt *RealTimeRobot) onSlotSetting(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.TGW_MSG_Setting)

	rbt.BetMoneys = make([]int64, 0)
	for _, row := range data.GetBetCount() {
		rbt.BetMoneys = append(rbt.BetMoneys, int64(row))
	}

	rbt.BetMoney = rbt.BetMoneys[0]

}

func (rbt *RealTimeRobot) onGameStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_GameStatus)

	rbt.GameState = slot.GameState(data.GetState())

	if data.GetState() == 0 {

		go func() {
			time.Sleep(time.Duration(int64(rbt.RobotConfig.GameActionWaitSpanMin)) * time.Second)

			if rbt.PickGame != nil {
				rbt.actionPickGame(robot)
				return
			}

			if rbt.LeftTimes == 0 {
				rbt.GameLogic.UserLeave(robot.UserInfo)
				rbt.cancel()
				return
			}

			if rbt.LeftTimes%10000 == 0 {
				fmt.Printf("[%s] 剩餘 %d 次, %v\n", robot.UserInfo.NickName, rbt.LeftTimes, time.Since(rbt.LastTime))
				rbt.LastTime = time.Now()
			}

			rbt.LeftTimes--
			rbt.actionSpin(robot)
			return
		}()
	}

}

func (rbt *RealTimeRobot) onBetStatus(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *RealTimeRobot) onJapotPool(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *RealTimeRobot) onGameResult(robot *robotcomm.Robot, m *msg.Message) {

	if data, ok := m.MsgData.(*netproto.Slot_MSG_GameResult); ok {
		if data.GetJackpotMoney() > 0 {
			if data.GetBetMoney() == 0 {
				rbt.FreeGameJPTimes++
				rbt.FreeGameJPWin += data.GetJackpotMoney()
			} else {
				rbt.MainGameJPTimes++
				rbt.MainGameJPWin += data.GetJackpotMoney()
			}
		}
	}

	if _, ok := m.MsgData.(*netproto.Slot_MSG_StepGameResult); ok {

	}
}

func (rbt *RealTimeRobot) onPickGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.TGW_MSG_BonusGame)

	rbt.PickGame = data
}

func (rbt *RealTimeRobot) actionSpin(robot *robotcomm.Robot) {

	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_ActionSpinID)

	m := new(netproto.TGW_MSG_Spin)

	var money int64
	if rbt.RandomBet {
		{
			r := new(prob.Random)
			for _, row := range rbt.BetMoneys {
				r.Add(row, 1)
			}
			t := r.Choose()
			money = t.(int64)
		}
	} else {
		money = rbt.BetMoney
	}

	m.Money = proto.Int64(money)

	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

func (rbt *RealTimeRobot) actionPickGame(robot *robotcomm.Robot) {

	records := rbt.PickGame.GetRecords()
	rbt.PickGame = nil

	for i, row := range records {
		if row.GetLocation() == 0 {
			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_ActionPickID)

			detail := &netproto.Slot_MSG_PickBonus{}

			m := new(netproto.Slot_BonusGameRecord)
			m.ID = proto.Int32(row.GetID())
			m.Location = proto.Int32(int32(i + 1))
			m.RewareValue = proto.Int64(int64(row.GetRewareValue()))

			detail.Record = m

			b, _ := proto.Marshal(detail)

			_m.MsgData = b

			rbt.SendGameMessageNow(robot.UserInfo, _m)
		}
	}

}
