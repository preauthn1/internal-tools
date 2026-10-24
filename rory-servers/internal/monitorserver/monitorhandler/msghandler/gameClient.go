package msghandler

import (
	"errors"
	"fmt"
	"net"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/clientframework"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// GameClient ?
type GameClient struct {
	*ClientBase
}

// NewGameClient ?
func NewGameClient(addr string) *GameClient {
	gc := new(GameClient)
	gc.ClientBase = new(ClientBase)
	gc.RegisterMsgHandler()
	gc.StopSig = make(chan bool, 1)
	var err error
	gc.TCPAddr, err = net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		mlog.Error("解析tcpaddr %s失敗", addr)
	}

	gc.Timeout = 10
	return gc
}

// RequestGameVer ?
func (gc *GameClient) RequestGameVer() error {
	msgparser := gc.MsgParser
	msgparser.RegisterParser(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_RequestGameVerID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_GameVerInfoID), msg.NewProtoMsgParser())

	sig := make(chan bool, 1)
	handler := clientframework.NewSimpleMsgHandler(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_GameVerInfoID), func(m *msg.Message) {
		//解析MsgData詳細格式與內容
		gameVer := &netproto.GameVerInfo{}
		byteMsgData := m.MsgData.([]byte) //須將interface轉[]byte
		proto.Unmarshal(byteMsgData, gameVer)
		//mlog.System("============%v", gameVer)
		if gameVer.GetAndroidVersion() == "" || gameVer.GetIOSVersion() == "" {
			sig <- false
		} else {
			sig <- true
		}
	})
	gc.MsgHandler.AddHandler(handler)
	defer func() {
		gc.MsgHandler.RemoveSimpleHandler(handler)
	}()
	msg := msg.NewMessage(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_RequestGameVerID), nil)
	if err := gc.SendMessage(msg); err != nil {
		return err
	}
	select {
	case b := <-sig:
		if b {
			mlog.Debug("收到遊戲版本消息")
		} else {
			errmsg := fmt.Sprint("Request GameVersion fail")
			mlog.Error(errmsg)
			return errors.New(errmsg)
		}
	case <-time.After(time.Duration(30) * time.Second):
		errmsg := fmt.Sprintf("Request GameVersion timeout, %s", gc.TCPAddr.String())
		mlog.Error(errmsg)
		return errors.New(errmsg)
	}
	return nil
}

// LoginGameID 登入遊戲
func (gc *GameClient) LoginGameID(loginRet *netproto.UserLoginRet, enable int) error {
	now := time.Now().UTC()
	gc.EnableHeartBeat = false

	sig := make(chan *netproto.LoginGameRoomRet, 1)
	_, sig, err := gc.userLoginRetHandler(sig, loginRet)
	if err != nil {
		mlog.Error("發送登入遊戲消息錯誤, %s", err.Error())
		return err
	}

	defer func() {
		//gc.MsgHandler.RemoveSimpleHandler(handler)
	}()

	select {
	case val := <-sig:
		errmsg := fmt.Sprintf("收到登入 %v 成功 , 花費時間[%v]", val.RoomInfo.GetGameName(), time.Since(now))
		mlog.System(errmsg)
		if enable == 1 {
			return errors.New(errmsg)
		}
		select {}
	case <-time.After(time.Duration(30) * time.Second):
		errmsg := fmt.Sprint("Request 登入遊戲消息 timeout")
		mlog.Error(errmsg)
		return errors.New(errmsg)
	}

	return nil
}

func (gc *GameClient) userLoginRetHandler(sig chan *netproto.LoginGameRoomRet, loginRet *netproto.UserLoginRet) (*clientframework.SimpleMsgHandler, chan *netproto.LoginGameRoomRet, error) {
	handler := clientframework.NewSimpleMsgHandler(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_LoginRoomRetID), func(m *msg.Message) {
		loginRoom := &netproto.LoginGameRoomRet{}
		byteMsgData := m.MsgData.([]byte)
		proto.Unmarshal(byteMsgData, loginRoom)
		gc.LoginGameRoomRet = loginRoom
		sig <- loginRoom
	})
	gc.MsgHandler.AddHandler(handler)

	msgData := new(netproto.LoginGameRoomInfo)
	msgData.Cer = proto.String(loginRet.GetCer())
	msgData.UserID = proto.Int32(loginRet.GetUserID())
	msgData.HDCode = proto.String(loginRet.GetHDCode())
	msgData.HDType = proto.Int32(loginRet.GetHDType())
	msg1 := msg.NewMessage(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_LoginRoomID), msgData)
	if err := gc.SendMessage(msg1); err != nil {
		mlog.Error("發送登入遊戲消息錯誤, %s", err.Error())
		sig <- nil
		return nil, sig, err
	}

	return handler, sig, nil
}

// ByLeaveTipHandler ?
func (gc *GameClient) ByLeaveTipHandler() {
	handler := clientframework.NewSimpleMsgHandler(int32(netproto.MessageBClassID_Game), int32(netproto.BY_GameMessageClassID_BYLeaveTip), func(m *msg.Message) {
		loginRoom := &netproto.BY_LeaveMsg{}
		byteMsgData := m.MsgData.([]byte)
		proto.Unmarshal(byteMsgData, loginRoom)
		mlog.System("%v", loginRoom.GetTipMessage())
	})
	gc.MsgHandler.AddHandler(handler)
}
