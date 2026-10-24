package notifyhandler

import (
	proto "example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
)

// AttrNotifyHandler ?
type AttrNotifyHandler struct{}

// NewAttrNotifyHandler ?
func NewAttrNotifyHandler() *AttrNotifyHandler {
	attr := new(AttrNotifyHandler)

	return attr
}

// OnNotifyMessage ?
func (ann *AttrNotifyHandler) OnNotifyMessage(m *msg.Message, ctx *module.MouduleCtx, player *player.Manager) {

	mmbyte := m.MsgData.([]byte)

	mm := &netproto.AttrChangeList{}
	err := proto.Unmarshal(mmbyte, mm)

	if err != nil {
		ctx.Error("解析用戶屬性變化消息錯誤,%v", err)
		return
	}

	for _, attr := range mm.Attrs {
		urid := attr.GetUserID()

		// 構造屬性改變消息
		attrmsg := new(msg.Message)
		attrmsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
		attrmsg.SClassID = int32(netproto.NotifyServerClassID_AttrChangeID)

		attrdata := new(netproto.AttrChangeList)
		attrdata.Attrs = append(attrdata.Attrs, attr)
		attrmsg.MsgData = attrdata

		// 廣播屬性變化
		if urid == 0 {
			ctx.Server.SendAllMessage(attrmsg)
		} else {
			pr := player.GetPlayerByUserID(urid)
			if pr == nil {
				continue
			}

			// 發送消息
			ctx.Server.SendMessageByConnID(pr.Conn.ConnID(), attrmsg)
		}

		example.invalid("通知用戶%d屬性變化,%v", attr.GetUserID(), attr)
	}
}

// IsHook ?
func (ann *AttrNotifyHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_AttrChangeID)
}
