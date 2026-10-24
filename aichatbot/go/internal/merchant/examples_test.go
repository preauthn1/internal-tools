package merchant

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"aichatbot/internal/store"
)

// 比對規則和 eval/check_overlap.py 一樣：轉小寫，只留字母和數字（標點、空白、全形標點、底線都去掉）
func TestNormalize(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"Hello, World!", "helloworld"},
		{"請問 幾點開？", "請問幾點開"},
		{"「你好」～！", "你好"},
		{"a_b-c 1.5", "abc15"},
		{"ＡＢＣ１２３", "ａｂｃ１２３"},
	} {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// 題庫的句子比對後找得到是哪個題庫；只差標點、大小寫也算
func TestTestSetFind(t *testing.T) {
	tests, err := LoadTestSet(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ text, want string }{
		{"你們幾點開", "intent_set"},
		{"你們幾點開？！", "intent_set"},
		{"what are your opening hours on sunday", "faq_set"},
		// 第二家測試商家（美髮沙龍，#66）的題組
		{"染髮會不會傷髮質", "salon/intent_set"},
		{"Mia 這週六下午有空嗎", "salon/testset"},
		{"請問你是哪位", "salon/flow_tests"},
	} {
		if got := tests.Find(c.text); !slices.Contains(got, c.want) {
			t.Errorf("Find(%q) = %v, want 包含 %s", c.text, got, c.want)
		}
	}
	if got := tests.Find("你們的貓叫什麼名字"); len(got) != 0 {
		t.Errorf("不在題庫的句子：%v", got)
	}
}

// 和驗證題組重複時只說「驗證題組」，不說是哪一個檔（驗證題組不能被看出內容）
func TestTestSetHidesHoldout(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "holdout"), 0o755)
	os.MkdirAll(filepath.Join(dir, "stores", "salon"), 0o755)
	for name, body := range map[string]string{
		"intent_set.json":              `{"items": []}`,
		"testset.json":                 `{"cases": []}`,
		"faq_set.json":                 `{"cases": []}`,
		"flow_tests.json":              `{"conversations": []}`,
		"stores/salon/intent_set.json": `{"items": []}`,
		"stores/salon/testset.json":    `{"cases": []}`,
		"stores/salon/flow_tests.json": `{"conversations": []}`,
		"holdout/intent_holdout.json":  `{"items": [{"text": "祕密的句子"}]}`,
		"holdout/flow_holdout.json":    `{"conversations": [{"turns": [{"user": "祕密的句子！"}]}]}`,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	tests, err := LoadTestSet(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := tests.Find("祕密的句子"); !slices.Equal(got, []string{"驗證題組"}) {
		t.Errorf("Find = %v, want [驗證題組]", got)
	}
}

// 讀不到題庫就回錯誤（不能默默當作沒有重複）
func TestLoadTestSetMissing(t *testing.T) {
	if _, err := LoadTestSet(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("題庫資料夾不存在應該失敗")
	}
}

// 每個意圖都有繁中、英、日的白話名稱，代號不重複；內建案例用到的意圖都在表裡
func TestIntentLabels(t *testing.T) {
	seen := map[string]bool{}
	for _, in := range Intents {
		if seen[in.Code] {
			t.Errorf("意圖代號重複：%s", in.Code)
		}
		seen[in.Code] = true
		for _, lang := range []string{"zh-TW", "en-US", "ja-JP"} {
			if strings.TrimSpace(in.Labels[lang]) == "" {
				t.Errorf("%s 沒有 %s 的名稱", in.Code, lang)
			}
		}
	}
	if l := intentLabel(t, "booking", "zh-TW"); l != "想預約" {
		t.Errorf("booking 的繁中名稱：%q", l)
	}
	// 11 個內建意圖（CONTEXT.md）
	var codes []string
	for _, in := range Intents {
		codes = append(codes, in.Code)
	}
	if got := strings.Join(codes, ","); got != "booking,cancel,follow_up,claim,identity,self_intro,human_request,complaint,off_topic,unclear,general" {
		t.Errorf("意圖：%s", got)
	}
	if l := intentLabel(t, "human_request", "zh-TW"); l != "想找真人" {
		t.Errorf("human_request 的繁中名稱：%q", l)
	}
	examples, err := readExamples(filepath.Join(evalDir, "examples.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range examples {
		if _, ok := IntentByCode(e.Intent); !ok {
			t.Errorf("內建案例的意圖 %s 不在 Intents 裡", e.Intent)
		}
	}
}

// 匯入 eval/ 的內建案例：舊代號改成 general，分類當下的狀態（state）存進資料庫、讀得回來
func TestEvalExamplesKeepStateAndMergeCodes(t *testing.T) {
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	db := openDB(t)
	if _, err := Save(context.Background(), db.SQL(), c); err != nil {
		t.Fatal(err)
	}
	got, err := Load(context.Background(), db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	find := func(text string) Example {
		for _, e := range got.Examples {
			if e.Text == text {
				return e
			}
		}
		t.Fatalf("找不到案例 %s", text)
		return Example{}
	}
	if e := find("怎麼還沒回"); e.Intent != "follow_up" || e.State != "ticket" || !e.Builtin {
		t.Errorf("怎麼還沒回：%+v", e)
	}
	if e := find("五位"); e.Intent != "booking" || e.State != "booking:people" {
		t.Errorf("五位：%+v", e)
	}
	if e := find("你們開到幾點"); e.Intent != "general" || e.State != "" {
		t.Errorf("你們開到幾點：%+v", e)
	}
}

func intentLabel(t *testing.T, code, lang string) string {
	t.Helper()
	in, ok := IntentByCode(code)
	if !ok {
		t.Fatalf("找不到意圖 %s", code)
	}
	return in.Labels[lang]
}

// reimportFixture：匯入森日茶屋，後台新增一句商家案例
func reimportFixture(t *testing.T, text string) (*store.DB, Config) {
	t.Helper()
	ctx := context.Background()
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	db := openDB(t)
	if _, err := Save(ctx, db.SQL(), c); err != nil {
		t.Fatal(err)
	}
	tests, err := LoadTestSet(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := AddExample(ctx, db.SQL(), tests, text, "general"); err != nil {
		t.Fatal(err)
	}
	return db, c
}

// splitExamples 把案例庫分成內建和商家新增的
func splitExamples(t *testing.T, db *store.DB) (builtin, added []Example) {
	t.Helper()
	all, err := ListExamples(context.Background(), db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range all {
		x := Example{Text: e.Text, Intent: e.Intent, State: e.State, Builtin: e.Builtin}
		if e.Builtin {
			builtin = append(builtin, x)
		} else {
			added = append(added, x)
		}
	}
	return builtin, added
}

// 重新匯入：商家新增的案例保留，內建案例和檔案一致
func TestReimportKeepsMerchantExamples(t *testing.T) {
	db, c := reimportFixture(t, "你們的貓叫什麼名字")
	c.Examples = append(c.Examples[1:], Example{Text: "店裡可以帶狗嗎", Intent: "general", Builtin: true})
	removed, err := Save(context.Background(), db.SQL(), c)
	if err != nil {
		t.Fatal(err)
	}
	if len(removed) != 0 {
		t.Errorf("沒有重複，不該刪掉商家的案例：%v", removed)
	}
	builtin, added := splitExamples(t, db)
	if !slices.Equal(builtin, c.Examples) {
		t.Errorf("內建案例和檔案不同：%d 句 vs %d 句", len(builtin), len(c.Examples))
	}
	if want := []Example{{Text: "你們的貓叫什麼名字", Intent: "general"}}; !slices.Equal(added, want) {
		t.Errorf("商家案例：%v, want %v", added, want)
	}
}

// 商家的案例和新的內建案例重複（比對後相同）：只留內建的，回傳被刪掉的商家案例
func TestReimportDropsMerchantDuplicateOfBuiltin(t *testing.T) {
	db, c := reimportFixture(t, "你們的貓叫什麼名字")
	c.Examples = append(c.Examples, Example{Text: "你們的貓叫什麼名字！", Intent: "off_topic", Builtin: true})
	removed, err := Save(context.Background(), db.SQL(), c)
	if err != nil {
		t.Fatal(err)
	}
	if want := []Example{{Text: "你們的貓叫什麼名字", Intent: "general"}}; !slices.Equal(removed, want) {
		t.Errorf("被刪掉的：%v, want %v", removed, want)
	}
	builtin, added := splitExamples(t, db)
	if !slices.Equal(builtin, c.Examples) {
		t.Errorf("內建案例和檔案不同：%d 句 vs %d 句", len(builtin), len(c.Examples))
	}
	if len(added) != 0 {
		t.Errorf("重複的商家案例應該刪掉：%v", added)
	}
}

// 匯出的檔案含商家案例：匯入同一個資料庫不會變兩句，匯入新的資料庫會還原
func TestImportExportedMerchantExamples(t *testing.T) {
	ctx := context.Background()
	db, _ := reimportFixture(t, "你們的貓叫什麼名字")
	exported, err := Load(ctx, db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	want := []Example{{Text: "你們的貓叫什麼名字", Intent: "general"}}
	for _, target := range []*store.DB{db, openDB(t)} {
		removed, err := Save(ctx, target.SQL(), exported)
		if err != nil {
			t.Fatal(err)
		}
		if len(removed) != 0 {
			t.Errorf("不該刪掉商家的案例：%v", removed)
		}
		if _, added := splitExamples(t, target); !slices.Equal(added, want) {
			t.Errorf("商家案例：%v, want %v", added, want)
		}
	}
}

// 現有的 examples.jsonl 和題庫 0 句重複（和 check_overlap.py 的結果一致）
func TestEvalExamplesDoNotOverlap(t *testing.T) {
	tests, err := LoadTestSet(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	examples, err := readExamples(filepath.Join(evalDir, "examples.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range examples {
		if sets := tests.Find(e.Text); len(sets) > 0 {
			t.Errorf("%q 和題庫重複：%v", e.Text, sets)
		}
	}
}
