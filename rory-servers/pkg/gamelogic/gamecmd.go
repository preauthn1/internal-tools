package gamelogic

import (
	"bytes"
	"fmt"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameCmd 遊戲命令對象
type GameCmd struct {
	gamelogics *LogicManager
}

// NewGameCmd 創建遊戲命令對象
func NewGameCmd(logics *LogicManager) *GameCmd {
	gcmd := new(GameCmd)
	gcmd.gamelogics = logics

	return gcmd
}

// Name ?
func (cmd *GameCmd) Name() string {
	return "game"
}

// Help ?
func (cmd *GameCmd) Help() string {
	return "遊戲管理，eg:game stat | game list (0/1) | game find tableno"
}

// Exec ?
func (cmd *GameCmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error.seeref usage"
	}
	cmdname := args[0]

	switch cmdname {
	case "stat":
		gamecount := 0
		freecount := 0
		cmd.gamelogics.Range(func(g GameLogic) {
			if g.IsGaming() {
				gamecount++
			} else {
				freecount++
			}
		})
		// for _, g := range cmd.gamelogics.logics {
		// 	if g.IsGaming() {
		// 		gamecount++
		// 	} else {
		// 		freecount++
		// 	}
		// }

		return fmt.Sprintf("game TableManager count:%d,game count:%d,free count:%d", gamecount+freecount, gamecount, freecount)
	case "list":
		if len(args) != 2 {
			return "error,seeref usage"
		}

		stype := args[1]
		var buf bytes.Buffer

		cmd.gamelogics.Range(func(g GameLogic) {
			if g.IsGaming() {
				if stype == "1" {
					buf.WriteString(fmt.Sprintf("第%d桌,", g.TableNo()))
				}
			} else {
				if stype == "0" {
					buf.WriteString(fmt.Sprintf("第%d桌,", g.TableNo()))
				}
			}
		})

		// for _, g := range cmd.gamelogics.logics {
		// 	if g.IsGaming() {
		// 		if stype == "1" {
		// 			buf.WriteString(fmt.Sprintf("第%d桌,", g.TableNo()))
		// 		}
		// 	} else {
		// 		if stype == "0" {
		// 			buf.WriteString(fmt.Sprintf("第%d桌,", g.TableNo()))
		// 		}
		// 	}
		// }

		buf.WriteString("\r\n")

		return buf.String()
	case "find":
		if len(args) != 2 {
			return "error,seeref usage"
		}

		sno := args[1]

		if !util.IsNumber(sno) {
			return "error,桌號必須為數字"
		}

		tno, _ := strconv.Atoi(sno)

		glc := cmd.gamelogics.GetGameLogic(int32(tno))

		if glc == nil {
			return "error,桌號不存在"
		}

		return fmt.Sprintf("第%d桌,遊戲狀態:%v", glc.TableNo(), glc.IsGaming())

	default:
		return ""
	}

}
