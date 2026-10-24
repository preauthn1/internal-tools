package httpsvr

import (
	"net/http"

	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameHTTPServer http服務
type GameHTTPServer struct {
	addr        string
	eventChanel *event.Dispatcher
}

// NewGameHTTPServer 創建遊戲Http服務器
func NewGameHTTPServer(addr string, evtChannel *event.Dispatcher) *GameHTTPServer {
	svr := new(GameHTTPServer)
	svr.addr = addr
	svr.eventChanel = evtChannel

	return svr
}

// Start 启动服务器
func (svr *GameHTTPServer) Start() {
	http.HandleFunc("/game", svr.gameHandler)

	go func() {
		defer util.PanicErrStack()

		err := http.ListenAndServe(svr.addr, nil)
		if err != nil {
			mlog.Error("%v", err)
		}
	}()

	example.invalid("HTTP服務器啟動成功,偵聽在地址%s...", svr.addr)
}

// http事件響應接口
func (svr *GameHTTPServer) gameHandler(w http.ResponseWriter, req *http.Request) {
	mlog.Debug("接收到Http請求.%v", req)

	//构造http消息事件
	httpEvtData := make(map[string]interface{})

	callBackChan := util.NewSemaphore(0)
	httpEvtData[msgserver.HTTPMessageEventDataSEMAPHORE] = callBackChan
	httpEvtData[msgserver.HTTPMessageEventDataHTTPRESPONSEWRITER] = w
	httpEvtData[msgserver.HTTPMessageEventDataHTTPREQUEST] = req

	httpEvt := msgserver.NewHTTPMsgEvent(httpEvtData)
	svr.eventChanel.OnAsyncEvent(httpEvt, nil)

	callBackChan.Acquire()
	mlog.Debug("Http請求處理完成.%v", req)
}
