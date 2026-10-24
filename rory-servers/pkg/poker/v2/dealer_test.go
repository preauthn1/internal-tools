package poker

import (
	"reflect"
	"testing"
)

func TestNaturalDealer(t *testing.T) {
	t.Run("One Poker", func(t *testing.T) {
		deck := NewDeck(1)
		hand := NewEmptyDeck()
		t.Run("Deal 1", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 51 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 51, deck.Len())
			}
		})
		t.Run("Deal 2", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades, TwoSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 50 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 50, deck.Len())
			}
		})
		t.Run("Deal 3", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades, TwoSpades, ThreeSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 49 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 49, deck.Len())
			}
		})
		t.Run("Deal 4", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades, TwoSpades, ThreeSpades, FourSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 48 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 48, deck.Len())
			}
		})
		t.Run("Deal 5", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades, TwoSpades, ThreeSpades, FourSpades, FiveSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 47 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 47, deck.Len())
			}
		})
	})
	t.Run("Not enough", func(t *testing.T) {
		deck := NewEmptyDeck()
		deck.AppendCard(AceSpades)
		deck.AppendCard(TwoSpades)
		deck.AppendCard(ThreeSpades)
		deck.AppendCard(FourSpades)

		hand := NewEmptyDeck()

		t.Run("Deal 1", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 3 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 3, deck.Len())
			}

			expDeck := []Card{TwoSpades, ThreeSpades, FourSpades}

			if !reflect.DeepEqual(expDeck, deck.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck want %v, but %v", expDeck, deck)
			}
		})
		t.Run("Deal 2", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades, TwoSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 2 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 2, deck.Len())
			}

			expDeck := []Card{ThreeSpades, FourSpades}

			if !reflect.DeepEqual(expDeck, deck.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck want %v, but %v", expDeck, deck)
			}
		})
		t.Run("Deal 3", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades, TwoSpades, ThreeSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 1 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 1, deck.Len())
			}

			expDeck := []Card{FourSpades}

			if !reflect.DeepEqual(expDeck, deck.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck want %v, but %v", expDeck, deck)
			}
		})
		t.Run("Deal 4", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != nil {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", nil, err)
			}

			expHands := []Card{AceSpades, TwoSpades, ThreeSpades, FourSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 0 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 0, deck.Len())
			}

			expDeck := []Card{}

			if !reflect.DeepEqual(expDeck, deck.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck want %v, but %v", expDeck, deck)
			}
		})
		t.Run("Deal 5", func(t *testing.T) {
			err := NaturalDealer.Deal(deck, hand)
			if err != ErrCardNotEnough {
				t.Errorf("NaturalDealer.Deal(deck, hand) err want %v, but %v", ErrCardNotEnough, err)
			}

			expHands := []Card{AceSpades, TwoSpades, ThreeSpades, FourSpades}

			if !reflect.DeepEqual(expHands, hand.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) hand want %v, but %v", expHands, hand)
			}

			if deck.Len() != 0 {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck.Len() want %d, but %d", 0, deck.Len())
			}

			expDeck := []Card{}

			if !reflect.DeepEqual(expDeck, deck.Cards) {
				t.Errorf("NaturalDealer.Deal(deck, hand) deck want %v, but %v", expDeck, deck)
			}
		})
	})
}
