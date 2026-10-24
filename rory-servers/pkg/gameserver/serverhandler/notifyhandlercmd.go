package serverhandler

// NotifyHandlerCmd 服務器關閉命令
type NotifyHandlerCmd struct {
	ntfHandler *NotifyHandler
}

// NewNotifyHandlerCmd 創建服務器關閉命令
func NewNotifyHandlerCmd(handler *NotifyHandler) *NotifyHandlerCmd {
	cmd := new(NotifyHandlerCmd)
	cmd.ntfHandler = handler

	return cmd
}

// Name 指令名稱
func (cmd *NotifyHandlerCmd) Name() string {
	return "ntfhandler"
}

// Help 指令說明
func (cmd *NotifyHandlerCmd) Help() string {
	return "協調服務器cmd.usage:ntfhandler sendmessage | ntfhandler close"
}

// Exec 執行指令
func (cmd *NotifyHandlerCmd) Exec(args []string) string {
	if len(args) < 1 {
		return "error,seeref usage"
	}

	if args[0] == "sendmessage" {
		data := []byte("test")
		cmd.ntfHandler.TCPMsgClient.SendData(data)
	} else if args[0] == "close" {
		cmd.ntfHandler.TCPMsgClient.Close()
	} else {
		return "error,seeref usage"
	}

	return "ok..."
}
