package main

import (
	"example.invalid/ServerTeam/Servers/internal/phnum/dbserver"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

func main() {
	svr := msgserver.NewBgMsgServer(conf.DefaultEventThrottle)
	ms := module.NewServer(svr)

	// 獲取配置
	hc := ms.Config.GetConfig("db")

	addr := hc.GetConfigValue("tcp", "addr")

	// 初始化 mlog
	mlog.SetServerString("dbserver")

	cfg := network.NewServerCfg(addr, conf.TCPSERVERCONNDBSERVER, conf.TCPSERVERNAMEDBSERVER)
	svr.AddTCPServer(cfg)

	// 數據庫並發數
	// _ = hc.GetConfigInt("conn", "semaphore")
	semaphore := hc.GetConfigInt("conn", "semaphore")

	// 邏輯
	ms.ModuleGroup.Register(dbserver.NewDBLogicModule(semaphore))

	// 啟動服務器
	ms.Run()

	// command.AddCommand(new(msghandler.FastLogonStatCommand))

	mlog.System("dbserver start...")
	mlog.System("log off")
	mlog.GetCLogger().Enable(false)
	select {}
}
