package dchjk

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	proto "example.invalid/golang/protobuf/proto"

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

//机器人个性
type RobotCha struct {
	iRate_Hit           int // 小於17點時要牌機率
	iRate_Hit_ReachSoft int // 大於等於17點時要牌機率
	iRate_Split         int // 可分牌時分牌機率
	iRate_Double        int // 可加倍時加倍機率
	iRate_Insurrence    int // 可保險時保險機率
	iRate_Surrender     int // 莊家明牌有10或是A時投降的機率
	iBetRange_Min       int // 最小下注籌碼索引值
	iBetRange_Max       int // 最大下注籌碼索引值
	iDefaultBetCount    int // 預設下注數量
}

//创建机器个性对象
func NewRobotCha(cf *config.Config, cname string) *RobotCha {
	rcha := new(RobotCha)

	s1 := cf.GetConfigValue(cname, "Rate_Hit")
	rcha.iRate_Hit = rcha.getIntValue(s1)

	s1 = cf.GetConfigValue(cname, "Rate_Be17")
	rcha.iRate_Hit_ReachSoft = rcha.getIntValue(s1)

	s1 = cf.GetConfigValue(cname, "Rate_Split")
	rcha.iRate_Split = rcha.getIntValue(s1)

	s1 = cf.GetConfigValue(cname, "Rate_Double")
	rcha.iRate_Double = rcha.getIntValue(s1)

	s1 = cf.GetConfigValue(cname, "Rate_Insurrence")
	rcha.iRate_Insurrence = rcha.getIntValue(s1)

	s1 = cf.GetConfigValue(cname, "Rate_Surrender")
	rcha.iRate_Surrender = rcha.getIntValue(s1)

	s1 = cf.GetConfigValue(cname, "BetRange")
	if s1 != "" {
		strarr := strings.Split(s1, ",")
		rcha.iBetRange_Min = rcha.getIntValue(strarr[0])
		rcha.iBetRange_Max = rcha.getIntValue(strarr[1])
		rcha.iDefaultBetCount = rcha.getIntValue(strarr[2])
	}

	return rcha
}

func (rc *RobotCha) getIntValue(str string) int {
	v, _ := strconv.ParseInt(str, 10, 32)
	return int(v)
}

//机器人逻辑
type DchjkRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

//新建一个机器人逻辑
func NewDchjkRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *DchjkRobotLogic {
	pRobot := new(DchjkRobotLogic)
	pRobot.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	pRobot.parseRobotCha(configs)
	pRobot.gameChas = make(map[int32]*RobotCha)

	return pRobot
}

//解析机器人个性
func (lgc *DchjkRobotLogic) parseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig("robot")

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)
		cha := NewRobotCha(robotCfg, cname)
		lgc.robotChas = append(lgc.robotChas, cha)
	}
}

//获取机器人个性
func (lgc *DchjkRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *DchjkRobotLogic) getGameLogic() *DchjkLogic {
	return lgc.GameLogic.(*DchjkLogic)
}

func (lgc *DchjkRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

//游戏消息
func (lgc *DchjkRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

	//获取机器人个性
	cha := lgc.getRobotCha(robot.UserInfo.UserID)
	glc := lgc.getGameLogic()

	iGuestSeat := glc.GetGuestSeat(robot.UserInfo)
	if iGuestSeat != gamelogic.FreeSeatNo && m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_AllUserInfo) {
		lgc.SendGameMessage(robot.UserInfo, lgc.getSelectSeat())
		mlog.Debug("旁观者选桌")
		return
	}

	robotdata := glc.getUserBySeatNo(robot.UserInfo.SeatNo)
	if robotdata == nil {
		mlog.Debug("机器人数据是空")
		return
	}

	// 傳送牌局資訊
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_GameInfo) {
	}

	// 傳送玩家資訊
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_AllUserInfo) {

	}

	// 等待遊戲開始
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_WaitingStart) {
	}

	// 開始下注
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_StartBetting) {
		robotdata.iBetTimes = 0 //該機器人下注次數計次歸0

		go func() {
			time.Sleep(time.Duration(500) * time.Millisecond)
			lgc.SendGameMessageNow(robot.UserInfo, lgc.getBetOKMsg())
		}()

		var iBetCount = rand.Intn(cha.iDefaultBetCount) + 1

		go func(iBetCount int, robot *robotcomm.Robot, glc *DchjkLogic, cha *RobotCha, robotdata *DchjkUser) {
			for i := 1; i <= iBetCount; i++ {
				rand.Seed(time.Now().UnixNano())
				time.Sleep(time.Duration(100+rand.Intn(500)) * time.Millisecond)
				lgc.robotBetOnce(iBetCount, robot, glc, cha, robotdata)
			}
		}(iBetCount, robot, glc, cha, robotdata)

	}

	// 傳送玩家下注結果
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_BetResult) {
	}

	// 播放停止下注動畫
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_StopBetting) {
	}

	// 現在輪到哪一個玩家
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_WhosTurn) {
	}

	// 進入切牌階段
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_UserCutCards) {
		data := m.MsgData.(*netproto.DCHJK_Msg_UserCutCards)
		if data.GetSeatNo()-int32(JD21_IndexBase) != robot.UserInfo.SeatNo {
			return
		}

		go func() {
			time.Sleep(time.Duration(500) * time.Millisecond)
			lgc.SendGameMessageNow(robot.UserInfo, lgc.getCutCardsOKMsg())
		}()
		lgc.SendGameMessage(robot.UserInfo, lgc.getCutCardsMsg())
	}

	// 玩家切牌結果
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_CutCardsResult) {
	}

	// 發牌
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_Deal) {
		go func() {
			time.Sleep(time.Duration(500) * time.Millisecond)
			lgc.SendGameMessageNow(robot.UserInfo, lgc.getDealOKMsg())
		}()
	}

	// 最後一局
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_LastGame) {
	}

	// 詢問玩家是否買保險
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_InsuranceOrNot) {
		rand.Seed(time.Now().UnixNano())
		iRnd := rand.Intn(100)
		if iRnd < cha.iRate_Insurrence {
			lgc.SendGameMessage(robot.UserInfo, lgc.getTakeInsurance(true)) // 買保險
		} else {
			lgc.SendGameMessage(robot.UserInfo, lgc.getTakeInsurance(false)) // 不買保險
		}
	}

	// 玩家買保險結果
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_TakeInsuranceResult) {
	}

	// 莊家瞇牌
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_DealerPeeksCards) {
		go func() {
			time.Sleep(time.Duration(500) * time.Millisecond)
			lgc.SendGameMessageNow(robot.UserInfo, lgc.getDealerPeeksCardsDoneMsg())
		}()
	}

	// 保險結算
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_Settlement_Insurance) {
		go func() {
			time.Sleep(time.Duration(500) * time.Millisecond)
			lgc.SendGameMessageNow(robot.UserInfo, lgc.getSettlementInsuranceOKMsg())
		}()
	}

	// 通知玩家操作
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_StartPlaying) {
		data := m.MsgData.(*netproto.DCHJK_Msg_StartPlaying)
		if data.GetSeatNo()-int32(JD21_IndexBase) != robot.UserInfo.SeatNo {
			return
		}

		rand.Seed(time.Now().UnixNano())
		iRnd := rand.Intn(100)

		// 莊家明牌有10或是A時投降
		pCard := glc.pDealerCard.pCardSet.GetSingleCard(0, false)
		if len(robotdata.ary_pCardSet) == 1 && robotdata.CardSetUsedNow().pCardSet.CardCount() == JD21_BasicCardCount && (g_map_iCardNoValue[pCard.I32Number] == 1 || g_map_iCardNoValue[pCard.I32Number] == 10) {
			if iRnd < cha.iRate_Surrender && !robotdata.CardSetUsedNow().IsSoftReached() {
				lgc.SendGameMessage(robot.UserInfo, lgc.getUserPlayingMsg(JD21_UserAction_Surrender, robotdata))
				return
			}
		}

		// 可分牌時分牌
		if robotdata.i64Money >= robotdata.CardSetUsedNow().i64BetBase && robotdata.CardSetUsedNow().bCanSplit {
			if iRnd < cha.iRate_Split {
				//if glc.userAction_Split(robotdata) == JD21_Playing_Success {
				lgc.SendGameMessage(robot.UserInfo, lgc.getUserPlayingMsg(JD21_UserAction_Split, robotdata))
				return
				//}
			}
		}

		// 可加倍時加倍
		if robotdata.i64Money >= robotdata.CardSetUsedNow().i64BetBase && robotdata.CardSetUsedNow().bCanDouble {
			if iRnd < cha.iRate_Double {
				//if glc.userAction_Double(robotdata) == JD21_Playing_Success {
				lgc.SendGameMessage(robot.UserInfo, lgc.getUserPlayingMsg(JD21_UserAction_Double, robotdata))
				return
				//}
			}
		}

		if !robotdata.CardSetUsedNow().IsSoftReached() {
			// 小於17點時要牌
			if iRnd < cha.iRate_Hit {
				//if glc.userAction_Hit(robotdata) == JD21_Playing_Success {
				lgc.SendGameMessage(robot.UserInfo, lgc.getUserPlayingMsg(JD21_UserAction_Hit, robotdata))
				return
				//}
			}
		} else {
			// 大於等於17點時要牌機率
			if !robotdata.CardSetUsedNow().bIsBlackJack && !robotdata.CardSetUsedNow().bIsBusted {
				if iRnd < cha.iRate_Hit_ReachSoft {
					//if glc.userAction_Hit(robotdata) == JD21_Playing_Success {
					lgc.SendGameMessage(robot.UserInfo, lgc.getUserPlayingMsg(JD21_UserAction_Hit, robotdata))
					return
					//}
				}
			}
		}

		// 停牌
		lgc.SendGameMessage(robot.UserInfo, lgc.getUserPlayingMsg(JD21_UserAction_Stand, robotdata))
	}

	// 傳送玩家動作結果
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_UserPlayingResult) {
	}

	// 莊家開牌/結算
	if m.SClassID == int32(netproto.DCHJK_GameMessageClassID_DCHJK_Settlement_Dealer) {
		// glc.i64Timeout = 0
		go func() {
			time.Sleep(time.Duration(500) * time.Millisecond)
			lgc.SendGameMessageNow(robot.UserInfo, lgc.getSettlementDealerOKMsg())
		}()
	}

	// 遊戲結束
	if m.SClassID == int32(netproto.GameMessageClassID_GameEndID) {
		cn := glc.TableRealUserCount()
		if cn <= 0 {
			glc.UserFree(robot.UserInfo)
		}
	}

}

// 機器人選位
func (lgc *DchjkRobotLogic) getSelectSeat() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_SelectSeat)

	btdata := new(netproto.DCHJK_Msg_SelectSeat)

	glc := lgc.getGameLogic()

	iSeat := -1
	for i, pUser := range glc.ary_pUserData {
		if pUser.User == nil {
			iSeat = i
			break
		}
	}

	btdata.SeatNo = proto.Int32(int32(iSeat) + int32(JD21_IndexBase))

	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//機器人買保險
func (lgc *DchjkRobotLogic) getTakeInsurance(isTake bool) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_TakeInsurance)

	btdata := new(netproto.DCHJK_Msg_TakeInsurance)
	if isTake {
		btdata.Action = proto.Int32(1)
	} else {
		btdata.Action = proto.Int32(0)
	}

	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//莊家瞇牌動畫播放完成
func (lgc *DchjkRobotLogic) getSettlementInsuranceOKMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_Settlement_Insurance_OK)

	return btmsg
}

//莊家瞇牌動畫播放完成
func (lgc *DchjkRobotLogic) getDealerPeeksCardsDoneMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_DealerPeeksCardsDone)

	return btmsg
}

//機器人切牌動畫播放完成
func (lgc *DchjkRobotLogic) getDealOKMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_DealAnimationOK)

	return btmsg
}

//機器人切牌動畫播放完成
func (lgc *DchjkRobotLogic) getCutCardsOKMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_CutCardsAnimationOK)

	return btmsg
}

//機器人切牌
func (lgc *DchjkRobotLogic) getCutCardsMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_CutCards)

	rand.Seed(time.Now().UnixNano())
	cardIndex := rand.Int31n(int32(JD21_MaxCutCardIndex - JD21_MinCutCardIndex))
	cardIndex += int32(JD21_MinCutCardIndex)
	btdata := new(netproto.DCHJK_Msg_CutCards)
	btdata.CardIndex = proto.Int32(cardIndex)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//機器人開牌結算完成
func (lgc *DchjkRobotLogic) getSettlementDealerOKMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_Settlement_Dealer_OK)
	return btmsg
}

//機器人下注單次處理
func (lgc *DchjkRobotLogic) robotBetOnce(iBetCount int, robot *robotcomm.Robot, glc *DchjkLogic, cha *RobotCha, robotdata *DchjkUser) {
	if glc.iGameState != JD21_State_Betting {
		return
	}
	if robotdata.iBetTimes >= iBetCount {
		return
	}

	// fmt.Println(robotdata.i32UserID, glc.ary_pUserData[robotdata.i32SeatNo].i64Money, robotdata.i64Money, glc.ary_i32Chip)

	// 可下注籌碼
	ary_i64Chip := make([]int64, 0)
	for i := 0; i < len(glc.ary_i32Chip); i++ {
		if i < cha.iBetRange_Min-1 || i > cha.iBetRange_Max-1 {
			continue
		}

		if robotdata.i64Money >= int64(glc.ary_i32Chip[i]) {
			ary_i64Chip = append(ary_i64Chip, int64(glc.ary_i32Chip[i]))
		}
	}

	if len(ary_i64Chip) > 0 {
		iBetIdx := rand.Intn(len(ary_i64Chip))
		var betMoney = ary_i64Chip[iBetIdx]
		//iResult := glc.userAction_Bet(robotdata, ary_i64Chip[iBetIdx])
		//if iResult == JD21_BetResult_Success {
		robotdata.iBetTimes++ //增加該機器人下注次數計次
		lgc.SendGameMessageNow(robot.UserInfo, lgc.getBetMsg(betMoney))
		//}
	}
}

//機器人停止下注動畫播放完成
func (lgc *DchjkRobotLogic) getBetOKMsg() *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_StopBettingAnimationOK)

	return btmsg
}

//機器人下注
func (lgc *DchjkRobotLogic) getBetMsg(betMoney int64) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_Bet)

	btdata := new(netproto.DCHJK_Msg_Bet)
	btdata.Bet = proto.Int64(betMoney)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//機器人動作
func (lgc *DchjkRobotLogic) getUserPlayingMsg(action int, pUser *DchjkUser) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DCHJK_GameMessageClassID_DCHJK_UserPlaying)

	btdata := new(netproto.DCHJK_Msg_UserPlaying)
	btdata.Action = proto.Int32(int32(action))
	btdata.CardSetIndex = proto.Int32(int32(pUser.CardSetUsedNow().iIndex))
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}
