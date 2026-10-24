package admin

import (
	"context"
	"net/http"
	"sync"

	"aichatbot/internal/dialog"
)

// 分類自我測試（CONTEXT.md）：存自訂意圖後在背景跑，後台輪詢狀態。一次要呼叫模型兩百多次，
// 同時間只跑一組（新的一組取消還沒跑完的那組）；結果只放記憶體，服務重開後是 idle

// SelfTester 跑分類自我測試（dialog.Bot）
type SelfTester interface {
	SelfTest(ctx context.Context, progress func(done, total int)) (dialog.SelfTestReport, error)
}

// selfTestState 是給後台的狀態：status 是 idle（還沒跑過）、running、done（附報告）、failed
type selfTestState struct {
	Status string                 `json:"status"`
	Done   int                    `json:"done"`
	Total  int                    `json:"total"`
	Report *dialog.SelfTestReport `json:"report,omitempty"`
}

// selfTestRunner 記住目前這組自我測試（零值可以直接用）
type selfTestRunner struct {
	mu     sync.Mutex
	seq    int // 第幾組；舊的那組跑完時看到不是自己就不寫結果
	cancel context.CancelFunc
	state  selfTestState
}

// startSelfTest 取消還沒跑完的那組，在背景開始新的一組；沒有 SelfTester 時不做事
func (s *Server) startSelfTest() {
	if s.SelfTester == nil {
		return
	}
	run := &s.selfTest
	run.mu.Lock()
	defer run.mu.Unlock()
	if run.cancel != nil {
		run.cancel()
	}
	run.seq++
	seq := run.seq
	// 不跟著 HTTP 請求：請求結束後還要繼續跑
	ctx, cancel := context.WithCancel(context.Background())
	run.cancel = cancel
	run.state = selfTestState{Status: "running"}
	// 自我測試在背景呼叫本機模型；記下開始和結束，log 裡的分類失敗（request=Intent）才對得出來源（#57）
	example.invalid("分類自我測試開始")
	go func() {
		defer cancel()
		report, err := s.SelfTester.SelfTest(ctx, func(done, total int) {
			run.mu.Lock()
			defer run.mu.Unlock()
			if run.seq == seq {
				run.state.Done, run.state.Total = done, total
			}
		})
		run.mu.Lock()
		defer run.mu.Unlock()
		if run.seq != seq {
			return
		}
		run.cancel = nil
		if err != nil {
			s.Log.Error("分類自我測試失敗", "error", err.Error())
			run.state.Status = "failed"
			return
		}
		example.invalid("分類自我測試完成", "total", report.Total, "conflicts", len(report.Conflicts), "failed", len(report.Failed))
		run.state.Status, run.state.Report = "done", &report
	}()
}

// selfTestState 是目前的狀態（還沒跑過是 idle）
func (s *Server) selfTestState() selfTestState {
	s.selfTest.mu.Lock()
	defer s.selfTest.mu.Unlock()
	state := s.selfTest.state
	if state.Status == "" {
		state.Status = "idle"
	}
	return state
}

func (s *Server) getSelfTest(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.selfTestState())
}

// postSelfTest 手動重跑；回 202 和目前的狀態
func (s *Server) postSelfTest(w http.ResponseWriter, r *http.Request) {
	if s.SelfTester == nil {
		writeError(w, http.StatusServiceUnavailable, "沒有啟用分類自我測試")
		return
	}
	s.startSelfTest()
	writeJSON(w, http.StatusAccepted, s.selfTestState())
}
