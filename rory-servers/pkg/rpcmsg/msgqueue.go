package rpcmsg

import (
	"fmt"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
const (
	MsgSendedWaitReturn = iota // 已發送等待返回結果
	MsgDroped                  // 已丟棄，超時未收到返回結果
	MsgSendSuccess             // 發送並處理成功
)

// QueuedMsgInfo 隊列消息信息
type QueuedMsgInfo struct {
	Message      *msg.Message // 消息內容
	MsgID        int          // 隊列ID
	PushTime     time.Time    // 加入時間
	RetryTimes   int          // 重試次數
	LastSendTime time.Time    // 最後一次發送時間
	Status       int          // 狀態
	Safe         bool         // 是否安全消息，即永不丟棄
}

// 隊列消息字符串描述
func (msg *QueuedMsgInfo) String() string {
	statusdesc := ""
	if msg.Status == MsgSendedWaitReturn {
		statusdesc = "已發送"
	} else if msg.Status == MsgDroped {
		statusdesc = "已丟棄"
	} else if msg.Status == MsgSendSuccess {
		statusdesc = "發送成功"
	}

	return fmt.Sprintf("protoid:%s,msgid:%d,pushtime:%s,retrytimes:%d,lastsendtime:%s,status:%s,safe:%v", msg.Message, msg.MsgID, msg.PushTime, msg.RetryTimes, msg.LastSendTime, statusdesc, msg.Safe)
}

// MsgQueue 消息隊列
type MsgQueue struct {
	maxSize          int                    // 最大消息長度
	innerWaitMsgList map[int]*QueuedMsgInfo // 等待處理消息
	OnMessageSig     chan *QueuedMsgInfo    // 消息通道
	waitmsgLocker    *sync.RWMutex          // 等待處理消息鎖
	currentMsgID     int                    // 當前消息ID
	stopSig          chan bool              // 停止信號
	retryInterval    int                    // 重發時間間隔
	maxRetryTimes    int                    // 最大重發次數
}

// NewMsgQueue 創建消息隊列
func NewMsgQueue(maxSize int, throttle int, retryInterval int, maxRetryTimes int) *MsgQueue {
	mq := new(MsgQueue)
	mq.maxSize = maxSize
	mq.currentMsgID = 0
	mq.waitmsgLocker = &sync.RWMutex{}
	mq.OnMessageSig = make(chan *QueuedMsgInfo, throttle)
	mq.stopSig = make(chan bool, 1)
	mq.retryInterval = retryInterval
	mq.maxRetryTimes = maxRetryTimes
	mq.innerWaitMsgList = make(map[int]*QueuedMsgInfo)

	return mq
}

// Start 啟動
func (mq *MsgQueue) Start() {
	go func() {
		defer util.PanicErrStack()
		for {
			select {
			default:
				time.Sleep(time.Second * time.Duration(31))

				// 需要丟棄的消息
				dealmsg := make([]int, 0)
				// 需要再次發送的消息
				msgtosend := make([]*QueuedMsgInfo, 0)

				mq.waitmsgLocker.Lock()

				for _, msg := range mq.innerWaitMsgList {
					if msg.RetryTimes >= mq.maxRetryTimes && !msg.Safe {
						// 丟棄消息
						mlog.Error("消息%v超出最大重試次數,丟棄", msg)

						//					go func() {
						//						mq.DealMsg(msg.MsgID, MsgDroped)
						//					}()

						dealmsg = append(dealmsg, msg.MsgID)

						continue
					}

					if int(time.Now().Sub(msg.LastSendTime).Seconds()) < mq.retryInterval {
						continue
					}

					msg.LastSendTime = time.Now()
					msg.RetryTimes++

					msgtosend = append(msgtosend, msg)
				}
				mq.waitmsgLocker.Unlock()

				// 處理丟棄消息
				for _, msgid := range dealmsg {
					mq.DealMsg(msgid, MsgDroped)
				}
				// 處理需要再次發送的消息
				for _, m := range msgtosend {
					mq.OnMessageSig <- m
				}
			case <-mq.stopSig: // 停止信号
				mlog.Warn("消息隊列退出...")
				return
			}
		}
	}()
}

// Stop 停止
func (mq *MsgQueue) Stop() {
	mq.stopSig <- true
}

// Size 隊列數量
func (mq *MsgQueue) Size() int {
	mq.waitmsgLocker.RLock()
	defer mq.waitmsgLocker.RUnlock()

	return len(mq.innerWaitMsgList)
}

// PushMsg 添加消息
func (mq *MsgQueue) PushMsg(msg *msg.Message, safe bool) {
	mq.waitmsgLocker.Lock()
	if len(mq.innerWaitMsgList) >= mq.maxSize {
		mq.waitmsgLocker.Unlock()
		mlog.Error("排隊消息超出最大長度")
		return
	}

	// 消息ID
	msgid := mq.currentMsgID + 1
	mq.currentMsgID++

	// 隊列信息
	msgqinfo := new(QueuedMsgInfo)
	msgqinfo.MsgID = msgid
	msgqinfo.PushTime = time.Now()
	msgqinfo.Message = msg
	msgqinfo.RetryTimes = 1
	msgqinfo.LastSendTime = time.Now()
	msgqinfo.Status = MsgSendedWaitReturn
	msgqinfo.Safe = safe

	// 重試機制,進入隊列
	if mq.maxRetryTimes > 0 || safe {
		// 加入排隊隊列
		mq.innerWaitMsgList[msgid] = msgqinfo
	}
	mq.waitmsgLocker.Unlock()

	mlog.Trace("Push Message Queue, MsgID: %d, Data: %v", msgid, msg)

	// 通知隊列消息
	mq.OnMessageSig <- msgqinfo

}

// DealMsg 處理消息
func (mq *MsgQueue) DealMsg(msgid int, status int) {
	mq.waitmsgLocker.Lock()
	delete(mq.innerWaitMsgList, msgid)
	mq.waitmsgLocker.Unlock()
}

// Find 查找消息
func (mq *MsgQueue) Find(msgid int) int {
	mq.waitmsgLocker.RLock()
	defer mq.waitmsgLocker.RUnlock()

	if _, ok := mq.innerWaitMsgList[msgid]; ok {
		return MsgSendedWaitReturn
	}

	return MsgSendSuccess
}
