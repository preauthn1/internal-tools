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
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// RecordLightningData 更新LightningGame
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000053 (Slot_GameMessageClassID_Slot_LightningReadID)
func (logic *SlotLogicBase) RecordLightningData(user *gamelogic.GameUser) {

	if logic.LightningGame == nil {
		return
	}

	b, err := json.Marshal(logic.LightningGame)
	if err != nil {
		mlog.Error("封存LightningGame json.Marshal錯誤, %v", err)
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID)

	data := new(netproto.Slot_MSG_LightningRecord)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))
	data.PatternID = proto.Int32(logic.Proto.PatternID)
	data.Data = proto.String(string(b))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// ResetLightningData 更新LightningGame
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000053 (Slot_GameMessageClassID_Slot_LightningRecordID)
func (logic *SlotLogicBase) ResetLightningData(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID)

	data := new(netproto.Slot_MSG_LightningRecord)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))
	data.Data = proto.String("")

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// ReadLightningData 讀取LightningGame
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000051 (Slot_GameMessageClassID_Slot_LightningReadID)
func (logic *SlotLogicBase) ReadLightningData(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_LightningReadID)

	data := new(netproto.Slot_MSG_LightningRead)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

func (logic *SlotLogicBase) recvLightningData(m *msg.Message) {

	defer logic.CheckWaitGroup()
	mlog.Debug("Table %d, RecvLightningData", logic.TableNo())

	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_LightningRecord{}
		return mm
	})

	detail := rmm.(*netproto.Slot_MSG_LightningRecord)

	if logic.GameUser == nil {
		return
	}

	if logic.GameUser.UserID != detail.GetUserID() {
		return
	}

	if logic.TableManager.Room.GameInfo.GameID != detail.GetGameID() {
		return
	}

	if logic.TableManager.Room.Flag != strconv.Itoa(int(detail.GetFlagID())) {
		return
	}

	if len(detail.GetData()) == 0 {
		mlog.Debug("無Lightning資訊")
		return
	}

	b := []byte(detail.GetData())
	lg := new(LightningGame)

	if err := json.Unmarshal(b, lg); err != nil {
		mlog.Error("解析LightningGame錯誤, %v", err)
		return
	}

	lg.Logic = logic

	for _, step := range lg.Steps {
		var respin *protodata.RespinTable
		for _, row := range logic.Proto.RespinTable {
			if int32(row.GetId()) == step.RespinID {
				respin = row
			}
		}
		step.Origin = respin

		var rule *protodata.SlotPayTable
		for _, row := range logic.Proto.SlotPayTable {
			if int32(row.GetId()) == step.RuleID {
				rule = row
			}
		}
		step.Rule = rule
	}

	step := lg.GetLastStep()

	if step.LastChance == 0 {
		mlog.Debug("Lightning次數已用完")
		return
	}

	for _, row := range logic.Proto.RespinTable {
		if int32(row.GetId()) == lg.OriginID {
			lg.Origin = row
		}
	}

	// 查 Default 表
	if Proto != nil && (Proto.PatternID == detail.GetPatternID() || detail.GetPatternID() == 0) {
		for _, table := range Proto.RespinTable {
			if int32(table.GetId()) == lg.OriginID {
				lg.Origin = table
				logic.bufFeatureGames[DefaultPattern] = append(logic.bufFeatureGames[DefaultPattern], lg)
			}
		}
	}

	// 查 Release 表
	if ReleaseProto != nil && ReleaseProto.PatternID == detail.GetPatternID() {
		for _, table := range ReleaseProto.RespinTable {
			if int32(table.GetId()) == lg.OriginID {
				lg.Origin = table
				logic.bufFeatureGames[ReleasePattern] = append(logic.bufFeatureGames[ReleasePattern], lg)
			}
		}
	}

	// 查 Bite 表
	if BiteProto != nil && BiteProto.PatternID == detail.GetPatternID() {
		for _, table := range BiteProto.RespinTable {
			if int32(table.GetId()) == lg.OriginID {
				lg.Origin = table
				logic.bufFeatureGames[BitePattern] = append(logic.bufFeatureGames[BitePattern], lg)
			}
		}
	}

}
