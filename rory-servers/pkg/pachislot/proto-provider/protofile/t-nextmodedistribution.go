package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// NextModeDistributionProtoTransformer trans NextModeDistribution.bytes
type NextModeDistributionProtoTransformer struct {
	item *ppb.NextModeDistribution_ARRAY // 要轉換的物件
}

// NewNextModeDistributionProtoTransformer create NextModeDistributionProtoTransformer
func NewNextModeDistributionProtoTransformer() *NextModeDistributionProtoTransformer {
	t := new(NextModeDistributionProtoTransformer)
	t.item = new(ppb.NextModeDistribution_ARRAY)
	return t
}

// File 檔案名稱
func (t *NextModeDistributionProtoTransformer) File() string {
	return "NextModeDistribution.bytes"
}

// Item 要轉換的物件
func (t *NextModeDistributionProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *NextModeDistributionProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("NextModeDistributionProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.NextModeDistributionOption)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())
		item.Weight = int32(row.GetWeight())

		{ // ModeType
			v, ok := pb.ModeType_value[row.GetModeType()]
			if !ok {
				var strs []string
				for str := range pb.ModeType_value {
					strs = append(strs, str)
				}
				return fmt.Errorf("ModeType %s is not exist, accepted %v", row.GetModeType(), strs)
			}
			item.Mode = pb.ModeType(v)
		}

		p.NextModeDistribution = append(p.NextModeDistribution, item)
	}
	return nil
}
