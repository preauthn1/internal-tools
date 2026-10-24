package robot

import (
	"fmt"
	"io"
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

// AnalyseRobot DFDC robot logic
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

	PickGame *netproto.TGW_MSG_BonusGame

	RandomBet bool

	stop bool
}

// NewAnalysisRobot Create TGW robot logic
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

func (rbt *AnalyseRobot) onSlotSetting(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.TGW_MSG_Setting)

	rbt.BetMoneys = make([]int64, 0)
	for _, row := range data.GetBetCount() {
		rbt.BetMoneys = append(rbt.BetMoneys, int64(row))
	}

	rbt.BetMoney = rbt.BetMoneys[0]
}

func (rbt *AnalyseRobot) onGameStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_GameStatus)

	if data.GetState() == 0 {
		if rbt.PickGame != nil {
			rbt.actionPickGame(robot)
			return
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

}

func (rbt *AnalyseRobot) onBetStatus(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *AnalyseRobot) onJapotPool(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *AnalyseRobot) onGameResult(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *AnalyseRobot) onPickGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.TGW_MSG_BonusGame)

	rbt.PickGame = data
}

func (rbt *AnalyseRobot) actionSpin(robot *robotcomm.Robot) {

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

func (rbt *AnalyseRobot) actionPickGame(robot *robotcomm.Robot) {

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

// PrintReportByLogs 印出 SlotHistoryLog 結果
func (rbt *AnalyseRobot) PrintReportByLogs(w io.Writer, logs []*netproto.Slot_MSG_GameLogAdd) {

	var TotalWinMoney, TotalBetMoney, MainGameWinMoney, PickGameWinMoney int64
	var TotalTimes, MainGameTimes, PickGameTimes int64

	for _, row := range logs {
		TotalWinMoney += row.GetWinMoney()
		TotalBetMoney += row.GetBetMoney()

		switch int(row.GetGameType()) {
		case slot.HistoryTypeMAIN:
			TotalTimes++
			MainGameWinMoney += row.GetWinMoney()
			MainGameTimes++
		case slot.HistoryTypeBONUS:
			PickGameWinMoney += row.GetWinMoney()
			PickGameTimes++
		}
	}

	fmt.Fprintln(w, "===================== TGW =====================")
	if rbt.RandomBet {
		fmt.Fprintf(w, "Times: %d, RandomBet, RTP: %.4f%%\n", TotalTimes, float64(TotalWinMoney)/float64(TotalBetMoney)*100)
	} else {
		fmt.Fprintf(w, "Times: %d, EachBetMoney: %d, RTP: %.4f%% \n", rbt.TotalTimes, rbt.BetMoney*rbt.BetRate, float64(TotalWinMoney)/float64(TotalBetMoney)*100)
	}
	fmt.Fprintf(w, "BetMoneys: %v\n", rbt.BetMoneys)
	fmt.Fprintf(w, "BetRates: %v\n", rbt.BetRates)
	fmt.Fprintf(w, "BeginTime: %s\n", rbt.BeginTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Duration: %v\n", time.Since(rbt.BeginTime))
	fmt.Fprintln(w, "=================================================")
	fmt.Fprintf(w, "TotalBetMoney: %d, TotalWinMoney: %d\n", TotalBetMoney, TotalWinMoney)
	fmt.Fprintf(w, "MainGameTimes: %d, MainGameWinMoney: %d (RTP: %.2f%%)\n", MainGameTimes, MainGameWinMoney, float64(MainGameWinMoney)/float64(TotalBetMoney)*100)
	fmt.Fprintf(w, "PickGameTimes: %d, PickGameWinMoney: %d (RTP: %.2f%%)\n", PickGameTimes, PickGameWinMoney, float64(PickGameWinMoney)/float64(TotalBetMoney)*100)
	fmt.Fprintln(w, "=================================================")

}
