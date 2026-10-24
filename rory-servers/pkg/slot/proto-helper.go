package slot

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// FilterProtoDataByPattern 過濾指定Pattern
func FilterProtoDataByPattern(pd *protodata.ProtoData, pattern int32) (*protodata.ProtoData, error) {

	new := new(protodata.ProtoData)

	defer CompatiblePayTable(new)

	var SlotPattern *protodata.SlotPattern

	// SlotPattern
	for _, row := range pd.SlotPattern {
		if int32(row.GetId()) == pattern {
			SlotPattern = row
			new.SlotPattern = append(new.SlotPattern, row)
			break
		}
	}

	if SlotPattern == nil {
		return nil, fmt.Errorf("Could not find pattern %d", pattern)
	}

	new.PatternID = int32(SlotPattern.GetId())

	// 整理可能會有的 SymbolGroupID
	SymbolGroupIDMap := make(map[uint32]bool)

	// 整理可能會有的 GoldReelGroupID
	GoldReelGroupIDMap := make(map[uint32]bool)

	// 整理可能會有的 PickGameID
	PickGameIDMap := make(map[uint32]bool)

	// 整理可能會有的 BonusGroupID
	BonusGroupIDMap := make(map[uint32]bool)

	// 整理可能會有的 RespinGroupID
	RespinGroupIDMap := make(map[int32]bool)

	// 整理可能會有的 FreeSpinGroupID
	FreeSpinGroupIDMap := make(map[int32]bool)

	// 整理可能會有的 SelectionGroupID
	SelectionGroupIDMap := make(map[int32]bool)

	// MainGameTable
	for _, row := range pd.MainGameTable {
		if row.GetGroupId() != SlotPattern.GetMainGameTableGroupId() {
			continue
		}
		new.MainGameTable = append(new.MainGameTable, row)

		SymbolGroupIDMap[row.GetSymbolGroupId()] = true
		GoldReelGroupIDMap[uint32(row.GetGoldReelGroupId())] = true
	}

	// SymbolAttributes
	for _, row := range pd.SymbolAttributes {
		if row.GetGroupId() != SlotPattern.GetSymbolAttributesGroupId() {
			continue
		}
		new.SymbolAttributes = append(new.SymbolAttributes, row)
	}

	// SlotPayTable
	for _, row := range pd.SlotPayTable {
		if uint32(row.GetGroupId()) != SlotPattern.GetPayTableGroupId() {
			continue
		}
		new.SlotPayTable = append(new.SlotPayTable, row)

		switch PayType(row.GetPayType()) {
		case PayTypeRespin:
			id, _ := strconv.ParseFloat(row.GetPayout(), 64)
			RespinGroupIDMap[int32(id)] = true
		case PayTypePickGame:
			id, _ := strconv.ParseFloat(row.GetPayout(), 64)
			PickGameIDMap[uint32(id)] = true
		case PayTypeFreeSpin:
			id, _ := strconv.ParseFloat(row.GetPayout(), 64)
			FreeSpinGroupIDMap[int32(id)] = true
		case PayTypeSelectGameMode:
			id, _ := strconv.ParseFloat(row.GetPayout(), 64)
			FreeSpinGroupIDMap[int32(id)] = true
		case PayTypeLightning:
			id, _ := strconv.ParseFloat(row.GetPayout(), 64)
			RespinGroupIDMap[int32(id)] = true
		case PayTypeSelection:
			id, _ := strconv.ParseFloat(row.GetPayout(), 64)
			SelectionGroupIDMap[int32(id)] = true
		case PayTypeBonus:
			id, _ := strconv.ParseFloat(row.GetPayout(), 64)
			BonusGroupIDMap[uint32(id)] = true
		}
	}

	// SlotPayLine
	if SlotPattern.GetLineType() == PayLineModeFixed.String() {
		for _, row := range pd.SlotPayLine {
			if row.GetGroupId() != SlotPattern.GetPayLineGroupId() {
				continue
			}
			new.SlotPayLine = append(new.SlotPayLine, row)
		}
	}

	// SelectionTable
	for _, row := range pd.SelectionTable {
		if _, ok := SelectionGroupIDMap[int32(row.GetGroupId())]; !ok {
			continue
		}

		new.SelectionTable = append(new.SelectionTable, row)

		switch PayType(row.GetSelectionType()) {
		case PayTypeRespin:
			RespinGroupIDMap[int32(row.GetSelectionValue())] = true
		case PayTypeFreeSpin:
			FreeSpinGroupIDMap[int32(row.GetSelectionValue())] = true
		case PayTypeLightning:
			RespinGroupIDMap[int32(row.GetSelectionValue())] = true
		case PayTypeBonus:
			BonusGroupIDMap[row.GetSelectionValue()] = true
		}
	}

	// RespinTable
	for _, row := range pd.RespinTable {
		if _, ok := RespinGroupIDMap[row.GetGroupId()]; !ok {
			continue
		}
		new.RespinTable = append(new.RespinTable, row)

		SymbolGroupIDMap[row.GetRespinValue()] = true

		if row.GetIsLimit() == 1 {
			if row.GetLimitType() == "changeDistribution" {
				strs := strings.Split(row.GetLimitValue(), ":")
				if len(strs) != 2 {
					continue
				}
				f, _ := strconv.ParseFloat(strs[1], 64)
				SymbolGroupIDMap[uint32(f)] = true
			}
		}

	}

	// JackpotTable
	for _, row := range pd.JackpotTable {
		if row.GetGroupId() != int32(SlotPattern.GetJackpotGroup()) {
			continue
		}
		new.JackpotTable = append(new.JackpotTable, row)
	}

	// PickGameTable
	for _, row := range pd.PickGameTable {
		if _, ok := PickGameIDMap[row.GetId()]; !ok {
			continue
		}

		new.PickGameTable = append(new.PickGameTable, row)

		BonusGroupIDMap[row.GetRewardValue()] = true
	}

	// BonusTable
	for _, row := range pd.BonusTable {
		if _, ok := BonusGroupIDMap[uint32(row.GetGroupId())]; !ok {
			continue
		}

		new.BonusTable = append(new.BonusTable, row)

		switch PayType(row.GetRewardType()) {
		case PayTypeFreeSpin:
			id, _ := strconv.ParseFloat(row.GetRewardValue(), 64)
			FreeSpinGroupIDMap[int32(id)] = true
		case PayTypeLightning:
			id, _ := strconv.ParseFloat(row.GetRewardValue(), 64)
			RespinGroupIDMap[int32(id)] = true
		}
	}

	// FreeSpinTable
	for _, row := range pd.FreeSpinTable {
		if _, ok := FreeSpinGroupIDMap[row.GetGroupId()]; !ok {
			continue
		}

		new.FreeSpinTable = append(new.FreeSpinTable, row)

		SymbolGroupIDMap[row.GetFreeSymbolGroupId()] = true

		if row.GetIsGoldReel() == 1 {
			GoldReelGroupIDMap[uint32(row.GetGoldReelGroupId())] = true
		}
	}

	// SlotSymbolDistribution
	for _, row := range pd.SlotSymbolDistribution {
		if _, ok := SymbolGroupIDMap[row.GetGroupId()]; !ok {
			continue
		}
		new.SlotSymbolDistribution = append(new.SlotSymbolDistribution, row)
	}

	// ReelBgDistribution
	for _, row := range pd.ReelBgDistribution {
		if _, ok := GoldReelGroupIDMap[row.GetGroupId()]; !ok {
			continue
		}
		new.ReelBgDistribution = append(new.ReelBgDistribution, row)
	}

	// ExtraRule
	for _, row := range pd.ExtraRule {
		if row.GetGroupId() != int32(SlotPattern.GetExtraRuleGroup()) {
			continue
		}
		new.ExtraRule = append(new.ExtraRule, row)
	}

	return new, nil
}

// CompatiblePayTable 處理舊賠付表規則相容
func CompatiblePayTable(p *protodata.ProtoData) {

	// PayTypeSelectGameMode
	// 把payout指定的FreeSpinGroupID, 建立成selection, 再將原始的PayTable改成selection
	{
		selectionMap := make(map[uint32][]*protodata.SlotPayTable)
		for _, r := range p.SlotPayTable {
			switch PayType(r.GetPayType()) {
			case PayTypeSelectGameMode:
				f, _ := strconv.ParseFloat(r.GetPayout(), 64)
				selectionMap[uint32(f)] = append(selectionMap[uint32(f)], r)
			}
		}

		for fgGroup, rules := range selectionMap {
			selection := &protodata.SelectionTable{
				Id:             proto.Uint32(uint32(len(p.SelectionTable) + 1)),
				GroupId:        proto.Int32(int32(len(p.SelectionTable) + 1)),
				SelectionType:  proto.String(string(GameModeFreeSpin)),
				SelectionValue: proto.Uint32(fgGroup),
			}
			p.SelectionTable = append(p.SelectionTable, selection)

			for _, r := range rules {
				r.PayType = proto.String(string(PayTypeSelection))
				r.Payout = proto.String(strconv.Itoa(int(selection.GetGroupId())))
			}
		}
	}

}
