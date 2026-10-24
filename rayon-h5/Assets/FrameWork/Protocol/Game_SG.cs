

enum SG_GameMessageClassID {
	SettingID=10001,	//设置相关
	StateID=10002,		//阶段ID S2C
	SGQiangZhuangID=10003,	//抢庄ID C2S
	BroadCastZhuangID=10004,//广播庄家 S2C
	BetID=10005,			//下注ID C2S
	SGBroadCastBetID=10006,//广播下注 S2C
	SGGameResultID=10007,		//游戏结果 S2C
	SGGameSceneID=10008,	//游戏场景
	SGReadyTime=10009,
	SGUserQiangZhuang=10010,//用户抢庄广播
};
