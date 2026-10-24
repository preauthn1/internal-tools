package gamelogic

import (
	"fmt"
)

// Seat 座位
type Seat struct {
	SeatNo   int32     // 座位號
	SeatUser *GameUser // 用戶
}

// NewSeat 創建座位對象
func NewSeat(seatno int32) *Seat {
	st := new(Seat)
	st.SeatNo = seatno
	st.SeatUser = nil

	return st
}

// IsEmpty 沒人坐
func (seat *Seat) IsEmpty() bool {
	return seat.SeatUser == nil
}

func (seat *Seat) String() string {
	return fmt.Sprintf("seatno:%d, user:%v", seat.SeatNo, seat.SeatUser)
}
