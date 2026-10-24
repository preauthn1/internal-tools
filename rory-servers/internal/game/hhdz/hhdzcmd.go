package hhdz

import (
	"example.invalid/ServerTeam/Servers/pkg/util"
	"strconv"
)

type hhdzCmd struct {
	hhdzlogic *HHDZLogic
}

//创建用户命令
func NewLhdCmd(lhd *HHDZLogic) *hhdzCmd {
	cmd := new(hhdzCmd)
	cmd.hhdzlogic = lhd

	return cmd
}

func (lc *hhdzCmd) Name() string {
	return "syspool"
}

func (lc *hhdzCmd) Help() string {
	return "红黑大战.usage:syspool change (tableid) (money)/syspool listpool (tableid)/syspool listpool/syspool zhushui (tableid) (int32)/syspool he (tableid) (概率）"
}

func (lc *hhdzCmd) Exec(args []string) string {
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
			t := lc.hhdzlogic.GetTable().Room.GetTableByNo(int32(tableno))
			hhdz := t.GameLogic.(*HHDZLogic)
			hhdz.changeSysPool(int64(money))
		} else {
			if !util.IsNumber(args[1]) {
				return "error.改变后的钱必须为数字"
			}

			money, _ := strconv.Atoi(args[1])
			for i := 0; i < len(lc.hhdzlogic.GetTable().Room.Tables); i++ {
				t := lc.hhdzlogic.GetTable().Room.Tables[i]
				hhdz := t.GameLogic.(*HHDZLogic)
				hhdz.changeSysPool(int64(money))
			}
		}

		return "change success"
	case "listpool":
		if len(args) == 2 {
			if !util.IsNumber(args[1]) {
				return "error.房间必须为数字"
			}

			tableno, _ := strconv.Atoi(args[1])
			t := lc.hhdzlogic.GetTable().Room.GetTableByNo(int32(tableno))
			hhdz := t.GameLogic.(*HHDZLogic)
			hhdz.checkSysPool()
		} else {
			alltable := lc.hhdzlogic.GetTable().Room.Tables
			for _, t := range alltable {
				hhdz := t.GameLogic.(*HHDZLogic)
				hhdz.checkSysPool()
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
			t := lc.hhdzlogic.GetTable().Room.GetTableByNo(int32(tableno))
			hhdz := t.GameLogic.(*HHDZLogic)
			hhdz.changeZhuShui(int32(zhushui))
		} else {
			zhushui, _ := strconv.Atoi(args[1])
			for i := 0; i < len(lc.hhdzlogic.GetTable().Room.Tables); i++ {
				t := lc.hhdzlogic.GetTable().Room.Tables[i]
				hhdz := t.GameLogic.(*HHDZLogic)
				hhdz.changeZhuShui(int32(zhushui))
			}
		}

		return "zhushui success"
	case "he":
		if len(args) < 2 {
			return "error.eg kick 10001"
		}

		if len(args) == 3 {
			if !util.IsNumber(args[1]) {
				return "error.房间必须为数字"
			}

			if !util.IsNumber(args[2]) {
				return "error.改变后的下和的概率百分比必须为数字"
			}

			tableno, _ := strconv.Atoi(args[1])
			he, _ := strconv.Atoi(args[2])
			t := lc.hhdzlogic.GetTable().Room.GetTableByNo(int32(tableno))
			hhdz := t.GameLogic.(*HHDZLogic)
			hhdz.changeHe(int32(he))
		}

		return "zhushui success"
	}

	return "nil"
}
