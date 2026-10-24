package clientframework

import "example.invalid/ServerTeam/Servers/pkg/msgserver"

// CtxBase TCP
type CtxBase struct {
	ConnClt         *msgserver.TCPMsgClient
	Running         bool
	EnableHeartBeat bool
}

// NewCtxBase create CtxBase
func NewCtxBase(connClt *msgserver.TCPMsgClient, running bool, enableHeartBeat bool) *CtxBase {
	conn := new(CtxBase)
	conn.ConnClt = connClt
	conn.Running = running
	conn.EnableHeartBeat = enableHeartBeat
	return conn
}
