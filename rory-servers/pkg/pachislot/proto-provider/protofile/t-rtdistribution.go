package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// RTDistributionProtoTransformer trans RTDistribution.bytes
type RTDistributionProtoTransformer struct {
	item *ppb.RTDistribution_ARRAY // 要轉換的物件
}

// NewRTDistributionProtoTransformer create RTDistributionProtoTransformer
func NewRTDistributionProtoTransformer() *RTDistributionProtoTransformer {
	t := new(RTDistributionProtoTransformer)
	t.item = new(ppb.RTDistribution_ARRAY)
	return t
}

// File 檔案名稱
func (t *RTDistributionProtoTransformer) File() string {
	return "RTDistribution.bytes"
}

// Item 要轉換的物件
func (t *RTDistributionProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *RTDistributionProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("RTDistributionProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.RTDistributionOption)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())
		item.Min = int32(row.GetMin())
		item.Max = int32(row.GetMax())
		item.Weight = int32(row.GetWeight())

		p.RTDistribution = append(p.RTDistribution, item)
	}
	return nil
}
