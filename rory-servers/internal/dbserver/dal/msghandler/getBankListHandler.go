package msghandler

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetBankListHandler 獲取用戶郵件列表
type GetBankListHandler struct{}

// IsHook ?
func (rh *GetBankListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_GetBankListID)
}

// OnMessage ?
func (rh *GetBankListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo := rpcmsg.GetRPCInfo(m)

	if rpcinfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	ps := db.NewSQLParameters()

	proc := db.NewProcedure("PrPs_GetUnionpayBankCode", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetBankListRetID), nil)
	retdata := new(netproto.BankList)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_GetUnionpayBankCode错误, %v", err)
		ctx.Error(strMsg)
	} else {
		if dbret.GetRetTableCount() > 0 {
			tbData := dbret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				logDetail := new(netproto.BankDetail)

				logDetail.BankCode = proto.String(tbData.GetValueByColName(i, "Code").(string))
				logDetail.BankName = proto.String(tbData.GetValueByColName(i, "Bankname").(string))

				retdata.BankList = append(retdata.BankList, logDetail)
			}
		}
	}

	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
