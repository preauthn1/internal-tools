package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/game/infc"
	ppb "example.invalid/ServerTeam/Servers/internal/game/infc/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/internal/game/infc/protobuf-config"
)

// StageGroupTransformer trans StageGroup.bytes
type StageGroupTransformer struct {
	item *ppb.StageGroup_ARRAY // 要轉換的物件
}

// NewStageGroupTransformer create StageGroupTransformer
func NewStageGroupTransformer() *StageGroupTransformer {
	t := new(StageGroupTransformer)
	t.item = new(ppb.StageGroup_ARRAY)
	return t
}

// File 檔案名稱
func (t *StageGroupTransformer) File() string {
	return "StageGroup.bytes"
}

// Item 要轉換的物件
func (t *StageGroupTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *StageGroupTransformer) Transform(p *infc.Proto) error {
	if t.item == nil {
		return fmt.Errorf("StageGroupTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.StageGroup)
		item.StageGroupID = int32(row.GetStageGroupID())
		item.EmitterGroupID = int32(row.GetEmitterGroupID())
		item.OrderID = int32(row.GetOrderID())
		item.GroupDelayTime = int32(row.GetGroupDelayTime())

		p.StageGroups = append(p.StageGroups, item)
	}
	return nil
}
