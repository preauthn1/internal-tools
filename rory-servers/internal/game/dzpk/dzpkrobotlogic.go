// Package dzpk ?
// ==================================================================
// 文件名： dzpkrobotlogic.go
// 說明：   機器人邏輯
// 修改記錄：
//==================================================================
package dzpk

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

// ?
const (
	RobotLogicConfigFileName = "robot"
)

// RobotCha 機器人個性
type RobotCha struct {
	sData map[int32]*ValueInfo //c1 c2
}

// ValueInfo ?
type ValueInfo struct {
	vData   map[int32]int32
	strName string //标签名
}

// RobotLogic 機器人邏輯
type RobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha

	robotCards    map[int32]*poker.CardGroup
	mapFold       map[int32]int32
	lastGameState int32
}

// NewDZPKRobotLogic 創建機器人邏輯
func NewDZPKRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *RobotLogic {
	rblgc := new(RobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)
	rblgc.mapFold = make(map[int32]int32)
	rblgc.lastGameState = 0

	return rblgc
}

//解析机器人个性
func (lgc *RobotLogic) parseRobotCha(cfgs *config.Group) {

	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("size", "cn")
	sn := robotCfg.GetConfigInt("size", "sn")
	//mlog.Debug("解析机器人配置文件,cn[%d]", cn)
	for i := 1; i <= cn; i++ { //个性 //c1 c2 c3 ...
		cname := fmt.Sprintf("c%d", i)

		rData := new(RobotCha)
		rData.sData = make(map[int32]*ValueInfo, 0)
		for j := 1; j <= sn; j++ {
			sname := fmt.Sprintf("s%d", j)
			pstr := robotCfg.GetConfigValue(cname, sname)
			if pstr == "" {
				mlog.Error("配置文件错误,cname[%s],sname[%s]", cname, sname)
				break
			}
			vData := parseStrNewRate(pstr, ",")
			if vData != nil {
				rData.sData[int32(j)] = vData
				//mlog.Debug("保存机器人个性，名称[%s][%s]", cname, sname)
				continue
			}
			mlog.Error("保存个性失败，名称[%s][%s]", cname, sname)
		}
		//mlog.Debug("保存机器人个性，名称[%s]", cname)
		lgc.robotChas = append(lgc.robotChas, rData)
	}
}

// 解析机器人个性配置文件
func parseStrNewRate(str string, sep string) *ValueInfo {
	//mlog.Debug("解析配置文件")
	str = util.TrimAllSpace(str)
	sep = util.TrimAllSpace(sep)

	sarr := strings.Split(str, sep)

	data := new(ValueInfo)
	data.vData = make(map[int32]int32, 0)
	for index, strValue := range sarr {
		if index == 0 {
			data.strName = strValue
			continue
		}
		saveData, err := strconv.Atoi(strValue)
		if err != nil {
			mlog.Error("配置数据有误")
			return nil
		}
		data.vData[int32(index)] = int32(saveData)
	}
	return data

}

//获取机器人个性
func (lgc *RobotLogic) getRobotCha(userid int32) *RobotCha {
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

// 獲取遊戲邏輯
func (lgc *RobotLogic) getGameLogic() *DzpkLogic {
	return lgc.GameLogic.(*DzpkLogic)
}

// 獲取發牌器
func (lgc *RobotLogic) getPoker() *poker.FixedPoker {
	return lgc.getGameLogic().poker.(*poker.FixedPoker)
}

// OnGameMessage 遊戲消息
func (lgc *RobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	//底注阶段,收到底注倍数通知时

	//根据游戏阶段和公牌数量设置概率权重
	//金币不足，选择allin
	if lgc.isErrorNoMoney(robot, m) {
		lgc.robotAllin(robot.UserInfo)
		return
	}
	if lgc.isAction(robot, m) {
		lgc.robotAction(robot, m)
		return
	}

	//
}

// OnGameEnd 結束遊戲
func (lgc *RobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

//是否金币不足错误
func (lgc *RobotLogic) isErrorNoMoney(robot *robotcomm.Robot, m *msg.Message) bool {

	if int32(netproto.DZPK_GameMessageClassID_BroadCastErrorID) != int32(m.SClassID) {
		return false
	}
	if lgc.getGameLogic().currActionSeat != robot.UserInfo.SeatNo {
		return false
	}

	if m.MsgData == nil {
		mlog.Error("data is nil")
		return false
	}
	data := m.MsgData.(*netproto.DZPK_ErrorMsg)

	if *data.ErrorType == -10 {
		mlog.Debug("机器人操作金币不足，改为ALLIN")
		return true
	}
	return false
}

//操作
func (lgc *RobotLogic) isAction(robot *robotcomm.Robot, m *msg.Message) bool {

	if int32(netproto.DZPK_GameMessageClassID_BroadCastWaitOptID) != int32(m.SClassID) {
		return false
	}
	if lgc.getGameLogic().currActionSeat == robot.UserInfo.SeatNo {
		mlog.Debug("机器人操作开始，用户id[%d]", robot.UserInfo.UserID)
		return true
	}
	return false
}

// 机器人操作
func (lgc *RobotLogic) robotAction(robot *robotcomm.Robot, m *msg.Message) {
	//获取用户牌型
	ud := lgc.getGameLogic().getUserInfoBySeatNo(robot.UserInfo.SeatNo)
	if ud == nil {
		mlog.Error("机器人用户信息获取失败[robotAction],用户id[%d]", robot.UserInfo.UserID)
		return
	}
	handCards := ud.userCards
	publicCard := lgc.getGameLogic().savePublic
	if handCards == nil || publicCard == nil {
		mlog.Error("没有手牌或公牌,用户id[%d]", robot.UserInfo.UserID)
		return
	}
	if handCards.Len() < 2 || publicCard.Len() < 5 {
		mlog.Error("牌数量不对,用户id[%d]", robot.UserInfo.UserID)
		return
	}
	mlog.Debug("机器人牌手牌%s", handCards.String())
	robotCards := new(poker.CardGroup)
	robotCards.AddCardGroup(handCards)
	robotCards.AddCardGroup(publicCard)

	//解析用户牌型
	nType := GetCardType(robotCards)
	mlog.Debug("机器人牌型[%d]", nType)
	//获取机器人个性
	cha := lgc.getRobotCha(robot.UserInfo.UserID)
	if len(cha.sData) < 10 {
		mlog.Error("配置有误，机器人个性信息有误size[%d]", len(cha.sData))
		return
	}

	//选择操作方式
	nType = nType + 1
	if lgc.getGameLogic().betCurrMax == ud.currBet { //主动说话
		//根据概率选择动作，主动说话，可看牌，不可跟注
		actionType := lgc.getActionTypeActive(robot.UserInfo.SeatNo, nType, cha)
		lgc.actionSwitch(robot.UserInfo, actionType, nType, 1)
	} else {
		//根据概率选择动作，被动说话，不可看牌，可跟注
		actionType := lgc.getActionTypePassive(robot.UserInfo.SeatNo, nType, cha)
		lgc.actionSwitch(robot.UserInfo, actionType, nType, 0)
	}
}

//主动说话，可看牌，不可跟注
func (lgc *RobotLogic) getActionTypeActive(seatNo, cardType int32, r *RobotCha) (op int32) {
	weigth := int32(0)
	if lgc.getGameLogic().gameState <= GAMEPROGRESS_PREFLOP {
		weigth = 0
	} else if lgc.getGameLogic().gameState <= GAMEPROGRESS_FLOP {
		weigth = 60
	} else if lgc.getGameLogic().gameState <= GAMEPROGRESS_TURN {
		weigth = 80
	} else {
		weigth = 100
	}
	nFold := r.getSvData(cardType, INDEXA_FOLD)
	nCheck := r.getSvData(cardType, INDEXA_CHECK)
	nRaise := r.getSvData(cardType, INDEXA_RAISE)
	nAllin := r.getSvData(cardType, INDEXA_ALLIN)

	//mlog.Debug("被动说话,Fold[%d],Check[%d],Raise[%d],Allin[%d]", nFold, nCheck, nRaise, nAllin)

	nAllin = nAllin * weigth / 100

	lgc.reCalFold(seatNo, nFold)
	nFold = lgc.mapFold[seatNo]

	sum := nFold + nCheck + nRaise + nAllin
	if sum < 1 {
		sum = 1
	}
	random := int32(util.RndInt(0, int(sum)))
	defer func() {
		mlog.Debug("主动说话[%d],[%d/%d]Fold[%d],Check[%d],Raise[%d],Allin[%d]", op, random, sum, nFold, nCheck, nRaise, nAllin)
	}()

	if nFold != 0 && random <= nFold {
		return int32(FOLD) //弃牌
	}
	if nCheck != 0 && random <= nCheck+nFold {
		return int32(CHECK) //看牌
	}
	if nRaise != 0 && random <= nRaise+nCheck+nFold {
		return int32(RAISE)
	}
	if nAllin != 0 && random <= nAllin+nRaise+nCheck+nFold {
		return int32(ALLIN)
	}
	return int32(CHECK)

}

//被动说话，不可让牌，可跟注
func (lgc *RobotLogic) getActionTypePassive(seatNo, cardType int32, r *RobotCha) (op int32) {
	weigth := int32(0)
	if lgc.getGameLogic().gameState <= GAMEPROGRESS_PREFLOP {
		weigth = 0
	} else if lgc.getGameLogic().gameState <= GAMEPROGRESS_FLOP {
		weigth = 60
	} else if lgc.getGameLogic().gameState <= GAMEPROGRESS_TURN {
		weigth = 80
	} else {
		weigth = 100
	}
	nFold := r.getSvData(cardType, INDEXP_FOLD)
	nRaise := r.getSvData(cardType, INDEXP_RAISE)
	nAllin := r.getSvData(cardType, INDEXP_ALLIN)
	nCall := r.getSvData(cardType, INDEXP_CALL)

	nAllin = nAllin * weigth / 100

	lgc.reCalFold(seatNo, nFold)
	nFold = lgc.mapFold[seatNo]
	sum := nFold + nRaise + nAllin + nCall
	if sum < 1 {
		sum = 1
	}
	random := int32(util.RndInt(0, int(sum)))
	defer func() {
		mlog.Debug("被动说话[%d],[%d/%d],Fold[%d],Raise[%d],Allin[%d],Call[%d]", op, random, sum, nFold, nRaise, nAllin, nCall)
	}()
	if nFold != 0 && random <= nFold {
		return int32(FOLD) //弃牌
	}
	if nRaise != 0 && random <= nRaise+nFold {
		return int32(RAISE)
	}
	if nAllin != 0 && random <= nAllin+nRaise+nFold {
		return int32(ALLIN)
	}
	if nCall != 0 && random <= nCall+nAllin+nRaise+nFold {
		return int32(CALL)
	}
	return int32(CALL)
}

// 获取个性值
func (r *RobotCha) getSvData(a, b int32) int32 {
	if r == nil {
		mlog.Error("getSvData[%d][%d] error ,set default 0", a, b)
		return 0
	}
	vInfo := r.sData[a]
	if vInfo == nil {
		mlog.Error("getSvData[%d][%d] error ,set default 0", a, b)
		return 0
	}
	data := vInfo.vData[b]
	//mlog.Debug("获取数据[%d][%d]=[%d]", a, b, data)
	return data
}

// 选择操作动作
func (lgc *RobotLogic) actionSwitch(ur *gamelogic.GameUser, actionType, nType, attitude int32) {
	switch actionType {
	case int32(FOLD):
		lgc.robotFold(ur, attitude)
	case int32(CHECK):
		lgc.robotCheck(ur, attitude)
	case int32(CALL):
		lgc.robotCall(ur, attitude)
	case int32(RAISE):
		lgc.robotRaise(ur, attitude, nType)
	case int32(ALLIN):
		lgc.robotAllin(ur)
	default:
	}
}

// 机器人弃牌
func (lgc *RobotLogic) robotFold(ur *gamelogic.GameUser, attitude int32) {
	mlog.Debug("机器人操作弃牌")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_ActionFoldID)

	btdata := new(netproto.DZPK_Fold)
	btdata.Jetton = proto.Int64(0)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessage(ur, btmsg)
}

// 机器人看牌
func (lgc *RobotLogic) robotCheck(ur *gamelogic.GameUser, attitude int32) {

	ud := lgc.getGameLogic().getUserInfoBySeatNo(ur.SeatNo) //玩家下注额
	if ud == nil {
		mlog.Error("机器人用户数据不存在[robotCheck]，用户座位号[%d]", ur.SeatNo)
		return
	}

	if ud.currBet != lgc.getGameLogic().betCurrMax {
		lgc.robotFold(ur, attitude)
		return
	}

	mlog.Debug("机器人操作看牌")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_ActionCheckID)

	btdata := new(netproto.DZPK_Check)
	btdata.Jetton = proto.Int64(0)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessage(ur, btmsg)
}

// 机器人跟注
func (lgc *RobotLogic) robotCall(ur *gamelogic.GameUser, attitude int32) {

	ud := lgc.getGameLogic().getUserInfoBySeatNo(ur.SeatNo) //玩家下注额
	if ud == nil {
		mlog.Error("机器人用户数据不存在[robotCall]，用户座位号[%d]", ur.SeatNo)
		return
	}
	if lgc.getGameLogic().betCurrMax-ud.currBet > ud.userMoney {
		lgc.robotAllin(ur)
		return
	}

	mlog.Debug("机器人操作跟注")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_ActionCallID)

	btdata := new(netproto.DZPK_Call)
	btdata.Jetton = proto.Int64(0)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessage(ur, btmsg)
}

// 机器人加注
func (lgc *RobotLogic) robotRaise(ur *gamelogic.GameUser, attitude, nType int32) {
	mlog.Debug("机器人操作加注")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_ActionRaiseID)

	btdata := new(netproto.DZPK_Raise)

	cha := lgc.getRobotCha(ur.UserID)
	opValue := int64(0)

	//有限注調整
	if true {
		//获取快捷加注权重
		weightEasy := cha.getSvData(nType, INDEXA_EASYWEIGHT)
		//计算快捷加注金额,权重，底池
		opValue = lgc.calValueLimit(weightEasy, ur.SeatNo) //计算快捷加注

		ud := lgc.getGameLogic().getUserInfoBySeatNo(ur.SeatNo) //玩家下注额
		if ud == nil {
			mlog.Error("机器人数据不存在[robotRaise]，用户座位号[%d]", ur.SeatNo)
			return
		}
		if opValue > ud.userMoney {
			lgc.robotAllin(ur)
			return
		}
	} else {
		if attitude > 0 {
			//积极态度
			x := cha.getSvData(nType, INDEXA_MINMAX)
			//获取快捷加注权重
			weightEasy := cha.getSvData(nType, INDEXA_EASYWEIGHT)
			//获取自由加注权重
			weightFree := cha.getSvData(nType, INDEXA_FREEWEIGHT)
			//计算快捷加注金额,权重，底池，小盲
			nValueEasy := lgc.calValueEasy(weightEasy, ur.SeatNo) //计算快捷加注
			//计算自由加注金额
			nValueFree := lgc.calValueFree(weightFree, ur.SeatNo)
			if (nValueEasy-nValueFree)*int64(x) >= 0 { //这里需要乘以一个系数改变符号,x大于0取大的，x小于0取小的
				opValue = nValueEasy
			} else {
				opValue = nValueEasy
			}
			ud := lgc.getGameLogic().getUserInfoBySeatNo(ur.SeatNo) //玩家下注额
			if ud == nil {
				mlog.Error("机器人数据不存在[robotRaise]，用户座位号[%d]", ur.SeatNo)
				return
			}
			if opValue > ud.userMoney {
				lgc.robotAllin(ur)
				return
			}

		} else {
			//消极
			x := cha.getSvData(nType, INDEXP_MINMAX)
			//获取快捷加注权重
			weightEasy := cha.getSvData(nType, INDEXP_EASYWEIGHT)
			//获取自由加注权重
			weightFree := cha.getSvData(nType, INDEXP_FREEWEIGHT)
			//计算快捷加注金额,权重，底池，小盲
			nValueEasy := lgc.calValueEasy(weightEasy, ur.SeatNo) //计算快捷加注
			//计算自由加注金额
			nValueFree := lgc.calValueFree(weightFree, ur.SeatNo)
			if (nValueEasy-nValueFree)*int64(x) >= 0 { //这里需要乘以一个系数改变符号,x大于0取大的，x小于0取小的
				opValue = nValueEasy
			} else {
				opValue = nValueEasy
			}
			ud := lgc.getGameLogic().getUserInfoBySeatNo(ur.SeatNo) //玩家下注额
			if ud == nil {
				mlog.Error("机器人数据不存在[robotRaise]，用户座位号[%d]", ur.SeatNo)
				return
			}
			if opValue > ud.userMoney {
				lgc.robotAllin(ur)
				return
			}
		}
	}

	btdata.Jetton = proto.Int64(opValue)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessage(ur, btmsg)
}

// 机器人allin
func (lgc *RobotLogic) robotAllin(ur *gamelogic.GameUser) *msg.Message {
	mlog.Debug("机器人操作Allin")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_ActionAllinID)

	btdata := new(netproto.DZPK_AllIn)
	btdata.Jetton = proto.Int64(0)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessage(ur, btmsg)
	return btmsg

}

//计算快捷加注
func (lgc *RobotLogic) calValueEasy(weight, seatNo int32) int64 {
	//mlog.Debug("计算快捷加注金额")
	//快捷加注还需要乘以盲注加上底池
	nPool := lgc.getGameLogic().betTotal                 //底池
	nSmallBet := lgc.getGameLogic().betBase / 2          //小盲
	curMax := lgc.getGameLogic().betCurrMax              //当前牌局最大下注额
	ud := lgc.getGameLogic().getUserInfoBySeatNo(seatNo) //玩家下注额
	if ud == nil {
		mlog.Error("机器人数据不存在[calValueEasy]，用户座位号[%d]", seatNo)
		return 0
	}
	a1 := weight / 1000
	a2 := weight % 1000 / 100
	a3 := weight % 1000 % 100 / 10
	a4 := weight % 1000 % 100 % 10
	sum := a1 + a2 + a3 + a4
	if sum < 1 {
		sum = 1
	}
	randNum := int32(util.RndInt(int(0), int(sum)))
	if lgc.getGameLogic().potArray == nil {

		if randNum <= a1 {
			return (nPool+nSmallBet)/2 + curMax - ud.currBet
		}
		if randNum <= a1+a2 {
			return (nPool+nSmallBet)*2/3 + curMax - ud.currBet
		}
		if randNum <= a1+a2+a3 {
			return (nPool+nSmallBet)*3 + curMax - ud.currBet
		}
		return nPool/2 + curMax - ud.currBet
	}

	if randNum <= a1 {
		return nPool/2 + nSmallBet + curMax - ud.currBet
	}
	if randNum <= a1+a2 {
		return nPool*2/3 + nSmallBet + curMax - ud.currBet
	}
	if randNum <= a1+a2+a3 {
		return nPool*3 + nSmallBet + curMax - ud.currBet
	}
	return nPool/2 + nSmallBet + curMax - ud.currBet

	// return nPool/2 + nSmallBet + curMax - ud.currBet
}

//计算自由加注
func (lgc *RobotLogic) calValueFree(weight int32, seatNo int32) int64 {
	//mlog.Debug("计算自由加注金额")
	ud := lgc.getGameLogic().getUserInfoBySeatNo(seatNo) //玩家下注额
	if ud == nil {
		mlog.Error("机器人数据不存在[calValueFree]，用户座位号[%d]", seatNo)
		return 0
	}
	min := weight / 100
	max := weight % 100
	num := max - min
	if num < 1 {
		num = 1
	}
	rndNum := util.RndInt(0, int(num))
	rate := int64(min) + int64(rndNum)
	if min == max {
		rate = int64(min)
	}

	curMax := lgc.getGameLogic().betCurrMax //当前牌局最大下注额
	nRaise := rate * ud.userMoney / 100
	nSmallBet := lgc.getGameLogic().betBase / 2

	return curMax - ud.currBet + nRaise + nSmallBet
}

//計算有限注 快捷加注
func (lgc *RobotLogic) calValueLimit(weight, seatNo int32) int64 {
	//mlog.Debug("计算快捷加注金额")
	//快捷加注还需要乘以盲注加上底池
	nPool := lgc.getGameLogic().betTotal //底池
	//nSmallBet := lgc.getGameLogic().betBase / 2          //小盲
	curMax := lgc.getGameLogic().betCurrMax              //当前牌局最大下注额
	ud := lgc.getGameLogic().getUserInfoBySeatNo(seatNo) //玩家下注额
	if ud == nil {
		mlog.Error("机器人数据不存在[calValueEasy]，用户座位号[%d]", seatNo)
		return 0
	}
	a1 := weight / 1000
	a2 := weight % 1000 / 100
	a3 := weight % 1000 % 100 / 10
	a4 := weight % 1000 % 100 % 10
	sum := a1 + a2 + a3 + a4
	if sum < 1 {
		sum = 1
	}
	randNum := int32(util.RndInt(int(0), int(sum)))

	if randNum <= a1 {
		return (nPool+curMax-ud.currBet)/2 + curMax - ud.currBet
	}
	if randNum <= a1+a2 {
		return (nPool+curMax-ud.currBet)*2/3 + curMax - ud.currBet
	}
	if randNum <= a1+a2+a3 {
		return (nPool + curMax - ud.currBet) + curMax - ud.currBet
	}
	return (nPool+curMax-ud.currBet)/2 + curMax - ud.currBet
}

// 重新计算弃牌
func (lgc *RobotLogic) reCalFold(seatNo, nFold int32) {
	ud := lgc.getGameLogic().getUserInfoBySeatNo(seatNo)
	if ud == nil {
		mlog.Error("机器人数据不存在[reCalFold]，用户座位号[%d]", seatNo)
		return
	}

	publicCard := lgc.getGameLogic().savePublic
	robotCards := new(poker.CardGroup)
	robotCards.AddCardGroup(ud.userCards)
	robotCards.AddCardGroup(publicCard)

	if lgc.lastGameState != int32(lgc.getGameLogic().gameState) {
		lgc.mapFold[seatNo] += nFold
		lgc.lastGameState = int32(lgc.getGameLogic().gameState)
	}

	bigger := 0
	for _, ud := range lgc.getGameLogic().gameUserData {
		if ud.isFold {
			continue
		}
		userCards := new(poker.CardGroup)
		userCards.AddCardGroup(ud.userCards)
		userCards.AddCardGroup(publicCard)
		nRet := ComPareCards(robotCards, userCards) //
		if nRet < 0 {
			bigger++
			lgc.mapFold[seatNo] += nFold
		}
	}

	if bigger == 0 {
		lgc.mapFold[seatNo] = 0
	}
}
