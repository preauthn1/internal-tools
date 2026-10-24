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

// PayLineProtoTransformer trans PayLine.bytes
type PayLineProtoTransformer struct {
	item *ppb.PayLine_ARRAY // 要轉換的物件
}

// NewPayLineProtoTransformer create PayLineProtoTransformer
func NewPayLineProtoTransformer() *PayLineProtoTransformer {
	t := new(PayLineProtoTransformer)
	t.item = new(ppb.PayLine_ARRAY)
	return t
}

// File 檔案名稱
func (t *PayLineProtoTransformer) File() string {
	return "PayLine.bytes"
}

// Item 要轉換的物件
func (t *PayLineProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *PayLineProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("PayLineProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.PayLine)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())
		item.LineID = int32(row.GetLineID())
		strs := strings.Split(row.GetPositions(), ",")
		for _, s := range strs {
			pos := new(pb.Position)
			_strs := strings.Split(s, "_")
			if val, err := strconv.Atoi(_strs[0]); err == nil {
				pos.Reel = int32(val)
			} else {
				return err
			}
			if val, err := strconv.Atoi(_strs[1]); err == nil {
				pos.Row = int32(val)
			} else {
				return err
			}
			item.Positions = append(item.Positions, pos)
		}

		p.PayLine = append(p.PayLine, item)
	}
	return nil
}
