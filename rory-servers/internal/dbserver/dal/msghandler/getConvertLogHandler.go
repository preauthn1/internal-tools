package msghandler

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetConvertLogHandler 獲取用戶郵件列表
type GetConvertLogHandler struct{}

// IsHook ?
func (rh *GetConvertLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_GetConvertLogID)
}

// OnMessage ?
func (rh *GetConvertLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetConvertLog{}
		return mm
	})

	if rmm == nil || rpcinfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	now := time.Now()

	mm := rmm.(*netproto.GetConvertLog)
	loginBeginTime := now.Format("2006-01-02 15:04:05")
	loginEndTime := now.AddDate(0, 0, 1).Format("2006-01-02 15:04:05")
	tnyStatus := mm.GetTnyStatus()

	//1(一周),2(半個月),3(一個月),4(三個月)
	logRange := mm.GetLogRange()
	switch logRange {
	case 1:
		loginBeginTime = now.AddDate(0, 0, -7).Format("2006-01-02") + " 00:00:00"
		break
	case 2:
		loginBeginTime = now.AddDate(0, 0, -15).Format("2006-01-02") + " 00:00:00"
		break
	case 3:
		loginBeginTime = now.AddDate(0, -1, 0).Format("2006-01-02") + " 00:00:00"
		break
	case 4:
		loginBeginTime = now.AddDate(0, -3, 0).Format("2006-01-02") + " 00:00:00"
		break
	}

	userid := rpcinfo.GetUserID()

	ps := db.NewSQLParameters()
	ps.AddVarcharInput("logInBeginTime", loginBeginTime, 19)
	ps.AddVarcharInput("logInEndTime", loginEndTime, 19)
	// tnyStatus 99 不帶等於全部
	if tnyStatus != 99 {
		ps.AddIntInput("tnyStatus", tnyStatus)
	}
	ps.AddIntInput("intUserID", userid)
	proc := db.NewProcedure("PrPs_GetConvertLog", ps)
	//dbret, err := DbUtil.ExecProc(proc)
	var dbret db.RecordSet
	var err error
	dbret, err = dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetConvertLogRetID), nil)
	retdata := new(netproto.ConvertLogList)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_GetConvertLog错误, %v", err)
		ctx.Error(strMsg)
	} else {
		if dbret.GetRetTableCount() > 0 {
			tbData := dbret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				logDetail := new(netproto.ConvertLogDetail)

				logDetail.Amount = proto.Int64(tbData.GetValueByColName(i, "Amount").(int64))
				logDetail.TaxAmount = proto.Int64(tbData.GetValueByColName(i, "TaxAmount").(int64))
				logDetail.Status = proto.Int32(int32(tbData.GetValueByColName(i, "Status").(int64)))
				logDetail.Payment = proto.String(tbData.GetValueByColName(i, "Payment").(string))
				logDetail.CreateTime = proto.String(tbData.GetValueByColName(i, "CreateTime").(time.Time).Format("2006-01-02 15:04:05"))

				retdata.ConvertLogList = append(retdata.ConvertLogList, logDetail)
			}
		}
	}

	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
