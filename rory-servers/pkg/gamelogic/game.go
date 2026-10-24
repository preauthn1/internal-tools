package gamelogic

// Game 遊戲信息
type Game struct {
	GameID    int32  // 遊戲ID
	GameName  string // 遊戲名稱
	MinPlayer int32  // 最小玩家數
	MaxPlayer int32  // 最大玩家數
}

// NewGame 創建遊戲對象
func NewGame() *Game {
	g := new(Game)

	return g
}
