package mlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"example.invalid/mattn/go-colorable"

	"example.invalid/ServerTeam/Servers/pkg/conf"
)

// color setting
const (
	ColorBlack = uint8(iota + 90)
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// ConsoleLogger Console Logger
type ConsoleLogger struct {
	innerLogger *log.Logger
	enabled     bool
}

// NewConsoleLogger ceate console logger
func NewConsoleLogger() (*ConsoleLogger, error) {
	logger := new(ConsoleLogger)
	logger.innerLogger = log.New(colorable.NewColorable(os.Stdout), "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	logger.enabled = true

	return logger, nil
}

// Close remove innerLogger
func (logger *ConsoleLogger) Close() {
	logger.innerLogger = nil
}

// WriteLog write log
func (logger *ConsoleLogger) WriteLog(lvl int, format string, a ...interface{}) {
	if lvl < conf.CLogLevel {
		return
	}

	if !logger.enabled {
		return
	}

	if logger.innerLogger == nil {
		return
	}

	lvldesc := getLogLevelDesc(lvl)

	text := fmt.Sprintf(lvldesc+format, a...)

	if conf.ClearLog {
		logger.innerLogger.Output(5, text)
		return
	}

	logger.innerLogger.Output(5, logger.ColorText(lvl, text))
}

// Enabled 是否啟用
func (logger *ConsoleLogger) Enabled() bool {
	return logger.enabled
}

// Enable 設定
func (logger *ConsoleLogger) Enable(ret bool) {
	logger.enabled = ret
}

// ColorText color
func (logger *ConsoleLogger) ColorText(lvl int, text string) string {

	switch lvl {
	case DebugLevel:
		lvlColor := ColorBlue
		text = logger.HLightMessage(text, lvlColor)
		return fmt.Sprintf("\033[%dm%s\033[0m", lvlColor, text)
	case InfoLevel:
		lvlColor := ColorGreen
		text = logger.HLightMessage(text, lvlColor)
		return fmt.Sprintf("\033[%dm%s\033[0m", lvlColor, text)
	case WarnLevel:
		lvlColor := ColorYellow
		text = logger.HLightMessage(text, lvlColor)
		return fmt.Sprintf("\033[%dm%s\033[0m", lvlColor, text)
	case ErrorLevel:
		lvlColor := ColorRed
		text = logger.HLightMessage(text, lvlColor)
		return fmt.Sprintf("\033[%dm%s\033[0m", lvlColor, text)
	case FatalLevel:
		lvlColor := ColorMagenta
		text = logger.HLightMessage(text, lvlColor)
		return fmt.Sprintf("\033[%dm%s\033[0m", lvlColor, text)
	//case 100: //黑色
	//	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_black, text)
	//case 101: //白色
	//	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_white, text)
	default:
		text = logger.HLightMessage(text, 0)
		return text
	}
}

// HLightMessage Message 高亮
func (logger *ConsoleLogger) HLightMessage(text string, originColor uint8) string {
	r := regexp.MustCompile(",?\"Message\":(\".*?\"),")
	strs := r.FindStringSubmatch(text)
	replace := fmt.Sprintf(",\"Message\":%s,", fmt.Sprintf("\033[1m\033[%dm%s\033[%dm", originColor, strs[1], originColor))
	text = r.ReplaceAllString(text, replace)
	return text
}

// Writer 當前 Writer
func (logger *ConsoleLogger) Writer() io.Writer {
	return logger.innerLogger.Writer()
}

// SetWriter 更新 Writer
func (logger *ConsoleLogger) SetWriter(w io.Writer) {
	logger.innerLogger.SetOutput(w)
}
