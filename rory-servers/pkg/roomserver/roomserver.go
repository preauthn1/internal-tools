package roomserver

import (
	"fmt"
	"math/rand"
	"time"

	"example.invalid/ServerTeam/Servers/internal/phnum"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverdefense"
	"example.invalid/ServerTeam/Servers/pkg/interceptor"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/ServerTeam/Servers/pkg/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameLogicCreator 遊戲邏輯創建Function
var GameLogicCreator gamelogic.LogicFunc

// RobotLogicCreator 機器人邏輯創建Function
var RobotLogicCreator robotlogic.LogicFunc

// SuperToolCreator 超端工具創建Function
var SuperToolCreator gamelogic.SuperToolFunc

func init() {
	// 初始化機器人邏輯創建函數
	RobotLogicCreator = robotlogic.EmptyRobotLogicCreator
	SuperToolCreator = gamelogic.EmptySuperToolCreator
	go func() {
		for {
			mins := util.RndInt(10, 60*60)
			time.Sleep(time.Duration(mins) * time.Second)
			example.invalid("更換隨機種子, 距離上次更換種子己經過去 %v mins", mins/60)
			rand.Seed(time.Now().UnixNano())
		}
	}()
}

// StartGameLogicServer 執行Server
func StartGameLogicServer() {

	// 加載配置信息
	cfgs := config.GetConfigGroup(conf.ConfigDirName)

	StartGameLogicServerWithConfig(cfgs)
}

// StartGameLogicServerWithConfig 傳入Config執行Server
func StartGameLogicServerWithConfig(cfgs *config.Group) {

	gs, _, servercfg, err := SetupGameServer(cfgs)
	if err != nil {
		mlog.Error("初始化服務器錯誤, %v", err)
		return
	}

	// 啟動服務器
	err = gs.Start()
	if err != nil {
		mlog.Error("服務器啟動錯誤,%v", err)
		return
	}

	example.invalid("gameserver  start...")
	example.invalid(servercfg.String())

	mlog.System("log off")
	mlog.GetCLogger().Enable(false)
	select {}
}

// SetupGameServer 建立GameServer
func SetupGameServer(cfgs *config.Group) (*gameserver.GameServer, *gamelogic.LogicManager, *gamelogic.GameServerConfig, error) {
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)

	{ // 註冊當前營運平台
		platform := serverInfo.GetConfigValue("gameserver", "platform")
		if len(platform) > 0 {
			conf.ServicePlatform = conf.Platform(platform)
		}
	}
	{ // 金邊專用
		{ // 舊的聯運ID相容
			lianyunid := cfg.GetConfigInt("gameserver", "lianyunid")
			if lianyunid > 0 {
				conf.LianyunID = lianyunid
			}
		}
		{ // 新的聯運ID位置
			lianyunid := serverInfo.GetConfigInt("gameserver", "lianyunid")
			if lianyunid > 0 {
				conf.LianyunID = lianyunid
			}
		}
		{ // RealAddr
			realaddr := serverInfo.GetConfigValue("server", "realaddr")
			if len(realaddr) > 0 {
				conf.RealAddr = realaddr
			}
		}
	}
	{ // 註冊封包攔截器
		switch conf.ServicePlatform {
		case conf.PlatformPhnum:
			interceptor.RegisterMSGInterceptor(phnum.MSGInterceptor)
		}
	}

	// 獲取配置信息
	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)

	// 初始化 mlog
	mlog.SetServerString(fmt.Sprintf("%s:%d_%d", servercfg.GameName, servercfg.GameID, servercfg.ServerID))

	// gameserver偵聽地址
	gscfg := network.NewServerCfg(servercfg.ServerAddr, conf.TCPSERVERCONNROOMCLIENT, conf.TCPSERVERNAMEROOMCLIENT)

	servermode := serverInfo.GetConfigInt("defense", "servermode")
	defenseCfg := serverdefense.ParseDefenseConfig(serverInfo)
	gs := gameserver.NewGameServer([]*network.ServerCfg{gscfg}, servercfg.DBServerAddr, servercfg.NotifyServerAddr, servercfg.HTTPServerAddr, int(servercfg.ServerID), conf.ServerTypeGameServer, servermode, int(servercfg.GameID), defenseCfg, servercfg.HostAndBackup, servercfg.SuperToolCgf, servercfg.Flag)

	err := gs.Init()

	// 初始化mylog queue
	// mylog.InitDatabaseQueue(serverInfo.GetConfigValue("myconn", "connstring"), serverInfo.GetConfigInt("myconn", "pool"), serverInfo.GetConfigInt("myconn", "perwritecount"), int(servercfg.ServerID), mlog.NewDefaultLogger())
	mylog.InitMongoDatabaseQueue(serverInfo.GetConfigValue("mongoconn", "connstring"),
		serverInfo.GetConfigInt("mongoconn", "pool"),
		serverInfo.GetConfigInt("mongoconn", "perwritecount"),
		int(servercfg.ServerID), mlog.NewDefaultLogger(),
		serverInfo.GetConfigValue("mongoconn", "database"),
		serverInfo.GetConfigValue("mongoconn", "collection"))

	if err != nil {
		return nil, nil, nil, err
	}

	// 創建房間對象
	room := createRoom(servercfg, gs.GameServerRunCtx.PlayerManager)
	// 超端服務
	room.SuperToolService = SuperToolCreator(room, gs.GameServerRunCtx.TimerFuncManager, gs.GameServerRunCtx.DBMsgQueue)

	// 創建房間邏輯管理器
	roomligics := gamelogic.NewGameLogicManager(room, cfgs, gs.GameServerRunCtx.TimerFuncManager, GameLogicCreator)

	// 房間邏輯模塊
	roomlogicmodule := roomlogic.NewRoomLogicModule(gs, room, roomligics, RobotLogicCreator, cfgs)
	gs.GameServerRunCtx.ModuleServer.ModuleGroup.Register(roomlogicmodule)

	// hallserver 註冊
	hallregHandler := serverhandler.NewHallServerRegHandler(servercfg.HallServerAddr, room, servercfg.RegHallServerInterval, roomlogicmodule.GetRobotManager())

	gs.ServerHandler.AddHandler(hallregHandler)
	roomlogicmodule.SetHallregHandler(hallregHandler)
	// 停服
	gs.ServerHandler.AddHandler(serverhandler.NewServerStopHandler(roomlogicmodule.ServerStatus, 60, room))

	// 初始化RPC IP黑名單
	gs.InitRPCBlackIPList(roomlogicmodule.RPCMsgHandler)

	// 超端服務啟動
	room.SuperToolService.Start()

	if app.Context != nil {
		if err := app.Context.Provide(func() *gameserver.GameServer {
			return gs
		}); err != nil {
			mlog.Error("*gameserver.GameServer Register dig error, %v", err)
		}
		if err := app.Context.Provide(func() *gamelogic.Room {
			return room
		}); err != nil {
			mlog.Error("*gamelogic.Room Register dig error, %v", err)
		}
		if err := app.Context.Provide(func() *gamelogic.LogicManager {
			return roomligics
		}); err != nil {
			mlog.Error("*gamelogic.GameLogicManager Register dig error, %v", err)
		}
	}

	return gs, roomligics, servercfg, nil
}

// 創建房間對象
func createRoom(cfg *gamelogic.GameServerConfig, playerManager *player.Manager) *gamelogic.Room {

	rumanager := gamelogic.NewRoomUserManager(playerManager)
	room := gamelogic.NewRoom(cfg, rumanager)

	return room
}
