# Slot

Slot 遊戲的特色功能皆以企劃告知需要哪些而進行客制

* [遊戲流程](#遊戲流程)
* [元件](#元件)
* [設定配置](#設定配置)

### 遊戲流程

1. 玩家進入
   a. 幫玩家坐下
   b. 廣播進入封包
   c. 幫玩家準備
   d. 設定`遊戲狀態`為`GameStateWait`
   e. 各遊戲會去額外實作該遊戲有使用的`FeatureMode`到後端查詢是否有殘留資料(使用者資料復原)
      > 發送查詢後設定`sync.WaitGroup`，在接收到資訊的`func`裡執行`Done()`
   f. 發送遊戲設定相關資料
2. 等待玩家動作，依照各遊戲的流程優先順序，在各個遊戲的`logic.go`中的`CheckGameState()`，若有殘留特色遊戲，會先執行特色遊戲
   > * 一般狀態 `GameStateBet`: 玩家Spin (包含FreeSpin)
   > * 特色遊戲狀態 `GameStateUserAction`: 玩家有特色遊戲需要做操作，再依照不同的特色遊戲給封包
   >   * PickGame, JackpotGame: 都是翻牌遊戲，只是翻的內容不同
   >   * Selection: 選擇遊戲的機會，例如：東海(DHZW)，金吉(JJBX)，五龍相關遊戲(WL...)
   >   * Lightning: 進入特殊模式，鎖定符號重轉
   * Spin: 
     * 產生盤面
     * 計算盤面符合的規則
     * 計算規則對應的獎勵，並給獎 (給獎為特色遊戲的更新)
     * 寫分
     * SlotHistroyLog
     * PerGameLog
   * PickGame (`action-pick-game.go`) / JackpotGame (`action-pick-jackpot.go`):
     * 確認翻牌
     * 更新DB
     * 結束後：
       * 寫分
       * SlotHistroyLog
       * PerGameLog
       * 清除DB
       * 清除PickGame / JackpotGame
   * Selection: (`action-choose-game.go`)
     * 對照選擇內容，建立對應獎勵
     * 移除一次選擇機會
   * Lightning: (`action-lightning-spin.go`)
     * 產生盤面，繼承前次盤面金球
     * 更新DB
     * SlotHistroyLog
     * PerGameLog
     * 結束後：
       * 寫分
       * 清除 Lightning
      


### 元件
* [GameResultManager](#GameResultManager)
* [PaylineMatcher](#PaylineMatcher)
* [PayResult](#PayResult)
* [Jackpot](#Jackpot)
* [FreeSpin](#FreeSpin)
* [PickGame](#PickGame)
* [JackpotGame](#JackpotGame)
* [Selection](#Selection)
* [Lightning](#Lightning)
* [ExtraRule](#ExtraRule)
* [JPDefender](#http://example.invalid:1688/server-team/servers/wikis/老虎機/Jackpot-Defender)

#### GameResultManager
盤面與背景管理

##### 建立盤面 (`GenerateGameResult`)
產生指定群組ID的盤面
* 一個Reel一個Reel骰
* 先決定停在最中間的Row是哪個Order
* 再向上與向下骰出對應的Symbol

##### 建立背景 (`GenerateBackground`)(有背景的遊戲才會使用)
產生指定群組ID的背景
* 一個Reel一個Reel骰
* 骰背景是n格

#### PaylineMatcher
盤面比對器，依照機率表 `SlotPayTable` 的規則比對是否符合該規則

##### 全線賠付(PaylineAllMatcher)
##### 指定線賠付(PaylineFixedMatcher)

#### PayResult
賠付結果，紀錄是否透過百搭連線，與計算中獎金額，標記連線符號位置

#### Jackpot 
有設定Jackpot的遊戲都會配置 `map[JackpotLevel]*Jackpot`

由於早期沒有實作連線型Jackpot，都以單機型Jackpot製作

早期 Jackpot 規則是彩金會依照下注比例做加成
> 比如 總獎金為 100,000，最大下注金額為 100
> * 下注金額為1時: 贏得 100,000 * 1/100 = 10,00
> * 下注金額為10時: 贏得 100,000 * 10/100 = 10,000
> * 下注金額為100時: 贏得 100,000 * 100/100 = 100,000

後期依照金邊要求，Jackpot數字不能因下注倍率被改變，所以將 `proto` 資料的 `IsLinked` 作為判斷
* 連線: 金額不因下注倍率改變 (永遠為最大出獎值)
* 非連線: 金額會改變

#### FreeSpin

免費遊戲
> `Used` 欄位並未在DB開設欄位，所以無法回復

特殊規格 `FreeSpinExtAward`
> 這東西是[`JJBX`]所用到的，是在免費遊戲中的盤面，會因為黃球跟綠球出現的次數，額外給獎勵


#### PickGame

使用 `PickGameTable` 查詢該有哪些 `BonusTable` 項目

並依照遊戲建立所需翻牌次數的答案

之後存在 `BonusGameRecord` 這個表內

#### JackpotGame

將所有Jackpot項目選出一個作為獎項

產生每次翻牌的選項

之後存在 `BonusGameRecord` 這個表內

#### Selection

`select-chance.go`

選擇機會，可以累積次數，提供選項 (`SelectionTable`)

早期 `WLZB` 使用的是舊設定值，在 `proto-helper` 內的 `CompatiblePayTable` 將其整理為 `SelectionTable` 模式

#### Lightning

此模式會固定金球

每轉一次會是一個Step

初始建立時會將獲得Lightning的盤面設定為第一個Step

`IsLimit` 設定為true，就會依照臨界值(第幾次轉)，超過就換更換指定的`SymbolGroupID`來轉盤面

當玩家目前需要轉 `Lightning` 時

Server 會先傳送 `LightningPrepare` 封包，再傳送 `GameStateUserAction` 狀態

每次 Step 計算獎勵，依照 `SymbolAttributes` 判斷金球要給的獎勵為何

因為會計算獎勵，所以使用`slot/baselogic/game.go`的function，只能先擺在`baselogic`內，不然會造成 `package: import cycle not allowed`

#### ExtraRule
提供一些盤面異動規則

##### Condition
| Condition       | Condition Value Pattern                                     | Example        | Description     |
|:----------------|:------------------------------------------------------------|:---------------|-----------------|
| AfterSpin       | {LinkType}:{Symbol}:{Count}                                 | any:g1:15      | g1出現任意位置15次 |
| GetAward        | {Paytype}                                                   | lightning      | 贏得Lighting獎勵  |
| AppearSymbol    | {Symbol}\|{Any/Reel/Row}\|{Value1},{Value2}...              | w1\|Reel\|2,4  | w1出現在2或4Reel  |
| AppearSymbol    | {Symbol}\|{Any/Reel/Row}\|{Value1},{Value2}...              | w1\|Row\|2,3   | w1出現在2或3Row   |
| AppearSymbol    | {Symbol}\|{Any/Reel/Row}\|{Value1},{Value2}...              | w14\|Any       | w14出現任意位置    |
| AfterBackground | {LinkType}:{Symbol}:{Count}                                 | left:bg4:4     | 背景由左連續4個bg4 |

##### Action
| Action           | Action Value Pattern                                 | Example                      |
|:-----------------|:-----------------------------------------------------|:-----------------------------|
| Respin           | {機率}:{分母}                                         | 100:100000                    |
| RespinR          | {下注倍率}:{機率}:{分母},...                            | 8:10:1000,18:20:1000,...     |
| ConditionalAction| {Condition}:{Action}                                 | JPDefense:Respin              |
| UniGoldBall      | {Symbol}:{機率}:{分母}                                | g19:10:1000                   |
| RespinGroup      | {SymbolGroupID}                                      | 305                           |
| RespinGroup      | {SymbolGroupID}                                      | 308                            |
| SwitchReelRandom | {機率}:{分母}                                         | 100:100000                     |
| ExpandReel       | {Symbol}                                             | w14                            |
| ReplaceBG        | {Reel/Row}:{val}:{BGName}:{機率}:{分母}                | Reel:5:bg4:100:1000000        |
| ReplaceBGT       | {Reel/Row}:{val}:{BGName}:{機率}:{分母}                | Reel:5:bg4:100:1000000         |
| ReplaceBGR       | {下注倍率}:{Reel/Row}:{val}:{BGName}:{機率}:{分母},...  | 28:Reel:5:bg4:100:1000000,...  |
| ReplaceBGRT      | {下注倍率}:{Reel/Row}:{val}:{BGName}:{機率}:{分母},...  | 28:Reel:5:bg4:100:1000000,...  |

> `Respin`: `{機率}`為放行的機率(不會Respin)
> 
> `RespinR`: `{機率}`要乘以TotalBet(BetMoney * BetRate)



### 設定配置

#### Config
提供 `slotconfig.dat` 檔案

```
[slotsetting]
betcount=1,2,5,10,20    // 可下注倍率
betrate=40              

slotpattern=10003       // 指定哪一套 Pattern
bigwinrule=30000,40000,50000,500000
bigwinnotify=400

[description]
```

> `betrate` 比較特殊，
> * 指定線賠付: 多少條線就是要在乘多少倍
> * 全線賠付: 部分遊戲玩家可以自己選擇(`DFDC`, `S88F`...) 其餘都設定為1
