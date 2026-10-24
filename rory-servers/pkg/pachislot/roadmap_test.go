package pachislot

import (
	"testing"
	"time"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

func TestRoadMap_SpinWin(t *testing.T) {
	m := NewRoadMap(10)

	t.Run("Play 1000 Round", func(tt *testing.T) {
		cnt := 1000
		for i := 0; i < cnt; i++ {
			m.Inc()
		}

		if m.TotalCnt != cnt {
			t.Errorf("TotalCnt want %d, but %d", cnt, m.TotalCnt)
		}

		if m.SpinCnt != cnt {
			t.Errorf("SpinCnt want %d, but %d", cnt, m.SpinCnt)
		}
	})

	t.Run("Play 500 Round Again", func(tt *testing.T) {
		old := 1000
		cnt := 500
		for i := 0; i < cnt; i++ {
			m.Inc()
		}

		if m.TotalCnt != cnt+old {
			t.Errorf("TotalCnt want %d, but %d", cnt+old, m.TotalCnt)
		}

		if m.SpinCnt != cnt+old {
			t.Errorf("SpinCnt want %d, but %d", cnt+old, m.SpinCnt)
		}
	})

	t.Run("Win BB", func(tt *testing.T) {

		old := 1500

		s := m.Win(pb.Scenario_BB)

		if s.Feature != pb.Scenario_BB {
			t.Errorf("Snapshot Feature want %v, but %v", pb.Scenario_BB, s.Feature)
		}

		if s.Count != old {
			t.Errorf("Snapshot Count want %d, but %d", old, s.Count)
		}

		if m.TotalCnt != old {
			t.Errorf("TotalCnt want %d, but %d", old, m.TotalCnt)
		}

		if m.SpinCnt != 0 {
			t.Errorf("SpinCnt want %d, but %d", 0, m.SpinCnt)
		}
	})

	t.Run("Play 789 Round", func(tt *testing.T) {
		old := 1500
		cnt := 789

		for i := 0; i < cnt; i++ {
			m.Inc()
		}

		if m.TotalCnt != cnt+old {
			t.Errorf("TotalCnt want %d, but %d", cnt+old, m.TotalCnt)
		}

		if m.SpinCnt != cnt {
			t.Errorf("SpinCnt want %d, but %d", cnt, m.SpinCnt)
		}
	})

	t.Run("Win RB", func(tt *testing.T) {

		old := 2289
		last := 789

		s := m.Win(pb.Scenario_RB)

		if s.Feature != pb.Scenario_RB {
			t.Errorf("Snapshot Feature want %v, but %v", pb.Scenario_RB, s.Feature)
		}

		if s.Count != last {
			t.Errorf("Snapshot Count want %d, but %d", last, s.Count)
		}

		if m.TotalCnt != old {
			t.Errorf("TotalCnt want %d, but %d", old, m.TotalCnt)
		}

		if m.SpinCnt != 0 {
			t.Errorf("SpinCnt want %d, but %d", 0, m.SpinCnt)
		}
	})

	t.Log(m)
}

func TestRoadMap_add(t *testing.T) {

	expCap := 10

	m := NewRoadMap(expCap)

	for i := 0; i < 100; i++ {
		m.add(&RoadMapSnapshot{Feature: pb.Scenario_BB, Time: time.Now()})
		m.add(&RoadMapSnapshot{Feature: pb.Scenario_RB, Time: time.Now()})
	}

	if len(m.Snapshots) != expCap {
		t.Errorf("Snapshot length want %d, but %d", expCap, len(m.Snapshots))
	}

	if cap(m.Snapshots) != expCap {
		t.Errorf("Snapshot capacity want %d, but %d", expCap, cap(m.Snapshots))
	}

	for i := 0; i < len(m.Snapshots); i += 2 {
		if m.Snapshots[i].Feature != pb.Scenario_RB {
			t.Errorf("Snapshot[%d] Feature want %d, but %d", i, pb.Scenario_RB, m.Snapshots[i].Feature)
		}

		if m.Snapshots[i+1].Feature != pb.Scenario_BB {
			t.Errorf("Snapshot[%d] Feature want %d, but %d", i+1, pb.Scenario_RB, m.Snapshots[i+1].Feature)
		}
	}
}

// Append-None-Limit-Cap 	BenchmarkRoadMapSnapshot-4   	12388686	        84.7 ns/op	      90 B/op	       1 allocs/op
// All Shift 				BenchmarkRoadMapSnapshot-4   	22911481	        49.3 ns/op	      48 B/op	       1 allocs/op
func BenchmarkRoadMapSnapshot(b *testing.B) {

	expCap := 10

	m := NewRoadMap(expCap)

	for i := 0; i < b.N; i++ {
		m.add(&RoadMapSnapshot{Feature: pb.Scenario_BB})
	}
}
