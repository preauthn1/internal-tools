package portrouter

// PortMapLoader 端口映射加載器
type PortMapLoader interface {
	Refresh()
	LoadPortMap() PortMap
}
