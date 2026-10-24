package zjh

import (
	"example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"		
	"example.invalid/ServerTeam/Servers/pkg/util"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

import (
	"sort"
)

type ZjhVipUserPokerCard struct {
	CurrentCard string
	CurrentCardType string
}

type ZjhSuperToolHandler struct {
	Service gamelogic.ISuperToolService
	UserList map[int32]*netproto.SuperToolZjhVipUserData
	UserPokerCard *util.SafeMap
	//UserPokerCard map[int32]*ZjhVipUserPokerCard
	Proiority []int32
	HistoryCardType  []int
}

func NewZjhSuperToolHandler(service gamelogic.ISuperToolService) *ZjhSuperToolHandler {
	provider := new(ZjhSuperToolHandler)
	provider.Service = service
	provider.UserList = make(map[int32]*netproto.SuperToolZjhVipUserData)
	provider.UserPokerCard = util.NewSafeMap()
	return provider
}

func (provider *ZjhSuperToolHandler) GetVipUserOrg(userid int32) proto.Message {
	if ud, ok := provider.UserList[userid]; ok {
		return ud
	} else {
		return nil
	}
}
func (provider *ZjhSuperToolHandler) SetVipUserList(rmm proto.Message) {
	msg := rmm.(*netproto.SuperToolZjhVipUserDataList)
	provider.UserList = make(map[int32]*netproto.SuperToolZjhVipUserData)
	provider.Proiority = []int32{}
	for _, vipuser := range msg.DataList {
		provider.UserList[vipuser.GetUserID()] = vipuser
		provider.Proiority = append(provider.Proiority, vipuser.GetUserID())
	}
	sort.Slice(provider.Proiority, func(i, j int) bool {
		iuser := provider.GetVipUser(provider.Proiority[i])
		juser := provider.GetVipUser(provider.Proiority[j])
		if iuser.GetStatus() == 1 && iuser.GetFinish() == 0 {
			if iuser.GetNextWinLose() == 1  || juser.GetNextWinLose() == 0 {
				return true
			} else if juser.GetNextWinLose() == 1 || iuser.GetNextWinLose() == 0 {
				return false
			} else if iuser.GetNextWinLoseNum() > 0 || juser.GetNextWinLoseNum() < 0 {
				return true
			} else if juser.GetNextWinLoseNum() > 0 || iuser.GetNextWinLoseNum() < 0{
				return false
			}
		}
		r := util.RndInt(0, 2)
		return r == 1
	})
}

func (provider *ZjhSuperToolHandler) IsVipUserExists(userid int32) bool {
	if _, ok := provider.UserList[userid]; ok {
		return true
	} else {
		return false
	}
}

func (provider *ZjhSuperToolHandler) SaveVipUser(userid int32) {
	if vipuser := provider.GetVipUser(userid) ; vipuser != nil {
		//mlog.Debug("================================ %d", userid)
		serverid := provider.Service.GetRoomServerID()
		pmsg := new(msg.Message)
		pmsg.BClassID = int32(netproto.MessageBClassID_Game)
		pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolZjhUpdVipUserID)
		data := new(netproto.SuperToolZjhUpdVipUserData)
		data.UserID = proto.Int32(userid)
		data.Finish = proto.Int64(vipuser.GetFinish())
		data.FinishNextWinLoseNum = proto.Int64(vipuser.GetFinishNextWinLoseNum())
		data.NextTable = proto.Int32(vipuser.GetNextTable())
		data.NextType = proto.Int64(vipuser.GetNextType())
		pmsg.MsgData = data

		rpcinfo := netproto.CreateEmptyRPCInfo(serverid)
		//设置RPC消息桌号
		provider.Service.PushRPCMessage(rpcinfo, true, pmsg)
	}
}

func (provider *ZjhSuperToolHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_SuperToolZjhGetVipUserRetID)
}

func (provider *ZjhSuperToolHandler) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolZjhVipUserDataList{}
		return mm
	})
	retmsg := rmm.(*netproto.SuperToolZjhVipUserDataList)
	provider.SetVipUserList(retmsg)
	mlog.Debug("=====================update %d vip record", len(retmsg.DataList))
}


func (provider *ZjhSuperToolHandler) GetUserPokerCard(userid int32) *ZjhVipUserPokerCard {
	userPoker := provider.UserPokerCard.Get(userid)
	if userPoker != nil {
		return provider.UserPokerCard.Get(userid).(*ZjhVipUserPokerCard)
	} else {
		return nil
	}

	//if card, ok := provider.UserPokerCard[userid]; ok {
	//	return card
	//} else {
	//	return nil
	//}
}
func (provider *ZjhSuperToolHandler) AddUserPokerCard(userid int32, card string, cardtype string) {
	provider.UserPokerCard.Set(userid, &ZjhVipUserPokerCard{card, cardtype})
}
func (provider *ZjhSuperToolHandler) DelUserPokerCard(userid int32) {
	provider.UserPokerCard.Del(userid)
	//delete(provider.UserPokerCard, userid)
}

func (provider *ZjhSuperToolHandler) SaveCardType(cardtype int) {
	provider.HistoryCardType = append(provider.HistoryCardType, cardtype)
	size := len(provider.HistoryCardType)
	if size > 1000 {
		provider.HistoryCardType = provider.HistoryCardType[:size-1000]
	}

}

func (provider *ZjhSuperToolHandler) CountCardType() map[int]int{
	ret := make(map[int]int)
	for _, v := range provider.HistoryCardType {
		if _, ok := ret[v]; ok {
			ret[v]++
		} else {
			ret[v] = 1
		}
	}
	return ret
}

func (provider *ZjhSuperToolHandler) GetNextTable(userid int32) int32 {
	if vipuser := provider.GetVipUser(userid); vipuser != nil {
		if vipuser.GetStatus() == 1 {
			return vipuser.GetNextTable()
		}
	}
	return -1
}

func (provider *ZjhSuperToolHandler) ResetNextTable(userid int32) {
	if vipuser := provider.GetVipUser(userid); vipuser != nil {
		vipuser.NextTable = proto.Int32(-1)
	}
}

func (provider *ZjhSuperToolHandler) GetVipUser(userid int32) *netproto.SuperToolZjhVipUserData {
	ud := provider.GetVipUserOrg(userid)
	if ud != nil {
		return ud.(*netproto.SuperToolZjhVipUserData)
	}
	return nil
}

//獲取系統水池
func (provider *ZjhSuperToolHandler) GetSysPool() int64 { return -1 }

//獲取贏錢線
func (provider *ZjhSuperToolHandler) GetWinLine() int64 { return -1 }

//更新水池
func (provider *ZjhSuperToolHandler) UpdateSysPool(updatePool int64) { }

//啟動系統水池
func (provider *ZjhSuperToolHandler) SysPoolStart() {}

//
//type BlackOrWhiteUserProvider struct {
//	*VipUserProviderBase
//	ProiorityUserList []*netproto.SuperToolZjhVipUserData
//}
//
//func NewBlackOrWhiteUserProvider(vipUserProvider *VipUserProvider) *BlackOrWhiteUserProvider {
//	bowProvider := new(BlackOrWhiteUserProvider)
//	bowProvider.VipUserProviderBase = new(VipUserProviderBase)
//	bowProvider.UserList = make(map[int32]*netproto.SuperToolZjhVipUserData)
//	for _, vipuser := range vipUserProvider.UserList {
//		if vipuser.GetAssignAmount() != 0 && vipuser.GetComplete() == 0 && vipuser.GetStatus() == 1 {
//			bowProvider.UserList[vipuser.GetUserID()] = vipuser
//			bowProvider.ProiorityUserList = append(bowProvider.ProiorityUserList, vipuser)
//		}
//	}
//	//sort proiority
//	sort.Slice(bowProvider.ProiorityUserList, func(i, j int) bool {
//		if bowProvider.ProiorityUserList[i].GetAssignAmount() > 0 && bowProvider.ProiorityUserList[j].GetAssignAmount() > 0 { //白名單
//			if bowProvider.ProiorityUserList[i].GetAssignRate() == bowProvider.ProiorityUserList[j].GetAssignRate() {
//				return bowProvider.ProiorityUserList[i].GetAssignAmount() > bowProvider.ProiorityUserList[j].GetAssignAmount()
//			} else {
//				return bowProvider.ProiorityUserList[i].GetAssignRate() > bowProvider.ProiorityUserList[j].GetAssignRate()
//			}
//		} else if bowProvider.ProiorityUserList[i].GetAssignAmount() < 0 && bowProvider.ProiorityUserList[j].GetAssignAmount() < 0 { //黑名單
//			if bowProvider.ProiorityUserList[i].GetAssignRate() == bowProvider.ProiorityUserList[j].GetAssignRate() {
//				return bowProvider.ProiorityUserList[i].GetAssignAmount() > bowProvider.ProiorityUserList[j].GetAssignAmount()
//			} else {
//				return bowProvider.ProiorityUserList[i].GetAssignRate() < bowProvider.ProiorityUserList[j].GetAssignRate()
//			}
//		} else {
//			return bowProvider.ProiorityUserList[i].GetAssignAmount() > bowProvider.ProiorityUserList[j].GetAssignAmount()
//		}
//	})
//	return bowProvider
//}
//
//func (provider *BlackOrWhiteUserProvider) GetProiorityListString() string {
//	var s string
//	for _, vipuser := range provider.ProiorityUserList {
//		s += fmt.Sprintf("[%d|%d|%d],", vipuser.GetUserID(), vipuser.GetAssignAmount(), vipuser.GetAssignRate())
//	}
//	return s
//}