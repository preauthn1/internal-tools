package msgserver

// ServerLogicHandler 服務器處邏輯理器
type ServerLogicHandler interface {
	Init(server *MsgServer) // 初始化服務器
	Start()                 // 啟動服務器
	Stop()                  // 停止服務器
	Name() string           // 名稱
}

// ServerLogicGroup 服務器邏輯處理器組
type ServerLogicGroup struct {
	innerLogicHandlers []ServerLogicHandler
}

// NewServerLogicGroup 創建服務器組
func NewServerLogicGroup() *ServerLogicGroup {
	return new(ServerLogicGroup)
}

// RegisterHandler 註冊邏輯處理器
func (slg *ServerLogicGroup) RegisterHandler(h ServerLogicHandler) *ServerLogicGroup {
	slg.innerLogicHandlers = append(slg.innerLogicHandlers, h)

	return slg
}

// Init 初始化
func (slg *ServerLogicGroup) Init(server *MsgServer) {
	for _, h := range slg.innerLogicHandlers {
		h.Init(server)
	}
}

// Start 啟動
func (slg *ServerLogicGroup) Start() {
	for _, h := range slg.innerLogicHandlers {
		h.Start()
	}
}

// Stop 停止
func (slg *ServerLogicGroup) Stop() {
	for _, h := range slg.innerLogicHandlers {
		h.Stop()
	}
}

// Name 名稱
func (slg *ServerLogicGroup) Name() string {
	return "logic group"
}
