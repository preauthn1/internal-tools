package httputil

import (
	"bytes"
	"fmt"
	"path"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

var history = NewAPIQueryHistory()

// DebugMode ?
func DebugMode(b bool) {
	history.Debug = b
}

// APIQueryHistory ?
type APIQueryHistory struct {
	HistoryList []*QueryData
	queue       chan *QueryData
	Debug       bool
}

// QueryData ?
type QueryData struct {
	TimeStamp   time.Time
	ExecTime    time.Duration
	QueryString string
}

// ToString ?
func (d *QueryData) ToString() string {
	buf := bytes.NewBufferString("")
	buf.WriteString(d.TimeStamp.Format("2006-01-02 15:04:05.999999 "))
	buf.WriteString(d.ExecTime.String())
	buf.WriteString(" ")
	buf.WriteString(d.QueryString)
	buf.WriteString("\n")
	return buf.String()
}

// NewAPIQueryHistory ?
func NewAPIQueryHistory() *APIQueryHistory {
	h := new(APIQueryHistory)
	h.HistoryList = make([]*QueryData, 1000)
	h.queue = make(chan *QueryData, 100000)
	go func() {
		for data := range h.queue {
			h.HistoryList = h.HistoryList[1:]
			h.HistoryList = append(h.HistoryList, data)
		}
		mlog.Error("History queue over")
	}()
	h.Debug = false
	return h
}

// AddHistory ?
func (h *APIQueryHistory) AddHistory(s string, t time.Duration) {
	data := new(QueryData)
	data.TimeStamp = time.Now()
	data.QueryString = s
	data.ExecTime = t
	h.queue <- data
	if h.Debug {
		fmt.Print(s)
	}
	//h.HistoryList = h.HistoryList[1:]
	//h.HistoryList = append(h.HistoryList, data)
}

// GetStatus ?
func (h *APIQueryHistory) GetStatus() string {
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
	return fmt.Sprintf("APIQueryHistory : %s execute time avg %s, size %d, queue size %d, min %s, max %s", h.GetIOPS(), h.GetAvg(), len(h.HistoryList), len(h.queue), min, max)
}

// GetAvg ?
func (h *APIQueryHistory) GetAvg() time.Duration {
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

// GetIOPS ?
func (h *APIQueryHistory) GetIOPS() string {
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

// OutputHistory ?
func (h *APIQueryHistory) OutputHistory() {
	prefix := "ApiQuery"
	var t = time.Now()
	filename := fmt.Sprintf("%s_%d%02d%02d_%02d_%02d_%02d.log",
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
		buf.WriteString(d.ToString())
	}
	err := util.FilePutContents(path.Join("log", filename), buf.String())
	if err != nil {
		mlog.Error("output history error : %s", err)
	}
}
