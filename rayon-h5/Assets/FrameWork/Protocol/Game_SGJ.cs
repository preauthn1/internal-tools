

enum SGJ_GameMessageClassID {
	SGJSettingID=10001,	//设置相关 S2C
	SGJStateID=10002,		//阶段ID S2C
	SGJShangZhuangID=10003,	//上庄ID C2S
	SGJGetZhuangListID=10004,//获取庄的列表 S2C
	SGJBetID=10005,			//下注ID C2S
	SGJBroadCastBetID=10006,//广播下注 S2C
	SGJBroadCastJiangType=10007,//广播开奖类型 S2C
	SGJGameResultID=10008,		//游戏结果 S2C
	SGJLastGameResultID=10009,//上次开奖结果，在用户进入时发送
	SGJZhuangInfoID=10010,//庄家信息
	SGJShangZhuangListID=10011,//下发庄家信息
	SGJUserBetInfo=10012,//玩家下注消息断线重连
	SGJCannotBetID=10013,//钱不够不能下注了
	SGJShangzhuangRetID=10014,//上庄结果
	SGJClearUserBetID=10015,//清除下注数据 C2S
	SGJBetClearRetID=10016,//清除数据结果
};
