package hhdz

import (
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

var _HungUserID = int32(1)
var _HeiUserID = int32(2)

type hhdzPokerHelper struct {
	poker     poker.Poker
	hungStore *hhdzPokerStore
	heiStore  *hhdzPokerStore
}

func newPokerHelper(poker poker.Poker) *hhdzPokerHelper {
	s := new(hhdzPokerHelper)
	s.poker = poker
	return s
}

func (s *hhdzPokerHelper) initStorePokerType(hungtype int32, heitype int32) {
	s.hungStore = s.getCardGroupByType(_HungUserID, hungtype)
	s.heiStore = s.getCardGroupByType(_HeiUserID, heitype)
}

//專門給測試程式產生牌型用的
func (s *hhdzPokerHelper) TestInitStorePoker(hungCG *poker.CardGroup, heiCG *poker.CardGroup) {
	s.hungStore = new(hhdzPokerStore)
	s.hungStore.store = hungCG
	s.heiStore = new(hhdzPokerStore)
	s.heiStore.store = heiCG
}

func (s *hhdzPokerHelper) getPrettyFirstCards() {
	s.hungStore.getPrettyFirstCards()
	s.heiStore.getPrettyFirstCards()
}

func (s *hhdzPokerHelper) getRandomFirstCards() {
	s.hungStore.getRandomFirstCards()
	s.heiStore.getRandomFirstCards()
}

//牌型交換, 讓winuserid 能夠贏得本局
func (s *hhdzPokerHelper) getPokerCardByWinUserIDFromStore(winuserid int32) {
	var changeStore = func() {
		tmpStore := s.hungStore.store
		s.hungStore.store = s.heiStore.store
		s.heiStore.store = tmpStore
	}
	if winuserid == 0 {
		// 不用換牌
	} else if winuserid == _HungUserID && comPareCards(s.hungStore.store, s.heiStore.store) != 1 {
		changeStore()
	} else if winuserid == _HeiUserID && comPareCards(s.heiStore.store, s.hungStore.store) != 1 {
		changeStore()
	} else if winuserid != _HungUserID && winuserid != _HeiUserID {
		mlog.Error("贏家指定錯誤")
	}
	s.hungStore.cards = s.hungStore.store.Cards
	s.heiStore.cards = s.heiStore.store.Cards
	s.hungStore.store.Cards = []*poker.Card{}
	s.heiStore.store.Cards = []*poker.Card{}
}

// 牌型已在發牌時按水池及機率表決定且前兩張不會是對子、這邊只需要逆轉第三張牌。
// 如果不是預期的輸贏結果，則將兩家牌的第三張牌移除掉後重新找牌
// winuserid 0 按原牌型發牌
// Hung 1, Hei 2
func (s *hhdzPokerHelper) getPokerCardByWinUserID(winuserid int32) {
	//1. 取出最後一張牌
	s.hungStore.getSecondCards()
	s.heiStore.getSecondCards()
	if winuserid == 0 { //未指定贏家就不調整牌型
		return
	}
	//2. 轉換輸贏角色
	var winstore *hhdzPokerStore
	var losestore *hhdzPokerStore
	loseuserid := int32(0)
	if winuserid == _HungUserID {
		loseuserid = _HeiUserID
		winstore = s.hungStore
		losestore = s.heiStore

	} else {
		loseuserid = _HungUserID
		winstore = s.heiStore
		losestore = s.hungStore
	}
	wincg := poker.NewCardGroup(winstore.cards)
	losecg := poker.NewCardGroup(losestore.cards)
	//mlog.System("1----------%s %s", wincg, losecg)

	//mlog.System("2----------%s %s", wincg, losecg)
	//3. 不符合預期進行贏家組牌程序
	if comPareCards(wincg, losecg) != 1 {
		// restore third card
		c1 := wincg.Cards[2]
		wincg.DeleteCard(c1)
		s.poker.AddCardGroup(winuserid, poker.NewCardGroup([]*poker.Card{c1}))
		c2 := losecg.Cards[2]
		losecg.DeleteCard(c2)
		s.poker.AddCardGroup(loseuserid, poker.NewCardGroup([]*poker.Card{c2}))
		//獲取能夠組單牌的牌、再交叉比較
		windancards := s.takeZuDanPaiDePai(winuserid, wincg)
		losedancards := s.takeZuDanPaiDePai(loseuserid, losecg)
		hasResult := false
		for _, windancard := range windancards {
			wincg.AddCard(windancard)
			for _, losedancard := range losedancards {
				losecg.AddCard(losedancard)
				if comPareCards(wincg, losecg) == 1 {
					//mlog.System("3----------%s %s", wincg, losecg)
					hasResult = true
					break
				} else {
					losecg.DeleteCard(losedancard)
				}
			}
			if hasResult {
				break
			} else {
				wincg.DeleteCard(windancard)
			}
		}
		if hasResult {
			winstore.cards = wincg.Cards
			losestore.cards = losecg.Cards
			//mlog.System("%s %s", wincg, losecg)
		} else {
			losecg.AddCard(losedancards[0]) //輸的任取一個可以組成單牌的牌
			losestore.cards = losecg.Cards
			duizipais := s.takeZuDuiziDePai(_HungUserID, wincg) //贏家要組成對子
			rnd := util.RndInt(0, len(duizipais))
			wincg.AddCard(duizipais[rnd])
			winstore.cards = wincg.Cards
			//mlog.System("%s %s", winstore.cards, losestore.cards)
		}
	} else {
		//mlog.System("%s %s", winstore.cards, losestore.cards)
	}
}

func (s *hhdzPokerHelper) getZuheTypes() (zuhehung, zuhehei, zuhetype int) {
	zuhehung = s.hungStore.getZuheType()
	zuhehei = s.heiStore.getZuheType()
	zuhetype = zuhehung
	if zuhehei < zuhetype {
		zuhetype = zuhehei
	}
	return
}

func (s *hhdzPokerHelper) getJiangType() (jiangtype int) {
	hungcg := poker.NewCardGroup(s.hungStore.cards)
	heicg := poker.NewCardGroup(s.heiStore.cards)
	if comPareCards(hungcg, heicg) == 1 {
		return HUNG
	} else {
		return HEI
	}
}

func (s *hhdzPokerHelper) getCardGroupByType(userid int32, pokertype int32) *hhdzPokerStore {
	cg := poker.NewCardGroup(nil)
	if pokertype == BAOZI {
		cg = s.getBowziCards(userid)
	} else if pokertype == SHUNJIN {
		cg = s.getShunJinCards(userid)
	} else if pokertype == JINHUA {
		cg = s.getJinHuaCards(userid)
	} else if pokertype == SHUNZI {
		cg = s.getShunZiCards(userid)
	} else if pokertype == BIGDUIZI {
		cg = s.getBigDuiziCards(userid)
	} else if pokertype == SMALLDUIZI {
		cg = s.getSmallDuiziCards(userid)
	} else if pokertype == DANPAI {
		cg = s.getDanPaiCards(userid)
	} else {
		cg = s.getNatureCards(userid)
	}
	sort.Slice(cg.Cards, func(i, j int) bool { return cg.Cards[i].CardNumber < cg.Cards[j].CardNumber })
	store := new(hhdzPokerStore)
	store.store = cg
	return store
}

func (s *hhdzPokerHelper) getBowziCards(userid int32) *poker.CardGroup {
	for i := 0; i < 1000; i++ { //找一千次
		s.poker.Shuffle()
		cg := s.poker.Deal(userid, 1)
		if ok, card := s.poker.TakeCardNumber(userid, cg.Cards[0].CardNumber); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		if ok, card := s.poker.TakeCardNumber(userid, cg.Cards[0].CardNumber); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		return cg
	}
	s.poker.Shuffle()
	cg := s.poker.Deal(userid, 3)
	return cg
}

func (s *hhdzPokerHelper) getShunJinCards(userid int32) *poker.CardGroup {
	for i := 0; i < 1000; i++ { //找一千次
		s.poker.Shuffle()
		cg := s.poker.Deal(userid, 1)
		firstCard := cg.Cards[0]
		startNum := firstCard.CardNumber
		if firstCard.CardNumber == 13 {
			continue
		} else if firstCard.CardNumber == 14 { //特殊取牌123
			startNum = 1
		}
		// 開始遞增取牌
		if ok, card := s.poker.TakeCardNumberAndColor(userid, startNum+1, firstCard.Color); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		if ok, card := s.poker.TakeCardNumberAndColor(userid, startNum+2, firstCard.Color); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		return cg
	}
	s.poker.Shuffle()
	cg := s.poker.Deal(userid, 3)
	return cg
}

func (s *hhdzPokerHelper) getJinHuaCards(userid int32) *poker.CardGroup {
	for i := 0; i < 1000; i++ { //找一千次
		s.poker.Shuffle()
		cg := s.poker.Deal(userid, 1)
		if ok, card := s.poker.TakeCardColor(userid, cg.Cards[0].Color); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		if ok, card := s.poker.TakeCardColor(userid, cg.Cards[0].Color); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		if isSunZi(cg) { //金花不可以是順子
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		return cg
	}
	s.poker.Shuffle()
	cg := s.poker.Deal(userid, 3)
	return cg
}

func (s *hhdzPokerHelper) getShunZiCards(userid int32) *poker.CardGroup {
	for i := 0; i < 1000; i++ { //找一千次
		s.poker.Shuffle()
		cg := s.poker.Deal(userid, 1)
		firstCard := cg.Cards[0]
		startNum := firstCard.CardNumber
		if firstCard.CardNumber == 13 {
			continue
		} else if firstCard.CardNumber == 14 { //特殊取牌123
			startNum = 1
		}
		// 開始遞增取牌
		if ok, card := s.poker.TakeCardNumber(userid, startNum+1); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		if ok, card := s.poker.TakeCardNumber(userid, startNum+2); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		if isJinHua(cg) { //順子不可以是金花
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		return cg
	}
	s.poker.Shuffle()
	cg := s.poker.Deal(userid, 3)
	return cg
}

func (s *hhdzPokerHelper) getBigDuiziCards(userid int32) *poker.CardGroup {
	for i := 0; i < 1000; i++ { //找一千次
		s.poker.Shuffle()
		v := int32(util.RndInt(9, 15))
		cg := poker.NewCardGroup(nil)
		if ok, card := s.poker.TakeCardNumber(userid, v); ok {
			cg.AddCard(&card)
		} else {
			continue
		}
		if ok, card := s.poker.TakeCardNumber(userid, v); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		if ok, card := s.poker.TakeDiffCardNumber(userid, []int32{v}); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		return cg
	}
	s.poker.Shuffle()
	cg := s.poker.Deal(userid, 3)
	return cg
}

func (s *hhdzPokerHelper) getSmallDuiziCards(userid int32) *poker.CardGroup {
	for i := 0; i < 1000; i++ { //找一千次
		s.poker.Shuffle()
		v := int32(util.RndInt(2, 9))
		cg := poker.NewCardGroup(nil)
		if ok, card := s.poker.TakeCardNumber(userid, v); ok {
			cg.AddCard(&card)
		} else {
			continue
		}
		if ok, card := s.poker.TakeCardNumber(userid, v); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		if ok, card := s.poker.TakeDiffCardNumber(userid, []int32{v}); ok {
			cg.AddCard(&card)
		} else {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		return cg
	}
	s.poker.Shuffle()
	cg := s.poker.Deal(userid, 3)
	return cg
}

func (s *hhdzPokerHelper) getDanPaiCards(userid int32) *poker.CardGroup {
	for i := 0; i < 1000; i++ { //找一千次
		s.poker.Shuffle()
		cg := s.poker.Deal(userid, 3)
		sort.Slice(cg.Cards, func(i, j int) bool { return cg.Cards[i].CardNumber < cg.Cards[j].CardNumber })
		if GetType(cg) != DanPai {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		return cg
	}
	s.poker.Shuffle()
	cg := s.poker.Deal(userid, 3)
	return cg
}

func (s *hhdzPokerHelper) getNatureCards(userid int32) *poker.CardGroup {
	for i := 0; i < 1000; i++ { //找一千次
		s.poker.Shuffle()
		cg := s.poker.Deal(userid, 3)
		sort.Slice(cg.Cards, func(i, j int) bool { return cg.Cards[i].CardNumber < cg.Cards[j].CardNumber })
		if GetType(cg) == BaoZi {
			s.poker.AddCardGroup(userid, cg)
			continue
		}
		return cg
	}
	s.poker.Shuffle()
	cg := s.poker.Deal(userid, 3)
	return cg
}

//獲取能組成單牌的牌
func (s *hhdzPokerHelper) takeZuDanPaiDePai(userid int32, cg *poker.CardGroup) (rsCards []*poker.Card) {
	s.poker.Shuffle()
	tmpcg := s.poker.Deal(userid, 10)
	//mlog.System("%v", tmpcg)
	for _, card := range tmpcg.Cards {
		tmpcg := poker.NewCardGroup(nil)
		tmpcg.AddCardGroup(cg)
		tmpcg.AddCard(card)
		//mlog.System("========== %v", tmpcg)
		if GetType(tmpcg) != DanPai {
			s.poker.AddCard(userid, card)
		} else {
			rsCards = append(rsCards, card)
		}
	}
	if len(rsCards) == 0 {
		mlog.System("=========== %s %v", cg, tmpcg)
		tmpcg = s.poker.Deal(userid, 1)
		rsCards = append(rsCards, tmpcg.Cards[0])
	}
	return rsCards
}

//獲取能組成對子的牌
func (s *hhdzPokerHelper) takeZuDuiziDePai(userid int32, cg *poker.CardGroup) (rsCards []*poker.Card) {
	s.poker.Shuffle()
	if ok, card := s.poker.TakeCardNumber(userid, cg.Cards[0].CardNumber); ok {
		rsCards = append(rsCards, &card)
	}
	if ok, card := s.poker.TakeCardNumber(userid, cg.Cards[1].CardNumber); ok {
		rsCards = append(rsCards, &card)
	}
	return
}

type hhdzPokerStore struct {
	cards []*poker.Card    //實際要出現的牌
	store *poker.CardGroup //預存牌型
}

//獲取好看的兩張牌來顯示,
func (s *hhdzPokerStore) getPrettyFirstCards() {
	// 像順金的兩張牌
	if ok, cards := s.likeShunJinCardsFromStore(); ok {
		s.cards = cards
	} else if ok, cards := s.likeJinHuaCardsFromStore(); ok {
		s.cards = cards
	} else if ok, cards := s.likeShunZiCardsFromStore(); ok {
		s.cards = cards
	} else {
		s.cards = s.randomFirstCards()
	}
}

//隨機取兩張牌來顯示
func (s *hhdzPokerStore) getRandomFirstCards() {
	s.cards = s.randomFirstCards()
}

//取最後一張牌
func (s *hhdzPokerStore) getSecondCards() {
	card := s.store.Cards[0]
	s.store.DeleteCard(card)
	s.cards = append(s.cards, card)
}

//獲取牌型
func (s *hhdzPokerStore) getZuheType() int {
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
	case DanPai:
		zuhetype = DANPAI
	default:
		zuhetype = DANPAI
	}
	return zuhetype
}

func (s *hhdzPokerStore) randomFirstCards() []*poker.Card {
	cards := make([]*poker.Card, 0)
	if isDuiZi(s.store) {
		rnd := util.RndInt(0, len(s.store.Cards)) //0,1 選一張
		card1 := s.store.Cards[rnd]
		s.store.DeleteCard(card1)
		card2idx := 0
		for i := 0; i < len(s.store.Cards); i++ {
			card2idx = i
			if card1.CardNumber != s.store.Cards[card2idx].CardNumber {
				break
			}
		}
		card2 := s.store.Cards[card2idx]
		s.store.DeleteCard(card2)
		cards = append(cards, card1, card2)
	} else {
		for i := 0; i < 2; i++ {
			l := len(s.store.Cards)
			rnd := util.RndInt(0, l)
			cards = append(cards, s.store.Cards[rnd])
			s.store.DeleteCard(s.store.Cards[rnd])
		}
	}
	return cards
}

//取出像順金的兩張牌
func (s *hhdzPokerStore) likeShunJinCardsFromStore() (bool, []*poker.Card) {
	if isSunZi(s.store) && isJinHua(s.store) { //如果是順金牌型、隨意取兩張
		return true, s.randomFirstCards()
	}
	size := len(s.store.Cards)
	for i := 0; i < size; i++ {
		card1 := s.store.Cards[i]
		for j := i + 1; j < size; j++ {
			card2 := s.store.Cards[j]
			minNumber := card1.CardNumber - 2
			maxNumber := card1.CardNumber + 2
			if card1.CardNumber != card2.CardNumber && card1.Color == card2.Color &&
				card2.CardNumber >= minNumber && card2.CardNumber <= maxNumber {
				s.store.DeleteCard(card1)
				s.store.DeleteCard(card2)
				return true, []*poker.Card{card1, card2}
			}
		}
	}
	return false, []*poker.Card{}
}

//取出像順子的兩張牌
func (s *hhdzPokerStore) likeShunZiCardsFromStore() (bool, []*poker.Card) {
	if isSunZi(s.store) {
		return true, s.randomFirstCards()
	}
	size := len(s.store.Cards)
	for i := 0; i < size; i++ {
		card1 := s.store.Cards[i]
		for j := i + 1; j < size; j++ {
			card2 := s.store.Cards[j]
			minNumber := card1.CardNumber - 2
			maxNumber := card1.CardNumber + 2
			if card1.CardNumber != card2.CardNumber &&
				card2.CardNumber >= minNumber && card2.CardNumber <= maxNumber {
				s.store.DeleteCard(card1)
				s.store.DeleteCard(card2)
				return true, []*poker.Card{card1, card2}
			}
		}
	}
	return false, []*poker.Card{}
}

//取出像金花的兩張牌
func (s *hhdzPokerStore) likeJinHuaCardsFromStore() (bool, []*poker.Card) {
	if isJinHua(s.store) {
		return true, s.randomFirstCards()
	}
	size := len(s.store.Cards)
	for i := 0; i < size; i++ {
		card1 := s.store.Cards[i]
		for j := i + 1; j < size; j++ {
			card2 := s.store.Cards[j]
			if card1.Color == card2.Color {
				s.store.DeleteCard(card1)
				s.store.DeleteCard(card2)
				return true, []*poker.Card{card1, card2}
			}
		}
	}
	return false, []*poker.Card{}
}
