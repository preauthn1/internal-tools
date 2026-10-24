package mahjong

import "example.invalid/ServerTeam/Servers/pkg/gamelogic"

// UserPoker 用戶牌型
type UserPoker struct {
	User  *gamelogic.GameUser // 用戶
	Cards *CardGroup          // 牌型
}

// NewUserPoker 創建用戶牌型對象
func NewUserPoker(user *gamelogic.GameUser, cards *CardGroup) *UserPoker {
	up := new(UserPoker)
	up.User = user
	up.Cards = cards

	return up
}
