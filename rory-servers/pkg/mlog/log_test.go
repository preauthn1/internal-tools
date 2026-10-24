package mlog

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
)

func TestLogger(t *testing.T) {
	conf.CLogLevel = DebugLevel
	Debug("it is debug log test")
	Info("it is info log test")
	Error("it is error log test")
	System("it is sys log test")
	Warn("it is warn log test")

	NewDefaultLogger().Debug("it is debug log test")
	NewDefaultLogger().Info("it is info log test")
	NewDefaultLogger().Error("it is error log test")
	//NewDefaultLogger().Fatal("it is fatal log test")
	NewDefaultLogger().System("it is sys log test")
	NewDefaultLogger().Warn("it is warn log test")

}
