package baselogic

import (
	"strconv"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// ReadJackpotPoolNoLink DFDC類型為非連線模式，每桌一池
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 30021 (Slot_GameMessageClassID_Slot_JackpotNoLinkReadID)
func (logic *SlotLogicBase) ReadJackpotPoolNoLink() {

	//DFDC類型為非連線模式，每桌一池
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID)

	data := new(netproto.Slot_MSG_JackpotRead)
	example.invalid = new(netproto.Slot_JackpotData)
	example.invalid.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.Atoi(logic.TableManager.Room.Flag)
	example.invalid.FlagID = proto.Int32(int32(FlagID))
	example.invalid.TableNo = proto.Int32(logic.TableManager.TableNo)
	example.invalid.Level = proto.Int32(slot.JackpotSingle.Int())

	//讀取資料不須填值
	example.invalid.Contribution = proto.Int64(0)
	example.invalid.TotalContribution = proto.Int64(0)
	example.invalid.AwardTotal = proto.Int64(0)

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// UpdateJackpotPoolNoLink 更新獎池，所有數字都是DB幫你加好，只需要帶增加的量，獎池只有一個所以只更新Level 1
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 30023 (Slot_GameMessageClassID_Slot_JackpotNoLinkRecordID)
func (logic *SlotLogicBase) UpdateJackpotPoolNoLink(contribution, totalContribution, awardTotal int64, level slot.JackpotLevel) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkRecordID)

	data := new(netproto.Slot_MSG_JackpotRecord)

	detail := new(netproto.Slot_JackpotData)
	detail.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	FlagID, _ := strconv.Atoi(logic.TableManager.Room.Flag)
	detail.FlagID = proto.Int32(int32(FlagID))
	detail.TableNo = proto.Int32(logic.TableManager.TableNo)

	detail.Contribution = proto.Int64(contribution)
	detail.TotalContribution = proto.Int64(totalContribution)
	detail.AwardTotal = proto.Int64(awardTotal)
	detail.Level = proto.Int32(level.Int())

	data.Records = append(data.Records, detail)

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecvJackpotData 接收db server回傳的Jackpot資訊，並整理至Logic內
func (logic *SlotLogicBase) RecvJackpotData(m *msg.Message) {

	defer logic.CheckWaitGroup()
	mlog.Debug("Table %d, RecvJackpotData", logic.TableNo())

	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		return new(netproto.Slot_MSG_JackpotRecord)
	})

	data := rmm.(*netproto.Slot_MSG_JackpotRecord)

	for _, row := range data.Records {

		if row.GetGameID() != logic.TableManager.Room.GameInfo.GameID {
			continue
		}

		if strconv.Itoa(int(row.GetFlagID())) != logic.TableManager.Room.Flag {
			continue
		}

		if row.GetTableNo() != logic.TableManager.TableNo {
			continue
		}

		detail, exist := logic.Jackpots[slot.NewJackpotLevel(row.GetLevel())]

		if !exist {
			continue
		}

		detail.Contribution = row.GetContribution()
		detail.TotalContribution = row.GetTotalContribution()
		detail.AwardTotal = row.GetAwardTotal()

		logic.Jackpots[detail.Level] = detail
	}
}
