package bjl

import (
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Cmd 百家樂指令
type Cmd struct {
	bjllogic *BaiJialeLogic
}

// NewCmd 建立指令
func NewCmd(bjl *BaiJialeLogic) *Cmd {
	cmd := new(Cmd)
	cmd.bjllogic = bjl

	return cmd
}

// Name 指令名稱
func (bc *Cmd) Name() string {
	return "syspool"
}

// Help 指令說明
func (bc *Cmd) Help() string {
	return "百家樂.usage:syspool change (talbeid) (Money)/syspool listpool (talbeid)/syspool listpool/syspool zhushui (talbeid) (int32)"
}

// Exec 執行內容
func (bc *Cmd) Exec(args []string) string {
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
				return "error.房间必须为数字"
			}

			if !util.IsNumber(args[2]) {
				return "error.改变后的钱必须为数字"
			}

			tableno, _ := strconv.Atoi(args[1])
			money, _ := strconv.Atoi(args[2])
			t := bc.bjllogic.GetTable().Room.GetTableByNo(int32(tableno))
			bjl := t.GameLogic.(*BaiJialeLogic)
			bjl.changeSysPool(int64(money))
		} else {
			money, _ := strconv.Atoi(args[1])
			for i := 0; i < len(bc.bjllogic.GetTable().Room.Tables); i++ {
				t := bc.bjllogic.GetTable().Room.Tables[i]
				bjl := t.GameLogic.(*BaiJialeLogic)
				bjl.changeSysPool(int64(money))
			}
		}

		return "change success"
	case "listpool":
		if len(args) == 2 {
			if !util.IsNumber(args[1]) {
				return "error.房间必须为数字"
			}

			tableno, _ := strconv.Atoi(args[1])
			t := bc.bjllogic.GetTable().Room.GetTableByNo(int32(tableno))
			bjl := t.GameLogic.(*BaiJialeLogic)
			bjl.checkSysPool()
		} else {
			alltable := bc.bjllogic.GetTable().Room.Tables
			for _, t := range alltable {
				bjl := t.GameLogic.(*BaiJialeLogic)
				bjl.checkSysPool()
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
			t := bc.bjllogic.GetTable().Room.GetTableByNo(int32(tableno))
			bjl := t.GameLogic.(*BaiJialeLogic)
			bjl.changeZhuShui(int32(zhushui))
		} else {
			zhushui, _ := strconv.Atoi(args[1])
			for i := 0; i < len(bc.bjllogic.GetTable().Room.Tables); i++ {
				t := bc.bjllogic.GetTable().Room.Tables[i]
				bjl := t.GameLogic.(*BaiJialeLogic)
				bjl.changeZhuShui(int32(zhushui))
			}

		}

		return "zhushui success"
	}

	return "nil"
}
