package command

import (
	"os"
)

// QuitCommand 退出命令
type QuitCommand struct {
}

// Name 指令名稱
func (c *QuitCommand) Name() string {
	return "quit"
}

// Help 指令說明
func (c *QuitCommand) Help() string {
	return "退出命令，关闭应用程序"
}

// Exec 執行指令
func (c *QuitCommand) Exec(args []string) string {
	if len(args) == 1 {
		if args[0] == "yes" {
			os.Exit(1)
		}
	} else {
		return "请输入quit yes 强制退出"
	}

	return ""
}
