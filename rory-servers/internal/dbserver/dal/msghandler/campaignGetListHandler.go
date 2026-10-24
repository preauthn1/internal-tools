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

// // 大廳用戶活動介面
// type CampaignGetListHandler struct{}

// func (fsp *CampaignGetListHandler) IsHook(bclsid int32, sclsid int32) bool {
// 	return bclsid == int32(netproto.MessageBClassID_Hall) && (sclsid == int32(netproto.HallMsgClassID_CampaignGetListID))
// }

// func (fsp *CampaignGetListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
// 	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
// 		mm := &netproto.CampaignGetList{}
// 		return mm
// 	})

// 	if rmm == nil {
// 		return []*msg.Message{GetDBRetMsg(0)}
// 	}

// 	mm := rmm.(*netproto.CampaignGetList)
// 	ps := db.NewSqlParameters()
// 	ps.AddIntInput("intUserID", mm.GetUserID())
// 	proc := db.NewProcedure("PrGs_Campaign_GetList", ps)
// 	ret, err := DbUtil.ExecProc(proc)
// 	defer func() {
// 		if err := recover(); err != nil {
// 			errorStr := dbutil.GetStackInfo()
// 			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
// 			panic(err)
// 		}
// 	}()
// 	//retcode := int32(ret.GetReturnValue())

// 	if err != nil {
// 		ctx.Error("执行存储过程[PrGs_Campaign_GetList]错误, %v, %s", err, proc.GetSQLString())
// 		//retcode = int32(0)
// 	}

// 	var retmsg *msg.Message
// 	retmsg = msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CampaignGetListRetID), nil)

// 	//get CampaignGetListRet
// 	msgData := new(netproto.CampaignGetListRet)
// 	detailMsg := make([]*netproto.CampaignData, 0)

// 	if ret.GetRetTableCount() > 0 {
// 		tbData := ret[0]
// 		if ret.GetRetTableCount() > 0 {
// 			if len(tbData.Rows) > 0 {
// 				for i := 0; i < len(tbData.Rows); i++ {
// 					detail := new(netproto.CampaignData)
// 					detail.ID = proto.Int32(tbData.GetValueByColNameInt32(i, "Id"))
// 					detail.Picture = proto.String(tbData.GetValueByColNameString(i, "Picture"))
// 					detail.SeqID = proto.Int32(tbData.GetValueByColNameInt32(i, "SeqID"))
// 					detail.Title = proto.String(tbData.GetValueByColNameString(i, "Title"))
// 					detail.CampaignTag = proto.Int32(tbData.GetValueByColNameInt32(i, "CampaignTag"))
// 					detail.CampaignBehavior = proto.Int32(tbData.GetValueByColNameInt32(i, "CampaignBehavior"))
// 					detail.Link = proto.String(tbData.GetValueByColNameString(i, "Link"))
// 					detail.StartTime = proto.String(tbData.GetValueByColNameTime(i, "StartTime").Format("2006-01-02 15:04:05"))
// 					detail.EndTime = proto.String(tbData.GetValueByColNameTime(i, "EndTime").Format("2006-01-02 15:04:05"))
// 					detailMsg = append(detailMsg, detail)
// 				}
// 			}
// 		}
// 	}
// 	msgData.CampaignData = detailMsg
// 	retmsg.MsgData = msgData
// 	return []*msg.Message{retmsg}
// }
