package msg

import (
	"fmt"
	"sync"

	"example.invalid/golang/protobuf/proto"
)

// InstanceProvider 提供解析 Message Data 的工具
var InstanceProvider *instanceProvider

type instanceProvider struct {
	store sync.Map
}

func (p *instanceProvider) key(BClassID, SClassID int32) string {
	return fmt.Sprintf("%d:%d", BClassID, SClassID)
}

func (p *instanceProvider) Register(BClassID, SClassID int32, fn func() proto.Message) {
	p.store.Store(p.key(BClassID, SClassID), fn)
}

func (p *instanceProvider) Load(BClassID, SClassID int32) func() proto.Message {
	val, ok := p.store.Load(p.key(BClassID, SClassID))
	if !ok {
		return nil
	}
	return val.(func() proto.Message)
}

func init() {
	InstanceProvider = new(instanceProvider)
}
