package dbserver

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/msghandler"
	"example.invalid/ServerTeam/Servers/internal/phnum/dbserver/handler"
	slot "example.invalid/ServerTeam/Servers/pkg/slot/dbhandler"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// DBLogic DB Logic Module
type DBLogic struct {
	mouduleCtx *module.MouduleCtx
	msgHandler *module.SemaphoreHookMsgHandlerGroup
}

// NewDBLogicModule Create DBLogic
func NewDBLogicModule(n int) *DBLogic {
	db := new(DBLogic)
	db.msgHandler = module.NewSemaphoreHookMsgHandlerGroup(n)

	return db
}

// Name 模組名稱
func (dl *DBLogic) Name() string {
	return "DBLogic"
}

// Run Start DB Logic Module
func (dl *DBLogic) Run(ctx *module.MouduleCtx) error {
	// Setting DB Connection
	dbutil.ConnString = ctx.GetConfigData("db", "conn", "connstring")
	dbutil.ConnString = strings.Replace(dbutil.ConnString, "\"", "", -1)

	dbutil.DealLogConnString = ctx.GetConfigData("db", "conn", "deallogconnstring")
	dbutil.DealLogConnString = strings.Replace(dbutil.DealLogConnString, "\"", "", -1)
	dbutil.ClubConnString = ctx.GetConfigData("db", "conn", "clubconnstring")
	dbutil.ClubConnString = strings.Replace(dbutil.ClubConnString, "\"", "", -1)

	dbutil.SMSWebAPI = ctx.GetConfigData("db", "conn", "smswebapi")

	// Setting Connection Pool
	psize := ctx.GetConfigData("db", "conn", "pool")
	ipsize, _ := strconv.Atoi(psize)
	dbutil.PoolSize = ipsize

	dbutil.InitDatabase()

	dl.mouduleCtx = ctx

	dl.registerMsgParser(ctx)
	dl.registerMsgHandler(ctx)

	dl.registerEvtHandler(ctx)

	return nil
}

// Stop Stop DB Logic Module
func (dl *DBLogic) Stop() {
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("[DBLogic]dbserver Stop")).Error("HallServer壓測時")
	dl.mouduleCtx.Server.Stop()
}

func (dl *DBLogic) registerMsgParser(ctx *module.MouduleCtx) {
	ctx.RegisterMsgParser(msg.AllMsgBClassID, msg.AllMsgSClassID, rpcmsg.NewRPCMessageParser())
}

func (dl *DBLogic) registerEvtHandler(ctx *module.MouduleCtx) {

	ctx.OnConnect(func(clt *network.ConnInfo) bool {
		ctx.Trace("client %s connnect...", clt)
		return true
	})
	ctx.OnDisConnect(func(clt *network.ConnInfo) bool {
		ctx.Trace("client %s disconnnect...", clt)
		return true
	})
	ctx.OnMessage(func(clt *network.ConnInfo, m *msg.Message) bool {
		ctx.Trace("client %s msg:%d_%d", clt, m.BClassID, m.SClassID)

		//填充RPC信息
		fromrpcmsg := m.MsgData.(*rpcmsg.RPCMessage)

		retmsgs := dl.msgHandler.OnMessage(ctx, clt, m)

		if retmsgs != nil {
			for _, retmsg := range retmsgs {
				rmsg := new(rpcmsg.RPCMessage)
				rmsg.RPCInfo = fromrpcmsg.RPCInfo
				rmsg.MsgData, _ = proto.Marshal(retmsg.MsgData.(proto.Message))

				retmsg.MsgData = rmsg
				ctx.Trace("return msg:%d_%d", retmsg.BClassID, retmsg.SClassID)
				ctx.Server.SendMessage(clt, retmsg)
			}
		}
		//example.invalid("client %s msg:%d_%d dealed...", clt, m.BClassID, m.SClassID)

		return true
	})
}

func (dl *DBLogic) registerMsgHandler(ctx *module.MouduleCtx) {
	// 調整過封包
	dl.msgHandler.AddHandler(new(handler.GameServerStartHandler)) // 房间启动
	dl.msgHandler.AddHandler(new(handler.ServerRegHandler))       // 服务器注册
	dl.msgHandler.AddHandler(new(handler.RobotLogonHandler))      // 加载机器人
	dl.msgHandler.AddHandler(new(handler.RoomLogonHandler))       // 登录房间
	dl.msgHandler.AddHandler(new(handler.RoomLogoutHandler))      // 登出房间
	dl.msgHandler.AddHandler(new(handler.AddPerGameLogHandler))   // 写每局记录

	// 原始
	dl.msgHandler.AddHandler(new(msghandler.UpdateUserScoreHandler)) // 即时写分
	dl.msgHandler.AddHandler(new(msghandler.BlockIPHandler))         // 封锁IP
	dl.msgHandler.AddHandler(new(msghandler.LoadBlockIPListHandler)) // 加载封锁IP列表

	// Slot 專區
	dl.msgHandler.AddHandler(slot.NewSelectGameHandler(slot.DB))       // 寫入及讀取選擇遊戲機會
	dl.msgHandler.AddHandler(slot.NewFreeSpinHandler(slot.DB))         // 寫入及讀取免費遊戲
	dl.msgHandler.AddHandler(slot.NewJackpotNoLinkHandler(slot.DB))    // 寫入及讀取JackpotNoLink資訊
	dl.msgHandler.AddHandler(slot.NewGameLogHandler(slot.DB))          // 寫入及讀取遊戲歷程
	dl.msgHandler.AddHandler(slot.NewBonusGameHandler(slot.DB))        // 寫入及讀取翻牌遊戲資料
	dl.msgHandler.AddHandler(slot.NewDeferredAwardHandler(slot.DB))    // 寫入及讀取延遲性獎勵
	dl.msgHandler.AddHandler(slot.NewLightningHandler(slot.DB))        // 寫入及讀取Lightning遊戲資料
	dl.msgHandler.AddHandler(slot.NewUserAssignedModeHandler(slot.DB)) // 寫入及讀取UserAssignedMode遊戲資料

}
