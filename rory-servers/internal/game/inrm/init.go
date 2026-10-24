package inrm

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// GameLogicCreator 遊戲邏輯產生器
var GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
	return NewGameServer(t, tfManager, cfgs)
}

func init() {
	msg.InstanceProvider.Register(int32(netproto.MessageBClassID_Game), int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID), func() proto.Message { return new(netproto.INRM_ThrowCard) })
	msg.InstanceProvider.Register(int32(netproto.MessageBClassID_Game), int32(netproto.INRM_GameMessageClassID_INRM_DeclareID), func() proto.Message { return new(netproto.INRM_Declare) })
	msg.InstanceProvider.Register(int32(netproto.MessageBClassID_Game), int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID), func() proto.Message { return new(netproto.Poker_Hand) })
}
