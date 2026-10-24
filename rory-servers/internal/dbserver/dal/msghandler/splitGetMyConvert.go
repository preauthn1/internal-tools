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

// SplitGetMyConvertHandler 裂變2提現紀錄
type SplitGetMyConvertHandler struct{}

// IsHook ?
func (rh *SplitGetMyConvertHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:SplitGetMyConvertHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_SplitGetMyConvertID)
}

// OnMessage ?
func (rh *SplitGetMyConvertHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:SplitGetMyConvertHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SplitGetMyConvertMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SplitGetMyConvertMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID())       //使用者ID
	ps.AddIntInput("intPageIndex", mm.GetPageIndex()) // 第幾頁
	ps.AddIntInput("intPageSize", mm.GetPageSize())   // 顯示幾筆
	ps.AddIntOutput("intPageCount", 0)                // 回傳頁數
	ps.AddIntOutput("lngRewardAmount", 0)             // 總金額

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Split_GetMyConvert", ps)
	// mlog.System("HANDLER名稱:SplitGetMyConvertHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:SplitGetMyConvertHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyConvertRetID), nil)
	retdata := new(netproto.SplitGetMyConvertRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_Split_GetMyConvert 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第二張表
			tb := dbret[1]
			if len(tb.Rows) > 0 {
				retdata.PageCount = proto.Int32(int32(tb.GetValueByColName(0, "intPageCount").(int64)))
				retdata.RewardAmount = proto.Int32(int32(tb.GetValueByColName(0, "lngRewardAmount").(int64)))
			} else { // 一定要回傳
				retdata.PageCount = proto.Int32(int32(0))
				retdata.RewardAmount = proto.Int32(int32(0))
			}
			// 第二張表
			tb = dbret[0]
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.SplitGetMyConvertList) // 產生新物件，並APPEND
					subRetData.LogDate = proto.String(string(tb.GetValueByColName(i, "LogDate").(string)))
					subRetData.ConvertAmount = proto.Int32(int32(tb.GetValueByColName(i, "ConvertAmount").(int64)))
					retdata.SplitGetMyConvert = append(retdata.SplitGetMyConvert, subRetData)
				}
			} /* 陣列不回傳 else {
				subRetData := new(netproto.SplitGetMyConvertList) // 產生新物件，並APPEND
				subRetData.LogDate = proto.String("123")
				subRetData.ConvertAmount = proto.Int32(0)
				retdata.SplitGetMyConvert = append(retdata.SplitGetMyConvert, subRetData)
			}
			*/
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
