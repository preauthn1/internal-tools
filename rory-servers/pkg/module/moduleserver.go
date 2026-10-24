package module

import (
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
)

// 模塊日誌記錄器，組合默認的日誌記錄實現
type moduleLog struct{}

func (ml *moduleLog) WriteLog(lvl int, format string, a ...interface{}) {
	mlog.NewDefaultLogger().SetCallDepLen(5).WriteLog(lvl, format, a...)
	//mlog.WriteLog(lvl, format, a...)
}

func (ml *moduleLog) Close() {}

// Server 模塊服務器
type Server struct {
	ModuleGroup *Group
	Logger      mlog.Logger
	Config      *config.Group
	Server      *msgserver.MsgServer
	moduleCtx   *MouduleCtx
}

// NewServer 創建模塊服務器
func NewServer(msgserver *msgserver.MsgServer) *Server {
	mserver := new(Server)

	ctx := new(MouduleCtx)
	// Initial configuration
	ctx.Config = config.GetConfigGroup(conf.ConfigDirName)
	// 日誌
	ctx.Logger = new(moduleLog)
	// Event dispatch
	ctx.EvtDispatcher = msgserver.EventChannel
	// 添加事件分派器命令
	// 消息服務器
	ctx.Server = msgserver

	// 由於command引用了event模塊.為避免循環引用，在此加入事件命令
	command.AddCommand(event.NewEventCmd(ctx.EvtDispatcher))

	mserver.moduleCtx = ctx
	mserver.Server = msgserver
	mserver.Config = ctx.Config
	mserver.Logger = ctx.Logger

	mserver.ModuleGroup = NewGroup()

	return mserver
}

// Run 啟動服務器
func (ms *Server) Run() error {
	// 初始化消息服務器
	ms.Server.InitNetServer()

	// 啟動消息服務器
	var err error
	err = ms.Server.Start()
	if err != nil {
		return err
	}

	var grpcaddr string

	if val := ms.Config.GetConfigValue("serverinfo", "server", "grpcaddr"); len(val) > 0 {
		grpcaddr = val
	} else if val := ms.Config.GetConfigValue("api", "tcp", "grpcaddr"); len(val) > 0 {
		grpcaddr = val
	} else if val := ms.Config.GetConfigValue("db", "tcp", "grpcaddr"); len(val) > 0 {
		grpcaddr = val
	}

	// 命令
	if len(grpcaddr) == 0 {
		command.Run()
	}

	for _, m := range ms.ModuleGroup.modules {
		err = m.Run(ms.moduleCtx)

		if err != nil {
			mlog.Error("啟動模塊 %s 錯誤 %v", m.Name(), err)
		}
	}

	return nil
}

// Stop 停止服務器
func (ms *Server) Stop() {
	ms.ModuleGroup.Range(func(m Module) {
		m.Stop()
	})

	ms.Server.Stop()
	command.Stop()
}
