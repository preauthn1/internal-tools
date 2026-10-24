package baselogic

import (
	"strconv"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// ReadUserAssignedMode 讀取 UserAssignedMode 紀錄
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000071 (Slot_GameMessageClassID_Slot_UserAssignedModeReadID)
func (logic *SlotLogicBase) ReadUserAssignedMode(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadID)

	data := new(netproto.Slot_MSG_UserGameQuery)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// AssignUserMode UserAssignedMode 存入DB
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000073 (Slot_GameMessageClassID_Slot_UserAssignedModeRecordID)
func (logic *SlotLogicBase) AssignUserMode(uid int32, mode PatternMode, t time.Time) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeRecordID)

	data := new(netproto.Slot_MSG_UserAssignedModeRecord)
	data.UserID = proto.Int32(uid)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))

	data.Records = new(netproto.Slot_UserAssignedMode)
	data.Records.Mode = proto.Int32(int32(mode))
	data.Records.ExpireTime = proto.String(t.Format("2006-01-02 15:04:05.000"))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// UnassignUserMode UserAssignedMode 刪除
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 10000074 (Slot_GameMessageClassID_Slot_UserAssignedModeDeleteID)
func (logic *SlotLogicBase) UnassignUserMode(uid int32) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeDeleteID)

	data := new(netproto.Slot_MSG_UserGameQuery)
	data.UserID = proto.Int32(uid)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.ParseInt(logic.TableManager.Room.Flag, 10, 32)
	data.FlagID = proto.Int32(int32(FlagID))

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecvUserAssignedMode 接收db server回傳的 UserAssignedMode 資訊，並整理至Logic內
func (logic *SlotLogicBase) RecvUserAssignedMode(m *msg.Message) {
	defer func() {
		logic.ModeWait <- true
	}()
	mlog.Debug("Table %d, RecvUserAssignedMode", logic.TableNo())

	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_UserAssignedModeRecord{
			Records: &netproto.Slot_UserAssignedMode{},
		}
		return mm
	})

	detail := rmm.(*netproto.Slot_MSG_UserAssignedModeRecord)

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

	if detail.Records == nil {
		return
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05.000", detail.Records.GetExpireTime(), time.Local)
	if err != nil {
		mlog.Warn("RecvUserAssignedMode ExpireTime Parse Error, %v", err)
		return
	}

	logic.Mode = PatternMode(detail.Records.GetMode())
	logic.ModeExpire = t

}
