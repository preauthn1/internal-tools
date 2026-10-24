// handler單元測試
package msghandler

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/stretchr/testify/assert"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func TestGameServerStartHandler_ProtoError(t *testing.T) {
	InitTest()

	// 驗證 handler
	serverid := int32(365)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	// 驗證資料
	msgdata := new(netproto.GameServerInfo)
	msgdata.ServerID = proto.Int32(1)
	msgdata.GameID = proto.Int32(1)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData

	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GameServerStartID), rpcMessage)
	handler := new(GameServerStartHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_GameServerStartRetID) {
		t.Error("OnMessage 錯誤")
	}
	data := ret[0].MsgData.(*netproto.GameServerStartRet)

	assert.Equal(t, data.GetCode(), int32(-1))
	assert.Equal(t, data.GetMessage(), "服務器啟動失敗, 解析封包內容失敗")
}

func TestGameServerStartHandler_No_Server(t *testing.T) {
	InitTest()

	// 驗證 handler
	serverid := int32(365)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	// 驗證資料
	msgdata := new(netproto.GameServerStart)
	msgdata.ServerID = proto.Int32(1)
	msgdata.GameID = proto.Int32(1)
	msgdata.ClientAddr = proto.String("0.0.0.0:30101")
	msgdata.HttpAddr = proto.String("")
	msgdata.NotifyAddr = proto.String("0.0.0.0:9000")
	msgdata.LoginMoney = proto.Int32(1000)
	msgdata.ShowName = proto.String("1")
	msgdata.ServerName = proto.String("多福多財_初級場_測試站")

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData

	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GameServerStartID), rpcMessage)
	handler := new(GameServerStartHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_GameServerStartRetID) {
		t.Error("OnMessage 錯誤")
	}
	data := ret[0].MsgData.(*netproto.GameServerStartRet)

	assert.Equal(t, data.GetCode(), int32(-1))
	assert.Equal(t, data.GetMessage(), "服務器啟動失敗, Sorry , NotFoundServer")
}

func TestGameServerStartHandler_API_Connection_Error(t *testing.T) {
	InitTest()

	cfgs := new(config.Group)
	cfgs.Configs = make(map[string]*config.Config)
	cfg := new(config.Config)
	cfg.Sections = make(map[string]config.Section)
	sec := config.NewConfigSection()
	sec["schema"] = "http"
	sec["backendapi"] = "test"
	sec["api_cer"] = "qq"
	cfg.Sections["conn"] = sec
	cfgs.Configs["api"] = cfg

	base := service.NewBase(&module.MouduleCtx{
		Config: cfgs,
	})
	service.ServerSettingService.Base = base

	// 驗證 handler
	serverid := int32(365)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	// 驗證資料
	msgdata := new(netproto.GameServerStart)
	msgdata.ServerID = proto.Int32(1)
	msgdata.GameID = proto.Int32(1)
	msgdata.ClientAddr = proto.String("0.0.0.0:30101")
	msgdata.HttpAddr = proto.String("")
	msgdata.NotifyAddr = proto.String("0.0.0.0:9000")
	msgdata.LoginMoney = proto.Int32(1000)
	msgdata.ShowName = proto.String("1")
	msgdata.ServerName = proto.String("多福多財_初級場_測試站")

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData

	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GameServerStartID), rpcMessage)
	handler := new(GameServerStartHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_GameServerStartRetID) {
		t.Error("OnMessage 錯誤")
	}
	data := ret[0].MsgData.(*netproto.GameServerStartRet)

	assert.Equal(t, data.GetCode(), int32(-1))
	assert.Equal(t, data.GetMessage(), "服務器啟動失敗, GameAPI 連線失敗")
}
