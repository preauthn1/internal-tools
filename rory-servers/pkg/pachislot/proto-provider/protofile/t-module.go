package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// ModuleProtoTransformer trans Module.bytes
type ModuleProtoTransformer struct {
	item *ppb.Module_ARRAY // 要轉換的物件
}

// NewModuleProtoTransformer create ModuleProtoTransformer
func NewModuleProtoTransformer() *ModuleProtoTransformer {
	t := new(ModuleProtoTransformer)
	t.item = new(ppb.Module_ARRAY)
	return t
}

// File 檔案名稱
func (t *ModuleProtoTransformer) File() string {
	return "Module.bytes"
}

// Item 要轉換的物件
func (t *ModuleProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *ModuleProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("ModuleProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.Module)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())

		{ // Type
			v, ok := pb.Scenario_value[row.GetType()]
			if !ok {
				var strs []string
				for str := range pb.Scenario_value {
					strs = append(strs, str)
				}
				return fmt.Errorf("UseType %s is not exist, accepted %v", row.GetType(), strs)
			}
			item.Type = pb.Scenario(v)
		}

		item.ResultProbGroupID = int32(row.GetResultProbGroupID())
		item.Count = int32(row.GetCount())
		item.Cost = int32(row.GetCost())

		p.ModuleConfigs = append(p.ModuleConfigs, item)
	}
	return nil
}
