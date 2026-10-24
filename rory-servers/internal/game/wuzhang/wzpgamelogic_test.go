package wzp

import (
	"encoding/csv"
	"fmt"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"

	"os"
	"path"
	"strconv"
	"strings"
	"testing"
)

func initConfig() {
	betamounts := "200,300,500,800,1500"
	gameconfig = new(wzpGameConfig)
	gameconfig.baseAmount = 100
	gameconfig.betTimeOut = 8
	betAmountArr := strings.Split(betamounts, ",")
	for _, ba := range betAmountArr {
		iba, _ := strconv.Atoi(ba)
		gameconfig.betAmount = append(gameconfig.betAmount, int32(iba))
	}

}

//隨機下注
func Test_getBetAmount(t *testing.T) {
	conf.CLogLevel = mlog.FatalLevel //disable_log
	initConfig()
	lgc := new(WZPRobotLogic)
	betScale := []int32{0, 2, 2, 0, 0}
	chaitem := NewRobotChaItem(true, true, 0, betScale, false, 0, false)

	for i := 0; i < 100; i++ {
		val := lgc.getBetAmount(200, chaitem)
		if int32(val) != gameconfig.betAmount[1] && int32(val) != gameconfig.betAmount[2] {
			t.Error("Robot 隨機下注 failure") // 如果不是如预期的那么就报错
		}
	}
	t.Log("#1 隨機下注 success")
}

//隨機發捭, 梭哈
func TestFapai(t *testing.T) {
	var p poker.Poker
	user1 := gamelogic.NewGameUser()
	user1.UserID = 1
	allTableUsers := []*gamelogic.GameUser{user1}
	p = poker.NewRandomPoker(false, false, 1)
	//洗牌分牌
	p = poker.NewFixedPoker(allTableUsers, p.(*poker.RandomPoker), 5)
	fp := p.(*poker.FixedPoker)

	cards := fp.GetUserCards(1)
	t.Log("初始化後用戶將會拿到的牌", cards)
	origin := strings.Split(strings.TrimRight(cards.String(), ","), ",")

	foo := p.Deal(1, 2)
	t.Log("獲取兩張牌", foo)
	get2 := strings.Split(strings.TrimRight(foo.String(), ","), ",")

	cards = fp.GetUserCards(1)
	t.Log("用戶剩下的牌", cards)
	other := strings.Split(strings.TrimRight(cards.String(), ","), ",")

	for i, v := range get2 {
		if v != origin[i] {
			t.Errorf(fmt.Sprintf("比對原始牌與獲取的兩張有所不同, %v : 獲取:%v!=原始:%v", i, v, origin[i]))
		}
	}

	for i, v := range other {
		if v != origin[i+2] {
			t.Errorf("取兩張牌後, 比對原始牌與剩餘牌有所不同, %v, %v", origin, other)
		}
	}

	t.Log("#2 隨機發牌 success")
}

//固定發牌
func TestFapai2(t *testing.T) {
	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	//t.Log(cfgs)
	// 需要三個 user
	var allTableUsers []*gamelogic.GameUser
	for i := 0; i < 3; i++ {
		user := gamelogic.NewGameUser()
		user.UserID = int32(i)
		allTableUsers = append(allTableUsers, user)
	}

	var p poker.Poker
	p = poker.NewFilePoker(cfgs, allTableUsers)
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	valid_cards := []string{"红桃A,红桃8,红桃Q,红桃J,红桃10,", "黑桃A,黑桃7,黑桃Q,黑桃J,黑桃10,", "梅花A,梅花6,梅花Q,梅花J,梅花10,"}
	for i := 0; i < 3; i++ {
		cards := fp.GetUserCards(int32(i))
		if cards.String() != valid_cards[i] {
			t.Errorf("固定發牌結果與設定檔有所不同, 實際發牌:%v, 發牌設定:%v", cards, valid_cards[i])
		}
		t.Log("user ", i, " ", cards)
	}

	t.Log("#3 固定發牌 success")
}

func writeDataToCsv(data [][]string) {
	f, err := os.Create("Fapai_10000.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)      //写入数据
	w.Flush()
}

type WZPCards []*poker.Card

func (p WZPCards) Len() int {
	return len(p)
}

func (p WZPCards) Less(i, j int) bool {
	cardi, _ := strconv.Atoi(fmt.Sprintf("%d%d", p[i].CardNumber, p[i].Color))
	cardj, _ := strconv.Atoi(fmt.Sprintf("%d%d", p[j].CardNumber, p[j].Color))
	// fmt.Println(cardi, cardj)
	return cardi < cardj
}

func (p WZPCards) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

//梭哈 发牌10000副 比大小结果
func GenCSVFapai10000(t *testing.T) {
	//同花大顺 1 同花顺 2 四条 3 葫芦 4 同花 5 顺子 6 三条 7 两对 8 一对 9 散牌 10 弃牌 11
	typeMap := map[int32]string{
		1:  "同花大顺",
		2:  "同花顺",
		3:  "四条",
		4:  "葫芦",
		5:  "同花",
		6:  "顺子",
		7:  "三条",
		8:  "两对",
		9:  "一对",
		10: "散牌",
		11: "弃牌",
	}

	// 設定csv title
	records := [][]string{
		{"玩家1", "玩家2", "1牌型", "2牌型", "結果"},
	}

	user1 := gamelogic.NewGameUser()
	user2 := gamelogic.NewGameUser()
	user1.UserID = 1
	user2.UserID = 2
	allTableUsers := []*gamelogic.GameUser{user1, user2}

	for i := 0; i < 10000; i++ {
		p := poker.NewRandomPoker(false, false, 1)
		//洗牌分牌
		p1 := poker.NewFixedPoker(allTableUsers, p, 5)
		p2 := poker.NewFixedPoker(allTableUsers, p, 5)

		cards1 := p1.GetUserCards(1)
		c1 := make(WZPCards, len(cards1.Cards))
		copy(c1, cards1.Cards)
		sort.Sort(c1)
		t.Log(c1)
		cards2 := p2.GetUserCards(2)
		c2 := make(WZPCards, len(cards2.Cards))
		copy(c2, cards2.Cards)
		sort.Sort(c2)
		// sort.Sort(cards2)
		// t.Log("玩家1:", cards1)
		// t.Log("玩家2:", cards2)
		c1Str := ""
		for _, v := range c1 {
			c1Str += fmt.Sprintf("%v, ", v.String())
		}
		c2Str := ""
		for _, v := range c2 {
			c2Str += fmt.Sprintf("%v, ", v.String())
		}

		ret, type1, type2 := compareCard(cards1, cards2)

		if ret {
			// t.Log("玩家1贏", "玩家1:", typeMap[type1], "玩家2:", typeMap[type2])

			records = append(records, []string{fmt.Sprintf("%v", c1Str), fmt.Sprintf("%v", c2Str), fmt.Sprintf("%v", typeMap[type1]), fmt.Sprintf("%v", typeMap[type2]), "1贏"})
		} else {
			// t.Log("玩家2贏", "玩家1:", typeMap[type1], "玩家2:", typeMap[type2])

			records = append(records, []string{fmt.Sprintf("%v", c1Str), fmt.Sprintf("%v", c2Str), fmt.Sprintf("%v", typeMap[type1]), fmt.Sprintf("%v", typeMap[type2]), "2贏"})
		}

	}

	writeDataToCsv(records)
	// t.Error("AAA")
}

//同花大顺 1 同花顺 2 四条 3 葫芦 4 同花 5 顺子 6 三条 7 两对 8 一对 9 散牌 10 弃牌 11
func TestRandomFapai(t *testing.T) {
	var allTableUsers []*gamelogic.GameUser
	for i := 0; i < 3; i++ {
		user := gamelogic.NewGameUser()
		user.UserID = int32(i)
		allTableUsers = append(allTableUsers, user)
	}

	//結果統計
	result := make(map[int]int)
	var p poker.Poker
	for step := 0; step < 100; step++ {
		p = poker.NewRandomPoker(false, false, 1)
		//洗牌分牌
		p = poker.NewFixedPoker(allTableUsers, p.(*poker.RandomPoker), 5)

		fp := p.(*poker.FixedPoker)
		for i := 0; i < 3; i++ {
			cards := fp.GetUserCards(int32(i))
			size := getCardSize(cards)
			result[int(size)] += 1
		}
	}
	if len(result) < 5 {
		t.Errorf("隨機發牌產生牌型少於5種, %v", result)
	}
	t.Log("#4 隨機發牌統計牌型 success")
}
