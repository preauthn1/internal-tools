package qznn

import (
	"math/rand"
	"sort"
	"time"
)

//用户倍数
type userBet struct {
	userID   int32 //用户ID
	betTimes int32 //倍数
	sysSet   bool  //是否系统设定
}

//用户选择的加倍信息
type gameUserBet struct {
	userBets []*userBet
}

//检查用户倍数信息
func (ubc *gameUserBet) checkUserBet(userid int32) bool {
	for _, ub := range ubc.userBets {
		if ub.userID == userid {
			return true
		}
	}

	return false
}

//添加用户倍数信息
func (ubc *gameUserBet) addUserBet(userid int32, bettimes int32, sys bool) bool {
	if ubc.checkUserBet(userid) {
		return false
	}

	ub := new(userBet)
	ub.userID = userid
	ub.betTimes = bettimes
	ub.sysSet = sys

	ubc.userBets = append(ubc.userBets, ub)

	return true
}

//抢庄
type zhuangJiaBets struct {
	*gameUserBet
}

func newZhuangJiaBets() *zhuangJiaBets {
	zjbets := new(zhuangJiaBets)
	zjbets.gameUserBet = new(gameUserBet)

	return zjbets
}

//添加用户最小抢庄倍数
func (qz *zhuangJiaBets) addMinUserBet(userid int32) {
	qz.addUserBet(userid, getNNConfig().minZhuangTimes(), false)
}

//获取庄家信息
func (qz *zhuangJiaBets) getZhuang() *userBet {
	sort.Sort(qz)

	maxbettimes := qz.userBets[0].betTimes
	if qz.userBets[0].betTimes > qz.userBets[1].betTimes {

		return qz.userBets[0]
	}

	ubts := make([]*userBet, 0)

	for _, ubt := range qz.userBets {
		if ubt.betTimes == maxbettimes {
			ubts = append(ubts, ubt)
		}
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	ubt := ubts[rnd.Intn(len(ubts))]

	//用户放弃抢庄,但随机到庄家
	if ubt.betTimes == 0 {
		ubt.betTimes = getNNConfig().minZhuangTimes()
		ubt.sysSet = false
	}

	return ubt
}

func (qz *zhuangJiaBets) Len() int {
	return len(qz.userBets)
}

func (qz *zhuangJiaBets) Swap(i, j int) {
	qz.userBets[i], qz.userBets[j] = qz.userBets[j], qz.userBets[i]
}

func (qz *zhuangJiaBets) Less(i, j int) bool {
	return qz.userBets[i].betTimes > qz.userBets[j].betTimes
}

// 闲家倍数信息
type xianjiaBets struct {
	*gameUserBet
}

func newXianjiaBets() *xianjiaBets {
	xjbets := new(xianjiaBets)
	xjbets.gameUserBet = new(gameUserBet)

	return xjbets
}

//添加闲家最小倍数
func (xjb *xianjiaBets) addMinUserBet(userid int32) {
	xjb.addUserBet(userid, getNNConfig().minXianTimes(), true)
}
