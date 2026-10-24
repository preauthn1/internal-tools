package msghandler

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameVerHandler 遊戲服務器啟動
type GameVerHandler struct{}

// IsHook ?
func (rh *GameVerHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_GameVerID)
}

// OnMessage ?
func (rh *GameVerHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GameVer{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GameVer)

	gameid := mm.GetGameID()
	//serverid := mm.GetServerID()

	sql := fmt.Sprintf("SELECT GameID, GameName, AndroidVersion, IOSVersion, status FROM ServerInfo.GameList WHERE GameID = %d", gameid)

	dbret, err := dbutil.DBUtil.ExecSQL(sql)
	mlog.NewDefaultLogger().SetData(dbret).System("執行SQL %s", sql)

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(dbret).Error("發生錯誤%s, 執行SQL %s", errorStr, sql)
			panic(err)
		}
	}()

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_GameVerRetID), nil)
	retdata := new(netproto.GameVerInfo)

	if err != nil {
		mlog.NewDefaultLogger().Error("執行存儲過程PrGs_ServerStart錯誤, %v, %s", err, sql)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata.AndroidVersion = proto.String(dbret[0].GetSingleValue("AndroidVersion").(string))
	retdata.IOSVersion = proto.String(dbret[0].GetSingleValue("IOSVersion").(string))

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
