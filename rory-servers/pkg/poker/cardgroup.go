package poker

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/util"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// CardGroup 牌組
type CardGroup struct {
	Cards []*Card
}

// NewCardGroup 創建牌組對象
func NewCardGroup(cards []*Card) *CardGroup {
	cg := new(CardGroup)
	if cards != nil {
		cg.Cards = append(cg.Cards, cards...)
	}
	return cg
}

// LoadCardGroupFromString ?
func LoadCardGroupFromString(pokerString string) *CardGroup {
	cg := new(CardGroup)
	tmpArr := strings.Split(pokerString, "||")
	// 解析牌型
	sarr := strings.Split(tmpArr[0], "|")
	for _, s := range sarr {
		if util.TrimAllSpace(s) == "" {
			continue
		}

		pk := NewPokerCardByStr(s, ",")
		if pk == nil {
			mlog.Error("解析牌型錯誤,字符串%s格式不正確.%s", s, pokerString)
			continue
		}
		cg.AddCard(pk)
	}
	return cg
}

// NewEmptyCardGroup 創建空牌組對象
func NewEmptyCardGroup(cn int) *CardGroup {
	cg := new(CardGroup)
	for i := 0; i < cn; i++ {
		cg.AddCard(NewEmptyPokerCard())
	}

	return cg
}

// SortMinFirst 排序小的牌在前
func (cg *CardGroup) SortMinFirst() {
	if len(cg.Cards) == 0 {
		return
	}
	sort.SliceStable(cg.Cards, func(i, j int) bool {
		var number1 = cg.Cards[i].CardNumber
		var number2 = cg.Cards[j].CardNumber
		var card1 = cg.Cards[i]
		var card2 = cg.Cards[j]

		if number1 < number2 {
			return true
		} else if number1 > number2 {
			return false
		} else {
			if card1.IsMeiHua() {
				return true
			} else if card1.IsFangKuai() && !card2.IsMeiHua() {
				return true
			} else if card1.IsHongTao() && !card2.IsMeiHua() && !card2.IsFangKuai() {
				return true
			}
			return false
		}
	})
}

// SortMaxFirst 排序大的牌在前
func (cg *CardGroup) SortMaxFirst() {
	if len(cg.Cards) == 0 {
		return
	}
	sort.SliceStable(cg.Cards, func(i, j int) bool {
		var number1 = cg.Cards[i].CardNumber
		var number2 = cg.Cards[j].CardNumber
		var card1 = cg.Cards[i]
		var card2 = cg.Cards[j]

		if number1 < number2 {
			return false
		} else if number1 > number2 {
			return true
		} else {
			if card1.IsMeiHua() {
				return false
			} else if card1.IsFangKuai() && !card2.IsMeiHua() {
				return false
			} else if card1.IsHongTao() && !card2.IsMeiHua() && !card2.IsFangKuai() {
				return false
			}
			return true
		}
	})
}

// ShuffleOnce 洗牌
func (cg *CardGroup) ShuffleOnce() {
	mcards := make(map[int]*Card)

	for idx, c := range cg.Cards {
		mcards[idx] = c
	}

	cg.Cards = make([]*Card, 0)

	for _, c := range mcards {
		cg.Cards = append(cg.Cards, c)
	}
}

// AddCard 添加牌
func (cg *CardGroup) AddCard(pk *Card) {
	cg.Cards = append(cg.Cards, pk)
}

// AddCardGroup 添加牌组
func (cg *CardGroup) AddCardGroup(ncg *CardGroup) {
	cg.Cards = append(cg.Cards, ncg.Cards...)
}

// IsEqual 是否同樣的牌組
func (cg *CardGroup) IsEqual(ncg *CardGroup) bool {
	if cg.Len() != ncg.Len() {
		return false
	}

	sort.Sort(cg)
	sort.Sort(ncg)

	for idx, c := range cg.Cards {
		if !c.IsEqual(ncg.Cards[idx]) {
			return false
		}
	}

	return true
}

func (cg *CardGroup) Len() int {
	return len(cg.Cards)
}

func (cg *CardGroup) Swap(i, j int) {
	cg.Cards[i], cg.Cards[j] = cg.Cards[j], cg.Cards[i]
}

func (cg *CardGroup) Less(i, j int) bool {
	cardi, _ := strconv.Atoi(fmt.Sprintf("%d%d", cg.Cards[i].CardNumber, cg.Cards[i].Color))
	cardj, _ := strconv.Atoi(fmt.Sprintf("%d%d", cg.Cards[j].CardNumber, cg.Cards[j].Color))

	return cardi > cardj
}

func (cg *CardGroup) String() string {
	if cg == nil {
		return ""
	}
	var buf bytes.Buffer
	for _, c := range cg.Cards {
		buf.WriteString(fmt.Sprintf("%s,", c.String()))
	}

	return buf.String()
}

// FixedString ?
func (cg *CardGroup) FixedString() string {
	if cg == nil {
		return ""
	}
	var buf bytes.Buffer
	for _, c := range cg.Cards {
		buf.WriteString(fmt.Sprintf("%d,%d|", c.Color, c.CardNumber))
	}

	return buf.String()
}

// IsContent 檢測給定牌組是否屬於玩家手牌 檢測出牌人 出的牌是否合法
func (cg *CardGroup) IsContent(cards []*Card) bool {
	l := len(cards)
	var isContent bool = false

	for i := 0; i < l; i++ {
		isContent = false

		for _, v := range cg.Cards {
			if !v.IsEqual(cards[i]) {
				continue
			}
			isContent = true
			break
		}
		if !isContent {
			return false
		}
	}
	return isContent
}

// HaveSameCards 檢測牌組內是否有相同的牌
func (cg *CardGroup) HaveSameCards(compareGroup *CardGroup) bool {

	for _, v1 := range cg.Cards {
		for _, v2 := range compareGroup.Cards {
			if v1.IsEqual(v2) {
				return true
			}
		}
	}
	return false
}

// DeleteGroup ?
func (cg *CardGroup) DeleteGroup(group *CardGroup) {
	for _, v := range group.Cards {
		cg.DeleteCard(v)
	}
}

// DeleteCard ?
func (cg *CardGroup) DeleteCard(card *Card) {
	tmp := make([]*Card, 0)
	for _, v := range cg.Cards {
		if !v.IsEqual(card) {
			tmp = append(tmp, v)
		}
	}
	cg.Cards = tmp
}

// DeleteCardForPokers ?
func (cg *CardGroup) DeleteCardForPokers(card *Card) {
	for i := 0; i < len(cg.Cards); i++ {
		if card.IsEqual(cg.Cards[i]) {
			cg.Cards = append(cg.Cards[:i], cg.Cards[i+1:]...)
			break
		}
	}

	mlog.Debug("card: %d", len(cg.Cards))
}

// GetMinCardNumber 獲取最小牌點
func (cg *CardGroup) GetMinCardNumber() int32 {
	return cg.GetMinCard().CardNumber
}

// GetMinCard 獲取最小牌型
func (cg *CardGroup) GetMinCard() *Card {
	minCard := cg.Cards[0]

	for _, c := range cg.Cards {
		if c.CardNumber < minCard.CardNumber {
			minCard = c
		} else if c.CardNumber == minCard.CardNumber {
			if c.IsMeiHua() {
				minCard = c
			} else if c.IsFangKuai() && !minCard.IsMeiHua() {
				minCard = c
			} else if c.IsHongTao() && !minCard.IsMeiHua() && !minCard.IsFangKuai() {
				minCard = c
			}
		}
	}

	return minCard
}

// GetMaxCardNumber 獲取最大牌點
func (cg *CardGroup) GetMaxCardNumber() int32 {
	return cg.GetMaxCard().CardNumber
}

// GetMaxCard 獲取最大牌型
func (cg *CardGroup) GetMaxCard() *Card {
	maxCard := cg.Cards[0]

	for _, c := range cg.Cards {
		if c.CardNumber > maxCard.CardNumber {
			maxCard = c
		} else if c.CardNumber == maxCard.CardNumber {
			if c.IsHeiTao() {
				maxCard = c
			} else if c.IsHongTao() && !maxCard.IsHeiTao() {
				maxCard = c
			} else if c.IsFangKuai() && !maxCard.IsHeiTao() && !maxCard.IsHongTao() {
				maxCard = c
			}
		}
	}

	return maxCard
}

// HasA 是否有A
func (cg *CardGroup) HasA() bool {
	for _, c := range cg.Cards {
		if c.IsA() {
			return true
		}
	}

	return false
}

// HasTenToThirteen 牌10到K的張數
func (cg *CardGroup) HasTenToThirteen() int {
	count := int(0)
	for _, c := range cg.Cards {
		if c.CardNumber >= int32(10) && c.CardNumber <= int32(13) {
			count++
		}
	}
	return count
}

// HasKingToPizi ?
func (cg *CardGroup) HasKingToPizi() int {
	count := int(0)
	for _, c := range cg.Cards {
		if c.CardNumber >= int32(15) && c.CardNumber <= int32(17) {
			count++
		}
	}
	return count
}

// IsShunzi 是否順子
func (cg *CardGroup) IsShunzi() bool {
	// sort.Sort(cg)
	if len(cg.Cards) == 0 {
		return false
	}

	if cg.HasA() && cg.IsShunziFromOne() {
		return true
	}

	min := cg.GetMinCardNumber()
	max := cg.GetMaxCardNumber()

	for i := min; i <= max; i++ {
		if len(cg.GetPokerCardByNumber(i)) != 1 {
			return false
		}
	}

	return true
}

// IsShunziFromOne 是否1點開始的順子
func (cg *CardGroup) IsShunziFromOne() bool {
	if cg.HasA() {
		max := len(cg.Cards)
		for i := 1; i <= max; i++ {
			if len(cg.GetPokerCardByNumber(int32(i))) != 1 {
				return false
			}
		}

		return true
	}

	return false
}

// GetPokerCardByNumber 根據牌點獲取牌型
func (cg *CardGroup) GetPokerCardByNumber(cn int32) []*Card {
	ncg := make([]*Card, 0)

	for _, c := range cg.Cards {
		if c.CardNumber == cn {
			ncg = append(ncg, c)
		}

		if cn == 1 && c.IsA() {
			ncg = append(ncg, c)
		}
	}

	return ncg
}

// GetPokerCardByColor 根據花色獲取牌型
func (cg *CardGroup) GetPokerCardByColor(color int32) []*Card {
	ncg := make([]*Card, 0)

	for _, c := range cg.Cards {
		if c.Color == color {
			ncg = append(ncg, c)
		}
	}

	return ncg
}

// IsTonghua 是否同花
func (cg *CardGroup) IsTonghua() bool {
	if len(cg.Cards) == 0 {
		return false
	}

	// 檢查是否同花
	color := cg.Cards[0].Color
	for k := 1; k < len(cg.Cards); k++ {
		c := cg.Cards[k]

		if c.Color != color {
			return false
		}
	}

	return true
}

// IsTongHuaShun 是否同花順
func (cg *CardGroup) IsTongHuaShun() bool {
	return cg.IsTonghua() && cg.IsShunzi()
}

// CardStack 牌堆
type CardStack struct {
	Groups []*CardGroup
}

// NewCardStack 創建牌堆對象
func NewCardStack() *CardStack {
	cst := new(CardStack)
	return cst
}

// AddCardGroup 添加牌組
func (cs *CardStack) AddCardGroup(cg *CardGroup) {
	cs.Groups = append(cs.Groups, cg)
}

// GetCardGroup 獲取牌組
func (cs *CardStack) GetCardGroup() *CardGroup {
	cg := new(CardGroup)
	for _, g := range cs.Groups {
		cg.AddCardGroup(g)
	}

	return cg
}

// UserCardStack 用戶牌堆
type UserCardStack struct {
	User  *gamelogic.GameUser
	Stack *CardStack
}

// NewUserCardStack 創建用戶牌堆
func NewUserCardStack(user *gamelogic.GameUser, stack *CardStack) *UserCardStack {
	ustack := new(UserCardStack)
	ustack.User = user
	ustack.Stack = stack

	return ustack
}

// AllCards 所有的牌
func (ucs *UserCardStack) AllCards() *CardGroup {
	cg := NewCardGroup(nil)

	for _, ug := range ucs.Stack.Groups {
		cg.AddCardGroup(ug)
	}

	return cg
}

// GameCardStack 遊戲牌堆
type GameCardStack struct {
	UserStacks []*UserCardStack
}

// NewGameCardStack 創建遊戲牌型
func NewGameCardStack() *GameCardStack {
	gcs := new(GameCardStack)

	return gcs
}

// CheckUserStack 檢查用戶牌型
func (gcs *GameCardStack) CheckUserStack(userid int32) bool {
	for _, ust := range gcs.UserStacks {
		if ust.User.UserID == userid {
			return true
		}
	}

	return false
}

// AddUserCardStack 添加用戶牌型
func (gcs *GameCardStack) AddUserCardStack(ucs *UserCardStack) {
	if gcs.CheckUserStack(ucs.User.UserID) {
		return
	}

	gcs.UserStacks = append(gcs.UserStacks, ucs)
}
