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

	// LeftFreeGame *netproto.Slot_MSG_LeftFreeGameData
	// JackpotGame  *netproto.Slot_MSG_BonusGame

	LeftFreeGame *netproto.DFDC_FreeOrPayMsg
	JackpotGame  *netproto.DFDC_BonusGame

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

	// 遊戲結束
	if m.SClassID == int32(netproto.GameMessageClassID_GameEndID) {
		return
	}

	// 黃金符號設定
	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCGoldenSymbolID) {
		rbt.onGoldenSymbol(robot, m)
	}

	// 初始設定
	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCSettingID) {
		rbt.BeginTime = time.Now()
		rbt.LastTime = time.Now()
		rbt.onSlotSetting(robot, m)
		return
	}

	// 遊戲狀態
	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCStateID) {
		rbt.onGameStatus(robot, m)
		return
	}

	// 下注狀態
	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCBetRetID) {
		rbt.onBetStatus(robot, m)
		return
	}

	// Jackpot pool
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_JackpotPoolRetID) {
		rbt.onJapotPool(robot, m)
		return
	}

	// 遊戲結果
	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCGameResult) {
		rbt.onGameResult(robot, m)
		return
	}

	// 免費遊戲
	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCFreeOrPayID) {
		rbt.onLeftFreeGame(robot, m)
		return
	}

	// Jackpot遊戲
	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCBonusGameID) {
		rbt.onJackpotGame(robot, m)
		return
	}
}

func (rbt *RealTimeRobot) onGoldenSymbol(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_GoldenSymbolList)

	rbt.BetRates = make([]int64, 0)

	for _, row := range data.GetGoldenDetail() {
		rbt.BetRates = append(rbt.BetRates, int64(row.GetCredits()))
	}

	rbt.BetRate = rbt.BetRates[0]
}

func (rbt *RealTimeRobot) onSlotSetting(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_Setting)

	rbt.BetMoneys = make([]int64, 0)
	for _, row := range data.GetBetCount() {
		rbt.BetMoneys = append(rbt.BetMoneys, int64(row))
	}

	if rbt.BetMoney == 0 {
		rbt.BetMoney = int64(rbt.BetMoneys[0])
	}
	if rbt.BetRate == 0 {
		rbt.BetRate = int64(rbt.BetRates[0])
	}

}

func (rbt *RealTimeRobot) onGameStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_GameStatus)

	rbt.GameState = slot.GameState(data.GetState())

	if data.GetState() == 0 {

		go func() {
			time.Sleep(time.Duration(int64(rbt.RobotConfig.GameActionWaitSpanMin)) * time.Second)

			if rbt.JackpotGame != nil {
				rbt.actionPickJackpot(robot)
				return
			}

			if rbt.LeftFreeGame != nil {
				rbt.LeftFreeGame = nil
				rbt.actionSpin(robot)
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

func (rbt *RealTimeRobot) onBetStatus(robot *robotcomm.Robot, m *msg.Message) {
	// data := m.MsgData.(*netproto.Slot_MSG_BetStatus)
	// rbt.BetStatus = data.GetRet()
}

func (rbt *RealTimeRobot) onJapotPool(robot *robotcomm.Robot, m *msg.Message) {
	// data := m.MsgData.(*netproto.Slot_MSG_JackpotPool)

	// if rbt.JackpotPool == nil {
	// 	rbt.JackpotPool = make(map[slot.JackpotLevel]*netproto.Slot_JackpotInfo)
	// }

	// for _, row := range data.Jackpots {
	// 	lv := slot.JackpotLevel(row.GetName())
	// 	rbt.JackpotPool[lv] = row
	// }

	// return
}

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

func (rbt *RealTimeRobot) onLeftFreeGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_FreeOrPayMsg)
	if data.GetFreeCount() != 0 {
		rbt.LeftFreeGame = data
	} else {
		rbt.LeftFreeGame = nil
	}
}

func (rbt *RealTimeRobot) onJackpotGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_BonusGame)

	rbt.JackpotGame = data
}

func (rbt *RealTimeRobot) actionSpin(robot *robotcomm.Robot) {

	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCGameStartID)

	m := new(netproto.DFDC_StartBet)

	var money, rate int64
	if rbt.RandomBet {
		{
			r := new(prob.Random)
			for _, row := range rbt.BetMoneys {
				r.Add(row, 1)
			}
			t := r.Choose()
			money = t.(int64)
		}
		{
			r := new(prob.Random)
			for _, row := range rbt.BetRates {
				r.Add(row, 1)
			}
			t := r.Choose()
			rate = t.(int64)
		}
	} else {
		money = rbt.BetMoney
		rate = rbt.BetRate
	}

	m.Money = proto.Int64(money)
	m.GoldenSymbol = proto.Int32(int32(rate))

	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

func (rbt *RealTimeRobot) actionPickJackpot(robot *robotcomm.Robot) {

	records := rbt.JackpotGame.BonusDetail
	rbt.JackpotGame = nil

	for i, row := range records {
		if row.GetPlace() == 0 {
			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCBonusGameIconID)

			m := new(netproto.DFDC_BonusGameDetail)
			m.Id = proto.Int32(row.GetId())
			m.Place = proto.Int32(int32(i + 1))
			m.Icon = proto.Int32(row.GetIcon())

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			rbt.SendGameMessageNow(robot.UserInfo, _m)
		}
	}
}
