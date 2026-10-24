package zjh

import (
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

type zjhPokerStore struct {
	cards []*poker.Card    //實際要出現的牌
	store *poker.CardGroup //預存牌型
}
type zjhPokerHelper struct {
	poker      poker.Poker
	user1Store *zjhPokerStore
	user2Store *zjhPokerStore
}

const (
	USER1 = 1 //User1
	USER2 = 2 //User2
)

//獲取牌型
func (s *zjhPokerStore) getZuheType() int {
	zuhetype := 0
	cg := poker.NewCardGroup(s.cards)
	cardtype := GetType(cg)

	switch cardtype {
	case BaoZi:
		zuhetype = BAOZI
		break
	case SunJin:
		zuhetype = SHUNJIN
		break
	case JiaHua:
		zuhetype = JINHUA
		break
	case SunZi:
		zuhetype = SHUNZI
		break
	case DuiZi:
		number := int32(0)
		if cg.Cards[0].CardNumber == cg.Cards[1].CardNumber {
			number = cg.Cards[0].CardNumber
		} else if cg.Cards[1].CardNumber == cg.Cards[2].CardNumber {
			number = cg.Cards[1].CardNumber
		} else if cg.Cards[0].CardNumber == cg.Cards[2].CardNumber {
			number = cg.Cards[0].CardNumber
		}
		if number >= 9 {
			zuhetype = BIGDUIZI
		} else {
			zuhetype = SMALLDUIZI
		}
		break
	case SanPai:
		zuhetype = DANPAI
	default:
		zuhetype = DANPAI
	}
	return zuhetype
}

func newPokerHelper(poker poker.Poker) *zjhPokerHelper {
	s := new(zjhPokerHelper)
	s.poker = poker
	return s
}
func (s *zjhPokerHelper) TestInitStorePoker(user1CG *poker.CardGroup, user2CG *poker.CardGroup) {
	s.user1Store = new(zjhPokerStore)
	s.user1Store.store = user1CG
	s.user2Store = new(zjhPokerStore)
	s.user2Store.store = user2CG
	s.user1Store.cards = s.user1Store.store.Cards
	s.user2Store.cards = s.user2Store.store.Cards
	s.user1Store.store.Cards = []*poker.Card{}
	s.user2Store.store.Cards = []*poker.Card{}
}

func (s *zjhPokerHelper) getJiangType() (jiangtype int) {
	user1cg := poker.NewCardGroup(s.user1Store.cards)
	user2cg := poker.NewCardGroup(s.user2Store.cards)

	if comPareCards(user1cg, user2cg) == 1 {
		return USER1
	} else {
		return USER2
	}
}
