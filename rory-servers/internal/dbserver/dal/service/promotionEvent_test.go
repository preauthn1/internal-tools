package service

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"testing"
)

func TestPromotionEvent_QueryEventList(t *testing.T) {
	service := &PromotionEvent{}
	service.Base = newBase()
	rpcInfo := &netproto.RPCInfo{}
	rpcInfo.UserID = proto.Int(15267354)
	rpcInfo.Cer = proto.String("D12BFB7CC90444188BDF6795CBF9CE7A")
	rpcInfo.IPAddress = proto.String("0.0.0.0")
	if ok, body := service.GetEventList(rpcInfo); ok {
		t.Log(body)
	} else {
		t.Error("fail")
	}
}

func TestPromotionEvent_QueryEventInfo(t *testing.T) {
	httputil.DebugMode(true)
	service := &PromotionEvent{}
	service.Base = newBase()
	rpcInfo := &netproto.RPCInfo{}
	rpcInfo.UserID = proto.Int(15267354)
	rpcInfo.Cer = proto.String("D12BFB7CC90444188BDF6795CBF9CE7A")
	rpcInfo.IPAddress = proto.String("0.0.0.0")
	if ok, body := service.GetEventInfo(rpcInfo, 25014); ok {
		t.Log(body)
	} else {
		t.Error("fail")
	}
}

func TestPromotionEvent_PostLuckyRouletteRequest(t *testing.T) {
	service := &PromotionEvent{}
	service.Base = newBase()
	rpcInfo := &netproto.RPCInfo{}
	rpcInfo.UserID = proto.Int(15267354)
	rpcInfo.Cer = proto.String("D12BFB7CC90444188BDF6795CBF9CE7A")
	rpcInfo.IPAddress = proto.String("0.0.0.0")
	params := map[string]interface{}{
		"eventID": 25014,
		"userID": rpcInfo.GetUserID(),
		"certificate": rpcInfo.GetCer(),
	}
	if ok, body := service.PostLuckyRouletteRequest(rpcInfo, params); ok {
		t.Log(body)
	} else {
		t.Error("fail")
	}
}