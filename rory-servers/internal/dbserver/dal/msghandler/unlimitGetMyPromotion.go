package msghandler

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetMyPromotionHandler 無限級QRCODE
type GetMyPromotionHandler struct{}

// IsHook ?
func (rh *GetMyPromotionHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:GetMyPromotionHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetMyPromotionID)
}

// OnMessage ?
func (rh *GetMyPromotionHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:GetMyPromotionHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetMyPromotionMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetMyPromotionMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數，須和資料庫輸出欄位相同
	ps.AddIntInput("intUserID", mm.GetUserID()) //使用者ID

	//開始準備SP字串
	proc := db.NewProcedure("PrPs_Unlimit_GetMyPromotion", ps)
	mlog.System("HANDLER名稱:GetMyPromotionHandler,SP字串:%v", proc.GetSQLString()) // DEBUG專用
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	//資料庫回傳值
	retcode := int32(dbret.GetReturnValue()) //retcode必須為1
	// mlog.System("HANDLER名稱:GetMyPromotionHandler,retcode:%v", retcode) // DEBUG專用

	// 準備回傳給CLIENT的封包
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMyPromotionRetID), nil)
	retdata := new(netproto.GetMyPromotionRet) // 回傳封包

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程PrPs_PrPs_Unlimit_GetMyPromotion 错误 %v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 { // 準備回傳資料
			// 第一張表
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.BossID = proto.Int32(int32(tb.GetValueByColName(0, "BossID").(int64)))
				retdata.TeamMember = proto.Int32(int32(tb.GetValueByColName(0, "TeamMember").(int64)))
				retdata.MyMember = proto.Int32(int32(tb.GetValueByColName(0, "MyMember").(int64)))
				retdata.TodayNewMember = proto.Int32(int32(tb.GetValueByColName(0, "TodayNewMember").(int64)))
				retdata.TodayActiveMember = proto.Int32(int32(tb.GetValueByColName(0, "TodayActiveMember").(int64)))
				retdata.QrURL = proto.String(string(tb.GetValueByColName(0, "QrURL").(string)))
				retdata.InviteCode = proto.String(string(tb.GetValueByColName(0, "invitecode").(string)))
			} else {
				retdata.BossID = proto.Int32(int32(0))
				retdata.TeamMember = proto.Int32(int32(0))
				retdata.MyMember = proto.Int32(int32(0))
				retdata.TodayNewMember = proto.Int32(int32(0))
				retdata.TodayActiveMember = proto.Int32(int32(0))
				retdata.QrURL = proto.String("")
				retdata.InviteCode = proto.String("")
			}
		}
	}

	currentSystemTime := time.Now()
	currentSystemTimeFormed := currentSystemTime.Format("2006-01-02")
	retdata.SystemTime = proto.String(currentSystemTimeFormed)

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
