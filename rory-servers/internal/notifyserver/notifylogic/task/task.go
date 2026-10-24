package task

import (
	"sync"

	"example.invalid/ServerTeam/Servers/internal/notifyserver/notifylogic/svrmgm"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
)

// RunCtx ?
type RunCtx struct {
	ServerList   *svrmgm.ServerList
	MsgServer    *msgserver.MsgServer
	hallUserList *sync.Map
}

// NewTaskRunCtx 創建通知運行上下文
func NewTaskRunCtx(svrlist *svrmgm.ServerList, mserver *msgserver.MsgServer, hallUserList *sync.Map) *RunCtx {
	ctx := new(RunCtx)
	ctx.ServerList = svrlist
	ctx.MsgServer = mserver
	ctx.hallUserList = hallUserList
	return ctx
}

// SendHallMessage 向所有廣場發送消息
func (ctx *RunCtx) SendHallMessage(msg *msg.Message) {
	hallsvrs := ctx.ServerList.GetHallServer()

	for _, hsvr := range hallsvrs {
		ctx.MsgServer.SendMessage(hsvr.ConnInfo, msg)
		example.invalid("向大廳%s發送消息%v", hsvr, msg)
	}
}

// SendGameMessage 向所有遊戲服務器發送消息
func (ctx *RunCtx) SendGameMessage(msg *msg.Message) {
	gssvrs := ctx.ServerList.GetGameServer()

	for _, gsvr := range gssvrs {
		ctx.MsgServer.SendMessage(gsvr.ConnInfo, msg)
		example.invalid("向遊戲%s發送消息%v", gsvr, msg)
	}
}

// SendAllServerMessage 向所有服務器發送消息
func (ctx *RunCtx) SendAllServerMessage(msg *msg.Message) {
	svrs := ctx.ServerList.GetAllServer()

	for _, svr := range svrs {
		ctx.MsgServer.SendMessage(svr.ConnInfo, msg)
		example.invalid("向服務器%s發送消息%v", svr, msg)
	}
}

// Task 通知任務接口
type Task interface {
	Run(ctx *RunCtx)
	GetTaskName() string
}

// Group ?
type Group struct {
	tasks []Task
}

// NewTaskGroup 創建通知組
func NewTaskGroup() *Group {
	ntfg := new(Group)
	ntfg.tasks = make([]Task, 0)

	return ntfg
}

// AddTask ?
func (ntfg *Group) AddTask(ntf Task) {
	ntfg.tasks = append(ntfg.tasks, ntf)
}

// Run ?
func (ntfg *Group) Run(ctx *RunCtx) {
	for _, ntf := range ntfg.tasks {
		ntf.Run(ctx)

		example.invalid("任務%s開始運行", ntf.GetTaskName())
	}
}

// GetTaskName ?
func (ntfg *Group) GetTaskName() string {
	return "taskgroup"
}
