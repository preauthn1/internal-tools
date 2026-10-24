package apilogic

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/msghandler"
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	slot "example.invalid/ServerTeam/Servers/pkg/slot/dbhandler"
)

// APILogic API Logic Module
type APILogic struct {
	mouduleCtx *module.MouduleCtx
	msgHandler *module.SemaphoreHookMsgHandlerGroup
}

// NewAPILogicModule Create APILogic
func NewAPILogicModule(n int) *APILogic {
	logic := new(APILogic)
	logic.msgHandler = module.NewSemaphoreHookMsgHandlerGroup(n)

	return logic
}

// Name 模組名稱
func (l *APILogic) Name() string {
	return "APILogic"
}

// Run Start API Logic Module
func (l *APILogic) Run(ctx *module.MouduleCtx) error {

	ctx.BackendAPI = ctx.GetConfigData("api", "conn", "backendapi")
	ctx.APICertificate = ctx.GetConfigData("api", "conn", "api_cer")

	service.InitService(ctx)

	httputil.APIUtilCmdP = httputil.NewAPIUtilCmd()
	command.AddCommand(httputil.APIUtilCmdP)

	l.mouduleCtx = ctx

	l.registerMsgParser(ctx)
	l.registerMsgHandler(ctx)

	l.registerEvtHandler(ctx)

	return nil
}

// Stop API Logic Module
func (l *APILogic) Stop() {
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("[APILogic]apiserver Stop")).Error("APIServer壓測時")
	l.mouduleCtx.Server.Stop()
}

func (l *APILogic) registerMsgParser(ctx *module.MouduleCtx) {
	ctx.RegisterMsgParser(msg.AllMsgBClassID, msg.AllMsgSClassID, rpcmsg.NewRPCMessageParser())
}

func (l *APILogic) registerEvtHandler(ctx *module.MouduleCtx) {

	ctx.OnConnect(func(clt *network.ConnInfo) bool {
		example.invalid("client %s connnect...", clt)
		return true
	})
	ctx.OnDisConnect(func(clt *network.ConnInfo) bool {
		example.invalid("client %s disconnnect...", clt)
		return true
	})
	ctx.OnMessage(func(clt *network.ConnInfo, m *msg.Message) bool {
		example.invalid("client %s msg:%d_%d", clt, m.BClassID, m.SClassID)

		// 填充RPC信息
		fromrpcmsg := m.MsgData.(*rpcmsg.RPCMessage)

		retmsgs := l.msgHandler.OnMessage(ctx, clt, m)

		if retmsgs != nil {
			for _, retmsg := range retmsgs {
				rmsg := new(rpcmsg.RPCMessage)
				rmsg.RPCInfo = fromrpcmsg.RPCInfo
				rmsg.MsgData, _ = proto.Marshal(retmsg.MsgData.(proto.Message))

				retmsg.MsgData = rmsg
				example.invalid("return msg:%d_%d", retmsg.BClassID, retmsg.SClassID)
				ctx.Server.SendMessage(clt, retmsg)
			}
		}
		//example.invalid("client %s msg:%d_%d dealed...", clt, m.BClassID, m.SClassID)

		return true
	})
}

func (l *APILogic) registerMsgHandler(ctx *module.MouduleCtx) {
	l.msgHandler.AddHandler(new(msghandler.ServerRegHandler))            // 服務器註冊
	l.msgHandler.AddHandler(new(msghandler.RoomLogonHandler))            // 登錄房間
	l.msgHandler.AddHandler(new(msghandler.RoomLogoutHandler))           // 登出房間
	l.msgHandler.AddHandler(new(msghandler.UpdateUserScoreHandler))      // 即時寫分
	l.msgHandler.AddHandler(new(msghandler.GameServerStartHandler))      // 房間啟動
	l.msgHandler.AddHandler(new(msghandler.LoadBlockIPListHandler))      // 加載封鎖IP列表
	l.msgHandler.AddHandler(new(msghandler.AddPerGameLogHandler))        // 寫每局記錄
	l.msgHandler.AddHandler(new(msghandler.RobotLogonHandler))           // 加載機器人
	l.msgHandler.AddHandler(new(msghandler.SuperToolPoolWinLoseHandler)) // 系統水池記錄

	// Slot 專區
	l.msgHandler.AddHandler(slot.NewSelectGameHandler(slot.API))       // 寫入及讀取選擇遊戲機會
	l.msgHandler.AddHandler(slot.NewFreeSpinHandler(slot.API))         // 寫入及讀取免費遊戲
	l.msgHandler.AddHandler(slot.NewJackpotNoLinkHandler(slot.API))    // 寫入及讀取JackpotNoLink資訊
	l.msgHandler.AddHandler(slot.NewGameLogHandler(slot.API))          // 寫入及讀取遊戲歷程
	l.msgHandler.AddHandler(slot.NewBonusGameHandler(slot.API))        // 寫入及讀取翻牌遊戲資料
	l.msgHandler.AddHandler(slot.NewDeferredAwardHandler(slot.API))    // 寫入及讀取延遲性獎勵
	l.msgHandler.AddHandler(slot.NewLightningHandler(slot.API))        // 寫入及讀取Lightning遊戲資料
	l.msgHandler.AddHandler(slot.NewUserAssignedModeHandler(slot.API)) // 寫入及讀取UserAssignedMode遊戲資料
}
