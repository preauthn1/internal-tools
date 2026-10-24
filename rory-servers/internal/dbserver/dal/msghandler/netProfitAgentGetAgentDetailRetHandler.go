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

// NetProfitAgentGetAgentDetailRetHandler USER 推廣信息(淨利潤)
type NetProfitAgentGetAgentDetailRetHandler struct{}

// IsHook ?
func (rh *NetProfitAgentGetAgentDetailRetHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentGetAgentDetailID)
}

// OnMessage ?
func (rh *NetProfitAgentGetAgentDetailRetHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentGetAgentDetailMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentGetAgentDetailMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddVarcharOutput("chvMsg", "")

	//開始準備SP指令
	proc := db.NewProcedure("PrPs_NetProfitAgent_GetAgentDetail", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentDetailRetID), nil)

	retdata := new(netproto.NetProfitAgentGetAgentDetailRet) // 不是LIST，在這宣告
	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("執行存儲過程 PrPs_NetProfitAgent_GetAgentDetail 錯誤,%v", err)
		ctx.Error(strMsg)
	} else {

		if retcode == 1 {
			tb := dbret[1]
			retdata.ChvMsg = proto.String(string(tb.GetValueByColName(0, "ChvMsg").(string)))
			tb = dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Rate = proto.Int32(int32(tb.GetValueByColName(0, "rate").(int64)))
				retdata.SumbonusAmount = proto.Int32(int32(tb.GetValueByColName(0, "sumbonusAmount").(int64)))
				retdata.Poolamount = proto.Int32(int32(tb.GetValueByColName(0, "poolamount").(int64)))
				retdata.Newagenttoday = proto.Int32(int32(tb.GetValueByColName(0, "newagenttoday").(int64)))
				retdata.Newagentlast = proto.Int32(int32(tb.GetValueByColName(0, "newagentlast").(int64)))
				retdata.Totalagentuser = proto.Int32(int32(tb.GetValueByColName(0, "totalagentuser").(int64)))
				retdata.Newusertoday = proto.Int32(int32(tb.GetValueByColName(0, "newusertoday").(int64)))
				retdata.Chargelast = proto.Int32(int32(tb.GetValueByColName(0, "chargelast").(int64)))
				retdata.Convertlast = proto.Int32(int32(tb.GetValueByColName(0, "convertlast").(int64)))
				retdata.Bonusamountlast = proto.Int32(int32(tb.GetValueByColName(0, "bonusamountlast").(int64)))
				retdata.Chargetoday = proto.Int32(int32(tb.GetValueByColName(0, "chargetoday").(int64)))
				retdata.Converttoday = proto.Int32(int32(tb.GetValueByColName(0, "converttoday").(int64)))
				retdata.Bonusamounttoday = proto.Int32(int32(tb.GetValueByColName(0, "bonusamounttoday").(int64)))
				retdata.Payusertoday = proto.Int32(int32(tb.GetValueByColName(0, "payusertoday").(int64)))
				retdata.Payuserlast = proto.Int32(int32(tb.GetValueByColName(0, "payuserlast").(int64)))
				retdata.Storageamount = proto.Int32(int32(tb.GetValueByColName(0, "storageamount").(int64)))
				retdata.Qrurl = proto.String(string(tb.GetValueByColName(0, "qrurl").(string)))
				retdata.Withdrawsecs = proto.Int32(int32(tb.GetValueByColName(0, "withdrawsecs").(int64)))
			}
		} else {
			retdata.ChvMsg = proto.String("")
			retdata.Rate = proto.Int32(0)
			retdata.SumbonusAmount = proto.Int32(0)
			retdata.Poolamount = proto.Int32(0)
			retdata.Newagenttoday = proto.Int32(0)
			retdata.Newagentlast = proto.Int32(0)
			retdata.Totalagentuser = proto.Int32(0)
			retdata.Newusertoday = proto.Int32(0)
			retdata.Chargelast = proto.Int32(0)
			retdata.Convertlast = proto.Int32(0)
			retdata.Bonusamountlast = proto.Int32(0)
			retdata.Chargetoday = proto.Int32(0)
			retdata.Converttoday = proto.Int32(0)
			retdata.Bonusamounttoday = proto.Int32(0)
			retdata.Payusertoday = proto.Int32(0)
			retdata.Payuserlast = proto.Int32(0)
			retdata.Storageamount = proto.Int32(0)
			retdata.Qrurl = proto.String("")
			retdata.Withdrawsecs = proto.Int32(0)
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
