package hhf

import (
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// OnRecvRPCMessageFunc 接收DB資訊
func (logic *Logic) OnRecvRPCMessageFunc(m *msg.Message) {

	switch m.SClassID {
	// 從DB獲得的PickGame資訊
	case int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID):
		logic.RecvJackpotGame(m)
	}

}
