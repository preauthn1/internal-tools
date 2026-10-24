package analyser

import (
	"fmt"
	"io"
	"os"
	"time"

	"example.invalid/dustin/go-humanize"
	"example.invalid/ServerTeam/Servers/internal/game/gen"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
)

// Run 跑
func Run(w io.Writer, cfgs *config.Group, proto *pachislot.Proto, pid int32, count int) {

	now := time.Now()

	var totalCost, totalWin int64
	var totalReleaseRT, totalReplay, total1Ggot, total1Guse int
	var total1GBB0, total1GBB1, total1GBB2, total1GJAC int
	var totalBB, totalRB, totalJAC int
	resultMap := make(map[pb.Scenario]map[pb.ResultType]int)
	resultMap[pb.Scenario_NG] = make(map[pb.ResultType]int)
	resultMap[pb.Scenario_BB] = make(map[pb.ResultType]int)
	resultMap[pb.Scenario_RB] = make(map[pb.ResultType]int)
	resultMap[pb.Scenario_JAC] = make(map[pb.ResultType]int)
	costMap := make(map[pb.Scenario]int64)
	winMap := make(map[pb.Scenario]int64)
	modeMap := make(map[pb.ModeType]map[pb.Scenario]int)
	modeMap[pb.ModeType_A] = make(map[pb.Scenario]int)
	modeMap[pb.ModeType_B] = make(map[pb.Scenario]int)
	modeMap[pb.ModeType_Heaven] = make(map[pb.Scenario]int)
	releaseMap := make(map[pb.ResultType]int)

	tcfg := cfgs.GetConfigSection(gamelogic.ConfigFileNameGameServer, "table")
	tcfg["distribution"] = fmt.Sprintf("%d:%d:%d", 1, pid, 10)

	u, _ := fixtures.NewGameUser(1000, "測試玩家", 1e9)

	grs := fixtures.NewTestGameRoomServerWithConfig(cfgs)
	game := grs.GameLogicManager.GetGameLogic(0).(*gen.Game)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u, score, tax)
	}

	modeMap[game.PachiSlotGame.Mode][game.PachiSlotGame.Bonus]++
	game.PachiSlotGame.AddAfterInitFunc(func(g *pachislot.Game, r bool) {
		modeMap[g.Mode][g.Bonus]++
		if r {
			total1Guse++
		}
	})

	u.UserState.OnUserQueue(u, game.TableManager.Room)
	game.UserEnter(u, 0, 0)

	game.OnSendMessage(func(uid int32, m *msg.Message) {
		switch netproto.GEN_GameMessageClassID(m.SClassID) {
		case netproto.GEN_GameMessageClassID_GEN_EventBB0JACID:
			data := m.MsgData.(*netproto.GEN_Payout)
			switch data.Type {
			case netproto.GEN_PayType_GEN_PT1G:
				total1Ggot += int(data.Amount)
			}
			total1GBB0++
		case netproto.GEN_GameMessageClassID_GEN_EventBB1JACID:
			data := m.MsgData.(*netproto.GEN_Payout)
			switch data.Type {
			case netproto.GEN_PayType_GEN_PT1G:
				total1Ggot += int(data.Amount)
			}
			total1GBB1++
		case netproto.GEN_GameMessageClassID_GEN_EventBB2JACID:
			data := m.MsgData.(*netproto.GEN_Payout)
			switch data.Type {
			case netproto.GEN_PayType_GEN_PT1G:
				total1Ggot += int(data.Amount)
			}
			total1GBB2++
		case netproto.GEN_GameMessageClassID_GEN_EventJAC8BellID:
			data := m.MsgData.(*netproto.GEN_Payout)
			switch data.Type {
			case netproto.GEN_PayType_GEN_PT1G:
				total1Ggot += int(data.Amount)
			}
			total1GJAC++
		case netproto.GEN_GameMessageClassID_GEN_GameStatusID:
			data := m.MsgData.(*netproto.GEN_GameStatus)
			if data.Replay {
				totalReplay++
			}
		case netproto.GEN_GameMessageClassID_GEN_GameResultID:
			data := m.MsgData.(*netproto.GEN_GameResult)
			totalCost += data.BetMoney
			totalWin += data.WinMoney
			resultMap[pb.Scenario(data.Scenario)][pb.ResultType(data.Result)]++
			costMap[pb.Scenario(data.Scenario)] += data.BetMoney
			winMap[pb.Scenario(data.Scenario)] += data.WinMoney
			if data.ReleaseRT {
				releaseMap[pb.ResultType(data.Result)]++
				totalReleaseRT++
			}

			for _, p := range data.Payouts {
				switch p.Type {
				case netproto.GEN_PayType_GEN_PT1G:
					total1Ggot += int(p.Amount)
				case netproto.GEN_PayType_GEN_PTBonus:
					switch p.Bonus {
					case netproto.GEN_Scenario_GEN_Scenario_BB:
						totalBB++
					case netproto.GEN_Scenario_GEN_Scenario_RB:
						totalRB++
					case netproto.GEN_Scenario_GEN_Scenario_JAC:
						totalJAC++
					}
				}
			}
		}
	})

	begin := time.Now()
	for game.PachiSlotGame.RoadMap.TotalCnt < count {
		if game.PachiSlotGame.RoadMap.TotalCnt > 0 && game.PachiSlotGame.RoadMap.TotalCnt%1e5 == 0 {
			fmt.Fprintf(os.Stdout, "Spined: %s (Times)(%v)\n", humanize.Comma(int64(game.PachiSlotGame.RoadMap.TotalCnt)), time.Since(begin))
			begin = time.Now()
		}
		game.Play(1, u)
	}

	var ng, bb, rb, jac int
	var bb1g, rb1g, jac1g int
	for _, cnt := range resultMap[pb.Scenario_NG] {
		ng += cnt
	}
	for t, cnt := range resultMap[pb.Scenario_BB] {
		bb += cnt
		switch t {
		case pb.ResultType_BlueSevenLink, pb.ResultType_RedSevenLink:
			bb1g += cnt
		case pb.ResultType_DoubleSevenLink:
			bb1g += 2 * cnt
		case pb.ResultType_FeatureNone:
			bb1g += 3 * cnt
		}
	}
	for t, cnt := range resultMap[pb.Scenario_RB] {
		rb += cnt
		switch t {
		case pb.ResultType_JACNone:
			rb1g += cnt
		}
	}
	for t, cnt := range resultMap[pb.Scenario_JAC] {
		jac += cnt
		switch t {
		case pb.ResultType_JACNone:
			jac1g += cnt
		}
	}

	var modeA, modeB, modeH int
	for _, cnt := range modeMap[pb.ModeType_A] {
		modeA += cnt
	}
	for _, cnt := range modeMap[pb.ModeType_B] {
		modeB += cnt
	}
	for _, cnt := range modeMap[pb.ModeType_Heaven] {
		modeH += cnt
	}

	fmt.Fprintln(w, "=======================================================")
	fmt.Fprintf(w, "Pattern: %d\n", pid)
	fmt.Fprintf(w, "Total Spin: %s (Times)\n", humanize.Comma(int64(game.PachiSlotGame.RoadMap.TotalCnt)))
	fmt.Fprintf(w, "Total Cost: NG(%s) + BB(%s) + RB(%s) + JAC(%s) = %s (Coin)\n", humanize.Comma(costMap[pb.Scenario_NG]), humanize.Comma(costMap[pb.Scenario_BB]), humanize.Comma(costMap[pb.Scenario_RB]), humanize.Comma(costMap[pb.Scenario_JAC]), humanize.Comma(totalCost))
	fmt.Fprintf(w, "Total Win: NG(%s) + BB(%s) + RB(%s) + JAC(%s) = %s (Coin)\n", humanize.Comma(winMap[pb.Scenario_NG]), humanize.Comma(winMap[pb.Scenario_BB]), humanize.Comma(winMap[pb.Scenario_RB]), humanize.Comma(winMap[pb.Scenario_JAC]), humanize.Comma(totalWin))
	fmt.Fprintf(w, "Total Replay: %s (Times)\n", humanize.Comma(int64(totalReplay)))
	fmt.Fprintf(w, "Total ReleaseRT: %s (Times)\n", humanize.Comma(int64(totalReleaseRT)))
	for t, cnt := range releaseMap {
		fmt.Fprintf(w, "  %v: %s(%.3f%%), Release Rate(%.3f%%)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(totalReleaseRT)*100, float64(cnt)/float64(resultMap[pb.Scenario_NG][t])*100)
	}
	fmt.Fprintf(w, "Total 1G: Get %s (Times), Use %s (Times)\n", humanize.Comma(int64(total1Ggot)), humanize.Comma(int64(total1Guse)))
	fmt.Fprintf(w, "  Win 0 JAC in BB : %s (Times)(1G * %s)\n", humanize.Comma(int64(total1GBB0)), humanize.Comma(int64(total1GBB0*3)))
	fmt.Fprintf(w, "  Win 1 JAC in BB : %s (Times)(1G * %s)\n", humanize.Comma(int64(total1GBB1)), humanize.Comma(int64(total1GBB1*2)))
	fmt.Fprintf(w, "  Win 2 JAC in BB : %s (Times)(1G * %s)\n", humanize.Comma(int64(total1GBB2)), humanize.Comma(int64(total1GBB2*1)))
	fmt.Fprintf(w, "  Win 8 Bell in JAC : %s (Times)(1G * %s)\n", humanize.Comma(int64(total1GJAC)), humanize.Comma(int64(total1GJAC*1)))
	fmt.Fprintf(w, "Final RoadMap: %v\n", game.PachiSlotGame.RoadMap)
	fmt.Fprintf(w, "  Win BB : %s (Times)\n", humanize.Comma(int64(totalBB)))
	fmt.Fprintf(w, "  Win RB : %s (Times)\n", humanize.Comma(int64(totalRB)))
	fmt.Fprintf(w, "  Win JAC : %s (Times)\n", humanize.Comma(int64(totalJAC)))
	fmt.Fprintf(w, "Mode A(%s): BB(%s)(%.3f%%) RB(%s)(%.3f%%)\n", humanize.Comma(int64(modeA)), humanize.Comma(int64(modeMap[pb.ModeType_A][pb.Scenario_BB])), float64(modeMap[pb.ModeType_A][pb.Scenario_BB])/float64(modeA)*100, humanize.Comma(int64(modeMap[pb.ModeType_A][pb.Scenario_RB])), float64(modeMap[pb.ModeType_A][pb.Scenario_RB])/float64(modeA)*100)
	fmt.Fprintf(w, "Mode B(%s): BB(%s)(%.3f%%) RB(%s)(%.3f%%)\n", humanize.Comma(int64(modeB)), humanize.Comma(int64(modeMap[pb.ModeType_B][pb.Scenario_BB])), float64(modeMap[pb.ModeType_B][pb.Scenario_BB])/float64(modeB)*100, humanize.Comma(int64(modeMap[pb.ModeType_B][pb.Scenario_RB])), float64(modeMap[pb.ModeType_B][pb.Scenario_RB])/float64(modeB)*100)
	fmt.Fprintf(w, "Mode Heaven(%s): BB(%s)(%.3f%%) RB(%s)(%.3f%%)\n", humanize.Comma(int64(modeH)), humanize.Comma(int64(modeMap[pb.ModeType_Heaven][pb.Scenario_BB])), float64(modeMap[pb.ModeType_Heaven][pb.Scenario_BB])/float64(modeH)*100, humanize.Comma(int64(modeMap[pb.ModeType_Heaven][pb.Scenario_RB])), float64(modeMap[pb.ModeType_Heaven][pb.Scenario_RB])/float64(modeH)*100)
	fmt.Fprintf(w, "RTP: %.3f%%\n", float64(totalWin)/float64(totalCost)*100)
	fmt.Fprintf(w, "  NG: %.3f%% (%.3f%%)\n", float64(winMap[pb.Scenario_NG])/float64(costMap[pb.Scenario_NG])*100, float64(winMap[pb.Scenario_NG])/float64(totalCost)*100)
	fmt.Fprintf(w, "  BB: %.3f%% (%.3f%%)\n", float64(winMap[pb.Scenario_BB])/float64(costMap[pb.Scenario_BB])*100, float64(winMap[pb.Scenario_BB])/float64(totalCost)*100)
	fmt.Fprintf(w, "  RB: %.3f%% (%.3f%%)\n", float64(winMap[pb.Scenario_RB])/float64(costMap[pb.Scenario_RB])*100, float64(winMap[pb.Scenario_RB])/float64(totalCost)*100)
	fmt.Fprintf(w, "  JAC: %.3f%% (%.3f%%)\n", float64(winMap[pb.Scenario_JAC])/float64(costMap[pb.Scenario_JAC])*100, float64(winMap[pb.Scenario_JAC])/float64(totalCost)*100)
	fmt.Fprintln(w, "=======================================================")
	fmt.Fprintf(w, "Result in NG: (%s)\n", humanize.Comma(int64(ng)))
	for t, cnt := range resultMap[pb.Scenario_NG] {
		fmt.Fprintf(w, "  %v: %s(%.3f%%)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(ng)*100)
	}
	fmt.Fprintf(w, "Result in BB: (%s)(1G * %s)\n", humanize.Comma(int64(bb)), humanize.Comma(int64(bb1g)))
	for t, cnt := range resultMap[pb.Scenario_BB] {
		switch t {
		default:
			fmt.Fprintf(w, "  %v: %s(%.3f%%)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(bb)*100)
		case pb.ResultType_BlueSevenLink, pb.ResultType_RedSevenLink:
			fmt.Fprintf(w, "  %v: %s(%.3f%%), (1G * %s)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(bb)*100, humanize.Comma(int64(cnt*1)))
		case pb.ResultType_DoubleSevenLink:
			fmt.Fprintf(w, "  %v: %s(%.3f%%), (1G * %s)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(bb)*100, humanize.Comma(int64(cnt*2)))
		case pb.ResultType_FeatureNone:
			fmt.Fprintf(w, "  %v: %s(%.3f%%), (1G * %s)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(bb)*100, humanize.Comma(int64(cnt*3)))
		}
	}
	fmt.Fprintf(w, "Result in RB: (%s)(1G * %s)\n", humanize.Comma(int64(rb)), humanize.Comma(int64(rb1g)))
	for t, cnt := range resultMap[pb.Scenario_RB] {
		switch t {
		default:
			fmt.Fprintf(w, "  %v: %s(%.3f%%)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(rb)*100)
		case pb.ResultType_JACNone:
			fmt.Fprintf(w, "  %v: %s(%.3f%%), (1G * %s)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(rb)*100, humanize.Comma(int64(cnt*1)))
		}
	}
	fmt.Fprintf(w, "Result in JAC: (%s)(1G * %s)\n", humanize.Comma(int64(jac)), humanize.Comma(int64(jac1g)))
	for t, cnt := range resultMap[pb.Scenario_JAC] {
		switch t {
		default:
			fmt.Fprintf(w, "  %v: %s(%.3f%%)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(jac)*100)
		case pb.ResultType_JACNone:
			fmt.Fprintf(w, "  %v: %s(%.3f%%), (1G * %s)\n", t, humanize.Comma(int64(cnt)), float64(cnt)/float64(jac)*100, humanize.Comma(int64(cnt*1)))
		}
	}
	fmt.Fprintln(w, "=======================================================")

	fmt.Fprintf(os.Stdout, "Spend (%v)\n", time.Since(now))
}
