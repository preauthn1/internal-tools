package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/game/infc"
	ppb "example.invalid/ServerTeam/Servers/internal/game/infc/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/internal/game/infc/protobuf-config"
)

// FruitEmitterTransformer trans FruitEmitter.bytes
type FruitEmitterTransformer struct {
	item *ppb.FruitEmitter_ARRAY // 要轉換的物件
}

// NewFruitEmitterTransformer create FruitEmitterTransformer
func NewFruitEmitterTransformer() *FruitEmitterTransformer {
	t := new(FruitEmitterTransformer)
	t.item = new(ppb.FruitEmitter_ARRAY)
	return t
}

// File 檔案名稱
func (t *FruitEmitterTransformer) File() string {
	return "FruitEmitter.bytes"
}

// Item 要轉換的物件
func (t *FruitEmitterTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *FruitEmitterTransformer) Transform(p *infc.Proto) error {
	if t.item == nil {
		return fmt.Errorf("FruitEmitterTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.FruitEmitter)
		item.ID = int32(row.GetID())
		item.EmitterID = int32(row.GetEmitterID())
		item.EmitterGroupID = int32(row.GetEmitterGroupID())
		item.AngleMin = row.GetAngleMin()
		item.AngleMax = row.GetAngleMax()
		item.ForceMin = row.GetForceMin()
		item.ForceMax = row.GetForceMax()
		item.IntervalMin = row.GetIntervalMin()
		item.IntervalMax = row.GetIntervalMax()
		item.FruitID = row.GetFruitID()

		p.FruitEmitters = append(p.FruitEmitters, item)
	}
	return nil
}
