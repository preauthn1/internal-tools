package event

import (
	"strings"
)

// Handler 事件處理器
type Handler interface {
	// 同步觸發事件
	OnEvent(evt *Event, ret *Ret) bool

	// 處理器名稱
	HandlerName() string
}

// CallbackFun 回調函數
type CallbackFun func(evt *Event, ret *Ret)

// EmptyEventHandler 空的事件處理器
type EmptyEventHandler struct {
}

// OnEvent ?
func (h *EmptyEventHandler) OnEvent(evt *Event, ret *Ret) bool {
	return true
}

// HandlerName ?
func (h *EmptyEventHandler) HandlerName() string {
	return "EmptyEventHandler"
}

// CommHandler 公用的事件處理器
type CommHandler struct {
	handerName string
	handler    func(evt *Event, ret *Ret) bool
	hookEvts   []string
}

// OnEvent ?
func (h *CommHandler) OnEvent(evt *Event, ret *Ret) bool {
	for _, s := range h.hookEvts {
		if strings.ToLower(s) == strings.ToLower(evt.Name) {
			break
		}

		return true
	}

	return h.handler(evt, ret)
}

// HandlerName ?
func (h *CommHandler) HandlerName() string {
	return h.handerName
}

// NewCommonHandler ?
func NewCommonHandler(name string, f func(evt *Event, ret *Ret) bool, evts ...string) *CommHandler {
	ch := new(CommHandler)
	ch.handerName = name
	ch.handler = f
	ch.hookEvts = evts

	return ch
}
