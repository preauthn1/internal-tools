package baselogic

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// CalcMultiMode 計算加倍數
func (logic *SlotLogicBase) CalcMultiMode(mode slot.MultiMode, str string) (multi int64, detail map[slot.Position]int32) {

	multi = 1
	detail = make(map[slot.Position]int32)

	// 先檢查盤面是否出現百搭符號，決定加倍數量
	type MultiOpt struct {
		Value  int64
		Weight int
	}
	var opts []MultiOpt

	// 決定倍率選項
	{
		// fmt.Println(str)
		multis := strings.Split(str, ",")
		for _, row := range multis {
			vals := strings.Split(row, ":")

			m, _ := strconv.ParseFloat(vals[0], 64)
			w, _ := strconv.ParseFloat(vals[1], 64)

			multi := MultiOpt{Value: int64(m), Weight: int(w)}
			opts = append(opts, multi)
		}

		if len(opts) > 0 {
			// 排序由大到小
			sort.Slice(opts, func(i, j int) bool {
				return opts[i].Value > opts[j].Value
			})
		}
	}

	switch mode {
	default:
	case slot.MultiModeDiceOnce:

		var multis []int
		for _, opt := range opts {
			multis = append(multis, int(opt.Value))
		}

		if len(multis) == 0 {
			return
		}

		sort.Slice(multis, func(i, j int) bool {
			return multis[i] < multis[j]
		})

		r := new(prob.Random)
		for _, opt := range opts {
			r.Add(opt, opt.Weight)
		}
		t := r.Choose().(MultiOpt)

		logic.GameLog.AddGameStep("免費遊戲百搭加倍選項 %v 倍, 骰出 %v 倍", multis, t.Value)

		multi = t.Value

	case slot.MultiModeDiceEachReelSum:

	case slot.MultiModeDiceEachReelMulti:

		// 判斷每個Reel的加倍數量
		logic.GameResultManager.ReelMode()

		var multis []string
		for _, reel := range logic.GameResultManager.GameResult {
			var cnt int

			for _, symbol := range reel {
				// 檢查是否為百搭
				for _, attr := range logic.Proto.SymbolAttributes {
					if (attr.GetIsWild() == 1 || attr.GetIsWild() == 2) && symbol.GetSymbol() == attr.GetSymbol() {
						cnt++
					}
				}
			}

			// 若整條Reel都是百搭則使用最大倍數
			if cnt == int(logic.GameResultManager.Row) {
				multi = multi * opts[0].Value
				multis = append(multis, strconv.Itoa(int(opts[0].Value)))

				for _, symbol := range reel {
					for _, attr := range logic.Proto.SymbolAttributes {
						if attr.GetIsWild() == 1 && symbol.GetSymbol() == attr.GetSymbol() {
							logic.GameResultManager.MultiSymbol[int32(symbol.GetId())] = opts[0].Value
						}
					}
				}
			} else if cnt > 0 {
				r := new(prob.Random)
				for _, opt := range opts {
					r.Add(opt, opt.Weight)
				}
				t := r.Choose().(MultiOpt)

				multi = multi * t.Value
				multis = append(multis, strconv.Itoa(int(t.Value)))

				for _, symbol := range reel {
					for _, attr := range logic.Proto.SymbolAttributes {
						if attr.GetIsWild() == 1 && symbol.GetSymbol() == attr.GetSymbol() {
							logic.GameResultManager.MultiSymbol[int32(symbol.GetId())] = t.Value
						}
					}
				}
			} else {
				multis = append(multis, "1")
			}
		}
		var _multis []int
		for _, opt := range opts {
			_multis = append(_multis, int(opt.Value))
		}
		sort.Slice(_multis, func(i, j int) bool {
			return _multis[i] < _multis[j]
		})

		logic.GameLog.AddGameStep("百搭加倍選項 %v 倍, 總加倍: %s = %d 倍", _multis, strings.Join(multis, " * "), multi)

	case slot.MultiModeEachWildSum, slot.MultiModeEachWildMulti:

		// 亂數器
		r := new(prob.Random)
		for _, opt := range opts {
			r.Add(opt, opt.Weight)
		}

		// 遇到Wild就產生一次加倍選項
		logic.GameResultManager.ReelMode()
		for reelNum, reel := range logic.GameResultManager.GameResult {
			for rowNum, symbol := range reel {
				if logic.IsWild(symbol.GetSymbol()) {
					t := r.Choose().(MultiOpt)
					detail[slot.Position{Reel: int32(reelNum + 1), Row: int32(rowNum + 1)}] = int32(t.Value)
				}
			}
		}

		if len(detail) == 0 {
			return
		}

		var descs []string

		for pos, m := range detail {
			descs = append(descs, fmt.Sprintf("[%v](x%d)", pos, m))
		}

		logic.GameLog.AddGameStep("百搭加倍選項: %s", strings.Join(descs, ", "))
	}

	return
}

// CalcMultiLevel 計算階段型加倍
// @Param str pattern
// @Param winTimes 第n次贏
func (logic *SlotLogicBase) CalcMultiLevel(str string, winTimes int) (multi int64) {

	if winTimes <= 0 {
		return 1
	}

	strs := strings.Split(str, ",")

	var multis []int64

	for _, s := range strs {
		val, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return 1
		}
		multis = append(multis, int64(val))
	}

	if winTimes > len(multis) {
		winTimes = len(multis)
	}

	return multis[winTimes-1]
}

// IsWild 判斷符號是否為百搭
func (logic *SlotLogicBase) IsWild(s string) bool {

	for _, attr := range logic.Proto.SymbolAttributes {
		if attr.GetSymbol() != s {
			continue
		}
		switch attr.GetIsWild() {
		case 1, 2:
			return true
		}
	}

	return false
}
