package tfnet

// TFChannel 转发通道接口定义
type TFChannel interface {
	// 收到客户端消息
	OnMessage(conninfo *ConnInfo, data []byte)

	// 连接建立
	OnConnected(conninfo *ConnInfo)

	// 服务端连接关闭
	OnServerClose(conninfo *ConnInfo)

	// 客户端连接关闭
	OnClientClose(conninfo *ConnInfo)
}

// TFChannelBase 转发通道基类
type TFChannelBase struct {
}

// OnMessage OnMessage
func (tfcb *TFChannelBase) OnMessage(conninfo *ConnInfo, data []byte) {

}

// OnConnected OnConnected
func (tfcb *TFChannelBase) OnConnected(conninfo *ConnInfo) {

}

// OnServerClose OnServerClose
func (tfcb *TFChannelBase) OnServerClose(conninfo *ConnInfo) {

}

// OnClientClose OnClientClose
func (tfcb *TFChannelBase) OnClientClose(conninfo *ConnInfo) {

}

// TFChannelGroup 转发通道组
type TFChannelGroup struct {
	channels []TFChannel
}

// NewTFChannelGroup create TFChannelGroup
func NewTFChannelGroup() *TFChannelGroup {
	cg := new(TFChannelGroup)
	cg.channels = make([]TFChannel, 0)

	return cg
}

// OnMessage OnMessage
func (cg *TFChannelGroup) OnMessage(conninfo *ConnInfo, data []byte) {
	for _, c := range cg.channels {
		c.OnMessage(conninfo, data)
	}
}

// OnConnected OnConnected
func (cg *TFChannelGroup) OnConnected(conninfo *ConnInfo) {
	for _, c := range cg.channels {
		c.OnConnected(conninfo)
	}
}

// OnServerClose OnServerClose
func (cg *TFChannelGroup) OnServerClose(conninfo *ConnInfo) {
	for _, c := range cg.channels {
		c.OnServerClose(conninfo)
	}
}

// OnClientClose OnClientClose
func (cg *TFChannelGroup) OnClientClose(conninfo *ConnInfo) {
	for _, c := range cg.channels {
		c.OnClientClose(conninfo)
	}
}

// AddChannel AddChannel
func (cg *TFChannelGroup) AddChannel(ch TFChannel) {
	cg.channels = append(cg.channels, ch)
}

// ChannelCtx Channel Context
var ChannelCtx *TFChannelGroup

func init() {
	ChannelCtx = NewTFChannelGroup()
}
