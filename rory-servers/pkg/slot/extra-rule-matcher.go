package slot

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// MatchAfterSpin 依照Payout規則比對盤面
func MatchAfterSpin(result *GameResultManager, rule *protodata.ExtraRule, machter PaylineMatcher) (ok bool) {

	linkType, symbol, count, err := ParseExtraRuleConditionAfterSpin(rule.GetConditionValue())
	if err != nil {
		mlog.Debug("%v", err)
		return false
	}

	r := &protodata.SlotPayTable{
		LinkType: proto.String(linkType),
		Symbol:   proto.String(symbol),
		Count:    proto.Int32(int32(count)),
	}

	if ok, _ := machter.IsMatch(result, r); !ok {
		return false
	}

	return true
}

// MatchAfterBackground 依照Payout規則比對盤面
func MatchAfterBackground(result *GameResultManager, rule *protodata.ExtraRule, machter BackgroundMatcher) (ok bool) {

	linkType, symbol, count, err := ParseExtraRuleConditionAfterSpin(rule.GetConditionValue())
	if err != nil {
		mlog.Debug("%v", err)
		return false
	}

	matchRule := &protodata.SlotPayTable{
		LinkType: proto.String(linkType),
		Symbol:   proto.String(symbol),
		Count:    proto.Int32(int32(count)),
	}
	if match, _ := machter.IsMatch(result, matchRule); match {
		return true
	}

	return false
}

// MatchAppearSymbol 是否出現符號
func MatchAppearSymbol(result *GameResultManager, rule *protodata.ExtraRule) (ok bool, positions []Position) {

	symbol, ror, vals, err := ParseExtraRuleConditionAppearSymbol(rule.GetConditionValue())
	if err != nil {
		mlog.Debug("%v", err)
		return false, []Position(nil)
	}

	// 檢查盤面是否符合規則
	switch ror {
	case "Reel":
		result.ReelMode()
	case "Row":
		result.RowMode()
	}

	var appears []Position

	for id, row := range result.GameResult {
		var needCheck bool
		if ror == "Any" {
			needCheck = true
		} else {
			for _, i := range vals {
				if i == id+1 {
					needCheck = true
					break
				}
			}
		}

		if !needCheck {
			continue
		}

		var appear bool
		for i, s := range row {
			if symbol == s.GetSymbol() {
				appear = true
				pos := Position{}
				if result.IsReelMode() {
					pos.Reel = int32(id + 1)
					pos.Row = int32(i + 1)
				} else {
					pos.Reel = int32(i + 1)
					pos.Row = int32(id + 1)
				}
				appears = append(appears, pos)
			}
		}

		if !appear {
			continue
		}
	}

	// 沒中
	if len(appears) == 0 {
		return false, []Position(nil)
	}

	return true, appears
}

// HandleRespin 處理重轉盤面, 沿用GameResultManager內的SymbolGroupID重轉
// @param result 要重轉的輪軸
// @param prop 放行的機率(不會Respin)
// @param denominator 機率分母
// @param dists 所有Symbol機率表
// @param diced 是否已經決定重轉過
func HandleRespin(result *GameResultManager, prop, denominator int64, dists []*protodata.SlotSymbolDistribution, diced bool) bool {

	var respin bool = true

	mlog.Debug("[Respin] prop: %d, denominator: %d", prop, denominator)

	// 骰看看要不要重轉
	if diced == false {

		mlog.Debug("決定是否要重轉")

		r := new(prob.Random)
		if prop > 0 {
			r.Add(false, int(prop))
		}

		if denominator-prop > 0 {
			r.Add(true, int(denominator-prop))
		}

		t := r.Choose()

		r.PrintOptions()

		if t != nil {
			respin = t.(bool)
		}
	}

	if !respin {
		mlog.Debug("不重轉")
		return false
	}

	result.GenerateGameResult(result.SymbolGroupID, dists)

	return true
}

// HandleRespinBG 處理重轉盤面, 沿用GameResultManager內的SymbolGroupID重轉
// @param result 要重轉的輪軸
// @param prop 放行的機率(不會Respin)
// @param denominator 機率分母
// @param dists 所有Symbol機率表
// @param diced 是否已經決定重轉過
func HandleRespinBG(result *GameResultManager, prop, denominator int64, dists []*protodata.ReelBgDistribution, diced bool) bool {

	var respin bool = true

	mlog.Debug("[Respin] prop: %d, denominator: %d", prop, denominator)

	// 骰看看要不要重轉
	if diced == false {

		mlog.Debug("決定是否要重轉")

		r := new(prob.Random)
		if prop > 0 {
			r.Add(false, int(prop))
		}

		if denominator-prop > 0 {
			r.Add(true, int(denominator-prop))
		}

		t := r.Choose()

		r.PrintOptions()

		if t != nil {
			respin = t.(bool)
		}
	}

	if !respin {
		mlog.Debug("不重轉")
		return false
	}

	result.GenerateBackground(result.BackGroundGroupID, dists)

	return true
}

// HandleExpandReel 處理擴展直軸
func HandleExpandReel(result *GameResultManager, rule *protodata.ExtraRule, positions []Position) {
	for _, pos := range positions {
		if pos.Reel-1 < 0 {
			mlog.Warn("Error index from appear")
			continue
		}

		result.ReelMode()
		symbols := result.GameResult[pos.Reel-1]
		for _, s := range symbols {
			s.Symbol = proto.String(rule.GetActionValue())
		}
	}
}

// HandleSwitchReelRandom 處理交換符號, 傳入的符號位置與同Reel的符號隨機交換
func HandleSwitchReelRandom(result *GameResultManager, rule *protodata.ExtraRule, positions []Position) {

	for _, pos := range positions {
		if pos.Reel-1 < 0 {
			mlog.Warn("Error index from appear")
			continue
		}

		result.ReelMode()

		origin := result.GameResult[pos.Reel-1][pos.Row-1]

		r := new(prob.Random)
		for i := 0; i < int(result.Row); i++ {
			r.Add(Position{Reel: pos.Reel, Row: int32(i + 1)}, 1)
		}
		t := r.Choose().(Position)

		tmp := result.GameResult[t.Reel-1][t.Row-1]

		mlog.Debug("(%s | %v) switch to (%s | %v)", origin.GetSymbol(), pos, tmp.GetSymbol(), t)

		result.GameResult[t.Reel-1][t.Row-1] = origin
		result.GameResult[pos.Reel-1][pos.Row-1] = tmp
	}
}

// HandleReplaceBG 處理替代背景, 骰出是否需要替換背景, 替換指定的背景
func HandleReplaceBG(result *GameResultManager, ror string, target int, symbol string, prop, denominator int64) {

	// 決定是否替換背景
	var isReplace bool
	if prop > denominator {
		mlog.Debug("ReplaceBG 機率 %d 大於分母 %d, 直接替換", prop, denominator)
		isReplace = true
	} else {
		r := new(prob.Random)
		r.Add(true, int(prop))
		r.Add(false, int(denominator-prop))

		isReplace = r.Choose().(bool)

		mlog.Debug("ReplaceBG 機率 %d, 分母 %d, 骰出結果 %t", prop, denominator, isReplace)
	}

	if !isReplace {
		return
	}

	switch ror {
	case "Reel":
		var newbg []*protodata.ReelBgDistribution
		for i, dist := range result.BackgroundResult.Distributions {
			if i+1 == target {
				dist.BgEnum = proto.String(symbol)
			}
			newbg = append(newbg, dist)
		}

		// clean background
		result.BackgroundResult.Distributions = make([]*protodata.ReelBgDistribution, 0)

		// add dist
		for _, dist := range newbg {
			result.BackgroundResult.AddDistribution(dist)
		}
	}
}
