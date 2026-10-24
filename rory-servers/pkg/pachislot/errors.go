package pachislot

import "errors"

// ErrProtoNotFound 查無 Proto
var ErrProtoNotFound = errors.New("Proto NotFound")

// ErrPatternNotFound 查無 Pattern
var ErrPatternNotFound = errors.New("Pattern NotFound")

// ErrModeConfigNotFound 查無 ModeConfig
var ErrModeConfigNotFound = errors.New("ModeConfig NotFound")

// ErrModuleConfigNotFound 查無 ModuleConfig
var ErrModuleConfigNotFound = errors.New("ModuleConfig NotFound")

// ErrNextModeDistributionNotFound 查無 NextModeDistribution
var ErrNextModeDistributionNotFound = errors.New("NextModeDistribution NotFound")

// ErrRTDistributionNotFound 查無 RTDistribution
var ErrRTDistributionNotFound = errors.New("RTDistribution NotFound")

// ErrBRDistributionNotFound 查無 BRDistribution
var ErrBRDistributionNotFound = errors.New("BRDistribution NotFound")

// ErrResultDistributionNotFound 查無 ResultDistribution
var ErrResultDistributionNotFound = errors.New("ResultDistribution NotFound")

// ErrPayLineNotFound 查無 PayLine
var ErrPayLineNotFound = errors.New("PayLine NotFound")

// ErrResultSetNotGeneratable 牌型無法生產
var ErrResultSetNotGeneratable = errors.New("ResultSet Not Generatable")
