package ddz

type sItem struct {
	number int32
	count  int32
}

type sortCtx struct {
	arr []*sItem
}

func (sc *sortCtx) Len() int {
	return len(sc.arr)
}

func (sc *sortCtx) Swap(i, j int) {
	sc.arr[i], sc.arr[j] = sc.arr[j], sc.arr[i]
}

func (sc *sortCtx) Less(i, j int) bool {
	return sc.arr[i%sc.Len()].count > sc.arr[j%sc.Len()].count
}

type sortCtxByNumber struct {
	sortCtx
}

func (sc *sortCtxByNumber) Less(i, j int) bool {
	return sc.arr[i%sc.Len()].number < sc.arr[j].number
}

func (sc *sortCtx) GetByIndex(i int) int32 {
	return sc.arr[i%sc.Len()].count
}

func (sc *sortCtx) GetValueByIndex(i int) int32 {
	return sc.arr[i%sc.Len()].number
}

func (sc *sortCtx) GetMainCardMax() int32 {
	var tmp int32 = 0
	c := sc.arr[0].count
	for _, v := range sc.arr {
		if v.count != c {
			continue
		}
		if tmp < v.number {
			tmp = v.number
		}
	}
	return tmp
}

//所有的组合是否满足指定数量限制
func (sc *sortCtx) CountCondition(c int32) bool {
	for i := 0; i < sc.Len(); i++ {
		if sc.arr[i].count != c {
			return false
		}
	}
	return true
}

//返回指定数量的牌 的数量
func (sc *sortCtx) CdiCount(c int32) int {
	var n int = 0
	for i := 0; i < sc.Len(); i++ {
		if sc.arr[i].count == c {
			n++
		}
	}
	return n
}
