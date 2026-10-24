package slot

import (
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// JackpotGameChance 玩家可玩JackpotGame的次數
type JackpotGameChance struct {
	Rule              *protodata.SlotPayTable `json:"-"`
	RuleID            int32
	BetMoney          int64
	BetRate           int64
	GameResultManager *GameResultManager `json:"-"`
	Data              map[string]interface{}
}

// NewJackpotGameChance create JackpotGameChance
func NewJackpotGameChance(result *GameResultManager, r *protodata.SlotPayTable, betMoney, betRate int64) (*JackpotGameChance, error) {
	s := new(JackpotGameChance)

	s.GameResultManager = result
	s.Rule = r
	s.RuleID = r.GetId()
	s.Data = make(map[string]interface{})

	s.BetMoney = betMoney
	s.BetRate = betRate

	return s, nil
}
