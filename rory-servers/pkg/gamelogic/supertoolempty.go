package gamelogic

import (
	"net/http"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// EmptySuperToolService 空的超端服務
type EmptySuperToolService struct {
	superToolHandler ISuperToolHandler
}

// NewEmptySuperToolService 空的超端服務
func NewEmptySuperToolService() ISuperToolService {
	s := new(EmptySuperToolService)
	s.superToolHandler = new(EmptySuperToolHandler)
	return s
}

// Start ?
func (s *EmptySuperToolService) Start() {}

// Stop ?
func (s *EmptySuperToolService) Stop() {}

// GetSuperToolHandler ?
func (s *EmptySuperToolService) GetSuperToolHandler() ISuperToolHandler { return s.superToolHandler }

// RefreshVIPUser ?
func (s *EmptySuperToolService) RefreshVIPUser() {}

// GetRoomServerID ?
func (s *EmptySuperToolService) GetRoomServerID() int32 { return -1 }

// PushRPCMessage ?
func (s *EmptySuperToolService) PushRPCMessage(rpcinfo *netproto.RPCInfo, safe bool, msg *msg.Message) error {
	return nil
}

// HTTPService ?
func (s *EmptySuperToolService) HTTPService(w http.ResponseWriter, req *http.Request) {}

// GetGameID 獲取遊戲ID
func (s *EmptySuperToolService) GetGameID() int32 { return -1 }

// GetFlag ?
func (s *EmptySuperToolService) GetFlag() string { return "" }

// EmptySuperToolHandler ?
type EmptySuperToolHandler struct{}

// GetVipUserOrg ?
func (s *EmptySuperToolHandler) GetVipUserOrg(userid int32) proto.Message { return nil }

// SetVipUserList ?
func (s *EmptySuperToolHandler) SetVipUserList(list proto.Message) {}

// IsVipUserExists 用戶是否存在
func (s *EmptySuperToolHandler) IsVipUserExists(userid int32) bool { return false }

// SaveVipUser ?
func (s *EmptySuperToolHandler) SaveVipUser(userid int32) {}

// IsHook ?
func (s *EmptySuperToolHandler) IsHook(bclsid int32, sclsid int32) bool { return false }

// OnRPCMessage ?
func (s *EmptySuperToolHandler) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {}

// EmptySuperToolCreator ?
var EmptySuperToolCreator = func(room *Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) ISuperToolService {
	return NewEmptySuperToolService()
}

// GetNextTable ?
func (s *EmptySuperToolHandler) GetNextTable(userid int32) int32 { return -1 }

// ResetNextTable ?
func (s *EmptySuperToolHandler) ResetNextTable(userid int32) {}

// GetPoolWinLoss ?
func (s *EmptySuperToolHandler) GetPoolWinLoss() *netproto.SuperTool_Msg_GetSysPoolWinLoseRet {
	return nil
}

// GetSysPool 獲取系統水池
func (s *EmptySuperToolHandler) GetSysPool() int64 { return -1 }

// GetWinLine 獲取贏錢線
func (s *EmptySuperToolHandler) GetWinLine() int64 { return -1 }

// UpdateSysPool 更新水池
func (s *EmptySuperToolHandler) UpdateSysPool(updatePool int64) {}

// SysPoolStart 啟動系統水池
func (s *EmptySuperToolHandler) SysPoolStart() {}
