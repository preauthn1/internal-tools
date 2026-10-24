package fkeb

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

type bettype struct {
	fkebtype int32 //下注类型
	money    int64 //下注的钱
}

type fkebuserdata struct {
	user        *gamelogic.GameUser
	userbet     []*bettype
	money       int64 //用户开始的钱
	moneylast   int64 //用户最开始的钱，用来更新用户中途加钱
	win         int64 //用户赢的钱
	updatemoney int64 //用户更新的钱
	clearcount  int32 //清理次数
}

//新建用户数据实例
func newLongHuDouUserData(user *gamelogic.GameUser) *fkebuserdata {
	ud := new(fkebuserdata)
	ud.user = user
	ud.userbet = make([]*bettype, 0)
	ud.money = user.Score
	ud.moneylast = user.Score
	ud.win = 0
	ud.clearcount = 0
	ud.updatemoney = 0
	return ud
}

//清除数据
func (ud *fkebuserdata) clearData() {
	ud.userbet = make([]*bettype, 0)
	ud.money = ud.user.Score
	ud.moneylast = ud.user.Score
	ud.win = 0
}

//设置玩家下注的钱
func (ud *fkebuserdata) setUserBet(typevalue int32, money int32) {
	if ud.user.IsRobot() {
		mlog.Debug("机器人%d,下注%d,v.fkebtype=%d", ud.user.UserID, money, typevalue)
	}
	for _, v := range ud.userbet {
		if v.fkebtype == typevalue {
			v.money += int64(money)
			ud.money -= int64(money)
			return
		}
	}

	v := new(bettype)
	v.fkebtype = typevalue
	v.money = int64(money)
	ud.userbet = append(ud.userbet, v)
	ud.money -= int64(money)
}

//获取玩家下注
func (ud *fkebuserdata) getUserBet() []*bettype {
	return ud.userbet
}

func (ud *fkebuserdata) getUserBetByType(t int32) int64 {
	for _, bet := range ud.userbet {
		if bet.fkebtype == t {
			return bet.money
		}
	}
	return 0
}
