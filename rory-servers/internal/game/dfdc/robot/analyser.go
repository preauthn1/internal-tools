package robot

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
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

// Analyser Slot 分析機器人
type Analyser struct {
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

// NewAnalyser Create Analyser Robot
func NewAnalyser(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher, debug bool) (*Analyser, context.Context, context.CancelFunc) {
	rbt := new(Analyser)
	rbt.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)

	rbt.ctx, rbt.cancel = context.WithCancel(context.Background())

	if debug {
		go func(ctx context.Context) {
			for {
				select {
				default:
					fmt.Printf("GameState: %d\n", rbt.GameState)
					fmt.Printf("BetStatus: %d\n", rbt.BetStatus)
					fmt.Printf("Left Time: %d\n", rbt.LeftTimes)

					fmt.Printf("LeftFreeGame: %v\n", rbt.LeftFreeGame)
					fmt.Printf("JackpotGame: %v\n", rbt.JackpotGame)

					fmt.Printf("LastSendOpt: %v\n", rbt.LastSendOpt)

					time.Sleep(5 * time.Second)
				case <-ctx.Done():
					return
				}
			}
		}(rbt.ctx)
	}

	return rbt, rbt.ctx, rbt.cancel
}

// OnGameMessage 收到遊戲訊息
func (rbt *Analyser) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

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

func (rbt *Analyser) onGoldenSymbol(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_GoldenSymbolList)

	rbt.BetRates = make([]int64, 0)

	for _, row := range data.GetGoldenDetail() {
		rbt.BetRates = append(rbt.BetRates, int64(row.GetCredits()))
	}

	rbt.BetRate = rbt.BetRates[0]
}

func (rbt *Analyser) onSlotSetting(robot *robotcomm.Robot, m *msg.Message) {
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

func (rbt *Analyser) onGameStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_GameStatus)

	rbt.GameState = slot.GameState(data.GetState())

	if data.GetState() == 0 {

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
	}

}

func (rbt *Analyser) onBetStatus(robot *robotcomm.Robot, m *msg.Message) {
	// data := m.MsgData.(*netproto.Slot_MSG_BetStatus)
	// rbt.BetStatus = data.GetRet()
}

func (rbt *Analyser) onJapotPool(robot *robotcomm.Robot, m *msg.Message) {
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

func (rbt *Analyser) onGameResult(robot *robotcomm.Robot, m *msg.Message) {

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

func (rbt *Analyser) onLeftFreeGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_FreeOrPayMsg)
	if data.GetFreeCount() != 0 {
		rbt.LeftFreeGame = data
	} else {
		rbt.LeftFreeGame = nil
	}
}

func (rbt *Analyser) onJackpotGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.DFDC_BonusGame)

	rbt.JackpotGame = data
}

func (rbt *Analyser) actionSpin(robot *robotcomm.Robot) {

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

func (rbt *Analyser) actionPickJackpot(robot *robotcomm.Robot) {

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

// PrintReportByLogs 印出結果
func (rbt *Analyser) PrintReportByLogs(w io.Writer, logs []*netproto.Slot_MSG_GameLogAdd) {

	var printLineCount int = 4

	// basic
	var TotalWinMoney, TotalBetMoney, MainGameWinMoney, TotalTimes, MainGameTimes int64

	// freespin
	var FreeGameWinMoney, FreeGameExtAwardWinMoney, FreeGameTimes int64

	// jp
	var JackpotWinMoney int64

	JackpotWinMoney += rbt.MainGameJPWin
	JackpotWinMoney += rbt.FreeGameJPWin
	JackpotWinMoney += rbt.LightningJPWin

	winJackpotLevels := make(map[slot.JackpotLevel]int)
	symbolGroupMap := make(map[int32]int)
	GoldBallMap := make(map[string]int)
	for _, row := range rbt.GoldBall {
		GoldBallMap[row.GetSymbol()] = 0
	}

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
				symbolGroupMap = rbt.parseSymbolGroup(symbolGroupMap, row.GetGameData())
			}
		case slot.HistoryTypeFREE:
			FreeGameWinMoney += row.GetWinMoney()

			switch int(row.GetLogType()) {
			case slot.LogTypeREEL:
				FreeGameTimes++
				symbolGroupMap = rbt.parseSymbolGroup(symbolGroupMap, row.GetGameData())
			case slot.LogTypeTOTALCOIN:
				FreeGameExtAwardWinMoney += row.GetWinMoney()
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

			case slot.LogTypeBG:
				winJackpotLevels[slot.JackpotGoldReel]++
			}
		}
	}

	fmt.Fprintf(w, "====================== %s ======================\n", rbt.GameName)
	if rbt.RandomBet {
		fmt.Fprintf(w, "Times: %d, RandomBet, RTP: %.4f%%\n", TotalTimes, float64(TotalWinMoney)/float64(TotalBetMoney)*100)
	} else {
		fmt.Fprintf(w, "Times: %d, EachBetMoney: %d, RTP: %.4f%% \n", TotalTimes, rbt.BetMoney*rbt.BetRate, float64(TotalWinMoney)/float64(TotalBetMoney)*100)
	}
	fmt.Fprintf(w, "BetMoneys: %v\n", rbt.BetMoneys)
	fmt.Fprintf(w, "BetRates: %v\n", rbt.BetRates)
	fmt.Fprintf(w, "BeginTime: %s\n", rbt.BeginTime.Format(time.RFC3339))
	fmt.Fprintf(w, "Duration: %v\n", time.Since(rbt.BeginTime))
	fmt.Fprintln(w, "==================================================")

	fmt.Fprintf(w, "TotalBetMoney: %d, TotalWinMoney: %d\n", TotalBetMoney, TotalWinMoney)
	fmt.Fprintf(w, "MainGameTimes: %d, MainGameWinMoney: %d (RTP: %.2f%%)\n", MainGameTimes, MainGameWinMoney, float64(MainGameWinMoney)/float64(TotalBetMoney)*100)

	if rbt.Features&slot.FeatureFreeSpin != 0 {
		fmt.Fprintf(w, "FreeGameTimes: %d, FreeGameWinMoney: %d (RTP: %.2f%%), FreeGameExtAwardWinMoney: %d (RTP: %.2f%%)\n", FreeGameTimes, FreeGameWinMoney, float64(FreeGameWinMoney)/float64(TotalBetMoney)*100, FreeGameExtAwardWinMoney, float64(FreeGameExtAwardWinMoney)/float64(TotalBetMoney)*100)
	}

	if rbt.Features&(slot.FeatureJackpotGame|slot.FeatureJackpot) != 0 {
		fmt.Fprintf(w, "JackpotWin: %d (RTP: %.2f%%)\n", JackpotWinMoney, float64(JackpotWinMoney)/float64(TotalBetMoney)*100)
		if rbt.Features&slot.FeatureJackpot != 0 {
			fmt.Fprintf(w, "MainGameJP: %d, Win: %d (RTP: %.2f%%)\n", rbt.MainGameJPTimes, rbt.MainGameJPWin, float64(rbt.MainGameJPWin)/float64(TotalBetMoney)*100)
			if rbt.Features&slot.FeatureFreeSpin != 0 {
				fmt.Fprintf(w, "FreeGameJP: %d, Win: %d (RTP: %.2f%%)\n", rbt.FreeGameJPTimes, rbt.FreeGameJPWin, float64(rbt.FreeGameJPWin)/float64(TotalBetMoney)*100)
			}
			if rbt.Features&slot.FeatureLightning != 0 {
				fmt.Fprintf(w, "LightningJP: %d, Win: %d (RTP: %.2f%%)\n", rbt.LightningJPTimes, rbt.LightningJPWin, float64(rbt.LightningJPWin)/float64(TotalBetMoney)*100)
			}
		}

		if rbt.Features&slot.FeatureJackpotGame != 0 {
			for lv, cnt := range winJackpotLevels {
				fmt.Fprintf(w, " - Jackpot[%s]: %d\n", lv.String(), cnt)
			}
		}
	}

	fmt.Fprintln(w, "==================================================")
	fmt.Fprintln(w, "Symbol Group:")
	// for groupID, cnt := range symbolGroupMap {
	// 	fmt.Fprintf(w, "%d: %d\n", groupID, cnt)
	// }
	{
		type SymbolGroupOpt struct {
			ID    int32
			Times int
		}

		var SymbolGroupOpts []*SymbolGroupOpt

		for id, cnt := range symbolGroupMap {
			SymbolGroupOpts = append(SymbolGroupOpts, &SymbolGroupOpt{id, cnt})
		}
		sort.Slice(SymbolGroupOpts, func(i, j int) bool {
			return SymbolGroupOpts[i].ID < SymbolGroupOpts[j].ID
		})

		var sliceGBOpt [][]*SymbolGroupOpt = make([][]*SymbolGroupOpt, 0)
		var l = len(SymbolGroupOpts) / printLineCount
		var less int = len(SymbolGroupOpts) % printLineCount
		if less > 0 {
			l++
		}
		for i := 0; i < l; i++ {
			var tmp int = i*printLineCount + printLineCount
			if tmp > len(SymbolGroupOpts) {
				tmp = i*printLineCount + less
			}
			sliceGBOpt = append(sliceGBOpt, SymbolGroupOpts[i*printLineCount:tmp])
		}
		for _, grp := range sliceGBOpt {
			var strs []string
			for _, row := range grp {
				if row == nil {
					continue
				}
				strs = append(strs, fmt.Sprintf("%4d: %4d", row.ID, row.Times))
			}
			fmt.Fprintf(w, "  %s\n", strings.Join(strs, "\t"))
		}
	}
	fmt.Fprintln(w, "==================================================")
}

func (rbt *Analyser) parseSymbolGroup(m map[int32]int, data string) map[int32]int {

	if m == nil {
		m = make(map[int32]int)
	}

	strs := strings.Split(data, "#")
	if len(strs) != 2 {
		return m
	}

	f, _ := strconv.ParseFloat(strs[0], 32)

	m[int32(f)]++

	return m
}

func (rbt *Analyser) parseReelData(data string) string {

	strs := strings.Split(data, "#")
	if len(strs) != 2 {
		return ""
	}

	return strs[1]
}
