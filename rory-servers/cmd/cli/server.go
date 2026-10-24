package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	pb "example.invalid/ServerTeam/Servers/internal/cli-api-svc/cliapisvc"
)

func servCmd() *cobra.Command {
	servCmd := &cobra.Command{
		Use:                   "serv",
		Short:                 "Server管理",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	servCmd.AddCommand(verCmd())

	if os.Getenv("APP_TYPE") == "gameserver" {
		servCmd.AddCommand(cverCmd())
		servCmd.AddCommand(shotdownCmd())
	}

	return servCmd
}

func verCmd() *cobra.Command {
	verCmd := &cobra.Command{
		Use:                   "ver",
		Short:                 "查詢版號",
		DisableFlagsInUseLine: true,
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

			CLIClient = pb.NewCLIClient(conn)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := CLIClient.ServerVer(cmd.Context(), new(pb.Empty))
			if err != nil {
				return err
			}

			cmd.Printf("=== 版本資訊 ===\n")
			cmd.Printf("當前版本：%s GitHash：%s\n", v.Version, v.Hash)
			cmd.Printf("建置時間：%s\n", v.BuildTime)
			cmd.Printf("建置編號：%s\n", v.BuildName)

			return nil
		},
	}
	return verCmd
}

func cverCmd() *cobra.Command {
	examples := []string{
		"cver\t\t\t// 查詢 Client 版號",
		"cver ios 1.2.1\t\t// 更新 IOS 版號",
		"cver android 3.2.1\t// 更新 Android 版號",
	}

	cverCmd := &cobra.Command{
		Use:                   "cver <platform> <version>",
		Short:                 "查詢/更新 Client 版號",
		Example:               strings.Join(examples, "\n"),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				v, err := CLIClient.ClientVer(cmd.Context(), new(pb.Empty))
				if err != nil {
					return err
				}

				cmd.Printf("===客戶端版本資訊===\n")
				cmd.Printf("IOS Version：%s\n", v.IOS)
				cmd.Printf("Android Version：%s\n", v.Android)

				return nil
			}

			if len(args) != 2 {
				return fmt.Errorf("請輸入版本號")
			}

			switch strings.ToLower(args[0]) {
			default:
				return fmt.Errorf("Client 版本輸入錯誤, 請輸入 [ios/android]")
			case "ios":
				_, err := CLIClient.SetClientVer(cmd.Context(), &pb.CVersion{IOS: args[1]})
				if err != nil {
					return err
				}

				cmd.Println("變更 ios 版號成功")
			case "android":
				_, err := CLIClient.SetClientVer(cmd.Context(), &pb.CVersion{Android: args[1]})
				if err != nil {
					return err
				}

				cmd.Println("變更 android 版號成功")
			}

			return nil
		},
	}

	return cverCmd
}

func shotdownCmd() *cobra.Command {
	shotdownCmd := &cobra.Command{
		Use:   "shotdown <message>",
		Short: "關閉 Server",
		RunE: func(cmd *cobra.Command, args []string) error {

			var message string = "Shotdown Server"

			if len(args) > 0 {
				message = args[0]
			}

			dur, err := cmd.PersistentFlags().GetDuration("delay")
			if err != nil {
				return nil
			}

			if _, err := CLIClient.Shotdown(cmd.Context(), &pb.ShotdownRequest{Message: message, Delay: dur.String()}); err != nil {
				return err
			}

			cmd.Println("Server即將關閉")

			return nil
		},
	}
	shotdownCmd.PersistentFlags().DurationP("delay", "d", 1*time.Second, "延遲關閉時間")
	return shotdownCmd
}
