package scmj

import (
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mahjong"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

//游戏用户数据
type STUserData struct {
	User                  *gamelogic.GameUser //用户
	timerTimeTuoguanEvent *event.TimerFunc
	handCards             *mahjong.CardGroup        //玩家手牌
	riverCards            *mahjong.CardGroup        //玩家牌河
	listenCard            []*STListenCard           //玩家听牌
	pengCards             []*STPengCard             //玩家碰牌
	gangCards             []*STGangCard             //玩家杠牌
	userMoney             int64                     //玩家金钱拷贝
	winScore              int64                     //玩家当局赢得分数
	intimeScore           int64                     //玩家当前结算赢得分数
	backScore             int64                     //退杠分
	gangLast              int64                     //上次杠得分，杠上炮用于呼叫转移
	seatNo                int32                     //玩家座位号
	color                 int32                     //玩家定缺的花色
	rate                  int32                     //倍率
	moPai                 int32                     //摸到的牌
	isHu                  bool                      //玩家是否已经胡了牌
	isEnd                 bool                      //玩家是否已经结束游戏
	isHuaZhu              bool                      //玩家是否是花猪
	isCanTingPai          bool                      //玩家是否可以听牌
	canOpOther            int32                     //
	canOpEventInfo        *STCanOpEventInfo         //可操作信息
	intimeInfo            []*STScoreRecord          //即时结算记录
	isTuoGuan             bool                      //是否托管
	isDeleteRiver         bool                      //减掉河牌区牌
	eventFlag             bool                      //事件操作标志
	remindData            map[int32][]*STRemindInfo //提示
	huCards               []int32                   //胡牌
	canEarlyChange        bool                      //是否可提早離桌
	isClear               bool                      //是否已結算
}

//创建
func NewUserData(user *gamelogic.GameUser) *STUserData {
	userdata := new(STUserData)
	userdata.User = user
	userdata.userMoney = user.Score
	userdata.seatNo = user.SeatNo
	userdata.handCards = new(mahjong.CardGroup)
	userdata.riverCards = new(mahjong.CardGroup)
	userdata.listenCard = make([]*STListenCard, 0)
	userdata.pengCards = make([]*STPengCard, 0)
	userdata.gangCards = make([]*STGangCard, 0)
	userdata.winScore = 0
	userdata.intimeScore = 0
	userdata.backScore = 0
	userdata.gangLast = 0
	userdata.color = -1
	userdata.rate = 0
	userdata.moPai = 0
	userdata.isHu = false
	userdata.isEnd = false
	userdata.isHuaZhu = false
	userdata.isCanTingPai = false
	userdata.canOpEventInfo = new(STCanOpEventInfo)
	userdata.intimeInfo = make([]*STScoreRecord, 0)
	userdata.isTuoGuan = false
	userdata.isDeleteRiver = false
	userdata.eventFlag = false
	userdata.remindData = make(map[int32][]*STRemindInfo, 0)
	userdata.huCards = make([]int32, 0)
	userdata.canEarlyChange = false
	userdata.isClear = false
	return userdata
}

//玩家听牌
type STListenCard struct {
	card       []int32 //听牌
	rate       int64   //倍率
	cardType   int32   //牌型
	surplusNum int32   //可见剩余张数
}

//碰牌
type STPengCard struct {
	pengCard int32 //玩家碰牌
	seatNo   int32 //来源
	valid    bool  //是否有效
}

//玩家结算信息记录
type STScoreRecord struct {
	Type    int32   //结算类型
	Score   int64   //结算分数
	SeatNo  []int32 //联系座位号
	HuTypes []int32 //胡牌类型，只针对结算类型为胡牌的情况
	valid   bool    //是否有效
	card    int32   //胡牌
	Rate    int32   //倍率
}

//杠牌
type STGangCard struct {
	gangCard  int32   //玩家杠牌
	seatNo    int32   //来源
	gangType  int32   //杠类型
	valid     bool    //是否有效
	losSeatNo []int32 //输钱玩家座位号
}

type STCanOpEventInfo struct {
	OpInfo []*STOpInfo
}
type STOpInfo struct {
	nOpType int32
	nValue  []int32
}

//提示信息
type STRemindInfo struct {
	TingCard int32 //牌
	Number   int32 //剩余数量
	Rate     int32 //倍率
	OutCard  int32 //出牌
}

//获取杠牌倍率
func (ud *STUserData) getRateGang(rateGang, eventcard int32) int32 {
	rate := int32(0)
	// 杠牌的跟
	for _, data := range ud.gangCards {
		if data.valid {
			rate += rateGang
		}
	}
	// 碰牌和手牌也可以组成跟,胡牌和碰牌组成的跟
	for _, data := range ud.pengCards {
		if data.valid {

			if eventcard == data.pengCard {
				rate += rateGang
				continue
			}
			for _, cd := range ud.handCards.Cards {
				if cd == data.pengCard {
					rate += rateGang
					break
				}
			}
		}
	}
	// 手牌单独跟,胡牌和手牌组成的跟
	mapGen := make(map[int32]int32, 0)
	for _, cd := range ud.handCards.Cards {
		mapGen[cd/10]++
	}
	mapGen[eventcard/10]++
	for _, count := range mapGen {
		if count == 4 {
			rate += rateGang
		}
	}
	return rate
}

//玩家是否可以操作牌，如果没有相应的操作则可以轮到下一个玩家操作
func (ud *STUserData) updataCanOpInfo(cd int32, curOp, canGang bool) bool {
	flag := false
	//是否可以胡,点炮胡，抢杠胡，自摸胡
	if ud.canHu(cd) {
		stOpinfo := new(STOpInfo)
		stOpinfo.nOpType = OP_HU
		stOpinfo.nValue = append(stOpinfo.nValue, cd)
		ud.canOpEventInfo.OpInfo = append(ud.canOpEventInfo.OpInfo, stOpinfo)
		mlog.Debug("玩家座位号[%d]可以胡", ud.seatNo)
		flag = true
	}
	//是否可以杠
	if canGang {
		if curOp && ud.canYinGang(cd, true) {
			flag = true
		}
		if ud.canMingGang(true) || ud.canAnGang(true) {
			flag = true
		}
	}

	//是否可以碰
	if curOp {
		if ud.canPeng(cd, true) {
			flag = true
		}
	}
	return flag
}

//
func (ud *STUserData) resetCanOpInfo() {
	ud.canOpEventInfo = new(STCanOpEventInfo)
	ud.canOpEventInfo.OpInfo = make([]*STOpInfo, 0)
	ud.setEnevtFlag(false)
}

//能否胡牌
func (ud *STUserData) canHu(cd int32) (b bool) {
	// mlog.Debug("检测玩家SeatNo[%d]是否可以胡牌[%d]", ud.User.SeatNo, cd)

	if ud.color == cd/100 {
		//mlog.Debug("error color")
		return false
	}
	cg := ud.handCards.CopyCardGroup()
	if cg == nil {
		//mlog.Debug("error cg")
		return false
	}
	if cd < 100 {
		//mlog.Debug("error cd")
		return false
	}

	for _, card := range cg.Cards {
		if card/100 == ud.color {
			//mlog.Debug("error color2")
			return false
		}
	}
	if cg.CanHuBase(cd) {
		ud.canOpOther = cd
		return true
	}
	// Remarked by Joecool Wang at 2019-06-05
	// 這是判斷可否胡七對子的程序，此處應不需執行此程序
	//cg.AddCard(cd)
	//if cg.IsSevenPair() {
	//	ud.canOpOther = cd
	//	return true
	//}
	// Remarked by Joecool Wang at 2019-06-05 end
	//mlog.Debug("end can't hu")
	return false
}

//能否明杠，自己摸到的牌
func (ud *STUserData) canMingGang(bSave bool) bool {
	if ud.moPai < 100 {
		return false
	}
	nr := false
	for _, pg := range ud.pengCards {
		if !pg.valid {
			continue
		}
		if pg.pengCard == ud.moPai {
			if bSave {
				stOpinfo := new(STOpInfo)
				stOpinfo.nOpType = OP_MINGGANG
				stOpinfo.nValue = append(stOpinfo.nValue, ud.moPai)
				ud.canOpEventInfo.OpInfo = append(ud.canOpEventInfo.OpInfo, stOpinfo)
				mlog.Debug("玩家可以明杠,座位号[%d],摸牌[%d]", ud.moPai, ud.seatNo)
			}
			nr = true
			continue
		}
		for _, cd := range ud.handCards.Cards {
			if cd == pg.pengCard {
				if bSave {
					stOpinfo := new(STOpInfo)
					stOpinfo.nOpType = OP_MINGGANG
					stOpinfo.nValue = append(stOpinfo.nValue, cd)
					ud.canOpEventInfo.OpInfo = append(ud.canOpEventInfo.OpInfo, stOpinfo)
					mlog.Debug("玩家可以明杠,座位号[%d],牌[%d]", cd, ud.seatNo)
				}
				nr = true
				break
			}
		}
	}
	return nr
}

//能否引杠，自己摸到的牌或者其他玩家打出的牌
func (ud *STUserData) canYinGang(cd int32, bSave bool) bool {
	if ud.color == cd/100 {
		return false
	}
	if cd < 100 {
		return false
	}
	if ud.handCards.HaveThreeOf(cd) {
		if bSave {
			stOpinfo := new(STOpInfo)
			stOpinfo.nOpType = OP_YINGANG
			stOpinfo.nValue = append(stOpinfo.nValue, cd)
			ud.canOpEventInfo.OpInfo = append(ud.canOpEventInfo.OpInfo, stOpinfo)
			mlog.Debug("玩家座位号[%d]可以引杠", ud.seatNo)
		}
		return true
	}
	return false
}

//能否暗杠，自己摸到的牌
func (ud *STUserData) canAnGang(bSave bool) bool {
	if ud.handCards.HaveThreeOf(ud.moPai) {
		if ud.color == ud.moPai/100 {
			return false
		}
		if bSave {
			stOpinfo := new(STOpInfo)
			stOpinfo.nOpType = OP_ANGANG
			stOpinfo.nValue = append(stOpinfo.nValue, ud.moPai)
			ud.canOpEventInfo.OpInfo = append(ud.canOpEventInfo.OpInfo, stOpinfo)
			mlog.Debug("玩家座位号[%d]可以暗杠", ud.seatNo)
		}
		return true
	}
	anCg := ud.handCards.HaveFour()
	if ud.moPai != 0 && len(anCg) > 0 {
		anReal := make([]int32, 0)
		for _, cd := range anCg {
			if ud.color == cd/100 {
				continue
			}
			anReal = append(anReal, cd)
		}
		if len(anReal) <= 0 {
			return false
		}
		if bSave {
			stOpinfo := new(STOpInfo)
			stOpinfo.nOpType = OP_ANGANG
			stOpinfo.nValue = append(stOpinfo.nValue, anReal...)
			ud.canOpEventInfo.OpInfo = append(ud.canOpEventInfo.OpInfo, stOpinfo)
			mlog.Debug("玩家座位号[%d]可以暗杠", ud.seatNo)
		}
		return true
	}
	return false
}

//能否碰
func (ud *STUserData) canPeng(cd int32, bSave bool) bool {
	if ud.color == cd/100 {
		return false
	}
	if cd < 100 {
		return false
	}
	if ud.handCards.HavePairOf(cd) {
		if bSave {
			ud.canOpOther = cd
			stOpinfo := new(STOpInfo)
			stOpinfo.nOpType = OP_PENG
			stOpinfo.nValue = append(stOpinfo.nValue, cd)
			ud.canOpEventInfo.OpInfo = append(ud.canOpEventInfo.OpInfo, stOpinfo)
			mlog.Debug("玩家座位号[%d]可以碰", ud.seatNo)
		}

		return true
	}
	return false
}

func (ud *STUserData) autoUserOp() {
}

func (ud *STUserData) findOpInfo(cd, opType int32) int32 {
	for _, data := range ud.canOpEventInfo.OpInfo {
		for _, value := range data.nValue {
			if cd == value && (opType == data.nOpType || opType == data.nOpType/10) {
				return data.nOpType
			}
		}
	}
	return 0
}

func (ud *STUserData) saveInTimeRecord(saveType, cd, rate int32, score int64, seatNos, huTypes []int32) {
	//类型 类型详情 相关人员 得分
	recordLos := new(STScoreRecord)
	recordLos.Score = score
	recordLos.SeatNo = seatNos
	recordLos.Type = saveType
	recordLos.valid = true
	recordLos.HuTypes = huTypes
	recordLos.card = cd
	recordLos.Rate = rate
	ud.intimeInfo = append(ud.intimeInfo, recordLos)
}

// 自动出牌
func (ud *STUserData) calBestOutCard() int32 {
	opValue := int32(0)
	//先定缺
	if ud.moPai/100 == ud.color {
		opValue = ud.moPai
	} else {
		for _, cd := range ud.handCards.Cards {
			if cd/100 == ud.color {
				opValue = cd
				break
			}
		}
	}
	if opValue > 0 {
		return opValue
	}
	handAll := new(mahjong.CardGroup)
	for _, cd := range ud.handCards.Cards {
		handAll.AddCard(cd)
	}
	if ud.moPai > 0 {
		handAll.AddCard(ud.moPai)
	}
	// 选取最优走法
	needMin := int32(20)
	for _, deCd := range handAll.Cards {
		hands := new(mahjong.CardGroup)
		for _, cd := range ud.handCards.Cards {
			hands.AddCard(cd)
		}
		if ud.moPai > 0 {
			hands.AddCard(ud.moPai)
		}
		hands.DeleteCard(deCd, 1)
		mapCard := make(map[int32]int32, 0)
		for _, cd := range hands.Cards {
			mapCard[cd/10]++
		}
		x := mahjong.NeedLaizi(mapCard)
		if needMin > x {
			needMin = x
			opValue = deCd
		}
	}
	// 重新遍历，选取所有最优走法
	bestCardList := make([]int32, 0)
	for _, deCd := range handAll.Cards {
		hands := new(mahjong.CardGroup)
		for _, cd := range ud.handCards.Cards {
			hands.AddCard(cd)
		}
		if ud.moPai > 0 {
			hands.AddCard(ud.moPai)
		}
		hands.DeleteCard(deCd, 1)
		mapCard := make(map[int32]int32, 0)
		for _, cd := range hands.Cards {
			mapCard[cd/10]++
		}
		x := mahjong.NeedLaizi(mapCard)
		if needMin == x {
			bestCardList = append(bestCardList, deCd)
		}
	}
	// 计算所选牌权值，选择更优走法
	minW := int32(999)
	for _, bcl := range bestCardList {

		w := ud.calWeight(bcl)
		// mlog.Debug("bcl[%d],w[%d]", bcl, w)
		if minW > w {
			minW = w
			opValue = bcl
		}
	}
	if opValue <= 0 {
		opValue = ud.handCards.Cards[0]
	}
	return opValue
}

//计算权值
func (ud *STUserData) calWeight(bcl int32) int32 {
	bcltmp := bcl
	hands := new(mahjong.CardGroup)
	for _, cd := range ud.handCards.Cards {
		hands.AddCard(cd)
	}
	if ud.moPai > 0 {
		hands.AddCard(ud.moPai)
	}
	hands.DeleteCard(bcltmp, 1)
	mapCard := make(map[int32]int32, 0)
	for _, cd := range hands.Cards {
		mapCard[cd/10]++
	}
	w := int32(0)
	bcltmp = bcltmp / 10
	if mapCard[bcltmp] >= 2 {
		w += 30
	} else if mapCard[bcltmp] == 1 {
		w += 20
	}
	if mapCard[bcltmp+1] > 0 && mapCard[bcltmp+2] <= 0 {
		w += 10
	} else if mapCard[bcltmp+1] > 0 && mapCard[bcltmp+2] > 0 {
		w += 15
	} else if mapCard[bcltmp+2] > 0 {
		w += 5
	}
	if mapCard[bcltmp-1] > 0 && mapCard[bcltmp-2] <= 0 {
		w += 10
	} else if mapCard[bcltmp-1] > 0 && mapCard[bcltmp-2] > 0 {
		w += 15
	} else if mapCard[bcltmp-2] > 0 {
		w += 5
	}
	bcltmp = bcltmp % 10
	w -= (5 - bcltmp) * (5 - bcltmp) / 2
	return w
}

// 自动定缺
func (ud *STUserData) autoDingque() int32 {
	min, color := int32(15), int32(1)
	for i := mahjong.ColorTong; i <= mahjong.ColorWan; i++ {
		count := ud.handCards.CalColorCount(i)
		if min > count {
			min = count
			color = i
		}
	}
	return color
}

// 玩家操作后设置标志位，防止断连后重复提示
func (ud *STUserData) setEnevtFlag(flag bool) {
	ud.eventFlag = flag
}

// 检测玩家是否为相公牌
func (ud *STUserData) checkErrorHandSize() {
	size := ud.handCards.Len()
	if size%3 != 1 {
		mlog.Error("游戏逻辑出错,玩家为相公牌,SeatNo[%d],size[%d]", ud.seatNo, size)
	}
}

//玩家是否可以抢杠
func (ud *STUserData) updataCanQiang(cd int32) bool {
	if ud.isHu {
		return false
	}
	flag := false
	//是否可以胡,点炮胡，抢杠胡，自摸胡
	if ud.canHu(cd) {
		stOpinfo := new(STOpInfo)
		stOpinfo.nOpType = OP_HU
		stOpinfo.nValue = append(stOpinfo.nValue, cd)
		ud.canOpEventInfo.OpInfo = append(ud.canOpEventInfo.OpInfo, stOpinfo)
		mlog.Debug("玩家座位号[%d]可以胡", ud.seatNo)
		flag = true
	}

	return flag
}
