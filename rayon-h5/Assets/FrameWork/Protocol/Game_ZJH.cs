

enum ZJH_GameMessageClassID {
	BetBaseId=10001,		//下底 ZJH_BetBase
	FaPaiId=10002,			//发牌
	WaitOptID=10003,		//等待操作 ZJH_WaitOpt
	CurrOptID=10004,		//当前操作的椅子号 ZJH_WaitOptSeat
	CheckCardsID=10005,		//查看手牌 客户端向服务器发送
	CheckCardsBroadCastID=10006,//查看手牌广播
	CardsInfo=10007,		//手牌信息 ZJH_Poker
	CompareID=10008,		//比牌消息 客户端向服务器发送 ZJH_CompareCards
	CompareRetID=10009,		//比牌结果 ZJH_CompareCardsRet
	DropID=10010,			//弃牌消息 ZJH_Drop
	DropBroadCastID=10011,	//弃牌广播
	FollowID=10012,			//跟注		
	FollowBrodCastID=10013,	//服务器下发客户端 跟注广播 ZJH_Follow
	BetAddID=10014,			//加注 客户端发送 ZJH_BetAdd
	BetAddBrodCastID=10015, //加注广播 ZJH_BetAddBrodCast
	GameResultID=10016,		//游戏计算消息 ZJH_GameRet
	ComparePokerID=10017,	//和玩家比过牌的牌型和牌数据
	ZJHGameSceneID=10018,   //游戏场景恢复消息
	ReadyTime=10019,		//游戏结束准备时间
	FollowEndID=10020,		//跟到底
    EarlyChangeTable = 10021,//提早換桌
    AllinID = 10022,        //全押
    AllinRetID = 10023,     // 全押回覆
    EarlyLeaveTable = 10024,// 提早離場
	GameCurrturn=10025,		//現在第幾輪
	AllInSettlement=10026,	//全押結算
	UserCardType=10027,		//用戶牌型
};