package monitorutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	//"mlog"
)

// NotifyData notify data
type NotifyData struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	ThemeColor string `json:"themeColor"`
}

// Push push data
func Push(NotifyWebHook, errStr string, benchmarkTime time.Duration) {
	data := &NotifyData{
		Title:      "連線異常",
		Text:       fmt.Sprintf("執行時間 : %s<br />錯誤訊息%s", benchmarkTime.String(), errStr),
		ThemeColor: "red"}
	post, _ := json.Marshal(data)

	fmt.Println(string(post))
	resp, err := http.Post(NotifyWebHook, "application/json", bytes.NewBuffer(post))
	_ = resp
	if err != nil {
		//mlog.Error("push error : %v", err.Error())
	}

	//fmt.Println("resp ", resp)
}
