//==================================================================
//文件名： Game_SCMJ.proto
//作者：   dyb
//日期：   2017/12/22
//说明：   麻将游戏客户服务器通信协议
//修改记录：
//==================================================================

/*
// Card int32=100*Color+10*Number+index(第几副牌)
// color=card/100
// number=card%100/10
const (
    COLOR_TONG  int32 = 1 //筒子
    COLOR_TIAO  int32 = 2 //条子
    COLOR_WAN   int32 = 3 //万子
    COLOR_FENG  int32 = 4 //风牌
    COLOR_JIAN  int32 = 5 //箭牌
    COLOR_HUA   int32 = 6 //花牌
    COLOR_BAIDA int32 = 7 //百搭牌
)

const (
    //风牌
    NUMBER_EAST  int32 = 1 //东
    NUMBER_WEST  int32 = 2 //西
    NUMBER_SOUTH int32 = 3 //南
    NUMBER_NORTH int32 = 4 //北
    //箭牌
    NUMBER_ZHONG int32 = 1 //中
    NUMBER_FA    int32 = 2 //发
    NUMBER_BAI   int32 = 3 //白
    //花牌
    NUMBER_CHUN int32 = 1 //春
    NUMBER_XIA  int32 = 2 //夏
    NUMBER_QIU  int32 = 3 //秋
    NUMBER_DONG int32 = 4 //冬
    NUMBER_MEI  int32 = 5 //梅
    NUMBER_LAN  int32 = 6 //兰
    NUMBER_ZHU  int32 = 7 //竹
    NUMBER_JU   int32 = 8 //菊
    //序牌1~9
)

//游戏阶段
int32 = 1 //开始
int32 = 2 //丢色子
int32 = 3 //发牌
int32 = 4 //定缺
int32 = 5 //战斗
int32 = 6 //查大叫
int32 = 7 //查花猪
int32 = 8 //结算
int32 = 9 //结束

*/
enum ERMJ_GameMessageClassID {
    SCMJErrorId = 10001,                          //错误信息，SCMJ_Error
    SCMJThrowDiceId = 10002,                      //丢色子 MJ_Dice
    SCMJDealId = 10003,                           //发牌[摸牌] SCMJ_BroadcastDeal

    SCMJBroadCastChooseId = 10005,                //选定缺返回，SCMJ_BroadcastAction
    SCMJBroadCastChooseResultId = 10006,          //所有玩家定缺信息,SCMJ_BroadcastChooseResult

    SCNJActionId = 10007,                         //操作[定缺,出牌,吃碰杠胡],SCMJ_Action
    SCMJBroadCastActionId = 10008,                //操作返回，SCMJ_BroadcastAction

    SCMJScenceExchangeId = 10009,                 //吃碰杠胡通知,SCMJ_BroadCastUserOper

    SCMJWaitId = 10010,                           //当前操作用户，SCMJ_Wait

    SCMJSettleMentIntimeId = 10011,               //即时小结算，SCMJ_InTimeSettle
    SCMJTotalSettleMentId = 10012,                //大结算，SCMJ_TotalSettleMent

    SCMJGameSenceId = 10013,                      //场景恢复,SCMJ_GameScence
    SCMJGameSettingId = 10014,                    //游戏设置,Mj_GameSetting
    SCMJCanOpEventId = 10015,                     //提示操作,SCMJ_CanOpEvent

    S2C_SCMJLightId=10016,                        //亮牌,S2C_SCMJ_Light
    C2S_SCMJRemindId=10017,                       //客户端请求提示信息，C2S_SCMJ_Remind
    S2C_SCMJRemindId=10018,                       //服务端返回提示信息，S2C_SCMJ_Remind

    S2C_TuoGuanId=10019,                            //托管，S2C_SCMJ_TuoGuan
    C2S_TuoGuanId=10020,                            //托管，C2S_SCMJ_TuoGuan
    
};