package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// PatternProtoTransformer trans Pattern.bytes
type PatternProtoTransformer struct {
	item *ppb.Pattern_ARRAY // 要轉換的物件
}

// NewPatternProtoTransformer create PatternProtoTransformer
func NewPatternProtoTransformer() *PatternProtoTransformer {
	t := new(PatternProtoTransformer)
	t.item = new(ppb.Pattern_ARRAY)
	return t
}

// File 檔案名稱
func (t *PatternProtoTransformer) File() string {
	return "Pattern.bytes"
}

// Item 要轉換的物件
func (t *PatternProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *PatternProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("PatternProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.PatternOption)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())
		item.ModuleGroupID = int32(row.GetModuleGroupID())
		item.ModeGroupID = int32(row.GetModeGroupID())
		item.NGResultGroupID = int32(row.GetNGResultGroupID())
		item.RTReleaseGroupID = int32(row.GetRTReleaseGroupID())

		p.Pattern = append(p.Pattern, item)
	}
	return nil
}
