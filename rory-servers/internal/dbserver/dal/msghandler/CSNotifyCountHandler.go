package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// CSNotifyCountHandler 代理相關資訊
type CSNotifyCountHandler struct{}

// IsHook ?
func (rh *CSNotifyCountHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && (sclsid == int32(netproto.HallMsgClassID_CheckCSNotifyID))
}

// OnMessage ?
func (rh *CSNotifyCountHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo := rpcmsg.GetRPCInfo(m)
	userid := rpcinfo.GetUserID()
	retData := new(netproto.CSNotiftyCount)
	if m.SClassID == int32(netproto.HallMsgClassID_CheckCSNotifyID) { //查詢未讀客服消息數量
		ps := db.NewSQLParameters()
		ps.AddIntInput("intUserID", userid)
		proc := db.NewProcedure("PrNs_ServiceTalkNotifyCount", ps)
		ret, err := dbutil.DBUtil.ExecProc(proc)
		defer func() {
			if err := recover(); err != nil {
				errorStr := util.GetStackInfo()
				mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
				panic(err)
			}
		}()
		var notifyCount int32 = 0
		if err != nil {
			ctx.Error("执行存储过程PrNs_ServiceTalkNotifyCount错误, %v", err)
		} else {
			notifyCount = int32(ret.GetReturnValue())
		}
		retData.UserID = proto.Int32(userid)
		retData.Count = proto.Int32(notifyCount)
	}
	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retmsg.SClassID = int32(netproto.HallMsgClassID_CheckCSNotifyRetID)
	retmsg.MsgData = retData
	return []*msg.Message{retmsg}
}
