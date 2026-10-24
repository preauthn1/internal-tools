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

// SplitGetMyRewardHandler 裂變4馬上提現
type SplitGetMyRewardHandler struct{}

// IsHook ?
func (rh *SplitGetMyRewardHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:SplitGetMyRewardHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_SplitGetMyRewardID)
}

// OnMessage ?
func (rh *SplitGetMyRewardHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:SplitGetMyRewardHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SplitGetMyRewardMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SplitGetMyRewardMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Split_GetMyReward", ps)
	// mlog.System("HANDLER名稱:SplitGetMyRewardHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:SplitGetMyRewardHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyRewardRetID), nil)
	retdata := new(netproto.SplitGetMyRewardRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_Split_GetMyReward 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.MyReward = proto.Int32(int32(tb.GetValueByColName(0, "MyReward").(int64)))   // 可提現金額
				retdata.MyWait = proto.Int32(int32(tb.GetValueByColName(0, "WithdrawTime").(int64))) // 倒數計時
			} else {
				retdata.MyReward = proto.Int32(int32(0)) // 可提領金額
				retdata.MyWait = proto.Int32(int32(0))   // 倒數計時
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
