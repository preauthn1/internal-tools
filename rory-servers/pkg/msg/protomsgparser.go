package msg

import (
	"errors"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ProtoParser 基于protobuf的消息解析器
//消息内容兼容[]byte和proto.Message
type ProtoParser struct {
	*ByteParser
}

// NewProtoMsgParser 創建基於protobuf的消息解析器
func NewProtoMsgParser() *ProtoParser {
	pm := new(ProtoParser)

	return pm
}

// Marshal ?
func (mp *ProtoParser) Marshal(msg *Message) ([]byte, error) {
	b1 := util.IntToBytes(int32(msg.BClassID), true)
	b2 := util.IntToBytes(int32(msg.SClassID), true)

	var b3 []byte
	if msg.MsgData != nil {
		var err error

		if bm, ok := msg.MsgData.([]byte); ok {
			b3 = bm
		} else if s, ok := msg.MsgData.(string); ok {
			b3 = []byte(s)
		} else {
			d, ok := msg.MsgData.(proto.Message)
			if !ok {
				mlog.Error("序列化消息產生錯誤,消息數據類型不正確")
				return nil, errors.New("數據類型必須為proto.Message")
			}

			b3, err = proto.Marshal(d)
			if err != nil {
				mlog.NewDefaultLogger().SetCallDepLen(8).Error("序列化消息產生錯誤,%v", err)
				return nil, err
			}
		}
	}

	return append(append(b1, b2...), b3...), nil
}
