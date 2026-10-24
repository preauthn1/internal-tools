package robot

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"example.invalid/spf13/cobra"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
	"example.invalid/ServerTeam/Servers/pkg/slot/helper"
	"example.invalid/ServerTeam/Servers/pkg/slot/mock"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

type slotLogicFunc func(gamelogic.GameLogic) *baselogic.SlotLogicBase

// ProtoOptions load proto
var ProtoOptions []protodata.ProtoOption

// SlotLogicFunc 取得SlotLogicBase
var SlotLogicFunc slotLogicFunc

// GameName 指定遊戲名稱
var GameName string

// Features 特色遊戲
var Features uint

var confpath string
var protopath string
var pattern int32

var times int
var betm int64
var betr int64
var sidx int

var randomAll bool
var randomBet bool
var randomSelection bool

var debugmod bool

var timeout time.Duration

func init() {
	AnalyserCmd.PersistentFlags().StringVarP(&confpath, "config", "c", "config", "Config folder")
	AnalyserCmd.PersistentFlags().StringVar(&protopath, "proto", "protodata", "Protodata folder")
	AnalyserCmd.PersistentFlags().Int32VarP(&pattern, "pattern", "p", -1, "Slot pattern")
	AnalyserCmd.PersistentFlags().IntVarP(&times, "times", "t", 100, "Bet spin count")
	AnalyserCmd.PersistentFlags().Int64Var(&betm, "bet-money", 0, "Bet money")
	AnalyserCmd.PersistentFlags().Int64Var(&betr, "bet-rate", 0, "Bet rate, Choose the MainGame, it can be BetRate or MainGameID")
	AnalyserCmd.PersistentFlags().IntVar(&sidx, "select-idx", 0, "Selection Index, Choose number of selection option, first is 0")
	AnalyserCmd.PersistentFlags().BoolVarP(&randomAll, "random-all", "R", false, "Random Bet/Selection, if false, each Bet/Selection will bet min or choose first option")
	AnalyserCmd.PersistentFlags().BoolVar(&randomBet, "random-bet", false, "Random Bet, if false, each Bet will bet min")
	AnalyserCmd.PersistentFlags().BoolVar(&randomSelection, "random-selection", false, "Random Selection, if false, each Selection will choose first option")
	AnalyserCmd.PersistentFlags().BoolVar(&debugmod, "debug", false, "Debug Mode, print debug message")
	AnalyserCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Minute, "Analysis timeout, if over this time will break process")
}

// AnalyserCmd 分析指令
var AnalyserCmd = &cobra.Command{
	Use:   "analysis",
	Short: "分析工具",
	RunE: func(cmd *cobra.Command, args []string) error {

		if SlotLogicFunc == nil {
			return fmt.Errorf("SlotLogicFunc is empty")
		}
		if GameName == "" {
			return fmt.Errorf("GameName is empty")
		}

		if randomAll {
			randomBet = true
			randomSelection = true
		}

		debug.SetMaxStack(1e10)

		// Load Config
		data, err := config.LoadConfigGroup(confpath)
		if err != nil {
			return err
		}
		helper.AnalysisModeConfig(data)

		// Load Proto
		Proto, err := protodata.NewProtoData(protopath, ProtoOptions...)
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

		// 每 5 秒 更新一次種子
		go func() {
			for {
				rand.Seed(time.Now().UnixNano())
				time.Sleep(5 * time.Second)
			}
		}()

		conf.CLogLevel = mlog.WarnLevel
		conf.ClearLog = true
		conf.AnalysisMode = true

		var ctx context.Context
		// var cancel context.CancelFunc

		roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
			var rbt *Analyser
			rbt, ctx, _ = NewAnalyser(table, tfManager, robotCfg, configs, glc, evtChannel, debugmod)
			rbt.GameName = GameName
			rbt.TotalTimes = times
			rbt.LeftTimes = times
			rbt.RandomBet = randomBet
			rbt.RandomSelection = randomSelection
			if betm > 0 {
				rbt.BetMoney = betm
			}
			if betr > 0 {
				rbt.BetRate = betr
			}
			rbt.SelectIdx = sidx
			rbt.Features = Features

			return rbt
		}

		grs := fixtures.NewTestGameRoomServerWithConfig(data)
		_lg := grs.GameLogicManager.GetGameLogic(0)
		logic := SlotLogicFunc(grs.GameLogicManager.GetGameLogic(0))
		rbtlogic := grs.RobotService.GetRobotLogicManager().GetRobotLogic(0).(*Analyser)

		robot := fixtures.GetRobots()[0]
		user := robot.UserInfo
		user.Score *= 1000

		logic.SendMessgeFunc = func(userid int32, m *msg.Message) {
			rbtlogic.OnGameMessage(robot, m)
		}

		var logs []*netproto.Slot_MSG_GameLogAdd

		logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
			if int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID) == m.SClassID {
				mm := mock.EmptyJackpotPoolRPCMessage(logic, user)
				logic.OnRecvRPCMessage(mm)
				return
			}
			if int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadID) == m.SClassID {
				mm := mock.EmptyFreeSpinRPCMessage(logic, user)
				logic.OnRecvRPCMessage(mm)
				return
			}
			if int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadID) == m.SClassID {
				mm := mock.EmptySelectionRPCMessage(logic, user)
				logic.OnRecvRPCMessage(mm)
				return
			}
			if int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID) == m.SClassID {
				detail := mock.EmptyBonusGameDetail(logic, user)
				mm := mock.MarshalBonusGame(detail)
				logic.OnRecvRPCMessage(mm)
				return
			}
			if int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID) == m.SClassID {
				mm := mock.EmptyDeferredAwardRPCMessage(logic, user)
				logic.OnRecvRPCMessage(mm)
				return
			}
			if int32(netproto.Slot_GameMessageClassID_Slot_LightningReadID) == m.SClassID {
				mm := mock.EmptyLightningRPCMessage(logic, user)
				logic.OnRecvRPCMessage(mm)
			}
			if int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadID) == m.SClassID {
				mm := mock.EmptyUserAssignedModeRPCMessage(logic, user)
				logic.OnRecvRPCMessage(mm)
			}
			if int32(netproto.Slot_GameMessageClassID_Slot_GameLogBatchAddID) == m.SClassID {
				logs = append(logs, m.MsgData.(*netproto.Slot_MSG_BatchGameLogAdd).GetLogs()...)
				return
			}
			if int32(netproto.Slot_GameMessageClassID_Slot_GameLogAddID) == m.SClassID {
				logs = append(logs, m.MsgData.(*netproto.Slot_MSG_GameLogAdd))
				return
			}
		}

		go func() {
			_lg.UserEnter(user, logic.GetTable().TableNo, logic.GetTable().GetFreeSeat().SeatNo)
		}()

		select {
		case <-ctx.Done():
			key := fmt.Sprintf("[%s]-%s", robot.UserInfo.NickName, time.Now().Format("2006-01-02T150405Z0700"))

			rpf, err := os.Create(fmt.Sprintf("%s-%s-report.txt", GameName, key))
			if err != nil {
				return err
			}
			defer rpf.Close()

			csvf, err := os.Create(fmt.Sprintf("%s-%s-data.csv", GameName, key))
			if err != nil {
				return err
			}
			defer csvf.Close()

			w := io.MultiWriter(cmd.OutOrStdout(), rpf)

			rbtlogic.PrintReportByLogs(w, logs)
			helper.OutputHistory(csvf, logs)

			return nil
		case <-time.After(timeout):
			return fmt.Errorf("Timeout for %v", timeout)
		}
	},
}
