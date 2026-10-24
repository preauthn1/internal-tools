package msghandler

import (
	"testing"
	"time"

	"example.invalid/golang/mock/gomock"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/msghandler/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func TestAddPerGameLogHandler(t *testing.T) {
	mockDbUtil := fixtures.MockDBUtil(t)
	mockRet := fixtures.GetExpectAddPerGameLogRet()
	mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockRet, nil).AnyTimes()
	mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockRet, nil).AnyTimes()
	dbutil.DBUtil = mockDbUtil
	dbutil.DealLogDBUtil = mockDbUtil
	//initDBConnection()
	gameid := int32(999)
	serverid := int32(999)
	gameBeginTime := time.Now()
	gameEndTime := time.Now()

	rpcInfo := netproto.CreateEmptyRPCInfo(serverid) //bjl 3
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	deallogdata := new(netproto.DealLog)
	deallogdata.GameBeginTime = proto.String(util.FormatTimeStr(gameBeginTime))
	deallogdata.GameEndTime = proto.String(util.FormatTimeStr(gameEndTime))
	deallogdata.GameID = proto.Int32(gameid)
	deallogdata.ServerID = proto.Int32(serverid)
	deallogdata.TableNo = proto.Int32(0)
	deallogdata.UserGameRet = proto.String("10001:2000:-10:20180401,10002:1000:-10:20180401")
	deallogdata.GameData = proto.String("test add per game log handler GameData")
	deallogdata.TaxAmount = proto.Int64(5)
	record := new(netproto.NextJiangRecord)
	record.AdminUser = proto.String("adminuser")
	record.JiangType = proto.String("jiangtype")
	deallogdata.NextJiangRecord = record
	msgData, _ := proto.Marshal(deallogdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_DealLogID), rpcMessage)
	handler := new(AddPerGameLogHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID) {
		t.Error("OnMessage 錯誤")
	}
}
