package mock

import (
	"strconv"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// AppendJackpotGameChance JackpotGameChance
func AppendJackpotGameChance(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser, detail *netproto.Slot_MSG_DeferredAwardRecord) {

	var rule *protodata.SlotPayTable
	for _, row := range logic.Proto.SlotPayTable {
		switch slot.PayType(row.GetPayType()) {
		default:
		case slot.PayTypeMysteryJackpot,
			slot.PayTypeMysteryJackpotR,
			slot.PayTypeMysteryJPG,
			slot.PayTypeMysteryJPGR:
			rule = row
			break
		}
	}
	if rule == nil {
		return
	}

	detail.Records = append(detail.Records, &netproto.Slot_DeferredAward{
		AwardType: proto.String(string(slot.GameModeJackpotGame)),
		OriginID:  proto.Int32(rule.GetId()),
		Data:      proto.String("{\"RuleID\":" + strconv.Itoa(int(rule.GetId())) + ",\"BetMoney\":1,\"BetRate\":28,\"Data\":{}}"),
	})

}
