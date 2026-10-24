package robot

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"example.invalid/ServerTeam/Servers/internal/game/inrm"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

var cfg *config.Group

// InitCfg 抓出設定檔
func InitCfg() error {
	// 檢查config位置
	var data *config.Group
	_, file, _, _ := runtime.Caller(0)
	var defaultPaths = []string{"./config", "../config", "../../config", path.Dir(file) + "/config", path.Dir(file) + "/../../config"}

	data = new(config.Group)

	for _, path := range defaultPaths {
		fi, err := os.Stat(path)
		if err != nil {
			mlog.Debug("Find config folder %s Error", path)
			continue
		}

		if !fi.IsDir() {
			mlog.Debug("%s is not a folder", path)
			continue
		}

		data.LoadConfig(path, 0)

		if len(data.Configs) > 0 {
			mlog.Debug("Find config folder %s Success", path)
			break
		}
		mlog.Debug("Find config folder %s Error", path)
	}

	if data == nil {
		return fmt.Errorf("Can not find config folder, please put setting files on this folder")
	}

	if len(data.Configs) == 0 {
		return fmt.Errorf("Can not load configs, path: %s", data.DirName)
	}

	cfg = data

	return nil
}

func TestMain(m *testing.M) {
	conf.CLogLevel = mlog.WarnLevel
	conf.TimerFastMode = 1000
	conf.TestMode = true
	conf.ClearLog = true

	err := InitCfg()
	if err != nil {
		panic(err)
	}

	roomserver.GameLogicCreator = inrm.GameLogicCreator
	roomserver.RobotLogicCreator = RobotLogicCreator

	m.Run()
}
