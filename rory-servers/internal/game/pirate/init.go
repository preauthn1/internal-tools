package pirate

// import (
// 	"fmt"
// 	"os"
// 	"path"
// 	"runtime"

// 	"example.invalid/ServerTeam/Servers/pkg/config"
// 	"example.invalid/ServerTeam/Servers/pkg/mlog"
// 	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
// )

// var cfg *config.ConfigGroup
// var protodatas *protodata.ProtoData

// func init() {

// 	// 檢查config位置
// 	if err := initCfg(); err != nil {
// 		fmt.Fprintf(os.Stderr, "Load config files error. %v", err)
// 		os.Exit(1)
// 	}

// 	// 檢查protodata位置
// 	if err := initProtoData(); err != nil {
// 		fmt.Fprintf(os.Stderr, "Load protodata files error. %v", err)
// 		os.Exit(1)
// 	}

// }

// func initCfg() error {
// 	// 檢查config位置
// 	var data *config.ConfigGroup
// 	_, file, _, _ := runtime.Caller(0)
// 	var defaultPaths = []string{"./config", "../../config", path.Dir(file) + "/config", path.Dir(file) + "/../../config"}

// 	data = new(config.ConfigGroup)

// 	for _, path := range defaultPaths {
// 		fi, err := os.Stat(path)
// 		if err != nil {
// 			mlog.Debug("Find config folder %s Error", path)
// 			continue
// 		}

// 		if !fi.IsDir() {
// 			mlog.Debug("%s is not a folder", path)
// 			continue
// 		}

// 		data.LoadConfig(path, 0)

// 		if len(data.Configs) > 0 {
// 			mlog.Debug("Find config folder %s Success", path)
// 			break
// 		}
// 		mlog.Debug("Find config folder %s Error", path)
// 	}

// 	if data == nil {
// 		return fmt.Errorf("Can not find config folder, please put setting files on this folder")
// 	}

// 	if len(data.Configs) == 0 {
// 		return fmt.Errorf("Can not load configs, path: %s", data.DirName)
// 	}

// 	cfg = data

// 	return nil
// }

// func initProtoData() error {
// 	var data *protodata.ProtoData
// 	var err error
// 	var defaultPaths = []string{"./protodata", "../../protodata", "./protoconfig", "../../protoconfig"}
// 	var protoOptions = []protodata.ProtoOption{
// 		protodata.OptSlotPattern,
// 		protodata.OptMainGameTable,
// 		protodata.OptSlotPayLine,
// 		protodata.OptSlotPayTable,
// 		protodata.OptSlotSymbolDistribution,
// 		protodata.OptSymbolAttributes,
// 		protodata.OptFreeSpinTable,
// 		protodata.OptJackpotTable,
// 	}

// 	for _, path := range defaultPaths {
// 		data, err = protodata.NewProtoData(path, protoOptions...)
// 		if data != nil {
// 			mlog.Debug("Find protodata folder %s Success", path)
// 			break
// 		}
// 		if err != nil {
// 			mlog.Debug("Find protodata folder %s Error, %v", path, err)
// 		}
// 	}

// 	if data == nil {
// 		return fmt.Errorf("Can not find protodata folder, please put setting files on this folder")
// 	}

// 	err = data.Reload()
// 	if err != nil {
// 		return fmt.Errorf("Load protodata fail, %v", err)
// 	}

// 	protodatas = data

// 	return nil
// }
