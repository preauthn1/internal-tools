package sss

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker"
)

func TestSearchTable(t *testing.T) {
	// var ht = new(HandsTable)
	// ht.Init()
}

func TestSSSPokerType(t *testing.T) {
	var s = new(SSSTable)
	var deckTree = new(SSSDeckTree)
	cg := poker.LoadCardGroupFromString("1,2|2,2|3,3|2,4|2,5|3,7|2,8|3,8|3,9|1,11|0,12|0,13|2,13")
	deckTree.NewDeckTree(s, cg, 4) // 最多4片葉子
	deckTree.GrowUp(s)
	deckTree.Print()
	cg = poker.LoadCardGroupFromString("0,2|2,2|2,3|1,4|3,4|0,5|1,6|3,7|1,8|0,12|2,13|0,14|1,14")
	deckTree.NewDeckTree(s, cg, 4) // 最多4片葉子
	deckTree.GrowUp(s)
	deckTree.Print()

}

// GenSSSDeckSuggest
func GenSSSDeckSuggest(t *testing.T) {
	// 設定csv title
	records := [][]string{
		{"牌組id", "推薦id", "原始牌組", "頭墩牌型", "頭墩牌組", "中墩牌型", "中墩牌組", "尾墩牌型", "尾墩牌組"},
	}
	var cardTypeMap = map[int32]string{
		1:  "同花順",
		2:  "炸彈",
		3:  "葫蘆",
		4:  "同花",
		5:  "順子",
		6:  "三條",
		7:  "衝三",
		8:  "兩對",
		9:  "一對",
		10: "散牌",
	}
	count := 0
	for i := 0; i < 1; i++ {
		var s = new(SSSTable)
		// s.gamePoker = poker.NewGamePoker()
		// s.poker = poker.NewRandomPoker(false, false, 1)
		// s.poker.Shuffle()
		// cg := s.poker.Deal(0, 13)

		// 红桃3,黑桃3,黑桃4,梅花5,黑桃5,黑桃6,黑桃8,梅花9,方块10,梅花10,梅花Q,梅花K,红桃K,
		cards := make([]*poker.Card, 0)
		cards = append(cards, &poker.Card{CardNumber: 3, Color: poker.PokerColorHongTao})
		cards = append(cards, &poker.Card{CardNumber: 3, Color: poker.PokerColorHeiTao})
		cards = append(cards, &poker.Card{CardNumber: 4, Color: poker.PokerColorHeiTao})
		cards = append(cards, &poker.Card{CardNumber: 5, Color: poker.PokerColorMeiHua})
		cards = append(cards, &poker.Card{CardNumber: 5, Color: poker.PokerColorHeiTao})
		cards = append(cards, &poker.Card{CardNumber: 6, Color: poker.PokerColorHeiTao})
		cards = append(cards, &poker.Card{CardNumber: 8, Color: poker.PokerColorHeiTao})
		cards = append(cards, &poker.Card{CardNumber: 9, Color: poker.PokerColorMeiHua})
		cards = append(cards, &poker.Card{CardNumber: 10, Color: poker.PokerColorFangKuai})
		cards = append(cards, &poker.Card{CardNumber: 10, Color: poker.PokerColorMeiHua})
		cards = append(cards, &poker.Card{CardNumber: 12, Color: poker.PokerColorMeiHua})
		cards = append(cards, &poker.Card{CardNumber: 13, Color: poker.PokerColorMeiHua})
		cards = append(cards, &poker.Card{CardNumber: 13, Color: poker.PokerColorHongTao})
		cg := poker.NewCardGroup(cards)

		// 红桃A,梅花K,方块Q,方块J,梅花10,黑桃8,红桃7,方块6,方块5,梅花4,方块4,红桃3,黑桃2,
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 14, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 13, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 11, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 10, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 8, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 7, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 6, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 3, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_MeiHua})
		// cg := poker.NewCardGroup(cards)

		// 方块A,黑桃K,红桃K,黑桃Q,方块J,红桃9,方块7,黑桃6,红桃6,梅花6,方块6,梅花4,梅花2,
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 14, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 13, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 13, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 11, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 7, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 6, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 6, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 6, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 6, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_MeiHua})
		// cg := poker.NewCardGroup(cards)

		// 梅花Q,黑桃J,梅花J,梅花10,黑桃9,红桃9,梅花9,梅花8,方块8,黑桃7,方块5,黑桃2,红桃2,
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 11, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 11, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 10, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 8, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 8, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 7, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_HongTao})
		// cg := poker.NewCardGroup(cards)

		// 黑桃A,红桃K,方块K,黑桃Q,方块Q,红桃J,方块9,黑桃8,红桃8,黑桃7,方块5,梅花4,梅花2,
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 14, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 13, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 13, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 11, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 8, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 8, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 7, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_MeiHua})
		// cg := poker.NewCardGroup(cards)

		// 红桃A,黑桃K,方块K,黑桃Q,红桃10,梅花10,黑桃9,红桃9,黑桃8,梅花5,梅花4,方块4,红桃2,
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 14, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 13, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 13, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 10, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 10, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 8, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_HeiTao})
		// cg := poker.NewCardGroup(cards)

		// 方块Q,红桃J,方块J,红桃10,黑桃9,梅花9,黑桃6,方块6,黑桃5,梅花5,梅花4,方块4,梅花3,
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 11, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 11, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 10, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 6, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 6, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 3, Color: poker.PokerColor_MeiHua})
		// cg := poker.NewCardGroup(cards)

		// 黑桃A,红桃A,方块A,黑桃Q,黑桃10,黑桃9,红桃9,梅花9,红桃8,方块5,红桃3,梅花2,方块2,
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 14, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 14, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 14, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 10, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 8, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 3, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_FangKuai})
		// cg := poker.NewCardGroup(cards)

		// 红桃A,梅花Q,梅花J,黑桃9,红桃8,红桃6,黑桃5,梅花5,方块5,黑桃4,红桃4,方块4,红桃3
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 14, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 11, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 8, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 6, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 3, Color: poker.PokerColor_HongTao})
		// cg := poker.NewCardGroup(cards)

		// 红桃Q,方块Q,红桃10,红桃9,方块7,红桃5,梅花5,方块4,黑桃3,黑桃2,红桃2,梅花2,方块2,
		// cards := make([]*poker.PokerCard, 0)
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 12, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 10, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 9, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 7, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 5, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 4, Color: poker.PokerColor_FangKuai})
		// cards = append(cards, &poker.PokerCard{CardNumber: 3, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_HeiTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_HongTao})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_MeiHua})
		// cards = append(cards, &poker.PokerCard{CardNumber: 2, Color: poker.PokerColor_FangKuai})
		// cg := poker.NewCardGroup(cards)

		var deckTree = new(SSSDeckTree)
		deckTree.NewDeckTree(s, cg, 4) // 最多4片葉子
		deckTree.GrowUp(s)
		deckTree.Print()
		// assert.Equal(t, len(deckTree.GetSortNodes()), 5, "不等於5")
		// if len(deckTree.GetSortNodes()) != 5 {
		// 	deckTree.Print()
		// }
		count++
		for _, v := range deckTree.GetSortNodes() {
			records = append(records, []string{fmt.Sprintf("%d", count), fmt.Sprintf("%d", v.position+1), fmt.Sprintf("%v", deckTree.source),
				fmt.Sprintf("%v", cardTypeMap[v.data.GetHeadCardType()]), fmt.Sprintf("%v", v.data.head),
				fmt.Sprintf("%v", cardTypeMap[v.data.GetBodyCardType()]), fmt.Sprintf("%v", v.data.body),
				fmt.Sprintf("%v", cardTypeMap[v.data.GetFootCardType()]), fmt.Sprintf("%v", v.data.foot)})
		}
	}
	// writeDataToCsv(records)
}

func writeDataToCsv(data [][]string) {
	f, err := os.Create("Suggest_1000.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)      //写入数据
	w.Flush()
}
