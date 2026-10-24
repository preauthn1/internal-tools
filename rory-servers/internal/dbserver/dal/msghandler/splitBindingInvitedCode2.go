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

// SplitBindingInvitedCodeHandler2 裂變6輸入推廣碼
type SplitBindingInvitedCodeHandler2 struct{}

// IsHook ?
func (rh *SplitBindingInvitedCodeHandler2) IsHook(bclsid int32, sclsid int32) bool {
	// mlog.System("HANDLER名稱:SplitBindingInvitedCodeHandler,接收HOOK") // DEBUG專用
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_SplitBindingInvitedCodeID)
}

// OnMessage ?
func (rh *SplitBindingInvitedCodeHandler2) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	// mlog.System("HANDLER名稱:SplitBindingInvitedCodeHandler,進入OnMessage") // DEBUG專用

	// 從CLIENT接收傳入的封包，並讀取封包裡面的訊息
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SplitBindingInvitedCodeMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata := new(netproto.SplitBindingInvitedCodeRet) // 回傳封包
	retdata.ChvMsg = proto.String("假資料測試成功")

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_SplitBindingInvitedCodeRetID), nil)
	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
