package wljr

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID, bgGroupID int32) {
	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)
	manager.GenerateBackground(bgGroupID, logic.Proto.ReelBgDistribution)

	// 檢查背景是否需要整理
	logic.HandleAfterBackground(manager)

	// 紀錄開獎符號
	logic.LogBackgroundResult(manager)

	// 檢查盤面是否需要整理 AppearSymbol
	logic.HandleAppearSymbol(manager)

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

	var ChipWinMoney int64
	var HistoryType int32

	if logic.IsFreeSpin() {
		HistoryType = slot.HistoryTypeFREE
		// 免費遊戲檢查是否有龍加倍
		multi, d := logic.CalcMultiMode(slot.MultiModeDiceOnce, logic.FreeSpin.Origin.GetMultiplierValue())
		ChipWinMoney = logic.SlotLogicBase.CalcGameAwardByMultiType(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate, multi, slot.MultiTypeLinkWild, d)

		logic.FreeSpin.WinCount++
		logic.FreeSpin.WinMoney += ChipWinMoney
	} else {
		HistoryType = slot.HistoryTypeMAIN
		// 一般遊戲錢照算
		ChipWinMoney = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)
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
	if bgMoney == 0 {
		logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score, logic.GamePayMoney, ChipWinMoney, HistoryType, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	} else {
		logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score, logic.GamePayMoney, ChipWinMoney-bgMoney, HistoryType, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
		if !conf.AnalysisMode {
			time.Sleep(100 * time.Millisecond)
		}
		logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score+ChipWinMoney-bgMoney-logic.GamePayMoney, 0, bgMoney, HistoryType, slot.LogTypeBG, logic.GameResultManager.BackgroundResult.String())
	}
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
		case slot.PayTypeJackpot:

			_, win := logic.WinJackpot(result.Rule)
			result.WinMoney = win

			if logic.IsFreeSpin() {
				logic.FreeSpin.WinMoney += win
			}

			if !conf.AnalysisMode {
				time.Sleep(100 * time.Millisecond)
			}

			if result.FromBackground {
				logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score+logic.GameWinMoney-logic.GamePayMoney, 0, win, slot.HistoryTypeJPNOLINK, slot.LogTypeBG, logic.GameResultManager.BackgroundResult.String())
			} else {
				logic.GameResultManager.ReelMode()
				logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score+logic.GameWinMoney-logic.GamePayMoney, 0, win, slot.HistoryTypeJPNOLINK, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
				logic.GameResultManager.RowMode()
			}

			logic.GameWinMoney += win
			logic.GameWinJackpot = win

		case slot.PayTypeMysteryJPG, slot.PayTypeMysteryJPGR:

			c := logic.WinJackpotGameChance(result.Rule, true)

			logic.JackpotGameChances = append(logic.JackpotGameChances, c)

			logic.SendIncreaseGameTimesMessage(logic.GameUser, string(slot.GameModeJackpotGame), 0, 1, "")

			logic.RecordJackpotGameChance(logic.GameUser, c)

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
		}
	}
}
