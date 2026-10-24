package module

import (
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// HookMsgHandler 攔截消息處理器
type HookMsgHandler interface {
	// 處理消息
	OnMessage(ctx *MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message

	// 是否處理
	IsHook(bclsid int32, sclsid int32) bool
}

// HookMsgHandlerGroup 攔截消息處理器組
type HookMsgHandlerGroup struct {
	innerHandlers []HookMsgHandler
}

// NewHookMsgHandlerGroup 創建攔截消息處理器組
func NewHookMsgHandlerGroup() *HookMsgHandlerGroup {
	hmg := new(HookMsgHandlerGroup)

	return hmg
}

// IsHook 是否處理
func (hlg *HookMsgHandlerGroup) IsHook(bclsid int, sclsid int) bool {
	return true
}

// OnMessage 處理消息
func (hlg *HookMsgHandlerGroup) OnMessage(ctx *MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	var msgs []*msg.Message

	for _, h := range hlg.innerHandlers {
		if h.IsHook(m.BClassID, m.SClassID) {
			msg := h.OnMessage(ctx, clt, m)

			if msg != nil {
				msgs = append(msgs, msg...)
			}
		}
	}

	return msgs
}

// AddHandler 增加處理器
func (hlg *HookMsgHandlerGroup) AddHandler(h HookMsgHandler) *HookMsgHandlerGroup {
	for _, exh := range hlg.innerHandlers {
		if exh == h {
			return hlg
		}
	}

	hlg.innerHandlers = append(hlg.innerHandlers, h)

	return hlg
}

// SemaphoreHookMsgHandlerGroup 限流的攔截消息處理器組
type SemaphoreHookMsgHandlerGroup struct {
	*HookMsgHandlerGroup
	semaphore util.Semaphore
}

// NewSemaphoreHookMsgHandlerGroup 創建限流的攔截消息處理器組對象
func NewSemaphoreHookMsgHandlerGroup(n int) *SemaphoreHookMsgHandlerGroup {
	shlg := new(SemaphoreHookMsgHandlerGroup)
	shlg.semaphore = util.NewSemaphore(n)
	shlg.HookMsgHandlerGroup = new(HookMsgHandlerGroup)

	return shlg
}

// OnMessage 處理消息
func (hlg *SemaphoreHookMsgHandlerGroup) OnMessage(ctx *MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	hlg.semaphore.Acquire()
	defer hlg.semaphore.Release()

	return hlg.HookMsgHandlerGroup.OnMessage(ctx, clt, m)
}
