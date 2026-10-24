package sss

import (
	proto "example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
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

//机器人个性
type RobotCha struct {
}

//创建机器个性对象
func NewRobotCha() *RobotCha {
	rcha := new(RobotCha)
	return rcha
}

//机器人逻辑
type SSSRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

func NewSSSRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *SSSRobotLogic {
	rblgc := new(SSSRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)

	return rblgc
}

//解析机器人个性
func (lgc *SSSRobotLogic) parseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cha := NewRobotCha()
		lgc.robotChas = append(lgc.robotChas, cha)
	}
}

//获取机器人个性
func (lgc *SSSRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *SSSRobotLogic) getGameLogic() *SSSTable {
	return lgc.GameLogic.(*SSSTable)
}

func (lgc *SSSRobotLogic) getPoker() *poker.FixedPoker {
	return lgc.getGameLogic().poker.(*poker.FixedPoker)
}

//游戏消息
func (lgc *SSSRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	//底注阶段,收到底注倍数通知时
	glc := lgc.getGameLogic()
	if glc.step == netproto.SSSGameStep_SSSGameStepCardTypeChoose &&
		m.SClassID == int32(netproto.SSS_GameMessageClassID_NotifyOptTimeOutId) {
		upk := glc.gamePoker.GetUserPoker(robot.UserInfo.UserID)
		t := glc.getSpecialType(upk.Cards)
		special := false
		act := NewAutoCardType()
		glc.getNormalType(upk.Cards, act)
		if t != 0 {
			special = true
		}

		lgc.SendGameMessageNow(robot.UserInfo, lgc.getSpecilTypeMsg(t, special, act))
	}

	if glc.step == netproto.SSSGameStep_SSSGameStepCombine &&
		m.SClassID == int32(netproto.SSS_GameMessageClassID_NotifyOptTimeOutId) {
		upk := glc.gamePoker.GetUserPoker(robot.UserInfo.UserID)
		t := glc.getSpecialType(upk.Cards)
		if t != 0 {
			return
		}
		act := NewAutoCardType()

		glc.getNormalType(upk.Cards, act)
		lgc.SendGameMessage(robot.UserInfo, lgc.getCardCombMsg(act))
	}
}

func (lgc *SSSRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

func (lgc *SSSRobotLogic) getSpecilTypeMsg(t int32, special bool, act *AutoCardType) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SSS_GameMessageClassID_ChooseSpecialCardTypeId)

	btdata := new(netproto.ChooseSpecialCardType)
	btdata.SpecialType = proto.Int32(t)
	btdata.Special = proto.Bool(special)
	if !special {
		btdata.Hcs = act.recommandCt
	}

	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

func (lgc *SSSRobotLogic) getCardCombMsg(act *AutoCardType) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SSS_GameMessageClassID_CompleteCardCombId)

	btdata := new(netproto.CompleteCardComb)
	btdata.Hcs = act.recommandCt
	btdata.IsSpecial = proto.Bool(act.isSpecial)
	btdata.St = proto.Int32(act.sct)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}
