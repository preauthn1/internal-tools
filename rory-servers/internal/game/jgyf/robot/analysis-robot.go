package robot

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"
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

// AnalyseRobot WLJR robot logic
type AnalyseRobot struct {
	*robotlogic.LogicBase

	Terminate chan bool

	BeginTime time.Time
	LastTime  time.Time

	TotalTimes int
	LeftTimes  int

	BetMoneys []int64
	BetRates  []int64

	BetMoney int64
	BetRate  int64

	LeftFreeGame      *netproto.Slot_MSG_LeftFreeGameData
	SelectGameOptions *netproto.Slot_MSG_ChooseGameOptions
	JackpotGame       *netproto.Slot_MSG_BonusGame

	RandomBet bool

	stop bool
}

// NewAnalysisRobot Create WLJR robot logic
func NewAnalysisRobot(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher, times int, random bool) *AnalyseRobot {
	rbt := new(AnalyseRobot)
	rbt.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rbt.TotalTimes = times
	rbt.LeftTimes = times
	rbt.RandomBet = random
	rbt.Terminate = make(chan bool)

	return rbt
}

// OnGameMessage 收到遊戲訊息
func (rbt *AnalyseRobot) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	if rbt.stop {
		return
	}

	// 遊戲結束
	if m.SClassID == int32(netproto.GameMessageClassID_GameEndID) {
		return
	}

	// 初始設定
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SettingID) {
		rbt.BeginTime = time.Now()
		rbt.LastTime = time.Now()
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
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
		rbt.onSelection(robot, m)
		return
	}

	// Jackpot遊戲
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameID) {
		rbt.onJackpotGame(robot, m)
		return
	}
}

func (rbt *AnalyseRobot) onSlotSetting(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_Setting)

	rbt.BetMoneys = make([]int64, 0)
	for _, row := range data.GetBetMoneys() {
		rbt.BetMoneys = append(rbt.BetMoneys, int64(row))
	}

	rbt.BetRates = make([]int64, 0)
	for _, row := range data.GetBetRates() {
		rbt.BetRates = append(rbt.BetRates, int64(row))
	}

	rbt.BetMoney = int64(rbt.BetMoneys[0])
	rbt.BetRate = int64(rbt.BetRates[0])
}

func (rbt *AnalyseRobot) onGameStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_Status)

	if data.GetState() == 0 {

		if rbt.LeftFreeGame != nil {
			rbt.LeftFreeGame = nil
			rbt.actionSpin(robot)
		}

		if rbt.LeftTimes == 0 {
			rbt.Terminate <- true
			rbt.stop = true
			rbt.GameLogic.UserLeave(robot.UserInfo)
			return
		}

		if rbt.LeftTimes%10000 == 0 {
			fmt.Printf("[%s] 剩餘 %d 次, %v\n", robot.UserInfo.NickName, rbt.LeftTimes, time.Since(rbt.LastTime))
			rbt.LastTime = time.Now()
		}

		rbt.LeftTimes--
		rbt.actionSpin(robot)
		return
	}

	if data.GetState() == 5 {

		if rbt.JackpotGame != nil {
			rbt.actionPickJackpot(robot)
			return
		}

		if rbt.SelectGameOptions != nil {
			rbt.actionSelectGame(robot)
			return
		}
	}
}

func (rbt *AnalyseRobot) onBetStatus(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *AnalyseRobot) onJapotPool(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *AnalyseRobot) onGameResult(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *AnalyseRobot) onLeftFreeGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
	if data.GetCount() > 0 {
		rbt.LeftFreeGame = data
	} else {
		rbt.LeftFreeGame = nil
	}

}

func (rbt *AnalyseRobot) onSelection(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
	rbt.SelectGameOptions = data
}

func (rbt *AnalyseRobot) onJackpotGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_BonusGame)
	rbt.JackpotGame = data
}

func (rbt *AnalyseRobot) actionSpin(robot *robotcomm.Robot) {

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
	if rbt.RandomBet {
		r := new(prob.Random)
		for _, row := range rbt.BetRates {
			r.Add(row, 1)
		}
		t := r.Choose()
		rate = t.(int64)
	} else {
		rate = rbt.BetRate
	}

	m.Money = proto.Int64(money)
	m.Rate = proto.Int32(int32(rate))
	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

func (rbt *AnalyseRobot) actionPickJackpot(robot *robotcomm.Robot) {

	records := rbt.JackpotGame.Records
	rbt.JackpotGame = nil

	for i, row := range records {
		if row.GetLocation() == 0 {
			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_PickBonusID)

			m := new(netproto.Slot_MSG_PickBonus)
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

func (rbt *AnalyseRobot) actionSelectGame(robot *robotcomm.Robot) {
	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

	m := new(netproto.Slot_MSG_ChooseGame)

	if rbt.RandomBet {
		r := new(prob.Random)
		for _, row := range rbt.SelectGameOptions.Options {
			r.Add(row, 1)
		}
		t := r.Choose()

		m.Option = t.(*netproto.Slot_ChooseGameOption)
	} else {
		m.Option = rbt.SelectGameOptions.Options[0]
	}

	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SelectGameOptions = nil
	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

// PrintReportByLogs 印出結果
func (rbt *AnalyseRobot) PrintReportByLogs(w io.Writer, logs []*netproto.Slot_MSG_GameLogAdd) {

	var TotalWinMoney, TotalBetMoney, MainGameWinMoney, FreeGameWinMoney, JackpotWinMoney, BGWinMoney int64
	var TotalTimes, MainGameTimes, FreeGameTimes, BGTimes int64

	winJackpotLevels := make(map[slot.JackpotLevel]int)

	for _, row := range logs {
		TotalWinMoney += row.GetWinMoney()
		TotalBetMoney += row.GetBetMoney()

		switch int(row.GetGameType()) {
		case slot.HistoryTypeMAIN:
			MainGameWinMoney += row.GetWinMoney()

			switch int(row.GetLogType()) {
			case slot.LogTypeREEL:
				TotalTimes++
				MainGameTimes++
			case slot.LogTypeBG:
				BGTimes++
				BGWinMoney += row.GetWinMoney()
			}
		case slot.HistoryTypeFREE:
			FreeGameWinMoney += row.GetWinMoney()

			if row.GetLogType() == slot.LogTypeREEL {
				FreeGameTimes++
			}
		case slot.HistoryTypeJPNOLINK:
			JackpotWinMoney += row.GetWinMoney()

			switch int(row.GetLogType()) {
			case slot.LogTypeOPEN:
				m := make(map[string]int)
				for _, row := range strings.Split(row.GetGameData(), ",") {
					m[row]++
				}

				for lv, cnt := range m {
					if cnt == 3 {
						level, _ := strconv.Atoi(lv)
						winJackpotLevels[slot.NewJackpotLevel(int32(level))]++
						break
					}
				}
			}

		}
	}

	fmt.Fprintln(w, "====================== JGYF ======================")
	if rbt.RandomBet {
		fmt.Fprintf(w, "Times: %d, RandomBet, RTP: %.4f%%\n", TotalTimes, float64(TotalWinMoney)/float64(TotalBetMoney)*100)
	} else {
		fmt.Fprintf(w, "Times: %d, EachBetMoney: %d, RTP: %.4f%% \n", TotalTimes, rbt.BetMoney*rbt.BetRate, float64(TotalWinMoney)/float64(TotalBetMoney)*100)
	}
	fmt.Fprintf(w, "BetMoneys: %v\n", rbt.BetMoneys)
	fmt.Fprintf(w, "BetRates: %v\n", rbt.BetRates)
	fmt.Fprintf(w, "BeginTime: %s\n", rbt.BeginTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Duration: %v\n", time.Since(rbt.BeginTime))
	fmt.Fprintln(w, "==================================================")
	fmt.Fprintf(w, "TotalBetMoney: %d, TotalWinMoney: %d\n", TotalBetMoney, TotalWinMoney)
	fmt.Fprintf(w, "MainGameTimes: %d, MainGameWinMoney: %d (RTP: %.2f%%)\n", MainGameTimes, MainGameWinMoney, float64(MainGameWinMoney)/float64(TotalBetMoney)*100)
	fmt.Fprintf(w, "FreeGameTimes: %d, FreeGameWinMoney: %d (RTP: %.2f%%)\n", FreeGameTimes, FreeGameWinMoney, float64(FreeGameWinMoney)/float64(TotalBetMoney)*100)
	fmt.Fprintf(w, "JackpotWin: %d (RTP: %.2f%%)\n", JackpotWinMoney, float64(JackpotWinMoney)/float64(TotalBetMoney)*100)
	for lv, cnt := range winJackpotLevels {
		fmt.Fprintf(w, " - Jackpot[%s]: %d\n", lv.String(), cnt)
	}
	fmt.Fprintln(w, "==================================================")
}
