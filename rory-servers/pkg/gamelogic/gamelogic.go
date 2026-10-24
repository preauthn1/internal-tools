package gamelogic

import (
	"net/http"
	"sync"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// LogicFunc 遊戲邏輯
type LogicFunc func(table *Table, tfManager *event.TimerFuncManager, configs *config.Group) GameLogic

// SuperToolFunc 超端邏輯
type SuperToolFunc func(room *Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) ISuperToolService

// AsyncCallFunc 異步調用函數
type AsyncCallFunc func(asyncFunc func(), callBackFunc func())

// GameLogic 遊戲邏輯接口
type GameLogic interface {
	// 用戶進入桌子
	UserEnter(gameuser *GameUser, tableno int32, seatno int32)

	// 用戶舉手
	UserReady(gameuser *GameUser)

	// 用户站立
	UserFree(gameuser *GameUser)

	// 用戶斷線
	UserOffLine(gameuser *GameUser)

	// 用戶斷線重連
	UserRecover(gameuser *GameUser)

	// 用戶離開
	UserLeave(gameuser *GameUser) bool

	// 遊戲消息
	GameMessage(user *GameUser, gmsg *msg.Message)

	// 比賽消息
	PiSaiMessage(user *GameUser, gmsg *msg.Message)

	// 是否遊戲中
	IsGaming() bool

	// 桌子號
	TableNo() int32

	// 遊戲開始
	OnGameStart(func())

	// 遊戲結束
	OnGameEnd(func(gametime int32))

	// 遊戲解散
	OnGameCancel(func(s string))

	// 更新用戶積分
	OnUpdateUserScore(func(userid int32, score int64, taxamount int64))

	OnUpdateExperienceScore(func(userid int32, score int64, taxamount int64))

	OnUpdateGameJackpotPool(func(userid int32, amount int64))

	// 踢人
	OnKickUser(func(userid int32, msg string))

	// 發送消息
	OnSendMessage(func(userid int32, m *msg.Message))

	// GetSendMessageFunc 發送消息 function
	GetSendMessageFunc() func(userid int32, m *msg.Message)

	// 發送RPC消息
	OnSendRPCMessage(func(safe bool, m *msg.Message))

	// GetSendRPCMessageFunc 發送RPC消息 function
	GetSendRPCMessageFunc() func(safe bool, m *msg.Message)

	// 捕魚公告
	OnNotifyFishMessage(func(m *msg.Message))

	// 收到RPC消息
	OnRecvRPCMessage(m *msg.Message)

	// 收到Http消息
	OnRecvHTTPMessage(req *http.Request) proto.Message

	// 獲取一個空閒的座位
	GetFreeSeat() (int32, int32)

	// 設置牌局記錄器
	SetGameLog(glog gamelog.GameLog)
	//
	//SetVipUserProvider(provider *VipUserProvider)
	//SetVipUserService(provider *VipUserService)

	// 異步執行
	OnAsyncCall(f AsyncCallFunc)

	// 交换用户最大牌型
	ExchangeMaxCard(userid int32)

	// 提前換桌
	OnUserEarlyChangeTable(func(user *GameUser, early string))

	// 更新設定檔
	RefreshConfigGroup(congifs *config.Group)

	// 桌子上的用户数
	TableRealUserCount() int

	// 此遊戲邏輯目前的狀態描述
	Summary() string
}

// LogicManager 遊戲邏輯管理器
type LogicManager struct {
	logics *sync.Map
}

// NewGameLogicManager 創建遊戲邏輯管理器對象
func NewGameLogicManager(r *Room, cfgs *config.Group, tfManager *event.TimerFuncManager, gcfunc LogicFunc) *LogicManager {
	glm := new(LogicManager)
	glm.logics = new(sync.Map)

	var i int32
	for i = 0; i < r.TableCount; i++ {
		t := r.GetTableByNo(i)
		lg := gcfunc(t, tfManager, cfgs)
		t.GameLogic = lg

		glm.logics.Store(i, lg)
	}

	glccmd := NewGameCmd(glm)
	command.AddCommand(glccmd)

	return glm
}

// GetGameLogic 獲取遊戲邏輯
func (glm *LogicManager) GetGameLogic(tableno int32) GameLogic {
	lg, ok := glm.logics.Load(tableno)

	if ok {
		return lg.(GameLogic)
	}

	return nil
}

// Range 遍歷
func (glm *LogicManager) Range(f func(gl GameLogic)) {
	glm.logics.Range(func(key, value interface{}) bool {
		g, ok := value.(GameLogic)
		if !ok {
			return false
		}
		f(g)
		return true
	})
}
