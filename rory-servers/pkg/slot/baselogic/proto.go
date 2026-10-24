package baselogic

import (
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// PatternMode Pattern 模式
type PatternMode uint8

// PatternModes
const (
	// 預設 Pattern
	DefaultPattern PatternMode = iota
	// 放分模式
	ReleasePattern
	// 咬分模式
	BitePattern
)

func (mode PatternMode) String() string {
	switch mode {
	default:
		return "預設"
	case ReleasePattern:
		return "放分"
	case BitePattern:
		return "殺分"
	}
}

// Proto 所有Protodata
var Proto *protodata.ProtoData

// ReleaseProto 所有Protodata
var ReleaseProto *protodata.ProtoData

// BiteProto 所有Protodata
var BiteProto *protodata.ProtoData
