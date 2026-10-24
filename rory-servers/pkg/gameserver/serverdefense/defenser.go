package serverdefense

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverdefense/blackip"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// ConnListProvider 客戶端連接上下文
type ConnListProvider interface {
	RangeConn(f func(conn *network.ConnInfo))
}

// ServerConnListProvider 服務器客戶端連接遍歷包裝類
type ServerConnListProvider struct {
	msgsvr *msgserver.MsgServer
}

// NewServerConnListProvider 創建服務器連接容器對象
func NewServerConnListProvider(msgsvr *msgserver.MsgServer) *ServerConnListProvider {
	connprovider := new(ServerConnListProvider)
	connprovider.msgsvr = msgsvr

	return connprovider
}

// RangeConn range each conntection
func (svrconnlst *ServerConnListProvider) RangeConn(f func(conn *network.ConnInfo)) {
	svrconnlst.msgsvr.RangeConn(f)
}

// Defenser 防護接口
type Defenser interface {
	msgserver.LogicHandler

	// 設置黑名單列表容器
	SetIPListProvider(iplist blackip.ListProvider)
	// 設置客戶端連接容器
	SetConnProvider(connprovider ConnListProvider)
}

// DefenserBase 防護基類
type DefenserBase struct {
	BlackIPList   blackip.ListProvider        // ip黑名單
	ConnProvider  ConnListProvider            // 客戶端連接容器
	CloseConnFunc func(clt *network.ConnInfo) // 斷開客戶端連接
	BlockIPTime   int                         // 封鎖時間
	WhiteIPList   string                      // ip白名單
}

// NewDefenserBase 防護基類
func NewDefenserBase() *DefenserBase {
	dd := new(DefenserBase)
	dd.BlockIPTime = 15 * 60

	return dd
}

// SetIPListProvider set ip list provider
func (db *DefenserBase) SetIPListProvider(iplist blackip.ListProvider) {
	db.BlackIPList = iplist
}

// SetConnProvider set connection provider
func (db *DefenserBase) SetConnProvider(connprovider ConnListProvider) {
	db.ConnProvider = connprovider
}

// CloseClientFunc close client function
func (db *DefenserBase) CloseClientFunc(f func(clt *network.ConnInfo)) {
	db.CloseConnFunc = f
}

// CloseConn 踢出用戶
func (db *DefenserBase) CloseConn(clt *network.ConnInfo) {
	db.CloseConnFunc(clt)
}

// GetConnNumByIP 根據IP匯總連接數
func (db *DefenserBase) GetConnNumByIP(ip string) int {
	num := 0

	db.ConnProvider.RangeConn(func(conn *network.ConnInfo) {
		if conn.IPAddress == ip {
			num = num + 1
		}
	})

	return num
}

func (db *DefenserBase) parseDefenseLimit(limit string) *DefenserLimit {
	arr := strings.Split(limit, ",")

	if len(arr) != 5 {
		mlog.Error("解析防護配置錯誤,%s", limit)
		return nil
	}

	seclimit5, err := strconv.Atoi(arr[0])

	if err != nil {
		mlog.Error("解析防護配置錯誤,%s", limit)
		return nil
	}

	seclimit10, err := strconv.Atoi(arr[1])
	if err != nil {
		mlog.Error("解析防護配置錯誤,%s", limit)
		return nil
	}

	seclimit30, err := strconv.Atoi(arr[2])
	if err != nil {
		mlog.Error("解析防護配置錯誤,%s", limit)
		return nil
	}

	seclimit60, err := strconv.Atoi(arr[3])
	if err != nil {
		mlog.Error("解析防護配置錯誤,%s", limit)
		return nil
	}

	seclimit300, err := strconv.Atoi(arr[4])
	if err != nil {
		mlog.Error("解析防護配置錯誤,%s", limit)
		return nil
	}

	return NewDefenserLimit(seclimit5, seclimit10, seclimit30, seclimit60, seclimit300)
}

func (db *DefenserBase) checkWhiteIPList(ip string) bool {
	return strings.Index(db.WhiteIPList, ip+",") > 0
}
