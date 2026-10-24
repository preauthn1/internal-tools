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
)

// GetAreaConfigHandler Hall 取得AreaId的AreaConfig資訊
type GetAreaConfigHandler struct{}

// IsHook ?
func (rh *GetAreaConfigHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) &&
		sclsid == int32(netproto.DBServerClassID_GetAreaConfigID)
}

// OnMessage ?
func (rh *GetAreaConfigHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetHallAreaMsg{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetHallAreaMsg)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intAreaID", int32(mm.GetAreaID()))
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrGs_ServerReStart_Areaconfig_Get", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	mlogWithSQL := mlog.NewDefaultLogger().SetData(proc.GetSQLString())
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GetAreaConfigRetID), nil)
	message := "服务器启动成功"
	if err != nil {
		message = fmt.Sprintf("执行存储过程PrGs_ServerReStart_Areaconfig_Get错误, %v", err)
		mlogWithSQL.Error(message)
	}
	retdata := new(netproto.AreaConfig)

	if dbret.GetRetTableCount() >= 1 {
		retdata.AreaID = proto.Int64(dbret[0].GetSingleValue("AreaID").(int64))
		retdata.AreaConfig = proto.String(dbret[0].GetSingleValue("AreaConfig").(string))
		retdata.GiftMoney = proto.Int64(dbret[0].GetSingleValue("GiftMoney").(int64))
		retdata.DefaultMoney = proto.Int64(dbret[0].GetSingleValue("DefaultMoney").(int64))
	} else {
		mlogWithSQL.Error("未获取到平台版本信息 areaid[%v] ", mm.GetAreaID())
		retdata.AreaID = proto.Int64(0)
		retdata.AreaConfig = proto.String("")
	}

	retdata.NotifyRefresh = proto.Bool(mm.GetNotifyRefresh())

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
