package scmj

const (
	PROGRESS_NONE   int32 = iota //0
	PROGRESS_BEGIN               //开始动画
	PROGRESS_DICE                //丢色子
	PROGRESS_DEAL                //发牌
	PROGRESS_CHOOSE              //定缺
	PROGRESS_FIGHT               //战斗
	PROGRESS_LIGHT               //亮牌
	PROGRESS_PIG                 //花猪
	PROGRESS_DAJIAO              //大叫
	PROGRESS_SETTLE              //结算
	PROGRESS_END                 //结束
)

const (
	OP_NONE          int32 = 0
	OP_CHOOSEDINGQUE       = 1
	OP_OUTCARD             = 2
	OP_PENG                = 3
	OP_GANG                = 4
	OP_MINGGANG            = 40 //碰了牌以后，自己摸上的第4张牌，可以明杠(补杠)
	OP_YINGANG             = 41 //自己手里有一坎牌，其他玩家打出了第4张牌(明杠)
	OP_ANGANG              = 42 //自己手里摸到了4张一样的牌，可以暗杠(暗杠)
	OP_HU                  = 7
	OP_CHI                 = 8
	OP_PASS                = 9
	OP_CHOOSERATE          = 10 //选倍率
	OP_HU_DIANPAO          = 70 //点炮胡
	OP_HU_QIANGGANG        = 71 //抢杠胡
	OP_HU_ZIMO             = 72 //自模胡
	OP_HU_GANG             = 73 //杠上花
	OP_HU_GANGPAO          = 74 //杠上炮

)

const (
	ANIMATION_BEGIN    = "BEGIN"
	ANIMATION_DICE     = "DICE"
	ANIMATION_DEAL     = "DEAL"
	ANIMATION_END      = "END"
	ANIMATION_HUAZHU   = "HUAZHU"
	ANIMATION_DAJIAO   = "DAJIAO"
	ANIMATION_FIGHTEND = "FIGHTEND"
	ANIMATION_LIGHT    = "LIGHT"
	ANIMATION_HUPAI    = "HUPAI"
)

const (
	TYPE_SETTLEMENT_NONE       int32 = iota //
	TYPE_SETTLEMENT_GF                      //刮风
	TYPE_SETTLEMENT_XY                      //下雨
	TYPE_SETTLEMENT_HZ                      //花猪
	TYPE_SETTLEMENT_DJ                      //大叫
	TYPE_SETTLEMENT_HU                      //胡牌5
	TYPE_SETTLEMENT_HU_ZIMO                 //自摸6
	TYPE_SETTLEMENT_HU_DIANPAO              //点炮7
	TYPE_SETTLEMENT_HU_QG                   //抢杠8
	TYPE_SETTLEMENT_BACK                    //退杠9
	TYPE_SETTLEMENT_HJZY                    //呼叫转移10
)

const (
	ERROR_STATE int32 = -1 //游戏状态错误，非法操作
)
