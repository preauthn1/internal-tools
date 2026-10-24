package db

import (
	"fmt"
)

// SQLConnInfo 數據庫連接
type SQLConnInfo struct {
	ServerIP string // 服務器地址
	Port     int    // 端口
	UserID   string // 用戶名
	Password string // 密碼
	DbName   string
}

// NewSQLConnInfo 創建sql數據庫連接
func NewSQLConnInfo(serverip string, port int, userid string, password string, dbName string) *SQLConnInfo {
	cn := new(SQLConnInfo)
	cn.ServerIP = serverip
	cn.Port = port
	cn.UserID = userid
	cn.Password = [REDACTED]

	return cn
}

// GetConnString 獲取連接字符串
func (cn *SQLConnInfo) GetConnString() string {
	return fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;encrypt=disable", cn.ServerIP, cn.DbName, cn.UserID, cn.Password, cn.Port)
}
