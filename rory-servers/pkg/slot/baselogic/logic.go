package baselogic

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	apipb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// SlotLogicBase Slot Game 專用的基礎GameLoigc
type SlotLogicBase struct {
	*gamelogic.SlotGameLogicBase
	GameMessageFunc      func(user *gamelogic.GameUser, gmsg *msg.Message)
	PiSaiMessage         func(user *gamelogic.GameUser, gmsg *msg.Message)
	OnRecvRPCMessageFunc func(m *msg.Message)
	SkipRPCMessageIDs    []int32

	WG *sync.WaitGroup // 等待群組, 當有需要去DB查詢東西回來才做事情可以用

	// 基本上一定有的ProtoData (Default)
	Mode                  PatternMode          // 當前運行的 PatternMode
	Proto                 *protodata.ProtoData // 當前運行的 Pattern Proto
	ModeExpire            time.Time            // 特殊 PatternMode 的過期時間
	ModeWait              chan bool            // 等待UserAssignedMode查詢的chan
	waitTimer             *time.Timer
	bufFeatureGames       map[PatternMode][]interface{} // 玩家所有模式的FeatureGame暫存
	BeforeSwitchProtoFunc []func()
	AfterSwitchProtoFunc  []func()

	SlotPatternID    uint32 // 基本PatternID
	ReleasePatternID uint32 // 放分PatternID
	BitePatternID    uint32 // 咬分PatternID

	BetCounts     []int32                      // 遊戲押注數量
	BetRates      []int32                      // 遊戲押注倍率
	MaxBetCount   int32                        // 最大押注數量
	MaxBetRate    int32                        // 最大押注倍率
	BigWinRule    []int32                      // client動畫播放階段,以Bet為主
	BigWinNotify  int32                        // 當下注16, Total=6400(16*400) 得獎超過6400才廣播.
	EndActionTime int32                        // 結算動畫超時
	MinCarryMoney int64                        // 最低攜帶金額
	SymbolsDesc   map[string]*slot.SymbolDesc  // 符號說明
	PayLineDesc   string                       // 派彩線說明文字
	UnitPrefix    map[string]int64             // 單位倍率對照
	Descriptions  map[string]*slot.Description // 文字訊息
	MultiMode     slot.MultiMode               // 加倍模式

	DefaultSymbolGroupID int32 // 預設SymbolGroupID，以MainGameTable第一筆做基準
	DefaultBgGroupID     int32 // 預設BackgroundGroupID，以MainGameTable第一筆做基準
	DefaultReelNum       int32 // 預設Reel數量，以MainGameTable第一筆做基準
	DefaultRowNum        int32 // 預設Row數量，以MainGameTable第一筆做基準

	PayLineMode       slot.PayLineMode       // 賠付線模式
	PaylineMatcher    slot.PaylineMatcher    // 賠付線比對器
	BackgroundMatcher slot.BackgroundMatcher // 背景比對器

	// 不會變的資訊，預先整理好準備隨時發送用
	NetprotoSlotReelsSetting   map[int32][]*netproto.Slot_ReelSetting
	NetprotoSlotReelSetting    []*netproto.Slot_ReelSetting
	NetprotoSlotSymbolPayTable []*netproto.Slot_SymbolPayTable
	NetprotoSlotPayLine        []*netproto.Slot_PayLine
	NetprotoFreeGameInfo       []*netproto.Slot_FreeGameInfo
	NetprotoGoldBallInfo       []*netproto.Slot_GoldBallInfo

	Jackpots map[slot.JackpotLevel]*slot.Jackpot // Jackpot資料，map[{JackpotLevel}] *Jackpot

	GameState slot.GameState
	GameUser  *gamelogic.GameUser

	GameBetMoney int64 // 玩家下注金額
	GameBetRate  int64 // 玩家下注倍率(黃金符號)

	GameResultManager *slot.GameResultManager // 遊戲結果
	GamePayResult     []*slot.PayResult       // 遊戲中獎內容
	GameMultiDetail   map[slot.Position]int32 // 加倍資訊

	TransactionNo string
	GamePayMoney  int64 // 花費金額
	// GamePayCoupon  int64 // 花費禮碼
	GameWinMoney   int64 // 贏得金額
	GameWinJackpot int64 // 贏得Jackpot

	// 部分遊戲會有的功能
	FreeSpin           *slot.FreeSpin            // 當前免費遊戲
	JackpotGame        *slot.JackpotGame         // Jackpot翻牌遊戲
	PickGame           *slot.PickGame            // 翻牌遊戲
	SelectChances      []*slot.SelectChance      // 選擇機會
	JackpotGameChances []*slot.JackpotGameChance // 玩JackpotGameChance
	LightningGame      *LightningGame            // Lightning 遊戲

	// 檢查遊戲狀態動作
	CheckFreeGameFunc func()

	cacheLog            []*netproto.Slot_MSG_GameLogAdd
	saveHistoryInterval int32
	saveHistoryMutex    sync.RWMutex
}

// NewSlotLogicBase 建立Slot專用的GameLogic
func NewSlotLogicBase(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *SlotLogicBase {
	logic := new(SlotLogicBase)

	logic.SlotGameLogicBase = gamelogic.NewSlotGameLogicBase(t, tfManager, cfgs)

	logic.loadConfigs()
	logic.ModeWait = make(chan bool, 1)
	logic.bufFeatureGames = make(map[PatternMode][]interface{})

	logic.GameBetRate = 1
	if len(logic.BetRates) == 0 {
		logic.MaxBetRate = 1
	}

	logic.setProto(DefaultPattern)

	logic.RegisterAfterGameEndFunc(func() {
		if logic.CheckFreeGameFunc != nil {
			logic.CheckFreeGameFunc()
		} else {
			logic.CheckFreeGame()
		}

		// 檢查當前的遊戲模式
		if !logic.ExistFeatureGame() {
			logic.CheckModeProto()
		}

		// 檢查特殊模式是否結束
		logic.CheckModeExpire()
	})

	logic.cacheLog = make([]*netproto.Slot_MSG_GameLogAdd, 0)

	logic.saveHistoryInterval = logic.GetServerConfig().SaveHistoryInterval

	if logic.saveHistoryInterval <= 0 {
		logic.saveHistoryInterval = 5000
	}

	tfManager.NewNamedMilliSecTicker(logic.saveHistoryInterval, logic.SendCurrentSlotGameHistoryLogs, "SendCurrentSlotGameHistoryLogs")

	return logic
}

// loadConfigs 讀取.dat設定檔資料
func (logic *SlotLogicBase) loadConfigs() {

	if logic.Configs == nil {
		return
	}

	// Slot Configs
	{
		slotconfig := logic.Configs.GetConfig("slotconfig")

		// slotsetting
		{
			cfg := slotconfig.GetConfigSection("slotsetting")

			// SlotPattern
			{
				i, _ := strconv.ParseInt(cfg["slotpattern"], 10, 32)
				logic.SlotPatternID = uint32(i)
			}

			// ReleasePattern
			{
				i, _ := strconv.ParseInt(cfg["releasepattern"], 10, 32)
				logic.ReleasePatternID = uint32(i)
			}

			// BitePattern
			{
				i, _ := strconv.ParseInt(cfg["bitepattern"], 10, 32)
				logic.BitePatternID = uint32(i)
			}

			// Bet Count
			{
				betcounts := strings.Split(cfg["betcount"], ",")
				logic.BetCounts = make([]int32, 0, len(betcounts))
				for _, betcount := range betcounts {
					i, err := strconv.ParseInt(betcount, 10, 32)
					if err != nil {
						continue
					}
					logic.BetCounts = append(logic.BetCounts, int32(i))

					if logic.MaxBetCount < int32(i) {
						logic.MaxBetCount = int32(i)
					}
				}
			}

			// Bet Rate
			{
				betrates := strings.Split(cfg["betrate"], ",")
				logic.BetRates = make([]int32, 0, len(betrates))
				for _, betrate := range betrates {
					i, err := strconv.ParseInt(betrate, 10, 32)
					if err != nil {
						continue
					}
					logic.BetRates = append(logic.BetRates, int32(i))

					if logic.MaxBetRate < int32(i) {
						logic.MaxBetRate = int32(i)
					}
				}
			}

			// BigWinRule
			{
				strs := strings.Split(cfg["bigwinrule"], ",")
				logic.BigWinRule = make([]int32, 0, len(strs))
				for _, str := range strs {
					i, err := strconv.ParseInt(str, 10, 32)
					if err != nil {
						continue
					}
					logic.BigWinRule = append(logic.BigWinRule, int32(i))
				}
			}

			// BigWinNotify
			{
				str := cfg["bigwinnotify"]
				i, _ := strconv.ParseInt(str, 10, 32)
				logic.BigWinNotify = int32(i)
			}

			// EndActionTime
			{
				str := cfg["endactiontime"]
				i, _ := strconv.ParseInt(str, 10, 32)
				logic.EndActionTime = int32(i)
			}
		}

		// symbols
		{
			cfg := slotconfig.GetConfigSection("symbols")

			logic.SymbolsDesc = make(map[string]*slot.SymbolDesc)

			for key, val := range cfg {
				keys := strings.Split(key, ".")
				if len(keys) != 2 {
					continue
				}
				name := keys[0]
				field := keys[1]

				if _, ok := logic.SymbolsDesc[name]; !ok {
					logic.SymbolsDesc[name] = &slot.SymbolDesc{Name: name}
				}

				switch field {
				case "title":
					logic.SymbolsDesc[name].Title = val
				case "desc":
					logic.SymbolsDesc[name].Description = val
				}
			}
		}

		// description
		{
			cfg := slotconfig.GetConfigSection("description")

			logic.Descriptions = make(map[string]*slot.Description)
			for key, val := range cfg {
				keys := strings.Split(key, ".")

				if len(keys) == 1 {
					name := keys[0]

					if _, ok := logic.Descriptions[name]; !ok {
						logic.Descriptions[name] = &slot.Description{Name: name}
					}

					logic.Descriptions[name].Title = name
					logic.Descriptions[name].Description = val

				} else if len(keys) == 2 {
					name := keys[0]
					field := keys[1]

					if _, ok := logic.Descriptions[name]; !ok {
						logic.Descriptions[name] = &slot.Description{Name: name}
					}

					switch field {
					case "title":
						logic.Descriptions[name].Title = val
					case "desc":
						logic.Descriptions[name].Description = val
					}
				}
			}
		}

		// unit prefix
		{
			cfg := slotconfig.GetConfigSection("unit-prefix")

			logic.UnitPrefix = make(map[string]int64)
			for key, val := range cfg {
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					continue
				}
				logic.UnitPrefix[key] = i
			}
		}
	}

	// Game Configs
	{
		gameconfig := logic.Configs.GetConfig("gameserver")

		cfg := gameconfig.GetConfigSection("gameserver")

		// MinCarryMoney
		{
			i, _ := strconv.ParseInt(cfg["mincarrymoney"], 10, 64)
			logic.MinCarryMoney = i * 100
		}
	}
}

func (logic *SlotLogicBase) setProto(mode PatternMode) {

	for _, f := range logic.BeforeSwitchProtoFunc {
		f()
	}

	defer func() {
		for _, f := range logic.AfterSwitchProtoFunc {
			f()
		}

		if mode != DefaultPattern {
			d := logic.ModeExpire.Sub(time.Now())
			logic.waitTimer = time.AfterFunc(d, func() {
				logic.CheckModeExpire()
			})
		} else {
			if logic.waitTimer != nil {
				logic.waitTimer.Stop()
			}
		}
	}()

	switch mode {
	default:
		if logic.Proto != nil && logic.Proto.PatternID == Proto.PatternID {
			return
		}
		// tmp := *Proto
		// logic.Proto = &tmp
		logic.Proto = Proto
	case ReleasePattern:
		if logic.Proto != nil && logic.Proto.PatternID == ReleaseProto.PatternID {
			return
		}
		if ReleaseProto == nil {
			mlog.Error("尚未設定 [%v] 模式 Pattern", mode)
			return
		}
		// tmp := *ReleaseProto
		// logic.Proto = &tmp
		logic.Proto = ReleaseProto
	case BitePattern:
		if logic.Proto != nil && logic.Proto.PatternID == BiteProto.PatternID {
			return
		}
		if BiteProto == nil {
			mlog.Error("尚未設定 [%v] 模式 Pattern", mode)
			return
		}
		// tmp := *BiteProto
		// logic.Proto = &tmp
		logic.Proto = BiteProto
	}

	logic.SlotPatternID = uint32(logic.Proto.PatternID)
	logic.Mode = mode

	if len(logic.Proto.SlotPattern) == 0 {
		panic(fmt.Errorf("Could not find SlotPatten"))
	}

	if len(logic.Proto.MainGameTable) == 0 {
		panic(fmt.Errorf("Could not find MainGameTable"))
	}

	// 設定預設值
	logic.DefaultSymbolGroupID = int32(logic.Proto.MainGameTable[0].GetSymbolGroupId())
	logic.DefaultBgGroupID = int32(logic.Proto.MainGameTable[0].GetGoldReelGroupId())
	logic.DefaultReelNum = int32(logic.Proto.MainGameTable[0].GetReelNum())
	logic.DefaultRowNum = int32(logic.Proto.MainGameTable[0].GetRowNum())

	// 設定賠付線
	logic.PayLineMode = slot.PayLineMode(logic.Proto.SlotPattern[0].GetLineType())

	switch logic.PayLineMode {
	case slot.PayLineModeAll:
		logic.PaylineMatcher = slot.NewPaylineAllMatcher(logic.Proto.SymbolAttributes)
	case slot.PayLineModeFixed:
		logic.PaylineMatcher = slot.NewPaylineFixedMatcher(logic.Proto.SymbolAttributes, logic.Proto.SlotPayLine)
		logic.BetRates = []int32{int32(len(logic.Proto.SlotPayLine))}
		logic.MaxBetRate = int32(len(logic.Proto.SlotPayLine))
	}

	logic.BackgroundMatcher = new(slot.BasicBackgroundMatcher)

	// 建立MainGame排列帶
	logic.NetprotoSlotReelsSetting = make(map[int32][]*netproto.Slot_ReelSetting)
	for _, row := range logic.Proto.MainGameTable {
		settings := logic.BuildReelSetting(int32(row.GetSymbolGroupId()), int32(row.GetReelNum()), int32(row.GetRowNum()), logic.Proto.SlotSymbolDistribution)
		logic.NetprotoSlotReelsSetting[int32(row.GetId())] = settings
	}

	logic.NetprotoSlotReelSetting = logic.BuildReelSetting(logic.DefaultSymbolGroupID, logic.DefaultReelNum, logic.DefaultRowNum, logic.Proto.SlotSymbolDistribution)
	logic.NetprotoSlotSymbolPayTable = logic.BuildSymbolPayTable(logic.Proto.SlotPayTable, logic.SymbolsDesc)
	logic.NetprotoSlotPayLine = logic.BuildPayLines(logic.Proto.SlotPayLine)
	logic.NetprotoFreeGameInfo = logic.BuildFreeGameInfos(logic.Proto.FreeSpinTable)
	logic.NetprotoGoldBallInfo = logic.BuildGoldBallInfos(logic.Proto.SymbolAttributes)

	logic.extendJackpot()

	mlog.Debug("設定 [%v] 模式 Pattern: %d", mode, logic.SlotPatternID)
}

// initJackpot 依照proto設定取得Jackpot基本資料，初始時只建立共有多少Jackpot
func (logic *SlotLogicBase) initJackpot() {
	logic.Jackpots = make(map[slot.JackpotLevel]*slot.Jackpot)

	for _, row := range logic.Proto.JackpotTable {
		jackpot := new(slot.Jackpot)
		jackpot.Origin = row
		if len(row.GetLevel()) == 0 {
			jackpot.Level = slot.JackpotSingle
		} else {
			jackpot.Level = slot.JackpotLevel(row.GetLevel())
		}
		logic.Jackpots[jackpot.Level] = jackpot
	}
}

// extendJackpot 切換 Proto 時繼承原有的Jackpot, 只改變Origin表
func (logic *SlotLogicBase) extendJackpot() {
	if logic.Jackpots == nil || len(logic.Jackpots) == 0 {
		logic.initJackpot()
		return
	}

	for _, row := range logic.Proto.JackpotTable {
		var level slot.JackpotLevel
		if len(row.GetLevel()) == 0 {
			level = slot.JackpotSingle
		} else {
			level = slot.JackpotLevel(row.GetLevel())
		}

		if jp, ok := logic.Jackpots[level]; ok {
			jp.Origin = row
		}
	}
}

// Clean game data
func (logic *SlotLogicBase) Clean() {

	mlog.Debug("重置遊戲牌局資料")

	logic.GameBetMoney = 0
	if len(logic.BetRates) > 0 {
		logic.GameBetRate = int64(logic.BetRates[0])
	} else {
		logic.GameBetRate = 1
	}

	logic.GameResultManager = nil
	logic.GamePayResult = make([]*slot.PayResult, 0)

	logic.GamePayMoney = 0
	logic.GameWinMoney = 0
	logic.GameWinJackpot = 0
	logic.GameMultiDetail = nil
	logic.TransactionNo = ""
}

// ExistFeatureGame 當前是否存在特殊遊戲
func (logic *SlotLogicBase) ExistFeatureGame() bool {
	if logic.FreeSpin != nil {
		return true
	}

	if logic.JackpotGame != nil {
		return true
	}

	if logic.PickGame != nil {
		return true
	}

	if logic.LightningGame != nil {
		return true
	}

	if len(logic.SelectChances) > 0 {
		return true
	}

	if len(logic.JackpotGameChances) > 0 {
		return true
	}

	return false
}

// SwitchProto 切換 Slot Pattern
func (logic *SlotLogicBase) SwitchProto(mode PatternMode) *protodata.ProtoData {

	if logic.GameUser == nil {
		return nil
	}

	now := time.Now()
	isExpired := now.After(logic.ModeExpire)
	if isExpired && mode != DefaultPattern {
		// example.invalid("桌號(%d) 玩家(%d) 指定切換為 [%v] 模式, 因指定到期時間已過, 取消切換", logic.TableNo(), logic.GameUser.UserID, mode)
		return nil
	}

	if logic.ExistFeatureGame() {
		example.invalid("桌號(%d) 玩家(%d) 準備切換為 [%v] 模式, 因目前尚有特殊遊戲, 等待 1s 後再執行...", logic.TableNo(), logic.GameUser.UserID, mode)
		time.Sleep(1 * time.Second)
		logic.SwitchProto(mode)
		return nil
	}

	logic.setProto(mode)
	logic.Mode = mode

	example.invalid("桌號(%d) 玩家(%d) 切換為 [%v] 模式, Pattern: %d", logic.TableNo(), logic.GameUser.UserID, mode, logic.Proto.PatternID)

	return logic.Proto
}

// RegisterBeforeSwitchProtoFunc 註冊切換Proto前想做的事
func (logic *SlotLogicBase) RegisterBeforeSwitchProtoFunc(f func()) {
	logic.BeforeSwitchProtoFunc = append(logic.BeforeSwitchProtoFunc, f)
}

// RegisterAfterSwitchProtoFunc 註冊切換Proto後想做的事
func (logic *SlotLogicBase) RegisterAfterSwitchProtoFunc(f func()) {
	logic.AfterSwitchProtoFunc = append(logic.AfterSwitchProtoFunc, f)
}

// CheckModeProto 檢查Proto是否與Mode一致
func (logic *SlotLogicBase) CheckModeProto() {

	if logic.GameUser == nil {
		return
	}

	orgMode := logic.Mode

	if len(logic.bufFeatureGames[DefaultPattern]) > 0 {
		logic.setProto(DefaultPattern)
		logic.Mode = orgMode
		logic.bufFeatureGames[DefaultPattern] = make([]interface{}, 0)
		return
	}

	switch logic.Mode {
	default:
		if logic.Proto != nil && logic.Proto.PatternID == Proto.PatternID {
			return
		}
		logic.setProto(DefaultPattern)
	case ReleasePattern:
		if logic.Proto != nil && logic.Proto.PatternID == ReleaseProto.PatternID {
			return
		}
		now := time.Now()
		isExpired := now.After(logic.ModeExpire)
		if isExpired && logic.ExistFeatureGame() {
			example.invalid("桌號(%d) 玩家(%d) [%v] 模式已經過期但還有殘留獎勵遊戲, 獎勵遊戲結束後, 將切換為 [%v] 模式", logic.TableNo(), logic.GameUser.UserID, ReleasePattern, DefaultPattern)
		}
		logic.setProto(ReleasePattern)
		example.invalid("桌號(%d) 玩家(%d) 切換為 [%v] 模式, Pattern: %d, %s 過期", logic.TableNo(), logic.GameUser.UserID, ReleasePattern, logic.Proto.PatternID, logic.ModeExpire.Format("2006-01-02 15:04:05.000"))
	case BitePattern:
		if logic.Proto != nil && logic.Proto.PatternID == BiteProto.PatternID {
			return
		}
		now := time.Now()
		isExpired := now.After(logic.ModeExpire)
		if isExpired && logic.ExistFeatureGame() {
			example.invalid("桌號(%d) 玩家(%d) [%v] 模式已經過期但還有殘留獎勵遊戲, 獎勵遊戲結束後, 將切換為 [%v] 模式", logic.TableNo(), logic.GameUser.UserID, BitePattern, DefaultPattern)
		}
		logic.setProto(BitePattern)
		example.invalid("桌號(%d) 玩家(%d) 切換為 [%v] 模式, Pattern: %d, %s 過期", logic.TableNo(), logic.GameUser.UserID, BitePattern, logic.Proto.PatternID, logic.ModeExpire.Format("2006-01-02 15:04:05.000"))
	}

}

// CheckModeExpire 檢查特殊模式是否結束
// 已過期 + 沒有 Feature Game 就回復為預設
func (logic *SlotLogicBase) CheckModeExpire() {

	if logic.GameUser == nil {
		return
	}

	if logic.Mode == DefaultPattern {
		return
	}

	now := time.Now()
	isExpired := now.After(logic.ModeExpire)

	// fmt.Println("CheckModeExpire", logic.Mode, logic.ModeExpire, isExpired, logic.ExistFeatureGame())

	if isExpired && !logic.ExistFeatureGame() {
		example.invalid("桌號(%d) 玩家(%d) [%v] 模式已過期 切換為 [%v]模式", logic.TableNo(), logic.GameUser.UserID, logic.Mode, DefaultPattern)
		logic.SwitchProto(DefaultPattern)
		logic.UnassignUserMode(logic.GameUser.UserID)
	}
}

// UpdateState 更新遊戲狀態並傳送給玩家
func (logic *SlotLogicBase) UpdateState(state slot.GameState) {
	logic.GameState = state

	if logic.GameUser != nil {
		logic.SendStateMessageToUser(logic.GameUser)
	}
}

// UserEnter 玩家進入
func (logic *SlotLogicBase) UserEnter(user *gamelogic.GameUser, tableno, seatno int32) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserEnter 發生錯誤, %v\n%s", rec, string(debug.Stack()))
		}
	}()

	logic.bufFeatureGames = make(map[PatternMode][]interface{})

	logic.SlotGameLogicBase.UserEnter(user, tableno, seatno)
	logic.GameStartTime = time.Now()

	logic.GameLog.GameBegin()

	logic.GameLog.AddGameStep("玩家進入")

	logic.GameUser = user

	logic.Clean()

	// 調整狀態為等待
	logic.UpdateState(slot.GameStateWait)

	// 查詢 UserAssignedMode
	logic.ReadUserAssignedMode(user)
}

// UserRecover 玩家斷線重連
func (logic *SlotLogicBase) UserRecover(user *gamelogic.GameUser) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserRecover 發生錯誤, %v\n%s", rec, string(debug.Stack()))
		}
	}()

	logic.bufFeatureGames = make(map[PatternMode][]interface{})

	logic.SlotGameLogicBase.UserRecover(user)
	logic.GameLog.AddGameStep("玩家斷線重連")

	logic.GameUser = user

	logic.Clean()

	// 調整狀態為等待
	logic.UpdateState(slot.GameStateWait)

	// 查詢 UserAssignedMode
	logic.ReadUserAssignedMode(user)
}

// UserLeave 玩家離開
func (logic *SlotLogicBase) UserLeave(user *gamelogic.GameUser) bool {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("UserLeave 發生錯誤, %v\n%s", rec, string(debug.Stack()))
		}
	}()

	logic.bufFeatureGames = make(map[PatternMode][]interface{})

	logic.FreeSpin = nil
	logic.SelectChances = make([]*slot.SelectChance, 0)
	logic.JackpotGameChances = make([]*slot.JackpotGameChance, 0)
	logic.JackpotGame = nil
	logic.PickGame = nil

	logic.GameLog.AddGameStep("玩家離開")
	if logic.Mode != DefaultPattern && logic.Proto != nil && logic.Proto.PatternID != Proto.PatternID {
		logic.SwitchProto(DefaultPattern)
	}

	logic.GameUser = nil

	logic.Clean()

	return logic.SlotGameLogicBase.UserLeave(user)
}

// GameEndLog call GameLog GameEnd()
func (logic *SlotLogicBase) GameEndLog() {
	if logic.GameUser == nil {
		return
	}

	logic.GameUser.BetAmount += logic.GamePayMoney
	logic.GameUser.PayAmount += logic.GameWinMoney
	logic.GameUser.FlowAmount += logic.GamePayMoney

	ud := logic.GameLog.Context().Users[logic.GameUser.UserID]
	ud.TransactionNo = logic.TransactionNo

	logic.GameLog.Context().UserUpdateScore(logic.GameUser.UserID, logic.GameWinMoney-logic.GamePayMoney)
	logic.GameLog.Context().UserBetScore(logic.GameUser.UserID, logic.GamePayMoney)
	logic.GameLog.Context().UserPayoutScore(logic.GameUser.UserID, logic.GameWinMoney)
	logic.GameLog.Context().UserValidScore(logic.GameUser.UserID, logic.GamePayMoney)
	logic.GameLog.Context().UserTaxAmount(logic.GameUser.UserID, 0)

	logic.GameLog.GameEnd()
}

// CalcFlowAmount 計算流水，增加本次異動值(贏的錢-花的錢)，使用禮碼不累計
func (logic *SlotLogicBase) CalcFlowAmount(user *gamelogic.GameUser) {
	// var diff float64
	// diff = float64(logic.GameWinMoney) - float64(logic.GamePayMoney)
	// flowAmt := int64(math.Abs(diff))

	// user.FlowAmount += flowAmt
	// user.BetAmount += logic.GamePayMoney
	// user.PayAmount += logic.GameWinMoney
}

// UpdateUserScore 更新玩家分數
func (logic *SlotLogicBase) UpdateUserScore(user *gamelogic.GameUser) {
	if user.IsRobot() {
		return
	}

	logic.GameLog.AddGameStep("更新玩家金幣, 使用金幣: %d, 贏得獎金: %d", logic.GamePayMoney, logic.GameWinMoney)

	// 單一錢包已先扣過錢，結算不再扣
	// logic.UpdateUserScoreFunc(user.UserID, logic.GameWinMoney-logic.GamePayMoney, 0)

	// logic.UpdateUserScoreFunc(user.UserID, logic.GameWinMoney, 0)
	if err := app.Context.Invoke(func(c apipb.GameResultServiceClient) error {

		req := new(apipb.SettlementReq)
		req.PlayerID = user.UserID
		req.ServerID = logic.TableManager.Room.ServerID
		req.TradeID = user.TradeID
		req.UpdateAmount = logic.GameWinMoney
		req.TransactionNo = logic.TransactionNo
		req.BetAmount = logic.GamePayMoney
		req.ValidAmount = logic.GamePayMoney
		req.PayoutAmount = logic.GameWinMoney
		req.TaxAmount = 0

		res, err := c.SettlementV2(context.Background(), req)
		if err != nil {
			mlog.Error("更新玩家金幣錯誤 SettlemmentV2 Error, %v", err)
			return err
		}

		score := logic.GameWinMoney - logic.GamePayMoney

		// 原先 UpdateUserScoreFunc 內做得要加回來
		user.UpdateScore += score

		// 原先在 roomactionmanager.go 內做得要加回來
		if score > 0 {
			user.WinCount++
			user.UpdateWinCount++
			user.LastWinLose = 1
		} else if score == 0 {
			user.DrawCount++
			user.UpdateDrawCount++
			user.LastWinLose = 0
		} else {
			user.LoseCount++
			user.UpdateLoseCount++
			user.LastWinLose = 0
		}

		user.Score = res.AfterAmount

		return nil
	}); err != nil {
		mlog.Error("更新玩家金幣錯誤 %v", err)
		logic.KickUserFunc(user.UserID, "System Error ...")
	}
}

// UpdateUserScoreWithHistory 先新增History再更新玩家分數
func (logic *SlotLogicBase) UpdateUserScoreWithHistory(user *gamelogic.GameUser, gameType int32, logType int32, gameData string) {
	if user.IsRobot() {
		return
	}

	logic.AddSlotGameHistoryLog(user, gameType, logType, gameData)

	logic.GameLog.AddGameStep("更新玩家金幣, 使用金幣: %d, 贏得獎金: %d", logic.GamePayMoney, logic.GameWinMoney)

	logic.UpdateUserScoreFunc(user.UserID, logic.GameWinMoney-logic.GamePayMoney, 0)
}
