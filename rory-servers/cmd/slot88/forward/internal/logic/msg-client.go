package logic

import (
	"encoding/json"
	"sort"
	"time"

	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/cmd/slot88/forward/internal/network"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

func (logic *Logic) gameresult() (interface{}, func()) {

	if logic.StepResult != nil && logic.leftshow > 0 {

		idx := len(logic.StepResult.SlotResult) - logic.leftshow

		logic.GameResult = logic.transerStepResultToGameResult(logic.StepResult, idx)
		logic.leftshow--

		if logic.leftshow == 0 {
			logic.StepResult = nil
		}

		defer logic.clientUpdateMoney(logic.GameResult.GetFinalMoney(), true)
	}

	if logic.GameResult == nil {
		logic.GameResult = logic.genFackResult()
	}

	// Bet Rate
	var betrate int32
	if logic.Setting != nil {
		betrate = logic.Setting.GetBetRates()[0]
	}
	if logic.lastspin != nil {
		betrate = logic.lastspin.GetRate()
	}

	var betcount float64
	// Bet Money
	if logic.lastspin != nil {
		betcount = logic.MoneyTransfer.S2C(float64(logic.lastspin.GetMoney()))
	} else {
		betcount = logic.MoneyTransfer.S2C(float64(logic.Setting.BetMoneys[0]))
	}

	// Win Money
	var winmoney float64
	if logic.GameResult != nil {
		if logic.GameResult.GetFreeGameWinMoney() > 0 {
			winmoney = logic.MoneyTransfer.S2C(float64(logic.GameResult.GetFreeGameWinMoney()))
		} else {
			winmoney = logic.MoneyTransfer.S2C(float64(logic.GameResult.GetWinMoney()))
		}
	}

	// Free Spin Count
	var selectionidx interface{} = 1
	var freespincount int32 = -1
	if logic.LeftFreeGameData != nil {
		if logic.ServeState == ServeStateLogined {
			freespincount = logic.LeftFreeGameData.GetCount()
		} else if logic.incfree > 0 && int(logic.LeftFreeGameData.GetCount()) == logic.incfree && !logic.freedatasend {
			freespincount = -1
		} else {
			freespincount = logic.LeftFreeGameData.GetCount()
		}

		betcount = logic.MoneyTransfer.S2C(float64(logic.LeftFreeGameData.GetBetMoney()))
		betrate = int32(logic.LeftFreeGameData.GetBetRate())

		// 選擇遊戲紀錄
		{
			for i, row := range logic.Setting.FreeGames {
				if row.GetID() == logic.LeftFreeGameData.GetID() {
					selectionidx = []interface{}{i}
				}
			}
		}
	}

	// 整理盤面
	var maxReel, maxRow int32
	rows := make([]interface{}, 0)
	if logic.GameResult != nil {
		_result := make(map[int32]map[int32]string)
		for _, row := range logic.GameResult.ReelResult {
			if maxReel < row.Position.GetReel() {
				maxReel = row.Position.GetReel()
			}
			if maxRow < row.Position.GetRow() {
				maxRow = row.Position.GetRow()
			}

			if _, ok := _result[row.Position.GetRow()]; !ok {
				_result[row.Position.GetRow()] = make(map[int32]string)
				_result[row.Position.GetRow()][row.Position.GetReel()] = row.Symbol.GetSymbolName()
			} else {
				_result[row.Position.GetRow()][row.Position.GetReel()] = row.Symbol.GetSymbolName()
			}
		}

		for i := int32(1); i <= maxRow; i++ {
			var row []string
			for j := int32(1); j <= maxReel; j++ {
				row = append(row, _result[i][j])
			}
			rows = append(rows, row)
		}
	}

	// 整理賠線
	lines := make([]interface{}, 0)
	if logic.GameResult != nil {
		for _, row := range logic.GameResult.PayGroups {

			if row.GetWinMoney() == 0 {
				continue
			}

			// 找這位子的SymbolName
			var symbol string
			// pos := row.BlinkPositions[0]
			// for _, result := range logic.GameResult.ReelResult {
			// 	if reflect.DeepEqual(result.Position, pos) {
			// 		symbol = result.Symbol.GetSymbolName()
			// 		if strings.HasPrefix(symbol, "w") {
			// 			continue
			// 		}
			// 		break
			// 	}
			// }
			symbol = row.GetSymbolName()

			var lineid int32

			if row.PayLines.GetPayLineID() > 0 {
				lineid = row.PayLines.GetPayLineID()
			}

			var vect int32 = 1 // 方向
			switch slot.LinkType(row.GetLinkType()) {
			case slot.LinkTypeLeft:
				vect = 1
			case slot.LinkTypeRight:
				vect = -1
			}

			var cnt int32 = row.GetCount()
			if len(row.BlinkPositions) > 0 {
				rMap := make(map[int32]int)
				for _, pos := range row.BlinkPositions {
					rMap[pos.GetReel()]++
				}
				cnt = int32(len(rMap))
			}

			lines = append(lines, []interface{}{
				lineid,
				symbol,
				cnt,
				logic.MoneyTransfer.S2C(float64(row.GetWinMoney())),
				vect, // 方向
			})
		}
	}

	data := []interface{}{
		[]interface{}{1, 0, 1, freespincount, betrate, betcount, winmoney, logic.lastWin, betcount, 0, selectionidx},
		[]interface{}{
			rows,
			[]interface{}{0, 0, 0, 0, 0},
			[]interface{}{-1, 0},
		},
		lines,
		[]interface{}{1},
	}

	var callback = func() {
		if logic.LeftFreeGameData != nil && logic.LeftFreeGameData.GetCount() == 0 {
			// {"actionId":"2","eventData":[1],"id":"2"}

			_m := new(network.Message)
			_m.ID = 2
			_m.ActionID = 2
			_m.EventData = []interface{}{1}

			logic.cw.WriteClient(_m)

			logic.LeftFreeGameData = nil
			logic.incfree = 0
			logic.freedatasend = false
		} else if logic.incfree > 0 && !logic.freedatasend {
			// {"actionId":"2","eventData":[1,[5,1]],"id":"2"}
			_m := new(network.Message)
			_m.ID = 2
			_m.ActionID = 2
			_m.EventData = []interface{}{
				1,
				[]interface{}{logic.incfree, 1},
			}

			logic.cw.WriteClient(_m)
			logic.freedatasend = true
		}

		// PickGame
		if logic.PickGameRecord != nil && logic.GameState == slot.GameStateUserAction {
			// {"actionId":"2","eventData":[4,[1]],"id":"2"}
			_m := new(network.Message)
			_m.ID = 2
			_m.ActionID = 2
			_m.EventData = []interface{}{
				4,
				[]interface{}{logic.PickGameRecord.GetCount()},
			}
			logic.cw.WriteClient(_m)
		}

		// JackpotGame
		if logic.JackpotGameRecord != nil && logic.GameState == slot.GameStateUserAction {

			if logic.Round > 0 { // send event
				_m := new(network.Message)
				_m.ID = 2
				_m.ActionID = 2
				_m.EventData = []interface{}{7}

				logic.cw.WriteClient(_m)
			}

			logic.clientJackpotGameMode()
		}

		// Selection
		if logic.SelectGameOptions != nil && logic.GameState == slot.GameStateUserAction {
			_m := new(network.Message)
			_m.ID = 2
			_m.ActionID = 2
			_m.EventData = []interface{}{2}

			logic.cw.WriteClient(_m)
		}

	}

	if logic.GameResult != nil && logic.LeftFreeGameData == nil && logic.GameResult.GetWinMoney() > 0 {
		logic.lastWin = logic.MoneyTransfer.S2C(float64(logic.GameResult.GetWinMoney()))
	}

	return data, callback
}

func (logic *Logic) clientGameResult() {

	// {
	// 	"actionId":"9",
	// 	"id":"2",
	// 	"eventData":[
	// 		[1,0,1,-1,5,1,0,0,1,0,1],
	// 		[
	// 			[
	// 				["s7","s5","s1","s3","s6"],
	// 				["s5","s1","s7","s5","s7"],
	// 				["s6","s7","s5","s1","s1"]
	// 			],
	// 			[0,0,0,0,0],
	// 			[-1,0]
	// 		],
	// 		[],
	// 		[1]
	// 	]
	// }

	// {   押20 贏400
	// 	"id":"2",
	// 	"actionId":"9",
	// 	"eventData":[
	// 		[1,0,1,-1,5,20,400,0,20,0,1],
	// 		[
	// 			[
	// 				["s7","s7","s7","s4","s5"],
	// 				["s6","s3","s2","s3","s7"],
	// 				["s5","s4","s5","s7","s6"]
	// 			],
	// 			[0,0,0,0,0],
	// 			[-1,0]
	// 		],
	// 		[
	// 			[2,"s7",3,400,1]
	// 		],
	// 		[1,3]
	// 	]
	// }

	result, callback := logic.gameresult()
	defer callback()

	_m := new(network.Message)
	_m.ID = 2
	_m.ActionID = 9
	_m.EventData = result

	logic.cw.WriteClient(_m)
}

func (logic *Logic) clientPickGameItem(id int32, money int64) {

	win := logic.MoneyTransfer.S2C(float64(money))

	_m := new(network.Message)
	_m.ID = 2
	_m.ActionID = 9
	_m.EventData = []interface{}{
		[]interface{}{4, nil},
		[]interface{}{
			win,
			id,
			[]interface{}{win},
		},
		nil,
	}

	logic.cw.WriteClient(_m)
}

func (logic *Logic) clientJackpotGameMode() {

	{ // send data

		var status int
		if logic.Round == 0 {
			status = 1
		}

		var records []interface{}
		var winJP slot.JackpotLevel
		jpMap := make(map[slot.JackpotLevel]int)

		sort.Slice(logic.JackpotGameRecord.Records, func(i, j int) bool {
			return logic.JackpotGameRecord.Records[i].GetID() < logic.JackpotGameRecord.Records[j].GetID()
		})

		logic.pickcnt = 0

		for _, row := range logic.JackpotGameRecord.Records {
			lv := slot.NewJackpotLevel(int32(row.GetRewareValue()))
			jpMap[lv]++

			pos := row.GetLocation() - 1

			if pos > -1 {
				logic.pickcnt++
			}

			data := []interface{}{pos, logic.clientJPType(lv)}
			records = append(records, data)
		}

		for lv, cnt := range jpMap {
			if cnt == 3 {
				winJP = lv
			}
		}

		win := logic.MoneyTransfer.S2C(float64(logic.JackpotGameRecord.GetWinMoney()))

		_m := new(network.Message)
		_m.ID = 2
		_m.ActionID = 9
		_m.EventData = []interface{}{
			7,
			status,                               // status 0: 正常, 1: 斷線復原
			win,                                  // total payout
			win,                                  // balance payout
			0,                                    // given payout
			logic.JackpotGameRecord.GetBetRate(), // bet size
			logic.clientJPType(winJP),            // jackpot type
			records,                              // records
		}

		logic.cw.WriteClient(_m)
	}
}

func (logic *Logic) clientSelectionChoose(t int) {

	if logic.SelectGameOptions == nil {
		return
	}

	{ // send data

		opt := logic.SelectGameOptions.Options[t]
		data := make(map[string]interface{})

		if err := json.Unmarshal([]byte(opt.GetData()), &data); err != nil {
			mlog.Warn("Parse Selection Data Error, %v", err)
			return
		}

		_m := new(network.Message)
		_m.ID = 2
		_m.ActionID = 9
		_m.EventData = []interface{}{
			[]interface{}{2, 0},
			[]interface{}{1, int(data["count"].(float64))},
			nil,
		}

		logic.cw.WriteClient(_m)

		switch slot.GameMode(opt.GetType()) {
		case slot.GameModeFreeSpin:
			_m := new(network.Message)
			_m.ID = 2
			_m.ActionID = 2
			_m.EventData = []interface{}{
				1,
				[]interface{}{int(data["count"].(float64)), 1},
			}

			logic.cw.WriteClient(_m)
			logic.freedatasend = true
		}

	}
}

func (logic *Logic) clientJackpotPool(lvs ...slot.JackpotLevel) {

	{ // send data

		var data []interface{}

		// [ money , level , update]
		// [
		// 	[100,0,-1],
		// 	[1000,1,-1],
		// 	[10000,2,-1],
		// 	[100000,3,-1],
		// ]

		for _, lv := range lvs {
			if pool, ok := logic.JackpotPool[lv]; ok {
				_d := []interface{}{logic.MoneyTransfer.S2C(float64(pool.GetMoney())), lv.Int(), -1}
				data = append(data, _d)
			}
		}

		sort.Slice(data, func(i, j int) bool {
			return data[i].([]interface{})[1].(int32) < data[j].([]interface{})[1].(int32)
		})

		if len(data) == 0 {
			return
		}

		_m := new(network.Message)
		_m.ID = 2
		_m.ActionID = int(netproto.Slot_GameMessageClassID_Slot_JackpotPoolRetID)
		_m.EventData = data

		logic.cw.WriteClient(_m)
	}
}

func (logic *Logic) clientJPType(lv slot.JackpotLevel) int32 {
	switch lv {
	default:
		return -1
	case slot.JackpotGrand:
		return 3
	case slot.JackpotMajor:
		return 2
	case slot.JackpotMinor:
		return 1
	case slot.JackpotMini:
		return 0
	}
}

func (logic *Logic) clientUpdateMoney(money int64, delay bool) {
	// {"actionId":"3","eventData":[0,1000130],"id":"0"}
	_m := new(network.Message)
	_m.ID = 0
	_m.ActionID = 3

	var _delay int
	if delay {
		_delay = 1
	}
	_m.EventData = []interface{}{_delay, logic.MoneyTransfer.S2C(float64(money))}

	logic.cw.WriteClient(_m)
}

func (logic *Logic) clientSetting() {

	// [
	//   [ 1, 2, 5, 10, 20], // bet money
	//   [ [0, 10000, true, false], ...], // jackpot info [level, current money, 按下注比例, 自動跳增加]
	// ]

	var bets []float64
	var jpinfos []interface{}

	for _, money := range logic.Setting.BetMoneys {
		bets = append(bets, logic.MoneyTransfer.S2C(float64(money)))
	}
	sort.Float64s(bets)

	for lv, pool := range logic.JackpotPool {
		_d := []interface{}{
			lv.Int(),
			logic.MoneyTransfer.S2C(float64(pool.GetMoney())),
			!pool.GetIsLink(),
			false,
		}
		jpinfos = append(jpinfos, _d)
	}
	sort.Slice(jpinfos, func(i, j int) bool {
		return jpinfos[i].([]interface{})[0].(int32) < jpinfos[j].([]interface{})[0].(int32)
	})

	data := []interface{}{
		bets,
		jpinfos,
	}

	if len(data) == 0 {
		return
	}

	_m := new(network.Message)
	_m.ID = 2
	_m.ActionID = int(netproto.Slot_GameMessageClassID_Slot_SettingID)
	_m.EventData = data

	logic.cw.WriteClient(_m)
}

func (logic *Logic) clientLogin() {
	// // {"actionId":"1","eventData":[73810,"kgs800087",1000160,"2020/04/30 13:45:38",[[],null]],"id":"0"}
	m1 := new(network.Message)
	m1.ID = 0
	m1.ActionID = 1
	m1.EventData = []interface{}{
		logic.UserData.GetUserID(),
		logic.UserData.GetNickName(),
		logic.MoneyTransfer.S2C(float64(logic.UserData.GetCashAmount())),
		time.Now().Format(time.RFC3339),
		[]interface{}{
			[]interface{}{},
			nil,
		},
	}
	logic.cw.WriteClient(m1)
}

func (logic *Logic) clientStart(data interface{}) {
	// // {"actionId":"1","eventData":[[null,null,0,0],["kgs800087",1000160],[[1,0,1,-1,10,5,0,0,5,0,[]],[[["s5","s7","s3","b1","s1"],["s7","s5","s7","s1","s2"],["s2","s2","s1","s3","s6"]],[0,0,0,0,0],[-1,0]],[[3,"m1",3,100,1],[8,"s7",3,50,1]],[1]],-1],"id":"2"}
	m2 := new(network.Message)
	m2.ID = 2
	m2.ActionID = 1
	m2.EventData = make([]interface{}, 0)
	m2.EventData = append(m2.EventData.([]interface{}),
		[]interface{}{nil, nil, 0, 0},
		[]interface{}{logic.UserData.GetNickName(), logic.MoneyTransfer.S2C(float64(logic.UserData.GetCashAmount()))},
		data,
	)
	logic.cw.WriteClient(m2)
}

func (logic *Logic) clientEventFreeSpin(cnt int) {
	_m := new(network.Message)
	_m.ID = 2
	_m.ActionID = 2
	_m.EventData = []interface{}{
		1,
		[]interface{}{cnt, 1},
	}

	logic.cw.WriteClient(_m)
}

func (logic *Logic) transerStepResultToGameResult(steps *netproto.Slot_MSG_StepGameResult, idx int) *netproto.Slot_MSG_GameResult {
	res := steps.SlotResult[idx]

	before := make(map[slot.Position]netproto.Slot_ReelResult)
	after := make(map[slot.Position]netproto.Slot_ReelResult)
	for _, pos := range res.InitResult {
		p := slot.Position{
			Reel: pos.Position.GetReel(),
			Row:  pos.Position.GetRow(),
		}
		before[p] = netproto.Slot_ReelResult{
			Position: &netproto.Slot_ReelPosition{
				Reel: proto.Int32(pos.Position.GetReel()),
				Row:  proto.Int32(pos.Position.GetRow()),
			},
			Symbol: &netproto.Slot_Symbol{
				SymbolID:   proto.Int32(pos.Symbol.GetSymbolID()),
				SymbolName: proto.String(pos.Symbol.GetSymbolName()),
				Position: &netproto.Slot_ReelPosition{
					Reel: proto.Int32(pos.Position.GetReel()),
					Row:  proto.Int32(pos.Position.GetRow()),
				},
			},
			Blink: proto.Bool(pos.GetBlink()),
			Multi: proto.Int64(pos.GetMulti()),
		}
	}
	for _, pos := range res.SlotResult.ReelResult {
		p := slot.Position{
			Reel: pos.Position.GetReel(),
			Row:  pos.Position.GetRow(),
		}
		after[p] = netproto.Slot_ReelResult{
			Position: &netproto.Slot_ReelPosition{
				Reel: proto.Int32(pos.Position.GetReel()),
				Row:  proto.Int32(pos.Position.GetRow()),
			},
			Symbol: &netproto.Slot_Symbol{
				SymbolID:   proto.Int32(pos.Symbol.GetSymbolID()),
				SymbolName: proto.String(pos.Symbol.GetSymbolName()),
				Position: &netproto.Slot_ReelPosition{
					Reel: proto.Int32(pos.Position.GetReel()),
					Row:  proto.Int32(pos.Position.GetRow()),
				},
			},
			Blink: proto.Bool(pos.GetBlink()),
			Multi: proto.Int64(pos.GetMulti()),
		}
	}

	// 比對拓展前後，如果前後都是w1則替換為e1
	var newreel []*netproto.Slot_ReelResult
	for pos, r := range before {

		rr := &netproto.Slot_ReelResult{
			Position: &netproto.Slot_ReelPosition{
				Reel: proto.Int32(r.Position.GetReel()),
				Row:  proto.Int32(r.Position.GetRow()),
			},
			Symbol: &netproto.Slot_Symbol{
				SymbolID:   proto.Int32(r.Symbol.GetSymbolID()),
				SymbolName: proto.String(r.Symbol.GetSymbolName()),
				Position: &netproto.Slot_ReelPosition{
					Reel: proto.Int32(r.Position.GetReel()),
					Row:  proto.Int32(r.Position.GetRow()),
				},
			},
			Blink: proto.Bool(r.GetBlink()),
			Multi: proto.Int64(r.GetMulti()),
		}

		if r.Symbol.GetSymbolName() == after[pos].Symbol.GetSymbolName() && r.Symbol.GetSymbolName() == "w1" && res.GetIsExpand() {
			rr.Symbol.SymbolName = proto.String("e1")
		}

		newreel = append(newreel, rr)
	}

	sort.Slice(newreel, func(i, j int) bool {
		if newreel[i].GetPosition().GetReel() == newreel[j].GetPosition().GetReel() {
			return newreel[i].GetPosition().GetRow() < newreel[j].GetPosition().GetRow()
		}

		return newreel[i].GetPosition().GetReel() < newreel[j].GetPosition().GetReel()
	})

	newresult := *res.SlotResult
	newresult.ReelResult = newreel

	return &newresult
}
