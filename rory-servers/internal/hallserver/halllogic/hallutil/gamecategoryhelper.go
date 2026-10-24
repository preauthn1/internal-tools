package hallutil

import (
	"time"

	proto "example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// ?
const (
	TimerNameServerReg = "hallserver_hallutil_gamecategory"
)

var (
	gameCategoryHelper *GameCategoryHelper
)

// GameCategoryHelper hallserver 定時更新遊戲分類清單
type GameCategoryHelper struct {
	rpcmsgqueue    *rpcmsg.Queue
	serverid       int32
	updateInterval int32
	data           *netproto.GameCategoryListDB
}

// NewGameCategoryHelper ?
func NewGameCategoryHelper(dbmsgqueue *rpcmsg.Queue, serverid int32, interval int32) *GameCategoryHelper {
	gameCategoryHelper = new(GameCategoryHelper)
	gameCategoryHelper.serverid = serverid
	gameCategoryHelper.updateInterval = interval
	gameCategoryHelper.rpcmsgqueue = dbmsgqueue

	return gameCategoryHelper
}

// GetGameCategoryList ?
func GetGameCategoryList(areaid int64) *netproto.GameCategoryList {

	if gameCategoryHelper != nil {
		var allOptions = gameCategoryHelper.data.Options
		// var options = make([]*netproto.GameCategoryOption, 0)
		var ret = new(netproto.GameCategoryList)
		for _, v := range gameCategoryHelper.data.GetGameTabs() {
			ret.GameTabs = append(ret.GameTabs, &netproto.CategoryInfo{CategoryID: proto.Int32(v.GetCategoryID()), Name: proto.String(v.GetName())})
		}
		for _, v := range gameCategoryHelper.data.GetCategoryList() {
			ret.CategoryList = append(ret.CategoryList, &netproto.GameCategoryIDs{GameID: proto.Int32(v.GetGameID()), CategoryIDs: proto.String(v.GetCategoryIDs())})
		}
		for _, v := range allOptions {
			if v.GetAreaID() == areaid {
				ret.IsVisible = proto.Int32(v.GetIsVisible())
				return ret
			}
		}
	}
	return nil
}

// SendDBGetGameCategoryMsg 遊戲分類清單
func (h *GameCategoryHelper) SendDBGetGameCategoryMsg() {

	svrmsg := new(msg.Message)
	svrmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	svrmsg.SClassID = int32(netproto.DBServerClassID_GetGameCategoryListID)

	svrdata := new(netproto.GetGameCategoryListMsg)
	svrmsg.MsgData = svrdata
	h.rpcmsgqueue.PushMessage(netproto.CreateEmptyRPCInfo(h.serverid), false, svrmsg)
}

// IsHook ?
func (h *GameCategoryHelper) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_GetGameCategoryListRetID)
}

// OnRPCMessage ?
func (h *GameCategoryHelper) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GameCategoryListDB{}
		h.data = mm

		// h.updateTicker = h.tfManager.NewNamedTicker(h.updateInterval, h.SendDBGetGameCategoryMsg, TimerName_ServerReg) 這個寫法會有 goroutine leak
		time.AfterFunc(time.Duration(h.updateInterval)*time.Second, h.SendDBGetGameCategoryMsg)
		return mm
	})
}
