package sgj

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

type sgjBetData struct {
	betType int32 //下注类型
	bet     int64 //下注额度
}

type sgjUserData struct {
	user        *gamelogic.GameUser
	betData     []*sgjBetData
	money       int64 //用户下注总额
	win         int64 //用户赢钱
	updatemoney int64 //用户更新的钱
	clearcount  int32 //清理次数
}

func newUserData(user *gamelogic.GameUser) *sgjUserData {
	ud := new(sgjUserData)
	ud.user = user
	ud.betData = make([]*sgjBetData, 0)
	ud.money = 0
	ud.win = 0
	ud.clearcount = 0
	ud.updatemoney = 0
	return ud
}

func (ud *sgjUserData) clearData() {
	ud.betData = make([]*sgjBetData, 0)
	ud.money = 0
	ud.win = 0
}

func (ud *sgjUserData) bet(betType int32, bet int32) {
	ud.money += int64(bet)
	for _, betData := range ud.betData {
		if betData.betType == betType {
			betData.bet += int64(bet)
			return
		}
	}
	bd := new(sgjBetData)
	bd.betType = betType
	bd.bet = int64(bet)
	ud.betData = append(ud.betData, bd)
}

func (ud *sgjUserData) faJiang(betType int32, odd int32) int64 {
	for _, bd := range ud.betData {
		if bd.betType == betType {
			return bd.bet * int64(odd)
		}
	}

	return 0
}

func (ud *sgjUserData) canBet(bet int32) bool {
	if ud.user.Score >= ud.money+int64(bet) {
		return true
	}

	return false
}

func (ud *sgjUserData) setUserWinMoney(win int64) {
	ud.win = win
}

func (ud *sgjUserData) getUserBet() []*sgjBetData {
	return ud.betData
}

func (ud *sgjUserData) getUserBetByType(t int32) int64 {
	for _, bet := range ud.betData {
		if bet.betType == t {
			return bet.bet
		}
	}
	return 0
}
