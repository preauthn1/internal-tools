package dialog

import "time"

// Customer 是顧客的對話狀態，依「平台＋使用者 ID」整份存成 JSON（store 的 customer_state）。
// 少了的欄位用零值（之後加欄位，舊資料照樣讀得出來）
type Customer struct {
	Name           string            `json:"name,omitempty"`
	Lang           string            `json:"lang,omitempty"`
	Booking        map[string]string `json:"booking,omitempty"`
	AskingField    string            `json:"asking_field,omitempty"`     // 預約中：正在問哪個欄位（空的 = 沒有在預約）
	AskCount       int               `json:"ask_count,omitempty"`        // 預約中：正在問的欄位已經問了幾次（#53：問太多次轉人工）
	OpenTicket     bool              `json:"open_ticket,omitempty"`      // 有事情已轉給員工、還沒有回覆
	FollowUps      int               `json:"follow_ups,omitempty"`       // 這件事被追問了幾次
	OffTopicStreak int               `json:"off_topic_streak,omitempty"` // 連續無關或聲稱的次數
	AskedName      bool              `json:"asked_name,omitempty"`       // 已經問過稱呼（ask_name_soft），之後不再追問（#61）
	AskingName     bool              `json:"asking_name,omitempty"`      // 上一句剛問了稱呼：這一句是自我介紹句型就當成報名字（#61）
	Complained     bool              `json:"complained,omitempty"`       // 上一句回了抱怨的道歉（complaint）：再抱怨就轉人工（#62）
	Human          bool              `json:"human,omitempty"`            // 已轉人工，AI 不再回覆
	HandoffTicket  int64             `json:"handoff_ticket,omitempty"`   // 這次轉人工的待處理事項（舊通知的「結束轉人工」不結束新的這次）
	HumanSince     time.Time         `json:"human_since,omitzero"`       // 轉人工的時間（自動結束用）
	LastStaffReply time.Time         `json:"last_staff_reply,omitzero"`  // 員工最後一次回覆的時間（自動結束用）
}

// endBooking 結束預約流程（送出、取消或轉人工）：清掉收到的欄位和正在問的欄位
func (c *Customer) endBooking() { c.Booking, c.AskingField, c.AskCount = nil, "", 0 }
