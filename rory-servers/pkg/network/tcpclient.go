package network

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ConnClient 連線設定
type ConnClient interface {
	Connect(addr string) error
	WriteMsg(data []byte) error
	Close()
	IsConnected() bool
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

// TCPClient tcp客戶端
type TCPClient struct {
	writeDataSig chan []byte // 發送數據消息通道
	msgPacker    MsgPacker   // 消息分包處理器
	conn         *ConnInfo   // 客戶端連接
	isConnected  bool        // 連接是否建立
	evtChan      *ClientEventChan
	innerLocker  *sync.RWMutex
	stopSig      chan bool
}

// NewTCPClient Create TcpClient
func NewTCPClient(msgPacker MsgPacker, evt *ClientEventChan) *TCPClient {
	clt := new(TCPClient)

	clt.writeDataSig = make(chan []byte, 100)

	clt.msgPacker = msgPacker
	clt.evtChan = evt
	clt.isConnected = false
	clt.innerLocker = &sync.RWMutex{}

	clt.stopSig = make(chan bool, 1)

	return clt
}

// Connect 連接服務器
func (clt *TCPClient) Connect(addr string) error {
	clt.innerLocker.Lock()
	defer clt.innerLocker.Unlock()

	if clt.isConnected {
		mlog.Error("連接已經建立, 但是又進行了連接. %v", clt.conn)
		return errors.New("連接已經建立")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	tconn := NewTCPConn(conn, conf.READBUFLEN)
	clt.conn = NewConnInfo(tconn, "client", TCP)

	clt.isConnected = true
	clt.evtChan.ConnectedSig <- clt.conn
	clt.stopSig = make(chan bool, 1) // 重置stopsig
	//example.invalid("tcpclient connected...")

	// 發送數據
	go func() {
		defer util.PanicErrStack()
		mlog.Debug("TPC Client %s 發送數據線程開始", clt.conn.String())
	ForEnd:
		for {
			var data []byte

			select {
			case data = <-clt.writeDataSig:
			case <-clt.stopSig:
				break ForEnd
			}

			b := clt.msgPacker.Packet(data)

			if b != nil && len(b) > 0 {
				_, err := clt.conn.writeMsg(b)

				if err != nil {
					mlog.NewDefaultLogger().SetData(fmt.Sprintf("[tcpclient]發送數據錯誤, %v...\r\n 發送數據線程退出 conn[%v]", err, clt.conn)).Error("HallServer壓測時")
					clt.Close()
				}
			}
		}

		mlog.Warn("TCP Client %s 發送數據線程退出", clt.conn.String())
	}()

	// 讀取數據
	go func() {
		defer util.PanicErrStack()
		mlog.Debug("TCP Client %s 讀取數據線程開始", clt.conn.String())

		for {
			data, err := clt.conn.readMsg()
			if err != nil {
				mlog.Warn("[tcpclient]讀取數據錯誤, %v...\r\n 程式將關閉連接 conn[%v]", err, clt.conn)
				clt.Close()

				break
			}

			if data != nil && len(data) > 0 {
				// 通知收到數據
				clt.evtChan.DataSig <- data
				var b []byte
				var err error

				// 解包數據流
				b, err = clt.msgPacker.UnPack(data)

				if err != nil {
					mlog.Error("分包出現錯誤, %v", err)
					clt.Close()
					break
				}

				for b != nil && len(b) > 0 {
					if err != nil {
						clt.Close()
						break
					}

					if b != nil && len(b) > 0 {
						clt.evtChan.MsgSig <- b
						b, err = clt.msgPacker.UnPack(nil)
					}
				}

			} else {
				clt.Close()
				break
			}
		}

		mlog.Warn("[tcpclient] %s 讀取數據線程退出", clt.conn.String())
	}()

	return nil
}

// WriteMsg 發送消息
func (clt *TCPClient) WriteMsg(data []byte) error {
	clt.innerLocker.Lock()
	defer clt.innerLocker.Unlock()

	if !clt.isConnected {
		return errors.New("not connect")
	}
	//non-blocking send work because the channel has no buffer and there is no receiver.
	select {
	case clt.writeDataSig <- data:
	default:
		return errors.New("can not send message because the channel has no buffer")
	}

	return nil
}

// Close 關閉連接
func (clt *TCPClient) Close() {
	clt.innerLocker.Lock()
	defer clt.innerLocker.Unlock()
	if !clt.isConnected {
		mlog.Error("連接已經關閉, 但是又進行了關閉. %v", clt.conn)
		return
	}

	clt.stopSig <- true

	clt.conn.Close()
	clt.isConnected = false

	clt.evtChan.DisconnectedSig <- clt.conn
}

// IsConnected 是否連接
func (clt *TCPClient) IsConnected() bool {
	clt.innerLocker.Lock()
	defer clt.innerLocker.Unlock()

	return clt.isConnected
}

// LocalAddr 本地位址
func (clt *TCPClient) LocalAddr() net.Addr {
	clt.innerLocker.Lock()
	defer clt.innerLocker.Unlock()

	return clt.conn.GetConn().LocalAddr()
}

// RemoteAddr 遠端位址
func (clt *TCPClient) RemoteAddr() net.Addr {
	clt.innerLocker.Lock()
	defer clt.innerLocker.Unlock()

	return clt.conn.GetConn().RemoteAddr()
}

// ReConnectTCPClient 自動重連的TCP客戶端
type ReConnectTCPClient struct {
	innerTCPClient *TCPClient
	reconnInterval int
}

// NewReConnectTCPClient Create ReConnectTcpClient
func NewReConnectTCPClient(clt *TCPClient, reconnInterval int) *ReConnectTCPClient {
	rclt := new(ReConnectTCPClient)
	rclt.innerTCPClient = clt
	rclt.reconnInterval = reconnInterval

	return rclt
}

// Connect 連接服務器
func (clt *ReConnectTCPClient) Connect(addr string) error {
	err := clt.innerTCPClient.Connect(addr)

	go func() {
		defer util.PanicErrStack()
		for {
			time.Sleep(time.Second * time.Duration(clt.reconnInterval))

			if !clt.innerTCPClient.IsConnected() {
				mlog.Trace("開始自動重連服務器 %s", addr)
				clt.innerTCPClient.Connect(addr)
			}
		}
	}()

	if err != nil {
		return err
	}

	return nil
}

// WriteMsg 發送消息
func (clt *ReConnectTCPClient) WriteMsg(data []byte) error {
	return clt.innerTCPClient.WriteMsg(data)
}

// Close 關閉連接
func (clt *ReConnectTCPClient) Close() {
	clt.innerTCPClient.Close()
}

// IsConnected 是否連接
func (clt *ReConnectTCPClient) IsConnected() bool {
	return clt.innerTCPClient.isConnected
}

// LocalAddr 本地位址
func (clt *ReConnectTCPClient) LocalAddr() net.Addr {
	return clt.innerTCPClient.LocalAddr()
}

// RemoteAddr 遠端位址
func (clt *ReConnectTCPClient) RemoteAddr() net.Addr {
	return clt.innerTCPClient.RemoteAddr()
}
