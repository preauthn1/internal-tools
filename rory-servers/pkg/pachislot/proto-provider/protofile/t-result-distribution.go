package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// ResultDistributionProtoTransformer trans ResultDistribution.bytes
type ResultDistributionProtoTransformer struct {
	item *ppb.ResultDistribution_ARRAY // 要轉換的物件
}

// NewResultDistributionProtoTransformer create ResultDistributionProtoTransformer
func NewResultDistributionProtoTransformer() *ResultDistributionProtoTransformer {
	t := new(ResultDistributionProtoTransformer)
	t.item = new(ppb.ResultDistribution_ARRAY)
	return t
}

// File 檔案名稱
func (t *ResultDistributionProtoTransformer) File() string {
	return "ResultDistribution.bytes"
}

// Item 要轉換的物件
func (t *ResultDistributionProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *ResultDistributionProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("ResultDistributionProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.ResultDistributionOption)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())

		{ // ResultType
			_, ok := pb.ResultType_name[int32(row.GetResult())]
			if !ok {
				return fmt.Errorf("Result %d is not exist, accepted %v", row.GetResult(), pb.ResultType_name)
			}
			item.Result = pb.ResultType(int32(row.GetResult()))
		}

		item.Weight = int32(row.GetWeight())

		p.ResultDistribution = append(p.ResultDistribution, item)
	}
	return nil
}
