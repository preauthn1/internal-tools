package msg

import (
	"errors"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ByteParser 基於[]byte消息內容的消息解析器
type ByteParser struct {
}

// Marshal ?
func (bmp *ByteParser) Marshal(msg *Message) ([]byte, error) {
	b1 := util.IntToBytes(int32(msg.BClassID), true)
	b2 := util.IntToBytes(int32(msg.SClassID), true)

	b3, ok := msg.MsgData.([]byte)

	if !ok {
		return nil, errors.New("數據類型必須為[]byte")
	}

	return append(append(b1, b2...), b3...), nil
}

// Unmarshal ?
func (bmp *ByteParser) Unmarshal(data []byte) (*Message, error) {
	msg := new(Message)
	msg.BClassID = util.BytesToInt(data[:4], true)
	msg.SClassID = util.BytesToInt(data[4:8], true)

	msg.MsgData = data[8:]

	return msg, nil
}
