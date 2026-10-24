package msghandler

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetAgentInfoHandler 玩家/代理信息(淨利潤)
type GetAgentInfoHandler struct{}

// IsHook ?
func (rh *GetAgentInfoHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentGetAgentInfoID)
}

// OnMessage ?
func (rh *GetAgentInfoHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentGetAgentInfoMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentGetAgentInfoMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_NetProfitAgent_GetAgentInfo", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentInfoRetID), nil)
	retdata := new(netproto.NetProfitAgentGetAgentInfoRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_NetProfitAgent_GetAgentInfo 错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.UserID = proto.Int32(int32(tb.GetValueByColName(0, "UserID").(int64)))
				retdata.RegDate = proto.String(tb.GetValueByColName(0, "RegDate").(time.Time).Format("2006-01-02 15:04:05"))
				retdata.ChargeAmount = proto.Int32(int32(tb.GetValueByColName(0, "ChargeAmount").(int64)))
				retdata.SelfWithDraw = proto.Int32(int32(tb.GetValueByColName(0, "SelfWithDraw").(int64)))
				retdata.Status = proto.String(string(tb.GetValueByColName(0, "Status").(string)))
				retdata.TotalUsers = proto.Int32(int32(tb.GetValueByColName(0, "TotalUsers").(int64)))
				retdata.SumSubUsers = proto.Int32(int32(tb.GetValueByColName(0, "SumSubUsers").(int64)))
				retdata.SumSubAgent = proto.Int32(int32(tb.GetValueByColName(0, "SumSubAgent").(int64)))
				retdata.SumSubWithDrawAmount = proto.Int32(int32(tb.GetValueByColName(0, "SumSubWithDrawAmount").(int64)))
				retdata.QRCode = proto.String(tb.GetValueByColName(0, "QRCode").(string))
				retdata.InviteCode = proto.String(tb.GetValueByColName(0, "InviteCode").(string))
				retdata.UpperUserID = proto.Int32(int32(tb.GetValueByColName(0, "UpperUserID").(int64)))
			} else {
				retdata.UserID = proto.Int32(0)
				retdata.RegDate = proto.String("")
				retdata.ChargeAmount = proto.Int32(0)
				retdata.SelfWithDraw = proto.Int32(0)
				retdata.Status = proto.String("")
				retdata.TotalUsers = proto.Int32(0)
				retdata.SumSubUsers = proto.Int32(0)
				retdata.SumSubAgent = proto.Int32(0)
				retdata.SumSubWithDrawAmount = proto.Int32(0)
				retdata.QRCode = proto.String("")
				retdata.InviteCode = proto.String("")
				retdata.UpperUserID = proto.Int32(0)
			}
		} else {
			strMsg := fmt.Sprintf("执行存储过程 PrPs_NetProfitAgent_GetAgentInfo 错误,%v", err)
			ctx.Error(strMsg)
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
