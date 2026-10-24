package sss

import (
	"math"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/util"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

const (
	AnimotionFaPai       = "fapai"
	AnimotionHeadComb    = "toudunbipai"
	AnimotionBodyComb    = "zhongdunbipai"
	AnimotionFootComb    = "houdunbipai"
	AnimotionSpecialComb = "teshubipai"
	AnimotionGun         = "daqiang"
	AnimotionEnd         = "end"
	AnimotinBeginCompare = "ksbpdh" //开始比牌动画
)

type SSSDunType int32

const (
	SSSDunHead SSSDunType = iota
	SSSDunBody
	SSSDunFoot
	SSSDunMax
)

type CardGroupSortNumber struct {
	*poker.CardGroup
	dealSpecialA bool
}

func (cg *CardGroupSortNumber) Swap(i, j int) {
	cg.CardGroup.Swap(i, j)
}

func (cg *CardGroupSortNumber) Less(i, j int) bool {
	if cg.dealSpecialA {
		return getCardNumber(cg.CardGroup.Cards[i].CardNumber) > getCardNumber(cg.CardGroup.Cards[j].CardNumber)
	} else {
		return cg.CardGroup.Cards[i].CardNumber > cg.CardGroup.Cards[j].CardNumber
	}
}

func (cg *CardGroupSortNumber) Oderbynumber() {
	for i := 0; i < cg.Len(); i++ {
		for j := 0; j+1 < cg.Len()-i; j++ {
			if cg.desc(j, j+1) {
				cg.Swap(j, j+1)
			}
		}
	}
}

func (cg *CardGroupSortNumber) desc(i, j int) bool {
	if cg.dealSpecialA {
		return getCardNumber(cg.CardGroup.Cards[i].CardNumber) > getCardNumber(cg.CardGroup.Cards[j].CardNumber)
	} else {
		return cg.CardGroup.Cards[i].CardNumber > cg.CardGroup.Cards[j].CardNumber
	}
}

type NResultItem struct {
	Ct     netproto.SSSDunCardsNomalType
	result map[int32]int32 //key:seat value:result
}

func NewNResultItem() []*NResultItem {
	nnr := make([]*NResultItem, SSSDunMax)
	for i := 0; i < int(SSSDunMax); i++ {
		nnr[i] = new(NResultItem)
		nnr[i].result = make(map[int32]int32)
	}
	return nnr
}

type SResultItem struct {
	Ct   netproto.SSSDunCardsSpecialType
	Win  int32
	Lose int32
}

type WinInfo struct {
	seatNo int32 //我赢的人的座位号
	win    int32 //赢的水数
}

type AutoCardType struct {
	cgs         []*poker.CardGroup
	ct          []int32 //自动获取的牌型 是否是特殊牌型  netproto.SSSDunCardsSpecialType
	recommandCt *netproto.HandCards
	isSpecial   bool
	sct         int32 //netproto.SSSDunCardsSpecialType
	sctWin      int32 //特殊牌型赢的水数
	sctLose     int32 //特殊牌型输的水数
	sctWinInfo  []*WinInfo
}

func NewAutoCardType() *AutoCardType {
	act := new(AutoCardType)
	act.cgs = make([]*poker.CardGroup, SSSDunMax)
	act.ct = make([]int32, SSSDunMax)
	act.recommandCt = new(netproto.HandCards)
	act.sctWinInfo = make([]*WinInfo, 0)
	return act
}

func (act *AutoCardType) Cg2Dun(dun *netproto.Dun, t int32, count int) {
	if act.cgs == nil {
		mlog.Debug("act.cgs=nil")
		return
	}
	//mlog.Debug("act.cgs=%v", act.cgs)
	//if count > len(act.cgs[t].Cards) || act.cgs[t].Cards == nil {
	//	mlog.Debug("when count > cards.len, cards => %v, count => %v", act.cgs[t].Cards, count)
	//}
	for i := 0; i < count; i++ {
		c := new(netproto.SCard)
		c.Number = proto.Int32(act.cgs[t].Cards[i].CardNumber)
		c.Color = proto.Int32(act.cgs[t].Cards[i].Color)
		dun.Cards = append(dun.Cards, c)
	}
}

func (act *AutoCardType) Cg2Hc() {
	act.recommandCt.Head = new(netproto.Dun)
	act.recommandCt.Head.Ct = proto.Int32(act.ct[0])
	act.Cg2Dun(act.recommandCt.Head, 0, 3)
	act.recommandCt.Body = new(netproto.Dun)
	act.recommandCt.Body.Ct = proto.Int32(act.ct[1])
	act.Cg2Dun(act.recommandCt.Body, 1, 5)
	act.recommandCt.Foot = new(netproto.Dun)
	act.recommandCt.Foot.Ct = proto.Int32(act.ct[2])
	act.Cg2Dun(act.recommandCt.Foot, 2, 5)
}

func (act *AutoCardType) Hc2Cgs() {
	hcg := poker.NewCardGroup(nil)
	for _, v := range act.recommandCt.Head.Cards {
		c := poker.NewPokerCard(v.GetColor(), v.GetNumber())
		hcg.AddCard(c)
	}
	act.ct[0] = act.recommandCt.Head.GetCt()
	bcg := poker.NewCardGroup(nil)
	for _, v := range act.recommandCt.Body.Cards {
		c := poker.NewPokerCard(v.GetColor(), v.GetNumber())
		bcg.AddCard(c)
	}
	act.ct[1] = act.recommandCt.Body.GetCt()
	fcg := poker.NewCardGroup(nil)
	for _, v := range act.recommandCt.Foot.Cards {
		c := poker.NewPokerCard(v.GetColor(), v.GetNumber())
		fcg.AddCard(c)
	}
	act.ct[2] = act.recommandCt.Foot.GetCt()
}

//角色信息
type SSSRole struct {
	roleid  int32
	seat    int32
	isDo    bool
	ris     []*NResultItem
	gun     []int32 //被我打枪的玩家座位号
	act     *AutoCardType
	isZuPai bool
}

type callBack func()

type SSSTable struct {
	*gamelogic.DynamicGameLogicBase
	key          bool
	gameover     bool
	turn         int32
	step         netproto.SSSGameStep
	jump         netproto.SSSGameStep
	maxroles     int32
	difen        int32
	roles        []*SSSRole
	handlers     []callBack
	poker        poker.Poker
	gamePoker    *poker.GamePoker
	stepTimer    *event.TimerFunc
	animotionCur string
	timeOut      int32
	zupaitime    int32 //组牌时间
	//specialcount       int32 //有特殊牌型的人数统计
	//choosespecialcount int32 //选择了特殊牌型的人统计
}

func NewSSSGameLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *SSSTable {
	glc := new(SSSTable)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)
	glc.gamePoker = poker.NewGamePoker() //创建游戏牌型
	gsvrcfg := cfgs.GetConfig("gameserver")
	glc.timeOut = int32(gsvrcfg.GetConfigInt("gameserver", "readywaittime"))
	glc.difen = int32(gsvrcfg.GetConfigInt("gameserver", "basescore"))
	glc.zupaitime = int32(gsvrcfg.GetConfigInt("gameserver", "zupaitime"))
	glc.RegisterAfterGameStartFunc(func() {
		glc.init(t)
		glc.regist()
		glc.logicEngine()
	})

	glc.RegisterBeforeGameStartFunc(func() {
		if glc.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
			glc.poker = poker.NewRandomPoker(false, false, 1)
			glc.poker = poker.NewFixedPoker(glc.AllTableUser(), glc.poker.(*poker.RandomPoker), 13)
		} else { //发牌器发牌
			glc.poker = poker.NewFilePoker(glc.Configs, glc.AllTableUser())
			// glc.poker = poker.NewFixedFilePoker(glc.poker.(*poker.FilePoker))
			glc.poker = poker.NewFixedFileFillPoker(glc.poker.(*poker.FilePoker), glc.AllTableUser())
		}

		glc.gamePoker = poker.NewGamePoker() //创建游戏牌型
	})

	glc.RegisterAfterGameEndFunc(func() {
		glc.clearData()
	})
	return glc
}

func (s *SSSTable) timeOutCallBack() {
	if s.handlers[s.step%(netproto.SSSGameStep_SSSGameStepEnd+1)] != nil {
		s.handlers[s.step%(netproto.SSSGameStep_SSSGameStepEnd+1)]()
		//在此step進行下一步
		s.moveStep()
		// fmt.Println(fmt.Sprintf("现在是阶段：%v", s.step))
		mlog.Debug("现在是阶段：%v", s.step)
		if s.step == s.jump {
			// fmt.Println(fmt.Sprintf("跳过的阶段：%v", s.step))
			mlog.Debug("跳过的阶段：%v", s.step)
			s.jump = netproto.SSSGameStep_SSSGameStepNull
			s.timeOutCallBack()
			//s.moveStep()
		} else {
			s.logicEngine()
		}
		//s.logicEngine()
	}
}

func (s *SSSTable) RegistCallBackFnc(step netproto.SSSGameStep, f callBack) {
	s.handlers[step%(netproto.SSSGameStep_SSSGameStepEnd+1)] = f
}

func (s *SSSTable) moveStep() {
	s.step = (s.step + 1) % (netproto.SSSGameStep_SSSGameStepEnd + 1)
}

func (s *SSSTable) SetJump(step netproto.SSSGameStep) {
	s.jump = step
}

func (s *SSSTable) moveTag(step netproto.SSSGameStep) {
	s.step = step % (netproto.SSSGameStep_SSSGameStepEnd + 1)
}

func (s *SSSTable) stopLoop() {
	s.gameover = true
	if s.stepTimer != nil {
		s.stepTimer.Stop()
		s.stepTimer = nil
	}
}

func (s *SSSTable) restartLoop() {
	s.gameover = false
}

func (s *SSSTable) logicEngine() {
	if s.stepTimer != nil {
		s.stepTimer.Stop()
	}
	if s.gameover {
		return
	}

	s.resetRoleDo()

	timeOut := s.getStepTimeOut()
	s.broadcastOptTime()
	s.stepTimer = s.TimerFuncManager.NewNamedTimer(timeOut, s.timeOutCallBack, "")
	if s.step == netproto.SSSGameStep_SSSGameStepCombine {
		for _, v := range s.roles {
			if v.isZuPai {
				user := s.GetUserBySeatNo(v.seat)
				s.broadcastUserCardTypeChoose(user, v.act.isSpecial)
				v.isDo = true
			}
		}

		if s.CheckSpecailDo() {
			s.timeOutCallBack()
		}
	}
}

func (s *SSSTable) checkAllRoleDo() bool {
	for _, v := range s.roles {
		if !v.isDo {
			return false
		}
	}
	return true
}

func (s *SSSTable) resetRoleDo() {
	for _, v := range s.roles {
		v.isDo = false
	}
}

//初始化游戏数据
func (s *SSSTable) init(table *gamelogic.Table) {
	// s.handsTable = new(HandsTable)
	// s.handsTable.Init()
	s.gameover = false
	s.step = netproto.SSSGameStep_SSSGameStepBegin
	s.jump = netproto.SSSGameStep_SSSGameStepNull
	s.roles = make([]*SSSRole, 0)
	s.handlers = make([]callBack, netproto.SSSGameStep_SSSGameStepEnd+1)
	for i := 0; i < int(table.Room.SeatCount); i++ {
		v := s.GetUserBySeatNo(int32(i))
		if v == nil {
			mlog.Debug("init role info faild seat=%d", i)
			continue
		}

		role := new(SSSRole)
		role.roleid = v.UserID
		role.seat = v.SeatNo
		role.act = nil
		role.ris = NewNResultItem()
		role.isZuPai = false
		s.roles = append(s.roles, role)
	}
	// fmt.Println("SSSTable.init()")
}

//根据座位号获取角色信息
func (s *SSSTable) getRoleInfoBySeat(seat int32) *SSSRole {
	for _, v := range s.roles {
		if v.seat == seat {
			return v
		}
	}
	return nil
}

//清除数据
func (s *SSSTable) clearData() {
	s.jump = netproto.SSSGameStep_SSSGameStepNull
	for _, v := range s.roles {
		v.isDo = false
		v.ris = NewNResultItem()
		v.gun = make([]int32, 0)
		v.act = nil
	}
}

//发送房间配置信息
func (s *SSSTable) PushGameCfg(user *gamelogic.GameUser) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_PushGameCfgId)
	gsdata := new(netproto.PushGameCfg)
	gsdata.DiFen = proto.Int32(s.difen)
	gsdata.TimeOut = proto.Int32(s.timeOut)
	pmsg.MsgData = gsdata
	s.SendMessageToUser(user, pmsg)
}

//发送游戏信息
func (s *SSSTable) PushGameInfo(user *gamelogic.GameUser) {
	mlog.Debug("PushGameInfo")

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_PushGameInfoId)
	gsdata := new(netproto.PushGameInfo)
	if (s.step >= netproto.SSSGameStep_SSSGameStepCompareH && s.step <= netproto.SSSGameStep_SSSGameStepCompareF) ||
		s.step == netproto.SSSGameStep_SSSGameStepEnd {
		gsdata.Step = proto.Int32(int32(s.step - 1))
	} else {
		gsdata.Step = proto.Int32(int32(s.step))
	}

	gsdata.Rris = make([]*netproto.RoleRecoverInfo, 0)
	gsdata.Sts = make([]*netproto.SItem, 0)
	for _, v := range s.roles {
		ri := new(netproto.RoleRecoverInfo)
		ri.Sis = make([]*netproto.SettlementCell, 0)
		ri.GunReciver = make([]int32, 0)
		ri.SeatNo = proto.Int32(v.seat)
		ri.Hcs = new(netproto.HandCards)
		if v.act != nil {
			ri.IsSpecial = proto.Bool(v.act.isSpecial)
			ri.Sct = proto.Int32(v.act.sct)
		} else {
			ri.IsSpecial = proto.Bool(false)
			ri.Sct = proto.Int32(0)
		}
		if s.step < netproto.SSSGameStep_SSSGameStepCompareH && v.act != nil { //还未开始比牌，数据恢复时不发送其他玩家的手牌信息
			if user.SeatNo == v.seat && s.step >= netproto.SSSGameStep_SSSGameStepBegin {
				mlog.Debug("s.step=%v", s.step)
				ri.Hcs.Head = v.act.recommandCt.GetHead()
				ri.Hcs.Body = v.act.recommandCt.GetBody()
				ri.Hcs.Foot = v.act.recommandCt.GetFoot()

				// fmt.Println(user.UserID, s.step, ri.Hcs.Head, ri.Hcs.Body, ri.Hcs.Foot, ri.GetSis())
			}
		} else { //组牌阶段已过，可以发送其他玩家的所有手牌信息
			for i := 0; i < 3; i++ { //每墩牌的比较结果
				var lose, win int32 = 0, 0
				for _, v := range v.ris[i].result {
					if v < 0 {
						lose += v
					} else {
						win += v
					}
				}
				c := new(netproto.SettlementCell)
				c.Win = proto.Int32(win)
				c.Lose = proto.Int32(lose)
				ri.Sis = append(ri.Sis, c)
			}

			if v.act != nil {
				//玩家手牌
				//mlog.Debug("1s.step=%v", s.step)
				if v.act.isSpecial && s.step >= netproto.SSSGameStep_SSSGameStepCompareH {
					mlog.Debug("1s.step=%v", s.step)
					upk := s.gamePoker.GetUserPoker(v.roleid)
					cg := poker.NewCardGroup(upk.Cards.Cards)
					v.act.cgs[0] = s.GetSP(3, cg)
					v.act.cgs[1] = s.GetSP(5, cg)
					v.act.cgs[2] = s.GetSP(5, cg)
					v.act.Cg2Hc()
				}

				if v.act.recommandCt != nil {
					if s.step >= netproto.SSSGameStep_SSSGameStepCompareH {
						ri.Hcs.Head = v.act.recommandCt.GetHead()
					}

					if s.step >= netproto.SSSGameStep_SSSGameStepCompareB {
						ri.Hcs.Body = v.act.recommandCt.GetBody()
					}

					if s.step >= netproto.SSSGameStep_SSSGameStepCompareF {
						ri.Hcs.Foot = v.act.recommandCt.GetFoot()
					}
				}
			}

			//打包打枪信息
			ri.GunReciver = append(ri.GunReciver, v.gun...)
		}

		ri.ZuPaiOver = proto.Bool(v.isZuPai)
		mlog.Debug("组牌完成：%v", ri.ZuPaiOver)
		gsdata.Rris = append(gsdata.Rris, ri)
	}
	var lose, win, amount, coins int32 = 0, 0, 0, 0
	//打包结算信息
	gsdata.Sts = make([]*netproto.SItem, 0)
	for _, v := range s.roles {
		si := new(netproto.SItem)
		if s.step >= netproto.SSSGameStep_SSSGameStepCompareH {
			lose, win, amount, coins = s.sectionSettlement(v)
		} else {
			lose, win, amount, coins = 0, 0, 0, 0
		}
		si.Lose = proto.Int32(lose)
		si.Win = proto.Int32(win)
		si.Amount = proto.Int32(amount)
		si.Coins = proto.Int32(coins)
		si.SeatNo = proto.Int32(v.seat)
		gsdata.Sts = append(gsdata.Sts, si)
	}

	gsdata.TimeOut = proto.Int32(s.getStepTimeOut())
	if s.stepTimer != nil {
		gsdata.LastTime = proto.Int32(s.getStepTimeOut() - s.stepTimer.ElapsedSeconds())
	} else {
		gsdata.LastTime = proto.Int32(0)
	}

	if s.IsGameUser(user) {
		upk := s.gamePoker.GetUserPoker(user.UserID)
		for _, pk := range upk.Cards.Cards {
			card := new(netproto.SCard)
			card.Color = proto.Int32(pk.Color)
			card.Number = proto.Int32(pk.CardNumber)
			gsdata.Cards = append(gsdata.Cards, card)
		}
	}

	// fmt.Println("发送游戏信息", user.UserID, gsdata.GetStep(), gsdata.GetRris())

	pmsg.MsgData = gsdata
	mlog.Debug("发送PushGameInfo")
	s.SendMessageToUser(user, pmsg)
}

//广播游戏操作超时信息
func (t *SSSTable) broadcastOptTime() {
	t.BoardcastMessage(t.getOptMsg())
	mlog.Debug("=====================step = %d", t.step)
}

func (t *SSSTable) getOptMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_NotifyOptTimeOutId)
	gsdata := new(netproto.NotifyOptTimeOut)
	gsdata.Step = proto.Int32(int32(t.step))
	gsdata.TimeOut = proto.Int32(t.getStepTimeOut())
	pmsg.MsgData = gsdata
	return pmsg
}

//广播玩家牌型选择消息（是否选择特殊牌型）
func (t *SSSTable) broadcastUserCardTypeChoose(user *gamelogic.GameUser, choose bool) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_NotifyUserCardTypeChooseId)
	gsdata := new(netproto.NotifyUserCardTypeChoose)
	gsdata.SeatNo = proto.Int32(user.SeatNo)
	gsdata.Special = proto.Bool(choose)
	pmsg.MsgData = gsdata

	t.BoardcastMessage(pmsg)
}

//广播用户牌型 普通牌型比较
func (t *SSSTable) broadcastUserCards(nDun netproto.SSSDunType) {
	// fmt.Println("广播用户牌型", nDun)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_NotifyUserCardsId)
	gsdata := new(netproto.NotifyUserCards)
	gsdata.Ci = make([]*netproto.CompareItem, 0)
	for _, v := range t.roles {
		if v.act.isSpecial {
			continue
		}
		ci := new(netproto.CompareItem)
		if nDun == netproto.SSSDunType_DunTypeHead {
			ci.D = v.act.recommandCt.GetHead()
		} else if nDun == netproto.SSSDunType_DunTypeBody {
			ci.D = v.act.recommandCt.GetBody()
		} else {
			ci.D = v.act.recommandCt.GetFoot()
		}
		ci.D.Sdt = proto.Int32(int32(nDun))
		ci.D.Ct = proto.Int32(int32(v.ris[nDun-1].Ct))
		var lose, win int32 = 0, 0
		for _, v := range v.ris[nDun-1].result {
			if v < 0 {
				lose += v
			} else {
				win += v
			}
		}
		ci.Lose = proto.Int32(lose)
		ci.Win = proto.Int32(win)
		ci.SeatNo = proto.Int32(v.seat)
		gsdata.Ci = append(gsdata.Ci, ci)
		// fmt.Println("QQQ", v.seat, v.roleid, ci.D, ci.D.Sdt, ci.D.Ct)
	}
	pmsg.MsgData = gsdata

	// fmt.Println(gsdata.GetCi())
	t.BoardcastMessage(pmsg)
}

//广播特殊牌型比牌结果
func (t *SSSTable) broadcastSpecialCompareResult() {
	t.BoardcastMessage(t.getSpecialCompareResultMsg())
}

func (t *SSSTable) getSpecialCompareResultMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_NotifyUserSpecialId)
	gsdata := new(netproto.NotifyUserSpecial)
	gsdata.Sci = make([]*netproto.SpecialCompareItem, 0)
	for _, v := range t.roles {
		if !v.act.isSpecial {
			continue
		}
		ci := new(netproto.SpecialCompareItem)
		ci.SeatNo = proto.Int32(v.seat)
		ci.Lose = proto.Int32(v.act.sctLose)
		ci.Win = proto.Int32(v.act.sctWin)
		ci.SpecialType = proto.Int32(v.act.sct)
		for _, c := range t.gamePoker.GetUserPoker(v.roleid).Cards.Cards {
			sc := new(netproto.SCard)
			sc.Number = proto.Int32(c.CardNumber)
			sc.Color = proto.Int32(c.Color)
			ci.Cards = append(ci.Cards, sc)
		}
		gsdata.Sci = append(gsdata.Sci, ci)
	}
	pmsg.MsgData = gsdata

	return pmsg
}

//广播玩家打枪信息
func (t *SSSTable) broadcastUserGun() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_NotifyUserGunId)
	gsdata := new(netproto.NotifyUserGun)
	gsdata.Ugs = make([]*netproto.UserGun, 0)
	for _, v := range t.roles {
		size := len(v.gun)
		if size > 0 {
			g := new(netproto.UserGun)
			if size == 3 { //全垒打
				g.IsAll = proto.Bool(true)
			} else {
				g.IsAll = proto.Bool(false)
			}
			g.Attacter = proto.Int32(v.seat)
			for _, n := range v.gun {
				g.Recvers = append(g.Recvers, n)
			}
			gsdata.Ugs = append(gsdata.Ugs, g)
		}
	}
	pmsg.MsgData = gsdata
	t.BoardcastMessage(pmsg)
}

func (t *SSSTable) CaculateOddsAndGun() {
	var odds int32 = 0
	for _, role := range t.roles {
		if len(role.gun) == 3 { //全垒打
			odds = t.GetGunCfg(ALLGun)
			for i := 0; i < 3; i++ {
				result := role.ris[i].result
				for seat, m := range result {
					if seat == role.seat {
						continue
					}
					role.ris[i].result[seat] = m * odds
					r := t.getRoleInfoBySeat(seat)
					if r != nil {
						r.ris[i].result[role.seat] = -m * odds
					}
				}
			}
		} else {
			if len(role.gun) > 0 { //普通打枪
				odds = t.GetGunCfg(NormalGun)
				for _, v := range role.gun {
					r := t.getRoleInfoBySeat(v)
					if r != nil {
						for i := 0; i < 3; i++ {
							r.ris[i].result[role.seat] *= odds
							role.ris[i].result[v] *= odds
						}
					}
				}
			}
		}
	}
}

//计算玩家输赢 return :总输的水数 总赢的水数 总计
func (t *SSSTable) CaculateSettlement(role *SSSRole) (int32, int32, int32, int32) {
	var Coins, amount, Win, Lose int32 = 0, 0, 0, 0
	for _, r := range role.ris {
		for _, n := range r.result {
			if n > 0 {
				Win += n
			} else {
				Lose += n
				mlog.Debug("Lose:::%v", Lose)
			}
		}
	}

	var flowAmount, betAmount, payAmount int64
	Lose += role.act.sctLose
	Win += role.act.sctWin
	amount = (Lose + Win)
	Coins = amount * t.difen
	sendmoney, tax := t.GetUserWinAndTax(int64(Coins))
	t.UpdateUserScoreFunc(role.roleid, sendmoney, tax)
	user := t.GetUserByUserID(role.roleid)
	flowAmount = util.ABSInt64(sendmoney)
	betAmount = util.ABSInt64(sendmoney) + tax
	payAmount = betAmount + sendmoney
	retmsg := "用戶%d, 流水%d, 有效投注 %d, 返還金額 %d"
	if user.IsRobot() {
		retmsg = "[機器人]" + retmsg
	}
	t.GameLog.AddGameStep(retmsg, user.UserID, flowAmount, betAmount, payAmount)
	user.FlowAmount += flowAmount
	user.BetAmount += betAmount
	user.PayAmount += payAmount

	t.GameLog.Context().UserUpdateScore(user.UserID, sendmoney)
	t.GameLog.Context().UserBetScore(user.UserID, betAmount)
	t.GameLog.Context().UserPayoutScore(user.UserID, payAmount)
	t.GameLog.Context().UserValidScore(user.UserID, flowAmount)
	t.GameLog.Context().UserTaxAmount(user.UserID, tax)

	return Lose, Win, amount, int32(sendmoney)
}

func (t *SSSTable) sectionSettlement(role *SSSRole) (int32, int32, int32, int32) {
	var Coins, amount, Win, Lose int32 = 0, 0, 0, 0
	for _, r := range role.ris {
		for _, n := range r.result {
			if n > 0 {
				Win += n
			} else {
				Lose += n
			}
		}
	}
	Lose += role.act.sctLose
	Win += role.act.sctWin
	amount = (Lose + Win)
	Coins = amount * t.difen
	sendmoney, _ := t.GetUserWinAndTax(int64(Coins))
	return Lose, Win, amount, int32(sendmoney)
}

//广播结算消息
func (t *SSSTable) broadCastSettlement() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_NotifySSSSettlementId)
	gsdata := new(netproto.NotifySSSSettlement)
	gsdata.Sis = make([]*netproto.SItem, 0)
	uscoresRet := make(map[int32]int64)
	for _, v := range t.roles {
		mlog.Debug("roles:%v", v)
		si := new(netproto.SItem)
		lose, win, amount, coins := t.CaculateSettlement(v)
		si.Lose = proto.Int32(lose)
		si.Win = proto.Int32(win)
		si.Amount = proto.Int32(amount)
		si.Coins = proto.Int32(coins)
		si.SeatNo = proto.Int32(v.seat)
		gsdata.Sis = append(gsdata.Sis, si)
		uscoresRet[v.roleid] = int64(coins)
		if v.act.cgs[0] == nil {
			cg := poker.NewCardGroup(nil)
			for _, v2 := range v.act.recommandCt.GetHead().GetCards() {
				cg.AddCard(&poker.Card{Color: v2.GetColor(), CardNumber: v2.GetNumber()})
			}
			t.GameLog.AddGameStep("玩家%d选择头敦牌型%s", v.roleid, cg.String())
		} else {
			t.GameLog.AddGameStep("玩家%d选择头敦牌型%s", v.roleid, v.act.cgs[0].String())
		}

		if v.act.cgs[1] == nil {
			cg := poker.NewCardGroup(nil)
			for _, v2 := range v.act.recommandCt.GetBody().GetCards() {
				cg.AddCard(&poker.Card{Color: v2.GetColor(), CardNumber: v2.GetNumber()})
			}
			t.GameLog.AddGameStep("玩家%d选择中敦牌型%s", v.roleid, cg.String())
		} else {
			t.GameLog.AddGameStep("玩家%d选择中敦牌型%s", v.roleid, v.act.cgs[1].String())
		}

		if v.act.cgs[2] == nil {
			cg := poker.NewCardGroup(nil)
			for _, v2 := range v.act.recommandCt.GetFoot().GetCards() {
				cg.AddCard(&poker.Card{Color: v2.GetColor(), CardNumber: v2.GetNumber()})
			}
			t.GameLog.AddGameStep("玩家%d选择后敦牌型%s", v.roleid, cg.String())
		} else {
			t.GameLog.AddGameStep("玩家%d选择后敦牌型%s", v.roleid, v.act.cgs[2].String())
		}
	}

	t.GameLog.GameEnd()
	pmsg.MsgData = gsdata
	t.BoardcastMessage(pmsg)
}

//广播特殊牌型比较
func (t *SSSTable) broadSpecialCompareResult() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_NotifyUserSpecialId)
	gsdata := new(netproto.NotifyUserSpecial)
	gsdata.Sci = make([]*netproto.SpecialCompareItem, 0)
	for _, v := range t.roles {
		si := new(netproto.SpecialCompareItem)
		si.Cards = make([]*netproto.SCard, 0)
		si.Lose = proto.Int32(v.act.sctLose)
		si.Win = proto.Int32(v.act.sctWin)
		si.SpecialType = proto.Int32(v.act.sct)
		si.SeatNo = proto.Int32(v.seat)
		upk := t.gamePoker.GetUserPoker(v.roleid)
		for _, c := range upk.Cards.Cards {
			sc := new(netproto.SCard)
			sc.Color = proto.Int32(c.Color)
			sc.Number = proto.Int32(c.CardNumber)
			si.Cards = append(si.Cards, sc)
		}
		gsdata.Sci = append(gsdata.Sci, si)
	}
	pmsg.MsgData = gsdata
	t.BoardcastMessage(pmsg)
}

//选择特殊牌型消息响应
func (t *SSSTable) responseChooseSpecialCardTypeRet(code int32, strMsg string, seatNo int32, st int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_ChooseSpecialCardTypeRetId)
	gsdata := new(netproto.ChooseSpecialCardTypeRet)
	gsdata.Code = proto.Int32(code)
	gsdata.Msg = proto.String(strMsg)
	gsdata.SeatNo = proto.Int32(seatNo)
	gsdata.SpecialType = proto.Int32(st)
	pmsg.MsgData = gsdata
	t.BoardcastMessage(pmsg)
}

//组牌完成响应
func (t *SSSTable) responseCompleteCardCombRet(user *gamelogic.GameUser, code int32, strMsg string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_CompleteCardCombRetId)
	gsdata := new(netproto.CompleteCardCombRet)
	gsdata.Code = proto.Int32(code)
	gsdata.Msg = proto.String(strMsg)
	pmsg.MsgData = gsdata
	t.SendMessageToUser(user, pmsg)
	//t.BoardcastMessage(pmsg)
}

//根据牌点获取相同数目的牌数量是否存在
func GroupByNumber(dun *netproto.Dun, count int32) bool {
	m := make(map[int32]int32)
	for _, v := range dun.Cards {
		if _, ok := m[v.GetNumber()]; ok {
			m[v.GetNumber()] += 1
		} else {
			m[v.GetNumber()] = 1
		}
	}
	for _, v := range m {
		if v == count {
			return true
		}
	}
	return false
}

func GetSameNumberCount(dun *netproto.Dun, count int32) int32 {
	m := make(map[int32]int32)
	for _, v := range dun.Cards {
		if _, ok := m[v.GetNumber()]; ok {
			m[v.GetNumber()] += 1
		} else {
			m[v.GetNumber()] = 1
		}
	}
	var n int32 = 0
	for _, v := range m {
		if v == count {
			n++
		}
	}
	return n
}

func GetCardsMap(arr []*poker.Card) map[int32]int32 {
	m := make(map[int32]int32)
	for _, v := range arr {
		if _, ok := m[v.CardNumber]; ok {
			m[v.CardNumber] += 1
		} else {
			m[v.CardNumber] = 1
		}
	}
	return m
}

func GetColorMap(arr []*poker.Card) map[int32]int32 {
	m := make(map[int32]int32)
	for _, v := range arr {
		if _, ok := m[v.Color]; ok {
			m[v.Color] += 1
		} else {
			m[v.Color] = 1
		}
	}
	return m
}

//校验同花顺
func CheckDunTongHuaShun(dun *netproto.Dun) bool {
	if !CheckDunTongHua(dun) || !CheckDunShunZi(dun) {
		mlog.Error("不是同花顺牌型,%v", dun.Cards)
		return false
	}
	return true
}

//校验炸弹
func CheckDunZhaDan(dun *netproto.Dun) bool {
	if !GroupByNumber(dun, 4) {
		mlog.Error("不是炸弹牌型")
		return false
	}
	return true
}

//校验葫芦
func CheckDunHuLu(dun *netproto.Dun) bool {
	if !GroupByNumber(dun, 3) || !GroupByNumber(dun, 2) {
		mlog.Error("不是葫芦牌型")
		return false
	}
	return true
}

//校验同花
func CheckDunTongHua(dun *netproto.Dun) bool {
	var color int32 = 0
	for _, v := range dun.Cards {
		if color == 0 {
			color = v.GetColor()
		} else {
			if color != v.GetColor() {
				mlog.Error("不是同花牌型 %v", dun.String())
				return false
			}
		}
	}
	return true
}

//校验顺子
func CheckDunShunZi(dun *netproto.Dun) bool {
	cg := new(poker.CardGroup)
	for _, v := range dun.Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		cg.AddCard(c)
	}

	status := ckDunSantiao(cg)
	return status
}

func ckDunSantiao(cg *poker.CardGroup) bool {
	size := len(cg.Cards)
	b := false
	if HasA(cg.Cards) && IsAllSmall5(cg.Cards) {
		b = true
	}
	scg := CardGroupSortNumber{cg, b}
	sort.Sort(&scg)

	for i := 1; i < size; i++ {
		if b {
			if math.Abs(float64(getCardNumber(cg.Cards[i].CardNumber)-getCardNumber(cg.Cards[i-1].CardNumber))) != 1 {
				mlog.Error("不是顺子牌型 1 %v", cg.String())
				return false
			}
		} else {
			if math.Abs(float64(cg.Cards[i].CardNumber-cg.Cards[i-1].CardNumber)) != 1 {
				mlog.Error("不是顺子牌型 2 %v", cg.String())
				return false
			}
		}

	}

	return true
}

//校验三条
func CheckDunSantiao(dun *netproto.Dun) bool {
	if !GroupByNumber(dun, 3) || !GroupByNumber(dun, 1) {
		mlog.Error("不是三条牌型")
		return false
	}
	return true
}

//校验冲三
func CheckDunChongSan(dun *netproto.Dun) bool {
	if dun.GetSdt() != int32(netproto.SSSDunType_DunTypeHead) {
		mlog.Error("不是头墩，不是冲三牌型")
		return false
	}
	if !GroupByNumber(dun, 3) {
		mlog.Error("不是冲三牌型")
		return false
	}
	return true
}

//校验两对
func CheckDunLiangDui(dun *netproto.Dun) bool {
	if GetSameNumberCount(dun, 2) != 2 {
		mlog.Error("不是两对牌型")
		return false
	}
	return true
}

//校验一对
func CheckDunYiDui(dun *netproto.Dun) bool {
	if GetSameNumberCount(dun, 2) != 1 {
		mlog.Error("不是两对牌型")
		return false
	}
	if GetSameNumberCount(dun, 3) != 0 {
		mlog.Error("不是两对牌型2")
		return false
	}
	return true
}

//校验一对
func CheckDunSanPai(dun *netproto.Dun) bool {
	if dun.GetSdt() == int32(netproto.SSSDunType_DunTypeHead) {
		if GetSameNumberCount(dun, 1) != 3 {
			mlog.Error("不是散牌牌型1")
			return false
		}
	} else {
		if GetSameNumberCount(dun, 1) != 5 {
			mlog.Error("不是散牌牌型2")
			return false
		}
	}
	return true
}

//检测头墩类型是否合法
func CheckHeadDun(dun *netproto.Dun) bool {
	if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_ChongSan) {
		return CheckDunChongSan(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_YiDui) {
		return CheckDunYiDui(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_SanPai) {
		return CheckDunSanPai(dun)
	} else {
		mlog.Error("头墩牌型检测失败%s", dun.String())
		return false
	}
}

//检测墩类型是否合法
func CheckDun(dun *netproto.Dun) bool {
	if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
		return CheckDunTongHuaShun(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
		return CheckDunZhaDan(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_HuLu) {
		return CheckDunHuLu(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_TongHua) {
		return CheckDunTongHua(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_ShunZi) {
		return CheckDunShunZi(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_SanTiao) {
		return CheckDunSantiao(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_ChongSan) {
		return CheckDunChongSan(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_LiangDui) {
		return CheckDunLiangDui(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_YiDui) {
		return CheckDunYiDui(dun)
	} else if dun.GetCt() == int32(netproto.SSSDunCardsNomalType_SanPai) {
		return CheckDunSanPai(dun)
	} else {
		mlog.Error("墩牌型检测失败%s", dun.String())
		return false
	}
}

// var searchTable map[netproto.SSSDunCardsNomalType]WeightTable

//获取比较权值
func (s *SSSTable) GetCompareWeight(dun *netproto.Dun) int64 {
	var w int64 = 0
	var c int = 5
	var x int64 = 10000000000
	cg := new(poker.CardGroup)
	for _, v := range dun.Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		cg.AddCard(c)
	}
	scg := CardGroupSortNumber{cg, false}
	sort.Sort(&scg)
	if dun.GetSdt() == int32(netproto.SSSDunType_DunTypeHead) {
		for i := len(cg.Cards) - 1; i >= 0; i-- {
			w += int64(cg.Cards[i].CardNumber) * int64(math.Pow(float64(100), float64(c-i)))
		}
	}
	w += int64(s.getWeight(int32(dun.GetSdt()))) * x
	return w
}

//判断是否倒水
func CheckIsDaoShui(hcs *netproto.HandCards, s *SSSTable) bool {
	h1 := s.GetCompareWeight(hcs.GetHead())
	h2 := s.GetCompareWeight(hcs.GetBody())
	h3 := s.GetCompareWeight(hcs.GetFoot())
	if h3 >= h2 && h2 > h1 {
		return true
	}
	return false
}

//判断玩家发送上来的牌与服务器玩家手牌是否是被包含关系
func CheckShouPaiContent(hcs *netproto.HandCards, upks *poker.UserPoker) bool {
	cg := poker.NewCardGroup(upks.Cards.Cards)
	tmpcg := poker.NewCardGroup(nil)
	for _, v := range hcs.GetHead().Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		tmpcg.AddCard(c)
	}
	for _, v := range hcs.GetBody().Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		tmpcg.AddCard(c)
	}
	for _, v := range hcs.GetFoot().Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		tmpcg.AddCard(c)
	}

	cg.DeleteGroup(tmpcg)
	mlog.Debug("cg.DeleteGroup:%v", cg)

	if cg.Len() > 0 {
		return false
	}
	return true
}

//检测是否是同花
func CheckTongHua(arr []*poker.Card) bool {
	var color int32 = 0
	for _, v := range arr {
		if color == 0 {
			color = v.Color
		} else {
			if color != v.Color {
				return false
			}
		}
	}
	return true
}

func HasA(arr []*poker.Card) bool {
	size := len(arr)
	for i := 0; i < size; i++ {
		if arr[i] != nil {
			if arr[i].CardNumber == poker.PokerNumberA {
				return true
			}
		}
	}
	return false
}

func IsAllSmall5(arr []*poker.Card) bool {
	size := len(arr)
	for i := 0; i < size; i++ {
		if arr[i] != nil {
			if getCardNumber(arr[i].CardNumber) > 5 {
				return false
			}
		}
	}
	return true
}

//检测是否是顺子
func CheckShunZi(arr []*poker.Card) bool {
	cg := poker.NewCardGroup(arr)
	size := len(arr)
	b := false
	if HasA(arr) && IsAllSmall5(arr) { //特殊处理
		b = true
	}
	scg := CardGroupSortNumber{cg, b}
	sort.Sort(&scg)
	for i := 1; i < size; i++ {
		if b {
			if math.Abs(float64(getCardNumber(cg.Cards[i].CardNumber)-getCardNumber(cg.Cards[i-1].CardNumber))) != 1 {
				return false
			}
		} else {
			if math.Abs(float64(cg.Cards[i].CardNumber-cg.Cards[i-1].CardNumber)) != 1 {
				return false
			}
		}
	}
	return true
}

//校验同花十三水
func (s *SSSTable) CheckTongHuaShiSanShui(usercards *poker.CardGroup) bool {
	if !CheckShunZi(usercards.Cards) || !CheckTongHua(usercards.Cards) {
		return false
	}
	return true
}

//校验十三水
func (s *SSSTable) CheckShiSanShui(usercards *poker.CardGroup) bool {
	if !CheckShunZi(usercards.Cards) {
		return false
	}
	return true
}

//校验十二皇族
func (s *SSSTable) CheckShiErHuangZu(usercards *poker.CardGroup) bool {
	count := 0
	for _, v := range usercards.Cards {
		if getCardNumber(v.CardNumber) <= 10 {
			count++
			if count >= 2 {
				return false
			}
		}
	}
	return true
}

//校验三同花顺
func (s *SSSTable) CheckSanTongHuaShun(usercards *poker.CardGroup) bool {
	return CheckSanTongHuaShunPatterns(usercards)
}

//校验三分天下
func (s *SSSTable) CheckSanFenTianXia(usercards *poker.CardGroup) bool {
	m := GetCardsMap(usercards.Cards)
	var n int32 = 0
	for _, v := range m {
		if v == 4 {
			n++
		}
	}
	if n == 3 {
		return true
	}
	return false
}

//校验全大
func (s *SSSTable) CheckQuanDa(usercards *poker.CardGroup) bool {
	for _, v := range usercards.Cards {
		if v.CardNumber <= 7 {
			return false
		}
	}
	return true
}

//校验全小
func (s *SSSTable) CheckQuanXiao(usercards *poker.CardGroup) bool {
	for _, v := range usercards.Cards {
		if v.CardNumber >= 9 {
			return false
		}
	}
	return true
}

//校验凑一色
func (s *SSSTable) CheckCouYiSe(usercards *poker.CardGroup) bool {
	var cardCount = usercards.Len()
	var blackCount, redCount = 0, 0
	for _, v := range usercards.Cards {
		if v.IsHeiTao() || v.IsMeiHua() {
			blackCount++
		}
		if v.IsHongTao() || v.IsFangKuai() {
			redCount++
		}
	}

	if blackCount == cardCount || redCount == cardCount {
		return true
	}
	return false
}

//校验六对半
func (s *SSSTable) CheckLiuDuiBan(usercards *poker.CardGroup) bool {
	cg := poker.NewCardGroup(usercards.Cards)
	m := GetCardsMap(cg.Cards)
	n := 0
	for _, v := range m {
		if v == 2 {
			n++
		}
	}
	if n == 6 {
		return true
	}
	return false
}

func (s *SSSTable) sortCardsByColor(cards *poker.CardGroup) {
	arrlen := len(cards.Cards)
	for i := 0; i < arrlen; i++ {
		for j := 0; j+1 < arrlen-i; j++ {
			card1 := cards.Cards[j].CardNumber + cards.Cards[j].Color*100
			card2 := cards.Cards[j+1].CardNumber + cards.Cards[j+1].Color*100
			if card1 > card2 {
				cards.Swap(j, j+1)
			}
		}
	}

	//mlog.Debug("花色排序后的扑克：%v", cards)
}

func (s *SSSTable) sortCardsByColorNodesc(cards *poker.CardGroup) {
	s.sortCardsByColor(cards)
	arrlen := len(cards.Cards)
	for i := 0; i < arrlen; i++ {
		for j := 0; j+1 < arrlen-i; j++ {
			if cards.Cards[j].Color == cards.Cards[j+1].Color {
				card1 := cards.Cards[j].CardNumber
				card2 := cards.Cards[j+1].CardNumber
				if card1 < card2 {
					cards.Swap(j, j+1)
				}
			}
		}
	}
}

//校验三同花
func (s *SSSTable) CheckSanTongHua(usercards *poker.CardGroup) bool {
	// return CheckSanTongHuaPatterns(usercards)
	cg := poker.NewCardGroup(usercards.Cards)
	//sort.Sort(cg)
	s.sortCardsByColor(cg)
	foot := s.GetTongHua(5, cg)
	body := s.GetTongHua(5, cg)
	head := s.GetTongHua(3, cg)
	if foot == nil || body == nil || head == nil {
		return false
	}
	return true
}

//校验三顺子
func (s *SSSTable) CheckSanShunZi(usercards *poker.CardGroup) bool {
	// return CheckSanShunZiPatterns(usercards)
	cg := poker.NewCardGroup(usercards.Cards)

	foot := s.GetSZ(5, cg)
	body := s.GetSZ(5, cg)
	head := s.GetSZ(3, cg)

	if foot != nil && body != nil && head != nil {
		return true
	}

	return false
}

//校验双怪冲三 2对葫芦+1个对子+任意一张杂牌
func (s *SSSTable) CheckShuangGuaiChongSan(usercards *poker.CardGroup) bool {
	cg := poker.NewCardGroup(usercards.Cards)
	m := GetCardsMap(cg.Cards)
	twoPairCount, threeOfAKindCount := 0, 0
	for _, v := range m {
		if v == 3 {
			threeOfAKindCount++
		}
		if v == 2 {
			twoPairCount++
		}
	}
	if twoPairCount == 3 && threeOfAKindCount == 2 {
		return true
	}
	return false
}

//校验三冲加弹 三个三条，剩下四张是炸弹
func (s *SSSTable) CheckSanChongJiaDan(usercards *poker.CardGroup) bool {
	cg := poker.NewCardGroup(usercards.Cards)
	m := GetCardsMap(cg.Cards)
	zhaDanCount, threeOfAKindCount := 0, 0
	for _, v := range m {
		if v == 3 {
			threeOfAKindCount++
		}
		if v == 4 {
			zhaDanCount++
		}
	}
	if zhaDanCount == 1 && threeOfAKindCount == 3 {
		return true
	}
	return false
}

//校验四套三条
func (s *SSSTable) CheckSiTaoSanTiao(usercards *poker.CardGroup) bool {
	cg := poker.NewCardGroup(usercards.Cards)
	m := GetCardsMap(cg.Cards)
	n := 0
	for _, v := range m {
		if v >= 3 {
			n++
		}
	}
	if n == 4 {
		return true
	}
	return false
}

//校验五对三条
func (s *SSSTable) CheckWuDuiSanTiao(usercards *poker.CardGroup) bool {
	cg := poker.NewCardGroup(usercards.Cards)
	m := GetCardsMap(cg.Cards)
	twoPairCount, threeOfAKindCount := 0, 0
	for _, v := range m {
		if v == 3 {
			threeOfAKindCount++
		}
		if v == 2 {
			twoPairCount++
		}
	}
	if twoPairCount == 5 && threeOfAKindCount == 1 {
		return true
	}
	return false
}

func (s *SSSTable) checkSpecialType(user *gamelogic.GameUser, hcs *netproto.HandCards, special bool, st int32) bool {
	upk := s.gamePoker.GetUserPoker(user.UserID)
	if st == int32(netproto.SSSDunCardsSpecialType_TongHuaShiSanShui) {
		return s.CheckTongHuaShiSanShui(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_ShiSanShui) {
		return s.CheckShiSanShui(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_ShiErHuangZu) {
		return s.CheckShiErHuangZu(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_SanTongHuaShun) {
		return s.CheckSanTongHuaShun(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_SanFenTianXia) {
		return s.CheckSanFenTianXia(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_QuanDa) {
		return s.CheckQuanDa(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_QuanXiao) {
		return s.CheckQuanXiao(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_CouYiSe) {
		return s.CheckCouYiSe(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_SiTaoSantiao) {
		return s.CheckSiTaoSanTiao(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_WuDuiSanTiao) {
		return s.CheckWuDuiSanTiao(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_LiuDuiBan) {
		return s.CheckLiuDuiBan(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_SanTongHua) {
		return s.CheckSanTongHua(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_SanShunZi) {
		return s.CheckSanShunZi(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_ShuangGuaiChongSan) {
		return s.CheckShuangGuaiChongSan(upk.Cards)
	} else if st == int32(netproto.SSSDunCardsSpecialType_SanChongJiaDan) {
		return s.CheckSanChongJiaDan(upk.Cards)
	} else {
		mlog.Error("特殊牌型检测失败 %d", st)
	}

	return false
}

//检测玩家发送上来的牌型是否合法
func (s *SSSTable) verifyUserCardType(user *gamelogic.GameUser, hcs *netproto.HandCards, special bool, st int32) bool {
	if special { //特殊牌型检测
		upk := s.gamePoker.GetUserPoker(user.UserID)
		if st == int32(netproto.SSSDunCardsSpecialType_TongHuaShiSanShui) {
			return s.CheckTongHuaShiSanShui(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_ShiSanShui) {
			return s.CheckShiSanShui(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_ShiErHuangZu) {
			return s.CheckShiErHuangZu(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_SanTongHuaShun) {
			return s.CheckSanTongHuaShun(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_SanFenTianXia) {
			return s.CheckSanFenTianXia(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_QuanDa) {
			return s.CheckQuanDa(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_QuanXiao) {
			return s.CheckQuanXiao(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_CouYiSe) {
			return s.CheckCouYiSe(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_SiTaoSantiao) {
			return s.CheckSiTaoSanTiao(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_WuDuiSanTiao) {
			return s.CheckWuDuiSanTiao(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_LiuDuiBan) {
			return s.CheckLiuDuiBan(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_SanTongHua) {
			return s.CheckSanTongHua(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_SanShunZi) {
			return s.CheckSanShunZi(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_ShuangGuaiChongSan) {
			return s.CheckShuangGuaiChongSan(upk.Cards)
		} else if st == int32(netproto.SSSDunCardsSpecialType_SanChongJiaDan) {
			return s.CheckSanChongJiaDan(upk.Cards)
		} else {
			return false
		}
	} else {
		if !CheckShouPaiContent(hcs, s.gamePoker.GetUserPoker(user.UserID)) {
			mlog.Error("玩家发送上来的手牌错误")
			return false
		}
		if !CheckHeadDun(hcs.GetHead()) || !CheckDun(hcs.GetBody()) || !CheckDun(hcs.GetFoot()) {
			return false
		}
		if CheckIsDaoShui(hcs, s) {
			mlog.Error("倒水检测失败%v", hcs)
			return false
		}
	}
	return true
}

//检测是否是同花十三水
func (s *SSSTable) ACSpecialCheckTongHuaShiSanShui(userId int32) (netproto.SSSDunCardsSpecialType, bool) {
	return netproto.SSSDunCardsSpecialType_QuanXiao, false //暂时
}

//获取同花顺牌型
func (s *SSSTable) GetTongHuaShun(count int, cg *poker.CardGroup) *poker.CardGroup {
	if cg.Len() < count {
		return nil
	}

	m := make(map[int32]int32)
	for i := 0; i < cg.Len(); i++ {
		if _, ok := m[cg.Cards[i].CardNumber]; ok {
			m[cg.Cards[i].CardNumber] += 1
		} else {
			m[cg.Cards[i].CardNumber] = 1
		}
	}

	// 取得順子牌型
	shunziList := s.CheckSZPattern(cg)
	for _, cards := range shunziList {
		ret := s.GetTH(count, cards, true)
		if ret != nil {
			cg.DeleteGroup(ret)
			return ret
		}
	}

	return nil
}

//获取同花牌型
func (s *SSSTable) GetTongHua(count int, cg *poker.CardGroup) *poker.CardGroup {
	if cg.Len() < count {
		return nil
	}
	ret := poker.NewCardGroup(nil)
	var card *poker.Card = nil
	for i := 0; i < cg.Len(); i++ {
		if card == nil {
			card = cg.Cards[i]
			ret.AddCard(cg.Cards[i])
		} else {
			if cg.Cards[i].Color == card.Color {
				ret.AddCard(cg.Cards[i])
				card = cg.Cards[i]
			} else {
				card = cg.Cards[i]
				ret = poker.NewCardGroup(nil)
				ret.AddCard(card)
			}

			if ret.Len() >= count {
				break
			}
		}
	}
	if ret.Len() < count {
		ret = nil
	}
	if ret != nil {
		cg.DeleteGroup(ret)
	}
	return ret
}

//获取炸弹牌型
func (s *SSSTable) Get4(cg *poker.CardGroup) *poker.CardGroup {
	var count = 4
	if cg.Len() < count {
		return nil
	}

	m := make(map[int32]int32)
	for i := 0; i < cg.Len(); i++ {
		if _, ok := m[cg.Cards[i].CardNumber]; ok {
			m[cg.Cards[i].CardNumber] += 1
		} else {
			m[cg.Cards[i].CardNumber] = 1
		}
	}

	for number, c := range m {
		if c == 4 {
			ret := poker.NewCardGroup(cg.GetPokerCardByNumber(number))
			if ret != nil {
				cg.DeleteGroup(ret)
				return ret
			}
		}
	}

	return nil
}

//获取三条
func (s *SSSTable) Get3(cg *poker.CardGroup) *poker.CardGroup {
	var count = 3
	if cg.Len() < count {
		return nil
	}
	m := GetCardsMap(cg.Cards)
	threeOfAKindList := make([]int32, 0)
	for number, v := range m {
		if v >= 3 {
			threeOfAKindList = append(threeOfAKindList, number)
		}
	}
	// 取最大三條
	sort.SliceStable(threeOfAKindList, func(i, j int) bool { return threeOfAKindList[i] > threeOfAKindList[j] })

	ret := poker.NewCardGroup(nil)
	if len(threeOfAKindList) > 0 {
		selectthreeOfAKind := cg.GetPokerCardByNumber(threeOfAKindList[0])
		for _, v := range selectthreeOfAKind {
			ret.AddCard(v)
			if ret.Len() >= count {
				break
			}
		}
	}
	if ret.Len() < count {
		ret = nil
	}
	if ret != nil {
		cg.DeleteGroup(ret)
	}
	return ret
}

//获取一对
func (s *SSSTable) Get1d(cg *poker.CardGroup) *poker.CardGroup {
	var count = 2
	if cg.Len() < count {
		return nil
	}

	m := make(map[int32]int32)
	for i := 0; i < cg.Len(); i++ {
		if _, ok := m[cg.Cards[i].CardNumber]; ok {
			m[cg.Cards[i].CardNumber] += 1
		} else {
			m[cg.Cards[i].CardNumber] = 1
		}
	}

	n := 0
	for _, v := range m {
		if v == 4 {
			n = n + 2
		} else if v >= 2 {
			n++
		}
	}

	ret := poker.NewCardGroup(nil)
	if n <= 3 {
		for i := 0; i < cg.Len(); i++ {
			ret = poker.NewCardGroup(nil)
			ret.AddCard(cg.Cards[i])
			for j := 0; j < cg.Len(); j++ {
				if cg.Cards[j].CardNumber-cg.Cards[i].CardNumber == 0 && cg.Cards[i].Color != cg.Cards[j].Color {
					ret.AddCard(cg.Cards[j])
					if ret.Len() >= count {
						break
					}
				}
			}
			if ret.Len() >= count {
				break
			}
		}
	} else {
		// 先取超過2張的牌
		for _, v := range cg.Cards {
			if m[v.CardNumber] > 2 {
				if len(ret.GetPokerCardByNumber(v.CardNumber)) < 2 {
					ret.AddCard(v)
				}
				if ret.Len() >= count {
					cg.DeleteGroup(ret)
					return ret
				}
			}
		}

		for i := 0; i < cg.Len(); i++ {
			ret = poker.NewCardGroup(nil)
			ret.AddCard(cg.Cards[i])
			for j := 0; j < cg.Len(); j++ {
				if cg.Cards[j].CardNumber-cg.Cards[i].CardNumber == 0 && cg.Cards[i].Color != cg.Cards[j].Color {
					ret.AddCard(cg.Cards[j])
					if ret.Len() >= count {
						break
					}
				}
			}
			if ret.Len() >= count {
				break
			}
		}
	}

	if ret.Len() < count {
		ret = nil
	}
	if ret != nil {
		cg.DeleteGroup(ret)
	}
	return ret
}

//获取較小的对子
func (s *SSSTable) get1dForCards(cg *poker.CardGroup) *poker.CardGroup {
	var count = 2
	if cg.Len() < count {
		return nil
	}
	m := GetCardsMap(cg.Cards)
	pairList := make([]int32, 0)
	for number, v := range m {
		if v >= 2 {
			pairList = append(pairList, number)
		}
	}
	sort.SliceStable(pairList, func(i, j int) bool { return pairList[i] < pairList[j] })

	ret := poker.NewCardGroup(nil)
	if len(pairList) > 0 {
		selectPair := cg.GetPokerCardByNumber(pairList[0])
		for _, v := range selectPair {
			ret.AddCard(v)
			if ret.Len() >= count {
				break
			}
		}
	}

	if ret.Len() < count {
		ret = nil
	}
	if ret != nil {
		cg.DeleteGroup(ret)
	}
	return ret
}

//检测是否可以获取出葫芦牌型
func (s *SSSTable) CheckHuLu(cg *poker.CardGroup) bool {
	var check = 0 //三條計數器
	m := make(map[int32]int32)
	for i := 0; i < cg.Len(); i++ {
		if _, ok := m[cg.Cards[i].CardNumber]; ok {
			m[cg.Cards[i].CardNumber] += 1
		} else {
			m[cg.Cards[i].CardNumber] = 1
		}
	}
	k3, k2 := false, false
	for _, v := range m {
		if v == 3 {
			k3 = true
			check += 1 //第一組三條，計數器+1
			if check > 1 {
				k2 = true //第二組三條，與第一組湊成一對葫蘆
			}
		} else if v == 2 {
			k2 = true
		}
	}
	return k3 && k2
}

//获取葫芦牌型
func (s *SSSTable) GetHL(cg *poker.CardGroup) *poker.CardGroup {
	if cg.Len() < 5 || !s.CheckHuLu(cg) {
		return nil
	}

	ret := poker.NewCardGroup(nil)
	st := s.Get3(cg)
	d := s.get1dForCards(cg)
	if st != nil {
		ret.AddCardGroup(st)
	}

	if d != nil {
		ret.AddCardGroup(d)
	}

	return ret
}

// CheckSZPattern 取出可能順子牌型(可能超過5張)
func (s *SSSTable) CheckSZPattern(cg *poker.CardGroup) (results []*poker.CardGroup) {
	results = make([]*poker.CardGroup, 0)
	// 取出順子牌型的可能 patten
	var szNumberPattern = [][]int32{
		{2, 3, 4, 5, 6},
		{3, 4, 5, 6, 7},
		{4, 5, 6, 7, 8},
		{5, 6, 7, 8, 9},
		{6, 7, 8, 9, 10},
		{7, 8, 9, 10, 11},
		{8, 9, 10, 11, 12},
		{9, 10, 11, 12, 13},
		{10, 11, 12, 13, 14},
		{14, 2, 3, 4, 5},
	}

	// 是否有符合牌型
	for _, pattern := range szNumberPattern {
		group := poker.NewCardGroup(nil)
		var match1, match2, match3, match4, match5 bool
		for i := 0; i < cg.Len(); i++ {
			if cg.Cards[i].CardNumber == pattern[0] {
				// if !match1 {
				group.AddCard(cg.Cards[i])
				match1 = true
				// }
			}
			if cg.Cards[i].CardNumber == pattern[1] {
				// if !match2 {
				group.AddCard(cg.Cards[i])
				match2 = true
				// }
			}
			if cg.Cards[i].CardNumber == pattern[2] {
				// if !match3 {
				group.AddCard(cg.Cards[i])
				match3 = true
				// }
			}
			if cg.Cards[i].CardNumber == pattern[3] {
				// if !match4 {
				group.AddCard(cg.Cards[i])
				match4 = true
				// }
			}
			if cg.Cards[i].CardNumber == pattern[4] {
				// if !match5 {
				group.AddCard(cg.Cards[i])
				match5 = true
				// }
			}
		}
		if match1 && match2 && match3 && match4 && match5 {
			// fmt.Println(group)
			results = append(results, group)
		}
	}

	return results
}

//获取同花牌型
func (s *SSSTable) GetTH(count int, cg *poker.CardGroup, haveShunzi bool) *poker.CardGroup {
	if cg.Len() < count {
		return nil
	}

	ret := poker.NewCardGroup(nil)

	heiTaoCards := cg.GetPokerCardByColor(poker.PokerColorHeiTao)
	hongTaoCards := cg.GetPokerCardByColor(poker.PokerColorHongTao)
	fangKuaiCards := cg.GetPokerCardByColor(poker.PokerColorFangKuai)
	meiHuaCards := cg.GetPokerCardByColor(poker.PokerColorMeiHua)

	if len(heiTaoCards) >= count {
		cg := poker.NewCardGroup(heiTaoCards)
		ret.AddCard(cg.GetMaxCard())
		ret.AddCardGroup(s.getSPforCards(ret, count-1, cg))
	} else if len(hongTaoCards) >= count {
		cg := poker.NewCardGroup(hongTaoCards)
		ret.AddCard(cg.GetMaxCard())
		ret.AddCardGroup(s.getSPforCards(ret, count-1, cg))
	} else if len(fangKuaiCards) >= count {
		cg := poker.NewCardGroup(fangKuaiCards)
		ret.AddCard(cg.GetMaxCard())
		ret.AddCardGroup(s.getSPforCards(ret, count-1, cg))
	} else if len(meiHuaCards) >= count {
		cg := poker.NewCardGroup(meiHuaCards)
		ret.AddCard(cg.GetMaxCard())
		ret.AddCardGroup(s.getSPforCards(ret, count-1, cg))
	}

	if ret.Len() < count {
		ret = nil
	}

	if !haveShunzi && ret != nil && ret.IsShunzi() {
		ret = nil
	}

	if ret != nil {
		cg.DeleteGroup(ret)
	}
	return ret
}

//获取顺子牌型
func (s *SSSTable) GetSZ(count int, cg *poker.CardGroup) *poker.CardGroup {
	if cg.Len() < count {
		return nil
	}

	if count == 3 {
		if cg.IsShunzi() {
			return poker.NewCardGroup(cg.Cards)
		}
	} else {
		ret := poker.NewCardGroup(nil)
		shunziList := s.CheckSZPattern(cg)
		for _, shunzi := range shunziList {
			for _, card := range shunzi.Cards {
				if len(ret.GetPokerCardByNumber(card.CardNumber)) < 1 {
					ret.AddCard(card)
				}

				if ret.Len() >= count {
					cg.DeleteGroup(ret)
					return ret
				}
			}
		}
	}

	return nil
}

//获取顺子牌型
func (s *SSSTable) GetAllSZ(count int, cg *poker.CardGroup) *poker.CardGroup {
	if cg.Len() < count {
		return nil
	}

	ret := poker.NewCardGroup(nil)
	var card *poker.Card = nil
	for i := 0; i < cg.Len(); i++ {
		if card == nil {
			card = cg.Cards[i]
			ret.AddCard(cg.Cards[i])
		} else {
			if ret.Len() >= count {
				break
			}
			if math.Abs(float64(cg.Cards[i].CardNumber-card.CardNumber)) == 1 {
				ret.AddCard(cg.Cards[i])
				card = cg.Cards[i]
			} else if math.Abs(float64(cg.Cards[i].CardNumber-card.CardNumber)) == 0 {
				card = cg.Cards[i]
			} else {
				ret = poker.NewCardGroup(nil)
				card = cg.Cards[i]
				ret.AddCard(cg.Cards[i])
			}
		}
	}
	if ret.Len() < count {
		ret = nil
	}
	if ret != nil {
		cg.DeleteGroup(ret)
	}
	return ret
}

//检测是否有两对牌型
func (s *SSSTable) Check2d(cg *poker.CardGroup) bool {
	if cg.Len() < 4 {
		return false
	}
	m := make(map[int32]int32)
	for i := 0; i < cg.Len(); i++ {
		if _, ok := m[cg.Cards[i].CardNumber]; ok {
			m[cg.Cards[i].CardNumber] += 1
		} else {
			m[cg.Cards[i].CardNumber] = 1
		}
	}
	n := 0
	for _, v := range m {
		if v >= 2 {
			n++
		}
	}
	if n >= 2 {
		return true
	}
	return false
}

//获取两对牌型
func (s *SSSTable) Get2d(cg *poker.CardGroup) *poker.CardGroup {
	if !s.Check2d(cg) {
		return nil
	}

	d1 := s.Get1d(cg)
	d2 := s.get1dForCards(cg)
	ret := poker.NewCardGroup(nil)
	ret.AddCardGroup(d1)
	ret.AddCardGroup(d2)
	return ret
}

//获取count张散牌
func (s *SSSTable) GetSP(count int, cg *poker.CardGroup) *poker.CardGroup {
	if cg.Len() < count {
		mlog.Debug("剩余牌：%v", cg.Len())
		return nil
	}
	ret := poker.NewCardGroup(nil)
	m := make(map[int32]int32)
	for i := 0; i < cg.Len(); i++ {
		if _, ok := m[cg.Cards[i].CardNumber]; ok {
			m[cg.Cards[i].CardNumber] += 1
		} else {
			m[cg.Cards[i].CardNumber] = 1
		}
	}

	if count == 3 {
		for _, v := range cg.Cards {
			ret.AddCard(v)
			if ret.Len() >= count {
				break
			}
		}
	} else {
		isHeadHaveSameMaxNumberCard := false // 頭墩可能有相同點數牌
		// 先取一張最大
		ret.AddCard(cg.GetMaxCard())
		// 先取多張牌
		for i := cg.Len() - 1; i >= 0; i-- {
			if m[cg.Cards[i].CardNumber] >= 2 {
				if len(ret.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 {
					ret.AddCard(cg.Cards[i])
				}
				if ret.GetMaxCard().CardNumber == cg.Cards[i].CardNumber {
					isHeadHaveSameMaxNumberCard = true
				}
				if ret.Len() >= count {
					break
				}
			}
		}

		// 頭墩最大張有相同點數，找剩下最大的散牌
		if isHeadHaveSameMaxNumberCard && ret.Len() < count {
			cg.SortMaxFirst()
			var maxSanCard *poker.Card = nil
			for _, v := range cg.Cards {
				if len(ret.GetPokerCardByNumber(v.CardNumber)) == 0 {
					if maxSanCard == nil {
						maxSanCard = v
					} else if maxSanCard.CardNumber < v.CardNumber {
						maxSanCard = v
					}
				}
			}
			if maxSanCard != nil {
				ret.AddCard(maxSanCard)
			}
		}
		// 數量不足再放其他牌
		if ret.Len() < count {
			for i := cg.Len() - 1; i >= 0; i-- {
				if m[cg.Cards[i].CardNumber] < 2 && len(ret.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 {
					ret.AddCard(cg.Cards[i])
					if ret.Len() >= count {
						break
					}
				}
			}
		}
	}
	if ret != nil {
		cg.DeleteGroup(ret)
	}
	return ret
}

//抓取最小的單牌
func (s *SSSTable) getSPforCards(cards *poker.CardGroup, count int, cg *poker.CardGroup) *poker.CardGroup {
	if cg.Len() < count {
		mlog.Debug("剩余牌：%v", cg.Len())
		return nil
	}
	ret := poker.NewCardGroup(nil)

	m := make(map[int32]int32)
	for i := 0; i < cg.Len(); i++ {
		if _, ok := m[cg.Cards[i].CardNumber]; ok {
			m[cg.Cards[i].CardNumber] += 1
		} else {
			m[cg.Cards[i].CardNumber] = 1
		}
	}
	var countOnePair = 0
	for _, v := range m {
		if v >= 2 {
			countOnePair++
		}
	}
	// 依照對子數量，使用不同散牌取牌策略
	if countOnePair <= 3 {
		// 先取單張牌
		for i := cg.Len() - 1; i >= 0; i-- {
			if m[cg.Cards[i].CardNumber] < 2 &&
				len(cards.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 &&
				len(ret.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 {
				ret.AddCard(cg.Cards[i])
				if ret.Len() >= count {
					break
				}
			}
		}

		// 數量不足再放其他牌
		if ret.Len() < count {
			for i := cg.Len() - 1; i >= 0; i-- {
				if m[cg.Cards[i].CardNumber] >= 2 &&
					len(cards.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 &&
					len(ret.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 {
					ret.AddCard(cg.Cards[i])
					if ret.Len() >= count {
						break
					}
				}
			}
		}
	} else {
		// 先取多張牌
		for i := cg.Len() - 1; i >= 0; i-- {
			if m[cg.Cards[i].CardNumber] >= 2 &&
				len(cards.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 &&
				len(ret.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 {
				ret.AddCard(cg.Cards[i])
				if ret.Len() >= count {
					break
				}
			}
		}

		// 數量不足再放其他牌
		if ret.Len() < count {
			for i := cg.Len() - 1; i >= 0; i-- {
				if m[cg.Cards[i].CardNumber] < 2 &&
					len(cards.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 &&
					len(ret.GetPokerCardByNumber(cg.Cards[i].CardNumber)) == 0 {
					ret.AddCard(cg.Cards[i])
					if ret.Len() >= count {
						break
					}
				}

			}
		}
	}

	cg.DeleteGroup(ret)
	return ret
}

func (s *SSSTable) GetSpType(userId int32) int32 {
	min := int32(netproto.SSSDunCardsSpecialType_TongHuaShiSanShui)
	max := int32(netproto.SSSDunCardsSpecialType_SanShunZi)
	ur := s.GetUserByUserID(userId)
	u := s.getRoleInfoBySeat(ur.SeatNo)
	for i := min; i <= max; i++ {
		if s.verifyUserCardType(ur, u.act.recommandCt, true, i) {
			u.act.sct = i
			u.act.isSpecial = true
			return i
		}
	}
	return 0
}

func (s *SSSTable) GetSpTypeByPoker(userId int32) int32 {
	min := int32(netproto.SSSDunCardsSpecialType_TongHuaShiSanShui)
	max := int32(netproto.SSSDunCardsSpecialType_SanShunZi)
	ur := s.GetUserByUserID(userId)
	for i := min; i <= max; i++ {
		if s.checkSpecialType(ur, nil, true, i) {
			return i
		}
	}
	return 0
}

//随机获取一个非倒水的牌型
func (s *SSSTable) RadomCardType(seat int32) {
	role := s.getRoleInfoBySeat(seat)
	upk := s.gamePoker.GetUserPoker(role.roleid)
	cg := poker.NewCardGroup(upk.Cards.Cards)
	scg := CardGroupSortNumber{cg, false}
	sort.Sort(&scg)
	role.act = NewAutoCardType()
	//ur := s.GetUserByUserID(role.roleid)//默认为特殊牌型 注释部分
	//if s.GetSpType(ur.UserID) == 0 {
	for i := SSSDunMax - 1; i >= 0; i-- {
		c := 0
		if i == SSSDunHead {
			c = 3
		} else {
			c = 5
		}
		ret := s.GetTongHuaShun(c, scg.CardGroup)
		if ret != nil {
			role.act.cgs[i] = ret
			if c == 3 {
				role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
			} else {
				role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_TongHuaShun)
			}
			continue
		}
		ret = s.Get4(cg)
		if ret != nil {
			role.act.cgs[i] = ret
			role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_ZhaDan)
			role.act.cgs[i].AddCardGroup(s.getSPforCards(role.act.cgs[i], 1, cg)) //抓取最小的牌
			continue
		}
		ret = s.GetHL(cg)
		if ret != nil {
			role.act.cgs[i] = ret
			role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_HuLu)
			continue
		}
		ret = s.GetTH(c, cg, false)
		if ret != nil {
			role.act.cgs[i] = ret
			if c == 3 {
				role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
			} else {
				role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_TongHua)
			}
			continue
		}
		ret = s.GetSZ(c, cg)
		if ret != nil {
			role.act.cgs[i] = ret
			if c == 3 {
				role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
			} else {
				role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_ShunZi)
			}

			continue
		}
		ret = s.Get3(cg)
		if ret != nil {
			role.act.cgs[i] = ret
			if c != 3 {
				role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanTiao)
				role.act.cgs[i].AddCardGroup(s.getSPforCards(role.act.cgs[i], 2, cg))
			} else {
				role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_ChongSan)
			}

			continue
		}
		ret = s.Get2d(cg)
		if ret != nil {
			role.act.cgs[i] = ret
			role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_LiangDui)
			role.act.cgs[i].AddCardGroup(s.getSPforCards(role.act.cgs[i], 1, cg))
			continue
		}
		ret = s.Get1d(cg)
		if ret != nil {
			role.act.cgs[i] = ret
			role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_YiDui)
			if i != 0 {
				role.act.cgs[i].AddCardGroup(s.getSPforCards(role.act.cgs[i], 3, cg))
			} else {
				role.act.cgs[i].AddCardGroup(s.getSPforCards(role.act.cgs[i], 1, cg))
			}
			continue
		}
		if i == SSSDunHead {
			role.act.cgs[i] = s.GetSP(3, cg)
			mlog.Debug("获取三种单牌：%v", role.act.cgs[i])
			role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
		} else {
			role.act.cgs[i] = s.GetSP(5, cg)
			role.act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
		}
	}
	role.act.Cg2Hc()
}
