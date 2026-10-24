package tfnet

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/internal/tfserver/router"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// ConnInfo 连接信息
type ConnInfo struct {
	ClientConn     net.Conn  //客户端连接
	ServerConn     net.Conn  //服务端连接
	ConnectTime    time.Time //连接建立时间
	RouteData      router.RouteData
	OnRoute        func(*ConnInfo) //路由成功事件
	ServerMsgQueue chan []byte
	HasCheckToken  bool //有經過token驗證
	Token          string
	TokenUpdater   *time.Ticker
	OnClose        func(*ConnInfo) //路由關閉事件
	IsClose        bool
	locker         *sync.RWMutex
}

// SetToken set token
func (conninfo *ConnInfo) SetToken(t string) {
	conninfo.Token = t
}

// RunTokenUpdater run token updater
func (conninfo *ConnInfo) RunTokenUpdater(s int) {
	conninfo.TokenUpdater = time.NewTicker(time.Duration(s) * time.Second)
	go func() {
		for range conninfo.TokenUpdater.C {
			b := conninfo.GenHeartBeatPack()
			if conninfo.ClientConn != nil {
				if _, err := conninfo.ClientConn.Write(b); err != nil {
					conninfo.Close()
					break
				}
			} else {
				conninfo.Close()
				break
			}
			if conninfo.ServerConn != nil {
				if _, err := conninfo.ServerConn.Write(b); err != nil {
					conninfo.Close()
					break
				}
			} else {
				conninfo.Close()
				break
			}
			safeCipherKeyStore.Set(conninfo.Token, time.Now().Unix())
			//mlog.System("%s", server.netFlow)
		}
		mlog.System("結束 token ticker 線程")
	}()
}

// GenHeartBeatPack get heart beat pack
func (conninfo *ConnInfo) GenHeartBeatPack() []byte {
	m := msg.NewMessage(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_HeartBeatID), nil)
	mb, _ := msg.NewProtoMsgParser().Marshal(m)
	mb = network.NewHeadMsgPacker(conf.MSGPACKERBUF, true).Packet(mb)
	return mb
}

// 连接信息格式化字符串
func (conninfo *ConnInfo) String() string {
	msg := ""
	if conninfo.ClientConn != nil {
		msg += conninfo.ClientConn.RemoteAddr().String()
	}
	msg += ","
	if conninfo.ServerConn != nil {
		msg += conninfo.ServerConn.RemoteAddr().String()
	}
	return msg
}

// CloseClient 关闭客户端连接
func (conninfo *ConnInfo) CloseClient() {
	if conninfo.ClientConn != nil {
		conninfo.ClientConn.Close()
	}
}

// CloseServer 关闭服务端连接
func (conninfo *ConnInfo) CloseServer() {
	if conninfo.ServerConn != nil {
		conninfo.ServerConn.Close()
	}
}

// Close 关闭连接
func (conninfo *ConnInfo) Close() {
	conninfo.locker.Lock()
	defer conninfo.locker.Unlock()
	if !conninfo.IsClose {
		conninfo.IsClose = true
		conninfo.CloseClient()
		conninfo.CloseServer()
		if conninfo.TokenUpdater != nil {
			conninfo.TokenUpdater.Stop()
		}
		if conninfo.OnClose != nil {
			conninfo.OnClose(conninfo)
		}
	}
}

// SetServerConn 设置服务器连接,同时触发路由成功事件
func (conninfo *ConnInfo) SetServerConn(conn net.Conn) {
	conninfo.ServerConn = conn
	conninfo.OnRoute(conninfo)
}

// IsBound 是否绑定了前后端
func (conninfo *ConnInfo) IsBound() bool {
	return conninfo.ServerConn != nil
}

// WaitServerConn WaitServerConn
func (conninfo *ConnInfo) WaitServerConn() bool {
	c := make(chan bool, 1)
	stopSig := make(chan bool, 1)
	defer func() {
		if conninfo.ServerConn != nil {
			for {
				select {
				case qmsg := <-conninfo.ServerMsgQueue:
					mlog.Debug("send qmsg")
					conninfo.ServerConn.Write(qmsg)
				default:
					mlog.Debug("queue clear")
					return
				}
			}
		}
		stopSig <- true
	}()
	go func() {
		for {
			select {
			case <-stopSig:
				mlog.System("stop")
				return
			default:
				if conninfo.ServerConn != nil {
					c <- true
					return
				}
			}
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}()
	select {
	case r := <-c:
		return r
	case <-time.After(time.Duration(3) * time.Second):
		return false
		// the read from ch has timed out
	}
}

// NewConnInfo create ConnInfo
func NewConnInfo(cltconn net.Conn) *ConnInfo {
	cinfo := &ConnInfo{}
	cinfo.ServerConn = nil
	cinfo.ClientConn = cltconn
	cinfo.ConnectTime = time.Now()
	cinfo.ServerMsgQueue = make(chan []byte, 10)
	cinfo.locker = &sync.RWMutex{}
	addr := cltconn.LocalAddr().String()
	sArrAddr := strings.Split(addr, ":")

	sip := sArrAddr[0]
	sport := sArrAddr[1]
	isport, _ := strconv.Atoi(sport)
	cinfo.RouteData = router.NewRouteData(sip, isport)
	return cinfo
}
