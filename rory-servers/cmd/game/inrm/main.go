package main

import (
	"fmt"
	"os"

	"example.invalid/spf13/cobra"

	"example.invalid/ServerTeam/Servers/cmd"
	"example.invalid/ServerTeam/Servers/internal/game/inrm"
	"example.invalid/ServerTeam/Servers/internal/game/inrm/robot"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "inrm",
	Short: "印度拉米",
	RunE:  cmd.PokerCommandRunE,
}

func init() {
	cobra.MousetrapHelpText = ""

	cobra.MousetrapHelpText = ""

	roomserver.GameLogicCreator = inrm.GameLogicCreator
	roomserver.RobotLogicCreator = robot.RobotLogicCreator

	rootCmd.AddCommand(cmd.VersionCmd)

	rootCmd.PersistentFlags().StringP("config", "c", "config", "config folder")
}
