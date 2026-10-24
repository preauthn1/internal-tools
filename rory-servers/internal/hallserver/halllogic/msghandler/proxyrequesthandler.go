package msghandler

import (
	"io/ioutil"
	"net/http"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ProxyRequestHandler ?
type ProxyRequestHandler struct {
	PlayerManager *player.Manager
	evtChans      chan func()
}

// NewProxyRequestHandler ?
func NewProxyRequestHandler(PlayerManager *player.Manager) *ProxyRequestHandler {
	handler := new(ProxyRequestHandler)
	handler.PlayerManager = PlayerManager
	handler.evtChans = make(chan func(), 10000)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				callfunc := <-handler.evtChans
				if callfunc != nil {
					util.SafeCall(callfunc)
				}
			}
		}()
	}
	return handler
}

// OnMessage ?
func (h *ProxyRequestHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.ProxyData{}
		return mm
	})

	msgData := rmm.(*netproto.ProxyData)

	if pr := h.PlayerManager.GetPlayerByConnID(clt.ConnID()); pr != nil || !pr.IsGuest() {
		evtfunc := func() {
			url := msgData.GetData()
			body := h.proxyRequest(url)
			retmsg := new(msg.Message)
			retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
			retmsg.SClassID = int32(netproto.HallMsgClassID_ProxyRequestRetID)
			data := new(netproto.ProxyData)
			data.Data = proto.String(body)
			data.Type = proto.String(msgData.GetType())

			retmsg.MsgData = data
			ctx.Server.SendMessage(clt, retmsg)
		}
		h.evtChans <- evtfunc
	} else {
		ctx.Error("找不到連線用戶資料 %v", pr)
	}
	return nil
}

// IsHook ?
func (h *ProxyRequestHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_ProxyRequestID)
}

func (h *ProxyRequestHandler) proxyRequest(url string) string {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	//resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(codeId))

	if err != nil {
		mlog.Error("調用 %s 錯誤, %v", url, err)
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
		mlog.Error("調用 %s 錯誤, %v", url, err)
	}
	return string(body)
}
