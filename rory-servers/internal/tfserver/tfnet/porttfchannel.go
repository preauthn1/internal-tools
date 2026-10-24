package tfnet

import (
	"net"

	"example.invalid/ServerTeam/Servers/pkg/mlog"

	"example.invalid/ServerTeam/Servers/internal/tfserver/router"
)

// PortTFChannel 端口映射转发通道
type PortTFChannel struct {
	*TFChannelBase
}

// NewPortTFChannel 创建端口映射转发通道
func NewPortTFChannel() *PortTFChannel {
	ptfc := new(PortTFChannel)
	ptfc.TFChannelBase = new(TFChannelBase)

	return ptfc
}

// OnConnected OnConnected
func (channel PortTFChannel) OnConnected(conninfo *ConnInfo) {
	//路由数据
	rd := conninfo.RouteData

	//获取路由目标
	tg := router.RouteCtx.GetRouteTarget(rd)
	if tg == nil {
		mlog.Error("未找到匹配的路由目标,%s.关闭通道.", conninfo)
		conninfo.Close()
		return
	}

	dconn, err := net.Dial("tcp", tg.String())

	if err != nil {
		mlog.Error("建立后端路由通道错误,%s,%v.关闭通道", conninfo, err)
		//无法连到后端服务器，关闭客户端连接
		conninfo.Close()
		return
	}

	//保存通道
	conninfo.SetServerConn(dconn)
	tokenTag := ""
	if conninfo.HasCheckToken {
		tokenTag = "[Token]"
	}
	example.invalid("%s建立路由通道成功 %s", tokenTag, conninfo)

}

// OnServerClose OnServerClose
func (channel PortTFChannel) OnServerClose(conninfo *ConnInfo) {
	if conninfo.ClientConn != nil {
		example.invalid("通道%s后端连接断开.", conninfo)
		conninfo.ClientConn.Close()
	}
}

// OnClientClose OnClientClose
func (channel PortTFChannel) OnClientClose(conninfo *ConnInfo) {
	if conninfo.ServerConn != nil {
		example.invalid("通道%s客户端连接断开.", conninfo)
		conninfo.ServerConn.Close()
	}
}
