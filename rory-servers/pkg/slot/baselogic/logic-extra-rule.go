package baselogic

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/defender"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// GetExtraRuleByConditoin 依照條件取得額外條件
func (logic *SlotLogicBase) GetExtraRuleByConditoin(condi slot.ExtraRuleCondition, acts ...slot.ExtraRuleAction) []*protodata.ExtraRule {

	var rules []*protodata.ExtraRule

	for _, rule := range logic.Proto.ExtraRule {
		switch slot.UseType(rule.GetUseType()) {
		default:
			continue
		case slot.UseTypeALL:

		case slot.UseTypeMainGame:
			if logic.IsFreeSpin() {
				continue
			}
		case slot.UseTypeFreeGame:
			if !logic.IsFreeSpin() {
				continue
			}
		case slot.UseTypeRespin:

		}

		if slot.ExtraRuleCondition(rule.GetCondition()) != condi {
			continue
		}

		if len(acts) == 0 {
			rules = append(rules, rule)
			continue
		}

		for _, act := range acts {
			if slot.ExtraRuleAction(rule.GetAction()) == act {
				rules = append(rules, rule)
			}
		}
	}

	return rules
}

// HandleAfterSpin 處理盤面, 當出現符合條件的規則
func (logic *SlotLogicBase) HandleAfterSpin(manager *slot.GameResultManager) {

	rules := logic.GetExtraRuleByConditoin(slot.ExtraRuleConditionAfterSpin)

	for _, rule := range rules {

		if ok := slot.MatchAfterSpin(manager, rule, logic.PaylineMatcher); !ok {
			mlog.Debug("AfterSpin 不符合條件, %v", rule.GetConditionValue())
			continue
		}

		switch slot.ExtraRuleAction(rule.GetAction()) {
		case slot.ExtraRuleActionRespin, slot.ExtraRuleActionRespinR:

			mlog.Debug("%s", rule.GetAction())
			mlog.Debug("替換前輪軸結果")
			manager.RowMode()
			for i, row := range manager.GetSymbolString() {
				mlog.Debug("第 %d 行: %s", i+1, strings.Join(row, ", "))
			}

			// prop 為放行不Respin
			var prop, denominator int64
			var err error

			switch slot.ExtraRuleAction(rule.GetAction()) {
			case slot.ExtraRuleActionRespin:
				prop, denominator, err = slot.ParseExtraRuleActionRespin(rule.GetActionValue())
			case slot.ExtraRuleActionRespinR:
				prop, denominator, err = slot.ParseExtraRuleActionRespinR(rule.GetActionValue(), logic.GameBetRate)
				prop = prop * logic.GameBetMoney * logic.GameBetRate
			}

			if err != nil {
				continue
			}

			if slot.HandleRespin(manager, prop, denominator, logic.Proto.SlotSymbolDistribution, false) {
				for slot.MatchAfterSpin(manager, rule, logic.PaylineMatcher) {
					mlog.Debug("符合條件再重轉")
					if !slot.HandleRespin(manager, prop, denominator, logic.Proto.SlotSymbolDistribution, true) {
						break
					}
				}
			}

		case slot.ExtraRuleActionConditionalAction:
			mlog.Debug("ConditionalAction")
			// 解析要檢查的項目與動作
			condition, action, err := slot.ParseExtraRuleActionConditionalAction(rule.GetActionValue())
			if err != nil {
				mlog.Debug("Parse ConditionalAction Error, %v", err)
				continue
			}

			switch condition {
			case slot.ExtraRuleConditionJPDefense:
				if defender.JPDefender == nil {
					mlog.Debug("JPDefender not setting, skip")
					continue
				}

				mlog.Debug("JPDefender Checking")

				// 找出Jackpot
				var jackpot *slot.Jackpot

				for _, pay := range logic.CalcGameResult(manager) {
					pattern := fmt.Sprintf("%s:%s:%d", pay.Rule.GetLinkType(), pay.Rule.GetSymbol(), pay.Rule.GetCount())
					if pattern != rule.GetConditionValue() {
						continue
					}
					if pay.GetPayType() != slot.PayTypeJackpot {
						continue
					}

					f, _ := strconv.ParseFloat(pay.Rule.GetPayout(), 64)
					for _, jp := range logic.Jackpots {
						if jp.Origin.GetId() == int32(f) {
							jackpot = jp
						}
					}
				}

				if jackpot == nil {
					mlog.Debug("Can not find Jackpot")
					continue
				}

				poolFn := func() int64 {
					return jackpot.GetCurrentPoolByBet(logic.GameBetMoney*int64(logic.GameBetRate), int64(logic.MaxBetCount*logic.MaxBetRate))
				}
				if defender.JPDefender.Eligible(jackpot, poolFn) {
					continue
				}

				mlog.Debug("[JPDefender] Illegible Jackpot[%v], JPDefender Pool: %d, Win: %d", jackpot.Level, defender.JPDefender.Pool(), poolFn())

				switch action {
				case slot.ExtraRuleActionRespin:
					if slot.HandleRespin(manager, 0, 1, logic.Proto.SlotSymbolDistribution, true) {
						for slot.MatchAfterSpin(manager, rule, logic.PaylineMatcher) {
							mlog.Debug("符合條件再重轉")
							if !slot.HandleRespin(manager, 0, 1, logic.Proto.SlotSymbolDistribution, true) {
								break
							}
						}
					}
				}
			}
		}
	}
}

// HandleAppearSymbol 處理盤面, 當出現符合條件的符號
func (logic *SlotLogicBase) HandleAppearSymbol(manager *slot.GameResultManager) {

	rules := logic.GetExtraRuleByConditoin(slot.ExtraRuleConditionAppearSymbol)
	for _, rule := range rules {

		ok, appears := slot.MatchAppearSymbol(manager, rule)
		if !ok {
			mlog.Debug("AppearSymbol 不符合條件, %v", rule.GetConditionValue())
			continue
		}

		switch slot.ExtraRuleAction(rule.GetAction()) {
		case slot.ExtraRuleActionExpandReel:

			mlog.Debug("%s, %v", rule.GetAction(), appears)
			mlog.Debug("替換前輪軸結果")
			manager.RowMode()
			for i, row := range manager.GetSymbolString() {
				mlog.Debug("第 %d 行: %s", i+1, strings.Join(row, ", "))
			}

			slot.HandleExpandReel(manager, rule, appears)
		case slot.ExtraRuleActionSwitchReelRandom:

			mlog.Debug("%s, %v", rule.GetAction(), appears)
			mlog.Debug("替換前輪軸結果")
			manager.RowMode()
			for i, row := range manager.GetSymbolString() {
				mlog.Debug("第 %d 行: %s", i+1, strings.Join(row, ", "))
			}

			slot.HandleSwitchReelRandom(manager, rule, appears)
		}
	}
}

// HandleAfterBackground 處理背景, 當出現符合條件的規則
func (logic *SlotLogicBase) HandleAfterBackground(manager *slot.GameResultManager) {
	mlog.Debug("檢查換前背景結果, %p, %v", manager.BackgroundResult, manager.BackgroundResult)

	rules := logic.GetExtraRuleByConditoin(slot.ExtraRuleConditionAfterBackground)

	for _, rule := range rules {
		if !slot.MatchAfterBackground(manager, rule, logic.BackgroundMatcher) {
			mlog.Debug("AfterBackground 檢查背景不符合條件, %v", rule.GetConditionValue())
			continue
		}

		switch slot.ExtraRuleAction(rule.GetAction()) {
		case slot.ExtraRuleActionReplaceBG, slot.ExtraRuleActionReplaceBGT, slot.ExtraRuleActionReplaceBGR, slot.ExtraRuleActionReplaceBGTR:
			var ror, symbol string
			var target int
			var prop, denominator int64
			var err error

			switch slot.ExtraRuleAction(rule.GetAction()) {
			case slot.ExtraRuleActionReplaceBG, slot.ExtraRuleActionReplaceBGT:
				ror, target, symbol, prop, denominator, err = slot.ParseExtraRuleActionReplaceBG(rule.GetActionValue())
			case slot.ExtraRuleActionReplaceBGR, slot.ExtraRuleActionReplaceBGTR:
				ror, target, symbol, prop, denominator, err = slot.ParseExtraRuleActionReplaceBGR(rule.GetActionValue(), logic.GameBetRate)
			}

			if err != nil {
				mlog.Warn("%s 解析錯誤. %v", rule.GetAction(), err)
				continue
			}

			if slot.ExtraRuleAction(rule.GetAction()) == slot.ExtraRuleActionReplaceBGT || slot.ExtraRuleAction(rule.GetAction()) == slot.ExtraRuleActionReplaceBGTR {
				prop = prop * logic.GameBetRate * logic.GameBetRate
			}

			mlog.Debug("%s, %s, %d, %s, %d, %d", rule.GetAction(), ror, target, symbol, prop, denominator)
			mlog.Debug("替換前背景結果, %v", manager.BackgroundResult)
			for i := 0; i < int(manager.Row); i++ {
				var bg []string
				for j := 0; j < int(manager.Reel); j++ {
					bg = append(bg, string(manager.BackgroundResult.Result[slot.Position{Reel: int32(j + 1), Row: int32(i + 1)}]))
				}
				mlog.Debug("第 %d 行: %s", i+1, strings.Join(bg, ", "))
			}

			slot.HandleReplaceBG(manager, ror, target, symbol, prop, denominator)

		case slot.ExtraRuleActionConditionalAction:
			// 解析要檢查的項目與動作
			condition, action, err := slot.ParseExtraRuleActionConditionalAction(rule.GetActionValue())
			if err != nil {
				mlog.Debug("Parse ConditionalAction Error, %v", err)
				continue
			}

			mlog.Debug("Conditional: %s, Action: %s", condition, action)

			switch condition {
			default:
				mlog.Debug("No Condition")
				continue
			case slot.ExtraRuleConditionJPDefense:
				if defender.JPDefender == nil {
					mlog.Debug("JPDefender not setting, skip")
					continue
				}

				mlog.Debug("JPDefender Checking")

				// 找出Jackpot
				var jackpot *slot.Jackpot

				for _, pay := range logic.CalcGameResult(manager) {
					pattern := fmt.Sprintf("%s:%s:%d", pay.Rule.GetLinkType(), pay.Rule.GetSymbol(), pay.Rule.GetCount())
					if pattern != rule.GetConditionValue() {
						continue
					}
					if pay.GetPayType() != slot.PayTypeJackpot {
						continue
					}

					f, _ := strconv.ParseFloat(pay.Rule.GetPayout(), 64)
					for _, jp := range logic.Jackpots {
						if jp.Origin.GetId() == int32(f) {
							jackpot = jp
						}
					}
				}

				if jackpot == nil {
					mlog.Debug("Can not find Jackpot")
					continue
				}

				poolFn := func() int64 {
					return jackpot.GetCurrentPoolByBet(logic.GameBetMoney*int64(logic.GameBetRate), int64(logic.MaxBetCount*logic.MaxBetRate))
				}
				if defender.JPDefender.Eligible(jackpot, poolFn) {
					continue
				}

				mlog.Debug("[JPDefender] Illegible Jackpot[%v], JPDefender Pool: %d, Win: %d", jackpot.Level, defender.JPDefender.Pool(), poolFn())

				switch action {
				case slot.ExtraRuleActionRespin:
					if slot.HandleRespinBG(manager, 0, 1, logic.Proto.ReelBgDistribution, true) {
						for slot.MatchAfterBackground(manager, rule, logic.BackgroundMatcher) {
							mlog.Debug("符合條件再重轉")
							if !slot.HandleRespinBG(manager, 0, 1, logic.Proto.ReelBgDistribution, true) {
								break
							}
						}
					}
				}
			}
		}
	}
}

// HandleGetAwardLightning 處理取得Lightning獎勵後的特殊規則
// * UniGoldBall: 下一把塞入一顆金球
func (logic *SlotLogicBase) HandleGetAwardLightning(lightning *LightningGame) {

	if lightning == nil {
		return
	}

	rules := logic.GetExtraRuleByConditoin(slot.ExtraRuleConditionGetAward)

	for _, rule := range rules {
		if slot.GameMode(rule.GetConditionValue()) != slot.GameModeLightning {
			continue
		}
		switch slot.ExtraRuleAction(rule.GetAction()) {
		case slot.ExtraRuleActionUniGoldBall:

			// 只有初始化的LightningGame可以加入新金球
			if len(lightning.Steps) != 1 {
				continue
			}

			symbol, prop, denominator, err := slot.ParseExtraRuleActionUniGoldBall(rule.GetActionValue())
			if err != nil {
				mlog.Debug("%v", err)
				continue
			}

			r := new(prob.Random)
			if prop > 0 {
				r.Add(true, int(prop))
			}
			if denominator-prop > 0 {
				r.Add(false, int(denominator-prop))
			}

			t := r.Choose()

			r.PrintOptions()

			if t.(bool) {
				mlog.Debug("Add Gold Ball %s", symbol)
				step := lightning.GetLastStep()
				step.NextAdd = append(step.NextAdd, symbol)
			}

		}
	}

}
