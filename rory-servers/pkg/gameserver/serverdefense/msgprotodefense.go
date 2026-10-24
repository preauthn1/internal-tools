package serverdefense

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// MsgProtoDefense 消息協議防護
type MsgProtoDefense struct {
	*DefenserBase
	defaultProtoLimit *DefenserLimit            // 默認的協議發送頻率限制
	specialProtoLimit map[string]*DefenserLimit // 特殊協議發送頻率限制

	protoDefData    map[string]*DefenserData // 協議防護數據
	protoDataLocker *sync.RWMutex            // 協議數據鎖
}

// NewMsgProtoDefense 創建消息協議防護對象
func NewMsgProtoDefense(defaultLimit *DefenserLimit, specialLimit map[string]*DefenserLimit) *MsgProtoDefense {
	msgprotodef := new(MsgProtoDefense)
	msgprotodef.defaultProtoLimit = defaultLimit
	msgprotodef.specialProtoLimit = specialLimit

	msgprotodef.DefenserBase = NewDefenserBase()
	msgprotodef.protoDataLocker = &sync.RWMutex{}
	msgprotodef.protoDefData = make(map[string]*DefenserData)

	// 5秒協議頻率數據清理
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(5) * time.Second)

			msgprotodef.protoDataLocker.Lock()
			for _, dd := range msgprotodef.protoDefData {
				dd.ClearSecData5()
			}
			msgprotodef.protoDataLocker.Unlock()
		}
	}()

	// 10秒協議頻率數據清理
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(10) * time.Second)

			msgprotodef.protoDataLocker.Lock()
			for _, dd := range msgprotodef.protoDefData {
				dd.ClearSecData10()
			}
			msgprotodef.protoDataLocker.Unlock()
		}
	}()

	// 30秒協議頻率數據清理
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(30) * time.Second)

			msgprotodef.protoDataLocker.Lock()
			for _, dd := range msgprotodef.protoDefData {
				dd.ClearSecData30()
			}
			msgprotodef.protoDataLocker.Unlock()
		}
	}()

	// 60秒協議頻率數據清理
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(60) * time.Second)

			msgprotodef.protoDataLocker.Lock()
			for _, dd := range msgprotodef.protoDefData {
				dd.ClearSecData60()
			}
			msgprotodef.protoDataLocker.Unlock()
		}
	}()

	// 300秒協議頻率數據清理
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(300) * time.Second)

			msgprotodef.protoDataLocker.Lock()
			for _, dd := range msgprotodef.protoDefData {
				dd.ClearSecData300()
			}
			msgprotodef.protoDataLocker.Unlock()
		}
	}()

	// 15分鐘,清理數據
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(60*15) * time.Second)

			msgprotodef.protoDataLocker.Lock()
			msgprotodef.protoDefData = make(map[string]*DefenserData)
			msgprotodef.protoDataLocker.Unlock()
		}
	}()

	return msgprotodef
}

// AddProtoLimit 添加協議限制
func (dd *MsgProtoDefense) AddProtoLimit(bclsid int32, sclsid int32, limit *DefenserLimit) {
	protosign := fmt.Sprintf("%d_%d", bclsid, sclsid)

	dd.specialProtoLimit[protosign] = limit
}

// OnClientConnect 客戶端連接
func (dd *MsgProtoDefense) OnClientConnect(clt *network.ConnInfo) {}

// OnClientDisconnect 客戶端斷開
func (dd *MsgProtoDefense) OnClientDisconnect(clt *network.ConnInfo) {}

// OnClientData 收到數據
func (dd *MsgProtoDefense) OnClientData(data *network.ConnData) {}

// OnClientMessage 收到消息
func (dd *MsgProtoDefense) OnClientMessage(clt *network.ConnInfo, msg *msg.Message) {
	// ip白名單
	if dd.checkWhiteIPList(clt.IPAddress) {
		return
	}

	// 檢查特定協議頻率
	connsign := fmt.Sprintf("%s_%d_%d", clt.IPAddress, msg.BClassID, msg.SClassID)
	protosign := fmt.Sprintf("%d_%d", msg.BClassID, msg.SClassID)

	// 檢查特定協議
	if ddlimit, ok := dd.specialProtoLimit[protosign]; ok {
		dd.checkProtoData(clt, connsign, ddlimit, protosign)
		return
	}

	//	// 檢查大類協議
	//	bconnsign := fmt.Sprintf("%s_%d_*", clt.IPAddress, msg.BClassID)
	//	bprotosign := fmt.Sprintf("%d_*", msg.BClassID)

	//	// 檢查大類協議頻率
	//	if ddlimit, ok := dd.specialProtoLimit[bprotosign]; ok {
	//		dd.checkProtoData(clt, bconnsign, ddlimit, bprotosign)
	//		return
	//	}

	// 檢查默認協議
	ddlimit := dd.defaultProtoLimit
	dd.checkProtoData(clt, connsign, ddlimit, protosign)
}

func (dd *MsgProtoDefense) checkProtoData(clt *network.ConnInfo, connsign string, deflimit *DefenserLimit, protosign string) {
	// 獲取計數器
	var defdata *DefenserData
	dd.protoDataLocker.Lock()
	if pdd, ok := dd.protoDefData[connsign]; ok {
		defdata = pdd
	} else {
		defdata = NewDefenserData()
		dd.protoDefData[connsign] = defdata
	}
	dd.protoDataLocker.Unlock()

	// 計數器加1
	defdata.Increase(1)

	if deflimit == nil {
		return
	}

	// 檢查計數器是否超出設定
	if defdata.CheckLimit(deflimit) {
		dd.CloseConn(clt)

		dd.BlackIPList.BlockIP(clt.IPAddress, dd.BlockIPTime, fmt.Sprintf("單位時間發送協議%s超出設定次數.%s", protosign, defdata))
	}
}

// ParseDefaultProtoLimit 解析默認協議防護
func (dd *MsgProtoDefense) ParseDefaultProtoLimit(limit string) {
	lmt := dd.parseDefenseLimit(limit)
	if lmt == nil {
		return
	}

	dd.defaultProtoLimit = lmt
}

// ParseSpecialProtoLimit 解析特殊協議防護數據
func (dd *MsgProtoDefense) ParseSpecialProtoLimit(limit string) {
	defer func() {
		if err := recover(); err != nil {
			mlog.Error("解析特殊協議防護配置錯誤,%v", err)
		}
	}()

	dd.specialProtoLimit = make(map[string]*DefenserLimit)

	arr := strings.Split(limit, "|")

	for _, str := range arr {
		item := strings.Split(str, ":")

		// 協議號
		itemproto := item[0]
		// 防護配置
		itemlimit := item[1]

		itemprotoArr := strings.Split(itemproto, "_")
		bclsid, err := strconv.Atoi(itemprotoArr[0])

		if err != nil {
			mlog.Error("解析特殊協議防護配置錯誤,%s,%v", str, err)
			return
		}

		sclsid := itemprotoArr[1]

		protosign := fmt.Sprintf("%d_%s", bclsid, sclsid)
		protolimit := dd.parseDefenseLimit(itemlimit)

		if protolimit == nil {
			mlog.Error("解析特殊協議防護配置錯誤,%s", str)
			return
		}

		dd.specialProtoLimit[protosign] = protolimit
	}
}
