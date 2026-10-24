package main

import (
	"fmt"
	"os"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	"example.invalid/ServerTeam/Servers/cmd/cli/conf"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

// CLIClient 連線工具
var CLIClient pb.CLIClient

// LogClient Log工具
var LogClient pb.LogClient

var rootCmd = &cobra.Command{
	Use:     "cli",
	Short:   "Game Server 操作工具",
	Version: conf.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		addr := fmt.Sprintf("%s:%d", host, port)

		conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
		if err != nil {
			return fmt.Errorf("GRPC Dial Error to %s, %v", addr, err)
		}
		CLIClient = pb.NewCLIClient(conn)
		LogClient = pb.NewLogClient(conn)
		MessageClient = pb.NewMessageServiceClient(conn)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(logtailCmd())
	rootCmd.AddCommand(servCmd())

	if os.Getenv("APP_TYPE") == "tfserver" {
		rootCmd.AddCommand(tfConfigCmd())
	}

	if os.Getenv("APP_TYPE") == "apiserver" {
		rootCmd.AddCommand(apiCmd())
	}

	if os.Getenv("APP_TYPE") == "gameserver" {
		rootCmd.AddCommand(playerCmd())
		rootCmd.AddCommand(robotCmd())
		rootCmd.AddCommand(roomCmd())
		rootCmd.AddCommand(tableCmd())
		rootCmd.AddCommand(summaryCmd())
		rootCmd.AddCommand(messageCmd())
		rootCmd.AddCommand(timerManagerCmd())

		if os.Getenv("APP_GROUP") == "slot" {
			rootCmd.AddCommand(slotCmd())
		}
		if os.Getenv("APP_GROUP") == "pachislot" {
			rootCmd.AddCommand(pachiSlotCmd())
		}
	}

	rootCmd.PersistentFlags().String("host", "localhost", "Target server host")
	rootCmd.PersistentFlags().Int("port", 3100, "Target game server port")
	rootCmd.PersistentFlags().Int("port-api", 3101, "Target api server port")
}

func main() {
	rootCmd.Execute()
}
