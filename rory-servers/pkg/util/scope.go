package util

import (
	"sync"
	"time"
)

// Scope 表示一個範圍
type Scope struct {
	MinValue interface{}
	MaxValue interface{}
}

// IntScoper ?
type IntScoper interface {
	GetRandomValue() int32
}

// IntScope 表示一個int範圍
type IntScope struct {
	MinValue int32
	MaxValue int32
}

// NewIntScope 創建int範圍對象
func NewIntScope(min int32, max int32) *IntScope {
	s := new(IntScope)

	s.MinValue = min
	s.MaxValue = max

	return s
}

// GetRandomValue 獲取一個範圍內的int隨機值
func (s *IntScope) GetRandomValue() int32 {

	return int32(RndInt(int(s.MinValue), int(s.MaxValue)))
	/* equal
	r := rand.New(RndSeed)

	v := r.Intn(int(s.MaxValue - s.MinValue))

	return s.MinValue + int32(v)
	*/
}

// DynamicIntScope 動態的int範圍
type DynamicIntScope struct {
	*IntScope

	valueLock      *sync.RWMutex
	changeInterval int32
}

// NewDynamicIntScope 創建動態int範圍對象
func NewDynamicIntScope(f func() (int32, int32), changeinterval int32) *DynamicIntScope {
	dynacScope := new(DynamicIntScope)

	fminv, fmaxv := f()
	dynacScope.IntScope = NewIntScope(fminv, fmaxv)

	dynacScope.valueLock = &sync.RWMutex{}
	dynacScope.changeInterval = changeinterval

	go func() {
		for {
			time.Sleep(time.Duration(dynacScope.changeInterval) * time.Second)

			fminv, fmaxv := f()
			dynacScope.valueLock.Lock()
			dynacScope.IntScope = NewIntScope(fminv, fmaxv)
			dynacScope.valueLock.Unlock()
		}
	}()

	return dynacScope

}

// GetRandomValue ?
func (s *DynamicIntScope) GetRandomValue() int32 {
	s.valueLock.Lock()
	defer s.valueLock.Unlock()

	return s.IntScope.GetRandomValue()
}
