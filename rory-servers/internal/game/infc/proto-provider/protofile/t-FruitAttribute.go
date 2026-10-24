package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/game/infc"
	ppb "example.invalid/ServerTeam/Servers/internal/game/infc/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/internal/game/infc/protobuf-config"
)

// FruitAttributeTransformer trans FruitAttribute.bytes
type FruitAttributeTransformer struct {
	item *ppb.FruitAttribute_ARRAY // 要轉換的物件
}

// NewFruitAttributeTransformer create FruitAttributeTransformer
func NewFruitAttributeTransformer() *FruitAttributeTransformer {
	t := new(FruitAttributeTransformer)
	t.item = new(ppb.FruitAttribute_ARRAY)
	return t
}

// File 檔案名稱
func (t *FruitAttributeTransformer) File() string {
	return "FruitAttribute.bytes"
}

// Item 要轉換的物件
func (t *FruitAttributeTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *FruitAttributeTransformer) Transform(p *infc.Proto) error {
	if t.item == nil {
		return fmt.Errorf("FruitAttributeTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.FruitAttribute)
		item.ID = int32(row.GetID())
		item.GroupID = int32(row.GetGroupID())
		item.FruitName = row.GetFruitName()
		item.FruitWeight = row.GetFruitWeight()

		p.FruitAttributes = append(p.FruitAttributes, item)
	}
	return nil
}
