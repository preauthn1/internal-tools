package handler

import (
	"strings"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/netproto/phnum"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// AddPerGameLogHandler 寫遊戲記錄詳情
type AddPerGameLogHandler struct {
}

// IsHook 是否攔截
func (rh *AddPerGameLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_DealLogID)
}

// OnMessage 收到訊息
func (rh *AddPerGameLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &phnum.PhnumDealLog{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*phnum.PhnumDealLog)

	serverid := mm.GetServerID()
	gamedata := mm.GetGameData()
	usergameret := mm.GetUserGameRet()
	begintime := mm.GetGameBeginTime()
	endtime := mm.GetGameEndTime()
	taxamount := mm.GetTaxAmount()

	// 處理 usergameret
	// 原：urets += fmt.Sprintf("%d:%d:%d:%d,", uid, ud.Score, uwin, ud.TradeID)
	// 金邊：rets = append(rets, fmt.Sprintf("%d:%d:%d", uid, ud.Score, uwin))

	var tmprets []string
	rets := strings.Split(usergameret, ",")
	for _, ret := range rets {
		if len(ret) == 0 {
			continue
		}
		items := strings.Split(ret, ":")
		if len(items) == 4 {
			items = items[:3]
		}
		tmprets = append(tmprets, strings.Join(items, ":"))
	}

	ps := db.NewSQLParameters()
	ps.AddIntInput("intServerID", serverid)
	ps.AddVarcharInput("chvGameData", gamedata, -1)
	ps.AddVarcharInput("chvUserGameRet", strings.Join(tmprets, ","), -1)
	ps.AddVarcharInput("dtmGameBeginTime", begintime, 32)
	ps.AddVarcharInput("dtmGameEndTime", endtime, 32)
	ps.AddBigIntInput("lngTaxAmount", taxamount)
	ps.AddIntInput("intGroupID", mm.GetGroupID())

	proc := db.NewProcedure("PrGs_AddPerGameLog", ps)

	_, err := dbutil.DealLogDBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("執行存儲過程PrGs_AddPerGameLog錯誤, %v, %s", err, proc.GetSQLString())
		retcode = 0
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}

}
