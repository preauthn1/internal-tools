package protofile

import "testing"

func TestProvider(t *testing.T) {
	path := "fixtures"
	p, err := NewProvider(path)
	if err != nil {
		t.Errorf("err want %v, but %v", nil, err)
	}

	t.Run("Register Game", func(t *testing.T) {
		err := p.RegisterTransformer(NewGameProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register PayLine", func(t *testing.T) {
		err := p.RegisterTransformer(NewPayLineProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register PayTalbe", func(t *testing.T) {
		err := p.RegisterTransformer(NewPayTableProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register ReelStrips", func(t *testing.T) {
		err := p.RegisterTransformer(NewReelStripsProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register Pattern", func(t *testing.T) {
		err := p.RegisterTransformer(NewPatternProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register Module", func(t *testing.T) {
		err := p.RegisterTransformer(NewModuleProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register Mode", func(t *testing.T) {
		err := p.RegisterTransformer(NewModeProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register RTDistribution", func(t *testing.T) {
		err := p.RegisterTransformer(NewRTDistributionProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register BRDistribution", func(t *testing.T) {
		err := p.RegisterTransformer(NewBRDistributionProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register NextModeDistribution", func(t *testing.T) {
		err := p.RegisterTransformer(NewNextModeDistributionProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register RTRelease", func(t *testing.T) {
		err := p.RegisterTransformer(NewRTReleaseProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})
	t.Run("Register ResultDistribution", func(t *testing.T) {
		err := p.RegisterTransformer(NewResultDistributionProtoTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})

	t.Run("Proto", func(t *testing.T) {
		proto, err := p.Proto()
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}

		t.Log("=== Game ===")
		t.Log(proto.GameConfig)

		t.Log("=== PayLine ===")
		for _, row := range proto.PayLine {
			t.Log(row)
		}

		t.Log("=== PayTable ===")
		for _, row := range proto.PayTable {
			t.Log(row)
		}

		t.Log("=== ReelStrips ===")
		for _, row := range proto.ReelStrips {
			t.Log(row)
		}

		t.Log("=== Pattern ===")
		for _, row := range proto.Pattern {
			t.Log(row)
		}

		t.Log("=== Module ===")
		for _, row := range proto.ModuleConfigs {
			t.Log(row)
		}

		t.Log("=== Mode ===")
		for _, row := range proto.ModeConfigs {
			t.Log(row)
		}

		t.Log("=== RTDistribution ===")
		for _, row := range proto.RTDistribution {
			t.Log(row)
		}

		t.Log("=== BRDistribution ===")
		for _, row := range proto.BRDistribution {
			t.Log(row)
		}

		t.Log("=== NextModeDistribution ===")
		for _, row := range proto.NextModeDistribution {
			t.Log(row)
		}

		t.Log("=== RTReleaseDistribution ===")
		for _, row := range proto.RTReleaseDistribution {
			t.Log(row)
		}

		t.Log("=== ResultDistribution ===")
		for _, row := range proto.ResultDistribution {
			t.Log(row)
		}
	})
}
