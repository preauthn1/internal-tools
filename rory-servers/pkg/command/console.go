package command

import (
	"bufio"
	"os"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

var cmg Group

// AddCommand 增加命令
func AddCommand(c Command) error {
	return cmg.AddCommand(c)
}

func init() {
	cmg = make(map[string]Command)
	cmg.AddCommand(new(HelpCommand))
	cmg.AddCommand(new(ClsCommand))
	cmg.AddCommand(new(QuitCommand))
	cmg.AddCommand(new(TitleCommand))
	cmg.AddCommand(NewLogCommand())
	cmg.AddCommand(new(PProfCommand))
	cmg.AddCommand(new(VersionCommand))
	cmg.AddCommand(new(DebugCommand))
}

// Run 監聽Command指令
func Run() {
	mlog.System("Command runing...")
	go func() {
		defer util.PanicErrStack()
		for {
			reader := bufio.NewReader(os.Stdin)
			data, _, _ := reader.ReadLine()
			sline := string(data)

			sline = strings.ToLower(sline)
			sline = strings.TrimSpace(sline)

			sarr := strings.Split(sline, " ")
			cname := sarr[:1][0]
			cargs := sarr[1:]

			cmdret := cmg.Exec(cname, cargs)

			if cmdret != "" {
				mlog.Printf(cmdret)
			}
		}
	}()
}

// Stop 停止命令行
func Stop() {
	mlog.System("Command stop...")
}
