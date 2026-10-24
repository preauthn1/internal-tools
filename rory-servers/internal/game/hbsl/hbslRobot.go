package hbsl

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"

	"example.invalid/ServerTeam/Servers/pkg/util"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
)

// 機器人個性
type RobotCha struct {
	iPickRate      int   // 搶紅包機率
	iPickTiming    []int // 搶紅包時間點
	iMoneyTime     []int // 發紅包間隔時間
	ary_iMoneyPack []int // 紅包金額機率
	iRobotbombRate int   // 機器人作弊機率
	iLeaveRate     int   // 機器人離場率
}

// 創建機器人個性
func NewRobotCha(cf *config.Config, cname string) *RobotCha {
	rcha := new(RobotCha)

	s1 := cf.GetConfigValue(cname, "PickRate")
	rcha.iPickRate = rcha.GetIntValue(s1)

	s1 = cf.GetConfigValue(cname, "PickTiming")
	if s1 != "" {
		strarr := strings.Split(s1, ",")
		for i := 0; i < len(strarr); i++ {
			rcha.iPickTiming = append(rcha.iPickTiming, rcha.GetIntValue(strarr[i]))
		}
	}

	s1 = cf.GetConfigValue(cname, "MoneyTime")
	if s1 != "" {
		strarr := strings.Split(s1, ",")
		for i := 0; i < len(strarr); i++ {
			rcha.iMoneyTime = append(rcha.iMoneyTime, rcha.GetIntValue(strarr[i]))
		}
	}

	s1 = cf.GetConfigValue(cname, "MoneyPack")
	if s1 != "" {
		strarr := strings.Split(s1, ",")
		for i := 0; i < len(strarr); i++ {
			rcha.ary_iMoneyPack = append(rcha.ary_iMoneyPack, rcha.GetIntValue(strarr[i]))
		}
	}

	s1 = cf.GetConfigValue(cname, "RobotbombRate")
	rcha.iRobotbombRate = rcha.GetIntValue(s1)

	s1 = cf.GetConfigValue(cname, "LeaveRate")
	rcha.iLeaveRate = rcha.GetIntValue(s1)

	return rcha
}

func (rc *RobotCha) GetIntValue(str string) int {
	v, _ := strconv.ParseInt(str, 10, 32)
	return int(v)
}

// 機器人邏輯
type HbslRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
	robotFaHB map[int32]*event.TimerFunc
}

// 創建機器人邏輯
func NewHbslRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *HbslRobotLogic {
	rblgc := new(HbslRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.ParseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)
	rblgc.robotFaHB = make(map[int32]*event.TimerFunc)
	return rblgc
}

// 解析機器人個性
func (lgc *HbslRobotLogic) ParseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig("robot")

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)
		cha := NewRobotCha(robotCfg, cname)
		lgc.robotChas = append(lgc.robotChas, cha)
	}
}

// 獲取機器人個性
func (lgc *HbslRobotLogic) GetRobotCha(userid int32) *RobotCha {
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

func (lgc *HbslRobotLogic) GetGameLogic() *HbslLogic {
	return lgc.GameLogic.(*HbslLogic)
}

func (lgc *HbslRobotLogic) OnGameStart(robot *robotcomm.Robot) {

}

func (lgc *HbslRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

// 遊戲封包
func (lgc *HbslRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

	// 獲取機器人個性
	cha := lgc.GetRobotCha(robot.UserInfo.UserID)
	//glc := lgc.GetGameLogic()

	// S2C 牌局資訊，於玩家進入牌局或斷線重連時傳送
	if m.SClassID == int32(netproto.HBSL_GameMessageClassID_HBSL_GameInfo) {
		_, ok := lgc.robotFaHB[robot.UserInfo.UserID]
		if !ok {
			lgc.robotFaHB[robot.UserInfo.UserID] = lgc.TFManager.NewTimer(0, func() {
				lgc.FaHBLoop(robot)
			})
		}
	}

	// S2C 在線玩家結果
	if m.SClassID == int32(netproto.HBSL_GameMessageClassID_HBSL_AllUserInfoResult) {
	}

	// S2C 開始搶紅包
	if m.SClassID == int32(netproto.HBSL_GameMessageClassID_HBSL_StartGrab) {
		iRnd := rand.Intn(10000)
		if iRnd < cha.iPickRate {
			lgc.TFManager.NewMilliSecTimer(int32(cha.iPickTiming[0]+rand.Intn(cha.iPickTiming[1]-cha.iPickTiming[0])), func() {
				lgc.SendGameMessageNow(robot.UserInfo, lgc.GetGrabMsg())
			})
		}
	}

	// S2C 搶紅包結果
	if m.SClassID == int32(netproto.HBSL_GameMessageClassID_HBSL_GrabResult) {
	}

	// S2C 發紅包結果
	if m.SClassID == int32(netproto.HBSL_GameMessageClassID_HBSL_FaHBResult) {
	}

	// S2C 紅包列表結果
	if m.SClassID == int32(netproto.HBSL_GameMessageClassID_HBSL_AllHBInfoResult) {
	}

	// S2C 結算
	if m.SClassID == int32(netproto.HBSL_GameMessageClassID_HBSL_Settlement) {
	}

	if m.SClassID == int32(netproto.GameMessageClassID_GameEndID) {
		//cha := lgc.GetRobotCha(robot.UserInfo.UserID)
		//iRnd := rand.Intn(100)
		//// 收封包处里 OnGameEnd处里底层已过检查式，free之后会造成报错
		//if iRnd < cha.iLeaveRate {
		//	pUser := lgc.GetGameLogic().GetUserDataByID(robot.UserInfo.UserID)
		//
		//	if pUser != nil {
		//		if pUser.pGrabHB != nil {
		//			return
		//		}
		//
		//		if pUser.pFaHB != nil {
		//			return
		//		}
		//	}
		//
		//	if f, ok := lgc.robotFaHB[robot.UserInfo.UserID]; ok {
		//		f.Stop()
		//		delete(lgc.robotFaHB, robot.UserInfo.UserID)
		//	}
		//	lgc.GetGameLogic().UserFree(robot.UserInfo)
		//
		//	//fmt.Printf("機器人[%d] 離開\n", robot.UserInfo.UserID)
		//	mlog.Debug("機器人[%d] 離開", robot.UserInfo.UserID)
		//}
	}

}

// 機器人搶紅包
func (lgc *HbslRobotLogic) GetGrabMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.HBSL_GameMessageClassID_HBSL_Grab)

	return btmsg
}

// 機器人發紅包
func (lgc *HbslRobotLogic) GetFaHBMsg(cha *RobotCha) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.HBSL_GameMessageClassID_HBSL_FaHB)

	btdata := new(netproto.HBSL_Msg_FaHB)

	iRnd := rand.Intn(100)

	iIndex := 0
	iValue := 0

	for i := 0; i < len(cha.ary_iMoneyPack); i++ {
		iValue += cha.ary_iMoneyPack[i]
		if iRnd < iValue {
			iIndex = i
			break
		}
	}

	btdata.HBMoney = proto.Int64(lgc.GetGameLogic().ary_i64Chip[iIndex])

	iRnd = rand.Intn(10)
	btdata.LandmineNumber = proto.Int32(int32(iRnd))
	btmsg.MsgData, _ = proto.Marshal(btdata)

	return btmsg
}

// 定時發紅包
func (lgc *HbslRobotLogic) FaHBLoop(robot *robotcomm.Robot) {
	pGameUser := robot.UserInfo
	if _, ok := lgc.robotFaHB[pGameUser.UserID]; ok {
		ur := lgc.GetGameLogic().GetUserByUserID(pGameUser.UserID)
		if ur == nil || robot.IsLeave() {
			delete(lgc.robotFaHB, pGameUser.UserID)
			return
		}
		//mlog.System("用戶 %d 發紅包", pGameUser.UserID)
		cha := lgc.GetRobotCha(pGameUser.UserID)
		lgc.SendGameMessageNow(pGameUser, lgc.GetFaHBMsg(cha))
		iRand := rand.Intn(cha.iMoneyTime[1]-cha.iMoneyTime[0]) + cha.iMoneyTime[0]
		lgc.robotFaHB[pGameUser.UserID] = lgc.TFManager.NewTimer(int32(iRand), func() {
			lgc.FaHBLoop(robot)
		})
	}
}
