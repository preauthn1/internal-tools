package sg

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

type SgUserData struct {
	userCards         *poker.CardGroup    //用户手牌
	User              *gamelogic.GameUser //用户
	userMoney         int64               //用户钱币拷贝
	cardType          int32               //牌类型
	isQiangZhuang     bool                //是否抢庄
	CanQiangZhuangOpt bool                //是否已经进行抢庄操作
	bet               int32               //下注的钱数
	betodd            int32               //下注的倍数
	isBet             bool                //是否已经下注
}

//新建一个玩家数据
func NewSgUerData(user *gamelogic.GameUser) *SgUserData {
	userdata := new(SgUserData)
	userdata.User = user
	userdata.cardType = 0
	userdata.userMoney = user.Score
	userdata.userCards = nil
	userdata.isQiangZhuang = false
	userdata.CanQiangZhuangOpt = false
	userdata.isBet = false
	userdata.bet = 0
	userdata.betodd = 0

	return userdata
}

//设置玩家下注
func (ud *SgUserData) setUserBet(bet int32, betbase int32) {
	ud.betodd = bet
	ud.bet = bet * betbase
	ud.isBet = true
}

//设置玩家抢庄
func (ud *SgUserData) setUserQiangZhuang(qiangzhuang bool) {
	ud.isQiangZhuang = qiangzhuang
	ud.CanQiangZhuangOpt = true
}

//获取用户是否抢庄
func (ud *SgUserData) isQaingZhuang() bool {
	return ud.CanQiangZhuangOpt
}
