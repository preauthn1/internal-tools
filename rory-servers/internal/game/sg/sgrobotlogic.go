package sg

import (
	"fmt"
	"math/rand"
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

type qiangzhuang struct {
	danpai      bool //单牌是否抢庄
	teshu       bool //特殊牌是否抢庄
	hunsangong  bool //混三公是否抢庄
	xiaosangong bool //小三公是否抢庄
	dasangong   bool //大三公是否抢庄
}

type bet struct {
	danpai      []int32 //单牌下注数组
	teshu       []int32 //特殊牌下注数组
	hunsangong  []int32 //混三公下注数组
	xiaosangong []int32 //小三公下注数组
	dasangong   []int32 //大三公下注数组
}

//机器人个性
type RobotCha struct {
	zhuangodd qiangzhuang
	betodd    bet
}

//创建机器个性对象
func NewRobotCha(cf *config.Config, cname string) *RobotCha {
	rcha := new(RobotCha)
	s1 := cf.GetConfigValue(cname, "danpai")
	strarr := strings.Split(s1, ";")
	rcha.zhuangodd.danpai = rcha.getBoolValue(strarr[0])
	rcha.betodd.danpai = rcha.getValue(strarr[1])

	s1 = cf.GetConfigValue(cname, "teshu")
	strarr = strings.Split(s1, ";")
	rcha.zhuangodd.teshu = rcha.getBoolValue(strarr[0])
	rcha.betodd.teshu = rcha.getValue(strarr[1])

	s1 = cf.GetConfigValue(cname, "hunsangong")
	strarr = strings.Split(s1, ";")
	rcha.zhuangodd.hunsangong = rcha.getBoolValue(strarr[0])
	rcha.betodd.hunsangong = rcha.getValue(strarr[1])

	s1 = cf.GetConfigValue(cname, "xiaosangong")
	strarr = strings.Split(s1, ";")
	rcha.zhuangodd.xiaosangong = rcha.getBoolValue(strarr[0])
	rcha.betodd.xiaosangong = rcha.getValue(strarr[1])

	s1 = cf.GetConfigValue(cname, "dasangong")
	strarr = strings.Split(s1, ";")
	rcha.zhuangodd.dasangong = rcha.getBoolValue(strarr[0])
	rcha.betodd.dasangong = rcha.getValue(strarr[1])

	return rcha
}

func (rc *RobotCha) getBoolValue(str string) bool {
	if str == "1" {
		return true
	}
	return false
}

func (rc *RobotCha) getValue(str string) []int32 {
	strarr := strings.Split(str, ",")
	arr := make([]int32, 0)
	for i := 0; i < len(strarr); i++ {
		n, _ := strconv.ParseInt(strarr[i], 10, 32)
		arr = append(arr, int32(n))
	}

	return arr
}

//获取是否抢庄
func (rc *RobotCha) getQiangZhuangOdd(t int32) bool {
	if t == TeShu {
		return rc.zhuangodd.teshu
	}

	if t == HunSanGong {
		return rc.zhuangodd.hunsangong
	}

	if t == XiaoSanGong {
		return rc.zhuangodd.xiaosangong
	}

	if t == DaSanGong {
		return rc.zhuangodd.dasangong
	}

	return rc.zhuangodd.danpai
}

//获取下注倍数
func (rc *RobotCha) getBetOdd(t int32) int32 {
	if t == TeShu {
		cn := len(rc.betodd.teshu)
		return rc.betodd.teshu[rand.Intn(cn)]
	}

	if t == HunSanGong {
		cn := len(rc.betodd.hunsangong)
		return rc.betodd.hunsangong[rand.Intn(cn)]
	}

	if t == XiaoSanGong {
		cn := len(rc.betodd.xiaosangong)
		return rc.betodd.xiaosangong[rand.Intn(cn)]
	}

	if t == DaSanGong {
		cn := len(rc.betodd.dasangong)
		return rc.betodd.dasangong[rand.Intn(cn)]
	}

	cn := len(rc.betodd.danpai)
	return rc.betodd.danpai[rand.Intn(cn)]
}

//机器人逻辑
type SGRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

func NewSGRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *SGRobotLogic {
	rblgc := new(SGRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)

	return rblgc
}

//解析机器人个性
func (lgc *SGRobotLogic) parseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)

		cha := NewRobotCha(robotCfg, cname)
		lgc.robotChas = append(lgc.robotChas, cha)

	}
}

//获取机器人个性
func (lgc *SGRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *SGRobotLogic) getGameLogic() *SGGameLogic {
	return lgc.GameLogic.(*SGGameLogic)
}

func (lgc *SGRobotLogic) getPoker() *poker.FixedPoker {
	return lgc.getGameLogic().poker.(*poker.FixedPoker)
}

//游戏消息
func (lgc *SGRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	//底注阶段,收到底注倍数通知时
	glc := lgc.getGameLogic()
	if glc.state == QiangZhuang && m.SClassID == int32(netproto.SG_GameMessageClassID_StateID) {
		cha := lgc.getRobotCha(robot.UserInfo.UserID)
		rcards := lgc.getPoker().GetUserCards(robot.UserInfo.UserID)
		if rcards == nil || len(rcards.Cards) == 0 {
			mlog.Debug("没有牌？")
			return
		}
		_, t := glc.getCardsOdd(rcards)
		qiang := cha.getQiangZhuangOdd(t)
		lgc.SendGameMessage(robot.UserInfo, lgc.getQiangZhuangMsg(qiang))

		return
	}

	if glc.state == Bet && m.SClassID == int32(netproto.SG_GameMessageClassID_StateID) {
		cha := lgc.getRobotCha(robot.UserInfo.UserID)
		rcards := lgc.getPoker().GetUserCards(robot.UserInfo.UserID)
		if rcards == nil || len(rcards.Cards) == 0 {
			mlog.Debug("没有牌？")
			return
		}
		_, t := glc.getCardsOdd(rcards)
		odd := cha.getBetOdd(t)
		lgc.SendGameMessage(robot.UserInfo, lgc.getBetMsg(odd))
		return
	}
}

func (lgc *SGRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

//抢庄消息
func (lgc *SGRobotLogic) getQiangZhuangMsg(odd bool) *msg.Message {
	mlog.Debug("机器人抢庄")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SG_GameMessageClassID_SGQiangZhuangID)

	btdata := new(netproto.SG_QiangZhuang)
	btdata.QiangZhuang = proto.Bool(odd)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//下注消息
func (lgc *SGRobotLogic) getBetMsg(odd int32) *msg.Message {
	mlog.Debug("机器人选择倍数")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SG_GameMessageClassID_BetID)

	btdata := new(netproto.SG_Bet)
	btdata.Bet = proto.Int32(odd)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}
