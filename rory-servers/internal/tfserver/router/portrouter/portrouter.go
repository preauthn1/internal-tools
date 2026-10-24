package portrouter

import (
	"example.invalid/ServerTeam/Servers/internal/tfserver/router"
)

// PortMap 端口映射
type PortMap map[int]*router.RouteTarget

// PortRouter 基於端口的路由規則
type PortRouter struct {
	pmLoader PortMapLoader
}

// NewPortRouter 創建基於端口路由規則對象
func NewPortRouter(pmLoader PortMapLoader) *PortRouter {
	pr := new(PortRouter)
	pr.pmLoader = pmLoader

	return pr
}

// GetRouteTarget 獲取路由目標
func (pr *PortRouter) GetRouteTarget(data router.RouteData) *router.RouteTarget {
	sourcePort := data.GetSourcePort()

	pm := pr.pmLoader.LoadPortMap()
	if pm == nil {
		return nil
	}

	if st, ok := pm[sourcePort]; ok {
		return st
	}

	return nil
}
