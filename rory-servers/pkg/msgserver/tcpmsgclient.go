package msgserver

import (
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/interceptor"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// TCPMsgClient tcp消息客戶端
type TCPMsgClient struct {
	tcpclient network.ConnClient
	msgparser msg.Parser
	clt       *network.ConnInfo
	stopSig   chan bool
}

// MsgEvtChan 消息事件通道
type MsgEvtChan struct {
	DataSig         chan []byte            // 數據消息通道
	MsgSig          chan *msg.Message      // 完整數據包消息通道
	ConnectedSig    chan *network.ConnInfo // 連接成功消息通道
	DisconnectedSig chan *network.ConnInfo // 連接斷開消息通道
}

// NewMsgEvtChan 創建消息事件通道對象
func NewMsgEvtChan(size int) *MsgEvtChan {
	evtchan := new(MsgEvtChan)
	evtchan.DataSig = make(chan []byte, size)
	evtchan.MsgSig = make(chan *msg.Message, size)
	evtchan.ConnectedSig = make(chan *network.ConnInfo, size)
	evtchan.DisconnectedSig = make(chan *network.ConnInfo, size)
	return evtchan
}

// NewTCPMsgClient 創建tcp消息客戶端
func NewTCPMsgClient(msgPacker network.MsgPacker, evt *MsgEvtChan, msgparser msg.Parser, autoReconnect bool) *TCPMsgClient {
	tmc := new(TCPMsgClient)
	tmc.stopSig = make(chan bool, 1)
	tcpevt := network.NewClientEventChan(conf.EVTCHANSIZE)

	// 不自動連接
	if !autoReconnect {
		tmc.tcpclient = network.NewTCPClient(msgPacker, tcpevt)
	} else { // 自動連接
		tcpclient := network.NewTCPClient(msgPacker, tcpevt)
		tmc.tcpclient = network.NewReConnectTCPClient(tcpclient, 3)
	}
	tmc.msgparser = msgparser

	go func() {
		util.PanicErrStack()
	ForEnd:
		for {
			select {
			case clt := <-tcpevt.ConnectedSig:
				evt.ConnectedSig <- clt
				tmc.clt = clt
				//example.invalid("tcpmsgclient connected...")
			case clt := <-tcpevt.DisconnectedSig:
				evt.DisconnectedSig <- clt
			case data := <-tcpevt.DataSig:
				evt.DataSig <- data
			case msgdata := <-tcpevt.MsgSig:
				msg, err := tmc.msgparser.Unmarshal(msgdata)
				if msg == nil {
					mlog.Error("消息解析錯誤 %v, %v", err, msgdata)
					return
				}
				mlog.Trace("收到消息%d_%d %v", msg.BClassID, msg.SClassID, msg.MsgData)
				if err != nil {
					mlog.Error("連接%v出現錯誤%v,連接關閉", tmc.clt, err)
					tmc.Close()
					evt.DisconnectedSig <- tmc.clt
					break
				} else {
					evt.MsgSig <- msg
				}
			case <-tmc.stopSig:
				break ForEnd
			}
		}

		mlog.Debug("TCP連接%v消息處理線程退出", tmc.clt)
	}()

	return tmc
}

// Connect 連接遠程地址
func (tmc *TCPMsgClient) Connect(addr string) error {
	return tmc.tcpclient.Connect(addr)
}

// Close 關閉連接
func (tmc *TCPMsgClient) Close() {
	tmc.tcpclient.Close()
}

// CloseAndStop 結束連接 optimize by rory
func (tmc *TCPMsgClient) CloseAndStop() {
	tmc.tcpclient.Close()
	tmc.stopSig <- true
}

// IsConnected 是否連接
func (tmc *TCPMsgClient) IsConnected() bool {
	return tmc.tcpclient.IsConnected()
}

// SendMessage 發送消息
func (tmc *TCPMsgClient) SendMessage(msg *msg.Message) error {

	{ // 發送前增加攔截器處理
		defer func() {
			for _, i := range interceptor.GetMSGInterceptor() {
				i.AfterSend(msg)
			}
		}()

		for _, i := range interceptor.GetMSGInterceptor() {
			i.BeforeSend(msg)
		}
	}

	data, err := tmc.msgparser.Marshal(msg)
	if err != nil {
		return err
	}
	return tmc.tcpclient.WriteMsg(data)
}

// SendData 發送二進制數據
func (tmc *TCPMsgClient) SendData(data []byte) error {
	return tmc.tcpclient.WriteMsg(data)
}
