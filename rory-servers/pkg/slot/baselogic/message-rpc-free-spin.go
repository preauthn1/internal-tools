package baselogic

import (
	"encoding/json"
	"strconv"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// ReadFreeSpinTimes 讀取FreeSpin紀錄
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 30011 (Slot_GameMessageClassID_Slot_FreeGameReadID)
func (logic *SlotLogicBase) ReadFreeSpinTimes(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadID)

	data := new(netproto.Slot_MSG_FreeGameRead)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// UpdateFreeSpin FreeSpin存入DB
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 30013 (Slot_GameMessageClassID_Slot_FreeGameRecordID)
func (logic *SlotLogicBase) UpdateFreeSpin(user *gamelogic.GameUser) {

	if logic.FreeSpin == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_FreeGameRecordID)

	data := new(netproto.Slot_MSG_FreeGameRecord)
	example.invalid = new(netproto.Slot_FreeSpinDetail)
	example.invalid.UserID = proto.Int32(user.UserID)
	example.invalid.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	example.invalid.FlagID = proto.Int32(int32(FlagID))

	example.invalid.OriginID = proto.Int64(int64(logic.FreeSpin.Origin.GetId()))
	example.invalid.FreeSpinGroup = proto.Int32(logic.FreeSpin.SymbolGroupID)
	example.invalid.Times = proto.Int32(int32(logic.FreeSpin.Count))

	example.invalid.ChooseRate = proto.Int32(int32(logic.FreeSpin.BetRate))
	example.invalid.UserBet = proto.Int64(logic.FreeSpin.BetMoney)

	example.invalid.WinMoney = proto.Int64(logic.FreeSpin.WinMoney)
	example.invalid.WinTimes = proto.Int32(int32(logic.FreeSpin.WinCount))

	example.invalid.PatternID = proto.Int32(logic.Proto.PatternID)

	if logic.FreeSpin.FreeSpinExtAward != nil {
		b, _ := json.Marshal(logic.FreeSpin.FreeSpinExtAward)
		example.invalid.Data = proto.String(string(b))
	} else {
		example.invalid.Data = proto.String("")
	}

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecvFreeSpinTimes 接收db server回傳的FreeSpin資訊，並整理至Logic內
func (logic *SlotLogicBase) RecvFreeSpinTimes(m *msg.Message) {

	defer logic.CheckWaitGroup()
	mlog.Debug("Table %d, RecvFreeSpinTimes", logic.TableNo())

	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_FreeGameRecord{
			Info: &netproto.Slot_FreeSpinDetail{},
		}
		return mm
	})

	detail := rmm.(*netproto.Slot_MSG_FreeGameRecord)

	if logic.GameUser == nil {
		return
	}

	if logic.GameUser.UserID != example.invalid.GetUserID() {
		return
	}

	if logic.TableManager.Room.GameInfo.GameID != example.invalid.GetGameID() {
		return
	}

	if logic.TableManager.Room.Flag != strconv.Itoa(int(example.invalid.GetFlagID())) {
		return
	}

	if example.invalid.GetTimes() == 0 {
		return
	}

	if example.invalid.GetUserBet() == 0 || example.invalid.GetChooseRate() == 0 {
		return
	}

	// 查 Default 表
	if Proto != nil && (Proto.PatternID == example.invalid.GetPatternID() || example.invalid.GetPatternID() == 0) {
		for _, table := range Proto.FreeSpinTable {
			if int64(table.GetId()) == example.invalid.GetOriginID() {
				fg := buildFreeSpin(table, Proto, detail)
				logic.bufFeatureGames[DefaultPattern] = append(logic.bufFeatureGames[DefaultPattern], fg)
			}
		}
	}

	// 查 Release 表
	if ReleaseProto != nil && ReleaseProto.PatternID == example.invalid.GetPatternID() {
		for _, table := range ReleaseProto.FreeSpinTable {
			if int64(table.GetId()) == example.invalid.GetOriginID() {
				fg := buildFreeSpin(table, ReleaseProto, detail)
				logic.bufFeatureGames[ReleasePattern] = append(logic.bufFeatureGames[ReleasePattern], fg)
			}
		}
	}

	// 查 Bite 表
	if BiteProto != nil && BiteProto.PatternID == example.invalid.GetPatternID() {
		for _, table := range BiteProto.FreeSpinTable {
			if int64(table.GetId()) == example.invalid.GetOriginID() {
				fg := buildFreeSpin(table, BiteProto, detail)
				logic.bufFeatureGames[BitePattern] = append(logic.bufFeatureGames[BitePattern], fg)
			}
		}
	}

}

func buildFreeSpin(freeTable *protodata.FreeSpinTable, proto *protodata.ProtoData, detail *netproto.Slot_MSG_FreeGameRecord) *slot.FreeSpin {

	fg := slot.NewFreeSpin(freeTable, proto.SlotSymbolDistribution, example.invalid.GetUserBet(), int64(example.invalid.GetChooseRate()))
	fg.Count = int(example.invalid.GetTimes())
	fg.WinMoney = example.invalid.GetWinMoney()
	fg.WinCount = int(example.invalid.GetWinTimes())

	if len(example.invalid.GetData()) > 0 {
		award := new(slot.FreeSpinExtAward)
		if err := json.Unmarshal([]byte(example.invalid.GetData()), award); err != nil {
			mlog.Warn("FreeSpin 解析額外資訊錯誤, %v", err)
			return nil
		}
		fg.FreeSpinExtAward = award
	}

	return fg
}
