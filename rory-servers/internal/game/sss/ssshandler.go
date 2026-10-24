package sss

import (
	"sort"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

//注册各阶段超时处理函数
func (s *SSSTable) regist() {
	//游戏开始发牌播放动画等
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepBegin, s.gameBeginTimeOut)
	//牌型确定超时阶段
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepCardTypeChoose, s.CardTypeChooseTimeOut)
	//组牌超时阶段
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepCombine, s.CombineCardTimeOut)
	//比牌开始动画
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepAnimotion, s.AnimotionBegin)
	//头墩比较超时阶段
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepCompareH, s.CompareHTimeOut)
	//中墩比较超时阶段
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepCompareB, s.CompareBTimeOut)
	//后墩比较超时阶段
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepCompareF, s.CompareFTimeOut)
	//特殊牌型比较超时阶段
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepCompareSpecial, s.CompareSpecialTimeOut)
	//打枪阶段
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepUserGun, s.UserGunTimeOut)
	//结算超时阶段
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepSettlement, s.SettlementTimeOut)
	//游戏结束
	s.RegistCallBackFnc(netproto.SSSGameStep_SSSGameStepEnd, s.GameEnd)
}

//是否直接进入特殊牌型比较阶段
func (s *SSSTable) IsCompareDirect() bool {
	return false
	// n := 0
	// for _, v := range s.roles {
	// 	if !v.act.isSpecial {
	// 		n++
	// 	}
	// }

	// if n >= s.GetCurrentGameUserCount() {
	// 	return true
	// }

	// return false
}

//比牌开始动画
func (s *SSSTable) AnimotionBegin() {
	mlog.Debug("比牌开始动画")
	if s.IsCompareDirect() {
		mlog.Debug("跳过特殊比牌阶段")
		s.SetJump(netproto.SSSGameStep_SSSGameStepCompareSpecial)
	}

	if s.GetSpacialCount() >= int32(s.GetCurrentGameUserCount()) {
		s.SetJump(netproto.SSSGameStep_SSSGameStepCompareH)
	}
	s.stopLoop()
	s.animotionCur = s.RegisterAnimation(AnimotinBeginCompare, animotionTimeOut, s.ResumeLoop)
}

//头墩比牌超时回调
func (s *SSSTable) CompareHTimeOut() {
	mlog.Debug("比较头墩牌超时")
	s.CompareNormal()
	if s.GetSpacialCount() >= int32(s.GetCurrentGameUserCount()) {
		mlog.Debug("跳过此阶段")
		s.SetJump(netproto.SSSGameStep_SSSGameStepCompareB)
	} else {
		s.broadcastUserCards(netproto.SSSDunType_DunTypeHead)
		s.animotionCur = s.RegisterAnimation(AnimotionHeadComb, animotionTimeOut, s.ResumeLoop)
	}
	s.stopLoop()
}

//中墩比牌超时回调
func (s *SSSTable) CompareBTimeOut() {
	mlog.Debug("比较中墩牌超时")
	if s.GetSpacialCount() >= int32(s.GetCurrentGameUserCount()) {
		mlog.Debug("跳过此阶段")
		s.SetJump(netproto.SSSGameStep_SSSGameStepCompareF)
	} else {
		s.broadcastUserCards(netproto.SSSDunType_DunTypeBody)
		s.animotionCur = s.RegisterAnimation(AnimotionBodyComb, animotionTimeOut, s.ResumeLoop)
	}
	s.stopLoop()
}

//后墩比牌超时回调
func (s *SSSTable) CompareFTimeOut() {
	mlog.Debug("比较后墩牌超时")
	s.stopLoop()
	if s.GetSpacialCount() < int32(s.GetCurrentGameUserCount()) {
		s.broadcastUserCards(netproto.SSSDunType_DunTypeFoot)
		s.animotionCur = s.RegisterAnimation(AnimotionFootComb, animotionTimeOut, s.ResumeLoop)
	}
}

func (s *SSSTable) specialWinTree() {
	size := len(s.roles)
	for i := 0; i < size; i++ {
		if s.roles[i].act.isSpecial && len(s.roles[i].act.sctWinInfo) == 3 {
			s.roles[i].act.sctWin *= 2
			for j := 0; j < 3; j++ {
				r := s.getRoleInfoBySeat(s.roles[i].act.sctWinInfo[j].seatNo)
				r.act.sctLose -= s.roles[i].act.sctWinInfo[j].win
			}
		}
	}
}

//特殊牌型比较超时回调
func (s *SSSTable) CompareSpecialTimeOut() {
	mlog.Debug("比较特殊牌型阶段超时")
	if s.GetSpacialCount() > 0 {
		s.CompareSpeacil()
		s.CompareSpeacilAndNormal()
		s.specialWinTree()
		s.broadcastSpecialCompareResult()
		s.stopLoop()
		s.animotionCur = s.RegisterAnimation(AnimotionSpecialComb, 10, s.ResumeLoop)
		mlog.Debug("比较特殊牌型阶段超时")
	}
}

//打枪阶段
func (s *SSSTable) UserGunTimeOut() {
	mlog.Debug("进入打枪阶段")
	if !s.IsCompareDirect() {
		for _, v := range s.roles {
			if len(v.gun) > 0 {
				s.broadcastUserGun()
				s.stopLoop()
				s.animotionCur = s.RegisterAnimation(AnimotionGun, animotionTimeOut, s.ResumeLoop)
				break
			}
		}
	}
}

//结算超时回调
func (s *SSSTable) SettlementTimeOut() {
	mlog.Debug("进入结算阶段")
	s.CaculateOddsAndGun()
	s.broadCastSettlement()
	s.stopLoop()
	s.animotionCur = s.RegisterAnimation(AnimotionEnd, animotionTimeOut, s.ResumeLoop)
}

func (s *SSSTable) GameEnd() {
	s.GameEndFunc()
}

func cc(c1 *netproto.SCard, c2 *netproto.SCard) int32 {
	if c1.GetNumber() > c2.GetNumber() {
		return 1
	} else if c1.GetNumber() == c2.GetNumber() {
		return 0
	}
	return -1
}

func getSanTiao(c []*netproto.SCard) int32 {
	l := len(c)
	cardNumber1 := int32(0)
	count := 0
	for i := 0; i < l-2; i++ {
		if c[i].GetNumber() == c[i+1].GetNumber() {
			cardNumber1 = c[i].GetNumber()
			count++
		} else {
			count = 0
		}

		if count == 2 {
			break
		}
	}

	return cardNumber1
}

func getMaxDuiZi(c []*netproto.SCard) int32 {
	l := len(c)
	cardNumber1 := int32(0)
	for i := l - 1; i > 0; i-- {
		if c[i].GetNumber() == c[i-1].GetNumber() {
			cardNumber1 = c[i].GetNumber()
			break
		}
	}

	return cardNumber1
}

func getMinDuiZi(c []*netproto.SCard) int32 {
	l := len(c)
	cardNumber1 := int32(0)
	for i := 0; i < l-2; i++ {
		if c[i].GetNumber() == c[i+1].GetNumber() {
			cardNumber1 = c[i].GetNumber()
			break
		}
	}

	return cardNumber1
}

func compareHuLu(c1 []*netproto.SCard, c2 []*netproto.SCard) int32 {
	cardNumber1 := getSanTiao(c1)
	cardNumber2 := getSanTiao(c2)
	mlog.Debug("compareHuLu:c1=%v-c2=%v", c1, c2)
	mlog.Debug("compareHuLu:cardNumber1=%v,cardNumber2=%v", cardNumber1, cardNumber2)
	if cardNumber1 > cardNumber2 {
		return 1
	} else if cardNumber1 == cardNumber2 {
		return 0
	}
	return -1
}

func compareDuiZi(c1 []*netproto.SCard, c2 []*netproto.SCard) int32 {
	cardNumber1 := getMaxDuiZi(c1)
	cardNumber2 := getMaxDuiZi(c2)

	mlog.Debug("对子比较c1：%v,c2:%v", c1, c2)
	mlog.Debug("获取的对子牌为：cardNumber1=%v，cardNumber2=%v", cardNumber1, cardNumber2)
	if cardNumber1 > cardNumber2 {
		return 1
	} else if cardNumber1 == cardNumber2 {
		cardNumber3 := getMinDuiZi(c1)
		cardNumber4 := getMinDuiZi(c2)
		if cardNumber3 > cardNumber4 {
			return 1
		} else if cardNumber3 == cardNumber4 {
			return 0
		}

		return -1
	}
	return -1
}

func compareDun(d1 *netproto.Dun, d2 *netproto.Dun) int32 {
	if d1.GetCt() != d2.GetCt() {
		return -1 //墩类型不匹配 不能进行对比
	}

	if d1.GetCt() == int32(netproto.SSSDunCardsNomalType_HuLu) {
		return compareHuLu(d1.Cards, d2.Cards)
	}

	if d1.GetCt() == int32(netproto.SSSDunCardsNomalType_LiangDui) || d1.GetCt() == int32(netproto.SSSDunCardsNomalType_YiDui) {
		r := compareDuiZi(d1.Cards, d2.Cards)
		if r != 0 {
			return r
		}
	}

	l := len(d1.Cards)
	for i := l - 1; i >= 0; i-- {
		r := cc(d1.Cards[i], d2.Cards[i])
		if r != 0 {
			return r
		}
	}
	return 0
}

//dt 墩类型
func (s *SSSTable) CompareNormal() {
	size := len(s.roles)
	for i := 0; i < size; i++ {
		for n := 0; n < 3; n++ {
			s.roles[i].ris[n].Ct = netproto.SSSDunCardsNomalType(s.roles[i].act.ct[n])
		}
	}
	for dt := 0; dt < 3; dt++ {
		for i := 0; i < size; i++ {
			if s.roles[i].act.isSpecial == true {
				continue
			}
			for j := i + 1; j < size; j++ {
				if s.roles[j].act.isSpecial == true {
					continue
				}
				var shui int32 = 1
				if (dt + 1) == int(netproto.SSSDunType_DunTypeHead) { //头墩
					if s.roles[i].act.recommandCt.GetHead().GetCt() > s.roles[j].act.recommandCt.GetHead().GetCt() {
						if s.roles[j].act.recommandCt.GetHead().GetCt() == int32(netproto.SSSDunCardsNomalType_ChongSan) {
							shui = s.GetNsOdds(ChongSan)
						}
						s.roles[i].ris[dt].result[s.roles[j].seat] -= shui
						s.roles[j].ris[dt].result[s.roles[i].seat] += shui
					} else if s.roles[i].act.recommandCt.GetHead().GetCt() == s.roles[j].act.recommandCt.GetHead().GetCt() {
						r := compareDun(s.roles[i].act.recommandCt.GetHead(), s.roles[j].act.recommandCt.GetHead())
						if s.roles[i].act.recommandCt.GetHead().GetCt() == int32(netproto.SSSDunCardsNomalType_ChongSan) {
							shui = s.GetNsOdds(ChongSan)
						}
						if r > 0 {
							s.roles[i].ris[dt].result[s.roles[j].seat] += shui
							s.roles[j].ris[dt].result[s.roles[i].seat] -= shui
						} else if r < 0 {
							s.roles[i].ris[dt].result[s.roles[j].seat] -= shui
							s.roles[j].ris[dt].result[s.roles[i].seat] += shui
						} else {
							s.roles[i].ris[dt].result[s.roles[j].seat] -= 0
							s.roles[j].ris[dt].result[s.roles[i].seat] += 0
						}
					} else {
						if s.roles[i].act.recommandCt.GetHead().GetCt() == int32(netproto.SSSDunCardsNomalType_ChongSan) {
							shui = s.GetNsOdds(ChongSan)
						}
						s.roles[i].ris[dt].result[s.roles[j].seat] += shui
						s.roles[j].ris[dt].result[s.roles[i].seat] -= shui
					}
				} else if (dt + 1) == int(netproto.SSSDunType_DunTypeBody) { //中墩
					if s.roles[i].act.recommandCt.GetBody().GetCt() > s.roles[j].act.recommandCt.GetBody().GetCt() {
						if s.roles[j].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
							shui = s.GetNsOdds(ZhongDunTongHuaShun)
						} else if s.roles[j].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
							shui = s.GetNsOdds(ZhongDunZhaDan)
						} else if s.roles[j].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_HuLu) {
							shui = s.GetNsOdds(ZhongDunHuLu)
						}
						s.roles[i].ris[dt].result[s.roles[j].seat] -= shui
						s.roles[j].ris[dt].result[s.roles[i].seat] += shui
					} else if s.roles[i].act.recommandCt.GetBody().GetCt() == s.roles[j].act.recommandCt.GetBody().GetCt() {
						if s.roles[j].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
							shui = s.GetNsOdds(ZhongDunTongHuaShun)
						} else if s.roles[j].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
							shui = s.GetNsOdds(ZhongDunZhaDan)
						} else if s.roles[j].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_HuLu) {
							shui = s.GetNsOdds(ZhongDunHuLu)
						}
						r := compareDun(s.roles[i].act.recommandCt.GetBody(), s.roles[j].act.recommandCt.GetBody())
						if r > 0 {
							s.roles[i].ris[dt].result[s.roles[j].seat] += shui
							s.roles[j].ris[dt].result[s.roles[i].seat] -= shui
						} else if r < 0 {
							s.roles[i].ris[dt].result[s.roles[j].seat] -= shui
							s.roles[j].ris[dt].result[s.roles[i].seat] += shui
						} else {
							s.roles[i].ris[dt].result[s.roles[j].seat] -= 0
							s.roles[j].ris[dt].result[s.roles[i].seat] += 0
						}
					} else {
						if s.roles[i].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
							shui = s.GetNsOdds(ZhongDunTongHuaShun)
						} else if s.roles[i].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
							shui = s.GetNsOdds(ZhongDunZhaDan)
						} else if s.roles[i].act.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_HuLu) {
							shui = s.GetNsOdds(ZhongDunHuLu)
						}
						s.roles[i].ris[dt].result[s.roles[j].seat] += shui
						s.roles[j].ris[dt].result[s.roles[i].seat] -= shui
					}
				} else { //后墩
					if s.roles[i].act.recommandCt.GetFoot().GetCt() > s.roles[j].act.recommandCt.GetFoot().GetCt() {
						if s.roles[j].act.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
							shui = s.GetNsOdds(DiDunTongHuaShun)
						} else if s.roles[j].act.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
							shui = s.GetNsOdds(DiDunZhaDan)
						}
						s.roles[i].ris[dt].result[s.roles[j].seat] -= shui
						s.roles[j].ris[dt].result[s.roles[i].seat] += shui
					} else if s.roles[i].act.recommandCt.GetFoot().GetCt() == s.roles[j].act.recommandCt.GetFoot().GetCt() {
						if s.roles[j].act.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
							shui = s.GetNsOdds(DiDunTongHuaShun)
						} else if s.roles[j].act.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
							shui = s.GetNsOdds(DiDunZhaDan)
						}
						r := compareDun(s.roles[i].act.recommandCt.GetFoot(), s.roles[j].act.recommandCt.GetFoot())
						if r > 0 {
							s.roles[i].ris[dt].result[s.roles[j].seat] += shui
							s.roles[j].ris[dt].result[s.roles[i].seat] -= shui
						} else if r < 0 {
							s.roles[i].ris[dt].result[s.roles[j].seat] -= shui
							s.roles[j].ris[dt].result[s.roles[i].seat] += shui
						}
					} else {
						if s.roles[i].act.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
							shui = s.GetNsOdds(DiDunTongHuaShun)
						} else if s.roles[i].act.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
							shui = s.GetNsOdds(DiDunZhaDan)
						}
						s.roles[i].ris[dt].result[s.roles[j].seat] += shui
						s.roles[j].ris[dt].result[s.roles[i].seat] -= shui
					}
				}
			}
		}
	}
	//检测是否打枪
	for _, v := range s.roles {
		if v.act.isSpecial {
			continue
		}
		for i := 0; i < len(s.roles); i++ {
			if v.seat == s.roles[i].seat {
				continue
			}
			isGun := true
			n := 0
			for dt := 0; dt < 3; dt++ {
				if v.ris[int32(dt)].result[s.roles[i].seat] < 0 {
					isGun = false
					mlog.Debug("没有打枪：%v", v.seat)
					break
				}

				if v.ris[int32(dt)].result[s.roles[i].seat] == 0 {
					n++
				}
				mlog.Debug("seat:%v,打枪信息：%v", v.seat, v.ris[int32(dt)].result[s.roles[i].seat])
			}

			if n == 3 {
				isGun = false
			}
			if isGun {
				v.gun = append(v.gun, s.roles[i].seat)
			}
		}
	}
}

func (s *SSSTable) GetSpacialCount() int32 {
	var n int32 = 0
	size := len(s.roles)
	for i := 0; i < size; i++ {
		if s.roles[i].act.isSpecial {
			n++
		}
	}
	return n
}

//特殊牌型比较
func (s *SSSTable) CompareSpeacil() {
	size := len(s.roles)
	for i := 0; i < size; i++ {
		if !s.roles[i].act.isSpecial {
			mlog.Debug("s.roles[%v]:%v不是特殊牌型", i, s.roles[i].roleid)
			continue
		}

		for j := i + 1; j < size; j++ {
			if !s.roles[j].act.isSpecial {
				mlog.Debug("s.roles[%v]:%v不是特殊牌型", j, s.roles[j].roleid)
				continue
			}

			info := new(WinInfo)
			if s.roles[i].act.sct < s.roles[j].act.sct {
				s.roles[i].act.sctWin += s.getSpecialShui(s.roles[i].act.sct)
				s.roles[j].act.sctLose -= s.getSpecialShui(s.roles[i].act.sct)
				info.seatNo = s.roles[j].seat
				info.win = s.getSpecialShui(s.roles[i].act.sct)
				s.roles[i].act.sctWinInfo = append(s.roles[i].act.sctWinInfo, info)
			} else if s.roles[i].act.sct > s.roles[j].act.sct {
				s.roles[i].act.sctLose -= s.getSpecialShui(s.roles[j].act.sct)
				s.roles[j].act.sctWin += s.getSpecialShui(s.roles[j].act.sct)
				info.seatNo = s.roles[i].seat
				info.win = s.getSpecialShui(s.roles[j].act.sct)
				s.roles[j].act.sctWinInfo = append(s.roles[j].act.sctWinInfo, info)
			}

			mlog.Debug("获取到的特殊水数:%v", info.win)
			mlog.Debug("userid[%v].sctWin=%v---userid[%v].sctLose=%v", s.roles[i].seat, s.roles[i].act.sctWin, s.roles[j].seat, s.roles[j].act.sctLose)
		}
	}
}

//特殊牌型与普通牌型之间的比较
func (s *SSSTable) CompareSpeacilAndNormal() {
	size := len(s.roles)
	for i := 0; i < size; i++ {
		if !s.roles[i].act.isSpecial {
			continue
		}
		for j := 0; j < size; j++ {
			if s.roles[j].act.isSpecial {
				continue
			}
			s.roles[i].act.sctWin += s.getSpecialShui(s.roles[i].act.sct)
			s.roles[j].act.sctLose -= s.getSpecialShui(s.roles[i].act.sct)

			info := new(WinInfo)
			info.seatNo = s.roles[j].seat
			info.win = s.getSpecialShui(s.roles[i].act.sct)
			s.roles[i].act.sctWinInfo = append(s.roles[i].act.sctWinInfo, info)
		}
	}
}

//组牌超时回调
func (s *SSSTable) CombineCardTimeOut() {
	mlog.Debug("组牌超时回调")
}

//选择是特殊牌型还是普通牌型超时回调
func (s *SSSTable) CardTypeChooseTimeOut() {
	mlog.Debug("牌型确定阶段超时")
	for _, v := range s.roles {
		if v.act == nil {
			s.RadomCardType(v.seat)
		}
	}
	if s.IsCompareDirect() {
		s.SetJump(netproto.SSSGameStep_SSSGameStepAnimotion)
	}
}

//推荐牌型消息
func (s *SSSTable) getSuggestCardMsg(decks []*netproto.HandCards) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SSS_GameMessageClassID_SuggestCardId)
	gsdata := new(netproto.SuggestHandCards)
	gsdata.Hcs = decks

	pmsg.MsgData = gsdata
	return pmsg
}

//创建牌型消息超时回调
func (s *SSSTable) getPokerMsg(upk *poker.UserPoker) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_UserPokerID)
	pmsg.MsgData = upk.GetNetMessageData()
	return pmsg
}

func (s *SSSTable) ResumeLoop() {
	s.restartLoop()
	s.logicEngine()
}

//给扑克排一次序
func (s *SSSTable) sortPokerByNumber(cards *poker.CardGroup) {
	if cards == nil {
		return
	}
	arrlen := len(cards.Cards)
	for i := 0; i < arrlen; i++ {
		for j := 0; j+1 < arrlen-i; j++ {
			card1 := cards.Cards[j].CardNumber*10 + cards.Cards[j].Color
			card2 := cards.Cards[j+1].CardNumber*10 + cards.Cards[j+1].Color
			if card1 > card2 {
				cards.Swap(j, j+1)
			}
		}
	}
}

// isSameDeck
func (s *SSSTable) isSameDeck(deck1 *netproto.HandCards, deck2 *netproto.HandCards) bool {
	var isSame1, isSame2, isSame3 bool
	for _, v1 := range deck1.GetHead().Cards {
		isSame1 = false
		for _, v2 := range deck2.GetHead().Cards {
			if v1.GetColor() == v2.GetColor() && v1.GetNumber() == v2.GetNumber() {
				isSame1 = true
				break
			}
		}
	}
	for _, v1 := range deck1.GetBody().Cards {
		isSame2 = false
		for _, v2 := range deck2.GetBody().Cards {
			if v1.GetColor() == v2.GetColor() && v1.GetNumber() == v2.GetNumber() {
				isSame2 = true
				break
			}
		}
	}
	for _, v1 := range deck1.GetFoot().Cards {
		isSame3 = false
		for _, v2 := range deck2.GetFoot().Cards {
			if v1.GetColor() == v2.GetColor() && v1.GetNumber() == v2.GetNumber() {
				isSame3 = true
				break
			}
		}
	}
	return isSame1 && isSame2 && isSame3
}

func (s *SSSTable) createEmptyDeck() *netproto.HandCards {
	result := new(netproto.HandCards)
	result.Head = new(netproto.Dun)
	result.Head.Cards = make([]*netproto.SCard, 0)

	result.Body = new(netproto.Dun)
	result.Body.Cards = make([]*netproto.SCard, 0)

	result.Foot = new(netproto.Dun)
	result.Foot.Cards = make([]*netproto.SCard, 0)
	return result
}

// 與尾墩同牌型不同牌，中墩牌型較大或尾墩牌型較大
func (s *SSSTable) findNextDeckExtra(cg *poker.CardGroup, deck *netproto.HandCards) (nextDeck *netproto.HandCards) {
	// var headType = netproto.SSSDunCardsNomalType(deck.GetHead().GetCt())
	// var bodyType = netproto.SSSDunCardsNomalType(deck.GetBody().GetCt())
	var selectCard = poker.NewCardGroup(cg.Cards)
	var footType = netproto.SSSDunCardsNomalType(deck.GetFoot().GetCt())
	// if headType == netproto.SSSDunCardsNomalType_SanPai && bodyType == netproto.SSSDunCardsNomalType_SanPai && footType == netproto.SSSDunCardsNomalType_SanPai {
	// 	return nil
	// }
	footCards := new(poker.CardGroup)
	for _, v := range deck.GetFoot().GetCards() {
		footCards.AddCard(&poker.Card{Color: v.GetColor(), CardNumber: v.GetNumber()})
	}
	cardType, list := s.checkCards(selectCard, footCards)
	if len(list) > 0 && cardType == footType {
		var cardType2, cardType3 netproto.SSSDunCardsNomalType
		var cardGroups1, cardGroups2, cardGroups3 *poker.CardGroup = nil, nil, nil
		for i := 0; i < len(list); i++ {
			cardGroups1 = list[i]
			nextDeck = new(netproto.HandCards)
			// fmt.Println(cardGroups1)

			//選擇尾墩
			selectCard.DeleteGroup(cardGroups1)
			nextDeck.Foot = new(netproto.Dun)
			nextDeck.Foot.Ct = proto.Int32(int32(footType))
			nextDeck.Foot.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeFoot))
			nextDeck.Foot.Cards = make([]*netproto.SCard, 0)
			for _, v := range cardGroups1.Cards {
				card := new(netproto.SCard)
				card.Color = proto.Int32(v.Color)
				card.Number = proto.Int32(v.CardNumber)
				nextDeck.Foot.Cards = append(nextDeck.Foot.Cards, card)
			}

			//選擇中墩
			cardType2, cardGroups2 = s.selectCards(SSSDunBody, selectCard, nil)
			if cardGroups2 == nil {
				continue
			}
			// fmt.Println(2, cardType2, cardGroups2)
			nextDeck.Body = new(netproto.Dun)
			nextDeck.Body.Ct = proto.Int32(int32(cardType2))
			nextDeck.Body.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeBody))
			nextDeck.Body.Cards = make([]*netproto.SCard, 0)
			for _, v := range cardGroups2.Cards {
				card := new(netproto.SCard)
				card.Color = proto.Int32(v.Color)
				card.Number = proto.Int32(v.CardNumber)
				nextDeck.Body.Cards = append(nextDeck.Body.Cards, card)
			}

			// 如果尾墩跟中墩牌型相同且倒水，做交換
			if footType == cardType2 {
				if GetWeightByDun(nextDeck.GetFoot()) < GetWeightByDun(nextDeck.GetBody()) {
					nextDeck.Body.Cards, nextDeck.Foot.Cards = nextDeck.Foot.Cards, nextDeck.Body.Cards
				}
			}

			//選擇頭墩
			cardType3, cardGroups3 = s.selectCards(SSSDunHead, selectCard, nil)
			if cardGroups3 == nil {
				continue
			}
			// fmt.Println(3, cardType3, cardGroups3)
			nextDeck.Head = new(netproto.Dun)
			nextDeck.Head.Ct = proto.Int32(int32(cardType3))
			nextDeck.Head.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeHead))
			nextDeck.Head.Cards = make([]*netproto.SCard, 0)
			for _, v := range cardGroups3.Cards {
				card := new(netproto.SCard)
				card.Color = proto.Int32(v.Color)
				card.Number = proto.Int32(v.CardNumber)
				nextDeck.Head.Cards = append(nextDeck.Head.Cards, card)
			}

			if nextDeck != nil && netproto.SSSDunCardsNomalType(nextDeck.GetHead().GetCt()) >= netproto.SSSDunCardsNomalType(nextDeck.GetBody().GetCt()) {
				// fmt.Println(nextDeck)
				return nextDeck
			}
		}
	}

	return nil
}

// findNextDeck1 去除尾墩最大牌型找第二大牌型
func (s *SSSTable) findNextDeck1(cg *poker.CardGroup, deck *netproto.HandCards) (nextDeck *netproto.HandCards) {
	// var headType = netproto.SSSDunCardsNomalType(deck.GetHead().GetCt())
	// var bodyType = netproto.SSSDunCardsNomalType(deck.GetBody().GetCt())
	var footType = netproto.SSSDunCardsNomalType(deck.GetFoot().GetCt())
	// if headType == netproto.SSSDunCardsNomalType_SanPai && bodyType == netproto.SSSDunCardsNomalType_SanPai && footType == netproto.SSSDunCardsNomalType_SanPai {
	// 	return nil
	// }

	var deck1 *netproto.HandCards

	if footType != netproto.SSSDunCardsNomalType_SanPai {
		// 去除尾墩最大牌型找第二大牌型
		deck1 = s.cardTypeGroup(cg, map[SSSDunType]netproto.SSSDunCardsNomalType{
			SSSDunBody: footType,
			SSSDunFoot: footType,
		})
	}

	if deck1 != nil && netproto.SSSDunCardsNomalType(deck1.GetHead().GetCt()) >= netproto.SSSDunCardsNomalType(deck1.GetBody().GetCt()) {
		return deck1
	} else {
		return nil
	}
}

// findNextDeck2 去除中墩最大牌型找第二大牌型
func (s *SSSTable) findNextDeck2(cg *poker.CardGroup, deck *netproto.HandCards) (nextDeck *netproto.HandCards) {
	// var headType = netproto.SSSDunCardsNomalType(deck.GetHead().GetCt())
	var bodyType = netproto.SSSDunCardsNomalType(deck.GetBody().GetCt())
	var footType = netproto.SSSDunCardsNomalType(deck.GetFoot().GetCt())
	// if headType == netproto.SSSDunCardsNomalType_SanPai && bodyType == netproto.SSSDunCardsNomalType_SanPai && footType == netproto.SSSDunCardsNomalType_SanPai {
	// 	return nil
	// }

	if footType != netproto.SSSDunCardsNomalType_TongHuaShun {
		footType = footType - 1
	}

	var deck2 *netproto.HandCards

	if bodyType != netproto.SSSDunCardsNomalType_SanPai {
		// 去除中墩最大牌型找第二大牌型
		deck2 = s.cardTypeGroup(cg, map[SSSDunType]netproto.SSSDunCardsNomalType{
			SSSDunHead: bodyType,
			SSSDunBody: bodyType,
			SSSDunFoot: footType,
		})
	}

	if deck2 != nil && netproto.SSSDunCardsNomalType(deck2.GetHead().GetCt()) >= netproto.SSSDunCardsNomalType(deck2.GetBody().GetCt()) {
		return deck2
	} else {
		return nil
	}
}

// findNextDeck3 //去除頭墩最大牌型找第二大牌型
func (s *SSSTable) findNextDeck3(cg *poker.CardGroup, deck *netproto.HandCards) (nextDeck *netproto.HandCards) {
	var headType = netproto.SSSDunCardsNomalType(deck.GetHead().GetCt())
	var bodyType = netproto.SSSDunCardsNomalType(deck.GetBody().GetCt())
	var footType = netproto.SSSDunCardsNomalType(deck.GetFoot().GetCt())

	if headType == netproto.SSSDunCardsNomalType_SanPai {
		return nil
	}

	retDeck := new(netproto.HandCards)
	retDeck.Foot = new(netproto.Dun)
	retDeck.Foot.Ct = proto.Int32(int32(footType))
	retDeck.Foot.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeFoot))
	retDeck.Foot.Cards = make([]*netproto.SCard, 0)
	for _, v := range deck.GetFoot().GetCards() {
		card := new(netproto.SCard)
		card.Color = proto.Int32(v.GetColor())
		card.Number = proto.Int32(v.GetNumber())
		retDeck.Foot.Cards = append(retDeck.Foot.Cards, card)
	}

	retDeck.Body = new(netproto.Dun)
	retDeck.Body.Ct = proto.Int32(int32(bodyType))
	retDeck.Body.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeBody))
	retDeck.Body.Cards = make([]*netproto.SCard, 0)
	for _, v := range deck.GetBody().GetCards() {
		card := new(netproto.SCard)
		card.Color = proto.Int32(v.GetColor())
		card.Number = proto.Int32(v.GetNumber())
		retDeck.Body.Cards = append(retDeck.Body.Cards, card)
	}

	retDeck.Head = new(netproto.Dun)
	retDeck.Head.Ct = proto.Int32(int32(headType))
	retDeck.Head.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeHead))
	retDeck.Head.Cards = make([]*netproto.SCard, 0)
	for _, v := range deck.GetHead().GetCards() {
		card := new(netproto.SCard)
		card.Color = proto.Int32(v.GetColor())
		card.Number = proto.Int32(v.GetNumber())
		retDeck.Head.Cards = append(retDeck.Head.Cards, card)
	}

	if headType == netproto.SSSDunCardsNomalType_YiDui { // 頭墩的對子跟其他墩的散牌作交換
		headCards := retDeck.GetHead().GetCards()
		m := make(map[int32]int32)
		for _, v := range headCards {
			if _, ok := m[v.GetNumber()]; ok {
				m[v.GetNumber()] += 1
			} else {
				m[v.GetNumber()] = 1
			}
		}

		bodyCards := retDeck.GetBody().GetCards()
		m2 := make(map[int32]int32)
		for _, v := range bodyCards {
			if _, ok := m2[v.GetNumber()]; ok {
				m2[v.GetNumber()] += 1
			} else {
				m2[v.GetNumber()] = 1
			}
		}

		footCards := retDeck.GetFoot().GetCards()
		m3 := make(map[int32]int32)
		for _, v := range footCards {
			if _, ok := m3[v.GetNumber()]; ok {
				m3[v.GetNumber()] += 1
			} else {
				m3[v.GetNumber()] = 1
			}
		}

		if footType == netproto.SSSDunCardsNomalType_SanTiao || footType == netproto.SSSDunCardsNomalType_LiangDui || footType == netproto.SSSDunCardsNomalType_YiDui { //
			for _, v := range headCards {
				for _, v2 := range footCards {
					if m[v.GetNumber()] > 1 && m3[v2.GetNumber()] < 2 {
						// swap card
						tempNumber := proto.Int32(v2.GetNumber())
						tempColor := proto.Int32(v2.GetColor())
						v2.Number = proto.Int32(v.GetNumber())
						v2.Color = proto.Int32(v.GetColor())
						v.Number = tempNumber
						v.Color = tempColor
						retDeck.GetHead().Ct = proto.Int32(int32(netproto.SSSDunCardsNomalType_SanPai))
						return retDeck
					}
				}
			}
		}

		if bodyType == netproto.SSSDunCardsNomalType_SanTiao || bodyType == netproto.SSSDunCardsNomalType_LiangDui || bodyType == netproto.SSSDunCardsNomalType_YiDui { //
			for _, v := range headCards {
				for _, v2 := range bodyCards {
					if m[v.GetNumber()] > 1 && m2[v2.GetNumber()] < 2 {
						// swap card
						tempNumber := proto.Int32(v2.GetNumber())
						tempColor := proto.Int32(v2.GetColor())
						v2.Number = proto.Int32(v.GetNumber())
						v2.Color = proto.Int32(v.GetColor())
						v.Number = tempNumber
						v.Color = tempColor
						retDeck.GetHead().Ct = proto.Int32(int32(netproto.SSSDunCardsNomalType_SanPai))
						return retDeck
					}
				}
			}
		}
	}
	return nil
}

// containCardType 是否包含排除牌組
func (s *SSSTable) containExclusions(exclusions map[SSSDunType]netproto.SSSDunCardsNomalType, dunType SSSDunType, cardType netproto.SSSDunCardsNomalType) bool {
	if exclusions == nil {
		return false
	}

	for k, v := range exclusions {
		if k == dunType && v >= cardType {
			return true
		}
	}
	return false
}

// 檢查是否有最大的五張特殊牌型(只檢查:同花順, 炸彈, 葫蘆, 同花, 順子)
func (s *SSSTable) checkCards(source *poker.CardGroup, compareCards *poker.CardGroup) (cardType netproto.SSSDunCardsNomalType, maxCardTypeList []*poker.CardGroup) {
	weightList := make([]*WeightTable, 0)
	maxCardTypeList = make([]*poker.CardGroup, 0)
	// 同花順
	list := TongHuaShunPatterns()
	for _, v := range list {
		if source.IsContent(v.Cards.Cards) {
			if !v.Cards.IsEqual(compareCards) {
				weight := new(WeightTable)
				weight.Cards = poker.NewCardGroup(v.Cards.Cards)
				weight.Weight = v.Weight
				weightList = append(weightList, weight)
			}
		}
	}
	if len(weightList) > 0 {
		sort.SliceStable(weightList, func(i, j int) bool {
			return weightList[i].Weight > weightList[j].Weight
		})
		for _, v := range weightList {
			maxCardTypeList = append(maxCardTypeList, poker.NewCardGroup(v.Cards.Cards))
		}
		return netproto.SSSDunCardsNomalType_TongHuaShun, maxCardTypeList
	}

	// 炸彈(4支)
	list = ZhaDanPatterns()
	for _, v := range list {
		if source.IsContent(v.Cards.Cards) {
			if !v.Cards.IsEqual(compareCards) {
				weight := new(WeightTable)
				weight.Cards = poker.NewCardGroup(v.Cards.Cards)
				weight.Weight = v.Weight
				weightList = append(weightList, weight)
			}
		}
	}
	if len(weightList) > 0 {
		sort.SliceStable(weightList, func(i, j int) bool {
			return weightList[i].Weight > weightList[j].Weight
		})
		for _, v := range weightList {
			maxCardTypeList = append(maxCardTypeList, poker.NewCardGroup(v.Cards.Cards))
		}
		return netproto.SSSDunCardsNomalType_ZhaDan, maxCardTypeList
	}

	// 葫蘆
	list = HuluPatterns()
	for _, v := range list {
		if source.IsContent(v.Cards.Cards) {
			if !v.Cards.IsEqual(compareCards) {
				weight := new(WeightTable)
				weight.Cards = poker.NewCardGroup(v.Cards.Cards)
				weight.Weight = v.Weight
				weightList = append(weightList, weight)
			}
		}
	}
	if len(weightList) > 0 {
		sort.SliceStable(weightList, func(i, j int) bool {
			return weightList[i].Weight > weightList[j].Weight
		})
		for _, v := range weightList {
			maxCardTypeList = append(maxCardTypeList, poker.NewCardGroup(v.Cards.Cards))
		}
		return netproto.SSSDunCardsNomalType_HuLu, maxCardTypeList
	}

	// 同花
	list = TongHuaPatterns()
	for _, v := range list {
		if source.IsContent(v.Cards.Cards) {
			if !v.Cards.IsEqual(compareCards) {
				weight := new(WeightTable)
				weight.Cards = poker.NewCardGroup(v.Cards.Cards)
				weight.Weight = v.Weight
				weightList = append(weightList, weight)
			}
		}
	}
	if len(weightList) > 0 {
		sort.SliceStable(weightList, func(i, j int) bool {
			return weightList[i].Weight > weightList[j].Weight
		})
		for _, v := range weightList {
			maxCardTypeList = append(maxCardTypeList, poker.NewCardGroup(v.Cards.Cards))
		}
		return netproto.SSSDunCardsNomalType_TongHua, maxCardTypeList
	}

	// 順子
	list = ShunZiPatterns()
	for _, v := range list {
		if source.IsContent(v.Cards.Cards) {
			if !v.Cards.IsEqual(compareCards) {
				weight := new(WeightTable)
				weight.Cards = poker.NewCardGroup(v.Cards.Cards)
				weight.Weight = v.Weight
				weightList = append(weightList, weight)
			}
		}
	}
	if len(weightList) > 0 {
		sort.SliceStable(weightList, func(i, j int) bool {
			return weightList[i].Weight > weightList[j].Weight
		})
		for _, v := range weightList {
			maxCardTypeList = append(maxCardTypeList, poker.NewCardGroup(v.Cards.Cards))
		}
		return netproto.SSSDunCardsNomalType_ShunZi, maxCardTypeList
	}

	return netproto.SSSDunCardsNomalType_SanPai, nil //只是回傳個值，並不一定真的是散牌
}

// selectCards 根據 dunType 跟來源牌組，取得牌型跟牌組
func (s *SSSTable) selectCards(dunType SSSDunType, source *poker.CardGroup, exclusions map[SSSDunType]netproto.SSSDunCardsNomalType) (cardType netproto.SSSDunCardsNomalType, cardGroups *poker.CardGroup) {

	c := 5
	if dunType == SSSDunHead {
		c = 3
	}
	// 同花順
	if !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_TongHuaShun) {
		cardGroups = s.GetTongHuaShun(c, source)
		if cardGroups != nil {
			if dunType == SSSDunHead {
				cardType = netproto.SSSDunCardsNomalType_SanPai
			} else {
				cardType = netproto.SSSDunCardsNomalType_TongHuaShun
			}
			// fmt.Println("同花順", cardGroups)
			return
		}
	}

	// 炸彈(4支)
	if !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_ZhaDan) {
		cardGroups = s.Get4(source)
		if cardGroups != nil {
			cardGroups.AddCardGroup(s.getSPforCards(cardGroups, 1, source)) //抓取最小的牌
			cardType = netproto.SSSDunCardsNomalType_ZhaDan
			// fmt.Println("炸彈", cardGroups)
			return
		}
	}

	// 葫蘆
	if !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_HuLu) {
		cardGroups = s.GetHL(source)
		if cardGroups != nil {
			cardType = netproto.SSSDunCardsNomalType_HuLu
			// fmt.Println("葫蘆", cardGroups)
			return
		}
	}

	// 同花
	if !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_TongHua) {
		cardGroups = s.GetTH(c, source, false)
		if cardGroups != nil {
			if dunType == SSSDunHead {
				cardType = netproto.SSSDunCardsNomalType_SanPai
			} else {
				cardType = netproto.SSSDunCardsNomalType_TongHua
			}
			// fmt.Println("同花", cardGroups)
			return
		}
	}

	// 順子
	if !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_ShunZi) {
		cardGroups = s.GetSZ(c, source)

		if cardGroups != nil {
			if dunType == SSSDunHead {
				cardType = netproto.SSSDunCardsNomalType_SanPai
			} else {
				cardType = netproto.SSSDunCardsNomalType_ShunZi
			}
			// fmt.Println("順子", cardGroups)
			return
		}
	}

	// 三條
	if !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_SanTiao) || !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_ChongSan) && dunType == SSSDunHead {
		cardGroups = s.Get3(source)
		if cardGroups != nil {
			if dunType == SSSDunHead {
				cardType = netproto.SSSDunCardsNomalType_ChongSan // 衝三
			} else {
				cardGroups.AddCardGroup(s.getSPforCards(cardGroups, 2, source))
				cardType = netproto.SSSDunCardsNomalType_SanTiao
			}
			// fmt.Println("三條cardGroups", cardGroups)
			// fmt.Println("三條source", source)
			return
		}
	}

	// 2對
	if !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_LiangDui) {
		cardGroups = s.Get2d(source)
		if cardGroups != nil {
			cardGroups.AddCardGroup(s.getSPforCards(cardGroups, 1, source))
			cardType = netproto.SSSDunCardsNomalType_LiangDui
			// fmt.Println("2對", cardGroups)
			return
		}
	}

	// 1對
	if !s.containExclusions(exclusions, dunType, netproto.SSSDunCardsNomalType_YiDui) {
		cardGroups = s.Get1d(source)
		if cardGroups != nil {
			cardType = netproto.SSSDunCardsNomalType_YiDui
			if c != 3 {
				cardGroups.AddCardGroup(s.getSPforCards(cardGroups, 3, source))
			} else {
				cardGroups.AddCardGroup(s.getSPforCards(cardGroups, 1, source))
			}
			// fmt.Println("1對", cardGroups)
			return
		}
	}

	// 散牌
	if dunType == SSSDunHead {
		// fmt.Println("散牌source", source)
		cardGroups = s.GetSP(3, source)
		mlog.Debug("获取三种单牌：%v", cardGroups)
		// fmt.Println("散牌", cardGroups)
		cardType = netproto.SSSDunCardsNomalType_SanPai
	} else {
		cardGroups = s.GetSP(5, source)
		cardType = netproto.SSSDunCardsNomalType_SanPai
	}
	// fmt.Println("散牌", cardGroups)
	return
}

// cardTypeGroup 取得可能牌型 cardTypeGroup []netproto.SSSDunCardsNomalType
func (s *SSSTable) cardTypeGroup(group *poker.CardGroup, exclusions map[SSSDunType]netproto.SSSDunCardsNomalType) (deck *netproto.HandCards) {

	deck = new(netproto.HandCards)

	cg := poker.NewCardGroup(group.Cards)
	// fmt.Println("group", cg)
	// sortGroup := CardGroupSortNumber{cg, false}
	sort.Sort(cg)

	var cardType1, cardType2, cardType3 netproto.SSSDunCardsNomalType
	var cardGroups1, cardGroups2, cardGroups3 *poker.CardGroup = nil, nil, nil

	//選擇尾墩
	// fmt.Println("取得可能牌型", exclusions)
	cardType1, cardGroups1 = s.selectCards(SSSDunFoot, cg, exclusions)
	if cardGroups1 == nil {
		return
	}
	// fmt.Println(1, cardType1, cardGroups1)
	deck.Foot = new(netproto.Dun)
	deck.Foot.Ct = proto.Int32(int32(cardType1))
	deck.Foot.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeFoot))
	deck.Foot.Cards = make([]*netproto.SCard, 0)
	for _, v := range cardGroups1.Cards {
		card := new(netproto.SCard)
		card.Color = proto.Int32(v.Color)
		card.Number = proto.Int32(v.CardNumber)
		deck.Foot.Cards = append(deck.Foot.Cards, card)
	}

	//選擇中墩
	cardType2, cardGroups2 = s.selectCards(SSSDunBody, cg, exclusions)
	if cardGroups2 == nil {
		return
	}
	// fmt.Println(2, cardType2, cardGroups2)
	deck.Body = new(netproto.Dun)
	deck.Body.Ct = proto.Int32(int32(cardType2))
	deck.Body.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeBody))
	deck.Body.Cards = make([]*netproto.SCard, 0)
	for _, v := range cardGroups2.Cards {
		card := new(netproto.SCard)
		card.Color = proto.Int32(v.Color)
		card.Number = proto.Int32(v.CardNumber)
		deck.Body.Cards = append(deck.Body.Cards, card)
	}

	// 如果尾墩跟中墩牌型相同且倒水，做交換
	if cardType1 == cardType2 {
		if GetWeightByDun(deck.GetFoot()) < GetWeightByDun(deck.GetBody()) {
			// fmt.Println(GetWeightByDun(deck.GetFoot()), GetWeightByDun(deck.GetBody()))
			deck.Body.Cards, deck.Foot.Cards = deck.Foot.Cards, deck.Body.Cards
		}
	}

	//選擇頭墩
	cardType3, cardGroups3 = s.selectCards(SSSDunHead, cg, exclusions)
	if cardGroups3 == nil {
		return
	}
	// fmt.Println(3, cardType3, cardGroups3)
	deck.Head = new(netproto.Dun)
	deck.Head.Ct = proto.Int32(int32(cardType3))
	deck.Head.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeHead))
	deck.Head.Cards = make([]*netproto.SCard, 0)
	for _, v := range cardGroups3.Cards {
		card := new(netproto.SCard)
		card.Color = proto.Int32(v.Color)
		card.Number = proto.Int32(v.CardNumber)
		deck.Head.Cards = append(deck.Head.Cards, card)
	}
	// fmt.Println("")
	// fmt.Println("QQQQ", deck.Head.Cards, deck.Body.Cards, deck.Foot.Cards)
	return
}

//发牌
func (s *SSSTable) faPai(count int) {
	s.poker.Shuffle()
	s.stopLoop()
	//	hasSpecialType := false
	for _, ur := range s.AllGameUser() {
		cards := s.poker.Deal(ur.UserID, count)
		s.sortPokerByNumber(cards)
		upk := poker.NewUserPoker(ur, cards)
		s.gamePoker.AddUserPoker(upk)
		s.SendMessageToUser(ur, s.getPokerMsg(upk))
		// 優化 改成由server送推薦牌型給client 原本為client自己組
		if !ur.IsRobot() {

			var deckTree = new(SSSDeckTree)
			deckTree.NewDeckTree(s, cards, 4) // 最多4片葉子
			deckTree.GrowUp(s)
			deckTree.Print()

			var suggestDecks = make([]*netproto.HandCards, 0)
			for _, v := range deckTree.GetSortNodes() {
				suggestDecks = append(suggestDecks, v.GetHC())
			}

			s.SendMessageToUser(ur, s.getSuggestCardMsg(suggestDecks))
		} else {
			// var deckTree = new(SSSDeckTree)
			// deckTree.NewDeckTree(s, cards, 0) // 最多4片葉子
			// deckTree.Print()
		}

		s.GameLog.Fapai(ur.UserID, cards.String())
	}

	s.animotionCur = s.RegisterAnimation(AnimotionFaPai, animotionTimeOut, s.ResumeLoop)
}

//游戏开始超时回调
func (s *SSSTable) gameBeginTimeOut() {
	s.faPai(13)
	mlog.Debug("游戏开始阶段超时")
}

//牌型选择请求
func (s *SSSTable) OnChooseSpecialCardType(user *gamelogic.GameUser, data []byte) {
	if s.step != netproto.SSSGameStep_SSSGameStepCardTypeChoose && s.step != netproto.SSSGameStep_SSSGameStepCombine {
		s.responseChooseSpecialCardTypeRet(1, "不是牌型选择阶段", user.SeatNo, 0)
		mlog.Error("不是牌型选择阶段 seat=%d", user.SeatNo)
		return
	}
	ust := &netproto.ChooseSpecialCardType{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析玩家出牌消息错误,%v", err)
		s.KickUserFunc(user.UserID, "非法操作")
		return
	}
	u := s.getRoleInfoBySeat(user.SeatNo)
	u.act = NewAutoCardType()
	//	s.choosespecialcount++
	if ust.GetSpecial() == false {
		s.broadcastUserCardTypeChoose(user, false)
		if ust.GetHcs() != nil {
			//mlog.Debug("玩家选择普通牌型，获取推荐牌型%s", ust.GetHcs().String())
			if !s.verifyUserCardType(user, ust.GetHcs(), false, 0) {
				//mlog.Error("489玩家发送上来的手牌错误,%v", err)
				s.KickUserFunc(user.UserID, "玩家发送上来的手牌错误,非法操作")
				u.act = nil
				return
			}
			//後台deallog不能顯示 < > 內字元
			// desc := fmt.Sprintf("玩家%d选择普通牌型", user.UserID)
			// desc += fmt.Sprintf("Sdt:%v Ct:%v ", ust.GetHcs().GetHead().GetSdt(), ust.GetHcs().GetHead().GetCt())
			// for _, v := range ust.GetHcs().GetHead().Cards {
			// 	desc += fmt.Sprintf("(Number:%v Color:%v) ", v.GetNumber(), v.GetColor())
			// }
			// desc += fmt.Sprintf(", Sdt:%v Ct:%v ", ust.GetHcs().GetBody().GetSdt(), ust.GetHcs().GetBody().GetCt())
			// for _, v := range ust.GetHcs().GetBody().Cards {
			// 	desc += fmt.Sprintf("(Number:%v Color:%v) ", v.GetNumber(), v.GetColor())
			// }
			// desc += fmt.Sprintf(", Sdt:%v Ct:%v ", ust.GetHcs().GetFoot().GetSdt(), ust.GetHcs().GetFoot().GetCt())
			// for _, v := range ust.GetHcs().GetFoot().Cards {
			// 	desc += fmt.Sprintf("(Number:%v Color:%v) ", v.GetNumber(), v.GetColor())
			// }

			// s.GameLog.AddGameStep(desc)
			// fmt.Println(desc)

			u.act.recommandCt = ust.GetHcs()
			u.act.isSpecial = false
			u.act.Hc2Cgs()
		}
	} else {
		if !s.verifyUserCardType(user, u.act.recommandCt, true, ust.GetSpecialType()) {
			//mlog.Error("499玩家发送上来的手牌错误,%v", err)
			s.KickUserFunc(user.UserID, "玩家发送上来的手牌错误,非法操作")
			u.act = nil
			return
		}
		s.GameLog.AddGameStep("玩家%d选择特殊牌型%d", user.UserID, ust.GetSpecialType())
		s.broadcastUserCardTypeChoose(user, true)
		//mlog.Debug("玩家选择特殊牌型")
		u.act.isSpecial = true
		u.act.sct = ust.GetSpecialType()
		u.isZuPai = true
	}
	//其他阶段为推荐牌型只有特殊牌型才判断超时
	if s.step == netproto.SSSGameStep_SSSGameStepCardTypeChoose {
		u.isDo = true
		jumpstep := true
		//4个都选择了特殊牌型
		for _, v := range s.roles {
			if v.act == nil || !v.act.isSpecial {
				jumpstep = false
				break
			}
		}

		if jumpstep {
			//跳过组牌阶段
			s.SetJump(netproto.SSSGameStep_SSSGameStepCombine)
		}

		if s.checkAllRoleDo() {
			s.timeOutCallBack()
		}
	}
}

func (s *SSSTable) CheckSpecailDo() bool {
	n := 0
	for _, v := range s.roles {
		if v.isDo || (v.act != nil && v.act.isSpecial) {
			n++
			continue
		}

		if s.step == netproto.SSSGameStep_SSSGameStepCombine {
			if v.isZuPai {
				n++
				mlog.Debug("有一个组牌完成")
			}
		}
	}

	mlog.Debug("n=%v,len(s.roles)= %v", n, len(s.roles))
	return n >= len(s.roles)
}

//组牌完成请求
func (s *SSSTable) OnCompleteCombine(user *gamelogic.GameUser, data []byte) {
	if s.step != netproto.SSSGameStep_SSSGameStepCombine {
		s.responseCompleteCardCombRet(user, 1, "不是组牌阶段")
		mlog.Error("不是组牌阶段 seat=%d", user.SeatNo)
		return
	}
	ust := &netproto.CompleteCardComb{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析玩家出牌消息错误")
		s.KickUserFunc(user.UserID, "非法操作")
		return
	}
	u := s.getRoleInfoBySeat(user.SeatNo)
	if !s.verifyUserCardType(user, ust.GetHcs(), ust.GetIsSpecial(), ust.GetSt()) {
		mlog.Error("提交的牌型异常")
		return
	}
	//保存客户端发上来的牌型
	for i := 0; i < 3; i++ {
		u.act.cgs[i] = poker.NewCardGroup(nil)
	}
	for _, v := range ust.GetHcs().Head.Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		u.act.cgs[0].AddCard(c)
	}
	// s.GameLog.AddGameStep("玩家%d选择头敦牌型%s", u.roleid, u.act.cgs[0].String())
	u.act.ct[0] = ust.GetHcs().Head.GetCt()
	u.act.Cg2Dun(ust.GetHcs().Head, 0, 3)
	for _, v := range ust.GetHcs().Body.Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		u.act.cgs[1].AddCard(c)
	}

	// s.GameLog.AddGameStep("玩家%d选择中敦牌型%s", u.roleid, u.act.cgs[1].String())
	u.act.ct[1] = ust.GetHcs().Body.GetCt()
	u.act.Cg2Dun(ust.GetHcs().Body, 1, 5)
	for _, v := range ust.GetHcs().Foot.Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		u.act.cgs[2].AddCard(c)
	}

	// s.GameLog.AddGameStep("玩家%d选择后敦牌型%s", u.roleid, u.act.cgs[2].String())
	u.act.ct[2] = ust.GetHcs().Foot.GetCt()
	u.act.Cg2Dun(ust.GetHcs().Foot, 2, 5)
	if ust.GetIsSpecial() {
		u.act.isSpecial = true
		u.act.sct = ust.GetSt()
	}
	mlog.Debug("dun types =%v", u.act.ct)
	u.act.Cg2Hc()
	s.responseCompleteCardCombRet(user, 0, "Ok")
	s.broadcastUserCardTypeChoose(user, ust.GetIsSpecial())
	u.isDo = true
	u.isZuPai = true
	if s.CheckSpecailDo() {
		s.timeOutCallBack()
		//機器人邏輯在此階段結束後會多跑一次GameMessage(SendGameMessageNow)
	}
}

//断线重连
func (s *SSSTable) UserRecover(gameuser *gamelogic.GameUser) {
	s.DynamicGameLogicBase.UserRecover(gameuser)
	s.PushGameCfg(gameuser)
	if s.IsGaming() {
		sendmsg := false
		if s.step == netproto.SSSGameStep_SSSGameStepCombine || s.step == netproto.SSSGameStep_SSSGameStepCardTypeChoose {
			v := s.getRoleInfoBySeat(gameuser.SeatNo)
			if v != nil && v.isDo == false {
				sendmsg = true
			}

			if v != nil {
				v.isDo = true
				v.isZuPai = true
			}
		}
		s.PushGameInfo(gameuser)
		if sendmsg {
			s.broadcastUserCardTypeChoose(gameuser, false)
			if s.CheckSpecailDo() {
				s.timeOutCallBack()
			}
		}
	}
}

func (s *SSSTable) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	s.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)
	s.PushGameCfg(gameuser)
	if s.IsGaming() && !s.IsGameUser(gameuser) { //旁观
		s.PushGameInfo(gameuser)
	}
}

func (s *SSSTable) UserOffLine(gameuser *gamelogic.GameUser) {
	sendmsg := false
	if s.step == netproto.SSSGameStep_SSSGameStepCombine || s.step == netproto.SSSGameStep_SSSGameStepCardTypeChoose {
		v := s.getRoleInfoBySeat(gameuser.SeatNo)
		if v != nil && v.isDo == false {
			v.isDo = true
			v.isZuPai = true
			sendmsg = true
		}
	}

	if sendmsg {
		s.broadcastUserCardTypeChoose(gameuser, false)
		if s.CheckSpecailDo() {
			s.timeOutCallBack()
		}
	}

	s.DynamicGameLogicBase.UserOffLine(gameuser)
}

//客户端消息处理函数
func (s *SSSTable) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	//动画等待
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.GameMessageClassID_UserActionCommitID) {
		s.CheckUserAnimationCommited(user, gmsg, s.animotionCur, nil)
	}
	//选择牌型
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.SSS_GameMessageClassID_ChooseSpecialCardTypeId) {
		s.OnChooseSpecialCardType(user, gmsg.MsgData.([]byte))
	}
	//选择牌型
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.SSS_GameMessageClassID_CompleteCardCombId) {
		s.OnCompleteCombine(user, gmsg.MsgData.([]byte))
	}
}

func (s *SSSTable) ExchangeMaxCard(robotUserID int32) {
	s.LogicBase.ExchangeMaxCard(robotUserID)
	// 牌桌上每个机器人都会轮巡到这个 Method
	// 所以只需检查玩家的牌是否比机器人好
	// 如果比机器人好就跟玩家交换

	var playerUserIDList []int32

	for _, ur := range s.AllTableUser() {
		if !ur.IsRobot() {
			playerUserIDList = append(playerUserIDList, ur.UserID)
		}
	}

	var fpoker = s.poker.(*poker.FixedPoker)
	var robotCardGroup = fpoker.GetUserCards(robotUserID)
	var exchangeUserID = robotUserID
	for _, playerUserID := range playerUserIDList {
		var playerCardGroup = fpoker.GetUserCards(playerUserID)
		var playerSpecialType = s.getSpecialType(playerCardGroup)
		var robotSpecialType = s.getSpecialType(robotCardGroup)

		if robotSpecialType > playerSpecialType {
			exchangeUserID = robotUserID
			break
		}

		if playerSpecialType > 0 { //比较特殊牌型
			if playerSpecialType < robotSpecialType {
				exchangeUserID = playerUserID
			} else if playerSpecialType == robotSpecialType {

			}
		} else { // 普通牌型
			var playerAutoCardType = NewAutoCardType()
			s.getNormalType(playerCardGroup, playerAutoCardType)
			var robotAutoCardType = NewAutoCardType()
			s.getNormalType(robotCardGroup, robotAutoCardType)
			if s.compareNormal(playerAutoCardType, robotAutoCardType) {
				exchangeUserID = playerUserID
			}
		}
	}

	if exchangeUserID != robotUserID {
		example.invalid("交换用户%d和用户%d牌型", robotUserID, exchangeUserID)
		s.GameLog.AddGameStep("交换用户%d和用户%d牌型", robotUserID, exchangeUserID)
		fpoker.ExChange(robotUserID, exchangeUserID)
	}
}

func (s *SSSTable) getSpecialType(usercards *poker.CardGroup) int32 {
	if s.CheckTongHuaShiSanShui(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_TongHuaShiSanShui)
	}
	if s.CheckShiSanShui(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_ShiSanShui)
	}

	if s.CheckShiErHuangZu(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_ShiErHuangZu)
	}

	if s.CheckSanTongHuaShun(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_SanTongHuaShun)
	}

	if s.CheckSanFenTianXia(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_SanFenTianXia)
	}

	if s.CheckQuanDa(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_QuanDa)
	}

	if s.CheckQuanXiao(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_QuanXiao)
	}

	if s.CheckCouYiSe(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_CouYiSe)
	}

	if s.CheckSiTaoSanTiao(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_SiTaoSantiao)
	}

	if s.CheckWuDuiSanTiao(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_WuDuiSanTiao)
	}

	if s.CheckLiuDuiBan(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_LiuDuiBan)
	}

	if s.CheckSanTongHua(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_SanTongHua)
	}

	if s.CheckSanShunZi(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_SanShunZi)
	}

	if s.CheckShuangGuaiChongSan(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_ShuangGuaiChongSan)
	}

	if s.CheckSanChongJiaDan(usercards) {
		return int32(netproto.SSSDunCardsSpecialType_SanChongJiaDan)
	}

	return 0
}

func (s *SSSTable) getNormalType(usercards *poker.CardGroup, act *AutoCardType) {
	var deckTree = new(SSSDeckTree)
	deckTree.NewDeckTree(s, usercards, 0) // 最多4片葉子
	// deckTree.GetSortNodes()[0]
	// cg := poker.NewCardGroup(usercards.Cards)
	// scg := CardGroupSortNumber{cg, false}
	// sort.Sort(&scg)
	// for i := SSSDunMax - 1; i >= 0; i-- {
	// 	c := 0
	// 	if i == SSSDunHead {
	// 		c = 3
	// 	} else {
	// 		c = 5
	// 	}
	// 	ret := s.GetTongHuaShun(c, scg.CardGroup)
	// 	if ret != nil {
	// 		act.cgs[i] = ret
	// 		s.sortPokerByNumber(act.cgs[i])
	// 		if c == 3 {
	// 			act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
	// 		} else {
	// 			act.ct[i] = int32(netproto.SSSDunCardsNomalType_TongHuaShun)
	// 		}
	// 		continue
	// 	}
	// 	ret = s.Get4(cg)
	// 	if ret != nil {
	// 		act.cgs[i] = ret
	// 		act.ct[i] = int32(netproto.SSSDunCardsNomalType_ZhaDan)
	// 		act.cgs[i].AddCardGroup(s.getSPforCards(act.cgs[i], 1, cg))
	// 		continue
	// 	}
	// 	ret = s.GetHL(cg)
	// 	if ret != nil {
	// 		act.cgs[i] = ret
	// 		act.ct[i] = int32(netproto.SSSDunCardsNomalType_HuLu)
	// 		continue
	// 	}
	// 	tmpcg := poker.NewCardGroup(cg.Cards)
	// 	s.sortCardsByColor(tmpcg)
	// 	ret = s.GetTH(c, tmpcg, false)
	// 	if ret != nil {
	// 		act.cgs[i] = ret
	// 		cg.DeleteGroup(ret)
	// 		s.sortPokerByNumber(act.cgs[i])
	// 		if c == 3 {
	// 			act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
	// 		} else {
	// 			act.ct[i] = int32(netproto.SSSDunCardsNomalType_TongHua)
	// 		}
	// 		continue
	// 	}
	// 	ret = s.GetSZ(c, cg)
	// 	if ret != nil {
	// 		act.cgs[i] = ret
	// 		s.sortPokerByNumber(act.cgs[i])
	// 		if c == 3 {
	// 			act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
	// 		} else {
	// 			act.ct[i] = int32(netproto.SSSDunCardsNomalType_ShunZi)
	// 		}

	// 		continue
	// 	}
	// 	ret = s.Get3(cg)
	// 	if ret != nil {
	// 		act.cgs[i] = ret
	// 		if c != 3 {
	// 			act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanTiao)
	// 			cards := s.getSPforCards(act.cgs[i], 2, cg)
	// 			s.sortPokerByNumber(act.cgs[i])
	// 			act.cgs[i].AddCardGroup(cards)
	// 		} else {
	// 			act.ct[i] = int32(netproto.SSSDunCardsNomalType_ChongSan)
	// 		}

	// 		continue
	// 	}
	// 	ret = s.Get2d(cg)
	// 	if ret != nil {
	// 		act.cgs[i] = ret
	// 		s.sortPokerByNumber(act.cgs[i])
	// 		act.ct[i] = int32(netproto.SSSDunCardsNomalType_LiangDui)
	// 		act.cgs[i].AddCardGroup(s.getSPforCards(act.cgs[i], 1, cg))
	// 		continue
	// 	}

	// 	ret = s.Get1d(cg)
	// 	if ret != nil {
	// 		act.cgs[i] = ret
	// 		act.ct[i] = int32(netproto.SSSDunCardsNomalType_YiDui)
	// 		if i != 0 {
	// 			cards := s.getSPforCards(act.cgs[i], 3, cg)
	// 			s.sortPokerByNumber(act.cgs[i])
	// 			act.cgs[i].AddCardGroup(cards)
	// 		} else {
	// 			act.cgs[i].AddCardGroup(s.GetSP(1, cg))
	// 		}
	// 		continue
	// 	}

	// 	if i == SSSDunHead {
	// 		act.cgs[i] = s.GetSP(3, cg)
	// 		s.sortPokerByNumber(act.cgs[i])
	// 		act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
	// 	} else {
	// 		act.cgs[i] = s.GetSP(5, cg)
	// 		s.sortPokerByNumber(act.cgs[i])
	// 		act.ct[i] = int32(netproto.SSSDunCardsNomalType_SanPai)
	// 	}
	// }

	// for i := 0; i < 3; i++ {
	// 	mlog.Debug("顿类型%v", act.ct[i])
	// 	mlog.Debug("顿牌型%v", act.cgs[i])
	// 	fmt.Println(act.ct[i], act.cgs[i])
	// }
	var deck = deckTree.Root()

	// act.Cg2Hc()
	act.ct[SSSDunHead] = deck.data.GetHeadCardType()
	act.cgs[SSSDunHead] = deck.data.head
	act.ct[SSSDunBody] = deck.data.GetBodyCardType()
	act.cgs[SSSDunBody] = deck.data.body
	act.ct[SSSDunFoot] = deck.data.GetFootCardType()
	act.cgs[SSSDunFoot] = deck.data.foot

	// less := func(cards []*netproto.SCard, i int, j int){

	// }

	// headCards := v.GetHead().GetCards()
	// sort.SliceStable(deck.data.hc.Head.Cards, func(i, j int) bool {
	// 	var cg = new(poker.CardGroup)
	// 	for _, v := range deck.data.hc.Head.Cards {
	// 		cg.AddCard(&poker.PokerCard{Color: v.GetColor(), CardNumber: v.GetNumber()})
	// 	}
	// 	var number1 = deck.data.hc.Head.Cards[i].GetNumber()
	// 	var number2 = deck.data.hc.Head.Cards[j].GetNumber()
	// 	var count1 = len(cg.GetPokerCardByNumber(number1))
	// 	var count2 = len(cg.GetPokerCardByNumber(number2))
	// 	if count1 > count2 { // 牌多的牌前面
	// 		return true
	// 	} else if count1 < count2 { // 牌少的牌後面
	// 		return false
	// 	} else {
	// 		return number1 < number2 // 相同數量牌點數小的牌前面
	// 	}
	// })
	// fmt.Println(headCards)
	// bodyCards := v.GetBody().GetCards()
	// sort.SliceStable(deck.data.hc.Body.Cards, func(i, j int) bool {
	// 	var cg = new(poker.CardGroup)
	// 	for _, v := range deck.data.hc.Body.Cards {
	// 		cg.AddCard(&poker.PokerCard{Color: v.GetColor(), CardNumber: v.GetNumber()})
	// 	}
	// 	var number1 = deck.data.hc.Body.Cards[i].GetNumber()
	// 	var number2 = deck.data.hc.Body.Cards[j].GetNumber()
	// 	var count1 = len(cg.GetPokerCardByNumber(number1))
	// 	var count2 = len(cg.GetPokerCardByNumber(number2))
	// 	if count1 > count2 {
	// 		return true
	// 	} else if count1 < count2 {
	// 		return false
	// 	} else {
	// 		return number1 < number2
	// 	}
	// })
	// // fmt.Println(bodyCards)

	// // footCards := v.GetFoot().GetCards()
	// sort.SliceStable(deck.data.hc.Foot.Cards, func(i, j int) bool {
	// 	var cg = new(poker.CardGroup)
	// 	for _, v := range deck.data.hc.Foot.Cards {
	// 		cg.AddCard(&poker.PokerCard{Color: v.GetColor(), CardNumber: v.GetNumber()})
	// 	}
	// 	var number1 = deck.data.hc.Foot.Cards[i].GetNumber()
	// 	var number2 = deck.data.hc.Foot.Cards[j].GetNumber()
	// 	var count1 = len(cg.GetPokerCardByNumber(number1))
	// 	var count2 = len(cg.GetPokerCardByNumber(number2))
	// 	if count1 > count2 {
	// 		return true
	// 	} else if count1 < count2 {
	// 		return false
	// 	} else {
	// 		return number1 < number2
	// 	}
	// })
	act.recommandCt = deck.data.hc
	// act.recommandCt.Head.Sdt = proto.Int32(int32(netproto.SSSDunType_DunTypeHead))
}

func (s *SSSTable) compareNormal(act1 *AutoCardType, act2 *AutoCardType) bool {
	shui1 := int32(0)
	shui2 := int32(0)
	mlog.Debug("act1.recommandCt.Head.Cards=%v", act1.recommandCt.Head.Cards)
	mlog.Debug("act2.recommandCt.Head.Cards=%v", act2.recommandCt.Head.Cards)

	mlog.Debug("act1.recommandCt.Body.Cards=%v", act1.recommandCt.Body.Cards)
	mlog.Debug("act2.recommandCt.Body.Cards=%v", act2.recommandCt.Body.Cards)

	mlog.Debug("act1.recommandCt.Foot.Cards=%v", act1.recommandCt.Foot.Cards)
	mlog.Debug("act2.recommandCt.Foot.Cards=%v", act2.recommandCt.Foot.Cards)
	for dt := 0; dt < 3; dt++ {
		shui := int32(1)
		if dt == 0 {
			if act1.ct[dt] > act2.ct[dt] {
				if act1.ct[dt] == int32(netproto.SSSDunCardsNomalType_ChongSan) {
					shui = s.GetNsOdds(ChongSan)
				}

				mlog.Debug("shui2------------")
				shui2 += shui
				//shui2 -= shui
			} else if act1.ct[dt] < act2.ct[dt] {
				if act2.ct[dt] == int32(netproto.SSSDunCardsNomalType_ChongSan) {
					shui = s.GetNsOdds(ChongSan)
				}

				mlog.Debug("shui1------------")
				//shui1 -= shui
				shui1 += shui
			} else {
				r := compareDun(act1.recommandCt.GetHead(), act2.recommandCt.GetHead())
				if act1.recommandCt.GetHead().GetCt() == int32(netproto.SSSDunCardsNomalType_ChongSan) {
					shui = s.GetNsOdds(ChongSan)
				}

				if r > 0 {
					mlog.Debug("===shui1------------")
					shui1 += shui
					//shui2 -= shui
				} else if r < 0 {
					//shui1 -= shui
					mlog.Debug("===shui2------------")
					shui2 += shui
				}
			}
		} else if dt == 1 {
			//中墩
			if act1.recommandCt.GetBody().GetCt() > act2.recommandCt.GetBody().GetCt() {
				if act1.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
					shui = s.GetNsOdds(ZhongDunTongHuaShun)
				} else if act1.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
					shui = s.GetNsOdds(ZhongDunZhaDan)
				} else if act1.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_HuLu) {
					shui = s.GetNsOdds(ZhongDunHuLu)
				}

				mlog.Debug("shui2------------")
				shui2 += shui
				//shui2 -= shui
			} else if act1.recommandCt.GetBody().GetCt() == act2.recommandCt.GetBody().GetCt() {
				if act1.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
					shui = s.GetNsOdds(ZhongDunTongHuaShun)
				} else if act1.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
					shui = s.GetNsOdds(ZhongDunZhaDan)
				} else if act1.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_HuLu) {
					shui = s.GetNsOdds(ZhongDunHuLu)
				}
				r := compareDun(act1.recommandCt.GetBody(), act2.recommandCt.GetBody())
				if r > 0 {
					mlog.Debug("===shui1------------")
					shui1 += shui
					//shui2 -= shui
				} else if r < 0 {
					//shui1 -= shui
					mlog.Debug("=====shui2------------")
					shui2 += shui
				}
			} else {
				if act2.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
					shui = s.GetNsOdds(ZhongDunTongHuaShun)
				} else if act2.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
					shui = s.GetNsOdds(ZhongDunZhaDan)
				} else if act2.recommandCt.GetBody().GetCt() == int32(netproto.SSSDunCardsNomalType_HuLu) {
					shui = s.GetNsOdds(ZhongDunHuLu)
				}
				mlog.Debug("shui1------------")
				shui1 += shui
				//shui1 -= shui
			}
		} else {
			//后墩
			if act1.recommandCt.GetFoot().GetCt() > act2.recommandCt.GetFoot().GetCt() {
				if act1.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
					shui = s.GetNsOdds(DiDunTongHuaShun)
				} else if act1.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
					shui = s.GetNsOdds(DiDunZhaDan)
				}

				mlog.Debug("shui2------------")
				shui2 += shui
				//shui2 -= shui
			} else if act1.recommandCt.GetFoot().GetCt() == act2.recommandCt.GetFoot().GetCt() {
				if act1.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
					shui = s.GetNsOdds(DiDunTongHuaShun)
				} else if act1.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
					shui = s.GetNsOdds(DiDunZhaDan)
				}
				r := compareDun(act1.recommandCt.GetFoot(), act2.recommandCt.GetFoot())
				if r > 0 {
					mlog.Debug("====shui1------------")
					shui1 += shui
					//shui2 -= shui
				} else if r < 0 {
					mlog.Debug("====shui2------------")
					//shui1 -= shui
					shui2 += shui
				}
			} else {
				if act2.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_TongHuaShun) {
					shui = s.GetNsOdds(DiDunTongHuaShun)
				} else if act2.recommandCt.GetFoot().GetCt() == int32(netproto.SSSDunCardsNomalType_ZhaDan) {
					shui = s.GetNsOdds(DiDunZhaDan)
				}
				mlog.Debug("shui1------------")
				shui1 += shui
				//shui1 -= shui
			}
		}
	}

	mlog.Debug("shui1=%d,shui2=%d", shui1, shui2)
	return (shui1 > shui2)
}
