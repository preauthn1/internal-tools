package tbnn

import (
	"fmt"

	proto "example.invalid/golang/protobuf/proto"
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
)

const (
	RobotLogicConfigFileName = "robot"
)

// RobotCha 機器人個性
type RobotCha struct {
	Size_0_BetTimes   int32
	Size_1_9_BetTimes int32
	Size_10_BetTimes  int32
}

//创建机器个性对象
func NewRobotCha(size0 int32, size1_9 int32, size_10 int32) *RobotCha {
	rcha := new(RobotCha)
	rcha.Size_0_BetTimes = size0
	rcha.Size_1_9_BetTimes = size1_9
	rcha.Size_10_BetTimes = size_10

	return rcha
}

//机器人逻辑
type TBNNRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

func NewTBNNRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *TBNNRobotLogic {
	rblgc := new(TBNNRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)

	return rblgc
}

//解析机器人个性
func (lgc *TBNNRobotLogic) parseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)

		s1 := robotCfg.GetConfigInt32(cname, "size0")
		s2 := robotCfg.GetConfigInt32(cname, "size1_9")
		s3 := robotCfg.GetConfigInt32(cname, "size10")

		cha := NewRobotCha(s1, s2, s3)
		lgc.robotChas = append(lgc.robotChas, cha)

	}
}

//获取机器人个性
func (lgc *TBNNRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *TBNNRobotLogic) getGameLogic() *TBNNGameLogic {
	return lgc.GameLogic.(*TBNNGameLogic)
}

func (lgc *TBNNRobotLogic) getPoker() *poker.FixedPoker {
	return lgc.getGameLogic().poker.(*poker.FixedPoker)
}

//获取用户牌型
func (lgc *TBNNRobotLogic) getUserCards(userid int32) *poker.CardGroup {
	cards := poker.NewCardGroup(nil)

	uexpoker := lgc.getGameLogic().gamePoker.GetUserPoker(userid)
	if uexpoker != nil {
		uexcards := uexpoker.Cards
		cards.Cards = append(cards.Cards, uexcards.Cards...)
	}
	ucards := lgc.getPoker().GetUserCards(userid)

	cards.Cards = append(cards.Cards, ucards.Cards...)

	return cards
}

//游戏消息
func (lgc *TBNNRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	//底注阶段,收到底注倍数通知时
	if lgc.getGameLogic().gameStep == gameStep_BaseScore && m.SClassID == int32(netproto.TBNN_GameMessageClassID_TBNN_GameStepID) {
		//获取用户牌型
		ucards := lgc.getUserCards(robot.UserInfo.UserID)

		//解析用户牌型
		nnret := getNiuNiuRet(ucards)

		//获取机器人个性
		cha := lgc.getRobotCha(robot.UserInfo.UserID)

		var bt int32 = 1
		if nnret.niuRet == nnCard_WuNiu {
			bt = cha.Size_0_BetTimes
		} else if nnret.niuRet == nnCard_Niu {
			bt = cha.Size_1_9_BetTimes
		} else {
			bt = cha.Size_10_BetTimes
		}

		btmsg := new(msg.Message)
		btmsg.BClassID = int32(netproto.MessageBClassID_Game)
		btmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_UserBaseScoreID)

		btdata := new(netproto.TBNN_UserBaseScore)
		btdata.BaseScore = proto.Int32(bt)
		btmsg.MsgData, _ = proto.Marshal(btdata)

		lgc.SendGameMessage(robot.UserInfo, btmsg)

		mlog.Debug("机器人%v选择底注倍数%d", robot.UserInfo, bt)
	} else if lgc.getGameLogic().gameStep == gameStep_CardGroup && m.SClassID == int32(netproto.TBNN_GameMessageClassID_TBNN_GameStepID) {
		//组牌阶段
		cgmsg := new(msg.Message)
		cgmsg.BClassID = int32(netproto.MessageBClassID_Game)
		cgmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_UserCardGroupID)
		cgmsg.MsgData = make([]byte, 0)

		lgc.SendGameMessage(robot.UserInfo, cgmsg)

		mlog.Debug("机器人%v提交组牌消息", robot.UserInfo)

	}
}

func (lgc *TBNNRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}
