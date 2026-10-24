//==================================================================
//文件名： ttzdefine.go
//说明：   推筒子游戏常量定义
//修改记录：
//==================================================================
package ttz

//游戏阶段定义
const (
	GAMEPROGRESS_NONE         int32 = iota //空
	GAMEPROGRESS_DICE                      //丢骰子1
	GAMEPROGRESS_DEALONE                   //发牌一张2
	GAMEPROGRESS_ROBBUTTON                 //抢庄3
	GAMEPROGRESS_ROBANIMATION              //抢庄动画4
	GAMEPROGRESS_CHOOSEBET                 //闲家选择底注5
	GAMEPROGRESS_LIGHT                     //亮牌6
	GAMEPROGRESS_SETTLE                    //结算7
)

//游戏阶段名称
var GState []string = []string{0: "空", 1: "丢骰子", 2: "第一次发牌", 3: "抢庄", 4: "抢庄动画", 5: "下注", 6: "亮牌", 7: "结算"}

//操作动作定义
const (
	ACTION_NONE int32 = iota
	ACTION_ROB
	ACTION_BET
	ACTION_LIGHT
)

//错误信息定义
const (
	ERROR_USERNOTEXITST = -1 //用户不存在
	ERROR_REPEATOP      = -2 //重复操作
	ERROR_ACTIONTYPE    = -3 //错误操作类型
	ERROR_CLASSID       = -4 //错误class id
	ERROR_ILL           = -5 //非法操作
)

//定时器名称
const (
	TIME_NAME_TIMEOUTACTION = ""
)

//结算信息定义
type settleMentInfo struct {
	SeatNo      int32
	Score       int64
	ButtonScore int64
}

//获取游戏阶段名称
func getState(index int32) string {
	length := int32(len(GState))
	if index+1 <= length {
		return GState[index]
	}
	return "未知状态"
}
