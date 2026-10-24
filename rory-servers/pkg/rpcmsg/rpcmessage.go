package rpcmsg

import (
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// RPCMessage 遠程消息
type RPCMessage struct {
	RPCInfo *netproto.RPCInfo // RPC客戶端信息
	MsgData interface{}       // 消息內容
}
