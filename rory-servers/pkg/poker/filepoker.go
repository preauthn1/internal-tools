package poker

import (
	"os"
	"path"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
const (
	POKERFILENAME    = "poker"
	POKERSECTIONNAME = "userpoker"
)

// FilePoker 基於文件配置的牌型
// 用於特定牌型的發牌器
type FilePoker struct {
	cfg       *config.Config
	userCards map[int32]*CardGroup
}

// NewFilePoker 創建發牌器對象
func NewFilePoker(cfgs *config.Group, urs []*gamelogic.GameUser) *FilePoker {
	pokerFileName := POKERFILENAME
	fp := innerNewFilePoker(cfgs, pokerFileName)

	return fp
}

// NewTestFilePoker ?
func NewTestFilePoker(testFilePoker string) *FilePoker {
	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "testfixpoker")
	cfgs := new(config.Group)
	_ = cfgs.LoadConfig(configDir, 0)
	pokerFileName := testFilePoker
	fp := innerNewFilePoker(cfgs, pokerFileName)

	return fp
}

func innerNewFilePoker(cfgs *config.Group, pokerFileName string) *FilePoker {
	pcfg := cfgs.Configs[pokerFileName]
	fp := new(FilePoker)
	fp.cfg = pcfg
	fp.userCards = make(map[int32]*CardGroup)
	sec := pcfg.GetConfigSection("userpoker")
	for uidxs, pstr := range sec {
		if uidxs == "" || pstr == "" {
			continue
		}

		uidx, err := strconv.Atoi(uidxs)
		if err != nil {
			mlog.Error("解析用戶索引錯誤 %v, %s", err, uidxs)
			continue
		}
		//用户转换
		//var uid int32 = 0

		uid := int32(uidx)

		//创建用户牌型
		cg := NewCardGroup(nil)
		fp.userCards[uid] = cg
		//拆解牌型, 雙管線後面有特別的用途, 可以用於測試驗證牌型的結果
		tmpArr := strings.Split(pstr, "||")
		//解析牌型
		sarr := strings.Split(tmpArr[0], "|")

		for _, s := range sarr {
			if util.TrimAllSpace(s) == "" {
				continue
			}

			pk := NewPokerCardByStr(s, ",")
			if pk == nil {
				mlog.Error("解析牌型错误,字符串%s格式不正确.%s", s, pstr)
				continue
			}
			fp.userCards[uid].AddCard(pk)
		}
	}
	return fp
}

// GetUserCards ?
func (fp *FilePoker) GetUserCards() map[int32]*CardGroup {
	return fp.userCards
}

// GetConfig ?
func (fp *FilePoker) GetConfig() *config.Config {
	return fp.cfg
}

// Shuffle ?
func (fp *FilePoker) Shuffle() {
}

// UserCardsLen ?
func (fp *FilePoker) UserCardsLen(uid int32) int {
	return fp.userCards[uid].Len()
}

// ShuffleMore ?
func (fp *FilePoker) ShuffleMore(times int) {
	fp.Shuffle()
}

// Deal ?
func (fp *FilePoker) Deal(userid int32, cn int) *CardGroup {
	cuid := RestCardsUserID

	if _, ok := fp.userCards[userid]; ok {
		cuid = userid
	}

	cards := fp.userCards[cuid].Cards[:cn]
	fp.userCards[userid].Cards = fp.userCards[cuid].Cards[cn:]

	return NewCardGroup(cards)
}

//ExChange 換牌
func (fp *FilePoker) ExChange(userid int32, otheruserid int32) {

}

func (fp *FilePoker) String() string {
	return fp.cfg.String()
}

// TakeCardNumber ?
func (fp *FilePoker) TakeCardNumber(userid int32, cardNumber int32) (bool, Card) {
	for k, v := range fp.userCards[userid].Cards {
		if v.CardNumber == cardNumber {
			fp.userCards[userid].Cards = RemoveIndex(fp.userCards[userid].Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// TakeDiffCardNumber 指定取非存在的號碼
func (fp *FilePoker) TakeDiffCardNumber(userid int32, diffNums []int32) (bool, Card) {
	for k, v := range fp.userCards[userid].Cards {
		right := true
		for _, diffnum := range diffNums {
			if v.CardNumber == diffnum {
				right = false
			}
		}
		if right {
			fp.userCards[userid].Cards = RemoveIndex(fp.userCards[userid].Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// TakeCardColor ?
func (fp *FilePoker) TakeCardColor(userid int32, color int32) (bool, Card) {
	for k, v := range fp.userCards[userid].Cards {
		if v.Color == color {
			fp.userCards[userid].Cards = RemoveIndex(fp.userCards[userid].Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// TakeCardNumberAndColor ?
func (fp *FilePoker) TakeCardNumberAndColor(userid int32, cardNumber int32, color int32) (bool, Card) {
	for k, v := range fp.userCards[userid].Cards {
		if v.CardNumber == cardNumber && v.Color == color {
			fp.userCards[userid].Cards = RemoveIndex(fp.userCards[userid].Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// AddCardGroup ?
func (fp *FilePoker) AddCardGroup(userid int32, cg *CardGroup) {
	if len(cg.Cards) <= int(0) {
		return
	}

	fp.userCards[userid].Cards = append(fp.userCards[userid].Cards, cg.Cards...)
}

// AddCard ?
func (fp *FilePoker) AddCard(userid int32, card *Card) {
	if card == nil {
		return
	}
	fp.userCards[userid].Cards = append(fp.userCards[userid].Cards, card)
}
