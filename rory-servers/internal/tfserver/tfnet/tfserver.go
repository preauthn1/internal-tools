package tfnet

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/tfserver/tfconfig"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

//"gameserver/serverdefense"

//KEY: V1M5WE7X8WE2TWEK
//IV: DNW123EXGMTR83K7
var (
	locker             = &sync.RWMutex{}
	currentConn        = 0
	limit              = 100000
	lastConn           = 0 //上一次的連線數
	safeConnList       = util.NewSafeMap()
	safeCipherKeyStore = util.NewSafeMap()
)

// TFListenPortProvider 转发侦听端口接口
type TFListenPortProvider interface {
	LoadListenPort() []int
}

// TFServer 转发服务器
type TFServer struct {
	ListenPorts []int                //侦听端口
	listener    map[int]net.Listener //侦听器
}

// NewTFServer create TFServer
/**
whileiplist="0.0.0.0,0.0.0.0,"						//ip白名单
blocktime=900       							//封锁时间,单位秒
cc_ipmaxconn=21									//同IP最大连接数
cc_conntimes=500,1000,3000,5000,30000	 					//IP连接次数限制
cc_senddata=10240,20480,40960,61440,122880    	//IP发送数据包大小
cc_sendmsg=1000,2000,3000,5000,30000						//IP发送消息个数
proto_default=300,500,1000,3000,15000						//IP发送协议频率
*/
func NewTFServer(portProvider TFListenPortProvider) *TFServer {
	tfs := new(TFServer)
	tfs.ListenPorts = portProvider.LoadListenPort()
	tfs.listener = make(map[int]net.Listener)
	//cc防护
	//ccdef := serverdefense.NewCCDefense(100, nil, nil, nil)
	//ccdef.ParseConnTimesLimit("500,1000,3000,5000,30000")
	//ccdef.ParseSendDataLimit("10240,20480,40960,61440,122880")
	//ccdef.ParseSendMsgLimit("cc_sendmsg")
	//ccdef.BlockIPTime = blockiptime
	//ccdef.WhiteIPList = whiteiplist
	return tfs
}

// Start 启动服务器
func (server *TFServer) Start() {
	cfg := tfconfig.ConfigCtx
	ticker := time.NewTicker(time.Duration(cfg.ConnWatchInterval) * time.Second)
	go func() {
		for range ticker.C {
			diff := currentConn - lastConn
			if diff < 0 {
				diff = -diff
			}

			if diff > cfg.ConnWatchExceeded {
				errMsg := fmt.Sprintf("Alert : tfserver connections from %d to %d in %d seconds", lastConn, currentConn, cfg.ConnWatchInterval)
				tfconfig.AlertMsg(errMsg)
				mlog.Error(errMsg)
			}
			lastConn = currentConn
		}
	}()

	for _, p := range server.ListenPorts {
		server.listenPort(p)
	}

	example.invalid("服務器啟動成功...")
}

// 偵聽端口
func (server *TFServer) listenPort(p int) {

	if p == 0 {
		return
	}

	lisip := fmt.Sprintf(":%d", p)
	lis, err := net.Listen("tcp", lisip)
	if err != nil {
		mlog.Error("偵聽端口%d錯誤,%v", p, err)

		return
	}
	server.listener[p] = lis

	go func() {
		defer util.PanicErrStack()

		for {
			conn, err := lis.Accept()
			if err != nil {
				mlog.Error("端口侦听关闭,%v", err)
				break
			}
			if currentConn+1 > limit {
				mlog.Error("連線超過限制 %d, 拒絕連線", limit)
				conn.Close()
				continue
			}
			ip, _ := util.ParseIPAndPort(conn.RemoteAddr().String())
			nf := util.NewNetworkFilter(tfconfig.ConfigCtx.NeedCheckKeySourceIP)

			//创建连接上下文
			conninfo := NewConnInfo(conn)
			example.invalid("客户端%s建立连接", conninfo)
			if server.IsNeedCheckKey() && nf.CheckIPContains(ip) {
				go server.connCheckKey(conninfo)
			} else {
				go server.handleConn(conninfo)
			}

		}
	}()

	example.invalid("侦听端口%d成功...", p)
}

// 处理连接路由
func (server *TFServer) connCheckKey(conninfo *ConnInfo) {
	conn := conninfo.ClientConn
	defer util.PanicErrStack()
	sig := make(chan bool, 1)
	//透過客户端消息驗證連接合法性
	go func() {
		defer util.PanicErrStack()
		readbuff := make([]byte, 1024)

		for {
			_, err := conn.Read(readbuff)
			if err != nil {
				//mlog.System("[rory]error")
				//example.invalid("路由%s客户端连接断开,关闭通道", conninfo)
				conn.Close()
				return
			}
			msglen := util.BytesToInt(readbuff[:4], true)
			if msglen >= 10240 || msglen <= 0 {
				mlog.Warn("消息长度超出范围,%d 關閉連線", msglen)
				conn.Close()
				return
			}
			conninfo.ServerMsgQueue <- bytes.TrimRight(readbuff, "\x00")
			//管道处理消息
			if ok, token := checkKey(readbuff); ok {
				conninfo.SetToken(token)
				conninfo.RunTokenUpdater(60)
				sig <- true
				break
			}

		}
	}()

	select {
	case <-sig:
		conninfo.HasCheckToken = [REDACTED]
		server.handleConn(conninfo)
	case <-time.After(time.Duration(3) * time.Second):
		conninfo.Close()
		//conn.Close()
	}
}

// 处理连接路由
func (server *TFServer) handleConn(conninfo *ConnInfo) {
	defer util.PanicErrStack()
	//服务器消息路由
	fnSvrMsgRoute := func(conninfo *ConnInfo) {
		defer util.PanicErrStack()

		buff := make([]byte, 1024)
		for {
			n, err := conninfo.ServerConn.Read(buff)
			if err != nil {
				//example.invalid("路由%s通道,服务端连接断开,关闭通道", conninfo)
				closeConn(conninfo)

				onServerClose(conninfo)
				return
			}

			// 转发消息
			if conninfo.ClientConn != nil {
				conninfo.ClientConn.Write(buff[:n])
			}
		}
	}

	// 客户端IP上报
	fnCltIPMsg := func(conninfo *ConnInfo) {
		ipMsg := new(msg.Message)
		ipMsg.BClassID = int32(netproto.MessageBClassID_Common)
		ipMsg.SClassID = int32(netproto.PlatformCommonClassID_IPReportID)

		ipdata := new(netproto.IPReport)
		ip, _ := util.GetIPAndPort(conninfo.ClientConn.RemoteAddr())
		ipdata.IPAddress = proto.String(ip)
		ipMsg.MsgData = ipdata

		mb, _ := msg.NewProtoMsgParser().Marshal(ipMsg)
		if conninfo.ServerConn != nil {
			mb = network.NewHeadMsgPacker(conf.MSGPACKERBUF, true).Packet(mb)
			conninfo.ServerConn.Write(mb)
		}

	}

	// 路由成功事件
	conninfo.OnRoute = func(cinfo *ConnInfo) {
		// 是否代為上報IP
		if tfconfig.ConfigCtx.ReportIP == 1 {
			go fnCltIPMsg(cinfo)
		}

		go fnSvrMsgRoute(cinfo)
	}

	conninfo.OnClose = func(cinfo *ConnInfo) {
		locker.Lock()
		defer locker.Unlock()
		currentConn--
		example.invalid("連線斷開 %s", conninfo)
	}
	// 通知连接事件
	ChannelCtx.OnConnected(conninfo)

	locker.Lock()
	currentConn++
	locker.Unlock()

	// 转发客户端消息
	go func(conninfo *ConnInfo) {
		defer util.PanicErrStack()

		readbuff := make([]byte, 1024)
		for {
			n, err := conninfo.ClientConn.Read(readbuff)
			if err != nil {
				//example.invalid("路由%s客户端连接断开,关闭通道", conninfo)

				closeConn(conninfo)
				onClientClose(conninfo)
				return
			}

			// 管道处理消息
			ChannelCtx.OnMessage(conninfo, readbuff[:n])
			msglen := util.BytesToInt(readbuff[:4], true)
			//if msglen >= 10240 || msglen <= 0 {
			if msglen <= 0 {
				conninfo.Close()
				return
			}
			//mlog.System("msglen %d, bufflen %d", msglen, n)

			// 转发消息
			if conninfo.WaitServerConn() {
				conninfo.ServerConn.Write(readbuff[:n])
			}
		}
	}(conninfo)

}

// IsNeedCheckKey 是否需要檢查Key
func (server *TFServer) IsNeedCheckKey() bool {
	return tfconfig.ConfigCtx.NeedCheckKey == 1
}

// 驗證連接合法性
func checkKey(data []byte) (bool, string) {
	msgPacker := network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
	// 解包数据流
	b, err := msgPacker.UnPack(data)

	if err != nil {
		mlog.Error("解包数据流出現錯誤 %v", err)
		return false, ""
	}

	msgparser := msg.NewParserGroup()
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_ConnTokenID), msg.NewProtoMsgParser())
	msg, err := msgparser.Unmarshal(b)
	if msg == nil {
		mlog.Error("解析 CoonToken 出現錯誤 %v", err)
		return false, ""
	}
	var storeLife = 30 //通過驗證後在庫裡存活的時間 minutes
	var tokenLife = 3  //ConnToken 驗證的有效時間 minutes
	//mlog.System("%s 收到消息%v",conninfo, msg.MsgData)
	if msg.BClassID == int32(netproto.MessageBClassID_Common) && msg.SClassID == int32(netproto.PlatformCommonClassID_ConnTokenID) {
		connToken := new(netproto.ConnTokenMessage)
		byteMsgData := msg.MsgData.([]byte) //須將interface轉[]byte
		proto.Unmarshal(byteMsgData, connToken)
		ST := connToken.GetToken()
		val := safeCipherKeyStore.Get(ST)
		//example.invalid("取出 token %s, %v", ST, val != nil)
		if val != nil {
			tValid := time.Unix(val.(int64), int64(time.Duration(storeLife)*time.Minute))
			if time.Now().Before(tValid) {
				safeCipherKeyStore.Set(ST, time.Now().Unix())
				return true, ST
			}

			mlog.Warn("token 在 keystore 已超時 tValid %s, now %s", tValid, time.Now())
			safeCipherKeyStore.Del(ST)
			return false, ""
		}

		token, err := msgserver.ConnTokenDecryptor(ST)
		if err == nil {
			tValid := time.Unix(int64(token.Time), int64(time.Duration(tokenLife)*time.Minute))
			if time.Now().Before(tValid) {
				//example.invalid("存下token %s", ST)
				safeCipherKeyStore.Set(ST, time.Now().Unix())
				return true, ST
			}
			mlog.Warn("驗證連線 token 超時 tValid %s, now: %v", tValid, time.Now())
		} else {
			mlog.Error("驗證連線發生問題, %s", err)
		}
	}
	return false, ""
}

// DeniedClientConn DeniedClientConn
func DeniedClientConn(conninfo *ConnInfo) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Common)
	m.SClassID = int32(netproto.PlatformCommonClassID_TipMessageID)

	tm := new(netproto.TipMessage)
	tm.Message = proto.String("线路波动, 请稍候再试")
	tm.MsgType = proto.Int32(2)
	data, _ := proto.Marshal(tm)

	m.MsgData = data
	msgparser := msg.NewParserGroup()
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_TipMessageID), msg.NewProtoMsgParser())
	b, _ := msgparser.Marshal(m)
	msgPacker := network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
	b = msgPacker.Packet(b)
	conninfo.ClientConn.Write(b)
}

//关闭连接
//conninfo 连接信息
func closeConn(conninfo *ConnInfo) {
	//关闭连接
	conninfo.Close()
}

func onClientClose(conninfo *ConnInfo) {
	ChannelCtx.OnClientClose(conninfo)
}

func onServerClose(conninfo *ConnInfo) {
	ChannelCtx.OnServerClose(conninfo)
}

// AutoCheckPortTFServer 自动检查端口的转发服务器
type AutoCheckPortTFServer struct {
	innerTFServer *TFServer
	innerWSServer *WebSocketServer
	tfProvider    TFListenPortProvider
}

// NewAutoCheckPortTFServer 创建自动检查端口的转发服务器对象
func NewAutoCheckPortTFServer(portProvider TFListenPortProvider, checkInterval int) *AutoCheckPortTFServer {
	actfsvr := new(AutoCheckPortTFServer)
	if tfconfig.ConfigCtx.WebSocketServerPort != 0 {
		actfsvr.innerWSServer = NewWebSocketServer(portProvider)
	} else {
		actfsvr.innerTFServer = NewTFServer(portProvider)
		actfsvr.tfProvider = portProvider

		go func() {
			defer util.PanicErrStack()

			for {
				time.Sleep(time.Duration(checkInterval) * time.Second)

				example.invalid("开始检查转发端口...")
				actfsvr.checkPort()
			}
		}()
	}

	return actfsvr
}

//检查端口
func (server *AutoCheckPortTFServer) checkPort() {
	newpt := server.tfProvider.LoadListenPort()

	if newpt == nil || len(newpt) == 0 {
		return
	}

	portChanged := false

	//检查多余的端口
	for _, p := range server.innerTFServer.ListenPorts {
		if p == 0 {
			continue
		}

		//是否需要移除
		needremove := true
		for _, np := range newpt {
			if np == p {
				needremove = false
				break
			}
		}

		if needremove {
			if lis, lok := server.innerTFServer.listener[p]; lok {
				example.invalid("检查端口%d配置中已经不存在，关闭此端口", p)
				lis.Close()
				delete(server.innerTFServer.listener, p)

				//设置端口变化标志
				portChanged = true
			}
		}
	}

	//检查新增端口
	for _, inp := range newpt {
		exp := false

		for _, ep := range server.innerTFServer.ListenPorts {
			if inp == ep {
				exp = true
				break
			}
		}

		if !exp {
			server.innerTFServer.listenPort(inp)

			example.invalid("检查端口%d配置中新增，开启此端口", inp)

			//设置端口变化标志
			portChanged = true
		}
	}

	if portChanged {
		server.innerTFServer.ListenPorts = newpt
	} else {
		example.invalid("未检测到端口变化...")
	}
}

// Start 啟動
func (server *AutoCheckPortTFServer) Start() {
	if server.innerTFServer != nil {
		server.innerTFServer.Start()
	}
	if server.innerWSServer != nil {
		server.innerWSServer.Start()
	}

}
