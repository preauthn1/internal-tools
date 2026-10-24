package wlzbjp

import (
	"sort"
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

	// 檢查盤面是否需要整理 AfterSpin
	logic.HandleAfterSpin(manager)

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

	// 免費遊戲檢查是否有龍加倍
	if logic.FreeSpin != nil {

		var multi int64 = 1

		// 先檢查盤面是否出現百搭符號，決定加倍數量
		type MultiOpt struct {
			Value  int64
			Weight int
		}
		var opts []MultiOpt

		// 決定倍率選項
		{
			str := logic.FreeSpin.Origin.GetMultiplierValue()
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

				logic.GameLog.AddGameStep("免費遊戲百搭加倍選項 %v 倍, 骰出 %v 倍", multis, t.Value)

				multi = t.Value
			}
		}

		// 結算
		for _, result := range logic.GamePayResult {

			if !result.IsMatch || result.Rule == nil {
				continue
			}

			switch result.GetPayType() {
			case slot.PayTypeChip, slot.PayTypeChipXRate:
				// 計算贏得的金額
				result.Multi = int32(multi)
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

			logic.SendIncreaseGameTimesMessage(logic.GameUser, string(slot.GameModeJackpotGame), 0, 1, "")

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
