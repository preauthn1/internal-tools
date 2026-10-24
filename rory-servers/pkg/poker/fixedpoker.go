package poker

import (
	"bytes"
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
const (
	RestCardsUserID int32 = 0
)

// FixedPoker 根據用戶事先確定的牌型
type FixedPoker struct {
	userCards map[int32]*CardGroup
}

// NewFixedPoker 創建固定牌型對象
func NewFixedPoker(urs []*gamelogic.GameUser, randPoker *RandomPoker, userCardLen int) *FixedPoker {
	fp := new(FixedPoker)
	fp.userCards = make(map[int32]*CardGroup)

	// 洗牌
	randPoker.ShuffleMore(3)
	// 分牌
	fp.alloc(urs, randPoker, userCardLen)

	return fp
}

// NewFixedNoShufflePoker ?
func NewFixedNoShufflePoker(urs []*gamelogic.GameUser, noShuffle *NoShufflePoker, userCardLen int, allocTake map[int]int) *FixedPoker {
	fp := new(FixedPoker)
	fp.userCards = make(map[int32]*CardGroup)

	randPoker := new(RandomPoker)
	randPoker.Cards = append(randPoker.Cards, noShuffle.Cards...)

	// 分牌
	fp.allocTakeGroup(urs, randPoker, userCardLen, allocTake)

	return fp
}

// NewFixedFilePoker 創建基於發牌器的固定牌型對象
func NewFixedFilePoker(filePoker *FilePoker) *FixedPoker {
	fp := new(FixedPoker)
	fp.userCards = filePoker.userCards

	return fp
}

// NewFixedFileFillPoker 創建基於發牌器的固定牌型填滿角色對象(依照指定填滿，判斷順序依 座位 > UserID > 隨機填滿)
func NewFixedFileFillPoker(filePoker *FilePoker, urs []*gamelogic.GameUser) *FixedPoker {
	fp := new(FixedPoker)
	var userCards map[int32]*CardGroup
	userCards = make(map[int32]*CardGroup, 0)

	var cardList []*CardGroup
	for _, v := range filePoker.userCards {
		cardList = append(cardList, v)
	}
	for _, ur := range urs {
		if val, ok := filePoker.userCards[ur.SeatNo]; ok {
			userCards[ur.UserID] = val
		} else if val, ok := filePoker.userCards[ur.UserID]; ok {
			userCards[ur.UserID] = val
		} else {
			// var popFrontCard *CardGroup
			// if len(cardList) > 0 {
			// 	popFrontCard, cardList = cardList[0], cardList[1:] // pop front card
			// 	userCards[ur.UserID] = popFrontCard
			// } else {
			var randCards = NewFixedPoker(urs, NewRandomPoker(false, false, 13), 13).userCards
			userCards[ur.UserID] = randCards[ur.UserID]
			// }
		}
	}
	fp.userCards = userCards

	return fp
}

// NewZJHBigFixedPoker 第一種發牌規則
func NewZJHBigFixedPoker(userCardType map[int32]int, randPoker *RandomPoker) *FixedPoker {
	fp := new(FixedPoker)
	fp.userCards = make(map[int32]*CardGroup)
	// 洗牌
	randPoker.ShuffleMore(3)
	for userid, cardtype := range userCardType {
		cards := new(CardGroup)
		error := false
		if cardtype == 1 {
			cards, error = getBaoZi(randPoker, userid)
			if error {
				cardtype = 2
			}
		}
		if cardtype == 2 {
			cards, error = getShunJin(randPoker, userid)
			if error {
				cardtype = 3
			}
		}
		if cardtype == 3 {
			cards, error = getJinHua(randPoker, userid)
			if error {
				cardtype = 4
			}
		}
		if cardtype == 11 {
			cards, error = getBigJinHua(randPoker, userid)
			if error {
				cardtype = 4
			}
		}
		if cardtype == 12 {
			cards, error = getSmallJinHua(randPoker, userid)
			if error {
				cardtype = 4
			}
		}
		if cardtype == 4 {
			cards, error = getShunZi(randPoker, userid)
			if error {
				cardtype = 5
			}
		}
		if cardtype == 13 {
			cards, error = getDuiZi(randPoker, userid)
			if error {
				cardtype = 7
			}
		}
		if cardtype == 5 {
			cards, error = getBigDuiZi(randPoker, userid)
			if error {
				cardtype = 6
			}
		}
		if cardtype == 6 {
			cards, error = getSmallDuiZi(randPoker, userid)
			if error {
				cardtype = 7
			}
		}
		if cardtype == 7 || cards.Len() != 3 {
			cards = randPoker.Deal(userid, 3)
		}

		mlog.Debug("用戶 %d 的牌為: %v", userid, cards)
		fp.userCards[userid] = cards
	}
	return fp
}

// FixedFileAlloc 發牌器的固定牌型对象按順序分牌
func (fp *FixedPoker) FixedFileAlloc(urs []*gamelogic.GameUser) {
	idx := int32(1)
	for _, user := range urs {
		if _, ok := fp.userCards[user.UserID]; !ok {
			fp.userCards[user.UserID] = fp.userCards[idx]
			idx++
		}
	}
	mlog.Debug("=========== %v", fp.userCards)
}

// 分牌
func (fp *FixedPoker) alloc(urs []*gamelogic.GameUser, randPoker *RandomPoker, userCardLen int) {
	allocedLen := 0

	for _, ur := range urs {
		uid := ur.UserID

		fp.userCards[uid] = randPoker.Deal(uid, userCardLen)
		allocedLen += userCardLen

		mlog.Debug("userid: %d, card: %s", uid, fp.userCards[uid].String())
	}

	fp.userCards[RestCardsUserID] = NewCardGroup(randPoker.Cards[0:])

	mlog.Debug("牌型:%s", fp.String())
}

//指次每次分牌的張數
func (fp *FixedPoker) allocTakeGroup(urs []*gamelogic.GameUser, randPoker *RandomPoker, userCardLen int, group map[int]int) {
	allocedLen := 0
	takePoker := userCardLen
	rand := util.RndInt(0, len(group))
	take := group[rand] //設定幾張一組

	for takePoker > 0 {

		for _, ur := range urs {
			uid := ur.UserID

			if _, ok := fp.userCards[uid]; ok {
				fp.userCards[uid].Cards = append(fp.userCards[uid].Cards, randPoker.Deal(uid, take).Cards...)
			} else {
				fp.userCards[uid] = randPoker.Deal(uid, take)
			}

			allocedLen += take

		}

		takePoker -= take
		if takePoker < take {
			take = takePoker
			randPoker.Shuffle()
		}

	}

	fp.userCards[RestCardsUserID] = NewCardGroup(randPoker.Cards[0:])

	mlog.Debug("牌型:%s", fp.String())
}

// Shuffle 洗牌
func (fp *FixedPoker) Shuffle() {

}

// UserCardsLen ?
func (fp *FixedPoker) UserCardsLen(uid int32) int {
	return fp.userCards[uid].Len()
}

// UserCardLen ?
func (fp *FixedPoker) UserCardLen() int {
	return len(fp.userCards)
}

// ShuffleMore 洗多次牌
func (fp *FixedPoker) ShuffleMore(times int) {

}

// TakeCardNumber ?
func (fp *FixedPoker) TakeCardNumber(userid int32, cardNumber int32) (bool, Card) {
	for k, v := range fp.userCards[userid].Cards {
		if v.CardNumber == cardNumber {
			fp.userCards[userid].Cards = RemoveIndex(fp.userCards[userid].Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// TakeDiffCardNumber 指定取非存在的號碼
func (fp *FixedPoker) TakeDiffCardNumber(userid int32, diffNums []int32) (bool, Card) {
	for k, v := range fp.userCards[userid].Cards {
		right := true
		for _, diffnum := range diffNums {
			if v.CardNumber == diffnum {
				right = false
			}
		}
		if right {
			fp.userCards[userid].Cards = RemoveIndex(fp.userCards[userid].Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// TakeCardColor ?
func (fp *FixedPoker) TakeCardColor(userid int32, color int32) (bool, Card) {
	for k, v := range fp.userCards[userid].Cards {
		if v.Color == color {
			fp.userCards[userid].Cards = RemoveIndex(fp.userCards[userid].Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// TakeCardNumberAndColor ?
func (fp *FixedPoker) TakeCardNumberAndColor(userid int32, cardNumber int32, color int32) (bool, Card) {
	for k, v := range fp.userCards[userid].Cards {
		if v.CardNumber == cardNumber && v.Color == color {
			fp.userCards[userid].Cards = RemoveIndex(fp.userCards[userid].Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// Deal 發牌
func (fp *FixedPoker) Deal(userid int32, cn int) *CardGroup {
	cuid := RestCardsUserID

	if _, ok := fp.userCards[userid]; ok {
		cuid = userid
	}
	if len(fp.userCards[cuid].Cards) < cn {
		cn = len(fp.userCards[cuid].Cards)
	}

	cards := fp.userCards[cuid].Cards[:cn]
	fp.userCards[userid].Cards = fp.userCards[cuid].Cards[cn:]

	return NewCardGroup(cards)
}

// AddCardGroup ?
func (fp *FixedPoker) AddCardGroup(userid int32, cg *CardGroup) {
	if len(cg.Cards) <= int(0) {
		return
	}

	fp.userCards[userid].Cards = append(fp.userCards[userid].Cards, cg.Cards...)
}

// AddCard ?
func (fp *FixedPoker) AddCard(userid int32, card *Card) {
	if card == nil {
		return
	}

	fp.userCards[userid].Cards = append(fp.userCards[userid].Cards, card)
}

// ExChange ?
func (fp *FixedPoker) ExChange(userid int32, otheruserid int32) {
	if otheruserid == RestCardsUserID || userid == RestCardsUserID {
		return
	}

	fp.userCards[userid], fp.userCards[otheruserid] = fp.userCards[otheruserid], fp.userCards[userid]
}

// GetUserCards 獲取用戶牌型
func (fp *FixedPoker) GetUserCards(userid int32) *CardGroup {
	return fp.userCards[userid]
}

func (fp *FixedPoker) String() string {
	var buf bytes.Buffer
	for uid, cg := range fp.userCards {
		buf.WriteString(fmt.Sprintf("%d:%s\r\n", uid, cg.String()))
	}

	return buf.String()
}

func getBaoZi(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	takeCardGroup := new(CardGroup)
	fp.Shuffle()
	cg := fp.Deal(userID, 1)
	for i := 0; i < 2; i++ {
		ok, takeCard := fp.TakeCardNumber(userID, cg.Cards[0].CardNumber)
		if !ok {
			//找不到牌 往下找
			fp.AddCardGroup(userID, cg)
			fp.AddCardGroup(userID, takeCardGroup)
			return takeCardGroup, true
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	pokerCard5 := new(Card)
	pokerCard5.CardNumber = cg.Cards[0].CardNumber
	pokerCard5.Color = cg.Cards[0].Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard5)
	return takeCardGroup, false
}

func getShunJin(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	takeCardGroup := new(CardGroup)
	findCardNumber1 := int32(util.RndInt(2, 12))
	findCardNumber2 := findCardNumber1 + 1
	findCardNumber3 := findCardNumber2 + 1
	cardPoint := []int32{findCardNumber3, findCardNumber2, findCardNumber1}
	findColor := int32(util.RndInt(0, 4))
	for _, point := range cardPoint {

		if point == 1 {
			point = 14
		}
		ok, takeCard := fp.TakeCardNumberAndColor(userID, int32(point), findColor)
		if !ok {
			fp.AddCardGroup(userID, takeCardGroup)
			return takeCardGroup, true
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	return takeCardGroup, false
}

func getJinHua(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	takeCardGroup := new(CardGroup)
	fp.Shuffle()
	cg := fp.Deal(userID, 1)
	for i := 0; i < 2; i++ {
		ok, takeCard := fp.TakeCardColor(userID, cg.Cards[0].Color)
		if !ok {
			//補上再找其他牌的方法
			fp.AddCardGroup(userID, cg)
			fp.AddCardGroup(userID, takeCardGroup)
			return takeCardGroup, true
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	pokerCard5 := new(Card)
	pokerCard5.CardNumber = cg.Cards[0].CardNumber
	pokerCard5.Color = cg.Cards[0].Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard5)
	if isSunZi(takeCardGroup) { //若是順子則當作找不到
		return takeCardGroup, true
	}
	return takeCardGroup, false
}

func getBigJinHua(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	takeCardGroup := new(CardGroup)
	fp.Shuffle()
	rnd := []int{13, 14}
	rnd = util.IntArrayShuffle(rnd)
	var pokerCard *Card
	for _, r := range rnd {
		if ok, card := fp.TakeCardNumber(userID, int32(r)); ok {
			pokerCard = &card
			break
		}
	}
	cg := NewCardGroup([]*Card{pokerCard})
	for i := 0; i < 2; i++ {
		ok, takeCard := fp.TakeCardColor(userID, cg.Cards[0].Color)
		if !ok {
			//補上再找其他牌的方法
			fp.AddCardGroup(userID, cg)
			fp.AddCardGroup(userID, takeCardGroup)
			return takeCardGroup, true
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	pokerCard5 := new(Card)
	pokerCard5.CardNumber = cg.Cards[0].CardNumber
	pokerCard5.Color = cg.Cards[0].Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard5)
	if isSunZi(takeCardGroup) { //若是順子則當作找不到
		return takeCardGroup, true
	}
	return takeCardGroup, false
}

func getSmallJinHua(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	takeCardGroup := new(CardGroup)
	fp.Shuffle()
	rnd := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	rnd = util.IntArrayShuffle(rnd)
	var pokerCard *Card
	for _, r := range rnd {
		if ok, card := fp.TakeCardNumber(userID, int32(r)); ok {
			pokerCard = &card
			break
		}
	}
	cg := NewCardGroup([]*Card{pokerCard})
	for i := 0; i < 2; i++ {
		ok, takeCard := fp.TakeCardColorLessThen(userID, 13, cg.Cards[0].Color)
		if !ok {
			//補上再找其他牌的方法
			fp.AddCardGroup(userID, cg)
			fp.AddCardGroup(userID, takeCardGroup)
			return takeCardGroup, true
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	pokerCard5 := new(Card)
	pokerCard5.CardNumber = cg.Cards[0].CardNumber
	pokerCard5.Color = cg.Cards[0].Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard5)
	if isSunZi(takeCardGroup) { //若是順子則當作找不到
		return takeCardGroup, true
	}
	return takeCardGroup, false
}

func getShunZi(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	takeCardGroup := new(CardGroup)
	findCardNumber1 := int32(util.RndInt(2, 12))
	findCardNumber2 := findCardNumber1 + 1
	findCardNumber3 := findCardNumber2 + 1
	cardPoint := []int32{findCardNumber3, findCardNumber2, findCardNumber1}
	for _, point := range cardPoint {
		if point == 1 {
			point = 14
		}
		ok, takeCard := fp.TakeCardNumber(userID, int32(point))
		if !ok {
			//補上再找其他牌的方法
			fp.AddCardGroup(userID, takeCardGroup)
			return takeCardGroup, true
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	return takeCardGroup, false
}

func getBigDuiZi(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	takeCardGroup := new(CardGroup)
	fp.Shuffle()
	cg := fp.Deal(userID, 1) //先發一張 再挑兩張
	findCardNumber := int32(util.RndInt(9, 15))
	if cg.Cards[0].CardNumber == findCardNumber {
		findCardNumber = findCardNumber - 1
	}
	cardPoint := []int32{findCardNumber, findCardNumber}
	for _, point := range cardPoint {
		if point == 1 {
			point = 14
		}
		ok, takeCard := fp.TakeCardNumber(userID, int32(point))
		if !ok {
			//補上再找其他牌的方法
			fp.AddCardGroup(userID, takeCardGroup)
			fp.AddCardGroup(userID, cg)
			return takeCardGroup, true
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	pokerCard5 := new(Card)
	pokerCard5.CardNumber = cg.Cards[0].CardNumber
	pokerCard5.Color = cg.Cards[0].Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard5)
	return takeCardGroup, false
}

func getSmallDuiZi(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	takeCardGroup := new(CardGroup)
	fp.Shuffle()
	cg := fp.Deal(userID, 1) //先發一張 再挑兩張

	findCardNumber := int32(util.RndInt(2, 9))
	if cg.Cards[0].CardNumber == findCardNumber {
		findCardNumber = findCardNumber + 1
	}
	cardPoint := []int32{findCardNumber, findCardNumber}
	for _, point := range cardPoint {
		if point == 1 {
			point = 14
		}
		ok, takeCard := fp.TakeCardNumber(userID, int32(point))
		if !ok {
			//補上再找其他牌的方法
			fp.AddCardGroup(userID, takeCardGroup)
			fp.AddCardGroup(userID, cg)
			return takeCardGroup, true
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	pokerCard5 := new(Card)
	pokerCard5.CardNumber = cg.Cards[0].CardNumber
	pokerCard5.Color = cg.Cards[0].Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard5)
	return takeCardGroup, false
}

func getDuiZi(fp *RandomPoker, userID int32) (*CardGroup, bool) {
	fp.Shuffle()
	takeCardGroup := fp.Deal(userID, 2)                                         //先挑兩張
	if takeCardGroup.Cards[0].CardNumber == takeCardGroup.Cards[1].CardNumber { //己經成對, 則再發一張牌返回
		takeCardGroup.AddCardGroup(fp.Deal(userID, 1))
	} else {
		cardPoint := []int32{takeCardGroup.Cards[0].CardNumber, takeCardGroup.Cards[1].CardNumber}
		for _, point := range cardPoint {
			if point == 1 {
				point = 14
			}
			ok, takeCard := fp.TakeCardNumber(userID, int32(point))
			if ok {
				takeCardGroup.AddCard(&takeCard)
				break
			}
		}
	}

	if takeCardGroup.Len() != 3 { //把牌放回去
		fp.AddCardGroup(userID, takeCardGroup)
		return nil, true
	}
	return takeCardGroup, false
}

func isSunZi(Cards *CardGroup) bool {
	sortCard(Cards)
	if Cards.Cards[0].CardNumber+1 == Cards.Cards[1].CardNumber && Cards.Cards[1].CardNumber+1 == Cards.Cards[2].CardNumber {
		return true
	}

	//特殊123
	if Cards.Cards[0].CardNumber == 2 && Cards.Cards[1].CardNumber == 3 && Cards.Cards[2].CardNumber == 14 {
		return true
	}
	return false
}

//给牌排序
func sortCard(cards *CardGroup) {
	arrlen := len(cards.Cards)
	for i := 0; i < arrlen; i++ {
		for j := 0; j+1 < arrlen-i; j++ {
			card1 := cards.Cards[j].CardNumber*10 + cards.Cards[j].Color
			card2 := cards.Cards[j+1].CardNumber*10 + cards.Cards[j+1].Color
			if card1 > card2 {
				cards.Swap(j, j+1)
			}
		}
	}
}
