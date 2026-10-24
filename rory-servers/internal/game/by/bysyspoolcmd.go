package by

import (
	"fmt"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"strconv"
)

type bySyspoolCmd struct {
	by *Bylogic
}

func NewBySyspoolCmd(by *Bylogic) *bySyspoolCmd {
	cmd := new(bySyspoolCmd)
	cmd.by = by
	return cmd
}

func (byCmd *bySyspoolCmd) Name() string {
	return "syspool"
}

func (byCmd *bySyspoolCmd) Help() string {
	return "捕魚系統水池. eg:syspool getpool取得水池數值 | syspool addpool [正負數值] 設定水池 "
}


func (byCmd *bySyspoolCmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error. seeref usage"
	}

	cmdname := args[0]

	switch cmdname {
	case "getpool":
		result := fmt.Sprintf("%v", byCmd.by.TableManager.Room.FishSysPoolData.GetAfterAdjustSyspoolWinLoss()*-1)
		return result
	case "addpool":
		if len(args) != 2 {
			return "args error."
		}

		argPool, err := strconv.Atoi(args[1])
		if err != nil  {
			return "args pool number error."
		}

		ud := new(byuserdata)
		ud.user = new(gamelogic.GameUser)
		ud.user.UserID = 0
		ud.user.TableNo = 0
		byCmd.by.adjustSyspool(ud, int64(-argPool))
		return "Add Success"
	}

	return "No command."
}