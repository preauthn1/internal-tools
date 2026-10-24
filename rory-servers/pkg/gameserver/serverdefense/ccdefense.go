package serverdefense

import (
	"fmt"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// CCDefense cc防護
type CCDefense struct {
	*DefenserBase
	ipConnNumLimit    int            // 同IP連接限制
	connectTimesLimit *DefenserLimit // 連接次數
	sendDataLenLimit  *DefenserLimit // 發送數據包總長度
	sendMsgLimit      *DefenserLimit // 發送消息次數

	connTimesDefData   map[string]*DefenserData // 連接次數防護數據
	sendDataLenDefData map[string]*DefenserData // 發送數據大小防護數據
	sendMsgDefData     map[string]*DefenserData // 發送消息數量防護數據

	connTimesLocker   *sync.RWMutex // 連接次數防護鎖
	sendDataLenLocker *sync.RWMutex // 發送數據包長度防護鎖
	sendMsgLocker     *sync.RWMutex // 發送消息防護數鎖
}

// NewCCDefense 創建cc防護對象
func NewCCDefense(ipConnNumLimit int, connTimesLimit *DefenserLimit, sendDataLenLimit *DefenserLimit, sendMsgLimit *DefenserLimit) *CCDefense {
	ccdef := new(CCDefense)
	ccdef.connectTimesLimit = connTimesLimit
	ccdef.sendDataLenLimit = sendDataLenLimit
	ccdef.sendMsgLimit = sendMsgLimit
	ccdef.ipConnNumLimit = ipConnNumLimit

	ccdef.connTimesDefData = make(map[string]*DefenserData)
	ccdef.sendDataLenDefData = make(map[string]*DefenserData)
	ccdef.sendMsgDefData = make(map[string]*DefenserData)

	ccdef.connTimesLocker = &sync.RWMutex{}
	ccdef.sendDataLenLocker = &sync.RWMutex{}
	ccdef.sendMsgLocker = &sync.RWMutex{}

	ccdef.DefenserBase = NewDefenserBase()

	// 5秒數據
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(5) * time.Second)

			// 連接次數
			ccdef.connTimesLocker.Lock()
			for _, d := range ccdef.connTimesDefData {
				d.ClearSecData5()
			}
			ccdef.connTimesLocker.Unlock()

			// 數據包長度
			ccdef.sendDataLenLocker.Lock()
			for _, d := range ccdef.sendDataLenDefData {
				d.ClearSecData5()
			}
			ccdef.sendDataLenLocker.Unlock()

			// 新聞數量
			ccdef.sendMsgLocker.Lock()
			for _, d := range ccdef.sendMsgDefData {
				d.ClearSecData5()
			}
			ccdef.sendMsgLocker.Unlock()

		}
	}()

	// 10秒數據
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(10) * time.Second)

			// 連接次數
			ccdef.connTimesLocker.Lock()
			for _, d := range ccdef.connTimesDefData {
				d.ClearSecData10()
			}
			ccdef.connTimesLocker.Unlock()

			// 數據包長度
			ccdef.sendDataLenLocker.Lock()
			for _, d := range ccdef.sendDataLenDefData {
				d.ClearSecData10()
			}
			ccdef.sendDataLenLocker.Unlock()

			// 消息數
			ccdef.sendMsgLocker.Lock()
			for _, d := range ccdef.sendMsgDefData {
				d.ClearSecData10()
			}
			ccdef.sendMsgLocker.Unlock()
		}

	}()

	// 30秒數據
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(30) * time.Second)

			// 連接次數
			ccdef.connTimesLocker.Lock()
			for _, d := range ccdef.connTimesDefData {
				d.ClearSecData30()
			}
			ccdef.connTimesLocker.Unlock()

			// 數據包長度
			ccdef.sendDataLenLocker.Lock()
			for _, d := range ccdef.sendDataLenDefData {
				d.ClearSecData30()
			}
			ccdef.sendDataLenLocker.Unlock()

			// 消息數
			ccdef.sendMsgLocker.Lock()
			for _, d := range ccdef.sendMsgDefData {
				d.ClearSecData30()
			}
			ccdef.sendMsgLocker.Unlock()
		}

	}()

	// 60秒數據
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(60) * time.Second)

			// 連接次數
			ccdef.connTimesLocker.Lock()
			for _, d := range ccdef.connTimesDefData {
				d.ClearSecData60()
			}
			ccdef.connTimesLocker.Unlock()

			// 數據包長度
			ccdef.sendDataLenLocker.Lock()
			for _, d := range ccdef.sendDataLenDefData {
				d.ClearSecData60()
			}
			ccdef.sendDataLenLocker.Unlock()

			// 消息數
			ccdef.sendMsgLocker.Lock()
			for _, d := range ccdef.sendMsgDefData {
				d.ClearSecData60()
			}
			ccdef.sendMsgLocker.Unlock()
		}

	}()

	// 300秒數據
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(300) * time.Second)

			// 連接次數
			ccdef.connTimesLocker.Lock()
			for _, d := range ccdef.connTimesDefData {
				d.ClearSecData300()
			}
			ccdef.connTimesLocker.Unlock()

			// 數據包長度
			ccdef.sendDataLenLocker.Lock()
			for _, d := range ccdef.sendDataLenDefData {
				d.ClearSecData300()
			}
			ccdef.sendDataLenLocker.Unlock()

			// 消息數
			ccdef.sendMsgLocker.Lock()
			for _, d := range ccdef.sendMsgDefData {
				d.ClearSecData300()
			}
			ccdef.sendMsgLocker.Unlock()
		}

	}()

	// 15分鐘,清理所有數據
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(60*15) * time.Second)

			// 連接次數
			ccdef.connTimesLocker.Lock()
			ccdef.connTimesDefData = make(map[string]*DefenserData)
			ccdef.connTimesLocker.Unlock()

			// 數據包長度
			ccdef.sendDataLenLocker.Lock()
			ccdef.sendDataLenDefData = make(map[string]*DefenserData)
			ccdef.sendDataLenLocker.Unlock()

			// 消息數
			ccdef.sendMsgLocker.Lock()
			ccdef.sendMsgDefData = make(map[string]*DefenserData)
			ccdef.sendMsgLocker.Unlock()
		}
	}()

	return ccdef
}

// OnClientConnect 客戶端連接
func (d *CCDefense) OnClientConnect(clt *network.ConnInfo) {
	ip := clt.IPAddress
	cn := d.GetConnNumByIP(ip)

	// ip白名單
	if d.checkWhiteIPList(ip) {
		return
	}

	// 如果IP連接數超出上線
	if cn >= d.ipConnNumLimit {
		d.CloseConn(clt)

		example.invalid("IP%s並發連接數達到%d,加入黑名單", ip, cn)
		// 封鎖IP
		d.BlackIPList.BlockIP(ip, d.BlockIPTime, fmt.Sprintf("並發連接數達到%d", cn))

		return
	}

	var dd *DefenserData

	d.connTimesLocker.Lock()
	// 獲取連接次數計數
	if defdata, ok := d.connTimesDefData[ip]; ok {
		dd = defdata
	} else {
		dd = NewDefenserData()
		d.connTimesDefData[ip] = dd
	}
	d.connTimesLocker.Unlock()

	// 如果連接次數超出設定
	if dd.CheckLimit(d.connectTimesLimit) {
		d.CloseConn(clt)

		example.invalid("ip %s連接頻率超出設定,加入黑名單.%s", ip, dd)
		d.BlackIPList.BlockIP(ip, d.BlockIPTime, fmt.Sprintf("連接次數超出設定,%s", dd))
	} else {
		// 增加連接計數
		dd.Increase(1)
	}
}

// OnClientDisconnect 客戶端斷開
func (d *CCDefense) OnClientDisconnect(clt *network.ConnInfo) {}

// OnClientData 收到數據
func (d *CCDefense) OnClientData(data *network.ConnData) {
	// ip白名單
	if d.checkWhiteIPList(data.Conn.IPAddress) {
		return
	}

	var dd *DefenserData

	connsign := data.Conn.IPAddress

	d.sendDataLenLocker.Lock()
	// 獲取發送數據包大小計數
	if defdata, ok := d.sendDataLenDefData[connsign]; ok {
		dd = defdata
	} else {
		dd = NewDefenserData()
		d.sendDataLenDefData[connsign] = dd
	}
	d.sendDataLenLocker.Unlock()

	// 增加計數
	if data.Data != nil {
		dd.Increase(len(data.Data))
	}

	if dd.CheckLimit(d.sendDataLenLimit) {
		d.CloseConn(data.Conn)

		ip := data.Conn.IPAddress

		example.invalid("ip%s發送數據包長度超出設定,加入黑名單.%s", ip, dd)
		d.BlackIPList.BlockIP(ip, d.BlockIPTime, fmt.Sprintf("單位時間內,發送數據包大小超出設定,%s", dd))
	}
}

// OnClientMessage 收到消息
func (d *CCDefense) OnClientMessage(clt *network.ConnInfo, msg *msg.Message) {
	// ip白名單
	if d.checkWhiteIPList(clt.IPAddress) {
		return
	}

	var dd *DefenserData

	connsign := clt.IPAddress

	d.sendMsgLocker.Lock()
	// 獲取發送消息
	if defdata, ok := d.sendMsgDefData[connsign]; ok {
		dd = defdata
	} else {
		dd = NewDefenserData()
		d.sendMsgDefData[connsign] = dd
	}
	d.sendMsgLocker.Unlock()

	// 增加消息數量計數
	dd.Increase(1)

	if dd.CheckLimit(d.sendMsgLimit) {
		d.CloseConn(clt)

		ip := clt.IPAddress

		example.invalid("ip%s發送消息頻率超出設定,加入黑名單%s", ip, dd)
		d.BlackIPList.BlockIP(ip, d.BlockIPTime, fmt.Sprintf("單位時間內,發送消息數量超出設定,%s", dd))
	}
}

// ParseConnTimesLimit 解析連接次數防護數據
func (d *CCDefense) ParseConnTimesLimit(str string) {
	lmt := d.parseDefenseLimit(str)
	if lmt == nil {
		mlog.Error("解析連接頻率防護數據錯誤,%s", str)
		return
	}

	d.connectTimesLimit = lmt
}

// ParseSendDataLimit 解析發送數據包長度防護數據
func (d *CCDefense) ParseSendDataLimit(str string) {
	lmt := d.parseDefenseLimit(str)
	if lmt == nil {
		mlog.Error("解析發送數據包長度防護數據錯誤,%s", str)
		return
	}

	d.sendDataLenLimit = lmt
}

// ParseSendMsgLimit 解析發送消息防護數據
func (d *CCDefense) ParseSendMsgLimit(str string) {
	lmt := d.parseDefenseLimit(str)
	if lmt == nil {
		mlog.Error("解析發送協議頻率防護數據錯誤,%s", str)
		return
	}

	d.sendMsgLimit = lmt
}
