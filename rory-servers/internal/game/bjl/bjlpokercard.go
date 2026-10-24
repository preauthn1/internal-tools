package bjl

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func (bjl *BaiJialeLogic) getXDuiCards(isWin bool, isDui bool) ([]*card, bool) {
	cards := make([]*card, 0)
	card1 := bjl.getCardbyIndex(0)
	card2 := bjl.getCardbyIndex(1)
	if isDui {
		if !isWin && bjl.zuhetype&(1<<2) == 0 { //閒沒贏也不是和
			//mlog.Debug("card1 %v, card2 %v", card1.number, card2)
			if (card1.number == 4 || card1.number == 9) && card1.number == card2.number { //閒沒贏對子不可以是4或9
				return cards, false //重新發牌
			} else if card1.number == 4 || card1.number == 9 {
				bjl.delCard(card2)
				card1 = bjl.getCardByNumber(card2.number)
				bjl.delCard(card1)
			} else {
				bjl.delCard(card1)
				card2 = bjl.getCardByNumber(card1.number)
				bjl.delCard(card2)
			}
		}
		if card1.number != card2.number {
			bjl.delCard(card1)
			card2 = bjl.getCardByNumber(card1.number)
			bjl.delCard(card2)
		}

	} else {
		nArray := make([]int, 0)
		for i:=int32(1);i<=13;i++ {
			if card1.number != i {
				tmpcard := bjl.getCardByNumber(i)
				if tmpcard.number == 0 {
					continue
				}
				dian := bjl.getCardsDian([]*card{card1, tmpcard})
				if bjl.zuhetype&(1<<2) != 0 {
					nArray = append(nArray, int(i))
				} else if isWin {
					if dian > 4 && dian < 8 {
						nArray = append(nArray, int(i))
					}
				} else {
					if dian < 6 && dian < 8 {
						nArray = append(nArray, int(i))
					}
				}
			}
		}
		if len(nArray) < 1 {
			return cards, false
		}
		bjl.delCard(card1)
		n := util.IntArrayShuffle(nArray)[0]
		card2 = bjl.getCardByNumber(int32(n))
		bjl.delCard(card2)
	}

	cards = append(cards, card1, card2)
	xiandian := bjl.getCardsDian(cards)
	if xiandian < 6 { //閒補牌
		if bjl.zuhetype&(1<<1) != 0 { //閒贏 6以上
			n := bjl.getRndNeedDian(xiandian, 6, 10)
			card := bjl.getCardByNumber(n)
			bjl.delCard(card)
			cards = append(cards, card)
		} else { //閒輸開小一點
			n := bjl.getRndNeedDian(xiandian, 0, 5)
			card := bjl.getCardByNumber(n)
			bjl.delCard(card)
			cards = append(cards, card)
		}
		if cards[0].number == cards[1].number && cards[1].number == cards[2].number {
			return cards, false //重新發牌
		}
	}
	//如果有沒取到牌就直接返回錯誤
	restoreFlag := false
	for _, card := range cards {
		if card.number == 0 {
			restoreFlag = true
		}
	}
	if restoreFlag {
		for _, card := range cards {
			if card.number != 0 {
				bjl.restoreCard(card)
			}
		}
		return cards, false
	}

	return cards, true
}
func (bjl *BaiJialeLogic) getRndNeedDian(nowDian int32, minDian int, maxDian int) int32 {
	targetDian := int32(util.RndInt(minDian+10, maxDian+10))
	n := targetDian - nowDian
	if n > 10 {
		n = n - 10
	}
	return n
}

func (bjl *BaiJialeLogic) getZDuiCards(xiancards []*card, isWin bool, isDui bool) ([]*card, bool) {
	cards := make([]*card, 0)
	xiandian := bjl.getCardsDian(xiancards)
	hasThird := true
	if bjl.zuhetype&(1<<2) != 0 && (xiandian % 2) != 0 { //開和
		return cards, false
	}
	duiDianMap := map[int]int32{
		1:2, 2:4, 3:6, 4:8, 5:0, 6:2, 7:4, 8:6, 9:8, 10:0, 11:0, 12:0, 13:0,
	}
	if isDui {
		nArray := make([]int, 0)
		if len(xiancards) > 2 {
			for i:=1;i<=13;i++ {
				if isWin {
					if duiDianMap[i] > xiandian {
						nArray = append(nArray, i)
					}
				} else if bjl.zuhetype&(1<<2) != 0 {
					if duiDianMap[i] <= 2 || duiDianMap[i] == xiandian {
						nArray = append(nArray, i)
					}
				} else {
					if duiDianMap[i] < xiandian {
						nArray = append(nArray, i)
					}
				}
			}
			//閒有補牌不可以是8點以上
			if ok, idx := util.InArray(4, nArray);ok {
				util.IntArrayRemoveIndex(nArray, idx)
			}
			if ok, idx := util.InArray(9, nArray);ok {
				util.IntArrayRemoveIndex(nArray, idx)
			}
			if len(nArray) < 1 {
				return cards, false
			}
			n := util.IntArrayShuffle(nArray)[0]

			card1 := bjl.getCardByNumber(int32(n))
			bjl.delCard(card1)
			card2 := bjl.getCardByNumber(int32(n))
			bjl.delCard(card2)
			cards = append(cards, card1, card2)
			//閒有補牌需依點數判斷能不能補牌
			dian := bjl.getCardsDian(cards)
			if dian == 3 {
				if xiancards[2].number == 8 {
					hasThird = false
				}
			}
			if dian == 4 {
				if xiancards[2].number == 1 || xiancards[2].number == 8 ||
					xiancards[2].number == 9 || xiancards[2].number >= 10 {
					hasThird = false
				}
			}
			if dian == 5 {
				if xiancards[2].number == 1 || xiancards[2].number == 2 || xiancards[2].number == 3 ||
					xiancards[2].number == 8 ||	xiancards[2].number == 9 || xiancards[2].number >= 10 {
					hasThird = false
				}
			}
			if dian == 6 {
				if xiancards[2].number == 1 || xiancards[2].number == 2 || xiancards[2].number == 3 ||
					xiancards[2].number == 4 ||	xiancards[2].number == 5 ||	xiancards[2].number == 8 ||
					xiancards[2].number == 9 || xiancards[2].number >= 10 {
					hasThird = false
				}
			}
			if dian >= 7 {
				hasThird = false
			}

		} else {
			//閒沒補, 莊6點以上不能補牌
			for i:=1;i<=13;i++ {
				if isWin {
					if duiDianMap[i] > xiandian {
						nArray = append(nArray, i)
					}
				} else if bjl.zuhetype&(1<<2) != 0 {
					if duiDianMap[i] == xiandian || duiDianMap[i] < 6 {
						nArray = append(nArray, i)
					}
				} else {
					if duiDianMap[i] < xiandian {
						nArray = append(nArray, i)
					}
				}
			}
			if len(nArray) < 1 {
				return cards, false
			}
			n := util.IntArrayShuffle(nArray)[0]
			card1 := bjl.getCardByNumber(int32(n))
			bjl.delCard(card1)
			card2 := bjl.getCardByNumber(int32(n))
			bjl.delCard(card2)
			cards = append(cards, card1, card2)
			dian := bjl.getCardsDian(cards)
			if dian > 5 {
				hasThird = false
			}
		}
	} else {
		card1 := bjl.getCardbyIndex(0)
		nArray := make([]int, 0)
		for i:=int32(1);i<=13;i++ {
			if card1.number != i { //不能有對子
				tmpcard := bjl.getCardByNumber(i)
				if tmpcard.number == 0 {
					continue
				}
				dian := bjl.getCardsDian([]*card{card1, tmpcard})
				if len(xiancards) > 2 && dian > 7 { //閒有補牌點數不能是 8, 9
					continue
				}
				if isWin {
					if len(xiancards) > 2 { //閒有補牌
						if dian > xiandian && dian < 6{
							nArray = append(nArray, int(i))
						}
					} else {
						if dian > xiandian || dian < 6{
							nArray = append(nArray, int(i))
						}
					}
				} else if bjl.zuhetype&(1<<2) != 0 {
					if dian == xiandian {
						nArray = append(nArray, int(i))
					}
				} else {
					if dian < xiandian && dian < 6 {
						nArray = append(nArray, int(i))
					}
				}
			}
		}
		if len(nArray) < 1 {
			//mlog.Debug("==================%d , len %d, card1numer %d", xiandian, len(xiancards), card1.number)
			return cards, false
		}

		bjl.delCard(card1)
		n := util.IntArrayShuffle(nArray)[0]
		card2 := bjl.getCardByNumber(int32(n))
		bjl.delCard(card2)
		cards = append(cards, card1, card2)
		dian := bjl.getCardsDian(cards)
		if len(xiancards) > 2 {
			//閒有補牌需依點數判斷能不能補牌
			if dian == 3 {
				if xiancards[2].number == 8 {
					hasThird = false
				}
			}
			if dian == 4 {
				if xiancards[2].number == 1 || xiancards[2].number == 8 ||
					xiancards[2].number == 9 || xiancards[2].number >= 10 {
					hasThird = false
				}
			}
			if dian == 5 {
				if xiancards[2].number == 1 || xiancards[2].number == 2 || xiancards[2].number == 3 ||
					xiancards[2].number == 8 ||	xiancards[2].number == 9 || xiancards[2].number >= 10 {
					hasThird = false
				}
			}
			if dian == 6 {
				if xiancards[2].number == 1 || xiancards[2].number == 2 || xiancards[2].number == 3 ||
					xiancards[2].number == 4 ||	xiancards[2].number == 5 ||	xiancards[2].number == 8 ||
					xiancards[2].number == 9 || xiancards[2].number >= 10 {
					hasThird = false
				}
			}
			if dian >= 7 {
				hasThird = false
			}
		} else {
			dian := bjl.getCardsDian(cards)
			if dian > 5 {
				hasThird = false
			}
		}
	}
	if hasThird {
		var nArray []int
		for i:=1;i<=13;i++ {
			tmpcard := bjl.getCardByNumber(int32(i))
			if tmpcard.number == 0 {
				continue
			}
			dian := bjl.getCardsDian([]*card{cards[0], cards[1], tmpcard})
			if isWin {
				if dian > xiandian { //閒有補牌
					nArray = append(nArray, int(i))

				}
			} else if bjl.zuhetype&(1<<2) != 0 {
				if dian == xiandian {
					nArray = append(nArray, int(i))
				}
			} else {
				if dian < xiandian {
					nArray = append(nArray, int(i))
				}
			}
		}
		if len(nArray) < 1 {
			return cards, false
		}
		mlog.Debug("xiandian %d, %v, nArray %v", xiandian, cards, nArray)
		n := util.IntArrayShuffle(nArray)[0]
		card := bjl.getCardByNumber(int32(n))
		bjl.delCard(card)
		cards = append(cards, card)
	}
	mlog.Debug("xiandian %d, xiancards %s, zhuangcards %s", xiandian, xiancards, cards)
	//如果有沒取到牌就直接返回錯誤
	restoreFlag := false
	for _, card := range cards {
		if card.number == 0 {
			restoreFlag = true
		}
	}
	if restoreFlag {
		for _, card := range cards {
			if card.number != 0 {
				bjl.restoreCard(card)
			}
		}
		return cards, false
	}
	return cards, true

}