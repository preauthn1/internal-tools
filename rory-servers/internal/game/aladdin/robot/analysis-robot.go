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

// AnalyseRobot WLZB robot logic
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

	RandomBet bool

	stop bool
}

// NewAnalysisRobot Create WLZB robot logic
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

	// 遊戲結果
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameResultID) {
		rbt.onGameResult(robot, m)
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

		if rbt.LeftTimes == 0 {
			fmt.Printf("[%s] 剩餘 %d 次, %v\n", robot.UserInfo.NickName, rbt.LeftTimes, time.Since(rbt.LastTime))
			rbt.LastTime = time.Now()
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

func (rbt *AnalyseRobot) onGameResult(robot *robotcomm.Robot, m *msg.Message) {}

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

	m.Money = proto.Int64(money)
	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

// PrintReportByLogs 印出結果
func (rbt *AnalyseRobot) PrintReportByLogs(w io.Writer, logs []*netproto.Slot_MSG_GameLogAdd) {

	var TotalWinMoney, TotalBetMoney, MainGameWinMoney, RespinWinMoney int64
	var TotalTimes, MainGameTimes, RespinTimes int64

	for _, row := range logs {
		TotalWinMoney += row.GetWinMoney()
		TotalBetMoney += row.GetBetMoney()

		switch int(row.GetGameType()) {
		case slot.HistoryTypeMAIN:
			TotalTimes++
			MainGameWinMoney += row.GetWinMoney()
			MainGameTimes++
		case slot.HistoryTypeRESPIN:
			RespinWinMoney += row.GetWinMoney()
			RespinTimes++
		}
	}

	fmt.Fprintln(w, "===================== ALD =====================")
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
	fmt.Fprintf(w, "RespinTimes: %d, RespinWinMoney: %d (RTP: %.2f%%)\n", RespinTimes, RespinWinMoney, float64(RespinWinMoney)/float64(TotalBetMoney)*100)
	fmt.Fprintln(w, "==================================================")
}
