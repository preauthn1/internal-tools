package gameserver

import (
	"errors"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/httpsvr"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverdefense"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverdefense/blackip"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/supertool"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// Consts
const (
	GUESTCONNTIMEOUT = 60
	SUPERTOOLDYNAMIC = 1 // 配桌型
)

const ()

// GameServer 遊戲服務器對象
type GameServer struct {
	DBServerAddr   string                    // dbserver address
	NSServerAddr   string                    // nsserver address
	LocalAddr      []*network.ServerCfg      // Listening address
	HTTPServerAddr string                    // http Service listening address
	SuperToolCgf   gamelogic.SuperToolConfig // 超端設定檔

	GameServerRunCtx *serverhandler.GameServerCtx          // 服務器運行環境
	ServerHandler    *serverhandler.GameServerHandlerGroup // 服務器組件
	IsRuning         bool                                  // 是否運行中
	ServerID         int                                   // 服務器ID
	ServerType       int                                   // 服務器類型
	ServerMode       int                                   // 服務器模式
	GameID           int                                   // 遊戲ID
	FlagID           string

	msgServer  *msgserver.MsgServer         // 消息服務器
	defenseCfg *serverdefense.DefenseConfig // 防護配置

	msgthrottle     int // 消息並發數
	rpciplstHandler *blackip.RPCIPBlackListHandler
	rpciplst        *blackip.RPCBlackIPList
	NotifyHandler   *serverhandler.NotifyHandler

	gameHTTPSvr      *httpsvr.GameHTTPServer // 遊戲http服務
	gameSuperToolSvr supertool.ISuperTool    // 遊戲超端服務
}

// NewGameServer 創建gameserver對象
func NewGameServer(cfg []*network.ServerCfg, dbserveraddr string, nsserveraddr string, httpaddr string, serverid int, servertype int, servermode int, gameid int, defenseCfg *serverdefense.DefenseConfig, hostAndBackup int32, superToolCfg gamelogic.SuperToolConfig, flagID string) *GameServer {
	// 固定服務器為單線程環境
	msgthrottle := 1

	gs := new(GameServer)
	gs.DBServerAddr = dbserveraddr
	gs.NSServerAddr = nsserveraddr
	gs.LocalAddr = cfg
	gs.HTTPServerAddr = httpaddr
	gs.SuperToolCgf = superToolCfg

	gs.IsRuning = false
	gs.ServerHandler = serverhandler.NewGameServerHandlerGroup()
	gs.ServerID = serverid
	gs.ServerType = servertype
	gs.ServerMode = servermode
	gs.defenseCfg = defenseCfg
	gs.FlagID = flagID

	// 加入服務器邏輯處理器
	gs.ServerHandler.AddHandler(serverhandler.NewLogonCheckHandler(GUESTCONNTIMEOUT))
	gs.ServerHandler.AddHandler(serverhandler.NewGameServerRegHandler(300, int32(serverid), int32(servertype)))

	// 註冊協調服務器
	ntfHandler := serverhandler.NewNotifyHandler(nsserveraddr, serverid, servertype, gameid, hostAndBackup, gs.LocalAddr[0].Addr, flagID)
	command.AddCommand(serverhandler.NewNotifyHandlerCmd(ntfHandler))
	gs.ServerHandler.AddHandler(ntfHandler)
	gs.NotifyHandler = ntfHandler
	gs.msgthrottle = msgthrottle

	return gs
}

// Init 初始化
func (gs *GameServer) Init() error {
	if gs.GameServerRunCtx != nil {
		return errors.New("服務器已經初始化")
	}

	// 創建消息服務器
	svr := msgserver.NewFacadeMsgServer(gs.msgthrottle, gs.ServerMode)
	gs.msgServer = svr

	// 加入TCP服務器
	for _, adr := range gs.LocalAddr {
		svr.AddTCPServer(adr)
	}

	// db消息隊列
	rpcdbqueue := rpcmsg.NewRPCMsgQueue()

	err := rpcdbqueue.Connect(gs.DBServerAddr)

	if err != nil {
		mlog.Error("連接到 RPC Message Queue 失敗,%v,服務器啟動失敗", err)
		return err
	}
	mlog.System("連接 RPC Message Queue 成功.")
	rpcdbqueue.Start()
	mlog.System("RPC Message Queue 啟動.")

	// 用戶管理器
	playerManager := player.NewPlayerManager(svr)
	// 模塊服務器
	ms := module.NewServer(svr)

	gs.GameServerRunCtx = new(serverhandler.GameServerCtx)

	gs.GameServerRunCtx.DBMsgQueue = rpcdbqueue
	gs.GameServerRunCtx.MsgServer = svr
	gs.GameServerRunCtx.ModuleServer = ms
	gs.GameServerRunCtx.PlayerManager = playerManager
	gs.GameServerRunCtx.TimerFuncManager = event.NewTimerFuncManager(svr.EventChannel)
	gs.GameServerRunCtx.Throttle = gs.msgthrottle

	// 初始化服務器防護方案
	gs.initServerDefense()

	// 加入rpc ip黑名單命令
	command.AddCommand(blackip.NewRPCBlackIPListCommand(gs.rpciplst))

	// 初始化Http服務
	if gs.HTTPServerAddr != "" {
		gs.gameHTTPSvr = httpsvr.NewGameHTTPServer(gs.HTTPServerAddr, gs.msgServer.EventChannel)
	}

	// 初始化超端服務
	if gs.SuperToolCgf.SuperToolModeType == SUPERTOOLDYNAMIC {
		gs.gameSuperToolSvr = supertool.NewDynamicSuperTool(rpcdbqueue, gs.ServerID, gs.GameID, gs.FlagID)
	}

	if app.Context != nil {
		if err := app.Context.Provide(func() *module.Server {
			return ms
		}); err != nil {
			mlog.Error("*module.ModuleServer Register dig error, %v", err)
		}
		if err := app.Context.Provide(func() *player.Manager {
			return playerManager
		}); err != nil {
			mlog.Error("*player.PlayerManager Register dig error, %v", err)
		}
		if err := app.Context.Provide(func() *msgserver.MsgServer {
			return svr
		}); err != nil {
			mlog.Error("*msgserver.MsgServer Register dig error, %v", err)
		}
		if err := app.Context.Provide(func() *rpcmsg.Queue {
			return rpcdbqueue
		}); err != nil {
			mlog.Error("*rpcmsg.RPCMsgQueue Register dig error, %v", err)
		}
		if err := app.Context.Provide(func() *event.TimerFuncManager {
			return gs.GameServerRunCtx.TimerFuncManager
		}); err != nil {
			mlog.Error("*event.TimerFuncManager Register dig error, %v", err)
		}
	}

	return nil

}

// Start 啟動服務器
func (gs *GameServer) Start() error {
	if gs.IsRuning {
		return errors.New("服務器已經啟動")
	}

	err := gs.GameServerRunCtx.ModuleServer.Run()
	if err != nil {
		return err
	}

	gs.ServerHandler.Start(gs.GameServerRunCtx)
	gs.IsRuning = true

	// 啟動http服務
	if gs.gameHTTPSvr != nil {
		gs.gameHTTPSvr.Start()
	}

	// 啟動超端服務
	if gs.gameSuperToolSvr != nil {
		gs.gameSuperToolSvr.Start()
	}

	return nil
}

// Stop 結束
func (gs *GameServer) Stop() {
	if !gs.IsRuning {
		return
	}

	gs.IsRuning = false
	gs.ServerHandler.Stop()
	gs.GameServerRunCtx.DBMsgQueue.Stop()
	gs.GameServerRunCtx.ModuleServer.Stop()
	gs.GameServerRunCtx.MsgServer.Stop()

}

// 初始化服務器防護
func (gs *GameServer) initServerDefense() {
	if !gs.defenseCfg.Enable {
		return
	}

	blockiptime := gs.defenseCfg.BlockIPTime
	ccMaxconn := gs.defenseCfg.MaxConn
	ccConntimes := gs.defenseCfg.ConnTimesLimit
	ccSenddata := gs.defenseCfg.SendDataLimint
	ccSendmsg := gs.defenseCfg.SendMsgLimit

	protoDefault := gs.defenseCfg.ProtoLimitDefault
	protoSpecial := gs.defenseCfg.ProtoLimitSpecial

	whiteiplist := gs.defenseCfg.WhiteIPList

	// 當前網絡連接遍歷
	connprovider := serverdefense.NewServerConnListProvider(gs.msgServer)

	// ip黑名單
	// rpc ip黑名單處理器
	gs.rpciplstHandler = blackip.NewRPCIPBlackListHandler(gs.GameServerRunCtx.DBMsgQueue, int32(gs.ServerID), gs.defenseCfg.ServerDesc)

	loadipfunc := func(f func(map[string]string)) {
		gs.rpciplstHandler.LoadBlockIPList(f)
	}
	blockipfunc := func(ip string, sec int, memo string) {
		gs.rpciplstHandler.BlockIP(ip, sec, memo)
	}
	iplistprovider := blackip.NewRPCBlackIPList(60, loadipfunc, blockipfunc)
	gs.rpciplst = iplistprovider

	// cc 防護
	ccdef := serverdefense.NewCCDefense(ccMaxconn, nil, nil, nil)
	ccdef.ParseConnTimesLimit(ccConntimes)
	ccdef.ParseSendDataLimit(ccSenddata)
	ccdef.ParseSendMsgLimit(ccSendmsg)
	ccdef.BlockIPTime = blockiptime
	ccdef.WhiteIPList = whiteiplist

	ccdef.SetConnProvider(connprovider)
	ccdef.SetIPListProvider(iplistprovider)
	mlog.System("初始化CC防護成功")

	// 協議防護
	protodef := serverdefense.NewMsgProtoDefense(nil, nil)
	protodef.ParseDefaultProtoLimit(protoDefault)
	protodef.ParseSpecialProtoLimit(protoSpecial)
	protodef.BlockIPTime = blockiptime
	protodef.WhiteIPList = whiteiplist

	protodef.SetConnProvider(connprovider)
	protodef.SetIPListProvider(iplistprovider)
	mlog.System("初始化協議防護成功")

	// 黑名單防護
	blockipdef := serverdefense.NewBlackIPDefense()
	blockipdef.WhiteIPList = whiteiplist

	blockipdef.SetConnProvider(connprovider)
	blockipdef.SetIPListProvider(iplistprovider)
	mlog.System("初始化黑名單防護成功")

	// 註冊消息服務器邏輯組件
	gs.msgServer.GetLogicHandlerGroup().RegisterHandler(ccdef).RegisterHandler(protodef).RegisterHandler(blockipdef)

}

// InitRPCBlackIPList 初始化rpc ip黑名單
func (gs *GameServer) InitRPCBlackIPList(rpchkHGroup *logicmodule.RPCHookHandlerGroup) {
	if gs.defenseCfg.Enable {
		rpchkHGroup.AddHandler(gs.rpciplstHandler)
		gs.rpciplst.StartLoadIPList()
	}
}
