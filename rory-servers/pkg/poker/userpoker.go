package poker

import (
	"bytes"
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// UserPoker 用户牌型
type UserPoker struct {
	User  *gamelogic.GameUser // 用户
	Cards *CardGroup          // 牌型
}

// NewUserPoker 創建用戶牌型對象
func NewUserPoker(user *gamelogic.GameUser, cards *CardGroup) *UserPoker {
	up := new(UserPoker)
	up.User = user
	up.Cards = cards

	return up
}

// NewUserEmptyPoker ?
func NewUserEmptyPoker(user *gamelogic.GameUser) *UserPoker {
	up := new(UserPoker)
	up.User = user
	up.Cards = new(CardGroup)

	return up
}

func (up *UserPoker) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%v.", up.User))

	for _, c := range up.Cards.Cards {
		buf.WriteString(fmt.Sprintf("%s,", c))
	}

	return buf.String()
}

// GetNetMessageData 獲取用戶牌型網絡消息數據
func (up *UserPoker) GetNetMessageData() *netproto.UserPoker {
	pdata := new(netproto.UserPoker)
	pdata.UserID = proto.Int32(up.User.UserID)
	pdata.SeatNo = proto.Int32(up.User.SeatNo)

	for _, c := range up.Cards.Cards {
		pdata.Cards = append(pdata.Cards, c.GetNetMessageData())
	}

	return pdata
}

// GamePoker 牌局牌型
type GamePoker struct {
	UserPokers []*UserPoker
}

// NewGamePoker 創建遊戲牌型
func NewGamePoker() *GamePoker {
	gpk := new(GamePoker)
	return gpk
}

func (gp *GamePoker) String() string {
	var buf bytes.Buffer
	for _, p := range gp.UserPokers {
		buf.WriteString(fmt.Sprintf("%s\r\n", p))
	}

	return buf.String()
}

// AddUserPoker 添加用戶牌
func (gp *GamePoker) AddUserPoker(upk *UserPoker) {
	for _, p := range gp.UserPokers {
		if p.User.UserID == upk.User.UserID {
			p.Cards.AddCardGroup(upk.Cards)

			return
		}
	}

	gp.UserPokers = append(gp.UserPokers, upk)
}

// GetUserPoker 獲取用戶牌型
func (gp *GamePoker) GetUserPoker(userid int32) *UserPoker {
	for _, upk := range gp.UserPokers {
		if upk.User.UserID == userid {
			return upk
		}
	}

	return nil
}

// GetNetMessageData 獲取遊戲牌型網絡消息數據
func (gp GamePoker) GetNetMessageData() *netproto.GamePoker {
	ngpk := new(netproto.GamePoker)
	for _, upk := range gp.UserPokers {
		ngpk.Poker = append(ngpk.Poker, upk.GetNetMessageData())
	}

	return ngpk
}
