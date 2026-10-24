package scorestorage

import "example.invalid/ServerTeam/Servers/pkg/gamelogic"

// UserScoreStorage 用戶積分存儲接口
type UserScoreStorage interface {
	UpdateUserScore(user *gamelogic.GameUser, score int64, tax int64)       // 更新用戶積分
	OnUserLeave(user *gamelogic.GameUser)                                   // 用戶離開
	OnUserRecover(user *gamelogic.GameUser, dbscore int64)                  // 用戶斷線重連
	OnUserMoneyAdd(user *gamelogic.GameUser, changeScore int64)             // 用戶遊戲中到帳
	OnUserReady(user *gamelogic.GameUser)                                   // 用戶準備遊戲
	UpdateExperienceScore(user *gamelogic.GameUser, score int64, tax int64) // 更新體驗金幣
	UpdateGameJackpotPool(user *gamelogic.GameUser, amount int64)           // 更新遊戲獎金池
}
