package httputil

import (
	"bytes"
	"fmt"
	"time"

	"example.invalid/spf13/cobra"
)

// APIUtilCmdP 數據轉發指令
var APIUtilCmdP *APIUtilCmd

// APIUtilCmd 數據轉發指令
type APIUtilCmd struct{}

// NewAPIUtilCmd 創建數據轉發指令
func NewAPIUtilCmd() *APIUtilCmd {
	cvcmd := new(APIUtilCmd)
	return cvcmd
}

// Name 指令名稱
func (c *APIUtilCmd) Name() string {
	return "apiutil"
}

// Help 指令說明
func (c *APIUtilCmd) Help() string {
	return "接口調用工具 : apiutil <status|output|debug>"
}

// Exec 執行指令
func (c *APIUtilCmd) Exec(args []string) string {
	cmd := c.CobraRoot()
	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ""
}

// CobraRoot 取得 cobra command
func (c *APIUtilCmd) CobraRoot() *cobra.Command {
	cmd := c.rootCommand()
	cmd.AddCommand(c.debugCommand())
	return cmd
}

func (c *APIUtilCmd) rootCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "apiutil <status|debug>",
		Short:                 "日誌相關工具",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			default:
				return fmt.Errorf("輸入參數錯誤")
			case "status":
				var buf bytes.Buffer
				buf.WriteString(history.GetStatus())
				fmt.Println(buf.String())
			case "output":
				start := time.Now()
				history.OutputHistory()
				fmt.Printf("output done, used %s\n", time.Since(start))
			}
			return nil
		},
	}
}

func (c *APIUtilCmd) debugCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "debug <on/off>",
		Short:                 "除錯模式開關",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			default:
				return fmt.Errorf("輸入參數錯誤")
			case "on":
				history.Debug = true
				cmd.Println("開啟除錯模式")
			case "off":
				history.Debug = false
				cmd.Println("關閉除錯模式")
			}

			return nil
		},
	}
}

// IOPS get IOPS
func (c *APIUtilCmd) IOPS() string {
	return history.GetIOPS()
}

// Avg get avarige time
func (c *APIUtilCmd) Avg() time.Duration {
	return history.GetAvg()
}

// Size get queue size
func (c *APIUtilCmd) Size() int {
	return len(history.queue)
}

// MinMax get min and max time
func (c *APIUtilCmd) MinMax() (min, max time.Duration) {
	max = time.Duration(0)
	min = time.Duration(0)

	for _, data := range history.HistoryList {
		if data == nil {
			continue
		}
		if min == 0 {
			min = data.ExecTime
		} else if data.ExecTime < min {
			min = data.ExecTime
		}
		if data.ExecTime > max {
			max = data.ExecTime
		}
	}
	return min, max
}

// Output history
func (c *APIUtilCmd) Output() {
	history.OutputHistory()
}
