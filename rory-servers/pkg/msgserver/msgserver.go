package msgserver

import (
	"errors"
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 服務器運行模式
const (
	ServerModeDirect = 0 // 直連
	ServerModeRoute  = 1 // 轉發
)

// 斷線提示類型
const (
	TipMessageTypeIDHall int32 = 2 // 退出大廳
)

// MsgServer 服務器
type MsgServer struct {
	MsgParser    msg.Parser                  // 網絡消息解析器
	NetServer    network.NetServer           // 網絡服務器
	EventChannel *event.Dispatcher           // 事件分派器
	MsgPackerFun func() network.MsgPacker    // 消息分包處理器
	NetEventChan *network.ServerNetEventChan // 網絡服務器事件通道
	LogicHandler LogicHandler                // 服務器邏輯處理組件
	ServerMode   int                         // 服務器模式
}

// NewBgMsgServer 創建後端消息服務器
func NewBgMsgServer(throttle int) *MsgServer {
	return NewMsgServer(throttle, ServerModeDirect, 1024*500)
}

// NewFacadeMsgServer 創建客戶端消息服務器
func NewFacadeMsgServer(throttle int, servermode int) *MsgServer {
	return NewMsgServer(throttle, servermode, conf.MSGPACKERBUF)
}

// NewMsgServer 創建客戶端消息服務器
func NewMsgServer(throttle int, servermode int, msgMaxLen int32) *MsgServer {
	svr := new(MsgServer)

	// 消息解析器
	msgparser := msg.NewParserGroup()
	svr.MsgParser = msgparser
	svr.ServerMode = servermode

	// 事件分派器
	evtdispatcher := event.NewAsyncEventDispatcher(throttle)
	svr.EventChannel = evtdispatcher

	// 網絡服務器事件通道
	svr.NetEventChan = network.NewNetEventChan(conf.EVTCHANSIZE)

	// 消息分包處理器
	svr.MsgPackerFun = func() network.MsgPacker {
		return network.NewHeadMsgPacker(msgMaxLen, true)
	}

	// 服務器組
	svrg := network.NewNetServerGroup("svrgroup")
	svr.NetServer = svrg

	// 邏輯處理器
	logic := NewMsgServerLogicGroup()
	svr.LogicHandler = logic
	logic.closCltFunc = func(clt *network.ConnInfo) {

		svr.sendTipMessage(clt, "系统检测到您有异常的操作,拒绝了您当前的请求.请稍后再试!")

		go func() {
			time.Sleep(time.Duration(2) * time.Second)

			svr.NetServer.CloseConn(clt)
		}()
	}

	// 添加命令
	cmd := NewMsgServerCmd(svr)
	command.AddCommand(cmd)

	return svr
}

// InitNetServer 初始化服務器
func (s *MsgServer) InitNetServer() {}

// AddTCPServer 添加TCP服務器
func (s *MsgServer) AddTCPServer(cfg *network.ServerCfg) network.NetServer {
	// tcp服務器
	tcpsvr, _ := network.NewTCPServer(s.MsgPackerFun, cfg)
	s.GetServerGroup().AddServer(tcpsvr)

	return tcpsvr
}

// Start 啟動服務器
func (s *MsgServer) Start() error {

	go func() {
		for {
			select {
			case clt := <-s.NetEventChan.ConnectSig:
				// 直連模式,直接通知服務器邏輯處理.轉發模式需要等到IP上報
				if s.ServerMode == ServerModeDirect {
					s.LogicHandler.OnClientConnect(clt)
				}

				connevt := NewClientConnectEvt(clt)
				s.EventChannel.OnAsyncEvent(connevt, nil)

			case data := <-s.NetEventChan.DataSig:

				s.LogicHandler.OnClientData(data)

				dataevt := NewClientDataEvent(data.Conn, data.Data)
				s.EventChannel.OnAsyncEvent(dataevt, nil)

			case data := <-s.NetEventChan.MsgSig:
				msg, err := s.MsgParser.Unmarshal(data.Data)
				if err == nil && msg != nil {
					// 處理 Token 協議
					// 2019.3.8從tfapi改用遊戲盾因此不再需要檢查token.
					// 2019.11.29 重新開始使用token
					s.checkTokenMsg(data.Conn, msg)
					// 客戶端連接在轉發模式下工作,處理IP上報協議
					if s.ServerMode == ServerModeRoute && !data.Conn.IsServerConn(conf.TCPSERVERNAMEHALLGAMESERVER) {
						s.checkIPReportMsg(data.Conn, msg)
						// 轉發模式,IP沒上報,不處理任何消息
						if !data.Conn.IPRouted {
							mlog.Warn("客戶端%v未上報IP", data.Conn)
							//s.sendTipMessage(data.Conn, "非法操作.")
							continue
						}
					}

					s.LogicHandler.OnClientMessage(data.Conn, msg)

					msgevt := NewClientMessageEvent(data.Conn, msg)
					s.EventChannel.OnAsyncEvent(msgevt, nil)
				} else {
					mlog.Error("解析消息錯誤%s,關閉客戶端%v連接", err.Error(), data.Conn)
					util.AsyncSafeCall(func() {
						s.CloseConn(data.Conn)
					})
				}

			case clt := <-s.NetEventChan.DisConnectSig:
				s.LogicHandler.OnClientDisconnect(clt)

				disevt := NewClientDisconnectEvt(clt)
				s.EventChannel.OnAsyncEvent(disevt, nil)
			}
		}
	}()

	err := s.NetServer.Start(s.NetEventChan)
	if err != nil {
		return err
	}

	return nil
}

// 檢查是否為合法IP, 正確回傳 nil
func (s *MsgServer) checkValidIP(ip string) (err error) {
	_, _, err = util.ParseIPv4orIPv6(ip)
	return err
}

// 檢查IP上報協議
func (s *MsgServer) checkIPReportMsg(clt *network.ConnInfo, msg *msg.Message) {
	defer util.PanicErrStack()
	if msg.BClassID == int32(netproto.MessageBClassID_Common) && msg.SClassID == int32(netproto.PlatformCommonClassID_IPReportID) {
		dd := msg.MsgData.([]byte)
		mm := &netproto.IPReport{}

		err := proto.Unmarshal(dd, mm)

		if err != nil {
			mlog.Error("客戶端%v發送IP上報數據包解析錯誤,%s", clt, err.Error())
			s.NetServer.CloseConn(clt)
			return
		}

		// 如果IP已經上報
		if clt.IPRouted {
			return
		}

		// 驗證是否為合法IP
		var ip = mm.GetIPAddress()
		err = s.checkValidIP(ip)
		if err != nil {
			mlog.Warn("客戶端%v上報IP驗證錯誤,%s", clt, err.Error())
			// 上報失敗不做關閉連線
			// s.NetServer.CloseConn(clt)
			// return
		}

		// 設置客戶端真實IP
		clt.IPAddress = ip

		// 設置IP上報標記
		clt.IPRouted = true

		mlog.Debug("客戶端%v設置IP%s", clt, clt.IPAddress)

		// 通知客戶端連接
		s.LogicHandler.OnClientConnect(clt)

	}
}

// 處理token上報協議
// 從tfapi改用遊戲盾因此不再需要檢查token.
// 2019.11.29 重新開始使用token
func (s *MsgServer) checkTokenMsg(clt *network.ConnInfo, msg *msg.Message) {
	defer util.PanicErrStack()
	if msg.BClassID == int32(netproto.MessageBClassID_Common) && msg.SClassID == int32(netproto.PlatformCommonClassID_ConnTokenID) {
		dd := msg.MsgData.([]byte)
		mm := &netproto.ConnTokenMessage{}

		var err error
		var token *Token
		var ret *TFDataRet
		err = proto.Unmarshal(dd, mm)

		if err != nil {
			mlog.Error("客戶端%v發送Token數據包解析錯誤,%v", clt, err)
			s.NetServer.CloseConn(clt)
			return
		}
		// 解密連線口令
		token, err = ConnTokenDecryptor(mm.GetToken())
		if err != nil {
			mlog.Error("驗證連線發生問題, %s, tokenstr %s", err, mm.GetToken())
			//s.NetServer.CloseConn(clt) //先不掛掉連線
			return
		}

		var tValid = time.Unix(int64(token.Time), +int64(1*time.Minute))
		if time.Now().Before(tValid) {
			clt.ConnToken = [REDACTED]
			mlog.Debug("客戶端%v設置ConnToken%s", clt, clt.ConnToken)
		} else {
			mlog.Warn("驗證連線 token 超時 tValid %s", tValid)
		}

		ret, err = NodeDefDecryptor(mm.GetToken())
		if err != nil {
			return
		}
		var ip = ret.Data.UserIP
		err = s.checkValidIP(ip)
		if err != nil {
			mlog.Error("客戶端%v解析Token上報IP驗證錯誤,%s", clt, err.Error())
		}

		// 設置客戶端真實IP
		clt.IPAddress = ip
		// 設置IP上報標記
		clt.IPRouted = true
		mlog.Debug("客戶端%v設置IP%s", clt, clt.IPAddress)
		// 通知客戶端連接
		s.LogicHandler.OnClientConnect(clt)
	}
}

// 向客戶端發送斷線提示消息
func (s *MsgServer) sendTipMessage(clt *network.ConnInfo, str string) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Common)
	m.SClassID = int32(netproto.PlatformCommonClassID_TipMessageID)

	tm := new(netproto.TipMessage)
	tm.Message = proto.String(str)
	tm.MsgType = proto.Int32(TipMessageTypeIDHall)
	data, err := proto.Marshal(tm)

	if err != nil {
		mlog.Error("服務器防護,序列化斷線提示消息錯誤%v", err)
		return
	}

	m.MsgData = data

	s.SendMessage(clt, m)
}

// OnRPCMessageEvent 觸發RPC消息事件
func (s *MsgServer) OnRPCMessageEvent(msg *msg.Message) {
	evt := NewRPCMessageEvent(msg)

	s.EventChannel.OnAsyncEvent(evt, nil)
}

// Stop 停止服務器
func (s *MsgServer) Stop() {

	s.NetServer.Stop()
}

// SendMessage 發送消息
func (s *MsgServer) SendMessage(clt *network.ConnInfo, msg *msg.Message) (int, error) {
	data, err := s.MsgParser.Marshal(msg)

	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, errors.New("序列化消息錯誤,消息數據為空")
	}

	return s.NetServer.SendMessage(clt, data)
}

// SendData 發送數據
func (s *MsgServer) SendData(clt *network.ConnInfo, data []byte) (int, error) {
	return s.NetServer.SendMessage(clt, data)
}

// SendMessageByConnID 根據連接ID發送消息
func (s *MsgServer) SendMessageByConnID(connid string, msg *msg.Message) (int, error) {
	data, err := s.MsgParser.Marshal(msg)

	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, errors.New("序列化消息錯誤,消息數據為空")
	}

	clt := s.NetServer.GetConnByID(connid)

	if clt == nil {
		return 0, fmt.Errorf("連接%s不存在", connid)
	}

	return s.NetServer.SendMessage(clt, data)
}

// SendDataByConnID 根據連接ID發送數據
func (s *MsgServer) SendDataByConnID(connid string, data []byte) (int, error) {
	clt := s.NetServer.GetConnByID(connid)

	if clt == nil {
		return 0, fmt.Errorf("連接%s不存在", connid)
	}

	return s.NetServer.SendMessage(clt, data)
}

// SendAllMessage 廣播消息
func (s *MsgServer) SendAllMessage(msg *msg.Message) error {
	data, err := s.MsgParser.Marshal(msg)

	if err != nil {
		return err
	}

	if data == nil {
		return errors.New("序列化消息錯誤,消息數據為空")
	}

	s.NetServer.SendAllMessage(data)

	return nil
}

// SendAllData 廣播數據
func (s *MsgServer) SendAllData(data []byte) error {
	s.NetServer.SendAllMessage(data)

	return nil
}

// GetClientByID 獲取客戶端
func (s *MsgServer) GetClientByID(cid string) *network.ConnInfo {
	return s.NetServer.GetConnByID(cid)
}

// RangeConn 遍歷客戶端
func (s *MsgServer) RangeConn(f func(clt *network.ConnInfo)) {
	s.NetServer.RangeConn(f)
}

// CloseConn 關閉客戶端連接
func (s *MsgServer) CloseConn(clt *network.ConnInfo) {
	s.NetServer.CloseConn(clt)
}

// CloseConnByConnID 關閉客戶端連接
func (s *MsgServer) CloseConnByConnID(connid string) {
	clt := s.GetClientByID(connid)

	if clt != nil {
		s.CloseConn(clt)
	}
}

// GetServerGroup 獲取服務器組
func (s *MsgServer) GetServerGroup() *network.NetServerGroup {
	return s.NetServer.(*network.NetServerGroup)
}

// GetMsgParserGroup 獲取消息解析器集合
func (s *MsgServer) GetMsgParserGroup() *msg.ParserGroup {
	return s.MsgParser.(*msg.ParserGroup)
}

// GetLogicHandlerGroup 獲取邏輯處理組件集合
func (s *MsgServer) GetLogicHandlerGroup() *LogicHandlerGroup {
	return s.LogicHandler.(*LogicHandlerGroup)
}
