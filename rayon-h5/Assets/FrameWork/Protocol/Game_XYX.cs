//==================================================================
//文件名： Game_XYX.proto
//作者：   weihao
//日期：   2018/2/5
//说明：   小遊戲客户服务器通信协议
//修改记录：
//==================================================================


enum XYX_GameMessageClassID {
	XYX_GameStepID				= 10001,
	XYX_GameSettingID			= 10002,	//游戏配置
	XYX_GameRetID				= 10003,	//遊戲結果
	XYX_NotifyUserBaseScoreID	= 10004,	//使用者設定的倍數
	XYX_NotifyGameBaseScoreID 	= 10005,	//廣播所有使用者倍數, 若未選擇也會預設為１倍
	XYX_UserBaseScoreID			= 10006,	//Client User Bet 
	XYX_GameSceneID				= 10007,	//游戏场景 断线重连	
};