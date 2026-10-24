/*
基於時間定時切割的日誌類
*/

package mlog

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// NewTimeRollingLogger 初始化
func NewTimeRollingLogger(rollingInterval int, pathname string) (*TimeRollingLogger, error) {
	// 創建日誌對象函數
	var f = func() *FileLogger {
		var t = time.Now()
		filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
			t.Year(),
			t.Month(),
			t.Day(),
			t.Hour(),
			t.Minute(),
			t.Second())

		ff, _ := NewFileLogger(path.Join(pathname, filename))

		return ff
	}

	var trlogger = new(TimeRollingLogger)

	trlogger.filelastacctime = time.Now()
	trlogger.filelocker = new(sync.RWMutex)
	trlogger.logintervalmin = rollingInterval
	trlogger.loggerCreateFunc = f
	// 創建日誌對象實例
	trlogger.innerLogger = trlogger.loggerCreateFunc()

	// 檢查日誌切割
	go func() {

		defer util.PanicErrStack()

		for {

			time.Sleep(time.Minute * time.Duration(trlogger.logintervalmin))

			ll := trlogger.loggerCreateFunc()
			trlogger.filelocker.Lock()
			if trlogger.innerLogger != nil {
				trlogger.innerLogger.Close()
			}

			trlogger.innerLogger = ll
			trlogger.filelocker.Unlock()
		}
	}()

	return trlogger, nil
}

// TimeRollingLogger 基於時間定時切割的日誌類
type TimeRollingLogger struct {
	filelocker       *sync.RWMutex      // 同步鎖
	filelastacctime  time.Time          // 最後一次日誌切割時間
	logintervalmin   int                // 日誌切割時間間隔，單位分鐘
	loggerCreateFunc func() *FileLogger // 創建日誌實例函數
	innerLogger      *FileLogger        // 內部的日誌對象實例
}

// WriteLog 輸出日誌
func (logger *TimeRollingLogger) WriteLog(lvl int, format string, a ...interface{}) {
	if lvl < conf.FLogLevel {
		return
	}
	logger.filelocker.RLock()
	defer logger.filelocker.RUnlock()

	if logger.innerLogger == nil {
		return
	}

	logger.innerLogger.WriteLog(lvl, format, a...)
}

// Close ?
func (logger *TimeRollingLogger) Close() {
	if logger.innerLogger != nil {
		logger.filelocker.Lock()
		defer logger.filelocker.Unlock()

		logger.innerLogger.Close()
	}
}

// File 取得使用檔案
func (logger *TimeRollingLogger) File() *os.File {
	return logger.innerLogger.innerFile
}
