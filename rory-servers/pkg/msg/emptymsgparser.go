package msg

// EmptyMsgParser ?
type EmptyMsgParser struct{}

// Marshal ?
func (parser *EmptyMsgParser) Marshal(msg *Message) ([]byte, error) {
	return nil, nil
}

// Unmarshal ?
func (parser *EmptyMsgParser) Unmarshal(data []byte) (*Message, error) {
	return nil, nil
}
