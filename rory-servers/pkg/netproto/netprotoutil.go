package netproto

import proto "example.invalid/golang/protobuf/proto"

//创建空的RPC信息
func CreateEmptyRPCInfo(serverid int32) *RPCInfo {
	rpcinfo := new(RPCInfo)
	rpcinfo.ConnID = proto.String("")
	rpcinfo.IPAddress = proto.String("local")
	rpcinfo.Cer = proto.String("")
	rpcinfo.RouteServerID = proto.Int32(serverid)

	return rpcinfo
}
