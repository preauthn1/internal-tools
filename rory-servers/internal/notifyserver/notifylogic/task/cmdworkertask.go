package task

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// CmdWorkerData ?
type CmdWorkerData struct {
	CommandType int32
	SParam      string
}

// CmdWorkerTask 命令工作通知
type CmdWorkerTask struct {
	pushInterval int
}

// NewCmdWorkerTask 創建命令工作通知對象
func NewCmdWorkerTask(pushInterval int) *CmdWorkerTask {
	ret := new(CmdWorkerTask)
	ret.pushInterval = pushInterval
	return ret
}

// GetTaskName ?
func (task *CmdWorkerTask) GetTaskName() string {
	return "命令工作通知"
}

// Run ?
func (task *CmdWorkerTask) Run(ctx *RunCtx) {
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(task.pushInterval) * time.Second)
			start := time.Now()
			cmdWorkerList := task.getCmdWorkerListFromDB(ctx)
			if cmdWorkerList != nil && len(cmdWorkerList) > 0 {
				hasRefreshHallVerAll := false
				hasRefreshGameVerAll := false
				for _, cmd := range cmdWorkerList {
					// 全大廳刷新的命令只需要處理一次
					if cmd.CommandType == 0 && !hasRefreshHallVerAll {
						hasRefreshHallVerAll = true
						refreshVerMsg := new(msg.Message)
						refreshVerMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
						refreshVerMsg.SClassID = int32(netproto.NotifyServerClassID_RefreshVersionID)

						refreshVerMsg.MsgData = nil

						hsLst := ctx.ServerList.GetHallServer()
						if hsLst != nil {
							for _, hs := range hsLst {
								example.invalid("向大廳%v發送更新版本號消息", hs)
								ctx.MsgServer.SendMessage(hs.ConnInfo, refreshVerMsg)
							}
						}
						mlog.System("向全大廳發送更新版本號消息")
					}
					// 全遊戲刷新的命令只需要處理一次
					if cmd.CommandType == 1 && !hasRefreshGameVerAll {
						hasRefreshGameVerAll = true
						refreshVerMsg := new(msg.Message)
						refreshVerMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
						refreshVerMsg.SClassID = int32(netproto.NotifyServerClassID_RefreshVersionID)

						refreshVerMsg.MsgData = nil

						gsLst := ctx.ServerList.GetGameServer()
						if gsLst != nil {
							for _, gs := range gsLst {
								example.invalid("向遊戲%v發送更新版本號消息", gs)
								ctx.MsgServer.SendMessage(gs.ConnInfo, refreshVerMsg)
							}
						}
						mlog.System("向全遊戲發送更新版本號消息")
					}
				}
			} else {
				example.invalid("命令工作通知, 但當前沒有命令工作通知")
			}
			example.invalid("CmdWorkerTask 執行時間 %s", time.Since(start))
			//refreshVerMsg := new(msg.Message)
			//refreshVerMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
			//refreshVerMsg.SClassID = int32(netproto.NotifyServerClassID_RefreshVersionID)
			//
			//refreshVerMsg.MsgData = nil
			//
			//hsLst := ctx.ServerList.GetHallServer()
			//if hsLst != nil {
			//	for _, hs := range hsLst {
			//		example.invalid("向大厅%v发送更新版本号消息", hs)
			//		ctx.MsgServer.SendMessage(hs.ConnInfo, refreshVerMsg)
			//	}
			//}

		}
	}()
}

func (task *CmdWorkerTask) getCmdWorkerListFromDB(ctx *RunCtx) []*CmdWorkerData {
	defer util.PanicErrStack()

	ps := db.NewSQLParameters()

	proc := db.NewProcedure("PrNs_GetCmdWorkerList", ps)
	dbret, err := DbUtil.ExecProc(proc)

	if err != nil {
		mlog.Error("執行存儲過程PrNs_GetCmdWorkerList錯誤, %v", err)
		return nil
	}
	list := []*CmdWorkerData{}
	if dbret.GetRetTableCount() > 0 {
		tbAttr := dbret[0]
		for i := 0; i < len(tbAttr.Rows); i++ {
			data := new(CmdWorkerData)
			data.CommandType = int32(tbAttr.GetValueByColName(i, "CommandType").(int64))
			data.SParam = tbAttr.GetValueByColName(i, "SParam").(string)
			list = append(list, data)
		}
	}

	return list
}
