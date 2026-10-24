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

// UserAgentResultListHandler 我的業績(全民代理)
type UserAgentResultListHandler struct{}

// IsHook ?
func (rh *UserAgentResultListHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:UserAgentResultListHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserAgentResultListID)
}

// OnMessage ?
func (rh *UserAgentResultListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:UserAgentResultListHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserAgentResultListMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserAgentResultListMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID())       //使用者ID
	ps.AddIntInput("intPageIndex", mm.GetPageIndex()) //第幾頁
	ps.AddIntInput("intPageSize", mm.GetPageSize())   //共幾頁
	ps.AddIntOutput("intPageCount", 0)

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Unlimit_GetMyPerformance", ps)
	// mlog.System("HANDLER名稱:UserAgentResultListHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:UserAgentResultListHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentResultListRetID), nil)
	retdata := new(netproto.UserAgentResultListRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_Unlimit_GetMyPerformance 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Result = proto.Int32(int32(tb.GetValueByColName(0, "WeekMyFlows").(int64)))           // 使用參數0
				retdata.Groupresult = proto.Int32(int32(tb.GetValueByColName(0, "WeekTeamFlows").(int64)))    // 使用參數0
				retdata.Directresult = proto.Int32(int32(tb.GetValueByColName(0, "WeekMemberFlows").(int64))) // 使用參數0
				retdata.Bonus = proto.Int32(int32(tb.GetValueByColName(0, "Reward").(int64)))                 // 使用參數0
			} else { // 一定要回傳
				retdata.Result = proto.Int32(int32(0))
				retdata.Groupresult = proto.Int32(int32(0))
				retdata.Directresult = proto.Int32(int32(0))
				retdata.Bonus = proto.Int32(int32(0))
			}
			// 第二張表
			tb = dbret[1]
			if len(tb.Rows) > 0 { // 可以不用回傳
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.UserAgentResultListRetList)                               // 產生新物件，並APPEND
					subRetData.Time = proto.String(string(tb.GetValueByColName(i, "StatTime").(string))) // 使用參數i
					// subRetData.Tim = proto.String(tb.GetValueByColName(i, "Time").(time.Time).Format("2006-01-02 15:04:05"))
					subRetData.Daydirectresult = proto.Int32(int32(tb.GetValueByColName(i, "WeekMemberFlows").(int64))) // 使用參數i
					subRetData.Daygroupresult = proto.Int32(int32(tb.GetValueByColName(i, "WeekTeamFlows").(int64)))    // 使用參數i
					subRetData.Dayresult = proto.Int32(int32(tb.GetValueByColName(i, "WeekMyFlows").(int64)))           // 使用參數i
					retdata.UserAgentResultList = append(retdata.UserAgentResultList, subRetData)
				}
			}
			// 第三張表
			tb = dbret[2]
			if len(tb.Rows) > 0 {
				retdata.PageCount = proto.Int32(int32(tb.GetValueByColName(0, "intPageCount").(int64)))
			} else { // 一定要回傳
				retdata.PageCount = proto.Int32(0)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
