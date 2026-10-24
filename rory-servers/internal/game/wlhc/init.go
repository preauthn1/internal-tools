package wlhc

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// Protodatas 所有Protodata
var Protodatas *protodata.ProtoData

// ProtoOptions 此遊戲要讀取的Proto表
var ProtoOptions = []protodata.ProtoOption{
	protodata.OptSlotPattern,
	protodata.OptMainGameTable,
	protodata.OptSlotPayTable,
	protodata.OptSlotPayLine,
	protodata.OptSlotSymbolDistribution,
	protodata.OptSymbolAttributes,
	protodata.OptFreeSpinTable,
	protodata.OptJackpotTable,
	protodata.OptRespinTable,
	protodata.OptSelectionTable,
	protodata.OptExtraRule,
	protodata.OptReelBgDistribution,
	protodata.OptBonusTable,
	protodata.OptPickGameTable,
}

// Features 特色功能
var Features = slot.FeatureSpin | slot.FeatureFreeSpin | slot.FeatureJackpotGame | slot.FeaturePickGame | slot.FeatureSelection

// InitProtoData 抓出機率表
func InitProtoData() error {
	var data *protodata.ProtoData
	var err error
	var defaultPaths = []string{"./protodata", "../../protodata", "./protoconfig", "../../protoconfig"}

	for _, path := range defaultPaths {
		data, err = protodata.NewProtoData(path, ProtoOptions...)
		if data != nil {
			mlog.Debug("Find protodata folder %s Success", path)
			break
		}
		if err != nil {
			mlog.Debug("Find protodata folder %s Error, %v", path, err)
		}
	}

	if data == nil {
		return fmt.Errorf("Can not find protodata folder, please put setting files on this folder")
	}

	err = data.Reload()
	if err != nil {
		return fmt.Errorf("Load protodata fail, %v", err)
	}

	Protodatas = data

	return nil
}
