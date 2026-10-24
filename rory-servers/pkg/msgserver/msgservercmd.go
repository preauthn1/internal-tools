package msgserver

import (
	"fmt"
)

// Cmd 消息服务器命令
type Cmd struct {
	msgserver *MsgServer
}

// NewMsgServerCmd 創建消息服務器命令
func NewMsgServerCmd(msgserver *MsgServer) *Cmd {
	cmd := new(Cmd)
	cmd.msgserver = msgserver

	return cmd
}

// Name ?
func (cmd *Cmd) Name() string {
	return "msg"
}

// Help ?
func (cmd *Cmd) Help() string {
	return "消息服務器查看，eg:msg listserver | msg listlogic "
}

// Exec ?
func (cmd *Cmd) Exec(args []string) string {
	if len(args) == 0 {
		return "參數錯誤,usage:msg listserver/msg listlogic"
	}

	cmdname := args[0]
	switch cmdname {
	case "listserver":
		return cmd.msgserver.NetServer.String()
	case "listlogic":
		str := ""
		for _, lg := range cmd.msgserver.GetLogicHandlerGroup().innerLogicHandlers {
			str += fmt.Sprintf("%T", lg)
			str += ","
		}

		return str
	default:
		return "參數錯誤,usage:msg listserver/msg list logic"
	}

}
