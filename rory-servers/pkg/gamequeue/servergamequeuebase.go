package gamequeue

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// GameEndAction
const (
	GameEndActionNone  = 0 // 遊戲結束後 不動作
	GameEndActionReady = 1 // 遊戲結束後 自動準備
)

// ServerGameQueueBase 服務器遊戲撮合基類
type ServerGameQueueBase struct {
	evtDispatcher       *event.Dispatcher
	actionController    gamelogic.ActionController
	Configs             *config.Group
	gameendaction       int  // 遊戲結束後行為
	againstSameIP       bool // 拒絕同IP
	againstLastGameUser bool // 拒絕曾經遊戲過的用戶
}

// NewServerGameQueueBase 創建服務器排隊對象基類
func NewServerGameQueueBase(evtDispatcher *event.Dispatcher, actionController gamelogic.ActionController, cfgs *config.Group) *ServerGameQueueBase {
	svrqueue := new(ServerGameQueueBase)
	svrqueue.evtDispatcher = evtDispatcher
	svrqueue.actionController = actionController
	svrqueue.Configs = cfgs

	cfg := svrqueue.Configs.GetConfig(gamelogic.ConfigFileNameGameServer)

	svrqueue.gameendaction = cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "gameendaction")
	svrqueue.againstSameIP = (cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "againstsameip") == 1)
	svrqueue.againstLastGameUser = (cfg.GetConfigInt(gamelogic.ConfigSectionNameGameQueue, "againstlastgame") == 1)

	return svrqueue
}

// Start 開始運行
func (svrq *ServerGameQueueBase) Start() {
	svrq.registerMsgHandler()
}

func (svrq *ServerGameQueueBase) registerMsgHandler() {
	h := event.NewCommonHandler("gamematch", func(evt *event.Event, evtret *event.Ret) bool {
		urs := svrq.actionController.GetQueueUser()
		if len(urs) == 0 {
			return true
		}

		for _, ur := range urs {

			// 排除的IP
			exIP := ""
			// 排除的桌號
			var exGameTableNo int32 = gamelogic.FreeTableNo
			if svrq.againstSameIP {
				exIP = ur.ConnInfo.IPAddress
			}

			if svrq.againstLastGameUser {
				exGameTableNo = ur.LastTableNo
			}

			tno, sno := svrq.actionController.GetFreeSeatNew(ur, exIP, exGameTableNo)
			mlog.Debug("為用戶[%d]分配座位 Table[%d] Seat[%d]", ur.UserID, tno, sno)
			// msglog := fmt.Sprintf("userid[%v]玩家配桌到 tableno[%v] seatno[%v]", ur.UserID, tno, sno)
			// mylog.MongoLog(ur.UserID, 0, tno, msglog)

			// 如果沒有空閒的座位
			if tno == gamelogic.FreeTableNo || sno == gamelogic.FreeSeatNo {
				msglog := fmt.Sprintf("userid[%v]玩家配桌到 沒有空閒的座位", ur.UserID)
				// mylog.MongoLog(ur.UserID, 0, tno, msglog)
				mlog.System(msglog)
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

func (svrq *ServerGameQueueBase) publishGameMatch() {
	evt := NewGameMatchEvent()
	svrq.evtDispatcher.OnAsyncEvent(evt, nil)
}

// UserQueue 用戶排隊
func (svrq *ServerGameQueueBase) UserQueue(user *gamelogic.GameUser) {}

// UserGameEnd 用戶遊戲結束
func (svrq *ServerGameQueueBase) UserGameEnd(user *gamelogic.GameUser) {

	if svrq.gameendaction == GameEndActionNone {

	} else if svrq.gameendaction == GameEndActionReady {
		svrq.actionController.UserReady(user)
	}
}
