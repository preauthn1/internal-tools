package logic

// MoneyTransfer 金錢轉換
type MoneyTransfer interface {
	S2C(float64) float64 // Server to Client
	C2S(float64) float64 // Client to Server
}

// ScoreTransfer 以分數為主
type ScoreTransfer struct{}

// S2C Server 轉換 Client
func (*ScoreTransfer) S2C(f float64) float64 {
	return f
}

// C2S Server 轉換 Client
func (*ScoreTransfer) C2S(f float64) float64 {
	return f
}

// DollarTransfer 以錢幣為主
type DollarTransfer struct{}

// S2C Server 轉換 Client
func (*DollarTransfer) S2C(f float64) float64 {
	return f / 100
}

// C2S Server 轉換 Client
func (*DollarTransfer) C2S(f float64) float64 {
	return f * 100
}
