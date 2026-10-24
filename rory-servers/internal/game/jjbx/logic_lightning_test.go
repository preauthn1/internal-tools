package jjbx

import (
	"fmt"
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestJJBXLogic_SelectGame_Rate8_Lightning_15JackpotGrand(t *testing.T) {

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

	var opts []*netproto.Slot_ChooseGameOption
	var Prepare *netproto.Slot_MSG_LightningPrepare

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
			data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
			opts = data.Options
		}
		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LightningPrepareID) {
			time.Sleep(10 * time.Millisecond)
			data := m.MsgData.(*netproto.Slot_MSG_LightningPrepare)
			Prepare = data
		}
	}

	t.Run("Win", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("g4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("g5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("g6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("g4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("g5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("g6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},

				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("g4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("g5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("g6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("g4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("g5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("g6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			}

			logic.Proto.ExtraRule = []*protodata.ExtraRule{
				{Condition: proto.String("AfterSpin"), ConditionValue: proto.String("any:g1:15"), Action: proto.String("RespinR"), ActionValue: proto.String("8:8:10,18:8:10"), UseType: proto.String("all")},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		m.Rate = proto.Int32(8)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)
	})

	t.Run("Choose Lightning", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

		m := new(netproto.Slot_MSG_ChooseGame)

		var opt *netproto.Slot_ChooseGameOption
		for _, row := range opts {
			if row.GetType() == "lightning" {
				opt = row
			}
		}

		m.Option = opt

		tt.Log("Send:", m)

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if logic.LightningGame == nil {
			tt.Errorf("Game should gen Lightn Game")
		}
	})

	for Prepare != nil {
		t.Run(fmt.Sprintf("Show Step %d", Prepare.GetLightningStep()), func(tt *testing.T) {

			tmp := Prepare
			Prepare = nil

			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionLightningSpinID)

			m := new(netproto.Slot_MSG_LightningShow)
			m.LightningStep = proto.Int32(tmp.GetLightningStep())

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			logic.GameMessage(logic.GameUser, _m)
		})
	}
}

func TestJJBXLogic_SelectGame_Rate8_ChooseLightning_LimitToEnd(t *testing.T) {

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

	var opts []*netproto.Slot_ChooseGameOption
	var Prepare *netproto.Slot_MSG_LightningPrepare

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameOptionsID) {
			data := m.MsgData.(*netproto.Slot_MSG_ChooseGameOptions)
			opts = data.Options
		}
		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LightningPrepareID) {
			time.Sleep(10 * time.Millisecond)
			data := m.MsgData.(*netproto.Slot_MSG_LightningPrepare)
			Prepare = data
		}
	}

	t.Run("Win Select Game", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("g4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("g5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("g6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(701), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},

				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("g4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("g5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("g6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(721), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},

				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("g4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("g5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("g6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("g1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("g2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("g3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("g4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("g8"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(726), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("g7"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		m.Rate = proto.Int32(8)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	if len(opts) == 0 {
		t.Errorf("Game should has Select Game options")
	}

	t.Run("Choose Lightning", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

		m := new(netproto.Slot_MSG_ChooseGame)

		var opt *netproto.Slot_ChooseGameOption
		for _, row := range opts {
			if row.GetType() == "lightning" {
				opt = row
			}
		}

		m.Option = opt

		tt.Log("Send:", m)

		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	for Prepare != nil {
		t.Run(fmt.Sprintf("Show Step %d", Prepare.GetLightningStep()), func(tt *testing.T) {

			tmp := Prepare
			Prepare = nil

			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionLightningSpinID)

			m := new(netproto.Slot_MSG_LightningShow)
			m.LightningStep = proto.Int32(tmp.GetLightningStep())

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			logic.GameMessage(logic.GameUser, _m)

			time.Sleep(10 * time.Millisecond)
		})
	}
}
