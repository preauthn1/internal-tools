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

// UserAgentWinAmountListHandler 我的獎勵(全民代理)
type UserAgentWinAmountListHandler struct{}

// IsHook ?
func (rh *UserAgentWinAmountListHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:UserAgentWinAmountListHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserAgentWinAmountListID)
}

// OnMessage ?
func (rh *UserAgentWinAmountListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:UserAgentWinAmountListHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserAgentWinAmountListMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserAgentWinAmountListMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID
	// ps.AddVarcharInput("dtStattime", mm.GetDateTime(), 32) //時間dtStattime
	// ps.AddVarcharOutput("chvMsg", "") //chvMsg的輸出
	// ps.AddIntInput("intArea_ID", mm.GetArea_ID()) //頻道ID

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Unlimit_GetMyReward", ps)
	// mlog.System("HANDLER名稱:UserAgentWinAmountListHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:UserAgentWinAmountListHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentWinAmountListRetID), nil)
	retdata := new(netproto.UserAgentWinAmountListRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_Unlimit_GetMyReward 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.UserAgentWinAmountListRetList)                                            // 產生新物件，並APPEND
					subRetData.SendTime = proto.String(string(tb.GetValueByColName(i, "StatTime").(string)))             // 使用參數i
					subRetData.Weekresult = proto.Int32(int32(tb.GetValueByColName(i, "WeekMyFlows").(int64)))           // 使用參數i
					subRetData.Weekgroupresult = proto.Int32(int32(tb.GetValueByColName(i, "WeekTeamFLows").(int64)))    // 使用參數i
					subRetData.Weekdirectresult = proto.Int32(int32(tb.GetValueByColName(i, "WeekMemberFlows").(int64))) // 使用參數i
					subRetData.Bonus = proto.Int32(int32(tb.GetValueByColName(i, "Reward").(int64)))                     // 使用參數i
					subRetData.Status = proto.Int32(int32(tb.GetValueByColName(i, "Status").(int64)))                    // 使用參數i
					retdata.UserAgentWinAmount = append(retdata.UserAgentWinAmount, subRetData)
				}
			} // 可以不回傳資料
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
