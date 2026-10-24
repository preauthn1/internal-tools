package hhnn

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

func (hhnn *HHNNLogic) updateVipUserWin() {
	bowBet := hhnn.getBowBet()
	superToolHandler := hhnn.SuperToolService.GetBrGameSuperToolHandler()
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
		hhnn.SuperToolService.RefreshVIPUser()
	}
}

func (hhnn *HHNNLogic) getAssignReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailID)
	pdata := new(netproto.SuperToolGetAssignRewardDetail)
	pdata.GameID = proto.Int32(hhnn.GetServerConfig().GameID)
	pdata.FlagID = proto.Int32(util.StrToInt32(hhnn.GetServerConfig().Flag))
	pdata.TableNo = proto.Int32(hhnn.TableNo())
	pmsg.MsgData = pdata

	hhnn.SendRPCMessageFunc(false, pmsg)
}

func (hhnn *HHNNLogic) setAssignReward(m *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolAssignRewardDetail{}
		return mm
	})
	detail := rmm.(*netproto.SuperToolAssignRewardDetail)
	hhnn.assignRewardDetail = detail
}

func (hhnn *HHNNLogic) isAssignReward() (bool, []int32) {
	var ret []int32
	if hhnn.assignRewardDetail != nil {
		ss := strings.Split(hhnn.assignRewardDetail.GetDescription(), ",")
		for _, s:= range ss {
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
			if len(ret) > 1 && ret[1] != 3{
				return false, []int32{}
			}
			return true, ret
		} else {
			return false, []int32{}
		}
	}
	return false, []int32{}
}

func (hhnn *HHNNLogic) assignRewardDetailResult(result string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdAssignRewardDetailID)
	pdata := hhnn.assignRewardDetail
	pdata.Reward = proto.String(result)
	if hhnn.execAssignReward {
		pdata.Status = proto.Int32(2)
	} else {
		pdata.Status = proto.Int32(3)
	}
	pmsg.MsgData = pdata

	hhnn.SendRPCMessageFunc(true, pmsg)
}

func (hhnn *HHNNLogic) saveLastGameUserBet() {
	hhnn.lastGameUserBet = hhnn.lastGameUserBet[:0] //clear last game user bet
	for _, ud := range hhnn.userdata {
		try := int64(0)
		for _, b := range ud.userbet {
			try += b.money
		}
		if try > 0 {
			hhnn.lastGameUserBet = append(hhnn.lastGameUserBet, ud)
			if userhis, ok := HHNNUserHistoryStore[ud.user.UserID];ok {
				HHNNUserHistoryStore[ud.user.UserID] = append([]*hhnnuserdata{ud}, userhis...)
			} else {
				HHNNUserHistoryStore[ud.user.UserID] = []*hhnnuserdata{ud}
			}
			if len(HHNNUserHistoryStore[ud.user.UserID]) > 100 {
				HHNNUserHistoryStore[ud.user.UserID] = HHNNUserHistoryStore[ud.user.UserID][:100]
			}
		}
	}
}

//獲取 bowbet
func (hhnn *HHNNLogic) getBowBet() (vipBet map[int32]*hhnnuserdata){
	vipBet = make(map[int32]*hhnnuserdata)
	for _, ud := range hhnn.userdata{
		if vipuser := hhnn.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(ud.user.UserID); vipuser != nil {
			if vipuser.GetStatus() == 1 && vipuser.GetComplete() == 0 && vipuser.GetAssignAmount() != 0 {
				vipBet[ud.user.UserID] = ud
			}
		}
	}
	return
}
//監控放殺分
func (hhnn *HHNNLogic) bowJiangType() bool{
	bowBet := hhnn.getBowBet()
	betlimit := hhnn.ctrlbetlimit
	//先執行放分
	for userid, ubet := range bowBet {
		vipuser := hhnn.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() < 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		example.invalid("用戶%d, 放分機率%d, %d", userid, vipuser.GetAssignRate(), rnd)
		if rnd < vipuser.GetAssignRate(){ //隨機0-99 小於勝率則執行放分
			//1. 獲取下注的項目
			betlong := ubet.getUserBetByType(HUNG)
			bethu := ubet.getUserBetByType(HEI)

			//2. 決定開獎項目
			if betlong >= betlimit {
				hhnn.jiangtype = HUNG
			} else if bethu >= betlimit {
				hhnn.jiangtype = HEI
			}
			if hhnn.jiangtype != 0 { // 獲取到開獎就跳出
				break
			}
		}
	}
	if hhnn.jiangtype != 0 {
		example.invalid("執行放分開獎 : %v", hhnn.jiangtype)
		return true
	}
	//殺分
	for userid, ubet := range bowBet {
		vipuser := hhnn.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() > 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		example.invalid("用戶%d, 殺分機率%d, %d", userid, vipuser.GetAssignRate(), rnd)
		if rnd >= vipuser.GetAssignRate(){ //隨機0-99 大於等於勝率則執行殺分
			//1. 獲取下注的項目
			betlong := ubet.getUserBetByType(HUNG)
			bethu := ubet.getUserBetByType(HEI)

			//2. 決定開獎項目
			if betlong >= betlimit {
				hhnn.jiangtype = HEI
			} else if bethu >= betlimit {
				hhnn.jiangtype = HUNG
			}
			if hhnn.jiangtype != 0 { // 獲取到開獎就跳出
				break
			}

		} else {
			continue
		}
	}
	if hhnn.jiangtype != 0 {
		example.invalid("執行殺分開獎 : %v", hhnn.jiangtype)
		return true
	}
	return false
}


// 獲取當局下注用戶列表
func (hhnn *HHNNLogic) responseCurrentUserBetList() (list []*netproto.Http_HHNN_UserInfo) {
	for _, ud := range hhnn.userdata {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_HHNN_UserInfo)
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
		userinfo.IsVIP = proto.Bool(hhnn.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
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
func (hhnn *HHNNLogic) responseLastGameUserBetList() (list []*netproto.Http_HHNN_UserInfo) {
	tmparr := hhnn.lastGameUserBet
	if len(hhnn.lastGameUserBet) > 50 {
		tmparr = hhnn.lastGameUserBet[:50]
	}

	for _, ud := range tmparr {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_HHNN_UserInfo)
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
		userinfo.IsVIP = proto.Bool(hhnn.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
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


func (hhnn *HHNNLogic) getBetTypeWinRate() (rate []int32){

	list := hhnn.resultcardlist.GetCardList()
	size := len(list)
	if size == 0 {
		return []int32{0,0,0}
	}
	wincount := []int{0,0,0}
	for _, p := range list {
		jiangtype := p.GetWin()
		zuhetype := int32(0)
		if HUNG == jiangtype {
			wincount[0]++
			zuhetype = p.GetHungPoker().GetNiuNiuNumber()
		} else {
			wincount[2]++
			zuhetype = p.GetHeiPoker().GetNiuNiuNumber()
		}
		if zuhetype != 0 && zuhetype >= WUXIAONIU && zuhetype <= KNIUNIU{
			wincount[1]++
		}
	}
	rate = append(rate, int32( (float64(wincount[0]) / float64(size)) * 100 ))
	rate = append(rate, int32( (float64(wincount[1]) / float64(size)) * 100 ))
	rate = append(rate, int32( (float64(wincount[2]) / float64(size)) * 100 ))
	return
}


func (hhnn *HHNNLogic) setSysPool(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		for _, table := range hhnn.GetTable().Room.Tables {
			logic := table.GameLogic.(*HHNNLogic)
			logic.changeSysPool(pool)
		}
	} else {
		t := hhnn.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*HHNNLogic)
			logic.changeSysPool(pool)
		} else {
			return hhnn.getRetMsg(false, "修改系统水池失败")
		}
	}
	hhnn.addChangeLog(tableid, pool, operator, 1)
	return hhnn.getRetMsg(true, "修改系统水池成功")
}

func (hhnn *HHNNLogic) setFangFen(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		return hhnn.getRetMsg(false, "修改放分金額失败")
	} else {
		t := hhnn.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*HHNNLogic)
			logic.fangFenTarget = pool
			logic.fangFenComplete = 0
		} else {
			return hhnn.getRetMsg(false, "修改放分金額失败")
		}
	}
	hhnn.addChangeLog(tableid, pool, operator, 2)
	return hhnn.getRetMsg(true, "修改放分金額成功")
}

func (hhnn *HHNNLogic) addChangeLog(tableid int32, pool int64, operator string, typeid int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolAddChangeLogID)
	pdata := new(netproto.SuperToolAddChangeLog)
	pdata.GameID = proto.Int32(hhnn.GetServerConfig().GameID)
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
	hhnn.SendRPCMessageFunc(true, pmsg)
}


// 是否致命一擊牌型
func (hhnn *HHNNLogic) isPowerJiangType(zuhetype int32) bool {
	if hhnn.zuhetype != 0 && hhnn.zuhetype > KNIUNIU {
		return false
	} else if hhnn.zuhetype == 0 {
		return false
	} else {
		return true
	}
}

// 獲取致命一擊牌型
func (hhnn *HHNNLogic) getPowerJiangType() int32 {
	pokertype := []int32{WUXIAONIU, ZHADAN, WUHUANIU, KNIUNIU}
	return hhnn.weight[1].getJiangType(pokertype)
}