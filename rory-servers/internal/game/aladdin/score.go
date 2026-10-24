package aladdin

import (
	"context"

	apipb "example.invalid/ServerTeam/Servers/internal/api-service/apisvc"
	"example.invalid/ServerTeam/Servers/pkg/app"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// UpdateUserScore 更新玩家分數
func (logic *Logic) UpdateUserScore(user *gamelogic.GameUser, step *StepResult) {
	if user.IsRobot() {
		return
	}

	logic.GameLog.AddGameStep("更新玩家金幣, 下注金幣: %d, 贏得獎金: %d", step.GamePayMoney, step.GameWinMoney)

	// logic.UpdateUserScoreFunc(user.UserID, step.GameWinMoney-step.GamePayMoney, 0)

	if err := app.Context.Invoke(func(c apipb.GameResultServiceClient) error {

		req := new(apipb.SettlementReq)
		req.PlayerID = user.UserID
		req.ServerID = logic.TableManager.Room.ServerID
		req.TradeID = user.TradeID
		req.UpdateAmount = step.GameWinMoney
		req.TransactionNo = logic.TransactionNo
		req.BetAmount = step.GamePayMoney
		req.ValidAmount = step.GamePayMoney
		req.PayoutAmount = step.GameWinMoney
		req.TaxAmount = 0

		res, err := c.SettlementV2(context.Background(), req)
		if err != nil {
			mlog.Error("更新玩家金幣錯誤 SettlemmentV2 Error, %v", err)
			return err
		}

		score := logic.GameWinMoney - logic.GamePayMoney

		// 原先 UpdateUserScoreFunc 內做得要加回來
		user.UpdateScore += score

		// 原先在 roomactionmanager.go 內做得要加回來
		if score > 0 {
			user.WinCount++
			user.UpdateWinCount++
			user.LastWinLose = 1
		} else if score == 0 {
			user.DrawCount++
			user.UpdateDrawCount++
			user.LastWinLose = 0
		} else {
			user.LoseCount++
			user.UpdateLoseCount++
			user.LastWinLose = 0
		}

		user.Score = res.AfterAmount

		return nil
	}); err != nil {
		mlog.Error("更新玩家金幣錯誤 %v", err)
		logic.KickUserFunc(user.UserID, "System Error ...")
	}
}
