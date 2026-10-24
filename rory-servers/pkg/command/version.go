package command

import (
	"example.invalid/ServerTeam/Servers/pkg/conf"
)

// VersionCommand 顯示版本指令
type VersionCommand struct {
}

// Name 指令名稱
func (c *VersionCommand) Name() string {
	return "version"
}

// Help 指令說明
func (c *VersionCommand) Help() string {
	return "顥示目前版本資訊"
}

// Exec 執行指令
func (c *VersionCommand) Exec(args []string) string {
	return conf.Version()
}
