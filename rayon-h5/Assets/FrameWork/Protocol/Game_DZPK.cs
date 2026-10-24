//==================================================================
//文件名： Game_DZPK.proto
//作者：   dyb
//日期：   2017/9/25
//说明：   德州扑克游戏客户服务器通信协议
//修改记录：
//==================================================================


enum DZPK_GameMessageClassID {
	//BroadCastBeginID=10000,			//服务器推送：开始游戏
	BroadCastDealID=10001,				//服务器推送：发牌,DZPK_PrivatePoker
	BroadCastWaitOptID=10002,			//服务器推送：等待操作,DZPK_WaitOptSeat
	BroadCastCurrOptID=10003,			//服务器推送：当前操作的椅子号
	BroadCastCardsInfoID=10004,			//服务器推送：手牌信息
	
	ActionBetID=10005,					//客户端请求：押注,DZPK_Bet
	ActionCallID=10006,					//客户端请求：跟注,DZPK_Call
	ActionRaiseID=10007,				//客户端请求：加注,DZPK_Raise
	ActionCheckID=10008,				//客户端请求：过，看牌,DZPK_Check
	ActionAllinID=10009,				//客户端请求：allin,DZPK_AllIn
	ActionFoldID=10010,					//客户端请求：弃牌,DZPK_Fold

	RequestStatisticsID=10011,			//客户端请求：查看统计数据
	RequestRePlayID=10012,				//客户端请求：查看回放数据

	BroadCastBetBaseID=10013,			//服务器推送：下底广播,DZPK_BetBase
	BroadCastBetID=10014,				//服务器推送：押注广播,DZPK_BrodCastBet
	BroadCastCallID=10015,				//服务器推送：跟注广播,DZPK_BrodCastCall
	BroadCastRaiseID=10016,				//服务器推送：加注广播,DZPK_BrodCastRaise
	BroadCastCheckID=10017,				//服务器推送：过，看牌广播,DZPK_BrodCastCheck
	BroadCastAllinID=10018,				//服务器推送：allin广播,DZPK_BrodCastAllIn
	BroadCastFoldID=10019,				//服务器推送：弃牌广播,DZPK_BrodCastFold

	BroadCastSettlementID=10020,		//服务器推送：单局结算广播,DZPK_SettleMentRound
	BroadCastPotInfoID=10021,			//服务器推送：彩池信息广播,DZPK_BroadCastPotInfo
	BroadCastSenceRecoverID=10022,		//服务器推送：场景恢复广播,DZPK_GameSence
	BroadCastDealPublicID = 10023,		//服务器推送：发公共牌,DZPK_BrodCastCardsPoker
	BroadCastCardTypeInfoID = 10024,	//服务器推送：用户牌型信息,DZPK_BroadCastCardType
	BroadCastCheckCardsID = 10025,		//服务器推送：看牌比牌,DZPK_UserCardsInfo
	BroadCastGameSettingID = 10026,		//服务器推送：游戏设置,DZPK_GameSetting
	BroadCastGameSenceID = 10027,		//服务器推送：游戏场景,DZPK_GameSence
	BroadCastPreOprationID = 10029,		//服务器推送：预操作,DZPK_PreOpration

	BroadCastErrorID = 10028,			//服务器推送：操作失败DZPK_ErrorMsg
	ActionGMCommond = 10030,			//测试GM命令:DZPK_GM
};