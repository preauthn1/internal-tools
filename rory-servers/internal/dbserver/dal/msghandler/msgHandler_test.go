package msghandler

// @todo feature database bypass
import (
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

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

func TestContactServiceHandler(t *testing.T) {
	initDBConnection()
	cldata := new(netproto.ContactService)
	cldata.QuestionType = proto.Int32(1)
	cldata.Msg = proto.String("test")
	cldata.PlatformID = proto.Int32(1)
	cldata.AreaID = proto.Int32(1)
	msgData, _ := proto.Marshal(cldata)
	rpcInfo := mockRPCInfo(0)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ContactServiceID), rpcMessage)
	handler := new(ContactServiceHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_ContactServiceRetID) {
		t.Error("OnMessage 錯誤")
	}
	t.Log(ret)

}

func TestGetConvertLogHandler(t *testing.T) {
	initDBConnection()
	logRange := int32(4)
	tnyStatus := int32(99)

	cldata := new(netproto.GetConvertLog)
	cldata.LogRange = proto.Int32(logRange)
	cldata.TnyStatus = proto.Int32(tnyStatus)
	msgData, _ := proto.Marshal(cldata)
	rpcInfo := mockRPCInfo(1021216)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetConvertLogID), rpcMessage)

	handler := new(GetConvertLogHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetConvertLogRetID) {
		t.Error("OnMessage 錯誤")
	}
	t.Log(ret)
}

func TestGetBankListHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := mockRPCInfo(1021199)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetBankListID), rpcMessage)
	handler := new(GetBankListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetBankListRetID) {
		t.Error("OnMessage 錯誤")
	}
	t.Log(ret)
}

func TestGetMailListHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := mockRPCInfo(1021199)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMailListID), rpcMessage)
	handler := new(GetMailListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetMailListRetID) {
		t.Error("OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}

func initDBConnection() {
	_, filename, _, _ := runtime.Caller(0)

	//加载配置信息
	configDir := path.Dir(filename) + "/config"

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

func TestDBConnection(t *testing.T) {
	initDBConnection()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intServerID", 1)
	ps.AddIntInput("intCurrentRobotCount", 1)
	ps.AddBigIntInput("lngMinMoneyAmount", 1)
	ps.AddBigIntInput("lngMaxMoneyAmount", 1)
	ps.AddIntOutput("intLogoutRobotCount", 0)
	proc := db.NewProcedure("PrGs_RobotLogon", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	t.Log(dbret, err)
}

//func TestSendSMSCode(t *testing.T) {
//	initDBConnection()
//	SendSMSCode(3)
//}
