package roomlogic

import (
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RoomLogicCmdP ?
var RoomLogicCmdP *RLCmd

// RLCmd ?
type RLCmd struct {
	roomlogic *RoomLogic
}

// NewRoomLogicCmd ?
func NewRoomLogicCmd(rl *RoomLogic) *RLCmd {
	cmd := new(RLCmd)
	cmd.roomlogic = rl
	return cmd
}

// Name ?
func (rl *RLCmd) Name() string {
	return "roomlogic"
}

// Help ?
func (rl *RLCmd) Help() string {
	return "RoomLogic.usage:roomlogic reload 重新載入房間設定檔logicconfig.dat "
}

// Exec ?
func (rl *RLCmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error.seeref usage"
	}

	cmdname := args[0]

	switch cmdname {
	case "reload":
		rl.ReloadConfig()
		return "done"
	case "host":
		if len(args) != 3 || args[1] != "change" {
			return "error, host usage"
		}

		if !util.IsNumber(args[2]) {
			return "error.主機狀態必須為數字"
		}

		newHost, _ := strconv.Atoi(args[2])
		rl.updateHostAndBackup(int32(newHost))

		return "update host done."
	}
	return "nil"
}

// ReloadConfig ?
func (rl *RLCmd) ReloadConfig() {
	rl.roomlogic.ReloadGameServerConfig()
}

func (rl *RLCmd) updateHostAndBackup(newHostType int32) {
	rl.roomlogic.UpdateHostAndBackUp(newHostType, "")
}

// GetRoomActionManager ?
func (rl *RLCmd) GetRoomActionManager() *RoomActionManager {
	return rl.roomlogic.roomactionmanager
}

// GetRoomActionController ?
func (rl *RLCmd) GetRoomActionController() *RoomActionController {
	return rl.roomlogic.roomactionController
}
