package by

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"time"
)

type boominfo struct {
	boomid    int32 //炸弹ID
	odd       int32 //炸弹倍率
	shoottime int64 //发射时间，测试用
	boomtype  int32
}

type byuserdata struct {
	user           *gamelogic.GameUser //用户
	money          int64               //当前的钱
	boomNumber     int32               //用户当前选择的倍率
	boomindex      int32               //当前倍率的索引值
	suoding        bool                //是否是锁定状态
	allwin         int64               //在捕鱼中总的输赢
	lastshoottime  int64               //上一次射击时间
	boomdata       map[int32]*boominfo //炮弹的ID
	boomcount      int32               //炮弹计数
	teshufishtype  int32               //死亡的特殊鱼的类型
	fishtype       int32               //下次捕中的鱼类型
	winmoney       int64               //本次进入房间输赢的钱
	rate           float64             //用户攻擊概率变化
	currodd        int32               //当前子弹的倍率
	bcreatetimer   bool                //是否建立了时间检测
	lastframemoney int64               //上一帧用户的钱，用来给客户端校正
	lastice        int64               //上一次放冰冻的时间
	fishwinloss    map[int32]int32     //魚貢獻玩家多少, fishwinloss[fishid]money
	teshufishdiedinfo
	sumbetmoney int64   //玩家總押注金額
	ispochan    bool    //玩家破產
	candiedrate float32 //磚頭砲大死魚的動態機率
}

type teshufishdiedinfo struct {
	bossfishid          int32
	diedbosstype        int32
	allFishOdd          int32
	teshufishfishodd    int32
	spmaxfishmulti      int32 //倍數門檻	N倍以下(含)的普魚是可以被打死的	-1表示不使用
	fightspmaxfishcount int32 //攻擊隻數
}

func newuserdata(user *gamelogic.GameUser) *byuserdata {
	ud := new(byuserdata)
	ud.user = user
	ud.money = 0
	ud.suoding = false
	ud.boomNumber = 0
	ud.lastshoottime = int64(time.Now().UnixNano())
	ud.boomdata = make(map[int32]*boominfo)
	ud.fishwinloss = make(map[int32]int32)
	ud.boomcount = 0
	ud.fishtype = 0
	ud.teshufishtype = 0
	ud.winmoney = 0
	ud.rate = 1.0
	ud.lastframemoney = user.Score
	ud.lastice = 0
	return ud
}

//用户射击一次消耗的钱
func (ud *byuserdata) userShoot(fid int64, boomMoney int32) bool {
	//钱不够了
	if (ud.user.Score + ud.money - int64(ud.currodd)) < int64(boomMoney) {
		mlog.Debug("用户%d钱不够了,boomMoney为：%v,用户当前的钱为：%v", ud.user.UserID, boomMoney, ud.user.Score+ud.money)
		return false
	}
	ud.lastshoottime = int64(time.Now().UnixNano())
	//mlog.Debug("用户射击一次消耗的钱 userShoot time[%v]", ud.lastshoottime)
	return true
}

//设置用户炸彈編號
func (ud *byuserdata) setUserBoomNumber(boomNumber int32, index int32) {
	ud.boomindex = index
	ud.boomNumber = boomNumber
}

//获取用户倍率
func (ud *byuserdata) getUserOdd() int32 {
	return ud.currodd
}

//更新玩家的钱
func (ud *byuserdata) updateUserMoney(money int32) {
	ud.money += int64(money)
	ud.allwin += int64(money)
	ud.winmoney += int64(money)
}
//更新流水, 有效投注金額, 返還金額
func (ud *byuserdata) updateBetAmount(flowAmount, betAmount, payAmount int32) {
	ud.user.FlowAmount += int64(flowAmount)
	ud.user.BetAmount += int64(betAmount)
	ud.user.PayAmount += int64(payAmount)
}

//更新玩家流水
//func (ud *byuserdata) updateFlowAmount(money int32) {
//	ud.user.FlowAmount += int64(money)
//}

//设置锁定
func (ud *byuserdata) setSuoDing() {
	ud.suoding = !ud.suoding
}

//移除一个炮弹
func (ud *byuserdata) removeBoom(boomid int32) bool {
	if f, ok := ud.boomdata[boomid]; ok {
		if f.odd > 0 {
			ud.currodd = f.odd
		}
		delete(ud.boomdata, boomid)
		ud.boomcount--
		mlog.Debug("用户%d移除一颗炮弹id[%d],剩余[%d] f.odd[%v]", ud.user.UserID, boomid, ud.boomcount, f.odd)
		return true
	}

	return false
}

//添加一个炮弹
func (ud *byuserdata) addBoom(boomid, boomOdd int32) {
	bd := new(boominfo)
	bd.boomid = boomid
	bd.odd = boomOdd
	//mlog.Debug("addBoom oddType[%v] boomid[%v]", bd.odd, bd.boomid)
	bd.shoottime = time.Now().Unix()
	ud.boomdata[boomid] = bd
	ud.boomcount++

	for _, v := range ud.boomdata {
		if (time.Now().Unix() - v.shoottime) >= 30 {
			mlog.Debug("用户%d，子弹ID%d,超过30秒未碰撞,发射时间：%d", ud.user.UserID, v.boomid, v.shoottime)
		}
	}
}

//设置下次是否必中死亡特殊鱼的类型
func (ud *byuserdata) setTeShuFishType(t int32) {
	ud.teshufishtype = t
}

//设置下次是否必中死亡鱼的类型
func (ud *byuserdata) setFishType(t int32) {
	ud.fishtype = t
}

func (ud *byuserdata) findBoomById(id int32) bool {
	if _, ok := ud.boomdata[id]; ok {
		return true
	}

	return false
}

//保存上一帧的钱
func (ud *byuserdata) setLastFrameMoney() {
	ud.lastframemoney = ud.user.Score + ud.money
}

func (ud *byuserdata) delAllBoom() {
	ud.boomdata = make(map[int32]*boominfo, 0)
	ud.boomcount = 0
}

////玩家释放冰冻技能
//func (ud *byuserdata) castIce(frozenresettime int64) bool {
//	lastcast := time.Now().UnixNano() - ud.lastice
//	if lastcast > frozenresettime {
//		//重置冰冻释放时间
//		ud.lastice = time.Now().UnixNano()
//		return true
//	}
//
//	return false
//}
