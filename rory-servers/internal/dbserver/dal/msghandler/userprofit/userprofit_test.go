package userprofit

import (
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func TestWithdrawRequestHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := mockRPCInfo(15267605)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.UserProfitWithdrawRequest)
	msgdata.Amount = proto.Int64(100)
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitWithdrawLogID), rpcMessage)
	handler := new(WithdrawRequestHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if (ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserProfitWithdrawLogRetID)) &&
		(ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID)) {
		t.Error("OnMessage 錯誤")
	}
}

func TestGetWithdrawLogHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := mockRPCInfo(15267605)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	msgdata := new(netproto.UserProfitDailySharingRequest)
	msgdata.LimitDays = proto.Int32(7)
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetWithdrawLogID), rpcMessage)
	handler := new(GetWithdrawLogHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if (ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserProfitGetWithdrawLogRetID)) &&
		(ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID)) {
		t.Error("OnMessage 錯誤")
	}
}

func TestGetDailySharingHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := mockRPCInfo(15267605)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	msgdata := new(netproto.UserProfitDailySharingRequest)
	msgdata.LimitDays = proto.Int32(7)
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetDailySharingID), rpcMessage)
	handler := new(GetDailySharingHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if (ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserProfitGetDailySharingRetID)) &&
		(ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID)) {
		t.Error("OnMessage 錯誤")
	}
}

func TestGetExecInfoHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := mockRPCInfo(15267605)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetExecInfoID), rpcMessage)
	handler := new(GetExecInfoHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if (ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserProfitGetExecInfoRetID)) &&
		(ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID)) {
		t.Error("OnMessage 錯誤")
	}
}

func TestGetSharingPercentageHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := mockRPCInfo(15267605)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetSharingPercentageID), rpcMessage)
	handler := new(GetSharingPercentageHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if (ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserProfitGetSharingPercentageRetID)) &&
		(ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID)) {
		t.Error("OnMessage 錯誤")
	}
}

func initDBConnection() {
	_, filename, _, _ := runtime.Caller(0)

	//加载配置信息
	configDir := path.Dir(filename) + "/../config"

	//將單元測試編譯成執行程式, 將下行configDir註解打開, 在console下指令go test -c , 即可找到msghandler.test檔案
	//在C建立dbserverTestTool資料夾, 新增log資料夾, 新增config資料夾與設定. 將檔案與資料放一起即可.
	//configDir = "C:/dbserverTestTool/config"

	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)

	ctx := new(module.MouduleCtx)
	//初始化配置
	ctx.Config = cfgs
	//设置数据库连接
	dbutil.ConnString = ctx.GetConfigData("db", "conn", "connstring")

	dbutil.ConnString = strings.Replace(dbutil.ConnString, "\"", "", -1)
	dbutil.DealLogConnString = ctx.GetConfigData("db", "conn", "deallogconnstring")
	dbutil.DealLogConnString = strings.Replace(dbutil.DealLogConnString, "\"", "", -1)
	dbutil.ClubConnString = ctx.GetConfigData("db", "conn", "clubconnstring")
	dbutil.ClubConnString = strings.Replace(dbutil.ClubConnString, "\"", "", -1)

	dbutil.SMSWebAPI = ctx.GetConfigData("db", "conn", "smswebapi")

	//设置连接池大小
	psize := ctx.GetConfigData("db", "conn", "pool")
	ipsize, _ := strconv.Atoi(psize)
	dbutil.PoolSize = ipsize

	dbutil.InitDatabase()
}

func mockRPCInfo(userID int) *netproto.RPCInfo {
	ipaddr := "0.0.0.0"
	userid := int32(userID)
	rpcInfo := new(netproto.RPCInfo)
	rpcInfo.IPAddress = proto.String(ipaddr)
	rpcInfo.UserID = proto.Int32(userid)
	rpcInfo.ConnID = proto.String("fake")
	rpcInfo.Cer = proto.String("665E301DD8A44BF38C7FFF319052FB3A")
	rpcInfo.RouteServerID = proto.Int32(1)
	return rpcInfo
}
