package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

func logtailCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "監看當前日誌",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			var port int
			host, _ := cmd.Flags().GetString("host")
			port, _ = cmd.Flags().GetInt("port")

			if os.Getenv("APP_TYPE") == "apiserver" {
				port, _ = cmd.Flags().GetInt("port-api")
			} else {
				port, _ = cmd.Flags().GetInt("port")
			}

			addr := fmt.Sprintf("%s:%d", host, port)

			conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
			if err != nil {
				return fmt.Errorf("GRPC Dial Error to %s, %v", addr, err)
			}
			LogClient = pb.NewLogClient(conn)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			verbose, _ := cmd.PersistentFlags().GetBool("verbose")
			force, _ := cmd.PersistentFlags().GetBool("force")
			level, _ := cmd.PersistentFlags().GetInt32("level")

			r := new(pb.LogRequest)
			r.Verbose = verbose
			r.Force = force
			r.Setting = new(pb.LogState)
			r.Setting.Level = level

			stream, err := LogClient.Tail(cmd.Context(), r)
			if err != nil {
				return err
			}

			for {
				select {
				default:
					resp, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							return nil
						}
						return err
					}
					cmd.Printf("%s", resp.Message)
				case <-cmd.Context().Done():
					cmd.Println("Command Close")
				case <-stream.Context().Done():
					cmd.Println("Server Close Connect")
					return nil
				}
			}
		},
	}

	cmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose")
	cmd.PersistentFlags().BoolP("force", "f", false, "持續輸出日誌")
	cmd.PersistentFlags().Int32("level", 2, "日誌等級控制\t0: Trace, 1: Debug, 2: Info, 3: Warn, 4: Error, 9: Fatal, 99: System")

	return cmd
}
