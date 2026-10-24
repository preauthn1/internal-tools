package router

// RouteTable 路由表
type RouteTable struct {
	routers []Router
}

// AddRouter add Router to RouteTable
func (rt *RouteTable) AddRouter(r Router) {
	rt.routers = append(rt.routers, r)
}

// GetRouteTarget 取得路由目標
func (rt *RouteTable) GetRouteTarget(data RouteData) *RouteTarget {
	for _, r := range rt.routers {
		ret := r.GetRouteTarget(data)
		if ret != nil {
			return ret
		}
	}

	return nil
}

// RouteCtx Route Context
var RouteCtx *RouteTable

func init() {
	RouteCtx = new(RouteTable)
}
