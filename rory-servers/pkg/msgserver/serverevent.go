package msgserver

import (
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// Names
const (
	EventClientConnect        = "sys_cltconnect"        // 客戶端連接
	EventClientDisconnect     = "sys_cltdisconnect"     // 客戶端斷開連接
	EventClientMessage        = "sys_cltmsg"            // 客戶端消息
	EventClientData           = "sys_cltdata"           // 客戶端數據
	EventRPCMessage           = "sys_rpcmessage"        // rpc消息
	EventNotifyMessage        = "sys_notifymessage"     // notify 消息
	EventAsyncCall            = "sys_asynccall"         // 異步調用消息
	EventHTTPMessage          = "sys_httpmsg"           // http消息
	ArgsNameClientInfo        = "sys_clt"               // 事件參數：客戶端信息
	ArgsNameClientData        = "sys_data"              // 事件參數,客戶端數據
	ArgsNameClientMsg         = "sys_msg"               // 事件參數,客戶端消息
	ArgsNameRPCMsg            = "sys_rpcmsg"            // 事件參數,rpc消息
	ArgsNameNotifyMsg         = "sys_notifymsg"         // 事件參數,notify消息
	ArgsNameAsyncCallBackFunc = "sys_asynccallbackfunc" // 事件參數,異步調用
	ArgsNameHTTPMsg           = "sys_httpmsg"           // 事件採納數,http消息
)

// HTTPMessageEventData
const (
	HTTPMessageEventDataSEMAPHORE          = "HTTPMSGDATA_SEMAPHORE"
	HTTPMessageEventDataHTTPRESPONSEWRITER = "HTTPMSGEVTDATA_HTTPRESPONSEWRITER"
	HTTPMessageEventDataHTTPREQUEST        = "HTTPMSGEVTDATA_HTTPREQUEST"
)

// NewClientConnectEvt 客戶端連接事件
func NewClientConnectEvt(clt *network.ConnInfo) *event.Event {
	args := event.NewEventArgs()
	args[ArgsNameClientInfo] = clt

	evt := event.NewEvent(EventClientConnect, args)

	return evt
}

// NewClientDisconnectEvt 客戶端斷開事件
func NewClientDisconnectEvt(clt *network.ConnInfo) *event.Event {
	args := event.NewEventArgs()
	args[ArgsNameClientInfo] = clt

	evt := event.NewEvent(EventClientDisconnect, args)

	return evt
}

// NewClientMessageEvent 客戶端消息事件
func NewClientMessageEvent(clt *network.ConnInfo, msg *msg.Message) *event.Event {
	args := event.NewEventArgs()
	args[ArgsNameClientMsg] = msg
	args[ArgsNameClientInfo] = clt

	evt := event.NewEvent(EventClientMessage, args)

	return evt
}

// NewClientDataEvent 客戶端數據事件
func NewClientDataEvent(clt *network.ConnInfo, data []byte) *event.Event {
	args := event.NewEventArgs()
	args[ArgsNameClientData] = data
	args[ArgsNameClientInfo] = clt

	evt := event.NewEvent(EventClientData, args)

	return evt
}

// NewRPCMessageEvent RPC消息事件
func NewRPCMessageEvent(msg *msg.Message) *event.Event {
	args := event.NewEventArgs()
	args[ArgsNameRPCMsg] = msg

	evt := event.NewEvent(EventRPCMessage, args)

	return evt
}

// NewNotifyMessageEvent 通知消息事件
func NewNotifyMessageEvent(msg *msg.Message) *event.Event {
	args := event.NewEventArgs()
	args[ArgsNameNotifyMsg] = msg

	evt := event.NewEvent(EventNotifyMessage, args)

	return evt
}

// NewAsyncCallEvent 異步調用消息事件
func NewAsyncCallEvent(callbackFunc func()) *event.Event {
	args := event.NewEventArgs()
	args[ArgsNameAsyncCallBackFunc] = callbackFunc

	evt := event.NewEvent(EventAsyncCall, args)

	return evt
}

// NewHTTPMsgEvent http消息事件
func NewHTTPMsgEvent(data map[string]interface{}) *event.Event {
	args := event.NewEventArgs()
	args[ArgsNameHTTPMsg] = data
	evt := event.NewEvent(EventHTTPMessage, args)

	return evt
}

// GetClientInfo 從事件參數中獲取客戶端信息
func GetClientInfo(args event.Args) *network.ConnInfo {
	return args[ArgsNameClientInfo].(*network.ConnInfo)
}

// GetClientData 從事件參數中獲取數據
func GetClientData(args event.Args) []byte {
	return args[ArgsNameClientData].([]byte)
}

// GetClientMsg 從事件參數中獲取消息
func GetClientMsg(args event.Args) *msg.Message {
	return args[ArgsNameClientMsg].(*msg.Message)
}

// GetRPCMsg 從事件參數中獲取RPC消息
func GetRPCMsg(args event.Args) *msg.Message {
	return args[ArgsNameRPCMsg].(*msg.Message)
}

// GetNotifyMsg 從事件參數中獲取Notify消息
func GetNotifyMsg(args event.Args) *msg.Message {
	return args[ArgsNameNotifyMsg].(*msg.Message)
}

// GetAsyncEventCallBackFunc 從事件參數中獲取異步回調函數
func GetAsyncEventCallBackFunc(args event.Args) func() {
	return args[ArgsNameAsyncCallBackFunc].(func())
}

// GetHTTPMsgData 從Http消息事件中獲取http參數
func GetHTTPMsgData(args event.Args) map[string]interface{} {
	return args[ArgsNameHTTPMsg].(map[string]interface{})
}
