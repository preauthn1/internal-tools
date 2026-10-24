package embed

import (
	"os"
	"path/filepath"
	"testing"
)

// 沒有模型檔就跳過（模型不進 git，見 README「多語言句子向量」）
func load(t *testing.T) *Model {
	t.Helper()
	dir := ModelDir()
	if _, err := os.Stat(filepath.Join(dir, "model.onnx")); err != nil {
		t.Skipf("沒有模型檔（%s），跳過", dir)
	}
	m, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { m.Close() })
	return m
}

func vectors(t *testing.T, m *Model, texts ...string) [][]float32 {
	t.Helper()
	v, err := m.Embed(texts)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != len(texts) {
		t.Fatalf("%d 句只算出 %d 個向量", len(texts), len(v))
	}
	return v
}

func TestParaphraseCloserThanUnrelated(t *testing.T) {
	m := load(t)
	v := vectors(t, m, "你們幾點開門", "請問營業時間是幾點到幾點", "幫我寫一首關於海的詩")
	same, other := Cosine(v[0], v[1]), Cosine(v[0], v[2])
	if same <= other {
		t.Errorf("同義句 %.3f 應該比無關句 %.3f 像", same, other)
	}
}

// 一句的耗時：go -C go test -bench . -run '^$' ./internal/embed/
func BenchmarkEmbedOne(b *testing.B) {
	dir := ModelDir()
	if _, err := os.Stat(filepath.Join(dir, "model.onnx")); err != nil {
		b.Skipf("沒有模型檔（%s），跳過", dir)
	}
	m, err := Load(dir)
	if err != nil {
		b.Fatal(err)
	}
	defer m.Close()
	for b.Loop() {
		if _, err := m.Embed([]string{"請問你們週末幾點開門？"}); err != nil {
			b.Fatal(err)
		}
	}
}

// 日文、英文訊息查中文案例：意思相同的那句排第一（Apple 的向量只能比同一種文字系統，做不到）。
// 很短、語氣類的句子不一定找得到（「キャンセルでお願いします」最像「我想訂位」），抽查結果見 docs/研究總結.md 第 15 節
func TestCrossLanguageFindsChineseExample(t *testing.T) {
	m := load(t)
	chinese := []string{"我想訂位", "謝謝你", "你們幾點開門", "有停車位嗎", "股票會漲嗎"}
	queries := map[string]string{
		"予約したいです":                  "我想訂位",
		"ありがとうございます":               "謝謝你",
		"What time do you open?":   "你們幾點開門",
		"Is there parking nearby?": "有停車位嗎",
		"株価は上がりますか":                "股票會漲嗎",
	}
	cv := vectors(t, m, chinese...)
	for q, want := range queries {
		qv := vectors(t, m, q)[0]
		best := 0
		for i := range cv {
			if Cosine(qv, cv[i]) > Cosine(qv, cv[best]) {
				best = i
			}
		}
		if chinese[best] != want {
			t.Errorf("%q 最像的是 %q，應該是 %q", q, chinese[best], want)
		}
	}
}
