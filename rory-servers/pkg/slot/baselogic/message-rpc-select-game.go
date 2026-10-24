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

// RecordSelectGameTimes 更新SelectGame
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 30003 (Slot_GameMessageClassID_Slot_SelectGameRecordID)
func (logic *SlotLogicBase) RecordSelectGameTimes(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_SelectGameRecordID)

	data := new(netproto.Slot_MSG_SelectGameRecord)
	example.invalid = new(netproto.Slot_SelectGameInfo)
	example.invalid.UserID = proto.Int32(user.UserID)
	example.invalid.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	example.invalid.FlagID = proto.Int32(int32(FlagID))

	example.invalid.PatternID = proto.Int32(logic.Proto.PatternID)

	if len(logic.SelectChances) > 0 {
		sc := logic.SelectChances[0]
		example.invalid.OriginID = proto.Int32(sc.Rule.GetId())
		example.invalid.Times = proto.Int32(int32(len(logic.SelectChances)))

		b, _ := json.Marshal(sc)
		example.invalid.Data = proto.String(string(b))
	} else {
		example.invalid.Times = proto.Int32(0)
		example.invalid.OriginID = proto.Int32(0)
		example.invalid.Data = proto.String("")
	}

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// ReadSelectGameTimes 讀取SelectGame
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 30001 (Slot_GameMessageClassID_Slot_SelectGameReadID)
func (logic *SlotLogicBase) ReadSelectGameTimes(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadID)

	data := new(netproto.Slot_MSG_SelectGameRead)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecvSelectGameTimes 接收 SelectGame
func (logic *SlotLogicBase) RecvSelectGameTimes(m *msg.Message) {

	defer logic.CheckWaitGroup()
	mlog.Debug("Table %d, RecvSelectGameTimes", logic.TableNo())

	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_SelectGameRecord{
			Info: &netproto.Slot_SelectGameInfo{},
		}
		return mm
	})

	detail := rmm.(*netproto.Slot_MSG_SelectGameRecord)

	if logic.GameUser == nil {
		return
	}

	if example.invalid == nil {
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

	if example.invalid.GetOriginID() == 0 || example.invalid.GetTimes() == 0 {
		return
	}

	// 查 Default 表
	if Proto != nil && (Proto.PatternID == example.invalid.GetPatternID() || example.invalid.GetPatternID() == 0) {
		for _, table := range Proto.SlotPayTable {
			if int32(table.GetId()) == example.invalid.GetOriginID() {
				sg := buildSelectChance(table, Proto, detail)

				var selections []*slot.SelectChance

				for i := 0; i < int(example.invalid.GetTimes()); i++ {
					selections = append(selections, sg)
				}

				logic.bufFeatureGames[DefaultPattern] = append(logic.bufFeatureGames[DefaultPattern], selections)
			}
		}
	}

	// 查 Release 表
	if ReleaseProto != nil && ReleaseProto.PatternID == example.invalid.GetPatternID() {
		for _, table := range ReleaseProto.SlotPayTable {
			if int32(table.GetId()) == example.invalid.GetOriginID() {
				sg := buildSelectChance(table, ReleaseProto, detail)

				var selections []*slot.SelectChance

				for i := 0; i < int(example.invalid.GetTimes()); i++ {
					selections = append(selections, sg)
				}

				logic.bufFeatureGames[ReleasePattern] = append(logic.bufFeatureGames[ReleasePattern], selections)
			}
		}
	}

	// 查 Bite 表
	if BiteProto != nil && BiteProto.PatternID == example.invalid.GetPatternID() {
		for _, table := range BiteProto.SlotPayTable {
			if int32(table.GetId()) == example.invalid.GetOriginID() {
				sg := buildSelectChance(table, BiteProto, detail)

				var selections []*slot.SelectChance

				for i := 0; i < int(example.invalid.GetTimes()); i++ {
					selections = append(selections, sg)
				}

				logic.bufFeatureGames[BitePattern] = append(logic.bufFeatureGames[BitePattern], selections)
			}
		}
	}

}

func buildSelectChance(rule *protodata.SlotPayTable, proto *protodata.ProtoData, detail *netproto.Slot_MSG_SelectGameRecord) *slot.SelectChance {

	if len(example.invalid.GetData()) == 0 {
		return nil
	}

	data := new(slot.SelectChance)

	if err := json.Unmarshal([]byte(example.invalid.GetData()), data); err != nil {
		mlog.Error("解析選擇機會資料錯誤: %v", err)
		return nil
	}

	sg, err := slot.NewSelectChance(proto, data.GameResultManager, rule, data.BetMoney, data.BetRate)
	if err != nil {
		mlog.Error("建立選擇機會錯誤: %v", err)
		return nil
	}
	sg.Data = data.Data

	return sg
}
