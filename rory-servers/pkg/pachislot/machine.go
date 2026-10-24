package pachislot

// MachineConfig 機台設定
type MachineConfig struct {
	BetOdds int64
	Pattern int32
}

// MachineDistributionManager 機台分布管理器
type MachineDistributionManager struct {
	Distribution  map[MachineConfig]int
	DefaultConfig MachineConfig
}

// NewMachineDistributionManager 建立機台分布管理器
func NewMachineDistributionManager(defBetOdds int64, defPattern int32) *MachineDistributionManager {
	m := new(MachineDistributionManager)
	m.Distribution = make(map[MachineConfig]int)
	m.DefaultConfig = MachineConfig{BetOdds: defBetOdds, Pattern: defPattern}
	return m
}

// RegisterDistribution 註冊一組設定
func (m *MachineDistributionManager) RegisterDistribution(BetOdds int64, Pattern int32, Count int) {
	cfg := MachineConfig{BetOdds: BetOdds, Pattern: Pattern}
	m.Distribution[cfg] = Count
}

// Pick 取一組設定
func (m *MachineDistributionManager) Pick() (BetOdds int64, Pattern int32) {
	for cfg, cnt := range m.Distribution {
		if cnt > 0 {
			m.Distribution[cfg]--
			return cfg.BetOdds, cfg.Pattern
		}
	}
	return m.DefaultConfig.BetOdds, m.DefaultConfig.Pattern
}
