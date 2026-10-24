package zjh

import (
	"encoding/base64"
	"io"
	"net/http"
	"sort"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

type ZjhSuperToolService struct {
	Room             *gamelogic.Room
	tfManager        *event.TimerFuncManager
	DBMsgQueue       *rpcmsg.Queue
	UpdateTicker     *event.TimerFunc
	superToolHandler gamelogic.ISuperToolHandler
	fapaiRule        int32 //發牌規則 1 大牌規則, 2 一般規則
}

func NewZjhSuperToolService(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) *ZjhSuperToolService {
	service := new(ZjhSuperToolService)
	service.Room = room
	service.tfManager = tfm
	service.DBMsgQueue = dbMsgQueue
	service.superToolHandler = NewZjhSuperToolHandler(service)
	service.fapaiRule = 2
	return service
}

func (service *ZjhSuperToolService) Start() {
	service.RefreshVIPUser()
	service.UpdateTicker = service.tfManager.NewTicker(10, service.RefreshVIPUser)
}

func (service *ZjhSuperToolService) Stop() {
	service.UpdateTicker.Stop()
}

func (service *ZjhSuperToolService) GetSuperToolHandler() gamelogic.ISuperToolHandler {
	return service.superToolHandler
}

func (service *ZjhSuperToolService) GetZjhSuperToolHandler() *ZjhSuperToolHandler {
	return service.superToolHandler.(*ZjhSuperToolHandler)
}

func (service *ZjhSuperToolService) RefreshVIPUser() {
	//m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_SuperToolGetVipUserID), rpcMessage)
	serverid := service.Room.ServerID
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolZjhGetVipUserID)
	pmsg.MsgData = nil

	rpcinfo := netproto.CreateEmptyRPCInfo(serverid)
	//设置RPC消息桌号
	service.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func (service *ZjhSuperToolService) GetRoomServerID() int32 {
	return service.Room.ServerID
}

func (service *ZjhSuperToolService) PushRPCMessage(rpcinfo *netproto.RPCInfo, safe bool, msg *msg.Message) error {
	return service.DBMsgQueue.PushMessage(rpcinfo, safe, msg)
}

func (service *ZjhSuperToolService) HTTPService(w http.ResponseWriter, req *http.Request) {
	//mlog.Debug("==================== httpservice")
	action := req.PostFormValue("action")
	tableno := util.StrToInt32(req.PostFormValue("tableno"))
	var ret proto.Message
	if action == "getserverinfo" {
		//mlog.Debug("================getserverinfo")
		ret = service.getServerInfo(tableno)
	} else if action == "getassigntablelist" {
		ret = service.getAssignTableList()
	} else if action == "setfapairule" {
		rule := util.StrToInt32(req.PostFormValue("rule"))
		service.setFapaiRule(rule)
	}
	if ret != nil {
		data, _ := proto.Marshal(ret)
		//mlog.Debug("================%s", data)
		if data != nil {
			encodeStr := base64.StdEncoding.EncodeToString(data)
			//mlog.Debug("+++++++++++++++++++ %s", encodeStr)
			io.WriteString(w, encodeStr)
		}
	} else {
		io.WriteString(w, "ok")
	}

}

//獲取遊戲ID
func (h *ZjhSuperToolService) GetGameID() int32 {
	return h.Room.GameInfo.GameID
}

//獲取房間等級
func (h *ZjhSuperToolService) GetFlag() string {
	return h.Room.Flag
}

func (service *ZjhSuperToolService) setFapaiRule(rule int32) {
	example.invalid("setFapaiRule %d", rule)
	if rule == 2 {
		service.fapaiRule = 2
	} else {
		service.fapaiRule = 1
	}
}

func (service *ZjhSuperToolService) getServerInfo(tableno int32) *netproto.Http_ZJH_ServerInfo {
	serverInfo := new(netproto.Http_ZJH_ServerInfo)
	for _, t := range service.Room.Tables {
		serverInfo.TableList = append(serverInfo.TableList, t.TableNo)
	}

	serverInfo.OnlineUserList = service.getOnlineUserList(tableno)

	cardTypeMap := service.getCardTypeCount()
	//mlog.Debug("==================%v", cardTypeMap)
	serverInfo.BaoZiCount = cardTypeMap[5]
	serverInfo.SunJinCount = cardTypeMap[4]
	serverInfo.JiaHuaCount = cardTypeMap[3]
	serverInfo.SunZiCount = cardTypeMap[2]
	serverInfo.DuiZiCount = cardTypeMap[1]
	serverInfo.SanPaiCount = cardTypeMap[0]

	serverInfo.FapaiRule = proto.Int32(service.fapaiRule)
	return serverInfo
}

func (service *ZjhSuperToolService) getAssignTableList() *netproto.Http_ZJH_AssignTableList {
	assignTableList := new(netproto.Http_ZJH_AssignTableList)
	start := service.Room.TableCount - 1
	for i := start; i >= service.Room.TableCount-20; i-- {
		assignTableList.TableList = append(assignTableList.TableList, i)
	}
	return assignTableList
}

func (service *ZjhSuperToolService) getOnlineUserList(tableno int32) []*netproto.Http_ZJH_UserInfo {
	var userlist []*netproto.Http_ZJH_UserInfo
	provider := service.GetZjhSuperToolHandler()
	for _, gameuser := range service.Room.RoomUserManager.GetAllUser() {
		if gameuser.TableNo != tableno && tableno != -999 {
			continue
		}
		pokercard := provider.GetUserPokerCard(gameuser.UserID)
		ud := provider.GetVipUser(gameuser.UserID)
		userinfo := new(netproto.Http_ZJH_UserInfo)
		userinfo.UserID = proto.Int32(gameuser.UserID)
		ut := int32(1)
		if gameuser.IsRobot() {
			ut = 0
		}
		userinfo.UserType = proto.Int32(ut)
		userinfo.NickName = proto.String(gameuser.NickName)
		userinfo.GameScore = proto.Int64(gameuser.UpdateScore)
		userinfo.InitMoney = proto.Int64(gameuser.LoginScore)
		userinfo.CurrentMoney = proto.Int64(gameuser.Score)
		userinfo.SumCharge = proto.Int64(gameuser.SumCharge)
		userinfo.SumConvert = proto.Int64(gameuser.SumConvert)
		userinfo.SumWinLose = proto.Int64((gameuser.SumConvert - gameuser.SumCharge) + gameuser.Score)
		if ud != nil {
			userinfo.TodayCharge = proto.Int64(ud.GetTodaySumcharge())
			userinfo.TodayConvert = proto.Int64(ud.GetTodaySumConvert())
			userinfo.TodayWinLose = proto.Int64(ud.GetTodayUserWinLose())
		} else {
			userinfo.TodayCharge = proto.Int64(0)
			userinfo.TodayConvert = proto.Int64(0)
			userinfo.TodayWinLose = proto.Int64(0)
		}
		userinfo.IsVIP = proto.Bool(provider.IsVipUserExists(gameuser.UserID))
		userinfo.LoginName = proto.String(gameuser.LoginName)
		userinfo.AlipayAccount = proto.String(gameuser.AlipayAccount)
		if pokercard != nil {
			userinfo.CurrentCard = proto.String(pokercard.CurrentCard)
			userinfo.CurrentCardType = proto.String(pokercard.CurrentCardType)
		} else {
			userinfo.CurrentCard = proto.String("")
			userinfo.CurrentCardType = proto.String("")
		}
		userinfo.TableNo = proto.Int64(int64(gameuser.TableNo))
		userlist = append(userlist, userinfo)
	}
	sort.Slice(userlist, func(i, j int) bool {
		if userlist[i].GetIsVIP() && userlist[j].GetIsVIP() {
			return userlist[i].GetSumWinLose() > userlist[j].GetSumWinLose()
		} else if userlist[i].GetIsVIP() {
			return true
		} else if userlist[j].GetIsVIP() {
			return false
		}
		if userlist[i].GetSumWinLose() == userlist[j].GetSumWinLose() {
			return userlist[i].GetUserID() < userlist[j].GetUserID()
		}
		return userlist[i].GetSumWinLose() > userlist[j].GetSumWinLose()
	})
	return userlist
}

func (service *ZjhSuperToolService) getCardTypeCount() map[int]*netproto.Http_ZJH_CardTypeInfo {
	cardTypeMap := make(map[int]*netproto.Http_ZJH_CardTypeInfo)
	for i := 0; i < 6; i++ {
		data := new(netproto.Http_ZJH_CardTypeInfo)
		data.Count = proto.Int32(0)
		data.Rate = proto.Int32(0)
		cardTypeMap[i] = data
	}
	ret := service.GetZjhSuperToolHandler().CountCardType()
	total := 0
	for _, c := range ret {
		total += c
	}
	//mlog.Debug("==================total %d", total)
	for i := 0; i < 6; i++ {
		count := int32(ret[i])
		rate := int32(0)
		if count > 0 {
			rate = int32((float32(ret[i]) / float32(total)) * 100)
		}

		cardTypeMap[i].Count = proto.Int32(count)
		cardTypeMap[i].Rate = proto.Int32(rate)
	}
	return cardTypeMap
}
