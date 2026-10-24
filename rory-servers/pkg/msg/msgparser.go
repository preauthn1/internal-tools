package msg

// Parser 消息解析接口,消息對象和二進制數據之間的轉換
type Parser interface {
	Marshal(msg *Message) ([]byte, error)    // 序列化消息
	Unmarshal(data []byte) (*Message, error) // 反序列化消息
}

// ParserInfo 消息解析器信息
type ParserInfo struct {
	MsgParser Parser
	BClassID  int32
	SClassID  int32
}

// NewParserInfo 創建消息解析器信息對象
func NewParserInfo(bclsid int32, sclsid int32, parser Parser) *ParserInfo {
	mpi := new(ParserInfo)
	mpi.BClassID = bclsid
	mpi.SClassID = sclsid
	mpi.MsgParser = parser

	return mpi
}
