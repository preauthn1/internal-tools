package network

import "example.invalid/ServerTeam/Servers/pkg/msg"

// ServerWriter write message to server
type ServerWriter interface {
	WriteServer(*msg.Message)
}

// ClientWriter write message to client
type ClientWriter interface {
	WriteClient(*Message)
	WriteError(*Netty)
}
