package router

import (
	"fmt"
)

// RouteTarget 路由目標
type RouteTarget struct {
	TargetIP   string // 目標IP
	TargetPort int    // 目標端口
}

func (rt *RouteTarget) String() string {
	return fmt.Sprintf("%s:%d", rt.TargetIP, rt.TargetPort)
}

// Router interface
type Router interface {
	// 取得路由目標
	GetRouteTarget(routeData RouteData) *RouteTarget
}
