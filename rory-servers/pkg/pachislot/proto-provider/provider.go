package protos

import (
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
)

// Provider PahciSlot Proto 提供者
type Provider interface {
	Proto() (*pachislot.Proto, error) // 取得 Proto
}
