package baselogic

import (
	"context"
	"fmt"

	"example.invalid/golang/protobuf/proto"
	apipb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// ActionLightningSpin 玩家表演Lightning
func (logic *SlotLogicBase) ActionLightningSpin(user *gamelogic.GameUser, gmsg *msg.Message) {
	if logic.LightningGame == nil {
		mlog.Error("[玩家演完Lightning] 沒有Lightning遊戲")
		return
	}

	data := new(netproto.Slot_MSG_LightningShow)

	if err := proto.Unmarshal(gmsg.MsgData.([]byte), data); err != nil {
		mlog.Error("[玩家演完Lightning] 接收玩家封包解析錯誤. %v", err)
		return
	}

	logic.GameLog.GameBegin()
	logic.Clean()

	if logic.LightningGame.GetLastStep().LastChance == 0 {
		mlog.Error("[玩家演完Lightning] Lightning 遊戲已結束")
		return
	}

	logic.LightningGame.ShowStep = int(data.GetLightningStep())

	if logic.LightningGame.ShowStep > len(logic.LightningGame.Steps) {
		mlog.Error("[玩家演完Lightning] 超出範圍的Step")
		return
	}

	logic.LightningGame.Spin()

	prev := logic.LightningGame.Steps[logic.LightningGame.ShowStep-1]
	step := logic.LightningGame.Steps[logic.LightningGame.ShowStep]

	// 紀錄開獎符號
	logic.LogGameResult(step.GameResultManager)
	logic.GameLog.AddGameStep("Lightning 剩餘 %d 次, 累計金額: %d", prev.LastChance, logic.LightningGame.WinMoney)

	// 若有 Jackpot 就重置 Jackpot
	for _, award := range step.Awards {
		if award.Type != slot.PayTypeJackpot {
			continue
		}
		if jackpot, ok := logic.Jackpots[award.JackpotType]; ok {
			logic.CalcJackpotByMoney(jackpot, logic.LightningGame.BetMoney*logic.LightningGame.BetRate, int64(logic.MaxBetCount*logic.MaxBetRate))
		}
	}

	// 紀錄中獎結果
	if len(step.Awards) > 0 && step.LastChance > 0 {
		logic.GameLog.AddGameStep("Lightning 出現新錢幣:")
		step.PrintAward(logic)
	}

	step.PrintPayResult(logic)

	// 檢查是否表演結束
	if step.LastChance == 0 {
		logic.GameLog.AddGameStep("Lightning 遊戲結束, 開始結算")

		for _, s := range logic.LightningGame.Steps {
			s.PrintAward(logic)
		}

		logic.GameWinMoney += logic.LightningGame.WinMoney

		// 計算流水
		logic.CalcFlowAmount(user)

		// 寫分
		step.GameResultManager.ReelMode()
		logic.AddSlotGameHistoryLog(user, slot.HistoryTypeLIGHTNING, slot.LogTypeREEL, step.GameResultManager.ToSymbolString())
		step.GameResultManager.RowMode()

		if err := app.Context.Invoke(func(c apipb.GameBetServiceClient, pm *player.Manager) error {
			p := pm.GetPlayerByUserID(user.UserID)
			req := new(apipb.PlaceBetReq)
			req.AccessToken = [REDACTED]
			req.ServerID = logic.TableManager.Room.ServerID
			req.TradeID = user.TradeID
			req.BetAmount = logic.GamePayMoney

			req.Memo = fmt.Sprintf("Lightning")

			res, err := c.PlaceBet(context.Background(), req)
			if err != nil {
				return err
			}

			logic.TransactionNo = res.TransactionNo

			return nil
		}); err != nil {
			mlog.Error("Place Bet Error with Lightning, %v", err)
			logic.KickUserFunc(user.UserID, "System Error ...")
			return
		}

		logic.UpdateUserScore(user)
	} else {
		step.GameResultManager.ReelMode()
		logic.AddSlotGameHistoryLog(user, slot.HistoryTypeRESPIN, slot.LogTypeREEL, step.GameResultManager.ToSymbolString())
		step.GameResultManager.RowMode()
	}

	logic.SendLightningResultMessage(user)

	logic.RecordLightningData(user)

	// 檢查是否表演結束
	if step.LastChance == 0 {
		logic.LightningGame = nil
	}

	logic.GameEndLog()

	// 調整狀態為等待
	logic.UpdateState(slot.GameStateWait)

	logic.GameEndFunc()
}
