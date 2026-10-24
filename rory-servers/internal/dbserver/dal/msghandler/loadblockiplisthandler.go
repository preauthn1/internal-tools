package msghandler

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// LoadBlockIPListHandler 加載封鎖IP
type LoadBlockIPListHandler struct{}

// IsHook ?
func (rh *LoadBlockIPListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_LoadBlockIPID)
}

// OnMessage ?
func (rh *LoadBlockIPListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	ps := db.NewSQLParameters()
	proc := db.NewProcedure("PrSys_LoadBlockIPList", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(dbret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_BlockIPListID), nil)
	retdata := new(netproto.BlockIPList)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrSys_LoadBlockIPList错误, %v", err)
		ctx.Error(strMsg)
	} else {
		if dbret.GetRetTableCount() > 0 {
			tbIPLst := dbret[0]
			for i := 0; i < len(tbIPLst.Rows); i++ {
				ips := tbIPLst.GetValueByColNameString(i, "IPAddress")

				retdata.IPAddress = append(retdata.IPAddress, ips)
			}
		}
	}

	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
