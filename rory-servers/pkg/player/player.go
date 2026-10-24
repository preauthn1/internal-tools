package player

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/network"
)

// 設備類型
const (
	UnKnown = iota
	HDTypeAndroid
	HDTypeIOS
	HDTypePC
	HDTypeH5
)

// ?
const (
	GuestUserID  = 0  // 遊客用戶ID
	GuestCerData = "" // 遊客證書數據
)

// Player 遊戲玩家
type Player struct {
	Conn     *network.ConnInfo // 連接信息
	UserID   int32             // 用戶ID
	Cer      string            // 用戶證書
	HDType   int32             // 設備類型
	HDCode   string            // 設備碼
	UserData interface{}       // 用戶數據
	IsOnline bool              // 是否保持連接
}

// NewGuest 創建遊客對象
func NewGuest(conn *network.ConnInfo) *Player {
	player := new(Player)
	player.Conn = conn
	player.UserID = GuestUserID
	player.Cer = GuestCerData
	player.IsOnline = true

	return player
}

// IsGuest 是否遊客
func (p *Player) IsGuest() bool {
	return p.UserID == GuestUserID && p.Cer == GuestCerData
}

func (p *Player) String() string {
	return fmt.Sprintf("userid:%d,cer:%s,hdtype:%d,hdcode:%s,conn:%v,userdata:%v,offline:%t", p.UserID, p.Cer, p.HDType, p.HDCode, p.Conn, p.UserData, p.IsOffLine())
}

// SimpleString ?
func (p *Player) SimpleString() string {
	return fmt.Sprintf("userid:%d,cer:%s,hdtype:%d,hdcode:%s,conn:%v,offline:%v", p.UserID, p.Cer, p.HDType, p.HDCode, p.Conn, p.IsOffLine())
}

// IsOffLine 用戶是否斷線
func (p *Player) IsOffLine() bool {
	return !p.IsOnline
}
