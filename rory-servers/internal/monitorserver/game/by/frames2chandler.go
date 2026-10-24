package by

import (
	"math/rand"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/clientframework"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// FrameS2CHandler ?
type FrameS2CHandler struct {
	userBoom []*UserBoomInfo
}

// UserBoomInfo ?
type UserBoomInfo struct {
	userID int32
	boomID int32
}

// NewFrameS2CHandler ?
func NewFrameS2CHandler() *FrameS2CHandler {
	f := new(FrameS2CHandler)
	f.userBoom = make([]*UserBoomInfo, 0)
	return f
}

// IsHook ?
func (by *FrameS2CHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.BY_GameMessageClassID_FrameS2CID)
}

// OnMessage ?
func (by *FrameS2CHandler) OnMessage(client *clientframework.CtxBase, m *msg.Message) {
	frame := &netproto.BY_FrameS2C{}
	byteMsgData := m.MsgData.([]byte)
	proto.Unmarshal(byteMsgData, frame)
	r := util.RndInt(0, 50)
	userid := int32(0)
	//if client.ClientBase.LoginGameRoomRet != nil {
	//	userid = client.ClientBase.LoginGameRoomRet.GetUserData().GetUserID()
	//}

	if r >= 0 && len(frame.GetNewFish()) > int(0) {
		//
		//	if len(frame.GetUserInfo()) == 0 {
		//		return
		//	}
		for _, u := range frame.GetUserInfo() {
			if u.GetUserID() == userid {
				userBoom := by.CheckGameUser(u.GetUserID())
				for i := 1; i <= 4; i++ {
					shootMsg := by.getShootMsg(userBoom, int32(i), frame.GetNewFish()[0].GetFishID())
					client.ConnClt.SendMessage(shootMsg)
				}
			}

		}
	}

}

// CheckGameUser ?
func (by *FrameS2CHandler) CheckGameUser(userid int32) *UserBoomInfo {
	for _, v := range by.userBoom {
		if v.userID == userid {
			return v
		}
	}

	newUser := new(UserBoomInfo)
	newUser.userID = userid
	newUser.boomID = 10000000
	by.userBoom = append(by.userBoom, newUser)
	return newUser
}

// SetCallback /
func (by *FrameS2CHandler) SetCallback(cb func(m *msg.Message)) {}

func (by *FrameS2CHandler) getShootMsg(ur *UserBoomInfo, seatno int32, fishid int32) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.BY_GameMessageClassID_FrameC2SID)

	btdata := new(netproto.BY_FrameC2S)
	btdata.FrameID = proto.Int64(0)
	shoot := new(netproto.BY_ShootC2S)
	shoot.Pos = new(netproto.BY_Position)
	suoDing := new(netproto.BY_SuoDingC2S)
	suoDing.SuoDing = proto.Bool(true)
	btdata.SuoDing = suoDing
	userBoom := new(netproto.BY_BoomC2S)

	if ur.boomID == int32(0) {
		ur.boomID = int32(10000000)
	}

	ur.boomID++
	userBoom.BoomID = proto.Int32(ur.boomID)
	userBoom.BoomRate = proto.Int32(1)
	userBoom.FishID = []int32{fishid}
	userBoom.UserID = proto.Int32(ur.userID)
	btdata.Boom = userBoom

	y := float32(1.0)
	//座位2和3在上面需要向下射击
	if seatno == 2 || seatno == 3 {
		y = -1.0
	}

	x := rand.Float32()
	if rand.Int31n(2) == 1 {
		x = -x
	}

	shoot.Pos.X = proto.Float32(x)
	shoot.Pos.Y = proto.Float32(y)
	btdata.Shoot = shoot
	btmsg.MsgData, _ = proto.Marshal(btdata)

	return btmsg
}
