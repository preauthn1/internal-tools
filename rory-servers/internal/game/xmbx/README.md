# 熊貓報喜 (XMBX)

### Feature
* Lightning
* FreeSpin
* Jackpot

### 流程

``` mermaid
sequenceDiagram
    participant Client
    participant Game Server
    participant DB Server

    Client->>Game Server: UserEnter

    activate Game Server

        Game Server-->>Client: Game State [20002] 4
        Game Server->>DB Server: Find Data
        
        activate DB Server
            DB Server-->>Game Server: Jackpot
            DB Server-->>Game Server: FreeSpin
            DB Server-->>Game Server: Lightning
        deactivate DB Server

        Game Server-->>Client: Setting [20001]
        Game Server-->>Client: Jackpot [22004]
        
        alt Normal
            Game Server-->>Client: Game State [20002] 0
        else Has FreeSpin
            Game Server-->>Client: FreeSpin [20005]
            Game Server-->>Client: Game State [20002] 0
        else Has Lightning
            Game Server-->>Client: Lightning Prepare [10001]
            Game Server-->>Client: Game State [20002] 5
        end

    deactivate Game Server

    opt Normal Spin (State = 0)
        Client->>Game Server: Spin [21001]

        activate Game Server
            Note right of Game Server: Check User Money
            Game Server-->>Client: Bet State [20003] 0 / 1
            
            alt Bet Success
                Game Server->>DB Server: Jackpot Contribution

                Game Server->>DB Server: GameLog [MainGame(1)]

                Game Server-->>Client: Game Result [20004]
                Game Server->>DB Server: Update Score

                opt Win FreeSpin
                    Game Server->>DB Server: Save FreeSpin Data
                    Game Server-->>Client: FreeSpin [20005]
                end

                opt Win Lightning
                    Game Server->>DB Server: Save Lightning Data
                    Game Server-->>Client: Lightning Prepare [10001]
                end

                alt Has Lightning
                    Game Server-->>Client: Game State [20002] 5
                else
                    Game Server-->>Client: Game State [20002] 0
                end
            end
        deactivate Game Server
    end

    opt Free Spin (State = 0)
        Client->>Game Server: Spin [21001]

        activate Game Server
            Note right of Game Server: Check FreeSpin Times
            Game Server-->>Client: Bet State [20003] 0 / 1
            
            alt Bet Success
                Game Server->>DB Server: Update FreeSpin Times

                Game Server->>DB Server: GameLog [FreeGame(3)]

                Game Server-->>Client: Game Result [20004]
                Game Server->>DB Server: Update Score

                opt Win FreeSpin
                    Game Server->>DB Server: Update FreeSpin Times
                    Game Server-->>Client: FreeSpin [20005]
                end

                opt Win Lightning
                    Game Server->>DB Server: Save Lightning Data
                    Game Server-->>Client: Lightning Prepare [10001]
                end

                alt Has Lightning
                    Game Server-->>Client: Game State [20002] 5
                else
                    Game Server-->>Client: Game State [20002] 0
                end
            end
        deactivate Game Server
    end

    opt Lightning (State = 5)
        Client->>Game Server: Show Step [11001]

        activate Game Server
            Note right of Game Server: Check Lightning Exist

            Game Server->>DB Server: Update Lightning Show Step

            alt Has Unshowed Step
                Game Server->>DB Server: GameLog [Respin(6)]
                Game Server-->>Client: Lightning Prepare [10001]
                Game Server-->>Client: Game State [20002] 5
            else Is Last Step
                Game Server->>DB Server: GameLog [Lightning(7)]

                alt Win Jackpot
                    Game Server->>DB Server: Reset Jackpot
                end

                Game Server-->>Client: Lightning Result [10002]
                Game Server-->>Client: Game State [20002] 0
            end
        deactivate Game Server
    end
    
```