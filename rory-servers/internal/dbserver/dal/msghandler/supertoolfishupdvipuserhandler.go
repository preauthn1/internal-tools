package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// SuperToolFishUpdVIPUserHandler 更新個人捕魚放分 Fish_UpdVIPUser
type SuperToolFishUpdVIPUserHandler struct{}

// IsHook ?
func (fsu *SuperToolFishUpdVIPUserHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardID)
}

// OnMessage ?
func (fsu *SuperToolFishUpdVIPUserHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolfishUpdVIPUserRewardDetail{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperToolfishUpdVIPUserRewardDetail)
	var retcode int32 = 1
	var ret db.RecordSet
	var err error
	for _, v := range mm.FishUserRewardList {
		ps := db.NewSQLParameters()
		ps.AddIntInput("intGameID", mm.GetGameID())
		ps.AddIntInput("intUserID", v.GetUserID())
		ps.AddBigIntInput("intCompleteAmount", v.GetFinishAmount())
		proc := db.NewProcedure("PrGs_SuperTool_Fish_UpdVIPUser", ps)

		ret, err = dbutil.DBUtil.ExecProc(proc)

		if err != nil {
			ctx.Error("执行存储过程PrGs_SuperTool_UpdVIPUser错误, %v, %s", err, proc.GetSQLString())
			retcode = 0
		}
		retcode = int32(ret.GetReturnValue())
	}

	var retmsg *msg.Message
	retmsg = msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardRetID), nil)
	retdata := new(netproto.SuperToolFishGetVIPUserRewardDetail)

	if retcode == 1 {
		if ret.GetRetTableCount() > 0 {
			tbData := ret[0]
			for i := 0; i < len(tbData.Rows); i++ {
				vipUserList := new(netproto.FishVIPUserRewardDataList)
				vipUserList.UserID = proto.Int32(int32(tbData.GetValueByColName(i, "UserID").(int64)))
				vipUserList.AssignAmount = proto.Int64((tbData.GetValueByColName(i, "AssignAmount").(int64)))
				vipUserList.Rate = proto.Int32(int32(tbData.GetValueByColName(i, "AssignRate").(int64)))
				vipUserList.FinishAmount = proto.Int64(tbData.GetValueByColName(i, "CompleteAmount").(int64))
				vipUserList.Status = proto.Int32(int32(tbData.GetValueByColName(i, "Status").(int64)))
				vipUserList.SumAssignAmount = proto.Int64(tbData.GetValueByColName(i, "SumAssignAmount").(int64))

				retdata.FishVIPUserList = append(retdata.FishVIPUserList, vipUserList)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
