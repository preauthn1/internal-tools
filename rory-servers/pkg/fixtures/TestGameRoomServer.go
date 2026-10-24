package fixtures

import (
	"os"
	"path"

	"example.invalid/ServerTeam/Servers/internal/api-service"
	apipb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gamequeue"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverdefense"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/robot"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
)

// TestGameRoomServer ?
type TestGameRoomServer struct {
	*logicmodule.GameServerLogicModule
	GameServer           *gameserver.GameServer
	GameLogicManager     *gamelogic.LogicManager
	Room                 *gamelogic.Room
	RoomActionController *roomlogic.RoomActionController
	RoomActionManager    *roomlogic.RoomActionManager
	RobotService         *robot.Service
	sstorage             scorestorage.UserScoreStorage
}

// NewTestGameRoomServer ?
func NewTestGameRoomServer() (tgrs *TestGameRoomServer) {
	tgrs = new(TestGameRoomServer)
	tgrs.GameServerLogicModule = logicmodule.NewGameServerLogicModule()
	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)

	// 獲取配置信息
	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
	servercfg.DealLogType = gamelogic.DealLogTypeConsole
	// gameserver偵聽地址
	gscfg := network.NewServerCfg(servercfg.ServerAddr, conf.TCPSERVERCONNROOMCLIENT, conf.TCPSERVERNAMEROOMCLIENT)

	servermode := cfg.GetConfigInt("defense", "servermode")
	defenseCfg := serverdefense.ParseDefenseConfig(cfg)
	tgrs.GameServer = NewGameServerByTest([]*network.ServerCfg{gscfg}, servercfg.DBServerAddr, servercfg.NotifyServerAddr, servercfg.HTTPServerAddr, int(servercfg.ServerID), conf.ServerTypeGameServer, servermode, int(servercfg.GameID), defenseCfg, servercfg.Flag)
	tgrs.PlayerManager = tgrs.GameServer.GameServerRunCtx.PlayerManager
	// createRoom 創建房間對象
	rumanager := gamelogic.NewRoomUserManager(tgrs.GameServer.GameServerRunCtx.PlayerManager)
	tgrs.Room = gamelogic.NewRoom(servercfg, rumanager)

	// 超端服務
	tgrs.Room.SuperToolService = roomserver.SuperToolCreator(tgrs.Room, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, tgrs.GameServer.GameServerRunCtx.DBMsgQueue)
	tgrs.Room.SuperToolService.Start()

	// 創建房間邏輯管理器
	tgrs.GameLogicManager = gamelogic.NewGameLogicManager(tgrs.Room, cfgs, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, roomserver.GameLogicCreator)
	//console 寫分
	tgrs.sstorage = new(scorestorage.ConsoleUserScoreStorage)

	// 房間邏輯模塊
	tgrs.RoomActionManager = roomlogic.NewRoomActionManager(tgrs.Room.RoomUserManager, tgrs.Room, tgrs.GameLogicManager, nil, tgrs.GameServer, tgrs.sstorage, nil, servercfg, tgrs.ServerStatus)

	// 構造遊戲撮合器
	tgrs.RoomActionController = roomlogic.NewRoomActionController(tgrs.Room.RoomUserManager, tgrs.RoomActionManager, tgrs.Room, cfgs)

	tgrs.RobotService = robot.NewRobotService(tgrs.Room, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, tgrs.sstorage, tgrs.GameServer.GameServerRunCtx.DBMsgQueue, tgrs.GameLogicManager, cfgs, tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, roomserver.RobotLogicCreator, tgrs.RoomActionController, tgrs.ServerStatus)
	//tgrs.robotService = robot.NewRobotService(tgrs.Room, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, tgrs.sstorage, tgrs.GameServer.GameServerRunCtx.DBMsgQueue, tgrs.GameLogicManager, cfgs, tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, RobotLogicCreator, tgrs.RoomActionController, tgrs.ServerStatus)
	tgrs.RoomActionManager.SetRobotService(tgrs.RobotService)

	tgrs.RobotService.StartRobotService()

	//tgrs.robotService = robot.NewRobotService(tgrs.Room, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, tgrs.sstorage, tgrs.GameServer.GameServerRunCtx.DBMsgQueue, tgrs.GameLogicManager, cfgs, tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, RobotLogicCreator, tgrs.RoomActionController, tgrs.ServerStatus)
	//tgrs.RoomActionManager.SetRobotService(tgrs.RobotService)

	queuetype := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer).GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "queuetype")
	// 定時撮合
	if queuetype == roomlogic.QueueTypeTimeInterval {
		gmq := gamequeue.NewIntervalMatchGameQueue(tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, tgrs.RoomActionController, cfgs)
		tgrs.RoomActionManager.SetGameQueue(gmq)
	} else if queuetype == roomlogic.QueueTypeQueueUserCount { //固定人数撮合
		gmq := gamequeue.NewUserCountGameQueue(tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, tgrs.RoomActionController, cfgs)
		tgrs.RoomActionManager.SetGameQueue(gmq)
	} else if queuetype == roomlogic.QueueTypeWinLoseUser {
		gmq := gamequeue.NewWinLoseUserGameQueue(tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, tgrs.RoomActionController, cfgs)
		tgrs.RoomActionManager.SetGameQueue(gmq)
	} else if queuetype == roomlogic.QueueTypeUserChoose {
		gmq := gamequeue.NewUserChooseGameQueue()
		tgrs.RoomActionManager.SetGameQueue(gmq)
	}
	// 設置桌子對應的邏輯實例
	for _, t := range tgrs.Room.Tables {
		t.GameLogic = tgrs.GameLogicManager.GetGameLogic(t.TableNo)
	}
	//_ = roomlogic.NewRoomLogicModule(gs, tgrs.Room, tgrs.GameLogicManager, roomserver.RobotLogicCreator, cfgs)
	return
}

// NewTestGameRoomServerWithConfig ?
func NewTestGameRoomServerWithConfig(cfgs *config.Group) (tgrs *TestGameRoomServer) {
	tgrs = new(TestGameRoomServer)
	tgrs.GameServerLogicModule = logicmodule.NewGameServerLogicModule()
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)

	// 獲取配置信息
	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
	servercfg.DealLogType = gamelogic.DealLogTypeConsole
	// gameserver偵聽地址
	gscfg := network.NewServerCfg(servercfg.ServerAddr, conf.TCPSERVERCONNROOMCLIENT, conf.TCPSERVERNAMEROOMCLIENT)

	servermode := cfg.GetConfigInt("defense", "servermode")
	defenseCfg := serverdefense.ParseDefenseConfig(cfg)
	tgrs.GameServer = NewGameServerByTest([]*network.ServerCfg{gscfg}, servercfg.DBServerAddr, servercfg.NotifyServerAddr, servercfg.HTTPServerAddr, int(servercfg.ServerID), conf.ServerTypeGameServer, servermode, int(servercfg.GameID), defenseCfg, servercfg.Flag)
	tgrs.PlayerManager = tgrs.GameServer.GameServerRunCtx.PlayerManager
	// createRoom 創建房間對象
	rumanager := gamelogic.NewRoomUserManager(tgrs.GameServer.GameServerRunCtx.PlayerManager)
	tgrs.Room = gamelogic.NewRoom(servercfg, rumanager)

	// 超端服務
	tgrs.Room.SuperToolService = roomserver.SuperToolCreator(tgrs.Room, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, tgrs.GameServer.GameServerRunCtx.DBMsgQueue)
	tgrs.Room.SuperToolService.Start()

	// 創建房間邏輯管理器
	tgrs.GameLogicManager = gamelogic.NewGameLogicManager(tgrs.Room, cfgs, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, roomserver.GameLogicCreator)
	// console 寫分
	tgrs.sstorage = new(scorestorage.ConsoleUserScoreStorage)

	// 房間邏輯模塊
	tgrs.RoomActionManager = roomlogic.NewRoomActionManager(tgrs.Room.RoomUserManager, tgrs.Room, tgrs.GameLogicManager, nil, tgrs.GameServer, tgrs.sstorage, nil, servercfg, tgrs.ServerStatus)

	// 構造遊戲撮合器
	tgrs.RoomActionController = roomlogic.NewRoomActionController(tgrs.Room.RoomUserManager, tgrs.RoomActionManager, tgrs.Room, cfgs)

	tgrs.RobotService = robot.NewRobotService(tgrs.Room, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, tgrs.sstorage, tgrs.GameServer.GameServerRunCtx.DBMsgQueue, tgrs.GameLogicManager, cfgs, tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, roomserver.RobotLogicCreator, tgrs.RoomActionController, tgrs.ServerStatus)
	//tgrs.robotService = robot.NewRobotService(tgrs.Room, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, tgrs.sstorage, tgrs.GameServer.GameServerRunCtx.DBMsgQueue, tgrs.GameLogicManager, cfgs, tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, RobotLogicCreator, tgrs.RoomActionController, tgrs.ServerStatus)
	tgrs.RoomActionManager.SetRobotService(tgrs.RobotService)

	tgrs.RobotService.StartRobotService()

	//tgrs.robotService = robot.NewRobotService(tgrs.Room, tgrs.GameServer.GameServerRunCtx.TimerFuncManager, tgrs.sstorage, tgrs.GameServer.GameServerRunCtx.DBMsgQueue, tgrs.GameLogicManager, cfgs, tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, RobotLogicCreator, tgrs.RoomActionController, tgrs.ServerStatus)
	//tgrs.RoomActionManager.SetRobotService(tgrs.RobotService)

	queuetype := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer).GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "queuetype")
	// 定時撮合
	if queuetype == roomlogic.QueueTypeTimeInterval {
		gmq := gamequeue.NewIntervalMatchGameQueue(tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, tgrs.RoomActionController, cfgs)
		tgrs.RoomActionManager.SetGameQueue(gmq)
	} else if queuetype == roomlogic.QueueTypeQueueUserCount { // 固定人数撮合
		gmq := gamequeue.NewUserCountGameQueue(tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, tgrs.RoomActionController, cfgs)
		tgrs.RoomActionManager.SetGameQueue(gmq)
	} else if queuetype == roomlogic.QueueTypeWinLoseUser {
		gmq := gamequeue.NewWinLoseUserGameQueue(tgrs.GameServer.GameServerRunCtx.MsgServer.EventChannel, tgrs.RoomActionController, cfgs)
		tgrs.RoomActionManager.SetGameQueue(gmq)
	} else if queuetype == roomlogic.QueueTypeUserChoose {
		gmq := gamequeue.NewUserChooseGameQueue()
		tgrs.RoomActionManager.SetGameQueue(gmq)
	}
	// 設置桌子對應的邏輯實例
	for _, t := range tgrs.Room.Tables {
		t.GameLogic = tgrs.GameLogicManager.GetGameLogic(t.TableNo)
	}
	//_ = roomlogic.NewRoomLogicModule(gs, tgrs.Room, tgrs.GameLogicManager, roomserver.RobotLogicCreator, cfgs)

	app.Context.Provide(func() *player.Manager {
		return tgrs.PlayerManager
	})
	app.Context.Provide(func() apipb.GameBetServiceClient {
		return new(api.MockGameBetServiceServer)
	})
	return
}
