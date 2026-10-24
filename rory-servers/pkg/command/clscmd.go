package command

import (
	"os"
	"os/exec"
	"runtime"
)

// ClsCommand 清除屏幕命令
type ClsCommand struct{}

// Name 指令名稱
func (c *ClsCommand) Name() string {
	return "cls"
}

// Help 指令說明
func (c *ClsCommand) Help() string {
	return "清理当前屏幕"
}

// Exec 執行
func (c *ClsCommand) Exec(args []string) string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	return ""
}
