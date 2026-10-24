package slot

import (
	"reflect"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// PayResult 比對賠付結果
type PayResult struct {
	IsMatch        bool                    // 是否比對成功
	HasWild        bool                    // 是否有使用百搭
	FromBackground bool                    // 是否從背景盤面取得的獎勵
	WayCount       int32                   // 連線數量
	NormalWayCount int32                   // 一般連線數量
	WildWayCount   int32                   // 碰到Wild連線數量
	PayCount       int32                   // 共幾個Reel
	Multi          int32                   // 加乘倍數
	SymbolName     string                  // 符號
	SymbolIDs      []int32                 // 中獎的符號ID
	BlinkPositions []Position              // 中獎符號位置
	BlinkDetails   [][]Position            // 每個連線的位置
	SymbolPayLines []*SymbolPayLine        // 中獎的連線與符號
	Rule           *protodata.SlotPayTable // 中獎的規則
	WinMoney       int64                   // 本規則贏的金額
	MultiMode      MultiMode               // 加倍類型
	MultiType      MultiType               // 加倍類型
	MultiDetail    map[Position]int32      // 位置加倍資訊
}

// SymbolPayLine 中獎的連線與符號
type SymbolPayLine struct {
	HasWild        bool     // 是否有使用百搭
	SymbolIDs      []int32  // 中獎的符號ID
	PayLine        *PayLine // 中獎的連線
	SymbolNames    []string // 中獎的符號
	MultipleStep   []int32  // 加乘階段 w2 w3 = > [2, 3]
	BlinkPositions []Position
	Win            int64
}

func (line *SymbolPayLine) String() string {
	return line.PayLine.String()
}

// GetMultipleValue 取得加乘倍數
func (line *SymbolPayLine) GetMultipleValue(mode MultiMode) int32 {

	switch mode {
	case MultiModeEachWildSum:
		var multiVal int32
		for _, multi := range line.MultipleStep {
			multiVal += multi
		}
		if multiVal == 0 {
			multiVal = 1
		}
		return multiVal
	case MultiModeEachWildMulti:
		var multiVal int32 = 1
		for _, multi := range line.MultipleStep {
			multiVal *= multi
		}
		return multiVal
	}

	return 1
}

// GetPayType 取得中獎項目
func (result *PayResult) GetPayType() PayType {
	return PayType(result.Rule.GetPayType())
}

// CalcWinMoney 計算贏得的金額
func (result *PayResult) CalcWinMoney(betMoney, betRate int64) (base, win int64) {
	if !result.IsMatch || result.Rule == nil {
		return
	}

	// 計算賠付項
	switch result.GetPayType() {
	default:
		return
	case PayTypeChip:
		f, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)
		base = betMoney * int64(f)
	case PayTypeChipXRate:
		f, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)
		base = betMoney * betRate * int64(f)
	}

	// 計算加倍項
	switch result.MultiType {
	default:
		if len(result.SymbolPayLines) > 0 {
			for _, line := range result.SymbolPayLines {
				win += base * int64(line.GetMultipleValue(result.MultiMode))
			}
		} else {
			win = base * int64(result.WayCount)
		}
	case MultiTypeApearWild:
		if len(result.SymbolPayLines) > 0 {
			for _, line := range result.SymbolPayLines {
				win += base * int64(line.GetMultipleValue(result.MultiMode))
			}
		} else {
			win = base * int64(result.WayCount) * int64(result.Multi)
		}
	case MultiTypeLinkWild:
		if len(result.SymbolPayLines) > 0 {
			for _, line := range result.SymbolPayLines {
				win += base * int64(line.GetMultipleValue(result.MultiMode))
			}
		} else {
			win = base*int64(result.NormalWayCount) + base*int64(result.WildWayCount)*int64(result.Multi)
		}
	case MultiTypeLinkWildTarget:
		if len(result.SymbolPayLines) > 0 {
			return
		}
		result.Multi = 1
		var multipos []Position
		for _, pos := range result.BlinkPositions {
			if m, ok := result.MultiDetail[pos]; ok {
				result.Multi = result.Multi * m
				multipos = append(multipos, pos)
			}
		}

		for _, poies := range result.BlinkDetails {
			l := new(PayLine)
			l.Positions = poies

			sl := new(SymbolPayLine)
			sl.PayLine = l
			for _, pos := range poies {
				if m, ok := result.MultiDetail[pos]; ok {
					sl.MultipleStep = append(sl.MultipleStep, m)
					sl.HasWild = true
				}
			}
			result.SymbolPayLines = append(result.SymbolPayLines, sl)
		}

		if len(result.SymbolPayLines) > 0 {
			for _, line := range result.SymbolPayLines {
				_win := base * int64(line.GetMultipleValue(result.MultiMode))
				win += _win
				line.Win = _win
			}
		} else {
			win = base*int64(result.NormalWayCount) + base*int64(result.WildWayCount)*int64(result.Multi)
		}

	}
	return
}

// IsWinBonusGame 是否贏得BonusGame
func (result *PayResult) IsWinBonusGame() bool {
	if !result.IsMatch || result.Rule == nil {
		return false
	}

	switch result.Rule.GetPayType() {
	default:
		return false
	case PayTypePickGame.String():
		return true
	}
}

// IsWinJackpot 是否贏得Jackpot
func (result *PayResult) IsWinJackpot() bool {
	if !result.IsMatch || result.Rule == nil {
		return false
	}

	switch result.Rule.GetPayType() {
	default:
		return false
	case PayTypeJackpot.String():
		return true
	}
}

// RemoveLine 移除既有的線
func (result *PayResult) RemoveLine(lineID uint32) {
	if !result.IsMatch || result.Rule == nil {
		return
	}

	if len(result.SymbolPayLines) == 0 {
		return
	}

	var tmp []*SymbolPayLine
	for _, line := range result.SymbolPayLines {
		if lineID == line.PayLine.Origin.GetLineId() {
			continue
		}
		tmp = append(tmp, line)
	}

	result.SymbolPayLines = tmp
	result.WayCount = int32(len(result.SymbolPayLines))
}

// AddBlinkPosition add blink symbol position
func (result *PayResult) AddBlinkPosition(pos Position) {
	var exsit bool
	for _, p := range result.BlinkPositions {
		if reflect.DeepEqual(p, pos) {
			exsit = true
		}
	}
	if !exsit {
		result.BlinkPositions = append(result.BlinkPositions, pos)
	}

	if len(result.BlinkDetails) == 0 {
		result.BlinkDetails = append(result.BlinkDetails, []Position{pos})
	} else {

		var tmp [][]Position

		for i := range result.BlinkDetails {
			if len(result.BlinkDetails[i]) == int(pos.Reel)-1 {
				result.BlinkDetails[i] = append(result.BlinkDetails[i], pos)
				continue
			}
			tmp = append(tmp, result.BlinkDetails[i])
		}

		for i := range tmp {
			row := append([]Position(nil), tmp[i]...)
			row[len(row)-1] = pos

			var _exsit bool
			for _, line := range result.BlinkDetails {
				if reflect.DeepEqual(line, row) {
					_exsit = true
				}
			}
			if !_exsit {
				result.BlinkDetails = append(result.BlinkDetails, row)
			}
		}

	}

}
