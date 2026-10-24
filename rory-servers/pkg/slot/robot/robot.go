package robot

import (
	"context"
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

	BetMoneys   []int64
	BetRates    []int64
	MainGameIDs []int32

	BetMoney  int64
	BetRate   int64
	SelectIdx int

	GameState   slot.GameState
	BetStatus   int32
	LastSendOpt *netproto.Slot_ChooseGameOption

	LeftFreeGame       *netproto.Slot_MSG_LeftFreeGameData
	SelectGameOptions  *netproto.Slot_MSG_ChooseGameOptions
	SelectGameOptions2 *netproto.WLZB_MSG_ChooseGameOptions
	Lightning          *netproto.Slot_MSG_LightningPrepare
	JackpotGame        *netproto.Slot_MSG_BonusGame
	PickGame           *netproto.Slot_MSG_BonusGame

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

	Debug bool
}

// NewRealTimeRobot Create Robot
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

	// 輪軸設定
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameReelSettingID) {
		rbt.onGameReelSetting(robot, m)
	}

	// 初始設定
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SettingID) {
		rbt.onSlotSetting(robot, m)
		return
	}

	// 遊戲狀態
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_StatusID) {
		rbt.onGameStatus(robot, m)
		return
	}

	// 下注狀態
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BetStatusID) {
		rbt.onBetStatus(robot, m)
		return
	}

	// Jackpot pool
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_JackpotPoolRetID) {
		rbt.onJapotPool(robot, m)
		return
	}

	// 遊戲結果
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameResultID) {
		rbt.onGameResult(robot, m)
		return
	}

	// 免費遊戲
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
		rbt.onLeftFreeGame(robot, m)
		return
	}

	// 選擇遊戲
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) || m.SClassID == int32(netproto.WLZB_GameMessageClassID_WLZB_GameOptionsID) {
		rbt.onSelection(robot, m)
		return
	}

	// LightningPrepare
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LightningPrepareID) {
		rbt.onLightningPrepare(robot, m)
		return
	}

	// LightningResult
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LightningResultID) {
		rbt.onLightningResult(robot, m)
		return
	}

	// Jackpot遊戲
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameID) {
		rbt.onJackpotGame(robot, m)
		rbt.onPickGame(robot, m)
		return
	}
}

func (rbt *RealTimeRobot) onGameReelSetting(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_ReelSetting)

	rbt.MainGameIDs = make([]int32, 0)
	for _, row := range data.GameReels {
		rbt.MainGameIDs = append(rbt.MainGameIDs, row.GetID())
	}
}

func (rbt *RealTimeRobot) onSlotSetting(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_Setting)

	rbt.BetMoneys = make([]int64, 0)
	for _, row := range data.GetBetMoneys() {
		rbt.BetMoneys = append(rbt.BetMoneys, int64(row))
	}

	rbt.BetRates = make([]int64, 0)
	for _, row := range data.GetBetRates() {
		rbt.BetRates = append(rbt.BetRates, int64(row))
	}

	if rbt.BetMoney == 0 {
		rbt.BetMoney = int64(rbt.BetMoneys[0])
	}
	if rbt.BetRate == 0 {
		rbt.BetRate = int64(rbt.BetRates[0])
	}

	rbt.GoldBall = data.GetGoldBallSetting()

	rbt.NetprotoFreeGameInfos = data.FreeGames
}

func (rbt *RealTimeRobot) onGameStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_Status)

	rbt.GameState = slot.GameState(data.GetState())

	go func() {
		time.Sleep(time.Duration(int64(rbt.RobotConfig.GameActionWaitSpanMin)) * time.Second)

		if data.GetState() == 0 {

			if rbt.LeftFreeGame != nil {
				rbt.LeftFreeGame = nil
				rbt.actionSpin(robot)
			}

			rbt.actionSpin(robot)
			return
		}

		if data.GetState() == 5 {

			if rbt.JackpotGame != nil {
				rbt.actionPickJackpot(robot)
				return
			}

			if rbt.PickGame != nil {
				rbt.actionPickGame(robot)
				return
			}

			if rbt.SelectGameOptions != nil || rbt.SelectGameOptions2 != nil {
				rbt.actionSelectGame(robot)
				return
			}

			if rbt.Lightning != nil {
				rbt.actionLightning(robot)
				return
			}
		}
	}()

}

func (rbt *RealTimeRobot) onBetStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_BetStatus)
	rbt.BetStatus = data.GetRet()
}

func (rbt *RealTimeRobot) onJapotPool(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_JackpotPool)

	if rbt.JackpotPool == nil {
		rbt.JackpotPool = make(map[slot.JackpotLevel]*netproto.Slot_JackpotInfo)
	}

	for _, row := range data.Jackpots {
		lv := slot.JackpotLevel(row.GetName())
		rbt.JackpotPool[lv] = row
	}

	return
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
	data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
	if data.GetCount() > 0 {
		rbt.LeftFreeGame = data
	} else {
		rbt.LeftFreeGame = nil
	}
}

func (rbt *RealTimeRobot) onSelection(robot *robotcomm.Robot, m *msg.Message) {

	switch m.SClassID {
	case int32(netproto.WLZB_GameMessageClassID_WLZB_GameOptionsID):
		defer func() {
			if rec := recover(); rec != nil {
				return
			}
		}()
		data := m.MsgData.(*netproto.WLZB_MSG_ChooseGameOptions)
		if data.GetTimes() > 0 {
			rbt.SelectGameOptions2 = data
		} else {
			rbt.SelectGameOptions2 = nil
		}

	case int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID):
		data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
		rbt.SelectGameOptions = data

	}

}

func (rbt *RealTimeRobot) onLightningPrepare(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_LightningPrepare)

	if data.GetCount() > 0 {
		rbt.Lightning = data
	} else {
		rbt.Lightning = nil
	}
}

func (rbt *RealTimeRobot) onLightningResult(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_LightningResult)

	if data.GetJackpotMoney() > 0 {
		rbt.LightningJPTimes++
		rbt.LightningJPWin += data.GetJackpotMoney()
	}
}

func (rbt *RealTimeRobot) onJackpotGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_BonusGame)
	if data.GetType() == string(slot.GameModeJackpotGame) {
		rbt.JackpotGame = data
	}
}

func (rbt *RealTimeRobot) onPickGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_BonusGame)
	if data.GetType() == string(slot.GameModePickGame) {
		rbt.PickGame = data
	}
}

func (rbt *RealTimeRobot) actionSpin(robot *robotcomm.Robot) {

	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

	m := new(netproto.Slot_MSG_Spin)

	var money int64
	if rbt.RandomBet {
		r := new(prob.Random)
		for _, row := range rbt.BetMoneys {
			r.Add(row, 1)
		}
		t := r.Choose()
		money = t.(int64)
	} else {
		money = rbt.BetMoney
	}

	var rate int64
	var id int32
	if rbt.RandomBet {
		r := new(prob.Random)
		for _, row := range rbt.BetRates {
			r.Add(row, 1)
		}
		t := r.Choose()
		rate = t.(int64)
	} else {
		rate = rbt.BetRate
		id = rbt.MainGameIDs[0]
	}

	m.Money = proto.Int64(money)
	m.Rate = proto.Int32(int32(rate))
	m.ID = proto.Int32(id)

	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

func (rbt *RealTimeRobot) actionPickJackpot(robot *robotcomm.Robot) {

	records := rbt.JackpotGame.Records
	rbt.JackpotGame = nil

	for i, row := range records {
		if row.GetLocation() == 0 {
			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_PickBonusID)

			m := new(netproto.Slot_MSG_PickBonus)
			m.Type = proto.String(string(slot.GameModeJackpotGame))
			m.Record = new(netproto.Slot_BonusGameRecord)
			m.Record.ID = proto.Int32(row.GetID())
			m.Record.Location = proto.Int32(int32(i + 1))
			m.Record.RewareValue = proto.Int64(row.GetRewareValue())

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			rbt.SendGameMessageNow(robot.UserInfo, _m)
		}
	}
}

func (rbt *RealTimeRobot) actionPickGame(robot *robotcomm.Robot) {

	records := rbt.PickGame.Records
	rbt.PickGame = nil

	for i, row := range records {
		if row.GetLocation() == 0 {
			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_PickBonusID)

			m := new(netproto.Slot_MSG_PickBonus)
			m.Type = proto.String(string(slot.GameModePickGame))
			m.Record = new(netproto.Slot_BonusGameRecord)
			m.Record.ID = proto.Int32(row.GetID())
			m.Record.Location = proto.Int32(int32(i + 1))
			m.Record.RewareValue = proto.Int64(row.GetRewareValue())

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			rbt.SendGameMessageNow(robot.UserInfo, _m)
		}
	}
}

func (rbt *RealTimeRobot) actionSelectGame(robot *robotcomm.Robot) {
	_m := new(msg.Message)

	if rbt.SelectGameOptions != nil {
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

		m := new(netproto.Slot_MSG_ChooseGame)

		if rbt.RandomSelection {
			r := new(prob.Random)
			for _, row := range rbt.SelectGameOptions.Options {
				r.Add(row, 1)
			}
			t := r.Choose()

			m.Option = t.(*netproto.Slot_ChooseGameOption)
		} else {
			m.Option = rbt.SelectGameOptions.Options[rbt.SelectIdx]
		}

		rbt.LastSendOpt = m.Option

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		rbt.SelectGameOptions = nil
	} else if rbt.SelectGameOptions2 != nil {
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.WLZB_GameMessageClassID_WLZB_ChooseGameID)

		m := new(netproto.WLZB_MSG_ChooseGame)

		if rbt.RandomBet {
			r := new(prob.Random)
			for _, row := range rbt.NetprotoFreeGameInfos {
				r.Add(row, 1)
			}
			t := r.Choose()

			m.ID = proto.Int32(t.(*netproto.Slot_FreeGameInfo).GetID())
		} else {
			m.ID = proto.Int32(rbt.NetprotoFreeGameInfos[rbt.SelectIdx].GetID())
		}

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		rbt.SelectGameOptions2 = nil
	}

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

func (rbt *RealTimeRobot) actionLightning(robot *robotcomm.Robot) {

	tmp := rbt.Lightning
	rbt.Lightning = nil

	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionLightningSpinID)

	m := new(netproto.XMBX_MSG_LightningShow)
	m.LightningStep = proto.Int32(tmp.GetLightningStep())

	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}
