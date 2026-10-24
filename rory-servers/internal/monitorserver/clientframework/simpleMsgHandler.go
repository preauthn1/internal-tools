package clientframework

import "example.invalid/ServerTeam/Servers/pkg/msg"

// SimpleMsgHandler a simple msg handler
type SimpleMsgHandler struct {
	bclsid   int32
	sclsid   int32
	callback func(m *msg.Message)
}

// NewSimpleMsgHandler create SimpleMsgHandler
func NewSimpleMsgHandler(bclsid int32, sclsid int32, callback func(m *msg.Message)) *SimpleMsgHandler {
	h := new(SimpleMsgHandler)
	h.bclsid = bclsid
	h.sclsid = sclsid
	h.callback = callback
	return h
}

// IsHook check is match
func (h *SimpleMsgHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == h.bclsid && sclsid == h.sclsid
}

// OnMessage on handler message
func (h *SimpleMsgHandler) OnMessage(client *CtxBase, m *msg.Message) {
	if h.callback != nil {
		h.callback(m)
	}
}

// SetCallback setup call back function
func (h *SimpleMsgHandler) SetCallback(f func(m *msg.Message)) {
	h.callback = f
}
