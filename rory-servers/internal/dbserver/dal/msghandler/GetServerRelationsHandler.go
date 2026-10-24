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
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GetServerRelationsHandler 取得遊戲渠道, 只撈取渠道看不到的ServerID
type GetServerRelationsHandler struct{}

// IsHook ?
func (rh *GetServerRelationsHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_GetAreaServerRelationsID)
}

// OnMessage ?
func (rh *GetServerRelationsHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	ps := db.NewSQLParameters()
	proc := db.NewProcedure("PrGs_GetServerRelations", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(dbret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()

	retdata := new(netproto.GetAreaShowGameFlagListMsg)
	message := "服务器启动成功"
	if err != nil {
		message = fmt.Sprintf("执行存储过程PrPs_ServerStart错误, %v", err)
		mlog.Error(message)
	} else {
		if dbret.GetRetTableCount() == 1 {
			tbData := dbret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				pareaShowGame := new(netproto.AreaShowGameFlag)
				pareaShowGame.AreaID = proto.Int32(int32(tbData.GetValueByColName(i, "AreaID").(int64)))
				pareaShowGame.GameID = proto.Int32(int32(tbData.GetValueByColName(i, "GameID").(int64)))
				pareaShowGame.FlagID = proto.Int32(int32(tbData.GetValueByColName(i, "FlagID").(int64)))
				pareaShowGame.ServerID = proto.Int32(int32(tbData.GetValueByColName(i, "ServerID").(int64)))
				statusBool := tbData.GetValueByColName(0, "Status").(bool)
				status := 1
				if statusBool == false {
					status = 0
				}
				pareaShowGame.Status = proto.Int32(int32(status))
				retdata.AreaShowGameFlagList = append(retdata.AreaShowGameFlagList, pareaShowGame)
			}
		}
	}

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GetAreaServerRelationsRetID), nil)
	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
