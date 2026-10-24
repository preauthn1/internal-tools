package network

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ConnSet 客戶端連接集合
type ConnSet map[string]*ConnInfo

// Conn 客戶端連接
type Conn interface {
	ReadMsg() ([]byte, error)          // 接收數據
	WriteMsg(data []byte) (int, error) // 寫入數據
	LocalAddr() net.Addr               // 本地地址
	RemoteAddr() net.Addr              // 遠程地址
	Close()                            // 關閉連接
	Destory()
}

// 連接類型
const (
	TCP       = iota // TCP
	WebSocket        // WebSocket
	Mixed            // TCP+WebSocket
)

// ConnInfo 客戶端連接信息
type ConnInfo struct {
	IPAddress  string      // IP
	Port       int         // 端口
	conn       Conn        // 連接
	ServerID   string      // 所對應的服務器ID
	ConnTime   time.Time   // 連接建立時間
	connID     string      // 鏈接ID
	connLocker *sync.Mutex // 鎖
	ClientType int         // 客戶端類型
	IPRouted   bool        // ip是否已經上報
	ConnToken  string      // 連線口令
	AreaID     int64       // 渠道id
}

func (cinfo *ConnInfo) String() string {
	var addr string = "?"
	if cinfo.conn != nil {
		addr = cinfo.conn.RemoteAddr().String()
	}
	return fmt.Sprintf("Type: %s, Addr: %v", cinfo.ServerID, addr)
}

// ConnID 客戶端ID
func (cinfo *ConnInfo) ConnID() string {
	if cinfo == nil {
		return ""
	}
	return cinfo.connID
}

// NewConnInfo 建立客戶端連接信息
func NewConnInfo(conn Conn, sid string, ctype int) *ConnInfo {
	cinfo := new(ConnInfo)
	if conn != nil && conn.RemoteAddr() != nil {
		cinfo.IPAddress, cinfo.Port = util.GetIPAndPort(conn.RemoteAddr())
	} else {
		cinfo.IPAddress = "unknown"
		cinfo.Port = -1
	}
	cinfo.ServerID = sid
	cinfo.conn = conn
	cinfo.ConnTime = time.Now()
	cinfo.connID = fmt.Sprintf("%s:%d_%d", cinfo.IPAddress, cinfo.Port, cinfo.ConnTime.UnixNano()) // 客戶端ID,唯一標識客戶端
	cinfo.connLocker = new(sync.Mutex)
	cinfo.ClientType = ctype
	cinfo.IPRouted = false
	cinfo.ConnToken = ""

	return cinfo
}

// 發送消息
func (cinfo *ConnInfo) writeMsg(data []byte) (int, error) {
	cinfo.connLocker.Lock()
	defer cinfo.connLocker.Unlock()

	return cinfo.conn.WriteMsg(data)
}

func (cinfo *ConnInfo) readMsg() ([]byte, error) {
	return cinfo.conn.ReadMsg()
}

// Close 關閉連接
func (cinfo *ConnInfo) Close() {
	cinfo.connLocker.Lock()
	defer cinfo.connLocker.Unlock()

	cinfo.conn.Close()
}

// IsServerConn 是否指定服務器的連接
func (cinfo *ConnInfo) IsServerConn(serverid string) bool {
	return strings.ToLower(cinfo.ServerID) == strings.ToLower(serverid)
}

// GetConn 獲取連接
func (cinfo *ConnInfo) GetConn() Conn {
	return cinfo.conn
}

// ConnData 客戶端數據
type ConnData struct {
	Conn *ConnInfo
	Data []byte
}

// NetClient 網絡客戶端接口定義
type NetClient interface {
	Connect(addr string) error
	WriteMsg(data []byte) error
	Close()
	IsConnected() bool
}
