package slot

import (
	"bytes"
	"fmt"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// GameResult 遊戲結果
//
// 預設排法:
//     [
//         [ Symbole, Symbole, Symbole, Symbole, Symbole ],	// Row
//         [ Symbole, Symbole, Symbole, Symbole, Symbole ],	// Row
//         [ Symbole, Symbole, Symbole, Symbole, Symbole ],	// Row
//     ]
type GameResult [][]*protodata.SlotSymbolDistribution

// Transfor 翻轉排列方式
//
//     [
//         [ Symbole, Symbole, Symbole ],	// Reel
//         [ Symbole, Symbole, Symbole ],	// Reel
//         [ Symbole, Symbole, Symbole ],	// Reel
//         [ Symbole, Symbole, Symbole ],	// Reel
//         [ Symbole, Symbole, Symbole ],	// Reel
//     ]
func (r GameResult) Transfor() GameResult {
	if len(r) == 0 {
		return r
	}

	new := GameResult{}

	for i := 0; i < len(r[0]); i++ {
		symbols := make([]*protodata.SlotSymbolDistribution, 0, len(r))
		for j := 0; j < len(r); j++ {
			symbols = append(symbols, r[j][i])
		}
		new = append(new, symbols)
	}

	return new
}

// Reverse 將前後顛倒，請確保當前模式
func (r GameResult) Reverse() GameResult {
	if len(r) == 0 {
		return r
	}

	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return r
}

func (r GameResult) String() string {

	buf := bytes.NewBuffer(nil)

	for _, row := range r {

		format := make([]string, 0, len(row))
		rowsymbols := make([]interface{}, 0, len(row))

		for _, symbol := range row {
			format = append(format, "%5s")
			rowsymbols = append(rowsymbols, symbol.GetSymbol())
		}

		fmt.Fprintf(buf, strings.Join(format, " "), rowsymbols...)
		fmt.Fprintln(buf)
	}

	return buf.String()
}

// GetSymbolIDByRow 取得指定Row的符號ID
func (r GameResult) GetSymbolIDByRow(rowID uint32) []int32 {

	var ans []int32

	if len(r) < int(rowID) || rowID == 0 {
		return ans
	}

	for _, row := range r[rowID-1] {
		ans = append(ans, int32(row.GetId()))
	}

	return ans
}

// GetReelSymbol 取得map[Reel]map[Row]模式資料
func (r GameResult) GetReelSymbol() map[uint32]map[uint32][]*protodata.SlotSymbolDistribution {

	ans := make(map[uint32]map[uint32][]*protodata.SlotSymbolDistribution)

	for i := 0; i < len(r); i++ {
		for j := 0; j < len(r[i]); j++ {
			if ans[uint32(j+1)] == nil {
				ans[uint32(j+1)] = make(map[uint32][]*protodata.SlotSymbolDistribution)
			}

			ans[uint32(j+1)][uint32(i+1)] = []*protodata.SlotSymbolDistribution{r[i][j]}
		}
	}

	return ans
}

// GetSymbolString 取得符號二維陣列
func (r GameResult) GetSymbolString() [][]string {

	var strs [][]string

	for _, row := range r {
		var symbols []string
		for _, symbol := range row {
			symbols = append(symbols, symbol.GetSymbol())
		}
		strs = append(strs, symbols)
	}

	return strs
}

// ToSymbolString 取得 p1,p2,p3|p2,p3,p4|p3,p4,p5|p4,p5,p6|p5,p6,p7 字樣
func (r GameResult) ToSymbolString() string {

	strs1 := r.GetSymbolString()
	strs2 := make([]string, 0, len(strs1))

	for _, row := range strs1 {
		strs2 = append(strs2, strings.Join(row, ","))
	}

	return strings.Join(strs2, "|")
}
