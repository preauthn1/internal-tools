package model

import (
	"context"
	"sync"
	"testing"
	"time"
)

// slowModel 每次呼叫花 d，記下每次呼叫的開始和結束時間
type slowModel struct {
	d         time.Duration
	mu        sync.Mutex
	intervals [][2]time.Time
}

func (m *slowModel) Respond(ctx context.Context, r Request) (*Reply, Usage) {
	start := time.Now()
	time.Sleep(m.d)
	m.mu.Lock()
	m.intervals = append(m.intervals, [2]time.Time{start, time.Now()})
	m.mu.Unlock()
	return &Reply{Values: map[string]Value{}}, Usage{}
}

func (m *slowModel) overlapped() bool {
	a, b := m.intervals[0], m.intervals[1]
	return a[0].Before(b[1]) && b[0].Before(a[1])
}

func callTwice(m ChatModel) {
	var wg sync.WaitGroup
	for range 2 {
		wg.Go(func() { m.Respond(context.Background(), Request{Name: "Slow"}) })
	}
	wg.Wait()
}

// 本機模型同時兩個呼叫會排隊：第二個等第一個結束才開始；沒包排隊的（Haiku）會同時跑
func TestQueuedRunsOneAtATime(t *testing.T) {
	slow := &slowModel{d: 200 * time.Millisecond}
	callTwice(NewQueued(slow))
	if slow.overlapped() {
		t.Errorf("排隊的呼叫重疊了：%v", slow.intervals)
	}

	parallel := &slowModel{d: 200 * time.Millisecond}
	callTwice(parallel)
	if !parallel.overlapped() {
		t.Errorf("沒排隊的呼叫應該同時跑：%v", parallel.intervals)
	}
}

// 排隊等的時候 ctx 結束就放棄，回 nil、不呼叫模型
func TestQueuedWaitCancelled(t *testing.T) {
	slow := &slowModel{d: 300 * time.Millisecond}
	q := NewQueued(slow)
	go q.Respond(context.Background(), Request{Name: "First"})
	time.Sleep(50 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if reply, _ := q.Respond(ctx, Request{Name: "Second"}); reply != nil {
		t.Error("等排隊逾時應該回 nil")
	}
	time.Sleep(400 * time.Millisecond)
	slow.mu.Lock()
	defer slow.mu.Unlock()
	if len(slow.intervals) != 1 {
		t.Errorf("模型被呼叫了 %d 次", len(slow.intervals))
	}
}
