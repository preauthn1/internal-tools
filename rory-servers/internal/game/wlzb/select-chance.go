package wlzb

import (
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// SelectChance 玩家可選擇免費遊戲
type SelectChance struct {
	Rule     *protodata.SlotPayTable    `json:"-"`
	Options  []*protodata.FreeSpinTable `json:"-"`
	GroupID  int32
	BetMoney int64 `json:"BetMoney"`
	BetRate  int64 `json:"BetRate"`
}

// NewSelectChance create SelectChance
func NewSelectChance(r *protodata.SlotPayTable, opts []*protodata.FreeSpinTable, betMoney, betRate int64) (*SelectChance, error) {
	s := new(SelectChance)

	s.Rule = r

	f, err := strconv.ParseFloat(s.Rule.GetPayout(), 64)
	if err != nil {
		return nil, err
	}

	s.GroupID = int32(f)

	for _, opt := range opts {
		if opt.GetGroupId() != s.GroupID {
			continue
		}
		s.Options = append(s.Options, opt)
	}

	s.BetMoney = betMoney
	s.BetRate = betRate

	return s, nil
}

// ChooseByID get selection
func (sc *SelectChance) ChooseByID(id int32) *protodata.FreeSpinTable {
	for _, opt := range sc.Options {
		if opt.GetId() == uint32(id) {
			return opt
		}
	}
	return nil
}
