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

// ApplyListHandler 審核代理申請列表(淨利潤)
type ApplyListHandler struct{}

// IsHook ?
func (rh *ApplyListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentApplyListID)
}

// OnMessage ?
func (rh *ApplyListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentApplyListMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentApplyListMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_NetProfitAgent_AgentApplyList", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentApplyListRetID), nil)
	retdata := new(netproto.NetProfitAgentApplyListRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_NetProfitAgent_AgentApplyList 错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.NetProfitAgentApplyListList)
					subRetData.ApplyUser = proto.Int32(int32(tb.GetValueByColName(i, "ApplyUser").(int64)))
					subRetData.NickName = proto.String(string(tb.GetValueByColName(i, "NickName").(string)))
					subRetData.TopUserID = proto.Int32(int32(tb.GetValueByColName(i, "TopUserID").(int64)))
					subRetData.TopUserName = proto.String(string(tb.GetValueByColName(i, "TopUserName").(string)))
					subRetData.CreateTime = proto.String(tb.GetValueByColName(i, "CreateTime").(time.Time).Format("2006-01-02 15:04:05"))
					subRetData.Status = proto.String(string(tb.GetValueByColName(i, "Status").(string)))
					subRetData.UpdateTime = proto.String(string(tb.GetValueByColName(i, "UpdateTime").(string)))
					// subRetData.UpdateTime = proto.String(tb.GetValueByColName(i, "UpdateTime").(time.Time).Format("2006-01-02 15:04:05"))
					retdata.ApplyListList = append(retdata.ApplyListList, subRetData)
				}
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
