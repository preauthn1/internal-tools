package jgyf

import (
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID, reelNum, rowNum int32) {
	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(reelNum, rowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)

	logic.LogGameResult(manager)

	logic.GameResultManager = manager

	logic.CalcGameResult()

	logic.CalcGameAward()
}

// CalcGameResult 計算遊戲結果
func (logic *Logic) CalcGameResult() {
	logic.GamePayResult = logic.SlotLogicBase.CalcGameResult(logic.GameResultManager)
}

// CalcGameAward 計算遊戲獎勵
func (logic *Logic) CalcGameAward() {

	var ChipWinMoney int64 = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)
	var HistoryType int32

	if logic.IsFreeSpin() {
		HistoryType = slot.HistoryTypeFREE
		logic.FreeSpin.WinMoney += ChipWinMoney
		if ChipWinMoney > 0 {
			logic.FreeSpin.WinCount++
		}
	} else {
		HistoryType = slot.HistoryTypeMAIN
	}

	logic.GameWinMoney = ChipWinMoney

	// 找出背景贏得的獎勵拆開記Log
	var bgMoney int64
	for _, result := range logic.GamePayResult {
		if !result.FromBackground {
			continue
		}
		bgMoney += result.WinMoney
	}

	logic.GameResultManager.ReelMode()
	logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score, logic.GamePayMoney, ChipWinMoney, HistoryType, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	logic.GameResultManager.RowMode()

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
				logic.SendSelectChanceMessage(logic.GameUser)
			}

			logic.RecordSelectGameTimes(logic.GameUser)

		case slot.PayTypeAddTimes:

			cnt := logic.WinFreeTimes(result.Rule, logic.FreeSpin)

			logic.SendIncreaseGameTimesMessage(logic.GameUser, string(slot.GameModeFreeSpin), result.Rule.GetId(), int32(cnt), "")

			logic.UpdateFreeSpin(logic.GameUser)
		}
	}
}
