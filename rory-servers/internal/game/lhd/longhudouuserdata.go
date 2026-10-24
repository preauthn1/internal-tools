package lhd

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	//	"mlog"
)

var LhdUserHistoryStore = make(map[int32][]*longhudouuserdata)

type bettype struct {
	longhudoutype int32 //下注类型
	money         int64 //下注的钱
}

type longhudouuserdata struct {
	user        *gamelogic.GameUser
	userbet     []*bettype
	money       int64 //用户开始的钱
	moneylast   int64 //用户最开始的钱，用来更新用户中途加钱
	win         int64 //用户赢的钱
	updatemoney int64 //用户更新的钱
	clearcount  int32 //清理次数
	betAmount   int64
	winAmount   int64
	flowAmount  int64
}

//新建用户数据实例
func newLongHuDouUserData(user *gamelogic.GameUser) *longhudouuserdata {
	ud := new(longhudouuserdata)
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
func (ud *longhudouuserdata) clearData() {
	ud.userbet = make([]*bettype, 0)
	ud.money = ud.user.Score
	ud.moneylast = ud.user.Score
	ud.win = 0
}

//设置玩家下注的钱
func (ud *longhudouuserdata) setUserBet(typevalue int32, money int32) {
	for _, v := range ud.userbet {
		if v.longhudoutype == typevalue {
			v.money += int64(money)
			ud.money -= int64(money)
			return
		}
	}

	v := new(bettype)
	v.longhudoutype = typevalue
	v.money = int64(money)
	ud.userbet = append(ud.userbet, v)
	ud.money -= int64(money)
}

//获取玩家下注
func (ud *longhudouuserdata) getUserBet() []*bettype {
	return ud.userbet
}

func (ud *longhudouuserdata) getUserBetByType(t int32) int64 {
	for _, bet := range ud.userbet {
		if bet.longhudoutype == t {
			return bet.money
		}
	}
	return 0
}

func (ud *longhudouuserdata) getUserBetTotal() int64 {
	sumbet := int64(0)
	for _, b := range ud.userbet {
		sumbet += b.money
	}
	return sumbet
}

func (ud *longhudouuserdata) calUserWinRate(size int32) int64 {
	if userhis, ok := LhdUserHistoryStore[ud.user.UserID]; ok {
		wincount := int32(0)
		arr := userhis
		if len(userhis) > int(size) {
			arr = userhis[:size]
		} else {
			size = int32(len(arr))
		}

		for _, his := range arr {
			if his.win > 0 {
				wincount++
			}
		}
		return int64((float32(wincount) / float32(size)) * 100)
	}
	return 0
}

func (ud *longhudouuserdata) getShowWinLose() int64 {
	if ud.win == 0 {
		return -ud.getUserBetTotal()
	} else {
		return ud.win
	}
}
