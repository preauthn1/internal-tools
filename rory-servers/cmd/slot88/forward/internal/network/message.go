package network

// Message slot88 傳遞封包格式
type Message struct {
	Data      []byte      `json:"-"`
	ID        int         `json:"id,string"`
	ActionID  int         `json:"actionId,string"`
	EventData interface{} `json:"eventData,omitempty"`
}

// Netty slot88 錯誤封包
type Netty struct {
	Code    int
	Encrypt bool
	Data    interface{}
}
