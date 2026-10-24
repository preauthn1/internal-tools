package model

import "context"

// Queued 讓呼叫一次只跑一個（本機模型同時跑兩個，品質會明顯變差）。只包 Apple 本機模型，Haiku 可以同時呼叫。
// 等排隊的時候 ctx 結束就放棄（回 nil）。模型服務那邊還有 Swift 版的檔案鎖，和 Swift 的量測程式一起排隊
type Queued struct {
	base ChatModel
	slot chan struct{}
}

func NewQueued(base ChatModel) *Queued {
	return &Queued{base: base, slot: make(chan struct{}, 1)}
}

func (q *Queued) Respond(ctx context.Context, r Request) (*Reply, Usage) {
	select {
	case q.slot <- struct{}{}:
	case <-ctx.Done():
		return nil, Usage{}
	}
	defer func() { <-q.slot }()
	return q.base.Respond(ctx, r)
}
