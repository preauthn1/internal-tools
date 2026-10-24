package brnn

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

//廣播第一手三張牌，並开始下注，下注時間結束才發第二手兩張牌並接著開獎
func (brnn *BRNNLogic) startFirstCard() {
	brnn.fapeiFirstCardGroupNew()
	//廣播第一手三張牌
	brnn.boardcastFirstCardGroup()

}

func (brnn *BRNNLogic) fapeiFirstCardGroupNew() {
	var tenLimitCount int
	var kingLimitCount int
	if brnn.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		brnn.poker.Shuffle()
		tenLimitCount = 2  //十點以上的牌總共只能發出三張並且第一個手牌最多兩張十點以上
		kingLimitCount = 1 //大小王＋皮子只能三選一發出
	}

	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
		cg := brnn.poker.Deal(0, 3)                 //固定userid 0 只為了取得隨機三張牌
		countKingToPizi := cg.HasKingToPizi()       //計算大小王皮子的數量
		countTenToThirteen := cg.HasTenToThirteen() //大於10的牌數

		if kingLimitCount > 0 && countKingToPizi > kingLimitCount { //不要有王和皮子同時在一起
			i--
			brnn.poker.AddCardGroup(0, cg)
			brnn.poker.Shuffle()
			continue
		}

		if tenLimitCount > 0 && countTenToThirteen > tenLimitCount { //第一個手牌最多兩張十點以上
			i--
			brnn.poker.AddCardGroup(0, cg)
			brnn.poker.Shuffle()
			continue
		}

		urPokerCard := new(netproto.BRNN_BrodcastUserPokerCard)
		urPokerCard.SitNo = proto.Int32(i)

		for _, cd := range cg.Cards {
			card := new(netproto.PokerCard)
			card.Color = proto.Int32(cd.Color)
			card.CardNumber = proto.Int32(cd.CardNumber)
			urPokerCard.CardGroup = append(urPokerCard.CardGroup, card)
		}

		brnn.allKindPoker[int(i)] = cg
		brnn.GameLog.AddGameStep("index: %v 第一手手牌:%v ", i, cg.Cards)
	}

}

//產生第一手三張牌, 棄用
//func (brnn *BRNNLogic) fapeiFirstCardGroup() {
//	brnn.poker.Shuffle()
//	tenLimitCount := int(3) //十點以上的牌總共只能發出三張並且停一個手牌最多兩張十點以上
//	kingCount := int(1)     //大小王＋皮子只能三選一發出
//
//	for i := int32(0); i <= int32(len(brnn.betvalue.alltype)); i++ {
//		cg := brnn.poker.Deal(0, 3)                 //固定userid 0 只為了取得隨機三張牌
//		countKingToPizi := cg.HasKingToPizi()       //計算大小王皮子的數量
//		countTenToThirteen := cg.HasTenToThirteen() //大於10的牌數
//
//		if countKingToPizi > 0 && countTenToThirteen > 0 { //不要有王和皮子同時在一起
//			i--
//			brnn.poker.AddCardGroup(0, cg)
//			brnn.poker.Shuffle()
//			continue
//		}
//
//		if countTenToThirteen > 0 {
//			if tenLimitCount <= 0 || countTenToThirteen > 2 {
//				i--
//				brnn.poker.AddCardGroup(0, cg)
//				brnn.poker.Shuffle()
//				continue
//			}
//			tenLimitCount -= countTenToThirteen
//		}
//
//		if countKingToPizi > 0 {
//			if kingCount <= 0 || countKingToPizi > 1 {
//				i--
//				brnn.poker.AddCardGroup(0, cg)
//				brnn.poker.Shuffle()
//				continue
//			}
//			kingCount -= countKingToPizi
//		}
//
//		urPokerCard := new(netproto.BRNN_BrodcastUserPokerCard)
//		urPokerCard.SitNo = proto.Int32(i)
//
//		for _, cd := range cg.Cards {
//			card := new(netproto.PokerCard)
//			card.Color = proto.Int32(cd.Color)
//			card.CardNumber = proto.Int32(cd.CardNumber)
//			urPokerCard.CardGroup = append(urPokerCard.CardGroup, card)
//		}
//
//		brnn.allKindPoker[int(i)] = cg
//		brnn.GameLog.AddGameStep("index: %v 第一手手牌:%v ", i, cg.Cards)
//	}
//
//}

//廣播第一手牌
func (brnn *BRNNLogic) boardcastFirstCardGroup() {
	brnn.BoardcastMessage(brnn.getFirstCardGroup())
}
