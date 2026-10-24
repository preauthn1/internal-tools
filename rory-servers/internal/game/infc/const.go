package infc

// GameMode 遊戲機制
type GameMode uint

// GameModes
const (
	// 玩到死
	GameModeDeathGame GameMode = iota
	// 固定時間
	GameModeFixedTimeGame
)

// GameState 遊戲狀態
type GameState uint

// GameStates
const (
	GameStateWait GameState = iota
	GameStateGaming
	GameStateCalculating
)

// State 玩家狀態
type State uint

// States
const (
	StateOffline State = iota
	StateOnline
	StateGaming
	StateDeath
	StateGiveup
	StateLocked
)
