package hhnn

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
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

//机器人个性
type RobotCha struct {
	minbetcount int32
	maxbetcount int32
	bettype     []int32
	line        int32
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

	rcha.line = cf.GetConfigInt32(cname, "line")

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

//cha为正时表示龙多，为负时表示虎多
func (rc *RobotCha) getBetType(cha int32) int32 {
	v := cha
	if cha < 0 {
		v = -cha
	}

	if v > rc.line {
		if cha < 0 {
			//mlog.Debug("平衡龙")
			return 1
		} else {
			//mlog.Debug("平衡虎")
			return 3
		}
	}

	index := len(rc.bettype)
	n := rand.Int31n(int32(index))
	return rc.bettype[n]
}

func (rc *RobotCha) getBetCount() int32 {
	i := rc.maxbetcount - rc.minbetcount + 1
	n := rand.Int31n(i) + rc.minbetcount
	return n * 100
}

//机器人逻辑
type HHNNRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

//新建一个机器人逻辑实例
func NewHHNNRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *HHNNRobotLogic {
	rblgc := new(HHNNRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)

	return rblgc
}

//解析机器人个性
func (lgc *HHNNRobotLogic) parseRobotCha(cfgs *config.Group) {
	rand.Seed(time.Now().UnixNano())
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)
		cha := NewRobotCha(robotCfg, cname)
		lgc.robotChas = append(lgc.robotChas, cha)
	}
}

//获取机器人个性
func (lgc *HHNNRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *HHNNRobotLogic) getGameLogic() *HHNNLogic {
	return lgc.GameLogic.(*HHNNLogic)
}

//游戏消息
func (lgc *HHNNRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	glc := lgc.getGameLogic()
	if glc.gamestate != Bet || glc.zhuangID == robot.UserInfo.UserID {
		return
	}

	if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) ||
		m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNBetRetID) {
		l := glc.getBetMoneyByType(1)
		h := glc.getBetMoneyByType(3)
		c := (l - h) / 100
		//mlog.Debug("c=%d", c)
		//玩家当庄的时候不平衡
		if glc.zhuangID != 0 && !glc.GetUserByUserID(glc.zhuangID).IsRobot() {
			c = 0
		}

		lgc.robotBet(robot.UserInfo, int32(c))
		return
	}
}

//机器人下注
func (lgc *HHNNRobotLogic) robotBet(ur *gamelogic.GameUser, cha int32) {
	mlog.Debug("机器人%d下注", ur.UserID)
	rcha := lgc.getRobotCha(ur.UserID)
	glc := lgc.getGameLogic()
	if glc.he > 0 && glc.zhuangID != 0 && !glc.GetUserByUserID(glc.zhuangID).IsRobot() {
		pre := rand.Int31n(100)
		if pre < glc.he {
			lgc.SendGameMessage(ur, lgc.getBetMsg(2, rcha.getBetCount()))
			return
		}
	}

	lgc.SendGameMessage(ur, lgc.getBetMsg(rcha.getBetType(cha), rcha.getBetCount()))
}

func (lgc *HHNNRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

//获取下注消息
func (lgc *HHNNRobotLogic) getBetMsg(bettype int32, betcount int32) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.HHNN_GameMessageClassID_HHNNBetID)
	btdata := new(netproto.HHNN_BetArr)
	betdata := new(netproto.HHNN_Bet)
	betdata.BetType = proto.Int32(bettype)
	betdata.Count = proto.Int32(betcount)
	btdata.BetArr = append(btdata.BetArr, betdata)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}
