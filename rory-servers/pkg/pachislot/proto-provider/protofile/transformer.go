package protofile

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
)

// Transformer 每個輸出為.bytes的檔案轉換器
type Transformer interface {
	File() string                     // 檔案名稱
	Item() proto.Message              // 要轉換的物件
	Transform(*pachislot.Proto) error // 轉換
}
