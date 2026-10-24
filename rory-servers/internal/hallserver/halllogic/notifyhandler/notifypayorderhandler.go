package notifyhandler

import (
	"fmt"

	proto "example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/roomutil"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// NotifyPayOrderHandler ?
type NotifyPayOrderHandler struct{}

// NewNotifyPayOrderHandler ?
func NewNotifyPayOrderHandler() *NotifyPayOrderHandler {
	attr := new(NotifyPayOrderHandler)

	return attr
}

// OnNotifyMessage ?
func (h *NotifyPayOrderHandler) OnNotifyMessage(m *msg.Message, ctx *module.MouduleCtx, player *player.Manager) {
	defer util.PanicErrStack()

	mmbyte := m.MsgData.([]byte)
	mm := &netproto.PayOrderList{}
	err := proto.Unmarshal(mmbyte, mm)

	if err != nil {
		mlog.Error("解析錯誤 HostAndBackUp")
	}
	if mm.GetServerID() == 0 {
		for _, info := range mm.GetPayOrderList() {
			if pr := player.GetPlayerByUserID(info.GetUserID()); pr != nil {
				m, _ := roomutil.GetTipMessage(fmt.Sprintf("亲爱的玩家:\n您充值的金币:%d已到帐\n( 如有任何疑问请联系客服)", info.GetAmount()/100), roomutil.MsgTypeIDTIPRECHARGE)
				ctx.Server.SendMessageByConnID(pr.Conn.ConnID(), m)
				example.invalid("用戶 %d 到帳金額 %d 完成", info.GetUserID(), info.GetAmount())
			} else {
				mlog.Warn("通知用戶到帳、但 %d 用戶不存在", info.GetUserID())
			}

		}
	} else {
		mlog.Error("NotifyServer 發錯服務器了, %d", mm.GetServerID())
	}
}

// IsHook ?
func (h *NotifyPayOrderHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_PayOrderID)
}
