package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// ReelStripsProtoTransformer trans ReelStrips.bytes
type ReelStripsProtoTransformer struct {
	item *ppb.ReelStrips_ARRAY // 要轉換的物件
}

// NewReelStripsProtoTransformer create ReelStripsProtoTransformer
func NewReelStripsProtoTransformer() *ReelStripsProtoTransformer {
	t := new(ReelStripsProtoTransformer)
	t.item = new(ppb.ReelStrips_ARRAY)
	return t
}

// File 檔案名稱
func (t *ReelStripsProtoTransformer) File() string {
	return "ReelStrips.bytes"
}

// Item 要轉換的物件
func (t *ReelStripsProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *ReelStripsProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("ReelStripsProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.ReelStripSymbol)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())
		item.Reel = int32(row.GetReel())
		item.Row = int32(row.GetRow())

		{ // Symbol
			v, ok := pb.Symbol_value[row.GetSymbol()]
			if !ok {
				var strs []string
				for str := range pb.Symbol_value {
					strs = append(strs, str)
				}
				return fmt.Errorf("Symbol %s is not exist, accepted %v", row.GetSymbol(), strs)
			}
			item.Symbol = pb.Symbol(v)
		}

		p.ReelStrips = append(p.ReelStrips, item)
	}
	return nil
}
