package slot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// JackpotGameDetail 翻牌紀錄
type JackpotGameDetail struct {
	ID    int32 // 順序 1-9
	Place int32 // 位置
	Icon  int32 // 符號(Level) 1: Grand, 2: Major, 3: Minor, 4: Mini
}

// JackpotGame 翻牌遊戲
type JackpotGame struct {
	Origin *protodata.SlotPayTable

	Jackpots []*Jackpot // 當前所有獎池

	MaxBetRate  int32 // 最大黃金倍率
	MaxBetMoney int32 // 最大押注數量

	GameBetRate  int32 // 下注符號倍率
	GameBetMoney int64 // 下注金額

	WinJackpot *Jackpot // 中獎的Jackpot
	WinMoney   int64    // 中獎的錢

	Process map[int32]*JackpotGameDetail
}

// NewJackpotGame 建立翻牌遊戲
func NewJackpotGame(jackpots []*Jackpot, maxRate, maxMoney, rate int32, money int64) *JackpotGame {
	return &JackpotGame{
		Jackpots:     jackpots,
		MaxBetRate:   maxRate,
		MaxBetMoney:  maxMoney,
		GameBetRate:  rate,
		GameBetMoney: money,
	}
}

// RandomJackpot 決定中了哪個Jackpot
func (game *JackpotGame) RandomJackpot() *Jackpot {

	r := new(prob.Random)
	for _, jp := range game.Jackpots {
		r.Add(jp, int(jp.Origin.GetPerProb()))
	}

	t := r.Choose()

	r.PrintOptions()

	if t != nil {
		return t.(*Jackpot)
	}

	return nil
}

// GenerateProcess 產生翻牌歷程
func (game *JackpotGame) GenerateProcess() map[int32]*JackpotGameDetail {

	process := make(map[int32]*JackpotGameDetail)

	rand.Seed(time.Now().UnixNano())

	if game.WinJackpot == nil {
		return process
	}

	// 檢視所有Jackpot, 知道要開哪些符號
	levels := make([]JackpotLevel, 0, len(game.Jackpots))

	for _, jackpot := range game.Jackpots {
		levels = append(levels, jackpot.Level)
	}
	levels = append(levels, levels...)

	// 洗亂
	for i, l := 0, len(levels); i < l; i++ {
		j := rand.Intn(l - 1)
		levels[i], levels[j] = levels[j], levels[i]
	}

	length := rand.Intn(6) + 3
	procLevels := make([]JackpotLevel, 0, length)
	for i := 0; i < length; i++ {
		procLevels = append(procLevels, levels[i])
	}

	// 最後補中獎的level
	cnt := 0
	for _, level := range procLevels {
		if level == game.WinJackpot.Level {
			cnt++
		}
	}

	for i := 3; i > cnt; i-- {
		procLevels = append(procLevels, game.WinJackpot.Level)
	}

	// 組合
	for idx, level := range procLevels {
		detail := &JackpotGameDetail{
			ID:    int32(idx) + 1,
			Place: 0,
			Icon:  level.Int(),
		}

		process[detail.ID] = detail
	}

	return process
}

// Play 直接選出一個Jackpot
func (game *JackpotGame) Play() {

	game.WinJackpot = game.RandomJackpot()
	game.Process = game.GenerateProcess()

	currentBet := int64(game.GameBetRate) * game.GameBetMoney
	maxBet := int64(game.MaxBetRate) * int64(game.MaxBetMoney)

	game.WinMoney = game.WinJackpot.GetCurrentPoolByBet(currentBet, maxBet)
}

// Choose 玩家選擇，更新Proccess
func (game *JackpotGame) Choose(id, place, icon int32) error {

	// 檢查前個位置有被更新
	if detail, ok := game.Process[id-1]; ok {
		if detail.Place == 0 {
			return fmt.Errorf("JackpotGame 階段 %d 尚未選擇", id)
		}
	}

	if detail, ok := game.Process[id]; ok {
		if detail.Place == 0 {
			detail.Place = place
			return nil
		}
		return fmt.Errorf("JackpotGame 階段 %d 已選擇位子 %d", id, detail.Place)
	}

	return fmt.Errorf("JackpotGame 階段 %d 不存在", id)
}

// IsFinished 判斷JackpotGame是否結束, 每個Process位置都不是0則結束
func (game *JackpotGame) IsFinished() bool {

	if game.Process == nil {
		return false
	}

	for _, detail := range game.Process {
		if detail.Place == 0 {
			return false
		}
	}

	return true
}

// GetProcessData 取得翻牌歷程紀錄
func (game *JackpotGame) GetProcessData() string {

	strs := make([]string, 0, len(game.Process))

	for i := 0; i < len(game.Process); i++ {
		if detail, ok := game.Process[int32(i+1)]; ok {
			strs = append(strs, strconv.Itoa(int(detail.Icon)))
		}
	}

	return strings.Join(strs, ",")
}
