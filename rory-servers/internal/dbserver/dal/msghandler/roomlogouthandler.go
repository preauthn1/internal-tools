package msghandler

import (
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RoomLogoutHandler 登出房間
type RoomLogoutHandler struct{}

// IsHook 是否處理
func (h *RoomLogoutHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_LogoutRoomID)
}

// OnMessage 收到處理需求
func (h *RoomLogoutHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.LogoutGameRoomInfo{}
		return mm
	})

	if rpc == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.LogoutGameRoomInfo)

	example.invalid("RoomLogoutHandler-OnMessage:UserId[%d],UpdateAmount[%d],LoginMoney[%d],TradeID[%d]", mm.GetUserID(), mm.GetUpdateAmount(), mm.GetLoginMoney(), mm.GetTradeId())

	ps := db.NewSQLParameters()
	ps.AddIntInput("intMsgID", rpc.GetQueueID())
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddBigIntInput("lngUpdateAmount", mm.GetUpdateAmount())
	ps.AddBigIntInput("lngTaxAmount", mm.GetTaxAmount())
	ps.AddBigIntInput("lngFlowAmount", mm.GetFlowAmount())
	ps.AddBigIntInput("lngBetAmount", mm.GetBetAmount())
	ps.AddBigIntInput("lngPayAmount", mm.GetPayAmount())
	ps.AddIntInput("intWinCount", mm.GetWinCount())
	ps.AddIntInput("intLoseCount", mm.GetLoseCount())
	ps.AddIntInput("intDrawCount", mm.GetDrawCount())

	intime := mm.GetLoginTime()
	outtime := mm.GetLogoutTime()

	if t, err := time.Parse(time.RFC3339, intime); err == nil {
		intime = t.Format("2006-01-02 15:04:05.999")
	}
	if t, err := time.Parse(time.RFC3339, outtime); err == nil {
		outtime = t.Format("2006-01-02 15:04:05.999")
	}

	ps.AddVarcharInput("dtmLoginTime", intime, 32)
	ps.AddVarcharInput("dtmLogoutTime", outtime, 32)
	ps.AddVarcharInput("chvIPAddress", mm.GetIPAddress(), 15)
	ps.AddIntInput("tnyHDType", mm.GetHDType())
	ps.AddVarcharInput("chvHDCode", mm.GetHDCode(), 64)
	ps.AddIntInput("intGameTime", mm.GetGameTime())
	ps.AddIntInput("sintGameID", mm.GetGameID())
	ps.AddIntInput("sintServerID", mm.GetServerID())
	ps.AddBigIntInput("lngLoginMoney", mm.GetLoginMoney())
	ps.AddIntInput("tnyScoreUpdateMode", mm.GetScoreUpdateMode())
	ps.AddBigIntInput("lngTradeID", mm.GetTradeId())
	ps.AddIntInput("intLastWinLose", mm.GetLastWinLose())
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrGs_UserLogout", ps)
	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().Error("發生錯誤%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()

	var retcode int32 = 1

	if _, err := dbutil.DBUtil.ExecProc(proc); err != nil {
		ctx.Error("執行存儲過程PrGs_UserLogout錯誤, %v", err)
		retcode = 0
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}
}
