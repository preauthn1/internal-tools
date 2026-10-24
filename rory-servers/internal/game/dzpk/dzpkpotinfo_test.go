package dzpk

import ()

//結算
//func TestSettlement(t *testing.T){
//	dzpk := new(DzpkLogic)
//
//	//公牌
//	dzpk.publicCards = new(poker.CardGroup)
//	pokerCardGroup := new(poker.CardGroup)
//	pokerCard1 := new(poker.PokerCard)
//	pokerCard1.CardNumber = 3
//	pokerCard1.Color = 0
//	pokerCardGroup.Cards = append(pokerCardGroup.Cards, pokerCard1)
//	pokerCard2 := new(poker.PokerCard)
//	pokerCard2.CardNumber = 14
//	pokerCard2.Color = 2
//	pokerCardGroup.Cards = append(pokerCardGroup.Cards, pokerCard2)
//	pokerCard3 := new(poker.PokerCard)
//	pokerCard3.CardNumber = 14
//	pokerCard3.Color = 1
//	pokerCardGroup.Cards = append(pokerCardGroup.Cards, pokerCard3)
//	pokerCard4 := new(poker.PokerCard)
//	pokerCard4.CardNumber = 3
//	pokerCard4.Color = 3
//	pokerCardGroup.Cards = append(pokerCardGroup.Cards, pokerCard4)
//	pokerCard5 := new(poker.PokerCard)
//	pokerCard5.CardNumber = 10
//	pokerCard5.Color = 2
//	pokerCardGroup.Cards = append(pokerCardGroup.Cards, pokerCard5)
//
//	//1玩家手牌
//	playerCardGroup := new(poker.CardGroup)
//	player1Card1 := new(poker.PokerCard)
//	player1Card1.CardNumber = 10
//	player1Card1.Color = 0
//	playerCardGroup.Cards = append(playerCardGroup.Cards, player1Card1)
//	player1Card2 := new(poker.PokerCard)
//	player1Card2.CardNumber = 6
//	player1Card2.Color = 0
//	playerCardGroup.Cards = append(playerCardGroup.Cards, player1Card2)
//
//
//	gameUserData := new(DzpkUserData)
//	gameUserData.userCards = playerCardGroup
//	gameUser := new(gamelogic.GameUser)
//	gameUserData.User = gameUser
//	gameUserData.User.SeatNo = 0
//	dzpk.gameUserData = append(dzpk.gameUserData, gameUserData)
//
//	//2玩家手牌
//	playerCardGroup2 := new(poker.CardGroup)
//	player2Card1 := new(poker.PokerCard)
//	player2Card1.CardNumber = 8
//	player2Card1.Color = 0
//	playerCardGroup2.Cards = append(playerCardGroup2.Cards, player2Card1)
//	player2Card2 := new(poker.PokerCard)
//	player2Card2.CardNumber = 7
//	player2Card2.Color = 0
//	playerCardGroup2.Cards = append(playerCardGroup2.Cards, player2Card2)
//
//
//	gameUserData2 := new(DzpkUserData)
//	gameUserData2.userCards = playerCardGroup2
//	gameUser2 := new(gamelogic.GameUser)
//	gameUserData2.User = gameUser2
//	gameUserData2.User.SeatNo = 1
//	dzpk.gameUserData = append(dzpk.gameUserData, gameUserData2)
//	dzpk.gameUserCount = 2
//
//
//	dzpk.potArray = new(potPool)
//	stJettonInfo1 := new(stJettonInfo)
//
//	stJettonInfo1.nJetton = int64(10)
//	stJettonInfo1.partin = append(stJettonInfo1.partin, int32(0))
//
//	dzpk.potArray.data = append(dzpk.potArray.data, stJettonInfo1)
//
//
//	stJettonInfo1.partin = append(stJettonInfo1.partin, int32(1))
//	//	pots.data[index].partin = append(pots.data[index].partin, user.seatNo)
//	dzpk.potArray.data = append(dzpk.potArray.data, stJettonInfo1)
//
//	dzpk.savePublic = pokerCardGroup
//	dzpk.publicCards.AddCardGroup(pokerCardGroup)
//
//
//	t.Logf("publicCards => %v", dzpk.publicCards)
//	settlement(dzpk)
//}
