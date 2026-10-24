package fkeb

import (
	"example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"	
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 更新超端資訊
func (fkeb *FKEBLogic) updateVipUserWin(user *gamelogic.GameUser, updatemoney int64) bool {

	// 如果玩家資料為空時，就不執行超端更新
	if user == nil {
		return false
	}

	superToolHandler := fkeb.SuperToolService.GetBrGameSuperToolHandler()

	// 如果在超端的玩家資料為空值時，Complete為1，或者是Status為0，就不執行超端更新
	vipuser := fkeb.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(user.UserID)
	if vipuser == nil || vipuser.GetComplete() == int32(1) || vipuser.GetStatus() == int32(0) {
		return false
	}

	// 此行移到這裡，狀況就變成只要玩家有下注，就更新份數
	vipuser.CompleteAmount = proto.Int64(vipuser.GetCompleteAmount() + updatemoney)

	// 如果updatemoney和CompleteAmout都同時大於零 或同時小於零，就直接更新分數
	if updatemoney > int64(0) && vipuser.GetAssignAmount() > int64(0) { // 在到達完成金額一瞬間，在放分狀態，完成金額和更新金額同時為正
		// vipuser.CompleteAmount = proto.Int64(vipuser.GetCompleteAmount() + updatemoney) // 此行移出來
		if vipuser.GetCompleteAmount() >= vipuser.GetAssignAmount() { // 如果超端運作達到完成金額，就更新超端玩家狀態
			vipuser.Status = proto.Int32(0)   // 狀態為停止
			vipuser.Complete = proto.Int32(1) // 獎態為完成
		}
	}

	if updatemoney < int64(0) && vipuser.GetAssignAmount() < int64(0) { // 在到達完成金額一瞬間，在殺分狀態，完成金額和更新金額同時為負
		// vipuser.CompleteAmount = proto.Int64(vipuser.GetCompleteAmount() + updatemoney) // 此行移出來
		if vipuser.GetCompleteAmount() <= vipuser.GetAssignAmount() { // 如果超端運作達到完成金額，就更新超端玩家狀態
			vipuser.Status = proto.Int32(0)   // 狀態為停止
			vipuser.Complete = proto.Int32(1) // 獎態為完成
		}
	}

	// 之後更新營運的超端畫面
	superToolHandler.SaveVipUser(user.UserID)
	fkeb.SuperToolService.RefreshVIPUser()

	return true
}

// 交換二八槓的牌
// 莊家POSTION為0，順為1，天為2，地為3，利用每個區域的POSTION去暗地裡交換每個區域的麻將
func (fkeb *FKEBLogic) exchangeCard(postition1, position2 int) {
	majiang1TempForExchange := fkeb.majiangresult[postition1].majiang1
	majiang2TempForExchange := fkeb.majiangresult[postition1].majiang2
	fkeb.majiangresult[postition1].majiang1 = fkeb.majiangresult[position2].majiang1
	fkeb.majiangresult[postition1].majiang2 = fkeb.majiangresult[position2].majiang2
	fkeb.majiangresult[position2].majiang1 = majiang1TempForExchange
	fkeb.majiangresult[position2].majiang2 = majiang2TempForExchange
}

//監控放殺分
func (fkeb *FKEBLogic) bowJiangType(user *gamelogic.GameUser) bool {
	result := false

	// 防呆，真人上莊或無玩家資料，不執行超端
	if !fkeb.GetUserByUserID(fkeb.zhuangID).IsRobot() || user == nil {
		return result
	}

	// 防呆，如果超端玩家沒資料，不執行超端
	vipuser := fkeb.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(user.UserID)
	if vipuser == nil || vipuser.GetStatus() == int32(0) || vipuser.GetComplete() == int32(1) {
		return result
	}

	// 防呆，如果玩家沒下注，不執行超端
	currentUser := fkeb.getUserDataByID(vipuser.GetUserID())
	if len(currentUser.userbet) <= 0 {
		return result
	}

	// 對userbet進行修正，有些區域的下注金額太小，不列入考量
	// 一開始要設定會被程式記錄的下注金額，請在gameserver.dat檔裡，新增以下內容
	// [supertool]
	// ctrlbetlimit = 500 // 監控名單會控制的最低注額(分)
	correctUserbet := []*bettype{}
	for i := 0; i < len(currentUser.userbet); i++ {
		if currentUser.userbet[i].money >= fkeb.ctrlbetlimit {
			correctUserbet = append(correctUserbet, currentUser.userbet[i])
		}
	}

	// 如果玩家沒有任何一個區域下注到一定金額，不執行超端
	if len(correctUserbet) <= 0 {
		return result
	}

	// 針對殺分同時又是下注二門進行特別處理
	// 做法是先把第二大的牌換給莊家
	// 再把最大的牌換給玩家沒有下注的區域
	// 在這程狀況下，開獎區地在玩家未下注的區域，玩家很容易賠錢
	if (len(correctUserbet)) == 2 && (vipuser.GetAssignAmount() < 0) { // 只針對殺分狀態進行處理
		// 把莊家的牌換成第二大的牌
		secondint, _ := fkeb.findMaxMinResultCard("second") // 找出第二大牌的位置
		fkeb.exchangeCard(0, secondint)
		for k := 1; k < 4; k++ { // 用迴圈去找玩家未下注的區域
			if k != int(correctUserbet[0].fkebtype) && k != int(correctUserbet[1].fkebtype) {
				maxint, _ := fkeb.findMaxMinResultCard("max") // 找出最大牌的位置
				fkeb.exchangeCard(k, maxint)
				return result
			}
		}
	}

	if vipuser.GetComplete() == 0 { // 如果狀態為未完成

		mlog.Debug("實做放分 userid[%v] AssignAmount[%v] assignRate[%v] Complete[%v] CompleteAmount[%v] Status[%v]",
			vipuser.GetUserID(), vipuser.GetAssignAmount(), vipuser.GetAssignRate(), vipuser.GetComplete(),
			vipuser.GetCompleteAmount(), vipuser.GetStatus())

		rnd1 := int32(util.RndInt(0, 100)) // 產生0(含)到99(含)
		// rnd1 = int32(150) // 強制執行超端(測試用)
		mlog.System("用戶%d 執行放分機率 %d, rnd[%d]", user.UserID, vipuser.GetAssignRate(), rnd1)

		// 產生亂數，此亂數小於被記錄的筆數
		rnd2 := 0
		if len(correctUserbet) > 0 { // rnd2為亂數，範圍在 0 到 下注被記錄的筆數
			rnd2 = int(util.RndInt(0, len(correctUserbet))) // 生成范围内的随机整数, 不包含 "max"
		}
        fkeb.GameLog.AddGameStep("用戶%d 執行放分機率 %d, 執行rnd1[%d] 控制門數rnd2[%d]", user.UserID, vipuser.GetAssignRate(), rnd1, rnd2)
		if rnd1 >= vipuser.GetAssignRate() { //隨機0-99 大於勝率則執行放分，這裡可以設定放分和殺分的速度

			// 如果玩家多門下注只開一門時
			// 做法是先把莊家的牌換成第二大的牌
			// 再把玩家的牌換成最大或最小的牌，形成只開一門的狀態
			if len(currentUser.userbet) > 1 {
				secondint, secondbool := fkeb.findMaxMinResultCard("second") // 找出第二大牌的位置
				if secondint != 5 && secondbool != false {
					// 把莊家的牌換成第二大的牌
					fkeb.exchangeCard(0, secondint) // 莊家的牌換成第二大牌
					// 再把玩家換成最大或最小的牌
					if vipuser.GetAssignAmount() > 0 {
						maxint, maxbool := fkeb.findMaxMinResultCard("max") // 找出最大牌的位置
						if maxbool == true {
							fkeb.exchangeCard(int(correctUserbet[rnd2].fkebtype), maxint) // 玩家手上牌和最大牌交換
						}
					}
					if vipuser.GetAssignAmount() < 0 {
						minint, minbool := fkeb.findMaxMinResultCard("min") // 找出最小牌的位置
						if minbool == true {
							fkeb.exchangeCard(int(correctUserbet[rnd2].fkebtype), minint) // 玩家手上牌和最小牌交換
						}
					}
					result = true // 有成功修改牌型
				}
			}

			// 在其餘狀況下，放分時，把最大的牌換到玩家手上
			if vipuser.GetAssignAmount() > 0 {
				// 在這裡編寫放分主要邏輯部份
				maxint, maxbool := fkeb.findMaxMinResultCard("max")
				if maxint != 5 && maxbool != false {
					// 把最大的牌換到玩家手上
					fkeb.exchangeCard(int(correctUserbet[rnd2].fkebtype), maxint)
					result = true // 有成功修改牌型
				}
				// 主要放分邏輯部份結束
			}

			// 在其餘狀況下，殺分時，把最小的牌換到玩家手上
			if vipuser.GetAssignAmount() < 0 {
				// 在這裡編寫殺分主要邏輯部份
				minint, minbool := fkeb.findMaxMinResultCard("min")
				if minint != 5 && minbool != false {
					// 把最小的牌換到玩家手上
					fkeb.exchangeCard(int(correctUserbet[rnd2].fkebtype), minint)
					result = true // 有成功修改牌型
				}
				// 主要殺分邏輯部份結束
			}
		}
	}

	// 到最後時，把開獎結果重新計算一次，並更新程式
	fkeb.jiangtype = []int32{}
	for j := 1; j < 4; j++ { // j 由 1 到 3
		if fkeb.getMajiangsType(fkeb.majiangresult[j]) > fkeb.getMajiangsType(fkeb.majiangresult[0]) {
			fkeb.jiangtype = append(fkeb.jiangtype, int32(j))
		}
	}

	// 如果沒有順天地都沒中獎，莊家輸贏
	// 開獎結果為0時，代表通殺(順天地三區域都沒中獎)，莊家贏錢
	if len(fkeb.jiangtype) == 0 {
		fkeb.jiangtype = append(fkeb.jiangtype, int32(0))
	}
	// 修改開獎結果完成

	return result
}

// 模擬無畫面監控放分 (只有放分部份)
func (fkeb *FKEBLogic) testbowJiangType() bool {
	AssignAmount := 10000 // 允許放分的金額
	Complete := 0         // 模擬超端設定值，在蓋牌前進行超端進行超端控制，開牌時就不能更改牌型
	AssignRate := 100     // 模擬超端設定值，分配送分時的牌型組合，50%比較合理，先100%測試
	// vipuser := fkeb.SuperToolService.GetBrGameSuperToolHandler().GetVipUser(user.UserID) // 第一步，先取到使用者資料
	vipuser := fkeb.getUserDataByID(int32(28651800))         // 模擬第一步，取到使用者資料
	if vipuser != nil && AssignAmount > 0 && Complete == 0 { // 模擬先通過過澽修件
		rnd1 := int32(util.RndInt(0, 100)) // 用機率去控制放分，不要放分太快
		if rnd1 < int32(AssignRate) {
			fkeb.jiangtype = []int32{vipuser.userbet[0].fkebtype}
			rnd2 := int32(util.RndInt(0, 100))
			if rnd2 < int32(50) { // 抽中小於某特定值的數字，將要執行送分
				fkeb.majiangresult[vipuser.userbet[0].fkebtype].majiang1 = 10 // 指定牌型
				fkeb.majiangresult[vipuser.userbet[0].fkebtype].majiang2 = 10 // 指定牌型
				return true                                                   // 有成功修改牌型
			}
		}
	}
	return false // 無成功修改牌型
}

// 找尋 最大 、 最小 和 第二大 的開牌結果
func (fkeb *FKEBLogic) findMaxMinResultCard(MaxOrMin string) (int, bool) {

	// 先在這裡找到最大和最小的牌
	maxNumber := int32(-600)
	maxPosition := 5
	minNumber := int32(600)
	minPosition := 5
	secondNumber := int32(600)
	secondPosition := 5

	// 先在這裡找到最大和最小的牌
	for i := 0; i < 4; i++ {
		if fkeb.getMajiangsType(fkeb.majiangresult[i]) >= maxNumber {
			// secondPosition = maxPosition
			// 儲存最大的牌
			maxNumber = fkeb.getMajiangsType(fkeb.majiangresult[i])
			maxPosition = i
		}
		if fkeb.getMajiangsType(fkeb.majiangresult[i]) <= minNumber {
			// thirdPosition = minPosition
			// 儲存最小的牌
			minNumber = fkeb.getMajiangsType(fkeb.majiangresult[i])
			minPosition = i
		}
	}

	if MaxOrMin == "max" {
		return maxPosition, true // 回傳最大牌的結果
	} else if MaxOrMin == "min" {
		return minPosition, true // 回傳最小牌的結果
	} else if MaxOrMin == "second" {
		for j := 0; j < 4; j++ {
			if (j != maxPosition) && (j != minPosition) { // 對不是最大和最小的牌進行處理
				if secondPosition == 5 {
					secondNumber = fkeb.getMajiangsType(fkeb.majiangresult[j])
					secondPosition = j
				}
				if (secondNumber < fkeb.getMajiangsType(fkeb.majiangresult[j])) && secondNumber != 600 {
					secondNumber = fkeb.getMajiangsType(fkeb.majiangresult[j])
					secondPosition = j
				}
			}
		}
		return secondPosition, true
	} else {
		return 5, false
	}
}
