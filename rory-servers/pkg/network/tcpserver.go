package network

import (
	"fmt"
	"io"
	"net"
	"sync"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// MesasgePackFun 消息分包器生成函數
type MesasgePackFun func() MsgPacker

// TCPServer tcp服務器
type TCPServer struct {
	netListener   net.Listener
	msgPackerFun  MesasgePackFun
	cfg           *ServerCfg
	evtChan       *ServerNetEventChan
	conns         ConnSet
	connLocker    *sync.RWMutex
	commpacker    MsgPacker
	stopSig       chan bool
	isServerStart bool
}

// NewTCPServer 創建TCP服務器對象
func NewTCPServer(msgPacker MesasgePackFun, cfg *ServerCfg) (*TCPServer, error) {
	svr := new(TCPServer)
	svr.msgPackerFun = msgPacker           // 創建消息分包方法
	svr.conns = make(map[string]*ConnInfo) // 當前客戶端連接集合
	svr.connLocker = new(sync.RWMutex)     // 客戶端連接集合同步鎖

	svr.cfg = cfg
	err := cfg.CheckCfg()

	if err != nil {
		return nil, err
	}

	// 默認的消息打包對象，用於發送消息的時候進行打包
	svr.commpacker = svr.msgPackerFun()

	svr.stopSig = make(chan bool, 1)

	return svr, nil
}

// Start 啟動服務器
func (server *TCPServer) Start(evtchan *ServerNetEventChan) error {
	var err error

	server.netListener, err = net.Listen("tcp", server.cfg.Addr)
	if err != nil {
		return err
	}

	server.evtChan = evtchan

	go func() {
		defer util.PanicErrStack()

		for {
			select {
			default:
				conn, err := server.netListener.Accept()
				if err != nil {
					continue
				}
				tconn := NewTCPConn(conn, conf.READBUFLEN) // 將net.Conn 轉換成 TcpConn
				cinfo := NewConnInfo(tconn, server.cfg.ServerID, TCP)

				go server.handleConn(cinfo)
			case <-server.stopSig:
				mlog.System("服務器 %s 偵聽 %s 線程退出.", server.cfg.ServerID, server.cfg.Addr)
				break
			}
		}
	}()

	mlog.System("TCP 服務器啟動成功, 偵聽地址 %s", server.cfg.Addr)
	// 設置服務器狀態
	server.isServerStart = true
	return nil
}

// 處理連接
func (server *TCPServer) handleConn(cinfo *ConnInfo) {
	defer func() {
		if err := recover(); err != nil {
			server.CloseConn(cinfo)

			mlog.Error("客戶端連接 %s 發生錯誤 %v", cinfo, err)
		}
	}()

	// 檢查連接是否超出限制
	server.connLocker.Lock()
	if len(server.conns) >= server.cfg.MaxConnNum {
		server.connLocker.Unlock()

		// 關閉客戶端連接
		cinfo.Close()

		mlog.Error("連接數量超出限制 %d, 關閉連接 %s, ServerID[%v]", server.cfg.MaxConnNum, cinfo.conn.RemoteAddr(), cinfo.ServerID)

		if cinfo.ServerID == conf.TCPSERVERNAMEHALLCLIENT {
			mlog.Fatal("超過連線負荷, 強制停止程序. 需手動重開.")
		}

		return
	}

	// 如果連接存在，先關閉
	if ecinfo, ok := server.conns[cinfo.connID]; ok {
		server.connLocker.Unlock()
		server.CloseConn(ecinfo)
		server.connLocker.Lock()
	}

	server.conns[cinfo.connID] = cinfo
	server.connLocker.Unlock()

	// 創建數據流分包處理器
	connpacker := server.msgPackerFun()

	// 通知客戶端連接
	mlog.Debug("客戶端 %v 連接 ...", cinfo.conn.RemoteAddr())
	//mlog.NewDefaultLogger().Info()
	server.evtChan.ConnectSig <- cinfo

	for {
		data, err := cinfo.conn.ReadMsg()
		if err != nil {
			server.CloseConn(cinfo)
			if err == io.EOF {
				example.invalid("客戶端 %v 關閉連接", cinfo.conn.RemoteAddr())
			} else {
				mlog.Error("客戶端 %v 讀取數據出現錯誤, 關閉連接, %v", cinfo.conn.RemoteAddr(), err)
			}

			return
		}

		if data != nil && len(data) > 0 {
			// 通知收到數據
			server.evtChan.DataSig <- &ConnData{cinfo, data}
			// 解包數據流
			var b []byte
			var err error
			b, err = connpacker.UnPack(data)

			if err != nil {
				mlog.Error("分包出現錯誤 %s, %v", cinfo, err)
				server.CloseConn(cinfo)
				return
			}

			for b != nil && len(b) > 0 {
				if b != nil && len(b) > 0 {
					if err != nil {
						mlog.Error("分包出現錯誤 %s, %v", cinfo, err)
						server.CloseConn(cinfo)
						return
					}

					server.evtChan.MsgSig <- &ConnData{cinfo, b}
					//example.invalid("收到客戶端 %v 數據包 %v", cinfo, b)
					b, err = connpacker.UnPack(nil)
				} else {
					b = nil
					err = nil
				}
			}
		} else {
			mlog.Error("接收到數據包長度為 0, 關閉客戶端 %s", cinfo)
			server.CloseConn(cinfo)
			return
		}

	}
}

// ServerID 服務器標識
func (server *TCPServer) ServerID() string {
	return server.cfg.ServerID
}

// CloseConn 關閉客戶端連接
func (server *TCPServer) CloseConn(cinfo *ConnInfo) {
	server.connLocker.Lock()
	if _, ok := server.conns[cinfo.connID]; !ok {
		server.connLocker.Unlock()
		return
	}
	delete(server.conns, cinfo.ConnID())
	server.connLocker.Unlock()

	cinfo.Close()

	server.evtChan.DisConnectSig <- cinfo
}

// CloseConnByConnID 關閉客戶端連接
func (server *TCPServer) CloseConnByConnID(connid string) {
	clt := server.GetConnByID(connid)

	if clt != nil {
		server.CloseConn(clt)
	}
}

// RangeConn 尋訪每個連接
func (server *TCPServer) RangeConn(f func(cinfo *ConnInfo)) {
	server.connLocker.RLock()
	defer server.connLocker.RUnlock()

	for _, ci := range server.conns {
		f(ci)
	}
}

// GetConnByID get connection by id
func (server *TCPServer) GetConnByID(connid string) *ConnInfo {
	server.connLocker.RLock()
	defer server.connLocker.RUnlock()

	return server.conns[connid]
}

// SendMessage 發送消息
func (server *TCPServer) SendMessage(cinfo *ConnInfo, data []byte) (int, error) {
	ndata := server.commpacker.Packet(data)
	n, err := cinfo.writeMsg(ndata)

	if err != nil {
		server.CloseConn(cinfo)
	}

	return n, err
}

// SendAllMessage send all message
func (server *TCPServer) SendAllMessage(data []byte) {
	ndata := server.commpacker.Packet(data)
	var errorConns []*ConnInfo

	server.connLocker.RLock()
	for _, ci := range server.conns {
		_, err := ci.writeMsg(ndata)
		if err != nil {
			errorConns = append(errorConns, ci)
		}
	}
	server.connLocker.RUnlock()

	if len(errorConns) > 0 {
		for _, ci := range errorConns {
			server.CloseConn(ci)
		}
	}
}

// ConnNum connection count
func (server *TCPServer) ConnNum() int {
	server.connLocker.RLock()
	defer server.connLocker.RUnlock()

	return len(server.conns)
}

// Stop terminet server
func (server *TCPServer) Stop() {
	if !server.isServerStart {
		return
	}

	server.connLocker.RLock()
	for _, ci := range server.conns {
		ci.Close()
	}
	server.connLocker.RUnlock()

	server.stopSig <- true

	server.netListener.Close()
	server.isServerStart = false
}

// ServerType return server type
func (server *TCPServer) ServerType() int {
	return TCP
}

func (server *TCPServer) String() string {
	return fmt.Sprintf("TCP Server Listen: %s, ServerID: %s, conn: %d", server.cfg.Addr, server.cfg.ServerID, server.ConnNum())
}
