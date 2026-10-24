package msghandler

import (
	"fmt"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// FastLogonStatCommand 是否需要兩次登入請求
type FastLogonStatCommand struct{}

// Name ?
func (c *FastLogonStatCommand) Name() string {
	return "fastlogonstat"
}

// Help ?
func (c *FastLogonStatCommand) Help() string {
	return "是否需要兩次登入請求, on:需要, off:不需要"
}

// Exec ?
func (c *FastLogonStatCommand) Exec(args []string) string {
	lenarg := len(args)
	switch lenarg {
	case 1:
		stype := args[0]
		if strings.ToLower(stype) == "on" {
			dbutil.FastLogonStat = true
			return fmt.Sprintf("FastLogonStat:%v", dbutil.FastLogonStat)
		} else if stype == "off" {
			dbutil.HDCodeMap = util.NewSafeMap()
			dbutil.FastLogonStat = false
			return fmt.Sprintf("FastLogonStat:%v", dbutil.FastLogonStat)
		} else {
			return "error args"
		}
	}
	return "nothing happen"
}
