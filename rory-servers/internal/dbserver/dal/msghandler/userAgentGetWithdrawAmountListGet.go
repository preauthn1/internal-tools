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

// UserAgentGetWithdrawAmountListHandler 申請提現例表(我的提現)
type UserAgentGetWithdrawAmountListHandler struct{}

// IsHook ?
func (rh *UserAgentGetWithdrawAmountListHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:UserAgentGetWithdrawAmountListHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserAgentGetWithdrawAmountListID)
}

// OnMessage ?
func (rh *UserAgentGetWithdrawAmountListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:UserAgentGetWithdrawAmountListHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserAgentGetWithdrawAmountListMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserAgentGetWithdrawAmountListMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID())       //使用者ID
	ps.AddIntInput("intPageIndex", mm.GetPageIndex()) //第幾頁
	ps.AddIntInput("intPageSize", mm.GetPageSize())   //共幾頁
	ps.AddIntOutput("intPageCount", 0)
	ps.AddIntOutput("intSumReward", 0)
	ps.AddIntOutput("intSumConvert", 0)
	ps.AddIntOutput("intRemainConvert", 0)

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Unlimit_GetMyConvert", ps)
	// mlog.System("HANDLER名稱:UserAgentGetWithdrawAmountListHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:UserAgentGetWithdrawAmountListHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithdrawAmountListRetID), nil)
	retdata := new(netproto.UserAgentGetWithdrawAmountListRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_Unlimit_GetMyConvert 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[1] // 第一張表
			if len(tb.Rows) > 0 {
				retdata.Userid = proto.Int32(int32(mm.GetUserID()))
				retdata.Pagecount = proto.Int32(int32(tb.GetValueByColName(0, "intPageCount").(int64))) // 使用數參數0
				retdata.Bonus = proto.Int32(int32(tb.GetValueByColName(0, "intSumReward").(int64)))     // 使用數參數0
				retdata.Rest = proto.Int32(int32(tb.GetValueByColName(0, "intSumConvert").(int64)))     // 使用數參數0
				retdata.Allow = proto.Int32(int32(tb.GetValueByColName(0, "intRemainConvert").(int64))) // 使用數參數0
			} else { // 這裡一定要回傳
				retdata.Userid = proto.Int32(int32(mm.GetUserID()))
				retdata.Pagecount = proto.Int32(int32(0)) // 使用數參數0
				retdata.Bonus = proto.Int32(int32(0))     // 使用數參數0
				retdata.Rest = proto.Int32(int32(0))      // 使用數參數0
				retdata.Allow = proto.Int32(int32(0))     // 使用數參數0
			}
			// 第二張表
			tb = dbret[0] // 第二張表
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.UserAgentGetWithdrawAmountListRetList) // 產生新物件，並APPEND
					subRetData.Datetime = proto.String(string(tb.GetValueByColName(i, "StatTime").(string)))
					// subRetData.Datetime = proto.String(tb.GetValueByColName(i, "datetime").(time.Time).Format("2006-01-02 15:04:05"))
					subRetData.Noid = proto.Int32(int32(tb.GetValueByColName(i, "OrderId").(int64)))
					subRetData.Amount = proto.Int32(int32(tb.GetValueByColName(i, "Amount").(int64)))
					subRetData.Type = proto.String(string(tb.GetValueByColName(i, "Type").(string)))
					subRetData.Account = proto.String(tb.GetValueByColName(i, "Account").(string))
					// subRetData.Account = proto.Int32(int32(tb.GetValueByColName(i, "Account").(int64)))
					subRetData.Status = proto.Int32(int32(tb.GetValueByColName(i, "Status").(int64)))
					retdata.UserAgentGetWithdrawAmountList = append(retdata.UserAgentGetWithdrawAmountList, subRetData)
				}
			} // 這裡可以不用回傳
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
