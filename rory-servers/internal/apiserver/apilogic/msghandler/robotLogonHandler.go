package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RobotLogonHandler 機器人登錄處理器
type RobotLogonHandler struct {
}

// IsHook ?
func (h *RobotLogonHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_LogonRobotID)
}

// OnMessage ?
func (h *RobotLogonHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.LogonRobot{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.LogonRobot)

	example.invalid("收到机器人登录数据%v", mm)

	query := map[string]interface{}{
		"serverID":          mm.GetServerID(),
		"currentRobotCount": mm.GetCurrentRobotCount(),
		"minMoneyAmount":    util.CentToMoney(mm.GetMinMoneyAmount()),
		"maxMoneyAmount":    util.CentToMoney(mm.GetMaxMoneyAmount()),
	}

	ok, robotListBody := service.ServerSettingService.GetRobots(rpc, query)
	if !ok {
		mlog.Error("机器人进房失败, ")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	// ps := db.NewSqlParameters()
	// ps.AddIntInput("intServerID", mm.GetServerID())
	// ps.AddIntInput("intCurrentRobotCount", mm.GetCurrentRobotCount())
	// ps.AddBigIntInput("lngMinMoneyAmount", mm.GetMinMoneyAmount())
	// ps.AddBigIntInput("lngMaxMoneyAmount", mm.GetMaxMoneyAmount())
	// ps.AddIntOutput("intLogoutRobotCount", 0)
	// proc := db.NewProcedure("PrGs_RobotLogon", ps)
	// dbret, err := dbutil.DbUtil.ExecProc(proc)
	// //mlog.Error(proc.GetSQLString())
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_LogonRobotRetID), nil)
	retdata := new(netproto.LogonRobotRet)

	logonRobotCount := robotListBody.Data.LogoutRobotCount
	if logonRobotCount == 0 && mm.GetCurrentRobotCount() == 0 {
		mlog.Warn("目前機器人數量0, 加載機器人數量0")
	} else {
		for _, v := range robotListBody.Data.Robots {
			robotData := new(netproto.RobotData)
			robotData.UserID = proto.Int32(v.PlayerID)
			robotData.UserType = proto.Int32(v.PlayerType)
			robotData.NickName = proto.String(v.NickName)
			robotData.FaceID = proto.Int32(v.FaceID)
			robotData.Sex = proto.Int32(v.Sex)
			robotData.CashAmount = proto.Int64(util.MoneyToCent(v.CashAmount))
			robotData.WinCount = proto.Int32(v.WinCount)
			robotData.LoseCount = proto.Int32(v.LoseCount)
			robotData.DrawCount = proto.Int32(v.DrawCount)
			robotData.RobotLevelID = proto.Int32(v.LevelID)

			retdata.UserData = append(retdata.UserData, robotData)
		}

		retdata.RobotCount = proto.Int32(logonRobotCount)
	}

	retmsg.MsgData = retdata

	example.invalid("返回机器人登录数据%v", retdata)

	return []*msg.Message{retmsg}
}
