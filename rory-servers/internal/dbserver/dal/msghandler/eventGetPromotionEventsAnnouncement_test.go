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
func TestEventGetPromotionAnnouncementHandler(t *testing.T) {
	
		initDBConnection()
		rpcInfo := netproto.CreateEmptyRPCInfo(1)
		rpcInfo.UserID = proto.Int32(1)
		rpcMessage := new(rpcmsg.RPCMessage)
		rpcMessage.RPCInfo = rpcInfo

		msgdata := new(netproto.PromotionAnnouncement)
		msgdata.Area_ID = proto.Int32(6001001)

		msgData, _ := proto.Marshal(msgdata)
		rpcMessage.MsgData = msgData
		m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_PromotionEventsAnnouncementID), rpcMessage)
		db2hall := new(EventGetPromotionAnnouncementHandler)
		if !db2hall.IsHook(m.BClassID, m.SClassID) {
			t.Error("hook 錯誤")
		}

		ctx := new(module.MouduleCtx)
		clt := new(network.ConnInfo)
		// db2hall.OnMessage(ctx, clt, m)

		
		ret := db2hall.OnMessage(ctx, clt, m)
		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_PromotionEventsAnnouncementRetID) {
			t.Error("OnMessage 錯誤")
		}

		// fmt.Println(ret[0].MsgData.(*netproto.PromotionAnnouncementRet).GetArea_ID())
		// fmt.Println(ret[0].MsgData.(*netproto.PromotionAnnouncementRet).GetAnnouncementList())

		// fmt.Println("AAA", ret[0].BClassID, ret[0].SClassID)
		// hall2client := new(logicmodule.RPCRetHookHandler)
		// if !hall2client.IsHook(ret[0].BClassID, ret[0].SClassID) {
		// 	t.Error("hall -> client hook 錯誤")
		// }
		// rpcInfo, msg, err := hall2client.OnRPCMessage(ret[0], ctx)
		// if err != nil {
		// 	t.Error("hall -> client OnRPCMessage 錯誤")
		// }

		// if msg[0].BClassID != int32(netproto.MessageBClassID_Hall) || msg[0].SClassID != int32(netproto.HallMsgClassID_PromotionEventsAnnouncementRetID) {
		// 	t.Error("hall -> client OnRPCMessage 錯誤2")
		// }


		// fmt.Println(rpcInfo.GetUserID(), rpcInfo.GetIPAddress(), rpcInfo.GetConnID())		

		// fmt.Println(msg[0].MsgData.(*netproto.PromotionAnnouncementRet).GetArea_ID())
		// fmt.Println(msg[0].MsgData.(*netproto.PromotionAnnouncementRet).GetAnnouncementList())

		// t.Log(ret[0].MsgData.(*netproto.PromotionPopupRet).GetAreaID())
		// t.Log(ret[0].MsgData.(*netproto.PromotionPopupRet).GetPopupList())
			// mlog.System("查看EventGetPromotionPopupHandler回傳值的結果 : %!v(MISSING)", ret[0].MsgData.(*netproto.EventChargeBonusAmountRet).Get())
	
}
