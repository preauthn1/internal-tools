//==================================================================
//文件名： ermjrobotlogic.go
//说明：   麻将机器人
//修改记录：
//==================================================================
package ermj

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mahjong"
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
	changeRate  int //换牌概率
	changeTimes int //单次换牌次数
	ziMoRate    int //听牌后自摸概率
}

type rateInfo struct {
}

//机器人逻辑
type SCMJRobotLogic struct {
	*robotlogic.LogicBase
	robotChas           []*RobotCha
	gameChas            map[int32]*RobotCha
	robotActionManagers map[int32]*RobotActionManager
}

func NewSCMJRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *SCMJRobotLogic {
	lgc := new(SCMJRobotLogic)
	lgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	lgc.gameChas = make(map[int32]*RobotCha, 0)
	lgc.robotActionManagers = make(map[int32]*RobotActionManager, 0)
	lgc.parseRobotCha(configs)
	return lgc
}

//解析机器人个性
func (lgc *SCMJRobotLogic) parseRobotCha(cfgs *config.Group) {
	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)
	cn := robotCfg.GetConfigInt("csize", "cn")
	mlog.Debug("解析机器人配置文件,cn[%d]", cn)
	lgc.robotChas = make([]*RobotCha, 0)
	for i := 1; i <= cn; i++ {
		data := new(RobotCha)
		cname := fmt.Sprintf("c%d", i)
		data.changeRate = robotCfg.GetConfigInt(cname, "robotechangrate")
		data.changeTimes = robotCfg.GetConfigInt(cname, "robotechangtime")
		data.ziMoRate = robotCfg.GetConfigInt(cname, "robotezimo")
		lgc.robotChas = append(lgc.robotChas, data)
	}
}

//获取机器人个性
func (lgc *SCMJRobotLogic) getRobotCha(userid int32) *RobotCha {
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

func (lgc *SCMJRobotLogic) getGameLogic() *MahjongSC {
	return lgc.GameLogic.(*MahjongSC)
}

func (lgc *SCMJRobotLogic) getRobotActionManager(robot *robotcomm.Robot) *RobotActionManager {
	if _, ok := lgc.robotActionManagers[robot.UserInfo.UserID]; !ok {
		lgc.robotActionManagers[robot.UserInfo.UserID] = NewRobotActionManager(lgc, robot)
	}
	return lgc.robotActionManagers[robot.UserInfo.UserID]
}
func (lgc *SCMJRobotLogic) deleteRobotActionManager(robot *robotcomm.Robot) {
	delete(lgc.robotActionManagers, robot.UserInfo.UserID)
}

//游戏消息
func (lgc *SCMJRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

	// 事件操作
	if lgc.isHaveCanOp(robot, m) {
		lgc.getRobotActionManager(robot).OnEvent(m)
		//lgc.actionEvent(robot, m)
		return
	}
	// 出牌
	if lgc.isOrdertoOut(robot, m) {
		lgc.getRobotActionManager(robot).OnWaitOp(m)
		//lgc.actionOutCard(robot, m)
		return
	}
	if lgc.isOrderDingque(robot, m) {
		lgc.getRobotActionManager(robot).OnDingque(m)
		//lgc.actionDingque(robot, m)
	}
}

func (lgc *SCMJRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha, 0)
	lgc.robotActionManagers = make(map[int32]*RobotActionManager, 0)
}

// 是否轮到机器人出牌
func (lgc *SCMJRobotLogic) isOrdertoOut(robot *robotcomm.Robot, m *msg.Message) bool {
	if int32(netproto.SCMJ_GameMessageClassID_SCMJWaitId) != int32(m.SClassID) || lgc.getGameLogic().gameState != PROGRESS_FIGHT {
		return false
	}
	if lgc.getGameLogic().currOpSeatNo != robot.UserInfo.SeatNo {
		return false
	}
	if lgc.getGameLogic().isEnd {
		return false
	}
	return true
}

// 是否轮到机器人定缺
func (lgc *SCMJRobotLogic) isOrderDingque(robot *robotcomm.Robot, m *msg.Message) bool {
	if int32(netproto.SCMJ_GameMessageClassID_SCMJWaitId) != int32(m.SClassID) || lgc.getGameLogic().gameState != PROGRESS_CHOOSE {
		return false
	}

	return true
}

// 是否轮到机器人可选操作
func (lgc *SCMJRobotLogic) isHaveCanOp(robot *robotcomm.Robot, m *msg.Message) bool {
	if int32(netproto.SCMJ_GameMessageClassID_SCMJCanOpEventId) != int32(m.SClassID) {
		return false
	}
	if m == nil {
		return false
	}

	return true
}

func (lgc *SCMJRobotLogic) actionEvent(robot *robotcomm.Robot, m *msg.Message) {
	data := m.MsgData.(*netproto.SCMJ_OpEventInfo)
	//mlog.System("%d 收到事件 hu[%v], gang[%v], peng[%v] gangvalue[%v] othervalue[%v]", robot.UserInfo.UserID, data.GetHu(), data.GetGang(), data.GetPeng(), data.GetGangValue(), data.GetOtherValue())
	if *data.Hu {
		opValue := *data.OtherValue
		lgc.doHU(robot.UserInfo, opValue)
		return
	}
	if *data.Gang {
		opValue := data.GangValue[0]
		lgc.doGang(robot.UserInfo, opValue)
		return
	}
	if *data.Peng {
		opValue := *data.OtherValue
		lgc.doPeng(robot.UserInfo, opValue)
		return
	}
	lgc.doPass(robot.UserInfo)
}

//actionDingque
func (lgc *SCMJRobotLogic) actionDingque(robot *robotcomm.Robot, m *msg.Message) {
	//mlog.Debug("定缺")
	ud := lgc.getGameLogic().getUserInfoBySeatNo(robot.UserInfo.SeatNo)
	if ud == nil {
		return
	}
	tong, tiao, wan := 0, 0, 0
	for _, cd := range ud.handCards.Cards {
		if cd/100 == mahjong.ColorTong {
			tong++
			continue
		}
		if cd/100 == mahjong.ColorTiao {
			tiao++
			continue
		}
		if cd/100 == mahjong.ColorWan {
			wan++
			continue
		}
	}

	opValue := tong
	opType := mahjong.ColorTong
	if tong >= tiao {
		opValue = tiao
		opType = mahjong.ColorTiao
	}
	if opValue > wan {
		opValue = wan
		opType = mahjong.ColorWan
	}
	//
	if lgc.getGameLogic().Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havechooserate") == 1 {
		nRand := util.RndInt(0, 100)
		if nRand > 50 {
			lgc.doRate(robot.UserInfo, 2)
		} else {
			lgc.doRate(robot.UserInfo, 1)
		}

	} else {
		lgc.doDingque(robot.UserInfo, opType)
	}

}

//机器人出牌
func (lgc *SCMJRobotLogic) actionOutCard(robot *robotcomm.Robot, m *msg.Message) {
	mlog.Debug("机器人出牌")
	ud := lgc.getGameLogic().getUserInfoBySeatNo(robot.UserInfo.SeatNo)
	if ud == nil {
		return
	}
	if ud.seatNo != lgc.getGameLogic().currOpSeatNo {
		mlog.Debug("没有轮到机器人操作,不需要出牌")
		return
	}
	//mlog.System("actionOutCard 桌子号[%d],%d 機器人,手牌[%s],原模牌[%d]", lgc.Table.TableNo, robot.UserInfo.UserID, ud.handCards.String(), ud.moPai)
	//startMoPai := ud.moPai
	// config := lgc.Configs.GetConfig("gameserver")
	// robotTimes := config.GetConfigInt("gameprivate", "robotechangtime")
	rob := lgc.getRobotCha(robot.UserInfo.UserID)
	if rob != nil {
		for i := 0; i < rob.changeTimes; i++ {
			lgc.ExchangeBestCardPihu(robot.UserInfo.SeatNo)
		}
	}

	if ud.canHu(ud.moPai) {
		ud.updataCanOpInfo(ud.moPai, false, false)
		//mlog.System("桌子号[%d],%d 機器人,手牌[%s],原模牌[%d],摸牌[%d],胡牌", lgc.Table.TableNo, robot.UserInfo.UserID, ud.handCards.String(), startMoPai, ud.moPai)
		// mlog.Debug("交换牌后可以胡牌")
		lgc.getGameLogic().curCanOpUser[ud.seatNo] = true
		lgc.doHU(robot.UserInfo, ud.moPai)
		return
	}
	opValue := ud.calBestOutCard()

	// mlog.Debug("手牌[%s],摸牌[%d],出牌[%d]", ud.handCards.String(), ud.moPai, opValue)
	//mlog.System("桌子号[%d],%d 機器人,手牌[%s],原模牌[%d],摸牌[%d],出牌[%d]", lgc.Table.TableNo, robot.UserInfo.UserID, ud.handCards.String(), startMoPai, ud.moPai, opValue)
	lgc.doOut(robot.UserInfo, opValue)
}
func (lgc *SCMJRobotLogic) doDingque(ur *gamelogic.GameUser, opValue int32) {
	mlog.Debug("机器人操作定缺")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCNJActionId)

	btdata := new(netproto.SCMJ_Action)
	btdata.OpType = proto.Int32(OP_CHOOSEDINGQUE)
	btdata.OpValue = append(btdata.OpValue, opValue)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessageNow(ur, btmsg)
}

// 选倍数
func (lgc *SCMJRobotLogic) doRate(ur *gamelogic.GameUser, opValue int32) {
	mlog.Debug("机器人操作选倍数")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCNJActionId)

	btdata := new(netproto.SCMJ_Action)
	btdata.OpType = proto.Int32(OP_CHOOSERATE)
	btdata.OpValue = append(btdata.OpValue, opValue)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessageNow(ur, btmsg)
}

//
func (lgc *SCMJRobotLogic) doHU(ur *gamelogic.GameUser, opValue int32) {
	mlog.Debug("机器人操作胡")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCNJActionId)

	btdata := new(netproto.SCMJ_Action)
	btdata.OpType = proto.Int32(OP_HU)
	btdata.OpValue = append(btdata.OpValue, opValue)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessageNow(ur, btmsg)
}
func (lgc *SCMJRobotLogic) doGang(ur *gamelogic.GameUser, opValue int32) {
	mlog.Debug("机器人操作杠")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCNJActionId)

	btdata := new(netproto.SCMJ_Action)
	btdata.OpType = proto.Int32(OP_GANG)
	btdata.OpValue = append(btdata.OpValue, opValue)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessageNow(ur, btmsg)
}
func (lgc *SCMJRobotLogic) doPeng(ur *gamelogic.GameUser, opValue int32) {
	mlog.Debug("机器人操作碰")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCNJActionId)

	btdata := new(netproto.SCMJ_Action)
	btdata.OpType = proto.Int32(OP_PENG)
	btdata.OpValue = append(btdata.OpValue, opValue)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessageNow(ur, btmsg)
}
func (lgc *SCMJRobotLogic) doOut(ur *gamelogic.GameUser, opValue int32) {
	mlog.Debug("机器人操作出牌")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCNJActionId)

	btdata := new(netproto.SCMJ_Action)
	btdata.OpType = proto.Int32(OP_OUTCARD)
	btdata.OpValue = append(btdata.OpValue, opValue)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessageNow(ur, btmsg)
}
func (lgc *SCMJRobotLogic) doPass(ur *gamelogic.GameUser) {
	mlog.Debug("机器人操作Pass")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCNJActionId)

	btdata := new(netproto.SCMJ_Action)
	btdata.OpType = proto.Int32(OP_PASS)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessageNow(ur, btmsg)
}

//换牌胡 屁胡
func (lgc *SCMJRobotLogic) ExchangeBestCardPihu(SeatNo int32) {
	// config := lgc.Configs.GetConfig("gameserver")
	// robotRate := config.GetConfigInt("gameprivate", "robotechangrate")

	ud := lgc.getGameLogic().getUserInfoBySeatNo(SeatNo)
	if ud == nil {
		return
	}
	rob := lgc.getRobotCha(ud.User.UserID)
	if rob == nil {
		return
	}
	nRand := util.RndInt(0, 100)
	if nRand > rob.changeRate {
		return
	}
	// mlog.Debug("****************************************************")
	// defer func() {
	// 	mlog.Debug("手牌数量[%d],牌池数量[%d]", ud.handCards.Len(), lgc.getGameLogic().randGroup.GetCardGroup().Len())
	// 	mlog.Debug("****************************************************")
	// }()
	// mlog.Debug("换牌前的牌[%s]", ud.handCards.String())
	// mlog.Debug("换牌前的牌数量[%d]", ud.handCards.Len())
	tmp := new(mahjong.CardGroup)
	for _, cd := range ud.handCards.Cards {
		tmp.AddCard(cd)
	}
	//mlog.System("ExchangeBestCardPihu %s", lgc.getGameLogic().randGroup.GetCardGroup().String())
	//获取用户手牌
	mapTmp := make(map[int32]int32, 0)
	for _, cd := range tmp.Cards {
		if cd/100 == ud.color {
			continue
		}
		mapTmp[cd/10]++
	}
	//出去一坎的
	for i := int32(10); i < 40; i++ {
		if i/10 == ud.color {
			continue
		}
		mapTmp = lgc.checkShunzi(mapTmp, i)
	}
	sum := int32(0)
	for _, count := range mapTmp {
		// mlog.Debug("count[%d]", count)
		sum += count
	}

	//获取牌池
	mapTmpCg := make(map[int32]int32, 0)
	cg := lgc.getGameLogic().randGroup.GetCardGroup()
	// mlog.Debug("换牌前的牌池数量[%d]", cg.Len())
	for _, cdcg := range cg.Cards {
		if cdcg/100 == ud.color {
			continue
		}
		mapTmpCg[cdcg/10]++
	}
	for _, cdcg := range ud.handCards.Cards {
		if cdcg/100 == ud.color {
			sum += 1
		}
	}
	//大于等于4换三张,小于4换胡
	// mlog.Debug("sum[%d]", sum)
	// mlog.Debug("牌池[%d]", cg.Len())
	if sum >= 4 {
		//在牌池中找一坎
		result := lgc.findThreeCard(mapTmpCg)
		if len(result) < 3 {
			//如果没有一坎 找对子
			result = lgc.findTwoCard(mapTmpCg)
		}
		if len(result) <= 0 {
			mlog.Debug("牌库无对应牌,换牌失败")
			return
		}

		cgNewResult := lgc.list2cardGroup(cg, result)

		siceOld := make([]int32, 0)
		for _, hd := range ud.handCards.Cards {
			if hd/100 == ud.color {
				siceOld = append(siceOld, hd/10)
			}
		}
		for cdOld, count := range mapTmp {
			if count == 1 {
				siceOld = append(siceOld, cdOld)
			} else if count == 2 {
				siceOld = append(siceOld, cdOld)
				siceOld = append(siceOld, cdOld)
			}
		}
		size := len(result)
		if len(siceOld) > size {
			siceOld = siceOld[:size]
		}
		cgOld := lgc.list2cardGroup(ud.handCards, siceOld)
		if cgNewResult.Len() != cgOld.Len() {
			mlog.Debug("交换牌数量不等1,换牌失败,NewSize[%d],OldSize[%d]", size, len(siceOld))
			mlog.Debug("交换牌数量不等1,换牌失败,cgNewResult[%d],cgOld[%d]", cgNewResult.Len(), cgOld.Len())
			return
		}

		cg.AddCardGroup(cgOld)
		cg.DeleteGroup(cgNewResult)
		//重新设置牌库
		if lgc.getGameLogic().randGroup.SetCardGroup(cg) {
			//交换牌
			ud.handCards.DeleteGroup(cgOld)
			ud.handCards.AddCardGroup(cgNewResult)
		} else {
			mlog.Debug("交换牌数量不等2,换牌失败")
			return
		}

	} else {
		huRand := util.RndInt(0, 100)
		if huRand > rob.ziMoRate {
			lgc.getGameLogic().randGroup.Shuffle()
			return
		}
		// 查找有用的牌siceHuCard
		siceHuCard := make([]int32, 0)
		mapTmpCg := make(map[int32]int32, 0)
		for _, cdcg := range cg.Cards {
			if cdcg/100 == ud.color {
				continue
			}
			mapTmpCg[cdcg/10]++
		}
		for mapcd, count := range mapTmpCg {
			if count >= 2 {
				siceHuCard = append(siceHuCard, mapcd)
				break
			}
		}
		if len(siceHuCard) == 1 && ud.moPai > 0 {
			// 获取没有用的牌siceOld
			siceOld := make([]int32, 0)
			for cdOld, count := range mapTmp {
				if count == 1 {
					siceOld = append(siceOld, cdOld)
					break
				}
			}
			// 转化为手牌
			cgOld := lgc.list2cardGroup(ud.handCards, siceOld)
			if cgOld.Len() != 1 {
				mlog.Debug("没有找到无用手牌 手牌[%s]，摸牌[%d]", ud.handCards.String(), ud.moPai)
				return
			}
			cgNew := lgc.list2cardGroup(cg, siceHuCard)
			if cgNew.Len() != 1 {
				mlog.Debug("没有找到有用牌池")
				return
			}
			cg.AddCardGroup(cgOld)
			cg.AddCard(ud.moPai)
			cg.DeleteGroup(cgNew)
			if lgc.getGameLogic().randGroup.SetCardGroup(cg) {
				//交换牌
				ud.handCards.DeleteGroup(cgOld)
				ud.handCards.AddCardGroup(cgNew)
				ud.moPai = cgNew.Cards[0]
			} else {
				mlog.Debug("交换牌数量不等3,换牌失败")
				return
			}
		}
	}
	lgc.getGameLogic().randGroup.Shuffle()
	//mlog.System("ExchangeBestCardPihu after Shuffle %s", lgc.getGameLogic().randGroup.GetCardGroup().String())
	lgc.getGameLogic().GameLog.AddGameStep("换牌后的牌[%s]", ud.handCards.String())
	// mlog.Debug("换牌后的牌[%s]", ud.handCards.String())
	// mlog.Debug("换牌成功 屁胡")
	return
}

//检测顺子并剔除
func (lgc *SCMJRobotLogic) checkShunzi(cd map[int32]int32, i int32) map[int32]int32 {
	if cd[i] >= 3 {
		cd[i] -= 3
	}
	if cd[i] > 0 && cd[i+1] > 0 && cd[i+2] > 0 {
		cd[i] -= 1
		cd[i+1] -= 1
		cd[i+2] -= 1
		cd = lgc.checkShunzi(cd, i)
		return cd
	}
	return cd
}

func (lgc *SCMJRobotLogic) findThreeCard(cd map[int32]int32) []int32 {
	r := make([]int32, 0)
	// 随机一个数换牌
	nRand := util.RndInt(10, 35)
	for i := int32(nRand); i < 40; i++ {
		if cd[i] >= 3 {
			r = append(r, i)
			r = append(r, i)
			r = append(r, i)
			return r
		}
		if cd[i] > 0 && cd[i+1] > 0 && cd[i+2] > 0 {
			r = append(r, i)
			r = append(r, i+1)
			r = append(r, i+2)
			return r
		}
	}
	// 如果随机没有找到 重新找一次
	for i := int32(10); i < 40; i++ {
		if cd[i] >= 3 {
			r = append(r, i)
			r = append(r, i)
			r = append(r, i)
			return r
		}
		if cd[i] > 0 && cd[i+1] > 0 && cd[i+2] > 0 {
			r = append(r, i)
			r = append(r, i+1)
			r = append(r, i+2)
			return r
		}
	}
	return r
}

func (lgc *SCMJRobotLogic) findTwoCard(mapCd map[int32]int32) []int32 {
	r := make([]int32, 0)
	for cd, count := range mapCd {
		if count >= 2 {
			r = append(r, cd)
			r = append(r, cd)
			break
		}
	}
	return r
}

//转化list为牌组
func (lgc *SCMJRobotLogic) list2cardGroup(cg *mahjong.CardGroup, cdlist []int32) *mahjong.CardGroup {
	cg2 := new(mahjong.CardGroup)
	cgcal := new(mahjong.CardGroup)
	for _, cgcd := range cg.Cards {
		cgcal.Cards = append(cgcal.Cards, cgcd)
	}
	for _, cd := range cdlist {
		for _, cgcd := range cgcal.Cards {
			if cd == cgcd/10 {
				cg2.AddCard(cgcd)
				cgcal.DeleteCard(cgcd, 1)
				break
			}
		}
	}
	if len(cdlist) != cg2.Len() {
		mlog.Debug("转化list出错")
	}
	return cg2
}

// 寻找一个和cd相同的牌
func (lgc *SCMJRobotLogic) findSameCard(cg *mahjong.CardGroup, cd int32) bool {
	for _, cgcd := range cg.Cards {
		if cgcd == cd {
			return true
		}
	}
	return false
}
