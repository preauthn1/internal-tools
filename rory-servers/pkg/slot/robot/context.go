package robot

// ReportContext 把所有紀錄整理至Context，提供 Print 資料
type ReportContext struct {
	// var printLineCount int = 4

	// basic
	TotalWinMoney    int64
	TotalBetMoney    int64
	MainGameWinMoney int64
	TotalTimes       int64
	MainGameTimes    int64

	// freespin
	FreeGameWinMoney         int64
	FreeGameExtAwardWinMoney int64
	FreeGameTimes            int64

	// bg
	BGWinMoney int64
	BGTimes    int64

	// jp
	JackpotWinMoney int64

	// lightning
	LightningWinMoney int64
	LightningTimes    int64

	// pick game
	PickGameWinMoney int64
	PickGameTimes    int64
}
