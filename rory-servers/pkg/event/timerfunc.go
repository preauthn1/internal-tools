package event

import (
	"fmt"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
const (
	EVTNAMETIMERFUNC     = "timerfunc" // 事件名稱
	EVTArgsNameTimerFunc = "timerinfo" // 參數名稱
)

// TimerStatus
const (
	TimerStatusStop = 0
	TimerStatusRun  = 1
)

// TimerType
const (
	TimerTypeTimer  = 0
	TimerTypeTicker = 1
)

// TimerFunc 延遲執行一次的函數
type TimerFunc struct {
	WaitSeconds      int32     // 延遲時間
	WaitMilliSeconds int32     // 延遲時間
	Func             func()    // 處理函數
	Status           int       // 狀態
	Loop             bool      // 是否循環執行
	Name             string    // 名稱
	BeginTime        time.Time // 開始時間
	TimerType        int       // timer類型
}

// Exec 執行延遲函數
func (tf *TimerFunc) Exec() {
	defer util.PanicErrStack()
	if tf.Status == TimerStatusRun {
		//example.invalid("执行timer func %v", tf)
		tf.Func()
	}
}

// ElapsedSeconds 消耗的秒數
func (tf *TimerFunc) ElapsedSeconds() int32 {
	return int32(time.Now().Sub(tf.BeginTime).Seconds())
}

// Stop ?
func (tf *TimerFunc) Stop() {
	tf.Status = TimerStatusStop
}

// IsStop ?
func (tf *TimerFunc) IsStop() bool {
	return tf.Status == TimerStatusStop
}

// TimerFuncManager 延遲函數管理器
type TimerFuncManager struct {
	evtDispatcher *Dispatcher
	timerFuncs    sync.Map //key:string value:*TimerFunc
}

// NewTimerFuncManager 創建延遲函數管理器對象
func NewTimerFuncManager(evtDispatcher *Dispatcher) *TimerFuncManager {
	tfm := new(TimerFuncManager)
	tfm.evtDispatcher = evtDispatcher

	tfm.registerTimeFuncEvtHandler()

	return tfm
}

// NewTimer ?
func (tfm *TimerFuncManager) NewTimer(wait int32, f func()) *TimerFunc {
	return tfm.NewNamedTimer(wait, f, "")
}

// NewNamedTimer 創建一個timer
func (tfm *TimerFuncManager) NewNamedTimer(wait int32, f func(), timername string) *TimerFunc {
	return tfm.NewNamedMilliSecTimer(wait*1000, f, timername)
}

// NewMilliSecTimer ?
func (tfm *TimerFuncManager) NewMilliSecTimer(wait int32, f func()) *TimerFunc {
	return tfm.NewNamedMilliSecTimer(wait, f, "")
}

// NewNamedMilliSecTimer 創建一個毫秒級的timer
func (tfm *TimerFuncManager) NewNamedMilliSecTimer(wait int32, f func(), timername string) *TimerFunc {
	tf := new(TimerFunc)
	tf.WaitSeconds = wait / 1000
	tf.WaitMilliSeconds = wait
	tf.Loop = false
	tf.Status = TimerStatusRun
	tf.Func = f
	tf.Name = timername
	tf.BeginTime = time.Now()
	tf.TimerType = TimerTypeTimer

	// 測試模式下的快轉倍數 TimerFastMode
	if conf.TestMode && conf.TimerFastMode > 1 {
		wait = wait / int32(conf.TimerFastMode)
		tf.WaitSeconds = wait / 1000
		tf.WaitMilliSeconds = wait
	}

	go func() {
		time.Sleep(time.Duration(wait) * time.Millisecond)
		if tf.Status == TimerStatusStop {
			return
		}
		tfm.publishTimerFuncEvent(tf)
	}()

	if timername != "" {
		tfm.timerFuncs.Store(timername, tf)
		mlog.Debug("創建timerfunc %s , status %v", timername, tf.Status)
	}

	return tf
}

// NewTicker ?
func (tfm *TimerFuncManager) NewTicker(wait int32, f func()) *TimerFunc {
	return tfm.NewNamedTicker(wait, f, "")
}

// NewNamedTicker 創建一個 tick
func (tfm *TimerFuncManager) NewNamedTicker(wait int32, f func(), timername string) *TimerFunc {
	return tfm.NewNamedMilliSecTicker(wait*1000, f, timername)
}

// NewMilliSecTicker ?
func (tfm *TimerFuncManager) NewMilliSecTicker(wait int32, f func()) *TimerFunc {
	return tfm.NewNamedMilliSecTicker(wait, f, "")
}

// NewNamedMilliSecTicker 創建一個毫秒級的 tick
func (tfm *TimerFuncManager) NewNamedMilliSecTicker(wait int32, f func(), timername string) *TimerFunc {
	tf := new(TimerFunc)
	tf.WaitSeconds = wait / 1000
	tf.WaitMilliSeconds = wait
	tf.Loop = true
	tf.Status = TimerStatusRun
	tf.Func = f
	tf.Name = timername
	tf.BeginTime = time.Now()
	tf.TimerType = TimerTypeTicker
	// 測試模式下的快轉倍數 TimerFastMode
	if conf.TestMode && conf.TimerFastMode > 1 {
		wait = wait / int32(conf.TimerFastMode)
		tf.WaitSeconds = wait / 1000
		tf.WaitMilliSeconds = wait
	}
	go func() {
		for {
			time.Sleep(time.Duration(wait) * time.Millisecond)
			if tf.Status == TimerStatusStop {
				break
			}

			tfm.publishTimerFuncEvent(tf)
		}
	}()

	if timername != "" {
		tfm.timerFuncs.Store(timername, tf)
	}

	return tf
}

// NewNamedRandomTicker 創建一個timer tick
func (tfm *TimerFuncManager) NewNamedRandomTicker(rndScope *util.IntScope, f func(), timername string) *TimerFunc {
	wait := int32(rndScope.GetRandomValue())
	tf := new(TimerFunc)
	tf.WaitSeconds = rndScope.MaxValue
	tf.WaitMilliSeconds = rndScope.MaxValue * 1000
	tf.Loop = true
	tf.Status = TimerStatusRun
	tf.Func = f
	tf.Name = timername
	tf.BeginTime = time.Now()
	tf.TimerType = TimerTypeTicker
	// 測試模式下的快轉倍數 TimerFastMode
	if conf.TestMode && conf.TimerFastMode > 1 {
		wait = wait / int32(conf.TimerFastMode)
		tf.WaitSeconds = wait / 1000
		tf.WaitMilliSeconds = wait
	}
	go func() {
		for {
			time.Sleep(time.Duration(wait) * time.Second)

			if tf.Status == TimerStatusStop {
				break
			}
			tfm.publishTimerFuncEvent(tf)
		}
	}()

	if timername != "" {
		tfm.timerFuncs.Store(timername, tf)
	}

	return tf
}

// GetTimerFunc 獲取延遲函數
func (tfm *TimerFuncManager) GetTimerFunc(timername string) *TimerFunc {
	if val, ok := tfm.timerFuncs.Load(timername); ok {
		var tf = val.(*TimerFunc)
		return tf
	}
	return nil
}

// PrintKey ?
func (tfm *TimerFuncManager) PrintKey() {
	tfm.timerFuncs.Range(func(key, value interface{}) bool {
		fmt.Println(key)
		return true
	})
}

// Keys 取得所有 TimerFunc的Key
func (tfm *TimerFuncManager) Keys() []string {
	var keys []string
	tfm.timerFuncs.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

// Len 取得所有 註冊名稱的 TimerFunc 數量
func (tfm *TimerFuncManager) Len() int {
	var cnt int
	tfm.timerFuncs.Range(func(key, value interface{}) bool {
		cnt++
		return true
	})
	return cnt
}

// RemoveTimer 移除timer
func (tfm *TimerFuncManager) RemoveTimer(timername string) *TimerFunc {
	if val, ok := tfm.timerFuncs.Load(timername); ok {
		var tf = val.(*TimerFunc)
		tf.Status = TimerStatusStop
		tfm.timerFuncs.Delete(timername)
		return tf
	}
	return nil
}

func (tfm *TimerFuncManager) publishTimerFuncEvent(tf *TimerFunc) {
	tfm.evtDispatcher.OnAsyncEvent(newTimerFuncEvent(tf), nil)
}

func (tfm *TimerFuncManager) registerTimeFuncEvtHandler() {
	h := NewCommonHandler("timerfunc", func(evt *Event, evtret *Ret) bool {

		tf := getTimerFuncEventArg(evt.Args)

		util.SafeCall(tf.Exec)

		//移除命名的timer
		if tf.TimerType == TimerTypeTimer && tf.Name != "" {
			tfm.RemoveTimer(tf.Name)
		}

		return true
	}, EVTNAMETIMERFUNC)

	tfm.evtDispatcher.PushHandler(h)
}

// 創建延遲函數事件
func newTimerFuncEvent(tf *TimerFunc) *Event {
	evt := new(Event)
	evt.Name = EVTNAMETIMERFUNC
	evt.Args = NewEventArgs()
	evt.Args[EVTArgsNameTimerFunc] = tf

	return evt
}

// 獲取延遲函數事件參數
func getTimerFuncEventArg(evtArgs Args) *TimerFunc {
	return evtArgs[EVTArgsNameTimerFunc].(*TimerFunc)
}
