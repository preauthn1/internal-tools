package lhd

import (
	"example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"	
	"example.invalid/ServerTeam/Servers/pkg/util"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/mlog"	
)

import (
	"sort"
	"fmt"
)

func (lh *LongHudouLogic) updateVipUserWin() {
	bowBet := lh.getBowBet()
	superToolHandler := lh.SuperToolService.GetBrGameSuperToolHandler()
	refresh := false
	for _, ud := range bowBet {
		if vipuser := superToolHandler.GetVipUser(ud.user.UserID); vipuser != nil {
			vipuser.CompleteAmount = proto.Int64(vipuser.GetCompleteAmount() + ud.updatemoney)
			if vipuser.GetAssignAmount() > 0 {
				if vipuser.GetCompleteAmount() >= vipuser.GetAssignAmount() {
					vipuser.Complete = proto.Int32(1)
				}
			} else {
				if vipuser.GetCompleteAmount() <= vipuser.GetAssignAmount() {
					vipuser.Complete = proto.Int32(1)
				}
			}
			example.invalid("update user completeamount %d, userid %d", vipuser.GetCompleteAmount(), vipuser.GetUserID())
			superToolHandler.SaveVipUser(ud.user.UserID)
			refresh = true
		}
	}
	if refresh { //有異動才需要馬上更新
		lh.SuperToolService.RefreshVIPUser()
	}

}
//
//func (lh *LongHudouLogic) updateBowUserWin() {
//	bowBet := lh.getBowBet()
//	for _, ud := range bowBet {
//		if vipuser := lh.bowProvider.GetUser(ud.user.UserID); vipuser != nil {
//			vipuser.CompleteAmount = proto.Int64(vipuser.GetCompleteAmount() + ud.updatemoney)
//			if vipuser.GetAssignAmount() > 0 {
//				if vipuser.GetCompleteAmount() >= vipuser.GetAssignAmount() {
//					vipuser.Complete = proto.Int32(1)
//				}
//			} else {
//				if vipuser.GetCompleteAmount() <= vipuser.GetAssignAmount() {
//					vipuser.Complete = proto.Int32(1)
//				}
//			}
//			mlog.Debug("update user completeamount %d, userid %d", vipuser.GetCompleteAmount(), vipuser.GetUserID())
//
//			pmsg := new(msg.Message)
//			pmsg.BClassID = int32(netproto.MessageBClassID_Game)
//			pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdVipUserID)
//			pmsg.MsgData = vipuser
//			lh.SendRPCMessageFunc(true, pmsg)
//		}
//	}
//}

func (lh *LongHudouLogic) getAssignReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailID)
	pdata := new(netproto.SuperToolGetAssignRewardDetail)
	pdata.GameID = proto.Int32(lh.GetServerConfig().GameID)
	pdata.FlagID = proto.Int32(util.StrToInt32(lh.GetServerConfig().Flag))
	pdata.TableNo = proto.Int32(lh.TableNo())
	pmsg.MsgData = pdata

	lh.SendRPCMessageFunc(false, pmsg)
}

func (lh *LongHudouLogic) setAssignReward(m *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolAssignRewardDetail{}
		return mm
	})
	detail := rmm.(*netproto.SuperToolAssignRewardDetail)
	lh.assignRewardDetail = detail
}

func (lh *LongHudouLogic) isAssignReward() (bool, int32) {
	if lh.assignRewardDetail != nil {
		ret := util.StrToInt32(lh.assignRewardDetail.GetDescription())
		return true, ret
	}
	return false, 0
}

func (lh *LongHudouLogic) assignRewardDetailResult(result string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdAssignRewardDetailID)
	pdata := lh.assignRewardDetail
	pdata.Reward = proto.String(result)
	if lh.execAssignReward {
		pdata.Status = proto.Int32(2)
	} else {
		pdata.Status = proto.Int32(3)
	}
	pmsg.MsgData = pdata

	lh.SendRPCMessageFunc(true, pmsg)
}

func (lh *LongHudouLogic) saveLastGameUserBet() {
	lh.lastGameUserBet = lh.lastGameUserBet[:0] //clear last game user bet
	for _, ud := range lh.userdata {
		try := int64(0)
		for _, b := range ud.userbet {
			try += b.money
		}
		if try > 0 {
			lh.lastGameUserBet = append(lh.lastGameUserBet, ud)
			if userhis, ok := LhdUserHistoryStore[ud.user.UserID];ok {
				LhdUserHistoryStore[ud.user.UserID] = append([]*longhudouuserdata{ud}, userhis...)
			} else {
				LhdUserHistoryStore[ud.user.UserID] = []*longhudouuserdata{ud}
			}
			if len(LhdUserHistoryStore[ud.user.UserID]) > 100 {
				LhdUserHistoryStore[ud.user.UserID] = LhdUserHistoryStore[ud.user.UserID][:100]
			}
		}
	}
}


//獲取需要殺放分的用戶投注記錄
func (lh *LongHudouLogic) getBowBet() (vipBet map[int32]*longhudouuserdata){
	vipBet = make(map[int32]*longhudouuserdata)
	for _, ud := range lh.userdata{
		if vipuser := lh.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(ud.user.UserID); vipuser != nil {
			if vipuser.GetStatus() == 1 && vipuser.GetComplete() == 0 && vipuser.GetAssignAmount() != 0 {
				vipBet[ud.user.UserID] = ud
			}
		}
	}
	return
}
//監控放殺分
func (lh *LongHudouLogic) bowJiangType() bool{
	bowBet := lh.getBowBet()
	betlimit := lh.ctrlbetlimit
	//先執行放分
	for userid, ubet := range bowBet {
		vipuser := lh.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() < 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		if rnd < vipuser.GetAssignRate(){ //隨機0-99 小於勝率則執行放分
			//1. 獲取下注的項目
			betlong := ubet.getUserBetByType(LONG)
			bethu := ubet.getUserBetByType(HU)
			bethe := ubet.getUserBetByType(HE)

			//2. 決定開獎項目
			if betlong >= betlimit {
				lh.jiangtype = LONG
			} else if bethu >= betlimit {
				lh.jiangtype = HU
			}
			//3. 放分開和條件及機率
			if bethe >= betlimit {
				r := (vipuser.GetAssignAmount() * 15 / 10) - vipuser.GetCompleteAmount()
				hepay := ubet.getUserBetByType(HE) * int64(lh.betvalue.getOddByType(HE))
				if r > hepay {
					rnd := util.RndInt(0, 1000)
					if rnd < 75 {
						lh.jiangtype = HE
					}
				}
			}

			if lh.jiangtype != 0 { // 獲取到開獎就跳出
				break
			}
		}
	}
	if lh.jiangtype != 0 {
		example.invalid("執行放分開獎 : %v", lh.jiangtype)
		return true
	}
	//殺分
	for userid, ubet := range bowBet {
		vipuser := lh.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() > 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		if rnd >= vipuser.GetAssignRate(){ //隨機0-99 大於等於勝率則執行殺分
			//1. 獲取下注的項目
			betlong := ubet.getUserBetByType(LONG)
			bethu := ubet.getUserBetByType(HU)

			//2. 決定開獎項目
			if betlong >= betlimit {
				lh.jiangtype = HU
			} else if bethu >= betlimit {
				lh.jiangtype = LONG
			}
			if lh.jiangtype != 0 { // 獲取到開獎就跳出
				break
			}

		} else {
			continue
		}
	}
	if lh.jiangtype != 0 {
		example.invalid("執行殺分開獎 : %v", lh.jiangtype)
		return true
	}
	return false
}


// 獲取當局下注用戶列表
func (lh *LongHudouLogic) responseCurrentUserBetList() (list []*netproto.Http_LHD_UserInfo) {
	for _, ud := range lh.userdata {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_LHD_UserInfo)
		userinfo.UserID = proto.Int32(ud.user.UserID)
		ut := int32(1)
		if ud.user.IsRobot() {
			ut = 0
		}
		userinfo.UserType = proto.Int32(ut)
		userinfo.NickName = proto.String(ud.user.NickName)
		userinfo.GameScore = proto.Int64(ud.user.UpdateScore)
		userinfo.InitMoney = proto.Int64(ud.user.LoginScore)
		userinfo.CurrentMoney = proto.Int64(ud.user.Score)
		userinfo.BetInfoDragon = proto.Int64(ud.getUserBetByType(LONG))
		userinfo.BetInfoTiger = proto.Int64(ud.getUserBetByType(HU))
		userinfo.BetInfoDraw = proto.Int64(ud.getUserBetByType(HE))
		userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
		userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
		userinfo.WinRate = proto.Int64(ud.calUserWinRate(100))
		userinfo.BetTotal = proto.Int64(ud.getUserBetTotal())
		userinfo.IsVIP = proto.Bool(lh.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
		list = append(list, userinfo)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].GetIsVIP() && list[j].GetIsVIP() {
			return list[i].GetBetTotal() > list[j].GetBetTotal()
		} else if list[i].GetIsVIP() {
			return true
		} else if list[j].GetIsVIP() {
			return false
		}
		if list[i].GetBetTotal() == list[j].GetBetTotal() {
			return list[i].GetUserID() < list[j].GetUserID()
		}
		return list[i].GetBetTotal() > list[j].GetBetTotal()
	})
	return
}

// 獲取上局的下注用戶列表
func (lh *LongHudouLogic) responseLastGameUserBetList() (list []*netproto.Http_LHD_UserInfo) {
	tmparr := lh.lastGameUserBet
	if len(lh.lastGameUserBet) > 50 {
		tmparr = lh.lastGameUserBet[:50]
	}

	for _, ud := range tmparr {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_LHD_UserInfo)
		userinfo.UserID = proto.Int32(ud.user.UserID)
		ut := int32(1)
		if ud.user.IsRobot() {
			ut = 0
		}
		userinfo.UserType = proto.Int32(ut)
		userinfo.NickName = proto.String(ud.user.NickName)
		userinfo.GameScore = proto.Int64(ud.user.UpdateScore)
		userinfo.InitMoney = proto.Int64(ud.user.LoginScore)
		userinfo.CurrentMoney = proto.Int64(ud.user.Score)
		userinfo.BetInfoDragon = proto.Int64(ud.getUserBetByType(LONG))
		userinfo.BetInfoTiger = proto.Int64(ud.getUserBetByType(HU))
		userinfo.BetInfoDraw = proto.Int64(ud.getUserBetByType(HE))
		userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
		userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
		userinfo.WinRate = proto.Int64(ud.calUserWinRate(100))
		userinfo.BetTotal = proto.Int64(ud.getUserBetTotal())
		userinfo.IsVIP = proto.Bool(lh.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
		list = append(list, userinfo)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].GetIsVIP() && list[j].GetIsVIP() {
			return list[i].GetBetTotal() > list[j].GetBetTotal()
		} else if list[i].GetIsVIP() {
			return true
		} else if list[j].GetIsVIP() {
			return false
		}
		if list[i].GetBetTotal() == list[j].GetBetTotal() {
			return list[i].GetUserID() < list[j].GetUserID()
		}
		return list[i].GetBetTotal() > list[j].GetBetTotal()
	})
	return
}


func (lh *LongHudouLogic) getBetTypeWinRate() (rate []int32){

	list := lh.resultcardlist.GetCardList()
	size := len(list)
	if size == 0 {
		return []int32{0,0,0}
	}
	wincount := []int{0,0,0}
	for _, p := range list {
		foo := p.GetCards()
		longCard := foo[0]
		huCard := foo[1]
		if longCard.GetNumber() == huCard.GetNumber() {
			wincount[1]++
		} else if longCard.GetNumber() > huCard.GetNumber() {
			wincount[0]++
		} else {
			wincount[2]++
		}
	}
	rate = append(rate, int32( (float64(wincount[0]) / float64(size)) * 100 ))
	rate = append(rate, int32( (float64(wincount[1]) / float64(size)) * 100 ))
	rate = append(rate, int32( (float64(wincount[2]) / float64(size)) * 100 ))
	return
}


func (lh *LongHudouLogic) setSysPool(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		for _, table := range lh.GetTable().Room.Tables {
			logic := table.GameLogic.(*LongHudouLogic)
			logic.changeSysPool(pool)
		}
	} else {
		t := lh.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*LongHudouLogic)
			logic.changeSysPool(pool)
		} else {
			return lh.getRetMsg(false, "修改系统水池失败")
		}
	}
	lh.addChangeLog(tableid, pool, operator, 1)
	return lh.getRetMsg(true, "修改系统水池成功")
}

func (lh *LongHudouLogic) setFangFen(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		return lh.getRetMsg(false, "修改放分金額失败")
	} else {
		t := lh.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*LongHudouLogic)
			logic.fangFenTarget = pool
			logic.fangFenComplete = 0
		} else {
			return lh.getRetMsg(false, "修改放分金額失败")
		}
	}
	lh.addChangeLog(tableid, pool, operator, 2)
	return lh.getRetMsg(true, "修改放分金額成功")
}

func (lh *LongHudouLogic) addChangeLog(tableid int32, pool int64, operator string, typeid int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolAddChangeLogID)
	pdata := new(netproto.SuperToolAddChangeLog)
	pdata.GameID = proto.Int32(lh.GetServerConfig().GameID)
	pdata.TableNo = proto.Int32(tableid)
	spool := fmt.Sprintf("%.2f", float64(pool) / 100)
	if typeid == 1 {
		pdata.Result = proto.String(fmt.Sprintf("修改系統水池 %s", spool))
	} else {
		pdata.Result = proto.String(fmt.Sprintf("系統指定放分 %s", spool))
	}

	pdata.Type = proto.Int32(typeid)
	pdata.Operator = proto.String(operator)

	pmsg.MsgData = pdata
	lh.SendRPCMessageFunc(true, pmsg)
}
