package inrm

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

func TestFirstLifeValidator(t *testing.T) {
	t.Run("Pass Spades A234", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.ThreeSpades, poker.FourSpades}
		ok := FirstLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Spades A23", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.ThreeSpades}
		ok := FirstLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Spades 234", func(t *testing.T) {
		cards := []poker.Card{poker.TwoSpades, poker.ThreeSpades, poker.FourSpades}
		ok := FirstLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Spades AJQK", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.JackSpades, poker.QueenSpades, poker.KingSpades}
		ok := FirstLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Spades AQK", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.QueenSpades, poker.KingSpades}
		ok := FirstLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Reject Length not 3 or 4 Spades 34", func(t *testing.T) {
		cards := []poker.Card{poker.ThreeSpades, poker.FourSpades}
		ok := FirstLifeValidator().Validate(cards)
		if ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
	t.Run("Pass Spades A2345", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.ThreeSpades, poker.FourSpades, poker.FiveSpades}
		ok := FirstLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Reject Diff Suit Spades A234", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.ThreeClubs, poker.FourSpades}
		ok := FirstLifeValidator().Validate(cards)
		if ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
	t.Run("Reject Exist Joker Spades A2\U0001F0DF4 ", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.RedJoker, poker.FourSpades}
		ok := FirstLifeValidator().Validate(cards)
		if ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
	t.Run("Reject Not Straight Spades A245 ", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.FourSpades, poker.FiveSpades}
		ok := FirstLifeValidator().Validate(cards)
		if ok {
			t.Errorf("FirstLifeValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
}

func TestSecondLifeValidator(t *testing.T) {
	t.Run("Pass Spades A234", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.ThreeSpades, poker.FourSpades}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Spades A23", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.ThreeSpades}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Spades 234", func(t *testing.T) {
		cards := []poker.Card{poker.TwoSpades, poker.ThreeSpades, poker.FourSpades}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Spades AJQK", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.JackSpades, poker.QueenSpades, poker.KingSpades}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Spades AQK", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.QueenSpades, poker.KingSpades}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Reject Length not 3 or 4 Spades 34", func(t *testing.T) {
		cards := []poker.Card{poker.ThreeSpades, poker.FourSpades}
		ok := SecondLifeValidator().Validate(cards)
		if ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
	t.Run("Pass Spades A2345", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.ThreeSpades, poker.FourSpades, poker.FiveSpades}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Diff Suit A234", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoHearts, poker.ThreeClubs, poker.FourDiamonds}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Exist Joker A2\U0001F0DF4", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoHearts, poker.RedJoker, poker.FourDiamonds}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Exist Joker A\U0001F0DF3", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.RedJoker, poker.ThreeClubs}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Exist Joker A\U0001F0DFQK", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.BlackJoker, poker.QueenSpades, poker.KingSpades}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Exist Joker AQ\U0001F0DF", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.QueenSpades, poker.WhiteJoker}
		ok := SecondLifeValidator().Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Exist ExtraRank 7 A274 ", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoHearts, poker.SevenClubs, poker.FourDiamonds}
		ok := SecondLifeValidator(poker.Seven).Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Exist ExtraRank 7 A7QK", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.SevenDiamonds, poker.QueenSpades, poker.KingSpades}
		ok := SecondLifeValidator(poker.Seven).Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Exist ExtraRank 7 AQ7", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.QueenSpades, poker.SevenHearts}
		ok := SecondLifeValidator(poker.Seven).Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Pass Special Case Exist ExtraRank 7 789", func(t *testing.T) {
		cards := []poker.Card{poker.SevenClubs, poker.EightClubs, poker.NineDiamonds}
		ok := SecondLifeValidator(poker.Seven).Validate(cards)
		if !ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Reject Not Straight Spades A245", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.TwoSpades, poker.FourSpades, poker.FiveSpades}
		ok := SecondLifeValidator().Validate(cards)
		if ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
	t.Run("Reject Not Straight Exist Joker Spades A35\U0001F0DF", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.ThreeClubs, poker.FiveHearts, poker.BlackJoker}
		ok := SecondLifeValidator(poker.Seven).Validate(cards)
		if ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
	t.Run("Reject Not Straight Exist Joker Spades A5\U0001F0DF", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.FiveHearts, poker.BlackJoker}
		ok := SecondLifeValidator(poker.Seven).Validate(cards)
		if ok {
			t.Errorf("SecondLifeValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
}

func TestSetValidator(t *testing.T) {
	t.Run("No Wild 4 card \u2660A\u2661A\u2662A\u2663A Pass", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.AceHearts, poker.AceDiamonds, poker.AceClubs}
		ok := SetValidator().Validate(cards)
		if !ok {
			t.Errorf("SetValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("No Wild 3 card \u2660A\u2662A\u2663A Pass", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.AceDiamonds, poker.AceClubs}
		ok := SetValidator().Validate(cards)
		if !ok {
			t.Errorf("SetValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("No Wild 2 card \u2660A\u2663A Fail", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.AceClubs}
		ok := SetValidator().Validate(cards)
		if ok {
			t.Errorf("SetValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
	t.Run("No Wild 5 card \u2660A\u2661A\u2662A\u2663A\u2662A Fail", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.AceHearts, poker.AceDiamonds, poker.AceClubs, poker.AceDiamonds}
		ok := SetValidator().Validate(cards)
		if ok {
			t.Errorf("SetValidator.Validate() ok want %v, but %v", false, ok)
		}
	})
	t.Run("Wild Joker 3 card \u2660A\u2661A\U0001F0DF Fail", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.AceHearts, poker.RedJoker}
		ok := SetValidator().Validate(cards)
		if !ok {
			t.Errorf("SetValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
	t.Run("Wild 7 3 card \u2660A\u2661A\u26607 Fail", func(t *testing.T) {
		cards := []poker.Card{poker.AceSpades, poker.AceHearts, poker.SevenSpades}
		ok := SetValidator(poker.Seven).Validate(cards)
		if !ok {
			t.Errorf("SetValidator.Validate() ok want %v, but %v", true, ok)
		}
	})
}
