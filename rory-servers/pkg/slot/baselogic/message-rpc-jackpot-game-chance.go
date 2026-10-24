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
)

// ReadJackpotGameChance 讀取 JackpotGameChance
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000061 (Slot_GameMessageClassID_Slot_DeferredAwardReadID)
func (logic *SlotLogicBase) ReadJackpotGameChance(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID)

	data := new(netproto.Slot_MSG_DeferredAwardRead)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.Atoi(logic.TableManager.Room.Flag)
	data.FlagID = proto.Int32(int32(FlagID))
	data.AwardType = proto.String(string(slot.GameModeJackpotGame))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecordJackpotGameChance 紀錄 JackpotGameChance
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000063 (Slot_GameMessageClassID_Slot_DeferredAwardRecordID)
func (logic *SlotLogicBase) RecordJackpotGameChance(user *gamelogic.GameUser, chances ...*slot.JackpotGameChance) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardRecordID)

	data := new(netproto.Slot_MSG_DeferredAwardRecord)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.Atoi(logic.TableManager.Room.Flag)
	data.FlagID = proto.Int32(int32(FlagID))

	for _, chance := range chances {
		record := new(netproto.Slot_DeferredAward)
		record.PatternID = proto.Int32(logic.Proto.PatternID)
		record.AwardType = proto.String(string(slot.GameModeJackpotGame))
		record.OriginID = proto.Int32(chance.RuleID)
		b, err := json.Marshal(chance)
		if err != nil {
			mlog.Debug("RecordJackpotGameChance Marshal JSON Error, %v", err)
			continue
		}
		record.Data = proto.String(string(b))

		data.Records = append(data.Records, record)
	}

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// PopJackpotGameChance 使用 JackpotGameChance
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000064 (Slot_GameMessageClassID_Slot_DeferredAwardPopID)
func (logic *SlotLogicBase) PopJackpotGameChance(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardPopID)

	data := new(netproto.Slot_MSG_DeferredAwardRead)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.Atoi(logic.TableManager.Room.Flag)
	data.FlagID = proto.Int32(int32(FlagID))
	data.AwardType = proto.String(string(slot.GameModeJackpotGame))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// ResetJackpotGameChance 重置 JackpotGameChance
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000064 (Slot_GameMessageClassID_Slot_DeferredAwardPopID)
func (logic *SlotLogicBase) ResetJackpotGameChance(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardResetID)

	data := new(netproto.Slot_MSG_DeferredAwardRead)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.Atoi(logic.TableManager.Room.Flag)
	data.FlagID = proto.Int32(int32(FlagID))
	data.AwardType = proto.String(string(slot.GameModeJackpotGame))

	pmsg.MsgData = data

	logic.JackpotGameChances = make([]*slot.JackpotGameChance, 0)

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecvJackpotGameChance 接收 JackpotGameChance 資料
func (logic *SlotLogicBase) RecvJackpotGameChance(m *msg.Message) {
	defer logic.CheckWaitGroup()
	mlog.Debug("Table %d, RecvJackpotGameChance", logic.TableNo())

	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_DeferredAwardRecord{}
		return mm
	})

	rec := rmm.(*netproto.Slot_MSG_DeferredAwardRecord)

	if logic.GameUser == nil {
		return
	}

	if rec.GetUserID() != logic.GameUser.UserID {
		mlog.Debug("Deferred Award UserID %d is not CurrentUser %d, Skip Message", rec.GetUserID(), logic.GameUser.UserID)
		return
	}

	if rec.GetGameID() != logic.TableManager.Room.GameInfo.GameID {
		mlog.Debug("Deferred Award GameID %d is not this Game %d, Skip Message", rec.GetGameID(), logic.TableManager.Room.GameInfo.GameID)
		return
	}

	FlagID, _ := strconv.Atoi(logic.TableManager.Room.Flag)
	if rec.GetFlagID() != int32(FlagID) {
		mlog.Debug("Deferred Award FlagID %d is not this Room %d, Skip Message", rec.GetFlagID(), FlagID)
		return
	}

	if len(rec.GetRecords()) == 0 {
		mlog.Debug("Deferred Award has no Record, Skip Message")
		return
	}

	patternMap := make(map[PatternMode][]*slot.JackpotGameChance)

	for _, record := range rec.GetRecords() {
		if slot.GameMode(record.GetAwardType()) != slot.GameModeJackpotGame {
			continue
		}

		chance := new(slot.JackpotGameChance)

		if err := json.Unmarshal([]byte(record.GetData()), chance); err != nil {
			mlog.Warn("Parse JackpotGame Chance Fail. %v", err)
			continue
		}

		// 查 Default 表
		if Proto != nil && (Proto.PatternID == record.GetPatternID() || record.GetPatternID() == 0) {
			for _, table := range Proto.SlotPayTable {
				if table.GetId() == chance.RuleID {
					chance.Rule = table
					patternMap[DefaultPattern] = append(patternMap[DefaultPattern], chance)
				}
			}
		}

		// 查 Release 表
		if ReleaseProto != nil && ReleaseProto.PatternID == record.GetPatternID() {
			for _, table := range ReleaseProto.SlotPayTable {
				if table.GetId() == chance.RuleID {
					chance.Rule = table
					patternMap[ReleasePattern] = append(patternMap[ReleasePattern], chance)
				}
			}
		}

		// 查 Bite 表
		if BiteProto != nil && BiteProto.PatternID == record.GetPatternID() {
			for _, table := range BiteProto.SlotPayTable {
				if table.GetId() == chance.RuleID {
					chance.Rule = table
					patternMap[BitePattern] = append(patternMap[BitePattern], chance)
				}
			}
		}
	}

	for mode, chancies := range patternMap {
		logic.bufFeatureGames[mode] = append(logic.bufFeatureGames[mode], chancies)
	}

}
