package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"example.invalid/spf13/cobra"

	"example.invalid/ServerTeam/Servers/cmd"

	"example.invalid/ServerTeam/Servers/internal/game/gen"
	"example.invalid/ServerTeam/Servers/internal/game/gen/analyser"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	protos "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:               "gen",
	Short:             "將軍",
	PersistentPreRunE: cmd.PachiSlotProtoFileLoader,
	RunE:              cmd.PachiSlotCommandRunE,
}

var analyseCmd = &cobra.Command{
	Use:          "analyse",
	Short:        "分析",
	Long:         "在指定Pattern下的遊戲在執行 n 場(NG)後的結果分析",
	Example:      "gen analyse -p 600101 -t 1e6 # 使用 600101 跑 100萬場 NG",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		conf.ClearLog = true
		conf.CLogLevel = mlog.ErrorLevel
		conf.AnalysisMode = true

		path, _ := cmd.Flags().GetString("config")
		pid, _ := cmd.Flags().GetInt32("pattern")
		times, _ := cmd.Flags().GetFloat64("times")

		if pid == 0 {
			return fmt.Errorf("Please assign pattern to analyse")
		}

		// Load Config
		cfg, err := config.LoadConfigGroup(path)
		if err != nil {
			return err
		}

		app.Context.Provide(func(p protos.Provider) (*pachislot.Proto, error) {
			return p.Proto()
		})

		app.Context.Provide(func() *pachislot.MachineDistributionManager {
			mdm := pachislot.NewMachineDistributionManager(1, pid)
			return mdm
		})

		if err := app.Context.Invoke(func(proto *pachislot.Proto) error {
			f, err := os.Create(fmt.Sprintf("%s-%d-%s-report.txt", "將軍", pid, time.Now().Format("2006-01-02T150405Z0700")))
			if err != nil {
				return err
			}
			w := io.MultiWriter(cmd.OutOrStdout(), f)
			analyser.Run(w, cfg, proto, pid, int(times))
			return nil
		}); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	cobra.MousetrapHelpText = ""

	roomserver.GameLogicCreator = gen.GameLogicCreator
	roomlogic.RoomTableInfoCreator = gen.RoomTableInfoCreator

	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(analyseCmd)

	rootCmd.PersistentFlags().StringP("config", "c", "config", "config folder")
	rootCmd.PersistentFlags().String("proto", "protodata", "Proto Configurations")

	analyseCmd.PersistentFlags().Float64P("times", "t", 100, "執行 n 場 NG")
	analyseCmd.PersistentFlags().Int32P("pattern", "p", 0, "指定的 PatternID")
}
