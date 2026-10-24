package command

import (
	"fmt"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Group 命令组
type Group map[string]Command

// AddCommand 添加命令
func (cmg Group) AddCommand(c Command) error {
	if _, ok := cmg[c.Name()]; ok {
		return fmt.Errorf("命令%s已经存在", c.Name())
	}

	cmg[c.Name()] = c

	return nil
}

// Exec 执行命令
func (cmg Group) Exec(name string, args []string) string {
	defer util.PanicErrStack()

	if args != nil && len(args) > 0 {
		for i, v := range args {
			nv := strings.ToLower(v)
			nv = util.TrimAllSpace(nv)

			args[i] = nv
		}
	}

	var c Command
	var ok bool
	if c, ok = cmg[name]; !ok {
		return "命令不存在,请使用help命令列出所有支持的命令"
	}

	return c.Exec(args)
}
