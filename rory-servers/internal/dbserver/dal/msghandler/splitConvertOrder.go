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

// SplitConvertOrderHandler 裂變5提現紀錄送出提現金額
type SplitConvertOrderHandler struct{}

// IsHook ?
func (rh *SplitConvertOrderHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:SplitConvertOrderHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_SplitConvertOrderID)
}

// OnMessage ?
func (rh *SplitConvertOrderHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SplitConvertOrderMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SplitConvertOrderMsg)

	if mm.GetLngAmount() > 0 {
		returnMessage := rh.OnMessageInner(ctx, clt, m)
		return returnMessage
	}

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitConvertOrderRetID), nil)
	retdata := new(netproto.SplitConvertOrderRet) // 回傳封包

	if mm.GetLngAmount() <= 0 {
		retdata.ChvMsg = proto.String("可提现余额不正确，金额需大于0")
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}

// OnMessageInner ?
func (rh *SplitConvertOrderHandler) OnMessageInner(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:SplitConvertOrderHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SplitConvertOrderMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SplitConvertOrderMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID())                  //使用者ID
	ps.AddVarcharInput("chvIPaddress", mm.GetChvIPaddress(), 32) //IP位置
	ps.AddIntInput("lngAmount", mm.GetLngAmount())               //提領金額 // 要做防呆
	ps.AddVarcharOutput("chvMsg", "")                            //chvMsg的輸出

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Split_ConvertOrder", ps)
	// mlog.System("HANDLER名稱:SplitConvertOrderHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:SplitConvertOrderHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitConvertOrderRetID), nil)
	retdata := new(netproto.SplitConvertOrderRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_Split_ConvertOrder 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.ChvMsg = proto.String(string(tb.GetValueByColName(0, "chvMsg").(string)))
			} else { // 必需回傳
				retdata.ChvMsg = proto.String("")
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
