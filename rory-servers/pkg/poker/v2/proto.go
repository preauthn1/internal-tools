package poker

import "example.invalid/ServerTeam/Servers/pkg/netproto"

// TransePokerCard 轉換 Card 成 proto
func TransePokerCard(c Card) *netproto.Poker_Card {
	data := new(netproto.Poker_Card)
	data.Suit = netproto.Poker_Suit(c.Suit())
	data.Rank = netproto.Poker_Rank(c.Rank())
	return data
}

// ReversePokerCard 轉換 proto 成 Card
func ReversePokerCard(c *netproto.Poker_Card) Card {
	if c == nil {
		return 0
	}
	return Card(uint8(c.Rank) + uint8(c.Suit))
}

// TransePokerDeck 轉換 Deck 成 proto
func TransePokerDeck(d *Deck) *netproto.Poker_Deck {
	if d == nil {
		return nil
	}
	data := new(netproto.Poker_Deck)
	d.Range(func(i int, c Card) bool {
		data.Cards = append(data.Cards, TransePokerCard(c))
		return true
	})
	return data
}

// ReversePokerDeck 轉換 proto 成 Deck
func ReversePokerDeck(d *netproto.Poker_Deck) *Deck {
	if d == nil {
		return nil
	}
	deck := NewEmptyDeck()
	for _, c := range d.Cards {
		deck.AppendCard(ReversePokerCard(c))
	}
	return deck
}
