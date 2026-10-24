package hbsl

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"example.invalid/stretchr/testify/assert"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 測試雷包紅包設定
func TestHbslLogic_SetBombHongBao(t *testing.T) {
	tgrs := fixtures.NewTestGameRoomServer()
	pGame := tgrs.GameLogicManager.GetGameLogic(0).(*HbslLogic)

	pGame.i32HBNum = 10

	pGame.pNowHB = new(HBRecord)
	pGame.pNowHB.pHB = new(HB)
	pGame.pNowHB.ary_i64SHB = make([]int64, 0)
	pGame.pNowHB.pHB.i64Chip = 10000

	_, i64Tax := pGame.GetUserWinAndTax(pGame.pNowHB.pHB.i64Chip)

	for i := 0; i < 50; i++ {
		pGame.pNowHB.pHB.i32BombNumber = int32(rand.Intn(10))

		bBomb := false
		iRnd := util.RndInt(0, 2)
		if iRnd == 0 {
			bBomb = false
		} else {
			bBomb = true
		}

		pGame.pNowHB.ary_i64SHB = pGame.SetBombHongBao(bBomb)
		iBombCount, i64Sum := 0, int64(0)
		for _, i64Value := range pGame.pNowHB.ary_i64SHB {
			i64Sum += i64Value
			if i64Value%10 == int64(pGame.pNowHB.pHB.i32BombNumber) {
				iBombCount++
			}
		}

		str := fmt.Sprintf("%+v, 紅包金額[%d], 稅[%d], 總和[%d], 含雷[%t], 雷號[%d], 雷包數[%d]", pGame.pNowHB.ary_i64SHB,
			pGame.pNowHB.pHB.i64Chip, i64Tax, i64Sum, bBomb, pGame.pNowHB.pHB.i32BombNumber, iBombCount)
		t.Log(str)

		assert.False(t, bBomb && iBombCount <= 0, fmt.Sprintf("bBomb[%t], iBombCount[%d]", bBomb, iBombCount))
		assert.False(t, !bBomb && iBombCount > 0, fmt.Sprintf("bBomb[%t], iBombCount[%d]", bBomb, iBombCount))
	}
}

// 測試產生作弊模式紅包
func TestHbslLogic_CheatHongBao(t *testing.T) {
	tgrs := fixtures.NewTestGameRoomServer()
	pGame := tgrs.GameLogicManager.GetGameLogic(0).(*HbslLogic)

	pGame.i32HBNum = 10

	pGame.pNowHB = new(HBRecord)
	pGame.pNowHB.pHB = new(HB)
	pGame.pNowHB.pHB.pUser = new(User)
	pGame.pNowHB.pHB.pUser.pGameUser = gamelogic.NewGameUser()
	pGame.pNowHB.ary_i64SHB = make([]int64, 0)
	pGame.pNowHB.pHB.i64Chip = 10000

	i32HBNum := int32(10)
	i64HBMoney, i64Tax := pGame.GetUserWinAndTax(pGame.pNowHB.pHB.i64Chip)

	gbCheatMode = true

	for a := 0; a < 2; a++ {
		if gbCheatMode {
			gbCheatMode = false
		} else {
			gbCheatMode = true
		}

		for i := 0; i < 50; i++ {
			pGame.pNowHB.pHB.i32BombNumber = int32(rand.Intn(10))

			iHBGiver := rand.Intn(3)
			strGiver, strBomb, strCheat := "", "", ""
			switch iHBGiver {
			case 0: // 系統
				pGame.pNowHB.pHB.i32UserID = -1
				strGiver = "系統"
				strBomb = "，應有雷包"
			case 1: // 機器人
				pGame.pNowHB.pHB.i32UserID = 1001
				pGame.pNowHB.pHB.pUser.pGameUser.UserType = gamelogic.UserTypeRebot
				strGiver = "機器人"
				strBomb = "，應有雷包"
			case 2: // 玩家
				pGame.pNowHB.pHB.i32UserID = 1002
				pGame.pNowHB.pHB.pUser.pGameUser.UserType = gamelogic.UserTypeNormal
				strGiver = "玩家"
				strBomb = "，不應有雷包"
			}

			if !gbCheatMode {
				strBomb = ""
			} else {
				strCheat = "【作弊模式】"
			}
			str := fmt.Sprintf("%s發紅包者：%s，雷號：%d，紅包總金額[%d]，稅[%d]%s", strCheat, strGiver, pGame.pNowHB.pHB.i32BombNumber,
				pGame.pNowHB.pHB.i64Chip, i64Tax, strBomb)
			t.Log(str)

			iBombCount, i64Sum := 0, int64(0)
			aryi64NewHongBao := newHongBao(i64HBMoney, i32HBNum)
			for _, i64Value := range aryi64NewHongBao {
				i64Sum += i64Value
				if i64Value%10 == int64(pGame.pNowHB.pHB.i32BombNumber) {
					iBombCount++
				}
			}
			str = fmt.Sprintf("原紅包【%+v，紅包總金額[%d]，雷包數[%d]】", aryi64NewHongBao, i64Sum, iBombCount)

			iBombCount, i64Sum = 0, int64(0)
			pGame.pNowHB.ary_i64SHB = pGame.CheatHongBao(aryi64NewHongBao)
			for _, i64Value := range pGame.pNowHB.ary_i64SHB {
				i64Sum += i64Value
				if i64Value%10 == int64(pGame.pNowHB.pHB.i32BombNumber) {
					iBombCount++
				}
			}
			str += fmt.Sprintf("新紅包【%+v，紅包總金額[%d]，雷包數[%d]】", pGame.pNowHB.ary_i64SHB, i64Sum, iBombCount)
			t.Log(str)

			if gbCheatMode {
				if iHBGiver == 0 || iHBGiver == 1 {
					// 系統或機器人發紅包
					assert.False(t, iBombCount <= 0, fmt.Sprintf("應有雷包，雷包數[%d]", iBombCount))
				} else {
					// 玩家發紅包
					assert.False(t, iBombCount > 0, fmt.Sprintf("不應有雷包，雷包數[%d]", iBombCount))
				}
			} else {
				bSame := true
				if len(aryi64NewHongBao) != len(pGame.pNowHB.ary_i64SHB) {
					bSame = false
				} else {
					for j := 0; j < len(aryi64NewHongBao); j++ {
						if aryi64NewHongBao[j] != pGame.pNowHB.ary_i64SHB[j] {
							bSame = false
							break
						}
					}
				}
				assert.False(t, !bSame, "紅包內容不相同")
			}
		}
		t.Log("----------------------------------------------")
	}
}

func newHongBao(i64HBMoney int64, i32HBNum int32) []int64 {
	if i64HBMoney <= 0 || i32HBNum <= 0 {
		return make([]int64, 0)
	}

	aryi64RetHB := make([]int64, 0)
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

	i64Money := i64HBMoney
	for _, i64Rnd := range aryi64Rand {
		i64Value := i64Money - i64Rnd
		aryi64RetHB = append(aryi64RetHB, i64Value)
		i64Money = i64Rnd
	}

	return aryi64RetHB
}

// 測試搶作弊紅包
func TestHbslLogic_GrabCheatHongBao(t *testing.T) {
	tgrs := fixtures.NewTestGameRoomServer()
	pGame := tgrs.GameLogicManager.GetGameLogic(0).(*HbslLogic)

	pGame.i32HBNum = 10

	pGame.pNowHB = new(HBRecord)
	pGame.pNowHB.pHB = new(HB)
	pGame.pNowHB.pHB.pUser = new(User)
	pGame.pNowHB.pHB.pUser.pGameUser = gamelogic.NewGameUser()
	pGame.pNowHB.ary_i64SHB = make([]int64, 0)
	pGame.pNowHB.pHB.i64Chip = 10000

	i32HBNum := int32(10)
	i64HBMoney, i64Tax := pGame.GetUserWinAndTax(pGame.pNowHB.pHB.i64Chip)

	gbCheatMode = true

	for i := 0; i < 50; i++ {
		pGame.pNowHB.pHB.i32BombNumber = int32(rand.Intn(10))

		iHBGiver := rand.Intn(3)
		strGiver := ""
		switch iHBGiver {
		case 0: // 系統
			pGame.pNowHB.pHB.i32UserID = -1
			strGiver = "系統"
		case 1: // 機器人
			pGame.pNowHB.pHB.i32UserID = 1001
			pGame.pNowHB.pHB.pUser.pGameUser.UserType = gamelogic.UserTypeRebot
			strGiver = "機器人"
		case 2: // 玩家
			pGame.pNowHB.pHB.i32UserID = 1002
			pGame.pNowHB.pHB.pUser.pGameUser.UserType = gamelogic.UserTypeNormal
			strGiver = "玩家"
		}
		str := fmt.Sprintf("發紅包者：%s，雷號：%d，紅包總金額[%d]，稅[%d]", strGiver, pGame.pNowHB.pHB.i32BombNumber,
			pGame.pNowHB.pHB.i64Chip, i64Tax)
		t.Log(str)

		aryi64NewHongBao := newHongBao(i64HBMoney, i32HBNum)
		pGame.pNowHB.ary_i64SHB = pGame.CheatHongBao(aryi64NewHongBao)
		iBombCount := 0
		for _, i64Value := range pGame.pNowHB.ary_i64SHB {
			if i64Value%10 == int64(pGame.pNowHB.pHB.i32BombNumber) {
				iBombCount++
			}
		}
		str = fmt.Sprintf("%+v 雷包數[%d]", pGame.pNowHB.ary_i64SHB, iBombCount)
		t.Log(str)

		pUser := NewUser(gamelogic.NewGameUser())
		iRnd := rand.Intn(2)
		if iRnd == 0 {
			pUser.pGameUser.UserType = gamelogic.UserTypeNormal
			str = "搶紅包者：玩家，"
		} else {
			pUser.pGameUser.UserType = gamelogic.UserTypeRebot
			str = "搶紅包者：機器人，"
		}
		i64Money := pGame.GrabCheatHongBao(pUser)
		str += fmt.Sprintf("金額[%d], %+v", i64Money, pGame.pNowHB.ary_i64SHB)
		t.Log(str)

		bBomb := false
		if i64Money%10 == int64(pGame.pNowHB.pHB.i32BombNumber) {
			bBomb = true
		}
		if iHBGiver == 0 || iHBGiver == 1 {
			assert.False(t, pUser.pGameUser.IsRobot() && bBomb, "機器人拿到雷包")
			assert.False(t, !pUser.pGameUser.IsRobot() && !bBomb, "玩家沒拿到雷包")
		} else {
			assert.False(t, bBomb, "搶紅包者拿到雷包")
		}

		t.Log("--------------------------------------------------")
	}
}
