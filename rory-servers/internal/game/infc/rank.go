package infc

import (
	"sort"
)

// Rank 排行榜
type Rank []*Player

// Append 加入玩家
func (r Rank) Append(p *Player) Rank {
	return append([]*Player(r), p)
}

// Player 取得玩家
func (r Rank) Player(uid int32) *Player {
	for _, p := range r {
		if p.UID == uid {
			return p
		}
	}
	return nil
}

// Rank 查詢排名
func (r Rank) Rank(uid int32) int {
	sort.Sort(r)
	type s struct {
		score   int64
		players []int32
	}
	var scores []*s
	for _, p := range r {
		var last *s
		if len(scores) > 0 {
			last = scores[len(scores)-1]
		}
		if last != nil && last.score == p.Score {
			last.players = append(last.players, p.UID)
			continue
		}
		data := new(s)
		data.score = p.Score
		data.players = append(data.players, p.UID)
		scores = append(scores, data)
	}
	for i, s := range scores {
		for _, p := range s.players {
			if p == uid {
				return i + 1
			}
		}
	}
	return 0
}

// Winner 取得贏家
func (r Rank) Winner() []*Player {
	sort.Sort(r)
	var wins []*Player
	for _, p := range r {
		if p.Score == r.TopScore() {
			wins = append(wins, p)
		}
	}
	return wins
}

// TopScore 最高分
func (r Rank) TopScore() int64 {
	var top int64
	for _, p := range r {
		if p.State == StateGiveup {
			continue
		}
		if p.Score > top {
			top = p.Score
		}
	}
	return top
}

// Survivors 活人
func (r Rank) Survivors() []*Player {
	var arr []*Player
	for _, p := range r {
		if p.State == StateGaming {
			arr = append(arr, p)
		}
	}
	return arr
}

// IsOthersIncalculable 其餘玩家是否都不可計算分數
func (r Rank) IsOthersIncalculable(p *Player) bool {
	for _, other := range r {
		if other == p {
			continue
		}
		if other.State != StateGiveup {
			return false
		}
	}
	return true
}

// Len is the number of elements in the collection.
func (r Rank) Len() int {
	return len(r)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (r Rank) Less(i, j int) bool {
	if r[i].State != StateGaming && r[j].State == StateGaming {
		return false
	}
	if r[i].State == StateGaming && r[j].State != StateGaming {
		return true
	}
	return r[i].Score > r[j].Score
}

// Swap swaps the elements with indexes i and j.
func (r Rank) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
