package slot

import (
	"sort"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// PaylineAllMatcher 全線賠付比對器
type PaylineAllMatcher struct {
	*BasicPayLineMatcher
}

// NewPaylineAllMatcher 建立全線賠付比對器，傳入所有SymbolAttributes讓Matcher抓出百搭設定
func NewPaylineAllMatcher(attrs []*protodata.SymbolAttributes) *PaylineAllMatcher {

	matcher := new(PaylineAllMatcher)
	matcher.BasicPayLineMatcher = NewBasicPayLineMatcher(attrs)

	matcher.BasicPayLineMatcher.MatchLeft = matcher.MatchLeft
	matcher.BasicPayLineMatcher.MatchRight = matcher.MatchRight

	return matcher
}

// MatchLeft 比對GameResult是否符合PayTable規則
func (matcher *PaylineAllMatcher) MatchLeft(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {

	result.ReelMode()

	paylineres := new(PayResult)
	paylineres.SymbolName = rule.GetSymbol()
	paylineres.Multi = 1

	reelCnts := make([]int, 0, result.Reel)
	nreelCnts := make([]int, 0, result.Reel)
	wreelCnts := make([]int, 0, result.Reel)

	for reelNum, reel := range result.GameResult {
		var repeat int
		var nrepeat int
		var wrepeat int

		for rowNum, symbol := range reel {

			if symbol == nil {
				continue
			}

			// 檢查符號相同
			if symbol.GetSymbol() == paylineres.SymbolName {
				repeat++
				nrepeat++
				paylineres.SymbolIDs = append(paylineres.SymbolIDs, int32(symbol.GetId()))
				// paylineres.BlinkPositions = append(paylineres.BlinkPositions, Position{Reel: int32(reelNum + 1), Row: int32(rowNum + 1)})
				paylineres.AddBlinkPosition(Position{Reel: int32(reelNum + 1), Row: int32(rowNum + 1)})
			}

			// 是否百搭
			for _, wild := range matcher.Wilds {
				if wild.GetSymbol() != symbol.GetSymbol() {
					continue
				}
				// 是否是例外
				isEx := false
				exts := strings.Split(wild.GetWildEx(), ",")
				for _, ext := range exts {
					if ext == paylineres.SymbolName {
						isEx = true
					}
				}
				// 不是例外則可搭
				if !isEx {
					repeat++
					if wild.GetIsWild() == 1 {
						nrepeat++
					} else {
						wrepeat++
					}
					paylineres.HasWild = true
					paylineres.SymbolIDs = append(paylineres.SymbolIDs, int32(symbol.GetId()))
					// paylineres.BlinkPositions = append(paylineres.BlinkPositions, Position{Reel: int32(reelNum + 1), Row: int32(rowNum + 1)})
					paylineres.AddBlinkPosition(Position{Reel: int32(reelNum + 1), Row: int32(rowNum + 1)})
				}
			}
		}

		if repeat+wrepeat == 0 {
			break
		}

		paylineres.PayCount++
		reelCnts = append(reelCnts, repeat)
		nreelCnts = append(nreelCnts, nrepeat)
		wreelCnts = append(wreelCnts, wrepeat)
	}

	if paylineres.PayCount != rule.GetCount() {
		return false, paylineres
	}

	// 檢查每個位置要有一個此規則的符號
	var hasRuleSymbol bool
	for _, pos := range paylineres.BlinkPositions {
		if rule.GetSymbol() == result.GameResult[pos.Reel-1][pos.Row-1].GetSymbol() {
			hasRuleSymbol = true
			break
		}
	}

	if !hasRuleSymbol {
		return false, paylineres
	}

	paylineres.IsMatch = true
	paylineres.WayCount = 1
	for _, cnt := range reelCnts {
		paylineres.WayCount *= int32(cnt)
	}
	paylineres.NormalWayCount = 1
	for _, cnt := range nreelCnts {
		paylineres.NormalWayCount *= int32(cnt)
	}
	paylineres.WildWayCount = paylineres.WayCount - paylineres.NormalWayCount
	paylineres.Rule = rule

	sort.Slice(paylineres.SymbolIDs, func(i, j int) bool {
		return paylineres.SymbolIDs[i] < paylineres.SymbolIDs[j]
	})

	return true, paylineres
}

// MatchRight 比對GameResult是否符合PayTable規則
func (matcher *PaylineAllMatcher) MatchRight(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {
	// 最後要轉回來
	defer func() {
		result = result.Reverse()
	}()
	result.ReelMode()
	// 轉過去
	result = result.Reverse()
	return matcher.MatchLeft(result, rule)
}
