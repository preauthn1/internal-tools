

enum BJL_GameMessageClassID {
	BJLSettingID=10001,	//设置相关 S2C
	BJLStateID=10002,		//阶段ID S2C
	BJLBetID=10003,			//下注ID C2S
	BJLBroadCastBetID=10004,//广播下注 S2C
	BJLGameResultPokerID=10005,		//游戏结果扑克 S2C
	BJLGameSceneID=10006,	//游戏场景 断线重连
	BJLGameResultWinID=10007,		//游戏结果赢钱 S2C
	BJLBetInfoID=10008,//用户断线重连已下注的消息
	BJLBetRetID=10009,//下注结果 S2C
	BJLShangZhuangID=10010,//上庄请求 C2S
	BJLShangZhuangRetID=10011,//上庄请求结果 S2C
	BJLGetShangZhuangListID=10012,//获取上庄列表 C2S
	BJLShangZhuangListInfoID=10013,//下发上庄列表 S2C
	BJLZhuangInfoID=10014,//庄家信息 S2C
	BJLBetClearID=10015,//玩家清除已下的注 C2S
	BJLResultCardListID=10016,//前20次开牌结果
	BJLBetClearRetID=10017,//清理下注结果
};
