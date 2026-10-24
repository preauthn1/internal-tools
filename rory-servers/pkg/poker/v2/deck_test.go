package poker

import (
	"reflect"
	"sort"
	"testing"
)

func TestDeck(t *testing.T) {
	deck := NewDeck(1)
	if deck.Len() != 52 {
		t.Errorf("deck length want %d, but %d", 52, deck.Len())
	}
	t.Logf("%p", deck)
	t.Log(deck)
	t.Log(deck.Len())

	sort.Sort(deck)

	t.Run("Exist AceDiamonds", func(t *testing.T) {
		idx, ok := deck.Exist(Diamonds, Ace)
		if idx != 26 {
			t.Errorf("AceDiamonds idx want %d, but %d", 26, idx)
		}
		if !ok {
			t.Errorf("AceDiamonds want %t, but %t", true, ok)
		}
	})

	t.Run("ExistCard RedJoker", func(t *testing.T) {
		idx, ok := deck.ExistCard(RedJoker)
		if idx != -1 {
			t.Errorf("RedJoker idx want %d, but %d", -1, idx)
		}
		if ok {
			t.Errorf("RedJoker want %t, but %t", false, ok)
		}
	})

	t.Run("Remove", func(t *testing.T) {
		err := deck.remove(100)
		if err != ErrCardNotEnough {
			t.Errorf("deck.remove(100) want %v, but %v", ErrCardNotEnough, err)
		}
	})

	t.Run("CountOf", func(t *testing.T) {
		if cnt := deck.CountOf(Diamonds, Jack); cnt != 1 {
			t.Errorf("deck.CountOf(Diamonds, Jack) want %d, but %d", 1, cnt)
		}
	})

	t.Run("CountOfCard", func(t *testing.T) {
		if cnt := deck.CountOfCard(TwoClubs); cnt != 1 {
			t.Errorf("deck.CountOfCard(TwoClubs) want %d, but %d", 1, cnt)
		}
		if cnt := deck.CountOfCard(RedJoker); cnt != 0 {
			t.Errorf("deck.CountOfCard(TwoClubs) want %d, but %d", 0, cnt)
		}
	})

	t.Run("CountOfSuit", func(t *testing.T) {
		if cnt := deck.CountOfSuit(Spades); cnt != 13 {
			t.Errorf("deck.CountOfSuit(Spades) want %d, but %d", 13, cnt)
		}
		if cnt := deck.CountOfSuit(Hearts); cnt != 13 {
			t.Errorf("deck.CountOfSuit(Hearts) want %d, but %d", 13, cnt)
		}
		if cnt := deck.CountOfSuit(Diamonds); cnt != 13 {
			t.Errorf("deck.CountOfSuit(Diamonds) want %d, but %d", 13, cnt)
		}
		if cnt := deck.CountOfSuit(Clubs); cnt != 13 {
			t.Errorf("deck.CountOfSuit(Clubs) want %d, but %d", 13, cnt)
		}
	})

	t.Run("CountOfRank", func(t *testing.T) {
		if cnt := deck.CountOfRank(Queen); cnt != 4 {
			t.Errorf("deck.CountOfRank(Queen) want %d, but %d", 4, cnt)
		}
		if cnt := deck.CountOfRank(Ten); cnt != 4 {
			t.Errorf("deck.CountOfRank(Ten) want %d, but %d", 4, cnt)
		}

		if cnt := deck.CountOfRank(Ace); cnt != 4 {
			t.Errorf("deck.CountOfRank(Ace) want %d, but %d", 4, cnt)
		}

		deck.AppendCard(BlackJoker)
		deck.AppendCard(BlackJoker)
		deck.AppendCard(BlackJoker)

		if cnt := deck.CountOfRank(Ace); cnt != 4 {
			t.Errorf("deck.CountOfRank(Ace) when exist BlackJoker want %d, but %d", 4, cnt)
		}
	})
}

func TestDeck_Append(t *testing.T) {
	deck := NewEmptyDeck()
	deck.Append(Diamonds, King)
	deck.AppendCard(RedJoker)

	if deck.Cards[0] != KingDiamonds {
		t.Errorf("deck.Cards[0] want %v, but %v", KingDiamonds, deck.Cards[0])
	}

	if deck.Cards[1] != RedJoker {
		t.Errorf("deck.Cards[1] want %v, but %v", RedJoker, deck.Cards[1])
	}

	if deck.Len() != 2 {
		t.Errorf("deck.Len() want %v, but %v", 2, deck.Len())
	}

	d2 := NewEmptyDeck(KingSpades, AceDiamonds)
	d3 := NewEmptyDeck(JackClubs, QueenHearts)

	deck.AppendDeck(d2, d3)

	if deck.Cards[0] != KingDiamonds {
		t.Errorf("deck.Cards[0] want %v, but %v", KingDiamonds, deck.Cards[0])
	}

	if deck.Cards[1] != RedJoker {
		t.Errorf("deck.Cards[1] want %v, but %v", RedJoker, deck.Cards[1])
	}

	if deck.Cards[2] != KingSpades {
		t.Errorf("deck.Cards[2] want %v, but %v", KingSpades, deck.Cards[2])
	}

	if deck.Cards[3] != AceDiamonds {
		t.Errorf("deck.Cards[3] want %v, but %v", AceDiamonds, deck.Cards[3])
	}

	if deck.Cards[4] != JackClubs {
		t.Errorf("deck.Cards[4] want %v, but %v", JackClubs, deck.Cards[4])
	}

	if deck.Cards[5] != QueenHearts {
		t.Errorf("deck.Cards[5] want %v, but %v", QueenHearts, deck.Cards[5])
	}

	if deck.Len() != 6 {
		t.Errorf("deck.Len() want %v, but %v", 6, deck.Len())
	}
}

func TestDeck_Stack(t *testing.T) {
	deck := NewEmptyDeck()
	deck.Stack(Diamonds, King)
	t.Log(deck)
	deck.StackCard(RedJoker)
	t.Log(deck)

	if deck.Cards[1] != KingDiamonds {
		t.Errorf("deck.Cards[1] want %v, but %v", KingDiamonds, deck.Cards[1])
	}

	if deck.Cards[0] != RedJoker {
		t.Errorf("deck.Cards[0] want %v, but %v", RedJoker, deck.Cards[0])
	}

	if deck.Len() != 2 {
		t.Errorf("deck.Len() want %v, but %v", 2, deck.Len())
	}
}

func TestDeck_Pick(t *testing.T) {
	t.Run("Exist Card 1-6", func(t *testing.T) {
		deck := NewEmptyDeck()
		t.Run("Add", func(t *testing.T) {
			deck.Append(Diamonds, King)
			deck.AppendCard(RedJoker)
			deck.Append(Diamonds, Queen)
			deck.Append(Spades, King)
			deck.Append(Diamonds, Ace)
			deck.Append(Clubs, King)

			expCards := []Card{KingDiamonds, RedJoker, QueenDiamonds, KingSpades, AceDiamonds, KingClubs}

			if !reflect.DeepEqual(expCards, deck.Cards) {
				t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
			}
		})
		t.Run("3st", func(t *testing.T) {
			c, err := deck.Pick(Diamonds, Queen)

			if c != QueenDiamonds {
				t.Errorf("Pick(Diamonds, Queen) card want %v, but %v", QueenDiamonds, c)
			}

			if err != nil {
				t.Errorf("Pick(Diamonds, Queen) err want %v, but %v", nil, err)
			}

			expCards := []Card{KingDiamonds, RedJoker, KingSpades, AceDiamonds, KingClubs}

			if !reflect.DeepEqual(expCards, deck.Cards) {
				t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
			}
		})
		t.Run("5st", func(t *testing.T) {
			c, err := deck.Pick(Diamonds, Ace)

			if c != AceDiamonds {
				t.Errorf("Pick(Diamonds, Ace) card want %v, but %v", AceDiamonds, c)
			}

			if err != nil {
				t.Errorf("Pick(Diamonds, Ace) err want %v, but %v", nil, err)
			}

			expCards := []Card{KingDiamonds, RedJoker, KingSpades, KingClubs}

			if !reflect.DeepEqual(expCards, deck.Cards) {
				t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
			}
		})
		t.Run("2st", func(t *testing.T) {
			c, err := deck.PickCard(RedJoker)

			if c != RedJoker {
				t.Errorf("PickCard(RedJoker) card want %v, but %v", RedJoker, c)
			}

			if err != nil {
				t.Errorf("PickCard(RedJoker) err want %v, but %v", nil, err)
			}

			expCards := []Card{KingDiamonds, KingSpades, KingClubs}

			if !reflect.DeepEqual(expCards, deck.Cards) {
				t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
			}
		})
		t.Run("6st", func(t *testing.T) {
			c, err := deck.PickCard(KingClubs)

			if c != KingClubs {
				t.Errorf("PickCard(KingClubs) card want %v, but %v", KingClubs, c)
			}

			if err != nil {
				t.Errorf("PickCard(KingClubs) err want %v, but %v", nil, err)
			}

			expCards := []Card{KingDiamonds, KingSpades}

			if !reflect.DeepEqual(expCards, deck.Cards) {
				t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
			}
		})
		t.Run("1st", func(t *testing.T) {
			c, err := deck.PickCard(KingDiamonds)

			if c != KingDiamonds {
				t.Errorf("PickCard(KingDiamonds) card want %v, but %v", KingDiamonds, c)
			}

			if err != nil {
				t.Errorf("PickCard(KingDiamonds) err want %v, but %v", nil, err)
			}

			expCards := []Card{KingSpades}

			if !reflect.DeepEqual(expCards, deck.Cards) {
				t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
			}
		})
		t.Run("4st", func(t *testing.T) {
			c, err := deck.PickCard(KingSpades)

			if c != KingSpades {
				t.Errorf("PickCard(KingSpades) card want %v, but %v", KingSpades, c)
			}

			if err != nil {
				t.Errorf("PickCard(KingSpades) err want %v, but %v", nil, err)
			}

			expCards := []Card{}

			if !reflect.DeepEqual(expCards, deck.Cards) {
				t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
			}
		})
	})

	t.Run("Not Exist", func(t *testing.T) {
		deck := NewEmptyDeck()
		deck.Append(Diamonds, King)
		deck.AppendCard(RedJoker)
		deck.Append(Diamonds, Queen)
		deck.Append(Spades, King)
		deck.Append(Diamonds, Ace)
		deck.Append(Clubs, King)

		t.Log(deck)

		c, err := deck.Pick(Diamonds, Three)

		if c != 0 {
			t.Errorf("deck.Pick(Diamonds, Three) card want %v, but %v", 0, c)
		}

		if err != ErrCardNotFound {
			t.Errorf("deck.Pick(Diamonds, Three) err want %v, but %v", ErrCardNotFound, nil)
		}

	})
}

func TestDeck_Pop(t *testing.T) {
	deck := NewEmptyDeck()
	deck.Append(Diamonds, King)
	deck.AppendCard(RedJoker)
	deck.Append(Diamonds, Queen)
	deck.Append(Spades, King)
	deck.Append(Diamonds, Ace)
	deck.Append(Clubs, King)

	t.Run("Pop 1", func(t *testing.T) {
		c, err := deck.Pop()

		if c != KingDiamonds {
			t.Errorf("deck.Pop() card want %v, but %v", KingDiamonds, c)
		}

		if err != nil {
			t.Errorf("deck.Pop() err want %v, but %v", nil, err)
		}

		expCards := []Card{RedJoker, QueenDiamonds, KingSpades, AceDiamonds, KingClubs}

		if !reflect.DeepEqual(expCards, deck.Cards) {
			t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
		}
	})
	t.Run("Pop 2", func(t *testing.T) {
		c, err := deck.Pop()

		if c != RedJoker {
			t.Errorf("deck.Pop() card want %v, but %v", RedJoker, c)
		}

		if err != nil {
			t.Errorf("deck.Pop() err want %v, but %v", nil, err)
		}

		expCards := []Card{QueenDiamonds, KingSpades, AceDiamonds, KingClubs}

		if !reflect.DeepEqual(expCards, deck.Cards) {
			t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
		}
	})
	t.Run("Pop 3", func(t *testing.T) {
		c, err := deck.Pop()

		if c != QueenDiamonds {
			t.Errorf("deck.Pop() card want %v, but %v", QueenDiamonds, c)
		}

		if err != nil {
			t.Errorf("deck.Pop() err want %v, but %v", nil, err)
		}

		expCards := []Card{KingSpades, AceDiamonds, KingClubs}

		if !reflect.DeepEqual(expCards, deck.Cards) {
			t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
		}
	})
	t.Run("Pop 4", func(t *testing.T) {
		c, err := deck.Pop()

		if c != KingSpades {
			t.Errorf("deck.Pop() card want %v, but %v", KingSpades, c)
		}

		if err != nil {
			t.Errorf("deck.Pop() err want %v, but %v", nil, err)
		}

		expCards := []Card{AceDiamonds, KingClubs}

		if !reflect.DeepEqual(expCards, deck.Cards) {
			t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
		}
	})
	t.Run("Pop 5", func(t *testing.T) {
		c, err := deck.Pop()

		if c != AceDiamonds {
			t.Errorf("deck.Pop() card want %v, but %v", AceDiamonds, c)
		}

		if err != nil {
			t.Errorf("deck.Pop() err want %v, but %v", nil, err)
		}

		expCards := []Card{KingClubs}

		if !reflect.DeepEqual(expCards, deck.Cards) {
			t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
		}
	})
	t.Run("Pop 6", func(t *testing.T) {
		c, err := deck.Pop()

		if c != KingClubs {
			t.Errorf("deck.Pop() card want %v, but %v", KingClubs, c)
		}

		if err != nil {
			t.Errorf("deck.Pop() err want %v, but %v", nil, err)
		}

		expCards := []Card{}

		if !reflect.DeepEqual(expCards, deck.Cards) {
			t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
		}
	})
	t.Run("Pop 7", func(t *testing.T) {
		c, err := deck.Pop()

		if c != Card(0) {
			t.Errorf("deck.Pop() card want %v, but %v", Card(0), c)
		}

		if err != ErrCardNotEnough {
			t.Errorf("deck.Pop() err want %v, but %v", ErrCardNotEnough, err)
		}

		expCards := []Card{}

		if !reflect.DeepEqual(expCards, deck.Cards) {
			t.Errorf("deck.Cards want %v, but %v", expCards, deck.Cards)
		}
	})
}

func TestDeck_PickSuit(t *testing.T) {
	deck := NewDeck(1)
	deck.AppendCard(RedJoker, BlackJoker)

	t.Run("Pick Spades", func(t *testing.T) {
		cards := deck.PickSuit(Spades)

		expCards := []Card{AceSpades, TwoSpades, ThreeSpades, FourSpades, FiveSpades, SixSpades, SevenSpades, EightSpades, NineSpades, TenSpades, JackSpades, QueenSpades, KingSpades}

		if deck.Len() != 41 {
			t.Errorf("deck.Len() want %d, but %d", 41, deck.Len())
		}

		if !reflect.DeepEqual(expCards, cards) {
			t.Errorf("deck.PickSuit(Spades) want %v, but %v", expCards, cards)
		}
	})

	t.Run("Pick Joker", func(t *testing.T) {
		cards := deck.PickSuit(Jokers)

		expCards := []Card{RedJoker, BlackJoker}

		if deck.Len() != 39 {
			t.Errorf("deck.Len() want %d, but %d", 39, deck.Len())
		}

		if !reflect.DeepEqual(expCards, cards) {
			t.Errorf("deck.PickSuit(Spades) want %v, but %v", expCards, cards)
		}
	})

	t.Run("Pick Hearts", func(t *testing.T) {
		cards := deck.PickSuit(Hearts)

		expCards := []Card{AceHearts, TwoHearts, ThreeHearts, FourHearts, FiveHearts, SixHearts, SevenHearts, EightHearts, NineHearts, TenHearts, JackHearts, QueenHearts, KingHearts}

		if deck.Len() != 26 {
			t.Errorf("deck.Len() want %d, but %d", 26, deck.Len())
		}

		if !reflect.DeepEqual(expCards, cards) {
			t.Errorf("deck.PickSuit(Spades) want %v, but %v", expCards, cards)
		}
	})

	t.Run("Pick Diamonds", func(t *testing.T) {
		cards := deck.PickSuit(Diamonds)

		expCards := []Card{AceDiamonds, TwoDiamonds, ThreeDiamonds, FourDiamonds, FiveDiamonds, SixDiamonds, SevenDiamonds, EightDiamonds, NineDiamonds, TenDiamonds, JackDiamonds, QueenDiamonds, KingDiamonds}

		if deck.Len() != 13 {
			t.Errorf("deck.Len() want %d, but %d", 13, deck.Len())
		}

		if !reflect.DeepEqual(expCards, cards) {
			t.Errorf("deck.PickSuit(Spades) want %v, but %v", expCards, cards)
		}
	})

	t.Run("Pick Clubs", func(t *testing.T) {
		cards := deck.PickSuit(Clubs)

		expCards := []Card{AceClubs, TwoClubs, ThreeClubs, FourClubs, FiveClubs, SixClubs, SevenClubs, EightClubs, NineClubs, TenClubs, JackClubs, QueenClubs, KingClubs}

		if deck.Len() != 0 {
			t.Errorf("deck.Len() want %d, but %d", 0, deck.Len())
		}

		if !reflect.DeepEqual(expCards, cards) {
			t.Errorf("deck.PickSuit(Spades) want %v, but %v", expCards, cards)
		}
	})
}

func TestDeck_PickRank(t *testing.T) {
	deck := NewDeck(1)
	deck.AppendCard(RedJoker, BlackJoker)

	t.Run("Pick King", func(t *testing.T) {
		cards := deck.PickRank(King)

		expCards := []Card{KingSpades, KingHearts, KingDiamonds, KingClubs}

		if deck.Len() != 50 {
			t.Errorf("deck.Len() want %d, but %d", 50, deck.Len())
		}

		if !reflect.DeepEqual(expCards, cards) {
			t.Errorf("deck.PickSuit(Spades) want %v, but %v", expCards, cards)
		}
	})

	t.Run("Pick Ace", func(t *testing.T) {
		cards := deck.PickRank(Ace)

		expCards := []Card{AceSpades, AceHearts, AceDiamonds, AceClubs}

		if deck.Len() != 46 {
			t.Errorf("deck.Len() want %d, but %d", 46, deck.Len())
		}

		if !reflect.DeepEqual(expCards, cards) {
			t.Errorf("deck.PickSuit(Spades) want %v, but %v", expCards, cards)
		}
	})
}

func TestDeck_PickCards(t *testing.T) {
	deck := NewDeck(2)
	deck.AppendCard(RedJoker, BlackJoker)

	targets := NewDeck(1)

	t.Run("Pick King", func(t *testing.T) {
		cards := deck.PickCards(targets.Cards...)

		if len(cards) != 52 {
			t.Errorf("picked length want %d, but %d", 52, len(cards))
		}

		if deck.Len() != 54 {
			t.Errorf("deck.Len() want %d, but %d", 54, deck.Len())
		}
	})

}
