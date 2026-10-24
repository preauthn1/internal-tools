package poker

import (
	"fmt"
	"strings"
)

// A Rank represents the rank of a card.
type Rank uint8

var (
	ranksStr = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	ranksMap = map[string]Rank{
		"A":  Ace,
		"2":  Two,
		"3":  Three,
		"4":  Four,
		"5":  Five,
		"6":  Six,
		"7":  Seven,
		"8":  Eight,
		"9":  Nine,
		"10": Ten,
		"J":  Jack,
		"Q":  Queen,
		"K":  King,
	}
)

// Ranks 16進制第一位
const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) aceLastIndexOf() int {
	for i, rank := range RegularRanksAceLast {
		if r == rank {
			return i
		}
	}
	return -1
}

func (r Rank) String() string {
	idx := r - 1
	if idx < 0 || int(idx) > len(ranksStr) {
		return "?"
	}
	return ranksStr[idx]
}

// Suit 花色 16進制第二位
type Suit uint8

// Suits
const (
	// 黑桃
	Spades Suit = (iota + 1) * 0x10
	// 紅心
	Hearts
	// 方塊
	Diamonds
	// 梅花
	Clubs
	// 鬼牌
	Jokers
)

var (
	suitsStr = []string{"\u2660", "\u2665", "\u2666", "\u2663", "\U0001F0DF"}
	suitsMap = map[string]Suit{
		"\u2660":     Spades,
		"\u2665":     Hearts,
		"\u2666":     Diamonds,
		"\u2663":     Clubs,
		"\U0001F0DF": Jokers,
	}
)

func (s Suit) String() string {
	idx := s/0x10 - 1
	if idx < 0 || int(idx) > len(suitsStr) {
		return "?"
	}
	return suitsStr[idx]
}

// Card 牌 16進制兩位
type Card uint8

// Spades Card
const (
	AceSpades Card = (iota + 1) + 0x10
	TwoSpades
	ThreeSpades
	FourSpades
	FiveSpades
	SixSpades
	SevenSpades
	EightSpades
	NineSpades
	TenSpades
	JackSpades
	QueenSpades
	KingSpades
)

// Hearts Card
const (
	AceHearts Card = (iota + 1) + 0x20
	TwoHearts
	ThreeHearts
	FourHearts
	FiveHearts
	SixHearts
	SevenHearts
	EightHearts
	NineHearts
	TenHearts
	JackHearts
	QueenHearts
	KingHearts
)

// Diamonds Card
const (
	AceDiamonds Card = (iota + 1) + 0x30
	TwoDiamonds
	ThreeDiamonds
	FourDiamonds
	FiveDiamonds
	SixDiamonds
	SevenDiamonds
	EightDiamonds
	NineDiamonds
	TenDiamonds
	JackDiamonds
	QueenDiamonds
	KingDiamonds
)

// Clubs Card
const (
	AceClubs Card = (iota + 1) + 0x40
	TwoClubs
	ThreeClubs
	FourClubs
	FiveClubs
	SixClubs
	SevenClubs
	EightClubs
	NineClubs
	TenClubs
	JackClubs
	QueenClubs
	KingClubs
)

// Jokers Card
const (
	RedJoker Card = iota + 0x50
	BlackJoker
	WhiteJoker
)

// Rank returns the rank of the card.
func (c Card) Rank() Rank {
	return Rank(c & 0x0F)
}

// Suit returns the suit of the card.
func (c Card) Suit() Suit {
	return Suit(c & 0xF0)
}

func (c Card) String() string {
	switch c.Suit() {
	default:
		return "?"
	case Jokers:
		switch c {
		default:
			return c.Suit().String()
		case RedJoker:
			return c.Suit().String() + "R"
		case BlackJoker:
			return c.Suit().String() + "B"
		case WhiteJoker:
			return c.Suit().String() + "W"
		}
	case Spades, Hearts, Diamonds, Clubs:
		return c.Suit().String() + c.Rank().String()
	}
}

// Unicode 使用Unicode顯示
//  Ref: https://example.invalid/charts/PDF/U1F0A0.pdf
func (c Card) Unicode() string {
	switch c {
	default:
		return "\U0001F0A0"
	// Spades 1F0A
	case AceSpades:
		return "\U0001F0A1"
	case TwoSpades:
		return "\U0001F0A2"
	case ThreeSpades:
		return "\U0001F0A3"
	case FourSpades:
		return "\U0001F0A4"
	case FiveSpades:
		return "\U0001F0A5"
	case SixSpades:
		return "\U0001F0A6"
	case SevenSpades:
		return "\U0001F0A7"
	case EightSpades:
		return "\U0001F0A8"
	case NineSpades:
		return "\U0001F0A9"
	case TenSpades:
		return "\U0001F0AA"
	case JackSpades:
		return "\U0001F0AB"
	case QueenSpades:
		return "\U0001F0AD"
	case KingSpades:
		return "\U0001F0AE"
	// Hearts 1F0B
	case AceHearts:
		return "\U0001F0B1"
	case TwoHearts:
		return "\U0001F0B2"
	case ThreeHearts:
		return "\U0001F0B3"
	case FourHearts:
		return "\U0001F0B4"
	case FiveHearts:
		return "\U0001F0B5"
	case SixHearts:
		return "\U0001F0B6"
	case SevenHearts:
		return "\U0001F0B7"
	case EightHearts:
		return "\U0001F0B8"
	case NineHearts:
		return "\U0001F0B9"
	case TenHearts:
		return "\U0001F0BA"
	case JackHearts:
		return "\U0001F0BB"
	case QueenHearts:
		return "\U0001F0BD"
	case KingHearts:
		return "\U0001F0BE"
	// Diamonds 1F0C
	case AceDiamonds:
		return "\U0001F0C1"
	case TwoDiamonds:
		return "\U0001F0C2"
	case ThreeDiamonds:
		return "\U0001F0C3"
	case FourDiamonds:
		return "\U0001F0C4"
	case FiveDiamonds:
		return "\U0001F0C5"
	case SixDiamonds:
		return "\U0001F0C6"
	case SevenDiamonds:
		return "\U0001F0C7"
	case EightDiamonds:
		return "\U0001F0C8"
	case NineDiamonds:
		return "\U0001F0C9"
	case TenDiamonds:
		return "\U0001F0CA"
	case JackDiamonds:
		return "\U0001F0CB"
	case QueenDiamonds:
		return "\U0001F0CD"
	case KingDiamonds:
		return "\U0001F0CE"
	// Clubs 1F0D
	case AceClubs:
		return "\U0001F0D1"
	case TwoClubs:
		return "\U0001F0D2"
	case ThreeClubs:
		return "\U0001F0D3"
	case FourClubs:
		return "\U0001F0D4"
	case FiveClubs:
		return "\U0001F0D5"
	case SixClubs:
		return "\U0001F0D6"
	case SevenClubs:
		return "\U0001F0D7"
	case EightClubs:
		return "\U0001F0D8"
	case NineClubs:
		return "\U0001F0D9"
	case TenClubs:
		return "\U0001F0DA"
	case JackClubs:
		return "\U0001F0DB"
	case QueenClubs:
		return "\U0001F0DD"
	case KingClubs:
		return "\U0001F0DE"
	// Jokers
	case RedJoker:
		return "\U0001F0BF"
	case BlackJoker:
		return "\U0001F0CF"
	case WhiteJoker:
		return "\U0001F0DF"
	}
}

var simpleCards []Card = []Card{
	AceSpades, TwoSpades, ThreeSpades, FourSpades, FiveSpades, SixSpades, SevenSpades, EightSpades, NineSpades, TenSpades, JackSpades, QueenSpades, KingSpades,
	AceHearts, TwoHearts, ThreeHearts, FourHearts, FiveHearts, SixHearts, SevenHearts, EightHearts, NineHearts, TenHearts, JackHearts, QueenHearts, KingHearts,
	AceDiamonds, TwoDiamonds, ThreeDiamonds, FourDiamonds, FiveDiamonds, SixDiamonds, SevenDiamonds, EightDiamonds, NineDiamonds, TenDiamonds, JackDiamonds, QueenDiamonds, KingDiamonds,
	AceClubs, TwoClubs, ThreeClubs, FourClubs, FiveClubs, SixClubs, SevenClubs, EightClubs, NineClubs, TenClubs, JackClubs, QueenClubs, KingClubs,
}

// RankOrder 點數排序
type RankOrder []Rank

// RegularRanks 所有正常點數
var RegularRanks RankOrder = RankOrder{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

// RegularRanksAceLast 所有正常點數
var RegularRanksAceLast RankOrder = RankOrder{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}

// RegularSuits 所有正常花色
var RegularSuits []Suit = []Suit{Spades, Hearts, Diamonds, Clubs}

// ParseCards 傳入卡片文字轉換
//  Ex: ♦2 ♠4 ♣5 ♠6 ♠7 ♦7 ♠8 ♦8 ♥9 ♦10 ♦Q ♥K ♣K
func ParseCards(s string) ([]Card, error) {
	strs := strings.Split(s, " ")

	var ans []Card

	for _, str := range strs {
		var runes []rune
		for _, r := range str {
			runes = append(runes, r)
		}

		// 一定要有花色跟點數兩個字元
		if len(runes) < 2 {
			return nil, fmt.Errorf("%s must has suit and rank", str)
		}

		var suitStr string = fmt.Sprintf("%c", runes[0])
		var rankStr string
		for _, r := range runes[1:] {
			rankStr += fmt.Sprintf("%c", r)
		}

		var suit Suit
		var rank Rank
		// 第一字元為花色
		if val, ok := suitsMap[suitStr]; ok {
			suit = val
		} else {
			return nil, fmt.Errorf("%s suit is not available", str)
		}

		switch suit {
		default:
			return nil, fmt.Errorf("%s suit is not available", str)
		case Spades, Hearts, Diamonds, Clubs:
			if val, ok := ranksMap[rankStr]; ok {
				rank = val
				ans = append(ans, Card(uint8(suit)+uint8(rank)))
			} else {
				return nil, fmt.Errorf("%s rank is not available", str)
			}
		case Jokers:
			switch rankStr {
			default:
				return nil, fmt.Errorf("%s joker is not available", str)
			case "R":
				ans = append(ans, RedJoker)
			case "B":
				ans = append(ans, BlackJoker)
			case "W":
				ans = append(ans, WhiteJoker)
			}
		}
	}

	return ans, nil
}

// ParseDeck 傳入卡片文字轉換
//  Ex: ♦2 ♠4 ♣5 ♠6 ♠7 ♦7 ♠8 ♦8 ♥9 ♦10 ♦Q ♥K ♣K
func ParseDeck(s string) (*Deck, error) {
	strs := strings.Split(s, " ")

	var ans []Card

	for _, str := range strs {
		var runes []rune
		for _, r := range str {
			runes = append(runes, r)
		}

		// 一定要有花色跟點數兩個字元
		if len(runes) < 2 {
			return nil, fmt.Errorf("%s must has suit and rank", str)
		}

		var suitStr string = fmt.Sprintf("%c", runes[0])
		var rankStr string
		for _, r := range runes[1:] {
			rankStr += fmt.Sprintf("%c", r)
		}

		var suit Suit
		var rank Rank
		// 第一字元為花色
		if val, ok := suitsMap[suitStr]; ok {
			suit = val
		} else {
			return nil, fmt.Errorf("%s suit is not available", str)
		}

		switch suit {
		default:
			return nil, fmt.Errorf("%s suit is not available", str)
		case Spades, Hearts, Diamonds, Clubs:
			if val, ok := ranksMap[rankStr]; ok {
				rank = val
				ans = append(ans, Card(uint8(suit)+uint8(rank)))
			} else {
				return nil, fmt.Errorf("%s rank is not available", str)
			}
		case Jokers:
			switch rankStr {
			default:
				return nil, fmt.Errorf("%s joker is not available", str)
			case "R":
				ans = append(ans, RedJoker)
			case "B":
				ans = append(ans, BlackJoker)
			case "W":
				ans = append(ans, WhiteJoker)
			}
		}
	}

	return NewEmptyDeck(ans...), nil
}
