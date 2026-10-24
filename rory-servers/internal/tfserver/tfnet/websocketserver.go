package tfnet

import (
	"net"
	"sync"
	"time"

	"example.invalid/gorilla/websocket"
	"example.invalid/ServerTeam/Servers/internal/tfserver/router"
	"example.invalid/ServerTeam/Servers/internal/tfserver/tfconfig"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"

	"net/http"
	"regexp"
	"strconv"
)

// CurrentConnInfos 現在的連線表
var CurrentConnInfos sync.Map

// GetCurrentWSConnCount 獲取當前WS連接數
func GetCurrentWSConnCount() int {
	length := 0

	CurrentConnInfos.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

// WSConnInfo 連線資訊
type WSConnInfo struct {
	remoteAddr  string
	clientConn  *websocket.Conn
	backendConn *network.TCPConn
}

// NewWSConnInfo 連線資訊實例
func NewWSConnInfo(clientConn *websocket.Conn) *WSConnInfo {
	return &WSConnInfo{
		clientConn.RemoteAddr().String(),
		clientConn,
		nil}
}

// Close 連線關閉
func (c *WSConnInfo) Close() {
	if c.clientConn != nil {
		c.clientConn.Close()
	}

	if c.backendConn != nil {
		c.backendConn.Close()
	}
	CurrentConnInfos.Delete(c.remoteAddr)
}

// WebSocketServer 服務器實作
type WebSocketServer struct {
	portProvider TFListenPortProvider
}

// NewWebSocketServer 實例化一個ws服務器
func NewWebSocketServer(portProvider TFListenPortProvider) *WebSocketServer {
	wss := new(WebSocketServer)
	wss.portProvider = portProvider

	return wss
}

// Start 啟動服務器
func (server *WebSocketServer) Start() {
	cfg := tfconfig.ConfigCtx
	upgrader := &websocket.Upgrader{
		//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	/*
		path example
		http://0.0.0.0/ws/10000
		http://0.0.0.0/ws/10101
	*/
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		clientConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			mlog.Error("upgrade to the WebSocket protocol fail : %v", err)
			return
		}
		re, _ := regexp.Compile("/ws/([0-9]+)")
		match := re.FindStringSubmatch(r.URL.Path)
		if len(match) < 1 {
			mlog.Error("path has not match %s error %v", re.String(), err)
			return
		}

		sourcePort, _ := strconv.Atoi(match[1])
		//ip, _ := util.ParseIPAndPort(r.RemoteAddr)
		rd := router.NewRouteData("0.0.0.0", sourcePort)
		//获取路由目标
		tg := router.RouteCtx.GetRouteTarget(rd)
		if tg == nil {
			mlog.Error("未找到匹配的路由目标 %s, %d.关闭通道.", r.RemoteAddr, sourcePort)
			return
		}

		go server.handleConn(clientConn, tg)
	})

	mlog.System("server start at :" + strconv.Itoa(cfg.WebSocketServerPort))
	mlog.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.WebSocketServerPort), nil).Error())
}

func (server *WebSocketServer) handleConn(clientConn *websocket.Conn, tg *router.RouteTarget) {
	v, _ := CurrentConnInfos.LoadOrStore(clientConn.RemoteAddr().String(), NewWSConnInfo(clientConn))
	connInfo := v.(*WSConnInfo)

	d := net.Dialer{Timeout: 3 * time.Second}
	backendConn, err := d.Dial("tcp", tg.String())
	if err != nil {
		mlog.Error("建立後端連線失敗 %s error %s", clientConn.RemoteAddr(), err.Error())
		connInfo.Close()
		return
	}
	backendTCPConn := network.NewTCPConn(backendConn, conf.READBUFLEN)
	connInfo.backendConn = backendTCPConn

	//處理服務器端轉發
	go func() {
		defer connInfo.Close()
		defer util.PanicErrStack()
		var legacy []byte
		//拆解消息函式
		msgSplitFunc := func(b []byte) (msgs [][]byte) {
			b = append(legacy, b...)
			legacy = legacy[:0]
			for {
				if len(b) == 0 {
					break
				}
				size := util.BytesToInt(b[0:4], true)
				if int(size+4) > len(b) {
					legacy = b
					mlog.Error("size %d, actual len %d", size+4, len(b))
					break
				}
				//mlog.Printf("size : %d, len : %d\n", size, len(b))
				msg := b[0 : size+4]
				b = b[size+4:]
				//msgid := fmt.Sprintf("%x:%x", binary.BigEndian.Uint32(msg[4:8]), binary.BigEndian.Uint32(msg[8:12]))
				//fmt.Printf("----ipaddr %v, msgid %v, size %v\n", clientConn.RemoteAddr(), msgid, fmt.Sprintf("%x", size))
				msgs = append(msgs, msg)
			}
			return
		}
		for {
			b, err := backendTCPConn.ReadMsg()
			if err != nil {
				mlog.Error("backend %s read error: %s and disconnect !!", backendConn.RemoteAddr(), err)
				break
			}
			//time.Sleep(100 * time.Millisecond)
			// 拆解消息
			msgs := msgSplitFunc(b)
			//if len(msgs) > 1 {
			//	fmt.Printf("---------合併噴發 %d-------------\n", len(msgs))
			//}

			// 转发消息
			for _, msg := range msgs {
				if err = clientConn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
					break
				}
				mlog.Debug("s2c ---%v", msg[:12])
			}
			if err != nil {
				mlog.Error("client %s write error: %s and disconnect !!", clientConn.RemoteAddr(), err)
				break
			}
		}
	}()

	//處理客戶端消息轉發
	go func() {
		defer connInfo.Close()
		defer util.PanicErrStack()
		for {
			mtype, msg, err := clientConn.ReadMessage()

			if mtype != websocket.BinaryMessage && mtype != -1 {
				mlog.Error("client %s message type error mtype [%d] and disconnect !!", clientConn.RemoteAddr(), mtype)
				break
			} else if mtype == -1 && err != nil {
				mlog.Error("read client message type as -1 and disconnect !! %v", err.Error())
				break
			}
			if err != nil {
				mlog.Error("client %s read error: %s and disconnect !!", clientConn.RemoteAddr(), err)
				break
			}
			msglen := util.BytesToInt(msg[:4], true)
			if msglen <= 0 {
				return
			}
			_, err = backendTCPConn.WriteMsg(msg)
			mlog.Debug("c2s ---%v", msg[:12])
			if err != nil {
				mlog.Error("backend %s write error: %s and disconnect !!", backendConn.RemoteAddr(), err)
				break
			}

		}
	}()
	//if err != nil {
	//	mlog.Error("建立后端路由通道错误,%s,%v.关闭通道", tg, err)
	//}

}
