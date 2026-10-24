# 0001 改用 Go 服務＋Vue Vben Admin 後台，Swift 只包 Apple 本機模型

- 日期：2026-10-08
- 狀態：已採用
- 推翻：[`開發計畫.md`](../開發計畫.md) 2026-10-07 版的「程式結構：Swift Package（BotCore）＋伺服器」

## 背景

測試機上線前，需求多了兩項：AI 模型要能切換成 Claude Haiku API；機器人要能在後台設定（商家資訊、FAQ、預約開關、轉人工設定、意圖案例庫、待補問題）。後台給不懂技術的商家用，手機和電腦都要能開。

原本的 Swift 版（#4～#14）已完成 Telegram webhook、SQLite、店員群組、預約按鈕和轉人工。Apple 本機模型（Foundation Models）只能從 Swift 呼叫。

## 決定

- **Go 服務**負責所有非 AI 模型的部分：Telegram webhook、平台層、對話邏輯（顧客狀態、分類後的分派、預約狀態機、組 AI 指示、檢查回覆）、SQLite、店員群組、後台 JSON API。
- **後台前端用 Vue Vben Admin**，前後端分離，建置成靜態檔由 Go 提供；介面支援繁中、英、日。
- **Swift 縮成一支小服務**，只包 Apple 本機模型：收到指示和 JSON schema，回傳 JSON。Haiku 由 Go 直接呼叫。
- **語意向量由 Go 自己算**：用開源的多語言 embedding 模型（例如 hugot＋`multilingual-e5-small`），不靠 Apple 的 `NLContextualEmbedding`。要先做技術驗證（Apple Silicon 上能不能跑、分類分數）。
- 同一個 repo：新增 `go/`、`web/`，`Sources/` 縮成模型服務；`eval/` 題庫共用。Swift 版打 tag `swift-final` 當參考。
- **Swift 版停止開發**，測試機等 Go 版完成才上。Go 版的完成條件：題庫分數不低於 Swift 版，而且 #9～#14 的實機驗收項目都通過。

## 理由

- 新需求加上「AI 讀知識庫寫回覆」（[ADR 0002](0002-ai-writes-replies-from-knowledge.md)），BotCore 大部分本來就要重寫；對話邏輯分在兩種語言之間最難維護，所以整段搬到 Go。
- Swift 只留一定要用 Swift 的部分（Apple 模型）。之後如果量完決定只用 Haiku，這支服務可以整個拿掉。
- Apple 的 `NLContextualEmbedding` 是逐字向量取平均，只能比同一種文字系統；多語言句子向量可以讓中文案例直接比對到日文、英文的顧客訊息，案例不用翻譯。

## 代價

- #9～#14 要用 Go 重寫（資料表結構和流程可以照搬）。
- 題庫執行工具（`run_flows.py`、`loop/score.py`）要改成呼叫 Go 服務；改進迴圈暫停到 Go 版能跑題庫為止。
- 開發機要多裝 Go 和 node（建置 Vben）。
