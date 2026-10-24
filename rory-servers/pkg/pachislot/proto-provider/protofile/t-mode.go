package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// ModeProtoTransformer trans Mode.bytes
type ModeProtoTransformer struct {
	item *ppb.Mode_ARRAY // 要轉換的物件
}

// NewModeProtoTransformer create ModeProtoTransformer
func NewModeProtoTransformer() *ModeProtoTransformer {
	t := new(ModeProtoTransformer)
	t.item = new(ppb.Mode_ARRAY)
	return t
}

// File 檔案名稱
func (t *ModeProtoTransformer) File() string {
	return "Mode.bytes"
}

// Item 要轉換的物件
func (t *ModeProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *ModeProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("ModeProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.Mode)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())

		{ // ModeType
			v, ok := pb.ModeType_value[row.GetModeType()]
			if !ok {
				var strs []string
				for str := range pb.ModeType_value {
					strs = append(strs, str)
				}
				return fmt.Errorf("ModeType %s is not exist, accepted %v", row.GetModeType(), strs)
			}
			item.ModeType = pb.ModeType(v)
		}

		item.NextModeGroupID = int32(row.GetNextModeGroupID())
		item.BRGroupID = int32(row.GetBRGroupID())
		item.RTGroupID = int32(row.GetRTGroupID())

		p.ModeConfigs = append(p.ModeConfigs, item)
	}
	return nil
}
