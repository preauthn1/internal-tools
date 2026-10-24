package network

// ClientEventChan 客戶端消息通道
type ClientEventChan struct {
	DataSig         chan []byte    // 數據消息通道
	MsgSig          chan []byte    // 完整數據包消息通道
	ConnectedSig    chan *ConnInfo // 連接成功消息通道
	DisconnectedSig chan *ConnInfo // 連接斷開消息通道
}

// NewClientEventChan 新建客戶端消息通道
func NewClientEventChan(sigSize int) *ClientEventChan {
	evt := new(ClientEventChan)

	evt.DataSig = make(chan []byte, sigSize)
	evt.MsgSig = make(chan []byte, sigSize)
	evt.ConnectedSig = make(chan *ConnInfo, sigSize)
	evt.DisconnectedSig = make(chan *ConnInfo, sigSize)

	return evt
}

// ServerNetEventChan 服務器網絡事件通道
type ServerNetEventChan struct {
	ConnectSig    chan *ConnInfo // 客戶端連接
	DisConnectSig chan *ConnInfo // 客戶端斷開
	DataSig       chan *ConnData // 收到數據
	MsgSig        chan *ConnData // 收到完整數據包
}

// NewNetEventChan 服務器網絡事件通道
func NewNetEventChan(chansize int) *ServerNetEventChan {
	evtchan := new(ServerNetEventChan)
	evtchan.ConnectSig = make(chan *ConnInfo, chansize)
	evtchan.DisConnectSig = make(chan *ConnInfo, chansize)
	evtchan.DataSig = make(chan *ConnData, chansize)
	evtchan.MsgSig = make(chan *ConnData, chansize)

	return evtchan
}
