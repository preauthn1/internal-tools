package rpchandler

import (
	"fmt"
	"strings"
	"sync"

	proto "example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// HallServerRegHandler 大廳啟動
type HallServerRegHandler struct {
	clientAddr            string
	notifyAddr            string
	httpAddr              string
	serverid              int32
	rpcmsgqueue           *rpcmsg.Queue
	hallServerRegRet      *netproto.HallServerStartRet
	AreaConfigMap         map[int64]*netproto.AreaConfig // 渠道配置
	AreaHallVersionMap    *sync.Map                      // 渠道版本號
	AreaPlatformMap       *sync.Map                      // 渠道平台表
	PlatformGameServerMap *sync.Map                      // 遊戲服務器平台表
	NotifyHandler         *serverhandler.NotifyHandler
	ChanAreaConfigStatus  chan bool
}

// NewHallServerRegHandler 創建大廳服務器註冊處理器對象
func NewHallServerRegHandler(clientaddr string, notifyaddr string, httpaddr string, serverid int32, dbmsgqueue *rpcmsg.Queue, nshandler *serverhandler.NotifyHandler) *HallServerRegHandler {
	hreg := new(HallServerRegHandler)
	hreg.clientAddr = clientaddr
	hreg.notifyAddr = notifyaddr
	hreg.httpAddr = httpaddr
	hreg.serverid = serverid
	hreg.rpcmsgqueue = dbmsgqueue
	hreg.hallServerRegRet = new(netproto.HallServerStartRet)
	hreg.AreaConfigMap = make(map[int64]*netproto.AreaConfig, 0)
	hreg.AreaHallVersionMap = new(sync.Map)
	hreg.AreaPlatformMap = new(sync.Map)
	hreg.PlatformGameServerMap = new(sync.Map)
	hreg.ChanAreaConfigStatus = make(chan bool, 1)
	hreg.NotifyHandler = nshandler

	return hreg
}

// IsHook ?
func (hsreg *HallServerRegHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) &&
		(sclsid == int32(netproto.DBServerClassID_HallServerStartRetID) ||
			sclsid == int32(netproto.DBServerClassID_GetAreaConfigRetID))
}

// OnRPCMessage ?
func (hsreg *HallServerRegHandler) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {

	if m.SClassID == int32(netproto.DBServerClassID_HallServerStartRetID) {
		_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
			mm := &netproto.HallServerStartRet{}
			return mm
		})

		mm := rmm.(*netproto.HallServerStartRet)
		// 重置
		hsreg.AreaConfigMap = make(map[int64]*netproto.AreaConfig, 0)
		hsreg.AreaHallVersionMap = new(sync.Map)
		hsreg.AreaPlatformMap = new(sync.Map)
		hsreg.PlatformGameServerMap = new(sync.Map)

		for _, v := range mm.GetAreaConfigList() {
			hsreg.AreaConfigMap[v.GetAreaID()] = v
		}

		for _, v := range mm.GetAreaHallVersionList() {
			hsreg.AreaHallVersionMap.Store(v.GetAreaID(), v)
		}

		for _, v := range mm.GetAreaPlatformList() {
			hsreg.AreaPlatformMap.Store(v.GetAreaID(), v.GetPlatformID())
		}
		tmpPlatformGameServerList := make(map[int64][]*netproto.PlatformGameServer)
		for _, v := range mm.GetPlatformGameServerList() {
			if _, ok := tmpPlatformGameServerList[v.GetPlatformID()]; !ok {
				tmpPlatformGameServerList[v.GetPlatformID()] = []*netproto.PlatformGameServer{}
			}
			tmpPlatformGameServerList[v.GetPlatformID()] = append(tmpPlatformGameServerList[v.GetPlatformID()], v)
		}
		for platformID, v := range tmpPlatformGameServerList {
			hsreg.PlatformGameServerMap.Store(platformID, v)
		}

		hsreg.hallServerRegRet = mm

		example.invalid("收到dbserver返回註冊信息.code:%d,msg:%s,android ver:%s,ios ver:%s,", mm.GetCode(), mm.GetMessage(), mm.GetAndroidVersion(), mm.GetIOSVersion())

		if mm.GetCode() == 0 {
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("[HallServerRegHandler] 向DBServer註冊失敗.%s 請儘速洽詢資訊人員", mm.GetMessage())).Error("HallServer壓測時")
			ctx.Error("向DBServer註冊失敗.%s", mm.GetMessage())
		} else {
			example.invalid("向DBServer註冊成功")
		}
	}

	if m.SClassID == int32(netproto.DBServerClassID_GetAreaConfigRetID) {
		_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
			mm := &netproto.AreaConfig{}
			return mm
		})

		mm := rmm.(*netproto.AreaConfig)

		hsreg.AreaConfigMap[mm.GetAreaID()] = mm

		if !mm.GetNotifyRefresh() {
			hsreg.ChanAreaConfigStatus <- true
		}
	}

}

// SendDBServerStartMsg 通知服務器啟動
func (hsreg *HallServerRegHandler) SendDBServerStartMsg() {
	svrmsg := new(msg.Message)
	svrmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	svrmsg.SClassID = int32(netproto.DBServerClassID_HallServerStartID)

	svrdata := new(netproto.HallServerStart)
	svrdata.ClientAddr = proto.String(hsreg.clientAddr)
	svrdata.NotifyAddr = proto.String(hsreg.notifyAddr)
	svrdata.ServerID = proto.Int32(hsreg.serverid)
	svrdata.HttpAddr = proto.String(hsreg.httpAddr)

	svrmsg.MsgData = svrdata
	hsreg.rpcmsgqueue.PushMessage(netproto.CreateEmptyRPCInfo(hsreg.serverid), false, svrmsg)

}

// SendAreaConfigRefreshMsg 刷新AreaConfig暫存資料
func (hsreg *HallServerRegHandler) SendAreaConfigRefreshMsg() {
	for _, v := range hsreg.AreaConfigMap {
		svrmsg := new(msg.Message)
		svrmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
		svrmsg.SClassID = int32(netproto.DBServerClassID_GetAreaConfigID)
		svrdata := new(netproto.GetHallAreaMsg)
		svrdata.AreaID = proto.Int64(v.GetAreaID())
		svrdata.NotifyRefresh = proto.Bool(true)
		svrmsg.MsgData = svrdata
		hsreg.rpcmsgqueue.PushMessage(netproto.CreateEmptyRPCInfo(hsreg.serverid), false, svrmsg)
	}
}

// RPCHallVerProvider ?
type RPCHallVerProvider struct {
	hallRegHandler *HallServerRegHandler
	PlayerManager  *player.Manager
	evtChans       chan func()
}

// NewRPCHallVerProvider 創建平台版本提供器對象
func NewRPCHallVerProvider(hallRegHandler *HallServerRegHandler, playerManager *player.Manager) *RPCHallVerProvider {
	rpch := new(RPCHallVerProvider)
	rpch.hallRegHandler = hallRegHandler
	rpch.PlayerManager = playerManager

	rpch.evtChans = make(chan func(), 10000)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				callfunc := <-rpch.evtChans
				if callfunc != nil {
					util.SafeCall(callfunc)
				}
			}
		}()
	}

	return rpch
}

// GetShowGameFlagByArea ?
func (p *RPCHallVerProvider) GetShowGameFlagByArea(areaID int64) []*netproto.AreaShowGameFlag {
	resultList := make([]*netproto.AreaShowGameFlag, 0)
	for _, v := range p.hallRegHandler.hallServerRegRet.AreaShowGameFlagList {
		if v.GetAreaID() == int32(areaID) {
			resultList = append(resultList, v)
		}
	}
	return resultList
}

// GetGameServerListByArea ?
func (p *RPCHallVerProvider) GetGameServerListByArea(areaid int64) []*netproto.PlatformGameServer {
	platformID := p.GetPlatformIDByAreaID(areaid)
	if gameServerList, ok := p.hallRegHandler.PlatformGameServerMap.Load(platformID); ok {
		return gameServerList.([]*netproto.PlatformGameServer)
	}
	return []*netproto.PlatformGameServer{}
}

// GetVer ?
func (p *RPCHallVerProvider) GetVer(ctx *module.MouduleCtx, clt *network.ConnInfo, areaID int64) *netproto.HallVerInfo {
	v := new(netproto.HallVerInfo)

	regret := p.hallRegHandler.hallServerRegRet

	if hallver, ok := p.hallRegHandler.AreaHallVersionMap.Load(areaID); ok && regret.GetCode() == 1 {
		hallver = hallver.(*netproto.AreaHallVersion)
		v.AndroidVersion = hallver.(*netproto.AreaHallVersion).AndroidVersion
		v.IOSVersion = hallver.(*netproto.AreaHallVersion).IOSVersion
		v.AndroidVersion2 = hallver.(*netproto.AreaHallVersion).AndroidVersion2
		v.IOSVersion2 = hallver.(*netproto.AreaHallVersion).IOSVersion2
		v.UpdateDomain = proto.String(strings.Split(regret.GetUpdateDomain(), ",")[0]) //原熱更域名, 用逗號切開取第一個
		v.BillDomain = regret.BillDomain
		v.ClubDomain = regret.ClubDomain
		v.CustomerDomain = regret.CustomerDomain
		v.UpdateDomain2 = regret.UpdateDomain //選線用, 用逗號分隔
	} else if !ok {
		mlog.Error("取得大廳版號失敗, 找不到渠道對應的版號, areaId=%d", areaID)
	} else {
		mlog.Error("大廳服務器注冊資訊不存在")
	}

	ok, areaConfigList := p.GetAreaConfig(areaID)

	if ok {
		v.AreaConfig = areaConfigList.AreaConfig
		v.GiftMoney = areaConfigList.GiftMoney
		v.DefaultMoney = areaConfigList.DefaultMoney
		return v
	}
	p.SendDBAreaIDMsg(areaID)

	if !ok {
		if pr := p.PlayerManager.GetPlayerByConnID(clt.ConnID()); pr != nil || !pr.IsGuest() {
			evtfunc := func() {
				select {
				case <-p.hallRegHandler.ChanAreaConfigStatus:
					ok, areaConfigList = p.GetAreaConfig(areaID)
					if ok {
						v.AreaConfig = areaConfigList.AreaConfig
					}

					vermsg := new(msg.Message)
					vermsg.BClassID = int32(netproto.MessageBClassID_Hall)
					vermsg.SClassID = int32(netproto.HallMsgClassID_HallVerInfoID)
					vermsg.MsgData = v

					ctx.Server.SendMessage(clt, vermsg)
				}
			}
			p.evtChans <- evtfunc
		}
	}

	return nil
}

// SendDBAreaIDMsg 大廳取得DB的AredID, AreaConfig
func (p *RPCHallVerProvider) SendDBAreaIDMsg(areaID int64) {
	svrmsg := new(msg.Message)
	svrmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	svrmsg.SClassID = int32(netproto.DBServerClassID_GetAreaConfigID)

	svrdata := new(netproto.GetHallAreaMsg)
	svrdata.AreaID = proto.Int64(areaID)
	svrdata.NotifyRefresh = proto.Bool(false)

	svrmsg.MsgData = svrdata
	p.hallRegHandler.rpcmsgqueue.PushMessage(netproto.CreateEmptyRPCInfo(p.hallRegHandler.serverid), false, svrmsg)
}

// GetAreaConfig 取得AreaConfig
func (p *RPCHallVerProvider) GetAreaConfig(areaID int64) (bool, *netproto.AreaConfig) {
	areaConfigList := new(netproto.AreaConfig)

	if val, ok := p.hallRegHandler.AreaConfigMap[areaID]; ok {
		return true, val
	}

	return false, areaConfigList
}

// GetPlatformIDByAreaID ?
func (p *RPCHallVerProvider) GetPlatformIDByAreaID(areaID int64) int64 {
	if platformID, ok := p.hallRegHandler.AreaPlatformMap.Load(areaID); ok {
		return platformID.(int64)
	}
	return 0
}

// RPCHallVerLoader ?
type RPCHallVerLoader struct {
	hallRegHandler *HallServerRegHandler
}

// NewRPCHallVerLoader 創建平台版本加載器對象
func NewRPCHallVerLoader(hallRegHandler *HallServerRegHandler) *RPCHallVerLoader {
	verLoader := new(RPCHallVerLoader)
	verLoader.hallRegHandler = hallRegHandler

	return verLoader
}

// LoadHallVer ?
func (ver *RPCHallVerLoader) LoadHallVer() {
	ver.hallRegHandler.SendDBServerStartMsg()
	ver.hallRegHandler.SendAreaConfigRefreshMsg()
}
