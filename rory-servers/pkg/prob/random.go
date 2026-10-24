package prob

import (
	"math/rand"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// Random 可Add物件及機率，再以Choose隨機選出一個
type Random struct {
	total  int // 總重量
	random int // 骰出得值
	target interface{}
	objs   []interface{}
	split  []int // 分水嶺
}

// Add 加入物件及機率
func (r *Random) Add(i interface{}, weight int) {
	if weight == 0 {
		return
	}
	r.total += weight
	r.objs = append(r.objs, i)
	r.split = append(r.split, r.total)
}

// Choose 隨機選出一個
func (r *Random) Choose() interface{} {
	if r.total == 0 {
		return nil
	}

	// 隨機產生該總重量
	r.random = rand.Intn(r.total)

	// 找出是哪個
	for i, row := range r.split {
		if r.random < row {
			r.target = r.objs[i]
			return r.objs[i]
		}
	}

	return nil
}

// PrintOptions 印出所有選項
func (r *Random) PrintOptions() {
	mlog.Trace("Random Options(%d), Dice: %d", len(r.objs), r.random)
	for i, row := range r.objs {
		mlog.Trace(" - %v(%d)", row, r.split[i])
	}
}
