package dchjk

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	Pokerset "example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

var cfg *config.Group

func init() {
	conf.CLogLevel = mlog.DebugLevel
	conf.FLogLevel = 0

	// 檢查config位置
	if err := initCfg(); err != nil {
		fmt.Fprintf(os.Stderr, "Load config files error. %v\n", err)
		os.Exit(1)
	}
}

func initCfg() error {
	// 檢查config位置
	var data *config.Group
	_, file, _, _ := runtime.Caller(0)
	var defaultPaths = []string{"./config", "../../config", path.Dir(file) + "/config", path.Dir(file) + "/../../config"}

	data = new(config.Group)

	for _, path := range defaultPaths {
		fi, err := os.Stat(path)
		if err != nil {
			mlog.Debug("Find config folder %s Error", path)
			continue
		}

		if !fi.IsDir() {
			mlog.Debug("%s is not a folder", path)
			continue
		}

		data.LoadConfig(path, 0)

		if len(data.Configs) > 0 {
			mlog.Debug("Find config folder %s Success", path)
			break
		}
		mlog.Debug("Find config folder %s Error", path)
	}

	if data == nil {
		return fmt.Errorf("Can not find config folder, please put setting files on this folder")
	}

	if len(data.Configs) == 0 {
		return fmt.Errorf("Can not load configs, path: %s", data.DirName)
	}

	cfg = data

	return nil
}

// 莊家與玩家比牌
func TestCompareCards(t *testing.T) {
	for i := 0; i < 20; i++ {
		pPokerSet := Pokerset.NewPokerCardSet()
		pPokerSet.NewSet(JD21_PokerSetCount, 0)
		pPokerSet.Shuffle()

		pPlayerSet := NewCardSetInfo(JD21_IndexBase)
		pPset := pPokerSet.Deal(2)
		pPlayerSet.AddCardSet(pPset)
		if !pPlayerSet.bIsBlackJack && !pPlayerSet.bIsBusted && !pPlayerSet.IsSoftReached() {
			for {
				pCardSet := pPokerSet.Deal(1)
				pPlayerSet.AddCardSet(pCardSet)
				if pPlayerSet.IsSoftReached() {
					break
				}
			}
		}

		pDealerSet := NewCardSetInfo(JD21_IndexBase)
		pDset := pPokerSet.Deal(2)
		pDealerSet.AddCardSet(pDset)
		if !pDealerSet.bIsBlackJack && !pDealerSet.bIsBusted && !pDealerSet.IsSoftReached() {
			for {
				pCardSet := pPokerSet.Deal(1)
				pDealerSet.AddCardSet(pCardSet)
				if pDealerSet.IsSoftReached() {
					break
				}
			}
		}

		map_strResult := map[int]string{
			JD21_Compare_Even:   "平局",
			JD21_Compare_Dealer: "莊家贏",
			JD21_Compare_Player: "玩家贏",
		}

		mlog.Debug("莊家牌：%s (%d, %d), BJ[%t}, Bust[%t]", pDealerSet.pCardSet.CodeString(), pDealerSet.i32Point_1, pDealerSet.i32Point_2,
			pDealerSet.bIsBlackJack, pDealerSet.bIsBusted)
		mlog.Debug("玩家牌：%s (%d, %d), BJ[%t}, Bust[%t]", pPlayerSet.pCardSet.CodeString(), pPlayerSet.i32Point_1, pPlayerSet.i32Point_2,
			pPlayerSet.bIsBlackJack, pPlayerSet.bIsBusted)
		iResult := CompareCards(pDealerSet, pPlayerSet)
		mlog.Debug("比牌結果：%s", map_strResult[iResult])
		mlog.Debug("----------------------------------------------------------------------------------")
	}
}
