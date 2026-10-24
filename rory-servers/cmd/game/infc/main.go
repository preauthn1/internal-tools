package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"example.invalid/spf13/cobra"

	"example.invalid/ServerTeam/Servers/cmd"

	cliapisvr "example.invalid/ServerTeam/Servers/internal/cli-api-svc"

	"example.invalid/ServerTeam/Servers/internal/game/infc"
	protos "example.invalid/ServerTeam/Servers/internal/game/infc/proto-provider"
	"example.invalid/ServerTeam/Servers/internal/game/infc/proto-provider/protofile"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "infc",
	Short: "印度切水果",
	RunE: func(cmd *cobra.Command, args []string) error {

		conf.ClearLog = true
		conf.CLogLevel = mlog.InfoLevel

		path, _ := cmd.Flags().GetString("config")
		proto, _ := cmd.Flags().GetString("proto")

		var cfg *config.Group
		if _cfg, err := config.LoadConfigGroup(path); err == nil {
			cfg = _cfg
		} else {
			return err
		}

		if err := app.Context.Provide(func() (protos.Provider, error) {
			provider, err := protofile.NewProvider(proto)
			if err != nil {
				return nil, err
			}
			for _, t := range protofile.DefaultTransformers {
				if err := provider.RegisterTransformer(t); err != nil {
					return nil, err
				}
			}
			return provider, nil
		}); err != nil {
			return err
		}

		if err := app.Context.Provide(func(provider protos.Provider) (*infc.Proto, error) {
			return provider.Proto()
		}); err != nil {
			return err
		}

		var mode infc.GameMode
		if val, err := strconv.ParseUint(cfg.GetConfigValue(gamelogic.ConfigFileNameGameServer, gamelogic.ConfigSectionNameGameServer, "gamemode"), 0, 8); err != nil {
			mode = infc.GameModeDeathGame
		} else {
			mode = infc.GameMode(val)
		}

		setupLogic(mode)

		gs, _, srvcfg, err := roomserver.SetupGameServer(cfg)
		if err != nil {
			mlog.Error("初始化服務器錯誤, %v", err)
			return err
		}

		// 啟動服務器
		err = gs.Start()
		if err != nil {
			mlog.Error("服務器啟動錯誤,%v", err)
			return err
		}

		if len(srvcfg.GRPCAddr) > 0 {
			lis, err := net.Listen("tcp", srvcfg.GRPCAddr)
			if err != nil {
				mlog.Error("初始化 GRPC 控制服務錯誤, %v", err)
				return err
			}

			grpcServer := cliapisvr.NewGRPC()
			mlog.System("啟動 GRPC 控制服務, %s", srvcfg.GRPCAddr)
			go grpcServer.Serve(lis)
		}

		select {}
	},
}

func init() {
	cobra.MousetrapHelpText = ""

	rootCmd.AddCommand(cmd.VersionCmd)

	rootCmd.PersistentFlags().StringP("config", "c", "config", "config folder")
	rootCmd.PersistentFlags().String("proto", "protodata", "protodata folder")
}

func setupLogic(mode infc.GameMode) {
	switch mode {
	case infc.GameModeDeathGame:
		roomserver.GameLogicCreator = infc.DeathGameCreater
	case infc.GameModeFixedTimeGame:
		roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
			return infc.NewFixedTimeGame(t, tfManager, cfgs)
		}
	}
}
