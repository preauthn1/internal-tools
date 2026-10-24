package serverhandler

// GameServerHandler 服務器處理器
type GameServerHandler interface {
	Name() string             // 處理器名稱
	Start(ctx *GameServerCtx) // 啟動處理器
	Stop()                    // 停止
}

// GameServerHandlerGroup gameserver處理器組
type GameServerHandlerGroup struct {
	innerHandlers []GameServerHandler
}

// NewGameServerHandlerGroup 創建服務器處理器組對象
func NewGameServerHandlerGroup() *GameServerHandlerGroup {
	gsg := new(GameServerHandlerGroup)

	return gsg
}

// AddHandler 添加處理器
func (gsg *GameServerHandlerGroup) AddHandler(h GameServerHandler) *GameServerHandlerGroup {
	gsg.innerHandlers = append(gsg.innerHandlers, h)

	return gsg
}

// Name 處理器名稱
func (gsg *GameServerHandlerGroup) Name() string {
	return "handler group"
}

// Start 啟動處理器
func (gsg *GameServerHandlerGroup) Start(ctx *GameServerCtx) {
	for _, h := range gsg.innerHandlers {
		go h.Start(ctx)
	}
}

// Stop 停止
func (gsg *GameServerHandlerGroup) Stop() {
	for _, h := range gsg.innerHandlers {
		h.Stop()
	}
}

// Range range all handlers
func (gsg *GameServerHandlerGroup) Range(f func(h GameServerHandler)) {
	for _, h := range gsg.innerHandlers {
		f(h)
	}
}
