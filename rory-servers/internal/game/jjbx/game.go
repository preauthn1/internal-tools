package jjbx

import (
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
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

	logic.GameWinMoney = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)

	if logic.FreeSpin != nil && logic.GameWinMoney > 0 {

		multi := logic.FreeSpin.GetMultipleValue()
		win := logic.GameWinMoney
		total := win * int64(multi)

		logic.FreeSpin.WinCount++
		logic.FreeSpin.WinMoney += total
		logic.GameWinMoney = total

		// logic.GameLog.AddGameStep("免費遊戲中第 %d 次中獎, 加乘 %d 倍, 原中獎金額: %d, 加乘後金額: %d", logic.FreeSpin.WinCount, multi, win, total)
	}

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

			logic.GameWinMoney += win
			logic.GameWinJackpot = win

		case slot.PayTypeMysteryJP, slot.PayTypeMysteryJPR:

			_, win := logic.WinMysteryJP(result.Rule)

			logic.GameWinMoney += win
			logic.GameWinJackpot = win

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

			var gMoney int64
			var gCnt int
			for _, reel := range selection.GameResultManager.GameResult {
				for _, symbol := range reel {
					symbol.GetSymbol()
					for _, attr := range logic.Proto.SymbolAttributes {
						if attr.GetIsGoldBall() == 1 && attr.GetSymbol() == symbol.GetSymbol() {
							gCnt++
							switch slot.PayType(attr.GetGoldballType()) {
							case slot.PayTypeChip:
								f, _ := strconv.ParseFloat(attr.GetGoldballValue(), 64)
								win := selection.BetMoney * int64(f)
								gMoney += win
							case slot.PayTypeChipXRate:
								f, _ := strconv.ParseFloat(attr.GetGoldballValue(), 64)
								win := selection.BetMoney * selection.BetRate * int64(f)
								gMoney += win
							}
						}
					}
				}
			}
			selection.Data["GMoney"] = gMoney

			logic.SelectChances = append(logic.SelectChances, selection)

			if logic.IsFreeSpin() {
				logic.GameLog.AddGameStep("因目前正於免費遊戲中, 將於免費遊戲結束後開始選擇, 剩餘 %d 次", len(logic.SelectChances))
			}

			logic.RecordSelectGameTimes(logic.GameUser)

		case slot.PayTypeAddTimes:

			if logic.FreeSpin == nil {
				mlog.Error("增加免費遊戲次數錯誤, 沒有免費遊戲")
				continue
			}

			f, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)

			logic.FreeSpin.Count += int(f)

			logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得免費遊戲次數, 增加 %d 次, 剩餘 %d 次", result.Rule.GetSymbol(), result.Rule.GetLinkType(), result.Rule.GetCount(), int(f), logic.FreeSpin.Count)

		}

	}

}
