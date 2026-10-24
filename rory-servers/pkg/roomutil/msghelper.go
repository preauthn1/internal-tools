package roomutil

import (
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// ?
const (
	MsgTypeIDRoomKick    int32 = 1
	MsgTypeIDHallKick    int32 = 2
	MsgTypeIDTIP         int32 = 3
	MsgTypeIDTIPRECHARGE int32 = 4
)

// ?
const (
	MsgParamMoneyNotEnough int32 = 1
	MsgParamMoneyOverLoad  int32 = 2
)

// GetTipMessage ?
func GetTipMessage(str string, msgtypeid int32) (*msg.Message, error) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Common)
	m.SClassID = int32(netproto.PlatformCommonClassID_TipMessageID)

	tm := new(netproto.TipMessage)
	tm.Message = proto.String(str)
	tm.MsgType = proto.Int32(msgtypeid)
	data, err := proto.Marshal(tm)

	m.MsgData = data

	if err != nil {
		return nil, err
	}

	return m, nil
}

// GetTipMessageWithParam ?
func GetTipMessageWithParam(str string, msgtypeid int32, param int32) (*msg.Message, error) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Common)
	m.SClassID = int32(netproto.PlatformCommonClassID_TipMessageID)

	tm := new(netproto.TipMessage)
	tm.Message = proto.String(str)
	tm.MsgType = proto.Int32(msgtypeid)
	tm.Param = proto.Int32(param)
	data, err := proto.Marshal(tm)

	m.MsgData = data

	if err != nil {
		return nil, err
	}

	return m, nil
}

// KickUser 踢人
func KickUser(msg string, connid string, kickType int32, msgserver *msgserver.MsgServer) {
	m, _ := GetTipMessage(msg, kickType)

	msgserver.SendMessageByConnID(connid, m)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)

		msgserver.CloseConnByConnID(connid)
	}()
}

// KickUserParam 踢人
func KickUserParam(msg string, connid string, kickType int32, param int32, msgserver *msgserver.MsgServer) {
	m, _ := GetTipMessageWithParam(msg, kickType, param)

	msgserver.SendMessageByConnID(connid, m)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)

		msgserver.CloseConnByConnID(connid)
	}()
}

// KickUserToNewRoom 切換AB服, 把玩家踢A服出並送到新的B服
func KickUserToNewRoom(msg string, connid string, kickType int32, param int32, msgserver *msgserver.MsgServer) {
	m, _ := GetTipMessageWithParam(msg, kickType, param)
	msgserver.SendMessageByConnID(connid, m)

	go func() {
		time.Sleep(time.Duration(2) * time.Second)

		msgserver.CloseConnByConnID(connid)
	}()
}
