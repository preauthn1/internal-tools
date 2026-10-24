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

// AppendJackpotGame JackpotGame
func AppendJackpotGame(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser, detail *netproto.Slot_MSG_BonusGameRecord) {

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
		Data:     proto.String("{\"Level\":\"Grand\",\"Win\":20002}"),
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
		Data:     proto.String("{\"Level\":\"Grand\",\"Win\":20002}"),
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
		Data:     proto.String("{\"Level\":\"Minor\",\"Win\":20002}"),
	})

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(4),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"Level\":\"Minor\",\"Win\":20002}"),
	})

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(5),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"Level\":\"Mini\",\"Win\":20002}"),
	})

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(6),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"Level\":\"Major\",\"Win\":20002}"),
	})

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(7),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"Level\":\"Mini\",\"Win\":20002}"),
	})

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(8),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"Level\":\"Major\",\"Win\":20002}"),
	})

	detail.Detail = append(detail.Detail, &netproto.Slot_BonusGameDetail{
		UserID:   proto.Int32(user.UserID),
		GameID:   proto.Int32(logic.TableManager.Room.GameInfo.GameID),
		FlagID:   proto.Int32(int32(f)),
		OriginID: proto.Int32(rule.GetId()),
		TurnID:   proto.Int32(9),
		Location: proto.Int32(0),
		BetMoney: proto.Int32(5),
		BetRate:  proto.Int32(60),
		Data:     proto.String("{\"Level\":\"Mini\",\"Win\":20002}"),
	})
}
