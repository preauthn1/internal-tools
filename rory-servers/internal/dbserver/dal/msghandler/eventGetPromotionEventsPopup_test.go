// handler單元測試
package msghandler

import (
	"example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

import (
	"testing"
)

// 充值返利活動  取得個人流水及限制
func TestEventGetPromotionPopupHandler(t *testing.T) {
	
		initDBConnection()
		rpcInfo := netproto.CreateEmptyRPCInfo(1)
		rpcInfo.UserID = proto.Int32(1)
		rpcMessage := new(rpcmsg.RPCMessage)
		rpcMessage.RPCInfo = rpcInfo

		msgdata := new(netproto.PromotionPopup)
		msgdata.Area_ID = proto.Int32(6001001)

		msgData, _ := proto.Marshal(msgdata)
		rpcMessage.MsgData = msgData
		m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_PromotionEventsPopupID), rpcMessage)
		handler := new(EventGetPromotionPopupHandler)
		if !handler.IsHook(m.BClassID, m.SClassID) {
			t.Error("hook 錯誤")
		}

		ctx := new(module.MouduleCtx)
		clt := new(network.ConnInfo)		
		ret := handler.OnMessage(ctx, clt, m)
		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_PromotionEventsPopupRetID) {
			t.Error("OnMessage 錯誤")
		}

		// fmt.Println(ret[0].MsgData.(*netproto.PromotionPopupRet).GetAreaID())
		// fmt.Println(ret[0].MsgData.(*netproto.PromotionPopupRet).GetPopupList())
		// t.Log(ret[0].MsgData.(*netproto.PromotionPopupRet).GetAreaID())
		// t.Log(ret[0].MsgData.(*netproto.PromotionPopupRet).GetPopupList())
			// mlog.System("查看EventGetPromotionPopupHandler回傳值的結果 : %!v(MISSING)", ret[0].MsgData.(*netproto.EventChargeBonusAmountRet).Get())
	
}
