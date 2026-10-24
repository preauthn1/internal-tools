package mahjong

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
const (
	POKERFILENAME    = "mahjong"
	POKERSECTIONNAME = "usercards"
)

// FileCards 基於文件配置的牌型
//  用於特定牌型的發牌器
type FileCards struct {
	Cards     []int32
	cfg       *config.Config
	filecards []int32
}

// NewFileCards 創建發牌器對象
func NewFileCards(cfgs *config.Group) *FileCards {
	pcfg := cfgs.Configs[POKERFILENAME]

	fp := new(FileCards)
	fp.cfg = pcfg

	for i := 0; i < 100; i++ {
		pstr := pcfg.GetConfigValue(POKERSECTIONNAME, fmt.Sprintf("item%d", i))

		if pstr == "" {
			break
		}

		sarr := strings.Split(pstr, ",")

		for _, s := range sarr {
			if util.TrimAllSpace(s) == "" {
				continue
			}

			str := util.TrimAllSpace(s)
			cd, err := strconv.Atoi(str)
			if err != nil {
				return nil
			}
			if cd < 110 {
				mlog.Error("配置牌錯誤CD[%d]", cd)
				continue
			}

			fp.filecards = append(fp.filecards, int32(cd))
		}

	}
	fp.Cards = make([]int32, 0)
	fp.Cards = append(fp.Cards, fp.filecards...)
	//fp.Shuffle()

	return fp
}

// Shuffle ?
func (fp *FileCards) Shuffle() {
	//fp.Cards = make([]int32, 0)
	//fp.Cards = append(fp.Cards, fp.filecards...)
}

// ShuffleMore ?
func (fp *FileCards) ShuffleMore(times int) {
	fp.Shuffle()
}

// Deal ?
func (fp *FileCards) Deal(cn int) *CardGroup {
	if len(fp.Cards) < cn {
		return nil
	}

	ret := fp.Cards[:cn]
	fp.Cards = fp.Cards[cn:]

	return NewCardGroup(ret)
}

// GetCardGroup 獲取牌
func (fp *FileCards) GetCardGroup() *CardGroup {
	return NewCardGroup(fp.Cards)
}

// SetCardGroup 設置牌
func (fp *FileCards) SetCardGroup(cg *CardGroup) bool {
	if cg == nil {
		return false
	}
	if cg.Len() != fp.Len() {
		return false
	}
	fp.Cards = make([]int32, 0)
	fp.Cards = append(fp.Cards, cg.Cards...)
	return true
}

// AddCardGroup ?
func (fp *FileCards) AddCardGroup(cg *CardGroup) {}

// Len ?
func (fp *FileCards) Len() int {
	return len(fp.Cards)
}

func (fp *FileCards) String() string {
	return fp.cfg.String()
}
