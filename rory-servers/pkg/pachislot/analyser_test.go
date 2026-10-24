package pachislot

import (
	"context"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

func TestAnalyse(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	conf.AnalysisMode = true
	conf.CLogLevel = mlog.SystemLevel

	a := new(analyser)
	a.deadline = 1e3

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	a.Run(context.Background(), game)
}
