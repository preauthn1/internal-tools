package main

import (
	"example.invalid/ServerTeam/Servers/internal/notifyserver/notifylogic"
	"example.invalid/ServerTeam/Servers/internal/notifyserver/notifylogic/svrmgm"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

func main() {
	svr := msgserver.NewBgMsgServer(1)
	ms := module.NewServer(svr)

	// 獲取配置
	hc := ms.Config.GetConfig(notifylogic.ConfigFileNameNS)

	addr := hc.GetConfigValue("tcp", "addr")

	// 初始化 mlog
	mlog.SetServerString("notifyserver")

	cfg := network.NewServerCfg(addr, conf.TCPSERVERCONNNOTIFYSERVER, conf.TCPSERVERNAMENOTIFYSERVER)
	svr.AddTCPServer(cfg)

	// 服務器信息
	svrlist := svrmgm.NewServerList()
	// 邏輯模塊
	ms.ModuleGroup.Register(notifylogic.NewNSLogic(svrlist))

	// 啟動服務器
	ms.Run()

	example.invalid("notify server start...")
	mlog.System("log off")
	mlog.GetCLogger().Enable(false)
	select {}
}
