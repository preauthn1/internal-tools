package event

import (
	"fmt"
)

// Cmd 事件命令
type Cmd struct {
	evtdispatcher *Dispatcher
}

// NewEventCmd 創建事件命令對象
func NewEventCmd(evtdispatcher *Dispatcher) *Cmd {
	cmd := new(Cmd)
	cmd.evtdispatcher = evtdispatcher

	return cmd
}

// Name 命令名稱
func (cmd *Cmd) Name() string {
	return "event"
}

// Help 命令說明
func (cmd *Cmd) Help() string {
	return "事件處理器，eg:event stat | event listhandler"
}

// Exec 執行命令
func (cmd *Cmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error,seeref usage"
	}

	evtname := args[0]
	switch evtname {
	case "stat":
		etype := "async"
		if cmd.evtdispatcher.asyncEvents == nil {
			etype = "sync"
		}

		return fmt.Sprintf("hander:%d,type:%s:event channel len:%d", len(cmd.evtdispatcher.innerHandlers), etype, len(cmd.evtdispatcher.asyncEvents))
	case "listhandler":
		str := ""
		for _, h := range cmd.evtdispatcher.innerHandlers {
			str += h.Handler.HandlerName()
			str += ","
		}
		return str
	default:
		return "error,seeref usage."
	}

}
