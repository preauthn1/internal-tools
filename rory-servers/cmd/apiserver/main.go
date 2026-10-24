package main

import (
	"net"
	"os"

	api "example.invalid/ServerTeam/Servers/internal/api-service"
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic"
	cliapisvr "example.invalid/ServerTeam/Servers/internal/cli-api-svc"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

func main() {

	conf.ClearLog = true
	conf.CLogLevel = mlog.WarnLevel

	svr := msgserver.NewBgMsgServer(conf.DefaultEventThrottle)
	ms := module.NewServer(svr)

	// 獲取配置
	hc := ms.Config.GetConfig("api")
	addr := hc.GetConfigValue("tcp", "addr")
	apiaddr := hc.GetConfigValue("tcp", "apiaddr")
	grpcaddr := hc.GetConfigValue("tcp", "grpcaddr")

	// 初始化 mlog
	mlog.SetServerString("apiserver")

	cfg := network.NewServerCfg(addr, conf.TCPSERVERCONNDBSERVER, conf.TCPSERVERNAMEDBSERVER)
	svr.AddTCPServer(cfg)

	// 數據庫並發數
	semaphore := hc.GetConfigInt("conn", "semaphore")

	// 邏輯
	ms.ModuleGroup.Register(apilogic.NewAPILogicModule(semaphore))

	// 啟動服務器
	ms.Run()

	if len(apiaddr) > 0 {
		lis, err := net.Listen("tcp", apiaddr)
		if err != nil {
			mlog.Error("初始化 API GRPC 服務錯誤, %v", err)
			os.Exit(1)
		}

		srv := api.NewGRPC()

		mlog.System("啟動 API GRPC 服務, %s", apiaddr)
		go srv.Serve(lis)
	}

	if len(grpcaddr) > 0 {
		lis, err := net.Listen("tcp", grpcaddr)
		if err != nil {
			mlog.Error("初始化 GRPC 控制服務錯誤, %v", err)
			os.Exit(1)
		}

		grpcServer := cliapisvr.NewGRPC()
		mlog.System("啟動 GRPC 控制服務, %s", grpcaddr)
		go grpcServer.Serve(lis)
	}

	mlog.System("API Server Start ...")
	// mlog.System("Log OFF")
	// mlog.GetCLogger().Enable(false)
	select {}
}
