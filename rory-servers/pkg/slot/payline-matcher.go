package slot

import (
	"sort"
	"strings"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// PaylineMatcher 賠付線比對器
type PaylineMatcher interface {
	IsMatch(*GameResultManager, *protodata.SlotPayTable) (bool, *PayResult)
}

// BasicPayLineMatcher 基本賠線比對器，實作 MatchAny，MatchStraight
type BasicPayLineMatcher struct {
	Wilds           []*protodata.SymbolAttributes
	Pluris          []*protodata.SymbolAttributes
	PlurisMap       map[string]*protodata.SymbolAttributes
	GoladBalls      []*protodata.SymbolAttributes
	GoladBallMap    map[string]*protodata.SymbolAttributes
	GoladBallPrefix string
	MatchLeft       func(*GameResultManager, *protodata.SlotPayTable) (bool, *PayResult)
	MatchRight      func(*GameResultManager, *protodata.SlotPayTable) (bool, *PayResult)
}

// NewBasicPayLineMatcher create BasicPayLineMatcher
func NewBasicPayLineMatcher(attrs []*protodata.SymbolAttributes) *BasicPayLineMatcher {
	matcher := new(BasicPayLineMatcher)

	matcher.GoladBallMap = make(map[string]*protodata.SymbolAttributes)
	matcher.PlurisMap = make(map[string]*protodata.SymbolAttributes)

	for _, attr := range attrs {
		if attr.GetIsWild() > int32(0) {
			matcher.Wilds = append(matcher.Wilds, attr)
		}
		if attr.GetIsGoldBall() == int32(1) {
			matcher.GoladBallPrefix = attr.GetSymbol()[0:1]
			matcher.GoladBalls = append(matcher.GoladBalls, attr)
			matcher.GoladBallMap[attr.GetSymbol()] = attr
		}
		if attr.GetIsPluri() == int32(1) {
			matcher.Pluris = append(matcher.Pluris, attr)
			matcher.PlurisMap[attr.GetSymbol()] = attr
		}
	}

	return matcher
}

// IsMatch 是否符合
func (m *BasicPayLineMatcher) IsMatch(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {
	if result == nil || rule == nil {
		return false, nil
	}

	switch LinkType(rule.GetLinkType()) {
	case LinkTypeAny:
		return m.MatchAny(result, rule)
	case LinkTypeLeft:
		return m.MatchLeft(result, rule)
	case LinkTypeBoth:
		return m.MatchBoth(result, rule)
	case LinkTypeRight:
		return m.MatchRight(result, rule)
	case LinkTypeStraight:
		return m.MatchStraight(result, rule)
	}

	return false, nil
}

// MatchAny ANY 規則
func (m *BasicPayLineMatcher) MatchAny(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {

	result.ReelMode()

	var count int
	var symbolIDs []int32
	var blinkPos []Position
	// 掃描整個GameResult有幾個對應的Symbol
	for reelID, reel := range result.GameResult {
		for rowID, symbol := range reel {
			// 一般符號比對
			if symbol.GetSymbol() == rule.GetSymbol() {
				count++
				symbolIDs = append(symbolIDs, int32(symbol.GetId()))
				blinkPos = append(blinkPos, Position{int32(reelID + 1), int32(rowID + 1)})
				continue
			}

			// Lightning類型，比對此符號是否為金球
			if strings.HasPrefix(rule.GetSymbol(), "g1") {
				// if strings.HasPrefix(rule.GetSymbol(), m.GoladBallPrefix) {
				if m.GoladBallMap[symbol.GetSymbol()] != nil {
					count++
					symbolIDs = append(symbolIDs, int32(symbol.GetId()))
					blinkPos = append(blinkPos, Position{int32(reelID + 1), int32(rowID + 1)})
					continue
				}
			}
		}
	}

	sort.Slice(symbolIDs, func(i, j int) bool {
		return symbolIDs[i] < symbolIDs[j]
	})

	payresult := new(PayResult)

	if int32(count) == rule.GetCount() {
		payresult.IsMatch = true
		payresult.SymbolName = rule.GetSymbol()
		payresult.Rule = rule
		payresult.WayCount = 1
		payresult.PayCount = int32(count)
		payresult.SymbolIDs = symbolIDs
		payresult.BlinkPositions = blinkPos
	}

	return payresult.IsMatch, payresult
}

// MatchStraight Straight 規則
func (m *BasicPayLineMatcher) MatchStraight(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {

	result.ReelMode()

	var count int
	var symbolIDs []int32

	var symbolMap = make(map[string]bool)
	symbols := strings.Split(rule.GetSymbol(), ",")
	for _, symbol := range symbols {
		symbolMap[symbol] = true
	}

	// 掃描整個GameResult有幾個對應的Symbol
	for _, reel := range result.GameResult {
		var hasDiff bool
		for _, symbol := range reel {
			if val, ok := symbolMap[symbol.GetSymbol()]; !ok || !val {
				hasDiff = true
				break
			}
		}
		if hasDiff {
			continue
		}

		for _, symbol := range reel {
			symbolIDs = append(symbolIDs, int32(symbol.GetId()))
		}

		count++
	}

	sort.Slice(symbolIDs, func(i, j int) bool {
		return symbolIDs[i] < symbolIDs[j]
	})

	payresult := new(PayResult)

	if int32(count) == rule.GetCount() {
		payresult.IsMatch = true
		payresult.SymbolName = rule.GetSymbol()
		payresult.Rule = rule
		payresult.WayCount = 1
		payresult.PayCount = int32(count)
		payresult.SymbolIDs = symbolIDs
	}

	return payresult.IsMatch, payresult
}

// MatchBoth Both 規則, 左線與右線均符合
func (m *BasicPayLineMatcher) MatchBoth(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {

	res := new(PayResult)
	res.IsMatch = true
	res.Rule = rule

	{
		r := new(protodata.SlotPayTable)
		r.LinkType = proto.String(LinkTypeLeft.String())
		r.Symbol = proto.String(rule.GetSymbol())
		r.Count = proto.Int32(rule.GetCount())

		if match, pr := m.MatchLeft(result, r); match {
			res.SymbolPayLines = append(res.SymbolPayLines, pr.SymbolPayLines...)
			res.HasWild = res.HasWild || pr.HasWild
			res.WayCount += pr.WayCount
			res.NormalWayCount += pr.NormalWayCount
			res.WildWayCount += pr.WildWayCount
		} else {
			return false, nil
		}
	}

	{
		r := new(protodata.SlotPayTable)
		r.LinkType = proto.String(LinkTypeRight.String())
		r.Symbol = proto.String(rule.GetSymbol())
		r.Count = proto.Int32(rule.GetCount())

		if match, pr := m.MatchRight(result, r); match {
			res.SymbolPayLines = append(res.SymbolPayLines, pr.SymbolPayLines...)
			res.HasWild = res.HasWild || pr.HasWild
			res.WayCount += pr.WayCount
			res.NormalWayCount += pr.NormalWayCount
			res.WildWayCount += pr.WildWayCount
		} else {
			return false, nil
		}
	}

	return true, res
}
