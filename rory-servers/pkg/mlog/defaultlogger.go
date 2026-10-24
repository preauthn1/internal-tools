package mlog

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

var glogger Logger
var clogger, _ = NewConsoleLogger()
var tlogger = CreateTimeRollingLogger()
var serverString string

// 初始化
func init() {
	glogger = BuildTimeRollingAndConsoleLogGroup()
}

// BuildTimeRollingAndConsoleLogGroup 構造定時切割和控制台日誌組
func BuildTimeRollingAndConsoleLogGroup() *GroupLogger {
	gglogger, _ := NewGroupLogger()

	gglogger.AddLogger(tlogger)
	gglogger.AddLogger(clogger)

	return gglogger
}

// CreateTimeRollingLogger create TimeRollingLogger
func CreateTimeRollingLogger() *TimeRollingLogger {
	tlogger, _ := NewTimeRollingLogger(1*60*24, "log")
	return tlogger
}

// SetServerString set server string
func SetServerString(server string) {
	serverString = server
}

// SetGLogger set group logger
func SetGLogger(l Logger) {
	glogger = l
}

// Trace write trace level message
func Trace(format string, a ...interface{}) {
	logger := NewDefaultLogger()
	logger.WriteLog(int(TraceLevel), format, a...)
}

// Debug write debug level message
func Debug(format string, a ...interface{}) {
	logger := NewDefaultLogger()
	logger.WriteLog(int(DebugLevel), format, a...)
}

// Info write info level message
func Info(format string, a ...interface{}) {
	logger := NewDefaultLogger()
	logger.WriteLog(int(InfoLevel), format, a...)
}

// Warn write warn level message
func Warn(format string, a ...interface{}) {
	logger := NewDefaultLogger()
	logger.WriteLog(int(WarnLevel), format, a...)
}

// Error write error level message
func Error(format string, a ...interface{}) {
	logger := NewDefaultLogger()
	logger.WriteLog(int(ErrorLevel), format, a...)
}

// Fatal write fatal level message
func Fatal(format string, a ...interface{}) {
	logger := NewDefaultLogger()
	logger.WriteLog(int(FatalLevel), format, a...)
}

// System write system level message
func System(format string, a ...interface{}) {
	logger := NewDefaultLogger()
	logger.WriteLog(int(SystemLevel), format, a...)
}

// WriteLog write message
func WriteLog(lvl int, format string, a ...interface{}) {
	logger := NewDefaultLogger()
	logger.WriteLog(lvl, format, a...)
}

// Printf 向屏幕輸出信息
func Printf(format string, a ...interface{}) {
	fmt.Printf(format+"\r\n", a...)
}

// GetCLogger 獲取控制台日誌
func GetCLogger() *ConsoleLogger {
	return clogger
}

// GetTLogger 基於時間定時切割的日誌類
func GetTLogger() *TimeRollingLogger {
	return tlogger
}

// DefaultLogger default logger
type DefaultLogger struct {
	StdLog       StdLog
	CallDepLen   int
	CallDepStart int
}

// StdLog std log
type StdLog struct {
	Server     string
	Message    string
	Data       interface{}
	StackTrace []StdStackTrace
}

// StdStackTrace stack trace info
type StdStackTrace struct {
	File     string
	FuncName string
}

/**
有資料想記錄的使用方法
mlog.NewDefaultLogger().SetData(gs.defenseCfg).Info("连接dbserver成功.")
3/14 強制調整為重複使用實體
4/27 重複使用實例雖可降低資源損耗, 但會造成多執行續的資料錯亂, 先改回來不優化
*/

// NewDefaultLogger create DefaultLogger
func NewDefaultLogger() *DefaultLogger {
	loggerInstance := new(DefaultLogger)
	loggerInstance.Reset()
	return loggerInstance
}

// Reset clean stack
func (d *DefaultLogger) Reset() {
	d.CallDepLen = 2
	d.CallDepStart = 2
	d.StdLog = StdLog{StackTrace: make([]StdStackTrace, 0)}
	//d.StdLog.StackTrace = make([]StdStackTrace, 0)
	d.StdLog.Server = serverString
}

// SetCallDepLen ?
func (d *DefaultLogger) SetCallDepLen(i int) *DefaultLogger {
	d.CallDepLen = i
	return d
}

// SetData ?
func (d *DefaultLogger) SetData(data interface{}) *DefaultLogger {
	d.StdLog.Data = data
	return d
}

// WriteLog write log
func (d *DefaultLogger) WriteLog(lvl int, format string, a ...interface{}) {
	defer util.PanicErrStack()
	if lvl < conf.CLogLevel && lvl < conf.FLogLevel {
		return
	}

	if conf.ClearLog {
		glogger.WriteLog(lvl, format, a...)
		d.Reset()
		if lvl == FatalLevel {
			os.Exit(1)
		}
		return
	}

	// 注入stdLog.Message
	d.StdLog.Message = fmt.Sprintf(format, a...)

	if !conf.ClearLog {
		// 注入stdLog.StackTrace
		for i := d.CallDepStart; i < d.CallDepStart+d.CallDepLen; i++ {
			pc, file, line, _ := runtime.Caller(i)
			f := runtime.FuncForPC(pc)

			st := StdStackTrace{
				File:     fmt.Sprintf("%v:%v", file, line),
				FuncName: f.Name(),
			}
			d.StdLog.StackTrace = append(d.StdLog.StackTrace, st)
		}

		// 將 stdLog.Data 轉成字串符
		ds, _ := json.Marshal(d.StdLog.Data)
		d.StdLog.Data = string(ds)
	}
	// 將 stdLog 轉成字串符
	dataStr, _ := json.Marshal(d.StdLog)
	newFormat := "%v"
	glogger.WriteLog(lvl, newFormat, string(dataStr))
	d.Reset()
	if lvl == FatalLevel {
		os.Exit(1)
	}
}

// GetCallStack code stack 程式func呼叫堆疊來源
func (d *DefaultLogger) GetCallStack() StdLog {
	for i := d.CallDepStart; i < d.CallDepStart+d.CallDepLen; i++ {
		pc, file, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)

		st := StdStackTrace{
			File:     fmt.Sprintf("%v:%v", file, line),
			FuncName: f.Name(),
		}

		if d.StdLog.StackTrace == nil {
			fmt.Println("is nil")

		}

		d.StdLog.StackTrace = append(d.StdLog.StackTrace, st)
	}
	//dataStr, _ := json.Marshal(d.StdLog)
	return d.StdLog
}

// Trace write trace level message
func (d *DefaultLogger) Trace(format string, a ...interface{}) {
	d.WriteLog(int(TraceLevel), format, a...)
}

// Debug write Debug level message
func (d *DefaultLogger) Debug(format string, a ...interface{}) {
	d.WriteLog(int(DebugLevel), format, a...)
}

// Info write Info level message
func (d *DefaultLogger) Info(format string, a ...interface{}) {
	d.WriteLog(InfoLevel, format, a...)
}

// Warn write Warn level message
func (d *DefaultLogger) Warn(format string, a ...interface{}) {
	d.WriteLog(WarnLevel, format, a...)
}

// Error write Error level message
func (d *DefaultLogger) Error(format string, a ...interface{}) {
	d.WriteLog(ErrorLevel, format, a...)
}

// Fatal write Fatal level message
func (d *DefaultLogger) Fatal(format string, a ...interface{}) {
	d.WriteLog(FatalLevel, format, a...)
}

// System write System level message
func (d *DefaultLogger) System(format string, a ...interface{}) {
	d.WriteLog(SystemLevel, format, a...)
}
