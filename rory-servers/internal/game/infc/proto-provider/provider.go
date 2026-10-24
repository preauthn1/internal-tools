package protos

import "example.invalid/ServerTeam/Servers/internal/game/infc"

// Provider PahciSlot Proto 提供者
type Provider interface {
	Proto() (*infc.Proto, error) // 取得 Proto
}
