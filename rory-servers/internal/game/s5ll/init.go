package s5ll

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
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
	protodata.OptSelectionTable,
}

// SlotBaseLogicFunc Slot遊戲 取得 SlotBaseLogic 方式
var SlotBaseLogicFunc = func(gl gamelogic.GameLogic) *baselogic.SlotLogicBase {
	return gl.(*Logic).SlotLogicBase
}

// Features 特色功能
var Features = slot.FeatureSpin | slot.FeatureFreeSpin | slot.FeatureSelection

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
