package fgsb

import (
	"fmt"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/fixtures"
)

func TestNewGame(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	fmt.Println(logic)
}
