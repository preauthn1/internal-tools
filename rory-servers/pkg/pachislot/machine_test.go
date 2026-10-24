package pachislot

import "testing"

func TestMachineDistributionManager(t *testing.T) {

	var defBetOdds int64 = 1000
	var defPattern int32 = 600101

	mdm := NewMachineDistributionManager(defBetOdds, defPattern)

	if mdm.DefaultConfig.BetOdds != defBetOdds {
		t.Errorf("DefaultConfig.BetOdds want %d, but %d", defBetOdds, mdm.DefaultConfig.BetOdds)
	}

	if mdm.DefaultConfig.Pattern != defPattern {
		t.Errorf("DefaultConfig.Pattern want %d, but %d", defPattern, mdm.DefaultConfig.Pattern)
	}

	t.Run("Pick with Default", func(t *testing.T) {
		odds, pattern := mdm.Pick()
		if odds != defBetOdds {
			t.Errorf("Pick() odds want %d, but %d", defBetOdds, odds)
		}
		if pattern != defPattern {
			t.Errorf("Pick() pattern want %d, but %d", defPattern, pattern)
		}
	})

	t.Run("Register", func(t *testing.T) {
		mdm.RegisterDistribution(2000, 600102, 100)
		mc := MachineConfig{BetOdds: 2000, Pattern: 600102}
		if mdm.Distribution[mc] != 100 {
			t.Errorf("Distribution[%v] want %d, but %d", mc, 100, mdm.Distribution[mc])
		}
	})

	t.Run("Pick to End", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			odds, pattern := mdm.Pick()
			if odds != 2000 {
				t.Errorf("Pick() odds want %d, but %d", 2000, odds)
			}
			if pattern != 600102 {
				t.Errorf("Pick() pattern want %d, but %d", 600102, pattern)
			}
		}
		odds, pattern := mdm.Pick()
		if odds != defBetOdds {
			t.Errorf("Pick() odds want %d, but %d", defBetOdds, odds)
		}
		if pattern != defPattern {
			t.Errorf("Pick() pattern want %d, but %d", defPattern, pattern)
		}
	})
}
