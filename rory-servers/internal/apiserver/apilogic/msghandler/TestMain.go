package msghandler

import (
	"testing"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
	"example.invalid/ServerTeam/Servers/pkg/module"
)

// GameLaunchBody ?
type GameLaunchBody struct {
	model.BodyBase
	Data model.GameLaunchViewModel `json:"data"`
}

// TestCtx ?
var TestCtx *module.MouduleCtx

// TestMain ?
func TestMain(m *testing.M) {
	m.Run()
}

// InitTest ?
func InitTest() {
	service.InitTestService()
}

// GetTestAccessToken ?
func GetTestAccessToken() string {
	_, token := service.GetTestAccessToken()
	return token
}
