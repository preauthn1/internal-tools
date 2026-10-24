package telegram

import "encoding/json"

// Update 是 Telegram 送來的 update，只解碼用得到的欄位。https://example.invalid/bots/api#update
type Update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"` // 員工按了通知上的按鈕
}

type Message struct {
	MessageID      int64    `json:"message_id"`
	Chat           Chat     `json:"chat"`
	From           *User    `json:"from"`             // 頻道貼文沒有 from
	Text           string   `json:"text"`             // 貼圖、照片沒有 text
	Caption        string   `json:"caption"`          // 照片、檔案的說明文字
	ReplyToMessage *Message `json:"reply_to_message"` // 員工「回覆」的那則
	Sticker        *Sticker `json:"sticker"`
	// 看不懂的非文字訊息：只看有沒有
	Photo     json.RawMessage `json:"photo"`
	Voice     json.RawMessage `json:"voice"`
	Audio     json.RawMessage `json:"audio"`
	Video     json.RawMessage `json:"video"`
	VideoNote json.RawMessage `json:"video_note"`
	Document  json.RawMessage `json:"document"`
	// 群組升級成超級群組時，舊群組會收到這則，之後要用新的 chat id
	MigrateToChatID int64 `json:"migrate_to_chat_id"`
}

type Sticker struct {
	Emoji string `json:"emoji"`
}

// mediaKind 是看不懂的非文字訊息的種類（給員工看）；文字、貼圖和其他訊息是空字串
func (m *Message) mediaKind() string {
	switch {
	case m.Photo != nil:
		return "圖片"
	case m.Voice != nil:
		return "語音"
	case m.Audio != nil:
		return "音訊"
	case m.Video != nil, m.VideoNote != nil:
		return "影片"
	case m.Document != nil:
		return "檔案"
	}
	return ""
}

type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"` // private、group、supergroup、channel
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// CallbackQuery 是按鈕回呼；Message 是按鈕所在的那則訊息（太舊的訊息可能沒有）
type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Message *Message `json:"message"`
	Data    string   `json:"data"`
}
