package dhzw

import (
	"strconv"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID, bgGroupID int32) {
	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)
	manager.GenerateBackground(bgGroupID, logic.Proto.ReelBgDistribution)

	logic.GameResultManager = manager

	// 檢查背景是否需要整理
	logic.HandleAfterBackground(manager)

	// 紀錄開獎符號
	logic.LogBackgroundResult(logic.GameResultManager)

	// 檢查盤面是否需要整理
	logic.HandleAppearSymbol(manager)

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

	var ChipWinMoney int64
	var HistoryType int32

	if logic.IsFreeSpin() {
		HistoryType = slot.HistoryTypeFREE
		// 免費遊戲檢查是否有龍加倍
		multi, d := logic.CalcMultiMode(slot.MultiModeDiceEachReelMulti, logic.FreeSpin.Origin.GetMultiplierValue())
		ChipWinMoney = logic.SlotLogicBase.CalcGameAwardByMultiType(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate, multi, slot.MultiTypeApearWild, d)

		if ChipWinMoney > 0 {
			logic.FreeSpin.WinCount++
			logic.FreeSpin.WinMoney += ChipWinMoney
		}
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
		case slot.PayTypeMysteryJPG, slot.PayTypeMysteryJPGR:

			// 在免費遊戲中僅給予機會
			if logic.IsFreeSpin() {

				c, err := slot.NewJackpotGameChance(logic.GameResultManager, result.Rule, logic.GameBetMoney, logic.GameBetRate)
				if err != nil {
					mlog.Warn("建立 Jackpot Game 機會錯誤 %v", err)
					continue
				}

				logic.JackpotGameChances = append(logic.JackpotGameChances, c)

				logic.SendIncreaseGameTimesMessage(logic.GameUser, string(slot.GameModeJackpotGame), 0, 1, "")

				logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得 JackpotGame, 正在免費遊戲中, 結束後將決定獎項, 剩餘 %d 次", result.SymbolName, result.Rule.GetLinkType(), result.PayCount, len(logic.JackpotGameChances))

				logic.RecordJackpotGameChance(logic.GameUser, c)

				continue
			}

			jpg := logic.WinMysteryJPG(result.Rule)

			if jpg == nil {
				continue
			}

			logic.JackpotGame = jpg

			logic.RecordJackpotGame(logic.GameUser)

			logic.SendIncreaseGameTimesMessage(logic.GameUser, string(slot.GameModeJackpotGame), 0, 1, "")

		case slot.PayTypeSelection:

			selection := logic.WinSelection(result.Rule)

			if selection == nil {
				continue
			}

			logic.SelectChances = append(logic.SelectChances, selection)

			logic.RecordSelectGameTimes(logic.GameUser)

		case slot.PayTypeAddTimes:

			if logic.FreeSpin == nil {
				mlog.Error("增加免費遊戲次數錯誤, 沒有免費遊戲")
				continue
			}

			f, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)

			logic.FreeSpin.Count += int(f)

			logic.GameLog.AddGameStep("開獎符號: %s:%s:%d, 贏得免費遊戲次數, 增加 %d 次, 剩餘 %d 次", result.SymbolName, result.Rule.GetLinkType(), result.PayCount, int(f), logic.FreeSpin.Count)

			logic.SendIncreaseGameTimesMessage(logic.GameUser, string(slot.GameModeFreeSpin), int32(logic.FreeSpin.Origin.GetId()), int32(f), "")

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
		}
	}
}
