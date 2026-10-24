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

// EventGetPromotionPopupHandler 活動公告系統 首頁彈出
type EventGetPromotionPopupHandler struct{}

// IsHook ?
func (rh *EventGetPromotionPopupHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:EventGetPromotionPopupHandler,接收HOOK") // DEBU專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_PromotionEventsPopupID)
}

// OnMessage ?
func (rh *EventGetPromotionPopupHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.PromotionPopup{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.PromotionPopup)
	ps := db.NewSQLParameters()
	ps.AddIntInput("AreaID", mm.GetArea_ID())

	proc := db.NewProcedure("PrGs_Event_GetPromotionEvents_Popup", ps)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_PromotionEventsPopupRetID), nil)
	retdata := new(netproto.PromotionPopupRet)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		strMsg := fmt.Sprintf("執行存儲過程PrGs_Event_GetPromotionEvents_Popup錯誤, %v", err)
		ctx.Error(strMsg)
		mlog.Error("執行存儲過程PrGs_Event_GetPromotionEvents_Popup錯誤, %s", err.Error())
		return []*msg.Message{retmsg}
	}

	if dbret.GetRetTableCount() > 0 {
		tb := dbret[0]
		if len(tb.Rows) > 0 {
			for i := 0; i < len(tb.Rows); i++ {
				pdata := new(netproto.PromotionPopupContent)
				pdata.EventID = proto.Int32(int32(tb.GetValueByColName(i, "EventID").(int64)))
				pdata.Title = proto.String(tb.GetValueByColName(i, "Title").(string))
				pdata.EventStartTime = proto.String(tb.GetValueByColName(i, "EventStartTime").(time.Time).Format("2006-01-02 15:04:05"))
				pdata.EventEndTime = proto.String(tb.GetValueByColName(i, "EventEndTime").(time.Time).Format("2006-01-02 15:04:05"))
				pdata.ImageUrl = proto.String(tb.GetValueByColName(i, "ImageUrl").(string))
				pdata.Sort = proto.Int32(int32(tb.GetValueByColName(i, "Sort").(int64)))
				pdata.BehaviorType = proto.Int32(int32(tb.GetValueByColName(i, "BehaviorType").(int64)))
				pdata.BehaviorSettingValue = proto.String(tb.GetValueByColName(i, "BehaviorSettingValue").(string))
				retdata.PopupList = append(retdata.PopupList, pdata)

				// fmt.Println(retdata)
			}
			retdata.Area_ID = proto.Int32(int32(tb.GetValueByColName(0, "AreaID").(int64)))
		}
	}

	// fmt.Println(retdata.GetArea_ID())
	// fmt.Println(retdata.GetPopupList())

	// mlog.System("popup %v", retdata.GetPopupList())

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
