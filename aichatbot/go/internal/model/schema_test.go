package model

import "testing"

// 欄位定義轉成 JSON schema：照欄位順序輸出（模型依序產生，順序會影響結果），選擇題是 enum，全部必填、不能多欄位
func TestSchema(t *testing.T) {
	r := Request{Name: "Booking", Fields: []Field{
		{Kind: KindChoice, Name: "intent", Description: "意圖", Options: []string{"booking", "other"}},
		{Kind: KindText, Name: "name", Description: "名字"},
		{Kind: KindInteger, Name: "people", Description: "人數"},
	}}
	want := `{"type":"object","properties":{` +
		`"intent":{"type":"string","description":"意圖","enum":["booking","other"]},` +
		`"name":{"type":"string","description":"名字"},` +
		`"people":{"type":"integer","description":"人數"}},` +
		`"required":["intent","name","people"],"additionalProperties":false}`
	if got := string(Schema(r)); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

// 沒有欄位時 required 是空陣列，不是 null
func TestSchemaNoFields(t *testing.T) {
	want := `{"type":"object","properties":{},"required":[],"additionalProperties":false}`
	if got := string(Schema(Request{Name: "Empty"})); got != want {
		t.Errorf("got %s", got)
	}
}
