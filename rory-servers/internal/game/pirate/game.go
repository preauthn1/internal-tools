package pirate

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
// 	"example.invalid/ServerTeam/Servers/pkg/mlog"
// 	"example.invalid/ServerTeam/Servers/pkg/slot"
// 	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
// )

// // Bet 判斷下注是否
// func (logic *Logic) Bet(user *gamelogic.GameUser, money int64) (ok bool, retcode int32, betmoney int64) {

// 	logic.GameBetMoney = money
// 	logic.GamePayMoney = money * int64(len(logic.SlotPayLine))

// 	if logic.FreeSpin != nil && logic.FreeSpin.Count > 0 {
// 		logic.GameBetMoney = logic.FreeSpin.BetMoney
// 		logic.GamePayMoney = 0
// 		logic.FreeSpin.Count--
// 	}

// 	// 檢查狀態是否可下注
// 	if logic.GameState != slot.GameStateBet {
// 		return false, slot.BetStateStepError, logic.GamePayMoney
// 	}

// 	// 檢查玩家金額
// 	if user.Score < logic.GamePayMoney {
// 		return false, slot.BetStateMoneyNotEnough, logic.GamePayMoney
// 	}

// 	// 檢查最低攜帶金額
// 	if user.Score < logic.MinCarryMoney {
// 		return false, slot.BetStateMoneyLessMin, logic.GamePayMoney
// 	}

// 	return true, slot.BetStateSuccess, logic.GamePayMoney
// }

// // Spin 玩家轉輪軸
// func (logic *Logic) Spin(symbolGroupID int32) {

// 	// 產生遊戲結果
// 	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
// 	manager.GenerateGameResult(symbolGroupID, logic.SlotSymbolDistribution)

// 	logic.GameResultManager = manager

// 	// 紀錄開獎符號
// 	logic.GameLog.AddGameStep("輪軸結果, 使用SymbolGroupID: %d", logic.GameResultManager.SymbolGroupID)
// 	mlog.Debug("輪軸結果, 使用SymbolGroupID: %d", logic.GameResultManager.SymbolGroupID)

// 	logic.GameResultManager.RowMode()
// 	for i, row := range logic.GameResultManager.GetSymbolString() {
// 		logic.GameLog.AddGameStep("第 %d 行: %s", i+1, strings.Join(row, ", "))
// 		mlog.Debug("第 %d 行: %s", i+1, strings.Join(row, ", "))
// 	}

// 	logic.CalcGameResult()

// 	logic.CalcGameAward()
// }

// // CalcGameResult 計算遊戲結果
// func (logic *Logic) CalcGameResult() {

// 	logic.GamePayResult = logic.SlotLogicBase.CalcGameResult(logic.GameResultManager)

// }

// // CalcGameAward 計算遊戲獎勵
// func (logic *Logic) CalcGameAward() {

// 	logic.GameWinMoney = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)

// 	if logic.FreeSpin != nil {

// 		multi := logic.FreeSpin.GetMultipleValue()
// 		win := logic.GameWinMoney
// 		total := win * int64(multi)

// 		logic.FreeSpin.WinCount++
// 		logic.FreeSpin.WinMoney += total
// 		logic.GameWinMoney = total

// 		logic.GameLog.AddGameStep("免費遊戲中第 %d 次中獎, 加乘 %d 倍, 原中獎金額: %d, 加乘後金額: %d", logic.FreeSpin.WinCount, multi, win, total)
// 		mlog.Debug("免費遊戲中第 %d 次中獎, 加乘 %d 倍, 原中獎金額: %d, 加乘後金額: %d", logic.FreeSpin.WinCount, multi, win, total)
// 	}

// 	for _, result := range logic.GamePayResult {

// 		if !result.IsMatch || result.Rule == nil {
// 			continue
// 		}

// 		if result.GetPayType() == slot.PayTypeJackpot {

// 			jackpot := logic.Jackpots[slot.JackpotSingle]

// 			win := jackpot.GetCurrentPoolByBet(logic.GameBetMoney, int64(logic.MaxBetCount))

// 			logic.GameWinMoney += win
// 			logic.GameWinJackpot = win

// 			jackpot.AwardTotal += win

// 			var paylines []string
// 			for _, line := range result.SymbolPayLines {
// 				paylines = append(paylines, strconv.Itoa(int(line.PayLine.Origin.GetLineId()))+" : "+line.PayLine.Origin.GetContent())
// 			}
// 			logic.GameLog.AddGameStep("贏得Jackpot; 金額: %d; 連線賠線: %s", win, strings.Join(paylines, " | "))
// 			mlog.Debug("贏得Jackpot; 金額: %d; 連線賠線: %s", win, strings.Join(paylines, " | "))

// 			if jackpot.IsNeedReset() {
// 				logic.UpdateJackpotPoolNoLink(0, -jackpot.TotalContribution, -jackpot.AwardTotal+win, jackpot.Level)
// 				jackpot.Reset()
// 			} else {
// 				logic.UpdateJackpotPoolNoLink(0, 0, win, jackpot.Level)
// 			}

// 			continue
// 		}

// 		if result.GetPayType() == slot.PayTypeFreeSpin {
// 			// 取得贏得FreeSpin的GroupID
// 			f, err := strconv.ParseFloat(result.Rule.GetPayout(), 64)
// 			if err != nil {
// 				mlog.Warn("SlotPayTable.Payout解析錯誤，%v", err)
// 				continue
// 			}
// 			freeSpinGroupID := int32(f)

// 			var freeSpinTable *protodata.FreeSpinTable
// 			// 取得FreeSpinTable
// 			for _, table := range logic.FreeSpinTable {
// 				if table.GetGroupId() == freeSpinGroupID {
// 					freeSpinTable = table
// 					break
// 				}
// 			}

// 			if freeSpinTable == nil {
// 				mlog.Warn("查無FreeSpinTable, FreeSpinGroupID: %d", freeSpinGroupID)
// 				continue
// 			}

// 			// 檢查本身是否是FreeSpin
// 			if logic.FreeSpin != nil {
// 				if logic.FreeSpin.Origin != freeSpinTable {
// 					mlog.Warn("贏得的FreeSpin與原本的FreeSpin不同")
// 					continue
// 				}

// 				logic.FreeSpin.Count += int(freeSpinTable.GetFreeSpinCount())

// 			} else {
// 				logic.FreeSpin = slot.NewFreeSpin(freeSpinTable, logic.SlotSymbolDistribution, logic.GameBetMoney, 1)
// 			}

// 			fmt.Println(logic.FreeSpin)

// 			logic.GameLog.AddGameStep("贏得FreeSpin, ID: %d, SymbolGroupID: %d, Count: %d, 剩餘次數: %d", logic.FreeSpin.Origin.GetId(), logic.FreeSpin.Origin.GetFreeSymbolGroupId(), logic.FreeSpin.Origin.GetFreeSpinCount(), logic.FreeSpin.Count)
// 			mlog.Debug("贏得FreeSpin, ID: %d, SymbolGroupID: %d, Count: %d, 剩餘次數: %d", logic.FreeSpin.Origin.GetId(), logic.FreeSpin.Origin.GetFreeSymbolGroupId(), logic.FreeSpin.Origin.GetFreeSpinCount(), logic.FreeSpin.Count)

// 			continue
// 		}

// 	}

// }
