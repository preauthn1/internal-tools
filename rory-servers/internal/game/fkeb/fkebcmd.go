package fkeb

import (
	"example.invalid/ServerTeam/Servers/pkg/util"
)

import (	
	"strconv"	
)

type fkebCmd struct {
	lhdlogic *FKEBLogic
}

//创建用户命令
func NewFKEBCmd(fkeb *FKEBLogic) *fkebCmd {
	cmd := new(fkebCmd)
	cmd.lhdlogic = fkeb

	return cmd
}

func (lc *fkebCmd) Name() string {
	return "syspool"
}

func (lc *fkebCmd) Help() string {
	return "疯狂二八.usage:syspool change (tableid) (money)/syspool listpool (tableid)/syspool listpool/syspool zhushui (tableid) (int32)/syspool he (tableid) (概率）"
}

func (lc *fkebCmd) Exec(args []string) string {
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
			fkeb := t.GameLogic.(*FKEBLogic)
			fkeb.changeSysPool(int64(money))
		} else {
			if !util.IsNumber(args[1]) {
				return "error.改变后的钱必须为数字"
			}

			money, _ := strconv.Atoi(args[1])
			for i := 0; i < len(lc.lhdlogic.GetTable().Room.Tables); i++ {
				t := lc.lhdlogic.GetTable().Room.Tables[i]
				fkeb := t.GameLogic.(*FKEBLogic)
				fkeb.changeSysPool(int64(money))
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
			fkeb := t.GameLogic.(*FKEBLogic)
			fkeb.checkSysPool()
		} else {
			alltable := lc.lhdlogic.GetTable().Room.Tables
			for _, t := range alltable {
				fkeb := t.GameLogic.(*FKEBLogic)
				fkeb.checkSysPool()
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
			fkeb := t.GameLogic.(*FKEBLogic)
			fkeb.changeZhuShui(int32(zhushui))
		} else {
			zhushui, _ := strconv.Atoi(args[1])
			for i := 0; i < len(lc.lhdlogic.GetTable().Room.Tables); i++ {
				t := lc.lhdlogic.GetTable().Room.Tables[i]
				fkeb := t.GameLogic.(*FKEBLogic)
				fkeb.changeZhuShui(int32(zhushui))
			}
		}

		return "zhushui success"
	}

	return "nil"
}
