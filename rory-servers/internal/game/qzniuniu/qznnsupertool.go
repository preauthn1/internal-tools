package qznn

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

//更新超端分數
func (qznn *QZNNGameLogic) updateVipUserWin(userid int32, winlose int64, jackpotmoney int64) {
	vipuser := qznn.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
	refresh := false
	superToolHandler := qznn.SuperToolService.GetBrGameSuperToolHandler()
	if vipuser != nil && vipuser.GetComplete() == 0 && vipuser.GetStatus() == 1 {
		vipuser.CompleteAmount = proto.Int64(vipuser.GetCompleteAmount() + winlose + jackpotmoney)

		if vipuser.GetAssignAmount() > 0 {
			if vipuser.GetCompleteAmount() >= vipuser.GetAssignAmount() {
				vipuser.Complete = proto.Int32(1)
			}
		} else {
			if vipuser.GetCompleteAmount() <= vipuser.GetAssignAmount() {
				vipuser.Complete = proto.Int32(1)
			}
		}

		mlog.System("update user completeamount %d, userid %d", vipuser.GetCompleteAmount(), vipuser.GetUserID())
		superToolHandler.SaveVipUser(userid)
		refresh = true

	}

	if refresh { //有異動才需要馬上更新
		qznn.SuperToolService.RefreshVIPUser()
	}
}

//監控放殺分
func (qznn *QZNNGameLogic) bowJiangType(user *gamelogic.GameUser, key int) (bool, *poker.UserPoker) {
	result := false
	pk := make([]*poker.Card, 0)
	upk := new(poker.UserPoker)

	//先執行放分
	vipuser := qznn.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(user.UserID)

	if vipuser != nil && vipuser.GetAssignAmount() > 0 && vipuser.GetComplete() == 0 {

		rnd := int32(util.RndInt(0, 100))
		mlog.System("用戶%d 執行放分機率 %d, rnd[%d]", user.UserID, vipuser.GetAssignRate(), rnd)
		if rnd < vipuser.GetAssignRate() { //隨機0-99 小於勝率則執行放分
			niuCardRet := new(niuniuCardRet)
			niuniuRnd := util.RndInt(0, 100)
			if niuniuRnd > 50 {
				niuCardRet = qznn.assignNewNiuNiuCard(qznn.gamePoker.UserPokers[key].Cards, nnCard_NiuNiu, 0, 30)
			} else {
				niuCardRet = qznn.assignNewNiuNiuCardLessThen(qznn.gamePoker.UserPokers[key].Cards, nnCard_Niu, 9, 30)
			}

			pk = niuCardRet.cards.Cards[4:5]
			cg := new(poker.CardGroup)
			cg.Cards = append(cg.Cards, pk...)
			upk = poker.NewUserPoker(user, cg)
			qznn.gamePoker.AddUserPoker(upk)

			result = true
		}
	}

	if vipuser != nil && vipuser.GetAssignAmount() < 0 && vipuser.GetComplete() == 0 {
		niuCardRet := qznn.assignNewNiuNiuCardLessThen(qznn.gamePoker.UserPokers[key].Cards, nnCard_Niu, 5, 30)
		pk = niuCardRet.cards.Cards[4:5]
		cg := new(poker.CardGroup)
		cg.Cards = append(cg.Cards, pk...)
		upk = poker.NewUserPoker(user, cg)
		qznn.gamePoker.AddUserPoker(upk)

		result = true
	}

	return result, upk
}
