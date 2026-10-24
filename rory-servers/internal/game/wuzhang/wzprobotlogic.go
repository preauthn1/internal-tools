package wzp

import (
	"fmt"
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
	//机器人配置文件名
	RobotLogicConfigFileName = "robot"
)

//机器人个性
type RobotCha struct {
	Discard     float32                 //弃牌概率累积值
	cardSizeAct map[int32]*RobotChaItem //牌型动作
}

// 创建机器人个性对象
func NewRobotCha() *RobotCha {
	cha := new(RobotCha)
	cha.cardSizeAct = make(map[int32]*RobotChaItem)

	return cha
}

//机器人个性子项
type RobotChaItem struct {
	FollowBet       bool    //是否跟注
	DoBet           bool    //是否下注
	AddBet          float32 //加注几率
	BetScale        []int32 //下注额分配比例
	FllowAfterAdd   bool    //有人加注后是否跟注
	Suoha           float32 //主动梭哈概率
	FllowAfterSuoha bool    //其他人梭哈后是否跟注
}

//创建机器人个性子项
func NewRobotChaItem(followBet bool, doBet bool, addBet float32, betScale []int32, fllowAfterAdd bool, suoha float32, followAfterSuoha bool) *RobotChaItem {
	chaitem := new(RobotChaItem)
	chaitem.FollowBet = followBet
	chaitem.DoBet = doBet
	chaitem.AddBet = addBet
	chaitem.BetScale = betScale
	chaitem.FllowAfterAdd = fllowAfterAdd
	chaitem.Suoha = suoha
	chaitem.FllowAfterSuoha = followAfterSuoha

	return chaitem
}

//是否梭哈
func (chaitem *RobotChaItem) CheckSuoha() bool {
	return chaitem.Suoha >= util.RndFloat()
}

//是否加注
func (chaitem *RobotChaItem) CheckAddBet() bool {
	return chaitem.AddBet >= util.RndFloat()
}

//五张牌机器人逻辑
type WZPRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha
}

func NewWZPRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *WZPRobotLogic {
	rblgc := new(WZPRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)

	return rblgc
}

func (lgc *WZPRobotLogic) getGameLogic() *WZPGameLogic {
	return lgc.GameLogic.(*WZPGameLogic)
}

func (lgc *WZPRobotLogic) getPoker() *poker.FixedPoker {
	return lgc.getGameLogic().poker.(*poker.FixedPoker)
}

//获取用户牌型
func (lgc *WZPRobotLogic) getUserCards(userid int32) *poker.CardGroup {
	cards := poker.NewCardGroup(nil)

	//获取用户已经发的牌型
	uexpoker := lgc.getGameLogic().gamePoker.GetUserPoker(userid)
	if uexpoker != nil {
		uexcards := uexpoker.Cards
		cards.Cards = append(cards.Cards, uexcards.Cards...)
	}

	//获取用户还没有发的牌型
	ucards := lgc.getPoker().GetUserCards(userid)
	if ucards != nil {
		cards.Cards = append(cards.Cards, ucards.Cards...)
	}

	return cards
}

//解析机器人个性
func (lgc *WZPRobotLogic) parseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")

	for i := 1; i <= cn; i++ {
		cname := fmt.Sprintf("c%d", i)

		cha := NewRobotCha()
		lgc.robotChas = append(lgc.robotChas, cha)

		//弃牌概率
		discard := robotCfg.GetConfigFloat32(cname, "discard")
		cha.Discard = discard

		//牌型动作
		for j := 1; j <= 10; j++ {
			csizename := fmt.Sprintf("c%d_size%d", i, j)

			chaitem := new(RobotChaItem)
			cha.cardSizeAct[int32(j)] = chaitem

			//跟注
			chaitem.FollowBet = robotCfg.GetConfigInt(csizename, "followbet") == 1
			//下注
			chaitem.DoBet = robotCfg.GetConfigInt(csizename, "dobet") == 1
			//加注
			chaitem.AddBet = robotCfg.GetConfigFloat32(csizename, "addbet")
			//投注比例
			betscales := robotCfg.GetConfigValue(csizename, "betscale")

			betscaleArr := strings.Split(betscales, ",")

			for _, bs := range betscaleArr {
				ibs, err := strconv.Atoi(bs)
				if err != nil {
					mlog.Error("解析机器人配置投注比例出现错误,%s,%v", betscales, err)
					return
				}

				chaitem.BetScale = append(chaitem.BetScale, int32(ibs))
			}
			//加注后是否跟注
			chaitem.FllowAfterAdd = robotCfg.GetConfigInt(csizename, "fllowafteradd") == 1
			//主动梭哈概率
			chaitem.Suoha = robotCfg.GetConfigFloat32(csizename, "suoha")
			//其他人梭哈是否跟
			chaitem.FllowAfterSuoha = robotCfg.GetConfigInt(csizename, "fllowaftersuoha") == 1
		}
	}
}

//获取机器人个性
func (lgc *WZPRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *WZPRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	//如果机器人已经弃牌
	if lgc.getGameLogic().turnManager.isGiveUp(robot.UserInfo.UserID) {
		return
	}

	ur := robot.UserInfo

	//游戏轮次消息
	if m.SClassID == int32(netproto.WZP_GameMessageClassID_UserGameTurnID) {
		//当前轮次由当前机器人操作
		if lgc.getGameLogic().turnManager.currentTurnUser.UserID != robot.UserInfo.UserID {
			return
		}
		//检查是否弃牌
		if lgc.checkDiscard(ur) {
			return
		}

		//获取机器人个性
		cha := lgc.getRobotCha(robot.UserInfo.UserID)

		//获取机器牌型
		cards := lgc.getUserCards(robot.UserInfo.UserID)
		csize := getCardSize(cards)

		//个性动作
		chaitem := cha.cardSizeAct[csize]

		//下注
		dobetFunc := func() {
			//下注金额
			amt := int64(gameconfig.minBetAmount())

			//加注
			if chaitem.CheckAddBet() {
				amt = lgc.getBetAmount(0, chaitem)
			}

			lgc.doBetAction(ur, actionID_DoBet, amt)
		}

		//跟注
		followbetFunc := func() {
			//判断上轮是否加注
			if lgc.getGameLogic().turnManager.addBet {
				if chaitem.FllowAfterAdd {
					//如果当前下注金额未达到上限
					currentBetAmt := lgc.getGameLogic().gameBetLog.currentTurnBetAmount
					if currentBetAmt < int64(gameconfig.maxBetAmount()) {
						if chaitem.CheckAddBet() {
							amt := lgc.getBetAmount(currentBetAmt, chaitem)

							lgc.doBetAction(ur, actionID_AddBet, amt)
						} else {
							lgc.doBetAction(ur, actionID_Followbet, 0)
						}
					} else {
						lgc.doBetAction(ur, actionID_Followbet, 0)
					}
				} else {
					lgc.discard(ur)
				}
			} else {
				if chaitem.FollowBet {
					//当前下注金额
					currentBetAmt := lgc.getGameLogic().gameBetLog.currentTurnBetAmount
					if currentBetAmt < int64(gameconfig.maxBetAmount()) {
						if chaitem.CheckAddBet() {
							amt := lgc.getBetAmount(currentBetAmt, chaitem)

							lgc.doBetAction(ur, actionID_AddBet, amt)
						} else {
							lgc.doBetAction(ur, actionID_Followbet, 0)
						}
					} else {
						lgc.doBetAction(ur, actionID_Followbet, 0)
					}
				} else {
					lgc.discard(ur)
				}
			}
		}

		//如果是下注用户
		if lgc.getGameLogic().turnManager.isFirstTurn() {
			if chaitem.DoBet {
				//不是最后一圈
				if lgc.getGameLogic().turnManager.turnNo != 4 {
					dobetFunc()
				} else {
					//最后一圈
					//是否主动梭哈
					if chaitem.CheckSuoha() {
						lgc.doBetAction(ur, actionID_Suoha, 0)
					} else {
						dobetFunc()
					}
				}
			} else {
				lgc.discard(ur)
			}
		} else {
			//如果已经梭哈
			if lgc.getGameLogic().turnManager.Suoha {
				if chaitem.FllowAfterSuoha {
					lgc.doBetAction(ur, actionID_Followbet, 0)
				} else {
					lgc.discard(ur)
				}
			} else {
				followbetFunc()
			}
		}
	}

}

//获取投注金额
func (lgc *WZPRobotLogic) getBetAmount(min int64, chaitem *RobotChaItem) int64 {

	mlog.Debug("获取下注金额%d,%v", min, chaitem)

	//计算筹码索引
	btidx := -1
	var nextAmt int64 = -1
	for i, amt := range gameconfig.betAmount {
		if int64(amt) > min {
			btidx = i
			nextAmt = int64(amt)
			break
		}
	}

	if btidx == -1 {
		return int64(gameconfig.minBetAmount())
	}

	var total int32 = 0

	for i, sc := range chaitem.BetScale {
		if i < btidx {
			continue
		}

		total += sc
	}

	if total == 0 {
		//example.invalid("计算金额1%d", nextAmt)
		return nextAmt
	}

	rndv := int32(util.RndInt(1, int(total+1)))

	total = 0
	idx := 0
	for i, sc := range chaitem.BetScale {
		if i < btidx {
			idx++
			continue
		}

		total += sc

		if total >= rndv {
			//example.invalid("计算金额2%d", gameconfig.betAmount[idx])
			return int64(gameconfig.betAmount[idx])
		}

		idx++
	}

	//example.invalid("计算金额3%d", gameconfig.betAmount[idx])
	return int64(gameconfig.betAmount[idx])
}

//弃牌
func (lgc *WZPRobotLogic) discard(ur *gamelogic.GameUser) {
	discardmsg := new(msg.Message)
	discardmsg.BClassID = int32(netproto.MessageBClassID_Game)
	discardmsg.SClassID = int32(netproto.WZP_GameMessageClassID_UserGameActionID)

	discarddata := new(netproto.WZP_UserGameAction)
	discarddata.ActionID = proto.Int32(actionID_Giveup)
	discarddata.BetAmount = proto.Int64(0)
	discardmsg.MsgData, _ = proto.Marshal(discarddata)

	lgc.SendGameMessage(ur, discardmsg)
}

//检查机器人是否弃牌
func (lgc *WZPRobotLogic) checkDiscard(ur *gamelogic.GameUser) bool {
	//获取机器人个性
	cha := lgc.getRobotCha(ur.UserID)

	//当前轮次
	turnno := lgc.getGameLogic().turnManager.turnNo
	//计算弃牌概率
	discardOdds := float32(turnno) * cha.Discard

	//计算最大牌型
	var maxcards *poker.CardGroup
	var maxCardsUserID int32 = 0

	for _, ur := range lgc.getGameLogic().AllGameUser() {
		if lgc.getGameLogic().turnManager.isGiveUp(ur.UserID) {
			continue
		}

		//获取当前用户桌面上的牌型
		ucards := lgc.getUserCards(ur.UserID)
		//ucards := lgc.getGameLogic().gamePoker.GetUserPoker(ur.UserID).Cards
		if maxcards == nil {
			maxcards = ucards
			maxCardsUserID = ur.UserID
		} else {
			ret, _, _ := compareCard(maxcards, ucards)

			if !ret {
				maxcards = ucards
				maxCardsUserID = ur.UserID
			}
		}
	}

	if maxCardsUserID != ur.UserID && maxCardsUserID != 0 {
		if discardOdds >= util.RndFloat() {
			lgc.discard(ur)
			mlog.Debug("牌型弃牌几率检查使得机器人%v弃牌", ur)
			return true
		}
	}

	return false
}

//投注
func (lgc *WZPRobotLogic) doBetAction(ur *gamelogic.GameUser, actID int32, amt int64) {
	//加注的特殊处理
	if actID == actionID_AddBet {

		// 2020.02.11 為了讓最這回合最後一個下注機器人用戶可以做加注動作，所以注解這裡 by Alex
		//如果是最后一个下注用户
		// if lgc.getGameLogic().turnManager.isLastTurn() {
		// 	actID = actionID_Followbet
		// }

		// fmt.Println("機器人想加注", lgc.getGameLogic().turnManager.turnNo, amt, lgc.getGameLogic().gameBetLog.currentTurnBetAmount, lgc.getGameLogic().gameBetLog.checkTurnBetAmount(amt, lgc.getGameLogic().turnManager.turnNo, ur))

		if !lgc.getGameLogic().gameBetLog.checkTurnBetAmount(amt, lgc.getGameLogic().turnManager.turnNo, ur) {
			actID = actionID_Followbet
		}

		//如果下注金额为0
		if amt == 0 {
			actID = actionID_Followbet
		}
	}

	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.WZP_GameMessageClassID_UserGameActionID)

	btdata := new(netproto.WZP_UserGameAction)
	btdata.ActionID = proto.Int32(actID)
	btdata.BetAmount = proto.Int64(amt)
	btmsg.MsgData, _ = proto.Marshal(btdata)

	lgc.SendGameMessage(ur, btmsg)
}

func (lgc *WZPRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}
