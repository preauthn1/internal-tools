package msghandler

// import (
// 	"example.invalid/golang/protobuf/proto"
// )

// import (
// 	"example.invalid/ServerTeam/Servers/pkg/db"
// 	"example.invalid/ServerTeam/Servers/pkg/mlog"
// 	"example.invalid/ServerTeam/Servers/pkg/module"
// 	"example.invalid/ServerTeam/Servers/pkg/msg"
// 	"example.invalid/ServerTeam/Servers/pkg/netproto"
// 	"example.invalid/ServerTeam/Servers/pkg/network"
// 	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
// 	"example.invalid/ServerTeam/Servers/pkg/dbutil"
// )

// // 修改VIP頭像框ID
// type ModifyAvatarFrameHandler struct {
// }

// func (rh *ModifyAvatarFrameHandler) IsHook(bclsid int32, sclsid int32) bool {
// 	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_ModifyAvatarFrameID)
// }

// func (rh *ModifyAvatarFrameHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
// 	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
// 		mm := &netproto.ModifyAvatarFrameMsg{}
// 		return mm
// 	})

// 	if rmm == nil || rpcinfo == nil {
// 		return []*msg.Message{GetDBRetMsg(0)}
// 	}

// 	mm := rmm.(*netproto.ModifyAvatarFrameMsg)

// 	i32UserID := mm.GetUserID()
// 	i32AvatarFrameID := mm.GetAvatarFrameID()

// 	ps := db.NewSqlParameters()
// 	ps.AddIntInput("intUserID", i32UserID)
// 	ps.AddIntInput("intAvatarFrame", i32AvatarFrameID)
// 	ps.AddVarcharOutput("chvErrMsg", "")

// 	proc := db.NewProcedure("PrGs_User_ModifyAvatarFrame", ps)
// 	var ret db.RecordSet
// 	var err error
// 	ret, err = DbUtil.ExecProc(proc)

// 	defer func() {
// 		if err := recover(); err != nil {
// 			errorStr := dbutil.GetStackInfo()
// 			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
// 			panic(err)
// 		}
// 	}()

// 	var i32RetCode int32 = 0
// 	chvRetMessage := ""
// 	if err != nil {
// 		ctx.Error("执行存储过程PrGs_User_ModifyAvatarFrame错误, %v", err)
// 		i32RetCode = 0
// 		chvRetMessage = "服务器错误"
// 	} else {
// 		i32RetCode = int32(ret.GetReturnValue())
// 		chvRetMessage = ret.GetOutputParamValue("chvErrMsg").(string)
// 	}

// 	retmsg := new(msg.Message)
// 	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
// 	retmsg.SClassID = int32(netproto.HallMsgClassID_ModifyAvatarFrameRetID)
// 	retmsg.MsgData = &netproto.ModifyAvatarFrameRetMsg{
// 		RetCode: proto.Int32(i32RetCode),
// 		ChvMsg:  proto.String(chvRetMessage),
// 	}

// 	return []*msg.Message{retmsg}
// }
