package serverdefense

import (
	"fmt"
	"sync"
)

// DefenserData 防護數據
type DefenserData struct {
	SecData5   int // 5秒數據
	SecData10  int // 10秒數據
	SecData30  int // 30秒數據
	SecData60  int // 60秒數據
	SecData300 int // 5分鐘數據

	innerLocker *sync.RWMutex
}

// NewDefenserData 創建防護數據對象
func NewDefenserData() *DefenserData {
	defdata := new(DefenserData)
	defdata.SecData5 = 0
	defdata.SecData10 = 0
	defdata.SecData30 = 0
	defdata.SecData60 = 0
	defdata.SecData300 = 0

	defdata.innerLocker = &sync.RWMutex{}

	return defdata
}

// Increase 增加計數器
func (dd *DefenserData) Increase(num int) {
	dd.innerLocker.Lock()
	defer dd.innerLocker.Unlock()

	dd.SecData5 += num
	dd.SecData10 += num
	dd.SecData30 += num
	dd.SecData60 += num
	dd.SecData300 += num
}

// Clear 清除計數器
func (dd *DefenserData) Clear() {
	dd.innerLocker.Lock()
	defer dd.innerLocker.Unlock()

	dd.SecData5 = 0
	dd.SecData10 = 0
	dd.SecData30 = 0
	dd.SecData60 = 0
	dd.SecData300 = 0
}

// CheckLimit 檢查防護數據是否超出設定
func (dd *DefenserData) CheckLimit(limit *DefenserLimit) bool {
	dd.innerLocker.RLock()
	defer dd.innerLocker.RUnlock()

	if dd.SecData5 >= limit.SecLimit5 {
		return true
	}
	if dd.SecData10 >= limit.SecLimit10 {
		return true
	}
	if dd.SecData30 >= limit.SecLimit30 {
		return true
	}
	if dd.SecData60 >= limit.SecLimit60 {
		return true
	}
	if dd.SecData300 >= limit.SecLimit300 {
		return true
	}

	return false
}

// ClearSecData5 5
func (dd *DefenserData) ClearSecData5() {
	dd.innerLocker.Lock()
	defer dd.innerLocker.Unlock()

	dd.SecData5 = 0
}

// ClearSecData10 10
func (dd *DefenserData) ClearSecData10() {
	dd.innerLocker.Lock()
	defer dd.innerLocker.Unlock()

	dd.SecData10 = 0
}

// ClearSecData30 30
func (dd *DefenserData) ClearSecData30() {
	dd.innerLocker.Lock()
	defer dd.innerLocker.Unlock()

	dd.SecData30 = 0
}

// ClearSecData60 60
func (dd *DefenserData) ClearSecData60() {
	dd.innerLocker.Lock()
	defer dd.innerLocker.Unlock()

	dd.SecData60 = 0
}

// ClearSecData300 300
func (dd *DefenserData) ClearSecData300() {
	dd.innerLocker.Lock()
	defer dd.innerLocker.Unlock()

	dd.SecData300 = 0
}

func (dd *DefenserData) String() string {
	dd.innerLocker.RLock()
	dd.innerLocker.RUnlock()

	return fmt.Sprintf("5:%d,10:%d,30:%d,60:%d,300:%d", dd.SecData5, dd.SecData10, dd.SecData30, dd.SecData60, dd.SecData300)
}
