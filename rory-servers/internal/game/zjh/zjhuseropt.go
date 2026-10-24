package zjh

import (
	proto "example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

import (
	"fmt"
)

type ZjhUserData struct {
	userCards         *poker.CardGroup    //用户手牌
	User              *gamelogic.GameUser //用户
	compareUserSeatNo []int32             //这些人和我比过牌
	win               int64               //用户赢钱
	userMoney         int64               //用户钱币拷贝
	isCheck           bool                //是否看牌
	isDrop            bool                //是否弃牌
	isFaild           bool                //是否比牌失败
	isFill            bool                //是否能跟牌
	isBetAdd          bool                //是否能加注
	isFollowEnd       bool                //是否跟到底
	isOpt             bool                //本轮是否操作过
	isClear           bool                //是否己結算
	isAllin           bool                //當局是否在全押狀態
	AllInUserId       int32               //全押玩家id
}

//创建一个用户数据
func NewZjhUserData(user *gamelogic.GameUser) *ZjhUserData {
	userdata := new(ZjhUserData)
	userdata.userMoney = user.Score
	userdata.isCheck = false
	userdata.isDrop = false
	userdata.isFaild = false
	userdata.isFill = true
	userdata.isBetAdd = true
	userdata.User = user
	userdata.userCards = nil
	userdata.isFollowEnd = false
	userdata.isOpt = false
	userdata.isClear = false
	userdata.isAllin = false

	return userdata
}

//用户是否可操作
func (ud *ZjhUserData) canOpt() bool {
	if ud.isFaild || ud.isDrop {
		return false
	}

	return true
}

//用戶可提前換桌
func (ud *ZjhUserData) canEarlyChange() bool {
	if ud.isFaild || ud.isDrop {
		return true
	}

	return false
}

//用戶可提前換桌
func (ud *ZjhUserData) String() string {
	return fmt.Sprintf("userid %d, tableno %d, seatno %d", ud.User.UserID, ud.User.TableNo, ud.User.SeatNo)
}

//获取用户的操作消息
func (ud *ZjhUserData) GetOptMsg(betadd []int32, time int32, compare bool, bFirst bool, bCheck bool, betlimit int64, ballin bool) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_WaitOptID)

	//已经看牌了
	var odd int32 = 1
	if ud.isCheck {
		odd = 2
	}

	//玩家的钱超过比牌所需钱上限了不能加注也不能跟注
	if ud.userMoney < betlimit {
		ud.isFill = false
		ud.isBetAdd = false
		ballin = false
	} else {
		ud.isFill = true
		ud.isBetAdd = true
	}
	if ud.isAllin {
		ud.isFill = true
		ud.isBetAdd = false
		ud.isFollowEnd = false
		ballin = false
		compare = false
		bFirst = false
	}

	arrLen := len(betadd)
	if arrLen <= 0 {
		ud.isBetAdd = false
	}

	check := bCheck
	if bCheck {
		check = !ud.isCheck
	}

	pdata := new(netproto.ZJH_WaitOpt)
	pdata.Compare = proto.Bool(compare)
	pdata.Check = proto.Bool(check)
	pdata.Fill = proto.Bool(ud.isFill)
	pdata.Drop = proto.Bool(!ud.isDrop)
	pdata.First = proto.Bool(bFirst)
	pdata.Betadd = proto.Bool(ud.isBetAdd)
	pdata.OptTime = proto.Int32(time)
	pdata.Allin = proto.Bool(ballin)

	for i := 0; i < arrLen; i++ {
		pdata.Betaddvalue = append(pdata.Betaddvalue, betadd[i]*odd)
	}

	pmsg.MsgData = pdata
	return pmsg
}

//判断用户是否能跟注
func (ud *ZjhUserData) userCanFollow(betlimit int64) bool {
	if ud.userMoney < betlimit {
		return false
	}

	return true
}

//设置用户的钱
func (ud *ZjhUserData) SetMoney(money int64) {
	ud.userMoney += money
}

//设置用户弃牌
func (ud *ZjhUserData) SetDrop() {
	ud.isDrop = true
}

//设置和玩家比牌的座位号
func (ud *ZjhUserData) setCompareUserSeatNo(seatno int32) {
	ud.compareUserSeatNo = append(ud.compareUserSeatNo, seatno)
}

//设置玩家跟到底
func (ud *ZjhUserData) setFollowEnd() {
	ud.isFollowEnd = !ud.isFollowEnd
	mlog.Debug("跟到底值为:%v", ud.isFollowEnd)
}

//设置玩家跟到底
func (ud *ZjhUserData) setAllin() {
	ud.isAllin = !ud.isAllin
	ud.AllInUserId = ud.User.UserID
	example.invalid("是否全押:%v", ud.isAllin)
}
