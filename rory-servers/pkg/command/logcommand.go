package command

import (
	"fmt"
	"strconv"

	"example.invalid/spf13/cobra"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// LogCommand 日誌控制命令
type LogCommand struct {
	logtype int
}

// NewLogCommand create LogCommand
func NewLogCommand() *LogCommand {
	cmd := new(LogCommand)
	return cmd
}

// Name 指令名稱
func (c *LogCommand) Name() string {
	return "log"
}

// Help 指令說明
func (c *LogCommand) Help() string {
	return "日誌相關工具\t輸入『log help』可查詢功能說明"
	// return "设置日志类型，eg:log on/off | log get | log set 0/1/2 0 控制台和定时切割的文本 1控制台 2定时切割的文本 | log level 0/1/2/3 0调试信息 1输出日志信息 2错误 3严重错误"
}

// Exec 執行指令
func (c *LogCommand) Exec(args []string) string {
	cmd := c.CobraRoot()
	cmd.SetArgs(args)
	cmd.Execute()
	return ""
}

// CobraRoot 取得 cobra command
func (c *LogCommand) CobraRoot() *cobra.Command {
	cmd := c.rootCommand()
	cmd.AddCommand(c.levelCommand())
	cmd.AddCommand(c.liteCommand())
	cmd.AddCommand(c.setCommand())
	return cmd
}

func (c *LogCommand) rootCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "log <on/off>",
		Short:                 "日誌相關工具",
		Args:                  cobra.ExactValidArgs(1),
		ValidArgs:             []string{"on", "off"},
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			switch args[0] {
			default:
				return fmt.Errorf("輸入參數錯誤")
			case "on":
				mlog.GetCLogger().Enable(true)
				cmd.Println("開啟日誌")
			case "off":
				mlog.GetCLogger().Enable(false)
				cmd.Println("關閉日誌")
			}

			return nil
		},
	}
}

func (c *LogCommand) levelCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "level <level number>",
		Short:                 "日誌等級控制\t0: Trace, 1: Debug, 2: Info, 3: Warn, 4: Error, 9: Fatal, 99: System",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			lv, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("日誌等級輸入錯誤")
			}

			lvstr := mlog.GetLogLevelString(lv)
			if len(lvstr) == 0 {
				return fmt.Errorf("日誌等級輸入錯誤")
			}

			conf.CLogLevel = lv
			cmd.Printf("日誌等級調整為 %s\n", lvstr)
			return nil
		},
	}
}

func (c *LogCommand) liteCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "lite <on/off>",
		Short:                 "輕量化日誌訊息",
		Args:                  cobra.ExactValidArgs(1),
		ValidArgs:             []string{"on", "off"},
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("輸入參數錯誤")
			}

			switch args[0] {
			default:
				return fmt.Errorf("輸入參數錯誤")
			case "on":
				conf.ClearLog = true
				cmd.Println("簡化日誌內容")
			case "off":
				conf.ClearLog = false
				cmd.Println("還原日誌內容")
			}

			return nil
		},
	}
}

func (c *LogCommand) setCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "set <mode>",
		Short:                 "設定日誌模式\t0: 控制台和定時切割的文本, 1: 控制台, 2: 定時切割的文本",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("日誌模式輸入錯誤")
			}

			switch mode {
			default:
				return fmt.Errorf("日誌模式輸入錯誤")
			case 0:
				mlog.SetGLogger(mlog.BuildTimeRollingAndConsoleLogGroup())
			case 1:
				logger, _ := mlog.NewConsoleLogger()
				mlog.SetGLogger(logger)
			case 2:
				mlog.SetGLogger(mlog.CreateTimeRollingLogger())
			}

			cmd.Println("日誌模式設定成功")
			return nil
		},
	}
}
