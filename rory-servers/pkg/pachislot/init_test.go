package pachislot

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

func TestMain(m *testing.M) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 1000 //40
	conf.TestMode = true
	conf.ClearLog = true

	// mlog.GetCLogger().Enable(false)

	m.Run()
}
