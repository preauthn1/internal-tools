package hhdz

import (
	"fmt"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
	"sort"
	"strings"
)

//更新監控放殺分數據
func (hhdz *HHDZLogic) updateVipUserWin() {
	bowBet := hhdz.getBowBet()
	superToolHandler := hhdz.SuperToolService.GetBrGameSuperToolHandler()
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
		hhdz.SuperToolService.RefreshVIPUser()
	}
}

//請求指定多局開獎
func (hhdz *HHDZLogic) getAssignReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailID)
	pdata := new(netproto.SuperToolGetAssignRewardDetail)
	pdata.GameID = proto.Int32(hhdz.GetServerConfig().GameID)
	pdata.FlagID = proto.Int32(util.StrToInt32(hhdz.GetServerConfig().Flag))
	pdata.TableNo = proto.Int32(hhdz.TableNo())
	pmsg.MsgData = pdata
	hhdz.SendRPCMessageFunc(false, pmsg)
}

//返回指定多局開獎
func (hhdz *HHDZLogic) setAssignReward(m *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolAssignRewardDetail{}
		return mm
	})
	detail := rmm.(*netproto.SuperToolAssignRewardDetail)
	hhdz.assignRewardDetail = detail
}

//是否執行指定多局開獎
func (hhdz *HHDZLogic) isAssignReward() (bool, []int32) {
	var ret []int32
	if hhdz.assignRewardDetail != nil {
		ss := strings.Split(hhdz.assignRewardDetail.GetDescription(), ",")
		for _, s := range ss {
			i := util.StrToInt32(s)
			if i != 0 {
				ret = append(ret, i)
			}
		}
		if len(ret) > 0 {
			//validate ret[0]
			if ret[0] != 1 && ret[0] != 2 {
				return false, []int32{}
			}
			//validate ret[1]
			if len(ret) > 1 && ret[1] != 3 {
				return false, []int32{}
			}
			return true, ret
		} else {
			return false, []int32{}
		}
	}
	return false, []int32{}
}

// 指定多局開獎結果回報
func (hhdz *HHDZLogic) assignRewardDetailResult(result string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdAssignRewardDetailID)
	pdata := hhdz.assignRewardDetail
	pdata.Reward = proto.String(result)
	if hhdz.execAssignReward {
		pdata.Status = proto.Int32(2)
	} else {
		pdata.Status = proto.Int32(3)
	}
	pmsg.MsgData = pdata

	hhdz.SendRPCMessageFunc(true, pmsg)
}

// 上局投注用戶
func (hhdz *HHDZLogic) saveLastGameUserBet() {
	hhdz.lastGameUserBet = hhdz.lastGameUserBet[:0] //clear last game user bet
	for _, ud := range hhdz.userdata {
		try := int64(0)
		for _, b := range ud.userbet {
			try += b.money
		}
		if try > 0 {
			hhdz.lastGameUserBet = append(hhdz.lastGameUserBet, ud)
			if userhis, ok := HhdzUserHistoryStore[ud.user.UserID]; ok {
				HhdzUserHistoryStore[ud.user.UserID] = append([]*hhdzuuserdata{ud}, userhis...)
			} else {
				HhdzUserHistoryStore[ud.user.UserID] = []*hhdzuuserdata{ud}
			}
			if len(HhdzUserHistoryStore[ud.user.UserID]) > 100 {
				HhdzUserHistoryStore[ud.user.UserID] = HhdzUserHistoryStore[ud.user.UserID][:100]
			}
		}
	}
}

//獲取監控殺放分名單
func (hhdz *HHDZLogic) getBowBet() (vipBet map[int32]*hhdzuuserdata) {
	vipBet = make(map[int32]*hhdzuuserdata)
	for _, ud := range hhdz.userdata {
		if vipuser := hhdz.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(ud.user.UserID); vipuser != nil {
			if vipuser.GetStatus() == 1 && vipuser.GetComplete() == 0 && vipuser.GetAssignAmount() != 0 {
				vipBet[ud.user.UserID] = ud
			}
		}
	}
	return
}

//監控放殺分
func (hhdz *HHDZLogic) bowJiangType() bool {
	bowBet := hhdz.getBowBet()
	betlimit := hhdz.ctrlbetlimit
	//先執行放分
	for userid, ubet := range bowBet {
		vipuser := hhdz.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() < 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		if rnd < vipuser.GetAssignRate() { //隨機0-99 小於勝率則執行放分
			//1. 獲取下注的項目
			bethung := ubet.getUserBetByType(HUNG)
			bethei := ubet.getUserBetByType(HEI)
			betpower := ubet.getUserBetByType(3)
			//2. 決定開獎項目
			if bethung >= betlimit {
				hhdz.jiangtype = HUNG
			} else if bethei >= betlimit {
				hhdz.jiangtype = HEI
			}
			//3. 決定組合牌型
			rnd := util.RndInt(0, 1000)
			var zuhetype []int32
			if betpower >= betlimit && rnd < 165 {
				r := (vipuser.GetAssignAmount() * 15 / 10) - vipuser.GetCompleteAmount()
				alltype := []int32{BAOZI, SHUNJIN, JINHUA, SHUNZI, BIGDUIZI}
				for _, zuhe := range alltype {
					zuhepay := betpower * int64(hhdz.betvalue.getOddByType(zuhe))
					if r > zuhepay {
						zuhetype = append(zuhetype, zuhe)
					}
				}
				if len(zuhetype) < 1 {
					zuhetype = []int32{SMALLDUIZI, DANPAI}
				}
				hhdz.zuhetype = hhdz.weight[1].getJiangType(zuhetype)
			}

			if hhdz.jiangtype != 0 { // 獲取到開獎就跳出
				break
			}
		}
	}
	if hhdz.jiangtype != 0 {
		//決定牌型
		pokertype := []int32{BAOZI, SHUNJIN, JINHUA, SHUNZI, BIGDUIZI, SMALLDUIZI, DANPAI}
		hhdz.zuhetype = hhdz.weight[1].getJiangType(pokertype)
		example.invalid("執行監控放分開獎 : %v, 牌型 %v", hhdz.jiangtype, hhdz.zuhetype)
		return true
	}
	//殺分
	for userid, ubet := range bowBet {
		vipuser := hhdz.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() > 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		if rnd >= vipuser.GetAssignRate() { //隨機0-99 大於等於勝率則執行殺分
			//1. 獲取下注的項目
			bethung := ubet.getUserBetByType(HUNG)
			bethei := ubet.getUserBetByType(HEI)

			//2. 決定開獎項目
			if bethung >= betlimit {
				hhdz.jiangtype = HEI
			} else if bethei >= betlimit {
				hhdz.jiangtype = HUNG
			}

			if hhdz.jiangtype != 0 { // 獲取到開獎就跳出
				break
			}
		} else {
			continue
		}
	}
	if hhdz.jiangtype != 0 {
		//決定牌型
		pokertype := []int32{BAOZI, SHUNJIN, JINHUA, SHUNZI, BIGDUIZI, SMALLDUIZI, DANPAI}
		hhdz.zuhetype = hhdz.weight[1].getJiangType(pokertype)
		example.invalid("執行監控殺分開獎 : %v, 牌型 %v", hhdz.jiangtype, hhdz.zuhetype)
		return true
	}
	return false
}

// 獲取當局下注用戶列表
func (hhdz *HHDZLogic) responseCurrentUserBetList() (list []*netproto.Http_HHDZ_UserInfo) {
	for _, ud := range hhdz.userdata {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_HHDZ_UserInfo)
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
		userinfo.BetInfoHung = proto.Int64(ud.getUserBetByType(HUNG))
		userinfo.BetInfoHei = proto.Int64(ud.getUserBetByType(HEI))
		userinfo.BetInfoPower = proto.Int64(ud.getUserBetByType(3))
		userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
		userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
		userinfo.WinRate = proto.Int64(ud.calUserWinRate(100))
		userinfo.BetTotal = proto.Int64(ud.getUserBetTotal())
		userinfo.IsVIP = proto.Bool(hhdz.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
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
func (hhdz *HHDZLogic) responseLastGameUserBetList() (list []*netproto.Http_HHDZ_UserInfo) {
	tmparr := hhdz.lastGameUserBet
	if len(hhdz.lastGameUserBet) > 50 {
		tmparr = hhdz.lastGameUserBet[:50]
	}

	for _, ud := range tmparr {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_HHDZ_UserInfo)
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
		userinfo.BetInfoHung = proto.Int64(ud.getUserBetByType(HUNG))
		userinfo.BetInfoHei = proto.Int64(ud.getUserBetByType(HEI))
		userinfo.BetInfoPower = proto.Int64(ud.getUserBetByType(3))
		userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
		userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
		userinfo.WinRate = proto.Int64(ud.calUserWinRate(100))
		userinfo.BetTotal = proto.Int64(ud.getUserBetTotal())
		userinfo.IsVIP = proto.Bool(hhdz.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
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

// 獲取投注區勝率
func (hhdz *HHDZLogic) getBetTypeWinRate() (rate []int32) {

	list := hhdz.resultcardlist.GetCardList()
	size := len(list)
	if size == 0 {
		return []int32{0, 0, 0}
	}
	wincount := []int{0, 0, 0}
	for _, p := range list {
		jiangtype := p.GetWin()
		zuhetype := int32(0)
		if HUNG == jiangtype {
			wincount[0]++
			zuhetype = p.GetHungZuHe()
		} else {
			wincount[2]++
			zuhetype = p.GetHeiZuHe()
		}
		if zuhetype != 0 && zuhetype != DANPAI && zuhetype != SMALLDUIZI {
			wincount[1]++
		}

	}
	rate = append(rate, int32((float64(wincount[0])/float64(size))*100))
	rate = append(rate, int32((float64(wincount[1])/float64(size))*100))
	rate = append(rate, int32((float64(wincount[2])/float64(size))*100))
	return
}

// 配置系統水池
func (hhdz *HHDZLogic) setSysPool(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		for _, table := range hhdz.GetTable().Room.Tables {
			logic := table.GameLogic.(*HHDZLogic)
			logic.changeSysPool(pool)
		}
	} else {
		t := hhdz.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*HHDZLogic)
			logic.changeSysPool(pool)
		} else {
			return hhdz.getRetMsg(false, "修改系统水池失败")
		}
	}
	hhdz.addChangeLog(tableid, pool, operator, 1)
	return hhdz.getRetMsg(true, "修改系统水池成功")
}

// 配置系統放分
func (hhdz *HHDZLogic) setFangFen(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		return hhdz.getRetMsg(false, "修改放分金額失败")
	} else {
		t := hhdz.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*HHDZLogic)
			logic.fangFenTarget = pool
			logic.fangFenComplete = 0
		} else {
			return hhdz.getRetMsg(false, "修改放分金額失败")
		}
	}
	hhdz.addChangeLog(tableid, pool, operator, 2)
	return hhdz.getRetMsg(true, "修改放分金額成功")
}

// 配置系統水池/系統放分記錄
func (hhdz *HHDZLogic) addChangeLog(tableid int32, pool int64, operator string, typeid int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolAddChangeLogID)
	pdata := new(netproto.SuperToolAddChangeLog)
	pdata.GameID = proto.Int32(hhdz.GetServerConfig().GameID)
	pdata.TableNo = proto.Int32(tableid)
	spool := fmt.Sprintf("%.2f", float64(pool)/100)
	if typeid == 1 {
		pdata.Result = proto.String(fmt.Sprintf("修改系統水池 %s", spool))
	} else {
		pdata.Result = proto.String(fmt.Sprintf("系統指定放分 %s", spool))
	}

	pdata.Type = proto.Int32(typeid)
	pdata.Operator = proto.String(operator)

	pmsg.MsgData = pdata
	hhdz.SendRPCMessageFunc(true, pmsg)
}

// 是否致命一擊牌型
func (hhdz *HHDZLogic) isPowerJiangType(zuhetype int32) bool {
	if hhdz.zuhetype != 0 && hhdz.zuhetype != DANPAI && hhdz.zuhetype != SMALLDUIZI {
		return true
	} else {
		return false
	}
}

// 獲取致命一擊牌型
func (hhdz *HHDZLogic) getPowerJiangType() int32 {
	pokertype := []int32{BAOZI, SHUNJIN, JINHUA, SHUNZI, BIGDUIZI}
	return hhdz.weight[1].getJiangType(pokertype)
}
