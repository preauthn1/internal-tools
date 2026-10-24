package poker

// Dealer 發牌器
type Dealer interface {
	// Deal 發牌
	//  @param src 要發牌的牌堆
	//  @param dest 要拿牌的牌堆
	//  @return err
	Deal(src *Deck, dest *Deck) error
}

// NaturalDealer 自然發牌器
//  從 src Pop() 一張卡牌加入 dest 牌堆
var NaturalDealer Dealer = new(naturalDealer)

type naturalDealer struct{}

func (dealer *naturalDealer) Deal(src *Deck, dest *Deck) error {

	card, err := src.Pop()
	if err != nil {
		return err
	}

	dest.AppendCard(card)

	return nil
}
