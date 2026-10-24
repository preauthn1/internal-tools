package hbsl

import (
	"example.invalid/ServerTeam/Servers/pkg/util"
	"math/rand"
	"sort"
)

// GameDataDB 遊戲輸贏結果 DB 儲存資料
type GameDataDB struct {
}

// Copy 複製遊戲輸贏結果 DB 儲存資料
func (pData *GameDataDB) Copy(pOther *GameDataDB) {
	if pOther == nil {
		return
	}
}

// LogString log 用字串
func (pData *GameDataDB) LogString() string {
	strLog := ""

	return strLog
}

//---------------------------------------------------------------------------

// SysPoolWinLose 遊戲輸贏結果
type SysPoolWinLose struct {
	i64PoolAmount     int64      // 目前水線
	i64LastPoolAmount int64      // 上一局結算後水線
	GameDBData        GameDataDB // 遊戲輸贏結果 DB 儲存資料
}

// Copy 複製遊戲輸贏結果
func (pWinLose *SysPoolWinLose) Copy(pOther *SysPoolWinLose) {
	if pOther == nil {
		return
	}

	pWinLose.i64PoolAmount = pOther.i64PoolAmount
	pWinLose.i64LastPoolAmount = pOther.i64LastPoolAmount
	pWinLose.GameDBData.Copy(&pOther.GameDBData)
}

//---------------------------------------------------------------------------
var gbEnablePoolControl = false				// 是否開啟水線控制功能
var gbGetSysPoolConfig = false				// 是否已取得水池 config
var gi64WinLine = int64(0)					// 贏錢線
var giBombRate = 0							// 作弊機率
var gbGetWinLoseData = false				// 是否已自 DB 取得遊戲輸贏統計資料
var gbCheatMode = false						// 是否進入作弊模式

var gpSysPoolWinLose *SysPoolWinLose		// 房間總輸贏結果

//---------------------------------------------------------------------------

// CheatHongBao 產生作弊模式紅包
// aryi64OriginalHB：原來的紅包
// 回傳［新的紅包資料］
func (glc *HbslLogic) CheatHongBao(aryi64OriginalHB []int64) []int64 {
	if !gbCheatMode {
		return aryi64OriginalHB
	}

	if int32(len(aryi64OriginalHB)) != glc.i32HBNum {
		return aryi64OriginalHB
	}

	aryi64RetHB := make([]int64, 0)

	// 檢查紅包內雷包數量
	iBombCount := 0
	for _, i64Value := range aryi64OriginalHB {
		if i64Value%10 == int64(glc.pNowHB.pHB.i32BombNumber) {
			iBombCount++
		}
	}

	if glc.pNowHB.pHB.i32UserID == -1 || glc.pNowHB.pHB.pUser.pGameUser.IsRobot() {
		// 系統或機器人發的紅包
		if iBombCount > 0 {
			return aryi64OriginalHB
		}

		aryi64RetHB = glc.SetBombHongBao(true)
	} else {
		// 玩家發的紅包
		if iBombCount <= 0 {
			return aryi64OriginalHB
		}

		aryi64RetHB = glc.SetBombHongBao(false)
	}

	return aryi64RetHB
}

// SetBombHongBao 設定雷包紅包
// bHasBomb：紅包是否要含雷包
// 回傳［新的紅包］
func (glc *HbslLogic) SetBombHongBao(bHasBomb bool) []int64 {
	aryi64RetHB := make([]int64, 0)

	i64HBMoney := glc.pNowHB.pHB.i64Chip
	i32HBNum := glc.i32HBNum

	// 扣稅
	_, tax := glc.GetUserWinAndTax(i64HBMoney)
	i64HBMoney -= tax

	for {
		aryi64TempHB := make([]int64, 0)
		aryi64Rand := make([]int64, 0)
		aryi64Rand = append(aryi64Rand, 0)
		for i := int32(0); i < i32HBNum-1; i++ {
			for {
				i64Rnd := rand.Int63n(i64HBMoney-1) + 1
				if exist, _ := util.InArray(i64Rnd, aryi64Rand); !exist {
					aryi64Rand = append(aryi64Rand, i64Rnd)
					break
				}
			}
		}

		sort.Slice(aryi64Rand, func(i, j int) bool {
			return aryi64Rand[i] > aryi64Rand[j]
		})

		i64Money, iBombCount := i64HBMoney, 0
		for _, i64Rnd := range aryi64Rand {
			i64Value := i64Money - i64Rnd
			if i64Value%10 == int64(glc.pNowHB.pHB.i32BombNumber) {
				iBombCount++
			}
			aryi64TempHB = append(aryi64TempHB, i64Value)
			i64Money = i64Rnd
		}

		if bHasBomb {
			if iBombCount > 0 {
				aryi64RetHB = append(aryi64RetHB, aryi64TempHB...)
				break
			}
		} else {
			if iBombCount <= 0 {
				aryi64RetHB = append(aryi64RetHB, aryi64TempHB...)
				break
			}
		}
	}

	return aryi64RetHB
}

// GrabCheatHongBao 搶作弊紅包
func (glc *HbslLogic) GrabCheatHongBao(pUser *User) int64 {
	if len(glc.pNowHB.ary_i64SHB) <= 0 {
		return 0
	}

	if pUser == nil {
		return 0
	}

	if !gbCheatMode {
		return glc.pNowHB.ary_i64SHB[0]
	}

	i64RetMoney := int64(0)

	bGet, iIdx := false, -1
	if glc.pNowHB.pHB.i32UserID == -1 || glc.pNowHB.pHB.pUser.pGameUser.IsRobot() {
		// 系統或機器人發的紅包
		for i := 0; i < len(glc.pNowHB.ary_i64SHB); i++ {
			bIsBomb := false
			if glc.pNowHB.ary_i64SHB[i]%10 == int64(glc.pNowHB.pHB.i32BombNumber) {
				bIsBomb = true
			}

			if !pUser.pGameUser.IsRobot() {
				// 玩家應拿到雷包
				if bIsBomb {
					bGet = true
					iIdx = i
					break
				}
			} else {
				// 機器人不應拿到雷包
				if bIsBomb {
					continue
				}

				bGet = true
				iIdx = i
				break
			}
		}
	}

	if !bGet {
		i64RetMoney = glc.pNowHB.ary_i64SHB[0]
		glc.pNowHB.ary_i64SHB = append(glc.pNowHB.ary_i64SHB[:0], glc.pNowHB.ary_i64SHB[1:]...)
	} else {
		i64RetMoney = glc.pNowHB.ary_i64SHB[iIdx]
		glc.pNowHB.ary_i64SHB = append(glc.pNowHB.ary_i64SHB[:iIdx], glc.pNowHB.ary_i64SHB[iIdx+1:]...)
	}

	return i64RetMoney
}