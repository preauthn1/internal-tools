package scorestorage

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// ConsoleUserScoreStorage Console 顯示
type ConsoleUserScoreStorage struct{}

// UpdateUserScore 寫分 : 用戶資料, 加減分數, 稅額
func (uss *ConsoleUserScoreStorage) UpdateUserScore(user *gamelogic.GameUser, score int64, tax int64) {
	user.Score += score
	user.UpdateScore += score
	user.TaxAmount += tax
	mlog.Debug("[玩家寫分] 最終金幣: %d, 更新金幣: %d, 稅金: %d", user.Score, score, tax)
}

// OnUserLeave 退房
func (uss *ConsoleUserScoreStorage) OnUserLeave(user *gamelogic.GameUser) {
	mlog.Debug("[玩家退房] 最終金幣: %d, 更新金幣: %d, 稅金: %d", user.Score, user.UpdateScore, user.TaxAmount)
}

// OnUserRecover 斷線重連時會呼叫此函式記錄斷線重連時的分數異動
func (uss *ConsoleUserScoreStorage) OnUserRecover(user *gamelogic.GameUser, dbscore int64) {
	//fmt.Printf("user score recover:%v, %d\r\n", user, dbscore)
}

// OnUserMoneyAdd 即時到帳, 有可能一局中要Add兩次以上
func (uss *ConsoleUserScoreStorage) OnUserMoneyAdd(user *gamelogic.GameUser, changeScore int64) {
	user.InGamingMoneyChange += changeScore
}

// OnUserReady 在遊戲結束時，每個用戶都會做的事
func (uss *ConsoleUserScoreStorage) OnUserReady(user *gamelogic.GameUser) {
	uss.beforeGameStart(user)
}

// UpdateExperienceScore 寫分 : 用戶資料, 加減分數, 稅額
func (uss *ConsoleUserScoreStorage) UpdateExperienceScore(user *gamelogic.GameUser, score int64, tax int64) {
}

// UpdateGameJackpotPool 遊戲獎金池寫分
func (uss *ConsoleUserScoreStorage) UpdateGameJackpotPool(user *gamelogic.GameUser, amount int64) {}

// 正在遊戲中有需要異動金額的部份等到下局開始時才生效
func (uss *ConsoleUserScoreStorage) beforeGameStart(user *gamelogic.GameUser) {
	if user.InGamingMoneyChange != 0 {
		user.Score += user.InGamingMoneyChange
		user.LoginScore += user.InGamingMoneyChange
		user.InGamingMoneyChange = 0
		if user.IsOffLine() {
			return
		}

	}
}
