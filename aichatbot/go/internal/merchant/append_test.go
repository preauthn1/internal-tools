package merchant

import (
	"context"
	"errors"
	"testing"
)

// AddFAQ 只在最後新增一題（代號接續沒用過的 faq-N），原本的 FAQ 不動；問題、答案必填
func TestAddFAQ(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	c := valid()
	c.FAQ = []FAQ{{Key: "faq-1", Question: "可以刷卡嗎", AltQuestions: []string{}, Answer: "可以"}}
	if _, err := Save(ctx, db, c); err != nil {
		t.Fatal(err)
	}
	key, err := AddFAQ(ctx, db, FAQ{Question: " 有停車場嗎 ", AltQuestions: []string{"哪裡停車", " "}, Answer: "對面有付費停車場"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := Load(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if key != "faq-2" || len(got.FAQ) != 2 || got.FAQ[0].Key != "faq-1" {
		t.Fatalf("key %s，FAQ %+v", key, got.FAQ)
	}
	if f := got.FAQ[1]; f.Key != "faq-2" || f.Question != "有停車場嗎" || len(f.AltQuestions) != 1 || f.AltQuestions[0] != "哪裡停車" || f.Answer != "對面有付費停車場" {
		t.Errorf("新增的：%+v", f)
	}
	var invalid *InvalidError
	if _, err := AddFAQ(ctx, db, FAQ{Question: "有包廂嗎"}); !errors.As(err, &invalid) {
		t.Errorf("答案空白：%v", err)
	}
}

// AppendKnowledge：標題是現有的段落就接在那段後面（換行），否則在最後新增一段；其他段落不動
func TestAppendKnowledge(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	c := valid()
	c.Knowledge = []Section{{Title: "營業時間", Body: "11 點開"}, {Title: "交通", Body: "捷運站走路 5 分鐘"}}
	if _, err := Save(ctx, db, c); err != nil {
		t.Fatal(err)
	}
	if err := AppendKnowledge(ctx, db, " 交通 ", " 對面有付費停車場 "); err != nil {
		t.Fatal(err)
	}
	if err := AppendKnowledge(ctx, db, "寵物", "可以帶狗"); err != nil {
		t.Fatal(err)
	}
	got, err := Load(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	want := []Section{{Title: "營業時間", Body: "11 點開"}, {Title: "交通", Body: "捷運站走路 5 分鐘\n對面有付費停車場"}, {Title: "寵物", Body: "可以帶狗"}}
	if len(got.Knowledge) != 3 {
		t.Fatalf("知識段落：%+v", got.Knowledge)
	}
	for i, s := range want {
		if got.Knowledge[i].Title != s.Title || got.Knowledge[i].Body != s.Body {
			t.Errorf("第 %d 段：%+v", i+1, got.Knowledge[i])
		}
	}
	var invalid *InvalidError
	if err := AppendKnowledge(ctx, db, "", "內容"); !errors.As(err, &invalid) {
		t.Errorf("標題空白：%v", err)
	}
}
