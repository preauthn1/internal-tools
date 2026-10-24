package wzp

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"

	proto "example.invalid/golang/protobuf/proto"
)

//游戏结果
type wzpGameRet struct {
	glc *WZPGameLogic
}

func newWZPGameRet(glc *WZPGameLogic) *wzpGameRet {
	gret := new(wzpGameRet)
	gret.glc = glc

	return gret
}

func (ret *wzpGameRet) getGameRet() *netproto.WZP_GameRet {
	wzpret := new(netproto.WZP_GameRet)
	wzpret.GameBetAmount = proto.Int64(ret.glc.gameBetLog.totalAmount)
	wzpret.GiveUp = proto.Int32(0)

	//游戏以弃牌结束
	if ret.glc.turnManager.isAllGiveUp() {
		wzpret.GiveUp = proto.Int32(1)

		mlog.Debug("gameret :%v", ret.glc.gameBetLog.userBets)

		for uid, amount := range ret.glc.gameBetLog.userBets {
			ur := ret.glc.GetUserByUserID(uid)
			uret := new(netproto.WZP_UserRet)
			uret.BetAmount = proto.Int64(ret.glc.gameBetLog.userBets[uid])
			uret.UserID = proto.Int32(uid)
			uret.SeatNo = proto.Int32(ur.SeatNo)

			wzpret.UserRets = append(wzpret.UserRets, uret)

			//弃牌
			if ret.glc.turnManager.isGiveUp(uid) {
				mlog.Debug("gameret:用户%d弃牌", uid)
				uret.WinLose = proto.Int64(-1 * amount)
				uret.CardType = proto.Int32(11)

				//				uret.UserPoker = ret.glc.gamePoker.GetUserPoker(uid).GetNetMessageData()
				//				for i, _ := range uret.UserPoker.Cards {
				//					uret.UserPoker.Cards[i] = poker.NewEmptyPokerCard().GetNetMessageData()
				//					uret.UserPoker.Cards[i].CardNumber = proto.Int32(1)
				//				}
			} else {
				mlog.Debug("gameret:用户%d未弃牌", uid)
				wzpret.WinUserID = proto.Int32(uid)
				uret.WinLose = proto.Int64(ret.glc.gameBetLog.totalAmount - ret.glc.gameBetLog.userBets[uid])
			}
			// fmt.Println("弃牌结束", uid, uret.GetBetAmount(), uret.GetWinLose(), ret.glc.gameBetLog.totalAmount, ret.glc.gameBetLog.userBets[uid])
		}

		return wzpret
	}

	//过滤弃牌用户
	gpk := poker.NewGamePoker()
	for _, upk := range ret.glc.gamePoker.UserPokers {
		if ret.glc.turnManager.isGiveUp(upk.User.UserID) {
			continue
		}

		gpk.UserPokers = append(gpk.UserPokers, upk)
	}

	max, ctypes := getGameMaxCard(gpk, 5, ret.glc)

	for uid, amount := range ret.glc.gameBetLog.userBets {
		ur := ret.glc.GetUserByUserID(uid)
		uret := new(netproto.WZP_UserRet)
		uret.BetAmount = proto.Int64(ret.glc.gameBetLog.userBets[uid])
		uret.UserID = proto.Int32(uid)
		uret.SeatNo = proto.Int32(ur.SeatNo)

		wzpret.UserRets = append(wzpret.UserRets, uret)

		//弃牌
		if ret.glc.turnManager.isGiveUp(uid) {
			uret.WinLose = proto.Int64(-1 * amount)
			uret.CardType = proto.Int32(11)

			//			uret.UserPoker = ret.glc.gamePoker.GetUserPoker(uid).GetNetMessageData()
			//			for i, _ := range uret.UserPoker.Cards {
			//				uret.UserPoker.Cards[i] = poker.NewEmptyPokerCard().GetNetMessageData()
			//				uret.UserPoker.Cards[i].CardNumber = proto.Int32(1)
			//			}

			continue
		}

		uret.CardType = proto.Int32(ctypes[uid])
		uret.UserPoker = ret.glc.gamePoker.GetUserPoker(uid).GetNetMessageData()

		if uid == max.User.UserID {
			uret.WinLose = proto.Int64(ret.glc.gameBetLog.totalAmount - ret.glc.gameBetLog.userBets[uid])
			wzpret.WinUserID = proto.Int32(uid)
		} else {
			uret.WinLose = proto.Int64(-1 * amount)
		}

		// fmt.Println("结束", uid, uret.GetBetAmount(), uret.GetWinLose(), ret.glc.gameBetLog.totalAmount, ret.glc.gameBetLog.userBets[uid])
	}

	return wzpret
}
