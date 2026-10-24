package msghandler

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetContestUserDataHandler ?
type GetContestUserDataHandler struct{}

// IsHook ?
func (rh *GetContestUserDataHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.DBServerClassID_GetContestUserDataID)
}

// OnMessage ?
func (rh *GetContestUserDataHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetConestUser{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetConestUser)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", mm.GetGameID())
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrGs_Event_GetContestUserData", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.DBServerClassID_GetContestUserDataRetID), nil)
	retdata := new(netproto.UpdateContestUser)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrGs_Event_GetContestUserData错误, %v", err)
		ctx.Error(strMsg)
	} else {
		if dbret.GetRetTableCount() > 0 && len(dbret[0].Rows) > 0 {
			retdata.GameID = proto.Int32(int32(dbret[0].GetValueByColName(0, "GameID").(int64)))
			retdata.UserID = proto.Int32(int32(dbret[0].GetValueByColName(0, "UserID").(int64)))
			retdata.ContestID = proto.Int32(int32(dbret[0].GetValueByColName(0, "ContestID").(int64)))
			retdata.ContestNo = proto.Int32(int32(dbret[0].GetValueByColName(0, "ContestNo").(int64)))
			retdata.FlagID = proto.Int32(int32(dbret[0].GetValueByColName(0, "FlagID").(int64)))
			retdata.TableNo = proto.Int32(int32(dbret[0].GetValueByColName(0, "TableNo").(int64)))
			retdata.PlayTime = proto.Int32(int32(dbret[0].GetValueByColName(0, "PlayTime").(int64)))
			retdata.RemainScore = proto.Int32(int32(dbret[0].GetValueByColName(0, "RemainScore").(int64)))
			retdata.Score = proto.Int32(int32(dbret[0].GetValueByColName(0, "Score").(int64)))
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
