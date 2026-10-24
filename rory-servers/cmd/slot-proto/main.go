package main

import (
	"fmt"
	"os"

	"example.invalid/spf13/cobra"

	"example.invalid/ServerTeam/Servers/pkg/slot/helper"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

var path string
var proto string
var pattern int32
var opt int

var rootCmd = &cobra.Command{
	Use:   "slot-proto",
	Short: "Slot protodata 工具",
	RunE: func(cmd *cobra.Command, args []string) error {

		option := protodata.ProtoOption(opt)

		// Load Proto
		Proto, err := protodata.NewProtoData(proto, option)
		if err != nil {
			return err
		}
		if err := Proto.Reload(); err != nil {
			return err
		}

		// 手動設定 Slot Pattern
		if pattern != -1 {
			_proto, err := helper.FilterProtoDataByPattern(Proto, pattern)
			if err != nil {
				return err
			}
			Proto = _proto
		}

		Proto.Print()

		return nil
	},
}

func init() {
	cobra.MousetrapHelpText = ""

	rootCmd.PersistentFlags().StringVar(&proto, "proto", "protodata", "protodata folder")
	rootCmd.PersistentFlags().Int32VarP(&pattern, "pattern", "p", -1, "slot pattern")
	rootCmd.PersistentFlags().IntVarP(&opt, "target", "t", 0, "proto option")
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
