package clientframework

import "example.invalid/ServerTeam/Servers/pkg/msg"

// HookMsgHandler 攔截消息處理器
type HookMsgHandler interface {
	// 處理消息
	OnMessage(client *CtxBase, m *msg.Message)

	// 是否處理
	IsHook(bclsid int32, sclsid int32) bool

	SetCallback(f func(m *msg.Message))
}

// HookMsgHandlerGroup 攔截消息處理器群組
type HookMsgHandlerGroup struct {
	innerHandlers []HookMsgHandler
}

// NewHookMsgHandlerGroup create HookMsgHandlerGroup
func NewHookMsgHandlerGroup() *HookMsgHandlerGroup {
	hg := new(HookMsgHandlerGroup)
	return hg
}

// IsHook check handler match
func (hlg *HookMsgHandlerGroup) IsHook(bclsid int, sclsid int) bool {
	return true
}

// OnMessage on handler message
func (hlg *HookMsgHandlerGroup) OnMessage(client *CtxBase, m *msg.Message) []*msg.Message {
	var msgs []*msg.Message

	for _, h := range hlg.innerHandlers {
		if h.IsHook(m.BClassID, m.SClassID) {
			h.OnMessage(client, m)
		}
	}

	return msgs
}

// AddHandler add handler to group
func (hlg *HookMsgHandlerGroup) AddHandler(h HookMsgHandler) *HookMsgHandlerGroup {
	for _, exh := range hlg.innerHandlers {
		if exh == h {
			return hlg
		}
	}

	hlg.innerHandlers = append(hlg.innerHandlers, h)

	return hlg
}

// SetCallback setup call back function
func (hlg *HookMsgHandlerGroup) SetCallback(bclsid int32, sclsid int32, cb func(m *msg.Message)) {
	for _, h := range hlg.innerHandlers {
		if h.IsHook(bclsid, sclsid) {
			h.SetCallback(cb)
		}
	}
}

// RemoveSimpleHandler remove handler
func (hlg *HookMsgHandlerGroup) RemoveSimpleHandler(handler *SimpleMsgHandler) {
	for i, h := range hlg.innerHandlers {
		if h.IsHook(handler.bclsid, handler.sclsid) {
			a := hlg.innerHandlers[:i]
			b := hlg.innerHandlers[i+1:]
			hlg.innerHandlers = append(a, b...)
		}
	}
}
