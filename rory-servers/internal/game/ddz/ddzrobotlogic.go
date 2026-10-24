package ddz

import (
	"fmt"
	"sort"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
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
)

const (
	RobotLogicConfigFileName = "robot"
)

const (
	RobotAIFlag1 = iota + 1
	RobotAIFlag2
	RobotAIFlag3
	RobotAIFlag4
	RobotAIFlag5
	RobotAIFlag6
)

//机器人个性
type RobotCha struct {
	// Data map[string]*rateInfo
	Data map[string]int32
}

// type rateInfo struct {
// 	robmin int
// 	robmax int
// 	betmin int
// 	betmax int
// }

//机器人逻辑
type DDZRobotLogic struct {
	*robotlogic.LogicBase
	robotChas   []*RobotCha
	gameChas    map[int32]*RobotCha
	robotUser   map[int32]*robotUserInfoLocal //机器人拆牌结果
	outInfo     *outCardInfo                  //出牌信息
	userRole    map[int32]*robotRole          //场面局势信息
	sandBoxUser map[int32]*robotUserInfoLocal //黑箱玩家, 提供AI自我演練情境
}

// 拆牌结果
type robotUserInfoLocal struct {
	ctg            *CtGroup           //拆牌得到的牌组
	cardTypesCount map[cardType]int32 //牌组剩余次数
}

// 游戏玩家角色
type robotRole struct {
	bMine    bool  //机器人是否是地主
	bLast    bool  //上家地主
	bNext    bool  //下家地主
	SeatMine int32 //机器人座位号
	SeatLast int32 //上家座位
	SeatNext int32 //下家座位
}

// 玩家出牌信息
type outCardInfo struct {
	outType    int32 //类型 0pass 1出牌
	SeatNo     int32 //座位号
	UserId     int32 //用户id
	CardType   int32 //牌型
	nMin       int32
	nMax       int32
	cg         *poker.CardGroup
	RealSeatNo int32 //有效玩家
}

// 玩家拆牌结构
type CtGroup struct {
	Ct          []*cardTypeGroup //数据
	SingleCards []int32
	hands       int32
	bigs        int32
}

func (cg *CtGroup) String() string {
	return fmt.Sprintf("[%d] %v, %v", cg.hands, cg.Ct, cg.SingleCards)
}

type cardTypeGroup struct {
	byteCardType cardType //牌型
	nMax         int32    //最大值
	nMin         int32    //最小值
}

func (ctg *cardTypeGroup) String() string {
	return fmt.Sprintf("type %d, max %d", ctg.byteCardType, ctg.nMax)
}

// 初始化
func (lgc *DDZRobotLogic) initOutInfo() {
	lgc.outInfo = new(outCardInfo)
	lgc.outInfo.SeatNo = -1
	lgc.outInfo.UserId = 0
	lgc.outInfo.outType = 0
	lgc.outInfo.nMin = 0
	lgc.outInfo.nMax = 0
	lgc.outInfo.RealSeatNo = -1
	lgc.outInfo.cg = new(poker.CardGroup)
	lgc.userRole = make(map[int32]*robotRole, 0)

}

// 创建机器人逻辑
func NewDDZRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *DDZRobotLogic {
	rblgc := new(DDZRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.robotUser = make(map[int32]*robotUserInfoLocal, 0)
	rblgc.sandBoxUser = make(map[int32]*robotUserInfoLocal, 0)
	rblgc.initOutInfo()
	rblgc.gameChas = make(map[int32]*RobotCha, 0)
	return rblgc
}

//解析机器人个性
func (lgc *DDZRobotLogic) parseRobotCha(cfgs *config.Group) {

	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("size", "cn")
	//mlog.Debug("解析机器人配置文件,cn[%d]", cn)
	for i := 1; i <= cn; i++ { //个性 //c1 c2 c3 ...
		cname := fmt.Sprintf("c%d", i)

		rData := new(RobotCha)
		rData.Data = make(map[string]int32, 0)
		rData.Data["danpai"] = robotCfg.GetConfigInt32(cname, "danpai")
		rData.Data["duizi"] = robotCfg.GetConfigInt32(cname, "duizi")
		rData.Data["santiao"] = robotCfg.GetConfigInt32(cname, "santiao")
		rData.Data["shunzi"] = robotCfg.GetConfigInt32(cname, "shunzi")
		rData.Data["shunziplus"] = robotCfg.GetConfigInt32(cname, "shunziplus")
		rData.Data["liandui"] = robotCfg.GetConfigInt32(cname, "liandui")
		rData.Data["lianduiplus"] = robotCfg.GetConfigInt32(cname, "lianduiplus")
		rData.Data["feiji"] = robotCfg.GetConfigInt32(cname, "feiji")
		rData.Data["feijiplus"] = robotCfg.GetConfigInt32(cname, "feijiplus")
		rData.Data["zhadan"] = robotCfg.GetConfigInt32(cname, "zhadan")
		rData.Data["huojian"] = robotCfg.GetConfigInt32(cname, "huojian")
		rData.Data["jiaodizhumax"] = robotCfg.GetConfigInt32(cname, "jiaodizhumax")
		rData.Data["jiaodizhumin"] = robotCfg.GetConfigInt32(cname, "jiaodizhumin")
		rData.Data["jiaodizhup"] = robotCfg.GetConfigInt32(cname, "jiaodizhup")
		rData.Data["jiaodizhuhuojian"] = robotCfg.GetConfigInt32(cname, "jiaodizhuhuojian")
		rData.Data["jiaodizhuzhadan"] = robotCfg.GetConfigInt32(cname, "jiaodizhuzhadan")
		rData.Data["jiaodizhudawang"] = robotCfg.GetConfigInt32(cname, "jiaodizhudawang")
		rData.Data["jiaodizhuxiaowang"] = robotCfg.GetConfigInt32(cname, "jiaodizhuxiaowang")
		rData.Data["jiaodizhu2"] = robotCfg.GetConfigInt32(cname, "jiaodizhu2")
		rData.Data["qiangdizhu"] = robotCfg.GetConfigInt32(cname, "qiangdizhu")
		rData.Data["dizhumingpai"] = robotCfg.GetConfigInt32(cname, "dizhumingpai")
		lgc.robotChas = append(lgc.robotChas, rData)
		//mlog.Debug("danpai=%d", rData.Data["danpai"])
	}
}

//获取机器人个性
func (lgc *DDZRobotLogic) getRobotCha(userid int32) *RobotCha {
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

// 获取游戏逻辑
func (lgc *DDZRobotLogic) getGameLogic() *DDZGameLogic {
	return lgc.GameLogic.(*DDZGameLogic)
}

// 获取发牌器
func (lgc *DDZRobotLogic) getPoker() *poker.FixedPoker {
	return lgc.getGameLogic().poker.(*poker.FixedPoker)
}

//游戏消息
func (lgc *DDZRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {

	if robot == nil {
		mlog.Debug("機器人是nil")
		return
	}

	if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyDiZhuInfo) {
		lgc.getGameLogic().GameLog.AddGameStep("机器人ID:[%d],智商:[%d]", robot.UserInfo.UserID, robot.RebotLevel)
		//牌局紀錄 機器人智商
		lgc.getGameLogic().record.setRobotLevel(robot.UserInfo.UserID, robot.RebotLevel)
	}

	if m.SClassID == int32(netproto.DDZ_GameMessageClassID_BeforeFapai) {
		lgc.getGameLogic().ExchangeMaxCard(robot.UserInfo.UserID)
	}
	//mlog.System("[%v]机器人收到消息[%d]***************************gamestep[%d]", robot.UserInfo.UserID, m.SClassID, lgc.getGameLogic().gamestep)
	lgc.receiveChupai(robot, m)
	//叫地主
	if m.SClassID == int32(netproto.DDZ_GameMessageClassID_ChuPaiRet) {
		//		data := m.MsgData.(*netproto.DDZ_ChuPaiRet)
		//mlog.Debug("出牌返回%v", data)
	}
	if lgc.canJiaodizhu(robot, m) {
		//mlog.Debug("轮到机器人叫地主")
		lgc.robotJiaoDizhu(robot)
		return
	}
	//抢地主
	if lgc.canRobDizhu(robot, m) {
		//mlog.Debug("轮到机器人抢地主")
		lgc.robotRobDizhu(robot)
		return
	}
	//不明牌DDZ_MingPai
	if lgc.canMinPai(robot, m) {
		//mlog.Debug("机器人 明牌")
		lgc.robotMingPai(robot)
		return
	}
	//收到3地主牌

	//出牌,出牌错误后点击 pass
	if lgc.canChupai(robot, m) {
		lgc.updateRoleInfo(robot)
		//if robot.RebotLevel < 6 {	//TODO 機器人是地主進入AI, robot.RebotLevel < 6
		lgc.addSandBoxAllUser()
		up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
		hands := lgc.sandBoxUser[robot.UserInfo.UserID].ctg.hands
		_, dizhu := lgc.getGameLogic().ext.getDiZhu()                                          //地主user
		if checkAIWithFlag(robot, RobotAIFlag2) && (hands <= 2 || len(up.Cards.Cards) == 20) { //TODO 機器人是地主進入AI, 手牌 < 6
			lgc.robotAction(robot)
		} else {
			dz := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
			if checkAIWithFlag(robot, RobotAIFlag5) && (dz.IsRobot() && dizhu.userId != robot.UserInfo.UserID) {
				lgc.robotAction(robot)
			} else {
				lgc.robotAIAction(robot)
			}
		}

		return
	}
	//接收玩家出牌消息DDZ_GameMessageClassID_NotifyChuPai

}

// 结束游戏
func (lgc *DDZRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
	lgc.robotUser = make(map[int32]*robotUserInfoLocal)
	lgc.sandBoxUser = make(map[int32]*robotUserInfoLocal)
	lgc.initOutInfo()
}

// 机器人是否可以叫地主
func (lgc *DDZRobotLogic) canJiaodizhu(robot *robotcomm.Robot, m *msg.Message) bool {
	if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameTimeOut) && lgc.getGameLogic().gamestep == gameProgressJiaoDiZhu {
		if m.MsgData == nil {
			mlog.Error("data is nil")
			return false
		}
		data := m.MsgData.(*netproto.DDZ_NotifyGameTimeOut)

		if *data.SeatNo == robot.UserInfo.SeatNo {
			return true
		}
	}
	return false
}

//是否可以打开明牌
func (lgc *DDZRobotLogic) canMinPai(robot *robotcomm.Robot, m *msg.Message) bool {
	if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameTimeOut) && lgc.getGameLogic().gamestep == gameProgressShow {
		if m.MsgData == nil {
			mlog.Error("data is nil")
			return false
		}
		data := m.MsgData.(*netproto.DDZ_NotifyGameTimeOut)
		if *data.SeatNo == robot.UserInfo.SeatNo {
			return true
		}
	}
	return false
}

// 是否可以抢地主
func (lgc *DDZRobotLogic) canRobDizhu(robot *robotcomm.Robot, m *msg.Message) bool {
	if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameTimeOut) && lgc.getGameLogic().gamestep == gameProgressRob {
		if m.MsgData == nil {
			mlog.Error("data is nil")
			return false
		}
		data := m.MsgData.(*netproto.DDZ_NotifyGameTimeOut)
		if *data.SeatNo == robot.UserInfo.SeatNo {
			return true
		}
	}
	return false
}

//收出牌消息
func (lgc *DDZRobotLogic) receiveChupai(robot *robotcomm.Robot, m *msg.Message) {

	//mlog.Debug("机器人收到用户出牌信息")
	if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyChuPai) && lgc.getGameLogic().gamestep == gameProgressStart {
		if m.MsgData == nil {
			mlog.Error("data is nil")
			return
		}
		// lgc.removeTuoguan(robot.UserInfo)
		data := m.MsgData.(*netproto.DDZ_NotifyChuPai)
		if data.GetMainCards() == 0 && data.GetMainCardsNumber() == 0 {
			lgc.outInfo.outType = 0
		} else {
			lgc.outInfo.RealSeatNo = data.GetSeatNo()
			lgc.outInfo.outType = 1
			lgc.outInfo.cg = new(poker.CardGroup)
			for _, cd := range data.Cards.Cards {
				cdtmp := new(poker.Card)
				cdtmp.Color = cd.GetColor()
				cdtmp.CardNumber = cd.GetNumber()
				lgc.outInfo.cg.AddCard(cdtmp)
			}
			lgc.outInfo.nMin, lgc.outInfo.nMax = lgc.analysisCards(lgc.outInfo.cg, data.GetCardType())
			lgc.outInfo.SeatNo = data.GetSeatNo()
			lgc.outInfo.CardType = data.GetCardType()
		}

	}
}

//出牌
func (lgc *DDZRobotLogic) canChupai(robot *robotcomm.Robot, m *msg.Message) bool {
	if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameTimeOut) && lgc.getGameLogic().gamestep == gameProgressStart {
		if m.MsgData == nil {
			mlog.Error("data is nil")
			return false
		}
		data := m.MsgData.(*netproto.DDZ_NotifyGameTimeOut)
		if *data.SeatNo == robot.UserInfo.SeatNo {
			mlog.NewDefaultLogger().SetCallDepLen(10).SetData(fmt.Sprintf("機器人出牌 A [%v]", robot.UserInfo.UserID)).Debug("主动出牌")
			return true
		}
		u := lgc.getGameLogic().ext.getExtInfoByTurn(lgc.getGameLogic().ext.turn)
		mlog.Debug("userId[%v] 出牌 u.isTuoGuan[%v] seatno[%v], userseat[%v] ",
			robot.UserInfo.UserID, u.isTuoGuan, u.userSeat, robot.UserInfo.SeatNo)
		if robot.UserInfo.SeatNo != u.userSeat {
			mlog.Error("userId[%v] 非法操作，不该您出牌 u.isTuoGuan[%v] seatno %v, userseat %v ", robot.UserInfo.UserID, u.isTuoGuan, u.userSeat, robot.UserInfo.SeatNo)
			return false
		}
	}

	return false
}

// 机器人抢地主
func (lgc *DDZRobotLogic) robotRobDizhu(robot *robotcomm.Robot) {
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		mlog.Error("系统错误robotJiaoDizhu，获取玩家牌型失败")
		return
	}
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_QiangDiZhu)

	btdata := new(netproto.DDZ_QiangDiZhu)
	if lgc.robotConsinderJiaoDizhu(up, 2) {
		//叫地主
		//mlog.Debug("机器人操作 抢地主")
		btdata.IsRob = proto.Bool(true)
		btmsg.MsgData, _ = proto.Marshal(btdata)
		lgc.SendGameMessage(robot.UserInfo, btmsg)
	} else {
		mlog.Debug("机器人[%v]操作 不抢地主", robot.UserInfo.UserID)
		btdata.IsRob = proto.Bool(false)
		btmsg.MsgData, _ = proto.Marshal(btdata)
		lgc.SendGameMessage(robot.UserInfo, btmsg)
	}
	//w := lgc.CalBigWeigth(up.Cards, robot.UserInfo.UserID)
	//if lgc.isCanRob(w, robot.UserInfo.UserID) {
	//	//叫地主
	//	mlog.Debug("机器人操作 抢地主")
	//	btdata.IsRob = proto.Bool(true)
	//
	//	btmsg.MsgData, _ = proto.Marshal(btdata)
	//	lgc.SendGameMessage(robot.UserInfo, btmsg)
	//} else {
	//	//叫地主
	//	mlog.Debug("机器人操作 不抢地主")
	//
	//	btdata.IsRob = proto.Bool(false)
	//
	//	btmsg.MsgData, _ = proto.Marshal(btdata)
	//	lgc.SendGameMessage(robot.UserInfo, btmsg)
	//}
}
func (lgc *DDZRobotLogic) robotConsinderJiaoDizhu(up *poker.UserPoker, d int) bool { //d 1是叫地主, 2是搶地主

	if conf.TestMode { //測試時機器人不叫地主
		return false
	}

	dc := NewDefineCard(lgc.getGameLogic())
	userCard := dc.Combination(up.Cards)
	weight, hands, _ := dc.CalWeight(userCard, userCard.SingleCards)
	//mlog.Debug("userid %d, weight %d, hands %d, poker %v", up.User.UserID, weight, hands, up.Cards)
	if dc.GameLogic.GetServerConfig().PokerType == gamelogic.PokerTypeShuffle {
		if lgc.getGameLogic().GetTable().TableRealRobotCount() >= 2 && weight > 20 && hands < 8 {
			return true
		} else {
			return false
		}
	} else {
		w := lgc.CalBigWeigth(up.Cards, up.User.UserID)
		//mlog.System("手牌權值 %d", w)
		jiaodizhu := false
		if d == 1 && lgc.isCanjiaodizhu(w, up.User.UserID) {
			jiaodizhu = true

		} else if d == 2 && lgc.isCanRob(w, up.User.UserID) {
			jiaodizhu = true
		}
		return jiaodizhu
	}
}

// 机器人叫地主
func (lgc *DDZRobotLogic) robotJiaoDizhu(robot *robotcomm.Robot) {
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		mlog.Error("系统错误robotJiaoDizhu，获取玩家牌型失败")
		return
	}
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_JiaoDiZhu)

	btdata := new(netproto.DDZ_JiaoDiZhu)

	if lgc.robotConsinderJiaoDizhu(up, 1) {
		//叫地主
		btdata.IsJiaoDiZhu = proto.Bool(true)
		btmsg.MsgData, _ = proto.Marshal(btdata)
		lgc.SendGameMessage(robot.UserInfo, btmsg)
	} else {
		//mlog.Debug("机器人操作 不叫地主")
		btdata.IsJiaoDiZhu = proto.Bool(false)
		btmsg.MsgData, _ = proto.Marshal(btdata)
		lgc.SendGameMessage(robot.UserInfo, btmsg)
	}

	//w := lgc.CalBigWeigth(up.Cards, robot.UserInfo.UserID)
	//if lgc.isCanjiaodizhu(w, robot.UserInfo.UserID) {
	//	//叫地主
	//	mlog.Debug("机器人操作 叫地主")
	//	btdata.IsJiaoDiZhu = proto.Bool(true)
	//	btmsg.MsgData, _ = proto.Marshal(btdata)
	//	lgc.SendGameMessage(robot.UserInfo, btmsg)
	//} else {
	//	//叫地主
	//	mlog.Debug("机器人操作 不叫地主")
	//	btdata.IsJiaoDiZhu = proto.Bool(false)
	//
	//	btmsg.MsgData, _ = proto.Marshal(btdata)
	//	lgc.SendGameMessage(robot.UserInfo, btmsg)
	//}
}

//转化list为牌组
func (lgc *DDZRobotLogic) list2cardGroup(cg *poker.CardGroup, cdlist []int32) *poker.CardGroup {
	mlog.Debug("robot 出牌[%v]", cdlist)
	cg2 := new(poker.CardGroup)
	cgcal := new(poker.CardGroup)
	for _, cgcd := range cg.Cards {
		cd := new(poker.Card)
		cd.Color = cgcd.Color
		cd.CardNumber = cgcd.CardNumber
		cgcal.Cards = append(cgcal.Cards, cd)
	}
	for _, cd := range cdlist {
		for _, cgcd := range cgcal.Cards {
			if cd == cgcd.CardNumber {
				cg2.AddCard(cgcd)
				cgcal.DeleteCard(cgcd)
				break
			}
		}
	}
	return cg2
}

// 更新角色信息
func (lgc *DDZRobotLogic) updateRoleInfo(robot *robotcomm.Robot) {
	roleInfo := new(robotRole)
	SeatNoDizhu := lgc.getGameLogic().ext.lastQiangDiZhu
	//本家身份
	roleInfo.SeatMine = robot.UserInfo.SeatNo
	if robot.UserInfo.SeatNo == SeatNoDizhu {
		//mlog.System("地主机器人操作 robot[%v]", robot.UserInfo.UserID)
		roleInfo.bMine = true
	} else {
		roleInfo.bMine = false
		//mlog.System("农民机器人操作 robot[%v]", robot.UserInfo.UserID)
	}

	//上家身份
	roleInfo.SeatLast = lgc.outInfo.SeatNo
	if lgc.outInfo.SeatNo == SeatNoDizhu {
		//mlog.System("上家上地主 , robot[%v]", robot.UserInfo.UserID)
		roleInfo.bLast = true
	} else {
		roleInfo.bLast = false
		//mlog.System("上家是农民 , robot[%v]", robot.UserInfo.UserID)
	}
	//下家身份
	SeatNoNext := int32(0)
	if robot.UserInfo.SeatNo < 2 {
		SeatNoNext = robot.UserInfo.SeatNo + 1
	} else {
		SeatNoNext = 0
	}
	roleInfo.SeatNext = SeatNoNext

	user := lgc.getGameLogic().ext.getExtInfoByTurn(SeatNoNext)
	if user != nil {
		if user.isDiZhu {
			//mlog.System("下家是地主 user[%v]", user.userId)
			roleInfo.bNext = true
			roleInfo.bLast = false
		} else {
			//mlog.System("下家是农民 user[%v]", user.userId)
			roleInfo.bNext = false
			roleInfo.bLast = true
		}
	}
	lgc.userRole[robot.UserInfo.UserID] = roleInfo
}

//取消托管
func (lgc *DDZRobotLogic) removeTuoguan(ur *gamelogic.GameUser) {
	//DDZ_GameMessageClassID_TuoGuan
	tuoguan := lgc.getGameLogic().ext.getExtInfoByTurn(ur.SeatNo)
	if tuoguan != nil {
		if tuoguan.isTuoGuan {
			// tuoguan.isTuoGuan = false
			mlog.Debug("机器人取消托管,SeatNo[%d]", ur.SeatNo)
			btmsg := new(msg.Message)
			btmsg.BClassID = int32(netproto.MessageBClassID_Game)
			btmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_TuoGuan)

			btdata := new(netproto.DDZ_TuoGuan)
			btmsg.MsgData, _ = proto.Marshal(btdata)
			lgc.SendGameMessage(ur, btmsg)
		}
	}
}

//机器人操作
func (lgc *DDZRobotLogic) robotAction(robot *robotcomm.Robot) {
	//lgc.updateRoleInfo(robot)
	roleInfo := lgc.userRole[robot.UserInfo.UserID]
	if roleInfo == nil {
		lgc.getGameLogic().GameLog.AddGameStep("robotPass4 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAction 取不到場面局勢")
		lgc.robotPass(robot)
		return
	}
	mlog.Debug("一般Robot操作 userid[%v] 上家出牌座位 %v , 本身座位[%v]", robot.UserInfo.UserID, lgc.outInfo.RealSeatNo, robot.UserInfo.SeatNo)

	mlog.Debug("robotAction userid[%v] seatno[%v] 本身地主[%v] 上家地主[%v] 下家地主[%v] 本身座位[%v] 上家座位[%v] 下家座位[%v] ",
		robot.UserInfo.UserID, robot.UserInfo.SeatNo,
		roleInfo.bMine, roleInfo.bLast, roleInfo.bNext, roleInfo.SeatMine, roleInfo.SeatLast, roleInfo.SeatNext)

	//主动出牌
	if lgc.outInfo.RealSeatNo == robot.UserInfo.SeatNo || lgc.outInfo.RealSeatNo == -1 {
		dc := NewDefineCard(lgc.GameLogic.(*DDZGameLogic))
		foo := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
		cardtype, ctg := dc.ActiveChuPai(foo.Cards, robot.UserInfo.UserID, robot)
		lgc.robotUser[robot.UserInfo.UserID] = ctg
		ucg := foo.Cards
		ocList := lgc.outCardByType(cardtype, robot.UserInfo.UserID)
		mlog.Debug("user %d, 主动出牌 %v, cardtype %d", robot.UserInfo.UserID, ocList, cardtype)
		mlog.NewDefaultLogger().SetCallDepLen(10).SetData(fmt.Sprintf("user %d, 出牌 %v, cardtype %d", robot.UserInfo.UserID, ocList, cardtype)).Debug("主动出牌")
		ocg := lgc.list2cardGroup(ucg, ocList)
		lgc.S2COutCard(ocg, robot.UserInfo)
		return
	}
	//跟牌
	if lgc.outInfo.RealSeatNo != -1 && lgc.outInfo.RealSeatNo != robot.UserInfo.SeatNo {
		mlog.Warn("机器人[%d]跟牌型[%d] 點數[%v]", robot.UserInfo.UserID, lgc.outInfo.CardType, lgc.outInfo.nMax)
		dc := NewDefineCard(lgc.GameLogic.(*DDZGameLogic))
		up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
		ctg := dc.Combination(up.Cards)
		rInfo := new(robotUserInfoLocal)
		rInfo.cardTypesCount = make(map[cardType]int32, 0)
		rInfo.ctg = ctg
		for _, data := range ctg.Ct {
			rInfo.cardTypesCount[data.byteCardType]++
		}
		rInfo.cardTypesCount[cardTypeSingle] += int32(len(ctg.SingleCards))
		lgc.robotUser[robot.UserInfo.UserID] = rInfo
		//mlog.Debug("=================== ctg.bigs %d", lgc.robotUser[robot.UserInfo.UserID].ctg.bigs)

		//TODO 當地主只剩兩手牌, 若機器人手上牌數超過3手,並且手上沒有最大炸彈, 則pass
		passTwoHand := false
		passTwoHand = lgc.passTwoHands(ctg, dc, robot, passTwoHand, roleInfo)
		if checkAIWithFlag(robot, RobotAIFlag3) && passTwoHand {
			//當地主手上沒單張, 只剩炸彈, 不能跟單張,這是陷阱.
			lastUr := lgc.getGameLogic().GetUserBySeatNo(lgc.outInfo.RealSeatNo)
			lastUrPoker := lgc.getGameLogic().gamePoker.GetUserPoker(lastUr.UserID)
			lastCtg := dc.Combination(lastUrPoker.Cards)
			if lgc.outInfo.CardType == int32(cardTypeSingle) && len(lastCtg.SingleCards) != 0 {
				lgc.FollowCards(robot, cardTypeSingle)
				return
			}
			lgc.getGameLogic().GameLog.AddGameStep("robotPass5 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAction checkAIWithFlag(robot, RobotAIFlag3) 跟牌，當地主只剩兩手牌, 若機器人手上牌數超過3手,並且手上沒有最大炸彈, 則pass")
			lgc.robotPass(robot)
			//mlog.Fatal("stop")
			return
		}

		//TODO 機器人是地主, 當農民只剩2手牌, 地主手上牌數超過3手,並且手上沒有最大炸彈, 則pass
		lgc.funcName(roleInfo, robot, ctg, dc, passTwoHand)

		//倒数第二手時的出牌邏輯
		if ctg.hands < 3 && ctg.bigs > 0 {
			//mlog.Debug("用戶 %d 倒數第二手 %d, 大牌數 %d, %s", robot.UserInfo.UserID, ctg.hands, ctg.bigs, ctg.String())
			for _, ct := range ctg.Ct {
				if cardTypeRocket == ct.byteCardType {
					cdlist := lgc.outRocket(robot.UserInfo.UserID)
					if len(cdlist) > 0 {
						cg := lgc.list2cardGroup(up.Cards, cdlist)
						lgc.S2COutCard(cg, robot.UserInfo)
						return
					}
				}

				if checkAIWithFlag(robot, RobotAIFlag4) {
					if cardTypeBomb == ct.byteCardType {
						if dc.GameLogic.isBiggestCT(ct, robot.UserInfo.UserID) || ctg.hands == 1 {
							mlog.Debug("倒數第二手時有場上最大的牌型 %d", ct.byteCardType)
							mlog.NewDefaultLogger().SetCallDepLen(8).Debug("倒數第二手時有場上最大的牌型 userid[%v]", robot.UserInfo.UserID)
							if ok, cg := lgc.FollowCardsBomb(robot); ok {
								lgc.S2COutCard(cg, robot.UserInfo)
								return
							}
						}
					}
				}
			}
		}

		buchu := false
		crazy := false
		glc := lgc.getGameLogic()

		//機器人邏輯, 自己是農民, 地主也是機器人, 機器人輔助地主, 儘量不出
		if checkAIWithFlag(robot, RobotAIFlag6) {
			if !roleInfo.bMine {
				_, dizhu := glc.ext.getDiZhu()
				if dizhu != nil {
					if ur := glc.GetUserByUserID(dizhu.userId); ur != nil && ur.IsRobot() {
						mlog.Debug("用戶 %d, 地主也是機器人", robot.UserInfo.UserID)
						cn := util.RndInt(0, 10) + 1
						if int32(cn) <= 9 { //九成不出, 一成瘋狂出
							mlog.Debug("用戶 %d, 被動出牌不出", robot.UserInfo.UserID)
							buchu = true
						} else {
							mlog.Debug("用戶 %d, 被動出牌瘋狂出", robot.UserInfo.UserID)
							crazy = true
						}
					}
				}
			}
		}
		//检测上家出牌,上家不是地主自己不是地主,上家出牌尽量不跟,不是瘋狂機器人
		//或者對機器人莊家儘量不押
		robot.BuchuBombCard = false
		mlog.Debug("用戶 %d, %v, %v ",
			robot.UserInfo.UserID, roleInfo.bLast == roleInfo.bMine, lgc.outInfo.RealSeatNo == roleInfo.SeatLast)

		_, dizhu := lgc.getGameLogic().ext.getDiZhu()
		dizhuUr := lgc.getGameLogic().GetUserByUserID(dizhu.userId)

		mlog.Debug("robot[%v] 地主是機器人[%v], 農民機器人才不壓  seatno[%v] RealSeatNo[%v] ",
			robot.UserInfo.UserID, dizhuUr.IsRobot(), dizhuUr.SeatNo, lgc.outInfo.RealSeatNo)
		//地主是機器人, 農民機器人才不壓
		if checkAIWithFlag(robot, RobotAIFlag6) {
			if (dizhuUr.IsRobot() && !roleInfo.bMine && roleInfo.bLast == roleInfo.bMine && lgc.outInfo.RealSeatNo == roleInfo.SeatLast && !crazy) || buchu {
				mlog.Debug("robot[%v] 儘量不押 bLast[%v] mine[%v] crazy[%v] buchu[%v] RealSeatNo[%v] SeatLast[%v]",
					robot.UserInfo.UserID, roleInfo.bLast, roleInfo.bMine, !crazy, buchu, lgc.outInfo.RealSeatNo, roleInfo.SeatLast)
				if checkAIWithFlag(robot, RobotAIFlag6) {
					robot.BuchuBombCard = true
				}
				if lgc.outInfo.CardType == int32(cardTypeDouble) || lgc.outInfo.CardType == int32(cardTypeSingle) {
					up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
					if up != nil {
						lgc.recombination(robot.UserInfo.UserID, up.Cards)
					}

					if (lgc.outInfo.nMax > 2 && lgc.outInfo.nMax < 12) || lgc.getUserHand(robot.UserInfo.UserID) <= 2 {
						//去跟牌
					} else {
						if !lgc.outNeverSingle(robot.UserInfo.UserID) {
							mlog.Debug("上次出牌的玩家为友方，PASS")
							lgc.getGameLogic().GameLog.AddGameStep("robotPass6 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAction checkAIWithFlag(robot, RobotAIFlag6) 跟牌，地主是機器人, 農民機器人才不壓")
							lgc.robotPass(robot)
							return
						}
					}
				} else {
					if !lgc.outNeverSingle(robot.UserInfo.UserID) {
						//mlog.Debug("上次出牌的玩家为友方，PASS")
						lgc.getGameLogic().GameLog.AddGameStep("robotPass7 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAction checkAIWithFlag(robot, RobotAIFlag6) 跟牌，地主是機器人, 農民機器人才不壓")
						lgc.robotPass(robot)
						return
					}
				}
			}
		}

		if lgc.outInfo.CardType == int32(cardTypeRocket) {
			lgc.FollowCards(robot, cardTypeRocket)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypeStraight) {
			lgc.FollowCards(robot, cardTypeStraight)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypeThreeStraight) {
			lgc.FollowCards(robot, cardTypeThreeStraight)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypeDouble) {
			lgc.FollowCards(robot, cardTypeDouble)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypeSingle) {
			lgc.FollowCards(robot, cardTypeSingle)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypeThree) || lgc.outInfo.CardType == int32(cardTypeThreeAndOne) {
			lgc.FollowCards(robot, cardTypeThree)
			return
		}

		if lgc.outInfo.CardType == int32(cardTypeDoubleStraights) {
			lgc.FollowCards(robot, cardTypePlaneAndWing)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypeDoubleStraightd) {
			lgc.FollowCards(robot, cardTypePlaneAndWing)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypePlaneAndWing) {
			lgc.FollowCards(robot, cardTypePlaneAndWing)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypeFourAndTwod) {
			lgc.FollowCards(robot, cardTypePlaneAndWing)
			return
		}
		if lgc.outInfo.CardType == int32(cardTypeFourAndTwos) {
			lgc.FollowCards(robot, cardTypePlaneAndWing)
			return
		}

		//消極狀態下找不到牌型不要用炸彈壓
		if robot.BuchuBombCard {
			mlog.Debug("userid[%v] 消極狀態下找不到牌型不要用炸彈壓, 手牌還有 %v ", robot.UserInfo.UserID, up.Cards)
			lgc.getGameLogic().GameLog.AddGameStep("robotPass8 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAction 跟牌，消極狀態下找不到牌型不要用炸彈壓")
			lgc.robotPass(robot)
		} else {
			//跟牌時當機器人手牌數比真人大就不出
			if lgc.robotHandMoreByPass(robot) || ctg.hands > 5 {
				mlog.Debug("userid[%v] 跟牌時當機器人手牌數比真人大就不出, 手牌數[%v], 手牌還有 %v",
					robot.UserInfo.UserID, ctg.hands, up.Cards)
				lgc.getGameLogic().GameLog.AddGameStep("robotPass9 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "robotAction 跟牌，跟牌時當機器人手牌數比真人大就不出")
				lgc.robotPass(robot)
			} else {
				lgc.FollowCards(robot, cardTypeBomb)
			}
		}

	}
	return
}

func (lgc *DDZRobotLogic) funcName(roleInfo *robotRole, robot *robotcomm.Robot, ctg *CtGroup, dc *DefineCard, passTwoHand bool) {

}

//當地主只剩兩手牌, 若機器人手上牌數超過3手,並且手上沒有最大炸彈, 則pass
//機器人是地主, 當農民只剩2手牌, 地主手上牌數超過3手,並且手上沒有最大炸彈, 則pass
func (lgc *DDZRobotLogic) passTwoHands(ctg *CtGroup, dc *DefineCard, robot *robotcomm.Robot, passTwoHand bool, roleInfo *robotRole) bool {
	if ctg.hands >= 3 {
		isBiggtest := false
		for _, ct := range ctg.Ct { //機器人手牌有最大牌
			if dc.GameLogic.isBiggestCT(ct, robot.UserInfo.UserID) {
				isBiggtest = true
				break
			}
		}

		_, dizhu := lgc.getGameLogic().ext.getDiZhu()
		uHands := lgc.getGameLogic().getUserPokerHands(dizhu.userId)
		dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)

		mlog.Debug("check 地主沒有組牌只剩單牌, 地主手牌數[%v] sigLen[%v]", uHands, len(lgc.sandBoxUser[dizhu.userId].ctg.SingleCards))
		if !dizhuUser.IsRobot() && uHands == 1 && dizhu.userId != robot.UserInfo.UserID &&
			len(lgc.sandBoxUser[dizhu.userId].ctg.SingleCards) == 1 && noGroupCard(lgc.sandBoxUser[dizhu.userId].ctg.Ct) {
			mlog.Debug("robot[%v] 地主[%v]是真人, 地主沒有組牌只剩單牌. 機器人要從最大單牌出", robot.UserInfo.UserID, dizhu.userId)
			passTwoHand = false
			return passTwoHand
		}

		mlog.Debug("[農民] if 地主userid[%v]只剩兩手牌, 出牌robot[%v]手上有沒有最大炸彈, 則pass, 機器人有無最大炸彈[%v] ",
			dizhu.userId, robot.UserInfo.UserID, isBiggtest)
		if uHands <= 2 && !isBiggtest && dizhu.userId != robot.UserInfo.UserID {
			mlog.Debug("[農民] robot[%v] pass . 真人[%v]手牌[%v], robot手牌數[%v]", robot.UserInfo.UserID, dizhu.userId, uHands,
				ctg.hands)
			passTwoHand = true
			//mlog.Fatal("stop")
			return passTwoHand
		}
	}

	if roleInfo.bMine { //當前機器人是地主
		twoHand := false
		for _, u := range lgc.getGameLogic().AllGameUser() {
			if u.UserID != robot.UserInfo.UserID {
				hands := lgc.getGameLogic().getUserPokerHands(u.UserID)
				if hands <= 2 {
					twoHand = true
				}
			}
		}

		isBiggtest := false
		for _, ct := range ctg.Ct {
			if dc.GameLogic.isBiggestCT(ct, robot.UserInfo.UserID) {
				isBiggtest = true
				break
			}
		}

		//最後一手要出牌
		if ctg.hands == 1 {
			return false
		}

		//mlog.Debug("[地主] robot[%v]是地主. 當農民只剩2手牌, 手牌[%v], robot手牌數[%v] isBiggtest[%v] pass",
		//	robot.UserInfo.UserID, twoHand, ctg.hands, isBiggtest)
		if twoHand && !isBiggtest {
			mlog.Debug("[地主] robot[%v]是地主. 當農民只剩2手牌, 手牌[%v], robot手牌數[%v] pass",
				robot.UserInfo.UserID, twoHand, ctg.hands)
			//mlog.Fatal("stop")
			passTwoHand = true
			return passTwoHand
		}
	}

	return passTwoHand
}

//跟牌時當機器人手牌數比真人大就不出
func (lgc *DDZRobotLogic) robotHandMoreByPass(robot *robotcomm.Robot) bool {

	ispass := false
	robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
	//跟牌時當機器人手牌數比真人少才使用炸彈
	seatlast := lgc.userRole[robot.UserInfo.UserID].SeatLast
	user := lgc.getGameLogic().GetUserBySeatNo(seatlast)
	if user == nil {
		return false
	}

	seatLastHands := lgc.getGameLogic().getUserPokerHands(user.UserID)
	mlog.Debug("userid[%v] 上家 手牌數[%v], robot手牌數[%v}", user.UserID, seatLastHands, robotHands)

	if robotHands < seatLastHands {
		ispass = false
	} else {
		ispass = true
	}

	//機器人剩下兩手牌就出牌
	rup := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	mlog.Debug("userid[%v] 機器人剩下炸彈就出牌, 手牌數[%v] , ispass[%v], 手牌還有 %v", robot.UserInfo.UserID, robotHands, ispass, rup.Cards)
	if robotHands <= 2 {
		ispass = false
		example.invalid("userid[%v] 機器人剩下炸彈就出牌, 手牌數[%v] noPass, 手牌還有 %v", robot.UserInfo.UserID, robotHands, rup.Cards)
	}

	return ispass
}

// 机器人选择pass
func (lgc *DDZRobotLogic) robotPass(robot *robotcomm.Robot) {
	mlog.Debug("userid[%v] 机器人操作 Pass", robot.UserInfo.UserID)
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_Pass)

	btdata := new(netproto.DDZ_Pass)

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessage(robot.UserInfo, btmsg)
}

//明牌
func (lgc *DDZRobotLogic) robotMingPai(robot *robotcomm.Robot) {
	mlog.Debug("机器人操作 明牌 userid[%v]", robot.UserInfo.UserID)
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_MingPai)

	btdata := new(netproto.DDZ_MingPai)

	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		btdata.IsMingPai = proto.Bool(false)
		btmsg.MsgData, _ = proto.Marshal(btdata)
		lgc.SendGameMessage(robot.UserInfo, btmsg)
		return
	}
	w := lgc.CalBigWeigth(up.Cards, robot.UserInfo.UserID)
	if lgc.isMingpai(w, robot.UserInfo.UserID) {
		btdata.IsMingPai = proto.Bool(true)
	} else {
		btdata.IsMingPai = proto.Bool(false)
	}

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessage(robot.UserInfo, btmsg)
}

//机器人出牌,主动出牌
func (lgc *DDZRobotLogic) robotOutCards(robot *robotcomm.Robot) []int32 {
	//机器人出牌
	outCards := make([]int32, 0)
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		mlog.Error("系统错误robotOutCards，获取玩家牌型失败")
		return outCards
	}
	lgc.recombination(robot.UserInfo.UserID, up.Cards)

	robotUser := lgc.robotUser[robot.UserInfo.UserID]
	var outCardType cardType = cardTypeNone

	totalHand := int32(0)
	for _, count := range robotUser.cardTypesCount {
		if count <= 0 {
			continue
		}
		totalHand += count
	}

	//如果敌人只剩1张牌 不出单牌
	//mlog.Debug("总手数[%d]", totalHand)
	for cardtype, count := range robotUser.cardTypesCount {
		if count <= 0 {
			continue
		}
		if lgc.outNeverSingle(robot.UserInfo.UserID) && cardtype == cardTypeSingle {
			continue
		}
		if len(robotUser.ctg.SingleCards) <= 0 && cardtype == cardTypeSingle {
			continue
		}
		if up.Cards.Len() > 3 && cardtype == cardTypeRocket {
			continue
		}
		if totalHand > 2 && cardtype == cardTypeBomb {
			continue
		}
		outCardType = cardtype
		break
	}
	// if outCardType == cardTypeNone {
	// 	outCardType = cardTypeSingle
	// }
	cgList := lgc.outCardByType(outCardType, robot.UserInfo.UserID)
	if len(cgList) < 1 || len(cgList) == 3 {
		lgc.recombination(robot.UserInfo.UserID, up.Cards)
		suffleMap := make(map[cardType]int32, 0)
		for _, dataTmp := range robotUser.ctg.Ct {
			if dataTmp.byteCardType == cardTypeNone {
				continue
			}
			suffleMap[dataTmp.byteCardType]++
		}

		_, cgList = lgc.outCardforever(suffleMap, cgList, robot.UserInfo.UserID, 0)
	}
	if lgc.outNeverSingle(robot.UserInfo.UserID) && len(cgList) <= 1 {
		//选择一张最大的手牌
		//mlog.Debug("机器人选中一张最大的手牌")
		maxList := make([]int32, 0)
		for _, cdMax := range up.Cards.Cards {
			maxList = append(maxList, cdMax.CardNumber)
		}
		maxList = lgc.sort(maxList)
		lenMax := int32(len(maxList))
		cgList = make([]int32, 0)
		maxOut := maxList[lenMax-1]
		if maxOut <= 14 && maxList[0] == 2 {
			maxOut = 2
		}
		cgList = append(cgList, maxOut)
	}
	return cgList
}

//机器人选择出牌
func (lgc *DDZRobotLogic) outCardforever(suffleMap map[cardType]int32, cgList []int32, userid, count int32) (a int32, b []int32) {
	count += 1
	up := lgc.getGameLogic().gamePoker.GetUserPoker(userid)
	if up == nil {
		mlog.Error("系统错误robotOutCards，获取玩家牌型失败")
		return count, cgList
	}
	var outCardType cardType = cardTypeSingle
	for t, count := range suffleMap {
		if count <= 0 {
			continue
		}
		if lgc.outNeverSingle(userid) && t == cardTypeSingle {
			continue
		}
		if up.Cards.Len() > 3 && t == cardTypeRocket {
			continue
		}
		outCardType = t
		break
	}
	cgList = lgc.outCardByType(outCardType, userid)

	if count > 20 || len(cgList) > 0 {
		return count, cgList
	} else {
		return lgc.outCardforever(suffleMap, cgList, userid, count)
	}
	// return count, cgList
}

//选择最小单牌
func (lgc *DDZRobotLogic) outMinSingle(robot *robotcomm.Robot) bool {
	//如果下家是友方且只剩一张牌，出最小单牌
	roleInfo := lgc.userRole[robot.UserInfo.UserID]
	if roleInfo != nil {
		if roleInfo.bNext == roleInfo.bMine {
			if roleInfo.SeatNext < 0 {
				return false
			}
			userNext := lgc.getGameLogic().ext.getExtInfoByTurn(roleInfo.SeatNext)
			if userNext != nil {
				// upNext.userId
				upNext := lgc.getGameLogic().gamePoker.GetUserPoker(userNext.userId)
				if upNext != nil {
					if upNext.Cards.Len() < 2 {
						//mlog.Debug("友方只剩一张单牌，机器人尽量出最小单")
						return true
					}
				}
			}
		}
	}
	return false
}

////如果敌人只剩1张牌 不出单牌
func (lgc *DDZRobotLogic) outNeverSingle(userid int32) bool {
	roleInfo := lgc.userRole[userid]
	if roleInfo != nil {
		if roleInfo.bMine {
			for i := int32(0); i <= 2; i++ {
				if i == roleInfo.SeatMine {
					continue
				}
				useri := lgc.getGameLogic().ext.getExtInfoByTurn(i)
				if useri != nil {
					upi := lgc.getGameLogic().gamePoker.GetUserPoker(useri.userId)
					if upi != nil {
						if upi.Cards.Len() < 2 {
							mlog.Debug("有敌人只剩一张单牌，机器人尽量避免出单")
							return true
						}
						//mlog.Debug("敌人剩牌[%d]", upi.Cards.Len())
					}
				}
			}
		} else {
			dizhuSeatNo := lgc.getGameLogic().ext.lastQiangDiZhu
			userDizhu := lgc.getGameLogic().ext.getExtInfoByTurn(dizhuSeatNo)
			if userDizhu != nil {
				uprDizhu := lgc.getGameLogic().gamePoker.GetUserPoker(userDizhu.userId)
				if uprDizhu != nil {
					if uprDizhu.Cards.Len() < 2 {
						example.invalid("userid[%v] 地主只剩一张单牌，机器人尽量避免出单", uprDizhu.User.UserID)
						return true
					}
					//mlog.Debug("地主剩牌[%d]", uprDizhu.Cards.Len())
				}
			}
		}

	}
	return false
}

//根据牌型出牌
func (lgc *DDZRobotLogic) outCardByType(outCardType cardType, userid int32) (a []int32) {
	defer func() {
		mlog.Debug("根据牌型出牌，选定的牌型为[%d],牌%v", int(outCardType), a)
	}()

	outCards := make([]int32, 0)

	if outCardType == cardTypeFourAndTwos || outCardType == cardTypeFourAndTwod { //四帶二單(雙)
		return lgc.outFourCards(userid)
	}
	//三带
	if outCardType == cardTypeThree { //三带
		return lgc.outThreeCards(userid)
	}
	//飞机
	if outCardType == cardTypePlaneAndWing {
		return lgc.outFeiji(userid, 5, 0)
	}
	//连对
	if outCardType == cardTypeThreeStraight {
		return lgc.outLiandui(userid)
	}
	//顺子
	if outCardType == cardTypeStraight {
		return lgc.outShunzi(userid)
	}
	//对子
	if outCardType == cardTypeDouble {
		return lgc.outDuizi(userid)
	}
	//炸弹
	if outCardType == cardTypeBomb {
		return lgc.outZhadan(userid)
	}
	//单牌
	if outCardType == cardTypeSingle {
		return lgc.outDanpai(userid)
	}
	//火箭
	if outCardType == cardTypeRocket {
		return lgc.outRocket(userid)
	}
	return outCards
}

//出牌 四帶
func (lgc *DDZRobotLogic) outFourCards(userid int32) []int32 {
	//四帶
	outCardType := cardTypeBomb
	outCards := make([]int32, 0)
	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		mlog.Error("机器人数据获取失败，四带")
		return outCards
	}

	outCardTypeGroup := lgc.getMinCardType(outCardType, 2, robotUser.ctg)
	if outCardTypeGroup == nil {
		if robotUser.ctg.hands < 3 {
			outCardTypeGroup = lgc.getMinCardType(outCardType, 0, robotUser.ctg)
			if outCardTypeGroup == nil {
				mlog.Error("机器人数据获取失败，四带")
				return outCards
			}
		} else {
			mlog.Error("机器人数据获取失败，四带")
			return outCards
		}
	}
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeBomb, 1)
	mlog.Debug("決定出四帶 single %d, double %d", robotUser.cardTypesCount[cardTypeSingle], robotUser.cardTypesCount[cardTypeDouble])

	singles, pairs := make([]int32, 0), make([]int32, 0)
	for _, s := range robotUser.ctg.SingleCards {
		if s <= 10 {
			singles = append(singles, s)
		}
	}
	for _, ct := range robotUser.ctg.Ct {
		if ct.byteCardType == cardTypeDouble && ct.nMax <= 10 && ct.nMax != 2 {
			pairs = append(pairs, ct.nMax)
			continue
		}
	}

	sort.Slice(singles, func(i, j int) bool { return cardMap[singles[i]] < cardMap[singles[j]] })
	sort.Slice(pairs, func(i, j int) bool { return cardMap[pairs[i]] < cardMap[pairs[j]] })
	var chudan, chupair bool
	if len(singles) > 1 && len(pairs) > 1 {
		if singles[0] < pairs[0] {
			chudan = true
		} else {
			chupair = true
		}
	} else if len(singles) > 1 {
		chudan = true
	} else if len(pairs) > 1 {
		chupair = true
	}

	if chudan {
		outCards = append(outCards, singles[0:2]...)
	} else if chupair {
		outCards = append(outCards, pairs[0], pairs[0])
		outCards = append(outCards, pairs[1], pairs[1])
	}
	return outCards
}

//出牌 三带
func (lgc *DDZRobotLogic) outThreeCards(userid int32) []int32 {
	//三带
	outCardType := cardTypeThree
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取三带牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		mlog.Error("机器人数据获取失败，三带")
		return outCards
	}

	outCardTypeGroup := lgc.getMinCardType(outCardType, 2, robotUser.ctg)
	if outCardTypeGroup == nil {
		if robotUser.ctg.hands < 3 {
			outCardTypeGroup = lgc.getMinCardType(outCardType, 0, robotUser.ctg)
			if outCardTypeGroup == nil {
				mlog.Error("机器人数据获取失败，三带")
				return outCards
			}
		} else {
			mlog.Error("机器人数据获取失败，三带")
			return outCards
		}
	}
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeThree, 1)
	mlog.Debug("決定出三帶 single %d, double %d", robotUser.cardTypesCount[cardTypeSingle], robotUser.cardTypesCount[cardTypeDouble])

	singleMin := int32(17)
	doubleMin := int32(17)
	for _, ct := range robotUser.ctg.Ct {
		if ct.byteCardType == cardTypeSingle && cardMap[ct.nMax] < cardMap[singleMin] {
			singleMin = ct.nMax
		}

		if ct.byteCardType == cardTypeDouble && cardMap[ct.nMax] < cardMap[doubleMin] {
			doubleMin = ct.nMax
		}
	}
	for _, i := range robotUser.ctg.SingleCards {
		if cardMap[i] < cardMap[singleMin] {
			singleMin = i
		}
	}

	if cardMap[singleMin] <= cardMap[doubleMin] && singleMin != 17 {
		//出三帶一, 不帶鬼跟2出或最後一手直接出
		if (singleMin != 15 && singleMin != 16 && singleMin != 2) || robotUser.ctg.hands <= 2 {
			outCards = append(outCards, singleMin)
			lgc.reduceCardType(robotUser, cardTypeSingle, 1)
			return outCards
		}
	} else if cardMap[doubleMin] <= cardMap[singleMin] && doubleMin != 17 {
		//出三帶二, 不帶2出或最後一手直接出
		if doubleMin != 2 || robotUser.ctg.hands <= 2 {
			double := lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
			outCards = append(outCards, doubleMin)
			outCards = append(outCards, doubleMin)
			lgc.eraseResult(robotUser.ctg, double)
			lgc.reduceCardType(robotUser, cardTypeDouble, 1)
			return outCards
		}
	}

	return outCards
}

//出牌 飞机
func (lgc *DDZRobotLogic) outFeiji(userid, sum, more int32) []int32 {
	outCardType := cardTypePlaneAndWing
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取飞机牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		// mlog.Error("机器人数据获取失败，飞机")
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, more, robotUser.ctg)
	if outCardTypeGroup == nil {
		mlog.Error("机器人数据获取失败，飞机")
		return outCards
	}
	if sum < 2 {
		return outCards //
	}
	size := outCardTypeGroup.nMax - outCardTypeGroup.nMin + 1
	if size > sum {
		size = sum
	}
	for i := outCardTypeGroup.nMax; i > outCardTypeGroup.nMax-size; i-- {
		outCards = append(outCards, i)
		outCards = append(outCards, i)
		outCards = append(outCards, i)
	}
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypePlaneAndWing, 1)
	//有size单牌
	bNBeedReset := false
	if robotUser.cardTypesCount[cardTypeSingle] >= size {
		singleList := make([]int32, 0)
		for i := int32(0); i < size; i++ {
			min := int32(0)
			robotUser.ctg.SingleCards, min = lgc.minOneSingle(robotUser.ctg.SingleCards)
			if min <= 0 {
				if lgc.getHand(userid) <= size+3 {
					singOne := lgc.getMinCardType(cardTypeSingle, 0, robotUser.ctg)
					if singOne == nil {
						bNBeedReset = true
						break
					}
					singleList = append(singleList, singOne.nMax)
					lgc.reduceCardType(robotUser, cardTypeSingle, 1)
					lgc.eraseResult(robotUser.ctg, singOne)
					continue
				} else {
					bNBeedReset = true
					break
				}

			}
			singleList = append(singleList, min)
			lgc.reduceCardType(robotUser, cardTypeSingle, 1)

		}
		if len(singleList) == int(size) {
			outCards = append(outCards, singleList...)
			return outCards
		}
	}
	//有两对的情况
	if robotUser.cardTypesCount[cardTypeDouble] >= size {
		doubleList := make([]int32, 0)
		for i := int32(0); i < size; i++ {
			double := new(cardTypeGroup)
			if lgc.getHand(userid) <= size+3 {
				double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
			} else {
				double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
			}
			if double == nil {
				mlog.Error("机器人数据获取失败，飞机")
				bNBeedReset = true
				break
			}
			doubleList = append(doubleList, double.nMax)
			doubleList = append(doubleList, double.nMax)
			lgc.reduceCardType(robotUser, cardTypeDouble, 1)
			lgc.eraseResult(robotUser.ctg, double)
		}
		if len(doubleList) == int(size)*2 {
			outCards = append(outCards, doubleList...)
			return outCards
		}

	}
	//检测手牌时候除了三条之外还有1~2张2 或一个王
	up := lgc.getGameLogic().gamePoker.GetUserPoker(userid)
	if up == nil {
		mlog.Error("系统错误robotOutCards，获取玩家牌型失败")
		return outCards
	}
	//有1对的情况
	if bNBeedReset {
		bNBeedReset = false
		lgc.recombination(userid, up.Cards)
		lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
		lgc.reduceCardType(robotUser, cardTypePlaneAndWing, 1)
	}
	if robotUser.cardTypesCount[cardTypeDouble]*2 >= size {
		doubleList := make([]int32, 0)
		x := size / 2 * 2
		if x < size {
			for i := int32(0); i < x/2; i++ {
				double := new(cardTypeGroup)
				if lgc.getHand(userid) <= size+3 {
					double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
				} else {
					double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
				}
				if double == nil {
					mlog.Error("机器人数据获取失败，飞机")
					bNBeedReset = true
					break
				}
				doubleList = append(doubleList, double.nMax)
				doubleList = append(doubleList, double.nMax)
				lgc.reduceCardType(robotUser, cardTypeDouble, 1)
				lgc.eraseResult(robotUser.ctg, double)
			}
			double := new(cardTypeGroup)
			if lgc.getHand(userid) <= size+3 {
				double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
			} else {
				double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
			}

			if double == nil {
				mlog.Error("机器人数据获取失败，飞机")
				bNBeedReset = true
			}
			doubleList = append(doubleList, double.nMax)
			lgc.reduceCardType(robotUser, cardTypeDouble, 1)
			lgc.eraseResult(robotUser.ctg, double)
		} else {
			for i := int32(0); i < size/2; i++ {
				double := new(cardTypeGroup)
				if lgc.getHand(userid) <= size+3 {
					double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
				} else {
					double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
				}
				if double == nil {
					mlog.Error("机器人数据获取失败，飞机")
					bNBeedReset = true
					break
				}
				doubleList = append(doubleList, double.nMax)
				doubleList = append(doubleList, double.nMax)
				lgc.reduceCardType(robotUser, cardTypeDouble, 1)
				lgc.eraseResult(robotUser.ctg, double)
			}
		}
		if len(doubleList) == int(size) {
			outCards = append(outCards, doubleList...)
			return outCards
		}
	}

	return outCards
}

//连对
func (lgc *DDZRobotLogic) outLiandui(userid int32) []int32 {
	outCardType := cardTypeThreeStraight
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取连对牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		mlog.Error("机器人数据获取失败，连对")
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, 0, robotUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("机器人数据获取失败，连对")
		return outCards
	}

	for i := outCardTypeGroup.nMin; i <= outCardTypeGroup.nMax; i++ {
		outCards = append(outCards, i)
		outCards = append(outCards, i)
	}
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeThreeStraight, 1)
	return outCards
}

//顺子
func (lgc *DDZRobotLogic) outShunzi(userid int32) []int32 {
	outCardType := cardTypeStraight
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取顺子牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		mlog.Error("机器人数据获取失败，顺子")
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, 0, robotUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("机器人数据获取失败，顺子")
		return outCards
	}
	for i := outCardTypeGroup.nMin; i <= outCardTypeGroup.nMax; i++ {
		outCards = append(outCards, i)
	}
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeStraight, 1)
	return outCards
}

//对子
func (lgc *DDZRobotLogic) outDuizi(userid int32) []int32 {
	outCardType := cardTypeDouble
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取对子牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		mlog.Error("获取机器人数据失败，对子")
		return outCards
	}
	mlog.Debug("出對子, 目前手數%d", robotUser.ctg.hands)
	//剩兩手時從最大的開始出
	if robotUser.ctg.hands < 3 {

		a := make([]int32, 0)
		max := int32(0)
		maxN := int32(0)
		for _, ct := range robotUser.ctg.Ct {
			if ct.byteCardType == cardTypeDouble {
				a = append(a, ct.nMax)
			}
		}
		for _, num := range a {
			if cardMap[num] > max {
				max = cardMap[num]
				maxN = num
			}
		}
		outCards = append(outCards, maxN)
		outCards = append(outCards, maxN)
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, 2, robotUser.ctg)
	if outCardTypeGroup == nil {
		if lgc.getHand(userid) < 3 || lgc.outNeverSingle(userid) {
			outCardTypeGroup = lgc.getMinCardType(outCardType, 0, robotUser.ctg)
			if outCardTypeGroup == nil {
				// mlog.Error("获取机器人数据失败，对子")
				return outCards
			}
		} else {
			// mlog.Error("获取机器人数据失败，对子")
			return outCards
		}
	}

	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeDouble, 1)
	return outCards
}

//炸弹
func (lgc *DDZRobotLogic) outZhadan(userid int32) []int32 {
	outCardType := cardTypeBomb
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取炸弹牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		mlog.Error("获取机器人数据失败，炸弹")
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, 0, robotUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("获取机器人数据失败，炸弹")
		return outCards
	}
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeBomb, 1)
	return outCards
}

//火箭
func (lgc *DDZRobotLogic) outRocket(userid int32) []int32 {
	outCardType := cardTypeRocket
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取火箭牌型%v", outCards)
	// }()

	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		//mlog.Error("获取机器人数据失败，火箭")
		return outCards
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, 0, robotUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("获取机器人数据失败，火箭")
		return outCards
	}
	//當大牌只有一組時要留在最後關頭才出
	if robotUser.ctg.bigs < 2 {
		lastUr := lgc.getGameLogic().GetUserBySeatNo(lgc.outInfo.RealSeatNo)
		enemyHands := lgc.getGameLogic().getUserPokerHands(lastUr.UserID)
		mlog.Debug("判斷要不要出王炸, enemyhands %d, selfhands %d", enemyHands, robotUser.ctg.hands)
		if enemyHands > 1 && robotUser.ctg.hands > 2 {
			return outCards
		}
	}
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMin)
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeRocket, 1)
	return outCards
}

//拆火箭成單牌
func (lgc *DDZRobotLogic) outRocketSplitOne(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	outCardType := cardTypeRocket
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取火箭牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[robot.UserInfo.UserID]
	if robotUser == nil {
		mlog.Error("获取机器人数据失败，火箭")
		return false, nil
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, 0, robotUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("获取机器人数据失败，火箭")
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}

	outCards = append(outCards, outCardTypeGroup.nMin)
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeRocket, 1)
	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//单牌
func (lgc *DDZRobotLogic) outDanpai(userid int32) []int32 {
	outCardType := cardTypeSingle
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取单牌牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[userid]
	if robotUser == nil {
		mlog.Error("获取机器人数据失败，单牌")
		return outCards
	}
	mlog.Debug("出單牌, 目前手數%d", robotUser.ctg.hands)
	//敵人只有一張牌就先出大牌
	var dizhuOpt *userDdzOpt
	descChu := false
	glc := lgc.getGameLogic()
	for _, opt := range glc.ext.userOpts {
		if opt.isDiZhu {
			dizhuOpt = opt
		}
	}
	if userid == dizhuOpt.userId { //如果是地主, 觀察其他用戶牌數
		for _, user := range glc.AllGameUser() {
			if user.UserID != userid {
				upk := glc.gamePoker.GetUserPoker(user.UserID)
				if upk.Cards.Len() < 2 {
					descChu = true
				}
			}
		}
	} else {
		upk := glc.gamePoker.GetUserPoker(dizhuOpt.userId)
		if upk.Cards.Len() < 2 {
			descChu = true
		}
	}
	//剩兩手時從最大的開始出
	if robotUser.ctg.hands < 3 {
		descChu = true
	}
	// 出最大的單牌
	if descChu {
		a := robotUser.ctg.SingleCards
		max := int32(0)
		maxN := int32(0)
		for _, ct := range robotUser.ctg.Ct {
			if ct.byteCardType == cardTypeSingle {
				a = append(a, ct.nMax)
			}
		}
		for _, num := range a {
			if cardMap[num] > max {
				max = cardMap[num]
				maxN = num
			}
		}
		outCards = append(outCards, maxN)
		return outCards
	}
	if robotUser.cardTypesCount[cardTypeSingle] > 0 {
		min := int32(0)
		robotUser.ctg.SingleCards, min = lgc.minOneSingle(robotUser.ctg.SingleCards)
		outCards = append(outCards, min)
		if len(outCards) > 0 {
			lgc.reduceCardType(robotUser, cardTypeSingle, 1)
			return outCards
		}
	}
	outCardTypeGroup := lgc.getMinCardType(outCardType, 0, robotUser.ctg)
	if outCardTypeGroup != nil {
		outCards = append(outCards, outCardTypeGroup.nMax)
		lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
		lgc.reduceCardType(robotUser, cardTypeSingle, 1)
		return outCards
	}
	return outCards
}

//获取不小于more的最小牌型
func (lgc *DDZRobotLogic) getMinCardType(findCardType cardType, more int32, ctg *CtGroup) (p *cardTypeGroup) {
	// defer func() {
	// 	if p != nil {
	// 		mlog.Debug("找到对应牌型[%d],more[%d],Info%v", int(findCardType), more, *p)
	// 	} else {
	// 		mlog.Debug("没有找到对应牌型[%d],more[%d]", int(findCardType), more)
	// 	}
	// }()

	minCt := new(cardTypeGroup)
	minCt.nMax = 20
	for _, ct := range ctg.Ct {
		if findCardType == ct.byteCardType {
			if minCt.nMax > ct.nMax && ct.nMax > more {
				minCt.byteCardType = ct.byteCardType
				minCt.nMax = ct.nMax
				minCt.nMin = ct.nMin
			}
		}
	}

	if minCt.nMax == 20 {
		return nil
	}
	return minCt
}

//移除对应牌组
func (lgc *DDZRobotLogic) eraseResult(ctg *CtGroup, minCt *cardTypeGroup) {
	for _, ct := range ctg.Ct {
		if minCt.byteCardType == ct.byteCardType && minCt.nMax == ct.nMax && minCt.nMin == ct.nMin {
			ct.byteCardType = cardTypeNone
			return
		}
	}
	return
}

//减少牌型
func (lgc *DDZRobotLogic) reduceCardType(robotUser *robotUserInfoLocal, bct cardType, num int32) {
	//mlog.Debug("减取牌型[%d],数量[%d]", int(bct), num)
	robotUser.cardTypesCount[bct] -= num
}

// 最小单牌
func (lgc *DDZRobotLogic) minOneSingle(a []int32) ([]int32, int32) {
	//mlog.Debug("单牌数量[%d],减少[1]", len(a))
	if len(a) <= 0 {
		return a, 0
	}
	a = lgc.sort(a)
	min := a[:1][0]
	a = a[1:]
	return a, min
}

// 解析牌组
func (lgc *DDZRobotLogic) analysisCards(cg *poker.CardGroup, outCardType int32) (nMin, nMax int32) {
	defer func() {
		mlog.Debug("收到上家出牌 ，牌型[%d],返回max[%d],min[%d]", outCardType, nMax, nMin)
	}()
	if cg == nil {
		return 0, 0
	}
	if cg.Len() <= 0 {
		return 0, 0
	}
	if outCardType == int32(cardTypeBomb) {
		return 0, cg.Cards[0].CardNumber
	}
	if outCardType == int32(cardTypePlaneAndWing) || outCardType == int32(cardTypeDoubleStraights) {

		List3 := make([]int32, 0)
		ListOther := make([]int32, 0)
		cardMap := make(map[int32]int32, 0)
		for _, cd := range cg.Cards {
			cardMap[cd.CardNumber]++
		}
		for number, count := range cardMap {
			if count == 3 {
				List3 = append(List3, number)
				continue
			}
			ListOther = append(ListOther, number)
		}
		List3 = lgc.sort(List3)
		if len(List3) <= 1 {
			return 0, 0
		}
		return List3[0], List3[1]
	}
	if outCardType == int32(cardTypeStraight) {
		strightList := make([]int32, 0)
		for _, cd := range cg.Cards {
			strightList = append(strightList, cd.CardNumber)
		}
		strightList = lgc.sort(strightList)
		lengthList := len(strightList)
		if lengthList < 5 {
			return 0, 0
		}
		return strightList[0], strightList[lengthList-1]
	}
	if outCardType == int32(cardTypeThreeStraight) || outCardType == 9 {
		strightList := make([]int32, 0)
		for _, cd := range cg.Cards {
			strightList = append(strightList, cd.CardNumber)
		}
		strightList = lgc.sort(strightList)
		lengthList := len(strightList)
		if lengthList < 6 {
			return 0, 0
		}
		return strightList[0], strightList[lengthList-1]
	}
	if outCardType == int32(cardTypeDouble) {
		return 0, cg.Cards[0].CardNumber
	}
	if outCardType == int32(cardTypeSingle) {
		return 0, cg.Cards[0].CardNumber
	}
	if outCardType == int32(cardTypeThreeAndOne) || outCardType == int32(cardTypeThree) {
		cardMap := make(map[int32]int32, 0)
		for _, cd := range cg.Cards {
			cardMap[cd.CardNumber]++
		}
		for number, count := range cardMap {
			if count == 3 {
				return 0, number
			}
		}
	}
	return 0, 0

}

//************************************************************************************跟牌
// 跟牌
func (lgc *DDZRobotLogic) FollowCards(robot *robotcomm.Robot, outCardType cardType) {
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		mlog.Error("系统错误FollowCards，获取玩家牌型失败")
		return
	}
	lgc.recombination(robot.UserInfo.UserID, up.Cards)
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		mlog.Debug("[%v] 沒有牌 robot pass", robot.UserInfo.UserID)
		lgc.getGameLogic().GameLog.AddGameStep("robotPass10 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 沒有牌")
		lgc.robotPass(robot)
		return
	}
	if outCardType == cardTypeRocket {
		lgc.getGameLogic().GameLog.AddGameStep("robotPass11 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是王炸")
		lgc.robotPass(robot)
		return
	}
	if outCardType == cardTypeBomb {
		var ok bool
		cg := new(poker.CardGroup)
		if ok, cg = lgc.FollowCardsBomb(robot); !ok {
			cdlist := lgc.outRocket(robot.UserInfo.UserID)
			pass := false
			if lgc.robotHandMoreByPass(robot) {
				pass = true
			}

			robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
			_, dizhu := lgc.getGameLogic().ext.getDiZhu()
			dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
			mlog.Debug("炸彈 len(cdlist)[%v] pass[%v] isrobot[%v] robotHands[%v]",
				len(cdlist), pass, dizhuUser.IsRobot(), robotHands)

			if len(cdlist) <= 0 || pass {
				lgc.getGameLogic().GameLog.AddGameStep("robotPass12 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是炸彈，沒有王炸")
				lgc.robotPass(robot)
				return
			} else {
				if !ok && len(cdlist) > 0 {
					cg := lgc.list2cardGroup(up.Cards, cdlist)
					lgc.S2COutCard(cg, robot.UserInfo)
					return
				}
			}
		}

		if ok {
			lgc.S2COutCard(cg, robot.UserInfo)
		} else {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass13 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是炸彈，沒有更大的牌")
			lgc.robotPass(robot)
		}

		return
	}
	if outCardType == cardTypePlaneAndWing {
		var ok bool
		cg := new(poker.CardGroup)
		if ok, cg = lgc.FollowCardsFeiji(robot); !ok {
			if ok, cg = lgc.FollowCardsBomb(robot); !ok && ctg.ctg.hands > 5 {
				cdlist := lgc.outRocket(robot.UserInfo.UserID)
				pass := false
				if lgc.robotHandMoreByPass(robot) {
					pass = true
				}

				robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
				_, dizhu := lgc.getGameLogic().ext.getDiZhu()
				dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
				mlog.Debug("炸彈 len(cdlist)[%v] pass[%v] isrobot[%v] robotHands[%v]",
					len(cdlist), pass, dizhuUser.IsRobot(), robotHands)

				if len(cdlist) <= 0 || pass {
					lgc.getGameLogic().GameLog.AddGameStep("robotPass14 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是飞机带翅膀，沒有王炸或跟牌時當機器人手牌數比真人大就不出")
					lgc.robotPass(robot)
					return
				} else {
					if !ok && len(cdlist) > 0 {
						cg := lgc.list2cardGroup(up.Cards, cdlist)
						lgc.S2COutCard(cg, robot.UserInfo)
						return
					}
				}
			}
		}

		if ok {
			lgc.S2COutCard(cg, robot.UserInfo)
		} else {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass15 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是飞机带翅膀，沒有更大的牌或炸彈")
			lgc.robotPass(robot)
		}

		return
	}
	if outCardType == cardTypeStraight {
		var ok bool
		cg := new(poker.CardGroup)
		if ok, cg = lgc.FollowCardsShunzi(robot); !ok {
			if ok, cg = lgc.FollowCardsBomb(robot); !ok && ctg.ctg.hands > 5 {
				cdlist := lgc.outRocket(robot.UserInfo.UserID)
				pass := false
				if lgc.robotHandMoreByPass(robot) {
					pass = true
				}

				robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
				_, dizhu := lgc.getGameLogic().ext.getDiZhu()
				dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
				mlog.Debug("炸彈 len(cdlist)[%v] pass[%v] isrobot[%v] robotHands[%v]",
					len(cdlist), pass, dizhuUser.IsRobot(), robotHands)

				if len(cdlist) <= 0 || pass {
					lgc.getGameLogic().GameLog.AddGameStep("robotPass16 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是单顺（顺子），沒有王炸或跟牌時當機器人手牌數比真人大就不出")
					lgc.robotPass(robot)
					return
				} else {
					if !ok && len(cdlist) > 0 {
						cg := lgc.list2cardGroup(up.Cards, cdlist)
						lgc.S2COutCard(cg, robot.UserInfo)
						return
					}
				}
			}
		}

		if ok {
			lgc.S2COutCard(cg, robot.UserInfo)
		} else {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass17 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是单顺（顺子），沒有更大的牌或炸彈")
			lgc.robotPass(robot)
		}
		return
	}
	if outCardType == cardTypeThreeStraight {
		var ok bool
		cg := new(poker.CardGroup)
		if ok, cg = lgc.FollowCardsLiandui(robot); !ok {
			if ok, cg = lgc.FollowCardsBomb(robot); !ok && ctg.ctg.hands > 5 {
				cdlist := lgc.outRocket(robot.UserInfo.UserID)
				pass := false
				if lgc.robotHandMoreByPass(robot) {
					pass = true
				}

				robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
				_, dizhu := lgc.getGameLogic().ext.getDiZhu()
				dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
				mlog.Debug("炸彈 len(cdlist)[%v] pass[%v] isrobot[%v] robotHands[%v]",
					len(cdlist), pass, dizhuUser.IsRobot(), robotHands)

				if len(cdlist) <= 0 || pass {
					lgc.getGameLogic().GameLog.AddGameStep("robotPass18 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是連對(對子順)，沒有王炸或跟牌時當機器人手牌數比真人大就不出")
					lgc.robotPass(robot)
					return
				} else {
					if !ok && len(cdlist) > 0 {
						cg := lgc.list2cardGroup(up.Cards, cdlist)
						lgc.S2COutCard(cg, robot.UserInfo)
						return
					}
				}
			}
		}

		if ok {
			lgc.S2COutCard(cg, robot.UserInfo)
		} else {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass19 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是連對(對子順)，沒有更大的牌或炸彈")
			lgc.robotPass(robot)
		}

		return
	}
	if outCardType == cardTypeThree {
		var ok bool
		cg := new(poker.CardGroup)
		if ok, cg = lgc.FollowCardsThree(robot); !ok {
			if ok, cg = lgc.FollowCardsBomb(robot); !ok && ctg.ctg.hands > 5 {
				cdlist := lgc.outRocket(robot.UserInfo.UserID)
				pass := false
				if lgc.robotHandMoreByPass(robot) {
					pass = true
				}

				robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
				_, dizhu := lgc.getGameLogic().ext.getDiZhu()
				dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
				mlog.Debug("炸彈 len(cdlist)[%v] pass[%v] isrobot[%v] robotHands[%v]",
					len(cdlist), pass, dizhuUser.IsRobot(), robotHands)

				if len(cdlist) <= 0 || pass {
					mlog.Debug("bomb[%v] pass[%v]", cg, pass)
					lgc.getGameLogic().GameLog.AddGameStep("robotPass20 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是三张，沒有王炸或跟牌時當機器人手牌數比真人大就不出")
					lgc.robotPass(robot)
					return
				} else {
					if !ok && len(cdlist) > 0 {
						cg := lgc.list2cardGroup(up.Cards, cdlist)
						lgc.S2COutCard(cg, robot.UserInfo)
						return
					}
				}
			}
		}

		if ok {
			lgc.S2COutCard(cg, robot.UserInfo)
		} else {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass21 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是三张，沒有更大的牌或炸彈")
			lgc.robotPass(robot)
		}
		return
	}
	if outCardType == cardTypeDouble {
		var ok bool
		cg := new(poker.CardGroup)
		if ok, cg = lgc.FollowCardsDuizi(robot); !ok {
			//拆连对
			if ok, cg = lgc.FollowCardsDuiziLian(robot); !ok {
				//拆三个的
				if ok, cg = lgc.FollowCardsDuiziUse3(robot); !ok {
					if ok, cg = lgc.FollowCardsBomb(robot); !ok {
						cdlist := lgc.outRocket(robot.UserInfo.UserID)
						pass := false
						if lgc.robotHandMoreByPass(robot) {
							pass = true
						}

						robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
						_, dizhu := lgc.getGameLogic().ext.getDiZhu()
						dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
						mlog.Debug("炸彈 len(cdlist)[%v] pass[%v] isrobot[%v] robotHands[%v]",
							len(cdlist), pass, dizhuUser.IsRobot(), robotHands)

						if len(cdlist) <= 0 || pass {
							lgc.getGameLogic().GameLog.AddGameStep("robotPass22 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是对子，沒有王炸或跟牌時當機器人手牌數比真人大就不出")
							lgc.robotPass(robot)
							return
						} else {
							if !ok && len(cdlist) > 0 {
								cg := lgc.list2cardGroup(up.Cards, cdlist)
								lgc.S2COutCard(cg, robot.UserInfo)
								return
							}
						}
					}
				}
			}
		}

		if ok {
			lgc.S2COutCard(cg, robot.UserInfo)
		} else {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass23 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是对子，沒有更大的牌或炸彈")
			lgc.robotPass(robot)
		}

		return
	}
	if outCardType == cardTypeSingle {
		var ok bool
		cg := new(poker.CardGroup)
		if ok, cg = lgc.FollowCardsDanpai(robot); !ok {
			if ok, cg = lgc.FollowCardsDanpaiUse2(robot); !ok {
				if ok, cg = lgc.FollowCardsDanpaiUse3(robot); !ok {
					//建議不拆王
					//if ok, cg = lgc.outRocketSplitOne(robot); !ok {
					robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
					if ok, cg = lgc.FollowCardsBomb(robot); !ok && robotHands > 3 {
						cdlist := lgc.outRocket(robot.UserInfo.UserID)
						pass := false
						if lgc.robotHandMoreByPass(robot) {
							pass = true
						}

						robotHands := lgc.getGameLogic().getUserPokerHands(robot.UserInfo.UserID)
						_, dizhu := lgc.getGameLogic().ext.getDiZhu()
						dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
						//dizhuHands := lgc.getGameLogic().getUserPokerHands(dizhuUser.UserID)

						mlog.Debug("炸彈 len(cdlist)[%v] pass[%v] isrobot[%v] 機器人手牌數[%v]",
							len(cdlist), pass, dizhuUser.IsRobot(), robotHands)

						if (len(cdlist) <= 0) || pass || (robotHands > 7 && dizhuUser.IsRobot()) {
							lgc.getGameLogic().GameLog.AddGameStep("robotPass24 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是单牌，沒有王炸或跟牌時當機器人手牌數比真人大就不出或剩大於7手而且地主是機器人")
							lgc.robotPass(robot)
							return
						} else {
							if !ok && len(cdlist) > 0 {
								cg := lgc.list2cardGroup(up.Cards, cdlist)
								lgc.S2COutCard(cg, robot.UserInfo)
								return
							}
						}
					}
					//}
				}
			}
		}

		if ok {
			lgc.S2COutCard(cg, robot.UserInfo)
		} else {
			lgc.getGameLogic().GameLog.AddGameStep("robotPass25 RobotID:[%d] Msg:[%s]", robot.UserInfo.UserID, "FollowCards 場面牌是单牌，沒有更大的牌或炸彈")
			lgc.robotPass(robot)
		}

		return
	}

	//if outCardType == cardTypeFourAndTwos {
	//	mlog.Error("未實作.")
	//}
}

//炸弹特殊处理
func (lgc *DDZRobotLogic) FollowCardsBomb(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	//defer func() {
	//	mlog.Debug("跟炸弹,是否成功[%t]", b)
	//}()

	if lgc.robotHandMoreByPass(robot) {
		mlog.Warn("機器人手牌超過真人, pass Bomb (待確認)")
		//return false, nil
	}

	//_ := lgc.passTwoHands(ctg, dc, robot, passTwoHand, roleInfo)

	robotUser := lgc.robotUser[robot.UserInfo.UserID]
	if robotUser == nil {
		return false, nil
	}
	//不出狀態不壓炸彈
	if robot.BuchuBombCard {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}

	if lgc.outInfo.nMax == 2 && lgc.outInfo.CardType == int32(cardTypeBomb) {
		return false, nil
	}
	var data *cardTypeGroup
	if lgc.outInfo.CardType == int32(cardTypeRocket) {
		return false, nil
	} else if lgc.outInfo.CardType == int32(cardTypeBomb) { //當最後出牌是炸彈, 手牌沒有更大點數則不出
		data = lgc.getMinCardType(cardTypeBomb, lgc.outInfo.nMax, robotUser.ctg)
	} else { //跟牌或壓牌, 當最後出牌不是炸彈時, 就用炸彈去壓
		data = lgc.getMinCardType(cardTypeBomb, 2, robotUser.ctg)
		mlog.Debug("robot[%v] 炸彈去壓 card[%v]", robot.UserInfo.UserID, data)
		if data == nil {
			data = lgc.getMinCardType(cardTypeBomb, 0, robotUser.ctg)
			mlog.Debug("robot[%v] 炸彈去壓B card[%v]", robot.UserInfo.UserID, data)
		}
	}

	if data != nil {
		//// PS. lgc.outInfo.userID 永遠都會是0
		lastUr := lgc.getGameLogic().GetUserBySeatNo(lgc.outInfo.RealSeatNo)
		enemyHands := lgc.getGameLogic().getUserPokerHands(lastUr.UserID)
		//若炸彈是2, 敵人手數不是2手, 且自己不是最後兩手的話不跟
		if data.nMax == 2 {
			mlog.Debug("為什麼要出炸彈2, enemyhands %d, selfhands %d, bigs %d", enemyHands, robotUser.ctg.hands, robotUser.ctg.bigs)
			if enemyHands > 1 && robotUser.ctg.hands > 2 {
				return false, nil
			}
		}
		//當大牌只有一組時要留在最後關頭才出
		//當地主不是玩家時才保留最後一手大牌
		dizhuIsRealUser := true
		_, dizhu := lgc.getGameLogic().ext.getDiZhu()
		if dizhu != nil {
			dizhuUser := lgc.getGameLogic().GetUserByUserID(dizhu.userId)
			if dizhuUser != nil && dizhuUser.IsRobot() {
				dizhuIsRealUser = false
			}
		}
		if robotUser.ctg.bigs < 2 && !dizhuIsRealUser {
			//最後兩手且不是場上最大的牌
			if robotUser.ctg.hands > 1 && robotUser.ctg.hands < 3 && !lgc.getGameLogic().isBiggestCT(data, robot.UserInfo.UserID) {
				return false, nil
			}
		}

		cdList := make([]int32, 0)
		cdList = append(cdList, data.nMax)
		cdList = append(cdList, data.nMax)
		cdList = append(cdList, data.nMax)
		cdList = append(cdList, data.nMax)
		cg := lgc.list2cardGroup(up.Cards, cdList)
		lgc.S2COutCard(cg, robot.UserInfo)
		return true, cg

	} else {
		return false, nil
	}
	// return false, nil
}

//飞机
func (lgc *DDZRobotLogic) FollowCardsFeiji(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	//defer func() {
	//	mlog.System("跟飞机,是否成功[%t]", cg)
	//}()
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	size := lgc.outInfo.nMax - lgc.outInfo.nMin + 1
	cards := lgc.outFeiji(robot.UserInfo.UserID, size, lgc.outInfo.nMax)
	if len(cards) != lgc.outInfo.cg.Len() && lgc.outInfo.cg.Len() != 0 {
		return false, nil
	}
	cg := lgc.list2cardGroup(up.Cards, cards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

// 获取手牌手数
func (lgc *DDZRobotLogic) getHand(userid int32) int32 {
	totalHand := int32(0)
	ctg := lgc.robotUser[userid]
	if ctg == nil {
		return 20
	}
	for _, count := range ctg.cardTypesCount {
		if count <= 0 {
			continue
		}
		totalHand += count
	}
	return totalHand
}

// 发送出牌结果
func (lgc *DDZRobotLogic) S2COutCard(cg *poker.CardGroup, ur *gamelogic.GameUser) {
	up := lgc.getGameLogic().gamePoker.GetUserPoker(ur.UserID)

	mlog.NewDefaultLogger().SetCallDepLen(10).SetData(fmt.Sprintf("userid[%v]机器人出牌[%s] 手牌還有 %v", ur.UserID, cg.String(), up.Cards)).Warn("跟牌")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_ChuPai)

	btdata := new(netproto.DDZ_ChuPai)

	for _, cd := range cg.Cards {
		cdTmp := new(netproto.Card)
		cdTmp.Color = proto.Int32(cd.Color)
		cdTmp.Number = proto.Int32(cd.CardNumber)
		btdata.Cards = append(btdata.Cards, cdTmp)
	}

	btmsg.MsgData, _ = proto.Marshal(btdata)
	lgc.SendGameMessage(ur, btmsg)
}

//三带
func (lgc *DDZRobotLogic) FollowCardsThree(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	// defer func() {
	// 	mlog.Debug("跟三带,是否成功[%t]", b)
	// }()
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	num := int32(0) //
	if lgc.outInfo.cg.Len() == 3 {
		num = 0
	} else if lgc.outInfo.cg.Len() == 4 {
		num = 1
	} else if lgc.outInfo.cg.Len() == 5 {
		num = 2
	}
	cards := lgc.outThreeNum(robot, num)
	if len(cards) < lgc.outInfo.cg.Len() {
		return false, nil
	}
	cg := lgc.list2cardGroup(up.Cards, cards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

// 出牌三带
func (lgc *DDZRobotLogic) outThreeNum(robot *robotcomm.Robot, num int32) []int32 {
	// mlog.Debug("找三带，上家三带[%d]", lgc.outInfo.nMax)
	//三带
	outCardType := cardTypeThree
	outCards := make([]int32, 0)
	// defer func() {
	// 	mlog.Debug("获取三带牌型%v", outCards)
	// }()
	robotUser := lgc.robotUser[robot.UserInfo.UserID]
	if robotUser == nil {
		mlog.Error("机器人数据获取失败，三带")
		return outCards
	}
	if lgc.outInfo.nMax == 2 {
		return outCards
	}
	outCardTypeGroup := new(cardTypeGroup)
	outCardTypeGroup = lgc.getMinCardType(outCardType, lgc.outInfo.nMax, robotUser.ctg)
	if outCardTypeGroup == nil {
		outCardTypeGroup = lgc.getMinCardType(outCardType, 0, robotUser.ctg)
		if outCardTypeGroup == nil {
			return outCards
		}
		// mlog.Debug("搜索三条，查看是否是三条 2==%d", outCardTypeGroup.nMax)
		if outCardTypeGroup.nMax != 2 {
			return outCards
		}
	}

	//如果最小的三帶是2, 敵人跟自己的手數都大於2則不出
	lastUr := lgc.getGameLogic().GetUserBySeatNo(lgc.outInfo.RealSeatNo)
	enemyHands := lgc.getGameLogic().getUserPokerHands(lastUr.UserID)
	if outCardTypeGroup.nMax == 2 {
		mlog.Debug("為什麼不出2三帶, enemyHands %d, hands %d", enemyHands, robotUser.ctg.hands)
		if enemyHands > 1 && robotUser.ctg.hands > 2 {
			return outCards
		}
	}

	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
	lgc.reduceCardType(robotUser, cardTypeThree, 1)
	if num == 0 {
		return outCards
	} else if num == 1 {
		//三带一
		if robotUser.cardTypesCount[cardTypeSingle] >= 1 {
			min := int32(0)
			robotUser.ctg.SingleCards, min = lgc.minOneSingle(robotUser.ctg.SingleCards)
			if min == 0 {
				singleTmp := lgc.getMinCardType(cardTypeSingle, 0, robotUser.ctg)
				if singleTmp != nil {
					min = singleTmp.nMax
					outCards = append(outCards, min)
				}
			} else {
				outCards = append(outCards, min)
			}
			lgc.reduceCardType(robotUser, cardTypeSingle, 1)
			lgc.eraseResult(robotUser.ctg, outCardTypeGroup)
			return outCards
		}
	} else if num == 2 {
		if robotUser.cardTypesCount[cardTypeDouble] >= 1 {
			double := new(cardTypeGroup)
			if lgc.getHand(robot.UserInfo.UserID) <= 3 {
				double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
			} else {
				double = lgc.getMinCardType(cardTypeDouble, 2, robotUser.ctg)
			}
			if double == nil {
				double = lgc.getMinCardType(cardTypeDouble, 0, robotUser.ctg)
				if double != nil {
					outCards = append(outCards, double.nMax)
					outCards = append(outCards, double.nMax)
					lgc.eraseResult(robotUser.ctg, double)
					lgc.reduceCardType(robotUser, cardTypeDouble, 1)
					return outCards
				}
			} else {
				outCards = append(outCards, double.nMax)
				outCards = append(outCards, double.nMax)
				lgc.eraseResult(robotUser.ctg, double)
				lgc.reduceCardType(robotUser, cardTypeDouble, 1)
				return outCards
			}
		}
	}
	return outCards
}

//出牌连对
func (lgc *DDZRobotLogic) FollowCardsLiandui(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	// defer func() {
	// 	mlog.Debug("跟连对,是否成功[%t]", b)
	// }()
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	size := lgc.outInfo.nMax - lgc.outInfo.nMin + 1

	outCardTypeGroup := lgc.getMinCardType(cardTypeThreeStraight, lgc.outInfo.nMax, ctg.ctg)
	if outCardTypeGroup == nil {
		return false, nil
	}
	if outCardTypeGroup.nMax-outCardTypeGroup.nMin+1 < size {
		return false, nil
	}
	outCards := make([]int32, 0)
	for i := lgc.outInfo.nMin + 1; i <= outCardTypeGroup.nMax; i++ {
		if i < outCardTypeGroup.nMin {
			continue
		}
		outCards = append(outCards, i)
		outCards = append(outCards, i)
		if len(outCards) == lgc.outInfo.cg.Len() {
			break
		}
	}
	lgc.eraseResult(ctg.ctg, outCardTypeGroup)
	lgc.reduceCardType(ctg, cardTypeThreeStraight, 1)

	if len(outCards) < lgc.outInfo.cg.Len() {
		return false, nil
	}
	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//跟牌顺子
func (lgc *DDZRobotLogic) FollowCardsShunzi(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	// defer func() {
	// 	mlog.Debug("跟顺子,是否成功[%t]", b)
	// }()
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	size := lgc.outInfo.nMax - lgc.outInfo.nMin + 1

	outCardTypeGroup := lgc.getMinCardType(cardTypeStraight, lgc.outInfo.nMax, ctg.ctg)
	if outCardTypeGroup == nil {
		return false, nil
	}
	if outCardTypeGroup.nMax-outCardTypeGroup.nMin+1 < size {
		return false, nil
	}
	outCards := make([]int32, 0)
	for i := lgc.outInfo.nMin + 1; i <= lgc.outInfo.nMin+size; i++ {
		if i < outCardTypeGroup.nMin {
			continue
		}
		outCards = append(outCards, i)
		if len(outCards) == lgc.outInfo.cg.Len() {
			break
		}
	}
	lgc.eraseResult(ctg.ctg, outCardTypeGroup)
	lgc.reduceCardType(ctg, cardTypeThreeStraight, 1)

	if len(outCards) < lgc.outInfo.cg.Len() {
		return false, nil
	}
	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//跟牌对子
func (lgc *DDZRobotLogic) FollowCardsDuizi(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	// defer func() {
	// 	mlog.Debug("跟对子,是否成功[%t]", b)
	// }()
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		mlog.Error("没有找到机器人信息userid[%d]", robot.UserInfo.UserID)
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		mlog.Error("没有找到机器人信息userid[%d]", robot.UserInfo.UserID)
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		// mlog.Debug("手牌不足接牌[%d]", robot.UserInfo.UserID)
		return false, nil
	}
	if lgc.outInfo.nMax == 2 {
		// mlog.Debug("没有比2更大的对子")
		return false, nil
	}
	outCardTypeGroup := lgc.getMinCardType(cardTypeDouble, lgc.outInfo.nMax, ctg.ctg)
	if outCardTypeGroup == nil {
		outCardTypeGroup = lgc.getMinCardType(cardTypeDouble, 0, ctg.ctg)
		if outCardTypeGroup == nil {
			return false, nil
		}
		if outCardTypeGroup.nMax != 2 {
			return false, nil
		}
	}

	outCards := make([]int32, 0)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)

	lgc.eraseResult(ctg.ctg, outCardTypeGroup)
	lgc.reduceCardType(ctg, cardTypeThreeStraight, 1)

	if len(outCards) < lgc.outInfo.cg.Len() {
		return false, nil
	}
	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//跟牌单牌
func (lgc *DDZRobotLogic) FollowCardsDanpai(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	outCards := make([]int32, 0)
	newSingle := make([]int32, 0)
	outTmp := int32(0)
	for _, cd := range ctg.ctg.SingleCards {
		if cd > lgc.outInfo.nMax && lgc.outInfo.nMax != 14 && lgc.outInfo.nMax != 2 {
			outTmp = cd
			outCards = append(outCards, cd)
			break
		}
	}
	havaSave := true
	for _, cd := range ctg.ctg.SingleCards {
		if havaSave && cd == outTmp {
			havaSave = false
			continue
		}
		newSingle = append(newSingle, cd)
	}
	ctg.ctg.SingleCards = newSingle

	if len(outCards) <= 0 {
		outCardTypeGroup := new(cardTypeGroup)

		if lgc.outInfo.nMax > 2 && lgc.outInfo.nMax <= 14 {
			outCardTypeGroup = lgc.getMinCardType(cardTypeSingle, 0, ctg.ctg)
		} else {
			outCardTypeGroup = lgc.getMinCardType(cardTypeSingle, lgc.outInfo.nMax, ctg.ctg)
		}
		if outCardTypeGroup == nil {
			return false, nil
		}
		lgc.eraseResult(ctg.ctg, outCardTypeGroup)
		outCards = append(outCards, outCardTypeGroup.nMax)
	}

	lgc.reduceCardType(ctg, cardTypeSingle, 1)

	if len(outCards) < lgc.outInfo.cg.Len() {
		return false, nil
	}

	//如果下家是敌人 并且下家牌只剩一张，出最大的牌
	roleInfo := lgc.userRole[robot.UserInfo.UserID]
	if roleInfo != nil {
		mlog.Debug("robot[%v] bnext[%v] bmine[%v] 如果下家是敌人 并且下家牌只剩一张，出最大的牌", robot.UserInfo.UserID, roleInfo.bNext, roleInfo.bMine)
		if roleInfo.bNext != roleInfo.bMine {
			user := lgc.getGameLogic().ext.getExtInfoByTurn(roleInfo.SeatNext)
			if user != nil {
				upNext := lgc.getGameLogic().gamePoker.GetUserPoker(user.userId)
				if upNext != nil {
					if upNext.Cards.Len() > 2 {
						outCards2 := make([]int32, 0)
						for _, cdOut := range up.Cards.Cards {
							if lgc.outInfo.nMax == poker.PokerNumberTwo && cdOut.CardNumber < poker.PokerNumberSmallKing {
								mlog.Debug("最後出手牌2, 選擇牌小於小王")
								continue
							}

							if cdOut.CardNumber > lgc.outInfo.nMax {
								outCards2 = append(outCards2, cdOut.CardNumber)
								break
							}
						}
						outCards2 = lgc.sort(outCards2)
						outCards = make([]int32, 0)
						lengthOut := int32(len(outCards2))
						if len(outCards2) > 0 {
							outCTmp := outCards2[lengthOut-1]
							outCards = append(outCards, outCTmp)
						}
					}

					if upNext.Cards.Len() < 2 {
						outCards2 := make([]int32, 0)
						for _, cdOut := range up.Cards.Cards {
							outCards2 = append(outCards2, cdOut.CardNumber)
						}
						outCards2 = lgc.sort(outCards2)
						outCards = make([]int32, 0)
						lengthOut := int32(len(outCards2))
						outCTmp := outCards2[lengthOut-1]
						if outCTmp <= 14 && outCards2[0] == 2 {
							outCTmp = 2
						}
						outCards = append(outCards, outCTmp)
					}
				}
			}
		}
	}

	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//拆连对配对子
func (lgc *DDZRobotLogic) FollowCardsDuiziLian(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	if lgc.outInfo.nMax == 2 {
		return false, nil
	}
	outCardTypeGroup := lgc.getMinCardType(cardTypeThreeStraight, lgc.outInfo.nMax, ctg.ctg)
	if outCardTypeGroup == nil {
		return false, nil
	}

	outCards := make([]int32, 0)
	for i := outCardTypeGroup.nMin; i <= outCardTypeGroup.nMax; i++ {
		if i > lgc.outInfo.nMax {
			outCards = append(outCards, i)
			outCards = append(outCards, i)
			break
		}
	}
	if len(outCards) < lgc.outInfo.cg.Len() {
		return false, nil
	}
	lgc.eraseResult(ctg.ctg, outCardTypeGroup)
	lgc.reduceCardType(ctg, cardTypeThree, 1)

	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//拆三个的配对子
func (lgc *DDZRobotLogic) FollowCardsDuiziUse3(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	if lgc.outInfo.nMax == 2 {
		return false, nil
	}
	outCardTypeGroup := lgc.getMinCardType(cardTypeThree, lgc.outInfo.nMax, ctg.ctg)
	if outCardTypeGroup == nil {
		if lgc.outInfo.nMax <= 14 {
			outCardTypeGroup = lgc.getMinCardType(cardTypeThree, 0, ctg.ctg)
		}
		if outCardTypeGroup == nil {
			return false, nil
		}
		if outCardTypeGroup.nMax != 2 {
			return false, nil
		}
	}

	outCards := make([]int32, 0)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)

	lgc.eraseResult(ctg.ctg, outCardTypeGroup)
	lgc.reduceCardType(ctg, cardTypeThree, 1)

	if len(outCards) < lgc.outInfo.cg.Len() {
		return false, nil
	}
	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//拆对子为单牌
func (lgc *DDZRobotLogic) FollowCardsDanpaiUse2(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	if lgc.outInfo.nMax == 2 {
		return false, nil
	}
	outCardTypeGroup := lgc.getMinCardType(cardTypeDouble, lgc.outInfo.nMax, ctg.ctg)
	if outCardTypeGroup == nil {
		if lgc.outInfo.nMax <= 14 {
			outCardTypeGroup = lgc.getMinCardType(cardTypeDouble, 0, ctg.ctg)
		}
		if outCardTypeGroup == nil {
			return false, nil
		}
		if outCardTypeGroup.nMax != 2 {
			return false, nil
		}
	}

	outCards := make([]int32, 0)
	outCards = append(outCards, outCardTypeGroup.nMax)

	lgc.eraseResult(ctg.ctg, outCardTypeGroup)
	lgc.reduceCardType(ctg, cardTypeDouble, 1)

	if len(outCards) < lgc.outInfo.cg.Len() {
		return false, nil
	}
	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//拆3为单牌
func (lgc *DDZRobotLogic) FollowCardsDanpaiUse3(robot *robotcomm.Robot) (bool, *poker.CardGroup) {
	ctg := lgc.robotUser[robot.UserInfo.UserID]
	if ctg == nil {
		return false, nil
	}
	up := lgc.getGameLogic().gamePoker.GetUserPoker(robot.UserInfo.UserID)
	if up == nil {
		return false, nil
	}
	if lgc.outInfo.cg.Len() > up.Cards.Len() {
		return false, nil
	}
	if lgc.outInfo.nMax == 2 {
		return false, nil
	}
	outCardTypeGroup := lgc.getMinCardType(cardTypeThree, lgc.outInfo.nMax, ctg.ctg)
	if outCardTypeGroup == nil {
		if lgc.outInfo.nMax <= 14 {
			outCardTypeGroup = lgc.getMinCardType(cardTypeThree, 0, ctg.ctg)
		}
		if outCardTypeGroup == nil {
			return false, nil
		}
		if outCardTypeGroup.nMax != 2 {
			return false, nil
		}
	}

	outCards := make([]int32, 0)
	outCards = append(outCards, outCardTypeGroup.nMax)

	lgc.eraseResult(ctg.ctg, outCardTypeGroup)
	lgc.reduceCardType(ctg, cardTypeThree, 1)

	if len(outCards) < lgc.outInfo.cg.Len() {
		return false, nil
	}
	cg := lgc.list2cardGroup(up.Cards, outCards)
	//lgc.S2COutCard(cg, robot.UserInfo)
	return true, cg
}

//******拆牌***************************************************************

//拆分手牌
func (lgc *DDZRobotLogic) Combination(cg *poker.CardGroup, userid int32) *CtGroup {
	mlog.Debug("开始拆分手牌")
	wMax, cMin := int32(0), int32(21)
	resultBest := new(CtGroup)
	resultBest.Ct = make([]*cardTypeGroup, 0)
	singleBest := make([]int32, 0)
	for i := 1; i <= 7; i++ {
		result := new(CtGroup)
		result.Ct = make([]*cardTypeGroup, 0)
		single := make([]int32, 0)
		mapCard := mapCardNumberCount(cg)
		mapCard = lgc.combinationKing(mapCard, result)                               //王
		mapCard = lgc.combinationTwo(mapCard, result)                                //2
		mapCard = lgc.combinationCardsIndex(mapCard, cardTypeBomb, result, int32(i)) //剩下单牌，

		erasePair := make([]int32, 0)
		eraseThree := make([]int32, 0)
		//mlog.Debug("补张开始")
		//把对子或者三放入单牌中看能否消掉一些，这里单牌还没有带
		pairList, threeList := lgc.getPairAndThree(result)

		for _, cd := range pairList {
			mapCard[cd] = 2
		}
		mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result) //顺子

		for _, cd := range threeList {
			mapCard[cd] = 3
			mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result) //顺子
			if mapCard[cd] != 3 {
				//标记，并且去除结果中的三条
				eraseThree = append(eraseThree, cd)
			} else {
				mapCard[cd] = 0 //如果没有匹配到，从重置为0
			}
		}

		for _, cd := range pairList {
			if mapCard[cd] != 2 {
				//标记，并且去除结果中的对子
				erasePair = append(erasePair, cd)
			} else {
				mapCard[cd] = 0
			}
		}
		mapCard = lgc.combinationCards(mapCard, cardTypeDouble, result) //对子

		//跟新result,erasePair,eraseThree
		result = lgc.updataResult(result, erasePair, eraseThree)

		mapCard, single = lgc.combinationSingle(mapCard) //单牌
		for number, count := range mapCard {
			if count > 0 {
				mlog.Error("错误数据 number[%d],count[%d]", number, count)
			}
		}
		//mlog.Debug("single%v", single)

		//mlog.Debug("牌组大小[%d]", len(result.Ct))
		w, c := lgc.CalWeight(result, single, userid)
		mlog.Debug("计算權重[w=%d,c=%d]", w, c)
		//mlog.Debug("------------------")
		if cMin > c {
			wMax = w
			cMin = c
			resultBest = result
			singleBest = single
		} else if c == cMin {
			if wMax < w {
				wMax = w
				cMin = c
				resultBest = result
				singleBest = single
				//singleBest = append(singleBest, single...)
			}
		}
	}
	resultBest.SingleCards = singleBest
	resultBest.hands = cMin
	//mlog.Debug("--------Best----------")
	mlog.Debug("牌组大小[%d]", len(resultBest.Ct))
	mlog.Debug("计算所得[w=%d,c=%d]", wMax, cMin)
	//mlog.Debug("single%v", resultBest.SingleCards)
	// for _, prinData := range resultBest.Ct {
	// 	mlog.Debug("resultBest%v", prinData)
	// }
	return resultBest
}

//用map計算出現牌點數出現的次數
func mapCardNumberCount(cg *poker.CardGroup) map[int32]int32 {
	mapCard := make(map[int32]int32, 0)
	for _, cd := range cg.Cards {
		mapCard[cd.CardNumber]++
	}
	return mapCard
}

//重新拆牌
func (lgc *DDZRobotLogic) recombination(userid int32, cg *poker.CardGroup) {
	//重新拆牌
	dc := NewDefineCard(lgc.getGameLogic())
	result := dc.Combination(cg)
	ctg := new(robotUserInfoLocal)
	ctg.cardTypesCount = make(map[cardType]int32, 0)
	ctg.ctg = result

	for _, data := range result.Ct {
		ctg.cardTypesCount[data.byteCardType]++
	}
	ctg.cardTypesCount[cardTypeSingle] += int32(len(result.SingleCards))
	mlog.Debug("機器人牌组大小[%d]", len(ctg.ctg.Ct))
	lgc.robotUser[userid] = ctg
}

//跟新result,erasePair,eraseThree
func (lgc *DDZRobotLogic) updataResult(result *CtGroup, erasePair, eraseThree []int32) *CtGroup {
	for _, cd := range erasePair {
		for _, ctg := range result.Ct {
			if ctg.byteCardType == cardTypeDouble && ctg.nMax == cd {
				ctg.byteCardType = cardTypeNone
			}
		}
	}
	for _, cd := range eraseThree {
		for _, ctg := range result.Ct {
			if ctg.byteCardType == cardTypeThree && ctg.nMax == cd {
				ctg.byteCardType = cardTypeNone
			}
		}
	}
	resultNew := new(CtGroup)
	resultNew.Ct = make([]*cardTypeGroup, 0)
	for _, ctg := range result.Ct {
		if ctg.byteCardType != cardTypeNone {
			ctg2 := new(cardTypeGroup)
			ctg2.byteCardType = ctg.byteCardType
			ctg2.nMax = ctg.nMax
			ctg2.nMin = ctg.nMin
			resultNew.Ct = append(resultNew.Ct, ctg2)
		}
	}
	return resultNew

}

// 随机拆牌
func (lgc *DDZRobotLogic) combinationCardsIndex(mapCard map[int32]int32, label cardType, result *CtGroup, index int32) map[int32]int32 {

	if index == 1 { //原本優先順序是先順子再炸彈
		mapCard = lgc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = lgc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = lgc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = lgc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = lgc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 2 {
		mapCard = lgc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = lgc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = lgc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = lgc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = lgc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 3 {
		mapCard = lgc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = lgc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = lgc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = lgc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = lgc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 4 {
		mapCard = lgc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = lgc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = lgc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = lgc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = lgc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 5 {
		mapCard = lgc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = lgc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = lgc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = lgc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = lgc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 6 {
		mapCard = lgc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = lgc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = lgc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = lgc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = lgc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		return mapCard
	}
	//原本炸彈放在第三位
	mapCard = lgc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
	mapCard = lgc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
	mapCard = lgc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
	mapCard = lgc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
	mapCard = lgc.combinationCards(mapCard, cardTypeThree, result)         //三条
	mapCard = lgc.combinationCards(mapCard, cardTypeDouble, result)        //对子

	return mapCard

}

//按类型组牌
func (lgc *DDZRobotLogic) combinationCards(mapCard map[int32]int32, label cardType, result *CtGroup) map[int32]int32 {
	if result == nil {
		return mapCard
	}
	if label == cardTypeDouble {
		for number, count := range mapCard {
			if count == 2 {
				//mlog.Debug("找到对子[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeDouble
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeThree {
		for number, count := range mapCard {
			if count == 3 {
				//mlog.Debug("找到三条[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeThree
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeBomb {
		for number, count := range mapCard {
			if count == 4 {
				//mlog.Debug("找到炸弹[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeBomb
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeStraight { //顺子
		lgc.checkStraightLong(mapCard, 1, 5, cardTypeStraight, result)
	} else if label == cardTypeThreeStraight { //连队
		lgc.checkStraightLong(mapCard, 2, 3, cardTypeThreeStraight, result)
	} else if label == cardTypePlaneAndWing { //飞机
		lgc.checkStraightLong(mapCard, 3, 2, cardTypePlaneAndWing, result)
	}

	return mapCard
}

//检测顺子并剔除
func (lgc *DDZRobotLogic) checkStraight(mapCard map[int32]int32, x, y int32, result *CtGroup) map[int32]int32 {

	bflag := true
	for i := x; i < x+y; i++ {
		if mapCard[i] < 1 {
			bflag = false
		}
	}
	if bflag {
		//mlog.Debug("找到一个顺子[%d~%d]", x, x+y-1)
		ctg := new(cardTypeGroup)
		ctg.byteCardType = cardTypeStraight
		ctg.nMax = x + y - 1
		ctg.nMin = x
		result.Ct = append(result.Ct, ctg)
		for i := x; i < x+y; i++ {
			mapCard[i] -= 1
		}
		mapCard = lgc.checkStraight(mapCard, x, y, result)
	}
	return mapCard
}

//连对,顺子，飞机,size表示至少要几连
func (lgc *DDZRobotLogic) checkStraightLong(mapCard map[int32]int32, zhang, size int32, bType cardType, result *CtGroup) map[int32]int32 {
	count := int32(0)
	min := int32(2)
	max := int32(2)
	for i := int32(2); i <= 15; i++ {
		//mlog.Debug("前count[%d],min[%d],max[%d],i[%d]", count, min, max, i)
		if mapCard[i] >= zhang {
			count++
			max = i
		} else {
			if count >= size {
				//mlog.Debug("找到[%d,%d]一个[%d~%d]", zhang, size, min, max)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = bType
				ctg.nMax = max
				ctg.nMin = min
				result.Ct = append(result.Ct, ctg)
				for j := min; j <= max; j++ {
					mapCard[j] -= zhang
				}
				mapCard = lgc.checkStraightLong(mapCard, zhang, size, bType, result)
			}
			//mlog.Debug("重置count")
			count = 0
			min = i + 1
		}
		//mlog.Debug("后count[%d],min[%d],max[%d],i[%d]", count, min, max, i)
	}
	return mapCard
}

//排序
func (lgc *DDZRobotLogic) sort(single []int32) []int32 {
	for i := 0; i < len(single); i++ {
		for j := i + 1; j < len(single); j++ {
			if single[i] > single[j] {
				single[i], single[j] = single[j], single[i]
			}
		}
	}
	return single
}

//拆王
func (lgc *DDZRobotLogic) combinationKing(mapCard map[int32]int32, result *CtGroup) map[int32]int32 {

	if mapCard[poker.PokerNumberSmallKing]+mapCard[poker.PokerNumberBigKing] == 2 {
		//mlog.Debug("找到火箭")
		ctg := new(cardTypeGroup)
		ctg.byteCardType = cardTypeRocket
		ctg.nMax = poker.PokerNumberBigKing
		ctg.nMin = poker.PokerNumberSmallKing
		result.Ct = append(result.Ct, ctg)
	} else {
		if mapCard[poker.PokerNumberSmallKing] > 0 {
			//mlog.Debug("找到小王")
			ctg := new(cardTypeGroup)
			ctg.byteCardType = cardTypeSingle
			ctg.nMax = poker.PokerNumberSmallKing
			result.Ct = append(result.Ct, ctg)
		}
		if mapCard[poker.PokerNumberBigKing] > 0 {
			//mlog.Debug("找到大王")
			ctg := new(cardTypeGroup)
			ctg.byteCardType = cardTypeSingle
			ctg.nMax = poker.PokerNumberBigKing
			result.Ct = append(result.Ct, ctg)
		}
	}
	mapCard[poker.PokerNumberSmallKing] = 0
	mapCard[poker.PokerNumberBigKing] = 0
	return mapCard
}

//拆2
func (lgc *DDZRobotLogic) combinationTwo(mapCard map[int32]int32, result *CtGroup) map[int32]int32 {
	//mlog.Debug("拆分2")
	if mapCard[2] > 0 {
		ctg := new(cardTypeGroup)
		ctg.nMax = 2
		if mapCard[2] == 1 {
			//mlog.Debug("找到一个2")
			ctg.byteCardType = cardTypeSingle
		} else if mapCard[2] == 2 {
			//mlog.Debug("找到一对2")
			ctg.byteCardType = cardTypeDouble
		} else if mapCard[2] == 3 {
			//mlog.Debug("找到三个2")
			ctg.byteCardType = cardTypeThree
		} else if mapCard[2] == 4 {
			//mlog.Debug("找到炸弹2")
			ctg.byteCardType = cardTypeBomb
		}
		result.Ct = append(result.Ct, ctg)
	}
	mapCard[2] = 0
	return mapCard
}

//获取手数最小的牌 ，如果手数一样获取权值最大的牌
func (lgc *DDZRobotLogic) combinationSingle(mapCard map[int32]int32) (mapReturn map[int32]int32, sg []int32) {
	single := make([]int32, 0)
	for number, count := range mapCard {
		if count == 1 {
			single = append(single, number)
			mapCard[number] = 0
		}
	}
	single = lgc.sort(single)
	return mapCard, single
}

//提取对子和三张
func (lgc *DDZRobotLogic) getPairAndThree(result *CtGroup) (pair, three []int32) {

	pairList := make([]int32, 0)
	threeList := make([]int32, 0)
	for _, ctg := range result.Ct {
		//三带
		if ctg.nMax == 2 {
			continue
		}
		if ctg.byteCardType == cardTypeThree {
			threeList = append(threeList, ctg.nMax)
			continue
		}
		//对子
		if ctg.byteCardType == cardTypeDouble {
			pairList = append(pairList, ctg.nMax)
			continue
		}
	}
	return pairList, threeList

}

//计算权值
func (lgc *DDZRobotLogic) CalWeight(result *CtGroup, single []int32, userid int32) (w, c int32) {

	totalweight := int32(0)
	need := int32(0)
	pairSize := int32(0)
	count := int32(0)
	//
	huojian := int32(7)
	zhadan := int32(7)
	feiji := int32(6)
	feijiplus := int32(3)
	liandui := int32(5)
	lianduiplus := int32(2)

	shunzi := int32(4)
	shunziplus := int32(1)
	santiao := int32(3)
	duizi := int32(2)
	danpai := int32(1)

	gx := lgc.getRobotCha(userid)
	if gx != nil {
		danpai = gx.Data["danpai"]
		duizi = gx.Data["duizi"]
		santiao = gx.Data["santiao"]
		shunzi = gx.Data["shunzi"]
		shunziplus = gx.Data["shunziplus"]
		liandui = gx.Data["liandui"]
		lianduiplus = gx.Data["lianduiplus"]
		feiji = gx.Data["feiji"]
		feijiplus = gx.Data["feijiplus"]
		zhadan = gx.Data["zhadan"]
		huojian = gx.Data["huojian"]
	}
	for _, ctg := range result.Ct {
		//火箭
		if ctg.byteCardType == cardTypeRocket {
			totalweight += huojian
			count++
			continue
		}
		//炸弹
		if ctg.byteCardType == cardTypeBomb {
			totalweight += zhadan
			count++
			continue
		}
		//飞机
		if ctg.byteCardType == cardTypePlaneAndWing {
			size := ctg.nMax - ctg.nMin + 1
			s := feiji + (size-2)*feijiplus
			totalweight += s
			need += 2
			count++
			continue
		}
		//连对
		if ctg.byteCardType == cardTypeThreeStraight {
			size := ctg.nMax - ctg.nMin + 1
			s := liandui + (size-3)*lianduiplus
			totalweight += s
			count++
			continue
		}
		//顺子
		if ctg.byteCardType == cardTypeStraight {
			size := ctg.nMax - ctg.nMin + 1
			s := shunzi + (size-5)*shunziplus
			totalweight += s
			count++
			continue
		}
		//三带
		if ctg.byteCardType == cardTypeThree {
			totalweight += santiao
			need += 1
			count++
			continue
		}
		//对子
		if ctg.byteCardType == cardTypeDouble {
			pairSize += 1
			continue
		}
		//2 王
		if ctg.byteCardType == cardTypeSingle {
			totalweight += danpai
			count++
			continue
		}
	}
	singleSize := int32(len(single))
	if singleSize-need >= 0 {
		totalweight += (singleSize - need) * danpai
		count += (singleSize - need)
		need = 0
	} else {
		need = need - singleSize
	}
	if pairSize-need >= 0 {
		totalweight += (pairSize - need) * duizi
		count += (pairSize - need)
	}
	return totalweight, count
}

//cal big weight计算大牌的权值分数
func (lgc *DDZRobotLogic) CalBigWeigth(cg *poker.CardGroup, userid int32) int32 {
	w := int32(0)

	jiaodizhuhuojian := int32(8)
	jiaodizhuzhadan := int32(6)
	jiaodizhudawang := int32(4)
	jiaodizhuxiaowang := int32(3)
	jiaodizhu2 := int32(2)
	gx := lgc.getRobotCha(userid)
	if gx != nil {
		jiaodizhuhuojian = gx.Data["jiaodizhuhuojian"]
		jiaodizhuzhadan = gx.Data["jiaodizhuzhadan"]
		jiaodizhudawang = gx.Data["jiaodizhudawang"]
		jiaodizhuxiaowang = gx.Data["jiaodizhuxiaowang"]
		jiaodizhu2 = gx.Data["jiaodizhu2"]
	}

	mapcard := make(map[int32]int32, 0)
	for _, cd := range cg.Cards {
		mapcard[cd.CardNumber]++
	}
	for card, count := range mapcard {
		if count == 4 {
			w += jiaodizhuzhadan
			continue
		}
		if card == 2 {
			w += (jiaodizhu2 * count)
			continue
		}
	}
	if mapcard[poker.PokerNumberSmallKing] > 0 && mapcard[poker.PokerNumberBigKing] > 0 {
		w += jiaodizhuhuojian
		return w
	}
	if mapcard[poker.PokerNumberSmallKing] > 0 {
		w += jiaodizhuxiaowang
		return w
	}
	if mapcard[poker.PokerNumberBigKing] > 0 {
		w += jiaodizhudawang
		return w
	}
	return w
}

//isCanRob
func (lgc *DDZRobotLogic) isCanRob(w, userid int32) bool {
	qiangdizhu := int32(9)
	gx := lgc.getRobotCha(userid)
	if gx != nil {
		qiangdizhu = gx.Data["qiangdizhu"]
	}
	if w >= qiangdizhu {
		return true
	}
	return false
}

//是否可以叫地主
func (lgc *DDZRobotLogic) isCanjiaodizhu(w, userid int32) bool {
	jiaodizhumax := int32(9)
	jiaodizhumin := int32(6)
	jiaodizhup := int32(50)
	gx := lgc.getRobotCha(userid)
	if gx != nil {
		jiaodizhumax = gx.Data["jiaodizhumax"]
		jiaodizhumin = gx.Data["jiaodizhumin"]
		jiaodizhup = gx.Data["jiaodizhup"]
	}
	if w >= jiaodizhumax {
		return true
	}
	if w >= jiaodizhumin {
		numRand := util.RndFloat() * 100
		if jiaodizhup >= int32(numRand) {
			return true
		}
		return false
	}
	return false
}

//明牌
func (lgc *DDZRobotLogic) isMingpai(w, userid int32) bool {
	dizhumingpai := int32(10)
	gx := lgc.getRobotCha(userid)
	if gx != nil {
		dizhumingpai = gx.Data["dizhumingpai"]
	}
	if w >= dizhumingpai {
		return true
	}
	return false
}

//明牌
func (lgc *DDZRobotLogic) getUserHand(userid int32) int32 {
	result := lgc.robotUser[userid]
	if result != nil {
		return result.ctg.hands
	}
	return 20
}

// 判斷機器人AIlevel與Flag關係
func checkAIWithFlag(robot *robotcomm.Robot, Flag int) bool {
	switch robot.RebotLevel {
	case 1:
		return false
	case 2:
		if Flag == RobotAIFlag1 {
			return true
		}
	case 3, 4:
		if Flag == RobotAIFlag1 || Flag == RobotAIFlag2 {
			return true
		}
	case 5, 6:
		if Flag == RobotAIFlag1 || Flag == RobotAIFlag2 || Flag == RobotAIFlag3 {
			return true
		}
	case 7, 8:
		if Flag == RobotAIFlag1 || Flag == RobotAIFlag2 || Flag == RobotAIFlag3 || Flag == RobotAIFlag4 {
			return true
		}
	case 9:
		if Flag == RobotAIFlag1 || Flag == RobotAIFlag2 || Flag == RobotAIFlag3 || Flag == RobotAIFlag4 || Flag == RobotAIFlag5 {
			return true
		}
	case 10:
		if Flag == RobotAIFlag1 || Flag == RobotAIFlag2 || Flag == RobotAIFlag3 || Flag == RobotAIFlag4 || Flag == RobotAIFlag5 || Flag == RobotAIFlag6 {
			return true
		}
	default:
		return false
	}

	return false
}
