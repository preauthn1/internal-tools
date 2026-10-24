package baselogic

import (
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// IsFreeSpin 檢查目前是否為FreeSpin狀態
func (logic *SlotLogicBase) IsFreeSpin() bool {
	if logic.FreeSpin != nil && logic.FreeSpin.Count >= 0 && logic.FreeSpin.Used > 0 {
		return true
	}

	return false
}

// CheckFreeGame 檢查 FreeSpin, 檢查免費遊戲是否結束沒次數
func (logic *SlotLogicBase) CheckFreeGame() {
	if logic.FreeSpin == nil {
		return
	}

	if logic.FreeSpin.Count > 0 {
		return
	}

	logic.SendLeftFreeGameToUser(logic.GameUser)
	logic.FreeSpin = nil
}

// CheckJackpotGameChances 檢查JackpotGameChances, 若有次數且沒有在免費遊戲狀態下, 建立JackpotGame
func (logic *SlotLogicBase) CheckJackpotGameChances() {

	if len(logic.JackpotGameChances) == 0 {
		return
	}

	if logic.FreeSpin != nil {
		return
	}

	if logic.JackpotGame != nil {
		return
	}

	chance := logic.JackpotGameChances[0]
	logic.GameBetMoney = chance.BetMoney
	logic.GameBetRate = chance.BetRate

	logic.PopJackpotGameChance(logic.GameUser)
	logic.JackpotGameChances = append([]*slot.JackpotGameChance(nil), logic.JackpotGameChances[1:]...)

	logic.GameLog.AddGameStep("自動開啟 Jackpot 翻牌遊戲, 底注: %d, 倍率: %d, 剩餘 %d 次", chance.BetMoney, chance.BetRate, len(logic.JackpotGameChances))

	switch slot.PayType(chance.Rule.GetPayType()) {
	case slot.PayTypeMysteryJPG, slot.PayTypeMysteryJPGR:
		logic.JackpotGame = logic.WinMysteryJPG(chance.Rule)
		logic.RecordJackpotGame(logic.GameUser)
	}

	// logic.SendJackpotGameDetailMessage(logic.GameUser)
}
