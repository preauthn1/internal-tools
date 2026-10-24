# Slot Basic Logic

Slot 遊戲邏輯，定義共用的參數與Function

### Config Datas

讀取 `config` 資料夾內的 `.dat` 檔案

#### `slotconfig.dat`

* SlotPatternID: `[slotsetting].slotpattern` 遊戲主設定
* BetCounts: `[slotsetting].betcount` 遊戲押注數量
* MaxBetCount: 最大押注數量
* BigWinRule: `[slotsetting].bigwinrule` client動畫播放階段,以Bet為主
* BigWinNotify: `[slotsetting].bigwinnotify` 當下注16, Total=6400(16*400) 得獎超過6400才廣播.
* EndActionTime: `[slotsetting].endactiontime` 結算動畫超時
* SymbolsDesc: `[symbols].{name}.title`, `[symbols].{name}.desc` 符號說明
* PayLineDesc: `[description].payline` 派彩線說明文字

#### `gameserver.dat`

* MinCarryMoney: `[gameserver].mincarrymoney` 最低攜帶金額

### Proto Datas

傳入 `slot.ProtoData`，依Config設定的SlotPatternID去抓取相對應的表

* SlotPattern: 遊戲主設定
* MainGameTable: 輪軸遊戲設定
* SymbolAttributes: 特殊符號設定
* SlotPayTable: 賠付表設定
* SlotSymbolDistribution: 符號機率設定
* SlotPayLine: 賠付線設定
* JackpotTable: 獎池設定
* RespinTable: 鎖軸重轉設定
* BonusTable: 獎勵遊戲設定
* PickGameTable: 選擇遊戲設定
