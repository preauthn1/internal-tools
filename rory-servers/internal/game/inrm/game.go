package inrm

import (
	"fmt"
	"math/rand"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

// GState 遊戲狀態
type GState uint8

// GStates
const (
	GStatePrepare GState = iota
	GStateGame
	GStateCalc
	GStateOver
)

func (s GState) String() string {
	switch s {
	default:
		return ""
	case GStatePrepare:
		return "準備中"
	case GStateGame:
		return "遊戲中"
	case GStateCalc:
		return "結算中"
	case GStateOver:
		return "遊戲結束"
	}
}

// Game 印度拉米(INRM)
type Game struct {
	HandCardLength int            // 正常狀態下手牌數量 Default: 13
	GiveupPoint    int            // 放棄遊戲時，輸分點數 Default: 40
	MaxPoint       int            // 最大點數 Default: 80
	MissLimit      int            // 錯過牌局次數上限 Default: 3
	ShuffleTimes   int            // 每局洗牌次數 Default: 3
	UsePokerNum    int            // 使用幾副牌 Default: 2
	Shuffler       poker.Shuffler // 洗牌器 Default: poker.RandShuffler	// 亂數種子 Default: time.Now().Unitnano()
	Dealer         poker.Dealer   // 發牌器 Default: poker.NaturalDealer

	RandJoker poker.Rank  // 隨機Joker
	Players   []*Player   // 玩家
	CloseDeck *poker.Deck // 蓋牌區
	OpenDeck  *poker.Deck // 亮牌區

	State GState // 遊戲狀態
	Round int    // 目前輪迴數
	Idx   int    // 目前玩家

	OnGameStart      func()
	OnRefreshHands   func(p *Player)
	OnPlayerPopClose func(c poker.Card, p *Player)
	OnPlayerPopOpen  func(c poker.Card, p *Player)
	OnPlayerThrow    func(c poker.Card, p *Player)
	OnPlayerGiveup   func(p *Player)
	OnPlayerMiss     func(p *Player)
	OnPlayerBluff    func(p *Player, c poker.Card, hand []*poker.Deck)
	OnPlayerVictory  func(p *Player, c poker.Card, hand []*poker.Deck)
	OnPlayerPass     func(prev, next *Player)
	OnDeckReset      func()
	OnGameOver       func()
}

func (game *Game) String() string {
	if game.State == GStatePrepare {
		return fmt.Sprintf("[%v]", game.State)
	}
	var open, close poker.Card
	if game.OpenDeck.Len() > 0 {
		open = game.OpenDeck.Cards[0]
	}
	if game.CloseDeck.Len() > 0 {
		close = game.CloseDeck.Cards[0]
	}
	return fmt.Sprintf("[%v] Round[%d] [%dP][%d] 百搭[%v] 蓋牌堆[%d][%v] 開牌堆[%d][%v]", game.State, game.Round, len(game.Players), game.Players[game.Idx].ID, game.RandJoker, game.CloseDeck.Len(), close, game.OpenDeck.Len(), open)
}

// NewGame 建立印度拉米(INRM)
//  @param pids 	玩家們id
//  @param options 	建立選項
func NewGame(pids []int32, options ...GameOpt) (*Game, error) {
	// 至少需求2位玩家才能開始牌局
	if len(pids) < 2 {
		return nil, ErrPlayerNotEnough
	}

	game := new(Game)
	game.OnGameStart = func() {}
	game.OnRefreshHands = func(p *Player) {}
	game.OnPlayerPopClose = func(c poker.Card, p *Player) {}
	game.OnPlayerPopOpen = func(c poker.Card, p *Player) {}
	game.OnPlayerThrow = func(c poker.Card, p *Player) {}
	game.OnPlayerGiveup = func(p *Player) {}
	game.OnPlayerMiss = func(p *Player) {}
	game.OnPlayerBluff = func(p *Player, c poker.Card, hand []*poker.Deck) {}
	game.OnPlayerVictory = func(p *Player, c poker.Card, hand []*poker.Deck) {}
	game.OnPlayerPass = func(prev, next *Player) {}
	game.OnDeckReset = func() {}
	game.OnGameOver = func() {}

	// 設定預設值
	game.HandCardLength = 13
	game.GiveupPoint = 40
	game.MaxPoint = 80
	game.MissLimit = 3
	game.ShuffleTimes = 3
	game.UsePokerNum = 2
	src := rand.NewSource(time.Now().UnixNano())
	game.Shuffler = poker.NewRandShuffler(src)
	game.Dealer = poker.NaturalDealer

	for _, opt := range options {
		opt.apply(game)
	}

	// 初始化玩家
	for _, pid := range pids {
		game.Players = append(game.Players, NewPlayer(pid))
	}

	return game, nil
}

// Start 遊戲開始
//  1. 遊戲狀態為 GStateGame
//  2. 決定一張牌為Joker
//  3. 用n副(52張)的牌組作自然發牌 + 2張鬼牌
//  4. 洗牌n次
//  5. 每人發n張牌, 整理手牌牌堆
//  6. 從蓋牌堆翻一張牌亮出放亮牌堆第一張
//  7. 指定起始玩家
func (game *Game) Start() error {
	defer game.OnGameStart()

	game.State = GStateGame

	// 隨機決定一張牌為Joker
	idx := rand.Intn(len(poker.RegularRanks))
	game.RandJoker = poker.RegularRanks[idx]

	// 用n副(52張+2張鬼牌)的牌組作自然發牌
	game.OpenDeck = poker.NewEmptyDeck()
	game.CloseDeck = poker.NewDeck(game.UsePokerNum)          // 用n副(52張+2張鬼牌)的牌組作自然發牌
	game.CloseDeck.AppendCard(poker.RedJoker, poker.RedJoker) // n組鬼牌

	// 洗牌n次
	game.Shuffler.Shuffle(game.CloseDeck, game.ShuffleTimes)

	// 每人發n張牌
	for i := 0; i < game.HandCardLength; i++ {
		for _, p := range game.Players {
			err := game.Dealer.Deal(game.CloseDeck, p.Deck)
			if err != nil {
				return err
			}
		}
	}

	// 整理牌堆
	for _, p := range game.Players {
		p.Deck.Sort(poker.SortSameRank)
		p.initHands()
		// RefreshHandsActionOpt(p.Deck).Do(game, p)
	}

	// 翻一張牌亮出放在open deck區
	c, err := game.CloseDeck.Pop()
	if err != nil {
		return err
	}
	game.OpenDeck.StackCard(c)

	// 指定起始玩家
	game.Players[game.Idx].State = PStateGameV1

	game.Round = 1

	return nil
}

// Action 玩家動作
func (game *Game) Action(uid int32, opt ActionOpt) error {
	p, err := game.Player(uid)
	if err != nil {
		return err
	}
	return opt.Do(game, p)
}

// Player 取得玩家
func (game *Game) Player(uid int32) (*Player, error) {
	for _, p := range game.Players {
		if p.ID == uid {
			return p, nil
		}
	}
	return nil, ErrPlayerNotFound
}

// AlivePlayerCount 取得存活玩家數
func (game *Game) AlivePlayerCount() int {
	var cnt int
	for _, p := range game.Players {
		if p.Alive() {
			cnt++
		}
	}
	return cnt
}

// CurrentPlayer 取得當前玩家
func (game *Game) CurrentPlayer() (*Player, error) {
	if game.State != GStateGame {
		return nil, ErrGameState
	}
	return game.Players[game.Idx], nil
}

// NextPlayer 把指針指向下一人，並回傳下一個人
func (game *Game) NextPlayer() (*Player, error) {
	// 先查都是活的
	if game.AlivePlayerCount() == 0 {
		return nil, ErrAllPlayerDead
	}

	var next = func() *Player {
		game.Idx++
		if game.Idx >= len(game.Players) {
			game.Round++
			game.Idx -= len(game.Players)
		}
		return game.Players[game.Idx]
	}

	var n *Player = next()

	// 下一人若不是活著就往下一個
	for ; !n.Alive(); n = next() {
	}

	return n, nil
}

// ValidateVectory 檢查牌型是否勝利
//  一定要是 13 張
//  需要至少有 1 個 1st Life
//  需要至少有 1 個 2st Life
//  其餘不可有散牌，一定要組為 Set 或 順
func (game *Game) ValidateVectory(hands []*poker.Deck) bool {

	// 一定要是 13 張
	deck := poker.NewEmptyDeck()
	deck.AppendDeck(hands...)
	if deck.Len() != game.HandCardLength {
		return false
	}

	r := NewPointResult(game.RandJoker, hands...)
	return r.IsValidDeclaration()
}

// CalcPlayerLossPoint 計算玩家分數
func (game *Game) CalcPlayerLossPoint(p *Player) {
	switch p.State {
	default:
		r := NewPointResult(game.RandJoker, p.Hands...)
		p.LossPoint = r.Points()
		if p.LossPoint > game.MaxPoint {
			p.LossPoint = game.MaxPoint
		}
		p.State = PStateLoss
	case PStateVictory:
		p.LossPoint = 0
	case PStateFirstGiveup:
		p.LossPoint = game.GiveupPoint / 2
	case PStateGiveup:
		p.LossPoint = game.GiveupPoint
	case PStateBluff:
		p.LossPoint = game.MaxPoint
	}
}

// Category 牌型
type Category int

// Categorys
const (
	FirstLife Category = iota + 1
	SecondLife
	Set
	None
)
