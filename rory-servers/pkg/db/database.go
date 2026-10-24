package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"path"
	"time"

	_ "example.invalid/denisenkom/go-mssqldb" // driver

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// DatabaseEngine db engine
type DatabaseEngine interface {
	ExecProc(*Procedure) (RecordSet, error)
	ExecSQL(string) (RecordSet, error)
	ConnString() string
}

// NewDatabaseEngine new db engine
func NewDatabaseEngine(dao DatabaseEngine) DatabaseEngine {
	return dao
}

// ProcError 錯誤返回
type ProcError struct {
	Code    int
	Message string
	PS      string
}

// ExecHistory history
type ExecHistory struct {
	HistoryList []*ExecData
	queue       chan *ExecData
}

// ExecData exec data
type ExecData struct {
	TimeStamp time.Time
	ExecTime  time.Duration
	SQLString string
}

// NewExecHistory histroy
func NewExecHistory() *ExecHistory {
	h := new(ExecHistory)
	h.HistoryList = make([]*ExecData, 1000)
	h.queue = make(chan *ExecData, 100000)
	go func() {
		for data := range h.queue {
			h.HistoryList = h.HistoryList[1:]
			h.HistoryList = append(h.HistoryList, data)
		}
		mlog.Error("History queue over")
	}()
	return h
}

// AddHistory add history
func (h *ExecHistory) AddHistory(sql string, t time.Duration) {
	data := new(ExecData)
	data.TimeStamp = time.Now()
	data.SQLString = sql
	data.ExecTime = t
	h.queue <- data
	//h.HistoryList = h.HistoryList[1:]
	//h.HistoryList = append(h.HistoryList, data)
}

// GetStatus get status
func (h *ExecHistory) GetStatus() string {
	max := time.Duration(0)
	min := time.Duration(0)
	for _, data := range h.HistoryList {
		if data == nil {
			continue
		}
		if min == 0 {
			min = data.ExecTime
		} else if data.ExecTime < min {
			min = data.ExecTime
		}
		if data.ExecTime > max {
			max = data.ExecTime
		}
	}
	return fmt.Sprintf("ExecHistory : %s execute time avg %s, size %d, queue size %d, min %s, max %s", h.GetIOPS(), h.GetAvg(), len(h.HistoryList), len(h.queue), min, max)
}

// GetAvg get avg
func (h *ExecHistory) GetAvg() time.Duration {
	size := int32(0)
	total := time.Duration(0)
	for _, data := range h.HistoryList {
		if data == nil {
			continue
		}
		size++
		total += data.ExecTime
	}
	if size == 0 {
		return 0
	}

	return time.Duration(int64(total) / int64(size))
}

// GetIOPS get iops
func (h *ExecHistory) GetIOPS() string {
	size := int(0)
	lastTime := time.Now()
	firstTime := time.Now()
	for _, data := range h.HistoryList {
		if data == nil {
			break
		}
		size++
		if lastTime.Before(data.TimeStamp) {
			lastTime = data.TimeStamp
		}
		if firstTime.After(data.TimeStamp) {
			firstTime = data.TimeStamp
		}
	}
	if size > 0 {
		subTime := lastTime.Sub(firstTime)
		nanos := subTime.Nanoseconds()
		if nanos < 1 {
			return fmt.Sprintf("%d iops", size)
		}
		return fmt.Sprintf("%.2f iops", float32(size)/float32(nanos)*1000000000)
	}
	return "0 iops"
}

// OutputHistory output history
func (h *ExecHistory) OutputHistory(prefix string) {
	var t = time.Now()
	filename := fmt.Sprintf("%s_sql_%d%02d%02d_%02d_%02d_%02d.log",
		prefix,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second())
	buf := bytes.NewBufferString("")
	for _, d := range h.HistoryList {
		if d == nil {
			continue
		}
		buf.WriteString(d.TimeStamp.Format("2006-01-02 15:04:05.999999 "))
		buf.WriteString(d.ExecTime.String())
		buf.WriteString(" ")
		buf.WriteString(d.SQLString)
		buf.WriteString("\n")
	}
	err := util.FilePutContents(path.Join("log", filename), buf.String())
	if err != nil {
		mlog.Error("output history error : %s", err)
	}
}

// Database 數據庫類
type Database struct {
	Connstring string
	Connpool   *ConnPool
	History    *ExecHistory
}

// NewDatabase 創建數據庫對象
func NewDatabase(connstring string, poolsize int) *Database {
	d := new(Database)
	d.Connstring = connstring
	d.Connpool = NewConnPool(poolsize, connstring)
	d.History = NewExecHistory()
	//d.TickerStatus()
	return d
}

// GetStatus get status
func (d *Database) GetStatus() string {
	s := fmt.Sprintf(d.History.GetStatus()+", %s\n", d.ConnString())
	s += fmt.Sprintf("ConnPool : current used %d / %d, %s\n", d.Connpool.currentSize, d.Connpool.poolSize, d.ConnString())
	return s
}

//func (d *Database) TickerStatus() {
//	go func() {
//		for {
//			time.Sleep(time.Duration(30) * time.Second)
//			mlog.Debug(d.History.GetStatus()+", %s", d.ConnString())
//			mlog.Debug("ConnPool : current used %d / %d, %s", d.Connpool.currentSize, d.Connpool.poolSize, d.ConnString())
//		}
//	}()
//}

// ConnString 數據庫連接字符串
func (d *Database) ConnString() string {
	return d.Connstring
}

// ExecSQL exec SQL
func (d *Database) ExecSQL(sql string) (RecordSet, error) {
	conn, err := d.Connpool.Acquire()
	if err != nil {
		return nil, err
	}

	defer d.Connpool.Release(conn)
	rows, err := conn.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	// 結果集
	retset := RecordSet{}
	rettable, err := parseRowsRet(rows)

	if err != nil {
		return nil, err
	}

	retset = append(retset, rettable)

	for rows.NextResultSet() {

		rettable, err := parseRowsRet(rows)

		if err != nil {
			return nil, err
		}

		retset = append(retset, rettable)
	}

	return retset, nil
}

// ExecProc 執行存儲過程
func (d *Database) ExecProc(p *Procedure) (RecordSet, error) {
	conn, err := d.Connpool.Acquire()
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		d.History.AddHistory(p.GetSQLString(), elapsed)
	}()

	if err != nil {
		return nil, err
	}

	defer d.Connpool.Release(conn)

	sql := p.GetSQLString()

	rows, err := conn.Query(sql)

	if err != nil {
		printError(err, sql)
		return nil, err
	}
	defer rows.Close()

	// 結果集
	retset := RecordSet{}
	rettable, err := parseRowsRet(rows)

	if err != nil {
		return nil, err
	}

	retset = append(retset, rettable)

	for rows.NextResultSet() {

		rettable, err := parseRowsRet(rows)

		if err != nil {
			return nil, err
		}

		retset = append(retset, rettable)
	}

	if rows.Err() != nil {
		printError(rows.Err(), sql)
		return nil, rows.Err()
	}

	return retset, nil
}
func printError(err error, sql string) {
	data := ProcError{500, fmt.Sprintf("%v", err), sql}
	mlog.NewDefaultLogger().SetCallDepLen(8).SetData(data).Error("執行存儲過程錯誤 %v", err)
}

// 解析行集
func parseRowsRet(rows *sql.Rows) (*RecordTable, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if cols == nil {
		return nil, nil
	}

	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	rettable := &RecordTable{}
	rettable.Cols = cols

	for rows.Next() {
		err = rows.Scan(vals...)

		if err != nil {
			return nil, err
		}

		retrow := RecordRow{}

		for _, r := range vals {
			v := *r.(*interface{})

			retrow = append(retrow, v)
		}

		rettable.Rows = append(rettable.Rows, retrow)
	}

	return rettable, err
}
