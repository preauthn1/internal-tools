package sgj

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/util"

	proto "example.invalid/golang/protobuf/proto"
)

const (
	RobotLogicConfigFileName = "robot"
)

var baseScore int32

//机器人个性
type RobotCha struct {
	minbetcount int32   //最少下多少注
	maxbetcount int32   //最多一次下注多少
	bettype     []int32 //下那些注
}

//创建机器个性对象
func NewRobotCha(cf *config.Config, cname string) *RobotCha {
	rcha := new(RobotCha)
	s1 := cf.GetConfigValue(cname, "betcount")
	strarr := strings.Split(s1, ",")
	min, _ := strconv.ParseInt(strarr[0], 10, 32)
	rcha.minbetcount = int32(min)
	max, _ := strconv.ParseInt(strarr[1], 10, 32)
	rcha.maxbetcount = int32(max)
	s1 = cf.GetConfigValue(cname, "bettype")
	rcha.bettype = rcha.getArr(s1)
	rand.Seed(time.Now().UnixNano())

	return rcha
}

func (rc *RobotCha) getArr(str string) []int32 {
	arr := make([]int32, 0)
	strarr := strings.Split(str, ",")
	for i := 0; i < len(strarr); i++ {
		m, _ := strconv.ParseInt(strarr[i], 10, 32)
		arr = append(arr, int32(m))
	}

	return arr
}

func (rc *RobotCha) getBetType() int32 {
	index := len(rc.bettype)
	n := rand.Int31n(int32(index))
	return rc.bettype[n]
}

func (rc *RobotCha) getBetCount() int32 {
	i := rc.maxbetcount - rc.minbetcount + 1
	n := rand.Int31n(i) + rc.minbetcount
	return n * baseScore
}

//机器人逻辑
type SGJRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
	startbet  bool
}

func NewSGJRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *SGJRobotLogic {
	rblgc := new(SGJRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)
	rblgc.startbet = false

	return rblgc
}

//解析机器人个性
func (lgc *SGJRobotLogic) parseRobotCha(cfgs *config.Group) {
	rand.Seed(time.Now().UnixNano())
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)
		cha := NewRobotCha(robotCfg, cname)
		lgc.robotChas = append(lgc.robotChas, cha)
	}
	glc := lgc.getGameLogic()
	baseScore = glc.betcount[0]
}

//获取机器人个性
func (lgc *SGJRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *SGJRobotLogic) getGameLogic() *ShuiGuoJiLogic {
	return lgc.GameLogic.(*ShuiGuoJiLogic)
}

//游戏消息
func (lgc *SGJRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	glc := lgc.getGameLogic()
	if glc.zhuangUserID == robot.UserInfo.UserID {
		return
	}

	if m.SClassID == int32(netproto.SGJ_GameMessageClassID_SGJStateID) &&
		glc.gameState == betState {
		lgc.robotBet(robot.UserInfo)
		return
	}

	if m.SClassID == int32(netproto.SGJ_GameMessageClassID_SGJCannotBetID) &&
		glc.gameState == betState {
		lgc.robotBet(robot.UserInfo)
		return
	}
}

func (lgc *SGJRobotLogic) robotBet(ur *gamelogic.GameUser) {
	rcha := lgc.getRobotCha(ur.UserID)
	lgc.SendGameMessage(ur, lgc.getBetMsg(rcha.getBetType(), rcha.getBetCount()))
}

func (lgc *SGJRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

func (lgc *SGJRobotLogic) getBetMsg(bettype int32, betcount int32) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SGJ_GameMessageClassID_SGJBetID)
	btdata := new(netproto.SGJ_BetTypeGroup)
	betdata := new(netproto.SGJ_BetType)
	betdata.BetType = proto.Int32(bettype)
	betdata.BetCount = proto.Int32(betcount)
	btdata.BetGroup = append(btdata.BetGroup, betdata)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}
