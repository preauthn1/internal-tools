package slot

import (
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// JackpotLevel define jackpot level
type JackpotLevel string

// Jackpot Levels
var (
	JackpotSingle   JackpotLevel = "Single"
	JackpotGrand    JackpotLevel = "Grand"
	JackpotMajor    JackpotLevel = "Major"
	JackpotMinor    JackpotLevel = "Minor"
	JackpotMini     JackpotLevel = "Mini"
	JackpotGoldReel JackpotLevel = "GoldReel"
)

func (jplevel JackpotLevel) String() string {
	return string(jplevel)
}

// Int get jackpot level number
func (jplevel JackpotLevel) Int() int32 {
	switch jplevel {
	default:
		return 0
	case JackpotGrand:
		return 1
	case JackpotMajor:
		return 2
	case JackpotMinor:
		return 3
	case JackpotMini:
		return 4
	case JackpotGoldReel:
		return 1001
	}
}

// NewJackpotLevel create JackpotLevel from int32
func NewJackpotLevel(level int32) JackpotLevel {
	switch level {
	case 1:
		return JackpotGrand
	case 2:
		return JackpotMajor
	case 3:
		return JackpotMinor
	case 4:
		return JackpotMini
	case 1001:
		return JackpotGoldReel
	default:
		return JackpotSingle
	}
}

// Jackpot 獎池
type Jackpot struct {
	Origin *protodata.JackpotTable

	// 本次貢獻值
	Contribution int64
	// 累積貢獻值
	TotalContribution int64
	// 累積出獎金額
	AwardTotal int64
	// 獎池等級
	Level JackpotLevel
}

// Reset 重置獎池金額
func (jp *Jackpot) Reset() {
	jp.Contribution = 0
	jp.TotalContribution = 0
	jp.AwardTotal = 0
}

// GetCurrentPool 取得當前獎池金額，以最大出獎金額計算
// 每次Spin貢獻值(contribution)=increment*toltalBet
// 現在水池實際金額=startUp+每次Spin貢獻值+已儲存累積的貢獻值-重置前累積JP出獎金額
func (jp *Jackpot) GetCurrentPool() int64 {

	inc := (jp.TotalContribution / 10000) - jp.AwardTotal

	if inc < 0 {
		return jp.Origin.GetStartUp()
	}

	award := jp.Origin.GetStartUp() + inc

	if jp.Origin.GetIsCeling() == 1 && int64(jp.Origin.GetCelingValue()) < award {
		return int64(jp.Origin.GetCelingValue())
	}

	return award
}

// GetCurrentPoolByBet 取得當前獎池金額，隨玩家當前下注金額變動
// 玩家實際看到的現在水池實際金額*CurrentBet/MAXBet
// 水池會根據變更投注額即時變動
// 如果是連線獎池回傳最大金額, 非連線獎池依照下注比例給獎金
func (jp *Jackpot) GetCurrentPoolByBet(currentBet, maxBet int64) int64 {
	if jp.Origin.GetIsLinked() == 1 {
		return jp.GetCurrentPool()
	}
	return jp.GetCurrentPool() * currentBet / maxBet
}

// Clean 清除獎池累積資訊
func (jp *Jackpot) Clean() {
	jp.Contribution = 0
	jp.TotalContribution = 0
	jp.AwardTotal = 0
}

// Increase 貢獻
//
// 若無貢獻則不可增加
//
// 若有上限，上限金額為當前可領取金額
func (jp *Jackpot) Increase(money int64) {
	if jp.Origin.GetIsProgress() != 1 {
		return
	}

	jp.TotalContribution += jp.GetIncrease(money)
}

// GetIncrease 取得本次貢獻值
func (jp *Jackpot) GetIncrease(money int64) int64 {
	if jp.Origin.GetIsProgress() != 1 {
		return 0
	}

	return int64(jp.Origin.GetIncrement()) * money
}

// Award 中獎
func (jp *Jackpot) Award(currentBet, maxBet int64) int64 {

	award := jp.GetCurrentPoolByBet(currentBet, maxBet)

	jp.AwardTotal += award

	return award
}

// IsNeedReset 檢查是否需要重置
func (jp *Jackpot) IsNeedReset() bool {
	inc := (jp.TotalContribution / 10000) - jp.AwardTotal
	return inc < 0
}
