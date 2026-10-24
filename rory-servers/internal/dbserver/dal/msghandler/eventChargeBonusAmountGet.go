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
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// EventChargeBonusAmountHandler 充值返利活動  取得個人流水及限制
type EventChargeBonusAmountHandler struct{}

// IsHook ?
func (rh *EventChargeBonusAmountHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:EventChargeBonusAmountHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_EventChargeBonusAmountID)
}

// OnMessage ?
func (rh *EventChargeBonusAmountHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:EventChargeBonusAmountHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.EventChargeBonusAmountMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.EventChargeBonusAmountMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID
	ps.AddVarcharOutput("chvMsg", "")           //chvMsg的輸出

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Event_ChargeBonusAmount_Get", ps)
	mlog.System("HANDLER名稱:EventChargeBonusAmountHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue())                                   //retcode必須為1
	mlog.System("HANDLER名稱:EventChargeBonusAmountHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_EventChargeBonusAmountRetID), nil)
	retdata := new(netproto.EventChargeBonusAmountRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程PrPs_Event_ChargeBonusAmount_Get 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Rate = proto.Int32(int32(tb.GetValueByColName(0, "Rate").(int64)))
				retdata.Times = proto.Int32(int32(tb.GetValueByColName(0, "FlowAmountMultiple").(int64)))
				retdata.NowFlowAmount = proto.Int32(int32(tb.GetValueByColName(0, "NowFlowAmount").(int64)))
				retdata.FlowAmountLimit = proto.Int32(int32(tb.GetValueByColName(0, "FlowAmountLimit").(int64)))
			} else {
				retdata.Rate = proto.Int32(int32(0))
				retdata.Times = proto.Int32(int32(0))
				retdata.NowFlowAmount = proto.Int32(int32(0))
				retdata.FlowAmountLimit = proto.Int32(int32(0))
			}
			// 第二張表
			tb = dbret[1]
			if len(tb.Rows) > 0 {
				retdata.ChvMsg = proto.String(string(tb.GetValueByColName(0, "chvMsg").(string)))
			} else {
				retdata.ChvMsg = proto.String("")
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
