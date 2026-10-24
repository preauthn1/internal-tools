package main

import (
	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic"
	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"
	hshandler "example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/serverhandler"

	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverdefense"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func initServer() (server *gameserver.GameServer, err error) {

	var cfgHelper = hallutil.GetHallServerConfig(config.GetConfigGroup(conf.ConfigDirName).GetConfig("hall"))

	//服务器配置
	var listenerAddr = []*network.ServerCfg{
		network.NewServerCfg(cfgHelper.HallAddr, conf.TCPSERVERCONNHALLCLIENT, conf.TCPSERVERNAMEHALLCLIENT),          //大厅侦听地址
		network.NewServerCfg(cfgHelper.GSMgmAddr, conf.TCPSERVERCONNHALLGAMESERVER, conf.TCPSERVERNAMEHALLGAMESERVER), //gameserver侦听地址
	}
	var dbAddr = cfgHelper.DBServerAddr
	var notifyAddr = cfgHelper.NotifyServerAddr
	var httpAddr = ""
	var serverID = int(cfgHelper.ServerID)
	var serverType = conf.ServerTypeHallServer
	var serverMode = config.GetConfigGroup(conf.ConfigDirName).GetConfig("hall").GetConfigInt("defense", "servermode")
	var gameID = -1
	var defenseCfg = serverdefense.ParseDefenseConfig(config.GetConfigGroup(conf.ConfigDirName).GetConfig("hall"))
	var hostAndBackup = int32(0)
	var superToolCfg = gamelogic.SuperToolConfig{}
	var flagID = ""

	//初始化 mlog
	mlog.SetServerString(fmt.Sprintf("%s:%d", "hallserver", serverID))

	//服务器建立
	server = gameserver.NewGameServer(listenerAddr, dbAddr, notifyAddr, httpAddr, serverID, serverType, serverMode, gameID, defenseCfg, hostAndBackup, superToolCfg, flagID)
	err = server.Init()

	//模組配置(必须在建立服务器后)
	setHandler(server, cfgHelper)

	return server, err
}

func setHandler(server *gameserver.GameServer, cfgHelper *hallutil.HallServerConfig) {
	//大厅逻辑
	hallmodule := halllogic.NewHallLogicModule(server, cfgHelper)
	server.GameServerRunCtx.ModuleServer.ModuleGroup.Register(hallmodule)

	//初始化RPC IP黑名单
	server.InitRPCBlackIPList(hallmodule.RPCMsgHandler)

	//停服广播服务器处理器
	server.ServerHandler.AddHandler(hshandler.NewServerStopHandler(hallmodule.ServerStatus, 30))
}

func main() {
	defer util.PanicErrStack()

	var gs *gameserver.GameServer
	var err error

	//服务器配置
	gs, err = initServer()
	if err != nil {
		mlog.Error("初始化服务器错误,%v", err)
		return
	}

	//启动服务器
	err = gs.Start()
	if err != nil {
		mlog.Error("服务器启动错误,%v", err)
		return
	}

	example.invalid("hallserver  start...")

	//測試 elk 收集日志用
	var data = db.ProcError{
		Code:    500,
		Message: "test message",
		PS:      "test ps",
	}
	mlog.NewDefaultLogger().SetData(data).Info("hallserver start")
	mlog.System("log off")
	mlog.GetCLogger().Enable(false)
	select {}
}
