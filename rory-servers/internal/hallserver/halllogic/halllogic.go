package halllogic

import (
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"
	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/msghandler"
	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/notifyhandler"
	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/rpchandler"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// HallLogic 大廳邏輯
type HallLogic struct {
	*logicmodule.GameServerLogicModule
	serverCtx       *serverhandler.GameServerCtx
	gameserverList  *msghandler.GameServerList
	server          *gameserver.GameServer
	serverCfg       *hallutil.HallServerConfig
	hallverProvider hallutil.HallVerProvider
	hallVerLoader   hallutil.HallVerLoader
	gserverProvider hallutil.GameServerProvider
}

// NewHallLogicModule 創建大廳邏輯對象
func NewHallLogicModule(gs *gameserver.GameServer, cfg *hallutil.HallServerConfig) *HallLogic {
	hl := new(HallLogic)
	hl.serverCtx = gs.GameServerRunCtx
	hl.server = gs
	hl.serverCfg = cfg

	hl.GameServerLogicModule = logicmodule.NewGameServerLogicModule()
	hl.gameserverList = msghandler.NewGameServerList()

	return hl
}

// Name name
func (hl *HallLogic) Name() string {
	return "halllogic"
}

// Run run
func (hl *HallLogic) Run(ctx *module.MouduleCtx) error {

	hl.Init(ctx, hl.server.GameServerRunCtx.PlayerManager, hl.server.GameServerRunCtx.DBMsgQueue, nil)

	hl.registerEvtHandler(ctx)

	hl.registerRCPMsgHandler(ctx)

	hl.registerMsgHandler(ctx)

	hl.registerNotifyMsgHandler(ctx)

	hl.startUpdateHallUserListService()
	return nil
}

// Stop stop
func (hl *HallLogic) Stop() {}

func (hl *HallLogic) registerEvtHandler(ctx *module.MouduleCtx) {

	ctx.OnConnect(func(clt *network.ConnInfo) bool {
		example.invalid("client %s connnect...", clt)

		if clt.IsServerConn(conf.TCPSERVERNAMEHALLCLIENT) {
			// 通知用戶管理器用戶連接
			hl.PlayerManager.SetPlayer(player.NewGuest(clt))
		}
		return true
	})
	ctx.OnDisConnect(func(clt *network.ConnInfo) bool {
		example.invalid("client %s disconnnect...", clt)

		// 客戶端連接
		if clt.IsServerConn(conf.TCPSERVERNAMEHALLCLIENT) {
			pr := hl.PlayerManager.GetPlayerByConnID(clt.ConnID())

			// 登錄用戶退出
			if pr != nil && !pr.IsGuest() {
				uid := pr.UserID
				cer := pr.Cer

				hl.onUserLogout(uid, cer)
			}

			// 通知用戶管理器用戶斷開連接
			hl.PlayerManager.RemovePlayer(clt.ConnID())
		}

		// gameserver連接
		if clt.IsServerConn(conf.TCPSERVERNAMEHALLGAMESERVER) {
			hl.gameserverList.RemoveGameServerByConn(clt)
		}

		return true
	})
	ctx.OnMessage(func(clt *network.ConnInfo, m *msg.Message) bool {
		// 檢查服務器狀態
		if clt.IsServerConn(conf.TCPSERVERNAMEHALLCLIENT) {
			if hl.ServerStatus.GetServerStatus() == conf.ServerStatusClose && !hl.ServerStatus.IsWhiteIP(clt.IPAddress) {
				hl.kickUser(clt.ConnID(), hl.ServerStatus.GetStopServerTipMsg())

				return false
			}
		}

		return true
	})
}

func (hl *HallLogic) registerMsgHandler(ctx *module.MouduleCtx) {
	// 平台邏輯轉發給dbserver進行處理
	msgid := []*msg.RPCMessageID{
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GuestLoginID), false),                            // 遊客登錄
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserLoginID), false),                             // 賬號登錄
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindZhifubaoID), true),                           // 綁定支付寶
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindBankCardID), true),                           // 绑定銀行卡
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SendPhoneVCodeID), false),                        // 發送驗證碼
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyPasswordID), true),                         // 修改密碼
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyFaceID), true),                             // 修改頭像
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_DepositMoneyID), true),                           // 存取款
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindGuestAccountID), false),                      // 綁定遊客賬號
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestUserHallInfoID), true),                    // 請求用戶大廳信息
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_FindSetPwdByPhoneID), false),                     // 請求找回賬號密碼
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyBankPasswordID), true),                     // 請求修改銀行密碼
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_FindSetBankPwdByPhoneID), true),                  // 請求找回銀行密碼
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMailListID), true),                            // 請求郵件列表
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMailDetailID), true),                          // 請求郵件詳情
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ContactServiceID), false),                        // 請求聯繫客服
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ConvertMoneyID), true),                           // 兌換金幣
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentListID), true),                           // 請求代理列表
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ApplyAgentID), true),                             // 請求申請代理
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ReportAgentID), true),                            // 舉報代理
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CreatePayOrderID), true),                         // 請求創建充值訂單
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetFAQID), true),                                 // 請求FAQ信息
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetConvertLogID), true),                          // 查詢兌換記錄
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetBankListID), true),                            // 查詢銀行列表
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestLoginBonusID), true),                      // 查詢是否可領取登入獎勵
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ApplyLoginBonusID), true),                        // 領取登入獎勵
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestAgentListID), true),                       // 查詢代理資訊
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_JoinAgentListID), true),                          // 申請代理下線
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ApplyAgentListID), true),                         // 申請成為代理
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CheckAgentID), true),                             // 查詢申請狀態
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CheckCSNotifyID), true),                          // 未讀客服消息數量
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ProxyRequestID), true),                           // 代理訪問地址
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ProxyPOSTRequestID), true),                       // 代理訪問地址
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UpdateUserNotifyIdxC2HID), true),                 // 代理訪問地址
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GameCompetitionID), true),                        // 查詢查詢比賽活動公告
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID), true),                 // 查詢查詢比賽活動規則
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BuyTicketID), true),                              // 買比賽活動票
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetJackpotPoolID), true),                         // 取得獎金池
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserQACodeID), true),                          // 查詢全民代理QRCode
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentExecID), true),                       // 取得推廣員
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentTaxLogID), true),                     // 取得今日稅收
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentListID), true),                       // 推廣人數
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentTaxBonusLogID), true),                    // 全民推廣 每日收益
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentTaxListID), true),                        // 全民推廣 收益規則
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentUserLevelInfoID), true),                  // 取得推廣等級
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetShanFuID), true),                              // 取得閃付
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetShanFuChatID), true),                          // 取得聊天
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetWithdrawListID), true),                        // 領取收益
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentInfoID), true),             // 玩家/代理信息(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetSubUserListID), true),           // 新增玩家(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentWithdrawAmountID), true),   // 代理提現紀錄(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawAmountID), true),        // 提現紀錄(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetBonusAmountID), true),           // 可提現金額查詢(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentWithDrawOrderID), true),            // 提現(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentApplyListID), true),                // 審核代理申請列表(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentDetailID), true),           // USER 推廣信息(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentReportLogID), true),        // USER 新增代理(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawDetailID), true),        // USER 新增代理(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentUserModeGetID), true),                   // USER 取得已開的模式
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetWithdrawOrderID), true),                       // 馬上提現
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountID), true),                 // 申請提現(全民代理)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountListID), true),             // 申請提現例表(全民代理)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithdrawAmountListID), true),         // 申請提現例表(我的提現)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentResultListID), true),                    // 我的業績(全民代理)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentBorrowListID), true),                    // 我的預提(全民代理)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentWinAmountListID), true),                 // 我的獎勵(全民代理)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMyPromotionID), true),                         // 無限級QRCODE
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetRewardTipsID), true),                          // 無限級獎勵說明
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitConvertOrderID), true),                    // 無限級按鈕1
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitPreConvertOrderID), true),                 // 無限級按鈕2
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyInfoID), true),                         // 裂變1推廣信息
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyConvertID), true),                      // 裂變2提現紀錄
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetRewardTipsID), true),                     // 裂變3說明
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyRewardID), true),                       // 裂變4馬上提現
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitConvertOrderID), true),                      // 裂變5提現紀錄送出提現金額
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitBindingInvitedCodeID), true),                // 裂變6輸入推廣碼
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_EventChargeBonusAmountID), true),                 // 充值返利活動  取得個人流水及限制
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitAddUserRelationID), true),                 // USER 加入下線(無限級)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentAddUserRelationID), true),          // USER 加入下線(淨利潤)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitAddUserRelationID), true),                   // USER 加入下線(裂變)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentAddUserRelationID), true),               // USER 加入下線(全民代理)
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_PromotionEventsPopupID), true),                   // 活動公告系統首頁彈出
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_PromotionEventsAnnouncementID), true),            // 活動公告系統郵件公告頁
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ThirdPartyGameLaunchID), true),                   // 三方遊戲登入
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ThirdPartyGameWalletAutoTransferID), true),       // 三方遊戲額度轉回
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NewPromotionEventsInfoID), true),                 // 活動公告系統首頁彈出
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NewPromotionEventsListID), true),                 // 活動公告系統郵件公告頁
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NewPromotionEventsLuckyRouletteRequestID), true), // 活動公告系統郵件公告頁

		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetSharingPercentageID), true), // 全民淨利潤-獲取設定(%) C2S
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetExecInfoID), true),          // 全民淨利潤-獲取累計收益及推廣人數 C2S
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetDailySharingID), true),      // 全民淨利潤-獲取每日收益列表 C2S
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetWithdrawLogID), true),       // 全民淨利潤-獲取領取收益明細 C2S
		msg.NewRPCMessageID(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitWithdrawLogID), true),          // 全民淨利潤-領取收益 C2S

	}

	rpchandler := logicmodule.NewRouteRPCMsgHandler(hl.DBMsgQueue, ctx, hl.PlayerManager, msgid, int32(hl.server.ServerID))
	hl.MsgHandler.AddHandler(rpchandler) // rpc 轉發

	gsmgmhandler := msghandler.NewGsMgmHandler(hl.gameserverList, hl.serverCfg.GSMgmAddr) // gameserver 管理
	hl.MsgHandler.AddHandler(gsmgmhandler)
	// 服務器集合信息
	hl.gserverProvider = msghandler.NewRemoteRegGameServerProvider(gsmgmhandler, msghandler.NewDefaultServerAllocStrategy(hl.serverCfg))

	// 平台版本號
	hl.MsgHandler.AddHandler(msghandler.NewGetHallVerHandler(hl.hallverProvider))
	hl.gameserverList.SetHallVerProvider(hl.hallverProvider)
	// 服務器列表

	hl.MsgHandler.AddHandler(msghandler.NewGameServerListHandler(hl.gserverProvider, hl.hallverProvider))
	// 服務器分配
	hl.MsgHandler.AddHandler(msghandler.NewAllocGameServerHandler(hl.gserverProvider))

	// 代理訪問支付
	hl.MsgHandler.AddHandler(msghandler.NewProxyRequestHandler(hl.PlayerManager))
	hl.MsgHandler.AddHandler(msghandler.NewProxyPOSTRequestHandler(hl.PlayerManager))
	// 更新用戶推播 idx
	hl.MsgHandler.AddHandler(msghandler.NewUpdateUserNotifyIdxHandler(hl.server, hl.PlayerManager))
}

func (hl *HallLogic) registerRCPMsgHandler(ctx *module.MouduleCtx) {
	// dbserver處理平台邏輯結果轉發給客戶端
	msgid := []*msg.MessageID{
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_LoginRetID)},                                  // 登錄結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindZhifubaoRetID)},                           // 綁定支付寶結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SendPhoneVCodeRetID)},                         // 發送驗證碼結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyPasswordRetID)},                         // 修改密碼結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SendPhoneVCodeRetID)},                         // 發送驗證碼結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyFaceRetID)},                             // 修改頭像結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_DepositMoneyRetID)},                           // 存取款結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindGuestAccountRetID)},                       // 綁定遊客賬號結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserHallInfoID)},                              // 返回用戶大廳信息
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_FindSetPwdByPhoneRetID)},                      // 返回找回密碼信息
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ModifyBankPasswordRetID)},                     // 返回找回密碼信息
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_FindSetBankPwdByPhoneRetID)},                  // 返回找回銀行密碼
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMailListRetID)},                            // 返回郵件列表
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMailDetailRetID)},                          // 返回郵件詳情
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ContactServiceRetID)},                         // 返回聯繫客服
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ConvertMoneyRetID)},                           // 返回兌換金幣
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentListRetID)},                           // 返回代理列表
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ApplyAgentRetID)},                             // 返回申請代理
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ReportAgentRetID)},                            // 返回舉報代理
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CreatePayOrderRetID)},                         // 返回創建充值訂單結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetFAQRetID)},                                 // 返回FAQ信息
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetBankListRetID)},                            // 返回銀行列表
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetConvertLogRetID)},                          // 返回兌換記錄
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BindBankCardRetID)},                           // 綁定銀行卡結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestLoginBonusRetID)},                      // 返回登入獎勵結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestAgentListRetID)},                       // 回傳代理資訊
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_JoinAgentListRetID)},                          // 回傳申請代理結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ApplyAgentListRetID)},                         // 回傳是否成功
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CheckAgentRetID)},                             // 回傳查詢狀態
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CheckCSNotifyRetID)},                          // 回傳未讀客服消息數量
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GameCompetitionRetID)},                        // 回傳比賽活動公告
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GameCompetioniInfoRuleRetID)},                 // 回傳比賽規則
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_BuyTicketRetID)},                              // 回傳買票結果
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetJackpotPoolRetID)},                         // 回傳獎金池
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserQACodeRetID)},                          // 回傳全民代理QRCode
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentExecRetID)},                       // 回傳推廣員
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentTaxLogRetID)},                     // 回傳今日稅收
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentListRetID)},                       // 回傳推廣人數
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentTaxBonusLogRetID)},                    // 回傳全民推廣 每日收益
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentTaxListRetID)},                        // 回傳全民推廣 收益規則
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentUserLevelInfoRetID)},                  // 回傳推廣等級
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetShanFuRetID)},                              // 回傳閃付
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetShanFuChatRetID)},                          // 回傳閃付聊天
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetWithdrawListRetID)},                        // 領取收益
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentInfoRetID)},             // 玩家/代理信息(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetSubUserListRetID)},           // 新增玩家(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentWithdrawAmountRetID)},   // 代理提現紀錄(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawAmountRetID)},        // 提現紀錄(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetBonusAmountRetID)},           // 可提現金額查詢(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentWithDrawOrderRetID)},            // 提現(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentApplyListRetID)},                // 審核代理申請列表(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentDetailRetID)},           // USER 推廣信息(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentReportLogRetID)},        // USER 新增代理(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawDetailRetID)},        // USER 新增代理(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentUserModeGetRetID)},                   // USER 取得已開的模式
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetWithdrawOrderRetID)},                       // 馬上提現
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountRetID)},                 // 申請提現(全民代理)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountListRetID)},             // 申請提現例表(全民代理)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithdrawAmountListRetID)},         // 申請提現例表(我的提現)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentResultListRetID)},                    // 我的業績(全民代理)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentBorrowListRetID)},                    // 我的預提(全民代理)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentWinAmountListRetID)},                 // 我的獎勵(全民代理)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMyPromotionRetID)},                         // 無限級QRCODE
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetRewardTipsRetID)},                          // 無限級獎勵說明
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitConvertOrderRetID)},                    // 無限級按鈕1
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitPreConvertOrderRetID)},                 // 無限級按鈕2
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyInfoRetID)},                         // 裂變1推廣信息
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyConvertRetID)},                      // 裂變2提現紀錄
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetRewardTipsRetID)},                     // 裂變3說明
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyRewardRetID)},                       // 裂變4馬上提現
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitConvertOrderRetID)},                      // 裂變5提現紀錄送出提現金額
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitBindingInvitedCodeRetID)},                // 裂變6輸入推廣碼
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_EventChargeBonusAmountRetID)},                 // 充值返利活動  取得個人流水及限制
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitAddUserRelationRetID)},                 // USER 加入下線(無限級)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentAddUserRelationRetID)},          // USER 加入下線(淨利潤)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitAddUserRelationRetID)},                   // USER 加入下線(裂變)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentAddUserRelationRetID)},               // USER 加入下線(全民代理)
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_PromotionEventsPopupRetID)},                   // 活动公告系统首页弹出
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_PromotionEventsAnnouncementRetID)},            // 活动公告系统邮件公告页
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ThirdPartyGameLaunchRetID)},                   // 三方遊戲登入
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NewPromotionEventsInfoRetID)},                 // 三方遊戲登入
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NewPromotionEventsListRetID)},                 // 三方遊戲登入
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NewPromotionEventsLuckyRouletteRequestRetID)}, // 三方遊戲登入
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetSharingPercentageRetID)},         // 全民淨利潤-獲取設定(%) S2C
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetExecInfoRetID)},                  // 全民淨利潤-獲取累計收益及推廣人數 S2C
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetDailySharingRetID)},              // 全民淨利潤-獲取每日收益列表 S2C
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetWithdrawLogRetID)},               // 全民淨利潤-獲取領取收益明細 S2C
		{int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitWithdrawLogRetID)},                  // 全民淨利潤-領取收益 S2C
	}
	hl.RPCMsgHandler.AddHandler(logicmodule.NewRPCRetHookHandler(msgid)) // 轉發客戶端消息

	hl.RPCMsgHandler.AddHandler(rpchandler.NewLogonRetHandler(hl.PlayerManager, hl.onUserLogout, hl.kickUser)) // 用戶登錄處理

	// 大廳啟動向dbserver註冊
	hallreghandler := rpchandler.NewHallServerRegHandler(hl.serverCfg.HallAddr, hl.serverCfg.NotifyServerAddr, "", int32(hl.server.ServerID), hl.DBMsgQueue, hl.server.NotifyHandler)
	hl.RPCMsgHandler.AddHandler(hallreghandler)
	hallreghandler.SendDBServerStartMsg()

	// 更新遊戲分類清單
	gamecategoryhandler := hallutil.NewGameCategoryHelper(hl.DBMsgQueue, int32(hl.server.ServerID), hl.serverCfg.RefreshCategoryTime)
	hl.RPCMsgHandler.AddHandler(gamecategoryhandler)
	gamecategoryhandler.SendDBGetGameCategoryMsg()

	// 平台版本號管理
	hl.hallverProvider = rpchandler.NewRPCHallVerProvider(hallreghandler, hl.PlayerManager)
	// 平台版本號加載器
	hl.hallVerLoader = rpchandler.NewRPCHallVerLoader(hallreghandler)
	command.AddCommand(rpchandler.NewHallVerCmd(hl.hallverProvider))
}

// 註冊通知處理器
func (hl *HallLogic) registerNotifyMsgHandler(ctx *module.MouduleCtx) {
	// 公告
	hl.NotifyMsgHandler.AddHandler(notifyhandler.NewAnnNotifyHandler())
	// 捕魚公告
	hl.NotifyMsgHandler.AddHandler(notifyhandler.NewAnnNotifyFishHandler())
	// 用戶屬性變化
	hl.NotifyMsgHandler.AddHandler(notifyhandler.NewAttrNotifyHandler())
	// 強制刷新版本號
	hl.NotifyMsgHandler.AddHandler(notifyhandler.NewRefreshVerHandler(hl.hallVerLoader))
	// 大廳用戶到帳通知
	hl.NotifyMsgHandler.AddHandler(notifyhandler.NewNotifyPayOrderHandler())
	// 大廳用戶推播消息
	hl.NotifyMsgHandler.AddHandler(notifyhandler.NewUserNotifyHandler())
}

// 用戶登出大廳
func (hl *HallLogic) onUserLogout(uid int32, cer string) {
	logoutHallMsg := new(msg.Message)
	logoutHallMsg.BClassID = int32(netproto.MessageBClassID_Hall)
	logoutHallMsg.SClassID = int32(netproto.HallMsgClassID_UserLogoutID)

	logoutHallMsgData := new(netproto.UserLogout)
	logoutHallMsgData.UserID = proto.Int32(uid)
	logoutHallMsgData.Cer = proto.String(cer)
	logoutHallMsg.MsgData = logoutHallMsgData

	hl.DBMsgQueue.PushMessage(netproto.CreateEmptyRPCInfo(int32(hl.server.ServerID)), false, logoutHallMsg)
}

// 踢出用户
func (hl *HallLogic) kickUser(connid string, msgstr string) {
	example.invalid("%s,断开客户端%s连接", msgstr, connid)

	hl.sendTipMsg(connid, msgstr)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)

		hl.MsgServer.CloseConnByConnID(connid)
	}()

}

// 發送提示消息
func (hl *HallLogic) sendTipMsg(connid string, msgstr string) {
	m, _ := hallutil.GetTipMessage(msgstr, hallutil.MsgTypeIDHallKick)

	hl.serverCtx.MsgServer.SendMessageByConnID(connid, m)
}

// 消息通知服務
// 定時向 NotifyServer 更新目前大廳在線的用戶
func (hl *HallLogic) startUpdateHallUserListService() {
	go func() {
		defer util.PanicErrStack()
		for {
			time.Sleep(time.Duration(5) * time.Second)
			HallUserList := new(netproto.HallUserList)
			HallUserList.ServerID = proto.Int32(hl.serverCfg.ServerID)
			hl.PlayerManager.Range(func(p *player.Player) {
				if p.IsOnline {
					HallUserList.UserIDList = append(HallUserList.UserIDList, int32(p.UserID))
				}
			})
			msg := new(msg.Message)
			msg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
			msg.SClassID = int32(netproto.NotifyServerClassID_UpdateHallUserListID)
			msg.MsgData = HallUserList
			hl.server.NotifyHandler.TCPMsgClient.SendMessage(msg)
		}
		mlog.Error("---退出更新大廳用戶線程---")
	}()
}
