package infc

import (
	"reflect"
	"sort"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

func TestRank(t *testing.T) {

	var r Rank

	p1 := NewPlayer(&gamelogic.GameUser{UserID: 1})
	p1.Score = 1100
	p1.State = StateGaming

	p2 := NewPlayer(&gamelogic.GameUser{UserID: 2})
	p2.Score = 1200
	p2.State = StateDeath

	p3 := NewPlayer(&gamelogic.GameUser{UserID: 3})
	p3.Score = 1300
	p3.State = StateGaming

	p4 := NewPlayer(&gamelogic.GameUser{UserID: 4})
	p4.Score = 1400
	p4.State = StateGaming

	t.Run("Add", func(tt *testing.T) {

		r = r.Append(p1)
		r = r.Append(p2)
		r = r.Append(p3)
		r = r.Append(p4)

		exp := Rank([]*Player{p1, p2, p3, p4})

		if !reflect.DeepEqual(r, exp) {
			tt.Errorf("Rank Append want be %+v, but %+v", exp, r)
		}
	})

	t.Run("Sort Score", func(tt *testing.T) {

		sort.Sort(r)

		exp := Rank([]*Player{p4, p3, p1, p2})

		if !reflect.DeepEqual(r, exp) {
			tt.Errorf("Rank Sorted want be %+v, but %+v", exp, r)
		}

		for _, p := range r {
			t.Logf("%+v", p)
		}
	})

	t.Run("Player", func(tt *testing.T) {

		if p1 != r.Player(p1.UID) {
			tt.Errorf("Rank Player(%d) want %+v, but %+v", p1.UID, p1, r.Player(p1.UID))
		}

		if p2 != r.Player(p2.UID) {
			tt.Errorf("Rank Player(%d) want %+v, but %+v", p2.UID, p2, r.Player(p2.UID))
		}

		if p3 != r.Player(p3.UID) {
			tt.Errorf("Rank Player(%d) want %+v, but %+v", p3.UID, p3, r.Player(p3.UID))
		}

		if p4 != r.Player(p4.UID) {
			tt.Errorf("Rank Player(%d) want %+v, but %+v", p4.UID, p4, r.Player(p4.UID))
		}
	})

	t.Run("Winner", func(tt *testing.T) {
		if p4 != r.Winner()[0] {
			tt.Errorf("Rank Winner want %+v, but %+v", p4, r.Winner()[0])
		}
	})

	t.Run("TopScore", func(tt *testing.T) {
		if 1400 != r.TopScore() {
			tt.Errorf("Rank TopScore want %d, but %d", p4.Score, r.TopScore())
		}
	})

	t.Run("Survivors", func(tt *testing.T) {

		survivors := r.Survivors()
		exp := []*Player{p4, p3, p1}

		if !reflect.DeepEqual(exp, survivors) {
			tt.Errorf("Rank Survivors want %+v, but %+v", exp, survivors)
		}
	})

	t.Run("Rank Diff Score", func(t *testing.T) {
		if r.Rank(1) != 3 {
			t.Errorf("Rank User1 want %d, but %d", 3, r.Rank(1))
		}
		if r.Rank(2) != 4 {
			t.Errorf("Rank User2 want %d, but %d", 4, r.Rank(2))
		}
		if r.Rank(3) != 2 {
			t.Errorf("Rank User3 want %d, but %d", 2, r.Rank(3))
		}
		if r.Rank(4) != 1 {
			t.Errorf("Rank User4 want %d, but %d", 1, r.Rank(4))
		}
	})

	t.Run("Rank Same Score", func(t *testing.T) {
		p3.Score = 1400
		if r.Rank(1) != 2 {
			t.Errorf("Rank User1 want %d, but %d", 2, r.Rank(1))
		}
		if r.Rank(2) != 3 {
			t.Errorf("Rank User2 want %d, but %d", 3, r.Rank(2))
		}
		if r.Rank(3) != 1 {
			t.Errorf("Rank User3 want %d, but %d", 1, r.Rank(3))
		}
		if r.Rank(4) != 1 {
			t.Errorf("Rank User4 want %d, but %d", 1, r.Rank(4))
		}
	})
}
