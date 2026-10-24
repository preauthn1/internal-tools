
    enum HHDZ_GameMessageClassID
    
    {
        HHDZSettingID = 10001,	//设置相关 S2C
	HHDZStateID = 10002,		//阶段ID S2C
	HHDZBetID = 10003,			//下注ID C2S
	HHDZBroadCastBetID = 10004,//广播下注 S2C
	HHDZGameResultPokerID = 10005,		//游戏结果扑克 S2C
	HHDZGameSceneID = 10006,	//游戏场景 断线重连
	HHDZGameResultWinID = 10007,		//游戏结果赢钱 S2C
	HHDZBetInfoID = 10008,//用户断线重连已下注的消息
	HHDZBetRetID = 10009,//下注结果 S2C
	HHDZShangZhuangID = 10010,//上庄请求 C2S
	HHDZShangZhuangRetID = 10011,//上庄请求结果 S2C
	HHDZGetShangZhuangListID = 10012,//获取上庄列表 C2S
	HHDZShangZhuangListInfoID = 10013,//下发上庄列表 S2C
	HHDZZhuangInfoID = 10014,//庄家信息 S2C
	HHDZBetClearID = 10015,//玩家清除已下的注 C2S
	HHDZResultCardListID = 10016,//前20次开牌结果
	HHDZBetClearRetID = 10017,//清理下注结果
	HHDZRichestListID = 10018,//大富豪列表的消息
	HHDZLuckyStarListID = 10019,//明燈玩家的消息
	HHDZRoundID = 10020,//本局局數
    }

