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

type selectionOptDesc struct {
	id   int32
	mode slot.GameMode
	cnt  int
}

func (desc *selectionOptDesc) String() string {
	return fmt.Sprintf("%s(%d)", desc.mode, desc.id)
}

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

	choosedOptions map[int32]*selectionOptDesc

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
	rbt.choosedOptions = make(map[int32]*selectionOptDesc)

	if debug {
		go func(ctx context.Context) {
			for {
				select {
				default:
					fmt.Printf("GameState: %d\n", rbt.GameState)
					fmt.Printf("BetStatus: %d\n", rbt.BetStatus)
					fmt.Printf("Left Time: %d\n", rbt.LeftTimes)

					fmt.Printf("LeftFreeGame: %v\n", rbt.LeftFreeGame)
					fmt.Printf("SelectGameOptions: %v\n", rbt.SelectGameOptions)
					fmt.Printf("JackpotGame: %v\n", rbt.JackpotGame)
					fmt.Printf("PickGame: %v\n", rbt.PickGame)
					fmt.Printf("Lightning: %v\n", rbt.Lightning)

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

	// 輪軸設定
	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameReelSettingID) {
		rbt.onGameReelSetting(robot, m)
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

func (rbt *Analyser) onGameReelSetting(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_ReelSetting)

	rbt.MainGameIDs = make([]int32, 0)
	for _, row := range data.GameReels {
		rbt.MainGameIDs = append(rbt.MainGameIDs, row.GetID())
	}
}

func (rbt *Analyser) onSlotSetting(robot *robotcomm.Robot, m *msg.Message) {
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

func (rbt *Analyser) onGameStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_Status)

	rbt.GameState = slot.GameState(data.GetState())

	if data.GetState() == 0 {

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
}

func (rbt *Analyser) onBetStatus(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_BetStatus)
	rbt.BetStatus = data.GetRet()
}

func (rbt *Analyser) onJapotPool(robot *robotcomm.Robot, m *msg.Message) {
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
	data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
	if data.GetCount() > 0 {
		rbt.LeftFreeGame = data
	} else {
		rbt.LeftFreeGame = nil
	}
}

func (rbt *Analyser) onSelection(robot *robotcomm.Robot, m *msg.Message) {

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

		for _, row := range rbt.NetprotoFreeGameInfos {
			if _, ok := rbt.choosedOptions[row.GetID()]; !ok {
				rbt.choosedOptions[row.GetID()] = &selectionOptDesc{id: row.GetID(), mode: slot.GameModeFreeSpin}
			}
		}

	case int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID):
		data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
		rbt.SelectGameOptions = data

		for _, row := range data.Options {
			if _, ok := rbt.choosedOptions[row.GetID()]; !ok {
				rbt.choosedOptions[row.GetID()] = &selectionOptDesc{id: row.GetID(), mode: slot.GameMode(row.GetType())}
			}
		}
	}

}

func (rbt *Analyser) onLightningPrepare(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_LightningPrepare)

	if data.GetCount() > 0 {
		rbt.Lightning = data
	} else {
		rbt.Lightning = nil
	}
}

func (rbt *Analyser) onLightningResult(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_LightningResult)

	if data.GetJackpotMoney() > 0 {
		rbt.LightningJPTimes++
		rbt.LightningJPWin += data.GetJackpotMoney()
	}
}

func (rbt *Analyser) onJackpotGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_BonusGame)
	if data.GetType() == string(slot.GameModeJackpotGame) {
		rbt.JackpotGame = data
	}
}

func (rbt *Analyser) onPickGame(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.Slot_MSG_BonusGame)
	if data.GetType() == string(slot.GameModePickGame) {
		rbt.PickGame = data
	}
}

func (rbt *Analyser) actionSpin(robot *robotcomm.Robot) {

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

func (rbt *Analyser) actionPickJackpot(robot *robotcomm.Robot) {

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

func (rbt *Analyser) actionPickGame(robot *robotcomm.Robot) {

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

func (rbt *Analyser) actionSelectGame(robot *robotcomm.Robot) {
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

		if desc, ok := rbt.choosedOptions[m.Option.GetID()]; ok {
			desc.cnt++
		} else {
			rbt.choosedOptions[m.Option.GetID()] = &selectionOptDesc{id: m.Option.GetID(), mode: slot.GameMode(m.Option.GetType()), cnt: 1}
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

		if desc, ok := rbt.choosedOptions[m.GetID()]; ok {
			desc.cnt++
		} else {
			rbt.choosedOptions[m.GetID()] = &selectionOptDesc{id: m.GetID(), mode: slot.GameModeFreeSpin, cnt: 1}
		}

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		rbt.SelectGameOptions2 = nil
	}

	rbt.SendGameMessageNow(robot.UserInfo, _m)
}

func (rbt *Analyser) actionLightning(robot *robotcomm.Robot) {

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

// PrintReportByLogs 印出結果
func (rbt *Analyser) PrintReportByLogs(w io.Writer, logs []*netproto.Slot_MSG_GameLogAdd) {

	var printLineCount int = 4

	// basic
	var TotalWinMoney, TotalBetMoney, MainGameWinMoney, TotalTimes, MainGameTimes int64

	// expand
	var ExpandTimes int64

	// freespin
	var FreeGameWinMoney, FreeGameExtAwardWinMoney, FreeGameTimes int64

	// bg
	var BGWinMoney, BGTimes int64

	// jp
	var JackpotWinMoney int64

	// lightning
	var LightningWinMoney, LightningTimes int64

	var PickGameWinMoney, PickGameTimes int64

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
			case slot.LogTypeBG:
				BGTimes++
				BGWinMoney += row.GetWinMoney()
			}
		case slot.HistoryTypeRESPIN:
			switch int(row.GetLogType()) {
			case slot.LogTypeREEL:
				ExpandTimes++
				symbolGroupMap = rbt.parseSymbolGroup(symbolGroupMap, row.GetGameData())
			}
		case slot.HistoryTypeFREE:
			FreeGameWinMoney += row.GetWinMoney()

			switch int(row.GetLogType()) {
			case slot.LogTypeREEL:
				FreeGameTimes++
				symbolGroupMap = rbt.parseSymbolGroup(symbolGroupMap, row.GetGameData())
			case slot.LogTypeBG:
				BGTimes++
				BGWinMoney += row.GetWinMoney()
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
		case slot.HistoryTypeLIGHTNING:
			LightningWinMoney += row.GetWinMoney()
			LightningTimes++

			var symbols []string

			for _, reel := range strings.Split(rbt.parseReelData(row.GetGameData()), "|") {
				for _, s := range strings.Split(reel, ",") {
					symbols = append(symbols, s)
				}
			}

			for _, s := range symbols {
				if _, ok := GoldBallMap[s]; ok {
					GoldBallMap[s]++
				}
			}
		case slot.HistoryTypeBONUS:
			PickGameWinMoney += row.GetWinMoney()
			PickGameTimes++
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

	if rbt.Features&slot.FeatureExpand != 0 {
		fmt.Fprintf(w, "ExpandTimes: %d\n", ExpandTimes)
	}

	if rbt.Features&slot.FeatureBG != 0 {
		fmt.Fprintf(w, "BGWinTimes: %d, BGWinMoney: %d (RTP: %.2f%%)\n", BGTimes, BGWinMoney, float64(BGWinMoney)/float64(TotalBetMoney)*100)
	}

	if rbt.Features&slot.FeatureFreeSpin != 0 {
		fmt.Fprintf(w, "FreeGameTimes: %d, FreeGameWinMoney: %d (RTP: %.2f%%), FreeGameExtAwardWinMoney: %d (RTP: %.2f%%)\n", FreeGameTimes, FreeGameWinMoney, float64(FreeGameWinMoney)/float64(TotalBetMoney)*100, FreeGameExtAwardWinMoney, float64(FreeGameExtAwardWinMoney)/float64(TotalBetMoney)*100)
	}

	if rbt.Features&slot.FeatureLightning != 0 {
		fmt.Fprintf(w, "LightningTimes: %d, LightningWinMoney: %d (RTP: %.2f%%)\n", LightningTimes, LightningWinMoney, float64(LightningWinMoney)/float64(TotalBetMoney)*100)
		fmt.Fprintln(w, "Lightning GoalBall Appear:")

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

		var sliceGBOpt [][]*GoldBallOpt = make([][]*GoldBallOpt, 0)
		var l = len(GoldBallOpts) / printLineCount
		if len(GoldBallOpts)%printLineCount > 0 {
			l++
		}
		for i := 0; i < l; i++ {
			sliceGBOpt = append(sliceGBOpt, GoldBallOpts[i*printLineCount:i*printLineCount+printLineCount])
		}
		for _, grp := range sliceGBOpt {
			var strs []string
			for _, row := range grp {
				if row == nil {
					continue
				}
				strs = append(strs, fmt.Sprintf("%4s: %4d", row.Symbol, row.Times))
			}
			fmt.Fprintf(w, "  %s\n", strings.Join(strs, "\t"))
		}
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

	if rbt.Features&(slot.FeaturePickGame) != 0 {
		fmt.Fprintf(w, "PickGameTimes: %d, PickGameWinMoney: %d (RTP: %.2f%%)\n", PickGameTimes, PickGameWinMoney, float64(PickGameWinMoney)/float64(TotalBetMoney)*100)
	}

	if rbt.Features&slot.FeatureSelection != 0 {
		var cnt int

		var descs []*selectionOptDesc

		for _, desc := range rbt.choosedOptions {
			cnt += desc.cnt
			descs = append(descs, desc)
		}

		sort.Slice(descs, func(i, j int) bool {
			return descs[i].id < descs[j].id
		})

		fmt.Fprintf(w, "SelectTimes: %d\n", cnt)

		var sods [][]*selectionOptDesc = make([][]*selectionOptDesc, 0)
		var l = len(descs) / printLineCount
		var less int = len(descs) % printLineCount
		if less > 0 {
			l++
		}
		for i := 0; i < l; i++ {
			var tmp int = i*printLineCount + printLineCount
			if tmp > len(descs) {
				tmp = i*printLineCount + less
			}
			sods = append(sods, descs[i*printLineCount:tmp])
		}
		for _, grp := range sods {
			var strs []string
			for _, row := range grp {
				if row == nil {
					continue
				}
				strs = append(strs, fmt.Sprintf("%10s: %4d", fmt.Sprintf("%v", row), row.cnt))
			}
			fmt.Fprintf(w, "  %s\n", strings.Join(strs, "\t"))
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
