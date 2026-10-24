//==================================================================
//文件名： ttzrobotlogic.go
//说明：   推筒子发牌
//修改记录：
//==================================================================
package ttz

import (
	"fmt"
	"strconv"
	"strings"

	proto "example.invalid/golang/protobuf/proto"
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

//机器人个性
type RobotCha struct {
	Data map[string]*rateInfo
}

type rateInfo struct {
	robmin int
	robmax int
	betmin int
	betmax int
}

//机器人逻辑
type TTZRobotLogic struct {
	*robotlogic.LogicBase
	robotChas []*RobotCha
	gameChas  map[int32]*RobotCha

	robotCards map[int32]*poker.CardGroup
}

// 创建机器人逻辑
func NewTTZRobotLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *TTZRobotLogic {
	rblgc := new(TTZRobotLogic)
	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)
	rblgc.parseRobotCha(configs)
	rblgc.gameChas = make(map[int32]*RobotCha)
	rblgc.robotCards = make(map[int32]*poker.CardGroup)

	return rblgc
}

// 加载玩家牌信息
func (lgc *TTZRobotLogic) loadCards(userid int32) {
	data := lgc.getPoker().GetUserCards(userid)
	lgc.robotCards[userid] = data
}

//解析机器人个性
func (lgc *TTZRobotLogic) parseRobotCha(cfgs *config.Group) {

	robotCfg := cfgs.GetConfig(RobotLogicConfigFileName)

	cn := robotCfg.GetConfigInt("csize", "cn")
	//mlog.Debug("解析机器人配置文件,cn[%d]", cn)
	for i := 1; i <= cn; i++ { //个性 //c1 c2 c3 ...
		cname := fmt.Sprintf("c%d", i)

		rData := new(RobotCha)
		rData.Data = make(map[string]*rateInfo, 0)
		for j := 1; j <= 3; j++ {
			sname := fmt.Sprintf("s%d", j)
			pstr := robotCfg.GetConfigValue(cname, sname)
			if pstr == "" {
				mlog.Error("配置文件错误,cname[%s],sname[%s]", cname, sname)
				break
			}
			rateInfoData := parseStrNewRate(pstr, ",")
			if rateInfoData != nil {
				rData.Data[sname] = rateInfoData
				//mlog.Debug("保存机器人个性，名称[%s][%s]", cname, sname)
				continue
			}
			mlog.Error("保存个性失败，名称[%s][%s]", cname, sname)
		}
		//mlog.Debug("保存机器人个性，名称[%s]", cname)
		lgc.robotChas = append(lgc.robotChas, rData)
	}
}

// 解析机器人个性
func parseStrNewRate(str string, sep string) *rateInfo {
	//mlog.Debug("解析配置文件")
	str = util.TrimAllSpace(str)
	sep = util.TrimAllSpace(sep)

	sarr := strings.Split(str, sep)

	robmin, err := strconv.Atoi(sarr[0])
	if err != nil {
		mlog.Error("配置有误1")
		return nil
	}

	robmax, err := strconv.Atoi(sarr[1])
	if err != nil {
		mlog.Error("配置有误2")
		return nil
	}
	if robmin > robmax {
		mlog.Error("配置有误，最小值大于最大值")
		return nil
	}

	betmin, err := strconv.Atoi(sarr[2])
	if err != nil {
		mlog.Error("配置有误3")
		return nil
	}

	betmax, err := strconv.Atoi(sarr[3])
	if err != nil {
		mlog.Error("配置有误4")
		return nil
	}

	if betmin > betmax {
		mlog.Error("配置有误，最小值大于最大值")
		return nil
	}
	rate := new(rateInfo)
	rate.robmin = robmin
	rate.robmax = robmax
	rate.betmin = betmin
	rate.betmax = betmax
	//mlog.Debug("robmin[%d],robmax[%d],betmin[%d],betmax[%d]", robmin, robmax, betmin, betmax)
	return rate

}

//获取机器人个性
func (lgc *TTZRobotLogic) getRobotCha(userid int32) *RobotCha {
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
func (lgc *TTZRobotLogic) getGameLogic() *TtzLogic {
	return lgc.GameLogic.(*TtzLogic)
}

// 获取游戏牌池
func (lgc *TTZRobotLogic) getPoker() *poker.FixedPoker {
	return lgc.getGameLogic().pokerGroup.(*poker.FixedPoker)
}

//游戏消息
func (lgc *TTZRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {
	//底注阶段,收到底注倍数通知时
	//example.invalid("机器人收到信息。GameState[%s],SC[%d]", getState(lgc.getGameLogic().nGameState), int32(netproto.TTZ_GameMessageClassID_BroadCastGameState))
	if lgc.getGameLogic().nGameState == GAMEPROGRESS_ROBBUTTON && m.SClassID == int32(netproto.TTZ_GameMessageClassID_BroadCastGameState) {
		lgc.robotAction(robot, m, "抢庄", ACTION_ROB)
	} else if lgc.getGameLogic().nGameState == GAMEPROGRESS_CHOOSEBET && m.SClassID == int32(netproto.TTZ_GameMessageClassID_BroadCastGameState) && robot.UserInfo.SeatNo != lgc.getGameLogic().nButtonSeatNo {
		lgc.robotAction(robot, m, "押注", ACTION_BET)
	} else if lgc.getGameLogic().nGameState == GAMEPROGRESS_LIGHT && m.SClassID == int32(netproto.TTZ_GameMessageClassID_BroadCastGameState) {
		lgc.robotAction(robot, m, "亮牌", ACTION_LIGHT)
	}
	/* else if lgc.getGameLogic().nGameState == GAMEPROGRESS_NONE {
		lgc.loadCards(robot.UserInfo.UserID)
	}
	*/
}

// 游戏结束
func (lgc *TTZRobotLogic) OnGameEnd(robot *robotcomm.Robot) {
	lgc.gameChas = make(map[int32]*RobotCha)
}

// 机器人操作逻辑
func (lgc *TTZRobotLogic) robotAction(robot *robotcomm.Robot, m *msg.Message, name string, action int32) {
	//获取用户牌型
	ucards := lgc.getPoker().GetUserCards(robot.UserInfo.UserID)
	//ucards := lgc.robotCards[robot.UserInfo.UserID]
	if ucards == nil {
		mlog.Error("没有手牌")
		return
	}
	totalCards := new(poker.CardGroup)
	totalCards.AddCardGroup(ucards)
	if ucards.Len() < 2 {
		ud := lgc.getGameLogic().getUserInfoBySeatNo(robot.UserInfo.SeatNo)
		if ud == nil {
			mlog.Error("机器人用户信息获取失败")
			return
		}
		totalCards.AddCardGroup(ud.userCards)
	}
	mlog.Debug("机器人牌%s", totalCards.String())
	//解析用户牌型
	nType, _ := getCardTypeAndCount(totalCards)

	//获取机器人个性
	cha := lgc.getRobotCha(robot.UserInfo.UserID)

	if len(cha.Data) < 3 {
		mlog.Error("配置有误，机器人个性信息有误size[%d]", len(cha.Data))
		return
	}

	var bt int = 0
	if action == ACTION_ROB {
		if nType == HIGH {
			if cha.Data["s1"].robmax-cha.Data["s1"].robmin > 0 {
				bt = util.RndInt(cha.Data["s1"].robmin, cha.Data["s1"].robmax)
			} else {
				bt = cha.Data["s1"].robmax
			}
		} else if nType == PAIR {
			if cha.Data["s2"].robmax-cha.Data["s2"].robmin > 0 {
				bt = util.RndInt(cha.Data["s2"].robmin, cha.Data["s2"].robmax)
			} else {
				bt = cha.Data["s2"].robmax
			}
		} else if nType == LEOPARD {
			if cha.Data["s3"].robmax-cha.Data["s3"].robmin > 0 {
				bt = util.RndInt(cha.Data["s3"].robmin, cha.Data["s3"].robmax)
			} else {
				bt = cha.Data["s3"].robmax
			}
		}
		if bt > 4 {
			bt = 4
		}
	} else if action == ACTION_BET {
		bt = 2
		if nType == HIGH {
			if cha.Data["s1"].betmax-cha.Data["s1"].betmin > 0 {
				bt = util.RndInt(cha.Data["s1"].betmin, cha.Data["s1"].betmax)
			} else {
				bt = cha.Data["s1"].betmax
			}
		} else if nType == PAIR {
			if cha.Data["s2"].betmax-cha.Data["s2"].betmin > 0 {
				bt = util.RndInt(cha.Data["s2"].betmin, cha.Data["s2"].betmax)
			} else {
				bt = cha.Data["s2"].betmax
			}
		} else if nType == LEOPARD {
			if cha.Data["s3"].betmax-cha.Data["s3"].betmin > 0 {
				bt = util.RndInt(cha.Data["s3"].betmin, cha.Data["s3"].betmax)
			} else {
				bt = cha.Data["s3"].betmax
			}
		}
		bt = bt / 2
		bt = bt * 2
		if bt < 2 {
			bt = 2
		}
		if bt > 8 {
			bt = 8
		}

	} else if action == ACTION_LIGHT {

	} else {
		mlog.Error("操作类型错误[%d]", action)
		return
	}

	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_ActionOprateID)

	btdata := new(netproto.TTZ_ActionInfo)
	btdata.ActionType = proto.Int32(action)
	btdata.ActionValue = proto.Int32(int32(bt))
	btmsg.MsgData, _ = proto.Marshal(btdata)

	lgc.SendGameMessage(robot.UserInfo, btmsg)

	mlog.Debug("机器人%v%s倍数%d", robot.UserInfo, name, bt)
}
