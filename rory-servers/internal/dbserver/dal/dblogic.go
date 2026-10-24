package dal

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/msghandler"
	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/msghandler/api"
	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/msghandler/userprofit"
	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/service"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	slot "example.invalid/ServerTeam/Servers/pkg/slot/dbhandler"
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
	service.InitService(ctx)

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

		// 填充RPC信息
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
	dl.msgHandler.AddHandler(new(msghandler.GuestLogonHandler))                      // 遊客登錄
	dl.msgHandler.AddHandler(new(msghandler.UserLogonHandler))                       // 用戶登錄
	dl.msgHandler.AddHandler(new(msghandler.UserLeaveHallHandler))                   // 登出廣場
	dl.msgHandler.AddHandler(new(msghandler.ServerRegHandler))                       // 服務器註冊
	dl.msgHandler.AddHandler(new(msghandler.RoomLogonHandler))                       // 登錄房間
	dl.msgHandler.AddHandler(new(msghandler.RoomLogoutHandler))                      // 登出房間
	dl.msgHandler.AddHandler(new(msghandler.UpdateUserScoreHandler))                 // 即時寫分
	dl.msgHandler.AddHandler(new(msghandler.UpdateExperienceScoreHandler))           // 體驗金幣
	dl.msgHandler.AddHandler(new(msghandler.HallServerStartHandler))                 // 大廳啟動
	dl.msgHandler.AddHandler(new(msghandler.GameServerStartHandler))                 // 房間啟動
	dl.msgHandler.AddHandler(new(msghandler.BindZhifubaoHandler))                    // 綁定支付寶
	dl.msgHandler.AddHandler(new(msghandler.BindBankCardHandler))                    // 绑定銀行卡
	dl.msgHandler.AddHandler(new(msghandler.ModifyPasswordHandler))                  // 修改密碼
	dl.msgHandler.AddHandler(new(msghandler.SendPhoneVCodeHandler))                  // 發送驗證碼
	dl.msgHandler.AddHandler(new(msghandler.ModifyFaceHandler))                      // 修改頭像
	dl.msgHandler.AddHandler(new(msghandler.MoneyDepositHandler))                    // 存取款
	dl.msgHandler.AddHandler(new(msghandler.BindGuestAccountHandler))                // 綁定遊客賬號
	dl.msgHandler.AddHandler(new(msghandler.RequestUserInfoHandler))                 // 請求用戶大廳信息
	dl.msgHandler.AddHandler(new(msghandler.FindSetPwdByPhoneHandler))               // 請求找回賬號密碼
	dl.msgHandler.AddHandler(new(msghandler.ModifyBankPasswordHandler))              // 請求修改銀行密碼
	dl.msgHandler.AddHandler(new(msghandler.FindSetBankPwdByPhoneHandler))           // 請求找回銀行賬號密碼
	dl.msgHandler.AddHandler(new(msghandler.GetMailListHandler))                     // 請求郵件列表
	dl.msgHandler.AddHandler(new(msghandler.GetMailDetailHandler))                   // 請求郵件詳情
	dl.msgHandler.AddHandler(new(msghandler.ContactServiceHandler))                  // 發送聯繫客服信息
	dl.msgHandler.AddHandler(new(msghandler.ConvertMoneyHandler))                    // 兌換金幣
	dl.msgHandler.AddHandler(new(msghandler.GetAgentListHandler))                    // 請求代理列表
	dl.msgHandler.AddHandler(new(msghandler.ApplyAgentHandler))                      // 申請代理
	dl.msgHandler.AddHandler(new(msghandler.ReportAgentHandler))                     // 舉報代理
	dl.msgHandler.AddHandler(new(msghandler.CreatePayOrderHandler))                  // 創建充值訂單
	dl.msgHandler.AddHandler(new(msghandler.BlockIPHandler))                         // 封鎖IP
	dl.msgHandler.AddHandler(new(msghandler.LoadBlockIPListHandler))                 // 加載封鎖IP列表
	dl.msgHandler.AddHandler(new(msghandler.AddPerGameLogHandler))                   // 寫每局記錄
	dl.msgHandler.AddHandler(new(msghandler.RobotLogonHandler))                      // 加載機器人
	dl.msgHandler.AddHandler(new(msghandler.GetFAQHandler))                          // 加載機器人
	dl.msgHandler.AddHandler(new(msghandler.GetUserLuckScoreHandler))                // 獲取用戶幸運值
	dl.msgHandler.AddHandler(new(msghandler.GetConvertLogHandler))                   // 查詢兌換記錄
	dl.msgHandler.AddHandler(new(msghandler.GetBankListHandler))                     // 查詢銀行列表
	dl.msgHandler.AddHandler(new(msghandler.GameVerHandler))                         // 更新遊戲版本號
	dl.msgHandler.AddHandler(new(msghandler.RequestLoginBonusHandler))               // 登入獎勵操作
	dl.msgHandler.AddHandler(new(msghandler.RequestAgentDetailHandler))              // 取得代理資訊
	dl.msgHandler.AddHandler(new(msghandler.CSNotifyCountHandler))                   // 未讀客服消息數量
	dl.msgHandler.AddHandler(new(msghandler.SuperToolAddAssignRewardHandler))        // 新增指定多局開獎 <--- 多此一舉
	dl.msgHandler.AddHandler(new(msghandler.SuperToolGetAssignRewardDetailHandler))  // 獲取指定開獎項目
	dl.msgHandler.AddHandler(new(msghandler.SuperToolUpdAssignRewardDetailHandler))  // 更新指定開獎項目結果
	dl.msgHandler.AddHandler(new(msghandler.SuperToolAddChangeLogHandler))           // 操作獎池及指定放分的記錄
	dl.msgHandler.AddHandler(new(msghandler.SuperToolGetVipUserHandler))             // 獲取監控用戶列表
	dl.msgHandler.AddHandler(new(msghandler.SuperToolUpdVIPUserHandler))             // 更新監控用戶數據
	dl.msgHandler.AddHandler(new(msghandler.SuperToolFishUpdSystemRewardHandler))    // 捕魚更新系統放分和取得系統放分
	dl.msgHandler.AddHandler(new(msghandler.SuperToolFishUpdVIPUserHandler))         // 捕魚個人放分
	dl.msgHandler.AddHandler(new(msghandler.SuperToolZjhGetVipUserHandler))          // ZJH獲取監控用戶列表
	dl.msgHandler.AddHandler(new(msghandler.SuperToolZjhUpdVIPUserHandler))          // ZJH更新監控用戶數據
	dl.msgHandler.AddHandler(new(msghandler.FishSysPoolHandler))                     // Fish 新增/更新系統水池输赢记录
	dl.msgHandler.AddHandler(new(msghandler.FishPlayerPoolHandler))                  // 新增/更新玩家个人水池输赢记录
	dl.msgHandler.AddHandler(new(msghandler.GetBeginnerOfferHandler))                // 取得玩家新手優惠
	dl.msgHandler.AddHandler(new(msghandler.UpdateBeginnerOfferHandler))             // 更新玩家新手優惠
	dl.msgHandler.AddHandler(new(msghandler.GameCompetitionHandler))                 // 查詢比賽活動公告
	dl.msgHandler.AddHandler(new(msghandler.GameCompetitionRuleHandler))             // 查詢比賽活動規則
	dl.msgHandler.AddHandler(new(msghandler.UpdatePiSaiScoreHandler))                // 比賽積分写分
	dl.msgHandler.AddHandler(new(msghandler.GameCompetitionStartHandler))            // 報名參加比賽活動
	dl.msgHandler.AddHandler(new(msghandler.GetContestUserDataHandler))              // 取得比賽活動用戶資料
	dl.msgHandler.AddHandler(new(msghandler.UpdateContestUserDataHandler))           // 更新比賽活動用戶資料
	dl.msgHandler.AddHandler(new(msghandler.Cd21SysPoolHandler))                     // 21點 新增/更新系統水池输赢记录
	dl.msgHandler.AddHandler(new(msghandler.GetGameJackpotPoolHandler))              // 取得遊戲獎金池
	dl.msgHandler.AddHandler(new(msghandler.UpdateGameJackpotPoolHandler))           // 更新遊戲獎金池
	dl.msgHandler.AddHandler(new(msghandler.UpdateUserJackpotBoxHandler))            // 遊戲獎金寫入保險箱
	dl.msgHandler.AddHandler(new(msghandler.GetUserQRCodeHandler))                   // 全民推廣 我的推廣碼
	dl.msgHandler.AddHandler(new(msghandler.GetUserAngentExecInfoHandler))           // 全民推廣員Agent Exec
	dl.msgHandler.AddHandler(new(msghandler.GetUserAgentListHandler))                // 全民推廣 推廣人數
	dl.msgHandler.AddHandler(new(msghandler.GetAgentTaxBonusLogHandler))             // 全民推廣 每日收益
	dl.msgHandler.AddHandler(new(msghandler.GetAgentTaxListHandler))                 // 全民推廣 收益規則
	dl.msgHandler.AddHandler(new(msghandler.GetUserAngentTaxLogHandler))             // 全民推廣今日稅收 Tax Log
	dl.msgHandler.AddHandler(new(msghandler.GetShanFuHandler))                       // 閃付
	dl.msgHandler.AddHandler(new(msghandler.GetShanFuChatHandler))                   // 閃付聊天
	dl.msgHandler.AddHandler(new(msghandler.GetUserAngentUserLevelInfoHandler))      // 全民推廣今日稅收 Tax Log
	dl.msgHandler.AddHandler(new(msghandler.GetWithdrawListIHandler))                // 全民推廣領取收益
	dl.msgHandler.AddHandler(new(msghandler.GetAgentInfoHandler))                    // 玩家/代理信息(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.GetSubUserListRetHandler))               // 新增玩家(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.GetAgentWithdrawAmountHandler))          // 代理提現紀錄(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.GetWithdrawAmountHandler))               // 提現紀錄(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.GetBonusAmountHandler))                  // 可提現金額查詢(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.WithDrawOrderIDHandler))                 // 可提現金額查詢(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.ApplyListHandler))                       // 審核代理申請列表(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.NetProfitAgentGetAgentDetailRetHandler)) // USER 推廣信息(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.NetProfitAgentGetAgentReportLogHandler)) // USER 新增代理(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.NetProfitAgentGetWithdrawDetailHandler)) // USER 新增代理(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.UserAgentUserModeGetHandler))            // USER 取得已開的模式
	dl.msgHandler.AddHandler(new(msghandler.GetWithdrawOrderHandler))                // 馬上提現
	dl.msgHandler.AddHandler(new(msghandler.UserAgentGetWithAmountHandler))          // 申請提現(全民代理)
	dl.msgHandler.AddHandler(new(msghandler.UserAgentGetWithAmountListHandler))      // 申請提現例表(全民代理)
	dl.msgHandler.AddHandler(new(msghandler.UserAgentGetWithdrawAmountListHandler))  // 申請提現例表(我的提現)
	dl.msgHandler.AddHandler(new(msghandler.UserAgentResultListHandler))             // 我的業績(全民代理)
	dl.msgHandler.AddHandler(new(msghandler.UserAgentBorrowListHandler))             // 我的預提(全民代理)
	dl.msgHandler.AddHandler(new(msghandler.UserAgentWinAmountListHandler))          // 我的獎勵(全民代理)
	dl.msgHandler.AddHandler(new(msghandler.GetMyPromotionHandler))                  // 無限級QRCODE
	dl.msgHandler.AddHandler(new(msghandler.GetRewardTipsHandler))                   // 無限級獎勵說明
	dl.msgHandler.AddHandler(new(msghandler.UnlimitConvertOrderHandler))             // 無限級按鈕1
	dl.msgHandler.AddHandler(new(msghandler.UnlimitPreConvertOrderHandler))          // 無限級按鈕2
	dl.msgHandler.AddHandler(new(msghandler.SplitGetMyInfoHandler))                  // 裂變1推廣信息
	dl.msgHandler.AddHandler(new(msghandler.SplitGetMyConvertHandler))               // 裂變2提現紀錄
	dl.msgHandler.AddHandler(new(msghandler.SplitGetRewardTipsHandler))              // 裂變3說明
	dl.msgHandler.AddHandler(new(msghandler.SplitGetMyRewardHandler))                // 裂變4馬上提現
	dl.msgHandler.AddHandler(new(msghandler.SplitConvertOrderHandler))               // 裂變5提現紀錄送出提現金額
	dl.msgHandler.AddHandler(new(msghandler.SplitBindingInvitedCodeHandler))         // 裂變6輸入推廣碼
	dl.msgHandler.AddHandler(new(msghandler.SlotJackpotLinkHandler))                 // Slot 連接型Jackpot紀錄寫入及讀取
	dl.msgHandler.AddHandler(new(msghandler.EventChargeBonusAmountHandler))          // 充值返利活動  取得個人流水及限制
	dl.msgHandler.AddHandler(new(msghandler.GetAreaConfigHandler))                   // Hall 取得AreaId的AreaConfig資訊
	dl.msgHandler.AddHandler(new(msghandler.UnlimitAddUserRelationHandler))          // USER 加入下線(無限級)
	dl.msgHandler.AddHandler(new(msghandler.NetProfitAgentAddUserRelationHandler))   // USER 加入下線(淨利潤)
	dl.msgHandler.AddHandler(new(msghandler.SplitAddUserRelationHandler))            // USER 加入下線(裂變)
	dl.msgHandler.AddHandler(new(msghandler.UserAgentAddUserRelationHandler))        // USER 加入下線(全民代理)
	dl.msgHandler.AddHandler(new(msghandler.EventGetPromotionPopupHandler))          // 活動公告系統 首頁彈出
	dl.msgHandler.AddHandler(new(msghandler.EventGetPromotionAnnouncementHandler))   // 活動公告系統 郵件公告頁
	dl.msgHandler.AddHandler(new(msghandler.GetServerRelationsHandler))              // 遊戲渠道
	dl.msgHandler.AddHandler(new(msghandler.JackpotNoLinkHandler))                   // 非連接型Jackpot紀錄寫入及讀取
	dl.msgHandler.AddHandler(new(msghandler.JackpotLinkHandler))                     // 連接型Jackpot紀錄寫入及讀取
	dl.msgHandler.AddHandler(new(msghandler.GetGameCategoryListHandler))             // 遊戲分類清單
	//dl.msgHandler.AddHandler(new(msghandler.SuperToolJLBPKHandler))                  // 加勒比海撲克超端
	dl.msgHandler.AddHandler(new(msghandler.SystemPoolHandler)) // 水線控制

	// API Section
	dl.msgHandler.AddHandler(new(api.NewPromotionEventHandler))       // 幸運輪盤相關接口
	dl.msgHandler.AddHandler(new(api.ThirdPartyGameInterfaceHandler)) // 三方遊戲相關接口

	// Slot 專區
	dl.msgHandler.AddHandler(slot.NewSelectGameHandler(slot.DB))       // 寫入及讀取選擇遊戲機會
	dl.msgHandler.AddHandler(slot.NewFreeSpinHandler(slot.DB))         // 寫入及讀取免費遊戲
	dl.msgHandler.AddHandler(slot.NewJackpotNoLinkHandler(slot.DB))    // 寫入及讀取JackpotNoLink資訊
	dl.msgHandler.AddHandler(slot.NewGameLogHandler(slot.DB))          // 寫入及讀取遊戲歷程
	dl.msgHandler.AddHandler(slot.NewBonusGameHandler(slot.DB))        // 寫入及讀取翻牌遊戲資料
	dl.msgHandler.AddHandler(slot.NewDeferredAwardHandler(slot.DB))    // 寫入及讀取延遲性獎勵
	dl.msgHandler.AddHandler(slot.NewLightningHandler(slot.DB))        // 寫入及讀取Lightning遊戲資料
	dl.msgHandler.AddHandler(slot.NewUserAssignedModeHandler(slot.DB)) // 寫入及讀取UserAssignedMode遊戲資料

	// User Profit
	dl.msgHandler.AddHandler(new(userprofit.GetDailySharingHandler))      // 全民淨利潤-獲取每日收益列表
	dl.msgHandler.AddHandler(new(userprofit.GetSharingPercentageHandler)) // 全民淨利潤-獲取設定(%)
	dl.msgHandler.AddHandler(new(userprofit.GetExecInfoHandler))          // 全民淨利潤-獲取累計收益及推廣人數
	dl.msgHandler.AddHandler(new(userprofit.GetWithdrawLogHandler))       // 全民淨利潤-獲取領取收益明細
	dl.msgHandler.AddHandler(new(userprofit.WithdrawRequestHandler))      // 全民淨利潤-領取收益
}
