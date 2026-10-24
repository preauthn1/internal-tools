package serverdefense

import (
	"fmt"
)

// DefenserLimit 防護信息
type DefenserLimit struct {
	SecLimit5   int // 5秒閥值
	SecLimit10  int // 10秒閥值
	SecLimit30  int // 30秒閥值
	SecLimit60  int // 60秒閥值
	SecLimit300 int // 5分鐘閥值
}

// NewDefenserLimit 創建防護信息對象
func NewDefenserLimit(secLimit5 int, secLimit10 int, secLimit30 int, secLimit60 int, secLimit300 int) *DefenserLimit {
	definfo := new(DefenserLimit)
	definfo.SecLimit5 = secLimit5
	definfo.SecLimit10 = secLimit10
	definfo.SecLimit30 = secLimit30
	definfo.SecLimit60 = secLimit60
	definfo.SecLimit300 = secLimit300

	return definfo
}

func (limit *DefenserLimit) String() string {
	return fmt.Sprintf("5:%d,10:%d,30:%d,60:%d,300:%d", limit.SecLimit5, limit.SecLimit10, limit.SecLimit30, limit.SecLimit60, limit.SecLimit300)
}
