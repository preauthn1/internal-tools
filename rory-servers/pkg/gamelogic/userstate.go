package gamelogic

// UserState 用戶狀態
type UserState interface {
	// 用戶站立
	OnUserFree(user *GameUser, room *Room) bool

	// 用戶坐下
	OnUserSit(user *GameUser, tableno int32, seatno int32, room *Room) bool

	// 用戶排隊
	OnUserQueue(user *GameUser, room *Room) bool

	// 用戶準備
	OnUserReady(user *GameUser, room *Room) bool

	// 用戶開始遊戲
	OnUserGame(user *GameUser, room *Room) bool

	// 用戶遊戲結束
	OnGameEnd(user *GameUser, room *Room) bool

	// 用戶離開
	OnUserLeave(user *GameUser, room *Room) bool

	// 用戶斷線
	OnUserOffLine(user *GameUser, room *Room) bool

	// 狀態名稱
	StateName() string

	// 狀態值
	StateValue() int32
}

// StateName
const (
	StateNameFree    = "free"
	StateNameSit     = "sit"
	StateNameQueue   = "queue"
	StateNameReady   = "ready"
	StateNameGame    = "game"
	StateNameOffline = "offline"
)

// 用戶狀態
const (
	UserStateFree    = 0 // 用戶站立狀態
	UserStateQueue   = 1 // 用戶排隊狀態
	UserStateSit     = 2 // 用戶坐下狀態
	UserStateReady   = 3 // 用戶舉手狀態
	UserStateGame    = 4 // 用戶遊戲狀態
	UserStateOffline = 5 // 用戶斷線狀態
)

var freeState UserState
var sitState UserState
var queueState UserState
var readyState UserState
var gameState UserState
var offlineState UserState

func init() {
	freeState = new(UserFreeState)
	queueState = new(UserQueueState)
	sitState = new(UserSitState)
	readyState = new(UserReadyState)
	gameState = new(UserGameState)
	offlineState = new(UserOfflineState)
}

// UserStateNameByValue 取得狀態名稱
func UserStateNameByValue(v int) string {
	switch v {
	default:
		return "?"
	case UserStateFree:
		return freeState.StateName()
	case UserStateSit:
		return sitState.StateName()
	case UserStateQueue:
		return queueState.StateName()
	case UserStateReady:
		return readyState.StateName()
	case UserStateGame:
		return gameState.StateName()
	case UserStateOffline:
		return offlineState.StateName()
	}
}
