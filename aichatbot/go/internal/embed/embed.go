// Package embed 算多語言句子向量（intfloat/multilingual-e5-small 的 ONNX 版，用 hugot；預設純 Go 後端，不用 cgo、不用 onnxruntime），
// 給意圖案例庫找「和顧客訊息最像的幾句」。中文、日文、英文在同一個向量空間，不同語言的句子也能直接比。
// 技術驗證的量測結果見 docs/研究總結.md。
package embed

import (
	"math"
	"os"
	"path/filepath"

	"example.invalid/knights-analytics/hugot"
	"example.invalid/knights-analytics/hugot/pipelines"
)

// ModelDir 是模型資料夾：BOT_EMBED_MODEL，沒設就用 ~/.cache/aichatbot/models/multilingual-e5-small
// （裡面要有 model.onnx、tokenizer.json，就是 Hugging Face 上該模型的 onnx/ 資料夾）
func ModelDir() string {
	if d := os.Getenv("BOT_EMBED_MODEL"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "aichatbot", "models", "multilingual-e5-small")
}

type Model struct {
	session  *hugot.Session
	pipeline *pipelines.FeatureExtractionPipeline
}

// Load 載入模型（後端見 session_go.go、session_ort.go）
func Load(dir string) (*Model, error) {
	session, err := newSession()
	if err != nil {
		return nil, err
	}
	pipeline, err := hugot.NewPipeline(session, hugot.FeatureExtractionConfig{
		ModelPath:    dir,
		Name:         "e5",
		OnnxFilename: "model.onnx",
		Options:      []hugot.FeatureExtractionOption{pipelines.WithNormalization()},
	})
	if err != nil {
		session.Destroy()
		return nil, err
	}
	return &Model{session: session, pipeline: pipeline}, nil
}

func (m *Model) Close() error { return m.session.Destroy() }

// Embed 算句子向量（mean pooling 後 L2 正規化）。e5 的對稱任務（句子比句子）兩邊都加「query: 」前綴
func (m *Model) Embed(texts []string) ([][]float32, error) {
	inputs := make([]string, len(texts))
	for i, t := range texts {
		inputs[i] = "query: " + t
	}
	out, err := m.pipeline.RunPipeline(inputs)
	if err != nil {
		return nil, err
	}
	return out.Embeddings, nil
}

// Cosine 是兩個向量的 cosine 相似度
func Cosine(a, b []float32) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		na += float64(a[i]) * float64(a[i])
		nb += float64(b[i]) * float64(b[i])
	}
	return dot / (math.Sqrt(na)*math.Sqrt(nb) + 1e-9)
}
