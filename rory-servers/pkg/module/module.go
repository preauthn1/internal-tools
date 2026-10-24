package module

// Module 模塊接口定義
//  組織一個功能項,提供基礎的功能組件
type Module interface {
	Run(ctx *MouduleCtx) error // 啟動模塊
	Stop()                     // 停止模塊
	Name() string              // 模塊名稱
}
