

enum SSS_GameMessageClassID {
	ChooseSpecialCardTypeId 		= 10001,//牌型选择请求 是否选择特殊牌型
	ChooseSpecialCardTypeRetId		= 10002,//牌型选择响应
	CompleteCardCombId				= 10003,//组牌完成请求
	CompleteCardCombRetId			= 10004,//组牌完成消息响应
	
	PushGameInfoId					= 10005,//推送游戏信息
	PushGameCfgId					= 10006,//房间配置信息
	
	NotifyUserCardTypeChooseId 		= 10101,//广播用户牌型选择 是特殊牌型还是普通牌型
	NotifyUserCardsId				= 10102,//广播用户牌型 普通牌型比较
	NotifyUserSpecialId				= 10103,//广播用户特殊牌型比较 
	NotifyUserGunId					= 10104,//广播用户打枪消息
	NotifySSSSettlementId			= 10105,//广播结算消息
	NotifyOptTimeOutId				= 10106,//广播游戏阶段及超时时间
};