package baselogic

import (
	"fmt"
	"strings"
)

// IncreaseJackpotPool 累積Jackpot獎池
func (logic *SlotLogicBase) IncreaseJackpotPool() {

	var logs []string

	for _, jackpot := range logic.Jackpots {
		if jackpot.Origin.GetIsProgress() != 1 {
			continue
		}

		inc := jackpot.GetIncrease(logic.GamePayMoney)
		// logs = append(logs, fmt.Sprintf("N=%d,AwardAdd=%d,jpAmount=%d", inc, jackpot.TotalContribution, jackpot.GetCurrentPool()))

		// 手動更新獎池資訊
		jackpot.Increase(logic.GamePayMoney)
		logs = append(logs, fmt.Sprintf("[%s] %d:%d:%d", jackpot.Level.String(), jackpot.GetCurrentPool(), inc, jackpot.AwardTotal))

		// 更新Jackpot
		logic.UpdateJackpotPoolNoLink(0, inc, 0, jackpot.Level)
	}

	if len(logs) > 0 {
		logic.GameLog.AddGameStep("水池: %s", strings.Join(logs, "|"))
	}
}
