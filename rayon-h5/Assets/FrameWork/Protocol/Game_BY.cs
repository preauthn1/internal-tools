

enum BY_GameMessageClassID {
	BYSettingID=10001,	// 设置相关保留位置S2C
	FrameC2SID=10100,	// 客户端帧数据
	FrameS2CID=10500,		// 帧数据
	UserRecoverRetID = 10501,	// 断线重连
	RobotID=10800,              //机器人消息
    BYLeaveTip = 10900,     //踢出房提示訊息
    BYNotifyFish = 10901,       //公告獲得魚
    BYGameStartC2SID = 10902,   //Client通知Server可以準備開始
    BYGameStartS2CID = 10903,   //Server通知Client可以開始
    BYGameCompetitionEndID = 10904,	//比賽結束
};
