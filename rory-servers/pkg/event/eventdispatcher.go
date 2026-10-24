package event

import (
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// HandlerInfo 事件處理器信息
type HandlerInfo struct {
	Handler Handler
	Sort    int
}

// Dispatcher 事件處理器組
type Dispatcher struct {
	innerHandlers []*HandlerInfo       // 處理器集合
	asyncEvents   chan *AsyncEventInfo // 異步事件
	throttle      int                  // 並發數
}

// AsyncEventInfo 異步事件信息
type AsyncEventInfo struct {
	Evt      *Event
	Callback CallbackFun
}

// NewEventDispatcher 創建事件處理器組
func NewEventDispatcher() *Dispatcher {
	evg := new(Dispatcher)

	//evg.asyncEvents = make(chan *AsyncEventInfo, 1000)

	return evg
}

// NewAsyncEventDispatcher 創建異步事件處理器組
func NewAsyncEventDispatcher(throttle int) *Dispatcher {
	evg := new(Dispatcher)
	evg.asyncEvents = make(chan *AsyncEventInfo, 10000)
	evg.throttle = throttle

	for i := 0; i < throttle; i++ {
		//开启处理异步事件线程
		go func() {
			defer util.PanicErrStack()

			for {
				asinfo := <-evg.asyncEvents
				if asinfo != nil {
					ret := evg.OnEvent(asinfo.Evt)

					if asinfo.Callback != nil {
						callbackf := func() {
							asinfo.Callback(asinfo.Evt, ret)
						}

						util.SafeCall(callbackf)
					}
				}
			}
		}()
	}

	return evg
}

// Len sort 接口实现
func (evg *Dispatcher) Len() int {
	return len(evg.innerHandlers)
}

func (evg *Dispatcher) Less(i, j int) bool {
	ih := evg.innerHandlers[i]
	jh := evg.innerHandlers[j]

	return ih.Sort > jh.Sort
}

func (evg *Dispatcher) Swap(i, j int) {
	evg.innerHandlers[i], evg.innerHandlers[j] = evg.innerHandlers[j], evg.innerHandlers[i]
}

// AddHandler 添加事件處理器,可控制排序
func (evg *Dispatcher) AddHandler(h Handler, sortIndex int) *Dispatcher {
	if h == nil {
		panic("处理器不能为空")
	}

	hi := &HandlerInfo{
		Handler: h,
		Sort:    sortIndex}
	evg.innerHandlers = append(evg.innerHandlers, hi)

	//排序
	sort.Sort(evg)

	return evg
}

// PushHandler 添加事件處理器
func (evg *Dispatcher) PushHandler(h Handler) *Dispatcher {
	return evg.AddHandler(h, 0)
}

// RangeSort 遍歷處理器
func (evg *Dispatcher) RangeSort(f func(int, Handler) bool) {
	for _, v := range evg.innerHandlers {
		ret := f(v.Sort, v.Handler)

		if !ret {
			return
		}
	}
}

// Range 遍歷處理器
func (evg *Dispatcher) Range(f func(Handler) bool) {
	for _, v := range evg.innerHandlers {
		ret := f(v.Handler)

		if !ret {
			return
		}
	}
}

// OnEvent 同步觸發事件
func (evg *Dispatcher) OnEvent(evt *Event) (ret *Ret) {
	ret = NewEventRet()

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetCallDepLen(8).Error("觸發同步事件%v發生錯誤%v, errstack:%s", evt, err, errorStr)
			ret.Error = err
		}
	}()

	for _, h := range evg.innerHandlers {
		goon := h.Handler.OnEvent(evt, ret)

		if !goon {
			return
		}
	}

	return
}

// OnAsyncEvent 異步觸發事件
func (evg *Dispatcher) OnAsyncEvent(evt *Event, callback CallbackFun) {
	asinfo := new(AsyncEventInfo)
	asinfo.Evt = evt
	asinfo.Callback = callback

	evg.asyncEvents <- asinfo
}
