package slot

import (
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
	"sort"
)

// BackgroundMatcher 背景結果比對器
type BackgroundMatcher interface {
	IsMatch(*GameResultManager, *protodata.SlotPayTable) (bool, *PayResult)
}

// BasicBackgroundMatcher 基本背景比對器
type BasicBackgroundMatcher struct {
	// MatchLeft  func(*GameResultManager, *protodata.SlotPayTable) (bool, *PayResult)
	// MatchRight func(*GameResultManager, *protodata.SlotPayTable) (bool, *PayResult)
}

// IsMatch 是否符合
func (m *BasicBackgroundMatcher) IsMatch(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {
	if result == nil || rule == nil {
		return false, nil
	}

	if result.BackgroundResult == nil {
		return false, nil
	}

	switch LinkType(rule.GetLinkType()) {
	case LinkTypeAny:
		return m.MatchAny(result, rule)
	case LinkTypeLeft:
		return m.MatchLeft(result, rule)
	case LinkTypeRight:
		return m.MatchRight(result, rule)
	case LinkTypeAdjacent:
		return m.MatchAdjacent(result, rule)
	}

	return false, nil
}

// MatchAny ANY 規則
func (m *BasicBackgroundMatcher) MatchAny(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {

	var count int

	var pos []Position

	for _, dist := range result.BackgroundResult.Distributions {
		if dist.GetBgEnum() == rule.GetSymbol() {
			count++
			pos = append(pos, Position{Reel: int32(dist.GetReelNum())})
		}
	}

	payresult := new(PayResult)

	if int32(count) == rule.GetCount() {
		payresult.FromBackground = true
		payresult.IsMatch = true
		payresult.SymbolName = rule.GetSymbol()
		payresult.Rule = rule
		payresult.Multi = 1
		payresult.WayCount = 1
		payresult.PayCount = int32(count)
		payresult.BlinkPositions = pos
	}

	return payresult.IsMatch, payresult
}

// MatchLeft LEFT 規則
func (m *BasicBackgroundMatcher) MatchLeft(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {

	var count int

	sort.Slice(result.BackgroundResult.Distributions, func(i, j int) bool {
		return result.BackgroundResult.Distributions[i].GetReelNum() < result.BackgroundResult.Distributions[j].GetReelNum()
	})

	var pos []Position

	for _, dist := range result.BackgroundResult.Distributions {
		if dist.GetBgEnum() == rule.GetSymbol() {
			count++
			pos = append(pos, Position{Reel: int32(dist.GetReelNum())})
		} else {
			break
		}
	}

	payresult := new(PayResult)

	if int32(count) == rule.GetCount() {
		payresult.FromBackground = true
		payresult.IsMatch = true
		payresult.SymbolName = rule.GetSymbol()
		payresult.Rule = rule
		payresult.Multi = 1
		payresult.WayCount = 1
		payresult.PayCount = int32(count)
		payresult.BlinkPositions = pos
	}

	return payresult.IsMatch, payresult
}

// MatchRight Right 規則
func (m *BasicBackgroundMatcher) MatchRight(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {

	var count int

	sort.Slice(result.BackgroundResult.Distributions, func(i, j int) bool {
		return result.BackgroundResult.Distributions[i].GetReelNum() > result.BackgroundResult.Distributions[j].GetReelNum()
	})

	var pos []Position

	for _, dist := range result.BackgroundResult.Distributions {
		if dist.GetBgEnum() == rule.GetSymbol() {
			count++
			pos = append(pos, Position{Reel: int32(dist.GetReelNum())})
		} else {
			break
		}
	}

	payresult := new(PayResult)

	if int32(count) == rule.GetCount() {
		payresult.FromBackground = true
		payresult.IsMatch = true
		payresult.SymbolName = rule.GetSymbol()
		payresult.Rule = rule
		payresult.Multi = 1
		payresult.WayCount = 1
		payresult.PayCount = int32(count)
		payresult.BlinkPositions = pos
	}

	return payresult.IsMatch, payresult
}

// MatchAdjacent 連續相鄰
func (m *BasicBackgroundMatcher) MatchAdjacent(result *GameResultManager, rule *protodata.SlotPayTable) (bool, *PayResult) {
	var count int

	sort.Slice(result.BackgroundResult.Distributions, func(i, j int) bool {
		return result.BackgroundResult.Distributions[i].GetReelNum() < result.BackgroundResult.Distributions[j].GetReelNum()
	})

	var pos []Position

	for _, dist := range result.BackgroundResult.Distributions {
		if dist.GetBgEnum() == rule.GetSymbol() {
			count++
			pos = append(pos, Position{Reel: int32(dist.GetReelNum())})
		} else if count == 0 {
			continue
		} else {
			break
		}
	}

	payresult := new(PayResult)

	if int32(count) == rule.GetCount() {
		payresult.FromBackground = true
		payresult.IsMatch = true
		payresult.SymbolName = rule.GetSymbol()
		payresult.Rule = rule
		payresult.Multi = 1
		payresult.WayCount = 1
		payresult.PayCount = int32(count)
		payresult.BlinkPositions = pos
	}

	return payresult.IsMatch, payresult
}
