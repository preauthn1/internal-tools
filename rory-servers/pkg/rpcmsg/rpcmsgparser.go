package rpcmsg

import (
	"errors"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// CreateProtoMsgFun ?
type CreateProtoMsgFun func() proto.Message

// RPCMessageParser 遠程消息解析器
type RPCMessageParser struct{}

// NewRPCMessageParser ?
func NewRPCMessageParser() *RPCMessageParser {
	p := new(RPCMessageParser)

	return p
}

// Marshal ?
func (rp *RPCMessageParser) Marshal(msg *msg.Message) ([]byte, error) {
	// 一級協議號
	b1 := util.IntToBytes(int32(msg.BClassID), true)
	// 二級協議號
	b2 := util.IntToBytes(int32(msg.SClassID), true)

	// 消息轉換為rpcmsg
	rpcmsg, ok := msg.MsgData.(*RPCMessage)
	if !ok {
		return nil, errors.New("消息內容不正確,必須為*RPCMessage類型")
	}

	rpcdata, err := proto.Marshal(rpcmsg.RPCInfo)
	if err != nil {
		return nil, err
	}

	var msgdata []byte

	if ptmdd, ok := rpcmsg.MsgData.(proto.Message); ok {
		msgdata, err = proto.Marshal(ptmdd)

		if err != nil {
			return nil, err
		}
	}

	if btmdd, ok := rpcmsg.MsgData.([]byte); ok {
		msgdata = btmdd
	}

	//	if msgdata == nil {
	//		return nil, errors.New("數據格式不正確.")
	//	}

	// 消息頭
	bh := append(b1, b2...)
	// rpc 消息頭
	rpclen := int32(len(rpcdata))
	ret := append(bh, util.IntToBytes(rpclen, true)...)
	ret = append(ret, rpcdata...)
	if msgdata != nil {
		ret = append(ret, msgdata...)
	}

	return ret, nil
}

// Unmarshal ?
func (rp *RPCMessageParser) Unmarshal(data []byte) (*msg.Message, error) {
	m1 := new(msg.Message)
	m1.BClassID = util.BytesToInt(data[:4], true)
	m1.SClassID = util.BytesToInt(data[4:8], true)

	msgdata := data[8:]
	rpcdatalen := util.BytesToInt(msgdata[:4], true)

	if rpcdatalen+4 > int32(len(msgdata)) {
		return nil, errors.New("解析消息錯誤,rpc數據頭長度不正確")
	}

	rpcdata := msgdata[4 : 4+rpcdatalen]

	rpcinfo := &netproto.RPCInfo{}
	err := proto.Unmarshal(rpcdata, rpcinfo)

	if err != nil {
		return nil, err
	}

	rpcmsg := new(RPCMessage)
	rpcmsg.RPCInfo = rpcinfo
	rpcmsg.MsgData = msgdata[4+rpcdatalen:]

	m1.MsgData = rpcmsg

	return m1, nil

}
