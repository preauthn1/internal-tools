package rpcmsg

import (
	"bytes"
	"fmt"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RPCCommand RPC 命令
type RPCCommand struct {
	queue *Queue
}

// NewRPCCommand 創建RPC命令對象
func NewRPCCommand(queue *Queue) *RPCCommand {
	cmd := new(RPCCommand)
	cmd.queue = queue

	return cmd
}

// Name ?
func (cmd *RPCCommand) Name() string {
	return "rpcmsg"
}

// Help ?
func (cmd *RPCCommand) Help() string {
	return "RPC隊列消息信息，eg:rpcmsg stat | find 0 | list "
}

// Exec ?
func (cmd *RPCCommand) Exec(args []string) string {
	if len(args) == 0 {
		return fmt.Sprintf("error,%s", cmd.Help())
	}

	switch args[0] {
	case "stat":
		return fmt.Sprintf("wait:%d, current id:%d", cmd.queue.innerQueue.Size(), cmd.queue.innerQueue.currentMsgID)
	case "find":
		d := args[1]

		if !util.IsNumber(d) {
			return "error,必須為整數"
		}

		msgid, _ := strconv.Atoi(d)

		status := cmd.queue.innerQueue.Find(msgid)
		return fmt.Sprintf("msg status:%d", status)
	case "list":
		var buf bytes.Buffer

		for _, m := range cmd.queue.innerQueue.innerWaitMsgList {
			buf.WriteString(fmt.Sprintf("%v\r\n", m))
		}

		return buf.String()
	default:
		return fmt.Sprintf("error,%s", cmd.Help())
	}

}
