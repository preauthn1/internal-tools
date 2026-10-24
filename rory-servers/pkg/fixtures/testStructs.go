package fixtures

import (
	"os"
	"path"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverdefense"
	game_serverhandler "example.invalid/ServerTeam/Servers/pkg/gameserver/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// NewGameServerByTest test
func NewGameServerByTest(cfg []*network.ServerCfg, dbserveraddr string, nsserveraddr string, httpaddr string, serverid int, servertype int, servermode int, gameid int, defenseCfg *serverdefense.DefenseConfig, flagID string) *gameserver.GameServer {
	gs := new(gameserver.GameServer)
	// 固定服務器為單線程環境
	msgthrottle := 1

	gs.DBServerAddr = dbserveraddr
	gs.NSServerAddr = nsserveraddr
	gs.LocalAddr = cfg
	gs.HTTPServerAddr = httpaddr
	gs.IsRuning = false
	gs.ServerHandler = game_serverhandler.NewGameServerHandlerGroup()
	gs.ServerID = serverid
	gs.ServerType = servertype
	gs.ServerMode = servermode

	// 加入服務器邏輯處理器
	gs.ServerHandler.AddHandler(game_serverhandler.NewLogonCheckHandler(gameserver.GUESTCONNTIMEOUT))
	gs.ServerHandler.AddHandler(game_serverhandler.NewGameServerRegHandler(300, int32(serverid), int32(servertype)))
	// 註冊協調服務器
	gs.ServerHandler.AddHandler(game_serverhandler.NewNotifyHandler(nsserveraddr, serverid, servertype, gameid, 0, gs.LocalAddr[0].Addr, flagID))

	// 創建消息服務器
	svr := msgserver.NewFacadeMsgServer(1, 0)
	// 用戶管理器
	playerManager := player.NewPlayerManager(svr)
	// 模塊服務器
	ms := module.NewServer(svr)

	gs.GameServerRunCtx = new(game_serverhandler.GameServerCtx)

	gs.GameServerRunCtx.MsgServer = svr
	gs.GameServerRunCtx.ModuleServer = ms
	gs.GameServerRunCtx.PlayerManager = playerManager
	gs.GameServerRunCtx.TimerFuncManager = event.NewTimerFuncManager(svr.EventChannel)
	gs.GameServerRunCtx.Throttle = msgthrottle
	// 攔截RPC隊列
	gs.GameServerRunCtx.DBMsgQueue = NewDBMsgQueueByTest()

	return gs
}

// NewDBMsgQueueByTest 攔截RPC隊列
func NewDBMsgQueueByTest() *rpcmsg.Queue {
	rt := rpcmsg.NewRPCMsgQueue()
	return rt
}

// StartGameLogicByTest ?
func StartGameLogicByTest() (glm *gamelogic.LogicManager, room *gamelogic.Room) {

	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)

	// 獲取配置信息
	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
	servercfg.DealLogType = 0
	// gameserver偵聽地址
	gscfg := network.NewServerCfg(servercfg.ServerAddr, conf.TCPSERVERCONNROOMCLIENT, conf.TCPSERVERNAMEROOMCLIENT)

	servermode := cfg.GetConfigInt("defense", "servermode")
	defenseCfg := serverdefense.ParseDefenseConfig(cfg)
	gs := NewGameServerByTest([]*network.ServerCfg{gscfg}, servercfg.DBServerAddr, servercfg.NotifyServerAddr, servercfg.HTTPServerAddr, int(servercfg.ServerID), conf.ServerTypeGameServer, servermode, int(servercfg.GameID), defenseCfg, servercfg.Flag)

	// createRoom 創建房間對象
	rumanager := gamelogic.NewRoomUserManager(gs.GameServerRunCtx.PlayerManager)
	room = gamelogic.NewRoom(servercfg, rumanager)
	// 創建房間邏輯管理器
	glm = gamelogic.NewGameLogicManager(room, cfgs, gs.GameServerRunCtx.TimerFuncManager, roomserver.GameLogicCreator)
	// 房間邏輯模塊
	_ = roomlogic.NewRoomLogicModule(gs, room, glm, roomserver.RobotLogicCreator, cfgs)
	return
}
