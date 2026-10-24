# Deployments

內有各個環境的設定內容與部署樣板

### Folder Structure

```
--- Servers
 | --- deployments
    | --- dd-dev                // Dev Env Settings
       | --- protodata          // Slot Proto Data for Dev
       | --- values             // Values Setting for Dev
    | --- dd-qa                 // QA Env
    | --- target-dev               // target Dev Env
    |
    | --- game-server           // Helm Chart
```

### `game-server` 為部署樣板不要動到

其他為各個環境的設定內容

# 環境設定

### `protodata`

Script會去掃描此資料夾內的`子資料夾`作為`{名稱}`建立資料

在設定slot遊戲時，需要指定此`{名稱}`，掛載到該遊戲的 `protodata` 資料夾

### `values`

* `db.yaml`: 設定 DB Server 相關參數使用
* `game-{short-name}`: 各遊戲設定內容
  
### Game Value Setting

以 `game` 開頭帶遊戲簡稱為各遊戲的yaml設定檔

> 編寫 yaml 檔 請注意縮排階層

開頭必定要在 `game:` 結構下，再來設定該遊戲簡稱

```
game:
  {short-name}:
    id: ??
    name: ??
    group: ??         # slot/poker/old
    ...
```

### `Game` Object

| Key   | Type      | Required | Default | Description          |
|:------|:---------:|:--------:|:-------:|:---------------------| 
| id    | int       | Y        |         | Game ID              |
| name  | str       | Y        |         | Game Name            |
| group | str       | Y        |         | slot/poker/old       |
| room  | array     | Y        |         | Game Rooms Array     |

> `group` 為 slot 才會設定 protodata

### `Room` Object

| Key       | Type      | Required | Default | Description                |
|:----------|:---------:|:--------:|:-------:|:---------------------------| 
| port      | int       | Y        |         | Game ID                    |
| serverid  | int       | Y        |         | Game Name                  |
| proto     | str       | OPT      |         | protodata's {name}         |
| configs   | objs      | Y        |         | Game Config files content  |

### `Config` Object

`key`直接帶原本在config資料夾內的檔案名稱，內容則為檔案內文字

#### 加上 `|` 後可下一行縮排後直接編輯檔案

```
...
configs:
  gameserver.dat: |
    [gameserver]
    gameid=??
    gamename=??
    ...
    ...
  slotconfig.dat: |
    ...
    ...
    ...
...
```