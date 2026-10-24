package msghandler

import (
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// AddPerGameLogHandler 寫遊戲記錄詳情
type AddPerGameLogHandler struct{}

// IsHook 是否處理
func (rh *AddPerGameLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_DealLogID)
}

// OnMessage 收到處理需求
func (rh *AddPerGameLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.DealLog{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.DealLog)
	gameid := mm.GetGameID()
	serverid := mm.GetServerID()
	gamedata := mm.GetGameData()
	usergameret := mm.GetUserGameRet()
	begintime := mm.GetGameBeginTime()
	endtime := mm.GetGameEndTime()

	if t, err := time.Parse(time.RFC3339, begintime); err == nil {
		begintime = t.Format("2006-01-02 15:04:05.999")
	}
	if t, err := time.Parse(time.RFC3339, endtime); err == nil {
		endtime = t.Format("2006-01-02 15:04:05.999")
	}

	taxamount := mm.GetTaxAmount()
	tableno := mm.GetTableNo()
	njrecord := mm.GetNextJiangRecord()
	winPrizeID := int64(0)
	if njrecord != nil {
		ps := db.NewSQLParameters()
		ps.AddIntInput("intGameID", gameid)
		ps.AddIntInput("intServerID", serverid)
		ps.AddIntInput("intTableNo", tableno)
		ps.AddVarcharInput("chvAdminUser", njrecord.GetAdminUser(), 20)
		ps.AddVarcharInput("chvPrizeType", njrecord.GetJiangType(), 50)
		proc := db.NewProcedure("PrGs_AddWinPrizeLog", ps)
		ret, err := dbutil.DealLogDBUtil.ExecProc(proc)

		if err != nil {
			ctx.Error("執行存儲過程PrGs_AddWinPrizeLog錯誤, %v, %s", err, proc.GetSQLString())
			winPrizeID = 0
		}
		winPrizeID = ret.GetReturnValue()
	}
	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", gameid)
	ps.AddIntInput("intServerID", serverid)
	ps.AddVarcharInput("chvGameData", gamedata, -1)
	ps.AddVarcharInput("chvUserGameRet", usergameret, -1)
	ps.AddVarcharInput("dtmGameBeginTime", begintime, 32)
	ps.AddVarcharInput("dtmGameEndTime", endtime, 32)
	ps.AddBigIntInput("lngTaxAmount", taxamount)
	ps.AddIntInput("intWinPrizeID", int32(winPrizeID))

	proc := db.NewProcedure("PrGs_AddPerGameLog", ps)
	//mlog.System(proc.GetSQLString())
	_, err := dbutil.DealLogDBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("執行存儲過程PrGs_AddPerGameLog錯誤, %v, %s", err, proc.GetSQLString())
		retcode = 0
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}

	// return nil
}
