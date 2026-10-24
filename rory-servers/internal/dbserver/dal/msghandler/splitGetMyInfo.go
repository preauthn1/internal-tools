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

// SplitGetMyInfoHandler 裂變1推廣信息
type SplitGetMyInfoHandler struct{}

// IsHook ?
func (rh *SplitGetMyInfoHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:SplitGetMyInfoHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_SplitGetMyInfoID)
}

// OnMessage ?
func (rh *SplitGetMyInfoHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:SplitGetMyInfoHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SplitGetMyInfoMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SplitGetMyInfoMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID
	ps.AddIntOutput("lngRewardAmount", 0)       //可提現金額
	ps.AddIntOutput("lngSplitAmount", 0)        //裂變用戶額
	ps.AddIntOutput("intSuccessCount", 0)       //所有成功人數
	ps.AddIntOutput("intSuccessToday", 0)       //今日成功人數
	ps.AddIntOutput("intSuccessLastday", 0)     //昨日成功人數
	ps.AddVarcharOutput("chvQRCode", "")        //QRCODE
	ps.AddVarcharOutput("chvInviteCode", "")    //邀請碼
	ps.AddIntOutput("UpperUserid", 0)           //上級

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Split_GetMyInfo", ps)
	// mlog.System("HANDLER名稱:SplitGetMyInfoHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:SplitGetMyInfoHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitGetMyInfoRetID), nil)
	retdata := new(netproto.SplitGetMyInfoRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_Split_GetMyInfo 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.RewardAmount = proto.Int32(int32(tb.GetValueByColName(0, "lngRewardAmount").(int64)))
				retdata.SplitAmount = proto.Int32(int32(tb.GetValueByColName(0, "lngSplitAmount").(int64)))
				retdata.SuccessCount = proto.Int32(int32(tb.GetValueByColName(0, "intSuccessCount").(int64)))
				retdata.SuccessToday = proto.Int32(int32(tb.GetValueByColName(0, "intSuccessToday").(int64)))
				retdata.SuccessLastday = proto.Int32(int32(tb.GetValueByColName(0, "intSuccessLastday").(int64)))
				retdata.QrURL = proto.String(string(tb.GetValueByColName(0, "chvQRCode").(string)))
				retdata.Invited = proto.String(string(tb.GetValueByColName(0, "chvInviteCode").(string)))
				retdata.UpperUserID = proto.Int32(int32(tb.GetValueByColName(0, "UpperUserid").(int64)))
			} else {
				retdata.RewardAmount = proto.Int32(int32(0))
				retdata.SplitAmount = proto.Int32(int32(0))
				retdata.SuccessCount = proto.Int32(int32(0))
				retdata.SuccessToday = proto.Int32(int32(0))
				retdata.SuccessLastday = proto.Int32(int32(0))
				retdata.QrURL = proto.String("")
				retdata.Invited = proto.String("")
				retdata.UpperUserID = proto.Int32(0)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
