package slot

import (
	"fmt"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// SelectOption 選項
type SelectOption struct {
	Rule      *protodata.SelectionTable `json:"-"`
	RuleID    int32
	GameMode  GameMode
	GameTable interface{} `json:"-"`
	GameID    int32
}

func (opt *SelectOption) String() string {
	switch opt.GameMode {
	case GameModeFreeSpin:
		tb := opt.GameTable.(*protodata.FreeSpinTable)
		return fmt.Sprintf("FreeSpin(%d:%d)", tb.GetFreeSymbolGroupId(), tb.GetFreeSpinCount())
	case GameModeLightning:
		tb := opt.GameTable.(*protodata.RespinTable)
		return fmt.Sprintf("Lightning(%d:%d)", tb.GetRespinValue(), tb.GetRespinCount())
	case GameModeBonus:
		items := opt.GameTable.([]*protodata.BonusTable)
		return fmt.Sprintf("Bonus(%d)", len(items))
	}
	return ""
}

// SelectChance 玩家可選擇遊戲模式
type SelectChance struct {
	Rule              *protodata.SlotPayTable     `json:"-"`
	Selections        []*protodata.SelectionTable `json:"-"`
	GroupID           int32
	BetMoney          int64           `json:"BetMoney"`
	BetRate           int64           `json:"BetRate"`
	Options           []*SelectOption `json:"-"`
	GameResultManager *GameResultManager
	Data              map[string]interface{}
}

// NewSelectChance create SelectChance
func NewSelectChance(protos *protodata.ProtoData, result *GameResultManager, r *protodata.SlotPayTable, betMoney, betRate int64) (*SelectChance, error) {
	s := new(SelectChance)

	s.GameResultManager = result
	s.Rule = r
	s.Data = make(map[string]interface{})

	f, err := strconv.ParseFloat(s.Rule.GetPayout(), 64)
	if err != nil {
		return nil, err
	}

	s.GroupID = int32(f)

	for _, selection := range protos.SelectionTable {
		if selection.GetGroupId() != s.GroupID {
			continue
		}
		if selection.GetIsGoldsymbol() == 1 {
			switch selection.GetGoldCondition() {
			case ConditionRate:
				f, _ := strconv.ParseFloat(selection.GetConditionValue(), 64)
				if betRate != int64(f) {
					continue
				}
			}
		}
		s.Selections = append(s.Selections, selection)

		switch GameMode(selection.GetSelectionType()) {
		case GameModeFreeSpin:
			for _, row := range protos.FreeSpinTable {
				if row.GetGroupId() != int32(selection.GetSelectionValue()) {
					continue
				}
				opt := new(SelectOption)
				opt.Rule = selection
				opt.RuleID = int32(selection.GetId())
				opt.GameMode = GameModeFreeSpin
				opt.GameTable = row
				opt.GameID = int32(row.GetId())
				s.Options = append(s.Options, opt)
			}
		case GameModeLightning:
			for _, row := range protos.RespinTable {
				if row.GetGroupId() != int32(selection.GetSelectionValue()) {
					continue
				}
				opt := new(SelectOption)
				opt.Rule = selection
				opt.RuleID = int32(selection.GetId())
				opt.GameMode = GameModeLightning
				opt.GameTable = row
				opt.GameID = int32(row.GetId())
				s.Options = append(s.Options, opt)
			}
		case GameModeBonus:
			opt := new(SelectOption)
			opt.Rule = selection
			opt.RuleID = int32(selection.GetId())
			opt.GameMode = GameModeBonus

			var items []*protodata.BonusTable
			for _, row := range protos.BonusTable {
				if row.GetGroupId() != int32(selection.GetSelectionValue()) {
					continue
				}
				items = append(items, row)
			}

			if len(items) > 0 {
				opt.GameTable = items
				s.Options = append(s.Options, opt)
			}
		}
	}

	s.BetMoney = betMoney
	s.BetRate = betRate

	return s, nil
}
