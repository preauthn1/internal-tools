package robotaction

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
)

// RobotKicker ?
type RobotKicker interface {
	LogoutRobotUser(user *gamelogic.GameUser)
}

// RobotLoader 機器人加載器
type RobotLoader struct {
	rpcMsgQueue           *rpcmsg.Queue                 // db消息隊列
	robotManager          *robotcomm.RobotManager       // 機器人管理器
	loadRobotInterval     int32                         // 加載時間間隔
	timeFuncManager       *event.TimerFuncManager       // timer管理器
	minMoneyAmount        int64                         // 機器人最小金幣數
	maxMoneyAmount        int64                         // 機器人最大金幣數
	serverID              int32                         // 服務器ID
	scoreStorage          scorestorage.UserScoreStorage // 計分器
	robotActionController gamelogic.ActionController    // 用戶行為控制器
	serverStatus          logicmodule.ServerStatusProvider
}

// NewRobotLoader 創建機器人加載器對象
func NewRobotLoader(serverStatus logicmodule.ServerStatusProvider) *RobotLoader {
	rloader := new(RobotLoader)
	rloader.serverStatus = serverStatus

	return rloader
}

// Start 啟動機器人加載器
func (rl *RobotLoader) Start(ctx *HandlerCtx) {
	rl.rpcMsgQueue = ctx.RPCMsgQueue
	rl.robotManager = ctx.RobotManager
	rl.loadRobotInterval = ctx.RobotCfg.LogonInterval
	rl.timeFuncManager = ctx.TfManager
	rl.minMoneyAmount = ctx.RobotCfg.MinMoneyAmount
	rl.maxMoneyAmount = ctx.RobotCfg.MaxMoneyAmount
	rl.serverID = ctx.Room.ServerID
	rl.scoreStorage = ctx.ScoreStorage
	rl.robotActionController = ctx.ActionController
	mlog.System("加載機器人開始")

	RobotCmdP = NewRobotCmd(rl.robotManager, rl)
	// 添加機器人命令
	command.AddCommand(RobotCmdP)

	rl.balanceRobot()

	rl.timeFuncManager.NewTicker(rl.loadRobotInterval, func() {
		if ctx.ServerStatus.IsServerRuning() {
			rl.balanceRobot()
		}
	})
}

// 平衡机器人
func (rl *RobotLoader) balanceRobot() {
	// 如果服務器已經停止
	if !rl.serverStatus.IsServerRuning() {
		rl.logoutRobot(int32(rl.robotManager.RobotCount()))

		return
	}

	crcount := int32(rl.robotManager.RobotCount())

	// example.invalid("开始加载机器人,当前机器人数%d", crcount)

	logonRobotMsg := new(msg.Message)
	logonRobotMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	logonRobotMsg.SClassID = int32(netproto.DBServerClassID_LogonRobotID)

	logonRobotMsgData := new(netproto.LogonRobot)
	logonRobotMsgData.CurrentRobotCount = proto.Int32(crcount)
	logonRobotMsgData.MinMoneyAmount = proto.Int64(rl.minMoneyAmount)
	logonRobotMsgData.MaxMoneyAmount = proto.Int64(rl.maxMoneyAmount)
	logonRobotMsgData.ServerID = proto.Int32(rl.serverID)
	logonRobotMsg.MsgData = logonRobotMsgData

	rl.rpcMsgQueue.PushMessage(netproto.CreateEmptyRPCInfo(rl.serverID), false, logonRobotMsg)
}

// OnRPCMessage ?
func (rl *RobotLoader) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.LogonRobotRet{}
		return mm
	})

	if rmm == nil {
		return
	}

	mm := rmm.(*netproto.LogonRobotRet)

	if mm == nil {
		return
	}
	// 獲取需要加載的機器人數量
	robotCount := mm.GetRobotCount()

	// 不做處理
	if robotCount == 0 {
		mlog.Debug("收到機器人加載數據.但當前不需要登入登出機器人.")
		return
	}
	// 登出機器人
	if robotCount < 0 {
		logoutCount := -1 * robotCount

		example.invalid("登出 %d 個機器人", logoutCount)

		rl.logoutRobot(logoutCount)
	} else { // 登入機器人
		logonCount := robotCount
		// mlog.NewDefaultLogger().System("登入%d个机器人", logonCount)
		example.invalid("登入 %d 個機器人", logonCount)

		rl.logonRobot(logonCount, mm.GetUserData())
	}

}

// 登入機器人
func (rl *RobotLoader) logonRobot(cn int32, users []*netproto.RobotData) {
	for i := 0; i < int(cn); i++ {
		if len(users) < i+1 {
			example.invalid("服務器返回需要登錄 %d 個機器人, 但機器人數量只有 %d", cn, len(users))
			return
		}

		urdata := users[i]

		if rl.robotManager.GetRobot(urdata.GetUserID()) != nil {
			mlog.Warn("加載機器人 %v, 但機器人已經在當前服務器中", urdata)

			continue
		}

		nur := gamelogic.NewGameUser()
		nur.ConnInfo = nil
		nur.WinCount = urdata.GetWinCount()
		nur.DrawCount = urdata.GetDrawCount()
		nur.LoseCount = urdata.GetLoseCount()

		nur.HDCode = "000000"
		nur.HDType = 1

		nur.UserID = urdata.GetUserID()
		nur.NickName = urdata.GetNickName()
		nur.Score = urdata.GetCashAmount()
		nur.LoginScore = urdata.GetCashAmount()
		nur.UserSex = urdata.GetSex()
		// 設定用戶類型為機器人
		nur.UserType = int32(gamelogic.UserTypeRebot) // urdata.GetUserType()
		nur.FaceID = urdata.GetFaceID()
		nur.TradeID = urdata.GetTradeID()
		// nur.RobotLevel = urdata.GetRobotLevelID()
		// nur.AvatarFrameID = urdata.GetAvatarFrameID()

		// if urdata.GetUserType() != int32(gamelogic.UserTypeRebot) {
		// 	mlog.Warn("加载机器人信息%v,但用户类型数据不正确", nur)
		// }

		nur.ConnInfo = robotcomm.NewRobotConnInfo()

		robot := new(robotcomm.Robot)
		robot.UserInfo = nur
		robot.RebotLevel = urdata.GetRobotLevelID()

		rl.robotManager.AddRobot(robot)

		//example.invalid("机器人%v进入房间", robot)
	}
}

// 登出機器人
func (rl *RobotLoader) logoutRobot(cn int32) {
	canLogoutRobot := rl.robotManager.GetNotGamingRobot()
	mlog.Debug("需登出的機器人數量 %v, 可被登出的機器人數量 %v", cn, len(canLogoutRobot))
	i := int32(0)
	for _, r := range canLogoutRobot {
		if r.UserInfo.IsFree() {
			continue
		}
		rl.robotActionController.UserFree(r.UserInfo)
		i++
		if i >= cn {
			break
		}
	}
	logoutStop := int32(0)
	for i := int32(0); i < cn; i++ {
		robot := rl.robotManager.GetFreeRobot(gamelogic.FreeTableNo)

		if robot == nil {
			logoutStop = i
			break
		}

		rl.LogoutRobotUser(robot.UserInfo)
	}

	prepareLeaveRobotCount := int32(0)
	for _, r := range rl.robotManager.GetAllRobot() {
		if r.IsLeave() {
			prepareLeaveRobotCount++
		}
	}

	robotTotal := rl.robotManager.RobotCount()
	needMore := cn - logoutStop - prepareLeaveRobotCount
	if needMore > 0 && robotTotal > 0 {
		//example.invalid("目前機器人數量%d, 還需登出 %d 個機器人這些機器人將在遊戲結束後登出", robotTotal, needMore)
		for _, r := range rl.robotManager.GetAllRobot() {
			if r.IsLeave() {
				continue
			}
			if needMore < 1 {
				break
			}
			r.SetLeave()
			needMore--
		}
	}
}

// IsHook ?
func (rl *RobotLoader) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_LogonRobotRetID)
}

// LogoutRobotUser 登出機器人
func (rl *RobotLoader) LogoutRobotUser(user *gamelogic.GameUser) {
	// 機器熱在遊戲中,不登出
	if user.IsGameing() {
		return
	}

	// 如果用戶坐下狀態,讓用戶站立
	if user.TableNo != gamelogic.FreeTableNo {
		rl.robotActionController.UserFree(user)
	}

	// 移除機器人
	rl.robotManager.RemoveRobot(user.UserID)

	// 登出機器人
	rl.scoreStorage.OnUserLeave(user)

	// mlog.System("机器人%d离开房间", user.UserID)
}
