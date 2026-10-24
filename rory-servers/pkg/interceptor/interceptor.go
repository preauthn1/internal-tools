package interceptor

import "example.invalid/ServerTeam/Servers/pkg/msg"

var msgInterceptor []MessageInterceptor

// MessageInterceptor 封包攔截器
type MessageInterceptor interface {
	Name() string
	BeforeSend(m *msg.Message)
	AfterSend(m *msg.Message)
	BeforeReceive(m *msg.Message)
	AfterReceive(m *msg.Message)
}

// RegisterMSGInterceptor 註冊Message封包攔截器
func RegisterMSGInterceptor(i MessageInterceptor) {
	msgInterceptor = append(msgInterceptor, i)
}

// GetMSGInterceptor 取得註冊過的Message封包攔截器
func GetMSGInterceptor() []MessageInterceptor {
	return msgInterceptor
}
