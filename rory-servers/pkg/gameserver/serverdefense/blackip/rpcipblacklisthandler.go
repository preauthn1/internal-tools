package blackip

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// RPCIPBlackListHandler rpc ip黑名單處理器
type RPCIPBlackListHandler struct {
	msgqueue               *rpcmsg.Queue
	loadIPListCallBackFunc func(map[string]string)
	serverid               int32
	serverDesc             string
}

// NewRPCIPBlackListHandler 創建IP黑名單處理器對象
func NewRPCIPBlackListHandler(msgqueue *rpcmsg.Queue, svrid int32, serverDesc string) *RPCIPBlackListHandler {
	rpciplist := new(RPCIPBlackListHandler)

	rpciplist.msgqueue = msgqueue
	rpciplist.serverid = svrid
	rpciplist.serverDesc = serverDesc
	return rpciplist
}

// BlockIP 封鎖IP
func (rpciplist *RPCIPBlackListHandler) BlockIP(ip string, sec int, memo string) {
	memo = memo + "," + rpciplist.serverDesc

	blockmsg := new(msg.Message)
	blockmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	blockmsg.SClassID = int32(netproto.DBServerClassID_BlockIPID)

	blockdata := new(netproto.BlockIP)
	blockdata.IPAddress = proto.String(ip)
	blockdata.Sec = proto.Int32(int32(sec))
	blockdata.Memo = proto.String(memo)
	blockmsg.MsgData = blockdata

	rpciplist.msgqueue.PushMessage(netproto.CreateEmptyRPCInfo(rpciplist.serverid), false, blockmsg)
}

// LoadBlockIPList load
func (rpciplist *RPCIPBlackListHandler) LoadBlockIPList(callbackFunc func(map[string]string)) {
	rpciplist.loadIPListCallBackFunc = callbackFunc
	rpciplist.loadBlockIP()
}

// 加載IP黑名單
func (rpciplist *RPCIPBlackListHandler) loadBlockIP() {
	loadMsg := new(msg.Message)
	loadMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	loadMsg.SClassID = int32(netproto.DBServerClassID_LoadBlockIPID)
	loadMsg.MsgData = nil

	rpciplist.msgqueue.PushMessage(netproto.CreateEmptyRPCInfo(rpciplist.serverid), false, loadMsg)
}

// OnRPCMessage rpc
func (rpciplist *RPCIPBlackListHandler) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.BlockIPList{}
		return mm
	})

	mm := rmm.(*netproto.BlockIPList)

	if mm == nil {
		return
	}

	iplst := make(map[string]string)

	for _, sip := range mm.GetIPAddress() {
		iplst[sip] = sip
	}

	if rpciplist.loadIPListCallBackFunc != nil {
		rpciplist.loadIPListCallBackFunc(iplst)
	}
}

// IsHook hook
func (rpciplist *RPCIPBlackListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_BlockIPListID)
}
