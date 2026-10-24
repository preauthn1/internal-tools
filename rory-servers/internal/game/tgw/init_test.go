package tgw

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
)

var cfg *config.Group

// InitCfg 抓出設定檔
func InitCfg() error {
	// 檢查config位置
	var data *config.Group
	_, file, _, _ := runtime.Caller(0)
	var defaultPaths = []string{"./config", "../../config", path.Dir(file) + "/config", path.Dir(file) + "/../../config"}

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
	conf.TimerFastMode = 1000 //40
	conf.TestMode = true
	conf.ClearLog = true
	conf.AnalysisMode = false

	InitCfg()
	InitProtoData()

	// Default Slot Pattern
	{
		pid := cfg.GetConfigValue("slotconfig", "slotsetting", "slotpattern")
		f, _ := strconv.ParseFloat(pid, 32)
		pd, _ := slot.FilterProtoDataByPattern(Protodatas, int32(f))
		baselogic.Proto = pd
	}

	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		logic := NewLogic(t, tfManager, cfgs)
		tmp := *logic.Proto
		logic.Proto = &tmp
		return logic
	}

	m.Run()
}

// func TestProto(t *testing.T) {
// 	Protodatas.Print()
// }
