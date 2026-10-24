package srnn

import (
	proto "example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

//牛牛牌局结果
type nnGameRet struct {
	nnlgc *SRNNGameLogic
}

//创建牛牛游戏结果对象
func newNNGameRet(nnglc *SRNNGameLogic) *nnGameRet {
	nnret := new(nnGameRet)
	nnret.nnlgc = nnglc

	return nnret
}

//获取游戏结果
func (nnret *nnGameRet) getGameRet() *netproto.SRNN_GameRet {
	gret := new(netproto.SRNN_GameRet)

	//庄家牌型
	zjuid := nnret.nnlgc.zhuangInfo.userID

	//庄家牌型
	zjpoker := nnret.nnlgc.gamePoker.GetUserPoker(zjuid)

	//庄家输赢
	var zhuangwinlose int64 = 0

	zjcardret := getNiuNiuRet(zjpoker.Cards)
	gret.ZhuangRet = new(netproto.SRNN_UserRet)
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
		xjwinlose := int64(nnret.nnlgc.zhuangInfo.betTimes) * int64(xjbt.betTimes) * int64(betret) * int64(nnret.nnlgc.GetServerConfig().BaseScore)

		//庄家输赢金币数
		zhuangwinlose += -1 * xjwinlose

		xjuserRet := new(netproto.SRNN_UserRet)
		xjuserRet.WinLose = proto.Int64(xjwinlose)
		xjuserRet.User = nnret.nnlgc.GetUserByUserID(xjbt.userID).GetProtoBufGameUser()
		xjuserRet.NiuNiuRet = proto.Int32(int32(xjcardret.niuRet))
		xjuserRet.NiuNiuRetNumber = proto.Int32(int32(xjcardret.number))
		xjuserRet.NiuNiuBetTimes = proto.Int32(betret)

		xjuserRet.CardGroup = xjcardret.getNetCardGroup()

		gret.UserRet = append(gret.UserRet, xjuserRet)
	}

	gret.ZhuangRet.WinLose = proto.Int64(zhuangwinlose) //庄家输赢

	return gret
}
