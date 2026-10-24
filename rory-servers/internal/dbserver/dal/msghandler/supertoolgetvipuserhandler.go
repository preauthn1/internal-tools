package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// SuperToolGetVipUserHandler 寫遊戲記錄詳情
type SuperToolGetVipUserHandler struct{}

// IsHook ?
func (rh *SuperToolGetVipUserHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_SuperToolGetVipUserID)
}

// OnMessage ?
func (rh *SuperToolGetVipUserHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolGetVipUserRequest{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperToolGetVipUserRequest)
	gameid := mm.GetGameID()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", gameid)
	proc := db.NewProcedure("PrGs_SuperTool_GetVIPUser", ps)
	//mlog.System(proc.GetSQLString())
	ret, err := dbutil.DBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("执行存储过程PrGs_SuperTool_GetVIPUser错误, %v, %s", err, proc.GetSQLString())
		retcode = 0
	} else {
		retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_SuperToolGetVipUserRetID), nil)
		retdata := new(netproto.SuperToolVipUserDataList)
		for i := 0; i < len(ret[0].Rows); i++ {
			vipuser := new(netproto.SuperToolVipUserData)

			vipuser.UserID = proto.Int32(int32(ret[0].GetValueByColName(i, "UserID").(int64)))
			vipuser.AssignAmount = proto.Int64(ret[0].GetValueByColName(i, "AssignAmount").(int64))
			vipuser.AssignRate = proto.Int32(int32(ret[0].GetValueByColName(i, "AssignRate").(int64)))
			vipuser.CompleteAmount = proto.Int64(ret[0].GetValueByColName(i, "CompleteAmount").(int64))
			vipuser.Complete = proto.Int32(int32(ret[0].GetValueByColName(i, "Complete").(int64)))
			vipuser.Status = proto.Int32(int32(ret[0].GetValueByColName(i, "Status").(int64)))
			vipuser.TotalAll = proto.Int32(0)
			vipuser.TotalUserWinLose = proto.Int32(0)
			vipuser.TodayUserWinLose = proto.Int32(0)
			retdata.DataList = append(retdata.DataList, vipuser)
		}
		retmsg.MsgData = retdata
		return []*msg.Message{retmsg}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}

}
