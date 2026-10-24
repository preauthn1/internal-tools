package router

import "fmt"

// RouteDataNames
const (
	RouteDataNameSourceIP   = "routedata_sourceip"
	RouteDataNameSourcePort = "routedata_sourceport"
)

// RouteData 路由数据
type RouteData map[string]interface{}

func (rdata RouteData) String() string {
	return fmt.Sprintf("%s:%d", rdata.GetSourceIP(), rdata.GetSourcePort())
}

// GetSourceIP 取得來源 IP
func (rdata RouteData) GetSourceIP() string {
	return rdata[RouteDataNameSourceIP].(string)
}

// GetSourcePort 取得來源Port號
func (rdata RouteData) GetSourcePort() int {
	return rdata[RouteDataNameSourcePort].(int)
}

// NewRouteData create RouteData
func NewRouteData(ip string, port int) RouteData {
	rd := make(map[string]interface{})
	rd[RouteDataNameSourceIP] = ip
	rd[RouteDataNameSourcePort] = port

	return rd
}
