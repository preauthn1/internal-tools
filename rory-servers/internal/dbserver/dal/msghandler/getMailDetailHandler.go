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

// GetMailDetailHandler 獲取用戶郵件列表
type GetMailDetailHandler struct{}

// IsHook ?
func (rh *GetMailDetailHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_GetMailDetailID)
}

// OnMessage ?
func (rh *GetMailDetailHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.RequestMailDetail{}
		return mm
	})

	if rpc == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}
	mm := rmm.(*netproto.RequestMailDetail)

	userid := rpc.GetUserID()
	cer := rpc.GetCer()
	id := mm.GetID()
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userid)
	ps.AddIntInput("intMsgID", id)
	ps.AddVarcharInput("chvCer", cer, 32)
	proc := db.NewProcedure("PrPs_Mail_GetMailByID", ps)
	//dbret, err := DbUtil.ExecProc(proc)
	var dbret db.RecordSet
	var err error
	dbret, err = dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMailDetailRetID), nil)
	retdata := new(netproto.MailDetail)
	if err != nil {
		strMsg := fmt.Sprintf("執行存儲過程PrPs_Mail_GetMailByID錯誤, %v", err)
		ctx.Error(strMsg)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	if dbret.GetRetTableCount() > 0 && len(dbret[0].Rows) > 0 {
		retdata.ID = proto.Int32(dbret[0].GetValueByColNameInt32(0, "ID"))
		retdata.Content = proto.String(dbret[0].GetValueByColNameString(0, "Content"))
		retdata.Title = proto.String(dbret[0].GetValueByColNameString(0, "title"))
		retdata.IsRead = proto.Bool(dbret[0].GetValueByColNameInt64(0, "readtype") == 1)
		retdata.SendTime = proto.String(dbret[0].GetValueByColNameString(0, "SendTime"))
	}

	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
