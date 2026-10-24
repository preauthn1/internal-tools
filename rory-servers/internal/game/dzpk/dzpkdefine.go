//==================================================================
//文件名： dzpkdefine.go
//说明：   德州扑克游戏常量定义
//修改记录：
//==================================================================
package dzpk

//游戏阶段定义
type gameProgress int32

const (
	GAMEPROGRESS_NONE    gameProgress = 1
	GAMEPROGRESS_BEGIN                = 2 //游戏开始,此阶段会初始化开局自动押注
	GAMEPROGRESS_DEAL                 = 3 //发牌阶段，此阶段发两张手牌
	GAMEPROGRESS_PREFLOP              = 4 //底牌圈，发两张私牌后开始押注
	GAMEPROGRESS_FLOP                 = 5 //翻牌圈，发3张公牌后开始押注
	GAMEPROGRESS_TURN                 = 6 //转牌圈，发4张公牌后开始押注
	GAMEPROGRESS_RIVER                = 7 //河牌圈，发5张公牌后开始押注
	GAMEPROGRESS_END                  = 8 //结束结算
)

// 定义游戏阶段名称
var GState []string = []string{0: "空", 1: "游戏开始", 2: "自动押注", 3: "发两张手牌", 4: "底牌圈", 5: "翻牌圈", 6: "转牌圈", 7: "河牌圈", 8: "结束结算"}

// 获取游戏阶段名称
func getState(index gameProgress) string {
	length := int32(len(GState))
	if int32(index)+1 <= length {
		return GState[index]
	}
	return "未知状态"
}

//玩家位置角色
type playerRole int32

const (
	ROLE_NONE        playerRole = iota //空
	ROLE_BUTTON                        //庄家
	ROLE_SMALLBLIND                    //小盲
	ROLE_BIGBLIND                      //大盲
	ROLE_UNDERTHEGUN                   //枪口
	ROLE_CUTOFF                        //关位
	ROLE_BUTTONSMALL                   //玩家是小盲和庄家
)

//玩法类型
type playType int32

const (
	GREEN  playType = 1 //新手场
	MIDDLE          = 2 //中级场
	HIGH            = 3 //高级场
	RICH            = 4 //富豪场
)

// 定义操作类型
type ACTIONTYPE int32

const (
	NONE    ACTIONTYPE = iota //空
	BET                       //1
	CALL                      //2
	RAISE                     //3
	CHECK                     //4
	ALLIN                     //5
	FOLD                      //6
	CALLANY                   //7跟任何注
)

// 定义错误类型
const (
	ERROR_OPRATE        = -1
	ERROR_CLASSID       = -2
	ERROR_TESTDATAERROR = -3  //用户GM命令有误
	ERROR_NOMONEY       = -10 //金币不足
)

const (
	ACTION_NAME_GAMEEND = "GAME_END"
)

// 定义定时器名称
const opttimerstr = ""
const opttimerWait = ""
const opttimerNextGmae = ""

// 定义机器人个性配置索引
const (
	INDEXP_FOLD       = 1
	INDEXP_CALL       = 2
	INDEXP_RAISE      = 3
	INDEXP_EASYWEIGHT = 4
	INDEXP_ALLIN      = 5
	INDEXP_FREEWEIGHT = 6
	INDEXP_MINMAX     = 7

	INDEXA_FOLD       = 1
	INDEXA_CHECK      = 8
	INDEXA_RAISE      = 9
	INDEXA_EASYWEIGHT = 10
	INDEXA_ALLIN      = 11
	INDEXA_FREEWEIGHT = 12
	INDEXA_MINMAX     = 13
)
