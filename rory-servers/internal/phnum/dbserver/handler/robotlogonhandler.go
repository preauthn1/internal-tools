package handler

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/netproto/phnum"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// RobotLogonHandler 機器人登錄處理器
type RobotLogonHandler struct{}

// IsHook 是否攔截
func (h *RobotLogonHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_LogonRobotID)
}

// OnMessage 收到訊息
func (h *RobotLogonHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &phnum.PhnumLogonRobot{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*phnum.PhnumLogonRobot)

	example.invalid("收到机器人登录数据%v", mm)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intServerID", mm.GetServerID())
	ps.AddIntInput("intCurrentRobotCount", mm.GetCurrentRobotCount())
	ps.AddBigIntInput("lngMinMoneyAmount", mm.GetMinMoneyAmount())
	ps.AddBigIntInput("lngMaxMoneyAmount", mm.GetMaxMoneyAmount())
	ps.AddIntInput("intGroupID", mm.GetGroupID())
	ps.AddIntOutput("intLogoutRobotCount", 0)

	proc := db.NewProcedure("PrGs_RobotLogon", ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_LogonRobotRetID), nil)
	retdata := new(netproto.LogonRobotRet)

	mlog.Debug("%s", proc.GetSQLString())

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrGs_RobotLogon错误, %v", err)
		mlog.Error(strMsg)
	} else {
		logonRobotCount := dbret.GetOutputParamValueInt32("intLogoutRobotCount")
		if logonRobotCount == 0 && mm.GetCurrentRobotCount() == 0 {
			mlog.Warn("目前機器人數量0, 加載機器人數量0, proc:%s", proc.GetSQLString())
		}

		if logonRobotCount > 0 && dbret.GetRetTableCount() > 1 {
			tbRobot := dbret[0]
			for i := 0; i < len(tbRobot.Rows); i++ {
				robotData := new(netproto.RobotData)
				robotData.UserID = proto.Int32(tbRobot.GetValueByColNameInt32(i, "UserID"))
				robotData.UserType = proto.Int32(tbRobot.GetValueByColNameInt32(i, "UserType"))
				robotData.NickName = proto.String(tbRobot.GetValueByColNameString(i, "NickName"))
				robotData.FaceID = proto.Int32(tbRobot.GetValueByColNameInt32(i, "FaceID"))
				robotData.Sex = proto.Int32(tbRobot.GetValueByColNameInt32(i, "Sex"))
				robotData.CashAmount = proto.Int64(tbRobot.GetValueByColName(i, "CashAmount").(int64))
				robotData.WinCount = proto.Int32(tbRobot.GetValueByColNameInt32(i, "WinCount"))
				robotData.LoseCount = proto.Int32(tbRobot.GetValueByColNameInt32(i, "LoseCount"))
				robotData.DrawCount = proto.Int32(tbRobot.GetValueByColNameInt32(i, "DrawCount"))
				robotData.RobotLevelID = proto.Int32(tbRobot.GetValueByColNameInt32(i, "LevelID"))
				retdata.UserData = append(retdata.UserData, robotData)
			}
		}

		retdata.RobotCount = proto.Int32(logonRobotCount)
	}

	retmsg.MsgData = retdata

	example.invalid("返回机器人登录数据%v", retdata)

	return []*msg.Message{retmsg}

}
