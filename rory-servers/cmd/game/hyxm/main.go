package main

import (
	"fmt"
	"os"

	"example.invalid/spf13/cobra"

	"example.invalid/ServerTeam/Servers/cmd"
	"example.invalid/ServerTeam/Servers/internal/game/xmbx"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/slot/robot"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "hyxm",
	Short: "鴻運熊貓",
	RunE:  cmd.SlotCommandRunE,
}

var path string
var proto string
var pattern int32

func init() {
	cobra.MousetrapHelpText = ""

	cmd.SlotProtoOptions = xmbx.ProtoOptions
	cmd.SlotBaseLogicFunc = xmbx.SlotBaseLogicFunc

	robot.GameName = "鴻運熊貓"
	robot.ProtoOptions = xmbx.ProtoOptions
	robot.Features = xmbx.Features
	robot.SlotLogicFunc = xmbx.SlotBaseLogicFunc

	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return xmbx.NewLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		logic, _, _ := robot.NewRealTimeRobot(table, tfManager, robotCfg, configs, glc, evtChannel)
		return logic
	}

	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(robot.AnalyserCmd)

	rootCmd.PersistentFlags().StringP("config", "c", "config", "config folder")
	rootCmd.PersistentFlags().String("proto", "protodata", "protodata folder")
	rootCmd.PersistentFlags().Int32P("pattern", "p", -1, "slot pattern")
}
