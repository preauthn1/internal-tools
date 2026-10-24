package hallutil

import (
	"fmt"

	proto "example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// MsgTypeIDHallKick ?
const (
	MsgTypeIDHallKick int32 = 2
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
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[GetTipMessage] err : %v", err)).Error("HallServer壓測時")
		return nil, err
	}

	return m, nil
}
