package mlog

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
)

// Log Levels
const (
	TraceLevel  = iota // Trace
	DebugLevel         // Debug
	InfoLevel          // Info
	WarnLevel          // Warn
	ErrorLevel         // Error
	FatalLevel  = 9    // Fatal, Stop Server
	SystemLevel = 99   // System
)

// GetLogLevelString 取得 Log Level 文字
func GetLogLevelString(lv int) string {
	switch lv {
	default:
		return ""
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case SystemLevel:
		return "system"
	}
}

// GetLogLevelByString 取得 Log Level
func GetLogLevelByString(lv string) int {
	switch lv {
	default:
		return 2
	case "trace":
		return TraceLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	case "system":
		return SystemLevel
	}
}

// Logger 日誌接口定義
type Logger interface {
	WriteLog(lvl int, format string, a ...interface{})
	Close()
}

//获取日志等级描述字符串
/*
日誌使用規範
開發調試的日誌皆為 debug level
用戶操作或遊戲信息為 info level
預期的錯誤為 warning level
非預期錯誤為 error level
系統信息為 system level
基本上不會調用 fatal level
*/
func getLogLevelDesc(lv int) string {

	switch lv {
	default:
		return ""
	case TraceLevel:
		return "[trace] "
	case DebugLevel:
		return "[debug] "
	case InfoLevel:
		return "[info] "
	case WarnLevel:
		return "[warn] "
	case ErrorLevel:
		return "[error] "
	case FatalLevel:
		return "[fatal] "
	case SystemLevel:
		return "[system] "
	}

}

func init() {
	if conf.Uptime {
		roomStart := time.Now()
		go func() {
			for {
				elapsed := time.Since(roomStart)
				System("Uptime %v", elapsed)
				time.Sleep(time.Duration(1) * time.Minute)
			}
		}()
	}
}
