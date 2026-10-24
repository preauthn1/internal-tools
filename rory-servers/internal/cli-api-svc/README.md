# CLI API Server

將 Server 啟動 API 控制服務

基本現有 `Command` 的功能實作

透過此 API 來處理

* `GET` `/api/server/state` 查詢 Server 狀態
  * Server版號
  * Client版號
  * Log 狀態
  * 玩家資訊Summary
  * 機器人資訊Summary
* `GET` `/api/server/shotdown` 關閉 Server
* `GET` `/api/server/ver` 查詢Server版號
* `GET` `/api/server/cver` 查詢Client版號
* `PUT` `/api/server/cver/{platform}/{ver}` 更新Client版號
* `GET` `/api/log` 查詢 Log 狀態
  > 開關，Lite開關，Level
* `PUT` `/api/log` 開啟 Log
* `DELETE` `/api/log` 關閉 Log
* `PUT` `/api/log/lite` 開啟 Log Lite
* `DELETE ` `/api/log/lite` 關閉 Log Lite
* `PUT` `/api/log/level/{level}` 設定 Log Level [0:trace, 1: debug, 2: info, 3: warn, 4: error, 99: system]
* `GET` `/api/player` 列出玩家資訊Summary
* `GET` `/api/player/{playerID}` 查詢指定玩家狀態
* `DELETE` `/api/player/{playerID}` 踢出指定玩家
* `GET` `/api/player/list` 列出所有玩家資訊
* `GET` `/api/robot` 列出機器人資訊Summary
* `GET` `/api/slot/info` 查詢老虎機遊戲狀態
* `GET` `/api/slot/proto` 查詢老虎機proto