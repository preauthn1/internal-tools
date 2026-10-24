/*
日志类
输出调试、日志、错误、崩溃信息
不提供对外的此对象的访问，可以通过各种包装类去输出日志
*/

package mlog

import (
	"fmt"
	"log"
	"os"
)

// FileLogger 日誌對象
type FileLogger struct {
	innerLogger *log.Logger
	innerFile   *os.File
}

// NewFileLogger 創建新的日誌對象
func NewFileLogger(filename string) (*FileLogger, error) {
	var innerLogger *log.Logger
	var innerFile *os.File

	file, err := os.Create(filename)

	if err != nil {
		return nil, err
	}

	innerLogger = log.New(file, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	innerFile = file

	logger := new(FileLogger)
	logger.innerLogger = innerLogger
	logger.innerFile = innerFile

	return logger, nil
}

// Close 關閉日誌
func (logger *FileLogger) Close() {
	if logger.innerLogger != nil {
		logger.innerFile.Close()
	}

	logger.innerLogger = nil
	logger.innerFile = nil
}

// WriteLog 寫日誌
func (logger *FileLogger) WriteLog(lvl int, format string, a ...interface{}) {
	if logger.innerLogger == nil {
		return
	}

	lvldesc := getLogLevelDesc(lvl)
	// 輸出日誌
	logger.innerLogger.Output(6, fmt.Sprintf(lvldesc+format+"\r\n", a...))
}

// File 取得使用檔案
func (logger *FileLogger) File() *os.File {
	return logger.innerFile
}
