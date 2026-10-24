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

// GetUserLuckScoreHandler ?
type GetUserLuckScoreHandler struct{}

// IsHook ?
func (rh *GetUserLuckScoreHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_UserLuckScoreID)
}

// OnMessage ?
func (rh *GetUserLuckScoreHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserLuckScore{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserLuckScore)

	serverid := rpc.GetRouteServerID()
	userid := mm.GetUserID()
	score := mm.GetScore()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intServerID", serverid)
	ps.AddIntInput("intUserID", userid)
	ps.AddBigIntInput("lngScore", score)
	ps.AddVarcharOutput("chvValue", "")

	proc := db.NewProcedure("PrGs_GetUserLuckScore", ps)

	//ret, err := DbUtil.ExecProc(proc)
	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	if err != nil {
		ctx.Error("執行存儲過程PrGs_ServerReg錯誤, %v", err)
		return nil
	}

	luckValue := ret.GetOutputParamValue("chvValue").(string)
	retData := new(netproto.UserLuckScoreRet)
	retData.Value = proto.String(luckValue)
	retData.UserID = proto.Int32(userid)
	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Game)
	retmsg.SClassID = int32(netproto.GameMessageClassID_UserLuckScoreRetID)
	retmsg.MsgData = retData
	return []*msg.Message{retmsg}

}
