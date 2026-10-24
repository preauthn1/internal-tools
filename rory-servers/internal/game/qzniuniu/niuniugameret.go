package qznn

import (
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

import (
	"fmt"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

//牛牛牌局结果
type nnGameRet struct {
	nnlgc *QZNNGameLogic
}

//创建牛牛游戏结果对象
func newNNGameRet(nnglc *QZNNGameLogic) *nnGameRet {
	nnret := new(nnGameRet)
	nnret.nnlgc = nnglc

	return nnret
}

//获取游戏结果
func (nnret *nnGameRet) getGameRet(qznn *QZNNGameLogic) *netproto.QZNN_GameRet {
	gret := new(netproto.QZNN_GameRet)

	//庄家牌型
	zjuid := nnret.nnlgc.zhuangInfo.userID

	//庄家牌型
	zjpoker := nnret.nnlgc.gamePoker.GetUserPoker(zjuid)

	//庄家输赢
	var zhuangwinlose int64 = 0

	zjcardret := getNiuNiuRet(zjpoker.Cards)
	gret.ZhuangRet = new(netproto.QZNN_UserRet)
	gret.ZhuangRet.NiuNiuRet = proto.Int32(int32(zjcardret.niuRet))       //庄家牛牛结果
	gret.ZhuangRet.NiuNiuRetNumber = proto.Int32(int32(zjcardret.number)) //庄家牛牛点数
	gret.ZhuangRet.NiuNiuBetTimes = proto.Int32(0)                        //庄家倍数
	gret.ZhuangRet.WinLose = proto.Int64(zhuangwinlose)                   //庄家输赢
	gret.ZhuangRet.User = nnret.nnlgc.GetUserByUserID(zjuid).GetProtoBufGameUser()
	gret.ZhuangRet.CardGroup = zjcardret.getNetCardGroup()

	//闲家牌型
	for _, xjbt := range nnret.nnlgc.xianjiaBets.userBets {
		xjpoker := nnret.nnlgc.gamePoker.GetUserPoker(xjbt.userID)

		xjcardret := getNiuNiuRet(xjpoker.Cards)
		betret := compareCards(xjcardret, zjcardret)

		//闲家输赢金币数
		//xjwinlose := int64(nnret.nnlgc.zhuangInfo.betTimes) * int64(xjbt.betTimes) * int64(betret) * int64(nnret.nnlgc.GetServerConfig().BaseScore)
		xjwinlose, bonus := nnret.getWinLoseOrJackpot(qznn.jackpotTotalPool, xjcardret, xjbt, betret, qznn.zhuangInfo.betTimes)
		qznn.updateVipUserWin(xjbt.userID, xjwinlose, bonus)
		xjuser := nnret.nnlgc.GetUserByUserID(xjbt.userID)
		if !xjuser.IsRobot() {
			qznn.UpdateGameJackpotPoolFunc(xjbt.userID, -bonus)
			qznn.sendJackpotToBox(xjbt.userID, bonus) //TODO xjwinlose + bonus 送到保險箱 and 郵件
		}

		//庄家输赢金币数
		zhuangwinlose += -1 * xjwinlose
		xjuserRet := new(netproto.QZNN_UserRet)
		xjuserRet.WinLose = proto.Int64(xjwinlose)
		xjuserRet.User = nnret.nnlgc.GetUserByUserID(xjbt.userID).GetProtoBufGameUser()
		xjuserRet.NiuNiuRet = proto.Int32(int32(xjcardret.niuRet))
		xjuserRet.NiuNiuRetNumber = proto.Int32(int32(xjcardret.number))
		xjuserRet.NiuNiuBetTimes = proto.Int32(betret)
		xjuserRet.JackpotBonus = proto.Int64(bonus)
		xjuserRet.CardGroup = xjcardret.getNetCardGroup()

		gret.UserRet = append(gret.UserRet, xjuserRet)
		if xjuser := nnret.nnlgc.GetUserByUserID(xjbt.userID); xjuser != nil {
			var flowAmount, betAmount, payAmount int64
			amount, tax := qznn.GetUserWinAndTax(xjwinlose)
			flowAmount = util.ABSInt64(amount) //流水=派彩
			betAmount = util.ABSInt64(amount) + tax //有效投注=輸贏+稅
			payAmount = betAmount + amount //返邏金額=有效投注+派彩
			xjuser.FlowAmount += flowAmount
			xjuser.BetAmount += betAmount
			xjuser.PayAmount += payAmount
			retmsg := "用戶%d, 流水%d, 有效投注 %d, 返還金額 %d"
			if xjuser.IsRobot() {
				retmsg = "[機器人]" + retmsg
			}
			nnret.nnlgc.GameLog.AddGameStep(retmsg, xjbt.userID, flowAmount, betAmount, payAmount)
		}

	}

	_, zjbonus := nnret.getWinLoseOrJackpot(qznn.jackpotTotalPool, zjcardret, nil, 0, qznn.zhuangInfo.betTimes)
	qznn.updateVipUserWin(zjuid, zhuangwinlose, zjbonus)
	xjuser := nnret.nnlgc.GetUserByUserID(zjuid)
	if !xjuser.IsRobot() {
		qznn.UpdateGameJackpotPoolFunc(zjuid, -zjbonus)
		qznn.sendJackpotToBox(zjuid, zjbonus)
	}

	gret.ZhuangRet.WinLose = proto.Int64(zhuangwinlose) //庄家输赢
	gret.ZhuangRet.JackpotBonus = proto.Int64(zjbonus)
	gret.UserRet = append(gret.UserRet, gret.ZhuangRet)
	if zjuser := nnret.nnlgc.GetUserByUserID(zjuid); zjuser != nil {
		var flowAmount, betAmount, payAmount int64
		//if zhuangwinlose > 0 {
		//	flowAmountAfterTax, _ := qznn.GetUserWinAndTax(zhuangwinlose)
		//	flowAmount = flowAmountAfterTax
		//} else {
		//	flowAmount = -zhuangwinlose
		//}
		amount, tax := qznn.GetUserWinAndTax(zhuangwinlose)
		flowAmount = util.ABSInt64(amount) //流水=派彩
		betAmount = util.ABSInt64(amount) + tax //有效投注=輸贏+稅
		payAmount = betAmount + amount //返邏金額=有效投注+派彩
		zjuser.FlowAmount += flowAmount
		zjuser.BetAmount += betAmount
		zjuser.PayAmount += payAmount
		retmsg := "用戶%d, 流水%d, 有效投注 %d, 返還金額 %d"
		if zjuser.IsRobot() {
			retmsg = "[機器人]" + retmsg
		}
		nnret.nnlgc.GameLog.AddGameStep(retmsg, zjuid, flowAmount, betAmount, payAmount)
	}

	return gret
}

//計算輸贏與jackpot得獎
func (nnret *nnGameRet) getWinLoseOrJackpot(jackpotTotalPool int64, xjcardret *niuniuCardRet, xjbt *userBet, betret int32, zhuangBetTimes int32) (int64, int64) {

	if xjbt == nil {
		xjbt = new(userBet)
		xjbt.betTimes = 0
	}

	if xjcardret.niuRet >= nnCard_Niu {
		return int64(xjbt.betTimes) * int64(betret) * int64(zhuangBetTimes) , 0
	}

	jpKind := new(JPKind)
	rndJackpotKind := util.RndInt(0,100)	//隨機獎金種類
	sumPoolRate := int(0)
	for _, v := range nnconfig.jackpotConfig.jackpotKind {
		sumPoolRate += int(v.getsyspoolrate)
		if rndJackpotKind <= sumPoolRate {
			jpKind = v
			break
		}
	}

	bonusScale := int32(0)
	rndBonus := util.RndInt(0, 100)	//隨機獎金比例
	sumBonusRate := int(0)
	for _, v := range jpKind.bonusking {
		sumBonusRate += int(v.getrate)
		if rndBonus <= sumBonusRate {
			bonusScale = v.scale
			break
		}
	}

	c1 := jackpotTotalPool * int64(jpKind.syspoolrate) / 100
	c2 := c1 * int64(bonusScale) / 100
	newbonus := c2 * int64(nnconfig.jackpotConfig.roomamountrate) / 100

	if newbonus >= int64(nnconfig.jackpotConfig.maxbonus) {
		newbonus = int64(nnconfig.jackpotConfig.maxbonus)
	}

	newwinlose := int64(xjbt.betTimes) * int64(betret) * int64(zhuangBetTimes)

	log1 := fmt.Sprintf("用戶%d, 分配前彩金池總額[%v] 選中彩金池比例[%v] 選中獎金比例[%v] 房間彩金額度[%v] 最大獎金限制[%v]",
		xjbt.userID, jackpotTotalPool, jpKind.syspoolrate, bonusScale, nnconfig.jackpotConfig.roomamountrate, nnconfig.jackpotConfig.maxbonus)
	log2 := fmt.Sprintf("用戶%d, 輸贏[%v] 獲得彩金[%v]", xjbt.userID, newwinlose, newbonus)

	nnret.nnlgc.GameLog.AddGameStep(log1)
	nnret.nnlgc.GameLog.AddGameStep(log2)
	mlog.Debug("%v", log1)
	mlog.Debug("%v", log2)

	return newwinlose, newbonus
}

//玩家獲得的獎金存到保險箱
func (qznn *QZNNGameLogic) sendJackpotToBox(userid int32, bonus int64) {
	if bonus <= 0 {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_GameUpdateUserJackpotBox)
	ust := new(netproto.UserJackpotToBox)
	ust.UserID = proto.Int32(userid)
	ust.UpdateAmount = proto.Int64(bonus)
	pmsg.MsgData = ust


	logmsg := fmt.Sprintf("userid[%v]玩家獲得的獎金存到保險箱[%v]", userid, bonus)
	//nnret.nnlgc.GameLog.AddGameStep(logmsg)
	mlog.System(logmsg)

	qznn.SendRPCMessageFunc(true, pmsg)
}
