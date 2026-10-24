package ddz

import (
	"fmt"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"math"
	"sort"
)

//
var cardMap = []int32{0, 0, 20, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 21, 22, 23}

//斗地主玩家倍率选择
type userDdzOpt struct {
	userId       int32         //玩家id
	money        int32         //结算输赢金币
	userSeat     int32         //玩家座位号
	JiaoDiZhu    JiaoDiZhuEnum //是否叫牌
	isQiangDizhu int32         //是否抢地主 0无状态  1不抢 2抢
	isMingPai    bool          //是否名牌
	isSendCard   bool          //是否出过牌
	isLastDiZhu  bool          //上一把是否是地主
	isTuoGuan    bool          //是否托管
	isDo         bool          //用于判断当前阶段用户是否操作
	isDiZhu      bool          //是否是地主
	isWin        bool          //输赢
	lt           loseEnum      //输的类型
	action       *userAction   //用户当前动作
}

type userAction struct {
	userSeat int32
	cg       *poker.CardGroup
	ct       cardType
	mCard    int32 //主牌xx副
	maxCard  int32 //主牌点数
}

//倍率变化记录
type record struct {
	userId int32
	ot     netproto.OddsChangeType
	odds   int32
}

type osr struct {
	rs []*record
}

func (o *osr) push(r *record) {
	o.rs = append(o.rs, r)
}

func (o *osr) CaculateOdds(calMethod int) map[OddsType]int32 {
	var calfunc func(org int32, add int32) int32
	m := make(map[OddsType]int32)
	if calMethod == 0 { //決定倍率算法
		calfunc = o.plus
	} else {
		calfunc	= o.multiply
	}

	for _, v := range o.rs {
		if v.ot == netproto.OddsChangeType_null {
			m[DdzOddsTypeBase] = calfunc(m[DdzOddsTypeBase], v.odds)
		} else if v.ot == netproto.OddsChangeType_jiaodizhu {
			m[DdzOddsTypeJdz] = calfunc(m[DdzOddsTypeJdz], v.odds)
		} else if v.ot == netproto.OddsChangeType_qiangdizhu {
			m[DdzOddsTypeQdz] = calfunc(m[DdzOddsTypeQdz], v.odds)
		} else if v.ot == netproto.OddsChangeType_mingpai {
			m[DdzOddsTypeMingPai] = calfunc(m[DdzOddsTypeMingPai], v.odds)
		} else if v.ot == netproto.OddsChangeType_zhadan {
			m[DdzOddsTypeZhaDan] = calfunc(m[DdzOddsTypeZhaDan], v.odds)
		} else if v.ot == netproto.OddsChangeType_chuntian {
			m[DdzOddsTypeChunTian] = calfunc(m[DdzOddsTypeChunTian], v.odds)
		} else if v.ot == netproto.OddsChangeType_rocket {
			m[DdzOddsTypeRocket] = calfunc(m[DdzOddsTypeRocket], v.odds)
		} else if v.ot == netproto.OddsChangeType_fanchuntian {
			m[DdzOddsTypeFanChunTian] = calfunc(m[DdzOddsTypeFanChunTian], v.odds)
		} else {
			m[DdzOddsTypeDiZhuPai] = calfunc(m[DdzOddsTypeDiZhuPai], v.odds)
		}
	}
	return m
}

func (o *osr) plus(org int32, add int32) int32 {
	return org + add
}

func (o *osr) multiply(org int32, add int32) int32 {
	if org == 0 {
		org = 1
	}
	return org * add
}


func (ua *userAction) IsBig(a *userAction) bool {
	l := int32(len(cardMap))
	if cardMap[ua.maxCard%l] > cardMap[a.maxCard%l] {
		return true
	}
	return false
}

type history struct {
	actions []*userAction
}

func doubleCheck(cg *poker.CardGroup, ctx *sortCtx, ua *userAction) cardType {
	if !cg.Cards[0].IsSameCardNumber(cg.Cards[1]) {
		if cg.Cards[0].IsKing() && cg.Cards[1].IsKing() {
			ua.maxCard = ctx.GetMainCardMax()
			ua.mCard = 1
			return cardTypeRocket
		} else {
			return cardTypeNone
		}
	} else {
		ua.maxCard = ctx.GetMainCardMax()
		ua.mCard = 1
		return cardTypeDouble
	}
}

func threeCheck(cg *poker.CardGroup, ctx *sortCtx, ua *userAction) cardType {
	l := len(cg.Cards)
	for i := 1; i < l; i++ {
		if !cg.Cards[i-1].IsSameCardNumber(cg.Cards[i]) {
			return cardTypeNone
		}
	}
	ua.maxCard = ctx.GetMainCardMax()
	ua.mCard = 1
	return cardTypeThree
}

func fourCheck(cg *poker.CardGroup, ctx *sortCtx, ua *userAction) cardType {
	m := make(map[int32]int32)
	l := len(cg.Cards)
	for i := 0; i < l; i++ {
		m[cg.Cards[i].CardNumber] += 1
	}
	for _, v := range m {
		if v == 3 {
			ua.maxCard = ctx.GetMainCardMax()
			ua.mCard = 1
			return cardTypeThreeAndOne
		}
		if v == 4 {
			ua.maxCard = ctx.GetMainCardMax()
			ua.mCard = 1
			return cardTypeBomb
		}
	}
	return cardTypeNone
}

func fiveCheck(cg *poker.CardGroup, ctx *sortCtx, ua *userAction) cardType {

	sort.Sort(cg)
	for _, cd := range cg.Cards {
		if cd.CardNumber == 15 || cd.CardNumber == 16 {
			return cardTypeNone
		}
	}
	l := len(cg.Cards)
	m := make(map[int32]int32)
	for i := 1; i < l; i++ {
		if math.Abs(float64(cg.Cards[i-1].CardNumber-cg.Cards[i].CardNumber)) != 1 && (cg.Cards[i-1].CardNumber == 2 || cg.Cards[i].CardNumber == 2) {
			for n := 0; n < l; n++ {
				m[cg.Cards[n].CardNumber] += 1
			}
			return cardTypeNone
		}
	}
	ua.maxCard = ctx.GetMainCardMax()
	ua.mCard = 1
	return cardTypeStraight
}

func otherCheckShuangShun(ctx *sortCtx, ua *userAction) cardType {
	if ctx.Len() <= 2 { //无效牌
		return cardTypeNone
	}
	scn := sortCtxByNumber{*ctx}
	sort.Sort(&scn)
	i := 0
	index := -1
	for {
		if i >= scn.Len() {
			break
		}
		if scn.GetByIndex(i) == 2 && scn.GetValueByIndex(i) != 2 {
			if index != -1 {
				if math.Abs(float64(scn.GetValueByIndex(i)-scn.GetValueByIndex(index))) != 1 {
					mlog.Debug("value1=%d,value2=%d", scn.GetValueByIndex(i), scn.GetValueByIndex(index))
					return cardTypeNone
				}
			}
			index = i
		}
		i++
	}
	ua.maxCard = scn.GetMainCardMax()
	ua.mCard = int32(scn.Len())
	return cardTypeThreeStraight
}

func otherCheckSanShun(ctx *sortCtx, ua *userAction) cardType {
	scn := sortCtxByNumber{*ctx}
	sort.Sort(&scn)
	i := 0
	index := -1
	for {
		if i >= scn.Len() {
			break
		}
		if scn.GetByIndex(i) == 3 && scn.GetValueByIndex(i) != 2 {
			if index != -1 {
				if math.Abs(float64(scn.GetValueByIndex(i)-scn.GetValueByIndex(index))) != 1 {
					return cardTypeNone
				}
			}
			index = i
		}
		i++
	}
	ua.maxCard = scn.GetMainCardMax()
	ua.mCard = int32(scn.Len() / 3)
	return cardTypeThreeStraight
}

//飞机带翅膀 带单情况
func otherCheckFeiJiDaiChibangDan(l int, ctx *sortCtx, ua *userAction) cardType {
	scn := sortCtxByNumber{*ctx}
	sort.Sort(&scn)
	i, k, n, index := 0, 1, l/4, -1
	if scn.CdiCount(3) == (l/4)+1 {
		for {
			if i >= scn.Len() {
				break
			}
			if scn.GetByIndex(i) == 3 {
				if index != -1 {
					if math.Abs(float64(scn.GetValueByIndex(i)-scn.GetValueByIndex(index))) != 1 {
						ua.maxCard = scn.GetValueByIndex(index)
						break
					} else {
						k++
						if k == n {
							ua.maxCard = scn.GetValueByIndex(index)
						}
					}
				}
				index = i
			}
			i++
		}
		if k < n {
			return cardTypeNone
		}
		ua.mCard = int32(n)
		return cardTypeDoubleStraights
	}
	for {
		if i >= scn.Len() {
			break
		}
		if scn.GetValueByIndex(i) != 2 && scn.GetByIndex(i) > 2 {
			if index != -1 {
				if math.Abs(float64(scn.GetValueByIndex(i)-scn.GetValueByIndex(index))) != 1 {
					ua.maxCard = scn.GetValueByIndex(index)
					break
				} else {
					k++
					if k == n {
						ua.maxCard = scn.GetValueByIndex(index)
					}
				}
			}
			index = i
		}
		i++
	}
	if k < n {
		return cardTypeNone
	}
	ua.mCard = int32(n)
	mlog.Debug("==================mCard = %d,maxCard=%d", ua.mCard, ua.maxCard)
	return cardTypeDoubleStraights
}

//飞机带翅膀 带双情况
func otherCheckFeiJiDaiChibangShuang(l int, ctx *sortCtx, ua *userAction) cardType {
	scn := sortCtxByNumber{*ctx}
	sort.Sort(&scn)

	if scn.CdiCount(3)*3+scn.CdiCount(2)*2+scn.CdiCount(4)*4 == l {
		i := 0
		index := -1
		for {
			if i >= scn.Len() {
				break
			}
			if scn.GetByIndex(i) == 3 && scn.GetValueByIndex(i) != 2 {
				if index != -1 {
					if math.Abs(float64(scn.GetValueByIndex(i)-scn.GetValueByIndex(index))) != 1 {
						return cardTypeNone
					}
				}
				index = i
			}
			i++
		}

		sort.Sort(ctx)
		ua.maxCard = scn.GetMainCardMax()
		ua.mCard = int32(scn.CdiCount(3))

		if ua.mCard == int32(0) {
			return cardTypeNone
		}

		return cardTypeDoubleStraightd
	}
	return cardTypeNone
}

func otherCheckSiDaiEr(l int, ctx *sortCtx, ua *userAction) cardType {
	if ctx.CdiCount(4)*2 == ctx.CdiCount(2) || ctx.CdiCount(4) == ctx.CdiCount(2) || ctx.CdiCount(4)*4 == l {
		ua.maxCard = ctx.GetMainCardMax()
		ua.mCard = int32(ctx.CdiCount(4))
		return cardTypeFourAndTwod
	}
	if ctx.CdiCount(4)*2 == ctx.CdiCount(1) {
		ua.maxCard = ctx.GetMainCardMax()
		ua.mCard = int32(ctx.CdiCount(4))
		return cardTypeFourAndTwos
	}
	return cardTypeNone
}

func otherCheck(cg *poker.CardGroup, ctx *sortCtx, ua *userAction) cardType {
	l := len(cg.Cards)
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("otherCheck ctx[%v]", ctx)).Debug("託管")
	if ctx.GetByIndex(0) == 2 && ctx.CountCondition(2) { //顺对 这里还需要排除2点及双王 未做
		return otherCheckShuangShun(ctx, ua)
	}
	if ctx.GetByIndex(0) == 3 && ctx.Len() >= 2 && ctx.CountCondition(3) {
		ret := otherCheckSanShun(ctx, ua)
		if ret == cardTypeNone {
			ret = otherCheckFeiJiDaiChibangDan(l, ctx, ua)
		}
		return ret
	}
	if ctx.GetByIndex(0) == 3 && ctx.Len() == 2 && l == 5 {
		ua.maxCard = ctx.GetMainCardMax()
		ua.mCard = int32(ctx.CdiCount(2))
		return cardTypeThreeAndOne
	}
	if l%4 == 0 && l > 4 && !ctx.CountCondition(1) {
		if ctx.CdiCount(4) == 2 && l == 8 || (ctx.CdiCount(4) == 1 && ctx.CdiCount(2) == 2) {
			return otherCheckSiDaiEr(l, ctx, ua)
		}
		return otherCheckFeiJiDaiChibangDan(l, ctx, ua)
	}
	if l%5 == 0 && !ctx.CountCondition(1) {
		if ctx.Len()%2 == 0 && ctx.GetByIndex(0) == 3 { //飞机带翅膀  带双牌的情况
			return otherCheckFeiJiDaiChibangShuang(l, ctx, ua)
		}
	}
	if ((ctx.Len()%2 == 0 || ctx.Len()%3 == 0) && ctx.GetByIndex(0) == 4 && l%6 == 0) || (ctx.Len()%3 == 0 && ctx.GetByIndex(0) == 4 && l%8 == 0) { //四带二
		return otherCheckSiDaiEr(l, ctx, ua)
	}
	if ctx.Len() >= 5 {
		return fiveCheck(cg, ctx, ua)
	}
	return cardTypeNone
}

func (h *history) cardTypeCheck(cg *poker.CardGroup, ua *userAction) cardType {
	sort.Sort(cg)
	l := len(cg.Cards)
	ctx := new(sortCtx)
	ctx.arr = make([]*sItem, 0)
	m := make(map[int32]int32)
	for i := 0; i < l; i++ {
		m[cg.Cards[i].CardNumber] += 1
	}
	for k, v := range m {
		item := new(sItem)
		item.number = k
		item.count = v
		ctx.arr = append(ctx.arr, item)
	}
	sort.Sort(ctx)
	if l == 1 {
		ua.maxCard = ctx.GetMainCardMax()
		ua.mCard = 1
		return cardTypeSingle
	} else if l == 2 {
		return doubleCheck(cg, ctx, ua)
	} else if l == 3 {
		return threeCheck(cg, ctx, ua)
	} else if l == 4 {
		return fourCheck(cg, ctx, ua)
	} else {
		return otherCheck(cg, ctx, ua)
	}
}

//将用户出牌放入出牌记录中
func (h *history) Push(ua *userAction) bool {
	h.actions = append(h.actions, ua)
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("userSeat[%v] 将用户出牌放入出牌记录中 len(h.actions)[%v] h.actions[%v] ",
		ua.userSeat, len(h.actions), &h.actions)).Debug("託管")
	return false
}

//获取最近一条记录
func (h *history) GetBack() *userAction {
	l := len(h.actions)
	if l <= 0 {
		return nil
	}
	return h.actions[l-1]
}

type userExtInfo struct {
	userOpts       []*userDdzOpt    //玩家选项
	turn           int32            //当前起始位置
	curTurn        int32            //当前局数
	baseScore      int32            //底分
	his            *history         //出牌记录
	or             *osr             //倍率记录
	lastQiangDiZhu int32            //最后一个抢地主的玩家的userId
	odds           int32            //当前所有的倍率综合
	isChunTian     bool             //是否是春天
	isFanChunTian  bool             //莊家是否被春天
	dizhuCards     *poker.CardGroup //地主牌
	difen          int32            //底分
	timeOut        int32            //一局完了后  客户端等待超时时间
	readytime      int32            //准备时间
	opttime        int32            //操作时间
	maxodd         int32            //最大倍数
	bindaccount    bool             //綁定帳號 0關閉 1開啟
	payminmoney    int32            //充值金額大於等於指定的數值
}

//初始化玩家信息
func (ue *userExtInfo) init(t *gamelogic.Table) {
	for _, v := range t.GetTableUser() {
		ue.userOpts[v.SeatNo] = new(userDdzOpt)
		ue.userOpts[v.SeatNo].JiaoDiZhu = JDZNone
		ue.userOpts[v.SeatNo].isMingPai = false
		ue.userOpts[v.SeatNo].isQiangDizhu = 0
		ue.userOpts[v.SeatNo].isSendCard = false
		ue.userOpts[v.SeatNo].isTuoGuan = false
		ue.userOpts[v.SeatNo].isDo = false
		ue.userOpts[v.SeatNo].isLastDiZhu = false
		ue.userOpts[v.SeatNo].isDiZhu = false
		ue.userOpts[v.SeatNo].userId = v.UserID
		ue.userOpts[v.SeatNo].userSeat = v.SeatNo
		ue.userOpts[v.SeatNo].isWin = false
		ue.userOpts[v.SeatNo].lt = loseEnumNormal
		ue.userOpts[v.SeatNo].money = 0
		//ue.userOpts[v.SeatNo].action = new(userAction)
	}
	ue.isChunTian = true
	ue.isFanChunTian = true
	ue.his = new(history)
	ue.his.actions = make([]*userAction, 0)
	ue.or = new(osr)
	ue.or.rs = make([]*record, 0)
	ue.dizhuCards = new(poker.CardGroup)
}

//重置玩家信息
func (ue *userExtInfo) reset() {
	for i := 0; i < len(ue.userOpts); i++ {
		ue.userOpts[i].JiaoDiZhu = JDZNone
		ue.userOpts[i].isMingPai = false
		ue.userOpts[i].isQiangDizhu = 0
		ue.userOpts[i].isSendCard = false
		ue.userOpts[i].isTuoGuan = false
		ue.userOpts[i].isDiZhu = false
		ue.userOpts[i].isDo = false
		ue.userOpts[i].isWin = false
		ue.userOpts[i].lt = loseEnumNormal
		ue.userOpts[i].money = 0
		//ue.userOpts[i].action = new(userAction)
	}
	if ue.his != nil {
		ue.his.actions = make([]*userAction, 0)
	}
	if ue.or != nil {
		ue.or.rs = make([]*record, 0)
	}
	ue.dizhuCards = new(poker.CardGroup)
	ue.isChunTian = true
	ue.isFanChunTian = true
	ue.odds = 0
	ue.turn = ue.lastQiangDiZhu
}

func (ue *userExtInfo) getLastDiZhu() (int, *userDdzOpt) {
	for i, user := range ue.userOpts {
		if user.userSeat == ue.lastQiangDiZhu {
			return i, user
		}
	}
	return 0, nil
}

func (ue *userExtInfo) resetLastDiZhu() {
	for _, user := range ue.userOpts {
		user.isLastDiZhu = false
	}
	return
}

func (ue *userExtInfo) getDiZhu() (int, *userDdzOpt) {
	for i, user := range ue.userOpts {
		if user.isDiZhu == true {
			return i, user
		}
	}
	return 0, nil
}

func (ue *userExtInfo) getJiaoDiZhu() (int, *userDdzOpt) {
	for i, user := range ue.userOpts {
		if user.JiaoDiZhu == JDZJiao {
			return i, user
		}
	}
	return 0, nil
}

func (ue *userExtInfo) getExtInfoByTurn(trun int32) *userDdzOpt {
	//mlog.NewDefaultLogger().SetData(fmt.Sprintf("trun[%v] 取得下個出牌的位置 座位[%v]", trun, trun%int32(len(ue.userOpts)) )).Debug("託管")
	return ue.userOpts[trun%int32(len(ue.userOpts))]
}

func (ue *userExtInfo) getNowSeat() int32 {
	return ue.turn % int32(len(ue.userOpts))
}
