package roomlogic

import (
	"fmt"
	"strconv"

	"example.invalid/golang/protobuf/proto"
	"gopkg.in/mgo.v2/bson"

	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// SuperToolHandler ?
type SuperToolHandler struct {
	DBMsgQueue   *rpcmsg.Queue // 數據庫消息隊列
	gameID       int32
	serverID     int32
	room         *gamelogic.Room
	superToolCfg gamelogic.SuperToolConfig
	tfManager    *event.TimerFuncManager
	rewardtimer  *event.TimerFunc // 放分時間管理棄
	flagID       string
}

// NewSuperToolHandler ?
func NewSuperToolHandler(room *gamelogic.Room, dbMsgQueue *rpcmsg.Queue, serverID int32, superToolCfg gamelogic.SuperToolConfig, tfManager *event.TimerFuncManager, gameID int32, flagID string) *SuperToolHandler {
	handler := new(SuperToolHandler)
	handler.DBMsgQueue = dbMsgQueue
	handler.gameID = gameID
	handler.serverID = serverID
	handler.room = room
	handler.superToolCfg = superToolCfg
	handler.tfManager = tfManager
	handler.flagID = flagID
	return handler
}

// IsHook ?
func (stServer *SuperToolHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) &&
		(sclsid == int32(netproto.GameMessageClassID_SuperToolFishUpdSystemRewardRetID) ||
			sclsid == int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardRetID) ||
			sclsid == int32(netproto.GameMessageClassID_FishSysPoolRetID) ||
			sclsid == int32(netproto.GameMessageClassID_FishPlayerPoolRetID) ||
			sclsid == int32(netproto.GameMessageClassID_GetBeginnerOfferRetID))
}

// OnRPCMessage ?
func (stServer *SuperToolHandler) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	defer util.PanicErrStack()
	// 系統放分
	if m.SClassID == int32(netproto.GameMessageClassID_SuperToolFishUpdSystemRewardRetID) {
		stServer.SuperToolFishUpdSystemRewardRetID(m, ctx)
	}

	// 個人放分
	if m.SClassID == int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardRetID) {
		stServer.SuperToolFishUpdVIPuserRewardRetID(m, ctx)
	}

	if m.SClassID == int32(netproto.GameMessageClassID_FishSysPoolRetID) {
		stServer.FishSysPoolRetID(m, ctx)
	}

	if m.SClassID == int32(netproto.GameMessageClassID_FishPlayerPoolRetID) {
		stServer.FishPlayerPoolRetID(m, ctx)
	}

	if m.SClassID == int32(netproto.GameMessageClassID_GetBeginnerOfferRetID) {
		stServer.FishBeginnerOfferRetID(m, ctx)
	}

}

// SuperToolFishUpdSystemRewardRetID 系統放分
func (stServer *SuperToolHandler) SuperToolFishUpdSystemRewardRetID(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolFishGetSystemRewardDetail{}
		return mm
	})

	if rmm == nil {
		return
	}

	dbRetMsg := rmm.(*netproto.SuperToolFishGetSystemRewardDetail)

	if stServer.room.Fish == nil || stServer.room.SysRewardList == nil {
		stServer.room.SysRewardList = dbRetMsg
		for _, v := range dbRetMsg.GetFishRewardList() {
			mlog.Debug("init 取得DB sys超端資料成功 GetFishRewardList fisid : %v ", v.GetFishID())
		}

		updatescop := util.NewIntScope(stServer.superToolCfg.SyncInterval-2, stServer.superToolCfg.SyncInterval+2)
		stServer.rewardtimer = stServer.tfManager.NewNamedRandomTicker(updatescop, stServer.syncFishSysReward, "sysReward")
	} else {
		//mlog.Debug("系統放分")
		updateRoomSystemReward(dbRetMsg, stServer)
		insertRoomSystemReward(dbRetMsg, stServer)
		deleteRoomSystemReward(dbRetMsg, stServer)
	}
}

// SuperToolFishUpdVIPuserRewardRetID 個人放分
func (stServer *SuperToolHandler) SuperToolFishUpdVIPuserRewardRetID(m *msg.Message, ctx *module.MouduleCtx) {
	defer util.PanicErrStack()
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperToolFishGetVIPUserRewardDetail{}
		return mm
	})

	if rmm == nil {
		return
	}

	retmsg := rmm.(*netproto.SuperToolFishGetVIPUserRewardDetail)

	if stServer.room.Fish == nil || stServer.room.VIPUserRewardList == nil {
		stServer.room.VIPUserRewardList = retmsg
		for _, v := range retmsg.GetFishVIPUserList() {
			mlog.Debug("init 取得DB 個人超端資料成功 GetFishRewardList userid : %v status[%v] ", v.GetUserID(), v.GetStatus())
		}

		updatescop := util.NewIntScope(stServer.superToolCfg.SyncInterval-2, stServer.superToolCfg.SyncInterval+2)
		stServer.rewardtimer = stServer.tfManager.NewNamedRandomTicker(updatescop, stServer.syncFishVIPUserReward, "vipUserReward")
	} else {
		//mlog.Debug("個人放分")
		updateRoomVIPUserReward(retmsg, stServer)
		insertRoomVIPUserReward(retmsg, stServer)
		deleteRoomVIPUserReward(retmsg, stServer)
	}
}

//新增
func insertRoomSystemReward(dbRetMsg *netproto.SuperToolFishGetSystemRewardDetail, stServer *SuperToolHandler) {
	for _, dbFish := range dbRetMsg.FishRewardList {
		hasData := false
		for _, gameFish := range stServer.room.SysRewardList.GetFishRewardList() {
			if dbFish.GetFishID() == gameFish.GetFishID() {
				hasData = true
				break
			}
		}

		if !hasData {
			mlog.Debug("新增game超端list fishid[%v] ", dbFish.GetFishID())
			stServer.room.SysRewardList.FishRewardList = append(stServer.room.SysRewardList.FishRewardList, dbFish)
		}
	}
}

// 新增
func insertRoomVIPUserReward(retmsg *netproto.SuperToolFishGetVIPUserRewardDetail, st *SuperToolHandler) {
	for _, db := range retmsg.FishVIPUserList {
		hasData := false
		for _, gameFish := range st.room.VIPUserRewardList.GetFishVIPUserList() {
			if db.GetUserID() == gameFish.GetUserID() {
				hasData = true
				break
			}
		}

		if !hasData {
			hasData = false
			for _, t := range st.room.Tables {
				for _, s := range t.Seats {
					if s.SeatUser != nil && s.SeatUser.UserID == db.GetUserID() {
						hasData = true
						break
					}
				}
			}

			if hasData {
				mlog.Debug("新增game 在房間 個人超端list fishid[%v] ", db.GetUserID())
				st.room.VIPUserRewardList.FishVIPUserList = append(st.room.VIPUserRewardList.FishVIPUserList, db)
			}
		}
	}
}

// 更新系統放分
func updateRoomSystemReward(dbRetMsg *netproto.SuperToolFishGetSystemRewardDetail, stServer *SuperToolHandler) {
	for _, gameFish := range stServer.room.SysRewardList.GetFishRewardList() {
		for _, dbFish := range dbRetMsg.FishRewardList {
			if gameFish.GetFishID() == dbFish.GetFishID() {
				gameFish.Amount = dbFish.Amount
				gameFish.Rate = dbFish.Rate
				gameFish.Status = dbFish.Status
				example.invalid("更新系統放分 fishid[%v] 指定金額[%v] 增減殺方率[%v] 狀態[%v]",
					gameFish.GetFishID(), gameFish.GetAmount(), gameFish.GetRate(), gameFish.GetStatus())
			}
		}
	}
}

//更新個人放分
func updateRoomVIPUserReward(retmsg *netproto.SuperToolFishGetVIPUserRewardDetail, st *SuperToolHandler) {
	for _, gameFish := range st.room.VIPUserRewardList.GetFishVIPUserList() {
		for _, db := range retmsg.FishVIPUserList {
			if gameFish.GetUserID() == db.GetUserID() {
				gameFish.AssignAmount = db.AssignAmount
				gameFish.SumAssignAmount = db.SumAssignAmount
				gameFish.Rate = db.Rate
				gameFish.Status = db.Status
				example.invalid("更新個人放分 userid[%v] 指定金額[%v] 完成金額[%v] status[%v]", gameFish.GetUserID(), gameFish.GetAssignAmount(), gameFish.GetSumAssignAmount(), gameFish.GetStatus())
			}
		}
	}
}

// 刪除不存在的user
func deleteRoomVIPUserReward(retmsg *netproto.SuperToolFishGetVIPUserRewardDetail, st *SuperToolHandler) {
	defer util.PanicErrStack()
	for i, gameFish := range st.room.VIPUserRewardList.GetFishVIPUserList() {
		hasData := false
		for _, db := range retmsg.FishVIPUserList {
			if gameFish.GetUserID() == db.GetUserID() {
				hasData = true
				break
			}
		}

		if hasData {
			hasData = false
			for _, t := range st.room.Tables {
				for _, s := range t.Seats {
					if s.SeatUser != nil && s.SeatUser.UserID == gameFish.GetUserID() {
						hasData = true
						break
					}
				}
			}

			if !hasData {
				if st.room.VIPUserRewardList == nil ||
					st.room.VIPUserRewardList.FishVIPUserList == nil ||
					len(st.room.VIPUserRewardList.GetFishVIPUserList()) == int(1) {
					st.room.VIPUserRewardList = new(netproto.SuperToolFishGetVIPUserRewardDetail)
					st.room.VIPUserRewardList.FishVIPUserList = make([]*netproto.FishVIPUserRewardDataList, 0)
					continue
				}
				mlog.System("[weihao] len %v , i[%v]", len(st.room.VIPUserRewardList.FishVIPUserList), i)
				st.room.VIPUserRewardList.FishVIPUserList = append(st.room.VIPUserRewardList.FishVIPUserList[:i], st.room.VIPUserRewardList.FishVIPUserList[i+1:]...)
			}
		} else {
			st.room.VIPUserRewardList.FishVIPUserList = append(st.room.VIPUserRewardList.FishVIPUserList[:i], st.room.VIPUserRewardList.FishVIPUserList[i+1:]...)
		}
	}
}

func deleteRoomSystemReward(dbRetMsg *netproto.SuperToolFishGetSystemRewardDetail, stServer *SuperToolHandler) {
	for i, gameFish := range stServer.room.SysRewardList.GetFishRewardList() {
		hasData := false
		for _, db := range dbRetMsg.FishRewardList {
			if gameFish.GetFishID() == db.GetFishID() {
				hasData = true
				break
			}
		}

		if !hasData {
			stServer.room.SysRewardList.FishRewardList = append(stServer.room.SysRewardList.FishRewardList[:i], stServer.room.SysRewardList.FishRewardList[i+1:]...)
		}
	}
}

func (stServer *SuperToolHandler) syncFishSysReward() {
	if stServer.room.HostAndBackup == HostTypeSlave {
		return
	}

	if stServer.room.SysRewardList == nil || stServer.room.SysRewardList.FishRewardList == nil || len(stServer.room.SysRewardList.GetFishRewardList()) == int(0) {
		stServer.initFishSysReward()
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolFishUpdSystemRewardID)
	pdata := new(netproto.SuperToolFishUpdSystemRewardDetail)
	flagID, _ := strconv.Atoi(stServer.flagID)
	pdata.FlagID = proto.Int32(int32(flagID))
	for _, v := range stServer.room.SysRewardList.GetFishRewardList() {
		dataList := new(netproto.SuperToolFishUpdSystemRewardDataList)
		example.invalid("更新系統放分前 fishid[%v] 完成金額[%v]", v.GetFishID(), v.GetFinishAmount())
		dataList.FishID = proto.Int32(v.GetFishID())
		dataList.FinishAmount = proto.Int64(v.GetFinishAmount())
		pdata.FishRewardList = append(pdata.FishRewardList, dataList)
	}

	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(stServer.serverID)
	stServer.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func (stServer *SuperToolHandler) initFishSysReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolFishUpdSystemRewardID)
	pdata := new(netproto.SuperToolFishUpdSystemRewardDetail)
	flagID, _ := strconv.Atoi(stServer.flagID)
	pdata.FlagID = proto.Int32(int32(flagID))
	dataList := new(netproto.SuperToolFishUpdSystemRewardDataList)
	dataList.FishID = proto.Int32(999)
	dataList.FinishAmount = proto.Int64(0)

	pdata.FishRewardList = append(pdata.FishRewardList, dataList)
	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(stServer.serverID)
	stServer.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func (stServer *SuperToolHandler) syncFishVIPUserReward() {
	if stServer.room.HostAndBackup == HostTypeSlave {
		return
	}

	if stServer.room.VIPUserRewardList == nil || stServer.room.VIPUserRewardList.FishVIPUserList == nil || len(stServer.room.VIPUserRewardList.GetFishVIPUserList()) == int(0) {
		stServer.initFishVIPUserReward()
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardID)
	pdata := new(netproto.SuperToolfishUpdVIPUserRewardDetail)
	pdata.GameID = proto.Int32(stServer.gameID)

	for _, t := range stServer.room.Tables {
		for _, s := range t.Seats {
			if s.SeatUser != nil {
				for _, v := range stServer.room.VIPUserRewardList.FishVIPUserList {
					if s.SeatUser.UserID == v.GetUserID() {
						dataList := new(netproto.SuperToolFishUpdVIPUserRewardDataList)
						example.invalid("更新前 VIPUserRewardList userid[%v] finishAmount[%v] rate[%v] ",
							v.GetUserID(), v.GetFinishAmount(), v.GetRate())
						dataList.UserID = proto.Int32(v.GetUserID())
						dataList.FinishAmount = proto.Int64(v.GetFinishAmount())

						pdata.FishUserRewardList = append(pdata.FishUserRewardList, dataList)
					}
				}
			}
		}
	}

	if len(pdata.GetFishUserRewardList()) == int(0) {
		stServer.initFishVIPUserReward()
		return
	}

	pmsg.MsgData = pdata

	for _, v := range pdata.FishUserRewardList {
		mlog.Debug("更新DB 個人放分 userid[%v] finishAmount[%v] ", v.GetUserID(), v.GetFinishAmount())
	}

	rpcinfo := netproto.CreateEmptyRPCInfo(stServer.serverID)
	stServer.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func (stServer *SuperToolHandler) initFishVIPUserReward() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardID)
	pdata := new(netproto.SuperToolfishUpdVIPUserRewardDetail)
	pdata.GameID = proto.Int32(int32(stServer.gameID))
	dataList := new(netproto.SuperToolFishUpdVIPUserRewardDataList)
	dataList.UserID = proto.Int32(0)
	dataList.FinishAmount = proto.Int64(0)
	pdata.FishUserRewardList = append(pdata.FishUserRewardList, dataList)
	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(int32(stServer.serverID))
	stServer.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

// FishSysPoolRetID ?
func (stServer *SuperToolHandler) FishSysPoolRetID(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.FishSysPoolData{}
		return mm
	})

	if rmm == nil {
		return
	}

	dbRetMsg := rmm.(*netproto.FishSysPoolData)

	if stServer.room.FishSysPoolData == nil {
		findData := bson.M{"message.seq": bson.M{"$exists": true}}
		resultData := mylog.BYQuery(findData)
		f, ok := resultData.(mylog.QueryMessage)
		if ok {
			stServer.room.Seq = proto.Int64(f.Seq)
		} else {
			stServer.room.Seq = proto.Int64(0)
		}

		stServer.room.FishSysPoolData = dbRetMsg
		updatescop := util.NewIntScope(stServer.superToolCfg.SyncInterval-2, stServer.superToolCfg.SyncInterval+2)
		stServer.syncFishSysPoolRetID()
		stServer.rewardtimer = stServer.tfManager.NewNamedRandomTicker(updatescop, stServer.syncFishSysPoolRetID, "sysReward")
	} else {
		//updateSysPool(dbRetMsg, stServer)
	}
}

// FishPlayerPoolRetID ?
func (stServer *SuperToolHandler) FishPlayerPoolRetID(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.FishPlayerPoolDetail{}
		return mm
	})

	if rmm == nil {
		return
	}

	dbRetMsg := rmm.(*netproto.FishPlayerPoolDetail)

	if stServer.room.FishPlayerRewardDetail == nil || len(stServer.room.FishPlayerRewardDetail) == 0 {
		stServer.room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)
		for _, v := range dbRetMsg.FishRewardList {
			stServer.room.FishPlayerRewardDetail[v.GetUserID()] = v
		}

		updatescop := util.NewIntScope(stServer.superToolCfg.SyncInterval-2, stServer.superToolCfg.SyncInterval+2)
		//mlog.Debug("FishPlayerPoolDetail is nil updatescop[%v] dbRetMsg[%+v]", updatescop, stServer.room.FishPlayerRewardDetail)
		stServer.syncFishPlayerPoolRetID()
		stServer.rewardtimer = stServer.tfManager.NewNamedRandomTicker(updatescop, stServer.syncFishPlayerPoolRetID, "vipUserReward")
	} else {
		//mlog.Debug("insertPlayerPool")
		insertPlayerPool(dbRetMsg, stServer)
	}
}

// FishBeginnerOfferRetID ?
func (stServer *SuperToolHandler) FishBeginnerOfferRetID(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.FishBeginnerOfferDetail{}
		return mm
	})

	if rmm == nil {
		return
	}

	if stServer.room.GameType == gamelogic.PiSai {
		return
	}

	dbRetMsg := rmm.(*netproto.FishBeginnerOfferDetail)
	if stServer.room.FisnNewPlayerDetail == nil {
		stServer.room.FisnNewPlayerDetail = make(map[int32]*netproto.FishBeginnerOffer, 0)
		for _, v := range dbRetMsg.NewPlayerList {
			stServer.room.FisnNewPlayerDetail[v.GetUserID()] = v
			//mlog.Debug("新手資料 userid[%v]", v.GetUserID())
			stServer.NewPlayerEnter(v.GetUserID(), dbRetMsg.GetNewPlayerAdjustSyspool())
		}

	} else {
		for _, v := range dbRetMsg.NewPlayerList {
			if _, ok := stServer.room.FisnNewPlayerDetail[v.GetUserID()]; !ok {
				//mlog.Debug("notExsit getDB add server newPlayer userid[%v]", v.GetUserID())
				stServer.room.FisnNewPlayerDetail[v.GetUserID()] = v
				stServer.NewPlayerEnter(v.GetUserID(), dbRetMsg.GetNewPlayerAdjustSyspool())
			} else {
				//mlog.Debug("exsit NewPlayer userid[%v]", v.GetUserID())
			}
		}
	}
}

// NewPlayerEnter ?
func (stServer *SuperToolHandler) NewPlayerEnter(userid, adjustSyspool int32) {
	if stServer.room.FishPlayerRewardDetail == nil {
		stServer.room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)
	}
	if stServer.room.NewPlayerPoolRecord == nil {
		stServer.room.NewPlayerPoolRecord = make(map[int32]int32, 0)
	}
	pFishPlayerReward, okFishPlayerReward := stServer.room.FishPlayerRewardDetail[userid]
	if okFishPlayerReward {
		//mlog.Debug("FishPlayerRewardDetail")
	} else {
		stServer.adjustSyspool(userid, adjustSyspool)

		newPlayer := new(netproto.FishPlayerPoolDetailList)
		newPlayer.UserID = proto.Int32(userid)
		newPlayer.SpendMoney = proto.Int64(int64(0))
		newPlayer.HitMoney = proto.Int64(int64(0))
		stServer.room.FishPlayerRewardDetail[userid] = newPlayer
		stServer.room.NewPlayerPoolRecord[userid] = userid
		return
	}

	if p, ok := stServer.room.FisnNewPlayerDetail[userid]; ok && p.GetStatus() == 0 || (pFishPlayerReward.GetSpendMoney() == 0) {
		stServer.adjustSyspool(userid, adjustSyspool)
		stServer.room.NewPlayerPoolRecord[userid] = userid
	} else {
		//mlog.Warn("userid[%v] no FisnNewPlayerDetail", userid)
	}

}

func (stServer *SuperToolHandler) adjustSyspool(userid int32, money int32) {

	if stServer.room.FishSysPoolData == nil {
		mlog.Warn("增減系統水池FishSysPoolData 是空的nil. 請確認異常.")
		return
	}
	//系統角度, 正的是系統輸錢, 負的是系統贏錢, money正的是要輸, money負的是贏
	stServer.room.FishSysPoolData.AfterAdjustSyspoolWinLoss = proto.Int64(stServer.room.FishSysPoolData.GetSysSumWinLoss() + int64(money))
	stServer.room.FishSysPoolData.PoolAmount = proto.Int64(stServer.room.FishSysPoolData.GetPoolAmount() + int64(money))

	logMsg := fmt.Sprintf("新手進場 增減系統水池 調整前系統總輸贏[%v] 調整金額[%v] 調整後系統總輸贏[%v] ",
		-1*stServer.room.GetSysSumWinLoss(), -money, -stServer.room.GetAfterAdjustSyspoolWinLoss())
	mlog.Debug(logMsg)
	mylog.MongoLog(userid, *stServer.room.Seq, 0, logMsg)

	stServer.room.FishSysPoolData.SysSumWinLoss = proto.Int64(stServer.room.GetAfterAdjustSyspoolWinLoss())
}

func (stServer *SuperToolHandler) syncFishSysPoolRetID() {
	if stServer.room.HostAndBackup == HostTypeSlave {
		return
	}

	if stServer.room.FishSysPoolData == nil {
		stServer.initFishSysPool()
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_FishSysPoolID)
	pdata := new(netproto.FishSysPoolData)
	pdata.FlagID = proto.Int32(stServer.room.FishSysPoolData.GetFlagID())
	pdata.SumBetMoney = proto.Int64(stServer.room.FishSysPoolData.GetSumBetMoney())
	pdata.SumBullet = proto.Int64(stServer.room.FishSysPoolData.GetSumBullet())
	pdata.SysSumWinLoss = proto.Int64(0)
	pdata.PoolAmount = proto.Int64(stServer.room.FishSysPoolData.GetPoolAmount())
	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(stServer.serverID)
	stServer.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)

	stServer.room.FishSysPoolData.PoolAmount = proto.Int64(0)
}

func (stServer *SuperToolHandler) initFishSysPool() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_FishSysPoolID)
	pdata := new(netproto.FishSysPoolData)
	flagID, _ := strconv.Atoi(stServer.flagID)
	pdata.FlagID = proto.Int32(int32(flagID))
	pdata.SumBetMoney = proto.Int64(0)
	pdata.SumBullet = proto.Int64(0)
	pdata.SysSumWinLoss = proto.Int64(0)

	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(int32(stServer.serverID))
	stServer.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

//func insertSysPool(dbRetMsg *netproto.FishSysPoolData, stServer *SuperToolHandler) {
//	stServer.room.FishSysPoolData = dbRetMsg
//}

func (stServer *SuperToolHandler) syncFishPlayerPoolRetID() {
	if stServer.room.HostAndBackup == HostTypeSlave {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_FishPlayerPoolID)
	pdata := new(netproto.FishPlayerPoolDetail)
	flagID, _ := strconv.Atoi(stServer.flagID)
	pdata.FlagID = proto.Int32(int32(flagID))

	for _, t := range stServer.room.Tables {
		for _, s := range t.Seats {
			if s.SeatUser != nil {
				if v, ok := stServer.room.FishPlayerRewardDetail[s.SeatUser.UserID]; ok {
					detail := new(netproto.FishPlayerPoolDetailList)
					mlog.Debug("更新前 player pool userid[%v] GetHitMoney[%v] GetSpendMoney[%v] ",
						v.GetUserID(), v.GetHitMoney(), v.GetSpendMoney())
					detail.UserID = proto.Int32(v.GetUserID())
					detail.HitMoney = proto.Int64(v.GetHitMoney())
					detail.SpendMoney = proto.Int64(v.GetSpendMoney())
					pdata.FishRewardList = append(pdata.FishRewardList, detail)
				}
			}
		}
	}

	if len(pdata.GetFishRewardList()) == int(0) {
		stServer.initFishPlayerPool()
		return
	}

	pmsg.MsgData = pdata

	//for _, v := range pdata.FishRewardList {
	//	mlog.Debug("更新DB player pool userid[%v] GetSpendMoney[%v] hit[%v] ", v.GetUserID(), v.GetSpendMoney(), v.GetHitMoney())
	//}

	rpcinfo := netproto.CreateEmptyRPCInfo(stServer.serverID)
	stServer.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func (stServer *SuperToolHandler) initFishPlayerPool() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_FishPlayerPoolID)
	pdata := new(netproto.FishPlayerPoolDetail)
	flagID, _ := strconv.Atoi(stServer.flagID)
	pdata.FlagID = proto.Int32(int32(flagID))
	dataList := new(netproto.FishPlayerPoolDetailList)
	dataList.UserID = proto.Int32(0)
	dataList.SpendMoney = proto.Int64(0)
	dataList.HitMoney = proto.Int64(0)
	pdata.FishRewardList = append(pdata.FishRewardList, dataList)
	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(int32(stServer.serverID))
	stServer.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

func insertPlayerPool(retmsg *netproto.FishPlayerPoolDetail, st *SuperToolHandler) {

	if len(retmsg.FishRewardList) == int(0) {
		mlog.Error("DB 個人水池是空的.")
	}

	if len(retmsg.FishRewardList) == int(0) {
		st.room.FishPlayerRewardDetail = make(map[int32]*netproto.FishPlayerPoolDetailList, 0)
		st.room.FishPlayerPoolStatus = false
		return
	}

	st.room.FishPlayerPoolStatus = true

}
