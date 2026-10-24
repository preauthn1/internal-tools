package dbutil

import (
	"bytes"
	"fmt"
	"time"
)

// DBUtilcmd db util
type DBUtilcmd struct{}

// NewDBUtilCmd 创建數據轉發指令
func NewDBUtilCmd() *DBUtilcmd {
	cvcmd := new(DBUtilcmd)
	return cvcmd
}

// Name name
func (hv *DBUtilcmd) Name() string {
	return "dbutil"
}

// Help help
func (hv *DBUtilcmd) Help() string {
	return "數據轉發工具 : dbutil status | dbutil output <dbname> center,pergamelog,club"
}

// Exec exec
func (hv *DBUtilcmd) Exec(args []string) string {
	lenarg := len(args)
	switch lenarg {
	case 1:
		stype := args[0]
		if stype == "status" {
			var buf bytes.Buffer
			buf.WriteString(database.GetStatus())
			buf.WriteString(dealLogDataBase.GetStatus())
			buf.WriteString(clubDataBase.GetStatus())
			return buf.String()
		}
	case 2:
		stype := args[0]
		ltype := args[1]

		if stype == "output" {
			return hv.OutputHistory(ltype)
		}
	}
	return "error.seeref usage"
}

// OutputHistory history
func (hv *DBUtilcmd) OutputHistory(ltype string) string {
	start := time.Now()
	switch ltype {
	case "center":
		database.History.OutputHistory(ltype)
		return fmt.Sprintf("output done, used %s", time.Since(start))
	case "pergamelog":
		dealLogDataBase.History.OutputHistory(ltype)
		return fmt.Sprintf("output done, used %s", time.Since(start))
	case "club":
		clubDataBase.History.OutputHistory(ltype)
		return fmt.Sprintf("output done, used %s", time.Since(start))
	}
	return "error. database name not found"
}
