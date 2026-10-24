package infc

import (
	"encoding/json"

	pb "example.invalid/ServerTeam/Servers/internal/game/infc/protobuf-config"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// Proto 機率表設定
type Proto struct {
	Patterns        []*pb.Pattern
	FruitAttributes []*pb.FruitAttribute
	FruitEmitters   []*pb.FruitEmitter
	StageGroups     []*pb.StageGroup
}

// Clone 使用一樣的值複製一套出來
func (p *Proto) Clone() *Proto {
	dist := new(Proto)
	b, _ := json.Marshal(p)
	json.Unmarshal(b, dist)
	return dist
}

// TransforPattern 轉換 Pattern
func TransforPattern(src *pb.Pattern) *netproto.INFC_Proto_Pattern {
	if src == nil {
		return nil
	}
	dest := new(netproto.INFC_Proto_Pattern)
	dest.PatternID = src.PatternID
	dest.StageGroupID = src.StageGroupID
	dest.IsLoop = src.IsLoop
	return dest
}

// ReversePattern 轉換 Pattern
func ReversePattern(src *netproto.INFC_Proto_Pattern) *pb.Pattern {
	if src == nil {
		return nil
	}
	dest := new(pb.Pattern)
	dest.PatternID = src.PatternID
	dest.StageGroupID = src.StageGroupID
	dest.IsLoop = src.IsLoop
	return dest
}

// TransforFruitAttribute 轉換 FruitAttribute
func TransforFruitAttribute(src *pb.FruitAttribute) *netproto.INFC_Proto_FruitAttribute {
	if src == nil {
		return nil
	}
	dest := new(netproto.INFC_Proto_FruitAttribute)
	dest.ID = src.ID
	dest.GroupID = src.GroupID
	dest.FruitName = src.FruitName
	dest.FruitWeight = src.FruitWeight
	return dest
}

// ReverseFruitAttribute 轉換 FruitAttribute
func ReverseFruitAttribute(src *netproto.INFC_Proto_FruitAttribute) *pb.FruitAttribute {
	if src == nil {
		return nil
	}
	dest := new(pb.FruitAttribute)
	dest.ID = src.ID
	dest.GroupID = src.GroupID
	dest.FruitName = src.FruitName
	dest.FruitWeight = src.FruitWeight
	return dest
}

// TransforFruitEmitter 轉換 FruitEmitter
func TransforFruitEmitter(src *pb.FruitEmitter) *netproto.INFC_Proto_FruitEmitter {
	if src == nil {
		return nil
	}
	dest := new(netproto.INFC_Proto_FruitEmitter)
	dest.ID = src.ID
	dest.EmitterGroupID = src.EmitterGroupID
	dest.EmitterID = src.EmitterID
	dest.AngleMin = src.AngleMin
	dest.AngleMax = src.AngleMax
	dest.ForceMin = src.ForceMin
	dest.ForceMax = src.ForceMax
	dest.IntervalMin = src.IntervalMin
	dest.IntervalMax = src.IntervalMax
	dest.FruitID = src.FruitID
	return dest
}

// ReverseFruitEmitter 轉換 FruitEmitter
func ReverseFruitEmitter(src *netproto.INFC_Proto_FruitEmitter) *pb.FruitEmitter {
	if src == nil {
		return nil
	}
	dest := new(pb.FruitEmitter)
	dest.ID = src.ID
	dest.EmitterGroupID = src.EmitterGroupID
	dest.EmitterID = src.EmitterID
	dest.AngleMin = src.AngleMin
	dest.AngleMax = src.AngleMax
	dest.ForceMin = src.ForceMin
	dest.ForceMax = src.ForceMax
	dest.IntervalMin = src.IntervalMin
	dest.IntervalMax = src.IntervalMax
	dest.FruitID = src.FruitID
	return dest
}

// TransforStageGroup 轉換 StageGroup
func TransforStageGroup(src *pb.StageGroup) *netproto.INFC_Proto_StageGroup {
	if src == nil {
		return nil
	}
	dest := new(netproto.INFC_Proto_StageGroup)
	dest.StageGroupID = src.StageGroupID
	dest.EmitterGroupID = src.EmitterGroupID
	dest.OrderID = src.OrderID
	dest.GroupDelayTime = src.GroupDelayTime
	return dest
}

// ReverseStageGroup 轉換 StageGroup
func ReverseStageGroup(src *netproto.INFC_Proto_StageGroup) *pb.StageGroup {
	if src == nil {
		return nil
	}
	dest := new(pb.StageGroup)
	dest.StageGroupID = src.StageGroupID
	dest.EmitterGroupID = src.EmitterGroupID
	dest.OrderID = src.OrderID
	dest.GroupDelayTime = src.GroupDelayTime
	return dest
}
