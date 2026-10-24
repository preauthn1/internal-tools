package slot

// GameMode 遊戲模式
type GameMode string

// GameMode
const (
	GameModeMain        GameMode = "main"
	GameModeFreeSpin    GameMode = "free_spins"
	GameModeLightning   GameMode = "lightning"
	GameModeBonus       GameMode = "bonus"
	GameModePickGame    GameMode = "pickgame"
	GameModeJackpotGame GameMode = "jackpot_game"
	GameModeSelection   GameMode = "selection"
)

// Feature 特色功能
const (
	FeatureSpin uint = 1 << iota
	FeatureExpand
	FeatureBG
	FeatureFreeSpin
	FeatureBonus
	FeatureJackpotGame
	FeatureJackpot
	FeaturePickGame
	FeatureLightning
	FeatureSelection
)

// PayType 獎勵方式
type PayType string

// PayTypes
const (
	PayTypeChip            PayType = "chip"
	PayTypeChipXRate       PayType = "chipXrate"
	PayTypePickGame        PayType = "pickgame"
	PayTypeJackpot         PayType = "jackpot"
	PayTypeRespin          PayType = "respin"
	PayTypeFreeSpin        PayType = "free_spins"
	PayTypeMysteryJackpot  PayType = "mysteryjackpot"
	PayTypeMysteryJackpotR PayType = "mysteryjackpotR"
	PayTypeMysteryGrand    PayType = "mysterygrand"
	PayTypeMysteryGrandR   PayType = "mysterygrandR"
	PayTypeMysteryJP       PayType = "mysteryJP"
	PayTypeMysteryJPR      PayType = "mysteryJPR"
	PayTypeMysteryJPG      PayType = "mysteryJPG"
	PayTypeMysteryJPGR     PayType = "mysteryJPGR"
	PayTypeSelectGameMode  PayType = "selectgamemode"
	PayTypeLightning       PayType = "lightning"
	PayTypeSelection       PayType = "selection"
	PayTypeAddTimes        PayType = "addfreespins"
	PayTypeAccuMulate      PayType = "accumulate"
	PayTypeBonus           PayType = "bonus"
	PayTypeLock            PayType = "lock"
)

func (pt PayType) String() string {
	return string(pt)
}

// PayLineMode 賠付線方式
type PayLineMode string

// PayLineModes
const (
	PayLineModeAll   PayLineMode = "all"
	PayLineModeFixed PayLineMode = "fixed"
)

func (plm PayLineMode) String() string {
	return string(plm)
}

// Int get PayLineMode number
func (plm PayLineMode) Int() int32 {
	switch plm {
	default:
		return -1
	case PayLineModeAll:
		return 1
	case PayLineModeFixed:
		return 0
	}
}

// GameState 遊戲狀態
type GameState int32

// 階段狀態
const (
	GameStateBet         = iota // 下注狀態
	GameStateAward              // 開獎
	GameStateAwardAction        // 開獎動畫
	GameStateEndAction          // 結算動畫
	GameStateWait               // 等待遊戲開始
	GameStateUserAction         // 等待玩家動作(各遊戲可能有不同的事情要做)
)

func (s GameState) String() string {
	switch s {
	default:
		return "none"
	case GameStateBet:
		return "Bet"
	case GameStateAward:
		return "Award"
	case GameStateAwardAction:
		return "AwardAction"
	case GameStateEndAction:
		return "EndAction"
	case GameStateWait:
		return "Wait"
	case GameStateUserAction:
		return "UserAction"
	}
}

// BetStates
const (
	BetStateSuccess        int32 = 0  // 下注成功
	BetStateMoneyNotEnough int32 = -1 // 玩家金額不足
	BetStateBetTypeError   int32 = -2 // 下注類型未找到
	BetStateIsBanker       int32 = -3 // 莊家不能下注
	BetStateBetTypeLimit   int32 = -4 // 此類下注已達到上限
	BetStateStepError      int32 = -5 // 非下注階段
	BetStateMoneyLessMin   int32 = -6 // 玩家金額小於下注限制
	BetStateIllegal        int32 = -7 // 非法金額
)

// 賠率計算類型
const (
	ConditionRate   = "betRate"
	ConditionAmount = "betAmount"
)

// LinkType 賠線類型
type LinkType string

// LinkTypes
const (
	LinkTypeLeft     LinkType = "left"
	LinkTypeRight    LinkType = "right"
	LinkTypeBoth     LinkType = "both"
	LinkTypeAny      LinkType = "any"
	LinkTypeStraight LinkType = "straight"
	LinkTypeAdjacent LinkType = "adjacent"
)

func (t LinkType) String() string {
	return string(t)
}

// 歷程紀錄類型
const (
	HistoryTypeMAIN      = 1 // 主遊戲
	HistoryTypeBONUS     = 2 // BonusGame
	HistoryTypeFREE      = 3 // FreeGame
	HistoryTypeJPNOLINK  = 4 // 單機 JP
	HistoryTypeJPLINK    = 5 // 連線 JP
	HistoryTypeRESPIN    = 6 // 重轉
	HistoryTypeLIGHTNING = 7 // Lightning結算
	HistoryTypeCoin      = 8 // 金幣統計結算
)

// 遊戲歷程紀錄類型
const (
	LogTypeREEL      = 1 // 輪軸資訊
	LogTypeOPEN      = 2 // 翻圖案，三個一樣結束型
	LogTypeTOTALCOIN = 3 // 累計出現金幣統計
	LogTypeBG        = 4 // 背景
)

// UseType 應用情境
type UseType string

// UseType
const (
	UseTypeALL      UseType = "all"
	UseTypeMainGame UseType = "normal"
	UseTypeFreeGame UseType = "freespin"
	UseTypeRespin   UseType = "respin"
)

// ExtraRuleCondition 額外規則條件
type ExtraRuleCondition string

// ExtraRuleCondition
const (
	ExtraRuleConditionAfterSpin       ExtraRuleCondition = "AfterSpin"
	ExtraRuleConditionAfterBackground ExtraRuleCondition = "AfterBackground"
	ExtraRuleConditionGetAward        ExtraRuleCondition = "GetAward"
	ExtraRuleConditionAppearSymbol    ExtraRuleCondition = "AppearSymbol"
	ExtraRuleConditionJPDefense       ExtraRuleCondition = "JPDefense"
)

// ExtraRuleAction 額外規則動作
type ExtraRuleAction string

// ExtraRuleAction
const (
	ExtraRuleActionRespin            ExtraRuleAction = "Respin"
	ExtraRuleActionRespinR           ExtraRuleAction = "RespinR"
	ExtraRuleActionConditionalAction ExtraRuleAction = "ConditionalAction"
	ExtraRuleActionUniGoldBall       ExtraRuleAction = "UniGoldBall"
	ExtraRuleActionExpandReel        ExtraRuleAction = "ExpandReel"
	ExtraRuleActionSwitchReelRandom  ExtraRuleAction = "SwitchReelRandom"
	ExtraRuleActionReplaceBG         ExtraRuleAction = "ReplaceBG"
	ExtraRuleActionReplaceBGT        ExtraRuleAction = "ReplaceBGT"
	ExtraRuleActionReplaceBGR        ExtraRuleAction = "ReplaceBGR"
	ExtraRuleActionReplaceBGTR       ExtraRuleAction = "ReplaceBGTR"
)

// MultiMode 加倍模式
type MultiMode int

// MultiModes
const (
	MultiModeNone              MultiMode = iota // 沒有加倍
	MultiModeDiceOnce                           // 有出現Wild骰一次
	MultiModeDiceEachReelSum                    // 每一個Reel的Wild都骰一次, 所有Reel相加
	MultiModeDiceEachReelMulti                  // 每一個Reel的Wild都骰一次, 所有Reel相乘
	MultiModeLevel                              // 依照級別加倍
	MultiModeEachWildMulti                      // 每個Wild本身帶有倍數
	MultiModeEachWildSum                        // 每個Wild本身帶有倍數
)

// MultiType 加倍類型
type MultiType int

// MultiTypes
const (
	MultiTypeNone           MultiType = iota // 沒有加倍
	MultiTypeApearWild                       // 有出現Wild全部加倍
	MultiTypeLinkWild                        // 連到Wild才加倍
	MultiTypeLinkWildTarget                  // 連到Wild才加倍被連到的Wild
)
