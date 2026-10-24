package gamequeue

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// IntervalMatchGameQueue 定時撮合
type IntervalMatchGameQueue struct {
	scope util.IntScoper //随机撮合时间范围
	*ServerGameQueueBase
	stopSig chan bool
}

// NewIntervalMatchGameQueue 創建定時撮合邏輯
func NewIntervalMatchGameQueue(evtDispatcher *event.Dispatcher, actionController gamelogic.ActionController, cfgs *config.Group) *IntervalMatchGameQueue {
	intervalGameQueue := new(IntervalMatchGameQueue)
	intervalGameQueue.ServerGameQueueBase = NewServerGameQueueBase(evtDispatcher, actionController, cfgs)

	cfg := intervalGameQueue.Configs.GetConfig(gamelogic.ConfigFileNameGameServer)

	minwaittime := int32(cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "gamematchwaitmintime"))
	maxwaittime := int32(cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "gamematchwaitmaxtime"))
	intervalGameQueue.scope = util.NewIntScope(minwaittime, maxwaittime)
	intervalGameQueue.stopSig = make(chan bool, 1)

	return intervalGameQueue
}

// Start 開始
func (gq *IntervalMatchGameQueue) Start() {
	gq.ServerGameQueueBase.Start()

	go func() {
		defer util.PanicErrStack()

		example.invalid("定時撮合邏輯線程啟動")
		for {
			sleeptime := gq.scope.GetRandomValue()
			time.Sleep(time.Duration(sleeptime) * time.Second)

			select {
			case <-gq.stopSig:
				mlog.Error("定時撮合邏輯線程退出")
				break
			default:
			}

			gq.ServerGameQueueBase.publishGameMatch()
		}
	}()
}

// Stop 停止
func (gq *IntervalMatchGameQueue) Stop() {
	gq.stopSig <- true
}
