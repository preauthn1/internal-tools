package heartbeat

import (
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Info 心跳信息
type Info struct {
	ConnID       string    // 連接ID
	LastBeatTime time.Time // 上次心跳時間
}

// Manager 心跳包處理器
type Manager struct {
	beatLog    map[string]time.Time
	msgServer  *msgserver.MsgServer
	beatLocker *sync.RWMutex
}

// NewHeartBeatManager 創建心跳包處理器
func NewHeartBeatManager(msgServer *msgserver.MsgServer) *Manager {
	bmh := new(Manager)
	bmh.beatLog = make(map[string]time.Time)
	bmh.msgServer = msgServer
	bmh.beatLocker = &sync.RWMutex{}

	// 檢查心跳超時的客戶端
	go func() {
		defer util.PanicErrStack()
		for {
			time.Sleep(time.Duration(BeatInterval) * time.Second)
			bmh.checkBeatTimeout()
		}
	}()

	return bmh
}

func (bmh *Manager) checkBeatTimeout() {
	bmh.beatLocker.Lock()
	defer bmh.beatLocker.Unlock()

	for connid, t := range bmh.beatLog {
		tt := time.Now().Sub(t)

		if int32(tt.Seconds()) >= ConnLoseTimeSpan {
			mlog.Trace("客戶端%v心跳包超時,斷開連接", connid)

			// 處理連上馬上斷開的情況
			if bmh.msgServer.GetClientByID(connid) == nil {
				delete(bmh.beatLog, connid)
			}

			bmh.msgServer.CloseConnByConnID(connid)
		}
	}
}

// IsHook 是否處理
func (bmh *Manager) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Common) && sclsid == int32(netproto.PlatformCommonClassID_HeartBeatID)
}

// OnMessage 處理訊息
func (bmh *Manager) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	bmh.beatLocker.Lock()
	defer bmh.beatLocker.Unlock()

	connid := clt.ConnID()

	if _, ok := bmh.beatLog[connid]; ok {
		bmh.beatLog[connid] = time.Now()
	}

	ctx.Trace("收到客戶端%v心跳消息", clt)

	// 回覆心跳消息
	hbmsg := new(msg.Message)
	hbmsg.BClassID = int32(netproto.MessageBClassID_Common)
	hbmsg.SClassID = int32(netproto.PlatformCommonClassID_HeartBeatReturnID)

	return []*msg.Message{hbmsg}

}

// OnClientDisconnect 客戶端連接斷開
func (bmh *Manager) OnClientDisconnect(connid string) {
	bmh.beatLocker.Lock()
	defer bmh.beatLocker.Unlock()

	delete(bmh.beatLog, connid)
}

// OnClientConnect 客戶端連接
func (bmh *Manager) OnClientConnect(connid string) {
	bmh.beatLocker.Lock()
	defer bmh.beatLocker.Unlock()

	bmh.beatLog[connid] = time.Now()
}
