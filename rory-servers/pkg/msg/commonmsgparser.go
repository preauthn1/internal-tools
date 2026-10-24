package msg

// CommonParser Public message parser
type CommonParser struct {
	mf  func(msg *Message) ([]byte, error)
	umf func(data []byte) (*Message, error)
}

// Marshal ?
func (cmp *CommonParser) Marshal(msg *Message) ([]byte, error) {
	if cmp.mf != nil {
		return cmp.mf(msg)
	}
	return nil, nil
}

// Unmarshal ?
func (cmp *CommonParser) Unmarshal(data []byte) (*Message, error) {
	if cmp.umf != nil {
		return cmp.umf(data)
	}
	return nil, nil
}

// NewInputMsgParser 創建單向輸入解析的消息解析器
func NewInputMsgParser(umf func(data []byte) (*Message, error)) *CommonParser {
	cmp := new(CommonParser)
	cmp.mf = nil
	cmp.umf = umf

	return cmp
}

// NewOutputMsgParser 創建單向輸出解析的消息解析器
func NewOutputMsgParser(mf func(msg *Message) ([]byte, error)) *CommonParser {
	cmp := new(CommonParser)
	cmp.mf = mf
	cmp.umf = nil

	return cmp
}

// NewCommonMsgParser 創建消息解析器
func NewCommonMsgParser(mf func(msg *Message) ([]byte, error), umf func(data []byte) (*Message, error)) *CommonParser {
	cmp := new(CommonParser)
	cmp.mf = mf
	cmp.umf = umf

	return cmp
}
