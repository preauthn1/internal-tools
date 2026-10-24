package hljbp

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RobotCha 機器人個性
type RobotCha struct {
	betRate   int   // 下注機率 ,100 是為1%
	betRange  []int // 底注下注範圍,bet_range,1-7(籌碼編號),2(數量)
	moneyTime []int // 下注間隔時間 200,1000 為200-1000毫秒隨機
}

// NewRobotCha 創建機器人個性
func NewRobotCha(cf *config.Config, cname string) *RobotCha {
	rcha := new(RobotCha)

	s1 := cf.GetConfigValue(cname, "BetRate")
	rcha.betRate = rcha.GetIntValue(s1)

	s1 = cf.GetConfigValue(cname, "BetRange")
	if s1 != "" {
		strarr := strings.Split(s1, ",")
		for i := 0; i < len(strarr); i++ {
			rcha.betRange = append(rcha.betRange, rcha.GetIntValue(strarr[i]))
		}
	}

	s1 = cf.GetConfigValue(cname, "MoneyTime")
	if s1 != "" {
		strarr := strings.Split(s1, ",")
		for i := 0; i < len(strarr); i++ {
			rcha.moneyTime = append(rcha.moneyTime, rcha.GetIntValue(strarr[i]))
		}
	}

	return rcha
}

// GetIntValue 字串轉Int
func (rc *RobotCha) GetIntValue(str string) int {
	v, _ := strconv.ParseInt(str, 10, 32)
	return int(v)
}

// RobotLogic 機器人邏輯
type RobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

// NewRobotLogic 創建機器人邏輯
func NewRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *RobotLogic {
	rblgc := new(RobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.ParseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)
	return rblgc
}

// ParseRobotCha 解析機器人個性
func (lgc *RobotLogic) ParseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig("robot")

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)
		cha := NewRobotCha(robotCfg, cname)
		lgc.robotChas = append(lgc.robotChas, cha)
	}
}

// GetRobotCha 獲取機器人個性
func (lgc *RobotLogic) GetRobotCha(userid int32) *RobotCha {
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

// GetGameLogic 獲取遊戲邏輯
func (lgc *RobotLogic) GetGameLogic() *Logic {
	return lgc.GameLogic.(*Logic)
}

// OnGameEnd 遊戲結束
func (lgc *RobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

// OnGameMessage 遊戲封包
func (lgc *RobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

	//获取机器人个性
	cha := lgc.GetRobotCha(robot.UserInfo.UserID)
	glc := lgc.GetGameLogic()

	robotdata := glc.GetUserData(robot.UserInfo.UserID)
	if robotdata == nil {
		mlog.Debug("机器人数据是空")
		return
	}

	// S2C 牌局資訊，於玩家進入牌局或斷線重連時傳送
	if m.SClassID == int32(netproto.HLJBP_GameMessageClassID_HLJBP_GameInfo) {

	}

	// S2C 玩家資訊，當玩家人數或狀態有變動時傳送
	if m.SClassID == int32(netproto.HLJBP_GameMessageClassID_HLJBP_AllUserInfo) {
	}

	// S2C 玩家開始下注
	if m.SClassID == int32(netproto.HLJBP_GameMessageClassID_HLJBP_UserBettingStart) {
		rnd := rand.Intn(10000)

		if rnd < cha.betRate {
			betCount := rand.Intn(cha.betRange[2])
			//fmt.Printf("機器人[%d]押注次數:%d\n", robot.UserInfo.UserID, betCount+1)
			delay := int32(0)
			for i := 0; i <= betCount; i++ {
				delay += int32(cha.moneyTime[0] + rand.Intn(cha.moneyTime[1]-cha.moneyTime[0]))
				lgc.TFManager.NewMilliSecTimer(delay, func() {
					lgc.SendBetMsg(robot, glc, cha, robotdata)
				})
			}
		}
	}

	// S2C 總結算
	if m.SClassID == int32(netproto.HLJBP_GameMessageClassID_HLJBP_Settlement) {

	}

}

// SendBetMsg 傳送機器人下注
func (lgc *RobotLogic) SendBetMsg(robot *robotcomm.Robot, glc *Logic, cha *RobotCha, robotdata *User) {
	// 可下注籌碼
	chips := make([]int64, 0)
	for i := 0; i < len(glc.chip); i++ {
		if i < cha.betRange[0]-1 || i > cha.betRange[1]-1 {
			continue
		}

		if robotdata.totalBet+glc.chip[i] > glc.BetLimit && glc.BetLimit > 0 {
			break
		}

		if robotdata.money >= glc.chip[i] {
			chips = append(chips, glc.chip[i])
		}
	}

	if len(chips) > 0 {
		iBetIdx := rand.Intn(len(chips))
		var betMoney = chips[iBetIdx]
		lgc.SendGameMessageNow(robot.UserInfo, lgc.GetBetMsg(betMoney))
	} else {
		//fmt.Printf("機器人[%d]錢不夠\n", robot.UserInfo.UserID)
	}
}

// GetBetMsg 機器人下注訊息
func (lgc *RobotLogic) GetBetMsg(bet int64) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.HLJBP_GameMessageClassID_HLJBP_Bet)

	btdata := new(netproto.HLJBP_Msg_Bet)
	btdata.Bet = proto.Int64(bet)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}
