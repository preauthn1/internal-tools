package dhzw

import (
	"strconv"
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot/mock"
	"example.invalid/ServerTeam/Servers/pkg/slot/slot_test"
)

func TestDHZWNewLogic(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	slot_test.PrintLogic(t, logic.SlotLogicBase)
}

func TestDHZWLogic_UserEnter(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_StatusID) {
			data := m.MsgData.(*netproto.Slot_MSG_Status)
			if data.GetState() == 0 {
				c <- true
				close(c)
			}
			return
		}
	}

	user.UserState.OnUserQueue(user, logic.GetTable().Room)
	go logic.UserEnter(user, logic.GetTable().TableNo, logic.GetTable().GetFreeSeat().SeatNo)
	<-c
}

func TestDHZWLogic_Spin_Normal(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_StatusID) {
			data := m.MsgData.(*netproto.Slot_MSG_Status)
			if data.GetState() == 0 {
				c <- true
				close(c)
			}
			return
		}
	}

	user.UserState.OnUserQueue(user, logic.GetTable().Room)
	go logic.UserEnter(user, logic.GetTable().TableNo, logic.GetTable().GetFreeSeat().SeatNo)
	<-c

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("Random Rate 28", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		m.Rate = proto.Int32(28)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Random Rate 38", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		m.Rate = proto.Int32(38)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Random Rate 58", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		m.Rate = proto.Int32(58)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Random Rate 88", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		m.Rate = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func TestDHZWLogic_Spin_88(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_StatusID) {
			data := m.MsgData.(*netproto.Slot_MSG_Status)
			if data.GetState() == 0 {
				c <- true
				close(c)
			}
			return
		}
	}

	user.UserState.OnUserQueue(user, logic.GetTable().Room)
	go logic.UserEnter(user, logic.GetTable().TableNo, logic.GetTable().GetFreeSeat().SeatNo)
	<-c

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("Random Rate 88", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		m.Rate = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func userEnter(logic *Logic, user *gamelogic.GameUser, t *testing.T, safe bool, m *msg.Message) {
	handlers1 := make(map[int32]func())

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadID)] = func() {

		detail := new(netproto.Slot_MSG_FreeGameRecord)
		example.invalid = new(netproto.Slot_FreeSpinDetail)
		example.invalid.UserID = proto.Int32(user.UserID)
		example.invalid.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
		f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
		example.invalid.FlagID = proto.Int32(int32(f))
		example.invalid.UserBet = proto.Int64(0)
		example.invalid.FreeSpinGroup = proto.Int32(0)
		example.invalid.Times = proto.Int32(0)
		example.invalid.ChooseRate = proto.Int32(0)
		example.invalid.WinMoney = proto.Int64(0)
		example.invalid.OriginID = proto.Int64(0)

		b, _ := proto.Marshal(detail)

		_m := new(rpcmsg.RPCMessage)
		_m.MsgData = b

		mm := new(msg.Message)
		mm.BClassID = int32(netproto.MessageBClassID_Game)
		mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadRetID)
		mm.MsgData = _m

		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID)] = func() {
		detail := new(netproto.Slot_MSG_JackpotRecord)

		jp1 := new(netproto.Slot_JackpotData)
		jp1.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
		f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
		jp1.FlagID = proto.Int32(int32(f))
		jp1.TableNo = proto.Int32(logic.TableManager.TableNo)
		jp1.Contribution = proto.Int64(0)
		jp1.TotalContribution = proto.Int64(0)
		jp1.AwardTotal = proto.Int64(0)
		jp1.Level = proto.Int32(1)

		jp2 := new(netproto.Slot_JackpotData)
		jp2.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
		f, _ = strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
		jp2.FlagID = proto.Int32(int32(f))
		jp2.TableNo = proto.Int32(logic.TableManager.TableNo)
		jp2.Contribution = proto.Int64(0)
		jp2.TotalContribution = proto.Int64(0)
		jp2.AwardTotal = proto.Int64(0)
		jp2.Level = proto.Int32(2)

		jp3 := new(netproto.Slot_JackpotData)
		jp3.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
		f, _ = strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
		jp3.FlagID = proto.Int32(int32(f))
		jp3.TableNo = proto.Int32(logic.TableManager.TableNo)
		jp3.Contribution = proto.Int64(0)
		jp3.TotalContribution = proto.Int64(0)
		jp3.AwardTotal = proto.Int64(0)
		jp3.Level = proto.Int32(3)

		detail.Records = append(detail.Records, jp1)
		detail.Records = append(detail.Records, jp2)
		detail.Records = append(detail.Records, jp3)

		b, _ := proto.Marshal(detail)

		_m := new(rpcmsg.RPCMessage)
		_m.MsgData = b

		mm := new(msg.Message)
		mm.BClassID = int32(netproto.MessageBClassID_Game)
		mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID)
		mm.MsgData = _m

		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadID)] = func() {

		detail := new(netproto.Slot_MSG_SelectGameRecord)
		example.invalid = new(netproto.Slot_SelectGameInfo)
		example.invalid.UserID = proto.Int32(user.UserID)
		example.invalid.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
		f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
		example.invalid.FlagID = proto.Int32(int32(f))
		example.invalid.Times = proto.Int32(0)

		b, _ := proto.Marshal(detail)

		_m := new(rpcmsg.RPCMessage)
		_m.MsgData = b

		mm := new(msg.Message)
		mm.BClassID = int32(netproto.MessageBClassID_Game)
		mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadRetID)
		mm.MsgData = _m

		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID)] = func() {
		detail := mock.EmptyBonusGameDetail(logic.SlotLogicBase, user)
		mm := mock.MarshalBonusGame(detail)
		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID)] = func() {

		f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)

		detail := new(netproto.Slot_MSG_DeferredAwardRecord)
		detail.UserID = proto.Int32(user.UserID)
		detail.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
		detail.FlagID = proto.Int32(int32(f))

		// detail.Records = append(detail.Records, &netproto.Slot_DeferredAward{
		// 	OriginID:  proto.Int32(412),
		// 	AwardType: proto.String(string(slot.GameModeJackpotGame)),
		// 	Data:      proto.String("{\"RuleID\":413,\"BetMoney\":10,\"BetRate\":28,\"GameResultManager\":{\"GameResult\":[[{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":4}]],\"MultiSymbol\":{},\"Reel\":5,\"Row\":4,\"SymbolGroupID\":805,\"BackGroundGroupID\":0,\"BackgroundResult\":{\"Reel\":5,\"Row\":4,\"Distributions\":[{\"groupId\":801,\"reelNum\":5,\"bgEnum\":\"bg0\",\"proportion\":1},{\"groupId\":801,\"reelNum\":4,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":3,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":2,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":1,\"bgEnum\":\"bg4\",\"proportion\":1}]}},\"Data\":{}}"),
		// })

		// detail.Records = append(detail.Records, &netproto.Slot_DeferredAward{
		// 	OriginID:  proto.Int32(412),
		// 	AwardType: proto.String(string(slot.GameModeJackpotGame)),
		// 	Data:      proto.String("{\"RuleID\":413,\"BetMoney\":10,\"BetRate\":28,\"GameResultManager\":{\"GameResult\":[[{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":4}]],\"MultiSymbol\":{},\"Reel\":5,\"Row\":4,\"SymbolGroupID\":805,\"BackGroundGroupID\":0,\"BackgroundResult\":{\"Reel\":5,\"Row\":4,\"Distributions\":[{\"groupId\":801,\"reelNum\":5,\"bgEnum\":\"bg0\",\"proportion\":1},{\"groupId\":801,\"reelNum\":4,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":3,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":2,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":1,\"bgEnum\":\"bg4\",\"proportion\":1}]}},\"Data\":{}}"),
		// })

		// detail.Records = append(detail.Records, &netproto.Slot_DeferredAward{
		// 	OriginID:  proto.Int32(412),
		// 	AwardType: proto.String(string(slot.GameModeJackpotGame)),
		// 	Data:      proto.String("{\"RuleID\":413,\"BetMoney\":10,\"BetRate\":28,\"GameResultManager\":{\"GameResult\":[[{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":1,\"symbol\":\"p1\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":2,\"symbol\":\"p2\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":3,\"symbol\":\"p3\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":4,\"symbol\":\"p4\",\"proportion\":1,\"order\":4}],[{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":1},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":2},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":3},{\"groupId\":805,\"reelNum\":5,\"symbol\":\"p5\",\"proportion\":1,\"order\":4}]],\"MultiSymbol\":{},\"Reel\":5,\"Row\":4,\"SymbolGroupID\":805,\"BackGroundGroupID\":0,\"BackgroundResult\":{\"Reel\":5,\"Row\":4,\"Distributions\":[{\"groupId\":801,\"reelNum\":5,\"bgEnum\":\"bg0\",\"proportion\":1},{\"groupId\":801,\"reelNum\":4,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":3,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":2,\"bgEnum\":\"bg4\",\"proportion\":1},{\"groupId\":801,\"reelNum\":1,\"bgEnum\":\"bg4\",\"proportion\":1}]}},\"Data\":{}}"),
		// })

		b, _ := proto.Marshal(detail)

		_m := new(rpcmsg.RPCMessage)
		_m.MsgData = b

		mm := new(msg.Message)
		mm.BClassID = int32(netproto.MessageBClassID_Game)
		mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID)
		mm.MsgData = _m

		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadID)] = func() {
		mm := mock.EmptyUserAssignedModeRPCMessage(logic.SlotLogicBase, user)
		logic.OnRecvRPCMessage(mm)
	}

	slot_test.PrintSendRPCMessgeFunc(t, safe, m, handlers1)

	// printSendRPCMessgeFunc(t, safe, m, handlers2)
}

func printSendMessgeFunc(t *testing.T, userID int32, m *msg.Message) {}

func printSendRPCMessgeFunc(t *testing.T, safe bool, m *msg.Message, handlers map[int32]func()) {
	if handlers != nil {
		if f, ok := handlers[m.SClassID]; ok {
			f()
		}
	}
}
