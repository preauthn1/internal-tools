package poker

import (
	"testing"
)

func TestCard_Ace(t *testing.T) {
	if AceSpades.Rank() != Ace {
		t.Errorf("AceSpades.Rank() want %v, but %v", Ace, AceSpades.Rank())
	}
	if AceHearts.Rank() != Ace {
		t.Errorf("AceHearts.Rank() want %v, but %v", Ace, AceHearts.Rank())
	}
	if AceDiamonds.Rank() != Ace {
		t.Errorf("AceDiamonds.Rank() want %v, but %v", Ace, AceDiamonds.Rank())
	}
	if AceClubs.Rank() != Ace {
		t.Errorf("AceClubs.Rank() want %v, but %v", Ace, AceClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", AceSpades, AceHearts, AceDiamonds, AceClubs)
	t.Log(AceSpades, AceHearts, AceDiamonds, AceClubs)
	t.Log(AceSpades.Unicode(), AceHearts.Unicode(), AceDiamonds.Unicode(), AceClubs.Unicode())
}

func TestCard_Two(t *testing.T) {
	if TwoSpades.Rank() != Two {
		t.Errorf("TwoSpades.Rank() want %v, but %v", Two, TwoSpades.Rank())
	}
	if TwoHearts.Rank() != Two {
		t.Errorf("TwoHearts.Rank() want %v, but %v", Two, TwoHearts.Rank())
	}
	if TwoDiamonds.Rank() != Two {
		t.Errorf("TwoDiamonds.Rank() want %v, but %v", Two, TwoDiamonds.Rank())
	}
	if TwoClubs.Rank() != Two {
		t.Errorf("TwoClubs.Rank() want %v, but %v", Two, TwoClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", TwoSpades, TwoHearts, TwoDiamonds, TwoClubs)
	t.Log(TwoSpades, TwoHearts, TwoDiamonds, TwoClubs)
	t.Log(TwoSpades.Unicode(), TwoHearts.Unicode(), TwoDiamonds.Unicode(), TwoClubs.Unicode())
}

func TestCard_Three(t *testing.T) {
	if ThreeSpades.Rank() != Three {
		t.Errorf("ThreeSpades.Rank() want %v, but %v", Three, ThreeSpades.Rank())
	}
	if ThreeHearts.Rank() != Three {
		t.Errorf("ThreeHearts.Rank() want %v, but %v", Three, ThreeHearts.Rank())
	}
	if ThreeDiamonds.Rank() != Three {
		t.Errorf("ThreeDiamonds.Rank() want %v, but %v", Three, ThreeDiamonds.Rank())
	}
	if ThreeClubs.Rank() != Three {
		t.Errorf("ThreeClubs.Rank() want %v, but %v", Three, ThreeClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", ThreeSpades, ThreeHearts, ThreeDiamonds, ThreeClubs)
	t.Log(ThreeSpades, ThreeHearts, ThreeDiamonds, ThreeClubs)
	t.Log(ThreeSpades.Unicode(), ThreeHearts.Unicode(), ThreeDiamonds.Unicode(), ThreeClubs.Unicode())
}

func TestCard_Four(t *testing.T) {
	if FourSpades.Rank() != Four {
		t.Errorf("FourSpades.Rank() want %v, but %v", Four, FourSpades.Rank())
	}
	if FourHearts.Rank() != Four {
		t.Errorf("FourHearts.Rank() want %v, but %v", Four, FourHearts.Rank())
	}
	if FourDiamonds.Rank() != Four {
		t.Errorf("FourDiamonds.Rank() want %v, but %v", Four, FourDiamonds.Rank())
	}
	if FourClubs.Rank() != Four {
		t.Errorf("FourClubs.Rank() want %v, but %v", Four, FourClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", FourSpades, FourHearts, FourDiamonds, FourClubs)
	t.Log(FourSpades, FourHearts, FourDiamonds, FourClubs)
	t.Log(FourSpades.Unicode(), FourHearts.Unicode(), FourDiamonds.Unicode(), FourClubs.Unicode())
}

func TestCard_Five(t *testing.T) {
	if FiveSpades.Rank() != Five {
		t.Errorf("FiveSpades.Rank() want %v, but %v", Five, FiveSpades.Rank())
	}
	if FiveHearts.Rank() != Five {
		t.Errorf("FiveHearts.Rank() want %v, but %v", Five, FiveHearts.Rank())
	}
	if FiveDiamonds.Rank() != Five {
		t.Errorf("FiveDiamonds.Rank() want %v, but %v", Five, FiveDiamonds.Rank())
	}
	if FiveClubs.Rank() != Five {
		t.Errorf("FiveClubs.Rank() want %v, but %v", Five, FiveClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", FiveSpades, FiveHearts, FiveDiamonds, FiveClubs)
	t.Log(FiveSpades, FiveHearts, FiveDiamonds, FiveClubs)
	t.Log(FiveSpades.Unicode(), FiveHearts.Unicode(), FiveDiamonds.Unicode(), FiveClubs.Unicode())
}

func TestCard_Six(t *testing.T) {
	if SixSpades.Rank() != Six {
		t.Errorf("SixSpades.Rank() want %v, but %v", Six, SixSpades.Rank())
	}
	if SixHearts.Rank() != Six {
		t.Errorf("SixHearts.Rank() want %v, but %v", Six, SixHearts.Rank())
	}
	if SixDiamonds.Rank() != Six {
		t.Errorf("SixDiamonds.Rank() want %v, but %v", Six, SixDiamonds.Rank())
	}
	if SixClubs.Rank() != Six {
		t.Errorf("SixClubs.Rank() want %v, but %v", Six, SixClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", SixSpades, SixHearts, SixDiamonds, SixClubs)
	t.Log(SixSpades, SixHearts, SixDiamonds, SixClubs)
	t.Log(SixSpades.Unicode(), SixHearts.Unicode(), SixDiamonds.Unicode(), SixClubs.Unicode())
}

func TestCard_Seven(t *testing.T) {
	if SevenSpades.Rank() != Seven {
		t.Errorf("SevenSpades.Rank() want %v, but %v", Seven, SevenSpades.Rank())
	}
	if SevenHearts.Rank() != Seven {
		t.Errorf("SevenHearts.Rank() want %v, but %v", Seven, SevenHearts.Rank())
	}
	if SevenDiamonds.Rank() != Seven {
		t.Errorf("SevenDiamonds.Rank() want %v, but %v", Seven, SevenDiamonds.Rank())
	}
	if SevenClubs.Rank() != Seven {
		t.Errorf("SevenClubs.Rank() want %v, but %v", Seven, SevenClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", SevenSpades, SevenHearts, SevenDiamonds, SevenClubs)
	t.Log(SevenSpades, SevenHearts, SevenDiamonds, SevenClubs)
	t.Log(SevenSpades.Unicode(), SevenHearts.Unicode(), SevenDiamonds.Unicode(), SevenClubs.Unicode())
}

func TestCard_Eight(t *testing.T) {
	if EightSpades.Rank() != Eight {
		t.Errorf("EightSpades.Rank() want %v, but %v", Eight, EightSpades.Rank())
	}
	if EightHearts.Rank() != Eight {
		t.Errorf("EightHearts.Rank() want %v, but %v", Eight, EightHearts.Rank())
	}
	if EightDiamonds.Rank() != Eight {
		t.Errorf("EightDiamonds.Rank() want %v, but %v", Eight, EightDiamonds.Rank())
	}
	if EightClubs.Rank() != Eight {
		t.Errorf("EightClubs.Rank() want %v, but %v", Eight, EightClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", EightSpades, EightHearts, EightDiamonds, EightClubs)
	t.Log(EightSpades, EightHearts, EightDiamonds, EightClubs)
	t.Log(EightSpades.Unicode(), EightHearts.Unicode(), EightDiamonds.Unicode(), EightClubs.Unicode())
}

func TestCard_Nine(t *testing.T) {
	if NineSpades.Rank() != Nine {
		t.Errorf("NineSpades.Rank() want %v, but %v", Nine, NineSpades.Rank())
	}
	if NineHearts.Rank() != Nine {
		t.Errorf("NineHearts.Rank() want %v, but %v", Nine, NineHearts.Rank())
	}
	if NineDiamonds.Rank() != Nine {
		t.Errorf("NineDiamonds.Rank() want %v, but %v", Nine, NineDiamonds.Rank())
	}
	if NineClubs.Rank() != Nine {
		t.Errorf("NineClubs.Rank() want %v, but %v", Nine, NineClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", NineSpades, NineHearts, NineDiamonds, NineClubs)
	t.Log(NineSpades, NineHearts, NineDiamonds, NineClubs)
	t.Log(NineSpades.Unicode(), NineHearts.Unicode(), NineDiamonds.Unicode(), NineClubs.Unicode())
}

func TestCard_Ten(t *testing.T) {
	if TenSpades.Rank() != Ten {
		t.Errorf("TenSpades.Rank() want %v, but %v", Ten, TenSpades.Rank())
	}
	if TenHearts.Rank() != Ten {
		t.Errorf("TenHearts.Rank() want %v, but %v", Ten, TenHearts.Rank())
	}
	if TenDiamonds.Rank() != Ten {
		t.Errorf("TenDiamonds.Rank() want %v, but %v", Ten, TenDiamonds.Rank())
	}
	if TenClubs.Rank() != Ten {
		t.Errorf("TenClubs.Rank() want %v, but %v", Ten, TenClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", TenSpades, TenHearts, TenDiamonds, TenClubs)
	t.Log(TenSpades, TenHearts, TenDiamonds, TenClubs)
	t.Log(TenSpades.Unicode(), TenHearts.Unicode(), TenDiamonds.Unicode(), TenClubs.Unicode())
}

func TestCard_Jack(t *testing.T) {
	if JackSpades.Rank() != Jack {
		t.Errorf("JackSpades.Rank() want %v, but %v", Jack, JackSpades.Rank())
	}
	if JackHearts.Rank() != Jack {
		t.Errorf("JackHearts.Rank() want %v, but %v", Jack, JackHearts.Rank())
	}
	if JackDiamonds.Rank() != Jack {
		t.Errorf("JackDiamonds.Rank() want %v, but %v", Jack, JackDiamonds.Rank())
	}
	if JackClubs.Rank() != Jack {
		t.Errorf("JackClubs.Rank() want %v, but %v", Jack, JackClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", JackSpades, JackHearts, JackDiamonds, JackClubs)
	t.Log(JackSpades, JackHearts, JackDiamonds, JackClubs)
	t.Log(JackSpades.Unicode(), JackHearts.Unicode(), JackDiamonds.Unicode(), JackClubs.Unicode())
}

func TestCard_Queen(t *testing.T) {
	if QueenSpades.Rank() != Queen {
		t.Errorf("QueenSpades.Rank() want %v, but %v", Queen, QueenSpades.Rank())
	}
	if QueenHearts.Rank() != Queen {
		t.Errorf("QueenHearts.Rank() want %v, but %v", Queen, QueenHearts.Rank())
	}
	if QueenDiamonds.Rank() != Queen {
		t.Errorf("QueenDiamonds.Rank() want %v, but %v", Queen, QueenDiamonds.Rank())
	}
	if QueenClubs.Rank() != Queen {
		t.Errorf("QueenClubs.Rank() want %v, but %v", Queen, QueenClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", QueenSpades, QueenHearts, QueenDiamonds, QueenClubs)
	t.Log(QueenSpades, QueenHearts, QueenDiamonds, QueenClubs)
	t.Log(QueenSpades.Unicode(), QueenHearts.Unicode(), QueenDiamonds.Unicode(), QueenClubs.Unicode())
}

func TestCard_King(t *testing.T) {
	if KingSpades.Rank() != King {
		t.Errorf("KingSpades.Rank() want %v, but %v", King, KingSpades.Rank())
	}
	if KingHearts.Rank() != King {
		t.Errorf("KingHearts.Rank() want %v, but %v", King, KingHearts.Rank())
	}
	if KingDiamonds.Rank() != King {
		t.Errorf("KingDiamonds.Rank() want %v, but %v", King, KingDiamonds.Rank())
	}
	if KingClubs.Rank() != King {
		t.Errorf("KingClubs.Rank() want %v, but %v", King, KingClubs.Rank())
	}
	t.Logf("%#v %#v %#v %#v", KingSpades, KingHearts, KingDiamonds, KingClubs)
	t.Log(KingSpades, KingHearts, KingDiamonds, KingClubs)
	t.Log(KingSpades.Unicode(), KingHearts.Unicode(), KingDiamonds.Unicode(), KingClubs.Unicode())
}

func TestCard_Jocker(t *testing.T) {
	t.Logf("%#v %#v %#v", RedJoker, BlackJoker, WhiteJoker)
	t.Log(RedJoker, BlackJoker, WhiteJoker)
	t.Log(RedJoker.Unicode(), BlackJoker.Unicode(), WhiteJoker.Unicode())

	t.Log(Card(3).Unicode())
}

func TestCard(t *testing.T) {
	t.Log(Card(uint8(Spades) + uint8(Ace)))
	t.Log(Card(0), Card(99))
	t.Log(Rank(0), Suit(0))
}

func TestParseCard(t *testing.T) {
	t.Log(ParseCards("♦2 ♠4 ♣5 ♠6 ♠7 ♦7 🃟R ♠8 ♦8 ♥9 ♦10 ♦Q ♥K ♣K"))
}
