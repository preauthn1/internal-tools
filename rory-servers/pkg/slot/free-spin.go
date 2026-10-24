package slot

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// FreeSpinExtAward 免費遊戲額外獎勵
type FreeSpinExtAward struct {
	TotalCount     int32                                          // 累計贏取次數
	TotalWinMoney  int64                                          // 累計贏取金額
	Count          int32                                          // 本次贏取次數
	WinMoney       int64                                          // 本次贏取金額
	EachCountMoney int64                                          // 單次金額
	IsFinished     bool                                           // 獎勵是否結算
	Symbols        map[Position]*protodata.SlotSymbolDistribution `json:"-"` // 獎勵的符號
}

// FreeSpin 免費遊戲
type FreeSpin struct {
	Origin        *protodata.FreeSpinTable            // 原始FreeSpinTable
	SymbolGroupID int32                               // 使用的SymbolGroup
	BgGroupID     int32                               // 使用的BackGroundGroup
	Symbols       []*protodata.SlotSymbolDistribution // 此次FreeSpin會用到的Symbols
	WinMultiple   []int32                             // FreeSpin中獎後會加乘
	Count         int                                 // 剩餘次數
	WinCount      int                                 // 累計中獎次數
	WinMoney      int64                               // 累計中獎金額
	BetMoney      int64                               // 下注金額
	BetRate       int64                               // 下注倍率

	Used int // 使用次數

	FreeSpinExtAward *FreeSpinExtAward // 額外獎勵
}

// NewFreeSpin 建立免費遊戲
func NewFreeSpin(origin *protodata.FreeSpinTable, symbols []*protodata.SlotSymbolDistribution, betMoney int64, betRate int64) *FreeSpin {

	freeSpin := new(FreeSpin)

	freeSpin.Origin = origin
	freeSpin.SymbolGroupID = int32(origin.GetFreeSymbolGroupId())
	freeSpin.BgGroupID = int32(origin.GetGoldReelGroupId())

	for _, symbol := range symbols {
		if symbol.GetGroupId() != origin.GetFreeSymbolGroupId() {
			continue
		}

		freeSpin.Symbols = append(freeSpin.Symbols, symbol)
	}

	freeSpin.Count = int(origin.GetFreeSpinCount())
	freeSpin.BetMoney = betMoney
	freeSpin.BetRate = betRate
	// freeSpin.IsCoupon = isCoupon

	if origin.GetMultiplierType() == "level" {

		multis := strings.Split(origin.GetMultiplierValue(), ",")

		for _, multi := range multis {
			f, err := strconv.ParseFloat(multi, 64)
			if err != nil {
				continue
			}
			freeSpin.WinMultiple = append(freeSpin.WinMultiple, int32(f))
		}
	}

	return freeSpin
}

// GetMultipleValue 取得目前中獎加乘倍率
func (freeSpin *FreeSpin) GetMultipleValue() int32 {

	length := len(freeSpin.WinMultiple)

	if length == 0 {
		return 1
	}

	if freeSpin.WinCount >= length {
		return freeSpin.WinMultiple[length-1]
	}

	return freeSpin.WinMultiple[freeSpin.WinCount]
}
