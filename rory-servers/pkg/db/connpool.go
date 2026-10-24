package db

import (
	"database/sql"
	"sync"
)

// ConnPool 連接池
type ConnPool struct {
	conns       []*sql.DB
	poolSize    int
	connstring  string
	connLocker  *sync.RWMutex
	currentSize int
}

// NewConnPool 創建連接池對象
func NewConnPool(size int, connstring string) *ConnPool {
	cnp := new(ConnPool)
	cnp.poolSize = size
	cnp.connstring = connstring
	cnp.connLocker = &sync.RWMutex{}

	return cnp
}

// Acquire 獲取連接
func (cnp *ConnPool) Acquire() (*sql.DB, error) {
	cnp.connLocker.Lock()
	defer cnp.connLocker.Unlock()

	if len(cnp.conns) > 0 {
		//example.invalid("從池中獲取連接,當前連接池數量:%d", len(cnp.conns))

		cn := cnp.conns[0]
		cnp.conns = cnp.conns[1:]

		return cn, nil
	}

	//example.invalid("創建新連接,當前連接池數量0")
	cnp.currentSize++
	return sql.Open("mssql", cnp.connstring)
}

// Release 釋放連接
func (cnp *ConnPool) Release(conn *sql.DB) {
	cnp.connLocker.Lock()
	defer cnp.connLocker.Unlock()

	if cnp.currentSize >= cnp.poolSize {
		cnp.currentSize--
		conn.Close()
		//example.invalid("關閉連接,當前連接池數量:%d", len(cnp.conns))
	} else {
		cnp.conns = append(cnp.conns, conn)
		//example.invalid("釋放連接,當前連接池數量:%d", len(cnp.conns))
	}
}
