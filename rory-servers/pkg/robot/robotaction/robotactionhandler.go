package robotaction

import (
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
)

// HandlerCtx 機器人動作處理器運行環境
type HandlerCtx struct {
	Room             *gamelogic.Room
	ActionController gamelogic.ActionController
	RobotManager     *robotcomm.RobotManager
	TfManager        *event.TimerFuncManager
	RobotCfg         *robotcomm.RobotConfig
	RPCMsgQueue      *rpcmsg.Queue
	ScoreStorage     scorestorage.UserScoreStorage
	GameLogics       *gamelogic.LogicManager
	EventChannel     *event.Dispatcher
	ServerStatus     logicmodule.ServerStatusProvider
}

// ?
const (
	MatchTablePoolingOneRobot          = 0 // 配桌方式, 輪詢每張桌子每次放一個機器人
	QueueTypePoolingNRateRobot         = 1 // 配桌方式, 輪詢每張桌子每次放N個機器人
	QueueTypeFixedRobot                = 2 // 2固定配機器人可配置上限
	MatchTablePoolingUserWaitTimeRobot = 3 // 當真人超過某數量時調整機器人配桌時間增加
	MatchTableSaveLossUser             = 4 // 拯救輸太多的玩家, 讓輸的玩家拿最好牌, 其他座位玩家都是總輸贏為正
)

// NewRobotActionHandlerCtx 創建機器人動作處理器環境對象
func NewRobotActionHandlerCtx() *HandlerCtx {
	ctx := new(HandlerCtx)

	return ctx
}

// Handler 機器人動作處理器接口定義
type Handler interface {
	// 啟動動作
	Start(ctx *HandlerCtx)
}

// HandlerGroup 機器人動作處理器組
type HandlerGroup struct {
	handlers []Handler
}

// NewRobotActionHandlerGroup 創建機器人動作處理器組
func NewRobotActionHandlerGroup() *HandlerGroup {
	actg := new(HandlerGroup)

	return actg
}

// AddHandler 添加處理器
func (rahg *HandlerGroup) AddHandler(h Handler) {
	rahg.handlers = append(rahg.handlers, h)
}

// Start 啟動處理器
func (rahg *HandlerGroup) Start(ctx *HandlerCtx) {
	for _, h := range rahg.handlers {
		h.Start(ctx)
	}
}
