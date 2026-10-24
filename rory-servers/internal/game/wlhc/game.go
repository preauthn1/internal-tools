package wlhc

import (
	"sort"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// PrepareSpin 取得SymbolGroupID, BgGroupID, 累積JackpotPool
func (logic *Logic) PrepareSpin(mainID int32) (symbolGroupID, bgGroupID int32) {

	logic.CheckModeExpire()

	if logic.FreeSpin != nil {
		symbolGroupID = logic.FreeSpin.SymbolGroupID
		bgGroupID = logic.FreeSpin.BgGroupID

		logic.GameLog.AddGameStep("玩家 %d; 底注: %d, 倍率: %d, 免費遊戲剩餘次數: %d", logic.GameUser.UserID, logic.FreeSpin.BetMoney, logic.FreeSpin.BetRate, logic.FreeSpin.Count)
	} else {
		symbolGroupID = logic.DefaultSymbolGroupID
		bgGroupID = logic.DefaultBgGroupID

		// 依照下注倍率選擇MainGameTable
		for _, table := range logic.Proto.MainGameTable {
			if mainID > 0 {
				if mainID == int32(table.GetId()) {
					symbolGroupID = int32(table.GetSymbolGroupId())
					bgGroupID = int32(table.GetGoldReelGroupId())
					logic.CurrentMain = table
					break
				}
				continue
			}

			f, _ := strconv.ParseFloat(table.GetConditionValue(), 64)
			if int64(f) == logic.GameBetRate {
				symbolGroupID = int32(table.GetSymbolGroupId())
				bgGroupID = int32(table.GetGoldReelGroupId())
				break
			}
		}

		logic.GameLog.AddGameStep("玩家 %d; 底注: %d, 倍率: %d, 花費: %d", logic.GameUser.UserID, logic.GameBetMoney, logic.GameBetRate, logic.GamePayMoney)

		logic.IncreaseJackpotPool()

		logic.SendJackpotMessageToUser(logic.GameUser)
	}
	return
}

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID int32) {
	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)

	logic.GameResultManager = manager

	// 檢查盤面是否需要整理
	{
		extRules := logic.GetExtraRuleByConditoin(slot.ExtraRuleConditionAppearSymbol)
		for _, rule := range extRules {

			ok, appears := slot.MatchAppearSymbol(manager, rule)
			if !ok {
				mlog.Debug("%s 不符合條件, %v", slot.ExtraRuleConditionAppearSymbol, rule.GetConditionValue())
				continue
			}

			switch slot.ExtraRuleAction(rule.GetAction()) {
			case slot.ExtraRuleActionSwitchReelRandom:

				mlog.Debug("%s, %v", rule.GetAction(), appears)
				mlog.Debug("替換前輪軸結果")
				manager.RowMode()
				for i, row := range manager.GetSymbolString() {
					mlog.Debug("第 %d 行: %s", i+1, strings.Join(row, ", "))
				}

				slot.HandleSwitchReelRandom(logic.GameResultManager, rule, appears)
			}
		}
	}

	logic.LogGameResult(logic.GameResultManager)

	logic.CalcGameResult()

	logic.CalcGameAward()
}

// CalcGameResult 計算遊戲結果
func (logic *Logic) CalcGameResult() {
	logic.GamePayResult = logic.SlotLogicBase.CalcGameResult(logic.GameResultManager)
}

// CalcGameAward 計算遊戲獎勵
func (logic *Logic) CalcGameAward() {

	var multi int64 = 1

	// 先檢查盤面是否出現百搭符號，決定加倍數量
	type MultiOpt struct {
		Value  int64
		Weight int
	}
	var opts []MultiOpt

	// 決定倍率選項
	{
		var str string
		if logic.IsFreeSpin() {
			str = logic.FreeSpin.Origin.GetMultiplierValue()
		} else {
			str = logic.CurrentMain.GetMultiplierValue()
		}

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

			var multis []int
			for _, opt := range opts {
				multis = append(multis, int(opt.Value))
			}
			sort.Slice(multis, func(i, j int) bool {
				return multis[i] < multis[j]
			})

			r := new(prob.Random)
			for _, opt := range opts {
				r.Add(opt, opt.Weight)
			}
			t := r.Choose().(MultiOpt)

			logic.GameLog.AddGameStep("百搭加倍選項 %v 倍, 骰出 %v 倍", multis, t.Value)

			multi = t.Value
		}
	}

	// 結算
	for _, result := range logic.GamePayResult {

		if !result.IsMatch || result.Rule == nil {
			continue
		}

		switch result.Rule.GetUseType() {
		case "normal":
			if logic.IsFreeSpin() {
				continue
			}
		case "freespin":
			if !logic.IsFreeSpin() {
				continue
			}
		}

		switch result.GetPayType() {
		case slot.PayTypeChip, slot.PayTypeChipXRate:
			// 計算贏得的金額
			result.Multi = int32(multi)
			result.MultiType = slot.MultiTypeLinkWild

			var base, win int64
			if logic.IsFreeSpin() {
				base, win = result.CalcWinMoney(logic.FreeSpin.BetMoney, logic.FreeSpin.BetRate)
				logic.FreeSpin.WinMoney += win
			} else {
				base, win = result.CalcWinMoney(logic.GameBetMoney, logic.GameBetRate)
			}

			logic.GameWinMoney += win
			logic.LogPayResultChip(result, base, win)

		case slot.PayTypeMysteryJPG, slot.PayTypeMysteryJPGR:

			jpg := logic.WinMysteryJPG(result.Rule)

			if jpg == nil {
				continue
			}

			logic.JackpotGame = jpg

			logic.RecordJackpotGame(logic.GameUser)

		case slot.PayTypeSelection:

			selection := logic.WinSelection(result.Rule)

			if selection == nil {
				continue
			}

			logic.SelectChances = append(logic.SelectChances, selection)

			if logic.IsFreeSpin() {
				logic.GameLog.AddGameStep("因目前正於免費遊戲中, 將於免費遊戲結束後開始選擇, 剩餘 %d 次", len(logic.SelectChances))
				logic.SendIncreaseGameTimesMessage(logic.GameUser, string(slot.GameModeSelection), 0, 1, "")
			}

			logic.RecordSelectGameTimes(logic.GameUser)

		case slot.PayTypePickGame:

			game := logic.WinPickGame(result.Rule, result.PayCount, true)

			if game == nil {
				continue
			}

			// 位置加工
			var positions []slot.Position
			logic.GameResultManager.RowMode()
			for rowID, row := range logic.GameResultManager.GameResult {
				for reelID, symbol := range row {
					if symbol.GetSymbol() == result.SymbolName {
						positions = append(positions, slot.Position{Reel: int32(reelID + 1), Row: int32(rowID + 1)})
					}
				}
			}
			game.Positions = positions

			logic.PickGame = game

			logic.RecordPickGame(logic.GameUser)
		}
	}

}
