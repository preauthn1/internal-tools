package brnn

import (
	"example.invalid/golang/protobuf/proto"
)

import (	
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/util"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

import (
	"regexp"
	"sort"
	"fmt"
	"math/rand"
)

func (brnn *BRNNLogic) updateVipUserWin() {
	bowBet := brnn.getBowBet()
	superToolHandler := brnn.SuperToolService.GetBrGameSuperToolHandler()
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
		brnn.SuperToolService.RefreshVIPUser()
	}
}

func (brnn *BRNNLogic) getAssignReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetAssignRewardDetailID)
	pdata := new(netproto.SuperToolGetAssignRewardDetail)
	pdata.GameID = proto.Int32(brnn.GetServerConfig().GameID)
	pdata.FlagID = proto.Int32(util.StrToInt32(brnn.GetServerConfig().Flag))
	pdata.TableNo = proto.Int32(brnn.TableNo())
	pmsg.MsgData = pdata

	brnn.SendRPCMessageFunc(false, pmsg)
}

func (brnn *BRNNLogic) setAssignReward(m *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolAssignRewardDetail{}
		return mm
	})
	detail := rmm.(*netproto.SuperToolAssignRewardDetail)
	brnn.assignRewardDetail = detail
}

func (brnn *BRNNLogic) isAssignReward() (bool, []int32) {
	if brnn.assignRewardDetail != nil {
		re := regexp.MustCompile("(\\d),(\\d),(\\d),(\\d)")
		matches := re.FindStringSubmatch(brnn.assignRewardDetail.GetDescription())
		if matches != nil {
			return true, []int32{util.StrToInt32(matches[1]), util.StrToInt32(matches[2]), util.StrToInt32(matches[3]), util.StrToInt32(matches[4])}
		}
	}
	return false, nil
}

func (brnn *BRNNLogic) assignRewardDetailResult(result string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdAssignRewardDetailID)
	pdata := brnn.assignRewardDetail
	pdata.Reward = proto.String(result)
	if brnn.execAssignReward {
		pdata.Status = proto.Int32(2)
	} else {
		pdata.Status = proto.Int32(3)
	}
	pmsg.MsgData = pdata

	brnn.SendRPCMessageFunc(true, pmsg)
}

func (brnn *BRNNLogic) saveLastGameUserBet() {
	brnn.lastGameUserBet = brnn.lastGameUserBet[:0] //clear last game user bet
	for _, ud := range brnn.userdata {
		try := int64(0)
		for _, b := range ud.userbet {
			try += b.money
		}
		if try > 0 {
			brnn.lastGameUserBet = append(brnn.lastGameUserBet, ud)
			if userhis, ok := BrnnUserHistoryStore[ud.user.UserID]; ok {
				BrnnUserHistoryStore[ud.user.UserID] = append([]*brnnuserdata{ud}, userhis...)
			} else {
				BrnnUserHistoryStore[ud.user.UserID] = []*brnnuserdata{ud}
			}
			if len(BrnnUserHistoryStore[ud.user.UserID]) > 100 {
				BrnnUserHistoryStore[ud.user.UserID] = BrnnUserHistoryStore[ud.user.UserID][:100]
			}
		}
	}
}

//獲取 bowbet
func (brnn *BRNNLogic) getBowBet() (vipBet map[int32]*brnnuserdata) {
	vipBet = make(map[int32]*brnnuserdata)
	for _, ud := range brnn.userdata {
		if vipuser := brnn.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(ud.user.UserID); vipuser != nil {
			if vipuser.GetStatus() == 1 && vipuser.GetComplete() == 0 && vipuser.GetAssignAmount() != 0 {
				vipBet[ud.user.UserID] = ud
			}
		}
	}
	return
}

//監控放殺分
func (brnn *BRNNLogic) bowJiangType() bool {
	bowBet := brnn.getBowBet()
	var jiangitem []int
	betlimit := brnn.ctrlbetlimit
	//先執行放分
	for userid, ubet := range bowBet {
		vipuser := brnn.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() < 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		example.invalid("用戶%d 執行放分機率 %d, %d", userid, vipuser.GetAssignRate(), rnd)
		if rnd < vipuser.GetAssignRate() { //隨機0-99 小於勝率則執行放分
			//1. 獲取下注的項目
			var betitem []int
			for i := int32(1); i < 5; i++ {
				bet := ubet.getUserBetByType(i)
				if bet >= betlimit {
					betitem = append(betitem, int(i))
				}
			}
			//2. 決定新獎項目
			var newitem []int
			for _, v := range betitem {
				if exists, _ := util.InArray(jiangitem, v); !exists {
					newitem = append(newitem, v)
				}
			}
			//3. 確認開獎項目, 4個以上要將新獎打亂去掉一個再合併獎項
			if (len(jiangitem) + len(newitem)) >= 4 {
				//打亂要加入的獎
				for i := range newitem {
					j := rand.Intn(i + 1)
					newitem[i], newitem[j] = newitem[j], newitem[i]
				}
				newitem = newitem[:len(newitem)-1] //打亂要加入的獎去掉最後一個
			}
			jiangitem = append(jiangitem, newitem...)

			if len(jiangitem) >= 3 { //jiangitem 收集三個以上就不跑了
				break
			}

		} else {
			continue
		}
	}
	if len(jiangitem) > 0 {
		brnn.jiangtype = assign
		brnn.nextjiangitem = util.ArrayToString(jiangitem, ",")
		example.invalid("執行放分開獎 : %v", brnn.nextjiangitem)
		return true
	}
	//殺分
	var notjiangitem []int
	for userid, ubet := range bowBet {
		vipuser := brnn.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(userid)
		if vipuser.GetAssignAmount() > 0 {
			continue
		}
		rnd := int32(util.RndInt(0, 100))
		if rnd >= vipuser.GetAssignRate() { //隨機0-99 大於等於勝率則執行殺分
			//1. 獲取下注的項目
			var betitem []int
			for i := int32(1); i < 5; i++ {
				bet := ubet.getUserBetByType(i)
				if bet >= betlimit {
					betitem = append(betitem, int(i))
				}
			}
			//2. 決定不開哪些獎
			var newitem []int
			for _, v := range betitem {
				if exists, _ := util.InArray(notjiangitem, v); !exists {
					newitem = append(newitem, v)
				}
			}
			//3. 確認不開哪些獎, 4個以上要將新項目打亂去掉一個再合併獎項
			if (len(notjiangitem) + len(newitem)) >= 4 {
				//打亂要加入的獎
				for i := range newitem {
					j := rand.Intn(i + 1)
					newitem[i], newitem[j] = newitem[j], newitem[i]
				}
				newitem = newitem[:len(newitem)-1] //打亂要加入的獎去掉最後一個
			}
			notjiangitem = append(notjiangitem, newitem...)

			if len(notjiangitem) >= 3 { //jiangitem 收集三個以上就不跑了
				break
			}

		} else {
			continue
		}
	}

	// 重新定義 jiangitem
	jiangitem = []int{1, 2, 3, 4}
	if len(notjiangitem) > 0 {
		for _, v := range notjiangitem {
			for idx, j := range jiangitem {
				if v == j {
					jiangitem = util.IntArrayRemoveIndex(jiangitem, idx)
					break
				}
			}
		}
		if len(jiangitem) > 1 { //如果要開的獎項大於一個, 則隨機開
			r := util.RndInt(0, len(jiangitem)) + 1
			//打亂要開的獎
			for i := range jiangitem {
				j := rand.Intn(i + 1)
				jiangitem[i], jiangitem[j] = jiangitem[j], jiangitem[i]
			}
			jiangitem = jiangitem[:r]
		}
		brnn.jiangtype = assign
		brnn.nextjiangitem = util.ArrayToString(jiangitem, ",")
		example.invalid("執行殺分開獎 : %v", brnn.nextjiangitem)
		return true
	}
	return false
}

// 獲取當局下注用戶列表
func (brnn *BRNNLogic) responseCurrentUserBetList() (list []*netproto.Http_BRNN_UserInfo) {
	for _, ud := range brnn.userdata {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_BRNN_UserInfo)
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
		userinfo.BetInfoChingLong = proto.Int64(ud.getUserBetByType(block1))
		userinfo.BetInfoBaiHu = proto.Int64(ud.getUserBetByType(block2))
		userinfo.BetInfoChuChueh = proto.Int64(ud.getUserBetByType(block3))
		userinfo.BetInfoHsuanWu = proto.Int64(ud.getUserBetByType(block4))
		userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
		userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
		userinfo.WinRate = proto.Int64(ud.calUserWinRate(9))
		userinfo.BetTotal = proto.Int64(ud.getUserBetTotal())
		userinfo.IsVIP = proto.Bool(brnn.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
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
func (brnn *BRNNLogic) responseLastGameUserBetList() (list []*netproto.Http_BRNN_UserInfo) {
	tmparr := brnn.lastGameUserBet
	if len(brnn.lastGameUserBet) > 50 {
		tmparr = brnn.lastGameUserBet[:50]
	}

	for _, ud := range tmparr {
		if ud.user.IsRobot() {
			continue
		}
		userinfo := new(netproto.Http_BRNN_UserInfo)
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
		userinfo.BetInfoChingLong = proto.Int64(ud.getUserBetByType(block1))
		userinfo.BetInfoBaiHu = proto.Int64(ud.getUserBetByType(block2))
		userinfo.BetInfoChuChueh = proto.Int64(ud.getUserBetByType(block3))
		userinfo.BetInfoHsuanWu = proto.Int64(ud.getUserBetByType(block4))
		userinfo.SumCharge = proto.Int64(ud.user.SumCharge)
		userinfo.SumConvert = proto.Int64(ud.user.SumConvert)
		userinfo.WinRate = proto.Int64(ud.calUserWinRate(9))
		userinfo.BetTotal = proto.Int64(ud.getUserBetTotal())
		userinfo.IsVIP = proto.Bool(brnn.SuperToolService.GetBrGameSuperToolHandler().IsVipUserExists(ud.user.UserID))
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

func (brnn *BRNNLogic) getBetTypeWinRate(size int) (rate []int32) {

	list := brnn.resultcardlist.GetPokerList()
	if len(list) > size {
		list = list[len(list)-9:]
	}
	size = len(list)
	if size == 0 {
		return []int32{0, 0, 0, 0}
	}
	wincount := []int{0, 0, 0, 0}
	for _, p := range list {
		for idx, r := range p.RetType {
			if r == 1 {
				wincount[idx]++
			}
		}
	}
	rate = append(rate, int32((float64(wincount[0])/float64(size))*100))
	rate = append(rate, int32((float64(wincount[1])/float64(size))*100))
	rate = append(rate, int32((float64(wincount[2])/float64(size))*100))
	rate = append(rate, int32((float64(wincount[3])/float64(size))*100))

	return
}

func (brnn *BRNNLogic) setSysPool(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		for _, table := range brnn.GetTable().Room.Tables {
			logic := table.GameLogic.(*BRNNLogic)
			logic.changeSysPool(pool)
		}
	} else {
		t := brnn.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*BRNNLogic)
			logic.changeSysPool(pool)
		} else {
			return brnn.getRetMsg(false, "修改系统水池失败")
		}
	}
	brnn.addChangeLog(tableid, pool, operator, 1)
	return brnn.getRetMsg(true, "修改系统水池成功")
}

func (brnn *BRNNLogic) setFangFen(tableid int32, pool int64, operator string) proto.Message {
	if tableid == -1 {
		return brnn.getRetMsg(false, "修改放分金額失败")
	} else {
		t := brnn.GetTable().Room.GetTableByNo(tableid)

		if t != nil {
			logic := t.GameLogic.(*BRNNLogic)
			logic.fangFenTarget = pool
			logic.fangFenComplete = 0
		} else {
			return brnn.getRetMsg(false, "修改放分金額失败")
		}
	}
	brnn.addChangeLog(tableid, pool, operator, 2)
	return brnn.getRetMsg(true, "修改放分金額成功")
}

func (brnn *BRNNLogic) addChangeLog(tableid int32, pool int64, operator string, typeid int32) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolAddChangeLogID)
	pdata := new(netproto.SuperToolAddChangeLog)
	pdata.GameID = proto.Int32(brnn.GetServerConfig().GameID)
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
	brnn.SendRPCMessageFunc(true, pmsg)
}