package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"example.invalid/spf13/cobra"

	"example.invalid/ServerTeam/Servers/cmd"
	"example.invalid/ServerTeam/Servers/internal/game/wlhc"

	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
	"example.invalid/ServerTeam/Servers/pkg/slot/defender"
	"example.invalid/ServerTeam/Servers/pkg/slot/helper"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
	"example.invalid/ServerTeam/Servers/pkg/slot/robot"
)

var path string
var proto string
var pattern int32

var rootCmd = &cobra.Command{
	Use:   "wlhc",
	Short: "五龍皇朝",
	RunE: func(cmd *cobra.Command, args []string) error {

		conf.ClearLog = true
		conf.CLogLevel = mlog.InfoLevel

		// Load Config
		data, err := config.LoadConfigGroup(path)
		if err != nil {
			return err
		}

		// Load Proto
		Proto, err := protodata.NewProtoData(proto, wlhc.ProtoOptions...)
		if err != nil {
			return err
		}
		if err := Proto.Reload(); err != nil {
			return err
		}

		// 手動設定 Slot Pattern
		if pattern != -1 {
			section := data.GetConfigSection("slotconfig", "slotsetting")
			section["slotpattern"] = strconv.Itoa(int(pattern))
		}

		// Default Slot Pattern
		{
			pid := data.GetConfigValue("slotconfig", "slotsetting", "slotpattern")
			f, err := strconv.ParseFloat(pid, 32)
			if err != nil {
				return fmt.Errorf("slotconfig.dat -> [slotsetting] -> slotpattern parse error, %v", err)
			}
			pd, err := slot.FilterProtoDataByPattern(Proto, int32(f))
			if err != nil {
				return fmt.Errorf("Load Proto Data Fail, %v", err)
			}
			baselogic.Proto = pd
		}

		// Release Slot Pattern
		{
			pid := data.GetConfigValue("slotconfig", "slotsetting", "releasepattern")
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
			pid := data.GetConfigValue("slotconfig", "slotsetting", "bitepattern")
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
			section := data.GetConfigSection("slotconfig", "jackpot-defender")
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

		roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
			return wlhc.NewLogic(t, tfManager, cfgs)
		}

		gs, roomlogic, _, err := roomserver.SetupGameServer(data)
		if err != nil {
			mlog.Error("初始化服務器錯誤, %v", err)
			return err
		}

		helpercmd := helper.NewSlotHelperCmd(gs, roomlogic)
		helpercmd.BaseLogicFunc = func(gl gamelogic.GameLogic) *baselogic.SlotLogicBase {
			return gl.(*wlhc.Logic).SlotLogicBase
		}
		command.AddCommand(helpercmd)

		// 啟動服務器
		err = gs.Start()
		if err != nil {
			mlog.Error("服務器啟動錯誤,%v", err)
			return err
		}

		// mlog.System("關閉 Logger")
		// mlog.GetCLogger().Enable(false)
		select {}
	},
}

func init() {
	cobra.MousetrapHelpText = ""

	robot.GameName = "WLHC"
	robot.ProtoOptions = wlhc.ProtoOptions
	robot.Features = wlhc.Features
	robot.SlotLogicFunc = func(gl gamelogic.GameLogic) *baselogic.SlotLogicBase {
		return gl.(*wlhc.Logic).SlotLogicBase
	}

	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return wlhc.NewLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		logic, _, _ := robot.NewRealTimeRobot(table, tfManager, robotCfg, configs, glc, evtChannel)
		return logic
	}

	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(robot.AnalyserCmd)

	rootCmd.PersistentFlags().StringVarP(&path, "config", "c", "config", "config folder")
	rootCmd.PersistentFlags().StringVar(&proto, "proto", "protodata", "protodata folder")
	rootCmd.PersistentFlags().Int32VarP(&pattern, "pattern", "p", -1, "slot pattern")
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
