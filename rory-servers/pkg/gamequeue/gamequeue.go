package gamequeue

import "example.invalid/ServerTeam/Servers/pkg/gamelogic"

// GameQueue 遊戲撮合接口
type GameQueue interface {

	// 用戶排隊
	UserQueue(user *gamelogic.GameUser)

	// 用戶遊戲結束
	UserGameEnd(user *gamelogic.GameUser)

	// 開始運行
	Start()

	// 停止運行
	Stop()
}
