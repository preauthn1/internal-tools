package msghandler

import (
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

// RequestAgentDetailHandler 代理相關資訊
type RequestAgentDetailHandler struct{}

// IsHook ?
func (rh *RequestAgentDetailHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && (sclsid == int32(netproto.HallMsgClassID_RequestAgentListID) || sclsid == int32(netproto.HallMsgClassID_JoinAgentListID) || sclsid == int32(netproto.HallMsgClassID_ApplyAgentListID) || sclsid == int32(netproto.HallMsgClassID_CheckAgentID))
}

// OnMessage ?
func (rh *RequestAgentDetailHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo := rpcmsg.GetRPCInfo(m)
	userid := rpcinfo.GetUserID()
	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	if m.SClassID == int32(netproto.HallMsgClassID_RequestAgentListID) { //查詢
		retData := new(netproto.AgentData)
		ps := db.NewSQLParameters()
		ps.AddIntInput("intUserID", userid)
		ps.AddVarcharOutput("chvErrMsg ", "")
		proc := db.NewProcedure("PrPs_Club_AgentInfo_Get", ps)
		ret, err := dbutil.DBUtil.ExecProc(proc)
		defer func() {
			if err := recover(); err != nil {
				errorStr := util.GetStackInfo()
				mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
				panic(err)
			}
		}()
		retDetailData := new(netproto.AgentDetailData)
		var code int32 = 0
		var class int32 = 0
		var status int32 = 0
		var level int32 = 0
		var inviteCode = ""
		var parentUserID int32 = 0
		var parentClass int32 = 0
		var totalFlows int32 = 0
		var totalIncome int32 = 0
		var remainAmount int32 = 0
		var canBeAgent int32 = 0
		var sumCharge int32 = 0
		var userType int32 = 0
		var officeUser1 = ""
		var officeUser2 = ""
		var officeUser3 = ""
		var user1Type int32 = 0
		var user2Type int32 = 0
		var user3Type int32 = 0
		var parentWeixin = ""
		var downloadURL = ""
		var currentIncome int32 = 0
		var qrurl = ""

		retData.DownSteamList = nil

		if err != nil {
			ctx.Error("执行存储过程PrPs_Club_AgentInfo_Get错误, %v", err)
			code = -99
		} else {
			code = int32(ret.GetReturnValue())
			if code > 0 && userid != 0 {
				class = ret[0].GetSingleValueInt32("Class")
				status = ret[0].GetSingleValueInt32("Status")
				level = ret[0].GetSingleValueInt32("Level")
				parentUserID = ret[0].GetSingleValueInt32("ParentUserID")
				parentClass = ret[0].GetSingleValueInt32("ParentClass")
				//distributeRate = ret[0].GetSingleValueInt32("DistributeRate")
				//nextDistributeRate = ret[0].GetSingleValueInt32("NextDistributeRate")
				//upgradeLimit = ret[0].GetSingleValueInt32("UpgradeLimit")
				totalFlows = ret[0].GetSingleValueInt32("TotalFlows")
				totalIncome = ret[0].GetSingleValueInt32("TotalIncome")
				remainAmount = ret[0].GetSingleValueInt32("RemainAmount")
				//yesterdayFlows = ret[0].GetSingleValueInt32("YesterdayFlows")
				//currentFlows = ret[0].GetSingleValueInt32("CurrentFlows")
				//userType = ret[0].GetSingleValueInt32("UserType")
				officeUser1 = ret[0].GetSingleValue("OfficialUser1").(string)
				officeUser2 = ret[0].GetSingleValue("OfficialUser2").(string)
				officeUser3 = ret[0].GetSingleValue("OfficialUser3").(string)
				user1Type = ret[0].GetSingleValueInt32("User1Type")
				user2Type = ret[0].GetSingleValueInt32("User2Type")
				user3Type = ret[0].GetSingleValueInt32("User3Type")
				parentWeixin = ret[0].GetSingleValue("ParentWeixin").(string)
				downloadURL = ret[0].GetSingleValue("DownloadURL").(string)
				qrurl = ret[0].GetSingleValue("QRUrl").(string)
				currentIncome = ret[0].GetSingleValueInt32("CurrentIncome")
				if class == 1 || class == 2 {
					inviteCode = ret[0].GetSingleValue("InviteCode").(string)
					//gradeName = ret[0].GetSingleValue("GradeName").(string)
					//nextGradeName = ret[0].GetSingleValue("NextGradeName").(string)

					if ret.GetRetTableCount() > 2 {
						tbMail := ret[1]
						for i := 0; i < len(tbMail.Rows); i++ {
							downStream := new(netproto.AgentDownSteam)
							downStream.UserID = proto.Int32(int32(tbMail.GetValueByColName(i, "userID").(int64)))
							downStream.Yesterday = proto.Int32(int32(tbMail.GetValueByColName(i, "Yesterday").(int64)))
							downStream.Today = proto.Int32(int32(tbMail.GetValueByColName(i, "Today").(int64)))
							retData.DownSteamList = append(retData.DownSteamList, downStream)
						}
					}
				}
				if parentClass != 2 && userType == 1 {
					sumCharge = ret[0].GetSingleValueInt32("SumCharge")
					if sumCharge >= 500000 { //充值累計達5000元
						canBeAgent = 1
					}
				}
			}
			retDetailData.UserID = proto.Int32(userid)
			retDetailData.Class = proto.Int32(class)
			retDetailData.Status = proto.Int32(status)
			retDetailData.Level = proto.Int32(level)
			retDetailData.InviteCode = proto.String(inviteCode)
			retDetailData.ParentUserID = proto.Int32(parentUserID)
			//retDetailData.GradeName = proto.String(gradeName)
			//retDetailData.NextGradeName = proto.String(nextGradeName)
			//retDetailData.DistributeRate = proto.Int32(distributeRate)
			//retDetailData.NextDistributeRate = proto.Int32(nextDistributeRate)
			//retDetailData.UpgradeLimit = proto.Int32(upgradeLimit)
			retDetailData.TotalFlows = proto.Int32(totalFlows)
			retDetailData.TotalIncome = proto.Int32(totalIncome)
			retDetailData.RemainAmount = proto.Int32(remainAmount)
			//retDetailData.YesterdayFlows = proto.Int32(yesterdayFlows)
			//retDetailData.CurrentFlows = proto.Int32(currentFlows)
			retDetailData.CanBeAgent = proto.Int32(canBeAgent)
			retDetailData.OfficeUser1 = proto.String(officeUser1)
			retDetailData.OfficeUser2 = proto.String(officeUser2)
			retDetailData.OfficeUser3 = proto.String(officeUser3)
			retDetailData.User1Type = proto.Int32(user1Type)
			retDetailData.User2Type = proto.Int32(user2Type)
			retDetailData.User3Type = proto.Int32(user3Type)
			retDetailData.ParentWeixin = proto.String(parentWeixin)
			retDetailData.DownloadURL = proto.String(downloadURL)
			retDetailData.CurrentIncome = proto.Int32(currentIncome)
			retDetailData.QRUrl = proto.String(qrurl)
		}

		retData = new(netproto.AgentData)
		retData.RetCode = proto.Int32(code)
		retData.UserID = proto.Int32(userid)
		retData.AgentDetailData = retDetailData

		retmsg.SClassID = int32(netproto.HallMsgClassID_RequestAgentListRetID)
		retmsg.MsgData = retData

	} else if m.SClassID == int32(netproto.HallMsgClassID_JoinAgentListID) { //申請
		rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
			mm := &netproto.JoinAgentList{}
			return mm
		})

		if rpc == nil {
			return []*msg.Message{retmsg}
		}
		mm := rmm.(*netproto.JoinAgentList)

		var retData *netproto.JoinAgentList
		ps := db.NewSQLParameters()
		ps.AddIntInput("intUserID", userid)
		ps.AddVarcharInput("chvInviteCode", mm.GetInviteCode(), 20)
		ps.AddVarcharOutput("chvErrMsg", "")
		proc := db.NewProcedure("PrPs_JoinUser", ps)
		ret, err := dbutil.ClubDBUtil.ExecProc(proc)
		defer func() {
			if err := recover(); err != nil {
				errorStr := util.GetStackInfo()
				mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
				panic(err)
			}
		}()
		var code int32 = 0
		var inviteCode = ""
		var msg = ""

		if err != nil {
			ctx.Error("执行存储过程PrPs_JoinUser错误, %v", err)
			code = -99
			inviteCode = mm.GetInviteCode()
		} else {
			code = int32(ret.GetReturnValue())
			inviteCode = mm.GetInviteCode()
			if code == -1 {
				msg = ret.GetOutputParamValue("chvErrMsg").(string)
			}
		}
		retData = new(netproto.JoinAgentList)
		retData.UserID = proto.Int32(userid)
		retData.RetCode = proto.Int32(code)
		retData.InviteCode = proto.String(inviteCode)
		retData.RetMsg = proto.String(msg)

		retmsg.SClassID = int32(netproto.HallMsgClassID_JoinAgentListRetID)
		retmsg.MsgData = retData

	} else if m.SClassID == int32(netproto.HallMsgClassID_ApplyAgentListID) { //申請成為代理
		rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
			mm := &netproto.RequestApplyAgent{}
			return mm
		})

		if rpc == nil {
			return []*msg.Message{retmsg}
		}
		mm := rmm.(*netproto.RequestApplyAgent)

		ps := db.NewSQLParameters()
		ps.AddIntInput("intUserID", userid)
		ps.AddVarcharInput("chvTel", mm.GetTel(), 20)
		ps.AddVarcharInput("chvWeixin", mm.GetWeixin(), 50)
		ps.AddVarcharInput("chvQQ", mm.GetQQ(), 50)
		ps.AddVarcharInput("chvEmail", mm.GetEmail(), 50)
		ps.AddVarcharInput("chvCode", mm.GetVCode(), 10)
		ps.AddVarcharOutput("chvErrMsg", "")

		proc := db.NewProcedure("PrPs_ApplyAgent", ps)
		ret, err := dbutil.ClubDBUtil.ExecProc(proc)
		defer func() {
			if err := recover(); err != nil {
				errorStr := util.GetStackInfo()
				mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
				panic(err)
			}
		}()
		var code int32 = 0
		var retMsg string = ""
		if err != nil {
			ctx.Error("执行存储过程PrPs_JoinUser错误, %v", err)
			code = -99
		} else {
			code = int32(ret.GetReturnValue())
			if code == -1 {
				retMsg = ret.GetOutputParamValue("chvErrMsg").(string)
			}
		}
		var retData *netproto.ApplyAgentList
		retData = new(netproto.ApplyAgentList)
		retData.UserID = proto.Int32(userid)
		retData.RetCode = proto.Int32(code)
		retData.RetMsg = proto.String(retMsg)

		retmsg.SClassID = int32(netproto.HallMsgClassID_ApplyAgentListRetID)
		retmsg.MsgData = retData

	} else if m.SClassID == int32(netproto.HallMsgClassID_CheckAgentID) { //查詢是否申請代理
		var retData *netproto.ApplyAgentList
		ps := db.NewSQLParameters()
		ps.AddIntInput("intUserID", userid)
		ps.AddVarcharOutput("chvErrMsg", "")
		proc := db.NewProcedure("PrPs_ApplyAgentCheck", ps)
		ret, err := dbutil.ClubDBUtil.ExecProc(proc)
		defer func() {
			if err := recover(); err != nil {
				errorStr := util.GetStackInfo()
				mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
				panic(err)
			}
		}()
		var code int32 = 0
		var retMsg string = ""
		if err != nil {
			ctx.Error("执行存储过程PrPs_ApplyAgentCheck错误, %v", err)
			code = -99
		} else {
			code = int32(ret.GetReturnValue())
			if code < 0 {
				retMsg = ret.GetOutputParamValue("chvErrMsg").(string)
			}
		}
		retData = new(netproto.ApplyAgentList)
		retData.UserID = proto.Int32(userid)
		retData.RetCode = proto.Int32(code)
		retData.RetMsg = proto.String(retMsg)

		retmsg.SClassID = int32(netproto.HallMsgClassID_CheckAgentRetID)
		retmsg.MsgData = retData

	}
	return []*msg.Message{retmsg}
}
