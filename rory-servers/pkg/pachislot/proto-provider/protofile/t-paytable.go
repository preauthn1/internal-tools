package protofile

import (
	"fmt"
	"strconv"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// PayTableProtoTransformer trans PayTable.bytes
type PayTableProtoTransformer struct {
	item *ppb.PayTable_ARRAY // 要轉換的物件
}

// NewPayTableProtoTransformer create PayTableProtoTransformer
func NewPayTableProtoTransformer() *PayTableProtoTransformer {
	t := new(PayTableProtoTransformer)
	t.item = new(ppb.PayTable_ARRAY)
	return t
}

// File 檔案名稱
func (t *PayTableProtoTransformer) File() string {
	return "PayTable.bytes"
}

// Item 要轉換的物件
func (t *PayTableProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *PayTableProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("PayTableProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.PayoutOption)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())

		{ // ResultType
			_, ok := pb.ResultType_name[int32(row.GetResultType())]
			if !ok {
				return fmt.Errorf("ResultType %d is not exist, accepted %v", row.GetResultType(), pb.ResultType_name)
			}
			item.ResultType = pb.ResultType(int32(row.GetResultType()))
		}

		{ // PayType
			v, ok := pb.PayType_value[row.GetPayType()]
			if !ok {
				var strs []string
				for str := range pb.PayType_value {
					strs = append(strs, str)
				}
				return fmt.Errorf("PayType %s is not exist, accepted %v", row.GetPayType(), strs)
			}
			item.PayType = pb.PayType(v)
		}

		{ // Payout
			if val, err := strconv.ParseFloat(row.GetPayout(), 64); err == nil {
				item.Payout = int32(val)
			} else {
				return fmt.Errorf("parse payout error, %v", err)
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

		p.PayTable = append(p.PayTable, item)
	}
	return nil
}
