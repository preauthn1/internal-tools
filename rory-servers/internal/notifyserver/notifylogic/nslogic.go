package notifylogic

import (
	"strings"
	"sync"

	"example.invalid/ServerTeam/Servers/internal/notifyserver/notifylogic/svrmgm"
	"example.invalid/ServerTeam/Servers/internal/notifyserver/notifylogic/task"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// 配置文件名
const (
	ConfigFileNameNS = "ns"
)

// 配置節名稱
const (
	ConfigSectionNameConn = "conn"
	ConfigSectionNameNS   = "ns"
)

// NSLogic 通知/協調服務器邏輯
type NSLogic struct {
	ctx          *module.MouduleCtx
	msgHandler   *module.HookMsgHandlerGroup
	serverList   *svrmgm.ServerList
	hallUserList *sync.Map
}

// NewNSLogic 創建通知程序邏輯
func NewNSLogic(severList *svrmgm.ServerList) *NSLogic {
	nslogic := new(NSLogic)
	nslogic.msgHandler = module.NewHookMsgHandlerGroup()
	nslogic.serverList = severList
	nslogic.hallUserList = new(sync.Map)
	// 加入命令行
	// 服務器查詢命令
	command.AddCommand(svrmgm.NewSvrMgmCmd(severList))

	return nslogic
}

// Run ?
func (nslgc *NSLogic) Run(ctx *module.MouduleCtx) error {
	// 設置數據庫連接
	task.ConnString = ctx.GetConfigData(ConfigFileNameNS, ConfigSectionNameConn, "connstring")
	task.ConnString = strings.Replace(task.ConnString, "\"", "", -1)
	// 初始化數據庫
	task.InitDatabase()

	// 註冊協議解析器
	nslgc.registerMsgParser(ctx)

	// 註冊事件處理器
	nslgc.registerEvtHandler(ctx)

	// 註冊消息處理器
	nslgc.registerMsgHandler(ctx)

	// 運行任務
	nslgc.runTask(ctx)

	// 停服命令
	command.AddCommand(svrmgm.NewSvrStopCmd(ctx.Server, nslgc.serverList))
	command.AddCommand(svrmgm.NewServerCtlCmd(ctx.Server, nslgc.serverList))
	// 關服命令
	command.AddCommand(svrmgm.NewSvrCloseCmd(ctx.Server, nslgc.serverList))
	// 刷新版本號命令
	command.AddCommand(svrmgm.NewRefreshVerCmd(ctx.Server, nslgc.serverList))
	// 停服白名單命令
	command.AddCommand(svrmgm.NewServerStopWhiteIPCmd(ctx.Server, nslgc.serverList))
	// AB服命令
	command.AddCommand(svrmgm.NewSvrMasterSlaveCmd(ctx.Server, nslgc.serverList))

	return nil
}

// Stop ?
func (nslgc *NSLogic) Stop() {
	nslgc.ctx.Server.Stop()
}

// Name ?
func (nslgc *NSLogic) Name() string {
	return "nslogic"
}

func (nslgc *NSLogic) registerMsgParser(ctx *module.MouduleCtx) {
	ctx.RegisterMsgParser(msg.AllMsgBClassID, msg.AllMsgSClassID, msg.NewProtoMsgParser())
}

func (nslgc *NSLogic) registerEvtHandler(ctx *module.MouduleCtx) {
	ctx.OnConnect(func(clt *network.ConnInfo) bool {
		example.invalid("client %s connnect...", clt)
		return true
	})
	ctx.OnDisConnect(func(clt *network.ConnInfo) bool {
		example.invalid("client %s disconnnect...", clt)
		server := nslgc.serverList.GetServerByConn(clt)
		if server.ServerType == conf.ServerTypeHallServer {
			nslgc.hallUserList.Delete(int32(server.ServerID))
		}

		//移除服务器信息
		nslgc.serverList.RemoveServerByConn(clt)

		return true
	})
	ctx.OnMessage(func(clt *network.ConnInfo, m *msg.Message) bool {
		example.invalid("client %s msg:%d_%d", clt, m.BClassID, m.SClassID)

		retmsgs := nslgc.msgHandler.OnMessage(ctx, clt, m)

		if retmsgs != nil {
			for _, retmsg := range retmsgs {
				ctx.Server.SendMessage(clt, retmsg)
			}
		}

		return true
	})
}

func (nslgc *NSLogic) registerMsgHandler(ctx *module.MouduleCtx) {
	nslgc.msgHandler.AddHandler(svrmgm.NewServerRegHandler(nslgc.serverList, ctx.Server))
	nslgc.msgHandler.AddHandler(svrmgm.NewUpdateHallUserListHandler(nslgc.hallUserList))
	nslgc.msgHandler.AddHandler(task.NewUpdateUserNotifyHandler())
	taskctx := task.NewTaskRunCtx(nslgc.serverList, ctx.Server, nslgc.hallUserList)
	nslgc.msgHandler.AddHandler(task.NewFishHandler(taskctx))

}

func (nslgc *NSLogic) runTask(ctx *module.MouduleCtx) {
	var taskProvider TaskProvider

	// 局域網
	if ctx.Config.GetConfig(ConfigFileNameNS).GetConfigValue(ConfigSectionNameNS, "mode") == "1" {
		taskProvider = new(InternalTaskProvider)
	} else { // 外網
		taskProvider = new(ExternalTaskPrvider)
	}

	taskctx := task.NewTaskRunCtx(nslgc.serverList, ctx.Server, nslgc.hallUserList)
	taskProvider.GetTask().Run(taskctx)
}
