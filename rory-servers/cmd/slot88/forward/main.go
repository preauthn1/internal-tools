package main

import (
	"fmt"
	"net/http"
	"os"

	"example.invalid/spf13/cobra"
	"example.invalid/spf13/viper"

	slot88 "example.invalid/ServerTeam/Servers/cmd/slot88/forward/internal"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

var (
	port int
)

var rootCmd = &cobra.Command{
	Use:   "slot88-forward",
	Short: "slot88-forward transfor slot88 client websocket meassge to game server with netproto message",
	RunE: func(cmd *cobra.Command, args []string) error {

		conf.CLogLevel = mlog.GetLogLevelByString(viper.GetString("logger.level"))
		conf.ClearLog = true

		serv := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: new(slot88.Server),
		}

		mlog.System("Start forwarding on %s", serv.Addr)
		if err := serv.ListenAndServe(); err != nil {
			return err
		}

		return nil
	},
}

func init() {

	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		mlog.Fatal("Read Config File Error, %v\n", err)
	}

	example.invalid("Using config file: %s", viper.ConfigFileUsed())

	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 10090, "Listen port")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
