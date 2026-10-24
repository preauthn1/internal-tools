package protofile

import (
	"io/ioutil"
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
)

func TestTGame(t *testing.T) {
	tfm := NewGameProtoTransformer()
	if tfm.item == nil {
		t.Errorf("item want *ppb.Game_ARRAY, but nil")
	}

	t.Run("File", func(t *testing.T) {
		if tfm.File() != "Game.bytes" {
			t.Errorf("File() want %s, but %s", "Game.bytes", tfm.File())
		}
	})

	t.Run("Item", func(t *testing.T) {
		if tfm.Item() != tfm.item {
			t.Errorf("Item() want %s, but %s", tfm.item, tfm.File())
		}
	})

	t.Run("Load", func(t *testing.T) {
		path := "fixtures/" + tfm.File()
		b, _ := ioutil.ReadFile(path)
		proto.Unmarshal(b, tfm.Item())
		for _, row := range tfm.item.Items {
			t.Log(row)
		}
	})

	t.Run("Transform", func(t *testing.T) {
		p := new(pachislot.Proto)
		err := tfm.Transform(p)
		if err != nil {
			t.Errorf("Transform() err want %v, but %v", nil, err)
		}
		t.Log(p.GameConfig)
	})
}
