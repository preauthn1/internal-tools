//==================================================================
//文件名： dzpuseropt.go
//说明：   德州扑克用户数据类的成员和方法
//修改记录：
//==================================================================
package dzpk

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"sort"
)

// 德州扑克玩家数据
type DzpkUserData struct {
	userCards         *poker.CardGroup    //用户手牌
	User              *gamelogic.GameUser //用户
	compareUserSeatNo []int32             //这些人和我比过牌
	betCount          int32               //下注手数计数
	userRole          playerRole          //玩家角色
	seatNo            int32               //玩家座位号
	cardType          int32               //用户牌型
	isAction          bool                //是否已经操作过
	curRoundSettle    int64               //当局获胜所得筹码
	pokerTotalGroup   *poker.CardGroup    //总牌
	pokerMaxCardGroup *poker.CardGroup    //最大牌型构成的牌

	currBet           int64     //当前小局下注数
	userMoney         int64     //用户钱币拷贝
	curTurnBet        int64     //当前轮下注数
	isFold            bool      //是否弃牌
	isAllin           bool      //是否allin
	isCanCall         bool      //是否能跟牌
	isCanRaise        bool      //是否能加注
	isCanCheck        bool      //可过
	isCanAllIn        bool      //可ALLIN
	userOprateInfo    []*OpInfo //玩家操作信息
	userPreOprateInfo []*OpInfo //玩家预操作信息
	curOpration       int32     //玩家当前操作动作
	tax               int64     //抽水

}
type OpInfo struct {
	nOprateType   int32
	nOprateJetton int64
}
type sliceUserData struct {
	data []*DzpkUserData
}

// 创建数据结构
func NewDzpkUserData(user *gamelogic.GameUser) *DzpkUserData {
	userdata := new(DzpkUserData)
	userdata.userMoney = user.Score
	userdata.betCount = 0

	userdata.userRole = ROLE_NONE
	userdata.User = user
	userdata.userCards = nil
	userdata.seatNo = user.SeatNo
	userdata.currBet = 0
	userdata.curTurnBet = 0
	userdata.isAction = false
	userdata.cardType = HIGHCARD
	userdata.curRoundSettle = 0
	//用户状态
	userdata.isCanCall = true
	userdata.isCanRaise = true
	userdata.isCanCheck = true
	userdata.isCanAllIn = false

	userdata.isFold = false
	userdata.isAllin = false
	userdata.userOprateInfo = nil
	userdata.tax = 0

	return userdata
}

//重置数据
func (ud *DzpkUserData) userResetData() {
	ud.betCount = 0

	ud.userRole = ROLE_NONE
	ud.userCards = nil
	ud.currBet = 0
	ud.curTurnBet = 0
	ud.isAction = false
	ud.cardType = HIGHCARD
	ud.curRoundSettle = 0
	//用户状态
	ud.isCanCall = true
	ud.isCanRaise = true
	ud.isCanCheck = true
	ud.isCanAllIn = false

	ud.isFold = false
	ud.isAllin = false
	ud.userOprateInfo = nil
	ud.tax = 0
}

//用户是否可操作
func (ud *DzpkUserData) userCanOprate() bool {
	if ud.isFold || ud.isAllin {
		return false
	}

	return true
}

// 设置金币
func (ud *DzpkUserData) SetMoney(money int64) {
	ud.userMoney += money
	ud.currBet -= money
	ud.curTurnBet -= money
}

// 设置弃牌
func (ud *DzpkUserData) SetDrop() {
	ud.isFold = true
}

//设置比较的座位号
func (ud *DzpkUserData) setCompareUserSeatNo(seatno int32) {
	ud.compareUserSeatNo = append(ud.compareUserSeatNo, seatno)
}

//获取用户角色
func (ud *DzpkUserData) getUserRole() playerRole {
	return ud.userRole
}

//更新用户操作信息
func (ud *DzpkUserData) updateUserOpInfo(nCurMaxJetton, nSmallBlind int64) {
	mlog.Debug("updateUserOpInfo,seatNo=%d,curbet=%d,CurMaxJetton=%d", ud.seatNo, ud.currBet, nCurMaxJetton)
	ud.updateUserState(nCurMaxJetton, nSmallBlind)
	ud.userOprateInfo = make([]*OpInfo, 0)
	infoFold := new(OpInfo)
	infoFold.nOprateType = int32(FOLD)
	infoFold.nOprateJetton = 0
	ud.userOprateInfo = append(ud.userOprateInfo, infoFold)
	if ud.isCanCall {
		info := new(OpInfo)
		info.nOprateType = int32(CALL)
		info.nOprateJetton = nCurMaxJetton - ud.currBet
		ud.userOprateInfo = append(ud.userOprateInfo, info)
	}
	if ud.isCanRaise {
		info := new(OpInfo)
		info.nOprateType = int32(RAISE)
		info.nOprateJetton = nCurMaxJetton - ud.currBet
		ud.userOprateInfo = append(ud.userOprateInfo, info)
	}
	if ud.isCanCheck {
		info := new(OpInfo)
		info.nOprateType = int32(CHECK)
		info.nOprateJetton = 0
		ud.userOprateInfo = append(ud.userOprateInfo, info)
	}
	if ud.isCanAllIn {
		info := new(OpInfo)
		info.nOprateType = int32(ALLIN)
		info.nOprateJetton = ud.userMoney
		ud.userOprateInfo = append(ud.userOprateInfo, info)
	}
}

// 更新玩家操作状态
func (ud *DzpkUserData) updateUserState(nCurMaxJetton, nSmallBlind int64) {
	mlog.Debug("updateUserState")
	if ud.userMoney >= nCurMaxJetton-ud.currBet && nCurMaxJetton-ud.currBet > 0 {
		ud.isCanCall = true
	} else {
		ud.isCanCall = false
	}
	//check
	if ud.currBet == nCurMaxJetton {
		ud.isCanCheck = true
	} else {
		ud.isCanCheck = false
	}
	//raise
	if ud.userMoney > nCurMaxJetton-ud.currBet+nSmallBlind {
		ud.isCanRaise = true
	} else {
		ud.isCanRaise = false
	}
	//allin
	if ud.userMoney <= nCurMaxJetton-ud.currBet+nSmallBlind {
		ud.isCanAllIn = true
	} else {
		ud.isCanAllIn = false
	}

}

//预操作
func (ud *DzpkUserData) updatePreOpration(dzpk *DzpkLogic) {
	ud.userPreOprateInfo = make([]*OpInfo, 0)
	//Fold
	{
		op := new(OpInfo)
		op.nOprateType = int32(FOLD)
		op.nOprateJetton = 0
		ud.userPreOprateInfo = append(ud.userPreOprateInfo, op)
	}
	//看牌
	if ud.currBet == dzpk.betCurrMax {
		op := new(OpInfo)
		op.nOprateType = int32(CHECK)
		op.nOprateJetton = 0
		ud.userPreOprateInfo = append(ud.userPreOprateInfo, op)
	}
	//跟注
	if ud.userMoney > dzpk.betCurrMax-ud.currBet && dzpk.betCurrMax-ud.currBet > 0 {
		op := new(OpInfo)
		op.nOprateType = int32(CALL)
		op.nOprateJetton = dzpk.betCurrMax - ud.currBet
		ud.userPreOprateInfo = append(ud.userPreOprateInfo, op)
	}
	//跟任何注CALLANY
	if ud.userMoney > dzpk.betCurrMax-ud.currBet {
		op := new(OpInfo)
		op.nOprateType = int32(CALLANY)
		op.nOprateJetton = dzpk.betCurrMax - ud.currBet
		ud.userPreOprateInfo = append(ud.userPreOprateInfo, op)
	}

}

//重置当前操作状态
func (ud *DzpkUserData) resetCurOpration() {
	ud.curOpration = int32(NONE)
}

// 设置初始下注操作
func (ud *DzpkUserData) setBeginOprateForBet(nCurMaxJetton int64) {
	mlog.Debug("setBeginOprateForBet")
	ud.userOprateInfo = make([]*OpInfo, 0)
	//弃牌
	{
		info := new(OpInfo)
		info.nOprateType = int32(FOLD)
		info.nOprateJetton = 0
		ud.userOprateInfo = append(ud.userOprateInfo, info)
	}
	//加注
	{
		info := new(OpInfo)
		info.nOprateType = int32(RAISE)
		info.nOprateJetton = nCurMaxJetton - ud.currBet
		ud.userOprateInfo = append(ud.userOprateInfo, info)
	}
	//看牌押注
	if nCurMaxJetton-ud.currBet > 0 {
		info := new(OpInfo)
		info.nOprateType = int32(CALL)
		info.nOprateJetton = 0
		ud.userOprateInfo = append(ud.userOprateInfo, info)
	} else {
		info := new(OpInfo)
		info.nOprateType = int32(CHECK)
		info.nOprateJetton = 0
		ud.userOprateInfo = append(ud.userOprateInfo, info)
	}
}

// 获取玩家最大牌组
func (ud *DzpkUserData) getMaxCardsGroup(cardPublic *poker.CardGroup) {

	ud.pokerMaxCardGroup = new(poker.CardGroup)
	cardsGroup := new(poker.CardGroup)
	if ud.userCards != nil {
		cardsGroup.Cards = append(cardsGroup.Cards, ud.userCards.Cards...)
	}
	if cardPublic != nil {
		cardsGroup.Cards = append(cardsGroup.Cards, cardPublic.Cards...)
	}
	mlog.Debug("getMaxCardsGroup begin,and cardsGroup size=%d.", cardsGroup.Len())
	ud.pokerMaxCardGroup = GetCardsByCardType(cardsGroup, ud.cardType)
}

// 打印玩家筹码状态
func (ud *DzpkUserData) PrintJettonInfo() {
	mlog.Debug("当前下注数：%d,当前剩余筹码：%d", ud.currBet, ud.userMoney)
}

//拼接手牌和公牌
func (ud *DzpkUserData) andUserCardAndPublic(pubCardGroup *poker.CardGroup) {
	ud.pokerTotalGroup = new(poker.CardGroup)
	if pubCardGroup != nil {
		ud.pokerTotalGroup.Cards = append(ud.pokerTotalGroup.Cards, pubCardGroup.Cards...)
	}
	if ud.userCards != nil {
		ud.pokerTotalGroup.Cards = append(ud.pokerTotalGroup.Cards, ud.userCards.Cards...)
	}
}
func (sud *sliceUserData) Len() int {
	return len(sud.data)
}

func (sud *sliceUserData) Swap(i, j int) {
	sud.data[i], sud.data[j] = sud.data[j], sud.data[i]
}

//按照从小到大顺序排列
func (sud *sliceUserData) Less(i, j int) bool {
	return sud.data[i].currBet < sud.data[j].currBet
}

//按照从小到大顺序排列
func (sud *sliceUserData) ReSort() {
	sort.Sort(sud)
}
