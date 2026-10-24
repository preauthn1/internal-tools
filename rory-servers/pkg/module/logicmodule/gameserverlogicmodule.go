package logicmodule

import (
	"net/http"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/module/heartbeat"
	"example.invalid/ServerTeam/Servers/pkg/module/info"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// HTTPServerConfig http服務器配置
type HTTPServerConfig struct {
	ServerAddr string // 服務器地址
	VCode      string // http校驗碼
}

// GameServerLogicModule 服務器邏輯模塊
type GameServerLogicModule struct {
	MsgHandler       *module.HookMsgHandlerGroup
	RPCMsgHandler    *RPCHookHandlerGroup
	NotifyMsgHandler *NotifyHookHandlerGroup
	HTTPMsgHandler   *HTTPMsgHookHandlerGroup
	PlayerManager    *player.Manager
	DBMsgQueue       *rpcmsg.Queue
	ModuleCtx        *module.MouduleCtx
	MsgServer        *msgserver.MsgServer
	heartBeater      *heartbeat.Manager
	ServerStatus     ServerStatusProvider
}

// NewGameServerLogicModule 創建服務器邏輯模塊
func NewGameServerLogicModule() *GameServerLogicModule {
	sm := new(GameServerLogicModule)
	// rpc 消息處理器
	sm.RPCMsgHandler = NewRPCHookHandlerGroup()
	// 消息處理器
	sm.MsgHandler = module.NewHookMsgHandlerGroup()
	// 通知消息處理器
	sm.NotifyMsgHandler = NewNotifyHookHandlerGroup()
	// http消息處理器
	sm.HTTPMsgHandler = NewHTTPMsgHookHandlerGroup()

	// 服務器狀態
	ssnh := NewServerStopNotifyHandler()
	sm.ServerStatus = ssnh

	return sm
}

// Init 初始化
func (sm *GameServerLogicModule) Init(ctx *module.MouduleCtx, playerManager *player.Manager, dbmsgQueue *rpcmsg.Queue, httpCfg *HTTPServerConfig) {
	sm.ModuleCtx = ctx

	sm.PlayerManager = playerManager
	sm.DBMsgQueue = dbmsgQueue
	sm.MsgServer = ctx.Server

	// 註冊公共消息處理器
	sm.registerCommonMsgHandler()
	// 註冊公共通知處理器
	sm.registerCommNotifyHandler()
	// 註冊公共的http消息處理器
	sm.registerCommonHTTPMsgHandler(httpCfg)

	// 註冊消息解析器
	sm.registerMsgParser(ctx)
	// 註冊消息事件處理器
	sm.registerMsgEvtHandler(ctx)
	// 註冊RPC消息事件處理器
	sm.registerRPCMsgEventHandler(ctx)
	// 註冊Notify消息事件處理器
	sm.registerNotifyMsgEventHandler(ctx)
	// 註冊異步調用事件處理器
	sm.registerAsyncCallBackEventHandler(ctx)
	// 註冊http消息事件處理器
	sm.registerHTTPMsgEventHandler(ctx)

	// 處理RPC 消息
	go sm.rpcCallBackMsgDispatcher(ctx)

}

// 註冊公共的消息處理器
func (sm *GameServerLogicModule) registerCommonMsgHandler() {
	sm.MsgHandler.AddHandler(heartbeat.NewReqBeatConfigMsgHandler()) //请求心跳包配置
	sm.heartBeater = heartbeat.NewHeartBeatManager(sm.MsgServer)

	sm.MsgHandler.AddHandler(sm.heartBeater) //心跳包

	sm.MsgHandler.AddHandler(info.NewServerInfoHandler())
}

// 註冊公共的通知處理器
func (sm *GameServerLogicModule) registerCommNotifyHandler() {
	sm.NotifyMsgHandler.AddHandler(sm.ServerStatus.(*ServerStopNotifyHandler))
}

// 註冊公共的http消息處理器
func (sm *GameServerLogicModule) registerCommonHTTPMsgHandler(httpCfg *HTTPServerConfig) {
	if httpCfg != nil && httpCfg.ServerAddr != "" {
		//http校验码
		sm.HTTPMsgHandler.AddHandler(NewValidateHTTPMsgHookHandler(httpCfg.VCode))
	}
}

func (sm *GameServerLogicModule) registerMsgParser(ctx *module.MouduleCtx) {
	ctx.RegisterMsgParser(msg.AllMsgBClassID, msg.AllMsgSClassID, msg.NewProtoMsgParser())
}

func (sm *GameServerLogicModule) registerMsgEvtHandler(ctx *module.MouduleCtx) {
	ctx.OnMessage(func(clt *network.ConnInfo, m *msg.Message) bool {
		//example.invalid("client %s msg:%d,%d", clt, m.BClassID, m.SClassID)

		// 檢查協議內容
		if !sm.checkClientMsgProto(m) {
			ctx.Error("客戶端發送消息%v,但不在允許的通訊協議中.關閉客戶端連接", m)
			ctx.Server.CloseConn(clt)
			return false
		}

		retmsgs := sm.MsgHandler.OnMessage(ctx, clt, m)

		if retmsgs != nil {
			for _, retmsg := range retmsgs {
				if retmsg != nil {
					ctx.Server.SendMessage(clt, retmsg)
				}
			}
		}

		return true
	})
	ctx.OnDisConnect(func(clt *network.ConnInfo) bool {
		// 如果是大廳或者房間的客戶端連接
		if clt.IsServerConn(conf.TCPSERVERNAMEHALLCLIENT) || clt.IsServerConn(conf.TCPSERVERNAMEROOMCLIENT) {
			sm.heartBeater.OnClientDisconnect(clt.ConnID())
		}
		return true
	})
	ctx.OnConnect(func(clt *network.ConnInfo) bool {
		if clt.IsServerConn(conf.TCPSERVERNAMEHALLCLIENT) || clt.IsServerConn(conf.TCPSERVERNAMEROOMCLIENT) {
			sm.heartBeater.OnClientConnect(clt.ConnID())
		}

		return true
	})
}

func (sm *GameServerLogicModule) registerRPCMsgEventHandler(ctx *module.MouduleCtx) {
	h := event.NewCommonHandler("rpcmsg", func(evt *event.Event, evtret *event.Ret) bool {
		m := msgserver.GetRPCMsg(evt.Args)

		sm.RPCMsgHandler.OnRPCMessage(m, ctx)

		return true
	}, msgserver.EventRPCMessage)

	ctx.EvtDispatcher.PushHandler(h)
}

func (sm *GameServerLogicModule) registerNotifyMsgEventHandler(ctx *module.MouduleCtx) {
	h := event.NewCommonHandler("notifymsg", func(evt *event.Event, evtret *event.Ret) bool {
		m := msgserver.GetNotifyMsg(evt.Args)

		sm.NotifyMsgHandler.OnNotifyMessage(m, ctx, sm.PlayerManager)

		return true
	}, msgserver.EventNotifyMessage)

	ctx.EvtDispatcher.PushHandler(h)
}

func (sm *GameServerLogicModule) registerAsyncCallBackEventHandler(ctx *module.MouduleCtx) {
	h := event.NewCommonHandler("asynccall", func(evt *event.Event, evtret *event.Ret) bool {
		callbackFunc := msgserver.GetAsyncEventCallBackFunc(evt.Args)

		if callbackFunc != nil {
			util.SafeCall(callbackFunc)
		}

		return true
	}, msgserver.EventAsyncCall)

	ctx.EvtDispatcher.PushHandler(h)
}

func (sm *GameServerLogicModule) registerHTTPMsgEventHandler(ctx *module.MouduleCtx) {
	h := event.NewCommonHandler("httpmsg", func(evt *event.Event, evtret *event.Ret) bool {
		httpMsgData := msgserver.GetHTTPMsgData(evt.Args)

		seamphore := httpMsgData[msgserver.HTTPMessageEventDataSEMAPHORE].(util.Semaphore)
		req := httpMsgData[msgserver.HTTPMessageEventDataHTTPREQUEST].(*http.Request)
		w := httpMsgData[msgserver.HTTPMessageEventDataHTTPRESPONSEWRITER].(http.ResponseWriter)

		sm.HTTPMsgHandler.OnHTTPMsg(seamphore, w, req)

		return true
	}, msgserver.EventHTTPMessage)

	ctx.EvtDispatcher.PushHandler(h)
}

func (sm *GameServerLogicModule) rpcCallBackMsgDispatcher(ctx *module.MouduleCtx) {
	defer util.PanicErrStack()

	for {
		m := <-sm.DBMsgQueue.OnRPCMesasge
		ctx.Server.OnRPCMessageEvent(m)
	}
}

// 檢查客戶端協議號
func (sm *GameServerLogicModule) checkClientMsgProto(m *msg.Message) bool {
	// 允許的一級協議號
	allowBClsids := []int32{int32(netproto.MessageBClassID_Common), int32(netproto.MessageBClassID_Hall), int32(netproto.MessageBClassID_Game), int32(netproto.MessageBClassID_GameRoom), int32(netproto.MessageBClassID_ServerMgm)}
	for _, aid := range allowBClsids {
		if m.BClassID == aid {
			return true
		}
	}

	return false
}
