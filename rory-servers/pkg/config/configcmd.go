package config

import (
	"fmt"
)

// Cmd 配置命令
type Cmd struct {
	configgroup *Group
}

// NewConfigCmd 創建命令對象
func NewConfigCmd(configgroup *Group) *Cmd {
	cmd := new(Cmd)
	cmd.configgroup = configgroup

	return cmd
}

// Name 指令名稱
func (cmd *Cmd) Name() string {
	return "config"
}

// Help 指令說明
func (cmd *Cmd) Help() string {
	return "配置數據，eg:config display 配置文件名 | config getsection 配置文件名 節點名 | config getvalue 配置文件名 節點名 鍵名"
}

// Exec 執行指令
func (cmd *Cmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error,seeref usage."
	}

	cmdname := args[0]
	switch cmdname {
	case "display":
		if len(args) == 1 {
			return cmd.configgroup.String()
		}
		cname := args[1]
		cfg := cmd.configgroup.GetConfig(cname)

		if cfg == nil {
			return "配置文件不存在"
		}

		return cfg.String()
	case "getsection":
		if len(args) != 3 {
			return "error,參數格式不正確.see useage"
		}
		fname := args[1]
		sname := args[2]

		return fmt.Sprintf("%v", cmd.configgroup.GetConfigSection(fname, sname))
	case "getvalue":
		if len(args) != 4 {
			return "error,參數格式不正確.see useage"
		}
		fname := args[1]
		sname := args[2]
		kname := args[3]

		return cmd.configgroup.GetConfigValue(fname, sname, kname)
	default:
		return "error,seeref usage."
	}
}
