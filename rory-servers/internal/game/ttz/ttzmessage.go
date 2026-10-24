//==================================================================
//文件名： ttzmessage.go
//说明：   推筒子广播和消息
//修改记录：
//==================================================================
package ttz

import (
	proto "example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"	
)

import (
	"fmt"
	"time"	
)

//发送游戏设置
func sendGameSettings(ttz *TtzLogic, user *gamelogic.GameUser) {
	mlog.Debug("G2C_message发送当前游戏设置,SeatNo[%d]", user.SeatNo)
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastTTZGameSettingsID)

	gsdata := new(netproto.TTZ_GameSettings)
	gsdata.ReadyWaitTime = proto.Int32(int32(ttz.ttzCfgs.GetValueInt("ReadyWaitTime")))
	gsdata.BaseScore = proto.Int32(int32(ttz.ttzCfgs.GetValueInt("BaseScore")))
	gsmsg.MsgData = gsdata
	ttz.SendMessageToUser(user, gsmsg)
}

//发送当前游戏场景
func sendGameScene(ttz *TtzLogic, user *gamelogic.GameUser) {
	mlog.Debug("G2C_message发送当前游戏场景,SeatNo[%d]", user.SeatNo)
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastTTZRecoverID)

	gsdata := new(netproto.TTZ_GameSence)
	{
		gsdata.CurGameState = proto.Int32(ttz.nGameState)
		nowTime := int32(time.Now().Unix())
		var nTimeValue int32 = 0
		if ttz.nGameState != GAMEPROGRESS_SETTLE {
			nTimeValue = ttz.nAutoBeginTime + ttz.getOutTime() - nowTime
		}
		if nTimeValue < 0 {
			nTimeValue = 0
		}
		gsdata.OutTime = proto.Int32(nTimeValue)
		var nButtonSeatNo int32 = ttz.nButtonSeatNo
		for _, ud := range ttz.mapGameUserData {
			if !ud.bIsRob {
				nButtonSeatNo = -1
			}
		}
		udTmp := ttz.getUserInfoBySeatNo(user.SeatNo)
		if udTmp == nil {
			mlog.Debug("用户为观战玩家,uderid[%d]", user.UserID)
			gsdata.MaxRobRate = proto.Int32(0)
			gsdata.MaxBetRate = proto.Int32(0)

		} else {
			gsdata.MaxRobRate = proto.Int32(udTmp.nMaxRob)
			gsdata.MaxBetRate = proto.Int32(udTmp.nMaxBet)
		}

		gsdata.ButtonSeatNo = proto.Int32(nButtonSeatNo)
		gsdata.RobRate = proto.Int32(ttz.nButtonValue)
		for seatNo, ud := range ttz.mapGameUserData {
			userData := new(netproto.TTZ_SenceUserInfo)
			userData.SeatNo = proto.Int32(seatNo)
			userData.CardType = proto.Int32(-1)
			userData.CardCount = proto.Int32(-1)
			userData.BetRate = proto.Int32(ud.nRateBet)
			userData.RobRate = proto.Int32(ud.nRateRob)
			userData.IsRob = proto.Bool(ud.bIsRob)
			userData.IsBet = proto.Bool(ud.bIsBet)
			userData.IsLight = proto.Bool(ud.bIsLight)
			userData.WinScore = proto.Int64(ud.nCurRoundWinMoney)
			if user.SeatNo == seatNo || ud.bIsLight || ttz.nGameState >= GAMEPROGRESS_LIGHT {
				userData.CardType = proto.Int32(ud.nCardType)
				userData.CardCount = proto.Int32(ud.nCardCount)
				mlog.Debug("card size[%d]", ud.userCards.Len())
				for _, cd := range ud.userCards.Cards {
					card := new(netproto.PokerCard)
					card.Color = proto.Int32(cd.Color)
					card.CardNumber = proto.Int32(cd.CardNumber)
					userData.Cards = append(userData.Cards, card)
				}
			}
			gsdata.UserInfo = append(gsdata.UserInfo, userData)
		}
	}
	mlog.Debug("sendGameScene[%v]", gsdata)
	gsmsg.MsgData = gsdata

	ttz.SendMessageToUser(user, gsmsg)
}

// 广播玩家操作错误信息
func BroadCastErrMsg(ttz *TtzLogic, nSeatNo, nErrType int32, strErrInfo string) {
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("[TtzLogic]战斗逻辑操作 失敗,G2C_message广播错误消息，Seatno：%d,ErrType:%d,ErrInfo:%s。", nSeatNo, nErrType, strErrInfo)).Error("HallServer壓測時")
	mlog.Debug("G2C_message广播错误消息，Seatno：%d,ErrType:%d,ErrInfo:%s。", nSeatNo, nErrType, strErrInfo)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DZPK_GameMessageClassID_BroadCastErrorID)

	pdata := new(netproto.DZPK_ErrorMsg)
	pdata.ErrorType = proto.Int32(nErrType)
	pdata.ErrorInfo = proto.String(strErrInfo)
	pmsg.MsgData = pdata
	ttz.SendMessgeToSeat(nSeatNo, pmsg)
}

//广播庄家信息
func BroadCastButtonInfo(ttz *TtzLogic) {
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	mlog.Debug("G2C_message广播庄家消息，Seatno：%d，Value：%d。", ttz.nButtonSeatNo, ttz.nButtonValue)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastRotResultID)

	pdata := new(netproto.TTZ_BroadCastRobResult)
	pdata.ButtonSeatNo = proto.Int32(ttz.nButtonSeatNo)
	pdata.ActionValue = proto.Int32(ttz.nButtonValue)
	pmsg.MsgData = pdata
	ttz.SendMessageToUsers(ttz.AllTableUser(), pmsg)
}

//发牌
func sendDealMsg(ttz *TtzLogic, seatNo int32, cg *poker.CardGroup) {
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	mlog.Debug("G2C_message,sendDealMsg,seatNo=%d,card size=%d.", seatNo, cg.Len())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastDealInfoID)

	pdata := new(netproto.TTZ_BroadCastDealInfo)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.CardSize = proto.Int32(int32(cg.Len()))
	pdata.No = proto.Int32(ttz.nDealIndex)
	for _, cd := range cg.Cards {
		card := new(netproto.PokerCard)
		card.Color = proto.Int32(cd.Color)
		card.CardNumber = proto.Int32(cd.CardNumber)
		pdata.Cards = append(pdata.Cards, card)
	}
	pmsg.MsgData = pdata
	ttz.SendMessgeToSeat(seatNo, pmsg)

}

// 广播玩家观战信息
func sendDealMsgWatch(ttz *TtzLogic, seatNo, cardSize int32) {
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	mlog.Debug("G2C_message,sendDealMsgWatch,seatNo=%d,card size=%d.", seatNo, cardSize)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastDealInfoID)

	pdata := new(netproto.TTZ_BroadCastDealInfo)
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.CardSize = proto.Int32(cardSize)
	pdata.No = proto.Int32(ttz.nDealIndex)
	pmsg.MsgData = pdata
	ttz.SendMessgeToSeat(seatNo, pmsg)
}

//广播牌型消息
func BroadCastUpdateCT(ttz *TtzLogic) {
	mlog.Debug("G2C_message,BroadCastUpdateCT")
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastUpdateCTID)

	pdata := new(netproto.TTZ_BroadCastUpdateCT)
	for seatNo, ud := range ttz.mapGameUserData {
		pdata.SeatNo = proto.Int32(seatNo)
		pdata.CardType = proto.Int32(ud.nCardType)
		pdata.CardCount = proto.Int32(ud.nCardCount)
		pmsg.MsgData = pdata
		ttz.SendMessageToUser(ud.User, pmsg)
	}
}

// 广播玩家亮牌信息
func BroadCastLight(ttz *TtzLogic, seatNo int32) {
	mlog.Debug("G2C_message,BroadCastLight")
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastLightCardsID)

	pdata := new(netproto.TTZ_BroadCastLightCards)
	ud := ttz.mapGameUserData[seatNo]
	if ud == nil {
		mlog.Error("ud is nil")
		return
	}
	pdata.SeatNo = proto.Int32(seatNo)
	pdata.CardType = proto.Int32(ud.nCardType)
	pdata.CardCount = proto.Int32(ud.nCardCount)
	for _, cd := range ud.userCards.Cards {
		card := new(netproto.PokerCard)
		card.Color = proto.Int32(cd.Color)
		card.CardNumber = proto.Int32(cd.CardNumber)
		pdata.Cards = append(pdata.Cards, card)
	}
	pmsg.MsgData = pdata
	ttz.SendMessageToUsers(ttz.AllTableUser(), pmsg)
}

// 广播结算信息
func BroadCastSettleMentTotal(ttz *TtzLogic) {
	mlog.Debug("G2C_message,BroadCastTotalSettleMent")
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastSettleMentTotal)

	pdata := new(netproto.TTZ_BrodCastSettleMentTotal)
	for seatNo, ud := range ttz.mapGameUserData {
		pSettleMentInfo := new(netproto.TTZ_SettleMentInfo)
		pSettleMentInfo.SeatNo = proto.Int32(seatNo)
		pSettleMentInfo.Score = proto.Int64(ud.nCurGameWinMoney)
		pdata.UserScoreInfo = append(pdata.UserScoreInfo, pSettleMentInfo)
	}
	pmsg.MsgData = pdata
	ttz.SendMessageToUsers(ttz.AllTableUser(), pmsg)

}

//新的一局开始
func BroadCastBeginAndDice() {

}

// 广播玩家操作信息
func BroadCastActionMsg(ttz *TtzLogic, nSeatNo, nActionType, nActionValue int32) {
	mlog.Debug("G2C_message,BroadCastActionMsg")
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastActionResultID)

	pdata := new(netproto.TTZ_BroadCastActionResult)
	pdata.ActionSeatNo = proto.Int32(nSeatNo)
	pdata.ActionType = proto.Int32(nActionType)
	pdata.ActionValue = proto.Int32(nActionValue)
	pmsg.MsgData = pdata
	ttz.SendMessageToUsers(ttz.AllTableUser(), pmsg)

}

//小局结算
func BroadCastTotalSettleRound(ttz *TtzLogic) {
	mlog.Debug("G2C_message,BroadCastTotalSettleRound")
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastSettleMentRound)

	pdata := new(netproto.TTZ_BrodCastSettleMentRound)
	pdata.ButtonSeatNo = proto.Int32(ttz.nButtonSeatNo)
	userButton := ttz.mapGameUserData[ttz.nButtonSeatNo]
	if userButton != nil {
		pdata.ButtonTotalScore = proto.Int64(userButton.nCurRoundWinMoney)
	}
	for _, stInfo := range ttz.stSettlementInfo {
		pSettleMentInfo := new(netproto.TTZ_SettleMentInfo)
		pSettleMentInfo.SeatNo = proto.Int32(stInfo.SeatNo)
		pSettleMentInfo.Score = proto.Int64(stInfo.Score)
		pSettleMentInfo.ButtonScore = proto.Int64(stInfo.ButtonScore)
		pdata.UserScoreInfo = append(pdata.UserScoreInfo, pSettleMentInfo)
	}
	pmsg.MsgData = pdata
	ttz.SendMessageToUsers(ttz.AllTableUser(), pmsg)
}

// 广播游戏状态
func BroadCastGameState(ttz *TtzLogic, outtime int32) {
	mlog.Debug("G2C_message,BroadCastGameState,state=%d,outtime=%d", ttz.nGameState, outtime)
	if ttz == nil {
		mlog.Error("ttz is nil.")
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastGameState)

	pdata := new(netproto.TTZ_BroadCastGameState)

	pdata.CurGameState = proto.Int32(ttz.nGameState)
	pdata.OutTime = proto.Int32(outtime)

	pmsg.MsgData = pdata

	for _, user := range ttz.AllTableUser() {
		ud := ttz.getUserInfoBySeatNo(user.SeatNo)
		if ud != nil {
			if ttz.nGameState == GAMEPROGRESS_LIGHT && ud.bIsLight {
				continue
			}
		}
		ttz.SendMessageToUser(user, pmsg)
	}

	//ttz.SendMessageToUsers(ttz.AllTableUser(), pmsg)
}

// 发送玩家最大操作倍率
func BroadCastUserMaxOpRate(ttz *TtzLogic, user *gamelogic.GameUser, nRate, nType int32) {
	mlog.Debug("G2C_message,BroadCastUserMaxOpRate，Seatno:%d, nRate:%d, nType:%d", user.SeatNo, nRate, nType)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.TTZ_GameMessageClassID_BroadCastTTZMaxOpRateID)

	pdata := new(netproto.TTZ_MaxOpRate)
	pdata.MaxRate = proto.Int32(nRate)
	pdata.Type = proto.Int32(nType)
	pmsg.MsgData = pdata
	ttz.SendMessageToUser(user, pmsg)
}
