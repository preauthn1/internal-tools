package msghandler

import (
	"fmt"

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

// HallServerStartHandler 大廳啟動
type HallServerStartHandler struct{}

// IsHook ?
func (rh *HallServerStartHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_HallServerStartID)
}

// OnMessage ?
func (rh *HallServerStartHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.HallServerStart{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.HallServerStart)

	serverid := mm.GetServerID()
	clientaddr := mm.GetClientAddr()
	notifyaddr := mm.GetNotifyAddr()
	httpaddr := mm.GetHttpAddr()

	ps := db.NewSQLParameters()
	ps.AddIntInput("sintServerID", serverid)
	ps.AddVarcharInput("chvClientAddr", clientaddr, 32)
	ps.AddVarcharInput("chvNotifyAddr", notifyaddr, 32)
	ps.AddVarcharInput("chvHttpAddr", httpaddr, 32)

	proc := db.NewProcedure("PrPs_ServerStart", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(dbret).Error("發生錯誤%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()

	mlogWithSQL := mlog.NewDefaultLogger().SetData(proc.GetSQLString())
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_HallServerStartRetID), nil)
	var code int32
	code = 1
	message := "服務器啟動成功"
	if err != nil {
		message = fmt.Sprintf("執行存儲過程PrPs_ServerStart錯誤, %v", err)
		mlogWithSQL.Error(message)
		code = 0
	}
	retdata := new(netproto.HallServerStartRet)
	retdata.Code = proto.Int32(code)
	retdata.Message = proto.String(message)

	if dbret.GetRetTableCount() >= 1 {
		retdata.AndroidVersion = proto.String(dbret[0].GetSingleValue("AndroidVersion").(string))
		retdata.IOSVersion = proto.String(dbret[0].GetSingleValue("IOSVersion").(string))
		retdata.AndroidVersion2 = proto.String(dbret[0].GetSingleValue("AndroidVersion2").(string))
		retdata.IOSVersion2 = proto.String(dbret[0].GetSingleValue("IOSVersion2").(string))
	} else {
		mlogWithSQL.Error("未獲取到平台版本信息")
		retdata.Code = proto.Int32(0)
		retdata.Message = proto.String("未獲取到平台版本信息")
	}

	if dbret.GetRetTableCount() >= 2 {
		retdata.UpdateDomain = proto.String(dbret[1].GetSingleValueString("UpdateDomain"))
		retdata.BillDomain = proto.String(dbret[1].GetSingleValueString("BillDomain"))
		retdata.ClubDomain = proto.String(dbret[1].GetSingleValueString("ClubDomain"))
		retdata.CustomerDomain = proto.String(dbret[1].GetSingleValueString("CustomerDomain"))
	} else {
		mlogWithSQL.Error("未獲取到其他域名信息")
	}

	if dbret.GetRetTableCount() >= 3 {
		tbData := dbret[2]
		for i := 0; i < len(tbData.Rows); i++ {
			config := new(netproto.AreaConfig)
			config.AreaID = proto.Int64(tbData.GetValueByColName(i, "AreaID").(int64))
			config.AreaConfig = proto.String(tbData.GetValueByColName(i, "AreaConfig").(string))
			config.GiftMoney = proto.Int64(tbData.GetValueByColName(i, "GiftMoney").(int64))
			config.DefaultMoney = proto.Int64(tbData.GetValueByColName(i, "DefaultMoney").(int64))
			config.NotifyRefresh = proto.Bool(false)
			retdata.AreaConfigList = append(retdata.AreaConfigList, config)
		}
	}

	if dbret.GetRetTableCount() >= 4 {
		tbData := dbret[3]
		for i := 0; i < len(tbData.Rows); i++ {
			pareaShowGame := new(netproto.AreaShowGameFlag)
			pareaShowGame.AreaID = proto.Int32(int32(tbData.GetValueByColName(i, "AreaID").(int64)))
			pareaShowGame.GameID = proto.Int32(int32(tbData.GetValueByColName(i, "GameID").(int64)))
			pareaShowGame.FlagID = proto.Int32(int32(tbData.GetValueByColName(i, "FlagID").(int64)))
			pareaShowGame.ServerID = proto.Int32(int32(tbData.GetValueByColName(i, "ServerID").(int64)))
			statusBool := tbData.GetValueByColName(0, "Status").(bool)
			status := 1
			if statusBool == false {
				status = 0
			}
			pareaShowGame.Status = proto.Int32(int32(status))
			retdata.AreaShowGameFlagList = append(retdata.AreaShowGameFlagList, pareaShowGame)
		}
	}

	if dbret.GetRetTableCount() >= 5 {
		tbData := dbret[4]
		for i := 0; i < len(tbData.Rows); i++ {
			areaHallVersion := new(netproto.AreaHallVersion)
			areaHallVersion.AreaID = proto.Int64(tbData.GetValueByColName(i, "AreaID").(int64))
			areaHallVersion.AndroidVersion = proto.String(tbData.GetValueByColName(i, "AndroidVersion").(string))
			areaHallVersion.IOSVersion = proto.String(tbData.GetValueByColName(i, "IOSVersion").(string))
			areaHallVersion.AndroidVersion2 = proto.String(tbData.GetValueByColName(i, "AndroidVersion2").(string))
			areaHallVersion.IOSVersion2 = proto.String(tbData.GetValueByColName(i, "IOSVersion2").(string))
			retdata.AreaHallVersionList = append(retdata.AreaHallVersionList, areaHallVersion)
		}
	}

	if dbret.GetRetTableCount() >= 6 {
		tbData := dbret[5]
		for i := 0; i < len(tbData.Rows); i++ {
			areaPlatform := new(netproto.AreaPlatform)
			areaPlatform.AreaID = proto.Int64(tbData.GetValueByColName(i, "AreaID").(int64))
			areaPlatform.PlatformID = proto.Int64(tbData.GetValueByColName(i, "PlatformID").(int64))
			retdata.AreaPlatformList = append(retdata.AreaPlatformList, areaPlatform)
		}
	}

	if dbret.GetRetTableCount() >= 7 {
		tbData := dbret[6]
		for i := 0; i < len(tbData.Rows); i++ {
			platformGameServer := new(netproto.PlatformGameServer)
			platformGameServer.PlatformID = proto.Int64(tbData.GetValueByColName(i, "PlatformID").(int64))
			platformGameServer.GameID = proto.Int64(tbData.GetValueByColName(i, "GameID").(int64))
			platformGameServer.ServerID = proto.Int64(tbData.GetValueByColName(i, "ServerID").(int64))
			platformGameServer.Status = proto.Bool(tbData.GetValueByColName(i, "Status").(bool))
			retdata.PlatformGameServerList = append(retdata.PlatformGameServerList, platformGameServer)
		}
	}

	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
