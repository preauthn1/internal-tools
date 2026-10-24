// handler單元測試
package msghandler

import (
	"testing"
)

// 充值返利活動  取得個人流水及限制
func TestEventChargeBonusAmountHandler(t *testing.T) {
	/*
		initDBConnection()
		rpcInfo := netproto.CreateEmptyRPCInfo(1)
		rpcInfo.UserID = proto.Int32(1)
		rpcMessage := new(rpcmsg.RPCMessage)
		rpcMessage.RPCInfo = rpcInfo

		msgdata := new(netproto.EventChargeBonusAmountMsg)
		msgdata.UserID = proto.Int32(56022519)

		msgData, _ := proto.Marshal(msgdata)
		rpcMessage.MsgData = msgData
		m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_EventChargeBonusAmountID), rpcMessage)
		handler := new(EventChargeBonusAmountHandler)
		if !handler.IsHook(m.BClassID, m.SClassID) {
			t.Error("hook 錯誤")
		}

		ctx := new(module.MouduleCtx)
		clt := new(network.ConnInfo)
		handler.OnMessage(ctx, clt, m)
		ret := handler.OnMessage(ctx, clt, m)
		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_EventChargeBonusAmountRetID) {
			t.Error("OnMessage 錯誤")
		}

		t.Log(ret[0].MsgData.(*netproto.EventChargeBonusAmountRet).GetRate())
		t.Log(ret[0].MsgData.(*netproto.EventChargeBonusAmountRet).GetTimes())
		t.Log(ret[0].MsgData.(*netproto.EventChargeBonusAmountRet).GetNowFlowAmount())
		t.Log(ret[0].MsgData.(*netproto.EventChargeBonusAmountRet).GetFlowAmountLimit())

			mlog.System("查看EventChargeBonusAmountHandler回傳值的結果 : %!v(MISSING)", ret[0].MsgData.(*netproto.EventChargeBonusAmountRet).Get())
	*/
}
