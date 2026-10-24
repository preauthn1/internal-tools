package command

import (
	"fmt"
	"strconv"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
)

// DebugCommand 日志控制命令
type DebugCommand struct {
}

// Name 指令名稱
func (c *DebugCommand) Name() string {
	return "debug"
}

// Help 指令說明
func (c *DebugCommand) Help() string {
	return "除錯模式, 將 debug 的日誌寫入file : ex: debug time 30 時間單位是秒 | debug on/off"
}

// Exec 執行指令
func (c *DebugCommand) Exec(args []string) string {
	lenarg := len(args)
	switch lenarg {
	case 1:
		stype := args[0]
		if stype == "on" {
			conf.FLogLevel = 0
			return "debug enable"
		} else if stype == "off" {
			conf.FLogLevel = 1
			return "debug disable"
		}
	case 2:
		stype := args[0]
		ltype := args[1]

		if stype == "time" {
			s, err := strconv.ParseInt(ltype, 10, 64)
			if s > 0 && err == nil {
				go func() {
					c.debugMode(s)
				}()
			}
		}
		return ""

	default:

	}

	return "除錯模式, 將 debug level 的 log 寫入file : ex: debug time 30 | debug on/off"
}

func (c *DebugCommand) debugMode(s int64) {
	if s <= 0 {
		return
	}

	fmt.Println("====== 除錯模式己啟動 ======")
	start := time.Now()
	defer func() {
		conf.FLogLevel = 1
		fmt.Println("====== 除錯模式己結束 ======")
		fmt.Printf("除錯模式共運行了 %s 時間\n", time.Since(start))
	}()
	conf.FLogLevel = 0
	ss := s * 100
	//每10%提示一次
	for i := 1; i <= 10; i++ {
		select {
		case <-time.After(time.Duration(ss) * time.Millisecond):
			fmt.Printf("除錯模式正在進行中... %d %%\n", i*10)
		}
	}

}
