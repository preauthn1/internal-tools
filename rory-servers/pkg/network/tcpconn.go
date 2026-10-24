package network

import (
	"fmt"
	"net"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// TCPConn tcp連接
type TCPConn struct {
	conn    net.Conn
	readBuf []byte
}

// NewTCPConn 建立TCP連接
func NewTCPConn(conn net.Conn, buflen int) *TCPConn {
	tcpconn := new(TCPConn)
	tcpconn.conn = conn
	tcpconn.readBuf = make([]byte, buflen)

	return tcpconn
}

// ReadMsg 讀取訊息
func (conn *TCPConn) ReadMsg() (returnData []byte, returnErr error) {
	defer func() {
		if err := recover(); err != nil {
			mlog.Error("conn read msg error: %v", err)
			returnData, returnErr = nil, fmt.Errorf("%v", err)
			return
		}
	}()

	n, err := conn.conn.Read(conn.readBuf)
	if err != nil {
		returnData, returnErr = nil, err
		return
	}
	if n > conf.READBUFLEN || n <= 0 {
		mlog.Error("conn read msg error:數據包長度超出範圍 %d", n)

		returnData, returnErr = nil, fmt.Errorf("數據包長度超出範圍 %d", n)

		return
	}

	returnData, returnErr = conn.readBuf[:n], nil
	return
}

// WriteMsg 寫入訊息
func (conn *TCPConn) WriteMsg(data []byte) (int, error) {
	if data == nil || len(data) == 0 {
		return 0, nil
	}

	return conn.conn.Write(data)
}

// LocalAddr 本地地址
func (conn *TCPConn) LocalAddr() net.Addr {
	return conn.conn.LocalAddr()
}

// RemoteAddr 遠端地址
func (conn *TCPConn) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}

// Close 關閉
func (conn *TCPConn) Close() {
	conn.conn.Close()
}

// Destory 毀掉
func (conn *TCPConn) Destory() {
	conn.Close()
	conn.readBuf = nil
}
