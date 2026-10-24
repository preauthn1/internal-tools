package msghandler

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameServerStartHandler 遊戲服務器啟動
type GameServerStartHandler struct{}

// IsHook ?
func (rh *GameServerStartHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_GameServerStartID)
}

// OnMessage ?
func (rh *GameServerStartHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GameServerStart{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GameServerStart)

	serverid := mm.GetServerID()
	clientaddr := mm.GetClientAddr()
	notifyaddr := mm.GetNotifyAddr()
	httpaddr := mm.GetHttpAddr()
	gameid := mm.GetGameID()
	servername := mm.GetServerName()
	showname := mm.GetShowName()
	loginmoney := mm.GetLoginMoney()

	ps := db.NewSQLParameters()
	ps.AddIntInput("sintServerID", serverid)
	ps.AddIntInput("sintGameID", gameid)
	ps.AddVarcharInput("chvServerName", servername, 64)
	ps.AddVarcharInput("chvShowName", showname, 64)
	ps.AddVarcharInput("chvClientAddr", clientaddr, 32)
	ps.AddVarcharInput("chvNotifyAddr", notifyaddr, 32)
	ps.AddVarcharInput("chvHttpAddr", httpaddr, 32)
	ps.AddIntInput("lngLoginMoney", loginmoney)
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrGs_ServerStart", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(dbret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()
	var code int32
	code = 1
	message := "服务器启动成功"

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GameServerStartRetID), nil)
	retdata := new(netproto.GameServerStartRet)

	if err != nil {
		message = fmt.Sprintf("执行存储过程PrGs_ServerStart错误, %v", err)
		ctx.Error(message)
		code = -1
	} else if dbret.GetReturnValue() == 0 {
		message = dbret.GetOutputParamValue("chvErrMsg").(string)
		code = 0
	} else {
		retdata.AndroidVersion = proto.String(dbret[0].GetSingleValueString("AndroidVersion"))
		retdata.IOSVersion = proto.String(dbret[0].GetSingleValueString("IOSVersion"))
		retdata.Status = proto.Int32(dbret[0].GetSingleValueInt32("Status"))
		retdata.StatusMessage = proto.String(dbret[0].GetSingleValueString("StatusMessage"))
		code = int32(dbret.GetReturnValue())
		message = dbret.GetOutputParamValueString("chvErrMsg")
	}

	retdata.Code = proto.Int32(code)
	retdata.Message = proto.String(message)
	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
