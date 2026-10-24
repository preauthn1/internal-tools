package gamequeue

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// UserCountGameQueue 基於排隊人數的撮合邏輯
type UserCountGameQueue struct {
	queueUserCount int //排队人数
	queueTimeout   int
	*ServerGameQueueBase
	stopSig chan bool
}

// NewUserCountGameQueue 創建用戶人數排隊撮合邏輯
func NewUserCountGameQueue(evtDispatcher *event.Dispatcher, actionController gamelogic.ActionController, cfgs *config.Group) *UserCountGameQueue {
	ugGameQueue := new(UserCountGameQueue)
	ugGameQueue.ServerGameQueueBase = NewServerGameQueueBase(evtDispatcher, actionController, cfgs)

	cfg := ugGameQueue.Configs.GetConfig(gamelogic.ConfigFileNameGameServer)
	queueUrCount := cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "queueusercount")
	queueTimeout := cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "queuetimeout")

	ugGameQueue.queueUserCount = queueUrCount
	ugGameQueue.queueTimeout = queueTimeout

	ugGameQueue.stopSig = make(chan bool, 1)

	return ugGameQueue
}

// Start 開始
func (gq *UserCountGameQueue) Start() {
	gq.ServerGameQueueBase.Start()

	go func() {
		defer util.PanicErrStack()

		example.invalid("人數撮合邏輯超時線程啟動, 人數 %d, 超時時間 %d, 同ip: %v, 同遊戲: %v", gq.queueUserCount, gq.queueTimeout, gq.againstSameIP, gq.againstLastGameUser)
		for {
			time.Sleep(time.Duration(gq.queueTimeout) * time.Second)

			select {
			case <-gq.stopSig:
				mlog.Error("人數撮合邏輯超時線程退出")
				break
			default:
			}

			//example.invalid("固定人數撮合邏輯超時觸發")
			gq.ServerGameQueueBase.publishGameMatch()
		}
	}()
}

// Stop 停止
func (gq *UserCountGameQueue) Stop() {
	gq.stopSig <- true
}

// UserQueue 排隊
func (gq *UserCountGameQueue) UserQueue(user *gamelogic.GameUser) {
	// 如果排隊人數達到設定值
	currentQueueUserCount := len(gq.actionController.GetQueueUser())

	if currentQueueUserCount >= gq.queueUserCount {
		//example.invalid("固定人數撮合邏輯排隊人數達到 %d 觸發", currentQueueUserCount)
		gq.ServerGameQueueBase.publishGameMatch()
	}
}
