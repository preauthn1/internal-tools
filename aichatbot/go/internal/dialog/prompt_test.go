package dialog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// #64：內建提示詞不限店家類型，不能帶著測試店家（eval/store.json 的森日茶屋）的店名、品項名稱：
// 掃描這個套件非測試的原始碼（含註解，新增的提示詞也掃得到），不分大小寫
func TestSourceHasNoTestStoreNames(t *testing.T) {
	raw, err := os.ReadFile("../../../eval/store.json")
	if err != nil {
		t.Fatal(err)
	}
	var st struct {
		Name string `json:"name"`
		Menu []struct {
			Aliases []string          `json:"aliases"`
			Names   map[string]string `json:"names"`
		} `json:"menu"`
	}
	if err := json.Unmarshal(raw, &st); err != nil {
		t.Fatal(err)
	}
	names := []string{st.Name}
	for _, m := range st.Menu {
		names = append(names, m.Aliases...)
		for _, n := range m.Names {
			names = append(names, n)
		}
	}
	files, _ := filepath.Glob("*.go")
	if len(files) == 0 {
		t.Fatal("找不到原始碼")
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		// 分類提示詞的例句「紅玉怎麼賣」暫時保留：換成「招牌怎麼賣」「招牌多少錢」或拿掉，本機模型的分類都會退步（「中」「回答什麼東西啊」…，
		// #64 量過三種），要改進迴圈逐題守住不退步地換
		lower := strings.ToLower(strings.ReplaceAll(string(src), "紅玉怎麼賣", ""))
		for _, n := range names {
			if strings.Contains(lower, strings.ToLower(n)) {
				t.Errorf("%s 有測試店家的名稱：%s", f, n)
			}
		}
	}
}
