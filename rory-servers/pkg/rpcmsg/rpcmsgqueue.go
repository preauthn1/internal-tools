package rpcmsg

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Queue 遠程消息隊列
type Queue struct {
	tcpMsgClient *msgserver.TCPMsgClient
	innerQueue   *MsgQueue
	evtchan      *msgserver.MsgEvtChan
	stopSig      chan bool
	rpcmsgParser msg.Parser
	OnRPCMesasge chan *msg.Message
}

// NewRPCMsgQueue 創建遠程消息隊列
func NewRPCMsgQueue() *Queue {

	msc := new(Queue)
	msc.rpcmsgParser = new(RPCMessageParser)

	evtchan := msgserver.NewMsgEvtChan(100)

	msc.tcpMsgClient = msgserver.NewTCPMsgClient(network.NewHeadMsgPacker(1024*500, true), evtchan, msc.rpcmsgParser, true)

	msc.innerQueue = NewMsgQueue(200000, 100000, 60, 0)
	msc.OnRPCMesasge = make(chan *msg.Message, 100)
	msc.evtchan = evtchan

	msc.stopSig = make(chan bool, 1)

	// 添加rpcmsgquque 命令
	command.AddCommand(NewRPCCommand(msc))

	return msc
}

// Connect 連接遠程地址
func (rq *Queue) Connect(addr string) error {
	err := rq.tcpMsgClient.Connect(addr)

	if err != nil {
		return err
	}

	go func() {
		defer util.PanicErrStack()

		for {
			select {
			case <-rq.stopSig:
				break
			default:
			}

			select {
			case clt := <-rq.evtchan.ConnectedSig:
				mlog.System("連接遠程RPC服務器成功, %v", clt)
			case clt := <-rq.evtchan.DisconnectedSig:
				mlog.System("遠程RPC服務器連接斷開, %v", clt)
			case <-rq.evtchan.DataSig:

			case msg := <-rq.evtchan.MsgSig:
				//example.invalid("收到遠程RPC服務器消息,%v", msg.MsgData)

				rpcmsg := msg.MsgData.(*RPCMessage)
				// 通知消息處理結果
				rq.innerQueue.DealMsg(int(rpcmsg.RPCInfo.GetQueueID()), MsgSendSuccess)

				// 通知回發消息
				rq.OnRPCMesasge <- msg
			}
		}
	}()

	return nil
}

// Start 啟動
func (rq *Queue) Start() {
	rq.innerQueue.Start()

	go func() {
		defer util.PanicErrStack()

		for {
			select {
			case <-rq.stopSig:
				break
			default:

			}

			// 檢查遠程連接狀態
			for {
				if rq.tcpMsgClient.IsConnected() {
					break
				} else {
					time.Sleep(time.Duration(100) * time.Millisecond)
				}
			}

			qmsg := <-rq.innerQueue.OnMessageSig
			if rq.tcpMsgClient.IsConnected() {

				// 設置隊列ID
				rpcmsg := qmsg.Message.MsgData.(*RPCMessage)
				qmsgid := int32(qmsg.MsgID)
				rpcmsg.RPCInfo.QueueID = &qmsgid

				err := rq.tcpMsgClient.SendMessage(qmsg.Message)
				//example.invalid("发送消息%v到RPC服务器.消息内容%v", qmsg, qmsg.Message)

				if err != nil {
					mlog.Error("發送消息 %v 到RPC服務器錯誤, %v", qmsg.Message, err)
				}

			}
		}
	}()
}

// Stop 停止
func (rq *Queue) Stop() {
	rq.innerQueue.Stop()

	rq.stopSig <- true
}

// PushMessage 添加消息到達送隊列
func (rq *Queue) PushMessage(rpcinfo *netproto.RPCInfo, safe bool, msg *msg.Message) error {
	rpcmsg := new(RPCMessage)
	rpcmsg.RPCInfo = rpcinfo
	rpcmsg.MsgData = msg.MsgData

	msg.MsgData = rpcmsg

	rq.innerQueue.PushMsg(msg, safe)

	return nil
}

// IsConnected 是否連結
func (rq *Queue) IsConnected() bool {
	return rq.tcpMsgClient.IsConnected()
}
