package gamequeue

import "example.invalid/ServerTeam/Servers/pkg/event"

// ?
const (
	EventNameGameMatch = "gamequeue_GameMatch"  // 事件名稱,用戶匹配
	ArgsNameMatchUsers = "gamequeue_matchusers" // 事件參數,匹配的用戶列表
)

// NewGameMatchEvent 創建用戶撮合事件
func NewGameMatchEvent() *event.Event {
	args := event.NewEventArgs()

	evt := event.NewEvent(EventNameGameMatch, args)

	return evt
}
