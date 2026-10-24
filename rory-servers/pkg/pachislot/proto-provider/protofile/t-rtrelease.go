package protofile

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// RTReleaseProtoTransformer trans RTRelease.bytes
type RTReleaseProtoTransformer struct {
	item *ppb.RTRelease_ARRAY // 要轉換的物件
}

// NewRTReleaseProtoTransformer create RTReleaseProtoTransformer
func NewRTReleaseProtoTransformer() *RTReleaseProtoTransformer {
	t := new(RTReleaseProtoTransformer)
	t.item = new(ppb.RTRelease_ARRAY)
	return t
}

// File 檔案名稱
func (t *RTReleaseProtoTransformer) File() string {
	return "RTRelease.bytes"
}

// Item 要轉換的物件
func (t *RTReleaseProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *RTReleaseProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("RTReleaseProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.RTReleaseOption)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())

		{ // ResultType
			_, ok := pb.ResultType_name[int32(row.GetResultType())]
			if !ok {
				return fmt.Errorf("ResultType %d is not exist, accepted %v", row.GetResultType(), pb.ResultType_name)
			}
			item.ResultType = pb.ResultType(int32(row.GetResultType()))
		}

		item.Rate = int32(row.GetRate())
		item.Denominator = int32(row.GetDenominator())

		{ // Times
			strs := strings.Split(row.GetTimes(), ",")
			for _, str := range strs {
				str = strings.TrimSpace(str)
				if val, err := strconv.ParseFloat(str, 64); err == nil {
					item.Times = append(item.Times, int32(val))
				} else {
					return fmt.Errorf("parse Times error, %v", err)
				}
			}
		}

		{ // UseType
			v, ok := pb.Scenario_value[row.GetUseType()]
			if !ok {
				var strs []string
				for str := range pb.Scenario_value {
					strs = append(strs, str)
				}
				return fmt.Errorf("UseType %s is not exist, accepted %v", row.GetUseType(), strs)
			}
			item.UseType = pb.Scenario(v)
		}

		p.RTReleaseDistribution = append(p.RTReleaseDistribution, item)
	}
	return nil
}
