package gamequeue

import (
	"math/rand"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// ServerWinLoseGameQueueBase 輸贏對戰,服務器遊戲撮合基類
type ServerWinLoseGameQueueBase struct {
	evtDispatcher       *event.Dispatcher
	actionController    gamelogic.ActionController
	Configs             *config.Group
	gameendaction       int  // 遊戲結束後行為
	againstSameIP       bool // 拒絕同IP
	againstLastGameUser bool // 拒絕曾經遊戲過的用戶
	randJumpTable       int32
	lessthanuser        int32
}

// NewServerWinLoseGameQueueBase 輸贏對戰,創建服務器排隊對象基類
func NewServerWinLoseGameQueueBase(evtDispatcher *event.Dispatcher, actionController gamelogic.ActionController, cfgs *config.Group) *ServerWinLoseGameQueueBase {
	svrqueue := new(ServerWinLoseGameQueueBase)
	svrqueue.evtDispatcher = evtDispatcher
	svrqueue.actionController = actionController
	svrqueue.Configs = cfgs

	cfg := svrqueue.Configs.GetConfig(gamelogic.ConfigFileNameGameServer)

	svrqueue.gameendaction = cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "gameendaction")
	svrqueue.againstSameIP = (cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "againstsameip") == 1)
	svrqueue.againstLastGameUser = (cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "againstlastgame") == 1)
	svrqueue.randJumpTable = cfg.GetConfigInt32(gamelogic.ConfigSectionNameGameQueue, "randjumptable")
	svrqueue.lessthanuser = cfg.GetConfigInt32(gamelogic.ConfigSectionNameGameQueue, "lessthanuser")

	return svrqueue
}

// Start 開始
func (svrq *ServerWinLoseGameQueueBase) Start() {
	svrq.registerMsgHandler()
}

// 輸贏對戰
func (svrq *ServerWinLoseGameQueueBase) registerMsgHandler() {
	h := event.NewCommonHandler("gamematch", func(evt *event.Event, evtret *event.Ret) bool {
		urs := svrq.actionController.GetQueueUser()
		userCount := svrq.actionController.TotalUserCount()
		if len(urs) == 0 {
			return true
		}
		// 排除的IP
		exIP := ""
		var exGameTableNo int32 = 0
		var islessthanuser bool = true
		for _, ur := range urs {
			// 如果用戶有被指定進入哪一個桌子
			if tno, sno := svrq.actionController.GetNextTableSeat(ur); tno != gamelogic.FreeTableNo || sno != gamelogic.FreeSeatNo {
				svrq.actionController.UserSit(ur, tno, sno)
				continue
			}
			if userCount > svrq.lessthanuser {
				islessthanuser = false
				// 隨機跳桌
				randNumber := int32(1)
				if svrq.randJumpTable != 0 {
					randNumber = svrq.randJumpTable
				}

				rand.Seed(time.Now().UnixNano())
				exGameTableNo = rand.Int31n(randNumber)
				if svrq.againstSameIP {
					exIP = ur.ConnInfo.IPAddress
				}

				if svrq.againstLastGameUser {
					exGameTableNo = ur.LastTableNo
				}
			}

			tno, sno := svrq.actionController.GetWinLoseFreeSeat(ur.LastTableNo, exIP, exGameTableNo, ur.LastWinLose, islessthanuser)
			example.invalid("為用戶%v分配座位.參數,%d,%s,%d,ret:%d,%d", ur, ur.LastTableNo, exIP, exGameTableNo, tno, sno)

			// 如果沒有空閒的座位
			if tno == gamelogic.FreeTableNo || sno == gamelogic.FreeSeatNo {
				return true
			}

			if tno != gamelogic.FreeTableNo && sno != gamelogic.FreeSeatNo {
				// 坐下
				svrq.actionController.UserSit(ur, tno, sno)

				// 如果10秒後還未準備,則自動準備

			}
		}

		return true
	}, EventNameGameMatch)

	svrq.evtDispatcher.PushHandler(h)
}

func (svrq *ServerWinLoseGameQueueBase) publishGameMatch() {
	evt := NewGameMatchEvent()
	svrq.evtDispatcher.OnAsyncEvent(evt, nil)
}

// UserQueue 玩家排隊
func (svrq *ServerWinLoseGameQueueBase) UserQueue(user *gamelogic.GameUser) {}

// UserGameEnd 玩家遊戲結束
func (svrq *ServerWinLoseGameQueueBase) UserGameEnd(user *gamelogic.GameUser) {

	if svrq.gameendaction == GameEndActionNone {

	} else if svrq.gameendaction == GameEndActionReady {
		svrq.actionController.UserReady(user)
	}
}

// GetServerConfig 取得 Server 設定
func (svrq *ServerWinLoseGameQueueBase) GetServerConfig() *gamelogic.GameServerConfig {
	return gamelogic.GetGameServerConfig(svrq.Configs.GetConfig(gamelogic.ConfigFileNameServerInfo), svrq.Configs.GetConfig(gamelogic.ConfigFileNameGameServer))
}
