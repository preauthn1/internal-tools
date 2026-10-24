package msghandler

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameServerStartHandler 遊戲服務器啟動
type GameServerStartHandler struct{}

// IsHook 是否匹配
func (h *GameServerStartHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_GameServerStartID)
}

// OnMessage 處理訊息
func (h *GameServerStartHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GameServerStart{}
		return mm
	})

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GameServerStartRetID), nil)
	retdata := new(netproto.GameServerStartRet)
	retmsg.MsgData = retdata

	if rmm == nil {
		var message = "服務器啟動失敗, 解析封包內容失敗"
		mlog.Error(message)
		retdata.Message = proto.String(message)
		retdata.Code = proto.Int32(-1)
		return []*msg.Message{retmsg}
	}

	mm := rmm.(*netproto.GameServerStart)

	//客制一個clientAddr, 從 0.0.0.0:13101 變成 0.0.0.0:10080/ws/13101
	ip, _ := util.ParseIPAndPort(mm.GetClientAddr())
	clientAddr := fmt.Sprintf("%s/ws/%d", ip, mm.GetServerID())
	loginMoney := util.CentToMoney(int64(mm.GetLoginMoney()))
	params := map[string]interface{}{
		"serverID":          mm.GetServerID(),
		"gameID":            mm.GetGameID(),
		"clientAddr":        clientAddr,
		"httpAddr":          mm.GetHttpAddr(),
		"notifyAddr":        mm.GetNotifyAddr(),
		"minAvailableMoney": loginMoney,
		"flagLevel":         mm.GetShowName(),
	}

	ok, serverStartProcessBody := service.ServerSettingService.PostStartProcess(rpc, params)
	if !ok {
		var message = "服務器啟動失敗"
		if serverStartProcessBody != nil && len(serverStartProcessBody.Message) > 0 {
			message = fmt.Sprintf("%s, %s", message, serverStartProcessBody.Message)
		} else {
			message = fmt.Sprintf("%s, GameAPI 連線失敗", message)
		}
		mlog.Error(message)
		retdata.Message = proto.String(message)
		retdata.Code = proto.Int32(-1)
	} else {
		var message = "服務器啟動成功"
		retdata.Message = proto.String(message)
		retdata.Code = proto.Int32(1)

		retdata.AndroidVersion = proto.String("")
		retdata.IOSVersion = proto.String("")
		retdata.Status = proto.Int32(serverStartProcessBody.Data.GameStatus)
	}

	return []*msg.Message{retmsg}
}
