package module

import (
	"fmt"
)

// Cmd 模塊命令
type Cmd struct {
	moduleGroup *Group
}

// NewCmd 創建模塊命令對象
func NewCmd(moduleGroup *Group) *Cmd {
	cmd := new(Cmd)
	cmd.moduleGroup = moduleGroup

	return cmd
}

// Name 指令名稱
func (cmd *Cmd) Name() string {
	return "module"
}

// Help 指令說明
func (cmd *Cmd) Help() string {
	return "模塊管理，eg:moudle list | module stop *** "
}

// Exec 執行指令
func (cmd *Cmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error,usage:moudle list | module stop ***"
	}

	cmdname := args[0]
	switch cmdname {
	case "list":
		str := ""
		for _, m := range cmd.moduleGroup.modules {
			str += m.Name()
			str += "\r\n"
		}

		return str
	case "stop":
		if len(args) != 2 {
			return "error,usage:module stop ***"
		}
		mname := args[2]

		um := cmd.moduleGroup.UnLoad(mname)

		if um == nil {
			return "error,模塊名稱不存在"
		}

		return fmt.Sprintf("module %s stop.", mname)

	default:
		return "error,usage:module list"
	}
}
