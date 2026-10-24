//==================================================================
//文件名： dzpkmessage.go
//说明：   德州扑克广播和消息
//修改记录：
//==================================================================
package dzpk

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"time"
)
import proto "example.invalid/golang/protobuf/proto"

//发送等待消息
func BroadCastWaitOptMsg(dzpk *DzpkLogic) {
	mlog.Debug("G2C_message发送玩家等待消息，座位号：%d，等待时间：%d。", dzpk.currActionSeat, dzpk.optTime)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastWaitOptID)

	//当前可操作玩家
	pdata := new(netproto.DZPK_WaitOptSeat)
	pdata.SeatNo = proto.Int32(dzpk.currActionSeat)
	pdata.OptTime = proto.Int32(dzpk.optTime + 2)
	pdata.TotalJetton = proto.Int64(dzpk.betTotal)
	user := dzpk.getUserInfoBySeatNo(dzpk.currActionSeat)
	if user == nil {
		mlog.Error("ERROR! BroadCastWaitOptMsg and user is nil.")
		return
	}
	for _, action := range user.userOprateInfo {
		ops := new(netproto.DZPK_Operability)
		ops.OpType = proto.Int32(action.nOprateType)
		ops.OpJetton = proto.Int64(action.nOprateJetton)
		pdata.Op = append(pdata.Op, ops)
	}
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//发送自动押注消息
func sendBetBaseMsg(dzpk *DzpkLogic) {
	mlog.Debug("G2C_message发送玩家开局信息")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastBetBaseID)

	pdata := new(netproto.DZPK_BetBase)
	for _, user := range dzpk.gameUserData {
		basedata := new(netproto.DZPK_BetBaseInfo)
		basedata.Betbase = proto.Int64(dzpk.getInitJetton(user))
		basedata.SeatNo = proto.Int32(user.seatNo)
		basedata.Role = proto.Int32(int32(user.getUserRole()))
		pdata.Betbasedata = append(pdata.Betbasedata, basedata)
	}
	pmsg.MsgData = pdata

	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//广播下注消息
func BroadCastBetMsg(dzpk *DzpkLogic, seatNo int32, jetton int64) {
	mlog.Debug("G2C_message广播下注消息，Seatno：%d，Jetton：%d。", seatNo, jetton)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastBetID)

	pdata := new(netproto.DZPK_BrodCastBet)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.Jetton = proto.Int64(jetton)
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//跟注
func BroadCastCallMsg(dzpk *DzpkLogic, seatNo int32, jetton int64) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastCallID)

	pdata := new(netproto.DZPK_BrodCastCall)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.Jetton = proto.Int64(jetton)
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//加注
func BroadCastRaiseMsg(dzpk *DzpkLogic, seatNo int32, jetton int64) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastRaiseID)

	pdata := new(netproto.DZPK_BrodCastRaise)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.Jetton = proto.Int64(jetton)
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//过，看牌
func BroadCastCheckMsg(dzpk *DzpkLogic, seatNo int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastCheckID)

	pdata := new(netproto.DZPK_BrodCastCheck)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.Jetton = proto.Int64(0)
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//Allin
func BroadCastAllinMsg(dzpk *DzpkLogic, seatNo int32, jetton int64) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastAllinID)

	pdata := new(netproto.DZPK_BrodCastAllIn)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.Jetton = proto.Int64(jetton)
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//弃牌
func BroadCastFoldMsg(dzpk *DzpkLogic, seatNo int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastFoldID)

	pdata := new(netproto.DZPK_BrodCastFold)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.Jetton = proto.Int64(0)
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//发牌message，私有牌
func sendDealMsg(dzpk *DzpkLogic) {

	mlog.Debug("G2C_message获取发牌消息")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastDealID)

	for _, user := range dzpk.AllTableUser() {
		pdataPoker := new(netproto.DZPK_PrivatePoker)
		for _, haveCards := range dzpk.gameUserData {
			pdataCard := new(netproto.DZPK_PrivateCard)
			pdataCard.SeatNo = proto.Int32(haveCards.seatNo)
			if user.SeatNo == haveCards.seatNo {
				for _, havaCard := range haveCards.userCards.Cards {
					card := new(netproto.PokerCard)
					card.Color = proto.Int32(havaCard.Color)
					card.CardNumber = proto.Int32(havaCard.CardNumber)
					pdataCard.Cards = append(pdataCard.Cards, card)
				}
			}
			pdataPoker.Cards = append(pdataPoker.Cards, pdataCard)
		}
		pmsg.MsgData = pdataPoker
		//通知用户发牌
		dzpk.SendMessageToUser(user, pmsg)
	}

}

//发公共牌
func BroadCastDealPublicCards(dzpk *DzpkLogic, cards *poker.CardGroup) {
	mlog.Debug("G2C_message发公共牌")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastDealPublicID)

	pdata := new(netproto.DZPK_BrodCastCardsPoker)
	pdata.Type = proto.Int32(1)
	for _, c := range cards.Cards {
		pdata.Cards = append(pdata.Cards, c.GetNetMessageData())
	}
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//广播牌信息
func BroadCastUserCardType(dzpk *DzpkLogic, usd *DzpkUserData) {
	mlog.Debug("G2C_message广播牌信息,cardType=%d", usd.cardType)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastCardTypeInfoID)

	pdata := new(netproto.DZPK_BroadCastCardType)
	pdata.Type = proto.Int32(usd.cardType)
	pmsg.MsgData = pdata
	dzpk.SendMessageToUser(usd.User, pmsg)
}

//彩池
func BroadCastPotMsg(dzpk *DzpkLogic) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastPotInfoID)

	pdata := new(netproto.DZPK_BroadCastPotInfo)
	if dzpk.potArray == nil {
		mlog.Error("BroadCastPotMsg,potArray is nil")
		return
	}
	for _, data := range dzpk.potArray.data {
		pdata.Jetton = append(pdata.Jetton, data.nJetton)
	}
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

// 广播结算信息
func BroadCastRoundSettleMent(dzpk *DzpkLogic) {
	mlog.Debug("G2C_message广播结算信息")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastSettlementID)

	pdata := new(netproto.DZPK_SettleMentRound)

	for nPotIndex, PotInfo := range dzpk.settlePotInfo {
		//pdata.PotData
		pPotData := new(netproto.DZPK_PotInfo)
		pPotData.PotIndex = proto.Int32(nPotIndex)
		for _, PotUserInfo := range PotInfo.UserInfo {
			pPotUserInfo := new(netproto.DZPK_PotUserInfo)
			pPotUserInfo.SeatNo = proto.Int32(PotUserInfo.SeatNo)
			user := dzpk.getUserInfoBySeatNo(PotUserInfo.SeatNo)
			if user != nil {
				if user.tax > 0 && PotUserInfo.WinJetton > user.tax {
					pPotUserInfo.WinJetton = proto.Int64(PotUserInfo.WinJetton - user.tax)
					user.tax = 0
				} else {
					pPotUserInfo.WinJetton = proto.Int64(PotUserInfo.WinJetton)
				}
			} else {
				pPotUserInfo.WinJetton = proto.Int64(PotUserInfo.WinJetton)
			}
			pPotUserInfo.CardType = proto.Int32(-1)
			if dzpk.nFoldNumber+1 < dzpk.gameUserCount {
				pPotUserInfo.CardType = proto.Int32(PotUserInfo.CardType)
				for _, cd := range PotUserInfo.CardGroup.Cards {
					card := new(netproto.PokerCard)
					card.Color = proto.Int32(cd.Color)
					card.CardNumber = proto.Int32(cd.CardNumber)
					pPotUserInfo.Cards = append(pPotUserInfo.Cards, card)
				}
			}
			pPotData.UserInfo = append(pPotData.UserInfo, pPotUserInfo)
		}
		pdata.PotData = append(pdata.PotData, pPotData)
	}
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

func BroadCastNewRound(dzpk *DzpkLogic) {

}

//发送游戏设置
func sendGameSettings(dzpk *DzpkLogic, user *gamelogic.GameUser) {
	mlog.Debug("G2C_message发送游戏设置,Seatno=%d,Bigblind=%d.", user.SeatNo, dzpk.betBase)
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastGameSettingID)

	gsdata := new(netproto.DZPK_GameSetting)
	gsdata.SeatNo = proto.Int32(user.SeatNo)
	gsdata.Bigblind = proto.Int64(dzpk.betBase)
	gsdata.ReadyTime = proto.Int32(int32(dzpk.dzpkCfgs.GetValueInt("ReadyTime")))
	gsmsg.MsgData = gsdata

	dzpk.SendMessageToUser(user, gsmsg)

}

//发送当前游戏场景
func sendGameScene(dzpk *DzpkLogic, user *gamelogic.GameUser) {
	mlog.Debug("G2C_message发送当前游戏场景")
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastGameSenceID)

	gsdata := new(netproto.DZPK_GameSence)

	{
		gsdata.GameState = proto.Int32(int32(dzpk.gameState))
		gsdata.CurOprationSeatNo = proto.Int32(dzpk.currActionSeat)
		gsdata.BeginOutTime = proto.Int32(dzpk.beginOutTime)
		nowTime := int32(time.Now().Unix())
		gsdata.OpWaitTime = proto.Int32(dzpk.beginOutTime + dzpk.optTime - nowTime)
		//公牌
		if dzpk.publicCards != nil {
			mlog.Debug("publicCards size[%d]", dzpk.publicCards.Len())
			for _, cd := range dzpk.publicCards.Cards {
				pubCd := new(netproto.PokerCard)
				pubCd.Color = proto.Int32(cd.Color)
				pubCd.CardNumber = proto.Int32(cd.CardNumber)
				gsdata.PubCards = append(gsdata.PubCards, pubCd)
			}
		}

		//奖池
		if dzpk.potArray != nil {
			for _, data := range dzpk.potArray.data {
				gsdata.Pots = append(gsdata.Pots, data.nJetton)
			}
		}

		//玩家信息
		for _, ud := range dzpk.gameUserData {
			userTmp := new(netproto.DZPK_UsersInfo)
			userTmp.SeatNo = proto.Int32(ud.seatNo)
			userTmp.Role = proto.Int32(int32(ud.userRole))
			userTmp.TotalJetton = proto.Int64(ud.userMoney)
			userTmp.CurBetJetton = proto.Int64(ud.currBet)
			userTmp.CurTurnBetJetton = proto.Int64(ud.curTurnBet)
			userTmp.UserState = proto.Int32(ud.curOpration)
			userTmp.NetProfit = proto.Int64(ud.curRoundSettle - ud.currBet)
			//userTmp.CurWinScore = proto.Int32(ud.curRoundSettle)
			if ud.isFold {
				userTmp.UserState = proto.Int32(int32(FOLD))
			}
			if ud.isAllin {
				userTmp.UserState = proto.Int32(int32(ALLIN))
			}
			userTmp.CardType = proto.Int32(-1)
			if (!ud.isFold && dzpk.gameState == GAMEPROGRESS_END) || (ud.seatNo == user.SeatNo) {
				if !ud.isFold || dzpk.gameState != GAMEPROGRESS_END {
					userTmp.CardType = proto.Int32(ud.cardType)
				}
				if ud.userCards != nil {
					for _, cd := range ud.userCards.Cards {
						cards := new(netproto.PokerCard)
						cards.Color = proto.Int32(cd.Color)
						cards.CardNumber = proto.Int32(cd.CardNumber)
						userTmp.Cards = append(userTmp.Cards, cards)
					}
				}
				op := make([]*OpInfo, 0)
				if ud.seatNo == dzpk.currActionSeat {
					op = ud.userOprateInfo
				} else {
					op = ud.userPreOprateInfo
				}
				for _, op := range op {
					opSend := new(netproto.DZPK_Operability)
					opSend.OpType = proto.Int32(op.nOprateType)
					opSend.OpJetton = proto.Int64(op.nOprateJetton)
					userTmp.Op = append(userTmp.Op, opSend)
				}
			}

			gsdata.Users = append(gsdata.Users, userTmp)
		}
	}
	mlog.Debug("sendGameScene[%v]", gsdata)
	gsmsg.MsgData = gsdata
	dzpk.SendMessageToUser(user, gsmsg)
}

// 关闭错误信息
func BroadCastErrMsg(dzpk *DzpkLogic, nSeatNo, nErrType int32, strErrInfo string) {
	mlog.Debug("G2C_message广播错误消息，Seatno：%d,trErrInfo：%s。", nSeatNo, strErrInfo)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastErrorID)

	pdata := new(netproto.DZPK_ErrorMsg)
	pdata.ErrorType = proto.Int32(nErrType)
	pdata.ErrorInfo = proto.String(strErrInfo)
	pmsg.MsgData = pdata
	dzpk.SendMessgeToSeat(nSeatNo, pmsg)
}

// 广播手牌
func BroadCastHandCards(dzpk *DzpkLogic) {
	mlog.Debug("G2C_message,BroadCastHandCards")
	if dzpk == nil {
		mlog.Error("dzpk is nil")
		return
	}
	if dzpk.nFoldNumber+1 == dzpk.gameUserCount {
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastCheckCardsID)

	pdata := new(netproto.DZPK_UserCardsInfo)
	for _, ud := range dzpk.gameUserData {
		if ud.isFold {
			continue
		}
		handCards := new(netproto.DZPK_UserHandCards)
		handCards.SeatNo = proto.Int32(ud.seatNo)
		handCards.CardType = proto.Int32(ud.cardType)
		for _, cd := range ud.userCards.Cards {
			card := new(netproto.PokerCard)
			card.Color = proto.Int32(cd.Color)
			card.CardNumber = proto.Int32(cd.CardNumber)
			handCards.Cards = append(handCards.Cards, card)
		}
		pdata.Data = append(pdata.Data, handCards)

	}
	pmsg.MsgData = pdata
	dzpk.SendMessageToUsers(dzpk.AllTableUser(), pmsg)
}

//广播预操作信息
func BroadCastPreOpration(dzpk *DzpkLogic) {
	mlog.Debug("G2C_message,BroadCastPreOpration")
	if dzpk == nil {
		mlog.Error("dzpk is nil")
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastPreOprationID)

	for _, ud := range dzpk.gameUserData {
		if ud.seatNo == dzpk.currActionSeat {
			continue
		}
		if ud.isFold || ud.isAllin {
			continue
		}
		pdata := new(netproto.DZPK_PreOpration)
		pdata.SeatNo = proto.Int32(ud.seatNo)
		for _, op := range ud.userPreOprateInfo {
			opSend := new(netproto.DZPK_Operability)
			opSend.OpType = proto.Int32(op.nOprateType)
			opSend.OpJetton = proto.Int64(op.nOprateJetton)
			pdata.Op = append(pdata.Op, opSend)
		}
		pmsg.MsgData = pdata
		if dzpk.currActionSeat != ud.seatNo {
			dzpk.SendMessageToUser(ud.User, pmsg)
		}

	}
}
