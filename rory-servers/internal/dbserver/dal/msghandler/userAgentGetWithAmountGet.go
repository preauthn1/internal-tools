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

// UserAgentGetWithAmountHandler 申請提現(全民代理)
type UserAgentGetWithAmountHandler struct{}

// IsHook ?
func (rh *UserAgentGetWithAmountHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:UserAgentGetWithAmountHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserAgentGetWithAmountID)
}

// OnMessage ?
func (rh *UserAgentGetWithAmountHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserAgentGetWithAmountMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserAgentGetWithAmountMsg)

	if mm.GetWithdrawAmount() > 0 {
		returnMessage := rh.OnMessageInner(ctx, clt, m, mm)
		return returnMessage
	}

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountRetID), nil)
	retdata := new(netproto.UserAgentGetWithAmountRet) // 回傳封包

	if mm.GetWithdrawAmount() <= 0 {
		retdata.Status = proto.Int32(0)
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}

// OnMessageInner ?
func (rh *UserAgentGetWithAmountHandler) OnMessageInner(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message, mm *netproto.UserAgentGetWithAmountMsg) []*msg.Message {
	// mlog.System("HANDLER名稱:UserAgentGetWithAmountHandler,進入OnMessage") // DEBUG專用

	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID())                 //使用者ID
	ps.AddIntInput("intWithdrawAmount", mm.GetWithdrawAmount()) //提取金額

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_UserAgentGetWithAmount_Get", ps)
	// mlog.System("HANDLER名稱:UserAgentGetWithAmountHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:UserAgentGetWithAmountHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountRetID), nil)
	retdata := new(netproto.UserAgentGetWithAmountRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_PrPs_UserAgentGetWithAmount_Get 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Status = proto.Int32(int32(tb.GetValueByColName(0, "Status").(int64)))
			} else {
				retdata.Status = proto.Int32(int32(0))
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
