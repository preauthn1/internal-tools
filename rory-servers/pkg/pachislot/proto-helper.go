package pachislot

import (
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// FindPatternByID 查詢Pattern
func FindPatternByID(proto *Proto, pid int32) (*pb.PatternOption, error) {
	if proto == nil {
		return nil, ErrProtoNotFound
	}
	for _, p := range proto.Pattern {
		if p.ID == pid {
			return p, nil
		}
	}
	return nil, ErrPatternNotFound
}

// FindModesByGroupID 查詢一組模式
func FindModesByGroupID(proto *Proto, gid int32) (modes []*pb.Mode, err error) {
	if proto == nil {
		return modes, ErrProtoNotFound
	}
	for _, m := range proto.ModeConfigs {
		if m.GroupID == gid {
			modes = append(modes, m)
		}
	}
	if len(modes) == 0 {
		err = ErrModeConfigNotFound
	}
	return
}

// FindModuleByID 查詢一個模組
func FindModuleByID(proto *Proto, id int32) (*pb.Module, error) {
	if proto == nil {
		return nil, ErrProtoNotFound
	}
	for _, m := range proto.ModuleConfigs {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, ErrModuleConfigNotFound
}

// FindModulesByGroupID 查詢一組模組
func FindModulesByGroupID(proto *Proto, gid int32) (modules []*pb.Module, err error) {
	if proto == nil {
		return modules, ErrProtoNotFound
	}
	for _, m := range proto.ModuleConfigs {
		if m.GroupID == gid {
			modules = append(modules, m)
		}
	}
	if len(modules) == 0 {
		err = ErrModuleConfigNotFound
	}
	return
}

// FindRTReleaseDistributionByGroupID 查詢一組RT解除機率表 (可為空)
func FindRTReleaseDistributionByGroupID(proto *Proto, gid int32) (options []*pb.RTReleaseOption, err error) {
	if proto == nil {
		return options, ErrProtoNotFound
	}
	for _, opt := range proto.RTReleaseDistribution {
		if opt.GroupID == gid {
			options = append(options, opt)
		}
	}
	return
}

// FindNextModeDistributionByGroupID 查詢一組下次模式機率表
func FindNextModeDistributionByGroupID(proto *Proto, gid int32) (options []*pb.NextModeDistributionOption, err error) {
	if proto == nil {
		return options, ErrProtoNotFound
	}
	for _, opt := range proto.NextModeDistribution {
		if opt.GroupID == gid {
			options = append(options, opt)
		}
	}
	if len(options) == 0 {
		err = ErrNextModeDistributionNotFound
	}
	return
}

// FindRTDistributionByGroupID 查詢一組RT分配機率表
func FindRTDistributionByGroupID(proto *Proto, gid int32) (options []*pb.RTDistributionOption, err error) {
	if proto == nil {
		return options, ErrProtoNotFound
	}
	for _, opt := range proto.RTDistribution {
		if opt.GroupID == gid {
			options = append(options, opt)
		}
	}
	if len(options) == 0 {
		err = ErrRTDistributionNotFound
	}
	return
}

// FindBRDistributionByGroupID 查詢一組大獎分配機率表
func FindBRDistributionByGroupID(proto *Proto, gid int32) (options []*pb.BRDistributionOption, err error) {
	if proto == nil {
		return options, ErrProtoNotFound
	}
	for _, opt := range proto.BRDistribution {
		if opt.GroupID == gid {
			options = append(options, opt)
		}
	}
	if len(options) == 0 {
		err = ErrBRDistributionNotFound
	}
	return
}

// FindResultDistributionByGroupID 查詢一組牌組機率表
func FindResultDistributionByGroupID(proto *Proto, gid int32) (options []*pb.ResultDistributionOption, err error) {
	if proto == nil {
		return options, ErrProtoNotFound
	}
	for _, opt := range proto.ResultDistribution {
		if opt.GroupID == gid {
			options = append(options, opt)
		}
	}
	if len(options) == 0 {
		err = ErrResultDistributionNotFound
	}
	return
}

// UniqResultTypes 取得唯一ResultType
func UniqResultTypes(payouts []*pb.PayoutOption) []pb.ResultType {
	resMap := make(map[pb.ResultType]bool)
	for _, d := range payouts {
		resMap[d.ResultType] = true
	}
	var types []pb.ResultType
	for rt := range resMap {
		types = append(types, rt)
	}
	return types
}
