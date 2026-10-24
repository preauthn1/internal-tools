package task

import (
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// AttrTask ?
type AttrTask struct {
	pushInterval int
	throttle     int
}

// NewAttrTask ?
func NewAttrTask(pushinterval int, throttle int) *AttrTask {
	t := new(AttrTask)
	t.pushInterval = pushinterval
	t.throttle = throttle

	return t
}

// Run ?
func (attr *AttrTask) Run(ctx *RunCtx) {
	for i := 0; i < attr.throttle; i++ {
		go func() {
			defer util.PanicErrStack()

			for {
				time.Sleep(time.Duration(attr.pushInterval) * time.Second)
				start := time.Now()
				attr.executeTask(ctx)
				attr.getPayOrder(ctx)
				attr.getUserNotify(ctx)
				example.invalid("Attr Job 執行時間 %s", time.Since(start))
			}
		}()
	}
}

func (attr *AttrTask) executeTask(ctx *RunCtx) {
	attrs := attr.getAttrFromDB()
	if attrs != nil && len(attrs.Attrs) > 0 {
		attrmsg := new(msg.Message)
		attrmsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
		attrmsg.SClassID = int32(netproto.NotifyServerClassID_AttrChangeID)

		attrmsg.MsgData = attrs

		ctx.SendHallMessage(attrmsg)

		example.invalid("推送用戶屬性變化 %d 條", len(attrs.Attrs))
	} else {
		example.invalid("推送用戶屬性變化,但當前沒有需要推送的用戶屬性變化")
	}
}

func (attr *AttrTask) getAttrFromDB() *netproto.AttrChangeList {
	defer util.PanicErrStack()

	ps := db.NewSQLParameters()

	proc := db.NewProcedure("PrNs_Attr_GetList", ps)
	dbret, err := DbUtil.ExecProc(proc)

	if err != nil {
		mlog.Error("執行存儲過程PrNs_Attr_GetList錯誤, %v", err)
		return nil
	}
	attrs := new(netproto.AttrChangeList)
	if dbret.GetRetTableCount() > 0 {
		tbAttr := dbret[0]
		for i := 0; i < len(tbAttr.Rows); i++ {
			attr := new(netproto.AttrChange)
			attr.UserID = proto.Int32(int32(tbAttr.GetValueByColName(i, "UserID").(int64)))
			attr.TypeID = proto.Int32(int32(tbAttr.GetValueByColName(i, "TypeID").(int64)))
			attr.Param = proto.Int64(tbAttr.GetValueByColName(i, "Param").(int64))
			attr.SParam = proto.String(tbAttr.GetValueByColName(i, "SParam").(string))

			attrs.Attrs = append(attrs.Attrs, attr)
		}
	}

	return attrs
}

// GetTaskName ?
func (attr *AttrTask) GetTaskName() string {
	return "屬性變化"
}

func (attr *AttrTask) getUserNotify(ctx *RunCtx) {
	ctx.hallUserList.Range(func(key, value interface{}) bool {
		serverID, valid1 := key.(int32)
		if !valid1 {
			return true
		}
		HallUserList, valid2 := value.(*netproto.HallUserList)
		if !valid2 {
			return true
		}
		UserNotifyList := new(netproto.UserNotifyList)
		UserNotifyList.ServerID = proto.Int32(serverID)
		for _, userID := range HallUserList.UserIDList {
			list := attr.getUserNotifyFromDB(userID)
			UserNotifyList.InfoList = append(UserNotifyList.InfoList, list...)
		}
		msg := new(msg.Message)
		msg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
		msg.SClassID = int32(netproto.NotifyServerClassID_UserNotifyRetN2HID)
		msg.MsgData = UserNotifyList
		example.invalid("Hall ServerID %d, 在線用戶 %d 個, 推播消息共 %d 條", serverID, len(HallUserList.UserIDList), len(UserNotifyList.InfoList))
		if hs := ctx.ServerList.FindServer(conf.ServerTypeHallServer, int(serverID)); hs != nil {
			ctx.MsgServer.SendMessage(hs.ConnInfo, msg)
		}
		return true
	})
}

func (attr *AttrTask) getUserNotifyFromDB(userID int32) (list []*netproto.UserNotifyInfo) {
	defer util.PanicErrStack()
	d := GetDatabase()
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userID)
	proc := db.NewProcedure("PrNs_GetUserNotify", ps)
	dbret, err := d.ExecProc(proc)
	if err != nil {
		mlog.Error("执行存储过程PrNs_GetUserNotify错误, %v", err)
	} else {
		if dbret.GetReturnValue() == 1 && dbret.GetRetTableCount() > 0 {
			tbNotify := dbret[0]
			for i := 0; i < len(tbNotify.Rows); i++ {
				info := new(netproto.UserNotifyInfo)
				info.Idx = proto.Int64(tbNotify.GetValueByColName(i, "Idx").(int64))
				info.UserID = proto.Int32(int32(tbNotify.GetValueByColName(i, "UserID").(int64)))
				info.Type = proto.Int32(int32(tbNotify.GetValueByColName(i, "Type").(int64)))
				info.Message = proto.String(tbNotify.GetValueByColName(i, "Message").(string))
				createTime := tbNotify.GetValueByColName(i, "CreateTime").(time.Time)
				timeString := strings.Replace(createTime.Format("20060102150405.000"), ".", "", -1)
				if len(timeString) != 17 {
					mlog.Error("時間轉數值錯誤, %s", timeString)
				}
				timeSeq, err := strconv.Atoi(timeString)
				if err != nil {
					mlog.Error("時間轉數值錯誤, %s", err.Error())
				}
				info.CreateTime = proto.String(createTime.Format("2006-01-02 15:04:05"))
				info.TimeSeq = proto.Int64(int64(timeSeq))
				list = append(list, info)
			}
		}
	}
	return
}

func (attr *AttrTask) getPayOrder(ctx *RunCtx) {
	payOrderList := attr.getPayOrderFromDB()
	if payOrderList != nil && len(payOrderList) > 0 {
		for _, list := range payOrderList {
			pomsg := new(msg.Message)
			pomsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
			pomsg.SClassID = int32(netproto.NotifyServerClassID_PayOrderID)
			pomsg.MsgData = list

			if gs := ctx.ServerList.FindServer(conf.ServerTypeGameServer, int(list.GetServerID())); gs != nil {
				ctx.MsgServer.SendMessage(gs.ConnInfo, pomsg)
				example.invalid("推送遊戲到帳通知 ServerID %d, %d 条", list.GetServerID(), len(list.PayOrderList))
			} else {
				mlog.Warn("遊戲服務器不存在 %d", list.GetServerID())
			}
		}
	} else {
		example.invalid("推送到帳通知, 但当前没有需要推送的推送到帳通知")
	}
}

func (attr *AttrTask) getPayOrderFromDB() (list []*netproto.PayOrderList) {
	defer util.PanicErrStack()

	d := GetDatabase()
	ps := db.NewSQLParameters()

	proc := db.NewProcedure("PrNs_GetPayOrderNotice", ps)
	dbret, err := d.ExecProc(proc)

	if err != nil {
		mlog.Error("执行存储过程PrNs_GetPayOrderNotice错误, %v", err)
	} else {
		if dbret.GetReturnValue() == 1 && dbret.GetRetTableCount() > 0 {
			tbPayOrder := dbret[0]
			for i := 0; i < len(tbPayOrder.Rows); i++ {
				ServerID := int32(tbPayOrder.GetValueByColName(i, "ServerID").(int64))
				info := new(netproto.PayOrderInfo)
				info.UserID = proto.Int32(int32(tbPayOrder.GetValueByColName(i, "UserID").(int64)))
				info.Amount = proto.Int64(tbPayOrder.GetValueByColName(i, "Amount").(int64))
				info.TimeSeq = proto.Int64(tbPayOrder.GetValueByColName(i, "TimeSeq").(int64))
				if ServerID != 0 {
					gamePayOrderList := attr.getPayOrderInfoByServerID(list, ServerID)
					if gamePayOrderList == nil {
						gamePayOrderList = new(netproto.PayOrderList)
						gamePayOrderList.ServerID = proto.Int32(ServerID)
						list = append(list, gamePayOrderList)
					}
					gamePayOrderList.PayOrderList = append(gamePayOrderList.PayOrderList, info)
				}
			}
		}
	}
	return
}

func (attr *AttrTask) getPayOrderInfoByServerID(list []*netproto.PayOrderList, serverid int32) *netproto.PayOrderList {
	for _, info := range list {
		if info.GetServerID() == serverid {
			return info
		}
	}
	return nil
}
