package aladdin

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(stepno int, isRespin bool, symbolGroupID int32, prev *StepResult) {

	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)

	step := NewStepResult(stepno, isRespin, logic.GameBetMoney, logic.GamePayMoney, logic.Proto.SymbolAttributes, manager, prev)

	// 紀錄開獎符號
	step.UnexpandResult.RowMode()
	logic.GameLog.AddGameStep("第 %d 次輪軸結果(%d) %s", step.Order, step.UnexpandResult.SymbolGroupID, step.UnexpandResult.ToSymbolString())

	if step.IsExpand {
		step.ExpandResult.RowMode()
		logic.GameLog.AddGameStep("出現拓展符號: %s, 拓展後 %s", step.ExpandSymbol, step.ExpandResult.ToSymbolString())
	}

	// 計算遊戲結果
	logic.CalcGameResult(step)

	// 計算遊戲獎勵
	logic.CalcGameAward(step)

	logic.StepResults = append(logic.StepResults, step)

	sort.Slice(logic.StepResults, func(i, j int) bool {
		return logic.StepResults[i].Order < logic.StepResults[j].Order
	})
}

// CalcGameResult 計算遊戲結果
func (logic *Logic) CalcGameResult(step *StepResult) {

	step.GamePayResult = logic.SlotLogicBase.CalcGameResult(step.ExpandResult)

}

// CalcGameAward 計算遊戲獎勵
func (logic *Logic) CalcGameAward(step *StepResult) {

	step.GameWinMoney = logic.SlotLogicBase.CalcGameAward(step.GamePayResult, step.GameBetMoney, 1)

	if !conf.AnalysisMode {
		time.Sleep(100 * time.Millisecond)
	}

	logic.AddSlotGameHistoryLogByStep(logic.GameUser, step)

	// 寫分
	logic.UpdateUserScore(logic.GameUser, step)

	logic.GameWinMoney += step.GameWinMoney

	for _, result := range step.GamePayResult {

		if !result.IsMatch || result.Rule == nil {
			continue
		}

		// 獲得再轉一次
		if result.GetPayType() == slot.PayTypeRespin {

			if len(step.ExpandReel) == 0 {
				continue
			}

			// 找出RespinTable對應的SymbolGroupID
			f, _ := strconv.ParseFloat(result.Rule.GetPayout(), 64)
			respinGroupID := int32(f)
			var respinTable *protodata.RespinTable

			for _, row := range logic.Proto.RespinTable {
				if row.GetGroupId() != respinGroupID {
					continue
				}

				var reels []int32

				reelstrs := strings.Split(row.GetAppearReel(), ",")
				for _, reelstr := range reelstrs {
					f, _ := strconv.ParseFloat(reelstr, 64)
					reels = append(reels, int32(f))
				}

				sort.Slice(reels, func(i, j int) bool {
					return reels[i] < reels[j]
				})

				if reflect.DeepEqual(reels, step.AllExpandReels) {
					respinTable = row
					break
				}
			}

			if respinTable != nil {
				logic.Spin(step.Order+1, true, int32(respinTable.GetRespinValue()), step)
			}

			continue
		}

	}

}
