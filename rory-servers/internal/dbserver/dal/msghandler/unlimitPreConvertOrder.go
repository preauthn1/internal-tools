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

// UnlimitPreConvertOrderHandler 無限級按鈕2
type UnlimitPreConvertOrderHandler struct{}

// IsHook ?
func (rh *UnlimitPreConvertOrderHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:UnlimitPreConvertOrderHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UnlimitPreConvertOrderID)
}

// OnMessage ?
func (rh *UnlimitPreConvertOrderHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UnlimitPreConvertOrderMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UnlimitPreConvertOrderMsg)

	if mm.GetAmount() > 0 {
		returnMessage := rh.OnMessageInner(ctx, clt, m, mm)
		return returnMessage
	}

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitPreConvertOrderRetID), nil)
	retdata := new(netproto.UnlimitPreConvertOrderRet) // 回傳封包
	if mm.GetAmount() <= 0 {
		retdata.ChvMsg = proto.String("可提現餘額不正確，金額需大於0")
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}

// OnMessageInner ?
func (rh *UnlimitPreConvertOrderHandler) OnMessageInner(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message, mm *netproto.UnlimitPreConvertOrderMsg) []*msg.Message {
	// mlog.System("HANDLER名稱:UnlimitPreConvertOrderHandler,進入OnMessage") // DEBUG專用

	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID
	ps.AddIntInput("lngAmount", mm.GetAmount()) //需要防呆
	ps.AddVarcharOutput("chvMsg", "")           //chvMsg的輸出

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Unlimit_PreConvertOrder", ps)
	// mlog.System("HANDLER名稱:UnlimitPreConvertOrderHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:UnlimitPreConvertOrderHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitPreConvertOrderRetID), nil)
	retdata := new(netproto.UnlimitPreConvertOrderRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("執行存儲過程 PrPs_Unlimit_PreConvertOrder 錯誤 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
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
