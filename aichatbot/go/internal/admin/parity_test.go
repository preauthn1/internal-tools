package admin

import (
	"bytes"
	"context"
	"testing"

	"aichatbot/internal/merchant"
	"aichatbot/internal/store"
)

// exported 讀出資料庫裡的商家設定，輸出成匯出檔的格式
func exported(t *testing.T, db *store.DB) []byte {
	t.Helper()
	c, err := merchant.Load(context.Background(), db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	data, err := merchant.Export(c)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// 最終驗收（#39）：從空的資料庫只透過後台 API 做出森日茶屋的設定，匯出後和直接匯入 eval/ 的完全相同
func TestBuildViaAPIMatchesEvalImport(t *testing.T) {
	ctx := context.Background()
	cfg, err := merchant.FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	imported := openDB(t)
	if _, err := merchant.Save(ctx, imported.SQL(), cfg); err != nil {
		t.Fatal(err)
	}
	viaAPI := openDB(t)
	if err := BuildViaAPI(ctx, viaAPI, evalDir, cfg); err != nil {
		t.Fatal(err)
	}
	if a, b := exported(t, imported), exported(t, viaAPI); !bytes.Equal(a, b) {
		t.Errorf("後台做出的設定和匯入的不同：\n匯入：\n%s\n後台：\n%s", a, b)
	}
}

// 第二家測試商家（美髮沙龍，#66）的匯出檔：後台讀得進來，只透過後台 API 也能做出同樣的設定
func TestBuildViaAPIMatchesSalonConfig(t *testing.T) {
	ctx := context.Background()
	cfg, err := merchant.ReadFile(evalDir + "/stores/salon/config.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.FAQ) < 10 || len(cfg.Knowledge) < 6 || !merchant.FeatureEnabled(cfg, "booking") || !merchant.FeatureEnabled(cfg, "handoff") {
		t.Errorf("FAQ %d 題、知識 %d 段、功能 %+v", len(cfg.FAQ), len(cfg.Knowledge), cfg.Features)
	}
	// #67：沙龍不問人數，服務項目必填、指定設計師選填（後台的預約設定 API 也要做得出來）
	if b := cfg.BookingFields; b == nil || b.PeopleEnabled || len(b.Custom) != 2 ||
		b.Custom[0].Key != "service" || !b.Custom[0].Required || b.Custom[1].Key != "designer" || b.Custom[1].Required {
		t.Errorf("預約欄位：%+v", cfg.BookingFields)
	}
	imported := openDB(t)
	if _, err := merchant.Save(ctx, imported.SQL(), cfg); err != nil {
		t.Fatal(err)
	}
	viaAPI := openDB(t)
	if err := BuildViaAPI(ctx, viaAPI, evalDir, cfg); err != nil {
		t.Fatal(err)
	}
	if a, b := exported(t, imported), exported(t, viaAPI); !bytes.Equal(a, b) {
		t.Errorf("後台做出的設定和匯入的不同：\n匯入：\n%s\n後台：\n%s", a, b)
	}
}

// 森日茶屋的檔案很多區塊是預設值（功能都開、中英日、沒有自訂意圖和商家新增的案例），
// 改成非預設的設定再比一次，每一支 API 都要真的存進去
func TestBuildViaAPIMatchesNonDefaultConfig(t *testing.T) {
	ctx := context.Background()
	cfg, err := merchant.FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Merchant.PreferredLanguage = "ja"
	cfg.Languages = []merchant.Language{{Code: "ja", Enabled: true}, {Code: "zh", Enabled: true}, {Code: "ko", Enabled: true}, {Code: "en", Enabled: true}, {Code: "th", Enabled: false}}
	cfg.Features = []merchant.Feature{{Name: "booking", Enabled: false}, {Name: "handoff", Enabled: false}}
	cfg.Replies = append(cfg.Replies,
		merchant.Reply{Key: "booking_off", Lang: "en", Text: "Please call us to book: {contact}"},
		merchant.Reply{Key: "booking_off", Lang: "ja", Text: "予約は電話でお願いします：{contact}"},
		merchant.Reply{Key: "booking_off", Lang: "ko", Text: "예약은 전화로 부탁드립니다: {contact}"},
		merchant.Reply{Key: "booking_off", Lang: "zh", Text: "預約請打電話：{contact}|請來電預約：{contact}"},
		merchant.Reply{Key: "greeting", Lang: "ko", Text: "안녕하세요~"},
		merchant.Reply{Key: "handoff", Lang: "ko", Text: "직원에게 연결해 드릴게요~"},
	)
	cfg.CustomIntents = []merchant.CustomIntent{{
		Code: "custom_1", Name: "包場", Description: "想包下整間店辦活動",
		Examples: []string{"整間店可以包下來嗎", "想包場辦生日派對", "公司聚餐能不能包場"},
		Action:   merchant.ActionFixedText, Text: "包場請來電洽詢喔～", Knowledge: []string{},
	}}
	cfg.Examples = append(cfg.Examples,
		merchant.Example{Text: "週末想整間包下來", Intent: "custom_1"},
		merchant.Example{Text: "我要訂下禮拜的位子喔", Intent: "booking"},
	)
	imported := openDB(t)
	if _, err := merchant.Save(ctx, imported.SQL(), cfg); err != nil {
		t.Fatal(err)
	}
	viaAPI := openDB(t)
	if err := BuildViaAPI(ctx, viaAPI, evalDir, cfg); err != nil {
		t.Fatal(err)
	}
	if a, b := exported(t, imported), exported(t, viaAPI); !bytes.Equal(a, b) {
		t.Errorf("後台做出的設定和匯入的不同：\n匯入：\n%s\n後台：\n%s", a, b)
	}
}
