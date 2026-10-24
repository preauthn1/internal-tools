package jjbx

import (
	"context"
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"

	apipb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// ActionSpin 玩家下注動作
func (logic *Logic) ActionSpin(user *gamelogic.GameUser, gmsg *msg.Message) {

	data := &netproto.Slot_MSG_Spin{}
	if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
		mlog.Error("[玩家Spin] 接收玩家封包解析錯誤. %v", err)
		return
	}

	ok, ret, money := logic.Bet(user, data.GetMoney(), int64(data.GetRate()))

	// 下注狀態
	logic.SendBetStatusMessageToUser(user, ret, money)

	if !ok {
		return
	}

	logic.GameLog.GameBegin()

	logic.UpdateState(slot.GameStateAward)

	symbolGroupID, _, _, _ := logic.PrepareSpin()

	logic.Spin(symbolGroupID)

	logic.GameResultManager.RowMode()

	if logic.FreeSpin != nil {
		logic.UpdateFreeSpin(user)
		logic.SendLeftFreeGameToUser(user)
	}

	// 紀錄
	logic.GameResultManager.ReelMode()
	if logic.IsFreeSpin() {
		logic.AddSlotGameHistoryLog(user, slot.HistoryTypeFREE, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	} else {
		logic.AddSlotGameHistoryLog(user, slot.HistoryTypeMAIN, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	}
	logic.GameResultManager.RowMode()

	if logic.FreeSpin != nil {
		if logic.FreeSpin.FreeSpinExtAward != nil {

			// 計算盤面g7數量
			var cnt int32
			symbols := make(map[slot.Position]*protodata.SlotSymbolDistribution)
			logic.GameResultManager.RowMode()
			for rowID, row := range logic.GameResultManager.GameResult {
				for reelID, symbol := range row {
					if symbol.GetSymbol() == "g7" {
						cnt++
						pos := slot.Position{Reel: int32(reelID + 1), Row: int32(rowID)}
						symbols[pos] = symbol
					}
				}
			}

			logic.FreeSpin.FreeSpinExtAward.TotalCount += cnt
			logic.FreeSpin.FreeSpinExtAward.TotalWinMoney += logic.FreeSpin.FreeSpinExtAward.EachCountMoney * int64(cnt)
			logic.FreeSpin.FreeSpinExtAward.Count = cnt
			logic.FreeSpin.FreeSpinExtAward.WinMoney = logic.FreeSpin.FreeSpinExtAward.EachCountMoney * int64(cnt)
			logic.FreeSpin.FreeSpinExtAward.Symbols = symbols

			logic.GameLog.AddGameStep("出現 %d 個『黃幣』, 額外獎勵增加 %d, 累積 %d", logic.FreeSpin.FreeSpinExtAward.Count, logic.FreeSpin.FreeSpinExtAward.WinMoney, logic.FreeSpin.FreeSpinExtAward.TotalWinMoney)
		}
	}

	// 計算流水
	logic.CalcFlowAmount(user)

	// 寫分
	logic.UpdateUserScore(user)

	var skipEndLog bool

	if logic.FreeSpin != nil {
		if logic.FreeSpin.FreeSpinExtAward != nil {

			// 檢查是否結束
			if logic.FreeSpin.Count == 0 {
				logic.GameEndLog()

				if !conf.AnalysisMode {
					time.Sleep(100 * time.Millisecond)
				}

				logic.GameLog.GameBegin()

				logic.FreeSpin.FreeSpinExtAward.IsFinished = true

				logic.GameLog.AddGameStep("免費遊戲結束, 開始結算額外獎勵, 累計出現『黃幣』%d 次, 每次獎勵 %d, 共獎勵 %d", logic.FreeSpin.FreeSpinExtAward.TotalCount, logic.FreeSpin.FreeSpinExtAward.EachCountMoney, logic.FreeSpin.FreeSpinExtAward.TotalWinMoney)

				tmpwin := logic.GameWinMoney

				logic.GameWinMoney = logic.FreeSpin.FreeSpinExtAward.TotalWinMoney

				logic.AddSlotGameHistoryLog(user, slot.HistoryTypeFREE, slot.LogTypeTOTALCOIN, fmt.Sprintf("%d|%d|%d", logic.FreeSpin.FreeSpinExtAward.TotalCount, logic.FreeSpin.FreeSpinExtAward.EachCountMoney, logic.FreeSpin.FreeSpinExtAward.TotalWinMoney))

				// 計算流水
				logic.CalcFlowAmount(user)

				if err := app.Context.Invoke(func(c apipb.GameBetServiceClient, pm *player.Manager) error {
					p := pm.GetPlayerByUserID(user.UserID)
					req := new(apipb.PlaceBetReq)
					req.AccessToken = [REDACTED]
					req.ServerID = logic.TableManager.Room.ServerID
					req.TradeID = user.TradeID
					req.BetAmount = 0

					req.Memo = fmt.Sprintf("免費遊戲結束, 開始結算額外獎勵, 累計出現『黃幣』%d 次, 每次獎勵 %d, 共獎勵 %d", logic.FreeSpin.FreeSpinExtAward.TotalCount, logic.FreeSpin.FreeSpinExtAward.EachCountMoney, logic.FreeSpin.FreeSpinExtAward.TotalWinMoney)

					res, err := c.PlaceBet(context.Background(), req)
					if err != nil {
						return err
					}

					logic.TransactionNo = res.TransactionNo

					return nil
				}); err != nil {
					mlog.Error("Place Bet Error with FreeSpinExtAward, %v", err)
					logic.KickUserFunc(user.UserID, "System Error ...")
					return
				}

				// 寫分
				logic.UpdateUserScore(user)

				logic.GameEndLog()

				skipEndLog = true

				logic.GameWinMoney += tmpwin
			}
		}
	}

	// 傳送開獎結果
	logic.SendGameResultMessage(user)

	if !skipEndLog {
		logic.GameEndLog()
	}

	// 調整狀態為等待
	logic.UpdateState(slot.GameStateWait)

	logic.GameEndFunc()
}
