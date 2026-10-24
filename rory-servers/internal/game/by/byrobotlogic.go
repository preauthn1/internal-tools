package by

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	RobotLogicConfigFileName = "robot"
)

const (
	right = iota
	left
)

//机器人个性
type RobotCha struct {
	chageoddmin       int32   //切换倍率的最小时间
	chageoddmax       int32   //切换倍率的最大时间
	changrand         int     //切换倍率機率
	leavetime         int32   //离开的时间
	leaverate         int32   //离开的概率
	boomnumber        []int32 //可用炮彈number
	shoottime         int64   //間隔發射時間(ms)
	fishmincount      int32   //最少多少魚才發射
	slowmotion        []int32 //左右發射移動速度倍數(1~100倍)
	slowmotionrand    int
	posX              float32
	rightleft         int32 //right=true, left=false
	fixpointtimearr   []int64
	fixpointtimerand  int
	fixpointtime      int64 //固定點發射幾秒
	addtimearr        []int64
	addtime           int64 //增加秒數
	sleeptime         int64 //休息n秒在發射
	startfispointtime int64
	startsleeptime    int64
	isSuoDing         bool    //是否瞄準
	suodingrand       int     //瞄準機率
	suodingtimearr    []int64 //瞄準時間清單
	suodingtime       int64   //指定瞄準時間
	suodingfishid     int32   //瞄準fishid
}

//创建机器个性对象
func NewRobotCha(cf *config.Config, cname string) *RobotCha {
	rcha := new(RobotCha)
	s1 := cf.GetConfigValue(cname, "chageodd")
	s2 := strings.Split(s1, "|")
	changrand, _ := strconv.Atoi(s2[0])
	rcha.changrand = changrand
	strarr := strings.Split(s2[1], ",")
	min, _ := strconv.ParseInt(strarr[0], 10, 32)
	rcha.chageoddmin = int32(min)
	max, _ := strconv.ParseInt(strarr[1], 10, 32)
	rcha.chageoddmax = int32(max)
	rcha.leavetime = cf.GetConfigInt32(cname, "leavetime")
	rcha.leaverate = cf.GetConfigInt32(cname, "leaverate")
	rcha.shoottime = cf.GetConfigInt64(cname, "lastshoottime")
	rcha.fishmincount = cf.GetConfigInt32(cname, "fishmincount")
	slowmotionstr := cf.GetConfigValue(cname, "slowmotion")
	slowmotionstr2 := strings.Split(slowmotionstr, "|")
	slowmotionrand, _ := strconv.Atoi(slowmotionstr2[0])
	rcha.slowmotionrand = slowmotionrand
	rcha.slowmotion = append(rcha.slowmotion, parseInt32Arr(slowmotionstr2[1])...)
	fixpointtimestr := cf.GetConfigValue(cname, "fixpointtime")
	fixpointtimestr2 := strings.Split(fixpointtimestr, "|")
	fixpointtimerand, _ := strconv.Atoi(fixpointtimestr2[0])
	rcha.fixpointtimerand = fixpointtimerand
	rcha.fixpointtimearr = append(rcha.fixpointtimearr, parseInt64Arr(fixpointtimestr2[1])...)
	addtimestr := cf.GetConfigValue(cname, "addtime")
	rcha.addtimearr = append(rcha.addtimearr, parseInt64Arr(addtimestr)...)

	suodingstr := cf.GetConfigValue(cname, "suoding")
	suodingstr2 := strings.Split(suodingstr, "|")
	suodingrand, _ := strconv.Atoi(suodingstr2[0])
	rcha.suodingrand = suodingrand
	rcha.suodingtimearr = append(rcha.suodingtimearr, parseInt64Arr(suodingstr2[1])...)

	rcha.sleeptime = cf.GetConfigInt64(cname, "sleeptime")

	rcha.posX = float32(rand.Int31n(2))
	rcha.rightleft = rand.Int31n(2)

	strboom := cf.GetConfigValue(cname, "boomnumber")
	arrboom := strings.Split(strboom, ",")
	for _, v := range arrboom {
		boomnumber, _ := strconv.ParseInt(v, 10, 32)
		rcha.boomnumber = append(rcha.boomnumber, int32(boomnumber))
	}

	return rcha
}

func parseInt32Arr(s string) []int32 {
	result := []int32{}
	strarr := strings.Split(s, ",")
	for _, v := range strarr {
		r, _ := strconv.Atoi(v)
		result = append(result, int32(r))
	}
	return result
}

func parseInt64Arr(s string) []int64 {
	result := []int64{}
	strarr := strings.Split(s, ",")
	for _, v := range strarr {
		r, _ := strconv.Atoi(v)
		result = append(result, int64(r))
	}
	return result
}

//机器人逻辑
type BYRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

func NewBYRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *BYRobotLogic {
	rblgc := new(BYRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)

	return rblgc
}

//解析机器人个性
func (lgc *BYRobotLogic) parseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)

		cha := NewRobotCha(robotCfg, cname)
		lgc.robotChas = append(lgc.robotChas, cha)

	}
}

//获取机器人个性
func (lgc *BYRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *BYRobotLogic) getGameLogic() *Bylogic {
	return lgc.GameLogic.(*Bylogic)
}

//游戏消息
func (lgc *BYRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	defer util.PanicErrStack()

	if m.SClassID == int32(netproto.BY_GameMessageClassID_RobotID) {

		ust := &netproto.BY_FrameS2C{}
		if m.MsgData != nil {
			ust = m.MsgData.(*netproto.BY_FrameS2C)
		}

		glc := lgc.getGameLogic()
		ud := glc.getuserdatabyuserid(robot.UserInfo.UserID)
		if !ud.bcreatetimer {
			lgc.changeOddTimer(robot.UserInfo.UserID)
			lgc.checkLeave(robot.UserInfo.UserID)
			ud.bcreatetimer = true
		}
		cha := lgc.getRobotCha(robot.UserInfo.UserID)

		if int64(time.Now().UnixNano()-cha.startfispointtime) >= int64(cha.fixpointtime*int64(time.Millisecond)) {
			rndAddTime := util.RndInt(0, len(cha.addtimearr))
			AddTime := cha.addtimearr[rndAddTime]

			if int64(time.Now().UnixNano()-cha.startsleeptime) <= int64((cha.sleeptime+AddTime)*int64(time.Millisecond)) {
				//mlog.Debug("return sleep startsleeptime[%v] add[%v]", int64(time.Now().UnixNano()-cha.startsleeptime), int64((cha.sleeptime+AddTime)*int64(time.Millisecond) ))
				return
			} else {
				//mlog.System("sleep startsleeptime[%v] add[%v]", int64(time.Now().UnixNano()-cha.startsleeptime), int64((cha.sleeptime+AddTime)*int64(time.Millisecond) ))
				cha.startsleeptime = int64(time.Now().UnixNano())
			}
		}

		if (int64(time.Now().UnixNano())-ud.lastshoottime) >= int64(cha.shoottime*int64(time.Millisecond)) && ud.boomcount < glc.maxboom {
			//mlog.Debug("上一次射击的时间是%d int64(300*time.Millisecond)[%v]", (int64(time.Now().UnixNano()) - ud.lastshoottime), int64(300*time.Millisecond))
			//mlog.Debug("机器人id：%d", robot.UserInfo.UserID)
			lgc.SendGameMessageNow(robot.UserInfo, lgc.getShootMsg(ud.user, ust))
		} else if (int64(time.Now().UnixNano()) - ud.lastshoottime) >= int64(25*time.Second) {
			//mlog.Debug("2上一次射击的时间是%d", (int64(time.Now().UnixNano()) - ud.lastshoottime))
			//如果机器人10秒钟未开炮则清除炸弹自动开炮
			ud.delAllBoom()
			lgc.SendGameMessageNow(robot.UserInfo, lgc.getShootMsg(ud.user, ust))
		}
	}
}

//机器人切换倍率
func (lgc *BYRobotLogic) changeOddTimer(userid int32) {
	glc := lgc.getGameLogic()
	cha := lgc.getRobotCha(userid)

	time := rand.Int31n(cha.chageoddmax-cha.chageoddmin) + cha.chageoddmin
	glc.TimerFuncManager.NewTimer(time, func() {
		user := glc.GetUserByUserID(userid)
		if user != nil && glc.IsGameUser(user) {
			rand := util.RndInt(0, 100)
			if rand > cha.changrand {
				lgc.SendGameMessageNow(user, lgc.getChangeRateMsg(user.UserID))
			}
			lgc.changeOddTimer(userid)
		}
	})
}

//检测机器人是否离开
func (lgc *BYRobotLogic) checkLeave(userid int32) {
	glc := lgc.getGameLogic()
	cha := lgc.getRobotCha(userid)
	glc.TimerFuncManager.NewTimer(cha.leavetime, func() {
		user := glc.GetUserByUserID(userid)
		if user != nil && glc.IsGameUser(user) {
			if rand.Int31n(100) < cha.leaverate {
				glc.ChangeUserStateEnd(user)
				glc.deluserdatabyuserid(userid)
				glc.KickUserFunc(userid, "")
			} else {
				lgc.checkLeave(userid)
			}
		}
	})
}

func (lgc *BYRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

//射击消息
func (lgc *BYRobotLogic) getShootMsg(gameuser *gamelogic.GameUser, ust *netproto.BY_FrameS2C) *msg.Message {

	btmsg := new(msg.Message)

	if ust != nil && ust.Shoot != nil && len(ust.Shoot) > 0 {
		if gameuser.SeatNo == 2 || gameuser.SeatNo == 3 {
			btmsg = lgc.getSuoDingShootMsg(gameuser, ust)
		} else {
			btmsg = lgc.normalShoot(gameuser)
		}
	} else {
		btmsg = lgc.normalShoot(gameuser)
	}

	return btmsg
}

//一般射擊
func (lgc *BYRobotLogic) normalShoot(gameuser *gamelogic.GameUser) *msg.Message {
	seatno := gameuser.SeatNo
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameC2SID)
	cha := lgc.getRobotCha(gameuser.UserID)
	btdata := new(netproto.BY_FrameC2S)
	btdata.FrameID = proto.Int64(0)
	shoot := new(netproto.BY_ShootC2S)
	shoot.Pos = new(netproto.BY_Position)
	if !cha.isSuoDing {
		suoDing := new(netproto.BY_SuoDingC2S)
		suoDing.SuoDing = proto.Bool(false)
		btdata.SuoDing = suoDing
		cha.isSuoDing = true
	} else {
		btdata.SuoDing = nil
	}
	y := float32(1.0)
	//座位2和3在上面需要向下射击
	if seatno == 2 || seatno == 3 {
		y = -1.0
	}
	x := lgc.slowMotionPosX(gameuser.UserID, seatno)
	//x := lgc.randPosX()
	shoot.Pos.X = proto.Float32(x)
	shoot.Pos.Y = proto.Float32(y)
	btdata.Shoot = shoot
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//瞄準射擊
func (lgc *BYRobotLogic) getSuoDingShootMsg(gameuser *gamelogic.GameUser, ust *netproto.BY_FrameS2C) *msg.Message {
	cha := lgc.getRobotCha(gameuser.UserID)

	rndAddTime := util.RndInt(0, len(cha.addtimearr))
	AddTime := cha.addtimearr[rndAddTime]
	rand := util.RndInt(0, 100)
	if rand > cha.suodingrand {
		rnd2 := util.RndInt(0, len(cha.suodingtimearr))
		cha.suodingtime = cha.suodingtimearr[rnd2] + AddTime
	}

	if int64(time.Now().UnixNano()-cha.startfispointtime) >= int64(cha.suodingtime*int64(time.Millisecond)) {
		cha.startfispointtime = int64(time.Now().UnixNano())
		cha.suodingfishid = ust.Shoot[0].GetFishID()
	}

	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameC2SID)
	btdata := new(netproto.BY_FrameC2S)
	btdata.FrameID = proto.Int64(0)
	shoot := new(netproto.BY_ShootC2S)
	shoot.Pos = new(netproto.BY_Position)
	shoot.Pos.X = proto.Float32(ust.Shoot[0].GetPos().GetX())
	shoot.Pos.Y = proto.Float32(ust.Shoot[0].GetPos().GetY() * -1)
	shoot.FishID = proto.Int32(cha.suodingfishid)
	btdata.Shoot = shoot
	btmsg.MsgData, _ = proto.Marshal(btdata)

	return btmsg
}

//左右發射移動速度倍數
func (lgc *BYRobotLogic) slowMotionPosX(userid int32, seatno int32) float32 {
	cha := lgc.getRobotCha(userid)
	rnd := util.RndInt(0, len(cha.slowmotion))
	rand := util.RndInt(0, 100)
	slowmotion := float32(0)
	if rand > cha.slowmotionrand {
		slowmotion = float32(cha.slowmotion[rnd]) / float32(1000)
	}
	rndAddTime := util.RndInt(0, len(cha.addtimearr))
	AddTime := cha.addtimearr[rndAddTime]

	if int64(time.Now().UnixNano()-cha.startfispointtime) >= int64(cha.fixpointtime*int64(time.Millisecond)) {
		if cha.rightleft == right {
			cha.posX = cha.posX + slowmotion
		} else {
			cha.posX = cha.posX - slowmotion
		}
		cha.startfispointtime = int64(time.Now().UnixNano())

		rand = util.RndInt(0, 100)
		if rand > cha.fixpointtimerand {
			rnd2 := util.RndInt(0, len(cha.fixpointtimearr))
			cha.fixpointtime = cha.fixpointtimearr[rnd2] + AddTime
		}

	}

	if seatno == 2 {
		if cha.posX >= float32(1.5) {
			cha.rightleft = left
		} else if cha.posX <= float32(-0.5) {
			cha.rightleft = right
		}
	}

	if seatno == 3 {
		if cha.posX >= float32(0.5) {
			cha.rightleft = left
		} else if cha.posX <= float32(-1.5) {
			cha.rightleft = right
		}
	}

	x := cha.posX

	return x
}

//亂射
func (lgc *BYRobotLogic) randPosX() float32 {
	x := rand.Float32()

	if rand.Int31n(2) == 1 {
		x = -x
	}

	return x
}

func (lgc *BYRobotLogic) getChangeRateMsg(userid int32) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameC2SID)

	btdata := new(netproto.BY_FrameC2S)
	btdata.FrameID = proto.Int64(0)
	changerate := new(netproto.BY_UserChangeRateC2S)
	//glc := lgc.getGameLogic()
	cha := lgc.getRobotCha(userid)
	index := rand.Int31n(int32(len(cha.boomnumber)))
	changerate.RateAfterChange = proto.Int32(cha.boomnumber[index])
	btdata.UserChangeRate = changerate
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}
