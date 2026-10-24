package rpcmsg

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// ParseRPCMessage 解析RPC消息
func ParseRPCMessage(m *msg.Message, umf func() proto.Message) (*netproto.RPCInfo, proto.Message) {
	if m == nil {
		mlog.Error("傳入的 Message 為 nil")
		return nil, nil
	}

	rmm, ok := m.MsgData.(*RPCMessage)
	if !ok {
		mlog.Error("傳入的 Message 無法轉型為 RPCMessage")
		return nil, nil
	}

	mmbyte := rmm.MsgData.([]byte)
	mm := umf()
	// umf() 回傳 nil 就不需要轉換 MsgData 的資料了
	if mm == nil {
		return rmm.RPCInfo, nil
	}
	err := proto.Unmarshal(mmbyte, mm)

	if err != nil {
		//mlog.NewDefaultLogger().SetData(fmt.Sprintf("[msg.Message] 解析消息出现错误,%v 要專出的型別[%v] 請儘速通知資訊人員", err, mm)).Error("HallServer壓測時")
		//mlog.Error("解析消息出现错误,%v, %v_%v", err, m.BClassID, m.SClassID)
		mlog.NewDefaultLogger().SetData(mm).Error("解析消息出現錯誤,%v bytes[%v], 要轉出的型別[%T], %v_%v", err, mmbyte, mm, m.BClassID, m.SClassID)
		return nil, nil
	}

	return rmm.RPCInfo, mm
}

// GetRPCInfo 獲取RPC信息
func GetRPCInfo(m *msg.Message) *netproto.RPCInfo {
	if m == nil {
		mlog.Error("傳入的 Message 為 nil")
		return nil
	}

	rmm, ok := m.MsgData.(*RPCMessage)
	if !ok {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[msg.Message] GetRPCInfo 獲取RPC信息錯誤,%v 要專出的型別[RPCMessage] 請儘速通知資訊人員", m)).Error("HallServer壓測時")
		return nil
	}

	return rmm.RPCInfo
}
