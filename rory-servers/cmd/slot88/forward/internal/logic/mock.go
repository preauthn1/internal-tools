package logic

import (
	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/prob"
)

func (logic *Logic) genFackResult(skipSymbol ...string) *netproto.Slot_MSG_GameResult {
	if logic.Setting == nil {
		return nil
	}

	res := new(netproto.Slot_MSG_GameResult)

	skipMap := make(map[string]bool)
	for _, s := range skipSymbol {
		skipMap[s] = true
	}

	for _, set := range logic.Setting.ReelSetting {
		reel := set.GetReelID()
		for row := int32(1); row <= logic.Setting.GetRowCounts(); row++ {
			var symbol string
			rand := new(prob.Random)
			for _, s := range set.Symbols {
				if _, ok := skipMap[s.GetSymbolName()]; ok {
					continue
				}
				rand.Add(s.GetSymbolName(), 1)
			}
			symbol = rand.Choose().(string)

			r := new(netproto.Slot_ReelResult)
			r.Position = &netproto.Slot_ReelPosition{Reel: proto.Int32(reel), Row: proto.Int32(row)}
			r.Symbol = &netproto.Slot_Symbol{Position: r.Position, SymbolName: proto.String(symbol)}
			res.ReelResult = append(res.ReelResult, r)
		}
	}

	return res
}
