package hhdz

import (
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"
	"testing"
)

func TestCardStorePretty(t *testing.T) { //從自然發牌中獲取比較好看的兩張牌
	findlimit := 1000
	powerCount := 0
	for i := 0; i < findlimit; i++ {
		rpoker := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		rpoker.Shuffle()
		helper := newPokerHelper(rpoker)
		helper.initStorePokerType(NATURE, NATURE)
		store1 := helper.hungStore
		store2 := helper.heiStore
		if util.RndInt(0, 2) == 0 {
			helper.getPrettyFirstCards()
		} else {
			helper.getRandomFirstCards()
		}
		if len(store1.store.Cards) > 1 || len(store2.store.Cards) > 1 {
			t.Fatalf("取了兩張牌、手牌應該只剩一張 %s %s", store1.store, helper.heiStore.store)
		}
		winuserid := util.RndInt(1, 3) // 1-2
		helper.getPokerCardByWinUserID(int32(winuserid))

		cg1 := poker.NewCardGroup(helper.hungStore.cards)
		cg2 := poker.NewCardGroup(helper.heiStore.cards)
		//t.Logf("winuserid %d %s %s", winuserid, store1.cards, store2.cards)
		if winuserid == 1 && comPareCards(cg1, cg2) != 1 {
			t.Fatalf("winuserid %d 不符合預期 %s %s", winuserid, cg1, cg2)
		} else if winuserid == 2 && comPareCards(cg2, cg1) != 1 {
			t.Fatalf("winuserid %d 不符合預期 %s %s", winuserid, cg1, cg2)
		}
		_, _, zuhetype := helper.getZuheTypes()
		if zuhetype != DANPAI && zuhetype != SMALLDUIZI {
			powerCount++
			//t.Logf("zuhetype %d, %s %s", helper.getZuheType(), cg1, cg2)
		}
	}
	t.Logf("致命一擊次數, %d", powerCount)
}

func TestCardStore(t *testing.T) { //獲取牌型測試
	findlimit := 1000

	for i := 0; i < findlimit; i++ {
		rpoker := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		rpoker.Shuffle()
		helper := newPokerHelper(rpoker)
		helper.initStorePokerType(BAOZI, BAOZI)
		store1 := helper.hungStore
		store2 := helper.heiStore
		//t.Logf("%s, %s", store1.store, store2.store)
		if !isBaozi(store1.store) || !isBaozi(store2.store){
			t.Fatalf("無法組出順金牌 %s, %s", store1.store, store2.store)
		}
	}

	for i := 0; i < findlimit; i++ {
		rpoker := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		rpoker.Shuffle()
		helper := newPokerHelper(rpoker)
		helper.initStorePokerType(SHUNJIN, SHUNJIN)
		store1 := helper.hungStore
		store2 := helper.heiStore
		//t.Logf("%s, %s", store1.store, store2.store)
		if !isSunZi(store1.store) || !isJinHua(store1.store) ||
			!isSunZi(store2.store) || !isJinHua(store2.store) {
			t.Fatalf("無法組出順金牌 %s, %s", store1.store, store2.store)
		}
	}

	for i := 0; i < findlimit; i++ {
		rpoker := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		rpoker.Shuffle()
		helper := newPokerHelper(rpoker)
		helper.initStorePokerType(JINHUA, JINHUA)
		store1 := helper.hungStore
		store2 := helper.heiStore
		if !isJinHua(store1.store) || !isJinHua(store2.store) {
			t.Fatalf("無法組出金花牌 %s, %s", store1.store, store2.store)
		}
		if isSunZi(store1.store) || isSunZi(store2.store) {
			t.Fatalf("金花組出順金牌 %s, %s", store1.store, store2.store)
		}
	}

	for i := 0; i < findlimit; i++ {
		rpoker := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		rpoker.Shuffle()
		helper := newPokerHelper(rpoker)
		helper.initStorePokerType(SHUNZI, SHUNZI)
		store1 := helper.hungStore
		store2 := helper.heiStore
		//t.Logf("%s, %s", store1.store, store2.store)
		if !isSunZi(store1.store) || !isSunZi(store2.store) {
			t.Fatalf("無法組出順子牌 %s, %s", store1.store, store2.store)
		}
		if isJinHua(store1.store) || isJinHua(store2.store) {
			t.Fatalf("順子組出順金牌 %s, %s", store1.store, store2.store)
		}
	}

	for i := 0; i < findlimit; i++ {
		rpoker := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		rpoker.Shuffle()
		helper := newPokerHelper(rpoker)
		helper.initStorePokerType(BIGDUIZI, BIGDUIZI)
		store1 := helper.hungStore
		store2 := helper.heiStore
		//t.Logf("%s, %s", store1.store, store2.store)
		if !isDuiZi(store1.store) || !isDuiZi(store2.store) {
			t.Fatalf("無法組出對子牌 %s, %s", store1.store, store2.store)
		}
		if isBaozi(store1.store) || isBaozi(store2.store) {
			t.Fatalf("對子組出豹子牌 %s, %s", store1.store, store2.store)
		}
	}

	for i := 0; i < findlimit; i++ {
		rpoker := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		rpoker.Shuffle()
		helper := newPokerHelper(rpoker)
		helper.initStorePokerType(SMALLDUIZI, SMALLDUIZI)
		store1 := helper.hungStore
		store2 := helper.heiStore
		//t.Logf("%s, %s", store1.store, store2.store)
		if !isDuiZi(store1.store) || !isDuiZi(store2.store) {
			t.Fatalf("無法組出對子牌 %s, %s", store1.store, store2.store)
		}
		if isBaozi(store1.store) || isBaozi(store2.store) {
			t.Fatalf("對子組出豹子牌 %s, %s", store1.store, store2.store)
		}
	}

	for i := 0; i < findlimit; i++ {
		rpoker := poker.NewRandomPoker(false, false, 1) //初始化一副牌
		rpoker.Shuffle()
		helper := newPokerHelper(rpoker)
		helper.initStorePokerType(DANPAI, DANPAI)
		store1 := helper.hungStore
		store2 := helper.heiStore
		//t.Logf("%s, %s", store1.store, store2.store)
		if GetType(store1.store) != DanPai || GetType(store2.store) != DanPai {
			t.Fatalf("無法組出單牌 %s, %s", store1.store, store2.store)
		}

	}
}
