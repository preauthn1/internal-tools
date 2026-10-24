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

// HallClient ?
type HallClient struct {
	*ClientBase
}

// NewHallClient ?
func NewHallClient(addr string) *HallClient {
	hc := new(HallClient)
	hc.ClientBase = new(ClientBase)
	hc.RegisterMsgHandler()
	hc.StopSig = make(chan bool, 1)
	var err error
	hc.TCPAddr, err = net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		mlog.Error("解析tcpaddr %s失敗", addr)
	}

	hc.Timeout = 10
	return hc
}

// RegisterMsgHandler 註冊所有消息handler
func (hc *HallClient) RegisterMsgHandler() {
	hc.MsgHandler = clientframework.NewHookMsgHandlerGroup()
	mh := hc.MsgHandler
	mh.AddHandler(NewHeartBeatConfigHandler())
	mh.AddHandler(NewTipMessageHandler())
}

// RequestHallVersion ?
func (hc *HallClient) RequestHallVersion() error {
	sig := make(chan bool, 1)
	handler := clientframework.NewSimpleMsgHandler(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_HallVerInfoID), func(m *msg.Message) {
		hallVerData := &netproto.HallVerInfo{}
		byteMsgData := m.MsgData.([]byte) //須將interface轉[]byte
		proto.Unmarshal(byteMsgData, hallVerData)
		//mlog.System("============%v", hallVerData)
		if hallVerData.GetAndroidVersion() == "" || hallVerData.GetIOSVersion() == "" {
			sig <- false
		} else {
			sig <- true
		}
	})
	hc.MsgHandler.AddHandler(handler)
	defer func() {
		hc.MsgHandler.RemoveSimpleHandler(handler)
	}()

	msg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestHallVerID), nil)
	if err := hc.SendMessage(msg); err != nil {
		return err
	}
	select {
	case b := <-sig:
		if b {
			//mlog.System("收到大廳版本消息")
		} else {
			errmsg := fmt.Sprint("Request HallVersion fail")
			mlog.Error(errmsg)
			return errors.New(errmsg)
		}
	case <-time.After(time.Duration(30) * time.Second):
		errmsg := fmt.Sprint("Request HallVersion timeout")
		mlog.Error(errmsg)
		return errors.New(errmsg)
	}
	return nil
}

// RequestGuestLogon ?
func (hc *HallClient) RequestGuestLogon(siteID int32, guestID int32) (*netproto.UserLoginRet, error) {
	sig := make(chan *netproto.UserLoginRet, 1)
	handler := clientframework.NewSimpleMsgHandler(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_LoginRetID), func(m *msg.Message) {
		userLoginRet := new(netproto.UserLoginRet)
		byteMsgData := m.MsgData.([]byte) //須將interface轉[]byte
		proto.Unmarshal(byteMsgData, userLoginRet)
		mlog.System("登入大廳 userid[%+v] CashAmount[%+v] BankAmount[%v]", userLoginRet.GetUserID(),
			userLoginRet.GetUserData().GetCashAmount(), userLoginRet.GetUserData().GetBankAmount())
		sig <- userLoginRet

	})
	hc.MsgHandler.AddHandler(handler)
	defer func() {
		hc.MsgHandler.RemoveSimpleHandler(handler)
	}()
	//fmt.Println(fmt.Sprintf("%d", guestId))
	//Client傳入資料的協議號
	msgData := new(netproto.GuestLogin)
	msgData.Version = proto.String("0.0.0")
	msgData.HDType = proto.Int32(1)
	msgData.HDCode = proto.String(fmt.Sprintf("%d", guestID))
	msgData.SiteID = proto.Int32(siteID)
	msgData.PlatformID = proto.Int32(0)
	msgData.ServerID = proto.Int32(1)
	msg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GuestLoginID), msgData)
	if err := hc.SendMessage(msg); err != nil {
		mlog.Error("發送登入消息錯誤, 錯誤信息:%s", err.Error())
		return nil, err
	}

	select {
	case userLoginRet := <-sig:
		if userLoginRet != nil {
			if userLoginRet.GetCode() != int32(1) {
				mlog.Error("收到登入失敗消息%v userid[%v] Code[%v] hdcode[%v] cash[%v] siteID[%v]", userLoginRet.GetMessage(), userLoginRet.GetUserID(), userLoginRet.GetCode(),
					userLoginRet.GetHDCode(), userLoginRet.GetUserData().GetCashAmount(), siteID)
			}
			return userLoginRet, nil
		}
		errmsg := fmt.Sprint("Request GuestLogon fail")
		mlog.Error(errmsg)
		return nil, errors.New(errmsg)
	case <-time.After(time.Duration(30) * time.Second):
		errmsg := fmt.Sprint("Request GuestLogon timeout")
		return nil, errors.New(errmsg)
	}
}

// RequestServerList ?
func (hc *HallClient) RequestServerList() (*netproto.AllGameServerInfo, error) {
	sig := make(chan *netproto.AllGameServerInfo, 1)
	now := time.Now().UTC()
	handler := clientframework.NewSimpleMsgHandler(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ServerListDataID), func(m *msg.Message) {
		serverInfo := new(netproto.AllGameServerInfo)
		byteMsgData := m.MsgData.([]byte) //須將interface轉[]byte
		proto.Unmarshal(byteMsgData, serverInfo)
		sig <- serverInfo
	})
	hc.MsgHandler.AddHandler(handler)
	defer func() {
		hc.MsgHandler.RemoveSimpleHandler(handler)
	}()

	//Client傳入資料的協議號
	msg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestServerListID), nil)
	if err := hc.SendMessage(msg); err != nil {
		mlog.Error("發送獲取服務器列表消息錯誤, 錯誤信息:%s", err.Error())
		return nil, err
	}
	select {
	case data := <-sig:
		if data != nil {
			//mlog.System("收到ServerList消息")
			mlog.System("serverlist來回時間:%v", time.Since(now))
			return data, nil
		}
		errmsg := fmt.Sprint("Request ServerList fail")
		mlog.Error(errmsg)
		return nil, errors.New(errmsg)
	case <-time.After(time.Duration(30) * time.Second):
		errmsg := fmt.Sprint("Request ServerList timeout")
		mlog.Error(errmsg)
		return nil, errors.New(errmsg)
	}
}

// RequestServerAddr ?
func (hc *HallClient) RequestServerAddr(gameid int32, flag string) (*netproto.GameServerAddr, error) {
	sig := make(chan *netproto.GameServerAddr, 1)
	handler := clientframework.NewSimpleMsgHandler(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GameServerAddrInfoID), func(m *msg.Message) {
		gameServerAddr := new(netproto.GameServerAddr)
		byteMsgData := m.MsgData.([]byte) //須將interface轉[]byte
		proto.Unmarshal(byteMsgData, gameServerAddr)
		sig <- gameServerAddr
	})
	hc.MsgHandler.AddHandler(handler)
	defer func() {
		hc.MsgHandler.RemoveSimpleHandler(handler)
	}()

	// Client傳入資料的協議號
	msgData := new(netproto.RequestGameServerAddrInfo)
	msgData.ServerFlag = proto.String(flag)
	msgData.GameID = proto.Int32(gameid) //炸金花
	msg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestGameServerAddrID), msgData)
	if err := hc.SendMessage(msg); err != nil {
		mlog.Error("發送登入消息錯誤, 錯誤信息:%s", err.Error())
		return nil, err
	}

	select {
	case data := <-sig:
		if data != nil {
			//mlog.System("收到GameServerAddr消息")
			return data, nil
		}
		errmsg := fmt.Sprint("Request GameServerAddr fail")
		mlog.Error(errmsg)
		return nil, errors.New(errmsg)
	case <-time.After(time.Duration(30) * time.Second):
		errmsg := fmt.Sprint("Request GameServerAddr timeout")
		mlog.Error(errmsg)
		return nil, errors.New(errmsg)
	}
}
