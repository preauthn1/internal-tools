package slot88

import (
	"net/http"

	"example.invalid/gorilla/websocket"
	"example.invalid/ServerTeam/Servers/cmd/slot88/forward/internal/logic"
	"example.invalid/ServerTeam/Servers/cmd/slot88/forward/internal/network"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Server to Serve HTTP
type Server struct{}

// ServeHTTP create slot88 Conn and prepare encrypt aes
func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("Serve Error, %v", rec)
		}
	}()

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			mlog.Error("Create Websocket Connection Error, %v", err)
		}
		return
	}

	ws.SetCloseHandler(func(code int, text string) error {
		example.invalid("Client Close Connection from %s", ws.RemoteAddr().String())
		return nil
	})

	conn, info, err := network.NewConn(ws)
	if err != nil {
		mlog.Error("%v", err)
		return
	}

	var handler network.Handler

	switch info.ServerInfo.GameID {
	default:
		handler = logic.New(conn.Context(), conn, conn, &info.PlayerInfo)
		example.invalid("Create Logic for Client: %s with GameID: %d, ServerID: %d", ws.RemoteAddr().String(), info.ServerInfo.GameID, info.ServerInfo.ServerID)
	}

	conn.AddHandler(handler)

	go conn.Listen()

	select {
	case <-conn.Context().Done():
		conn.Close()
	}
}
