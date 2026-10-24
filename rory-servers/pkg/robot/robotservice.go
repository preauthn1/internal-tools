package robot

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotaction"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Service 機器人服務
type Service struct {
	robotActionCtx    *robotaction.HandlerCtx
	robotLogics       *robotlogic.Manager
	gameLogics        *gamelogic.LogicManager
	robotKicker       robotaction.RobotKicker
	robotActionHander robotaction.Handler
	robotConfig       *robotcomm.RobotConfig
	room              *gamelogic.Room
	actionController  gamelogic.ActionController
	robotManager      *robotcomm.RobotManager
	serverStatus      logicmodule.ServerStatusProvider
	tfManager         *event.TimerFuncManager
}

// NewRobotService 創建機器人服務
func NewRobotService(room *gamelogic.Room, tfManager *event.TimerFuncManager, sstorage scorestorage.UserScoreStorage, rpcmsgQueue *rpcmsg.Queue, gameLogics *gamelogic.LogicManager, cfgs *config.Group, evtChannel *event.Dispatcher, rblgcf robotlogic.LogicFunc, actionController gamelogic.ActionController, serverStatus logicmodule.ServerStatusProvider) *Service {
	// 解析機器人配置
	robotCfg := robotcomm.ParseRobotConfig(cfgs.GetConfig(gamelogic.ConfigFileNameGameServer))
	// 機器人管理器
	robotManager := robotcomm.NewRobotManager()

	rbs := new(Service)

	// 機器人動作處理器上下文
	rbs.robotActionCtx = robotaction.NewRobotActionHandlerCtx()
	rbs.robotActionCtx.Room = room
	rbs.robotActionCtx.RobotManager = robotManager
	rbs.robotActionCtx.TfManager = tfManager
	rbs.robotActionCtx.ScoreStorage = sstorage
	rbs.robotActionCtx.RPCMsgQueue = rpcmsgQueue
	rbs.robotActionCtx.GameLogics = gameLogics
	rbs.robotActionCtx.RobotCfg = robotCfg
	rbs.robotActionCtx.EventChannel = evtChannel
	rbs.robotActionCtx.ActionController = actionController
	rbs.robotActionCtx.ServerStatus = serverStatus

	// 如果機器人未開啟
	if !robotCfg.Enable {
		rblgcf = robotlogic.EmptyRobotLogicCreator
	}

	rbs.robotLogics = robotlogic.NewRobotLogicManager(room, robotCfg, cfgs, tfManager, rblgcf, gameLogics, evtChannel, robotManager)
	rbs.gameLogics = gameLogics
	rbs.robotConfig = robotCfg
	rbs.room = room
	rbs.actionController = actionController
	rbs.robotManager = robotManager
	rbs.serverStatus = serverStatus
	rbs.tfManager = tfManager

	rbs.initRobotActionHandler(serverStatus)

	return rbs
}

// GetRobotActionCtx ?
func (rbs *Service) GetRobotActionCtx() *robotaction.HandlerCtx {
	return rbs.robotActionCtx
}

// SendGameMessage 發送遊戲消息
func (rbs *Service) SendGameMessage(ur *gamelogic.GameUser, m *msg.Message) {

	tableno := ur.TableNo
	rglc := rbs.robotLogics.GetRobotLogic(tableno)
	if rglc == nil {
		return
	}

	// 如果不是遊戲中的用戶
	if ur.IsReady() || ur.IsQueue() || ur.IsFree() || ur.IsSit() {
		return
	}

	robot := rbs.GetRobotByID(ur.UserID)
	if robot == nil {
		return
	}

	rglc.OnGameMessage(robot, m)
}

// KickUser 踢人
func (rbs *Service) KickUser(ur *gamelogic.GameUser) {
	//mlog.System("KickUser %d SeatNo %d", ur.UserID, ur.SeatNo)
	rbs.robotKicker.LogoutRobotUser(ur)
}

// GameStart 遊戲開始
func (rbs *Service) GameStart(ur *gamelogic.GameUser) {
	tableno := ur.TableNo
	rglc := rbs.robotLogics.GetRobotLogic(tableno)
	if rglc == nil {
		return
	}

	robot := rbs.GetRobotByID(ur.UserID)
	if robot == nil {
		return
	}

	// 判断机器人换牌概率
	// 智商 * rnd(0-1) >= 10 - 換牌設置
	// 5 * 0.9 >= 10 - 5
	// if int32(util.RndFloat()*float32(robot.RebotLevel)) >= 10-rbs.robotConfig.MaxCardsRndLvl {
	// 2020-3-12 更換算法
	// rnd 1-100 less then 智商 * 換牌設置
	rnd := util.RndInt(0, 100) + 1
	if rnd <= int(robot.RebotLevel*rbs.robotConfig.MaxCardsRndLvl) {
		glc := rbs.gameLogics.GetGameLogic(tableno)
		if glc != nil {
			glc.ExchangeMaxCard(ur.UserID)
		}
	}

	rglc.OnGameStart(robot)
}

// GameEnd 遊戲結束
func (rbs *Service) GameEnd(ur *gamelogic.GameUser) {
	if ur.IsReady() {
		return
	}

	//通知机器人逻辑游戏结束
	tableno := ur.TableNo
	rglc := rbs.robotLogics.GetRobotLogic(tableno)
	if rglc == nil {
		return
	}

	robot := rbs.GetRobotByID(ur.UserID)
	if robot == nil {
		return
	}

	rglc.OnGameEnd(robot)
	// 檢查機器人是否需要在遊戲後離開
	if robot.IsLeave() {
		rbs.actionController.UserFree(ur)
		if ur.UserState.StateValue() == gamelogic.UserStateFree {
			rbs.KickUser(ur)
		}
		return
	}

	// 檢查機器人金幣數口口
	if !robot.IsLeave() && ur.Score < rbs.robotConfig.MinMoneyAmount {
		mlog.Warn("Robot %d 金幣不夠讓他遊戲結束後離開", robot.UserInfo.UserID)
		robot.SetLeave()
		rbs.actionController.UserFree(ur)
		if ur.UserState.StateValue() == gamelogic.UserStateFree {
			rbs.KickUser(ur)
		}
		return
	}

	// 檢查服務器狀態
	if !rbs.serverStatus.IsServerRuning() {
		rbs.KickUser(ur)
		return
	}

	// 檢查機器人是否繼續遊戲
	rbs.checkRobotFree(ur)

}

// 檢查是否讓機器人站立
func (rbs *Service) checkRobotFree(ur *gamelogic.GameUser) {
	table := rbs.room.GetTableByNo(ur.TableNo)
	maxuser := rbs.robotConfig.TableMaxUserCount
	curusercount := int32(table.TableUserCount())

	odds := rbs.robotConfig.GameEndLeaveOdds
	if curusercount > maxuser {
		odds = rbs.robotConfig.MaxUserGameEndLeaveOdds
	}

	rndv := util.RndInt(0, 100)

	// 隨機概率站立
	if int(100*odds) >= rndv {
		rbs.robotFree(ur, nil)

		mlog.Debug("遊戲結束,機器人%v隨機站立", ur)
		return
	}

	// 檢查人數
	//	minuser := rbs.robotConfig.TableMinUserCount
	//	maxuser := rbs.robotConfig.TableMaxUserCount

	//	curusercount := int32(table.TableUserCount())

	//	if curusercount <= minuser || curusercount > maxuser {
	//		rbs.robotFree(ur, table)

	//		example.invalid("遊戲結束,機器人%v站立.當前桌子人數%d", ur, curusercount)
	//		return
	//	}

	// 準備
	rbs.tfManager.NewTimer(rbs.robotConfig.GameReadyRndWaitSpan(), func() {
		if ur.TableNo == table.TableNo {
			rbs.actionController.UserReady(ur)
		}
	})
}

func (rbs *Service) robotFree(ur *gamelogic.GameUser, cktable *gamelogic.Table) {
	rbs.tfManager.NewTimer(rbs.robotConfig.GameReadyRndWaitSpan(), func() {
		if ur.IsSit() || ur.IsReady() {
			// 需要檢查桌子上人數
			if cktable != nil {
				curusercount := int32(cktable.TableUserCount())
				minuser := rbs.robotConfig.TableMinUserCount
				maxuser := rbs.robotConfig.TableMaxUserCount

				if curusercount <= minuser || curusercount > maxuser {
					rbs.actionController.UserFree(ur)
				} else {
					rbs.actionController.UserReady(ur)
				}
			} else {
				rbs.actionController.UserFree(ur)
			}

		}
	})
}

// GetRobotByID 根據用戶id獲取機器人
func (rbs *Service) GetRobotByID(userid int32) *robotcomm.Robot {
	return rbs.robotManager.GetRobot(userid)
}

// StartRobotService 啟動機器人服務器
func (rbs *Service) StartRobotService() {
	if !rbs.robotConfig.Enable {
		return
	}

	rbs.robotActionHander.Start(rbs.robotActionCtx)

}

// GetRobotLoaderHookHandler 獲取機器人加載攔截器
func (rbs *Service) GetRobotLoaderHookHandler() logicmodule.RPCHookHandler {
	return rbs.robotKicker.(*robotaction.RobotLoader)
}

// 初始化機器人動作處理器
func (rbs *Service) initRobotActionHandler(serverStatus logicmodule.ServerStatusProvider) {

	rbahg := robotaction.NewRobotActionHandlerGroup()
	// 機器人加載
	robotloader := robotaction.NewRobotLoader(serverStatus)
	rbs.robotKicker = robotloader
	rbahg.AddHandler(robotloader)

	if rbs.robotConfig.MatchTableMode == robotaction.MatchTablePoolingOneRobot {
		// 機器人配桌
		rbsqueue := robotaction.NewRobotQueueHandler()
		rbahg.AddHandler(rbsqueue)

		// 機器人主動入桌
		rbash := robotaction.NewRobotAutoSitHandler()
		rbahg.AddHandler(rbash)
	} else if rbs.robotConfig.MatchTableMode == robotaction.QueueTypePoolingNRateRobot {
		// Ｎ個機器人配桌
		rbsqueue := robotaction.NewRobotRateQueueHandler()
		rbahg.AddHandler(rbsqueue)
	} else if rbs.robotConfig.MatchTableMode == robotaction.QueueTypeFixedRobot {
		// 固定配機器人可配置上限
		rbsqueue := robotaction.NewRobotFixedQueueHandler()
		rbahg.AddHandler(rbsqueue)
	} else if rbs.robotConfig.MatchTableMode == robotaction.MatchTablePoolingUserWaitTimeRobot {
		// 當真人超過某數量時調整機器人配桌時間增加
		rbsqueue := robotaction.NewRobotDynamicWaitTimeByRealUserHandler()
		rbahg.AddHandler(rbsqueue)
	}

	rbs.robotActionHander = rbahg
}

// GetRobotManager ?
func (rbs *Service) GetRobotManager() *robotcomm.RobotManager {
	return rbs.robotManager
}

// GetRobotLogicManager 獲取機器人邏輯管理
func (rbs *Service) GetRobotLogicManager() *robotlogic.Manager {
	return rbs.robotLogics
}
