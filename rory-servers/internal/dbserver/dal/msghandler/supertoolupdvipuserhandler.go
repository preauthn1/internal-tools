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

// SuperToolUpdVIPUserHandler 寫遊戲記錄詳情
type SuperToolUpdVIPUserHandler struct{}

// IsHook ?
func (rh *SuperToolUpdVIPUserHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_SuperToolUpdVipUserID)
}

// OnMessage ?
func (rh *SuperToolUpdVIPUserHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolVipUserData{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperToolVipUserData)
	userid := mm.GetUserID()
	completeamount := mm.GetCompleteAmount()
	complete := mm.GetComplete()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userid)
	ps.AddBigIntInput("intCompleteAmount", completeamount)
	ps.AddIntInput("intComplete", complete)
	proc := db.NewProcedure("PrGs_SuperTool_UpdVIPUser", ps)
	ret, err := dbutil.DBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("執行存儲過程PrGs_SuperTool_UpdVIPUser錯誤, %v, %s", err, proc.GetSQLString())
		retcode = 0
	} else {
		retcode = int32(ret.GetReturnValue())
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}

}
