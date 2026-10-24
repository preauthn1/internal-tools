package mahjong

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

// CardGroup 牌組
type CardGroup struct {
	Cards []int32
}

// NewCardGroup 創建牌組對象
func NewCardGroup(cards []int32) *CardGroup {
	cg := new(CardGroup)
	if cards != nil {
		cg.Cards = append(cg.Cards, cards...)
	}
	return cg
}

// NewTong 筒子
func NewTong(index int32) *CardGroup {
	cg := new(CardGroup)
	for j := 0; j < 4; j++ {
		for i := int32(1); i <= 9; i++ {
			cd := ColorTong*100 + 10*i + index%10
			cg.AddCard(cd)
		}
	}
	return cg
}

// NewTiao 條
func NewTiao(index int32) *CardGroup {
	cg := new(CardGroup)
	for j := 0; j < 4; j++ {
		for i := int32(1); i <= 9; i++ {
			cd := ColorTiao*100 + 10*i + index%10
			cg.AddCard(cd)
		}
	}
	return cg
}

// NewWan 萬
func NewWan(index int32) *CardGroup {
	cg := new(CardGroup)
	for j := 0; j < 4; j++ {
		for i := int32(1); i <= 9; i++ {
			cd := ColorWan*100 + 10*i + index%10
			cg.AddCard(cd)
		}
	}
	return cg
}

// NewFeng 風牌
func NewFeng(index int32) *CardGroup {
	cg := new(CardGroup)
	for j := 0; j < 4; j++ {
		for i := int32(1); i <= 4; i++ {
			cd := ColorFeng*100 + 10*i + index%10
			cg.AddCard(cd)
		}
	}
	return cg
}

// NewJian 箭牌
func NewJian(index int32) *CardGroup {
	cg := new(CardGroup)
	for j := 0; j < 4; j++ {
		for i := int32(1); i <= 3; i++ {
			cd := ColorJain*100 + 10*i + index%10
			cg.AddCard(cd)
		}
	}
	return cg
}

// NewHua 花牌
func NewHua(index int32) *CardGroup {
	cg := new(CardGroup)
	for j := 0; j < 4; j++ {
		for i := int32(1); i <= 8; i++ {
			cd := ColorHua*100 + 10*i + index%10
			cg.AddCard(cd)
		}
	}
	return cg
}

// NewBaida 百搭
func NewBaida(index int32) *CardGroup {
	cg := new(CardGroup)
	for j := 0; j < 4; j++ {
		for i := int32(1); i <= 4; i++ {
			cd := ColorHua*100 + 10*i + index%10
			cg.AddCard(cd)
		}
	}
	return cg
}

// AddCard 添加牌
func (cg *CardGroup) AddCard(pk int32) {
	if pk < 110 {
		return
	}
	cg.Cards = append(cg.Cards, pk)
}

// AddCardGroup 添加牌组
func (cg *CardGroup) AddCardGroup(ncg *CardGroup) {
	cg.Cards = append(cg.Cards, ncg.Cards...)
}

// HaveCardOf 是否是手牌中的牌
func (cg *CardGroup) HaveCardOf(cd int32) bool {
	for _, card := range cg.Cards {
		if card == cd {
			return true
		}
	}
	return false
}

// IsEqual 是否同樣的牌組
func (cg *CardGroup) IsEqual(ncg *CardGroup) bool {
	if cg.Len() != ncg.Len() {
		return false
	}

	sort.Sort(cg)
	sort.Sort(ncg)

	for idx, c := range cg.Cards {
		if c != ncg.Cards[idx] {
			return false
		}
	}

	return true
}

// CalColorCount 統計特定花色張數
func (cg *CardGroup) CalColorCount(color int32) int32 {
	count := int32(0)
	for _, cd := range cg.Cards {
		if cd/100 == color {
			count++
		}
	}
	return count
}

func (cg *CardGroup) Len() int {
	return len(cg.Cards)
}

func (cg *CardGroup) Swap(i, j int) {
	cg.Cards[i], cg.Cards[j] = cg.Cards[j], cg.Cards[i]
}

func (cg *CardGroup) Less(i, j int) bool {
	return cg.Cards[i] < cg.Cards[j]
}

func (cg *CardGroup) String() string {
	var buf bytes.Buffer
	for _, c := range cg.Cards {
		buf.WriteString(fmt.Sprintf("%s,", CardString(c)))
	}
	return buf.String()
}

// DeleteGroup 删除手牌
func (cg *CardGroup) DeleteGroup(group *CardGroup) {
	for _, v := range group.Cards {
		cg.DeleteCard(v, 1)
	}
}

// DeleteCard 删除手牌
func (cg *CardGroup) DeleteCard(card, count int32) {
	tmp := make([]int32, 0)
	for _, v := range cg.Cards {
		if count > 0 && card == v {
			count--
			continue
		}
		tmp = append(tmp, v)
	}
	cg.Cards = tmp
}

// HavePairOf 是否可以碰牌
func (cg *CardGroup) HavePairOf(cd int32) bool {
	var count int32 = 0
	for _, card := range cg.Cards {
		if card == cd {
			count++
		}
	}
	if count >= 2 {
		return true
	}
	return false
}

// HaveThreeOf 是否可以引槓
func (cg *CardGroup) HaveThreeOf(cd int32) bool {
	var count int32 = 0
	for _, card := range cg.Cards {
		if card == cd {
			count++
		}
	}
	if count >= 3 {
		return true
	}
	return false
}

// HaveFour 查詢暗槓牌,暗槓
func (cg *CardGroup) HaveFour() []int32 {
	pc := make([]int32, 0)
	mapCard := make(map[int32]int32)
	for _, card := range cg.Cards {
		mapCard[card]++
	}
	for card, num := range mapCard {
		if num == 4 {
			pc = append(pc, card)
		}
	}
	return pc
}

// HaveFourOf 是否可以槓牌,暗槓
func (cg *CardGroup) HaveFourOf(cd int32) bool {
	sum := 0
	for _, card := range cg.Cards {
		if card == cd {
			sum++
		}
	}
	if sum >= 4 {
		return true
	}
	return false
}

// HaveColor 檢測牌中是否有某種花色的牌
func (cg *CardGroup) HaveColor(color int32) bool {
	for _, card := range cg.Cards {
		if card/100 == color {
			return true
		}
	}
	return false
}

// CopyCardGroup 深度拷貝
func (cg *CardGroup) CopyCardGroup() *CardGroup {
	newCards := new(CardGroup)
	newCards.Cards = cg.Cards
	return newCards
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

// DeepCopy ?
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
