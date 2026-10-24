package slot

import (
	"sort"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// PaylineFixedMatcher 指定線賠付比對器
type PaylineFixedMatcher struct {
	*BasicPayLineMatcher
	NormalLine  []*PayLine
	ReverseLine []*PayLine
}

// NewPaylineFixedMatcher 建立指定線線賠付比對器，傳入所有SymbolAttributes讓Matcher抓出百搭設定
func NewPaylineFixedMatcher(attrs []*protodata.SymbolAttributes, lines []*protodata.SlotPayLine) *PaylineFixedMatcher {

	matcher := new(PaylineFixedMatcher)
	matcher.BasicPayLineMatcher = NewBasicPayLineMatcher(attrs)

	matcher.NormalLine = make([]*PayLine, 0, len(lines))
	matcher.ReverseLine = make([]*PayLine, 0, len(lines))

	for _, line := range lines {
		matcher.NormalLine = append(matcher.NormalLine, NewPayLine(line))
		matcher.ReverseLine = append(matcher.ReverseLine, NewPayLine(line).Reverse())
	}

	matcher.BasicPayLineMatcher.MatchLeft = matcher.MatchLeft
	matcher.BasicPayLineMatcher.MatchRight = matcher.MatchRight

	return matcher
}

func (matcher *PaylineFixedMatcher) match(result *GameResultManager, rule *protodata.SlotPayTable, lines []*PayLine) (bool, *PayResult) {
	result.ReelMode()

	paylineres := new(PayResult)
	paylineres.SymbolName = rule.GetSymbol()
	paylineres.Multi = 1

	symbolMap := make(map[string]bool)

	symbols := strings.Split(rule.GetSymbol(), ",")
	for _, symbol := range symbols {
		symbolMap[symbol] = true
	}

	for _, line := range lines {

		symbolPayLine := new(SymbolPayLine)
		symbolPayLine.PayLine = line

		var repeat int

		// 保險起見，比對前先把Position排序過
		sort.Slice(line.Positions, func(i, j int) bool {
			return line.Positions[i].Reel < line.Positions[j].Reel
		})

		// 檢查該位置是否為指定符號
		for _, position := range line.Positions {
			if position.Reel <= 0 || position.Row <= 0 {
				break
			}

			if result.Reel < position.Reel {
				break
			}

			if len(result.GameResult[position.Reel-1]) < int(position.Row) {
				break
			}

			symbol := result.GameResult[position.Reel-1][position.Row-1]

			if symbol == nil {
				break
			}

			var match bool

			var pluri *protodata.SymbolAttributes
			if attr, ok := matcher.PlurisMap[symbol.GetSymbol()]; ok {
				pluri = attr
			}

			switch LinkType(rule.GetLinkType()) {
			case LinkTypeRight:
				position.Reel = result.Reel - position.Reel + 1
			}

			for symbolName := range symbolMap {
				sbn := symbol.GetSymbol()
				if pluri != nil {
					sbn = pluri.GetSubstituteSymbol()
				}
				if sbn == symbolName {
					repeat++
					if pluri != nil {
						repeat += int(pluri.GetPluriCount()) - 1
					}
					symbolPayLine.SymbolIDs = append(symbolPayLine.SymbolIDs, int32(symbol.GetId()))
					symbolPayLine.SymbolNames = append(symbolPayLine.SymbolNames, symbol.GetSymbol())
					symbolPayLine.BlinkPositions = append(symbolPayLine.BlinkPositions, position)
					match = true
					break
				}
			}

			if match {
				continue
			}

			// 是否百搭
			for _, wild := range matcher.Wilds {
				sbn := symbol.GetSymbol()
				if pluri != nil {
					sbn = pluri.GetSubstituteSymbol()
				}
				if wild.GetSymbol() != sbn {
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
					if pluri != nil {
						repeat += int(pluri.GetPluriCount()) - 1
					}
					symbolPayLine.SymbolIDs = append(symbolPayLine.SymbolIDs, int32(symbol.GetId()))
					symbolPayLine.SymbolNames = append(symbolPayLine.SymbolNames, symbol.GetSymbol())
					symbolPayLine.BlinkPositions = append(symbolPayLine.BlinkPositions, position)
					match = true
					symbolPayLine.HasWild = true
				}
			}

			if !match {
				break
			}

		}

		if repeat == 0 {
			continue
		}

		// 檢查每個位置要有一個此規則的符號
		var hasRuleSymbol bool
		for _, pos := range symbolPayLine.BlinkPositions {
			if hasRuleSymbol {
				break
			}

			switch LinkType(rule.GetLinkType()) {
			case LinkTypeRight:
				pos.Reel = result.Reel - pos.Reel + 1
			}

			for sbn := range symbolMap {
				// 檢查是否有被替換符號

				if attr, ok := matcher.PlurisMap[result.GameResult[pos.Reel-1][pos.Row-1].GetSymbol()]; ok && attr.GetSubstituteSymbol() == sbn {
					hasRuleSymbol = true
					break
				}

				if sbn == result.GameResult[pos.Reel-1][pos.Row-1].GetSymbol() {
					hasRuleSymbol = true
					break
				}
			}
		}

		if repeat == int(rule.GetCount()) && hasRuleSymbol {

			for _, symbolName := range symbolPayLine.SymbolNames {
				for _, wild := range matcher.Wilds {
					if symbolName != wild.GetSymbol() {
						continue
					}
					if wild.GetMultiplier() > 1 {
						symbolPayLine.MultipleStep = append(symbolPayLine.MultipleStep, wild.GetMultiplier())
					}
				}
			}

			paylineres.SymbolPayLines = append(paylineres.SymbolPayLines, symbolPayLine)

			for _, pos := range symbolPayLine.BlinkPositions {
				paylineres.BlinkPositions = append(paylineres.BlinkPositions, pos)
			}
		}
	}

	blinkmap := make(map[Position]bool)
	tmpblink := make([]Position, 0)
	for _, pos := range paylineres.BlinkPositions {
		if _, ok := blinkmap[pos]; !ok {
			blinkmap[pos] = true
			tmpblink = append(tmpblink, pos)
		}
	}
	paylineres.BlinkPositions = tmpblink

	if len(paylineres.SymbolPayLines) == 0 {
		return false, paylineres
	}

	paylineres.IsMatch = true
	paylineres.Rule = rule
	paylineres.WayCount = int32(len(paylineres.SymbolPayLines))

	for _, line := range paylineres.SymbolPayLines {
		if line.HasWild {
			paylineres.WildWayCount++
		} else {
			paylineres.NormalWayCount++
		}
	}

	paylineres.PayCount = rule.GetCount()

	for _, row := range paylineres.SymbolPayLines {
		if row.HasWild {
			paylineres.HasWild = true
			break
		}
	}

	return true, paylineres
}

// MatchLeft 比對GameResult是否符合PayTable規則
func (matcher *PaylineFixedMatcher) MatchLeft(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {
	result.ReelMode()
	return matcher.match(result, rule, matcher.NormalLine)
}

// MatchRight 比對GameResult是否符合PayTable規則
func (matcher *PaylineFixedMatcher) MatchRight(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {
	// 最後要轉回來
	defer func() {
		result = result.Reverse()
	}()
	result.ReelMode()
	// 轉過去
	result = result.Reverse()
	return matcher.match(result, rule, matcher.ReverseLine)
}
