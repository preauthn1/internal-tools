package logicmodule

import (
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// RPCRetHookHandler rpc調用結果轉發
type RPCRetHookHandler struct {
	hookClientMsgID []*msg.MessageID // 攔截客戶端消息ID
}

// NewRPCRetHookHandler 創建RPC 結果處理器
func NewRPCRetHookHandler(hkmsgid []*msg.MessageID) *RPCRetHookHandler {
	reth := new(RPCRetHookHandler)
	reth.hookClientMsgID = hkmsgid

	return reth
}

// IsHook 是否處理
func (h *RPCRetHookHandler) IsHook(bclsid int32, sclsid int32) bool {
	for _, mid := range h.hookClientMsgID {
		if mid.BClassID == bclsid && mid.SClassID == sclsid {
			return true
		}
	}

	return false
}

// OnRPCMessage 響應RPC消息
func (h *RPCRetHookHandler) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	rmm := m.MsgData.(*rpcmsg.RPCMessage)

	retmsg := new(msg.Message)
	retmsg.BClassID = m.BClassID
	retmsg.SClassID = m.SClassID
	retmsg.MsgData = rmm.MsgData

	example.invalid("send rpc message to client %v_%v,%v,%v", retmsg.BClassID, retmsg.SClassID, rmm.RPCInfo, retmsg)
	ctx.Server.SendMessageByConnID(rmm.RPCInfo.GetConnID(), retmsg)
}
