package tgw

import (
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID int32) {
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)

	logic.GameResultManager = manager

	// 檢查盤面是否需要整理 AfterSpin
	logic.HandleAfterSpin(manager)

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

	var ChipWinMoney int64 = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)
	var HistoryType int32 = slot.HistoryTypeMAIN

	logic.GameWinMoney = ChipWinMoney

	logic.GameResultManager.ReelMode()
	logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score, logic.GamePayMoney, ChipWinMoney, HistoryType, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	logic.GameResultManager.RowMode()

	for _, result := range logic.GamePayResult {
		if !result.IsMatch || result.Rule == nil {
			continue
		}

		switch result.GetPayType() {
		case slot.PayTypePickGame:

			game := logic.WinPickGame(result.Rule, logic.PickGameSize, true)

			if game == nil {
				continue
			}

			logic.PickGame = game

			logic.RecordPickGame(logic.GameUser)

		case slot.PayTypeJackpot:

			_, win := logic.WinJackpot(result.Rule)

			logic.GameWinMoney += win
			logic.GameWinJackpot = win

		}

	}

}
