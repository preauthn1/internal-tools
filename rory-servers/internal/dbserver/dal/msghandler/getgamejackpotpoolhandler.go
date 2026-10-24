package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// GetGameJackpotPoolHandler 查询遊戲獎金池
type GetGameJackpotPoolHandler struct{}

// IsHook ?
func (rh *GetGameJackpotPoolHandler) IsHook(bclsid int32, sclsid int32) bool {
	return (bclsid == int32(netproto.MessageBClassID_Hall) ||
		bclsid == int32(netproto.MessageBClassID_Game)) &&
		sclsid == int32(netproto.HallMsgClassID_GetJackpotPoolID)
}

// OnMessage ?
func (rh *GetGameJackpotPoolHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	retmsg := new(msg.Message)
	retmsg.BClassID = m.BClassID
	retmsg.SClassID = int32(netproto.HallMsgClassID_GetJackpotPoolRetID)
	retData := new(netproto.JackpotPoolRet)

	ps1 := db.NewSQLParameters()
	ps1.AddIntInput("intFlagID", 0)
	proc2 := db.NewProcedure("PrGs_GetQZNN_SystemPool", ps1)
	dbret2, err := dbutil.DBUtil.ExecProc(proc2)
	retcode2 := int32(dbret2.GetReturnValue())
	if err != nil {
		ctx.Error("執行存儲過程PrGs_GetQZNN_SystemPool錯誤, %v, %s", err, proc2.GetSQLString())
		retcode2 = int32(0)
	}

	tbData0 := dbret2[0]
	if retcode2 == 1 {
		if len(tbData0.Rows) > 0 {
			retData.Amount = proto.Int64(tbData0.GetValueByColNameInt64(0, "Amount")) //獎金池金額
		} else {
			retData.Amount = proto.Int64(5000000)
			retmsg.MsgData = retData

			return []*msg.Message{retmsg}
		}
	}
	retmsg.MsgData = retData

	return []*msg.Message{retmsg}
}
