package poker

import "testing"

func TestSizeValidator(t *testing.T) {
	t.Run("Single Size Pass", func(t *testing.T) {
		v := SizeValidator(5)
		ok := v.Validate([]Card{AceSpades, AceDiamonds, AceClubs, AceHearts, TenClubs})
		if !ok {
			t.Errorf("SizeValidate(5) want %v, but %v", true, ok)
		}
	})
	t.Run("Single Size Fail", func(t *testing.T) {
		v := SizeValidator(5)
		ok := v.Validate([]Card{AceSpades, AceDiamonds})
		if ok {
			t.Errorf("SizeValidate(5) want %v, but %v", false, ok)
		}
	})
	t.Run("Multiple Size Pass 3", func(t *testing.T) {
		v := SizeValidator(3, 4)
		ok := v.Validate([]Card{AceSpades, AceDiamonds, AceClubs})
		if !ok {
			t.Errorf("SizeValidate(3, 4) want %v, but %v", true, ok)
		}
	})
	t.Run("Multiple Size Pass 4", func(t *testing.T) {
		v := SizeValidator(3, 4)
		ok := v.Validate([]Card{AceSpades, AceDiamonds, AceClubs, AceHearts})
		if !ok {
			t.Errorf("SizeValidate(3, 4) want %v, but %v", true, ok)
		}
	})
	t.Run("Multiple Size Fail", func(t *testing.T) {
		v := SizeValidator(3, 4)
		ok := v.Validate([]Card{AceSpades, AceDiamonds})
		if ok {
			t.Errorf("SizeValidate(3, 4) want %v, but %v", false, ok)
		}
	})
}

func TestSizeLessValidator(t *testing.T) {
	t.Run("Size Over Pass", func(t *testing.T) {
		v := SizeLessValidator(3)
		ok := v.Validate([]Card{AceSpades, AceDiamonds, AceClubs, AceHearts, TenClubs})
		if !ok {
			t.Errorf("SizeLessValidator(3) want %v, but %v", true, ok)
		}
	})
	t.Run("Size Equal Pass", func(t *testing.T) {
		v := SizeLessValidator(3)
		ok := v.Validate([]Card{AceSpades, AceDiamonds, AceClubs})
		if !ok {
			t.Errorf("SizeLessValidator(3) want %v, but %v", true, ok)
		}
	})
	t.Run("Size Less Fail", func(t *testing.T) {
		v := SizeLessValidator(3)
		ok := v.Validate([]Card{AceSpades, AceDiamonds})
		if ok {
			t.Errorf("SizeLessValidator(3) want %v, but %v", false, ok)
		}
	})
}

func TestSequenceValidator(t *testing.T) {
	t.Run("No Wild RegularRanks A2345", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, nil, nil, nil)
		ok := v.Validate([]Card{AceSpades, TwoClubs, ThreeHearts, FourClubs, FiveHearts})
		if !ok {
			t.Errorf("StraightValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("No Wild RegularRanks JQKA", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, nil, nil, nil)
		ok := v.Validate([]Card{JackClubs, QueenClubs, KingClubs, AceDiamonds})
		if ok {
			t.Errorf("StraightValidator() want %v, but %v", false, ok)
		}
	})
	t.Run("No Wild RegularRanksAceLast JQKA", func(t *testing.T) {
		v := SequenceValidator(RegularRanksAceLast, nil, nil, nil)
		ok := v.Validate([]Card{JackClubs, QueenClubs, KingClubs, AceDiamonds})
		if !ok {
			t.Errorf("StraightValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild Joker RegularRanks A2\U0001F0DF45", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, nil, []Suit{Jokers}, nil)
		ok := v.Validate([]Card{AceSpades, TwoClubs, RedJoker, FourClubs, FiveHearts})
		if !ok {
			t.Errorf("StraightValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild Joker RegularRanks J\U0001F0DFKA", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, nil, []Suit{Jokers}, nil)
		ok := v.Validate([]Card{JackClubs, RedJoker, KingClubs, AceDiamonds})
		if ok {
			t.Errorf("StraightValidator() want %v, but %v", false, ok)
		}
	})
	t.Run("Wild Joker RegularRanksAceLast J\U0001F0DFKA", func(t *testing.T) {
		v := SequenceValidator(RegularRanksAceLast, nil, []Suit{Jokers}, nil)
		ok := v.Validate([]Card{JackClubs, RedJoker, KingClubs, AceDiamonds})
		if !ok {
			t.Errorf("StraightValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild Joker RegularRanks A\U0001F0DF5", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, nil, []Suit{Jokers}, nil)
		ok := v.Validate([]Card{AceDiamonds, RedJoker, FiveClubs})
		if ok {
			t.Errorf("StraightValidator() want %v, but %v", false, ok)
		}
	})
	t.Run("Wild 7 RegularRanks A2457", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, nil, nil, []Rank{Seven})
		ok := v.Validate([]Card{AceSpades, TwoClubs, FourClubs, FiveHearts, SevenDiamonds})
		if !ok {
			t.Errorf("StraightValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild 7 RegularRanks 7JKA", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, nil, nil, []Rank{Seven})
		ok := v.Validate([]Card{SevenSpades, JackHearts, KingClubs, AceDiamonds})
		if ok {
			t.Errorf("StraightValidator() want %v, but %v", false, ok)
		}
	})
	t.Run("Wild 7 RegularRanksAceLast 7JKA", func(t *testing.T) {
		v := SequenceValidator(RegularRanksAceLast, nil, nil, []Rank{Seven})
		ok := v.Validate([]Card{SevenClubs, JackClubs, KingClubs, AceDiamonds})
		if !ok {
			t.Errorf("StraightValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild \u26607 RegularRanks A245\u26607", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, []Card{SevenSpades}, nil, nil)
		ok := v.Validate([]Card{AceSpades, TwoClubs, FourClubs, FiveHearts, SevenSpades})
		if !ok {
			t.Errorf("StraightValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild \u26607 RegularRanks A245\u26667", func(t *testing.T) {
		v := SequenceValidator(RegularRanks, []Card{SevenSpades}, nil, nil)
		ok := v.Validate([]Card{AceSpades, TwoClubs, FourClubs, FiveHearts, SevenDiamonds})
		if ok {
			t.Errorf("StraightValidator() want %v, but %v", false, ok)
		}
	})
}

func TestSameSuitValidator(t *testing.T) {
	t.Run("No Wild \u2660A\u26602\u26603\u26604\u26605 Pass", func(t *testing.T) {
		v := SameSuitValidator(nil, nil, nil)
		ok := v.Validate([]Card{AceSpades, TwoSpades, ThreeSpades, FourSpades, FiveSpades})
		if !ok {
			t.Errorf("SameSuitValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("No Wild \u26613\u26615\u266110\u2661J Pass", func(t *testing.T) {
		v := SameSuitValidator(nil, nil, nil)
		ok := v.Validate([]Card{ThreeHearts, FiveHearts, TenHearts, JackHearts})
		if !ok {
			t.Errorf("SameSuitValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("No Wild \u26613\u26635\u266110\u2661J Fail", func(t *testing.T) {
		v := SameSuitValidator(nil, nil, nil)
		ok := v.Validate([]Card{ThreeHearts, FiveClubs, TenHearts, JackHearts})
		if ok {
			t.Errorf("SameSuitValidator() want %v, but %v", false, ok)
		}
	})
	t.Run("Wild Joker \u26613\u266110\u2661J\U0001F0DF Pass", func(t *testing.T) {
		v := SameSuitValidator(nil, []Suit{Jokers}, nil)
		ok := v.Validate([]Card{ThreeHearts, TenHearts, JackHearts, RedJoker})
		if !ok {
			t.Errorf("SameSuitValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild 7 \u26613\u26637\u266110\u2661J Pass", func(t *testing.T) {
		v := SameSuitValidator(nil, nil, []Rank{Seven})
		ok := v.Validate([]Card{ThreeHearts, SevenDiamonds, TenHearts, JackHearts})
		if !ok {
			t.Errorf("SameSuitValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild \u26607 \u26613\u26607\u266110\u2661J Pass", func(t *testing.T) {
		v := SameSuitValidator([]Card{SevenSpades}, nil, nil)
		ok := v.Validate([]Card{ThreeHearts, SevenSpades, TenHearts, JackHearts})
		if !ok {
			t.Errorf("SameSuitValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild \u26607 \u26613\u26637\u266110\u2661J Fail", func(t *testing.T) {
		v := SameSuitValidator([]Card{SevenSpades}, nil, nil)
		ok := v.Validate([]Card{ThreeHearts, SevenDiamonds, TenHearts, JackHearts})
		if ok {
			t.Errorf("SameSuitValidator() want %v, but %v", false, ok)
		}
	})
}

func TestSameRankValidator(t *testing.T) {
	t.Run("No Wild \u2660A\u2661A\u2662A\u2663A Pass", func(t *testing.T) {
		v := SameRankValidator(nil, nil, nil)
		ok := v.Validate([]Card{AceSpades, AceHearts, AceDiamonds, AceClubs})
		if !ok {
			t.Errorf("SameRankValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("No Wild Exist BlackJoker \u2660A\u2661A\u2662A\U0001F0DF Fail", func(t *testing.T) {
		v := SameRankValidator(nil, nil, nil)
		ok := v.Validate([]Card{AceSpades, AceHearts, AceDiamonds, BlackJoker})
		if ok {
			t.Errorf("SameRankValidator() want %v, but %v", false, ok)
		}
	})
	t.Run("Wild Joker \u2660A\u2661A\u2662A\U0001F0DF Pass", func(t *testing.T) {
		v := SameRankValidator(nil, []Suit{Jokers}, nil)
		ok := v.Validate([]Card{AceSpades, AceHearts, AceDiamonds, BlackJoker})
		if !ok {
			t.Errorf("SameRankValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild 7 \u2660A\u2661A\u2662A\u26637 Pass", func(t *testing.T) {
		v := SameRankValidator(nil, nil, []Rank{Seven})
		ok := v.Validate([]Card{AceSpades, AceHearts, AceDiamonds, SevenDiamonds})
		if !ok {
			t.Errorf("SameRankValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild \u26607 \u2660A\u2661A\u2662A\u26607 Pass", func(t *testing.T) {
		v := SameRankValidator([]Card{SevenSpades}, nil, nil)
		ok := v.Validate([]Card{AceSpades, AceHearts, AceDiamonds, SevenSpades})
		if !ok {
			t.Errorf("SameSuitValidator() want %v, but %v", true, ok)
		}
	})
	t.Run("Wild \u26607 \u2660A\u2661A\u2662A\u26637 Fail", func(t *testing.T) {
		v := SameRankValidator([]Card{SevenSpades}, nil, nil)
		ok := v.Validate([]Card{AceSpades, AceHearts, AceDiamonds, SevenDiamonds})
		if ok {
			t.Errorf("SameSuitValidator() want %v, but %v", false, ok)
		}
	})
}
