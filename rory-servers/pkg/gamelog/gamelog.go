package gamelog

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 遊戲動作時間類型
const (
	GameStepTimeTypeSimpleTime int32 = 0 // 精簡時間
	GameStepTimeTypeFullTime   int32 = 1 // 完整時間
)

// GameStep 遊戲動作
type GameStep struct {
	StepTime string // 動作發生時間
	StepDesc string // 動作描述
}

// UserData 玩家資訊
type UserData struct {
	TransactionNo string // 交易序號
	Score         int64  // 當前積分
	TradeID       int64  // 每次進入遊戲房間會獲取一個唯一的值
	UpdateScore   int64  // 單局更新分數
	BetScore      int64  // 下注分數
	PayoutScore   int64  // 返還分數
	ValidScore    int64  // 流水
	TaxAmount     int64  // 稅額
}

// NewGameStep 創建遊戲動作對象
func NewGameStep(stepdesc string, timetype int32) *GameStep {
	step := new(GameStep)
	step.StepDesc = stepdesc
	if timetype == GameStepTimeTypeSimpleTime {
		step.StepTime = time.Now().Format("15:04:05Z07:00")
	} else {
		step.StepTime = util.FormatTimeStr(time.Now())
	}

	return step
}

// GameLog 牌局記錄
type GameLog interface {
	// 當下牌局context
	Context() *RoundContext

	// 開始遊戲記錄
	BeginGameLog() *RoundContext

	// 添加遊戲步驟
	AddGameStep(step *GameStep)

	// 遊戲步驟結束
	EndGameLog()

	// 配置指定開獎
	SetNextJiangRecord(record *netproto.NextJiangRecord)
}
