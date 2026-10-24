//==================================================================
//文件名： Game_TTZ.proto
//作者：   dyb
//日期：   2017/10/30
//说明：   推筒子游戏客户服务器通信协议
//修改记录：
//==================================================================


enum TTZ_GameMessageClassID {

    ActionOprateID = 10001,                 //客户端操作，TTZ_ActionInfo
    BroadCastActionResultID = 10002,        //广播操作结果，TTZ_BroadCastActionResult
    BroadCastRotResultID = 10003,           //广播抢庄结果，TTZ_BroadCastRobResult
    BroadCastDealInfoID = 10004,            //广播发牌，TTZ_BroadCastDealInfo
    BroadCastUpdateCTID = 10005,            //广播牌型，TTZ_BroadCastUpdateCT
    BroadCastLightCardsID = 10006,          //广播亮牌信息,TTZ_BroadCastLightCards
    BroadCastSettleMentRound = 10007,       //广播小局结算,TTZ_BrodCastSettleMentRound
    BroadCastSettleMentTotal = 10008,       //广播总结算,TTZ_BrodCastSettleMentTotal
    BroadCastGameState = 10009,             //广播游戏状态,TTZ_BroadCastGameState
    BroadCastTTZRecoverID=10010,            //服务器推送：场景恢复广播,TTZ_GameSence
    BroadCastTTZGameSettingsID=10011,       //游戏设置，TTZ_GameSettings
    BroadCastTTZMaxOpRateID=10012,             //广播可以抢庄的最大倍数,TTZ_MaxOpRate

    BroadCastTtzErrorID = 10028,        //服务器推送：操作失败TTZ_ErrorMsg
    ActionGMID=10029,                    //GM，TTZ_GMinfo
};