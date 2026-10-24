package pachislot

import (
	"reflect"
	"testing"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

func TestFindPatternByID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		p, err := FindPatternByID(DefaultProto, 1)

		if p != nil {
			tt.Errorf("pattern want %v, but %v", nil, p)
		}

		if err != ErrPatternNotFound {
			tt.Errorf("err want %v, but %v", ErrPatternNotFound, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		p, err := FindPatternByID(DefaultProto, 600101)

		if p != DefaultPattern[0] {
			tt.Errorf("pattern want %v, but %v", DefaultPattern[0], p)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}

func TestFindModesByGroupID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		modes, err := FindModesByGroupID(DefaultProto, 1)

		if len(modes) != 0 {
			tt.Errorf("Mode length want %d, but %d", 0, len(modes))
		}

		if err != ErrModeConfigNotFound {
			tt.Errorf("err want %v, but %v", ErrModeConfigNotFound, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		modes, err := FindModesByGroupID(DefaultProto, DefaultPattern[0].ModeGroupID)

		expModes := []*pb.Mode{DefaultModeConfigs[0], DefaultModeConfigs[1], DefaultModeConfigs[2]}

		if !reflect.DeepEqual(expModes, modes) {
			tt.Errorf("Modes want %v, but %v", expModes, modes)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}

func TestFindModuleByID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		module, err := FindModuleByID(DefaultProto, 1)

		if module != nil {
			tt.Errorf("module want %v, but %v", nil, module)
		}

		if err != ErrModuleConfigNotFound {
			tt.Errorf("err want %v, but %v", ErrModuleConfigNotFound, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		module, err := FindModuleByID(DefaultProto, DefaultModuleConfigs[0].ID)

		exp := DefaultModuleConfigs[0]

		if module != exp {
			tt.Errorf("module want %v, but %v", exp, module)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}

func TestFindModulesByGroupID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		modules, err := FindModulesByGroupID(DefaultProto, 1)

		if len(modules) != 0 {
			tt.Errorf("Modules length want %d, but %d", 0, len(modules))
		}

		if err != ErrModuleConfigNotFound {
			tt.Errorf("err want %v, but %v", ErrModuleConfigNotFound, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		modules, err := FindModulesByGroupID(DefaultProto, DefaultPattern[0].ModuleGroupID)

		expModules := []*pb.Module{DefaultModuleConfigs[0], DefaultModuleConfigs[1], DefaultModuleConfigs[2]}

		if !reflect.DeepEqual(expModules, modules) {
			tt.Errorf("Modules want %v, but %v", expModules, modules)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}

func TestFindRTReleaseDistributionByGroupID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		options, err := FindRTReleaseDistributionByGroupID(DefaultProto, 1)

		if len(options) != 0 {
			tt.Errorf("Options length want %d, but %d", 0, len(options))
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		options, err := FindRTReleaseDistributionByGroupID(DefaultProto, DefaultPattern[0].RTReleaseGroupID)

		expOptions := DefaultRTReleaseDistribution

		if !reflect.DeepEqual(expOptions, options) {
			tt.Errorf("Options want %v, but %v", expOptions, options)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}

func TestFindNextModeDistributionByGroupID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		options, err := FindNextModeDistributionByGroupID(DefaultProto, 1)

		if len(options) != 0 {
			tt.Errorf("Options length want %d, but %d", 0, len(options))
		}

		if err != ErrNextModeDistributionNotFound {
			tt.Errorf("err want %v, but %v", ErrNextModeDistributionNotFound, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		options, err := FindNextModeDistributionByGroupID(DefaultProto, DefaultModeConfigs[0].NextModeGroupID)

		expOptions := []*pb.NextModeDistributionOption{DefaultNextModeDistribution[0], DefaultNextModeDistribution[1], DefaultNextModeDistribution[2]}

		if !reflect.DeepEqual(expOptions, options) {
			tt.Errorf("Options want %v, but %v", expOptions, options)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}

func TestFindRTDistributionByGroupID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		options, err := FindRTDistributionByGroupID(DefaultProto, 1)

		if len(options) != 0 {
			tt.Errorf("Options length want %d, but %d", 0, len(options))
		}

		if err != ErrRTDistributionNotFound {
			tt.Errorf("err want %v, but %v", ErrRTDistributionNotFound, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		options, err := FindRTDistributionByGroupID(DefaultProto, DefaultModeConfigs[0].RTGroupID)

		var expOptions []*pb.RTDistributionOption
		for i := 0; i < 23; i++ {
			expOptions = append(expOptions, DefaultRTDistribution[i])
		}

		if !reflect.DeepEqual(expOptions, options) {
			tt.Errorf("Options want %v, but %v", expOptions, options)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}

func TestFindBRDistributionByGroupID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		options, err := FindBRDistributionByGroupID(DefaultProto, 1)

		if len(options) != 0 {
			tt.Errorf("Options length want %d, but %d", 0, len(options))
		}

		if err != ErrBRDistributionNotFound {
			tt.Errorf("err want %v, but %v", ErrBRDistributionNotFound, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		options, err := FindBRDistributionByGroupID(DefaultProto, DefaultModeConfigs[0].BRGroupID)

		var expOptions []*pb.BRDistributionOption
		for i := 0; i < 2; i++ {
			expOptions = append(expOptions, DefaultBRDistribution[i])
		}

		if !reflect.DeepEqual(expOptions, options) {
			tt.Errorf("Options want %v, but %v", expOptions, options)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}

func TestFindResultDistributionByGroupID(t *testing.T) {

	t.Run("NotFound", func(tt *testing.T) {
		options, err := FindResultDistributionByGroupID(DefaultProto, 1)

		if len(options) != 0 {
			tt.Errorf("Options length want %d, but %d", 0, len(options))
		}

		if err != ErrResultDistributionNotFound {
			tt.Errorf("err want %v, but %v", ErrResultDistributionNotFound, err)
		}
	})

	t.Run("Found", func(tt *testing.T) {
		options, err := FindResultDistributionByGroupID(DefaultProto, DefaultPattern[0].NGResultGroupID)

		var expOptions []*pb.ResultDistributionOption
		for i := 0; i < 7; i++ {
			expOptions = append(expOptions, DefaultResultDistribution[i])
		}

		if !reflect.DeepEqual(expOptions, options) {
			tt.Errorf("Options want %v, but %v", expOptions, options)
		}

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}
	})
}
