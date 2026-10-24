package helper

import (
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// FilterProtoDataByPattern 過濾指定Pattern
func FilterProtoDataByPattern(pd *protodata.ProtoData, pattern int32) (*protodata.ProtoData, error) {
	return slot.FilterProtoDataByPattern(pd, pattern)
}
