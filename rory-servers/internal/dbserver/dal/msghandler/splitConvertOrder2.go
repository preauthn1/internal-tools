package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// SplitConvertOrderHandler2 裂變5提現紀錄送出提現金額
type SplitConvertOrderHandler2 struct{}

// IsHook ?
func (rh *SplitConvertOrderHandler2) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:SplitConvertOrderHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_SplitConvertOrderID)
}

// OnMessage ?
func (rh *SplitConvertOrderHandler2) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:SplitConvertOrderHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SplitConvertOrderMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata := new(netproto.SplitConvertOrderRet) // 回傳封包
	retdata.ChvMsg = proto.String("123456")

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitConvertOrderRetID), nil)
	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
