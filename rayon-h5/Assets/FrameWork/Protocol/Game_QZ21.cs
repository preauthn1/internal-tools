

enum QZ21_GameMessageClassID {
	QZ21ChoseOddsId 			= 10001,//选择倍率
	QZ21ChoseOddsRetId			= 10002,//倍率选择操作返回
	QZ21RoleOptId				= 10003,//玩家操作信息 包括叫牌及完成
	QZ21RoleOptRetId			= 10004,//玩家操作返回
	
	QZ21GameInfoPushId			= 10005,//游戏内容推送（用于断线重连，场景恢复）
	QZ21GameDeskCfgId			= 10006,//游戏配置
	
	QZ21NotifyOddsId 			= 10101,//广播倍率相关信息
	QZ21NotifyTimeOutOptId 		= 10102,//广播超时相关信息
	QZ21NotifyChoseOddsId		= 10103,//广播玩家倍率选择相关信息
	QZ21NotifyRoleStatusId		= 10104,//广播玩家操作状态
	QZ21NotifyShowId			= 10105,//亮牌
	QZ21NotifySettlementId		= 10106,//广播游戏结算信息	
	QZ21NotifyBankerId			= 10107,//广播庄家信息
	QZ21NotifyRoleCardChangeId	= 10108,//广播玩家牌数量变化
};