//==================================================================
//文件名： ttzuseropt.go
//说明：   推筒子用户信息
//修改记录：
//==================================================================
package ttz

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"	
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// 玩家数据结构
type TtzUserData struct {
	User              *gamelogic.GameUser //用户
	userTotalMoney    int64               //用户总金钱
	userCards         *poker.CardGroup    //用户手牌
	isAction          bool
	nRateRob          int32 //抢庄倍数
	nRateBet          int32 //押注倍数
	nCardType         int32 //牌型
	nCardCount        int32 //牌点数
	nMaxRob           int32 //玩家可以抢庄的最大倍数
	nMaxBet           int32 //玩家可以押注的最大倍数
	nCurRoundWinMoney int64 //当前小局输赢金钱
	nCurGameWinMoney  int64 //当前赢得的总金币
	bIsRob            bool  //是否已经抢庄
	bIsBet            bool  //是否已经押注
	bIsLight          bool  //是否已经亮牌
}

// 操作信息结构
type OpInfo struct {
	nOprateType   int32
	nOprateJetton int32
}

// 玩家数据
type sliceUserData struct {
	data []*TtzUserData
}

// 创建玩家数据
func NewTtzUserData(user *gamelogic.GameUser) *TtzUserData {
	ud := new(TtzUserData)
	ud.userTotalMoney = user.Score
	ud.User = user
	ud.userCards = new(poker.CardGroup)
	ud.isAction = false
	ud.nRateRob = 0
	ud.nRateBet = 2
	ud.nCardType = 0
	ud.nCardCount = 0
	ud.nCurRoundWinMoney = 0
	ud.nCurGameWinMoney = 0
	ud.bIsRob = false
	ud.bIsBet = false
	ud.bIsLight = false
	ud.nMaxRob = -1
	ud.nMaxBet = -1
	return ud
}

// 重置玩家数据
func (ud *TtzUserData) reset() {
	ud.userCards = new(poker.CardGroup)
	ud.isAction = false
	ud.nRateRob = 0
	ud.nRateBet = 2
	ud.nCardType = 0
	ud.nCardCount = 0
	ud.nCurRoundWinMoney = 0
	ud.bIsRob = false
	ud.bIsBet = false
	ud.bIsLight = false
	ud.nMaxRob = -1
	ud.nMaxBet = -1
}

// 打印筹码池
func (ud *TtzUserData) PrintJettonInfo() {
	//mlog.Debug("当前下注数：%d,当前剩余筹码：%d", ud.currBet, ud.userMoney)
}

// 计算玩家数据长度
func (sud *sliceUserData) Len() int {
	return len(sud.data)
}

// 设置玩家操作值
func (ud *TtzUserData) setAction(b bool) {
	//mlog.Debug("setAction userid[%d],seatNo[%d],set[%v]", ud.User.UserID, ud.User.SeatNo, b)
	ud.isAction = b
}

// 设置抢庄值
func (ud *TtzUserData) setRob(b bool) {
	//mlog.Debug("setRob userid[%d],seatNo[%d],set[%v]", ud.User.UserID, ud.User.SeatNo, b)
	ud.bIsRob = b
}

// 设置下注值
func (ud *TtzUserData) setBet(b bool) {
	//mlog.Debug("setBet userid[%d],seatNo[%d],set[%v]", ud.User.UserID, ud.User.SeatNo, b)
	ud.bIsBet = b
}

// 设置亮牌值
func (ud *TtzUserData) setLight(b bool) {
	//mlog.Debug("setLight userid[%d],seatNo[%d],set[%v]", ud.User.UserID, ud.User.SeatNo, b)
	ud.bIsLight = b
}

// 获取是否操作
func (ud *TtzUserData) getAction() bool {
	return ud.isAction
}

// 获取是否抢庄
func (ud *TtzUserData) getRob() bool {
	return ud.bIsRob
}

// 获取是否下注
func (ud *TtzUserData) getBet() bool {
	return ud.bIsBet
}

// 获取是否亮牌
func (ud *TtzUserData) getLight() bool {
	return ud.bIsLight
}
