package inrm

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

func TestPointResult(t *testing.T) {
	t.Run("♥J [♣A ♣2 ♣3] [♦10 ♦J ♣Q ♦K] [♠3 ♥3 ♥J] [♠5 ♥5 ♦5] 1st 2st Life win, point 0", func(t *testing.T) {
		var hands []*poker.Deck

		d1, _ := poker.ParseDeck("♣A ♣2 ♣3")
		d2, _ := poker.ParseDeck("♦10 ♦J ♣Q ♦K")
		d3, _ := poker.ParseDeck("♠3 ♥3 ♥J")
		d4, _ := poker.ParseDeck("♠5 ♥5 ♦5")

		hands = append(hands, d1, d2, d3, d4)

		r := NewPointResult(poker.Jack, hands...)

		t.Logf("%+v", r)

		r.Valid()

		t.Logf("%+v", r)

		if !r.IsValidDeclaration() {
			t.Errorf("IsValidDeclaration want %t, but %t", true, r.IsValidDeclaration())
		}

		if r.Points() != 0 {
			t.Errorf("Points want %d, but %d", 0, r.Points())
		}
	})

	t.Run("♥J [♣A ♦2 ♣3] [♦10 ♦J ♣Q ♦K] [♠3 ♥3 ♥J] [♠5 ♥5 ♦5] no 1st Life no win, point 66", func(t *testing.T) {
		var hands []*poker.Deck

		d1, _ := poker.ParseDeck("♣A ♦2 ♣3")
		d2, _ := poker.ParseDeck("♦10 ♦J ♣Q ♦K")
		d3, _ := poker.ParseDeck("♠3 ♥3 ♥J")
		d4, _ := poker.ParseDeck("♠5 ♥5 ♦5")

		hands = append(hands, d1, d2, d3, d4)

		r := NewPointResult(poker.Jack, hands...)

		t.Logf("%+v", r)

		r.Valid()

		t.Logf("%+v", r)

		if r.IsValidDeclaration() {
			t.Errorf("IsValidDeclaration want %t, but %t", false, r.IsValidDeclaration())
		}

		// ♣A ♦2 ♣3 15
		// ♦10 ♦J ♣Q ♦K 30
		// ♠3 ♥3 ♥J 6
		// ♠5 ♥5 ♦5 15

		if r.Points() != 66 {
			t.Errorf("Points want %d, but %d", 66, r.Points())
		}
	})

	t.Run("♥J [♣A ♣2 ♣3] [♦10 ♦J ♣Q ♦K] [♠3 ♥3 ♥J] [♠5 ♥7 ♦5] 1st 2st Life not win, point 17", func(t *testing.T) {
		var hands []*poker.Deck

		d1, _ := poker.ParseDeck("♣A ♣2 ♣3")
		d2, _ := poker.ParseDeck("♦10 ♦J ♣Q ♦K")
		d3, _ := poker.ParseDeck("♠3 ♥3 ♥J")
		d4, _ := poker.ParseDeck("♠5 ♥7 ♦5")

		hands = append(hands, d1, d2, d3, d4)

		r := NewPointResult(poker.Jack, hands...)

		t.Logf("%+v", r)

		r.Valid()

		t.Logf("%+v", r)

		if r.IsValidDeclaration() {
			t.Errorf("IsValidDeclaration want %t, but %t", false, r.IsValidDeclaration())
		}

		// ♣A ♣2 ♣3 0
		// ♦10 ♦J ♣Q ♦K 0
		// ♠3 ♥3 ♥J 0
		// ♠5 ♥7 ♦5 17

		if r.Points() != 17 {
			t.Errorf("Points want %d, but %d", 17, r.Points())
		}
	})

	t.Run("♥J [♣A ♣2 ♣3] [♦10 ♦J ♣Q] [♠A ♥2 ♥3] [♠6 ♥7 ♦8] [♥7] 1st 2st Life not win, point 7", func(t *testing.T) {
		var hands []*poker.Deck

		d1, _ := poker.ParseDeck("♣A ♣2 ♣3")
		d2, _ := poker.ParseDeck("♦10 ♦J ♣Q")
		d3, _ := poker.ParseDeck("♠A ♥2 ♥3")
		d4, _ := poker.ParseDeck("♠6 ♥7 ♦8")
		d5, _ := poker.ParseDeck("♥7")

		hands = append(hands, d1, d2, d3, d4, d5)

		r := NewPointResult(poker.Jack, hands...)

		t.Logf("%+v", r)

		r.Valid()

		t.Logf("%+v", r)

		if r.IsValidDeclaration() {
			t.Errorf("IsValidDeclaration want %t, but %t", false, r.IsValidDeclaration())
		}

		// "♣A ♣2 ♣3" 0
		// "♦10 ♦J ♣Q" 0
		// "♠1 ♥2 ♥3" 0
		// "♠6 ♥7 ♦8" 0
		// "♥7" 7

		if r.Points() != 7 {
			t.Errorf("Points want %d, but %d", 7, r.Points())
		}
	})

	t.Run("♥J [♣A ♣2 ♣3] [♦10 ♦A ♣Q] [♠A ♥5 ♥3] [♠6 ♥A ♦8] [♥7] 1st no 2st Life not win, point 79", func(t *testing.T) {
		var hands []*poker.Deck

		d1, _ := poker.ParseDeck("♣A ♣2 ♣3")
		d2, _ := poker.ParseDeck("♦10 ♦A ♣Q")
		d3, _ := poker.ParseDeck("♠A ♥5 ♥3")
		d4, _ := poker.ParseDeck("♠6 ♥A ♦8")
		d5, _ := poker.ParseDeck("♥7")

		hands = append(hands, d1, d2, d3, d4, d5)

		r := NewPointResult(poker.Jack, hands...)

		t.Logf("%+v", r)

		r.Valid()

		t.Logf("%+v", r)

		if r.IsValidDeclaration() {
			t.Errorf("IsValidDeclaration want %t, but %t", false, r.IsValidDeclaration())
		}

		// "♣A ♣2 ♣3" 0
		// "♦10 A ♣Q" 30
		// "♠A ♥5 ♥3" 18
		// "♠6 ♥A ♦8" 24
		// "♥7" 7

		if r.Points() != 79 {
			t.Errorf("Points want %d, but %d", 7, r.Points())
		}
	})
}
