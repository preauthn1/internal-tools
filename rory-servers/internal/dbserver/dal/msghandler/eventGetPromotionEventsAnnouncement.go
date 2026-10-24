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

// EventGetPromotionAnnouncementHandler 活動公告系統 郵件公告頁
type EventGetPromotionAnnouncementHandler struct{}

// IsHook ?
func (rh *EventGetPromotionAnnouncementHandler) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:EventGetPromotionAnnouncementHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_PromotionEventsAnnouncementID)
}

// OnMessage ?
func (rh *EventGetPromotionAnnouncementHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.PromotionAnnouncement{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.PromotionAnnouncement)
	ps := db.NewSQLParameters()
	ps.AddIntInput("AreaID", mm.GetArea_ID())

	proc := db.NewProcedure("PrGs_Event_GetPromotionEvents_Announcement", ps)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_PromotionEventsAnnouncementRetID), nil)
	retdata := new(netproto.PromotionAnnouncementRet)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrGs_Event_GetPromotionEvents_Announcement错误, %v", err)
		ctx.Error(strMsg)
		mlog.Error("执行存储过程PrGs_Event_GetPromotionEvents_Announcement错误, %s", err.Error())
		return []*msg.Message{retmsg}
	}

	if dbret.GetRetTableCount() > 0 {
		tb := dbret[0]
		if len(tb.Rows) > 0 {
			for i := 0; i < len(tb.Rows); i++ {
				pdata := new(netproto.PromotionAnnouncementContent)
				pdata.EventID = proto.Int32(int32(tb.GetValueByColName(i, "EventID").(int64)))
				pdata.Title = proto.String(tb.GetValueByColName(i, "Title").(string))
				pdata.EventStartTime = proto.String(tb.GetValueByColName(i, "EventStartTime").(time.Time).Format("2006-01-02 15:04:05"))
				pdata.EventEndTime = proto.String(tb.GetValueByColName(i, "EventEndTime").(time.Time).Format("2006-01-02 15:04:05"))
				pdata.ImageUrl = proto.String(tb.GetValueByColName(i, "ImageUrl").(string))
				pdata.Sort = proto.Int32(int32(tb.GetValueByColName(i, "Sort").(int64)))
				pdata.Tags = proto.Int32(int32(tb.GetValueByColName(i, "Tags").(int64)))
				pdata.DisplayMode = proto.Int32(int32(tb.GetValueByColName(i, "DisplayMode").(int64)))
				pdata.Content = proto.String(tb.GetValueByColName(i, "Content").(string))
				retdata.AnnouncementList = append(retdata.AnnouncementList, pdata)
			}
			retdata.Area_ID = proto.Int32(int32(tb.GetValueByColName(0, "AreaID").(int64)))
		}
	}

	// fmt.Println(retdata.GetArea_ID())
	// fmt.Println(retdata.GetAnnouncementList())
	mlog.System("announcement %v", retdata.GetAnnouncementList())

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
