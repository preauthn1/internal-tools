package protofile

import "testing"

func TestProvider(t *testing.T) {
	path := "fixtures"
	p, err := NewProvider(path)
	if err != nil {
		t.Errorf("err want %v, but %v", nil, err)
	}

	t.Run("Register Pattern", func(t *testing.T) {
		err := p.RegisterTransformer(NewPatternTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})

	t.Run("Register FruitAttribute", func(t *testing.T) {
		err := p.RegisterTransformer(NewFruitAttributeTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})

	t.Run("Register FruitEmitter", func(t *testing.T) {
		err := p.RegisterTransformer(NewFruitEmitterTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})

	t.Run("Register StageGroup", func(t *testing.T) {
		err := p.RegisterTransformer(NewStageGroupTransformer())
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}
	})

	t.Run("Proto", func(t *testing.T) {
		proto, err := p.Proto()
		if err != nil {
			t.Errorf("err want %v, but %v", nil, err)
		}

		t.Log("=== Pattern ===")
		for _, row := range proto.Patterns {
			t.Log(row)
		}
		t.Log("=== FruitAttribute ===")
		for _, row := range proto.FruitAttributes {
			t.Log(row)
		}
		t.Log("=== FruitEmitter ===")
		for _, row := range proto.FruitEmitters {
			t.Log(row)
		}
		t.Log("=== StageGroup ===")
		for _, row := range proto.StageGroups {
			t.Log(row)
		}
	})
}
