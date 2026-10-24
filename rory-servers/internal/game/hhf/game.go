package hhf

import (
	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID int32) {
	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)

	manager.RegisterAfterGenerateFunc(func(m *slot.GameResultManager) {
		if logic.FreeSpin == nil {
			return
		}
		for _, l := range logic.locks {
			m.ReelMode()
			for _, pos := range l.positions {
				m.GameResult[int(pos.Reel-1)][int(pos.Row-1)].Symbol = proto.String(l.symbol)
			}
		}
	})

	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)

	logic.GameResultManager = manager

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

	ChipWinMoney = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)

	if logic.IsFreeSpin() {
		HistoryType = slot.HistoryTypeFREE

		logic.FreeSpin.WinCount++
		logic.FreeSpin.WinMoney += ChipWinMoney
	} else {
		HistoryType = slot.HistoryTypeMAIN
	}

	logic.GameWinMoney = ChipWinMoney

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
		case slot.PayTypeFreeSpin:

			logic.WinFreeSpin(result.Rule, true)

			logic.UpdateFreeSpin(logic.GameUser)

			continue
		case slot.PayTypeLock:

			locksymbol := make(map[string][]slot.Position)

			logic.GameResultManager.ReelMode()
			for _, pos := range result.BlinkPositions {
				sbn := logic.GameResultManager.GameResult[int(pos.Reel)-1][int(pos.Row)-1].GetSymbol()
				locksymbol[sbn] = append(locksymbol[sbn], pos)
			}

			for sbn, positions := range locksymbol {
				var l *lock
				for _, _l := range logic.locks {
					if _l.symbol == sbn {
						l = _l
					}
				}

				if l == nil {
					l = new(lock)
					l.symbol = sbn

					logic.locks = append(logic.locks, l)
				}

				l.positions = slot.MergePositions(l.positions, positions)
			}

			continue
		}

	}
}
