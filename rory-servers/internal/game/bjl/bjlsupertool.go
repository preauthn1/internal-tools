package bjl

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

func (bjl *BaiJialeLogic) updateVipUserWin() {
	bowBet := bjl.getBowBet()
	superToolHandler := bjl.SuperToolService.GetBrGameSuperToolHandler()
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
		bjl.SuperToolService.RefreshVIPUser()
	}
}

func (bjl *BaiJialeLogic) getAssignReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailID)
	pdata := new(netproto.SuperToolGetAssignRewardDetail)
	pdata.GameID = proto.Int32(bjl.GetServerConfig().GameID)
	pdata.FlagID = proto.Int32(util.StrToInt32(bjl.GetServerConfig().Flag))
	pdata.TableNo = proto.Int32(bjl.TableNo())
	pmsg.MsgData = pdata

	bjl.SendRPCMessageFunc(false, pmsg)
}

func (bjl *BaiJialeLogic) setAssignReward(m *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolAssignRewardDetail{}
		return mm
	})
	detail := rmm.(*netproto.SuperToolAssignRewardDetail)
	bjl.assignRewardDetail = detail
}

func (bjl *BaiJialeLogic) isAssignReward() (bool, []int32) {
	var ret []int32
	if bjl.assignRewardDetail != nil {
		ss := strings.Split(bjl.assignRewardDetail.GetDescription(), ",")
		example.invalid("%s", ss)
		for _, s:= range ss {
			i := util.StrToInt32(s)
			if i != 0 {
				ret = append(ret, i)
			}
		}
		if len(ret) > 0 {
			//validate ret[0]
			if ret[0] != 1 && ret[0] != 2 && ret[0] != 3 {
				return false, []int32{}
			}
			//validate ret[1]
			if len(ret) > 1 && ret[1] != 6 && ret[1] != 7 {
				return false, []int32{}
			}
			return true, ret
		} else {
			return false, []int32{}
		}
	}
	return false, []int32{}
}

func (bjl *BaiJialeLogic) assignRewardDetailResult(result string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdAssignRewardDetailID)
	pdata := bjl.assignRewardDetail
	pdata.Reward = proto.String(result)
	if bjl.execAssignReward {
		pdata.Status = proto.Int32(2)
	} else {
		pdata.Status = proto.Int32(3)
	}
	pmsg.MsgData = pdata

	bjl.SendRPCMessageFunc(true, pmsg)
}

func (bjl *BaiJialeLogic) saveLastGameUserBet() {
	bjl.lastGameUserBet = bjl.lastGameUserBet[:0] //clear last game user bet
	for _, ud := range bjl.userdata {
		try := int64(0)
		for _, b := range ud.userbet {
			try += b.money
		}
		if try > 0 {
			bjl.lastGameUserBet = append(bjl.lastGameUserBet, ud)
			if userhis, ok := BJLUserHistoryStore[ud.user.UserID];ok {
				BJLUserHistoryStore[ud.user.UserID] = append([]*baijialeuserdata{ud}, userhis...)
			} else {
				BJLUserHistoryStore[ud.user.UserID] = []*baijialeuserdata{ud}
			}
			if len(BJLUserHistoryStore[ud.user.UserID]) > 100 {
				BJLUserHistoryStore[ud.user.UserID] = BJLUserHistoryStore[ud.user.UserID][:100]
			}
		}
	}
}

//獲取 bowbet
func (bjl *BaiJialeLogic) getBowBet() (vipBet map[int32]*baijialeuserdata){
	vipBet = make(map[int32]*baijialeuserdata)
	for _, ud := range bjl.userdata{
		if vipuser := bjl.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(ud.user.UserID); vipuser != nil {
			if vipuser.GetStatus() == 1 && vipuser.GetComplete() == 0 && vipuser.GetAssignAmount() != 0 {
				vipBet[ud.user.UserID] = ud
			}
		}
	}
	return
}
//監控放殺分
func (bjl *BaiJialeLogic) bowJiangType() bool{
	bowBet := bjl.getBowBet()
	betlimit := bjl.ctrlbetlimit
	//先執行放分
	for userid, ubet := range bowBet {
		vipuser := bjl.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() < 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		if rnd < vipuser.GetAssignRate(){ //隨機0-99 小於勝率則執行放分
			//1. 獲取下注的項目
			betlong := ubet.getUserBetByType(ZHUANG)
			bethu := ubet.getUserBetByType(XIAN)
			bethe := ubet.getUserBetByType(HE)

			//2. 決定開獎項目
			if betlong >= betlimit {
				bjl.jiangtype = []int32{ZHUANG}
				bjl.zuhetype |= 1 << uint32(ZHUANG-1) //莊
			} else if bethu >= betlimit {
				bjl.jiangtype = []int32{XIAN}
				bjl.zuhetype |= 1 << uint32(XIAN-1) //閒
			}
			//3. 放分開和條件及機率
			if bethe >= betlimit {
				r := (vipuser.GetAssignAmount() * 15 / 10) - vipuser.GetCompleteAmount()
				hepay := ubet.getUserBetByType(HE) * int64(bjl.betvalue.getOddByType(HE))
				if r > hepay {
					rnd := util.RndInt(0, 1000)
					if rnd < 75 {
						bjl.jiangtype = []int32{HE}
						bjl.zuhetype |= 1 << uint32(HE-1) //和
					}
				}
			}

			if bjl.jiangtype != nil { // 獲取到開獎就跳出
				break
			}
		}
	}
	if bjl.jiangtype != nil {
		example.invalid("執行放分開獎 : %v", bjl.jiangtype)
		return true
	}
	//殺分
	for userid, ubet := range bowBet {
		vipuser := bjl.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() > 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		if rnd >= vipuser.GetAssignRate(){ //隨機0-99 大於等於勝率則執行殺分
			//1. 獲取下注的項目
			betlong := ubet.getUserBetByType(ZHUANG)
			bethu := ubet.getUserBetByType(XIAN)

			//2. 決定開獎項目
			if betlong >= betlimit {
				bjl.jiangtype = []int32{XIAN}
				bjl.zuhetype |= 1 << uint32(XIAN-1) //閒
			} else if bethu >= betlimit {
				bjl.jiangtype = []int32{ZHUANG}
				bjl.zuhetype |= 1 << uint32(ZHUANG-1) //莊
			}
			if bjl.jiangtype != nil { // 獲取到開獎就跳出
				break
			}

		} else {
			continue
		}
	}
	if bjl.jiangtype != nil {
		example.invalid("執行殺分開獎 : %v", bjl.jiangtype)
		return true
	}
	return false
}


// 獲取當局下注用戶列表
func (bjl *BaiJialeLogic) responseCurrentUserBetList() (list []*netproto.Http_BJL_UserInfo) {
	for _, ud := range bjl.userdata {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_BJL_UserInfo)
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
		userinfo.BetInfoZhuang = proto.Int64(ud.getUserBetByType(ZHUANG))
		userinfo.BetInfoXian = proto.Int64(ud.getUserBetByType(XIAN))
		userinfo.BetInfoHe = proto.Int64(ud.getUserBetByType(HE))
		userinfo.BetInfoZhuangDZ = proto.Int64(ud.getUserBetByType(ZHUANGDUI))
		userinfo.BetInfoXianDZ = proto.Int64(ud.getUserBetByType(XIANDUI))
		userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
		userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
		userinfo.WinRate = proto.Int64(ud.calUserWinRate(100))
		userinfo.BetTotal = proto.Int64(ud.getUserBetTotal())
		userinfo.IsVIP = proto.Bool(bjl.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
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
func (bjl *BaiJialeLogic) responseLastGameUserBetList() (list []*netproto.Http_BJL_UserInfo) {
	tmparr := bjl.lastGameUserBet
	if len(bjl.lastGameUserBet) > 50 {
		tmparr = bjl.lastGameUserBet[:50]
	}

	for _, ud := range tmparr {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_BJL_UserInfo)
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
		userinfo.BetInfoZhuang = proto.Int64(ud.getUserBetByType(ZHUANG))
		userinfo.BetInfoXian = proto.Int64(ud.getUserBetByType(XIAN))
		userinfo.BetInfoHe = proto.Int64(ud.getUserBetByType(HE))
		userinfo.BetInfoZhuangDZ = proto.Int64(ud.getUserBetByType(ZHUANGDUI))
		userinfo.BetInfoXianDZ = proto.Int64(ud.getUserBetByType(XIANDUI))
		userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
		userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
		userinfo.WinRate = proto.Int64(ud.calUserWinRate(100))
		userinfo.BetTotal = proto.Int64(ud.getUserBetTotal())
		userinfo.IsVIP = proto.Bool(bjl.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
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


func (bjl *BaiJialeLogic) getBetTypeWinRate() (rate []int32){

	list := bjl.resultcardlist.GetCardList()
	size := len(list)
	if size == 0 {
		return []int32{0,0,0,0,0}
	}
	wincount := []int{0,0,0,0,0}
	for _, p := range list {
		zuhe := p.GetZuHe()
		length := fmt.Sprintf("%b", zuhe)
		if length == "1" || length == "100001" || length == "1000001" || length == "1100001"{
			//莊；莊，莊對；莊，閒對
			wincount[0]++
		}
		if length == "10" || length == "100010" || length == "1000010" || length == "1100010"{
			//閒；閒，莊對；閒，閒對
			wincount[1]++
		}
		if length == "100" || length == "100100" || length == "1000100" || length == "1100100"{
			//和；和，莊對；和，閒對
			wincount[2]++
		}
		if length == "100001" || length == "100010" || length == "100100" || length == "1100001" || length == "1100010"|| length == "1100100"{
			//莊，莊對；閒，莊對；和，莊對
			wincount[3]++
		}
		if length == "1000001" || length == "1000010" || length == "1000100" || length == "1100001" || length == "1100010"|| length == "1100100"{
			//莊，閒對；閒，閒對；和，閒對
			wincount[4]++
		}
	}

	rate = append(rate, int32( (float64(wincount[0]) / float64(size)) * 100 ))
	rate = append(rate, int32( (float64(wincount[1]) / float64(size)) * 100 ))
	rate = append(rate, int32( (float64(wincount[2]) / float64(size)) * 100 ))
	rate = append(rate, int32( (float64(wincount[3]) / float64(size)) * 100 ))
	rate = append(rate, int32( (float64(wincount[4]) / float64(size)) * 100 ))

	return
}


func (bjl *BaiJialeLogic) setSysPool(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		for _, table := range bjl.GetTable().Room.Tables {
			logic := table.GameLogic.(*BaiJialeLogic)
			logic.changeSysPool(pool)
		}
	} else {
		t := bjl.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*BaiJialeLogic)
			logic.changeSysPool(pool)
		} else {
			return bjl.getRetMsg(false, "修改系统水池失败")
		}
	}
	bjl.addChangeLog(tableid, pool, operator, 1)
	return bjl.getRetMsg(true, "修改系统水池成功")
}

func (bjl *BaiJialeLogic) setFangFen(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		return bjl.getRetMsg(false, "修改放分金額失败")
	} else {
		t := bjl.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*BaiJialeLogic)
			logic.fangFenTarget = pool
			logic.fangFenComplete = 0
		} else {
			return bjl.getRetMsg(false, "修改放分金額失败")
		}
	}
	bjl.addChangeLog(tableid, pool, operator, 2)
	return bjl.getRetMsg(true, "修改放分金額成功")
}

func (bjl *BaiJialeLogic) addChangeLog(tableid int32, pool int64, operator string, typeid int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolAddChangeLogID)
	pdata := new(netproto.SuperToolAddChangeLog)
	pdata.GameID = proto.Int32(bjl.GetServerConfig().GameID)
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
	bjl.SendRPCMessageFunc(true, pmsg)
}