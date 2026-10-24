//==================================================================
//文件名： dzpkpotinfo.go
//说明：   德州扑克筹码彩池以及结算系统，主要负责筹码的计算分配
//修改记录：
//==================================================================

package dzpk

import (
	"bytes"
	"fmt"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

//彩池结构信息
type stJettonInfo struct {
	nJetton int64   //筹码
	nAllin  int64   //allin值
	partin  []int32 //参与玩家座位号
}

//彩池
type potPool struct {
	data []*stJettonInfo
}

// 彩池信息转换为string
func (pp *potPool) String() string {

	var buf bytes.Buffer
	for index, Info := range pp.data {
		fmtInfo := fmt.Sprintf("\n奖池编号[%d],池中筹码[%d],池中玩家座位号%v;", index+1, Info.nJetton, Info.partin)
		buf.WriteString(fmtInfo)
	}
	return buf.String()

}

//********奖池结构体**************************************
type potUserInfo struct {
	SeatNo    int32
	WinJetton int64
	CardType  int32
	CardGroup *poker.CardGroup
}
type potInfo struct {
	PotIndex int32
	PotTotal int64
	UserInfo []*potUserInfo
}

// 奖池信息转化为字符
func (pi *potInfo) String() string {

	var buf bytes.Buffer
	fmtInfo1 := fmt.Sprintf("奖池编号[%d],池中筹码[%d]->", pi.PotIndex+1, pi.PotTotal)
	buf.WriteString(fmtInfo1)
	for _, data := range pi.UserInfo {
		fmtInfo2 := fmt.Sprintf("赢家座位号[%d],税前筹码[%d],", data.SeatNo, data.WinJetton)
		buf.WriteString(fmtInfo2)
	}
	return buf.String()

}

type settleMentRound map[int32]*potInfo

//**********************************************

//计算彩池
func calPot(gameUserData []*DzpkUserData) *potPool {
	//mlog.Debug("计算彩池...")
	//获取
	pots := new(potPool) //彩池筹码
	var userCal sliceUserData
	var nCount int32 = 0
	var nAllin, nJettonTmp int64 = 0, 0
	userCal.data = gameUserData
	nTotalGameCount := int32(len(userCal.data))
	userCal.ReSort()

	for _, user := range userCal.data {
		//过滤相同数据
		if 0 == user.currBet-nAllin {
			nCount++
			continue
		}
		if user.isAllin {
			stJt := new(stJettonInfo)
			stJt.nJetton = nJettonTmp + (user.currBet-nAllin)*int64(nTotalGameCount-nCount)
			stJt.nAllin = user.currBet
			nAllin = user.currBet

			nJettonTmp = 0
			pots.data = append(pots.data, stJt)
		} else {
			nJettonTmp = nJettonTmp + (user.currBet - nAllin)
		}
		nCount++
	}
	//最后一个池子
	if nJettonTmp != 0 {
		stJt := new(stJettonInfo)
		stJt.nJetton = nJettonTmp
		stJt.nAllin = userCal.data[nTotalGameCount-1].currBet
		pots.data = append(pots.data, stJt)
	}
	//加入参与人
	for index, jettonInfo := range pots.data {
		for _, user := range userCal.data {
			if !user.isFold && user.currBet >= jettonInfo.nAllin {
				pots.data[index].partin = append(pots.data[index].partin, user.seatNo)
			}
		}
	}
	mlog.Debug("彩池大小=%d", len(pots.data))
	return pots
}

//结算
func settlement(dzpk *DzpkLogic) {

	dzpk.settlePotInfo = make(map[int32]*potInfo)
	dzpk.nAnimationDy = dzpk.actionGameEndTime
	for nPotIndex, stjetton := range dzpk.potArray.data {

		potInfoTmp := new(potInfo) //结算信息
		seatNoMax := stjetton.partin[0]
		//获取最大牌型和座位号
		for _, seatNo := range stjetton.partin {
			user := dzpk.getUserInfoBySeatNo(seatNo)
			userMax := dzpk.getUserInfoBySeatNo(seatNoMax)
			if user == nil || userMax == nil {
				return
			}
			user.andUserCardAndPublic(dzpk.publicCards)
			userMax.andUserCardAndPublic(dzpk.publicCards)
			if ComPareCards(userMax.pokerTotalGroup, user.pokerTotalGroup) < 0 {
				seatNoMax = seatNo
			}
		}
		//获取最大牌
		var sliceWinner []int32 //胜利玩家seatNo
		for _, seatNo := range stjetton.partin {
			user := dzpk.getUserInfoBySeatNo(seatNo)
			userMax := dzpk.getUserInfoBySeatNo(seatNoMax)
			if user == nil || userMax == nil {
				return
			}
			if ComPareCards(userMax.pokerTotalGroup, user.pokerTotalGroup) <= 0 {
				sliceWinner = append(sliceWinner, seatNo)
			}
		}
		if len(sliceWinner) == 0 {
			sliceWinner = append(sliceWinner, seatNoMax)
		}
		//分配筹码
		nAverage := stjetton.nJetton / int64(len(sliceWinner)) //平均筹码
		nSurplus := stjetton.nJetton % int64(len(sliceWinner)) //除去能够平分的剩余筹码
		for _, seatNo := range sliceWinner {
			potUserInfoTmp := new(potUserInfo)
			potUserInfoTmp.SeatNo = seatNo
			potUserInfoTmp.WinJetton = nAverage
			user := dzpk.getUserInfoBySeatNo(seatNo)
			if user != nil {
				user.curRoundSettle += nAverage
				dzpk.nMaxCardTypeSeatNo = seatNo
				user.getMaxCardsGroup(dzpk.publicCards)
				potUserInfoTmp.CardType = user.cardType
				potUserInfoTmp.CardGroup = user.pokerMaxCardGroup
				//mlog.Debug("CardType=%d,CardGroup=%d", potUserInfoTmp.CardType, potUserInfoTmp.CardGroup.Len())
			}
			dzpk.nAnimationDy += 3
			potInfoTmp.UserInfo = append(potInfoTmp.UserInfo, potUserInfoTmp)
		}
		//剩余几个筹码给某个玩家
		userOne := dzpk.getUserInfoBySeatNo(sliceWinner[0])
		if userOne != nil {
			userOne.curRoundSettle += nSurplus
			potInfoTmp.UserInfo[0].WinJetton += nSurplus
		}
		potInfoTmp.PotIndex = int32(nPotIndex)
		potInfoTmp.PotTotal = stjetton.nJetton
		dzpk.settlePotInfo[int32(nPotIndex)] = potInfoTmp
	}
	for _, data := range dzpk.settlePotInfo {
		dzpk.GameLog.AddGameStep("奖池分配情况：%s", data.String())
	}

	//汇总
	summarizingJetton(dzpk)

}

//汇总
func summarizingJetton(dzpk *DzpkLogic) {
	for _, user := range dzpk.gameUserData {
		var flowAmount, betAmount, payAmount int64
		amount, tax := dzpk.GetUserWinAndTax(user.curRoundSettle - user.currBet)
		if amount > 0 {
			user.userMoney += amount
			user.User.FlowAmount += amount
		} else {
			user.User.FlowAmount -= amount
		}
		flowAmount = util.ABSInt64(amount) //流水=派彩
		betAmount = util.ABSInt64(amount) + tax //有效投注=輸贏+稅
		payAmount = betAmount + amount //返邏金額=有效投注+派彩
		user.User.FlowAmount += flowAmount
		user.User.BetAmount += betAmount
		user.User.PayAmount += payAmount
		user.curRoundSettle = amount
		user.tax = tax
		dzpk.GameLog.AddGameStep("用戶[%d] 結算金額, win %d, 稅收 %d, 流水 %d, 有效投注 %d, 返還金額 %d", user.User.UserID, amount, tax, flowAmount, betAmount, payAmount)
		dzpk.UpdateUserScoreFunc(user.User.UserID, amount, tax)
		user.curTurnBet = 0
	}
}

////广播（总剩余筹码，当局赢得筹码）
