package robot

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// AnalyseRobot WLZBJP robot logic
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

	LeftFreeGame *netproto.Slot_MSG_LeftFreeGameData
	Lightning    *netproto.Slot_MSG_LightningPrepare

	MainGameJPTimes  int
	MainGameJPWin    int64
	FreeGameJPTimes  int
	FreeGameJPWin    int64
	LightningJPTimes int
	LightningJPWin   int64

	RandomBet bool

	GoldBall []*netproto.Slot_GoldBallInfo

	stop bool
}

// NewAnalysisRobot Create WLZBJP robot logic
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

	rbt.GoldBall = data.GetGoldBallSetting()
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

		if rbt.Lightning != nil {
			rbt.actionLightning(robot)
			return
		}
	}
}

func (rbt *AnalyseRobot) onBetStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_BetStatus)
	if data.GetRet() != 0 {
		mlog.Warn("Bet Status Error %d", data.GetRet())
	}
}

func (rbt *AnalyseRobot) onJapotPool(robot *robotcomm.Robot, m *msg.Message) {}

func (rbt *AnalyseRobot) onGameResult(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_GameResult)

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

func (rbt *AnalyseRobot) onLeftFreeGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)

	if data.GetCount() > 0 {
		rbt.LeftFreeGame = data
	} else {
		rbt.LeftFreeGame = nil
	}
}

func (rbt *AnalyseRobot) onLightningPrepare(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_LightningPrepare)

	if data.GetCount() > 0 {
		rbt.Lightning = data
	} else {
		rbt.Lightning = nil
	}
}

func (rbt *AnalyseRobot) onLightningResult(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_LightningResult)

	if data.GetJackpotMoney() > 0 {
		rbt.LightningJPTimes++
		rbt.LightningJPWin += data.GetJackpotMoney()
	}
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

	m.Money = proto.Int64(money)
	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

func (rbt *AnalyseRobot) actionLightning(robot *robotcomm.Robot) {

	tmp := rbt.Lightning
	rbt.Lightning = nil

	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionLightningSpinID)

	m := new(netproto.Slot_MSG_LightningShow)
	m.LightningStep = proto.Int32(tmp.GetLightningStep())

	b, _ := proto.Marshal(m)

	_m.MsgData = b

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

// PrintReportByLogs 印出結果
func (rbt *AnalyseRobot) PrintReportByLogs(w io.Writer, logs []*netproto.Slot_MSG_GameLogAdd) {

	var TotalWinMoney, TotalBetMoney, MainGameWinMoney, FreeGameWinMoney, JackpotWinMoney, LightningWinMoney int64
	var TotalTimes, MainGameTimes, FreeGameTimes, LightningTimes int64

	JackpotWinMoney += rbt.MainGameJPWin
	JackpotWinMoney += rbt.FreeGameJPWin
	JackpotWinMoney += rbt.LightningJPWin

	GoldBallMap := make(map[string]int)

	for _, row := range rbt.GoldBall {
		GoldBallMap[row.GetSymbol()] = 0
	}

	for _, row := range logs {
		TotalWinMoney += row.GetWinMoney()
		TotalBetMoney += row.GetBetMoney()

		switch int(row.GetGameType()) {
		case slot.HistoryTypeMAIN:
			TotalTimes++
			MainGameWinMoney += row.GetWinMoney()
			MainGameTimes++
		case slot.HistoryTypeFREE:
			FreeGameWinMoney += row.GetWinMoney()
			FreeGameTimes++
		case slot.HistoryTypeJPNOLINK:
			JackpotWinMoney += row.GetWinMoney()
		case slot.HistoryTypeLIGHTNING:
			LightningWinMoney += row.GetWinMoney()
			LightningTimes++

			var symbols []string
			for _, reel := range strings.Split(row.GetGameData(), "|") {
				for _, s := range strings.Split(reel, ",") {
					symbols = append(symbols, s)
				}
			}

			for _, s := range symbols {
				if _, ok := GoldBallMap[s]; ok {
					GoldBallMap[s]++
				}
			}
		}
	}

	type GoldBallOpt struct {
		Symbol string
		Times  int
	}

	var GoldBallOpts []*GoldBallOpt
	for s, cnt := range GoldBallMap {
		GoldBallOpts = append(GoldBallOpts, &GoldBallOpt{s, cnt})
	}
	sort.Slice(GoldBallOpts, func(i, j int) bool {

		a := GoldBallOpts[i].Symbol[1:]
		b := GoldBallOpts[j].Symbol[1:]

		_a, _ := strconv.ParseInt(a, 0, 8)
		_b, _ := strconv.ParseInt(b, 0, 8)

		return _a < _b
	})

	fmt.Fprintln(w, "====================== XMBX ======================")
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
	fmt.Fprintf(w, "LightningTimes: %d, LightningWinMoney: %d (RTP: %.2f%%)\n", LightningTimes, LightningWinMoney, float64(LightningWinMoney)/float64(TotalBetMoney)*100)
	fmt.Fprintf(w, "JackpotWin: %d (RTP: %.2f%%)\n", JackpotWinMoney, float64(JackpotWinMoney)/float64(TotalBetMoney)*100)
	fmt.Fprintf(w, "MainGameJP: %d, Win: %d (RTP: %.2f%%)\n", rbt.MainGameJPTimes, rbt.MainGameJPWin, float64(rbt.MainGameJPWin)/float64(TotalBetMoney)*100)
	fmt.Fprintf(w, "FreeGameJP: %d, Win: %d (RTP: %.2f%%)\n", rbt.FreeGameJPTimes, rbt.FreeGameJPWin, float64(rbt.FreeGameJPWin)/float64(TotalBetMoney)*100)
	fmt.Fprintf(w, "LightningJP: %d, Win: %d (RTP: %.2f%%)\n", rbt.LightningJPTimes, rbt.LightningJPWin, float64(rbt.LightningJPWin)/float64(TotalBetMoney)*100)
	fmt.Fprintln(w, "==================================================")
	fmt.Fprintln(w, "Lightning GoalBall Appear:")
	{
		var splited [][]*GoldBallOpt

		for i := 0; i < len(GoldBallOpts)/5+1; i++ {
			splited = append(splited, GoldBallOpts[5*i:5*(i+1)])
		}

		for _, row := range splited {
			for _, opt := range row {
				if opt == nil {
					continue
				}
				fmt.Fprintf(w, "%5s: %d\t", opt.Symbol, opt.Times)
			}
			fmt.Fprintln(w)
		}
	}
	fmt.Fprintln(w, "==================================================")
}
