package roomlogic

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gamequeue"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/robot"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomutil"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
	"example.invalid/ServerTeam/Servers/pkg/serverhandler"
)

// 遊戲選桌類型
const (
	QueueTypeTimeInterval   = iota // 撮合方式,定時撮合
	QueueTypeQueueUserCount        // 撮合方式,排隊人數
	QueueTypeWinLoseUser           // 輸贏對戰,真人排隊
	QueueTypeUserChoose            // 玩家自選
	QueueTypeUserReady             // 玩家確定後排隊
)

// AB服
const (
	HostTypeMaster        = 0 // 正式
	HostTypeSlaveToMaster = 1 // 備援轉正式
	HostTypeSlave         = 2 // 備援
)

// RoomLogic 大廳邏輯
type RoomLogic struct {
	*logicmodule.GameServerLogicModule
	room                 *gamelogic.Room
	roomuserManager      *gamelogic.RoomUserManager
	serverCfg            *gamelogic.GameServerConfig
	configs              *config.Group
	verProvider          GameVerProvider
	gameServer           *gameserver.GameServer
	gameLogics           *gamelogic.LogicManager
	robotService         *robot.Service
	sstorage             scorestorage.UserScoreStorage
	roomactionController *RoomActionController
	roomactionmanager    *RoomActionManager
	regHallHandler       *serverhandler.HallServerRegHandler
	//vipUserService       *gamelogic.VipUserService
}

// NewRoomLogicModule 創建房間邏輯對象
func NewRoomLogicModule(gs *gameserver.GameServer, room *gamelogic.Room, logics *gamelogic.LogicManager, robotLogicCreator robotlogic.LogicFunc, cfgs *config.Group) *RoomLogic {
	rl := new(RoomLogic)
	rl.GameServerLogicModule = logicmodule.NewGameServerLogicModule()

	rl.gameServer = gs
	rl.room = room
	rl.configs = cfgs
	rl.roomuserManager = room.RoomUserManager
	rl.serverCfg = gamelogic.GetGameServerConfig(cfgs.Configs[gamelogic.ConfigFileNameServerInfo], cfgs.Configs[gamelogic.ConfigFileNameGameServer])
	rl.gameLogics = logics

	dbMsgQueue := gs.GameServerRunCtx.DBMsgQueue
	if rl.serverCfg.ScoreMode == gamelogic.ScoreUpdateModeLogout {
		rl.sstorage = scorestorage.NewLogoutUserScoreStorage(room, dbMsgQueue)
	} else if rl.serverCfg.ScoreMode == gamelogic.ScoreUpdateModePiSai {
		rl.sstorage = scorestorage.NewPiSaiScoreStorage(room, dbMsgQueue)
	} else if rl.serverCfg.ScoreMode == gamelogic.ScoreUpdateModeConditional {
		rl.sstorage = scorestorage.NewConditionalUserScoreStorage(room, dbMsgQueue, rl.serverCfg)
	} else {
		rl.sstorage = scorestorage.NewNowUserScoreStorage(room, dbMsgQueue)
	}

	rl.roomactionmanager = NewRoomActionManager(rl.roomuserManager, rl.room, logics, nil, gs, rl.sstorage, nil, rl.serverCfg, rl.ServerStatus)

	// 構造遊戲撮合器
	rac := NewRoomActionController(rl.roomuserManager, rl.roomactionmanager, rl.room, cfgs)
	rl.roomactionController = rac
	app.Context.Provide(func() *RoomActionController { return rl.roomactionController })
	app.Context.Provide(func() gamelogic.ActionController { return rl.roomactionController })

	// 初始化機器人服務
	rl.initRobotService(robotLogicCreator)

	// 設置機器人邏輯
	rl.roomactionmanager.robotService = rl.robotService

	queuetype := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer).GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "queuetype")

	{ // 撮合
		var gmq gamequeue.GameQueue
		if queuetype == QueueTypeTimeInterval {
			gmq = gamequeue.NewIntervalMatchGameQueue(gs.GameServerRunCtx.MsgServer.EventChannel, rac, cfgs)
			rl.roomactionmanager.SetGameQueue(gmq)
		} else if queuetype == QueueTypeQueueUserCount { // 固定人數撮合
			gmq = gamequeue.NewUserCountGameQueue(gs.GameServerRunCtx.MsgServer.EventChannel, rac, cfgs)
			rl.roomactionmanager.SetGameQueue(gmq)
		} else if queuetype == QueueTypeWinLoseUser {
			gmq = gamequeue.NewWinLoseUserGameQueue(gs.GameServerRunCtx.MsgServer.EventChannel, rac, cfgs)
			rl.roomactionmanager.SetGameQueue(gmq)
		} else if queuetype == QueueTypeUserChoose {
			gmq = gamequeue.NewUserChooseGameQueue()
			rl.roomactionmanager.SetGameQueue(gmq)
		} else if queuetype == QueueTypeUserReady {
			gmq = gamequeue.NewUserReadyQueue(rac, cfgs)
			rl.roomactionmanager.SetGameQueue(gmq)
		}
		app.Context.Provide(func() gamequeue.GameQueue { return gmq })
	}

	// 設置監控及黑白名單服務
	//rl.vipUserService = gamelogic.NewVipUserService(room, gs.GameServerRunCtx.TimerFuncManager, gs.GameServerRunCtx.DBMsgQueue)
	//rl.vipUserService.Start()
	// 設置詐金花監控名單服務
	//gameid := cfgs.GetConfig(gamelogic.CONFIGFILENAME_GAMESERVER).GetConfigInt("gameserver", "gameid")
	//if gameid == 3 {
	//	rl.zjhVipUserService = gamelogic.NewZjhVipUserService(room, gs.GameServerRunCtx.TimerFuncManager, gs.GameServerRunCtx.DBMsgQueue)
	//	rl.zjhVipUserService.Start()
	//	room.SetZjhVipUserService(rl.zjhVipUserService)
	//	//Zjh超端
	//	rl.RPCMsgHandler.AddHandler(gamelogic.NewRefreshZjhVipUserHandler(rl.zjhVipUserService.GetVipUserProvider()))
	//}

	//gameid := cfgs.GetConfig(gamelogic.CONFIGFILENAME_GAMESERVER).GetConfigInt("gameserver", "gameid")
	// 設置桌子對應的邏輯實例
	// for _, t := range room.Tables {
	// 	t.GameLogic = logics.GetGameLogic(t.TableNo)
	//t.GameLogic.SetVipUserProvider(rl.vipUserService.GetVipUserProvider())
	//t.GameLogic.SetVipUserService(rl.vipUserService)
	//t.GameLogic.SetZjhVipUserService(rl.zjhVipUserService)
	// }

	if app.Context != nil {
		app.Context.Provide(func() *RoomActionManager { return rl.roomactionmanager })
		app.Context.Provide(func() *robot.Service { return rl.robotService })
		app.Context.Provide(func() scorestorage.UserScoreStorage { return rl.sstorage })
		app.Context.Provide(func() logicmodule.ServerStatusProvider { return rl.ServerStatus })
	}

	return rl
}

// SetRegHallHandler 將 regHallHandler 配置給 遊戲注冊器使用
func (rl *RoomLogic) SetRegHallHandler(regHallHandler *serverhandler.HallServerRegHandler) {
	rl.regHallHandler = regHallHandler
}

// Name 模組名稱
func (rl *RoomLogic) Name() string {
	return "roomlogic"
}

// Run 運行模組
func (rl *RoomLogic) Run(ctx *module.MouduleCtx) error {

	rl.Init(ctx, rl.gameServer.GameServerRunCtx.PlayerManager, rl.gameServer.GameServerRunCtx.DBMsgQueue, rl.serverCfg.GetHTTPServerCfg())

	// 啟動機器人
	rl.robotService.StartRobotService()

	rl.registerEvtHandler(ctx)

	rl.registerRCPMsgHandler(ctx)

	rl.registerMsgHandler(ctx)

	rl.registerNotifyMsgHandler(ctx)

	rl.reigsterHTTPMsgHandler(ctx)

	RoomLogicCmdP = NewRoomLogicCmd(rl)
	RoomStateCmdP = NewRoomStateCmd(rl)
	command.AddCommand(RoomLogicCmdP)
	command.AddCommand(RoomStateCmdP)

	return nil
}

// Stop 停止模組
func (rl *RoomLogic) Stop() {}

// SetHallregHandler 設定重新註冊大廳處理器
func (rl *RoomLogic) SetHallregHandler(hallregHandler *serverhandler.HallServerRegHandler) {
	rl.regHallHandler = hallregHandler
}

func (rl *RoomLogic) registerEvtHandler(ctx *module.MouduleCtx) {
	ctx.OnConnect(func(clt *network.ConnInfo) bool {
		mlog.Debug("client %s connnect...", clt)
		// 通知用戶管理器用戶連接
		rl.PlayerManager.SetPlayer(player.NewGuest(clt))
		return true
	})
	ctx.OnDisConnect(func(clt *network.ConnInfo) bool {
		mlog.Debug("client %s disconnnect...", clt)
		rl.roomactionmanager.OnUserLeave(clt.ConnID())
		return true
	})
	ctx.OnMessage(func(clt *network.ConnInfo, m *msg.Message) bool {
		// 檢查服務器狀態
		if clt.IsServerConn(conf.TCPSERVERNAMEROOMCLIENT) {
			if (rl.ServerStatus.GetServerStatus() == conf.ServerStatusClose || rl.ServerStatus.GetServerStatus() == conf.ServerStatusRecover) && !rl.ServerStatus.IsWhiteIP(clt.IPAddress) {
				roomutil.KickUser(rl.ServerStatus.GetStopServerTipMsg(), clt.ConnID(), roomutil.MsgTypeIDRoomKick, ctx.Server)
				return false
			}
		}
		return true
	})
}

func (rl *RoomLogic) registerMsgHandler(ctx *module.MouduleCtx) {
	// 房間邏輯處理
	rl.MsgHandler.AddHandler(rl.roomactionmanager)
	// 版本號
	rl.MsgHandler.AddHandler(NewGetGameVerHandler(rl.verProvider))
}

func (rl *RoomLogic) registerRCPMsgHandler(ctx *module.MouduleCtx) {
	rl.RPCMsgHandler.AddHandler(rl.roomactionmanager) // 房間dbserver返回消息

	// 遊戲啟動註冊
	ip, port := rl.serverCfg.GetServerIPPort()
	gsreg := NewGameServerReg(fmt.Sprintf("%s:%d", ip, port), rl.serverCfg.NotifyServerAddr, rl.serverCfg.HTTPServerAddr, rl.room.ServerID, rl.DBMsgQueue, rl.room, rl.ServerStatus)
	// 將 regHallHandler 配置給 遊戲注冊器使用
	gsreg.SetRegHallHandler(rl.regHallHandler)
	rl.RPCMsgHandler.AddHandler(gsreg)
	gsreg.SendDBServerStartMsg()
	rpcverprovider := NewRPCGameVerProvider(gsreg)
	rl.verProvider = rpcverprovider
	// 客戶端版本命令
	cvercmd := NewCVersionCmd(rl.verProvider)
	command.AddCommand(cvercmd)
	CVersionCmdP = cvercmd

	// 機器人登錄
	rl.RPCMsgHandler.AddHandler(rl.robotService.GetRobotLoaderHookHandler())

	// 版本號更新, 需要配置 regHallHandler
	rl.RPCMsgHandler.AddHandler(NewRefreshVerHandler(rpcverprovider, rl.regHallHandler))
	// VIP用戶更新
	// rl.RPCMsgHandler.AddHandler(gamelogic.NewRefreshVipUserHandler(rl.vipUserService.GetVipUserProvider()))
	// 超端啟動註冊
	rl.RPCMsgHandler.AddHandler(NewSuperToolHandler(gsreg.room, rl.DBMsgQueue, rl.room.ServerID, rl.gameServer.SuperToolCgf, rl.gameServer.GameServerRunCtx.TimerFuncManager, rl.room.GameInfo.GameID, rl.room.Flag))

	rl.RPCMsgHandler.AddHandler(rl.room.SuperToolService.GetSuperToolHandler())

	if app.Context != nil {
		app.Context.Provide(func() GameVerProvider { return rl.verProvider })
	}
}

// 註冊通知處理器
func (rl *RoomLogic) registerNotifyMsgHandler(ctx *module.MouduleCtx) {
	rl.NotifyMsgHandler.AddHandler(NewNotifyRefreshVerHandler(rl.DBMsgQueue, rl.room.GameInfo.GameID, rl.room.ServerID))
	rl.NotifyMsgHandler.AddHandler(NewRoomMasterSlaveHandler(rl, rl.DBMsgQueue, rl.room.GameInfo.GameID, rl.room.ServerID))
	rl.NotifyMsgHandler.AddHandler(NewNotifyPayOrderHandler(rl.room.GameInfo.GameID, rl.room.ServerID, rl.roomuserManager, rl.sstorage))
}

// 註冊Http消息處理器
func (rl *RoomLogic) reigsterHTTPMsgHandler(ctx *module.MouduleCtx) {
	// 如果要啟動Http服務
	if rl.serverCfg.HTTPServerAddr != "" {
		// http處理器
		rl.HTTPMsgHandler.AddHandler(rl.roomactionmanager)
	}
}

// 初始化機器人邏輯
func (rl *RoomLogic) initRobotService(robotLogicCreator robotlogic.LogicFunc) {

	// 初始化機器人服務
	rl.robotService = robot.NewRobotService(rl.room, rl.gameServer.GameServerRunCtx.TimerFuncManager, rl.sstorage, rl.gameServer.GameServerRunCtx.DBMsgQueue, rl.gameLogics, rl.configs, rl.gameServer.GameServerRunCtx.MsgServer.EventChannel, robotLogicCreator, rl.roomactionController, rl.ServerStatus)
}

// ReloadGameServerConfig 重新讀取設定檔
func (rl *RoomLogic) ReloadGameServerConfig() {
	cfgs := config.GetConfigGroup(conf.ConfigDirName)

	for i := int32(0); i < rl.room.TableCount; i++ {
		rl.gameLogics.GetGameLogic(i).RefreshConfigGroup(cfgs)
	}
}

// UpdateHostAndBackUp 更新AB服
func (rl *RoomLogic) UpdateHostAndBackUp(newHostType int32, gameServerAddr string) {
	rl.roomactionController.room.HostAndBackup = newHostType
	rl.roomactionmanager.room.HostAndBackup = newHostType
	rl.gameServer.NotifyHandler.HostAndBackup = newHostType
	rl.roomactionmanager.room.NewGameServerAddr = gameServerAddr
}

// GetRobotManager 獲取機器人管理器
func (rl *RoomLogic) GetRobotManager() *robotcomm.RobotManager {
	return rl.robotService.GetRobotManager()
}
