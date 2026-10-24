package scmj

import (
	"fmt"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mahjong"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"math"
	"time"
)

//通知客户端发牌信息
func (m *MahjongSC) BoardcastDeal(cards *mahjong.CardGroup, seatNo int32) {
	if cards == nil {
		mlog.Error("BoardcastDeal cards is nil.")
		return
	}
	mlog.Debug("MSG_S2C,广播发牌消息,桌子号[%d],座位号%d,%d张牌.", m.TableNo(), seatNo, cards.Len())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJDealId)

	for _, user := range m.AllTableUser() {
		if cards.Len() > 1 && user.SeatNo != seatNo {
			continue
		}
		if user.IsOffLine() {
			continue
		}
		pdata := new(netproto.SCMJ_BroadcastDeal)
		pdata.SeatNo = proto.Int32(seatNo)
		pdata.Number = proto.Int32(int32(cards.Len()))
		pdata.SurplusNum = proto.Int32(int32(m.randGroup.Len()))
		if seatNo == user.SeatNo {
			for _, cd := range cards.Cards {
				pdata.Cards = append(pdata.Cards, cd)
			}
			// if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havechooserate") == 1 && m.gameState == PROGRESS_DEAL {
			// 	// 二人麻将开始不给牌
			// } else {
			// 	for _, cd := range cards.Cards {
			// 		pdata.Cards = append(pdata.Cards, cd)
			// 	}
			// }
		}
		pmsg.MsgData = pdata
		m.SendMessageToUser(user, pmsg)
	}

}

//摸牌
func (m *MahjongSC) BoardcastTakeCard(cards *mahjong.CardGroup, seatNo int32) {
	if cards == nil {
		mlog.Error("BoardcastTakeCard cards is nil.")
		return
	}
	mlog.Debug("MSG_S2C,广播摸牌消息,桌子号[%d],座位号%d,%d张牌.", m.TableNo(), seatNo, cards.Len())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJDealId)

	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		pdata := new(netproto.SCMJ_BroadcastDeal)
		pdata.SeatNo = proto.Int32(seatNo)
		pdata.Number = proto.Int32(int32(cards.Len()))
		pdata.SurplusNum = proto.Int32(int32(m.randGroup.Len()))
		if seatNo == user.SeatNo {
			for _, cd := range cards.Cards {
				pdata.Cards = append(pdata.Cards, cd)
			}
		}
		pmsg.MsgData = pdata
		m.SendMessageToUser(user, pmsg)
	}

}

//通知客户端选定缺
func (m *MahjongSC) BoardcastChooseColor(seatNo int32) {
	mlog.Debug("MSG_S2C,通知客户端选定缺,桌子号[%d],seatNo[%d].", m.TableNo(), seatNo)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJBroadCastChooseId)

	pdata := new(netproto.SCMJ_BroadcastAction)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.OpType = proto.Int32(OP_CHOOSEDINGQUE)
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

//通知客户端选倍率
func (m *MahjongSC) BoardcastChooseRate(seatNo int32) {
	mlog.Debug("MSG_S2C,通知客户端选倍率,桌子号[%d],seatNo[%d].", m.TableNo(), seatNo)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_S2C_BroadCastRateId)

	pdata := new(netproto.SCMJ_BroadcastAction)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.OpType = proto.Int32(OP_CHOOSERATE)
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

//通知客户端操作信息
func (m *MahjongSC) BoardcastAction(seatNo, opType, huType int32, opValue []int32) {
	mlog.Debug("MSG_S2C,广播玩家操作消息,桌子号[%d],seatNo[%d],type[%d],opValue[%v].", m.TableNo(), seatNo, opType, opValue)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJBroadCastActionId)

	pdata := new(netproto.SCMJ_BroadCastUserOper)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.BeSeatNo = proto.Int32(seatNo)
	var nOpvalue int32 = 0
	if len(opValue) > 0 {
		nOpvalue = opValue[0]
	}
	pdata.BeCard = proto.Int32(nOpvalue)
	pdata.OperType = proto.Int32(opType)
	pdata.HuType = proto.Int32(huType)
	pdata.Cards = opValue
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

//通知客户端操作信息
func (m *MahjongSC) BoardScenceExchange(seatNo, beseatNo, opType, huType int32, opValue []int32) {
	mlog.Debug("MSG_S2C,广播玩家操作消息,桌子号[%d],seatNo[%d],type[%d],opValue[%v].", m.TableNo(), seatNo, opType, opValue)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJScenceExchangeId)

	pdata := new(netproto.SCMJ_BroadCastUserOper)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.BeSeatNo = proto.Int32(beseatNo)
	var nOpvalue int32 = 0
	if len(opValue) > 0 {
		nOpvalue = opValue[0]
	}
	pdata.BeCard = proto.Int32(nOpvalue)
	pdata.OperType = proto.Int32(opType)
	pdata.HuType = proto.Int32(huType)
	pdata.Cards = opValue
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

//广播当前游戏状态
func (m *MahjongSC) BoardcastGameState() {

}

//通知客户端错误信息
func (m *MahjongSC) rspError(user *gamelogic.GameUser, id int32, str string) {
	mlog.Error("MSG_S2C,错误信息,桌子号[%d],SeatNo[%d],ID[%d],Info[%s]", m.TableNo(), user.SeatNo, id, str)
	if user.IsOffLine() {
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJErrorId)

	pdata := new(netproto.SCMJ_Error)
	pdata.ErrorId = proto.Int32(id)
	pdata.ErrorInfo = proto.String(str)
	pmsg.MsgData = pdata

	m.SendMessageToUser(user, pmsg)
}

//通知客户端即时结算信息
func (m *MahjongSC) BoardcastIntimeSettle(nType, nRate int32) {
	mlog.Debug("MSG_S2C,广播玩家即时结算分数,桌子号[%d],nType[%d].", m.TableNo(), nType)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJSettleMentIntimeId)

	pdata := new(netproto.SCMJ_InTimeSettle)

	// 結算類型，寫 log 用 , Added by Joecool Wang at 2019-06-03
	map_strSettlementType := make(map[int32]string)
	map_strSettlementType[0] = "无"
	map_strSettlementType[1] = "刮风"
	map_strSettlementType[2] = "下雨"
	map_strSettlementType[3] = "查花猪"
	map_strSettlementType[4] = "查大叫"
	map_strSettlementType[5] = "胡"
	map_strSettlementType[6] = "自摸胡"
	map_strSettlementType[7] = "点炮胡"
	map_strSettlementType[8] = "抢杠胡"
	map_strSettlementType[9] = "退杠"
	map_strSettlementType[10] = "呼叫转移"

	for _, ud := range m.gameUserData {
		if ud.intimeScore == 0 {
			continue
		}
		usInfo := new(netproto.SCMJ_SettleInfo)
		usInfo.SeatNo = proto.Int32(ud.seatNo)
		usInfo.Type = proto.Int32(nType)
		usInfo.Rate = proto.Int32(nRate)
		usInfo.Score = proto.Int64(ud.intimeScore)
		example.invalid = append(example.invalid, usInfo)

		mlog.Debug("BoardcastIntimeSettle(), UserID[%d], SeatNo[%d], Type[%d(%s)], Rate[%d], Score[%d]】",
			ud.User.UserID, ud.seatNo, nType, map_strSettlementType[nType], nRate, ud.intimeScore)
	}
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

//通知客户端等待信息
func (m *MahjongSC) BroadcastWait() {
	if m.isEnd {
		return
	}
	mlog.Debug("MSG_S2C,BroadcastWait,桌子号[%d]", m.TableNo())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJWaitId)

	pdata := new(netproto.SCMJ_Wait)
	pdata.SeatNo = proto.Int32(m.currOpSeatNo)
	surplus := 0
	if m.randGroup != nil {
		surplus = m.randGroup.Len()
	}

	pdata.OutTime = proto.Int32(m.getOutTimeWait())
	pdata.SurplusNum = proto.Int32(int32(surplus))
	pdata.GameState = proto.Int32(m.gameState)
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}

}

//通知客户端玩家选定缺结果
func (m *MahjongSC) BroadCastChooseFlower() {
	mlog.Debug("MSG_S2C,BroadCastChooseFlower,桌子号[%d],", m.TableNo())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJBroadCastChooseResultId)

	pdata := new(netproto.SCMJ_BroadcastChooseResult)
	for _, ud := range m.gameUserData {
		Info := new(netproto.SCMJ_ChooseInfo)
		Info.SeatNo = proto.Int32(ud.seatNo)
		Info.Type = proto.Int32(ud.color)
		example.invalid = append(example.invalid, Info)
	}
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

//通知客户端玩家选倍率结果
func (m *MahjongSC) BroadCastChooseRate() {
	mlog.Debug("MSG_S2C,BroadCastChooseRate,桌子号[%d],", m.TableNo())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_S2C_BroadCastRateResultId)

	pdata := new(netproto.SCMJ_BroadcastChooseResult)
	for _, ud := range m.gameUserData {
		Info := new(netproto.SCMJ_ChooseInfo)
		Info.SeatNo = proto.Int32(ud.seatNo)
		Info.Type = proto.Int32(ud.rate)
		example.invalid = append(example.invalid, Info)
	}
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

//通知客户端提示信息
func (m *MahjongSC) BroadcastUserCanOpInfo(ud *STUserData) {
	mlog.Debug("MSG_S2C,BroadcastUserCanOpInfo,桌子号[%d],", m.TableNo())
	if ud.User.IsOffLine() {
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJCanOpEventId)

	pdata := new(netproto.SCMJ_OpEventInfo)

	pdata.Chi = proto.Bool(false)
	pdata.Peng = proto.Bool(false)
	pdata.Gang = proto.Bool(false)
	pdata.Hu = proto.Bool(false)
	pdata.Time = proto.Int32(m.getOutTimeWait())
	pdata.OtherValue = proto.Int32(0)
	pdata.GangValue = make([]int32, 0)
	for _, data := range ud.canOpEventInfo.OpInfo {
		mlog.Debug("桌子号[%d],event optype[%d],value[%d]", m.TableNo(), data.nOpType, data.nValue)
		if data.nOpType == OP_CHI {
			pdata.Chi = proto.Bool(true)
			pdata.OtherValue = proto.Int32(data.nValue[0])
			continue
		}
		if data.nOpType == OP_PENG {
			pdata.Peng = proto.Bool(true)
			pdata.OtherValue = proto.Int32(data.nValue[0])
			continue
		}
		if data.nOpType == OP_ANGANG || data.nOpType == OP_YINGANG || data.nOpType == OP_MINGGANG {
			pdata.Gang = proto.Bool(true)
			pdata.GangValue = append(pdata.GangValue, data.nValue...)
			continue
		}
		if data.nOpType == OP_HU || data.nOpType == OP_HU_DIANPAO || data.nOpType == OP_HU_ZIMO || data.nOpType == OP_HU_QIANGGANG {
			pdata.Hu = proto.Bool(true)
			pdata.OtherValue = proto.Int32(data.nValue[0])
			continue
		}
	}

	pmsg.MsgData = pdata

	m.SendMessageToUser(ud.User, pmsg)
}

//通知客户端提示信息@
func (m *MahjongSC) BroadcastUserCanQiangGang(ud *STUserData, cd int32) {
	mlog.Debug("MSG_S2C,BroadcastUserCanOpInfo,桌子号[%d]", m.TableNo())
	if ud.User.IsOffLine() {
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJCanOpEventId)

	pdata := new(netproto.SCMJ_OpEventInfo)

	pdata.Chi = proto.Bool(false)
	pdata.Peng = proto.Bool(false)
	pdata.Gang = proto.Bool(false)
	pdata.Hu = proto.Bool(true)
	pdata.Time = proto.Int32(m.getOutTimeWait())
	pdata.OtherValue = proto.Int32(cd)
	pdata.GangValue = make([]int32, 0)

	pmsg.MsgData = pdata

	m.SendMessageToUser(ud.User, pmsg)
}

//通知客户端色子
func (m *MahjongSC) BroadcastDice() {
	mlog.Debug("MSG_S2C,通知客户端色子信息,桌子号[%d]", m.TableNo())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJThrowDiceId)

	pdata := new(netproto.SCMJ_BroadcastDice)
	pdata.Dice = m.dice
	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

//通知客户端总结算
func (m *MahjongSC) BroadcastTotalSettleMent() {
	mlog.Debug("MSG_S2C,通知客户端总结算信息,桌子号[%d]", m.TableNo())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJTotalSettleMentId)

	pdata := new(netproto.SCMJ_TotalSettleMent)
	rate := math.Pow(float64(2), float64(m.rateSettleMent.RATEMAX))
	pdata.RateMax = proto.Int32(int32(rate))

	var str string

	// 結算類型，寫 log 用, Added by Joecool Wang at 2019-06-03
	map_strSettlementType := make(map[int32]string)
	map_strSettlementType[0] = "无"
	map_strSettlementType[1] = "刮风"
	map_strSettlementType[2] = "下雨"
	map_strSettlementType[3] = "查花猪"
	map_strSettlementType[4] = "查大叫"
	map_strSettlementType[5] = "胡"
	map_strSettlementType[6] = "自摸胡"
	map_strSettlementType[7] = "点炮胡"
	map_strSettlementType[8] = "抢杠胡"
	map_strSettlementType[9] = "退杠"
	map_strSettlementType[10] = "呼叫转移"

	// 胡牌類型，寫 log 用, Added by Joecool Wang at 2019-06-03
	map_strHuType := make(map[int32]string)
	map_strHuType[0] = "无"
	map_strHuType[1] = "屁胡"
	map_strHuType[2] = "大对子"
	map_strHuType[3] = "金勾吊"
	map_strHuType[4] = "七对子"
	map_strHuType[5] = "龙七对"
	map_strHuType[6] = "清一色"
	map_strHuType[7] = "带幺九"
	map_strHuType[8] = "天胡"
	map_strHuType[9] = "地胡"
	map_strHuType[10] = "海底"

	for _, ud := range m.gameUserData {
		pud := new(netproto.SCMJ_SettleTotal)
		pud.SeatNo = proto.Int32(ud.seatNo)
		amount, _ := m.GetUserWinAndTax(ud.winScore)
		pud.TotalScore = proto.Int64(amount)
		mlog.Debug("----- 總結算：桌號[%d], UserID[%d], SeatNo[%d] -----", m.TableNo(), ud.User.UserID, ud.seatNo)

		// 寫 log 用 begin, Added by Joecool Wang at 2019-06-03
		tPengCard := new(mahjong.CardGroup)
		for _, tCard := range ud.pengCards {
			tPengCard.AddCard(tCard.pengCard)
		}

		tGangCard := new(mahjong.CardGroup)
		for _, tCard := range ud.gangCards {
			tGangCard.AddCard(tCard.gangCard)
		}

		tHuCard := new(mahjong.CardGroup)
		for _, iCard := range ud.huCards {
			tHuCard.AddCard(iCard)
		}
		// 寫 log 用 end, Added by Joecool Wang at 2019-06-03
		mlog.Debug("碰牌：%v, 槓牌：%v, 手牌：%v, 胡牌：%v", tPengCard, tGangCard, ud.handCards, tHuCard)

		for _, info := range ud.intimeInfo {
			if !info.valid {
				continue
			}
			//mlog.Debug("桌子号[%d],SeatNo[%d],总结算[%v]", m.TableNo(), ud.seatNo, *info)
			//type STScoreRecord struct {
			//	Type    int32   //结算类型
			//	Score   int64   //结算分数
			//	SeatNo  []int32 //联系座位号
			//	HuTypes []int32 //胡牌类型，只针对结算类型为胡牌的情况
			//	valid   bool    //是否有效
			//	card    int32   //胡牌
			//	Rate    int32   //倍率
			//}
			pinfo := new(netproto.SCMJ_SettleMyInfo)

			pinfo.Type = proto.Int32(info.Type)
			pinfo.Rate = proto.Int32(info.Rate)
			pinfo.HuType = info.HuTypes
			pinfo.MyScore = proto.Int64(info.Score)
			pinfo.FromSeatNo = info.SeatNo
			pinfo.AddRate = proto.Int32(info.Rate)
			example.invalid = append(example.invalid, pinfo)

			str = fmt.Sprintf("Type[%d(%s)], Score[%d], SeatNo[",
				info.Type, map_strSettlementType[info.Type], info.Score)
			for _, iNo := range info.SeatNo {
				str += fmt.Sprintf("%d ", iNo)
			}
			str += "], HuType["
			for _, iType := range info.HuTypes {
				str += fmt.Sprintf("%d(%s) ", iType, map_strHuType[iType])
			}
			str += fmt.Sprintf("], Valid[%t], card[%d], Rate[%d]", info.valid, info.card, info.Rate)
			mlog.Debug("%s", str)
		}
		pdata.Total = append(pdata.Total, pud)
	}

	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

// 亮牌
func (m *MahjongSC) BroadcastLight() {
	mlog.Debug("MSG_S2C,通知客户端亮牌,桌子号[%d],", m.TableNo())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_S2C_SCMJLightId)

	pdata := new(netproto.S2C_SCMJ_Light)

	for _, ud := range m.gameUserData {
		pud := new(netproto.SCMJ_HandCards)
		pud.SeatNo = proto.Int32(ud.seatNo)
		pud.Cards = ud.handCards.Cards
		pdata.Data = append(pdata.Data, pud)
	}

	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

// 提示
func (m *MahjongSC) BroadcastRemind(seatNo int32, data map[int32][]*STRemindInfo) {
	mlog.Debug("MSG_S2C,通知客户端听牌提示,桌子号[%d],", m.TableNo())
	user := m.GetUserBySeatNo(seatNo)
	if user == nil || user.IsOffLine() {
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_S2C_SCMJRemindId)

	pdata := new(netproto.S2C_SCMJ_RemindMsg)

	for outcard, sice := range data {

		prd := new(netproto.SCMJ_RemindData)
		prd.OutCard = proto.Int32(outcard)

		for _, info := range sice {
			mlog.Debug("广播提示:%v", *info)
			pinfo := new(netproto.SCMJ_RemindInfo)
			pinfo.TingCard = proto.Int32(info.TingCard)
			pinfo.Number = proto.Int32(info.Number)
			pinfo.Rate = proto.Int32(info.Rate)
			example.invalid = append(example.invalid, pinfo)
		}
		pdata.Data = append(pdata.Data, prd)

	}

	pmsg.MsgData = pdata

	m.SendMessgeToSeat(seatNo, pmsg)
}

// 托管
func (m *MahjongSC) BroadcastTuoGuan(seatNo int32, state bool) {
	mlog.Debug("MSG_S2C,通知客户端托管,桌子号[%d],,SeatNo[%d]", m.TableNo(), seatNo)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_S2C_TuoGuanId)

	pdata := new(netproto.S2C_SCMJ_TuoGuan)

	pdata.SeatNo = proto.Int32(seatNo)
	pdata.State = proto.Bool(state)

	pmsg.MsgData = pdata
	// m.SendMessageToUsers(m.AllTableUser(), pmsg)
	for _, user := range m.AllTableUser() {
		if user.IsOffLine() {
			continue
		}
		m.SendMessageToUser(user, pmsg)
	}
}

// 游戏设置
func (m *MahjongSC) SendGameSettings(gameuser *gamelogic.GameUser) {
	mlog.Debug("MSG_S2C,通知客户端牌局设置,桌子号[%d]", m.TableNo())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJGameSettingId)

	pdata := new(netproto.Mj_GameSetting)

	pdata.ReadyTime = proto.Int32(m.readywaittime)
	pdata.BaseScore = proto.Int64(m.baseScore)

	pmsg.MsgData = pdata
	m.SendMessageToUser(gameuser, pmsg)
}

// 断线重连
func (m *MahjongSC) SendGameScene(gameuser *gamelogic.GameUser) {
	mlog.Debug("MSG_S2C,发送游戏场景,桌子号[%d]", m.TableNo())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SCMJ_GameMessageClassID_SCMJGameSenceId)

	pdata := new(netproto.SCMJ_GameScence)

	pdata.ButtonSeatNo = proto.Int32(m.buttonSeatNo)
	pdata.GameState = proto.Int32(m.gameState)
	pdata.SurplusNum = proto.Int32(int32(m.randGroup.Len()))

	surplusTime := int32(0)
	if m.timerTimeOutEvent != nil {
		beginTime := m.timerTimeOutEvent.BeginTime.Unix()
		delayTime := m.timerTimeOutEvent.WaitSeconds
		surplusTime = int32(beginTime) + int32(delayTime) - int32(time.Now().Unix())
		// mlog.Debug("事件剩余时间,beginTime[%d],delayTime[%d],surplusTime[%d]", beginTime, delayTime, surplusTime)
	} else if m.timerTimeOutAction != nil {
		beginTime := m.timerTimeOutAction.BeginTime.Unix()
		delayTime := m.timerTimeOutAction.WaitSeconds
		surplusTime = int32(beginTime) + int32(delayTime) - int32(time.Now().Unix())
		// mlog.Debug("动作剩余时间,beginTime[%d],delayTime[%d],surplusTime[%d]", beginTime, delayTime, surplusTime)
	}
	if surplusTime < 0 {
		surplusTime = 0
	}

	pdata.SurplusTime = proto.Int32(surplusTime)
	pdata.CurSeatNo = proto.Int32(m.currOpSeatNo)
	pdata.Remind = new(netproto.S2C_SCMJ_RemindMsg)
	for _, ud := range m.gameUserData {
		pud := new(netproto.SCMJ_UserInfo)
		pud.SeatNo = proto.Int32(ud.seatNo)
		pud.UserState = proto.Bool(ud.isHu)
		pud.DingQue = proto.Int32(ud.color)
		pud.RiverCard = ud.riverCards.Cards
		pud.TuoGuan = proto.Bool(ud.isTuoGuan)
		pud.HuaZhu = proto.Bool(ud.isHuaZhu)
		pud.Dajiao = proto.Bool(ud.isCanTingPai)
		pud.Rate = proto.Int32(ud.rate)
		sumCard := ud.handCards.Len()
		if ud.moPai > 0 {
			sumCard += 1
		}
		pud.HandCount = proto.Int32(int32(sumCard))
		mlog.Debug("游戏场景,桌子号[%d],SeatNo[%d],河牌[%s]", m.TableNo(), ud.seatNo, ud.riverCards.String())
		//碰
		for _, pd := range ud.pengCards {
			if !pd.valid {
				continue
			}
			ppd := new(netproto.SCMJ_OpGroup)
			ppd.Card = append(ppd.Card, pd.pengCard)
			ppd.Card = append(ppd.Card, pd.pengCard)
			ppd.Card = append(ppd.Card, pd.pengCard)
			ppd.Type = proto.Int32(OP_PENG)
			ppd.SeatNo = proto.Int32(pd.seatNo)
			pud.OpCard = append(pud.OpCard, ppd)
		}
		//杠
		for _, gd := range ud.gangCards {
			if !gd.valid {
				continue
			}
			pgd := new(netproto.SCMJ_OpGroup)
			pgd.Card = append(pgd.Card, gd.gangCard)
			pgd.Card = append(pgd.Card, gd.gangCard)
			pgd.Card = append(pgd.Card, gd.gangCard)
			pgd.Card = append(pgd.Card, gd.gangCard)
			pgd.Type = proto.Int32(gd.gangType)
			pgd.SeatNo = proto.Int32(gd.seatNo)
			pud.OpCard = append(pud.OpCard, pgd)
		}
		//胡
		for _, hd := range ud.intimeInfo {
			if hd.Type != TYPE_SETTLEMENT_HU && hd.Type != TYPE_SETTLEMENT_HU_DIANPAO && hd.Type != TYPE_SETTLEMENT_HU_ZIMO && hd.Type != TYPE_SETTLEMENT_HU_QG {
				continue
			}
			if hd.Score <= 0 {
				continue
			}
			phd := new(netproto.SCMJ_OpGroup)
			phd.Card = append(phd.Card, hd.card)
			phd.Type = proto.Int32(OP_HU)
			if hd.Type == TYPE_SETTLEMENT_HU_ZIMO {
				phd.SeatNo = proto.Int32(ud.seatNo)
			} else {
				phd.SeatNo = proto.Int32(hd.SeatNo[0])
			}
			pud.OpCard = append(pud.OpCard, phd)
		}
		mlog.Debug("游戏场景,桌子号[%d],SeatNo[%d],碰杠胡[%v]", m.TableNo(), ud.seatNo, pud.OpCard)
		//即时结算
		if gameuser.SeatNo == ud.seatNo {
			// 只显示自己的手牌
			// if m.Configs.GetConfig("gameserver").GetConfigInt("gameprivate", "havechooserate") == 1 && m.gameState <= PROGRESS_FIGHT {
			// 	//
			// } else {
			// 	pud.HandCard = ud.handCards.Cards
			// 	if ud.moPai > 0 {
			// 		pud.HandCard = append(pud.HandCard, ud.moPai)
			// 	}
			// }
			pud.HandCard = ud.handCards.Cards
			if ud.moPai > 0 {
				pud.HandCard = append(pud.HandCard, ud.moPai)
			}
			//即时结算
			for _, intimeSettlement := range ud.intimeInfo {
				settleMent := new(netproto.SCMJ_SettleMyInfo)

				settleMent.Type = proto.Int32(intimeSettlement.Type)
				settleMent.Rate = proto.Int32(intimeSettlement.Rate)
				settleMent.HuType = intimeSettlement.HuTypes
				settleMent.MyScore = proto.Int64(intimeSettlement.Score)
				settleMent.FromSeatNo = intimeSettlement.SeatNo
				settleMent.AddRate = proto.Int32(intimeSettlement.Rate)
				pdata.InTimeSettle = append(pdata.InTimeSettle, settleMent)
			}
			// 玩家可操作信息
			pCaninfo := new(netproto.SCMJ_OpEventInfo)
			pCaninfo.Chi = proto.Bool(false)
			pCaninfo.Peng = proto.Bool(false)
			pCaninfo.Gang = proto.Bool(false)
			pCaninfo.Hu = proto.Bool(false)
			pCaninfo.Time = proto.Int32(m.getOutTimeWait())
			pCaninfo.OtherValue = proto.Int32(0)
			pCaninfo.GangValue = make([]int32, 0)
			for _, data := range ud.canOpEventInfo.OpInfo {
				if ud.eventFlag {
					break
				}
				mlog.Debug("桌子号[%d],event optype[%d],value[%d]", m.TableNo(), data.nOpType, data.nValue)
				if data.nOpType == OP_CHI {

					pCaninfo.Chi = proto.Bool(true)
					pCaninfo.OtherValue = proto.Int32(data.nValue[0])
					continue
				}
				if data.nOpType == OP_PENG {
					pCaninfo.Peng = proto.Bool(true)
					pCaninfo.OtherValue = proto.Int32(data.nValue[0])
					continue
				}
				if data.nOpType == OP_ANGANG || data.nOpType == OP_YINGANG || data.nOpType == OP_MINGGANG {
					pCaninfo.Gang = proto.Bool(true)
					pCaninfo.GangValue = append(pCaninfo.GangValue, data.nValue...)
					continue
				}
				if data.nOpType == OP_HU || data.nOpType == OP_HU_DIANPAO || data.nOpType == OP_HU_ZIMO || data.nOpType == OP_HU_QIANGGANG {
					pCaninfo.Hu = proto.Bool(true)
					pCaninfo.OtherValue = proto.Int32(data.nValue[0])
					continue
				}
			}
			pdata.CanOp = pCaninfo
			// 提示信息
			if ud.handCards.Len()%3 == 2 || ud.moPai > 0 {
				RemindMsg := new(netproto.S2C_SCMJ_RemindMsg)
				for outcard, sice := range ud.remindData {
					prd := new(netproto.SCMJ_RemindData)
					prd.OutCard = proto.Int32(outcard)
					for _, info := range sice {
						pinfo := new(netproto.SCMJ_RemindInfo)
						pinfo.TingCard = proto.Int32(info.TingCard)
						pinfo.Number = proto.Int32(info.Number)
						pinfo.Rate = proto.Int32(info.Rate)
						example.invalid = append(example.invalid, pinfo)
					}
					RemindMsg.Data = append(RemindMsg.Data, prd)
				}
				pdata.Remind = RemindMsg
			}

		} else {
			if m.gameState > PROGRESS_FIGHT {
				pud.HandCard = ud.handCards.Cards
			}
		}
		pdata.UserInfo = append(pdata.UserInfo, pud)
	}

	pmsg.MsgData = pdata
	m.SendMessageToUser(gameuser, pmsg)
}
