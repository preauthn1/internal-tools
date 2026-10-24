
enum HHNN_GameMessageClassID

{
    HHNNSettingID = 10001,  //设置相关 S2C
    HHNNStateID = 10002,        //阶段ID S2C
    HHNNBetID = 10003,          //下注ID C2S
    HHNNBroadCastBetID = 10004,//广播下注 S2C
    HHNNGameResultPokerID = 10005,      //游戏结果扑克 S2C
    HHNNGameSceneID = 10006,    //游戏场景 断线重连
    HHNNGameResultWinID = 10007,        //游戏结果赢钱 S2C
    HHNNBetInfoID = 10008,//用户断线重连已下注的消息
    HHNNBetRetID = 10009,//下注结果 S2C
    HHNNShangZhuangID = 10010,//上庄请求 C2S
    HHNNShangZhuangRetID = 10011,//上庄请求结果 S2C
    HHNNGetShangZhuangListID = 10012,//获取上庄列表 C2S
    HHNNShangZhuangListInfoID = 10013,//下发上庄列表 S2C
    HHNNZhuangInfoID = 10014,//庄家信息 S2C
    HHNNBetClearID = 10015,//玩家清除已下的注 C2S
    HHNNResultCardListID = 10016,//前20次开牌结果
    HHNNBetClearRetID = 10017,//清理下注结果
    HHNNRichestListID = 10018,//大富豪列表的消息
    HHNNLuckyStarListID = 10019,//明燈玩家的消息
    HHNNRoundID = 10020,//本局局數
    HHNNBroadCastFapaiThreeCard = 10021,	//第一手兩張牌
}

