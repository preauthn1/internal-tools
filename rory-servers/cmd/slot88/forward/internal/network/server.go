package network

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"example.invalid/gorilla/websocket"
	"example.invalid/spf13/viper"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// SrvConn to game server connection
type SrvConn interface {
	Close()
	Write(*msg.Message)
	Listen(ctx context.Context, m chan *msg.Message, cancel context.CancelFunc)
}

// TCPServer use tcp
type TCPServer struct {
	evt    *msgserver.MsgEvtChan
	client *msgserver.TCPMsgClient
}

// NewTCPServer create TCPServer
func NewTCPServer(addr string) (*TCPServer, error) {

	packer := network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
	evt := msgserver.NewMsgEvtChan(conf.EVTCHANSIZE)
	parser := msg.NewProtoMsgParser()

	client := msgserver.NewTCPMsgClient(packer, evt, parser, false)

	if err := client.Connect(addr); err != nil {
		return nil, fmt.Errorf("Connect Game Server Error, %v", err)
	}

	srv := new(TCPServer)
	srv.evt = evt
	srv.client = client

	return srv, nil
}

// Close server
func (srv *TCPServer) Close() {
	if srv.client == nil {
		return
	}
	srv.client.Close()
}

// Write send message
func (srv *TCPServer) Write(m *msg.Message) {
	if srv.client == nil {
		return
	}
	if m == nil {
		return
	}
	srv.client.SendMessage(m)
}

// Listen message from server
func (srv *TCPServer) Listen(ctx context.Context, c chan *msg.Message, cancel context.CancelFunc) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				mlog.Debug("TCP Server Listen Done")
				return
			case m := <-srv.evt.MsgSig:
				if m == nil {
					continue
				}
				c <- m
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				mlog.Debug("TCP Server Disconnected Done")
				return
			case <-srv.evt.DisconnectedSig:
				mlog.Debug("TCP Server Disconnected")
				cancel()
				return
			}
		}
	}()
}

// WebsocketServer use websocket
type WebsocketServer struct {
	conn *websocket.Conn
	resp *http.Response

	packer *network.HeadMsgPacker
	parser *msg.ProtoParser
}

// NewWebsocketServer create WebsocketServer
func NewWebsocketServer(addr string) (*WebsocketServer, error) {
	return NewWebsocketServerConext(context.Background(), addr)
}

// NewWebsocketServerConext create WebsocketServer
func NewWebsocketServerConext(ctx context.Context, addr string) (*WebsocketServer, error) {

	u, _ := url.Parse(addr)
	u.Scheme = viper.GetString("scheme")
	if u.Scheme != "ws" && u.Scheme != "wss" {
		u.Scheme = "ws"
	}

	websocket.DefaultDialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Connect Game Server Error to %s, %v", u.String(), err)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		example.invalid("Server Connection Closed")
		return nil
	})

	srv := new(WebsocketServer)
	srv.conn = conn
	srv.resp = resp
	srv.packer = network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
	srv.parser = msg.NewProtoMsgParser()

	return srv, nil
}

// Close server
func (srv *WebsocketServer) Close() {
	if srv.conn == nil {
		return
	}
	srv.conn.Close()
}

// Write send message
func (srv *WebsocketServer) Write(m *msg.Message) {
	if srv.conn == nil || m == nil {
		return
	}
	b, _ := srv.parser.Marshal(m)
	srv.conn.WriteMessage(websocket.BinaryMessage, srv.packer.Packet(b))
}

// Listen message from server
func (srv *WebsocketServer) Listen(ctx context.Context, c chan *msg.Message, cancel context.CancelFunc) {
	go func() {

		defer func() {
			if rec := recover(); rec != nil {
				mlog.Error("Listen Conn Error, %v", rec)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				example.invalid("Websocket Server Listen Done")
				return
			default:
				t, b, err := srv.conn.ReadMessage()
				if err != nil {
					cancel()
					break
				}
				switch t {
				default:
					break
				case websocket.TextMessage, websocket.BinaryMessage:

					data, err := srv.packer.UnPack(b)
					if err != nil {
						mlog.Warn("Websocket Server Unpack Data Error, %v", err)
						continue
					}

					m, err := srv.parser.Unmarshal(data)
					if err != nil {
						mlog.Warn("Websocket Server Unmarshal Message Error, %v", err)
						continue
					}

					c <- m
				}
			}
		}

	}()

}
