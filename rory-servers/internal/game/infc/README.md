# 印度切水果(INFC)

``` mermaid
sequenceDiagram
    participant Client
    participant Game Server

    Client->>Game Server: Connect
    Client->>Game Server: Login [5:1]
    Game Server-->>Client: Login Result [5:2]
    Game Server-->>Client: UserStateChange [10001:6] State: 0 (Free)
    Client->>Game Server: Queue [5:3]
    Game Server-->>Client: UserStateChange [10001:6] State: 1 (Queue)
    
    activate Game Server
        Game Server -->> Game Server: Match Player
        Game Server-->>Client: UserStateChange [10001:6] State: 3 (Ready)
    deactivate Game Server

    Game Server-->>Client: GameStart [10001:7]
    Game Server-->>Client: UserStateChange [10001:6] State: 4 (Game)
    Game Server-->>Client: INFC_GameSettingID [10001:10000] (INFC_GameSetting)
    Game Server-->>Client: INFC_BetStepID [10001:10010]
    Game Server-->>Client: INFC_GameStartID [10001:10020]

    activate Game Server
        loop Report Rank (1s)
        Game Server -->> Client: INFC_GameRankID [10001:10030] (INFC_GameRank)
        end

        alt Update Score
            Client ->> Game Server: INFC_IncScoreID [10001:30000] (INFC_IncScore)
        end
        alt Miss
            Client ->> Game Server: INFC_MissID [10001:30010] (一個封包表示Miss一次)
        end
        alt Boom
            Client ->> Game Server: INFC_BoomID [10001:30010] (切到炸彈))
        end
    deactivate Game Server

    Game Server-->>Client: INFC_GameResultID [10001:10040] (INFC_GameResult)
    Game Server-->>Client: GameEnd [10001:8]
    Game Server-->>Client: UserStateChange [10001:6] State: 2 (Sit)
    Game Server-->>Client: UserStateChange [10001:6] State: 0 (Free)
```