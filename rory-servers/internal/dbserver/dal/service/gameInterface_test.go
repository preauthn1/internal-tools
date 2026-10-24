package service

/**
@todo if api url is empty need notice and don't panic
*/
/* don't test gi
func TestGameInterface_GetGameLaunch(t *testing.T) {
	service := &GameInterface{}
	service.Base = newBase()
	rpcInfo := &netproto.RPCInfo{}
	rpcInfo.UserID = proto.Int(15267354)
	rpcInfo.Cer = proto.String("D12BFB7CC90444188BDF6795CBF9CE7A")
	rpcInfo.IPAddress = proto.String("0.0.0.0")
	if ok, body := service.PostGameLaunch(rpcInfo, 777031); ok {
		t.Log(body)
	} else {
		t.Error("fail")
	}
}

func TestGameInterface_PostWalletAutoTransfer(t *testing.T) {
	service := &GameInterface{}
	service.Base = newBase()
	rpcInfo := &netproto.RPCInfo{}
	rpcInfo.UserID = proto.Int(15267354)
	rpcInfo.Cer = proto.String("D12BFB7CC90444188BDF6795CBF9CE7A")
	rpcInfo.IPAddress = proto.String("0.0.0.0")
	if service.PostWalletAutoTransfer(rpcInfo) {
		t.Log("ok")
	} else {
		t.Error("fail")
	}
}
*/
