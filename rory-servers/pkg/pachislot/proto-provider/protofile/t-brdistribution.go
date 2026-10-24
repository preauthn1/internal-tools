package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// BRDistributionProtoTransformer trans BRDistribution.bytes
type BRDistributionProtoTransformer struct {
	item *ppb.BRDistribution_ARRAY // 要轉換的物件
}

// NewBRDistributionProtoTransformer create BRDistributionProtoTransformer
func NewBRDistributionProtoTransformer() *BRDistributionProtoTransformer {
	t := new(BRDistributionProtoTransformer)
	t.item = new(ppb.BRDistribution_ARRAY)
	return t
}

// File 檔案名稱
func (t *BRDistributionProtoTransformer) File() string {
	return "BRDistribution.bytes"
}

// Item 要轉換的物件
func (t *BRDistributionProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *BRDistributionProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("BRDistributionProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.BRDistributionOption)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())
		item.Weight = int32(row.GetWeight())

		{ // Type
			v, ok := pb.Scenario_value[row.GetType()]
			if !ok {
				var strs []string
				for str := range pb.Scenario_value {
					strs = append(strs, str)
				}
				return fmt.Errorf("Type %s is not exist, accepted %v", row.GetType(), strs)
			}
			item.Type = pb.Scenario(v)
		}

		p.BRDistribution = append(p.BRDistribution, item)
	}
	return nil
}
