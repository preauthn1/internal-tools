package qznn

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/util"

	proto "example.invalid/golang/protobuf/proto"
)

const (
	RobotLogicConfigFileName = "robot"
)

//机器人个性
type RobotCha struct {
	Size_0_ZhuangBetTimes   int32
	Size_0_XianBetTimes     int32
	Size_1_9_ZhuangBetTimes int32
	Size_1_9_XianBetTimes   int32
	Size_10_ZhuangBetTimes  int32
	Size_10_XianBetTimes    int32
}

//创建机器人个性对象
func NewRobotCha(size0_Zhuang int32, size0_Xian int32, size1_9_Zhuang int32, size1_9_Xian int32, size10_Zhuang int32, size10_Xian int32) *RobotCha {
	cha := new(RobotCha)
	cha.Size_0_ZhuangBetTimes = size0_Zhuang
	cha.Size_0_XianBetTimes = size0_Xian
	cha.Size_1_9_ZhuangBetTimes = size1_9_Zhuang
	cha.Size_1_9_XianBetTimes = size1_9_Xian
	cha.Size_10_ZhuangBetTimes = size10_Zhuang
	cha.Size_10_XianBetTimes = size10_Xian

	return cha
}

//抢庄牛牛机器人逻辑
type QZNNRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

//创建抢庄牛牛机器人逻辑
func NewQZNNRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *QZNNRobotLogic {
	rblgc := new(QZNNRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)

	return rblgc
}

//解析机器人个性
func (lgc *QZNNRobotLogic) parseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)

		s1 := robotCfg.GetConfigValue(cname, "size0")
		s2 := robotCfg.GetConfigValue(cname, "size1_9")
		s3 := robotCfg.GetConfigValue(cname, "size10")

		s1Arr := strings.Split(s1, ",")
		s1_zhuang, _ := strconv.Atoi(s1Arr[0])
		s1_xian, _ := strconv.Atoi(s1Arr[1])

		s2Arr := strings.Split(s2, ",")
		s2_zhuang, _ := strconv.Atoi(s2Arr[0])
		s2_xian, _ := strconv.Atoi(s2Arr[1])

		s3Arr := strings.Split(s3, ",")
		s3_zhuang, _ := strconv.Atoi(s3Arr[0])
		s3_xian, _ := strconv.Atoi(s3Arr[1])

		cha := NewRobotCha(int32(s1_zhuang), int32(s1_xian), int32(s2_zhuang), int32(s2_xian), int32(s3_zhuang), int32(s3_xian))
		lgc.robotChas = append(lgc.robotChas, cha)

	}
}

//获取机器人个性
func (lgc *QZNNRobotLogic) getRobotCha(userid int32) *RobotCha {
	var cha *RobotCha
	if cha, ok := lgc.gameChas[userid]; ok {
		return cha
	}

	clen := len(lgc.robotChas)
	cn := util.RndInt(0, clen)
	cha = lgc.robotChas[cn]

	lgc.gameChas[userid] = cha

	return cha
}

func (lgc *QZNNRobotLogic) getGameLogic() *QZNNGameLogic {
	return lgc.GameLogic.(*QZNNGameLogic)
}

func (lgc *QZNNRobotLogic) getPoker() *poker.FixedPoker {
	return lgc.getGameLogic().poker.(*poker.FixedPoker)
}

//获取用户牌型
func (lgc *QZNNRobotLogic) getUserCards(userid int32) *poker.CardGroup {
	cards := poker.NewCardGroup(nil)

	uexcards := lgc.getGameLogic().gamePoker.GetUserPoker(userid).Cards
	ucards := lgc.getPoker().GetUserCards(userid)

	cards.Cards = append(cards.Cards, uexcards.Cards...)
	cards.Cards = append(cards.Cards, ucards.Cards...)

	return cards
}

//游戏消息
func (lgc *QZNNRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

	//底注阶段,收到底注倍数通知时
	if lgc.getGameLogic().gameStep == gameStep_Begin && m.SClassID == int32(netproto.QZNN_GameMessageClassID_GameStepID) {
		//获取用户牌型
		ucards := lgc.getUserCards(robot.UserInfo.UserID)

		//解析用户牌型
		nnret := getNiuNiuRet(ucards)

		mlog.Debug("nret:%s", nnret)

		//获取机器人个性
		cha := lgc.getRobotCha(robot.UserInfo.UserID)

		var bt int32 = 1
		if nnret.niuRet == nnCard_WuNiu {
			bt = lgc.getZhuangBetBtn(robot.UserInfo, cha.Size_0_ZhuangBetTimes)
		} else if nnret.niuRet == nnCard_Niu {
			bt = lgc.getZhuangBetBtn(robot.UserInfo, cha.Size_1_9_ZhuangBetTimes)
		} else {
			bt = lgc.getZhuangBetBtn(robot.UserInfo, cha.Size_10_ZhuangBetTimes)
		}

		btmsg := new(msg.Message)
		btmsg.BClassID = int32(netproto.MessageBClassID_Game)
		btmsg.SClassID = int32(netproto.QZNN_GameMessageClassID_QiangZhuangID)

		btdata := new(netproto.QZNN_QiangZhuang)
		btdata.Times = proto.Int32(bt)
		btmsg.MsgData, _ = proto.Marshal(btdata)

		lgc.SendGameMessage(robot.UserInfo, btmsg)

		mlog.Debug("机器人%v选择抢庄倍数%d", robot.UserInfo, bt)
	} else if lgc.getGameLogic().gameStep == gameStep_ZhuangConfirm && m.SClassID == int32(netproto.QZNN_GameMessageClassID_GameStepID) {
		//如果机器人是庄家
		if lgc.getGameLogic().zhuangInfo.userID == robot.UserInfo.UserID {
			return
		}

		//获取用户牌型
		ucards := lgc.getUserCards(robot.UserInfo.UserID)

		//解析用户牌型
		nnret := getNiuNiuRet(ucards)

		//获取一个机器人个性
		cha := lgc.getRobotCha(robot.UserInfo.UserID)

		var bt int32 = 1
		if nnret.niuRet == nnCard_WuNiu {
			bt = cha.Size_0_XianBetTimes
		} else if nnret.niuRet == nnCard_Niu {
			bt = cha.Size_1_9_XianBetTimes
		} else {
			bt = cha.Size_10_XianBetTimes
		}

		//闲家倍数阶段
		cgmsg := new(msg.Message)
		cgmsg.BClassID = int32(netproto.MessageBClassID_Game)
		cgmsg.SClassID = int32(netproto.QZNN_GameMessageClassID_XianjiaBetTimesID)
		xjbt := new(netproto.QZNN_XianJiaBetTimes)
		xjbt.Times = proto.Int64(int64(bt))
		cgmsg.MsgData, _ = proto.Marshal(xjbt)

		lgc.SendGameMessage(robot.UserInfo, cgmsg)

		mlog.Debug("机器人%v闲家倍数消息", robot.UserInfo)

	} else if lgc.getGameLogic().gameStep == gameStep_CardGroup && m.SClassID == int32(netproto.QZNN_GameMessageClassID_GameStepID) {
		//组牌阶段
		cgmsg := new(msg.Message)
		cgmsg.BClassID = int32(netproto.MessageBClassID_Game)
		cgmsg.SClassID = int32(netproto.QZNN_GameMessageClassID_CardGroupFinishID)
		cgmsg.MsgData = make([]byte, 0)

		lgc.SendGameMessageNow(robot.UserInfo, cgmsg)

		mlog.Debug("机器人%v提交组牌消息", robot.UserInfo)

	}
}

func (lgc *QZNNRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

//玩家金額/玩家最低注金額/特殊牌倍數(最高3)/閒家人數(4) 之後取整數 就是最高能搶倍數
func (qznn *QZNNRobotLogic) getZhuangBetBtn(user *gamelogic.GameUser, betTimes int32) int32 {
	maxBetOption := int32(int32(user.Score) / int32(getNNConfig().xianBetTimes[0]) / 3 / 4)

	if maxBetOption > betTimes {
		return betTimes
	} else {
		return 1
	}

	return 0
}
