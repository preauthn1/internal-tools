package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/game/infc"
	ppb "example.invalid/ServerTeam/Servers/internal/game/infc/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/internal/game/infc/protobuf-config"
)

// PatternTransformer trans Pattern.bytes
type PatternTransformer struct {
	item *ppb.Pattern_ARRAY // 要轉換的物件
}

// NewPatternTransformer create PatternTransformer
func NewPatternTransformer() *PatternTransformer {
	t := new(PatternTransformer)
	t.item = new(ppb.Pattern_ARRAY)
	return t
}

// File 檔案名稱
func (t *PatternTransformer) File() string {
	return "Pattern.bytes"
}

// Item 要轉換的物件
func (t *PatternTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *PatternTransformer) Transform(p *infc.Proto) error {
	if t.item == nil {
		return fmt.Errorf("PatternTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.Pattern)
		item.PatternID = int32(row.GetPatternID())
		item.StageGroupID = int32(row.GetStageGroupID())
		item.IsLoop = int32(row.GetIsLoop())

		p.Patterns = append(p.Patterns, item)
	}
	return nil
}
