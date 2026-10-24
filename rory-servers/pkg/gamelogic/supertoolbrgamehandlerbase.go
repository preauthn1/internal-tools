package gamelogic

import (
	"fmt"
	"sort"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// SuperToolBrGameHandlerBase 百人超端基類
type SuperToolBrGameHandlerBase struct {
	Service         ISuperToolService
	UserList        map[int32]*netproto.SuperToolVipUserData
	Proiority       []int32
	HistoryCardType []int
	SysPool         int64 // 水池
	WinLine         int64 // 贏錢線
}

// NewBrGameSuperToolHandlerBase 百人超端基類
func NewBrGameSuperToolHandlerBase(service ISuperToolService) *SuperToolBrGameHandlerBase {
	provider := new(SuperToolBrGameHandlerBase)
	provider.Service = service
	provider.UserList = make(map[int32]*netproto.SuperToolVipUserData)
	return provider
}

// GetVipUserOrg ?
func (h *SuperToolBrGameHandlerBase) GetVipUserOrg(userid int32) proto.Message {
	if ud, ok := h.UserList[userid]; ok {
		return ud
	}
	return nil
}

// SetVipUserList ?
func (h *SuperToolBrGameHandlerBase) SetVipUserList(rmm proto.Message) {
	// 因為是事件驅動的關係, 在此不用考慮多執行緒問題
	msg := rmm.(*netproto.SuperToolVipUserDataList)
	// 更新vip 監控名單
	h.UserList = make(map[int32]*netproto.SuperToolVipUserData)
	h.Proiority = []int32{}
	for _, vipuser := range msg.DataList {
		userid := vipuser.GetUserID()
		h.UserList[userid] = vipuser
		h.Proiority = append(h.Proiority, userid)
	}

	sort.Slice(h.Proiority, func(i, j int) bool {
		iuser := h.GetVipUser(h.Proiority[i])
		juser := h.GetVipUser(h.Proiority[j])
		if iuser.GetAssignAmount() > 0 && juser.GetAssignAmount() > 0 { //白名單
			if iuser.GetAssignRate() == juser.GetAssignRate() {
				return iuser.GetAssignAmount() > juser.GetAssignAmount()
			}
			return iuser.GetAssignRate() > juser.GetAssignRate()
		} else if iuser.GetAssignAmount() < 0 && juser.GetAssignAmount() < 0 { //黑名單
			if iuser.GetAssignRate() == juser.GetAssignRate() {
				return iuser.GetAssignAmount() > juser.GetAssignAmount()
			}
			return iuser.GetAssignRate() < juser.GetAssignRate()
		}
		return iuser.GetAssignAmount() > juser.GetAssignAmount()
	})
}

// IsVipUserExists ?
func (h *SuperToolBrGameHandlerBase) IsVipUserExists(userid int32) bool {
	if _, ok := h.UserList[userid]; ok {
		return true
	}
	return false
}

// SaveVipUser ?
func (h *SuperToolBrGameHandlerBase) SaveVipUser(userid int32) {
	serverid := h.Service.GetRoomServerID()
	if vipuser := h.GetVipUser(userid); vipuser != nil {
		mlog.Debug("update user completeamount %d, userid %d", vipuser.GetCompleteAmount(), vipuser.GetUserID())
		pmsg := new(msg.Message)
		pmsg.BClassID = int32(netproto.MessageBClassID_Game)
		pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdVipUserID)
		pmsg.MsgData = vipuser
		rpcinfo := netproto.CreateEmptyRPCInfo(serverid)
		h.Service.PushRPCMessage(rpcinfo, true, pmsg)
	}
}

// IsHook ?
func (h *SuperToolBrGameHandlerBase) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && (sclsid == int32(netproto.GameMessageClassID_SuperToolGetVipUserRetID) ||
		sclsid == int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLoseRet))
}

// OnRPCMessage ?
func (h *SuperToolBrGameHandlerBase) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	switch netproto.GameMessageClassID(m.SClassID) {
	case netproto.GameMessageClassID_SuperToolGetVipUserRetID:
		h.getVipUserRet(m, ctx)
	case netproto.GameMessageClassID_SuperToolGetSysPoolWinLoseRet:
		h.setPoolWinLose(m, ctx)
	}
}

func (h *SuperToolBrGameHandlerBase) setPoolWinLose(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperTool_Msg_GetSysPoolWinLoseRet{}
		return mm
	})
	retmsg := rmm.(*netproto.SuperTool_Msg_GetSysPoolWinLoseRet)
	mlog.System("收到系統水池資料 pool %d, winline %d", retmsg.GetPoolAmount(), retmsg.GetWinLine())
	h.SysPool = retmsg.GetPoolAmount()
	h.WinLine = retmsg.GetWinLine()
}

func (h *SuperToolBrGameHandlerBase) getVipUserRet(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolVipUserDataList{}
		return mm
	})
	retmsg := rmm.(*netproto.SuperToolVipUserDataList)
	h.SetVipUserList(retmsg)
	mlog.Debug("=====================update %d vip record", len(retmsg.DataList))
}

// GetNextTable ?
func (h *SuperToolBrGameHandlerBase) GetNextTable(userid int32) int32 {
	return -1
}

// ResetNextTable ?
func (h *SuperToolBrGameHandlerBase) ResetNextTable(userid int32) {}

// GetVipUser ?
func (h *SuperToolBrGameHandlerBase) GetVipUser(userid int32) *netproto.SuperToolVipUserData {
	ud := h.GetVipUserOrg(userid)
	if ud != nil {
		return ud.(*netproto.SuperToolVipUserData)
	}
	return nil
}

// GetProiorityListString ?
func (h *SuperToolBrGameHandlerBase) GetProiorityListString() string {
	var s string
	for _, userid := range h.Proiority {
		if vipuser := h.GetVipUser(userid); vipuser != nil {
			s += fmt.Sprintf("[%d|%d|%d],", vipuser.GetUserID(), vipuser.GetAssignAmount(), vipuser.GetAssignRate())
		}
	}
	return s
}

// GetSysPool 獲取系統水池
func (h *SuperToolBrGameHandlerBase) GetSysPool() int64 { return h.SysPool }

// GetWinLine 獲取贏錢線
func (h *SuperToolBrGameHandlerBase) GetWinLine() int64 { return h.WinLine }

// UpdateSysPool 更新水池
func (h *SuperToolBrGameHandlerBase) UpdateSysPool(updatePool int64) {
	h.SysPool += updatePool
	serverid := h.Service.GetRoomServerID()
	gameid := h.Service.GetGameID()
	flagID, _ := strconv.Atoi(h.Service.GetFlag())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose)
	data := new(netproto.SuperTool_Msg_UpdateSysPoolWinLose)
	data.GameID = proto.Int32(gameid)
	data.FlagID = proto.Int32(int32(flagID))
	data.PoolAmount = proto.Int64(updatePool) //棄用
	data.UpdatePoolAmount = proto.Int64(updatePool)
	data.WinLine = proto.Int64(0)
	data.WinLoseData = proto.String("{}")
	pmsg.MsgData = data
	rpcinfo := netproto.CreateEmptyRPCInfo(serverid)
	_ = h.Service.PushRPCMessage(rpcinfo, true, pmsg)
}

// SysPoolStart ?
func (h *SuperToolBrGameHandlerBase) SysPoolStart() {
	mlog.System("請求系統水池資料")
	serverid := h.Service.GetRoomServerID()
	gameid := h.Service.GetGameID()
	flagID, _ := strconv.Atoi(h.Service.GetFlag())
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLose)
	data := new(netproto.SuperTool_Msg_GetSysPoolWinLose)
	data.GameID = proto.Int32(gameid)
	data.FlagID = proto.Int32(int32(flagID))
	pmsg.MsgData = data
	rpcinfo := netproto.CreateEmptyRPCInfo(serverid)
	_ = h.Service.PushRPCMessage(rpcinfo, true, pmsg)
}
