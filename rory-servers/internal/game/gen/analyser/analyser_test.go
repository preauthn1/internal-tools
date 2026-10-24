package analyser

import (
	"os"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
)

func TestAnalyser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	conf.ClearLog = true
	conf.AnalysisMode = true

	proto := pachislot.DefaultProto

	Run(os.Stdout, cfg, proto, 600101, 1e5)
}
