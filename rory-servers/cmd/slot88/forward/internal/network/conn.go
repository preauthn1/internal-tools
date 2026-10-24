package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"

	"example.invalid/gorilla/websocket"
	"example.invalid/spf13/viper"

	"example.invalid/ServerTeam/Servers/cmd/slot88/auth/service"
	"example.invalid/ServerTeam/Servers/cmd/slot88/forward/internal/encrypt"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
)

// Conn slot88 Websocket 連線
type Conn struct {
	ctx    context.Context
	cancel context.CancelFunc

	// Client Connection
	ws  *websocket.Conn
	aes *encrypt.AES

	msgReadChan  chan *Message
	msgWriteChan chan *Message
	msgErrorChan chan *Netty

	srvMsgReadChan  chan *msg.Message
	srvMsgWriteChan chan *msg.Message

	// informations
	info *example.invalid

	// Server Connection
	srvConn SrvConn

	// 訊息消耗者
	handlers []Handler
}

// NewConn create slot88 connection
func NewConn(ws *websocket.Conn, handlers ...Handler) (*Conn, *example.invalid, error) {

	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("New Conn Error from client: %s, %v, %s", ws.RemoteAddr().String(), rec, debug.Stack())
		}
	}()

	conn := new(Conn)
	conn.ctx, conn.cancel = context.WithCancel(context.Background())
	conn.ws = ws
	conn.handlers = handlers

	example.invalid = new(example.invalid)

	var token string

	// read first text message as token
	{
		_, p, err := ws.ReadMessage()
		if err != nil {
			ws.Close()
			return nil, nil, fmt.Errorf("Read Token Error from client: %s, %v", ws.RemoteAddr().String(), err)
		}

		token = [REDACTED]

		aes, err := encrypt.NewAES(token)
		if err != nil {
			ws.Close()
			return nil, nil, fmt.Errorf("Create AES Error from client: %s with token %s, %v", ws.RemoteAddr().String(), token, err)
		}

		conn.aes = aes
	}

	// 取得了 GameID, ServerID, UserID, ServerAddr
	{
		host := viper.GetString("auth.endpoint")
		q := url.Values{
			"mD5": []string{token},
		}

		u := host + "?" + q.Encode()

		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			ws.Close()
			return nil, nil, fmt.Errorf("Auth Request Error from client: %s, %v", ws.RemoteAddr().String(), err)
		}

		req.Header.Add("X-API-Certificate", viper.GetString("auth.cer"))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ws.Close()
			return nil, nil, fmt.Errorf("Auth Error from client: %s, %v", ws.RemoteAddr().String(), err)
		}

		if resp.StatusCode != 200 {
			ws.Close()
			return nil, nil, fmt.Errorf("Auth Error from client: %s, %d", ws.RemoteAddr().String(), resp.StatusCode)
		}

		dec := json.NewDecoder(resp.Body)

		t := struct {
			Data example.invalid `json:"data"`
		}{}

		if err := dec.Decode(&t); err != nil {
			ws.Close()
			return nil, nil, fmt.Errorf("Parse Auth Message Error from client: %s, %v", ws.RemoteAddr().String(), err)
		}

		example.invalid = &t.Data
	}

	// 建立game server通道
	{
		var srv SrvConn
		switch example.invalid.ServerInfo.ConnType {
		case ConnTypeTCP:
			s, err := NewTCPServer(example.invalid.ServerInfo.ServerConnection)
			if err != nil {
				ws.Close()
				return nil, nil, fmt.Errorf("Connect Game Server Error from client: %s to server %s, %v", ws.RemoteAddr().String(), example.invalid.ServerInfo.ServerConnection, err)
			}
			srv = s
		case ConnTypeWS:
			s, err := NewWebsocketServerConext(conn.ctx, example.invalid.ServerInfo.ServerConnection)
			if err != nil {
				ws.Close()
				return nil, nil, fmt.Errorf("Connect Game Server Error from client: %s to server %s, %v", ws.RemoteAddr().String(), example.invalid.ServerInfo.ServerConnection, err)
			}
			srv = s
		}

		if srv == nil {
			return nil, nil, fmt.Errorf("Create Game Server Error from client: %s to server %s", ws.RemoteAddr().String(), example.invalid.ServerInfo.ServerConnection)
		}

		conn.srvConn = srv
		conn.srvMsgReadChan = make(chan *msg.Message, 50)
		conn.srvMsgWriteChan = make(chan *msg.Message, 50)

		go conn.srvWriter()
		go conn.srvDispatcher()
		go srv.Listen(conn.ctx, conn.srvMsgReadChan, conn.cancel)
	}

	conn.msgReadChan = make(chan *Message, 50)
	conn.msgWriteChan = make(chan *Message, 50)
	conn.msgErrorChan = make(chan *Netty, 50)

	go conn.writer()
	go conn.dispatcher()

	return conn, example.invalid, nil
}

// Context from *http.Request
func (conn *Conn) Context() context.Context {
	return conn.ctx
}

// ClientConn get client connection
func (conn *Conn) ClientConn() *websocket.Conn {
	return conn.ws
}

// SrvConn get game server connection
func (conn *Conn) SrvConn() SrvConn {
	return conn.srvConn
}

// Close connection
func (conn *Conn) Close() {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Error("Close Server Error, %v", rec)
		}
	}()

	example.invalid("Close Connection from %s", conn.ws.RemoteAddr().String())

	// conn.cancel()
	conn.ws.Close()
	conn.srvConn.Close()
	close(conn.msgReadChan)
	close(conn.msgWriteChan)
	close(conn.msgErrorChan)
	close(conn.srvMsgReadChan)
	close(conn.srvMsgWriteChan)
}

// AddHandler add handler to conn
func (conn *Conn) AddHandler(c Handler) {
	conn.handlers = append(conn.handlers, c)
}

// Listen recive request from client
func (conn *Conn) Listen() {

	defer conn.cancel()

	for {
		t, b, err := conn.ws.ReadMessage()
		if err != nil {
			break
		}

		switch t {
		default:
			break
		case websocket.TextMessage:
			mlog.Trace("[%s] Read: %s", conn.ws.RemoteAddr().String(), b)

			// 文字內容，解密
			b = conn.aes.Decrypt(b)
			// mlog.Debug("Decrypted Bytes: %v", b)
			mlog.Trace("[%s] Read Decrypted: %s", conn.ws.RemoteAddr().String(), b)
			msg := new(Message)
			if err := json.Unmarshal(b, msg); err != nil {
				mlog.Error("Unmarshal Test Message Error, %v", err)
				continue
			}
			msg.Data = b

			conn.msgReadChan <- msg
		}
	}
}

// WriteClient write message to client
func (conn *Conn) WriteClient(m *Message) {
	b, _ := json.Marshal(m)
	mlog.Trace("[%s] Write: %s", conn.ws.RemoteAddr().String(), b)
	conn.msgWriteChan <- m
}

// WriteError write error to client
func (conn *Conn) WriteError(err *Netty) {
	conn.msgErrorChan <- err
}

// WriteServer write message to client
func (conn *Conn) WriteServer(m *msg.Message) {
	conn.srvMsgWriteChan <- m
}

func (conn *Conn) writer() {

	for {
		select {
		case <-conn.ctx.Done():
			mlog.Trace("Writer Context Done")
			return
		case msg := <-conn.msgWriteChan:
			if msg == nil {
				continue
			}
			b, err := json.Marshal(msg)
			if err != nil {
				mlog.Error("Marshal Test Message Error, %v", err)
			}
			b = conn.aes.Encrypt(b)
			if err := conn.ws.WriteMessage(websocket.TextMessage, b); err != nil {
				mlog.Error("WriteMessage Error, %v", err)
			}
		case err := <-conn.msgErrorChan:
			if err == nil {
				continue
			}
			s := fmt.Sprintf("netty::%s", err.Data)

			// if err.Data != nil {
			// 	s = fmt.Sprintf("%s>>%v", s, err.Data)
			// }

			var b []byte
			if err.Encrypt {
				b = conn.aes.Encrypt([]byte(s))
			} else {
				b = []byte(s)
			}
			if err := conn.ws.WriteMessage(websocket.TextMessage, b); err != nil {
				mlog.Error("WriteMessage Error, %v", err)
			}
		}
	}
}

func (conn *Conn) dispatcher() {
	for {
		select {
		case <-conn.ctx.Done():
			mlog.Trace("Dispathcer Context Done")
			return
		case msg := <-conn.msgReadChan:
			if msg == nil {
				continue
			}

			// 送給所有handlers
			for _, c := range conn.handlers {
				c.OnMessage(msg)
			}
		}
	}
}

func (conn *Conn) srvWriter() {

	for {
		select {
		case <-conn.ctx.Done():
			mlog.Trace("Srv Writer Context Done")
			return
		case m := <-conn.srvMsgWriteChan:
			conn.srvConn.Write(m)
		}
	}
}

func (conn *Conn) srvDispatcher() {
	for {
		select {
		case <-conn.ctx.Done():
			mlog.Trace("Srv Dispathcer Context Done")
			return
		case m := <-conn.srvMsgReadChan:
			if m == nil {
				continue
			}
			for _, c := range conn.handlers {
				c.OnSrvMessage(m)
			}
		}
	}
}
