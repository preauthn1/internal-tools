package msghandler

import (
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// LoadBlockIPListHandler 加載封鎖IP
type LoadBlockIPListHandler struct{}

// IsHook ?
func (h *LoadBlockIPListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_LoadBlockIPID)
}

// OnMessage ?
func (h *LoadBlockIPListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_BlockIPListID), nil)
	retdata := new(netproto.BlockIPList)

	ok, loadBlockIPListBody := service.SystemSettingService.GetDenyIP()

	if !ok && loadBlockIPListBody.Data == nil {
		return nil
	}
	for _, item := range loadBlockIPListBody.Data {
		retdata.IPAddress = append(retdata.IPAddress, item.IPAddress)
	}

	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}

//
//func (h *LoadBlockIPListHandler) requestAPI(method string, path string, params map[string]interface{}, headers ...string ) (result *LoadBlockIPListBody, err error) {
//
//	var groupRequest = new(httputil.WaitGroupRequest)
//	groupRequest.AddTask("0", method, path, params, headers...)
//	var respnose = groupRequest.Send()
//	groupRequest.Output() // 印資訊
//
//	body, err := respnose["0"].Body, respnose["0"].Err
//	if err != nil {
//		return nil, err
//	}
//	var model LoadBlockIPListBody
//	err = json.Unmarshal([]byte(body), &model)
//	result = &model
//	return
//}
