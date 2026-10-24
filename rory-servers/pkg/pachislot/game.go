package pachislot

import (
	"fmt"
	"math/rand"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/prob"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// Game 遊戲主流程
//  每次決定一組大獎
//  大獎結束後，結束主流程
type Game struct {
	// 固定
	Pattern     *pb.PatternOption // 當前 Pattern
	CostPerCoin int64             // 一枚硬幣花費(這邊用不到，拿來借給呼叫PachiGameFunc的人知道)
	RoadMap     *RoadMap          // 轉數表
	Proto       *Proto

	// 當前資訊
	Scenario  pb.Scenario // 當前情境
	Bonus     pb.Scenario // 大獎類型 (BB/RB)
	Mode      pb.ModeType // 遊戲模式 (A/B/天國)
	RT        int         // 剩餘RT數
	Replay    bool        // 是否為重玩
	BonusSure bool        // Bonus確定

	// 特色遊戲
	JACTimes       int // JAC次數
	BonusTimesLeft int // Bonus剩餘Spin次數
	JACTimesLeft   int // JAC剩餘Spin次數

	// 指定資訊
	NextMode   pb.ModeType   // 指定下次遊戲模式
	BonusStore int           // Bonus庫存
	NextResult pb.ResultType // 指定牌型

	// 設定值 from Pattern
	ModeConfigs     map[pb.ModeType]*pb.Mode              // 可用的Mode設定
	ModuleConfigs   map[pb.Scenario]*pb.Module            // 可用的遊戲模組
	RTReleseConfigs map[pb.ResultType]*pb.RTReleaseOption // 解除RT規則設定

	beforeInitFunc   []func(*Game)                      // 初始化前處理事項
	afterInitFunc    []func(*Game, bool)                // 初始化後處理事項
	beforeSpinFunc   []func(*Game)                      // Spin前處理事項
	afterSpinFunc    []func(*Game, pb.ResultType, bool) // Spin後處理事項
	beforePayoutFunc []func(*Game)                      // Payout前處理事項
	afterPayoutFunc  []func(*Game, []Payout)            // Payout後處理事項
	beforeRewardFunc []func(*Game)                      // Reward前處理事項
	afterRewardFunc  []func(*Game, []Payout)            // Reward後處理事項
}

func (game *Game) String() string {
	var cnt int

	switch game.Scenario {
	default:
		cnt = -1
	case pb.Scenario_BB, pb.Scenario_RB:
		cnt = game.BonusTimesLeft
	case pb.Scenario_JAC:
		cnt = game.JACTimesLeft
	}
	return fmt.Sprintf("Pattern[%d] 當前[%v(%d)][%d] Mode[%v][BR:%d][RT:%d][NM:%d] 大獎[%v] RT[%d] 1G[%d] Bonus確[%t] Replay[%t]",
		game.Pattern.ID,
		game.Scenario, cnt, game.ModuleConfigs[game.Scenario].ResultProbGroupID,
		game.Mode, game.ModeConfigs[game.Mode].BRGroupID, game.ModeConfigs[game.Mode].RTGroupID, game.ModeConfigs[game.Mode].NextModeGroupID,
		game.Bonus, game.RT, game.BonusStore, game.BonusSure, game.Replay)
}

// NewGame 建立遊戲
func NewGame(proto *Proto, pid int32) (*Game, error) {
	p, err := FindPatternByID(proto, pid)
	if err != nil {
		return nil, fmt.Errorf("find pattern id(%d) error, %v", pid, err)
	}
	modes, err := FindModesByGroupID(proto, p.ModeGroupID)
	if err != nil {
		return nil, fmt.Errorf("find mode gid(%d) error, %v", p.ModeGroupID, err)
	}
	modules, err := FindModulesByGroupID(proto, p.ModuleGroupID)
	if err != nil {
		return nil, fmt.Errorf("find module gid(%d) error, %v", p.ModuleGroupID, err)
	}
	rtrelease, err := FindRTReleaseDistributionByGroupID(proto, p.RTReleaseGroupID)
	if err != nil {
		return nil, fmt.Errorf("find RTRelease gid(%d) error, %v", p.RTReleaseGroupID, err)
	}

	game := new(Game)
	game.Proto = proto
	game.Pattern = p
	game.RoadMap = NewRoadMap(DefaultRoadMapBuf)
	game.Scenario = DefaultModule
	game.Mode = DefaultMode
	game.RT = DefaultRT
	game.Bonus = DefaultBonus
	game.NextMode = pb.ModeType_Unassigned

	// Init Mode Config
	game.ModeConfigs = make(map[pb.ModeType]*pb.Mode)
	for _, m := range modes {
		game.ModeConfigs[m.ModeType] = m
	}

	// Init Module Config
	game.ModuleConfigs = make(map[pb.Scenario]*pb.Module)
	for _, m := range modules {
		game.ModuleConfigs[m.Type] = m
	}
	// 把 NG 製作上去
	game.ModuleConfigs[pb.Scenario_NG] = &pb.Module{ResultProbGroupID: p.NGResultGroupID, Count: -1, Cost: 3}

	// Init RT Release Config
	game.RTReleseConfigs = make(map[pb.ResultType]*pb.RTReleaseOption)
	for _, opt := range rtrelease {
		game.RTReleseConfigs[opt.ResultType] = opt
	}

	game.Init()

	return game, nil
}

// AddBeforeInitFunc 增加初始化前處理事項
func (game *Game) AddBeforeInitFunc(f func(*Game)) {
	game.beforeInitFunc = append(game.beforeInitFunc, f)
}

// AddAfterInitFunc 增加初始化後處理事項
func (game *Game) AddAfterInitFunc(f func(*Game, bool)) {
	game.afterInitFunc = append(game.afterInitFunc, f)
}

// AddBeforeSpinFunc 增加Spin前處理事項
func (game *Game) AddBeforeSpinFunc(f func(*Game)) {
	game.beforeSpinFunc = append(game.beforeSpinFunc, f)
}

// AddAfterSpinFunc 增加Spin後處理事項
func (game *Game) AddAfterSpinFunc(f func(*Game, pb.ResultType, bool)) {
	game.afterSpinFunc = append(game.afterSpinFunc, f)
}

// AddBeforePayoutFunc 增加Payout前處理事項
func (game *Game) AddBeforePayoutFunc(f func(*Game)) {
	game.beforePayoutFunc = append(game.beforePayoutFunc, f)
}

// AddAfterPayoutFunc 增加Payout後處理事項
func (game *Game) AddAfterPayoutFunc(f func(*Game, []Payout)) {
	game.afterPayoutFunc = append(game.afterPayoutFunc, f)
}

// AddBeforeRewardFunc 增加Reward前處理事項
func (game *Game) AddBeforeRewardFunc(f func(*Game)) {
	game.beforeRewardFunc = append(game.beforeRewardFunc, f)
}

// AddAfterRewardFunc 增加Reward後處理事項
func (game *Game) AddAfterRewardFunc(f func(*Game, []Payout)) {
	game.afterRewardFunc = append(game.afterRewardFunc, f)
}

// Init 依照模式初始化
//  * 下次模式
//  * RT次數
//  * 大獎
//  * 有Bonus庫存則解開
func (game *Game) Init() {

	var releaseBonusStore bool

	for _, f := range game.beforeInitFunc {
		f(game)
	}
	defer func() {
		for _, f := range game.afterInitFunc {
			f(game, releaseBonusStore)
		}
	}()

	// 當前資訊
	game.Scenario = pb.Scenario_NG
	game.RT = DefaultRT
	game.Replay = false
	game.BonusSure = false

	// 特色遊戲
	game.BonusTimesLeft = 0
	game.JACTimesLeft = 0
	game.JACTimes = 0

	// 有被指定直接使用指定模式
	if game.NextMode != pb.ModeType_Unassigned {
		example.invalid("指定模式[%v]", game.NextMode)
		game.Mode = game.NextMode
		game.NextMode = pb.ModeType_Unassigned
	}

	m := game.ModeConfigs[game.Mode]

	{ // 決定下次模式
		options, err := FindNextModeDistributionByGroupID(game.Proto, m.NextModeGroupID)
		if err != nil { // 查無下次可使用的模式
			mlog.Warn("Find Next Mode Distribution (%d) Error, %v", m.NextModeGroupID, err)
			game.NextMode = DefaultMode
		} else {
			r := new(prob.Random)
			for _, opt := range options {
				r.Add(opt.Mode, int(opt.Weight))
			}
			game.NextMode = r.Choose().(pb.ModeType)
		}
	}

	{ // 決定RT，從表中骰一組設定，再從最小值到最大值之間隨機一個數值
		options, err := FindRTDistributionByGroupID(game.Proto, m.RTGroupID)
		if err != nil { // 查無RT分配表
			mlog.Warn("Find RT Distribution (%d) Error, %v", m.RTGroupID, err)
			game.RT = DefaultRT
		} else {
			r := new(prob.Random)
			for _, opt := range options {
				r.Add(opt, int(opt.Weight))
			}
			option := r.Choose().(*pb.RTDistributionOption)

			if option.Max < option.Min {
				game.RT = DefaultRT
			} else {
				game.RT = int(rand.Int31n(option.Max+1-option.Min) + option.Min)
			}
		}
	}

	{ // 決定大獎
		options, err := FindBRDistributionByGroupID(game.Proto, m.BRGroupID)
		if err != nil { // 查無BR分配表
			mlog.Warn("Find BR Distribution (%d) Error, %v", m.BRGroupID, err)
			game.Bonus = DefaultBonus
		} else {
			r := new(prob.Random)
			for _, opt := range options {
				r.Add(opt.Type, int(opt.Weight))
			}
			game.Bonus = r.Choose().(pb.Scenario)
		}
	}

	{ // 檢查Bonus庫存
		if game.BonusStore > 0 {
			example.invalid("有Bonus庫存[%d]次", game.BonusStore)
			game.BonusStore--
			game.BonusSure = true
			releaseBonusStore = true
		}
	}
}

// BetCoin 取得Spin所需代幣數
func (game *Game) BetCoin() int64 {
	c := game.ModuleConfigs[game.Scenario]
	if game.Replay {
		return 0
	}
	return int64(c.Cost)
}

// Spin 拉霸
//  @return cost 花費倍數
//  @return result 牌型
//  @return rtrelesd 是否解除RT
func (game *Game) Spin() (cost int64, result pb.ResultType, rtrelesd bool) {

	for _, f := range game.beforeSpinFunc {
		f(game)
	}
	defer func() {
		for _, f := range game.afterSpinFunc {
			f(game, result, rtrelesd)
		}
	}()

	defer func() {
		game.Replay = false
	}()

	c := game.ModuleConfigs[game.Scenario]
	if game.Replay {
		cost = 0
	} else {
		cost = int64(c.Cost)
	}
	if game.NextResult != pb.ResultType_Empty {
		result = game.NextResult
		game.NextResult = pb.ResultType_Empty
	} else {
		result = RandResult(game.Proto, c.ResultProbGroupID)
	}

	switch game.Scenario {
	case pb.Scenario_NG:
		game.RoadMap.Inc()

		// Bonus確定時，50%機率骰出確定進入Bonus，牌型由對應的BonusID查找，骰一個
		if game.BonusSure {
			r := rand.Intn(100)
			if r < DefaultBonusSureWinRate {
				// 找出可完成這項大獎的Payout
				var payouts []*pb.PayoutOption
				for _, p := range game.Proto.PayTable {
					if p.UseType != pb.Scenario_NG {
						continue
					}
					if p.PayType != pb.PayType_PTBonus {
						continue
					}
					m, err := FindModuleByID(game.Proto, p.Payout)
					if err != nil {
						continue
					}
					if m.Type != game.Bonus {
						continue
					}
					payouts = append(payouts, p)
				}

				if len(payouts) == 0 {
					mlog.Error("Bonus確定，即將開出大獎[%s]，但是查無可開出大獎的組合，隨機開一個組合[%s]", game.Bonus, result)
					return
				}

				result = payouts[rand.Intn(len(payouts))].ResultType
				game.BonusSure = false
				return
			}
			return
		}

		// 達到RT數
		if game.RT <= 0 {
			example.invalid("到達RT次數，執行Bonus確定")
			game.BonusSure = true
			result = pb.ResultType_BonusSure
			return
		}

		game.RT--

		// 檢查 RT 是否須更新
		if ok, cnt := game.chkRTRelease(result); ok {
			rtrelesd = true
			game.RT = cnt
			example.invalid("RT解除，更新次數[%d]", game.RT)
		}

	case pb.Scenario_BB, pb.Scenario_RB:
		game.BonusTimesLeft--
	case pb.Scenario_JAC:
		game.JACTimesLeft--
	}
	return
}

// Payout 取得
func (game *Game) Payout(t pb.ResultType) (payouts []Payout) {
	for _, f := range game.beforePayoutFunc {
		f(game)
	}
	defer func() {
		for _, f := range game.afterPayoutFunc {
			f(game, payouts)
		}
	}()

	for _, opt := range game.Proto.PayTable {
		if opt.UseType != pb.Scenario_ALL && opt.UseType != game.Scenario {
			continue
		}
		if opt.ResultType != t {
			continue
		}
		switch opt.PayType {
		case pb.PayType_PTChip:
			payouts = append(payouts, NewPayoutChip(int64(opt.Payout)))
		case pb.PayType_PTReplay:
			payouts = append(payouts, NewPayoutReplay())
		case pb.PayType_PTBonus:
			m, err := FindModuleByID(game.Proto, opt.Payout)
			if err != nil {
				mlog.Error("建立Bonus錯誤，查詢 Module %d 錯誤, %v", opt.Payout, err)
				continue
			}
			s := pb.Scenario(m.Type)
			switch s {
			default:
				mlog.Error("建立Bonus錯誤，無法辨別Payout值 %d", opt.Payout)
			case pb.Scenario_BB, pb.Scenario_RB, pb.Scenario_JAC:
				payouts = append(payouts, NewPayoutBonus(s))
			}
		case pb.PayType_PT1G:
			payouts = append(payouts, NewPayout1G(int(opt.Payout)))
		}
	}
	return
}

// Reward 給獎，回傳實際出獎的獎勵
func (game *Game) Reward(t pb.ResultType) (payouts []Payout) {
	for _, f := range game.beforeRewardFunc {
		f(game)
	}
	defer func() {
		for _, f := range game.afterRewardFunc {
			f(game, payouts)
		}
	}()

	for _, p := range game.Payout(t) {
		_, ok := p.Reward(game)
		if ok {
			payouts = append(payouts, p)
		}
	}

	if game.chkFinished() {
		game.Init()
	}

	return
}

func (game *Game) chkFinished() bool {
	switch game.Scenario {
	default:
		return false
	case pb.Scenario_BB, pb.Scenario_RB:
		return game.BonusTimesLeft == 0
	case pb.Scenario_JAC:
		if game.JACTimesLeft == 0 {
			game.Scenario = game.Bonus
			return game.BonusTimesLeft == 0
		}
		return false
	}
}

// chkRTRelease 解除RT限制，依照結果判斷是否解除RT限制，解除將把RT次數更新
//  @param r 盤面結果
//  @return 是否解除, 新RT數
func (game *Game) chkRTRelease(t pb.ResultType) (bool, int) {
	if game.RT <= DefaultMinRTLockTime {
		return false, 0
	}

	// 如果是 ReleaseRT 依照現在大獎類型替換
	if t == pb.ResultType_ReleaseRT {
		switch game.Bonus {
		case pb.Scenario_BB:
			t = pb.ResultType_ReleaseRTonBB
		case pb.Scenario_RB:
			t = pb.ResultType_ReleaseRTonRB
		}
	}

	rc, ok := game.RTReleseConfigs[t]

	// 此牌型無設定
	if !ok {
		return false, 0
	}

	// 情境不適用
	if rc.UseType != pb.Scenario_ALL && game.Scenario != rc.UseType {
		return false, 0
	}

	// 次數設定為空
	if len(rc.Times) == 0 {
		return false, 0
	}

	// 決定是否解除
	r := rand.Int31n(rc.Denominator)
	if r >= rc.Rate { // 沒中，不解除
		return false, 0
	}

	return true, int(rc.Times[rand.Intn(len(rc.Times))])
}

// RandResult 隨機抽取一個牌組結果，若查無設定會回傳一般槓龜
//  @param gid 群組編號
//  @return 牌組結果
func RandResult(proto *Proto, gid int32) pb.ResultType {
	options, err := FindResultDistributionByGroupID(proto, gid)
	if err != nil {
		mlog.Warn("Find Result Distribution (%d) Error, %v", gid, err)
		return pb.ResultType_None
	}

	r := new(prob.Random)

	for _, opt := range options {
		r.Add(opt.Result, int(opt.Weight))
	}

	return r.Choose().(pb.ResultType)
}
