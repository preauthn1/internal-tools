package msg

import (
	"fmt"
)

// Message 網絡消息對象定義
type Message struct {
	BClassID int32       // 一級協議號
	SClassID int32       // 二級協議號
	MsgData  interface{} // 消息數據
}

func (msg *Message) String() string {
	return fmt.Sprintf("%d:%d", msg.BClassID, msg.SClassID)
}

// NewMessage 創建消息
func NewMessage(bclsID int32, sclsID int32, data interface{}) *Message {
	msg := new(Message)
	msg.BClassID = bclsID
	msg.SClassID = sclsID
	msg.MsgData = data

	return msg
}

// MessageID 消息ID
type MessageID struct {
	BClassID int32
	SClassID int32
}

// RPCMessageID 轉發消息ID
type RPCMessageID struct {
	MsgID     *MessageID // 消息ID
	NeedLogon bool       // 是否需要登錄
}

// NewRPCMessageID 創建RPCMessageID 對象
func NewRPCMessageID(bclsid int32, sclsid int32, needlogon bool) *RPCMessageID {
	rpcmsgid := new(RPCMessageID)
	rpcmsgid.MsgID = new(MessageID)
	rpcmsgid.MsgID.BClassID = bclsid
	rpcmsgid.MsgID.SClassID = sclsid
	rpcmsgid.NeedLogon = needlogon

	return rpcmsgid
}

// Equals 消息ID是否一致
func (msgid *MessageID) Equals(othermsgid *MessageID) bool {
	return msgid.BClassID == othermsgid.BClassID && msgid.SClassID == othermsgid.SClassID
}
