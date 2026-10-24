package slot

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// PickGameDetail 翻牌紀錄
type PickGameDetail struct {
	Origin   *protodata.BonusTable `json:"-"`
	BonusID  int32
	ID       int32 // 順序 1-9
	Location int32 // 位置
	Money    int64 // 獎勵的錢
	Positon  Position
}

// PickGame 翻牌遊戲
type PickGame struct {
	Origin *protodata.SlotPayTable

	BonusTable []*protodata.BonusTable
	Size       int32 // 共有幾個選項

	GameBetMoney int64 // 下注金額
	GameBetRate  int64 // 下注倍率

	GameCount int32 // 遊戲次數
	WinMoney  int64 // 中獎的錢

	Process map[int32]*PickGameDetail
	All     []int64 // 所有翻牌紀錄，把所有12個選項都列出

	Positions []Position

	picker *prob.Random
}

// NewPickGame 建立翻牌遊戲
func NewPickGame(BonusTable []*protodata.BonusTable, size, count int32, betMoney, betRate int64) *PickGame {

	game := new(PickGame)
	game.picker = new(prob.Random)
	for _, b := range BonusTable {
		game.picker.Add(b, int(b.GetPerProb()))
	}

	game.BonusTable = BonusTable
	game.GameCount = count
	game.GameBetMoney = betMoney
	game.GameBetRate = betRate
	game.Process = game.generateProcess()
	game.Size = size
	game.All = game.generateLaftProcess(size)

	return game
}

// RandomBonus 決定中了哪個Bonus
func (game *PickGame) randomBonus() *protodata.BonusTable {

	t := game.picker.Choose()

	if t == nil {
		return nil
	}

	return t.(*protodata.BonusTable)
}

// GenerateProcess 產生翻牌歷程
func (game *PickGame) generateProcess() map[int32]*PickGameDetail {

	process := make(map[int32]*PickGameDetail)

	for i := int32(0); i < game.GameCount; i++ {
		detail := new(PickGameDetail)
		detail.ID = i + 1
		detail.Origin = game.randomBonus()
		detail.BonusID = detail.Origin.GetId()

		switch PayType(detail.Origin.GetRewardType()) {
		case PayTypeChip:
			f, _ := strconv.ParseFloat(detail.Origin.GetRewardValue(), 64)
			detail.Money = int64(f) * game.GameBetMoney
		case PayTypeChipXRate:
			f, _ := strconv.ParseFloat(detail.Origin.GetRewardValue(), 64)
			detail.Money = int64(f) * game.GameBetMoney * game.GameBetRate
		}

		game.WinMoney += detail.Money
		process[detail.ID] = detail

		game.All = append(game.All, detail.Money)
	}

	return process
}

// GenerateLaftProcess 產生所有牌局選項
func (game *PickGame) generateLaftProcess(size int32) []int64 {
	ans := make([]int64, 0, game.Size)

	// 計算剩餘數量
	for _, row := range game.Process {
		ans = append(ans, row.Money)
		size--
	}

	// 把剩餘的數量產出
	for i := 0; i < int(size); i++ {
		bonus := game.randomBonus()
		f, _ := strconv.ParseFloat(bonus.GetRewardValue(), 64)
		ans = append(ans, int64(f)*game.GameBetMoney)
	}

	return ans
}

// Choose 玩家選擇，更新Proccess
func (game *PickGame) Choose(id, location int32) (*PickGameDetail, error) {

	// 檢查是否結束
	if game.IsFinished() {
		return nil, fmt.Errorf("BonusGame is over")
	}

	// 檢查選擇次數是否正常
	if id <= 0 {
		return nil, fmt.Errorf("Choose ID should be [1-9]")
	}

	// 檢查前面次數是否有為選擇
	if id > 1 && game.Process[id-1].Location == 0 {
		return nil, fmt.Errorf("Choose ID: %d, Location: %d, But Record[%d] has not choosed", id, location, id-1)
	}

	if !game.Exist(id) {
		return nil, fmt.Errorf("Choose ID %d is not exist", id)
	}

	if game.Process[id].Location != 0 {
		return nil, fmt.Errorf("Choose ID %d has been choosed", id)
	}

	game.Process[id].Location = location

	if len(game.Positions) > 0 && len(game.Positions) > int(location) {
		game.Process[id].Positon = game.Positions[int(location-1)]
	}

	return game.Process[id], nil
}

// Exist 判斷該ID存不存在
func (game *PickGame) Exist(id int32) bool {
	if _, ok := game.Process[id]; ok {
		return true
	}
	return false
}

// IsFinished 判斷JackpotGame是否結束, 每個Process位置都不是0則結束
func (game *PickGame) IsFinished() bool {

	if game.Process == nil {
		return false
	}

	for _, detail := range game.Process {
		if detail.Location == 0 {
			return false
		}
	}

	return true
}

// GetProcessData 取得翻牌歷程紀錄
func (game *PickGame) GetProcessData() string {

	strs := make([]string, 0, len(game.Process))

	for i := 0; i < len(game.Process); i++ {
		if detail, ok := game.Process[int32(i+1)]; ok {
			strs = append(strs, fmt.Sprintf("%d", detail.Money))
		}
	}

	return strings.Join(strs, ",")
}
