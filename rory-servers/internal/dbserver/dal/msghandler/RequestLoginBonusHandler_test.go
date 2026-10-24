package msghandler

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"testing"
)



func TestGetLoginBonusConfig(t *testing.T) {
	initDBConnection()

	handler := new(RequestLoginBonusHandler)
	ctx := new(module.MouduleCtx)
	rs := handler.getLoginBonusConfig(ctx, 100)
	mlog.System("======== %v", rs)
}