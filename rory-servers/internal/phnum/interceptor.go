package phnum

import (
	"runtime"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/interceptor"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/netproto/phnum"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// MSGInterceptor 金邊的封包攔截器
var MSGInterceptor interceptor.MessageInterceptor

func init() {
	MSGInterceptor = new(MessageInterceptor)
}

// MessageInterceptor 金邊封包攔截器
type MessageInterceptor struct{}

// Name 攔截器名稱
func (i *MessageInterceptor) Name() string {
	return "PhnumInterceptor"
}

// BeforeSend 發送前
func (i *MessageInterceptor) BeforeSend(m *msg.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Warn("%s [BeforeSend] %v", i.Name(), rec)
		}
	}()

	_, f, l, _ := runtime.Caller(1)
	mlog.Debug("%s [BeforeSend] %v, From %s:%d", i.Name(), m, f, l)

	{ // 定時註冊遊戲 3:1
		if m.BClassID == int32(netproto.MessageBClassID_ServerMgm) && m.SClassID == int32(netproto.ServerMgmClassID_GameServerRegister) {
			mlog.Debug("Before Change(%v): %v", m, m.MsgData)

			rpcm := m.MsgData.(*rpcmsg.RPCMessage)
			pb := rpcm.MsgData.(*netproto.ServerRegInfo)

			info := new(phnum.PhnumServerRegInfo)
			info.ServerID = proto.Int32(pb.GetServerID())
			info.ServerType = proto.Int32(pb.GetServerType())
			info.LianyunID = proto.Int32(int32(conf.LianyunID))

			rpcm.MsgData = info

			mlog.Debug("After Change(%v): %v", m, m.MsgData)
			return
		}
	}

	{ // 註冊大廳 3:2
		if m.BClassID == int32(netproto.MessageBClassID_ServerMgm) && m.SClassID == int32(netproto.ServerMgmClassID_HallServerRegister) {
			mlog.Debug("Before Change(%v): %v", m, m.MsgData)

			buf := m.MsgData.([]byte)

			pb := new(netproto.GameServerRegInfo)
			if err := proto.Unmarshal(buf, pb); err != nil {
				mlog.Error("Phnum MessageInterceptor [BeforeSend] Parse Message %v Error, %v", m, err)
				return
			}
			mlog.Debug("Before Change(%v) - Object: %v", m, pb)
			ip, port := util.ParseIPAndPort(conf.RealAddr)

			info := new(phnum.PhnumGameServerRegInfo)
			info.ServerID = proto.Int32(pb.GetServerID())
			info.GameID = proto.Int32(pb.GetGameID())
			info.ServerName = proto.String(pb.GetServerName())
			info.GameName = proto.String(pb.GetGameName())
			info.Flag = proto.String(pb.GetFlag())

			info.IPAddress = proto.String(ip)
			info.Port = proto.Int32(int32(port))
			info.UserCount = proto.Int32(pb.GetUserCount())
			info.LoginMoney = proto.Int32(pb.GetLoginMoney())
			info.BaseScore = proto.Int32(pb.GetBaseScore())
			info.RobotCount = proto.Int32(pb.GetRobotCount())
			info.AndroidVersion = proto.String(pb.GetAndroidVersion())
			info.IOSVersion = proto.String(pb.GetIOSVersion())
			info.ServerStatus = proto.Int32(1)
			info.NoSelRoom = proto.Int32(1)
			info.LianyunID = proto.Int32(int32(conf.LianyunID))

			mlog.Debug("After Change(%v) - Object: %v", m, info)

			if nbuf, err := proto.Marshal(info); err == nil {
				m.MsgData = nbuf
			} else {
				mlog.Error("Phnum MessageInterceptor [BeforeSend] Marshal Message %v Error, %v", info, err)
				return
			}

			mlog.Debug("After Change(%v): %v", m, m.MsgData)
			return
		}
	}

	{ // 遊戲起動 4:3
		if m.BClassID == int32(netproto.MessageBClassID_DBServer) && m.SClassID == int32(netproto.DBServerClassID_GameServerStartID) {
			mlog.Debug("Before Change(%v): %v", m, m.MsgData)

			rpcm := m.MsgData.(*rpcmsg.RPCMessage)
			pb := rpcm.MsgData.(*netproto.GameServerStart)

			svrdata := new(phnum.PhnumGameServerStart)
			svrdata.ClientAddr = proto.String(conf.RealAddr)
			svrdata.NotifyAddr = proto.String(pb.GetNotifyAddr())
			svrdata.ServerID = proto.Int32(pb.GetServerID())
			svrdata.GameID = proto.Int32(pb.GetGameID())
			svrdata.LoginMoney = proto.Int32(pb.GetLoginMoney())
			svrdata.ServerName = proto.String(pb.GetServerName())
			svrdata.ShowName = proto.String(pb.GetShowName())
			svrdata.HttpAddr = proto.String(pb.GetHttpAddr())
			svrdata.GroupID = proto.Int32(int32(conf.LianyunID))

			rpcm.MsgData = svrdata

			mlog.Debug("After Change(%v): %v", m, m.MsgData)
			return
		}
	}

	{ // PerGameLog 4:11
		if m.BClassID == int32(netproto.MessageBClassID_DBServer) && m.SClassID == int32(netproto.DBServerClassID_DealLogID) {
			mlog.Debug("Before Change(%v): %v", m, m.MsgData)

			rpcm := m.MsgData.(*rpcmsg.RPCMessage)
			pb := rpcm.MsgData.(*netproto.DealLog)

			deallogdata := new(phnum.PhnumDealLog)
			deallogdata.GameBeginTime = proto.String(pb.GetGameBeginTime())
			deallogdata.GameEndTime = proto.String(pb.GetGameEndTime())
			deallogdata.GameID = proto.Int32(pb.GetGameID())
			deallogdata.ServerID = proto.Int32(pb.GetServerID())
			deallogdata.TableNo = proto.Int32(pb.GetTableNo())
			if pb.GetNextJiangRecord() != nil {
				record := new(phnum.NextJiangRecord)
				record.AdminUser = proto.String(pb.GetNextJiangRecord().GetAdminUser())
				record.JiangType = proto.String(pb.GetNextJiangRecord().GetJiangType())
				deallogdata.NextJiangRecord = record
			}

			deallogdata.UserGameRet = proto.String(pb.GetUserGameRet())
			deallogdata.GameData = proto.String(pb.GetGameData())
			deallogdata.TaxAmount = proto.Int64(pb.GetTaxAmount())
			deallogdata.GroupID = proto.Int32(int32(conf.LianyunID))

			rpcm.MsgData = deallogdata

			mlog.Debug("After Change(%v): %v", m, m.MsgData)
			return
		}
	}

	{ // RobotLogon 4:12
		if m.BClassID == int32(netproto.MessageBClassID_DBServer) && m.SClassID == int32(netproto.DBServerClassID_LogonRobotID) {
			mlog.Debug("Before Change(%v): %v", m, m.MsgData)

			rpcm := m.MsgData.(*rpcmsg.RPCMessage)
			pb := rpcm.MsgData.(*netproto.LogonRobot)

			data := new(phnum.PhnumLogonRobot)
			data.CurrentRobotCount = proto.Int32(pb.GetCurrentRobotCount())
			data.MinMoneyAmount = proto.Int64(pb.GetMinMoneyAmount())
			data.MaxMoneyAmount = proto.Int64(pb.GetMaxMoneyAmount())
			data.ServerID = proto.Int32(pb.GetServerID())
			data.GroupID = proto.Int32(int32(conf.LianyunID))

			rpcm.MsgData = data

			mlog.Debug("After Change(%v): %v", m, m.MsgData)
			return
		}
	}
}

// AfterSend 發送後
func (i *MessageInterceptor) AfterSend(m *msg.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Warn("%s [AfterSend] %v", i.Name(), rec)
		}
	}()

	_, f, l, _ := runtime.Caller(1)
	mlog.Debug("%s [AfterSend] %v, From %s:%d", i.Name(), m, f, l)
}

// BeforeReceive 接收前
func (i *MessageInterceptor) BeforeReceive(m *msg.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Warn("%s [BeforeReceive] %v", i.Name(), rec)
		}
	}()

	_, f, l, _ := runtime.Caller(1)
	mlog.Debug("%s [BeforeReceive] %v, From %s:%d", i.Name(), m, f, l)

}

// AfterReceive 接收後
func (i *MessageInterceptor) AfterReceive(m *msg.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			mlog.Warn("%s [AfterReceive] %v", i.Name(), rec)
		}
	}()

	_, f, l, _ := runtime.Caller(1)
	mlog.Debug("%s [AfterReceive] %v, From %s:%d", i.Name(), m, f, l)

}
