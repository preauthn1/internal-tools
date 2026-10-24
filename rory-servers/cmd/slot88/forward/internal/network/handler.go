package network

import (
	"example.invalid/ServerTeam/Servers/pkg/msg"
)

// Handler 要處理訊息的人
type Handler interface {
	OnMessage(*Message)
	OnSrvMessage(*msg.Message)
}
