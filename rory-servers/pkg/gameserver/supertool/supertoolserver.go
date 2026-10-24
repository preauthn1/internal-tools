package supertool

import (
	"strconv"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GameSuperToolServer supper tool server
type GameSuperToolServer struct {
	DBMsgQueue *rpcmsg.Queue // 數據庫消息隊列
	ServerID   int
	GameID     int
	FlagID     string
}

// NewDynamicSuperTool 超端 創建配桌型
func NewDynamicSuperTool(dbMsg *rpcmsg.Queue, serverID int, gameID int, flagID string) *GameSuperToolServer {
	st := new(GameSuperToolServer)
	st.DBMsgQueue = dbMsg
	st.ServerID = serverID
	st.GameID = gameID
	st.FlagID = flagID
	return st
}

// Start 啟動服務器
func (st *GameSuperToolServer) Start() {
	st.syncFishSysReward()
	st.syncFishVIPUserReward()
	st.startFishSysPool()
	st.startPlayerPool()
}

// Stop 停止
func (st *GameSuperToolServer) Stop() {}

func (st *GameSuperToolServer) syncFishSysReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolFishUpdSystemRewardID)
	pdata := new(netproto.SuperToolFishUpdSystemRewardDetail)
	flagID, _ := strconv.Atoi(st.FlagID)
	pdata.FlagID = proto.Int32(int32(flagID))
	dataList := new(netproto.SuperToolFishUpdSystemRewardDataList)
	dataList.FishID = proto.Int32(999)
	dataList.FinishAmount = proto.Int64(0)

	pdata.FishRewardList = append(pdata.FishRewardList, dataList)
	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(int32(st.ServerID))
	st.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func (st *GameSuperToolServer) startFishSysPool() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_FishSysPoolID)
	pdata := new(netproto.FishSysPoolData)
	flagID, _ := strconv.Atoi(st.FlagID)
	pdata.FlagID = proto.Int32(int32(flagID))
	pdata.SumBetMoney = proto.Int64(0)
	pdata.SumBullet = proto.Int64(0)
	pdata.SysSumWinLoss = proto.Int64(0)

	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(int32(st.ServerID))
	st.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func (st *GameSuperToolServer) startPlayerPool() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_FishPlayerPoolID)
	pdata := new(netproto.FishPlayerPoolDetail)
	flagID, _ := strconv.Atoi(st.FlagID)
	pdata.FlagID = proto.Int32(int32(flagID))
	dataList := new(netproto.FishPlayerPoolDetailList)
	dataList.UserID = proto.Int32(0)
	dataList.HitMoney = proto.Int64(0)
	dataList.SpendMoney = proto.Int64(0)
	pdata.FishRewardList = append(pdata.FishRewardList, dataList)
	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(int32(st.ServerID))
	st.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func (st *GameSuperToolServer) syncFishVIPUserReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardID)
	pdata := new(netproto.SuperToolfishUpdVIPUserRewardDetail)
	pdata.GameID = proto.Int32(int32(st.GameID))
	dataList := new(netproto.SuperToolFishUpdVIPUserRewardDataList)
	dataList.UserID = proto.Int32(0)
	dataList.FinishAmount = proto.Int64(0)
	pdata.FishUserRewardList = append(pdata.FishUserRewardList, dataList)
	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(int32(st.ServerID))
	st.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}
