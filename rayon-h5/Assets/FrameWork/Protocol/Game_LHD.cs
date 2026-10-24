

enum LHD_GameMessageClassID {
	LHDSettingID=10001,	//设置相关 S2C
	LHDStateID=10002,		//阶段ID S2C
	LHDBetID=10003,			//下注ID C2S
	LHDBroadCastBetID=10004,//广播下注 S2C
	LHDGameResultPokerID=10005,		//游戏结果扑克 S2C
	LHDGameSceneID=10006,	//游戏场景 断线重连
	LHDGameResultWinID=10007,		//游戏结果赢钱 S2C
	LHDBetInfoID=10008,//用户断线重连已下注的消息
	LHDBetRetID=10009,//下注结果 S2C
	LHDShangZhuangID=10010,//上庄请求 C2S
	LHDShangZhuangRetID=10011,//上庄请求结果 S2C
	LHDGetShangZhuangListID=10012,//获取上庄列表 C2S
	LHDShangZhuangListInfoID=10013,//下发上庄列表 S2C
	LHDZhuangInfoID=10014,//庄家信息 S2C
	LHDBetClearID=10015,//玩家清除已下的注 C2S
	LHDResultCardListID=10016,//前20次开牌结果
	LHDBetClearRetID=10017,//清理下注结果
    LHDRichestListID = 10018,//大富豪列表的消息
    LHDLuckyStarListID = 10019,//明燈玩家的消息
};
