package protofile

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	ppb "example.invalid/ServerTeam/Servers/pkg/pachislot/proto-provider/protofile/protodata"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// GameProtoTransformer trans Game.bytes
type GameProtoTransformer struct {
	item *ppb.Game_ARRAY // 要轉換的物件
}

// NewGameProtoTransformer create GameProtoTransformer
func NewGameProtoTransformer() *GameProtoTransformer {
	t := new(GameProtoTransformer)
	t.item = new(ppb.Game_ARRAY)
	return t
}

// File 檔案名稱
func (t *GameProtoTransformer) File() string {
	return "Game.bytes"
}

// Item 要轉換的物件
func (t *GameProtoTransformer) Item() proto.Message {
	return t.item
}

// Transform 轉換
func (t *GameProtoTransformer) Transform(p *pachislot.Proto) error {
	if t.item == nil {
		return fmt.Errorf("GameProtoTransformer Item not set")
	}
	for _, row := range t.item.Items {
		item := new(pb.GameConfig)
		item.ID = int32(row.GetID())
		item.Name = row.GetName()
		item.PayLineGroupID = int32(row.GetPayLineGroupID())
		item.PayTableGroupID = int32(row.GetPayTableGroupID())
		item.PatternGroupID = int32(row.GetPatternGroupID())
		item.ReelStripsGroupID = int32(row.GetReelStripsGroupID())
		p.GameConfig = item
	}
	return nil
}
