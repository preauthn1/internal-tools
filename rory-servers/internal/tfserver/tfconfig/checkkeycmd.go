package tfconfig

import "fmt"

// CheckKeyCmd 服务器关闭命令
type CheckKeyCmd struct {
}

// NewCheckKeyCmd 创建服务器关闭命令
func NewCheckKeyCmd() *CheckKeyCmd {
	cmd := new(CheckKeyCmd)
	return cmd
}

// Name command name
func (cmd *CheckKeyCmd) Name() string {
	return "checkkey"
}

// Help command help
func (cmd *CheckKeyCmd) Help() string {
	return fmt.Sprintf("配置調整cmd. usage:checkkey on|off(%s) ,etc.. | checkkey sip %s,etc...", cmd.getCurrentStatus(), ConfigCtx.NeedCheckKeySourceIP)
}

func (cmd *CheckKeyCmd) getCurrentStatus() string {
	if ConfigCtx.NeedCheckKey == 1 {
		return "on"
	}
	return "off"
}

// Exec command execution
func (cmd *CheckKeyCmd) Exec(args []string) string {
	if len(args) < 1 {
		return "error,seeref usage"
	}

	if len(args) == 1 {
		if args[0] == "on" {
			ConfigCtx.NeedCheckKey = 1
		} else {
			ConfigCtx.NeedCheckKey = 0
		}
	} else if len(args) > 1 && args[0] == "sip" && args[1] != "" {
		ConfigCtx.NeedCheckKeySourceIP = args[1]
	} else {
		return "error,seeref usage"
	}

	return "ok..."
}
