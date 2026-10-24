package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/grpc"

	apipb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
	cliapisvr "example.invalid/ServerTeam/Servers/internal/cli-api-svc"

	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	protos "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider"
	"example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
	"example.invalid/ServerTeam/Servers/pkg/slot/defender"
	"example.invalid/ServerTeam/Servers/pkg/slot/helper"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func init() {
	conf.ClearLog = true
	conf.CLogLevel = mlog.InfoLevel
}

// PokerCommandRunE 棋牌遊戲基本指令
var PokerCommandRunE = func(cmd *cobra.Command, args []string) error {

	conf.ClearLog = true
	conf.CLogLevel = mlog.InfoLevel

	path, _ := cmd.Flags().GetString("config")

	// Load Config
	cfg, err := config.LoadConfigGroup(path)
	if err != nil {
		return err
	}

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
}

// SlotProtoOptions Slot遊戲 Proto 項目
var SlotProtoOptions []protodata.ProtoOption

// SlotBaseLogicFunc Slot遊戲 取得 SlotBaseLogic 方式
var SlotBaseLogicFunc func(gl gamelogic.GameLogic) *baselogic.SlotLogicBase

// SlotCommandRunE Slot遊戲基本指令
var SlotCommandRunE = func(cmd *cobra.Command, args []string) error {

	conf.ClearLog = true
	conf.CLogLevel = mlog.InfoLevel

	path, _ := cmd.Flags().GetString("config")
	proto, _ := cmd.Flags().GetString("proto")
	pattern, _ := cmd.Flags().GetInt32("pattern")

	// Load Config
	cfg, err := config.LoadConfigGroup(path)
	if err != nil {
		return err
	}

	// Load Proto
	Proto, err := protodata.NewProtoData(proto, SlotProtoOptions...)
	if err != nil {
		return err
	}
	if err := Proto.Reload(); err != nil {
		return err
	}

	// 手動設定 Slot Pattern
	if pattern != -1 {
		section := cfg.GetConfigSection("slotconfig", "slotsetting")
		section["slotpattern"] = strconv.Itoa(int(pattern))
	}

	// Default Slot Pattern
	{
		pid := cfg.GetConfigValue("slotconfig", "slotsetting", "slotpattern")
		f, err := strconv.ParseFloat(pid, 32)
		if err != nil {
			return fmt.Errorf("slotconfig.dat -> [slotsetting] -> slotpattern parse error, %v", err)
		}
		pd, err := slot.FilterProtoDataByPattern(Proto, int32(f))
		if err != nil {
			return fmt.Errorf("Load Proto Data Fail, %v", err)
		}
		baselogic.Proto = pd

		app.Context.Provide(func() *protodata.ProtoData { return pd })
	}

	// Release Slot Pattern
	{
		pid := cfg.GetConfigValue("slotconfig", "slotsetting", "releasepattern")
		f, err := strconv.ParseFloat(pid, 32)
		if err != nil {
			mlog.Warn("slotconfig.dat -> [slotsetting] -> releasepattern parse error, SKIP..., %v", err)
		} else {
			pd, err := slot.FilterProtoDataByPattern(Proto, int32(f))
			if err != nil {
				return fmt.Errorf("Load ReleaseProto Data Fail, %v", err)
			}
			baselogic.ReleaseProto = pd
		}
	}

	// Bite Slot Pattern
	{
		pid := cfg.GetConfigValue("slotconfig", "slotsetting", "bitepattern")
		f, err := strconv.ParseFloat(pid, 32)
		if err != nil {
			mlog.Warn("slotconfig.dat -> [slotsetting] -> bitepattern parse error, SKIP..., %v", err)
		} else {
			pd, err := slot.FilterProtoDataByPattern(Proto, int32(f))
			if err != nil {
				return fmt.Errorf("Load BiteProto Data Fail, %v", err)
			}
			baselogic.BiteProto = pd
		}
	}

	// 設定 Jackpot Defender
	{
		section := cfg.GetConfigSection("slotconfig", "jackpot-defender")
		if len(section) != 0 {
			var increment int64
			var target []slot.JackpotLevel
			if val, err := strconv.ParseFloat(section["increment"], 64); err == nil {
				increment = int64(val)
			} else {
				cmd.Printf("slotconfig.dat -> [jackpot-defender] -> increment parse error, %v\n", err)
			}
			strs := strings.Split(section["target"], ",")
			for _, s := range strs {
				target = append(target, slot.JackpotLevel(s))
			}
			if increment > 0 && len(target) > 0 {
				defender.JPDefender = defender.NewJackpotDefender(increment, target...)
				command.AddCommand(new(defender.JackpotDefenderCmd))
			}
		}
	}

	gs, roomlogic, srvcfg, err := roomserver.SetupGameServer(cfg)
	if err != nil {
		mlog.Error("初始化服務器錯誤, %v", err)
		return err
	}

	helper.SlotHelperCmdP = helper.NewSlotHelperCmd(gs, roomlogic)
	helper.SlotHelperCmdP.BaseLogicFunc = SlotBaseLogicFunc
	command.AddCommand(helper.SlotHelperCmdP)

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

	apiaddr := cfg.GetConfigValue(gamelogic.ConfigFileNameServerInfo, gamelogic.ConfigSectionNameServer, "apiaddr")
	if len(apiaddr) > 0 {
		conn, err := grpc.Dial(apiaddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
		if err != nil {
			return fmt.Errorf("API GRPC Dial Error to %s, %v", apiaddr, err)
		}

		app.Context.Provide(func() apipb.GameBetServiceClient {
			return apipb.NewGameBetServiceClient(conn)
		})

		app.Context.Provide(func() apipb.GameResultServiceClient {
			return apipb.NewGameResultServiceClient(conn)
		})

		app.Context.Provide(func() apipb.GamePlayerServiceClient {
			return apipb.NewGamePlayerServiceClient(conn)
		})
	}

	// fmt.Println(app.Context)

	select {}
}

// PachiSlotCommandRunE 博青斯洛遊戲基本指令
var PachiSlotCommandRunE = func(cmd *cobra.Command, args []string) error {

	path, _ := cmd.Flags().GetString("config")

	// Load Config
	cfg, err := config.LoadConfigGroup(path)
	if err != nil {
		return err
	}

	gscfg := cfg.GetConfig(gamelogic.ConfigFileNameGameServer)

	// Setup Machine Distribution
	{
		defBetOdds := gscfg.GetConfigInt64(gamelogic.ConfigSectionNameGameServer, "defaultbetodds")
		if defBetOdds == 0 {
			return fmt.Errorf("Default Bet Odds not set, please set on gameserver.dat -> [gameserver] -> defaultbetodds, ex: defaultbetodds=1000")
		}
		defPattern := gscfg.GetConfigInt32(gamelogic.ConfigSectionNameGameServer, "defaultpattern")
		if defPattern == 0 {
			return fmt.Errorf("Default Pattern not set, please set on gameserver.dat -> [gameserver] -> defaultpattern, ex: defaultpattern=600101")
		}

		mdm := pachislot.NewMachineDistributionManager(defBetOdds, defPattern)

		distConfig := gscfg.GetConfigValue("table", "distribution")
		distGrps := strings.Split(distConfig, ",")
		for _, grp := range distGrps {
			set := strings.Split(grp, ":")

			var odds int64
			var pattern int32
			var cnt int

			if val, err := strconv.ParseInt(set[0], 0, 64); err == nil {
				odds = val
			} else {
				continue
			}
			if val, err := strconv.ParseInt(set[1], 0, 64); err == nil {
				pattern = int32(val)
			} else {
				continue
			}
			if val, err := strconv.ParseInt(set[2], 0, 64); err == nil {
				cnt = int(val)
			} else {
				continue
			}
			mdm.RegisterDistribution(odds, pattern, cnt)
		}

		app.Context.Provide(func() *pachislot.MachineDistributionManager { return mdm })
	}

	// Setup Proto
	{
		err := app.Context.Provide(func(p protos.Provider) (*pachislot.Proto, error) {
			// return pachislot.DefaultProto
			return p.Proto()
		})
		if err != nil {
			return err
		}
	}

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
}

// PachiSlotProtoFileLoader 博青斯洛讀取遊戲proto指令
var PachiSlotProtoFileLoader = func(cmd *cobra.Command, args []string) error {
	proto, _ := cmd.Flags().GetString("proto")

	provider, err := protofile.NewProvider(proto)
	if err != nil {
		return err
	}

	for _, t := range protofile.DefaultTransformers {
		if err := provider.RegisterTransformer(t); err != nil {
			return err
		}
	}

	app.Context.Provide(func() protos.Provider {
		return provider
	})

	return nil
}
