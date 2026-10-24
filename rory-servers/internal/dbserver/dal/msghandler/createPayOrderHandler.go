package msghandler

import (
	"math"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// CreatePayOrderHandler 創建充值訂單
type CreatePayOrderHandler struct{}

// IsHook ?
func (rh *CreatePayOrderHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_CreatePayOrderID)
}

// OnMessage ?
func (rh *CreatePayOrderHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.CreatePayOrder{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.CreatePayOrder)
	paytypeid := mm.GetPayTypeID()
	amount := int64(math.Abs(float64(mm.GetAmount())))
	userid := rpcinfo.GetUserID()
	ip := rpcinfo.GetIPAddress()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intPayTypeID", paytypeid)
	ps.AddIntInput("intUserID", userid)
	ps.AddBigIntInput("lngAmount", amount)
	ps.AddVarcharInput("chvIPAddress", ip, 15)
	ps.AddIntOutput("intOrderID", 0)
	ps.AddVarcharOutput("chvErrMsg", "")
	proc := db.NewProcedure("PrPs_Pay_CreateOrder", ps)
	//fmt.Println(proc.GetSQLString())
	//ret, err := DbUtil.ExecProc(proc)
	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	var code int32 = 0
	var orderid int32 = 0
	retmessage := ""

	if err != nil {
		ctx.Error("执行存储过程PrPs_Pay_CreateOrder错误, %v", err)
		code = 0
		retmessage = "服务器错误"
	} else {
		if ret.GetReturnValue() == 1 {
			code = 1
			orderid = int32(ret.GetOutputParamValue("intOrderID").(int64))
		}
		retmessage = ret.GetOutputParamValue("chvErrMsg").(string)
	}

	retData := &netproto.CreatePayOrderRet{}
	retData.Message = proto.String(retmessage)
	retData.Code = proto.Int32(code)
	retData.OrderID = proto.Int32(orderid)

	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retmsg.SClassID = int32(netproto.HallMsgClassID_CreatePayOrderRetID)
	retmsg.MsgData = retData

	return []*msg.Message{retmsg}
}
