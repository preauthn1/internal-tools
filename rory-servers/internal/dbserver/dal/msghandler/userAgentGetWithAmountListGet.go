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

// UserAgentGetWithAmountListHandler 申請提現例表(全民代理)
type UserAgentGetWithAmountListHandler struct{}

// IsHook ?
func (rh *UserAgentGetWithAmountListHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:UserAgentGetWithAmountListHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserAgentGetWithAmountListID)
}

// OnMessage ?
func (rh *UserAgentGetWithAmountListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:UserAgentGetWithAmountListHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserAgentGetWithAmountListMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserAgentGetWithAmountListMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID
	if mm.GetTargetID() != 0 {
		ps.AddIntInput("intTargetID", mm.GetTargetID()) //查詢玩家ID
	}
	ps.AddIntInput("intPageIndex", mm.GetPageIndex()) //第幾頁
	ps.AddIntInput("intPageSize", mm.GetPageSize())   //共幾頁
	ps.AddIntOutput("intPageCount", 0)

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Unlimit_GetMyMember", ps)
	// mlog.System("HANDLER名稱:UserAgentGetWithAmountListHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:UserAgentGetWithAmountListHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountListRetID), nil)
	retdata := new(netproto.UserAgentGetWithAmountListRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_Unlimit_GetMyMember 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0] // 使用第一張表
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.UserAgentGetWithAmountListRetList)                                   // 產生新物件，並APPEND
					subRetData.Userid = proto.Int32(int32(tb.GetValueByColName(i, "userid").(int64)))               // 使用參數i
					subRetData.Nickname = proto.String(string(tb.GetValueByColName(i, "nickname").(string)))        // 使用參數i
					subRetData.Groupoffer = proto.Int32(int32(tb.GetValueByColName(i, "WeekTeamFlows").(int64)))    // 使用參數i
					subRetData.Playeroffer = proto.Int32(int32(tb.GetValueByColName(i, "WeekMemberFlows").(int64))) // 使用參數i
					subRetData.Membernumbers = proto.Int32(int32(tb.GetValueByColName(i, "TeamMember").(int64)))    // 使用參數i
					// subRetData.Agent = proto.Int32(int32(tb.GetValueByColName(i, "Agent").(int64)))                 // 使用參數i
					subRetData.Belongs = proto.Int32(int32(tb.GetValueByColName(i, "Member").(int64)))  // 使用參數i
					subRetData.Type = proto.Int32(int32(tb.GetValueByColName(i, "MemberType").(int64))) // 使用參數i
					retdata.UserAgentGetWithAmountList = append(retdata.UserAgentGetWithAmountList, subRetData)
				}
			} // 沒資料不回傳
			tb = dbret[1]
			if len(tb.Rows) > 0 {
				retdata.PageCount = proto.Int32(int32(tb.GetValueByColName(0, "intPageCount").(int64))) // 使用參數0
			} else {
				retdata.PageCount = proto.Int32(0)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
