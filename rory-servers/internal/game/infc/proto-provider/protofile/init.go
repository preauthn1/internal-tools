package protofile

// DefaultTransformers 預設所有轉換器
var DefaultTransformers = []Transformer{
	NewPatternTransformer(),
	NewFruitAttributeTransformer(),
	NewFruitEmitterTransformer(),
	NewStageGroupTransformer(),
}
