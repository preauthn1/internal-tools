# 金吉報喜 (JJBX)

### Feature
* Lightning 模式
* FreeSpin 模式
* Jackpot 獎池
* JackpotGame 獎池遊戲
* Selection 選擇模式

## 流程

#### 登入 / 重連 
``` mermaid
sequenceDiagram
    participant Client
    participant Game Server
    participant DB Server

    Client->>Game Server: UserEnter / UserRecover

    activate Game Server

        Game Server-->>Client: Game State [20002] 4
        Game Server->>DB Server: Find Data
        
        activate DB Server
            DB Server-->>Game Server: Jackpot
            DB Server-->>Game Server: FreeSpin
            DB Server-->>Game Server: Lightning
            DB Server-->>Game Server: BonusGame
            DB Server-->>Game Server: SelectGame
        deactivate DB Server

        Game Server-->>Client: ReelSetting [20010]
        Game Server-->>Client: Setting [20001]
        Game Server-->>Client: Jackpot [22004]
        
        alt Normal
            Game Server-->>Client: Game State [20002] 0
        else Has FreeSpin
            Game Server-->>Client: FreeSpin [20005]
            Game Server-->>Client: Game State [20002] 0
        else Has Lightning
            Game Server-->>Client: Lightning Prepare [20006]
            Game Server-->>Client: Game State [20002] 5
        else Has SelectGame
            Game Server-->>Client: SelectGameOptions [20008]
            Game Server-->>Client: Game State [20002] 5
        else Has BonusGame
            Game Server-->>Client: BonusGame [20009]
            Game Server-->>Client: Game State [20002] 5
        end

    deactivate Game Server
```

#### 一般遊戲
``` mermaid
sequenceDiagram
    participant Client
    participant Game Server
    participant DB Server

    Client->>Game Server: Spin [21001]

    activate Game Server
        Note right of Game Server: Check Bet Illegal
        Note right of Game Server: Check User Money
        alt Bet Success
            Game Server-->>Client: Bet State [20003] 0
            
            Game Server->>DB Server: Jackpot Contribution

            Game Server->>DB Server: GameLog [MainGame(1)]

            Game Server-->>Client: Game Result [20004]
            Game Server->>DB Server: Update Score

            opt Win SelectGame
                Game Server->>DB Server: Save SelectGame Data
                Game Server-->>Client: SelectGameOptions [20008]
            end

            opt Win JackpotGame (BonusGame)
                Game Server->>DB Server: Save BonusGame Data
                Game Server-->>Client: BonusGame [20009]
            end
        else Bet Fail
            Game Server-->>Client: Bet State [20003] -1 ~ -6
        end
    deactivate Game Server
```

#### JackpotGame (BonusGame) 模式
``` mermaid
sequenceDiagram
    participant Client
    participant Game Server
    participant DB Server

    Client->>Game Server: PickItem [22005]

    activate Game Server
        Game Server->>DB Server: Update BonusGame Location
        
        alt Finished
            Game Server->>DB Server: GameLog [BonusGame]
            Game Server->>DB Server: Update Score
        end

        Game Server-->>Client: Confirm Location [22006]
    deactivate Game Server
```

#### SelectGame 模式
``` mermaid
sequenceDiagram
    participant Client
    participant Game Server
    participant DB Server

    Client->>Game Server: ChooseGame [21003]

    activate Game Server
        alt FreeSpin
            Game Server->>DB Server: Save FreeSpin Data
            Game Server-->>Client: FreeSpin [20005]
            Game Server-->>Client: Game State [20002] 0
        else Lightning
            Game Server->>DB Server: Save Lightning Data
            Game Server-->>Client: Lightning Prepare [20006]
            Game Server-->>Client: Game State [20002] 5
        end
    deactivate Game Server
```

#### FreeSpin 模式
``` mermaid
sequenceDiagram
    participant Client
    participant Game Server
    participant DB Server

    Client->>Game Server: Spin [21001]

    activate Game Server
        Note right of Game Server: Check Bet Illegal
        Note right of Game Server: Check FreeSpin Count
        
        Game Server-->>Client: Game Result [20004]
        Game Server->>DB Server: GameLog [FreeSpin]
        Game Server->>DB Server: Update Score

    deactivate Game Server
```

#### Lightning 模式

``` mermaid
sequenceDiagram
    participant Client
    participant Game Server
    participant DB Server

    Client->>Game Server: LightningSpin [21002]

    activate Game Server
        
        Game Server->>DB Server: Update Lightning Count -1
        
        alt Win Count
            Game Server->>DB Server: Update Lightning Count +1
            Game Server->>DB Server: GameLog [Respin]
        end

        alt Not Finished
            Game Server-->>Client: Lightning Prepare [20006]
            Game Server-->>Client: Game State [20002] 5
        else Finished
            Game Server-->>Client: Lightning Result [20007]
            Game Server-->>Client: Game State [20002] 0
        end

        
    deactivate Game Server
```

## Proto 對照
| Scenario         | BClassID | SClassID | Src -> Dest      | Message                     | Description           |
|:-----------------|:---------|:---------|:-----------------|:----------------------------|:----------------------|
| General          | 10000    | 20001    | Server -> Client | Slot_MSG_Setting            | 遊戲初始設定            |
|                  | 10000    | 20010    | Server -> Client | Slot_MSG_ReelSetting        | 遊戲狀態               |
|                  | 10000    | 22002    | Server -> Client | Slot_MSG_Status             | 排列帶內容             |
| Bet & Spin       | 10000    | 21001    | Client -> Server | Slot_MSG_Spin               | 下注轉                 |
|                  | 10000    | 20003    | Server -> Client | Slot_MSG_BetStatus          | 下注結果               |
|                  | 10000    | 20004    | Server -> Client | Slot_MSG_GameResult         | 遊戲結果               |
| Jackpot          | 10000    | 22004    | Server -> Client | Slot_MSG_JackpotPool        | Jackpot獎池設定        |
| FreeSpin         | 10000    | 20005    | Server -> Client | Slot_MSG_LeftFreeGameData   | 免費遊戲資訊            |
| Lightning        | 10000    | 20006    | Server -> Client | Slot_MSG_LightningPrepare   | Lightning遊戲準備資訊   |
|                  | 10000    | 20007    | Server -> Client | Slot_MSG_LightningResult    | Lightning遊戲結果      |
|                  | 10000    | 21002    | Client -> Server | Slot_MSG_LightningShow      | Client表演第幾階段      |
| SelectGame       | 10000    | 20008    | Server -> Client | Slot_MSG_ChooseGameOptions  | 選擇遊戲選項表          |
|                  | 10000    | 21003    | Client -> Server | Slot_MSG_ChooseGame         | 選擇的遊戲選項          |
| BonusGame        | 10000    | 20009    | Server -> Client | Slot_MSG_BonusGame          | Bonus(Jackpot)遊戲資訊 |
|                  | 10000    | 22005    | Client -> Server | Slot_MSG_PickBonus          | 玩家選擇Icon位置        |
|                  | 10000    | 22006    | Server -> Client | Slot_MSG_PickBonus          | Server收到選擇Icon位置  |