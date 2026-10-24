package aladdin

import (
	"fmt"
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestALDLogic_CalcGameResult(t *testing.T) {

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("Spin", func(tt *testing.T) {

		for _, row := range logic.Proto.SlotPayLine {
			fmt.Println(row)
		}

		result := slot.GameResult{
			[]*protodata.SlotSymbolDistribution{
				{Symbol: proto.String("p5")},
				{Symbol: proto.String("p4")},
				{Symbol: proto.String("p6")},
				{Symbol: proto.String("p6")},
				{Symbol: proto.String("p7")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Symbol: proto.String("p5")},
				{Symbol: proto.String("p3")},
				{Symbol: proto.String("p6")},
				{Symbol: proto.String("p3")},
				{Symbol: proto.String("p6")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Symbol: proto.String("p6")},
				{Symbol: proto.String("p5")},
				{Symbol: proto.String("p6")},
				{Symbol: proto.String("p3")},
				{Symbol: proto.String("p4")},
			},
		}

		manager, _ := slot.NewGameResultManager(5, 3)
		manager.GameResult = result

		step := NewStepResult(1, false, 6, 60, logic.Proto.SymbolAttributes, manager, nil)

		logic.CalcGameResult(step)

		fmt.Println(step.GamePayResult)
		step.ExpandResult.RowMode()
		fmt.Println(step.ExpandResult)

	})
}
