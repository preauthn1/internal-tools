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

// AppendPickGame PickGame
func AppendPickGame(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser, detail *netproto.Slot_MSG_BonusGameRecord) {

	var rule *protodata.SlotPayTable
	for _, row := range logic.Proto.SlotPayTable {
		switch slot.PayType(row.GetPayType()) {
		default:
		case slot.PayTypePickGame:
			rule = row
			break
		}
	}

	if rule == nil {
		return
	}

	f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(1),
		Location: proto.Int32(5),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"AllPosition\":[{\"Reel\":2,\"Row\":1},{\"Reel\":3,\"Row\":1},{\"Reel\":4,\"Row\":1},{\"Reel\":2,\"Row\":2},{\"Reel\":3,\"Row\":2},{\"Reel\":4,\"Row\":2},{\"Reel\":2,\"Row\":3},{\"Reel\":3,\"Row\":3},{\"Reel\":4,\"Row\":3}],\"Bonus\":52,\"Position\":{\"Reel\":3,\"Row\":2},\"TotalWin\":9320}"),
	})

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(2),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"AllPosition\":[{\"Reel\":2,\"Row\":1},{\"Reel\":3,\"Row\":1},{\"Reel\":4,\"Row\":1},{\"Reel\":2,\"Row\":2},{\"Reel\":3,\"Row\":2},{\"Reel\":4,\"Row\":2},{\"Reel\":2,\"Row\":3},{\"Reel\":3,\"Row\":3},{\"Reel\":4,\"Row\":3}],\"Bonus\":52,\"Position\":{\"Reel\":0,\"Row\":0},\"TotalWin\":9320}"),
	})

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(3),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"AllPosition\":[{\"Reel\":2,\"Row\":1},{\"Reel\":3,\"Row\":1},{\"Reel\":4,\"Row\":1},{\"Reel\":2,\"Row\":2},{\"Reel\":3,\"Row\":2},{\"Reel\":4,\"Row\":2},{\"Reel\":2,\"Row\":3},{\"Reel\":3,\"Row\":3},{\"Reel\":4,\"Row\":3}],\"Bonus\":49,\"Position\":{\"Reel\":0,\"Row\":0},\"TotalWin\":9320}"),
	})
}

// AppendPickGame2 PickGame
func AppendPickGame2(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser, detail *netproto.Slot_MSG_BonusGameRecord) {

	var rule *protodata.SlotPayTable
	for _, row := range logic.Proto.SlotPayTable {
		switch slot.PayType(row.GetPayType()) {
		default:
		case slot.PayTypePickGame:
			rule = row
			break
		}
	}

	if rule == nil {
		return
	}

	f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(1),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String(`{"AllPosition":[{"Reel":2,"Row":1},{"Reel":4,"Row":1},{"Reel":4,"Row":2}],"Bonus":47,"Position":{"Reel":0,"Row":0},"TotalWin":9440}`),
	})
}
