package jgyf

import (
	"fmt"
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestJGYFLogic_FreeSpin_WinTimes(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
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

	var opts []*netproto.Slot_ChooseGameOption
	var leftFreeGame *netproto.Slot_MSG_LeftFreeGameData

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
			data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
			opts = data.Options
		}
		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
			data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
			leftFreeGame = data
		}
	}

	t.Run("Win", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(1)
		m.Rate = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if len(logic.SelectChances) == 0 {
			tt.Errorf("SelectChance length want > 0")
		}
	})

	if len(opts) == 0 {
		t.Errorf("Game should has Select Game options")
	}

	t.Run("Choose", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

		m := new(netproto.Slot_MSG_ChooseGame)

		var opt *netproto.Slot_ChooseGameOption
		for _, row := range opts {
			if row.GetType() == "free_spins" {
				opt = row
				break
			}
		}

		m.Option = opt

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}
	})

	logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
	}

	for i := 0; i < 5 && leftFreeGame != nil && leftFreeGame.GetCount() > 0; i++ {
		t.Run(fmt.Sprintf("Free Spin left %d", leftFreeGame.GetCount()), func(tt *testing.T) {

			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

			m := new(netproto.Slot_MSG_Spin)

			m.Money = proto.Int64(leftFreeGame.GetBetMoney())
			m.Rate = proto.Int32(int32(leftFreeGame.GetBetRate()))

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			logic.GameMessage(logic.GameUser, _m)

		})
	}
}

func TestJGYFLogic_FreeSpin_5x3(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
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

	var opts []*netproto.Slot_ChooseGameOption
	var leftFreeGame *netproto.Slot_MSG_LeftFreeGameData

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
			data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
			opts = data.Options
		}
		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
			data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
			leftFreeGame = data
		}
	}

	t.Run("Win", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(1)
		m.Rate = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if len(logic.SelectChances) == 0 {
			tt.Errorf("SelectChance length want > 0")
		}
	})

	if len(opts) == 0 {
		t.Errorf("Game should has Select Game options")
	}

	t.Run("Choose", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

		m := new(netproto.Slot_MSG_ChooseGame)

		var opt *netproto.Slot_ChooseGameOption
		for _, row := range opts {
			if row.GetType() == "free_spins" {
				opt = row
				break
			}
		}

		m.Option = opt

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}
	})

	logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(1), Order: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(2), Order: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(3), Order: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(4), Order: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1402), ReelNum: proto.Uint32(5), Order: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
	}

	for leftFreeGame != nil && leftFreeGame.GetCount() > 0 {
		t.Run(fmt.Sprintf("Free Spin left %d", leftFreeGame.GetCount()), func(tt *testing.T) {
			logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
				// slot_test.PrintSendMessgeFunc(tt, userID, m)

				if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
					data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
					leftFreeGame = data
				}
			}

			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

			m := new(netproto.Slot_MSG_Spin)

			m.Money = proto.Int64(leftFreeGame.GetBetMoney())
			m.Rate = proto.Int32(int32(leftFreeGame.GetBetRate()))

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			logic.GameMessage(logic.GameUser, _m)

		})
	}
}

func TestJGYFLogic_FreeSpin_5x4(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
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

	var opts []*netproto.Slot_ChooseGameOption
	var leftFreeGame *netproto.Slot_MSG_LeftFreeGameData

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
			data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
			opts = data.Options
		}
		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
			data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
			leftFreeGame = data
		}
	}

	t.Run("Win", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(1)
		m.Rate = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if len(logic.SelectChances) == 0 {
			tt.Errorf("SelectChance length want > 0")
		}
	})

	if len(opts) == 0 {
		t.Errorf("Game should has Select Game options")
	}

	t.Run("Choose", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

		m := new(netproto.Slot_MSG_ChooseGame)

		var opt *netproto.Slot_ChooseGameOption
		var skip int = 1
		for i, row := range opts {
			if i < skip {
				continue
			}
			if row.GetType() == "free_spins" {
				opt = row
				break
			}
		}

		m.Option = opt

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}
	})

	logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(1), Order: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(2), Order: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(3), Order: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(4), Order: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1403), ReelNum: proto.Uint32(5), Order: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
	}

	for leftFreeGame != nil && leftFreeGame.GetCount() > 0 {
		t.Run(fmt.Sprintf("Free Spin left %d", leftFreeGame.GetCount()), func(tt *testing.T) {
			logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
				// slot_test.PrintSendMessgeFunc(tt, userID, m)

				if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
					data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
					leftFreeGame = data
				}
			}

			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

			m := new(netproto.Slot_MSG_Spin)

			m.Money = proto.Int64(leftFreeGame.GetBetMoney())
			m.Rate = proto.Int32(int32(leftFreeGame.GetBetRate()))

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			logic.GameMessage(logic.GameUser, _m)

		})
	}
}

func TestJGYFLogic_FreeSpin_5x5(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
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

	var opts []*netproto.Slot_ChooseGameOption
	var leftFreeGame *netproto.Slot_MSG_LeftFreeGameData

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
			data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
			opts = data.Options
		}
		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
			data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
			leftFreeGame = data
		}
	}

	t.Run("Win", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(1)
		m.Rate = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if len(logic.SelectChances) == 0 {
			tt.Errorf("SelectChance length want > 0")
		}
	})

	if len(opts) == 0 {
		t.Errorf("Game should has Select Game options")
	}

	t.Run("Choose", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

		m := new(netproto.Slot_MSG_ChooseGame)

		var opt *netproto.Slot_ChooseGameOption
		var skip int = 2
		for i, row := range opts {
			if i < skip {
				continue
			}
			if row.GetType() == "free_spins" {
				opt = row
				break
			}
		}

		m.Option = opt

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}
	})

	logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(1), Order: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(2), Order: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(3), Order: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(4), Order: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(5), Order: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(1), Order: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(2), Order: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(3), Order: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(4), Order: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1404), ReelNum: proto.Uint32(5), Order: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
	}

	for leftFreeGame != nil && leftFreeGame.GetCount() > 0 {
		t.Run(fmt.Sprintf("Free Spin left %d", leftFreeGame.GetCount()), func(tt *testing.T) {
			logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
				// slot_test.PrintSendMessgeFunc(tt, userID, m)

				if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
					data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
					leftFreeGame = data
				}
			}

			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

			m := new(netproto.Slot_MSG_Spin)

			m.Money = proto.Int64(leftFreeGame.GetBetMoney())
			m.Rate = proto.Int32(int32(leftFreeGame.GetBetRate()))

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			logic.GameMessage(logic.GameUser, _m)

		})
	}
}

func TestJGYFLogic_FreeSpin_5x6(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
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

	var opts []*netproto.Slot_ChooseGameOption
	var leftFreeGame *netproto.Slot_MSG_LeftFreeGameData

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
			data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
			opts = data.Options
		}
		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
			data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
			leftFreeGame = data
		}
	}

	t.Run("Win", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1401), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(1)
		m.Rate = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if len(logic.SelectChances) == 0 {
			tt.Errorf("SelectChance length want > 0")
		}
	})

	if len(opts) == 0 {
		t.Errorf("Game should has Select Game options")
	}

	t.Run("Choose", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

		m := new(netproto.Slot_MSG_ChooseGame)

		var opt *netproto.Slot_ChooseGameOption
		var skip int = 3
		for i, row := range opts {
			if i < skip {
				continue
			}
			if row.GetType() == "free_spins" {
				opt = row
				break
			}
		}

		m.Option = opt

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}
	})

	logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(1), Order: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(2), Order: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(3), Order: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(4), Order: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(1405), ReelNum: proto.Uint32(5), Order: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
	}

	for leftFreeGame != nil && leftFreeGame.GetCount() > 0 {
		t.Run(fmt.Sprintf("Free Spin left %d", leftFreeGame.GetCount()), func(tt *testing.T) {
			logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
				// slot_test.PrintSendMessgeFunc(tt, userID, m)

				if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LeftFreeGameID) {
					data := m.MsgData.(*netproto.Slot_MSG_LeftFreeGameData)
					leftFreeGame = data
				}
			}

			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

			m := new(netproto.Slot_MSG_Spin)

			m.Money = proto.Int64(leftFreeGame.GetBetMoney())
			m.Rate = proto.Int32(int32(leftFreeGame.GetBetRate()))

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			logic.GameMessage(logic.GameUser, _m)

		})
	}
}
