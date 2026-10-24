package gamequeue

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// UserReadyQueue 玩家自選配桌邏輯
type UserReadyQueue struct {
	ActCtrl gamelogic.ActionController
	Configs *config.Group
}

// NewUserReadyQueue 建立玩家自選配桌邏輯
func NewUserReadyQueue(actionController gamelogic.ActionController, cfgs *config.Group) *UserReadyQueue {
	q := new(UserReadyQueue)
	q.ActCtrl = actionController
	q.Configs = cfgs
	return q
}

// UserQueue 用戶排隊
func (q *UserReadyQueue) UserQueue(user *gamelogic.GameUser) {
	example.invalid("玩家[%d] 排隊", user.UserID)

	t := q.ActCtrl.GetFreeTable()
	if t == nil {
		example.invalid("沒有空閒牌桌")
		return
	}

	users := q.ActCtrl.GetQueueUser()

	if len(users) < len(t.Seats) {
		example.invalid("目前排隊人數不足，需 %d 人才可開始", len(t.Seats))
		return
	}

	for i := 0; i < len(t.Seats); i++ {
		q.ActCtrl.UserSit(users[i], t.TableNo, t.Seats[i].SeatNo)
	}
}

// UserGameEnd 用戶遊戲結束
func (q *UserReadyQueue) UserGameEnd(user *gamelogic.GameUser) {
	q.ActCtrl.UserFree(user)
}

// Start 開始運行
func (q *UserReadyQueue) Start() {}

// Stop 停止運行
func (q *UserReadyQueue) Stop() {}
