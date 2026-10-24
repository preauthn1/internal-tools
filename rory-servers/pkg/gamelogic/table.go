package gamelogic

import (
	"bytes"
	"fmt"
	"reflect"
)

// Table 桌子
type Table struct {
	Room      *Room
	TableNo   int32     // 桌子號
	Seats     []*Seat   // 座位
	GameLogic GameLogic // 遊戲邏輯
}

// NewTable 創建桌子對象
func NewTable(seatcount int32, tno int32) *Table {
	t := new(Table)
	t.TableNo = tno

	var i int32
	for i = 0; i < seatcount; i++ {
		st := NewSeat(i)
		t.Seats = append(t.Seats, st)
	}

	return t
}

// AddNewSeat 新增座位
func (t *Table) AddNewSeat(iSeatNo int) {
	if iSeatNo < len(t.Seats) {
		return
	}

	pSeat := NewSeat(int32(iSeatNo))
	t.Seats = append(t.Seats, pSeat)
}

// GetSeatByNo 根據座位號獲取座位
func (t *Table) GetSeatByNo(seatno int32) *Seat {
	for _, s := range t.Seats {
		if s.SeatNo == seatno {
			return s
		}
	}

	return nil
}

// GetTableUser 獲取桌子上的用戶
func (t *Table) GetTableUser() []*GameUser {
	var urs []*GameUser
	for _, st := range t.Seats {
		if st.SeatUser != nil {
			urs = append(urs, st.SeatUser)
		}
	}

	return urs
}

// GetTableRealUser 獲取桌子上的真實用戶
func (t *Table) GetTableRealUser() []*GameUser {
	var urs []*GameUser
	for _, st := range t.Seats {
		if st.SeatUser != nil && !st.SeatUser.IsRobot() {
			urs = append(urs, st.SeatUser)
		}
	}

	return urs
}

// GetTableRealRobot 獲取桌子上的真實用戶
func (t *Table) GetTableRealRobot() []*GameUser {
	var urs []*GameUser
	for _, st := range t.Seats {
		if st.SeatUser != nil && st.SeatUser.IsRobot() {
			urs = append(urs, st.SeatUser)
		}
	}

	return urs
}

// TableUserCount 桌子上的用戶數
func (t *Table) TableUserCount() int {
	return len(t.GetTableUser())
}

// TableRealUserCount 桌子上的用戶數
func (t *Table) TableRealUserCount() int {
	return len(t.GetTableRealUser())
}

// TableRealRobotCount 桌子上的robot用戶數
func (t *Table) TableRealRobotCount() int {
	return len(t.GetTableRealRobot())
}

// ReadyUser 獲取桌子上舉手的用戶
func (t *Table) ReadyUser() []*GameUser {
	var urs []*GameUser
	for _, st := range t.Seats {
		if st.SeatUser != nil && st.SeatUser.IsReady() {
			urs = append(urs, st.SeatUser)
		}
	}

	return urs
}

// GetFreeSeat 獲取一個空閒的座位
func (t *Table) GetFreeSeat() *Seat {
	for _, st := range t.Seats {
		if st.SeatUser == nil {
			return st
		}
	}

	return nil
}

// GetBYFreeSeat 獲取BY一個空閒的座位
func (t *Table) GetBYFreeSeat() *Seat {
	for _, st := range t.Seats {
		if st.SeatNo <= 1 {
			continue
		}
		if st.SeatUser == nil {
			return st
		}
	}

	return nil
}

// GetUserByID 根據用戶id獲取用戶
func (t *Table) GetUserByID(userid int32) *GameUser {
	for _, st := range t.Seats {
		if st.SeatUser != nil && st.SeatUser.UserID == userid {
			return st.SeatUser
		}
	}

	return nil
}

func (t *Table) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("tableno:%d, logic: %v ", t.TableNo, reflect.TypeOf(t.GameLogic)))
	for _, st := range t.Seats {
		buf.WriteString(st.String() + "\r\n")
	}

	return buf.String()
}
