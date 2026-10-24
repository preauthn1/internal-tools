package logic

// ServeState 服務狀態
type ServeState int32

// ServeState
const (
	ServeStateUnlogin ServeState = iota
	ServeStateLogined
	ServeStateWaitBet
	ServeStateWaitResult
)

func (s ServeState) String() string {
	switch s {
	default:
		return "none"
	case ServeStateUnlogin:
		return "Unlogin"
	case ServeStateLogined:
		return "Logined"
	case ServeStateWaitBet:
		return "WaitBet"
	case ServeStateWaitResult:
		return "WaitResult"
	}
}
