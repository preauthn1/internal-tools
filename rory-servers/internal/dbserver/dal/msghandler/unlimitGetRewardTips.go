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

// GetRewardTipsHandler 無限級獎勵說明
type GetRewardTipsHandler struct{}

// IsHook ?
func (rh *GetRewardTipsHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:GetRewardTipsHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetRewardTipsID)
}

// OnMessage ?
func (rh *GetRewardTipsHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:GetRewardTipsHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息

	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetRewardTipsMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	// mm := rmm.(*netproto.GetRewardTipsMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	// ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Unlimit_GetRewardTips", ps)
	// mlog.System("HANDLER名稱:GetRewardTipsHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:GetRewardTipsHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetRewardTipsRetID), nil)
	retdata := new(netproto.GetRewardTipsRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("執行存儲過程PrPs_PrPs_Unlimit_GetRewardTips 錯誤 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Self = proto.String(string(tb.GetValueByColName(0, "Self").(string))) // 用0參數
				retdata.AgentA = proto.String(string(tb.GetValueByColName(0, "AgentA").(string)))
				retdata.AgentB = proto.String(string(tb.GetValueByColName(0, "AgentB").(string)))
				retdata.AgentC = proto.String(string(tb.GetValueByColName(0, "AgentC").(string)))
				retdata.AgentD = proto.String(string(tb.GetValueByColName(0, "AgentD").(string)))
				retdata.AgentBonus = proto.String(string(tb.GetValueByColName(0, "AgentBonus").(string)))
			} else {
				retdata.Self = proto.String("")
				retdata.AgentA = proto.String("")
				retdata.AgentB = proto.String("")
				retdata.AgentC = proto.String("")
				retdata.AgentD = proto.String("")
				retdata.AgentBonus = proto.String("")
			}
			tb = dbret[1]
			if len(tb.Rows) > 0 {
				tb := dbret[1]
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.GetRewardTipsRetList)                          // 產生新物件，並APPEND
					subRetData.Lv = proto.Int32(int32(tb.GetValueByColName(i, "lv").(int64))) // 用i參數
					subRetData.Name = proto.String(string(tb.GetValueByColName(i, "name").(string)))
					subRetData.Minamount = proto.Int32(int32(tb.GetValueByColName(i, "minamount").(int64)))
					subRetData.Maxamount = proto.Int32(int32(tb.GetValueByColName(i, "maxamount").(int64)))
					subRetData.BonusRate = proto.String(string(tb.GetValueByColName(i, "BonusRate").(string)))
					retdata.GetRewardTipsRetData = append(retdata.GetRewardTipsRetData, subRetData)
				}
			} // 沒資料不回傳
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
