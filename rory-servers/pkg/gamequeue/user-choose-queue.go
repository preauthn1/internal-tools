package gamequeue

import "example.invalid/ServerTeam/Servers/pkg/gamelogic"

// UserChooseGameQueue 玩家自選配桌邏輯
type UserChooseGameQueue struct{}

// NewUserChooseGameQueue 建立玩家自選配桌邏輯
func NewUserChooseGameQueue() *UserChooseGameQueue {
	q := new(UserChooseGameQueue)
	return q
}

// UserQueue 用戶排隊
func (q *UserChooseGameQueue) UserQueue(user *gamelogic.GameUser) {}

// UserGameEnd 用戶遊戲結束
func (q *UserChooseGameQueue) UserGameEnd(user *gamelogic.GameUser) {}

// Start 開始運行
func (q *UserChooseGameQueue) Start() {}

// Stop 停止運行
func (q *UserChooseGameQueue) Stop() {}
