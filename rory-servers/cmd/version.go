package cmd

import (
	"example.invalid/spf13/cobra"
	"example.invalid/ServerTeam/Servers/pkg/conf"
)

// VersionCmd 顯示版本資訊
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "顯示版本資訊",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(conf.Version())
	},
}
