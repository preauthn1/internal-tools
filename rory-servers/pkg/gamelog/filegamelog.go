package gamelog

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// FileGameLog 基於文件的日誌記錄器
type FileGameLog struct {
	*EmptyGameLog
	logDir   string
	serverid int32
}

// NewFileGameLog 創建基於文件的日誌記錄器
func NewFileGameLog(logdir string, serverid int32) *FileGameLog {
	flg := new(FileGameLog)
	flg.logDir = logdir
	flg.serverid = serverid
	flg.EmptyGameLog = NewEmptyGameLog()

	// 檢查目錄
	// flg.ensureFolder()

	return flg
}

// EndGameLog 遊戲步驟結束
func (l *FileGameLog) EndGameLog() {
	defer util.PanicErrStack()

	if l.ctx == nil {
		return
	}

	l.ctx.EndTime = time.Now()

	logdesc := l.ctx.GetGameLogDesc()

	if logdesc == "" {
		return
	}

	uids := ""
	for uid := range l.ctx.Users {
		uids += fmt.Sprintf("%d_", uid)
	}

	logfilename := fmt.Sprintf("%s\\%d\\%s%d.log", l.logDir, l.serverid, uids, time.Now().Unix())

	file, err := os.Create(logfilename)

	if err != nil {
		mlog.Error("創建牌局記錄文件錯誤,%v", err)
		return
	}

	innerLogger := log.New(file, "", log.Llongfile)
	innerLogger.Output(3, l.ctx.GetGameLogDesc())

	file.Close()
	file = nil
	innerLogger = nil

	// 清除牌局記錄數據
	l.ctx = nil
}

func (l *FileGameLog) ensureFolder() {
	defer util.PanicErrStack()

	var err error
	logdirname := path.Join(l.logDir, fmt.Sprintf("%d", l.serverid))

	if _, err = os.Stat(logdirname); os.IsNotExist(err) {
		err = os.MkdirAll(logdirname, os.ModePerm)
	}

	if err != nil {
		mlog.Error("創建牌局記錄目錄錯誤,%v", err)
	}

}
