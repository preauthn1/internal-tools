package logicmodule

import (
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
)

// RPCHookHandler 響應RPC消息處理器
type RPCHookHandler interface {
	// 響應RPC消息
	OnRPCMessage(msg *msg.Message, ctx *module.MouduleCtx)

	// 是否處理
	IsHook(bclsid int32, sclsid int32) bool
}

// RPCHookHandlerGroup 響應RPC消息處理器群組
type RPCHookHandlerGroup struct {
	innerHandlers []RPCHookHandler
}

// NewRPCHookHandlerGroup 響應RPC消息處理器群組
func NewRPCHookHandlerGroup() *RPCHookHandlerGroup {
	hg := new(RPCHookHandlerGroup)

	return hg
}

// IsHook 是否處理
func (hg *RPCHookHandlerGroup) IsHook(bclsid int32, sclsid int32) bool {
	return true
}

// OnRPCMessage 響應RPC消息
func (hg *RPCHookHandlerGroup) OnRPCMessage(msg *msg.Message, ctx *module.MouduleCtx) {
	for _, h := range hg.innerHandlers {
		if h.IsHook(msg.BClassID, msg.SClassID) {
			h.OnRPCMessage(msg, ctx)
		}
	}
}

// AddHandler 添加處理器
func (hg *RPCHookHandlerGroup) AddHandler(h RPCHookHandler) *RPCHookHandlerGroup {
	hg.innerHandlers = append(hg.innerHandlers, h)
	return hg
}
