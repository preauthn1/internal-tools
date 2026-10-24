package msghandler

import (
	proto "example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// GameServerListHandler ?
type GameServerListHandler struct {
	serverProvider hallutil.GameServerProvider
	hallProvider   hallutil.HallVerProvider
}

// NewGameServerListHandler 創建遊戲服務器列表處理器對象
func NewGameServerListHandler(svrProvider hallutil.GameServerProvider, hallProvider hallutil.HallVerProvider) *GameServerListHandler {
	gshandler := new(GameServerListHandler)
	gshandler.serverProvider = svrProvider
	gshandler.hallProvider = hallProvider

	return gshandler
}

// IsHook ?
func (gh *GameServerListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_RequestServerListID)
}

// OnMessage ?
func (gh *GameServerListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	//mmbyte := m.MsgData.([]byte)
	//rmm := &netproto.GetShowGameFlagByAreaMsg{}
	//proto.Unmarshal(mmbyte, rmm)

	mret := new(msg.Message)
	mret.BClassID = int32(netproto.MessageBClassID_Hall)
	mret.SClassID = int32(netproto.HallMsgClassID_ServerListDataID)

	mdata := new(netproto.AllGameServerInfo)
	//svrlist := gh.serverProvider.GetGameServerList()
	areaID := clt.AreaID
	svrlist := gh.serverProvider.GetGameServerListByAreaID(areaID)
	allochSvrList := make([]*netproto.GameServerInfo, 0) // 取得分配服務器可用

	// 取得分配服務器可用
	for _, v := range svrlist {
		svraddr := gh.serverProvider.GetGameServerAddrByAreaIDAndFlag(areaID, v.GetGameID(), v.GetFlag())
		if svraddr != nil && len(svraddr) != 0 {
			for _, svr := range svraddr {
				exsit := false
				for _, alo := range allochSvrList {
					if svr.GetServerID() == alo.GetServerID() {
						exsit = true
						break
					}
				}
				if !exsit {
					svrdata := new(netproto.GameServerInfo)
					//svrdata = v
					svrdata.GameID = proto.Int32(v.GetGameID())
					svrdata.GameName = proto.String(v.GetGameName())
					svrdata.LoginMoney = proto.Int32(v.GetLoginMoney())
					svrdata.Flag = proto.String(v.GetFlag())
					svrdata.Status = proto.String(v.GetStatus())
					svrdata.BaseScore = proto.Int32(v.GetBaseScore())
					svrdata.Single = proto.Int32(v.GetSingle())
					svrdata.AndroidVersion = proto.String(v.GetAndroidVersion())
					svrdata.IOSVersion = proto.String(v.GetIOSVersion())
					svrdata.KickOutMoney = proto.Int32(v.GetKickOutMoney())
					svrdata.MaxOdd = proto.Int32(v.GetMaxOdd())
					svrdata.GameType = proto.Int32(v.GetGameType())

					svrdata.ServerID = proto.Int32(svr.GetServerID())
					svrdata.Port = proto.Int32(svr.GetPort())
					svrdata.IPAddress = proto.String(svr.GetIPAddress())
					allochSvrList = append(allochSvrList, svrdata)
				}
			}
		}
	}

	mdata.ServerList = allochSvrList
	mdata.GameCategory = hallutil.GetGameCategoryList(clt.AreaID)
	// fmt.Println("IsVisible: ", mdata.GameCategory.GetIsVisible())

	mret.MsgData = mdata

	return []*msg.Message{mret}
}
