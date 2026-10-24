package cd21

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
	HIT        = 1
	STAND      = 2
	SPLITCARD  = 3
	DOUBLECARD = 4
)

const (
	HARDTYPE = 1
	SOFTTYPE = 2
	DUITYPE  = 3
)

type RobotCha struct {
	minbetcount int32
	maxbetcount int32
	minadddizu  int32
	maxadddizu  int32
}

type baoXianArr struct {
	xianJiaDian int32
	jilu        int32
}

type tuoBaoArr struct {
	zhuangDian int32
	xianDian   []int32
}

type buPaiArr struct {
	buPaiType int32
	xianDian  int32
	action    []int32
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
	s2 := cf.GetConfigValue(cname, "adddizu")
	strarr2 := strings.Split(s2, ",")
	minadd, _ := strconv.ParseInt(strarr2[0], 10, 32)
	rcha.minadddizu = int32(minadd)
	maxadd, _ := strconv.ParseInt(strarr2[1], 10, 32)
	rcha.maxadddizu = int32(maxadd)
	rand.Seed(time.Now().UnixNano())

	return rcha
}

//获取下注数组
func (rc *RobotCha) getBetCount() int32 {
	i := rc.maxbetcount - rc.minbetcount + 1
	n := int32(i)
	return n * 100
}

//获取加注数组
func (rc *RobotCha) getAddDiZu() int32 {
	n := int32(1)
	return n * 1
}

//机器人逻辑
type CD21RobotLogic struct {
	*robotlogic.LogicBase
	robotChas  []*RobotCha
	gameChas   map[int32]*RobotCha
	daxiao13   int32
	lucky20    int32
	duizi      int32
	zhuangbao  int32
	zhuangdian []int32
	baoXian    []*baoXianArr
	tuoBao     []*tuoBaoArr
	buPai      []*buPaiArr
}

func NewCD21RobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *CD21RobotLogic {
	rblgc := new(CD21RobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.getrobotConfig(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)

	return rblgc
}

//解析机器人个性
func (lgc *CD21RobotLogic) parseRobotCha(cfgs *config.Group) {
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
func (lgc *CD21RobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *CD21RobotLogic) getGameLogic() *CD21gameLogic {
	return lgc.GameLogic.(*CD21gameLogic)
}

//游戏消息
func (lgc *CD21RobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	glc := lgc.getGameLogic()
	//下注時間
	if m.SClassID == int32(netproto.CD21_GameMessageClassID_CD21StateID) &&
		glc.gamestate == Bet {
		lgc.robotBet(robot.UserInfo)
		return
	}

	//機器人切牌
	if m.SClassID == int32(netproto.CD21_GameMessageClassID_CD21CutpaiID) &&
		glc.firstuser.user.SeatNo == robot.UserInfo.SeatNo &&
		glc.gamestate == Cutpai {
		lgc.robotCut(robot.UserInfo)
		return
	}

	//投保時間
	if m.SClassID == int32(netproto.CD21_GameMessageClassID_CD21ZhuangTuoBaoID) &&
		glc.gamestate == BaoXiao {
		lgc.robotCancel(robot.UserInfo)
		return
	}

	//要牌階段
	if m.SClassID == int32(netproto.CD21_GameMessageClassID_CD21CurrOptID) &&
		glc.gamestate == MainProgress && glc.currSeat == robot.UserInfo.SeatNo {
		ud := glc.getUserBySeatNo(robot.UserInfo.SeatNo)
		if ud.canOpt() && ud.user.SeatNo == glc.currSeat {
			//robotData :=  glc.getUserBySeatNo(robot.UserInfo.SeatNo)
			//if glc.getCardPoint(robotData.userCards) >= 17 {
			//	lgc.robotStand(robotData)
			//} else {
			//	lgc.robotHit(robotData)
			//}
			//return

			//Old
			//確認要牌及補牌 (對子:dui 有A:軟式 其他:硬式)
			bupaiType := int32(0)
			//獲取機器人資料
			robotData := glc.getUserBySeatNo(robot.UserInfo.SeatNo)
			if len(robotData.userCards.Cards) == 2 && robotData.userCards.Cards[0].CardNumber == robotData.userCards.Cards[1].CardNumber {
				//對子(分牌還沒做 先當硬式)
				bupaiType = HARDTYPE
			} else {
				bupaiType = HARDTYPE
				for _, cards := range robotData.userCards.Cards {
					if cards.CardNumber == 14 {
						//軟式
						bupaiType = SOFTTYPE
						break
					}
				}
			}
			zhuangdianflag := true
			bupaiflag := true
			for i := 0; i < len(lgc.zhuangdian); i++ {
				if glc.zhungDian == lgc.zhuangdian[i] { //獲得莊家牌點順位
					zhuangdianflag = false
					for _, buPaiArr := range lgc.buPai {
						if glc.getCardPoint(robotData.userCards) == buPaiArr.xianDian && bupaiType == buPaiArr.buPaiType {
							bupaiflag = false
							if len(robotData.userCards.Cards) == 2 && buPaiArr.action[i] == DOUBLECARD {
								lgc.robotDouble(robotData)
								break
							} else if buPaiArr.action[i] == STAND {
								//用戶停牌
								lgc.robotStand(robotData)
								break
							} else if buPaiArr.action[i] == SPLITCARD {
								//用戶分牌還沒做
								lgc.robotHit(robotData)
								break
							} else {
								//用戶要牌
								lgc.robotHit(robotData)
								break
							}

						}
					}
				}
			}
			if zhuangdianflag || bupaiflag {
				lgc.robotStand(robotData)
			}
			return
		}
	}
}

//机器人下注
func (lgc *CD21RobotLogic) robotBet(ur *gamelogic.GameUser) {
	rcha := lgc.getRobotCha(ur.UserID)
	lgc.SendGameMessage(ur, lgc.getBetMsg(DIZHU, rcha.getBetCount()))
	daxiao13 := rand.Int31n(100)
	lucky20 := rand.Int31n(100)
	duizi := rand.Int31n(100)
	zhuangbao := rand.Int31n(100)

	if daxiao13 < lgc.daxiao13 {
		//下大小13
		if daxiao13%2 == 0 {
			lgc.SendGameMessage(ur, lgc.getBetMsg(BIGGER13, rcha.getBetCount()))
		} else {
			lgc.SendGameMessage(ur, lgc.getBetMsg(SMALLER13, rcha.getBetCount()))
		}

	}
	if lucky20 < lgc.lucky20 {
		//下幸運20
		lgc.SendGameMessage(ur, lgc.getBetMsg(LUCKY20, rcha.getBetCount()))
	}
	if duizi < lgc.duizi {
		//下對子
		lgc.SendGameMessage(ur, lgc.getBetMsg(DUIZI, rcha.getBetCount()))
	}
	if zhuangbao < lgc.zhuangbao {
		//下莊爆
		lgc.SendGameMessage(ur, lgc.getBetMsg(ZHUANGBAO, rcha.getBetCount()))
	}

}

//机器人取消
func (lgc *CD21RobotLogic) robotCancel(ur *gamelogic.GameUser) {
	lgc.SendGameMessageNow(ur, lgc.getCancelMsg())
}

//机器人切牌
func (lgc *CD21RobotLogic) robotCut(ur *gamelogic.GameUser) {
	lgc.SendGameMessage(ur, lgc.getCutMsg(ur))
}

func (lgc *CD21RobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

//获取取消消息
func (lgc *CD21RobotLogic) getCancelMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserCancel)
	return btmsg
}

//获取切牌消息
func (lgc *CD21RobotLogic) getCutMsg(ur *gamelogic.GameUser) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BroadCastCutpaiRetID)
	btdata := new(netproto.CD21_BroadCastCutpaiRetID)
	btdata.SitNo = proto.Int32(ur.SeatNo)
	place := util.RndInt(53, 261)
	btdata.Place = proto.Int64(int64(place))
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//获取下注消息
func (lgc *CD21RobotLogic) getBetMsg(bettype int32, betcount int32) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BetID)
	btdata := new(netproto.CD21_BetArr)
	betdata := new(netproto.CD21_Bet)
	betdata.BetType = proto.Int32(bettype)
	betdata.Count = proto.Int32(betcount)
	btdata.BetArr = append(btdata.BetArr, betdata)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//獲取機率配置
func (lgc *CD21RobotLogic) getrobotConfig(cfgs *config.Group) {
	rand.Seed(time.Now().UnixNano())
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	//中間區域下注機率
	lgc.daxiao13 = int32(robotCfg.GetConfigInt("probability", "daxiao13"))
	lgc.lucky20 = int32(robotCfg.GetConfigInt("probability", "lucky20"))
	lgc.duizi = int32(robotCfg.GetConfigInt("probability", "duizi"))
	lgc.zhuangbao = int32(robotCfg.GetConfigInt("probability", "zhuangbao"))

	//保險機率表
	bn := robotCfg.GetConfigInt("baoXian", "bn")
	for i := 10; i < 10+bn; i++ {
		b := new(baoXianArr)
		b.xianJiaDian = int32(i)
		b.jilu = int32(robotCfg.GetConfigInt("baoXian", strconv.Itoa(i)))
		lgc.baoXian = append(lgc.baoXian, b)
	}

	//投降機率表
	tn := robotCfg.GetConfigInt("baoXian", "tn")
	for j := 2; j < 2+tn; j++ {
		rcha := new(tuoBaoArr)
		s1 := robotCfg.GetConfigValue("tuoXion", strconv.Itoa(j))
		rcha.zhuangDian = int32(j)
		//解析數字
		sarr := strings.Split(s1, "|")
		for _, s := range sarr {
			if util.TrimAllSpace(s) == "" {
				continue
			}
			dian, _ := strconv.Atoi(s)
			rcha.xianDian = append(rcha.xianDian, int32(dian))
			lgc.tuoBao = append(lgc.tuoBao, rcha)
		}
	}

	//補牌規則
	zhuangDainStr := robotCfg.GetConfigValue("buPai", "zhuang")

	strarr := strings.Split(zhuangDainStr, ",")
	for i := 0; i < len(strarr); i++ {
		dian, _ := strconv.ParseInt(strarr[i], 10, 32)
		lgc.zhuangdian = append(lgc.zhuangdian, int32(dian))
	}

	//硬式補牌
	hard := robotCfg.GetConfigInt("buPai", "hard")
	for x := 5; x < 2+hard; x++ {
		xcha := new(buPaiArr)
		s1 := robotCfg.GetConfigValue("buPai", strconv.Itoa(x))
		xcha.buPaiType = HARDTYPE
		xcha.xianDian = int32(x)
		//解析數字
		sarr := strings.Split(s1, ",")
		for _, s := range sarr {
			if util.TrimAllSpace(s) == "" {
				continue
			}
			action, _ := strconv.Atoi(s)
			xcha.action = append(xcha.action, int32(action))

		}
		lgc.buPai = append(lgc.buPai, xcha)
	}

	//軟式補牌
	soft := robotCfg.GetConfigInt("buPai", "soft")
	for y := 2; y < 2+soft; y++ {
		ycha := new(buPaiArr)
		s1 := robotCfg.GetConfigValue("buPai", strconv.Itoa(y))
		ycha.buPaiType = SOFTTYPE
		ycha.xianDian = int32(y)
		//解析數字
		sarr := strings.Split(s1, ",")
		for _, s := range sarr {
			if util.TrimAllSpace(s) == "" {
				continue
			}
			action, _ := strconv.Atoi(s)
			ycha.action = append(ycha.action, int32(action))

		}
		lgc.buPai = append(lgc.buPai, ycha)
	}

	//對子補牌
	dui := robotCfg.GetConfigInt("buPai", "duizi")
	for z := 2; z < 2+dui; z++ {
		zcha := new(buPaiArr)
		s1 := robotCfg.GetConfigValue("buPai", strconv.Itoa(z))
		zcha.buPaiType = DUITYPE
		zcha.xianDian = int32(z)
		//解析數字
		sarr := strings.Split(s1, ",")
		for _, s := range sarr {
			if util.TrimAllSpace(s) == "" {
				continue
			}
			action, _ := strconv.Atoi(s)
			zcha.action = append(zcha.action, int32(action))
		}
		lgc.buPai = append(lgc.buPai, zcha)
	}

}

func (lgc *CD21RobotLogic) robotHit(robot *cd21userdata) {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserHit)
	//mlog.System("--------------------HIT")
	lgc.SendGameMessage(robot.user, btmsg)
}

func (lgc *CD21RobotLogic) robotStand(robot *cd21userdata) {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserStand)
	//mlog.System("--------------------Stand")
	lgc.SendGameMessage(robot.user, btmsg)
}

func (lgc *CD21RobotLogic) robotDouble(robot *cd21userdata) {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserDouble)
	//mlog.System("--------------------Double")
	lgc.SendGameMessage(robot.user, btmsg)
}
