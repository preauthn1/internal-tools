package protofile

// DefaultTransformers 預設所有轉換器
var DefaultTransformers = []Transformer{
	NewGameProtoTransformer(),
	NewPayLineProtoTransformer(),
	NewPayTableProtoTransformer(),
	NewReelStripsProtoTransformer(),
	NewPatternProtoTransformer(),
	NewModuleProtoTransformer(),
	NewModeProtoTransformer(),
	NewRTDistributionProtoTransformer(),
	NewBRDistributionProtoTransformer(),
	NewNextModeDistributionProtoTransformer(),
	NewRTReleaseProtoTransformer(),
	NewResultDistributionProtoTransformer(),
}
