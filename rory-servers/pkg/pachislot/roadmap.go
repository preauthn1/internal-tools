package pachislot

import (
	"fmt"
	"sync"
	"time"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// RoadMap 轉數表
type RoadMap struct {
	sync.Mutex

	TotalCnt int
	SpinCnt  int

	FeaturesCnt map[pb.Scenario]int

	Snapshots []*RoadMapSnapshot
}

func (m *RoadMap) String() string {
	return fmt.Sprintf("總轉數[%d] 當前轉數[%d] BB[%d] RB[%d] AT[%d]", m.TotalCnt, m.SpinCnt, m.FeaturesCnt[pb.Scenario_BB], m.FeaturesCnt[pb.Scenario_RB], 0)
}

// NewRoadMap 建立轉數表
//  @Param buff 紀錄歷史暫存空間
func NewRoadMap(buff int) *RoadMap {
	m := new(RoadMap)
	m.FeaturesCnt = make(map[pb.Scenario]int)
	m.Snapshots = make([]*RoadMapSnapshot, buff)
	return m
}

// Win 贏得特殊獎勵
//  SpinCnt 清空
//  紀錄此次獎勵在第幾轉
//  Snapshot 保持最新在前，並不超過最大容量
func (m *RoadMap) Win(s pb.Scenario) *RoadMapSnapshot {
	m.Lock()
	defer m.Unlock()

	m.FeaturesCnt[s]++

	snapshot := &RoadMapSnapshot{Feature: s, Count: m.SpinCnt, Time: time.Now()}
	m.SpinCnt = 0

	m.add(snapshot)

	return snapshot
}

// Inc 有效轉數+1 (有投幣)
func (m *RoadMap) Inc() {
	m.TotalCnt++
	m.SpinCnt++
}

func (m *RoadMap) add(s *RoadMapSnapshot) {
	// Append-None-Limit-Cap
	// m.Snapshots = append(m.Snapshots, s)

	// All shift
	for i := len(m.Snapshots) - 1; i > 0; i-- {
		m.Snapshots[i] = m.Snapshots[i-1]
	}
	m.Snapshots[0] = s
}

// RoadMapSnapshot 轉數表快照
type RoadMapSnapshot struct {
	Feature pb.Scenario
	Count   int
	Time    time.Time
}
