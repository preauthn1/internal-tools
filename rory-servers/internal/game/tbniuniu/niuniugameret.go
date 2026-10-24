package tbnn

import (
	proto "example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

//牛牛牌局结果
type nnGameRet struct {
	nnlgc *TBNNGameLogic
}

//创建牛牛游戏结果对象
func newNNGameRet(nnglc *TBNNGameLogic) *nnGameRet {
	nnret := new(nnGameRet)
	nnret.nnlgc = nnglc

	return nnret
}

func (gameret *nnGameRet) getGameRet() *netproto.TBNN_GameRet {
	nnUserRets := make([]*netproto.TBNN_UserBetRet, 0)

	//赢家结果
	var maxNNRet *netproto.TBNN_UserBetRet
	var maxRet *niuniuCardRet

	//计算每个用户的牌型
	for _, ur := range gameret.nnlgc.AllGameUser() {
		ret := getNiuNiuRet(gameret.nnlgc.gamePoker.GetUserPoker(ur.UserID).Cards)

		//牌局记录
		gameret.nnlgc.GameLog.AddGameStep("用户%d %s", ur.UserID, ret.String())

		nnRet := new(netproto.TBNN_UserBetRet)

		nnRet.NiuNiuRet = proto.Int32(int32(ret.niuRet))
		nnRet.NiuNiuBetTimes = proto.Int32(ret.getBetTimes())

		ubtimes := gameret.nnlgc.gameBet.getUserBetTimes(ur.UserID)
		nnRet.UserBetTimes = proto.Int32(ubtimes)

		nnRet.UserID = proto.Int32(ur.UserID)
		nnRet.NiuNiuNumber = proto.Int32(int32(ret.number))
		nnRet.CardGroup = ret.getNetCardGroup()

		nnUserRets = append(nnUserRets, nnRet)

		if maxNNRet == nil {
			maxNNRet = nnRet
			maxRet = ret
		} else {
			cmpRet := compareCards(maxRet, ret)

			if cmpRet < 0 {
				maxNNRet = nnRet
				maxRet = ret
			}
		}
	}

	//计算输赢
	var winamount int64 = 0

	for _, nnret := range nnUserRets {
		if nnret == maxNNRet {
			continue
		}

		//计算输赢
		var winloseAmount int64 = int64(nnret.GetUserBetTimes()) * int64(maxNNRet.GetNiuNiuBetTimes()) * int64(maxNNRet.GetUserBetTimes()) * int64(gameret.nnlgc.GetServerConfig().BaseScore)

		//设置输家输赢数量
		nnret.WinLose = proto.Int64(-1 * winloseAmount)

		//计算总赢的金额
		winamount += winloseAmount
	}

	//设置赢家输赢数量
	maxNNRet.WinLose = proto.Int64(winamount)

	//创建游戏结果
	nnGameRet := new(netproto.TBNN_GameRet)
	nnGameRet.UserRet = nnUserRets

	return nnGameRet
}
