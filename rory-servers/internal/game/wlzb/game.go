package wlzb

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID int32) {
	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)

	logic.GameResultManager = manager

	// 紀錄開獎符號
	logic.LogGameResult(manager)

	logic.CalcGameResult()

	logic.CalcGameAward()
}

// CalcGameResult 計算遊戲結果
func (logic *Logic) CalcGameResult() {
	logic.GamePayResult = logic.SlotLogicBase.CalcGameResult(logic.GameResultManager)
}

// CalcGameAward 計算遊戲獎勵
func (logic *Logic) CalcGameAward() {

	// 免費遊戲檢查是否有龍加倍
	if logic.FreeSpin != nil {

		var multi int64 = 1

		for _, result := range logic.GamePayResult {

			if !result.IsMatch || result.Rule == nil {
				continue
			}

			switch result.GetPayType() {
			case slot.PayTypeChip:
				// 計算贏得的金額

				var _multi int64 = 1

				// 若有百搭加乘
				if result.HasWild {
					if multi == 1 {
						str := logic.FreeSpin.Origin.GetMultiplierValue()

						multis := strings.Split(str, ",")

						r := new(prob.Random)

						for _, row := range multis {
							vals := strings.Split(row, ":")

							m, _ := strconv.ParseFloat(vals[0], 64)
							w, _ := strconv.ParseFloat(vals[1], 64)

							r.Add(int64(m), int(w))
						}

						multi = r.Choose().(int64)
					}

					_multi = multi
				}

				result.Multi = int32(_multi)
				result.MultiType = slot.MultiTypeLinkWild

				base, win := result.CalcWinMoney(logic.FreeSpin.BetMoney, logic.FreeSpin.BetRate)

				logic.GameWinMoney += win
				logic.FreeSpin.WinMoney += win

				logic.LogPayResultChip(result, base, win)

			}
		}

	} else {
		// 一般遊戲錢照算
		logic.GameWinMoney = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)
	}

	for _, result := range logic.GamePayResult {

		if !result.IsMatch || result.Rule == nil {
			continue
		}

		if result.GetPayType() == slot.PayTypeSelection {

			selection := logic.WinSelection(result.Rule)

			if selection == nil {
				continue
			}

			logic.SelectChances = append(logic.SelectChances, selection)

			if logic.IsFreeSpin() {
				logic.GameLog.AddGameStep("因目前正於免費遊戲中, 將於免費遊戲結束後開始選擇, 剩餘 %d 次", len(logic.SelectChances))
				logic.SendSelectChanceMessage(logic.GameUser)
			}

			logic.RecordSelectGameTimes(logic.GameUser)
		}
	}

}
