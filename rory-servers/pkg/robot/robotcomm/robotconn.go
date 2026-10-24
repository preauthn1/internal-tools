package robotcomm

import (
	"net"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/network"
)

// ?
const (
	ServerIDRobotConn = "robotconn"
)

// RobotConn 機器人客戶端連接
type RobotConn struct{}

// ReadMsg ?
func (rc *RobotConn) ReadMsg() ([]byte, error) {
	return nil, nil
}

// WriteMsg ?
func (rc *RobotConn) WriteMsg(data []byte) (int, error) {
	return 0, nil
}

// LocalAddr ?
func (rc *RobotConn) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr ?
func (rc *RobotConn) RemoteAddr() net.Addr {
	return nil
}

// Close ?
func (rc *RobotConn) Close() {}

// Destory ?
func (rc *RobotConn) Destory() {}

// NewRebotConn ?
func NewRebotConn() *RobotConn {
	rconn := new(RobotConn)

	return rconn
}

// NewRobotConnInfo 創建機器人客戶端信息
func NewRobotConnInfo() *network.ConnInfo {
	cinfo := network.NewConnInfo(NewRebotConn(), ServerIDRobotConn, network.TCP)
	return cinfo
}

// IsRobotConn 是否機器人客戶端
func IsRobotConn(connid string) bool {
	return strings.Index(connid, "unknown:-1_") != -1
}
