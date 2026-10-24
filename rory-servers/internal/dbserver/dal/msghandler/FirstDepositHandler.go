package msghandler

// import (
// 	"example.invalid/golang/protobuf/proto"
// )

// import (
// 	"example.invalid/ServerTeam/Servers/pkg/db"
// 	"example.invalid/ServerTeam/Servers/pkg/module"
// 	"example.invalid/ServerTeam/Servers/pkg/msg"
// 	"example.invalid/ServerTeam/Servers/pkg/netproto"
// 	"example.invalid/ServerTeam/Servers/pkg/network"
// 	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
// )

// import (
// 	"time"
// )

// // FirstDepositHandler 首充返利活動
// type FirstDepositHandler struct{}

// // IsHook 是否有對應到這隻handler
// func (h *FirstDepositHandler) IsHook(bclsid int32, sclsid int32) bool {
// 	if bclsid != int32(netproto.MessageBClassID_Hall) {
// 		return false
// 	}
// 	return sclsid == int32(netproto.HallMsgClassID_FirstDepositGetInfoID) || sclsid == int32(netproto.HallMsgClassID_FirstDepositCollectID)
// }

// // OnMessage 處理
// func (h *FirstDepositHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

// 	switch netproto.HallMsgClassID(m.SClassID) {
// 	case netproto.HallMsgClassID_FirstDepositGetInfoID:
// 		return h.GetInfo(ctx, clt, m)
// 	case netproto.HallMsgClassID_FirstDepositCollectID:
// 		return h.Collect(ctx, clt, m)
// 	}

// 	return []*msg.Message{}
// }

// // GetInfo 取得首充返利活動與該用戶相關資料
// func (h *FirstDepositHandler) GetInfo(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

// 	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
// 		mm := &netproto.FirstDepositGetInfoMsg{}
// 		return mm
// 	})

// 	if rmm == nil {
// 		return []*msg.Message{GetDBRetMsg(0)}
// 	}

// 	mm := rmm.(*netproto.FirstDepositGetInfoMsg)

// 	const ProcedureName = "PrGs_Event_FirstDeposit_GetInfo"

// 	ps := db.NewSqlParameters()

// 	ps.AddIntInput("intUserID", mm.GetUserID()) // 使用者ID

// 	proc := db.NewProcedure(ProcedureName, ps)

// 	retMsg := new(msg.Message)
// 	retMsg.BClassID = int32(netproto.MessageBClassID_Hall)
// 	retMsg.SClassID = int32(netproto.HallMsgClassID_FirstDepositGetInfoRetID)

// 	data := new(netproto.FirstDepositGetInfoRetMsg)

// 	dbret, err := DbUtil.ExecProc(proc)
// 	if err != nil {
// 		data.RetCode = proto.Int32(-1)
// 		retMsg.MsgData = data

// 		return []*msg.Message{retMsg}
// 	}

// 	if dbret.GetRetTableCount() != 2 {
// 		data.RetCode = proto.Int32(-1)
// 		retMsg.MsgData = data

// 		return []*msg.Message{retMsg}
// 	}

// 	if len(dbret[0].Rows) > 0 {
// 		info := new(netproto.FirstDepositInfo)

// 		var StartTime string
// 		var EndTime string
// 		var DepositLimit int64
// 		var WinAmount int64
// 		var ReturnAmount int64
// 		var Status int32
// 		var CurrentUserCount int32

// 		if val, ok := dbret[0].GetValueByColName(0, "StartTime").(time.Time); ok {
// 			StartTime = val.Format("2006-01-02 15:04:05")
// 		}

// 		if val, ok := dbret[0].GetValueByColName(0, "EndTime").(time.Time); ok {
// 			EndTime = val.Format("2006-01-02 15:04:05")
// 		}

// 		if val, ok := dbret[0].GetValueByColName(0, "DepositLimit").(int64); ok {
// 			DepositLimit = val
// 		}

// 		if val, ok := dbret[0].GetValueByColName(0, "WinAmount").(int64); ok {
// 			WinAmount = val
// 		}

// 		if val, ok := dbret[0].GetValueByColName(0, "ReturnAmount").(int64); ok {
// 			ReturnAmount = val
// 		}

// 		if val, ok := dbret[0].GetValueByColName(0, "Status").(int64); ok {
// 			Status = int32(val)
// 		}

// 		if val, ok := dbret[0].GetValueByColName(0, "CurrentUserCount").(int64); ok {
// 			CurrentUserCount = int32(val)
// 		}

// 		info.StartTime = proto.String(StartTime)
// 		info.EndTime = proto.String(EndTime)
// 		info.DepositLimit = proto.Int64(DepositLimit)
// 		info.WinAmount = proto.Int64(WinAmount)
// 		info.ReturnAmount = proto.Int64(ReturnAmount)
// 		info.Status = proto.Int32(Status)
// 		info.CurrentUserCount = proto.Int32(CurrentUserCount)

// 		data.FirstDepositInfo = info
// 	}

// 	if len(dbret[1].Rows) > 0 {
// 		record := new(netproto.FirstDepositUserRecord)

// 		var UserID int32
// 		var EventID int64
// 		var JoinTime string
// 		var DepositAmount int64
// 		var WinAmount int64
// 		var Flag int32
// 		var CollectTime string
// 		var CollectAmount int64
// 		var UserType int32

// 		if val, ok := dbret[1].GetValueByColName(0, "UserID").(int64); ok {
// 			UserID = int32(val)
// 		}

// 		if val, ok := dbret[1].GetValueByColName(0, "EventID").(int64); ok {
// 			EventID = val
// 		}

// 		if val, ok := dbret[1].GetValueByColName(0, "JoinTime").(time.Time); ok {
// 			JoinTime = val.Format("2006-01-02 15:04:05")
// 		}

// 		if val, ok := dbret[1].GetValueByColName(0, "DepositAmount").(int64); ok {
// 			DepositAmount = val
// 		}

// 		if val, ok := dbret[1].GetValueByColName(0, "WinAmount").(int64); ok {
// 			WinAmount = val
// 		}

// 		if val, ok := dbret[1].GetValueByColName(0, "Flag").(int64); ok {
// 			Flag = int32(val)
// 		}

// 		if val, ok := dbret[1].GetValueByColName(0, "CollectTime").(time.Time); ok {
// 			CollectTime = val.Format("2006-01-02 15:04:05")
// 		}

// 		if val, ok := dbret[1].GetValueByColName(0, "CollectAmount").(int64); ok {
// 			CollectAmount = val
// 		}

// 		if val, ok := dbret[1].GetValueByColName(0, "UserType").(int64); ok {
// 			UserType = int32(val)
// 		}

// 		record.UserID = proto.Int32(UserID)
// 		record.EventID = proto.Int64(EventID)
// 		record.JoinTime = proto.String(JoinTime)
// 		record.DepositAmount = proto.Int64(DepositAmount)
// 		record.WinAmount = proto.Int64(WinAmount)
// 		record.Flag = proto.Int32(Flag)
// 		record.CollectTime = proto.String(CollectTime)
// 		record.CollectAmount = proto.Int64(CollectAmount)
// 		record.UserType = proto.Int32(UserType)

// 		data.FirstDepositUserRecord = record
// 	}

// 	retMsg.MsgData = data

// 	return []*msg.Message{retMsg}
// }

// // Collect 領取首充返利回饋金
// func (h *FirstDepositHandler) Collect(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

// 	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
// 		mm := &netproto.FirstDepositCollectMsg{}
// 		return mm
// 	})

// 	if rmm == nil {
// 		return []*msg.Message{GetDBRetMsg(0)}
// 	}

// 	mm := rmm.(*netproto.FirstDepositCollectMsg)

// 	const ProcedureName = "PrGs_Event_FirstDeposit_Collect"

// 	ps := db.NewSqlParameters()

// 	ps.AddIntInput("intUserID", mm.GetUserID())          // 使用者ID
// 	ps.AddIntInput("intEventID", int32(mm.GetEventID())) //  對應Table 1的ID
// 	if mm.GetIPAddress() != "" {
// 		ps.AddVarcharInput("chvIPAddress", mm.GetIPAddress(), 40) // 使用者IP
// 	}
// 	ps.AddVarcharOutput("chvErrMsg", "")

// 	proc := db.NewProcedure(ProcedureName, ps)

// 	retMsg := new(msg.Message)
// 	retMsg.BClassID = int32(netproto.MessageBClassID_Hall)
// 	retMsg.SClassID = int32(netproto.HallMsgClassID_FirstDepositCollectRetID)

// 	data := new(netproto.FirstDepositCollectRetMsg)

// 	dbret, err := DbUtil.ExecProc(proc)
// 	if err != nil {
// 		data.RetCode = proto.Int32(-1)
// 		retMsg.MsgData = data
// 		data.RetMessage = proto.String("执行存储过程错误")

// 		return []*msg.Message{retMsg}
// 	}

// 	data.RetCode = proto.Int32(int32(dbret.GetReturnValue()))
// 	data.RetMessage = proto.String(dbret.GetOutputParamValueString("chvErrMsg"))

// 	retMsg.MsgData = data

// 	return []*msg.Message{retMsg}
// }
