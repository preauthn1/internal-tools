package brnn

import (
	"example.invalid/ServerTeam/Servers/pkg/util"
)

import (
	//	"mlog"
	"strconv"
)

type brnnCmd struct {
	lhdlogic *BRNNLogic
}

//创建用户命令
func NewBRNNCmd(fkeb *BRNNLogic) *brnnCmd {
	cmd := new(brnnCmd)
	cmd.lhdlogic = fkeb

	return cmd
}

func (lc *brnnCmd) Name() string {
	return "syspool"
}

func (lc *brnnCmd) Help() string {
	return "百人牛牛.usage:syspool change (tableid) (money)/syspool listpool (tableid)/syspool listpool/syspool zhushui (tableid) (int32)/syspool he (tableid) (概率）"
}

func (lc *brnnCmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error.seeref usage"
	}

	cmdname := args[0]

	switch cmdname {
	case "change":
		if len(args) < 2 {
			return "error.eg kick 10001"
		}

		if len(args) == 3 {
			if !util.IsNumber(args[1]) {
				return "error.桌子号必须为数字"
			}

			if !util.IsNumber(args[2]) {
				return "error.改变后的钱必须为数字"
			}

			tableno, _ := strconv.Atoi(args[1])
			money, _ := strconv.Atoi(args[2])
			t := lc.lhdlogic.GetTable().Room.GetTableByNo(int32(tableno))
			brnn := t.GameLogic.(*BRNNLogic)
			brnn.changeSysPool(int64(money))
		} else {
			if !util.IsNumber(args[1]) {
				return "error.改变后的钱必须为数字"
			}

			money, _ := strconv.Atoi(args[1])
			for i := 0; i < len(lc.lhdlogic.GetTable().Room.Tables); i++ {
				t := lc.lhdlogic.GetTable().Room.Tables[i]
				brnn := t.GameLogic.(*BRNNLogic)
				brnn.changeSysPool(int64(money))
			}
		}

		return "change success"
	case "listpool":
		if len(args) == 2 {
			if !util.IsNumber(args[1]) {
				return "error.房间必须为数字"
			}

			tableno, _ := strconv.Atoi(args[1])
			t := lc.lhdlogic.GetTable().Room.GetTableByNo(int32(tableno))
			brnn := t.GameLogic.(*BRNNLogic)
			brnn.checkSysPool()
		} else {
			alltable := lc.lhdlogic.GetTable().Room.Tables
			for _, t := range alltable {
				brnn := t.GameLogic.(*BRNNLogic)
				brnn.checkSysPool()
			}
		}

		return "listpool success"
	case "zhushui":
		if len(args) < 2 {
			return "error.eg kick 10001"
		}

		if len(args) == 3 {
			if !util.IsNumber(args[1]) {
				return "error.房间必须为数字"
			}

			if !util.IsNumber(args[2]) {
				return "error.改变后的注水百分比必须为数字"
			}

			tableno, _ := strconv.Atoi(args[1])
			zhushui, _ := strconv.Atoi(args[2])
			t := lc.lhdlogic.GetTable().Room.GetTableByNo(int32(tableno))
			brnn := t.GameLogic.(*BRNNLogic)
			brnn.changeZhuShui(int32(zhushui))
		} else {
			zhushui, _ := strconv.Atoi(args[1])
			for i := 0; i < len(lc.lhdlogic.GetTable().Room.Tables); i++ {
				t := lc.lhdlogic.GetTable().Room.Tables[i]
				brnn := t.GameLogic.(*BRNNLogic)
				brnn.changeZhuShui(int32(zhushui))
			}
		}

		return "zhushui success"
	}

	return "nil"
}
