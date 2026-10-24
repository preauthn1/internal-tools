package infc

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

// Player 玩家
type Player struct {
	user *gamelogic.GameUser

	UID   int32
	Score int64
	Miss  int64
	State State

	LastTime time.Time
}

// NewPlayer Create Player
func NewPlayer(u *gamelogic.GameUser) *Player {
	p := new(Player)
	p.user = u
	p.UID = u.UserID
	p.State = StateOffline
	p.LastTime = time.Now()

	return p
}
