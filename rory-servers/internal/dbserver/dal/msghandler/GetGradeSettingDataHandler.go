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
// 	"example.invalid/ServerTeam/Servers/pkg/dbutil"
// )

// type GetGradeSettingDataHandler struct {
// }

// func (rh *GetGradeSettingDataHandler) IsHook(bclsid int32, sclsid int32) bool {
// 	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_GetGradeSettingDataID)
// }

// func (rh *GetGradeSettingDataHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

// 	ps := db.NewSqlParameters()
// 	ps.AddVarcharOutput("chvErrMsg", "")

// 	proc := db.NewProcedure("PrGs_User_GetGradeSetting", ps)

// 	ret, err := DbUtil.ExecProc(proc)

// 	defer func() {
// 		if err := recover(); err != nil {
// 			errorStr := dbutil.GetStackInfo()
// 			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
// 			panic(err)
// 		}
// 	}()

// 	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetGradeSettingDataRetID), nil)
// 	retdata := new(netproto.GradeSettingDataList)
// 	if err != nil {
// 		ctx.Error("执行存储过程PrGs_User_GetGradeSetting错误, %v, %s", err, proc.GetSQLString())
// 		retdata.RetCode = proto.Int32(0)
// 		retdata.ChvMsg = proto.String("服务器错误")
// 	} else {
// 		retdata.RetCode = proto.Int32(int32(ret.GetReturnValue()))
// 		retdata.ChvMsg = proto.String(ret.GetOutputParamValue("chvErrMsg").(string))
// 		for i := 0; i < len(ret[0].Rows); i++ {
// 			stGradeData := new(netproto.GradeSettingData)

// 			stGradeData.ID = proto.Int32(ret[0].GetValueByColNameInt32(i, "Id"))
// 			stGradeData.Amount = proto.Int64(ret[0].GetValueByColNameInt64(i, "Amount"))
// 			stGradeData.PromotionAward = proto.Int32(ret[0].GetValueByColNameInt32(i, "PromotionAward"))
// 			stGradeData.AvatarFrame = proto.Int32(ret[0].GetValueByColNameInt32(i, "AvatarFrame"))
// 			retdata.DataList = append(retdata.DataList, stGradeData)
// 		}
// 	}

// 	retmsg.MsgData = retdata
// 	return []*msg.Message{retmsg}
// }
