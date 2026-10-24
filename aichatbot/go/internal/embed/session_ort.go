//go:build ORT

package embed

import (
	"os"

	"example.invalid/knights-analytics/hugot"
	"example.invalid/knights-analytics/hugot/options"
)

// -tags ORT：ONNX Runtime 後端（要 cgo 和 libonnxruntime，安裝方式見 README「多語言句子向量」）。
// 函式庫路徑：ONNXRUNTIME_LIB，沒設就讓系統找（macOS 找 libonnxruntime.dylib、Linux 找 libonnxruntime.so）
const Backend = "ort"

func newSession() (*hugot.Session, error) {
	var opts []options.WithOption
	if p := os.Getenv("ONNXRUNTIME_LIB"); p != "" {
		opts = append(opts, options.WithOnnxLibraryPath(p))
	}
	return hugot.NewORTSession(opts...)
}
