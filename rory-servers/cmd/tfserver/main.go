package main

import (
	"net"
	"os"

	"example.invalid/ServerTeam/Servers/internal/tfserver/router"
	"example.invalid/ServerTeam/Servers/internal/tfserver/router/portrouter"
	"example.invalid/ServerTeam/Servers/internal/tfserver/tfconfig"
	"example.invalid/ServerTeam/Servers/internal/tfserver/tfnet"

	cliapisvr "example.invalid/ServerTeam/Servers/internal/cli-api-svc"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

var tfProvider tfnet.TFListenPortProvider

func init() {
	conf.CLogLevel = mlog.InfoLevel
	conf.ClearLog = true
}

func main() {

	// 初始化路由表
	initRouteTable()

	// 初始化转发通道
	initTFChannel()

	// 命令
	cfgs := config.GetConfigGroup("config")
	grpcaddr := cfgs.GetConfigValue("tfserver", "centerserver", "grpcaddr")

	if len(grpcaddr) > 0 {
		lis, err := net.Listen("tcp", grpcaddr)
		if err != nil {
			mlog.Error("初始化 GRPC 控制服務錯誤, %v", err)
			os.Exit(1)
		}

		grpcServer := cliapisvr.NewGRPC()
		mlog.System("啟動 GRPC 控制服務, %s", grpcaddr)
		go grpcServer.Serve(lis)
	} else {
		command.AddCommand(tfconfig.NewCheckKeyCmd())
		command.Run()
	}

	// 啟動服務器
	startTFServer()

	select {}
}

// 初始化路由表
func initRouteTable() {
	// 端口映射加载器
	pml := portrouter.NewAutoRefreshRemoteConfigMapLoader(tfconfig.ConfigCtx.PortMapURL, tfconfig.ConfigCtx.RefreshPortMapInterval)
	// 創建端口映射路由
	prt := portrouter.NewPortRouter(pml)

	router.RouteCtx.AddRouter(prt)

	// 轉發端口提供
	tfProvider = NewRemotePortMapTFPortProvider(pml)

	app.Context.Provide(func() portrouter.PortMapLoader { return pml })
	app.Context.Provide(func() *portrouter.PortRouter { return prt })
}

// 初始化轉發通道
func initTFChannel() {
	// 路由通道
	ptfc := tfnet.NewPortTFChannel()
	tfnet.ChannelCtx.AddChannel(ptfc)

	// 狀態上報通道
	rpc := tfnet.NewReportChannel(tfconfig.ConfigCtx.ReportURL, tfconfig.ConfigCtx.ReportInterval)
	tfnet.ChannelCtx.AddChannel(rpc)
}

// 啟動轉發服務器
func startTFServer() {
	ptckInterval := tfconfig.ConfigCtx.CheckTFPortInterval

	tfSvr := tfnet.NewAutoCheckPortTFServer(tfProvider, ptckInterval)
	tfSvr.Start()
}

// RemotePortMapTFPortProvider 轉發端口包裝類,從轉發端口映射規則中解析需要偵聽的端口
// 實現TFListenPortProvider接口
type RemotePortMapTFPortProvider struct {
	pmLoader portrouter.PortMapLoader
}

// NewRemotePortMapTFPortProvider 實現TFListenPortProvider接口
func NewRemotePortMapTFPortProvider(pmLoader portrouter.PortMapLoader) *RemotePortMapTFPortProvider {
	rpmtf := new(RemotePortMapTFPortProvider)
	rpmtf.pmLoader = pmLoader

	return rpmtf
}

// LoadListenPort 偵聽端口
func (rpm *RemotePortMapTFPortProvider) LoadListenPort() []int {
	// 偵聽端口
	var listenPorts []int

	pm := rpm.pmLoader.LoadPortMap()
	if pm != nil && len(pm) > 0 {
		for pt := range pm {
			listenPorts = append(listenPorts, pt)
		}
	}

	return listenPorts
}
