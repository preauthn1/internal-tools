package logicmodule

import (
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/player"
)

// NotifyHookHandler 通知消息處理器
type NotifyHookHandler interface {
	// 響應Notify消息
	OnNotifyMessage(msg *msg.Message, ctx *module.MouduleCtx, player *player.Manager)

	// 是否處理
	IsHook(bclsid int32, sclsid int32) bool
}

// NotifyHookHandlerGroup 通知消息處理器組
type NotifyHookHandlerGroup struct {
	handlers []NotifyHookHandler
}

// NewNotifyHookHandlerGroup 通知消息處理器組
func NewNotifyHookHandlerGroup() *NotifyHookHandlerGroup {
	ntfhg := new(NotifyHookHandlerGroup)
	ntfhg.handlers = make([]NotifyHookHandler, 0)

	return ntfhg
}

// IsHook 是否處理
func (ntfg *NotifyHookHandlerGroup) IsHook(bclsid int32, sclsid int32) bool {
	return true
}

// OnNotifyMessage 響應Notify消息
func (ntfg *NotifyHookHandlerGroup) OnNotifyMessage(msg *msg.Message, ctx *module.MouduleCtx, player *player.Manager) {

	for _, h := range ntfg.handlers {
		if h.IsHook(msg.BClassID, msg.SClassID) {
			h.OnNotifyMessage(msg, ctx, player)
		}
	}
}

// AddHandler 添加處理器
func (ntfg *NotifyHookHandlerGroup) AddHandler(h NotifyHookHandler) *NotifyHookHandlerGroup {
	ntfg.handlers = append(ntfg.handlers, h)

	return ntfg
}
