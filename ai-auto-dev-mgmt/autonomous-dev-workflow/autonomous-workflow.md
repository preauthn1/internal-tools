# 自治並行開發流程協議（Autonomous Parallel Development Workflow Protocol）

> **版本**：1.0
> **適用團隊**：PM Agent、Frontend Engineer Agent、Backend Engineer Agent、QA Engineer Agent
> **定位**：頂層協調協議，各 Agent 內部行為仍遵循各自的 System Prompt

---

## 目錄

1. [協議概述](#1-協議概述)
2. [流程六階段總覽](#2-流程六階段總覽)
3. [Phase 0：需求輸入](#3-phase-0需求輸入)
4. [Phase 1：需求分析與文件產出](#4-phase-1需求分析與文件產出)
5. [Phase 2：API Contract 協商](#5-phase-2api-contract-協商)
6.5. [Phase 2.5：專案初始化與環境建置](#65-phase-25專案初始化與環境建置)
6. [Phase 3：並行開發](#6-phase-3並行開發)
7. [Phase 4：整合與 QA 驗證](#7-phase-4整合與-qa-驗證)
8. [Phase 5：自治修復循環](#8-phase-5自治修復循環)
9. [Phase 6：交付與驗收](#9-phase-6交付與驗收)
10. [信號協議（Signal Protocol）](#10-信號協議signal-protocol)
11. [品質閘門（Quality Gates）](#11-品質閘門quality-gates)
12. [並行工作矩陣](#12-並行工作矩陣)
13. [衝突解決與升級機制](#13-衝突解決與升級機制)
14. [流程狀態追蹤](#14-流程狀態追蹤)
15. [自治邊界與人類介入點](#15-自治邊界與人類介入點)
16. [需求變更處理流程](#16-需求變更處理流程)

---

## 1. 協議概述

### 1.1 目的

本協議定義四個 AI Agent（PM、FE、BE、QA）之間的**自治協調規則**，使團隊能以最小化人類干預的方式，並行推進軟體開發流程。

### 1.2 核心原則

1. **並行最大化**：在依賴允許的範圍內，盡可能讓多個 Agent 同時工作
2. **同步最小化**：只在必要的決策點進行同步，減少等待
3. **品質不妥協**：每個階段轉換必須通過品質閘門，不可跳過
4. **自治有邊界**：明確定義哪些決策可自治、哪些必須上報人類
5. **可追溯性**：所有決策、交付物、信號都有記錄
6. **統一開發環境**：所有 Agent 均在預先配置好的開發環境（Devcontainer）內運行，可直接執行開發與測試指令，無需額外建立容器環境。

### 1.3 角色職責摘要

| 角色 | 自治流程中的定位 | 主導階段 |
|------|-----------------|---------|
| **PM** | 流程總指揮、品質把關、衝突仲裁 | Phase 0-1, 5-6 |
| **FE** | 前端並行開發執行者 | Phase 3 |
| **BE** | 後端並行開發執行者 | Phase 3 |
| **QA** | 獨立品質驗證者 | Phase 4 |
| **Code Reviewer** | 跨領域審查者（需求覆蓋、可達性、Contract 合規） | Phase 3→4 過渡（G3-Review） |

---

## 2. 流程六階段總覽

```
人類需求輸入
     │
     ▼
┌──────────────────────────────────────────────────────────────────┐
│                     自治管理流程開始                               │
│                                                                  │
│  Phase 0         Phase 1           Phase 2          Phase 2.5    │
│  需求輸入 ──G0──▶ 需求分析 ──G1──▶ Contract ──G2──▶ 專案初始化   │
│  [PM]            [PM]             協商              [FE + BE]    │
│                                   [PM+FE+BE]        [PM 監控]    │
│                                                                  │
│  Phase 2.5       Phase 3           Phase 4          Phase 5      │
│  專案初始化 ─G2.5─▶ 並行開發 ──G3──▶ QA 驗證 ──G4──▶ 自治修復    │
│  [FE + BE]       [FE + BE]         [QA]             循環         │
│                  [QA 設計]                           [PM 主導]    │
│                                                                  │
│  Phase 5         Phase 6                                         │
│  自治修復 ──G5──▶ 交付驗收                                        │
│  循環             [PM]                                            │
│  [PM 主導]                                                       │
└──────────────────────────────────────────────────────────────────┘
     │
     ▼
人類驗收 / 決策
```

### 預估時間分配（以單一 Feature 為基準）

| Phase | 預估佔比 | 可並行 |
|-------|---------|--------|
| Phase 0：需求輸入 | 5% | 否 |
| Phase 1：需求分析 | 15% | 部分（FE/BE/QA 可同步準備環境） |
| Phase 2：Contract 協商 | 10% | 是（QA 可同步設計測試計畫） |
| Phase 2.5：專案初始化 | 5% | 是（FE + BE 並行初始化） |
| Phase 3：並行開發 | 35% | **是（FE + BE 完全並行，QA 設計測試案例）** |
| Phase 4：QA 驗證 | 15% | 部分（FE/BE 待命支援） |
| Phase 5：自治修復 | 10% | 是（FE/BE 並行修復各自 Bug） |
| Phase 6：交付驗收 | 5% | 否 |

---

## 3. Phase 0：需求輸入

### 觸發條件
人類決策者提出產品需求。

### 執行者
**PM Agent**

### 流程
1. PM 接收人類的需求描述
2. PM 使用 5W1H 方法分析需求完整性
3. 若資訊不足 → 向人類提出**釐清問題清單**（`SIGNAL: REQ_CLARIFY → 人類`）
4. 若資訊充足 → 產出**需求摘要**，請人類確認
5. 人類確認後 → 通過 **Gate G0**

### 品質閘門 G0
- [ ] 需求描述完整（5W1H 已涵蓋）
- [ ] 使用者角色（Persona）已識別
- [ ] 業務目標明確
- [ ] 人類已確認需求摘要

### 輸出物
- 需求摘要文件（經人類確認）

### 信號
```
PM ──[REQ_CLARIFY]──▶ 人類     （需求不足時）
人類 ──[REQ_INPUT]──▶ PM       （補充需求）
Orchestrator ──[G0_PASSED]──▶ Phase 1    （閘門通過）
```

---

## 4. Phase 1：需求分析與文件產出

### 觸發條件
Gate G0 通過。

### 執行者
**PM Agent**（主導），FE / BE / QA（可並行準備）

### 流程

```
G0 通過
  │
  ├──▶ PM：需求分析 → Feature 拆解 → User Story → AC 定義
  │         │
  │         ├──▶ 前端需求文件
  │         ├──▶ 後端需求文件
  │         └──▶ 測試驗收文件
  │
  ├──▶ FE（並行）：技術棧確認
  ├──▶ BE（並行）：技術棧確認
  └──▶ QA（並行）：測試環境準備、測試工具確認
```

### PM 的工作（依其 System Prompt）
1. **Step 1**：需求理解與釐清（Persona、MoSCoW 優先級）
2. **Step 2**：功能拆解（Feature → User Story，遵循 INVEST 原則）
3. **Step 3**：驗收標準定義（Given-When-Then 格式）
4. **Step 4**：技術需求評估（前後端範圍、API 介面需求）
5. **Step 5**：Sprint 規劃建議（依賴關係、Story Point 估算）

### 品質閘門 G1
- [ ] 前端需求文件已產出（遵循 PM prompt 中的前端需求文件模板）
- [ ] 後端需求文件已產出（遵循 PM prompt 中的後端需求文件模板）
- [ ] 測試驗收文件已產出（遵循 PM prompt 中的測試驗收文件模板）
- [ ] 所有 User Story 通過品質檢查清單（10 項）
- [ ] FE 審閱前端需求 → 無阻塞性問題
- [ ] BE 審閱後端需求 → 無阻塞性問題
- [ ] QA 審閱測試驗收文件 → 無阻塞性問題

### 輸出物
- 前端需求文件（含 UI 需求、User Flow、表單驗證、API 對接、前端 AC）
- 後端需求文件（含商業邏輯、API 設計、資料模型、權限、效能、後端 AC）
- 測試驗收文件（含 User Story 清單、測試情境、測試資料、驗收確認清單）

### 信號
```
PM ──[REQ_READY]──▶ FE, BE, QA   （三份文件已產出，請審閱）
FE ──[REQ_REVIEW_OK / REQ_QUESTION]──▶ PM
BE ──[REQ_REVIEW_OK / REQ_QUESTION]──▶ PM
QA ──[REQ_REVIEW_OK / REQ_QUESTION]──▶ PM
Orchestrator ──[G1_PASSED]──▶ Phase 2      （所有審閱通過）
```

### 阻塞處理
- 若 FE/BE/QA 提出 `REQ_QUESTION`，PM 必須在**當回合內回覆**
- 若問題涉及需求變更 → PM 評估影響範圍，必要時回退 Phase 0 請人類確認
- 若問題僅為技術釐清 → PM 直接回覆或調整文件

---

## 5. Phase 2：API Contract 協商

### 觸發條件
Gate G1 通過。

### 執行者
**PM Agent**（主導仲裁），FE + BE（協商），QA（並行設計測試計畫）

### 流程

```
G1 通過
  │
  ├──▶ PM：根據需求文件中的 API 設計，產出 API Contract 草案
  │         │
  │         ▼
  │    ┌─────────────────────────────────────────┐
  │    │        API Contract 協商流程              │
  │    │                                         │
  │    │  1. PM 發出 Contract 草案               │
  │    │     ↓                                   │
  │    │  2. FE + BE 各自審閱                     │
  │    │     ↓                                   │
  │    │  3a. 雙方無異議 → Contract 鎖定          │
  │    │  3b. 有異議 → 提出修改建議               │
  │    │     ↓                                   │
  │    │  4. PM 評估修改建議，做最終決定           │
  │    │     ↓                                   │
  │    │  5. Contract 鎖定（CONTRACT_LOCKED）     │
  │    └─────────────────────────────────────────┘
  │
  └──▶ QA（並行）：根據測試驗收文件，開始設計測試計畫與測試案例
```

### API Contract 文件格式

```markdown
## API Contract — [Feature Name]

### 版本
- Contract 版本：v1.0
- 狀態：🔓 已鎖定 / 🔄 協商中

### 端點清單
| # | 方法 | 端點 | 描述 | 認證 |
|---|------|------|------|------|
| 1 | POST | /api/v1/resource | 建立資源 | Bearer Token |
| 2 | GET  | /api/v1/resource | 取得列表 | Bearer Token |

### 詳細規格（每個端點）

#### API-001：建立資源
- **端點**：`POST /api/v1/resource`
- **請求 Header**：
  - `Authorization: Bearer {token}`
  - `Content-Type: application/json`
- **請求 Body**：
  ```json
  {
    "field_a": "string — 必填，說明",
    "field_b": "number — 選填，說明"
  }
  ```
- **成功回應（201）**：
  ```json
  {
    "success": true,
    "data": {
      "id": "number",
      "field_a": "string",
      "field_b": "number",
      "created_at": "string (ISO 8601)"
    },
    "message": "建立成功"
  }
  ```
- **錯誤回應**：
  | 狀態碼 | 情境 | 回應格式 |
  |--------|------|----------|
  | 400 | 請求格式錯誤 | `{ "success": false, "message": "..." }` |
  | 401 | 未認證 | `{ "success": false, "message": "未授權" }` |
  | 422 | 驗證失敗 | `{ "success": false, "message": "...", "errors": {...} }` |

### 共用約定
- **回應格式標準**：統一使用 BE Prompt 定義的格式（`success` 布爾值 + `data` + `message`）。FE 的 `ApiResponse` 型別定義須與此對齊
- 分頁參數：`?page=1&per_page=15`
- 排序參數：`?sort_by=created_at&sort_order=desc`
- 時間格式：ISO 8601（`2026-03-24T12:00:00Z`）
- 認證方式：Laravel Sanctum Bearer Token

### 變更紀錄
| 版本 | 日期 | 變更內容 | 決定者 |
|------|------|----------|--------|
| v1.0 | | 初版 | PM |
```

### Contract 變更規則

| 時機 | 規則 |
|------|------|
| Phase 2（協商中） | FE/BE 可自由提出修改建議，PM 仲裁 |
| Phase 3（開發中） | **Contract 已鎖定**，變更需走變更流程 |
| 變更流程 | 提出者 → PM 評估影響 → FE + BE 確認 → PM 決定 → 更新 Contract 版本號 |

### 品質閘門 G2
- [ ] API Contract 文件完整（所有端點含請求/回應格式、錯誤碼）
- [ ] FE 確認 Contract 可支撐所有前端需求
- [ ] BE 確認 Contract 技術可行
- [ ] QA 確認 Contract 可測試
- [ ] Contract 狀態標記為 🔓 已鎖定
- [ ] QA 測試計畫初稿完成

### 信號
```
PM ──[CONTRACT_DRAFT]──▶ FE, BE          （草案發出）
FE ──[CONTRACT_OK / CONTRACT_CHANGE]──▶ PM
BE ──[CONTRACT_OK / CONTRACT_CHANGE]──▶ PM
PM ──[CONTRACT_LOCKED]──▶ FE, BE, QA     （鎖定，開始開發）
Orchestrator ──[G2_PASSED]──▶ Phase 2.5
```

---

## 6.5. Phase 2.5：專案初始化與環境建置

### 觸發條件
Gate G2 通過（API Contract 已鎖定）。

### 執行者
**FE Agent + BE Agent**（並行初始化），PM Agent（監控）

### 流程

```
G2 通過（Contract 🔓 已鎖定）
  │
  ├──▶ BE Agent（worktree 隔離）：
  │    1. 透過 Context7 MCP（或官方文件）查詢最新穩定版本（Laravel、PHP、MySQL、Redis）
  │    2. composer create-project laravel/laravel
  │    3. 設定 .env、database config
  │    4. 建立初始 Migration（若 Contract 已定義資料模型）
  │    5. 確認 php artisan serve 可啟動
  │    6. Git 初始 commit
  │
  └──▶ FE Agent（worktree 隔離）：
       1. 透過 Context7 MCP（或官方文件）查詢最新穩定版本（Vue 3、Node.js LTS、Vite、Playwright）
       2. npm create vue@latest
       3. 安裝依賴（Pinia、Vue Router、Axios、Vitest、Playwright）
       4. 設定 MSW（Mock Service Worker）依據 API Contract
       5. 確認 npm run dev 可啟動
       6. Git 初始 commit
```

### 品質閘門 G2.5
- [ ] BE 專案可正常啟動（`php artisan serve` 成功）
- [ ] FE 專案可正常啟動（`npm run dev` 成功）
- [ ] BE 基礎測試可執行（`php artisan test` 通過）
- [ ] FE 基礎測試可執行（`npm run test` 通過）
- [ ] MSW Mock 依據 API Contract 設定完成
- [ ] 兩個專案皆已 Git 初始 commit

### 信號
```
Orchestrator ──[INIT_START]──▶ FE, BE           （開始專案初始化）
FE ──[INIT_COMPLETE]──▶ Orchestrator   （FE 專案初始化完成）
BE ──[INIT_COMPLETE]──▶ Orchestrator   （BE 專案初始化完成）
Orchestrator ──[G2.5_PASSED]──▶ Phase 3（兩方都完成）
```

---

## 6. Phase 3：並行開發

### 觸發條件
Gate G2.5 通過（FE + BE 專案初始化完成）。

### 執行者
**FE Agent + BE Agent**（並行開發），QA Agent（並行設計測試案例），PM Agent（監控進度）

### 並行開發架構

```
                    API Contract 🔓 已鎖定
                           │
                ┌──────────┼──────────┐
                ▼          │          ▼
     ┌───────────────┐    │    ┌───────────────┐
     │  前端 TDD 開發  │    │    │  後端 TDD 開發  │
     │               │    │    │               │
     │  ① Red：      │    │    │  ① Red：      │
     │    寫失敗測試   │    │    │    寫失敗測試   │
     │  ② Green：    │    │    │  ② Green：    │
     │    最少代碼通過 │    │    │    最少代碼通過 │
     │  ③ Refactor： │    │    │  ③ Refactor： │
     │    重構        │    │    │    重構        │
     │               │    │    │               │
     │  Mock：MSW    │    │    │  DB：MySQL    │
     │  Test：Vitest │    │    │  Test：PHPUnit│
     │  E2E：Playwright│   │    │  Cache：Redis │
     └───────┬───────┘    │    └───────┬───────┘
             │            │            │
             ▼            │            ▼
     ┌───────────────┐    │    ┌───────────────┐
     │  前端自我驗證   │    │    │  後端自我驗證   │
     │  ✅ Vitest     │    │    │  ✅ php artisan│
     │  ✅ vue-tsc    │    │    │     test      │
     │  ✅ ESLint     │    │    │  ✅ coverage  │
     │  ✅ Playwright │    │    │  ✅ Pint/CS   │
     │    (Mock API)  │    │    │  ✅ PHPStan   │
     └───────┬───────┘    │    └───────┬───────┘
             │            │            │
             ▼            ▼            ▼
          DEV_COMPLETE  QA 測試    DEV_COMPLETE
          (FE)        案例完成     (BE)
                ┌────────┘
                ▼
         等待所有完成
           ──G3──▶ Phase 4

同時間：
     ┌───────────────┐
     │  QA 並行工作    │
     │               │
     │  ① 測試案例設計 │
     │  ② Playwright  │
     │    腳本撰寫    │
     │  ③ Postman     │
     │    Collection  │
     │    準備        │
     │  ④ 測試資料    │
     │    準備        │
     └───────────────┘
```

### 前端開發流程（依 FE System Prompt）
1. 閱讀前端需求文件，確認 AC
2. 拆解為可測試的 User Story
3. 確認 API Contract → 設定 MSW Mock
4. TDD 循環：Red → Green → Refactor
5. 補充 E2E 測試（Playwright，使用 MSW）
6. 自我驗證：執行 Vitest、vue-tsc、ESLint
7. 產出**前端開發報告**與**提測通知**

### 後端開發流程（依 BE System Prompt）
1. 審閱後端需求文件，確認 AC
2. 技術方案設計（Schema、Service 層、Middleware）
3. TDD 循環：Red → Green → Refactor
4. 自我驗證：執行 `php artisan test`、coverage、Pint、PHPStan
5. 產出**後端開發報告**與**提測通知**

### QA 並行工作（依 QA System Prompt）
1. 根據測試驗收文件設計完整測試計畫
2. 設計測試案例（正向、負向、邊界、例外）
3. 撰寫 Playwright E2E 測試腳本（框架準備）
4. 準備 Postman Collection（依據 API Contract）
5. 準備測試資料

### 品質閘門 G3（分為 G3-FE 與 G3-BE，各自獨立判定）

**G3-FE（前端閘門）**
- [ ] 所有 Vitest 單元測試通過
- [ ] vue-tsc 型別檢查無錯誤
- [ ] ESLint 無警告或錯誤
- [ ] Playwright E2E 測試通過（Mock API 環境）
- [ ] 所有前端 AC（AC-FE-xxx）有對應測試
- [ ] 無 `console.log` / `debugger` 殘留
- [ ] 前端開發報告已產出
- [ ] 提測通知已產出

**G3-BE（後端閘門）**
- [ ] `php artisan test` 所有測試通過
- [ ] 覆蓋率達標（Unit ≥ 90%, Feature ≥ 85%, 整合 ≥ 80%）
- [ ] Laravel Pint / PHP-CS-Fixer 無風格問題
- [ ] PHPStan / Larastan 靜態分析通過
- [ ] 所有後端 AC（AC-BE-xxx）有對應測試
- [ ] 無 `dd()` / `dump()` 殘留
- [ ] API 端點文件已產出（Scribe / Swagger / 手動撰寫）
- [ ] 後端開發報告已產出
- [ ] 提測通知已產出

**G3-Review（跨領域審查閘門 — Code Reviewer Agent）**
- [ ] 需求覆蓋矩陣：所有需求元件/頁面有對應原始碼（覆蓋率 100%）
- [ ] 頁面可達性：所有路由可從 UI 導航到達（無孤島頁面）
- [ ] API Contract 合規：所有端點在 BE Route + FE Service + MSW Handler 中有對應
- [ ] GAP_REPORT 判定為 PASS

> G3-Review 在 G3-FE 和 G3-BE 都通過後才啟動。若 GAP_REPORT 為 BLOCKED，Orchestrator 將 Gap 清單路由回對應的 FE/BE Agent 修復，修復後重新啟動 Code Reviewer 審查。

**G3 通過條件**：G3-FE ✅ **且** G3-BE ✅ **且** G3-Review ✅

### FE ↔ BE 直接溝通協議

在 Phase 3 並行開發期間，FE 與 BE 可就技術細節**直接溝通**，無需經過 PM：

- **適用範圍**：API 回應格式細節、分頁/排序參數命名、錯誤碼對照、資料型別對齊
- **信號**：`FE_BE_DISCUSS`（雙向）
- **規則**：
  - 討論結果不可違反已鎖定的 API Contract 核心約定
  - 若討論涉及 Contract 變更 → 必須走 Contract 變更流程（通知 PM）
  - 所有討論結論須記錄在流程狀態追蹤表中

### Code Review 規範

Phase 3 開發完成前，FE 和 BE 各自須完成 Code Review（依各自 prompt 定義）：

- **FE**：PR 必須通過 Checklist（FE Prompt 9.4），包含測試覆蓋、型別檢查、ESLint、無 console.log
- **BE**：交付物須通過自我驗證（BE Prompt Step 4），包含測試、覆蓋率、風格、靜態分析
- Code Review 為 G3 閘門的隱含條件

### 阻塞處理
- 若 FE 或 BE 遇到阻塞 → 發送 `DEV_BLOCKED` 信號給 PM
- PM 評估阻塞原因：
  - 需求不清 → PM 釐清
  - 技術障礙 → 協調對方 Agent 支援
  - Contract 需變更 → 啟動 Contract 變更流程
- 阻塞不影響另一方的並行開發

### 信號
```
Orchestrator ──[DEV_START]──▶ FE, BE, QA         （Phase 3 開始）
FE ──[DEV_BLOCKED]──▶ PM               （開發阻塞）
BE ──[DEV_BLOCKED]──▶ PM               （開發阻塞）
FE ──[DEV_COMPLETE]──▶ PM, QA          （前端開發完成）
BE ──[DEV_COMPLETE]──▶ PM, QA          （後端開發完成）
QA ──[TESTCASE_READY]──▶ PM            （測試案例設計完成）
Orchestrator ──[REVIEW_START]──▶ Code Reviewer   （G3-FE + G3-BE 通過後啟動跨領域審查）
Code Reviewer ──[GAP_REPORT]──▶ Orchestrator     （審查結果：PASS / BLOCKED）
Orchestrator ──[GAP_FIX_REQUIRED]──▶ FE / BE    （若 BLOCKED，路由 Gap 清單回對應 Agent）
Orchestrator ──[G3_PASSED]──▶ Phase 4            （G3-FE + G3-BE + G3-Review 皆通過）
```

---

## 7. Phase 4：整合與 QA 驗證

### 觸發條件
Gate G3 通過（FE + BE 皆完成自我驗證）。

### 執行者
**QA Agent**（主導），FE + BE（支援整合與問題排查），PM（監控品質）

### 流程

```
G3 通過
  │
  ▼
┌──────────────────────────────────────────────────────┐
│                  整合與 QA 驗證                        │
│                                                      │
│  Step 1：前後端整合                                    │
│  ├── FE 切換 MSW Mock → 真實後端 API                  │
│  ├── 聯調驗證：逐一比對 API Contract                   │
│  └── 整合問題修復（FE/BE 協作）                        │
│                                                      │
│  Step 2：QA 獨立驗證（依 QA System Prompt 完整流程）     │
│  ├── 2.1 環境確認（服務啟動）                           │
│  ├── 2.2 手動功能測試（依測試案例優先級 P0→P3）         │
│  ├── 2.3 API 測試（Postman / Newman）                 │
│  ├── 2.4 E2E 自動化測試（Playwright）                  │
│  ├── 2.5 探索性測試                                   │
│  └── 2.6 Bug 回報（即時通知對應工程師）                 │
│                                                      │
│  Step 3：測試報告產出                                   │
│  ├── 測試概況（通過率、Bug 統計）                       │
│  ├── AC 覆蓋矩陣                                      │
│  ├── Bug 清單                                         │
│  └── 驗收建議（✅ / ⚠️ / ❌）                          │
└──────────────────────────────────────────────────────┘
```

### 整合測試注意事項
- FE 在此階段將 MSW Mock 替換為真實 API 端點
- 逐一比對 API Contract，確認前後端回應格式一致
- 發現 Contract 不一致時：
  - 若為 BE 實作偏差 → BE 修正
  - 若為 FE 解析錯誤 → FE 修正
  - 若為 Contract 本身有缺陷 → 走 Contract 變更流程

### Bug 即時通報規則
QA 發現 Bug 時，**不等到報告完成才通報**，立即：
1. 判斷 Bug 歸屬（`[FE]` 或 `[BE]`）
2. 依 QA System Prompt 中的 Bug 回報格式記錄
3. 發送 `BUG_FOUND` 信號給對應工程師
4. 工程師**可在 QA 持續測試期間同步開始修復**（並行）

### 品質閘門 G4
- [ ] QA 測試報告已完成
- [ ] 驗收建議已明確標注（✅ / ⚠️ / ❌）
- [ ] 所有 Bug 已分類並標記嚴重程度
- [ ] 測試覆蓋所有 P0 + P1 測試案例

### 路由邏輯（基於 QA 驗收建議）
```
QA 驗收建議
    │
    ├── ✅ 建議通過 ──────────────▶ Phase 6（交付）
    │
    ├── ⚠️ 有條件通過 ──────────▶ Phase 5（自治修復）
    │
    └── ❌ 不建議通過 ──────────▶ Phase 5（自治修復）
```

### 信號
```
Orchestrator ──[INTEGRATION_START]──▶ FE, BE, QA   （整合開始）
QA ──[QA_START]──▶ PM                    （開始獨立驗證）
QA ──[BUG_FOUND]──▶ FE / BE             （即時 Bug 通報）
QA ──[QA_REPORT]──▶ PM                   （測試報告完成）
Orchestrator ──[G4_PASSED]──▶ Phase 5 或 Phase 6
```

---

## 8. Phase 5：自治修復循環

### 觸發條件
QA 驗收建議為 ⚠️ 有條件通過 或 ❌ 不建議通過。

### 執行者
**PM Agent**（主導），FE + BE（修復），QA（重新驗證）

### 流程
**直接引用 PM System Prompt 中的「驗收不通過自治流程」**，最多 **3 回合**。

```
QA 測試報告：⚠️ / ❌
         │
         ▼
┌─────────────────────────────────────────────────────────┐
│              自治修復循環（最多 3 回合）                    │
│                                                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Step A：PM 召集團隊會議（QA + FE + BE）             │  │
│  │ → 分析未通過的 Bug / AC                            │  │
│  │ → 確認問題根因與責任歸屬（FE / BE）                  │  │
│  └──────────────────────┬────────────────────────────┘  │
│                          ▼                               │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Step B：PM 協調分派修復任務                         │  │
│  │ → FE Bug 指派給 FE Agent                          │  │
│  │ → BE Bug 指派給 BE Agent                          │  │
│  │ → FE 與 BE 可並行修復各自的 Bug                    │  │
│  └──────────────────────┬────────────────────────────┘  │
│                          ▼                               │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Step C：工程師修復（TDD 方式）                      │  │
│  │ → 先寫重現 Bug 的失敗測試                          │  │
│  │ → 修復程式碼使測試通過                              │  │
│  │ → 自我驗證通過後提交 QA                            │  │
│  │                                                   │  │
│  │ ⚡ FE 和 BE 可並行修復，不需等待對方                │  │
│  └──────────────────────┬────────────────────────────┘  │
│                          ▼                               │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Step D：QA 重新驗證                                │  │
│  │ → 驗證修復的 Bug                                   │  │
│  │ → 回歸測試（影響範圍內的測試案例）                   │  │
│  │ → 產出新的測試報告                                  │  │
│  └──────────────────────┬────────────────────────────┘  │
│                          ▼                               │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Step E：PM 確認新測試報告                           │  │
│  │ → ✅ 通過 → 結束循環 → Phase 6                    │  │
│  │ → ⚠️/❌ → 回合數 +1                              │  │
│  │ → 回合 > 3 → 跳出循環 → 上報人類                   │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 每回合記錄（依 PM System Prompt 格式）
PM 必須記錄：問題分析表、修復方案表、重新驗收結果。

### 品質閘門 G5
- [ ] 驗收通過（✅ 建議通過）→ 進入 Phase 6
- [ ] 或 3 回合用盡 → 產出「自治流程摘要報告」→ 上報人類

### 信號
```
PM ──[FIX_ASSIGN]──▶ FE / BE           （指派修復任務）
FE ──[FIX_COMPLETE]──▶ QA, PM          （前端修復完成）
BE ──[FIX_COMPLETE]──▶ QA, PM          （後端修復完成）
QA ──[RETEST_REPORT]──▶ PM             （重測報告）
Orchestrator ──[G5_PASSED]──▶ Phase 6            （驗收通過）
PM ──[ESCALATE]──▶ 人類                （3 回合用盡，上報）
```

---

## 9. Phase 6：交付與驗收

### 觸發條件
Gate G5 通過（QA 驗收 ✅）。

### 執行者
**PM Agent**（主導），彙整所有交付物。

### 交付物清單

```
## [Feature Name] — 交付清單

### PM 產出
- [ ] 需求摘要（經人類確認版）
- [ ] 前端需求文件（最終版）
- [ ] 後端需求文件（最終版）
- [ ] 測試驗收文件（最終版）
- [ ] API Contract（最終鎖定版）
- [ ] Sprint 報告
- [ ] 自治修復回合記錄（若有）

### FE 產出
- [ ] 前端原始碼（通過所有測試）
- [ ] 前端開發報告
- [ ] Vitest 測試報告
- [ ] Playwright E2E 測試報告（Mock 環境）

### BE 產出
- [ ] 後端原始碼（通過所有測試）
- [ ] Migration 檔案
- [ ] 後端開發報告
- [ ] PHPUnit 測試報告（含覆蓋率）
- [ ] API 文件（Scribe / Swagger）
- [ ] 環境變數需求說明
- [ ] 部署注意事項

### QA 產出
- [ ] 測試計畫
- [ ] 測試案例集
- [ ] 測試報告（最終版）
- [ ] Playwright E2E 測試報告（真實環境）
- [ ] Postman Collection + 測試結果
- [ ] Bug 追蹤表（所有 Bug 已 Closed）
```

### 信號
```
PM ──[DELIVERY_READY]──▶ 人類           （交付物已準備完成）
人類 ──[ACCEPTED / REJECTED]──▶ PM      （人類最終驗收）
PM ──[RELEASE_APPROVED]──▶ FE, BE, QA   （發布核准）
```

---

## 10. 信號協議（Signal Protocol）

### 信號總表

| 信號名稱 | 發送者 | 接收者 | 觸發時機 | 攜帶資訊 |
|---------|--------|--------|---------|---------|
| `REQ_CLARIFY` | PM | 人類 | 需求資訊不足 | 釐清問題清單 |
| `REQ_INPUT` | 人類 | PM | 補充需求 | 需求補充內容 |
| `G0_PASSED` | Orchestrator | ALL | Phase 0 閘門通過 | 需求摘要 |
| `REQ_READY` | PM | FE, BE, QA | 三份需求文件產出 | 文件連結/內容 |
| `REQ_REVIEW_OK` | FE/BE/QA | PM | 審閱通過 | 無 |
| `REQ_QUESTION` | FE/BE/QA | PM | 審閱有疑問 | 問題清單 |
| `G1_PASSED` | Orchestrator | ALL | Phase 1 閘門通過 | 確認的需求文件 |
| `CONTRACT_DRAFT` | PM | FE, BE | Contract 草案發出 | Contract 文件 |
| `CONTRACT_OK` | FE/BE | PM | 同意 Contract | 無 |
| `CONTRACT_CHANGE` | FE/BE | PM | 建議修改 Contract | 修改建議 |
| `CONTRACT_LOCKED` | PM | FE, BE, QA | Contract 鎖定 | 最終 Contract |
| `G2_PASSED` | Orchestrator | ALL | Phase 2 閘門通過 | 鎖定的 Contract |
| `INIT_START` | Orchestrator | FE, BE | 開始專案初始化 | 無 |
| `INIT_COMPLETE` | FE/BE | Orchestrator | 專案初始化完成 | 初始化報告 |
| `G2.5_PASSED` | Orchestrator | ALL | Phase 2.5 閘門通過 | 無 |
| `DEV_START` | Orchestrator | FE, BE, QA | 開始並行開發 | 無 |
| `DEV_BLOCKED` | FE/BE | PM | 開發遇阻塞 | 阻塞描述 |
| `DEV_COMPLETE` | FE/BE | PM, QA | 開發完成 | 開發報告 + 提測通知 |
| `TESTCASE_READY` | QA | PM | 測試案例完成 | 測試計畫 |
| `REVIEW_START` | Orchestrator | Code Reviewer | G3-FE + G3-BE 通過後 | 需求文件 + 原始碼路徑 |
| `GAP_REPORT` | Code Reviewer | Orchestrator | 跨領域審查完成 | GAP Report（PASS / BLOCKED） |
| `GAP_FIX_REQUIRED` | Orchestrator | FE/BE | GAP_REPORT 為 BLOCKED | Gap 清單 |
| `G3_PASSED` | Orchestrator | ALL | Phase 3 閘門通過（含 G3-Review） | 無 |
| `INTEGRATION_START` | Orchestrator | FE, BE, QA | 整合開始 | 無 |
| `QA_START` | QA | PM | QA 開始驗證 | 無 |
| `BUG_FOUND` | QA | FE/BE, PM | 發現 Bug | Bug Report |
| `QA_REPORT` | QA | PM | 測試報告完成 | 測試報告 |
| `G4_PASSED` | Orchestrator | ALL | Phase 4 閘門通過 | 驗收建議 |
| `FIX_ASSIGN` | PM | FE/BE | 指派修復 | 修復任務清單 |
| `FIX_COMPLETE` | FE/BE | QA, PM | 修復完成 | 修復報告 |
| `RETEST_REPORT` | QA | PM | 重測完成 | 重測報告 |
| `G5_PASSED` | Orchestrator | ALL | 驗收通過 | 無 |
| `ESCALATE` | PM | 人類 | 3 回合用盡 | 自治流程摘要報告 |
| `DELIVERY_READY` | PM | 人類 | 交付準備完成 | 交付清單 |
| `RELEASE_APPROVED` | PM | ALL | 發布核准 | 無 |
| `FE_BE_DISCUSS` | FE/BE | BE/FE | 前後端直接技術討論 | 討論主題與結論 |
| `REQ_CHANGE` | PM | FE, BE, QA | 需求變更通知 | 變更內容與影響評估 |
| `REQ_ISSUE` | QA | PM | 發現需求矛盾或遺漏 | 問題描述 |

### 信號規則

1. **信號必須有回應**：每個信號的接收者必須在合理時間內回應
2. **阻塞信號優先處理**：`DEV_BLOCKED`、`BUG_FOUND` 屬於高優先級信號
3. **閘門信號不可跳過**：`Gx_PASSED` 信號只有在閘門條件全部滿足時才能發出
4. **上報信號立即生效**：`ESCALATE` 發出後，自治流程立即暫停

---

## 11. 品質閘門（Quality Gates）

### 閘門總表

| 閘門 | 位置 | 判定者 | 通過條件摘要 |
|------|------|--------|------------|
| **G0** | Phase 0 → 1 | PM | 需求完整、人類確認 |
| **G1** | Phase 1 → 2 | PM | 三份文件產出、FE/BE/QA 審閱通過 |
| **G2** | Phase 2 → 2.5 | PM | API Contract 鎖定、全員確認 |
| **G2.5** | Phase 2.5 → 3 | Orchestrator | FE/BE 專案初始化完成、基礎測試通過、MSW 設定完成 |
| **G3-FE** | Phase 3 → 4 | FE（自評）+ PM（確認）| 測試通過、型別檢查、Lint、報告產出 |
| **G3-BE** | Phase 3 → 4 | BE（自評）+ PM（確認）| 測試通過、覆蓋率、靜態分析、報告產出 |
| **G3-Review** | Phase 3 → 4 | Code Reviewer + Orchestrator（確認）| 需求覆蓋 100%、路由可達、Contract 合規 |
| **G4** | Phase 4 → 5/6 | QA + PM | 測試報告完成、驗收建議明確 |
| **G5** | Phase 5 → 6 | PM | 驗收通過 或 3 回合上報 |

### 閘門失敗處理

| 閘門 | 失敗時處理 |
|------|-----------|
| G0 失敗 | PM 繼續向人類釐清需求 |
| G1 失敗 | PM 修改文件，FE/BE/QA 重新審閱 |
| G2 失敗 | PM 仲裁 Contract 衝突，重新協商 |
| G2.5 失敗 | 對應 Agent 修復初始化問題，重新驗證 |
| G3-FE 失敗 | FE 修復問題，重新自我驗證 |
| G3-BE 失敗 | BE 修復問題，重新自我驗證 |
| G3-Review 失敗 | 依 Gap 清單路由回對應 FE/BE Agent 修復，修復後 Code Reviewer 重新審查 |
| G4 失敗 | 進入 Phase 5 自治修復循環 |
| G5 失敗 | 上報人類決策者 |

---

## 12. 並行工作矩陣

### 各 Phase 的 Agent 工作分配

| Phase | PM | FE | BE | QA | Code Reviewer | 並行度 |
|-------|----|----|----|----|---------------|--------|
| **0：需求輸入** | 🔵 需求收集釐清 | ⬜ 待命 | ⬜ 待命 | ⬜ 待命 | ⬜ 待命 | 1 |
| **1：需求分析** | 🔵 需求分析+文件產出 | 🟡 環境準備 | 🟡 環境準備 | 🟡 環境準備 | ⬜ 待命 | 1+3 |
| **2：Contract** | 🔵 仲裁 Contract | 🟢 審閱+提案 | 🟢 審閱+提案 | 🟢 測試計畫設計 | ⬜ 待命 | 4 |
| **2.5：專案初始化** | 🟡 監控 | 🔵 **FE 初始化** | 🔵 **BE 初始化** | 🟡 測試案例設計 | ⬜ 待命 | 2+2 |
| **3：並行開發** | 🟡 監控進度 | 🔵 **TDD 開發** | 🔵 **TDD 開發** | 🟢 測試案例設計 | 🟢 **G3-FE+BE 後審查** | 3+1+1 |
| **4：QA 驗證** | 🟡 監控品質 | 🟡 支援整合 | 🟡 支援整合 | 🔵 **獨立驗證** | ⬜ 待命 | 1+3 |
| **5：修復循環** | 🔵 主導修復 | 🟢 修復 FE Bug | 🟢 修復 BE Bug | 🟢 重新驗證 | ⬜ 待命 | 4 |
| **6：交付** | 🔵 彙整交付 | ⬜ 待命 | ⬜ 待命 | ⬜ 待命 | ⬜ 待命 | 1 |

**圖例**：🔵 主導工作 | 🟢 並行工作 | 🟡 輔助/待命 | ⬜ 無任務

### 關鍵並行點

1. **Phase 2 並行**：FE + BE 同時審閱 Contract ＋ QA 同步設計測試計畫
2. **Phase 2.5 並行**：FE 初始化 + BE 初始化 ＋ QA 持續測試案例設計
3. **Phase 3 並行**（最大並行）：FE 開發 + BE 開發 + QA 測試案例設計 + Code Reviewer 審查（G3-FE+BE 通過後） ← 全流程最高效能區間
4. **Phase 4 並行**：QA 測試同時，FE/BE 即時修復 QA 通報的 Bug
5. **Phase 5 並行**：FE 修復 FE Bug + BE 修復 BE Bug + QA 重新驗證已修復項目

---

## 13. 衝突解決與升級機制

### 三層升級階梯

```
Level 1：Agent 間直接協商
    │    （FE ↔ BE 討論技術細節，限時 1 回合）
    │
    ├── 解決 → 記錄結果，繼續流程
    │
    └── 未解決 ↓
         │
Level 2：PM 仲裁
    │    （PM 根據業務需求做最終決定，限時 1 回合）
    │
    ├── 解決 → 記錄結果，繼續流程
    │
    └── 未解決 ↓
         │
Level 3：人類決策者介入
         （暫停自治流程，上報人類）
```

### 常見衝突情境與處理

| 衝突情境 | 處理方式 | 升級層級 |
|---------|---------|---------|
| FE/BE 對 API 回應格式有分歧 | FE/BE 直接討論，提出各自理由 | L1 → L2 |
| 需求文件有歧義 | 通知 PM 釐清 | L2 |
| 技術可行性疑慮 | 工程師提出替代方案，PM 決定 | L1 → L2 |
| 需求變更影響範圍過大 | PM 評估後上報人類 | L2 → L3 |
| 3 回合修復仍無法通過驗收 | PM 產出摘要報告，上報人類 | L3 |
| QA 發現需求本身有矛盾 | QA 通知 PM，PM 判斷是否需人類確認 | L2 → L3 |

### 衝突記錄格式

```markdown
### 衝突記錄 [CONFLICT-{序號}]

- **發生階段**：Phase [X]
- **衝突方**：[Agent A] vs [Agent B]
- **衝突描述**：[具體描述]
- **各方立場**：
  - [Agent A]：[主張及理由]
  - [Agent B]：[主張及理由]
- **解決層級**：L1 / L2 / L3
- **最終決定**：[決定內容]
- **決定者**：[PM / 人類]
- **影響**：[對後續流程的影響]
```

---

## 14. 流程狀態追蹤

### 全流程狀態追蹤表

每次流程啟動時建立此表，隨流程推進即時更新：

```markdown
## [Feature Name] — 流程狀態追蹤

### 基本資訊
- 功能名稱：
- 啟動時間：
- 當前 Phase：
- 當前回合（若在 Phase 5）：

### Phase 狀態
| Phase | 狀態 | PM | FE | BE | QA | 閘門 | 阻塞項 | 備註 |
|-------|------|----|----|----|----|------|--------|------|
| 0 | ✅ 完成 | ✅ | - | - | - | G0 ✅ | 無 | |
| 1 | ✅ 完成 | ✅ | ✅ 審閱 | ✅ 審閱 | ✅ 審閱 | G1 ✅ | 無 | |
| 2 | ✅ 完成 | ✅ | ✅ 確認 | ✅ 確認 | ✅ 計畫 | G2 ✅ | 無 | |
| 2.5 | ✅ 完成 | 🟡 | ✅ 初始化 | ✅ 初始化 | 🟡 設計中 | G2.5 ✅ | 無 | |
| 3 | 🔄 進行中 | 🟡 | 🔵 開發中 | ✅ 完成 | 🟢 設計中 | G3 ⬜ | 無 | BE 已完成，等待 FE |
| 4 | ⬜ 未開始 | - | - | - | - | G4 ⬜ | - | |
| 5 | ⬜ 未開始 | - | - | - | - | G5 ⬜ | - | |
| 6 | ⬜ 未開始 | - | - | - | - | - | - | |

### 狀態符號
- ✅ 完成
- 🔄 進行中
- 🔵 主導工作中
- 🟢 並行工作中
- 🟡 輔助/監控
- ⬜ 未開始
- 🔴 阻塞
- ❌ 失敗

### 信號紀錄
| 時序 | 信號 | 發送者 | 接收者 | 備註 |
|------|------|--------|--------|------|
| 1 | G0_PASSED | PM | ALL | 需求確認 |
| 2 | REQ_READY | PM | FE,BE,QA | 三份文件發出 |
| 3 | REQ_REVIEW_OK | FE | PM | 前端審閱通過 |
| 4 | REQ_REVIEW_OK | BE | PM | 後端審閱通過 |
| 5 | REQ_REVIEW_OK | QA | PM | QA 審閱通過 |
| 6 | G1_PASSED | PM | ALL | |
| ... | ... | ... | ... | ... |

### 衝突紀錄
（若有衝突，記錄於此）

### 自治修復紀錄
（若進入 Phase 5，記錄每回合狀態）
```

---

## 15. 自治邊界與人類介入點

### 可自治決策（不需人類介入）

| 決策類型 | 決策者 | 範例 |
|---------|--------|------|
| 技術實作細節 | FE / BE | 選擇設計模式、命名、重構方式 |
| API Contract 細節 | PM（仲裁 FE/BE） | 回應格式、分頁參數命名 |
| Bug 優先級排序 | PM + QA | 修復順序決定 |
| 測試策略調整 | QA | 增減測試案例、調整測試優先級 |
| 修復方案選擇 | FE / BE | 技術修復方式 |
| Sprint 內工作排序 | PM | Story 開發順序調整 |

### 必須人類介入的決策

| 決策類型 | 觸發條件 | 上報內容 |
|---------|---------|---------|
| 需求確認 | Phase 0 | 需求摘要 |
| 需求變更 | 任何 Phase 發現需求矛盾 | 變更影響評估 |
| 3 回合修復失敗 | Phase 5 第 3 回合後 | 自治流程摘要報告 + PM 建議方案 |
| 範圍外技術決策 | 需要新增未預期的第三方服務 | 技術評估報告 |
| 安全性疑慮 | 任何 Phase 發現安全風險 | 風險評估 |
| 最終發布核准 | Phase 6 | 完整交付清單 |

### 人類介入點時序圖

```
人類干預點                       自治流程
    │                              │
    ●─── 需求輸入 ───────────────▶ Phase 0
    │                              │
    ●─── 需求確認 ───────────────▶ Phase 0 → Phase 1
    │                              │
    │    （自治運行）                │
    │                              ├── Phase 1（PM 分析）
    │                              ├── Phase 2（Contract 協商）
    │                              ├── Phase 2.5（專案初始化）
    │                              ├── Phase 3（並行開發）
    │                              ├── Phase 4（QA 驗證）
    │                              ├── Phase 5（修復循環）
    │                              │
    ●◄── 上報（若修復失敗）◄─────── Phase 5（ESCALATE）
    │                              │
    ●◄── 交付驗收 ◄──────────────── Phase 6
    │                              │
    ●─── 發布核准 ───────────────▶ 完成
```

**正常情況下，人類只需介入 3 次**：
1. 提出需求 + 確認需求摘要
2. （若修復失敗）收到上報並做決定
3. 最終驗收與發布核准

---

## 16. 需求變更處理流程

任何 Phase 中發現需求問題或收到需求變更時，啟動以下流程：

### 觸發來源

| 來源 | 信號 | 說明 |
|------|------|------|
| 人類主動變更 | `REQ_CHANGE`（人類 → PM） | 客戶/決策者提出新需求或修改 |
| QA 發現需求矛盾 | `REQ_ISSUE`（QA → PM） | 測試中發現需求本身有問題 |
| FE/BE 技術疑慮 | `REQ_QUESTION`（FE/BE → PM） | 開發中發現需求不可行 |

### 處理流程（PM 主導，依 PM Prompt「當收到需求變更時」）

```
收到變更/問題
     │
     ▼
1. PM 評估變更影響範圍（哪些 User Story、哪些角色受影響）
     │
     ▼
2. PM 判斷變更規模
     │
     ├── 小幅調整（不影響 Contract）
     │   → PM 更新需求文件，標注 [變更]
     │   → 通知受影響的 Agent
     │   → 流程繼續
     │
     ├── 中幅變更（影響 Contract 或 AC）
     │   → PM 重新評估 Story Point 與 Sprint 計劃
     │   → 啟動 Contract 變更流程（若涉及 API）
     │   → 通知所有 Agent
     │
     └── 重大變更（影響架構或核心需求）
         → PM 產出變更影響評估報告
         → 上報人類決策者確認
         → 人類確認後，可能回退到對應 Phase 重新執行
```

### QA 發現需求問題的特殊處理

依 QA Prompt 定義：「發現需求層級的問題（非 Bug，而是需求本身的矛盾或遺漏），立即通知 PM」

- QA 發送 `REQ_ISSUE` → PM
- PM 判斷是需求遺漏還是矛盾
- 若為遺漏 → PM 補充需求，通知相關 Agent
- 若為矛盾 → PM 評估後可能需上報人類確認

---

## 附錄 A：快速啟動指引

### 啟動自治流程

1. 人類向 PM Agent 提出需求
2. PM Agent 自動進入 Phase 0，開始需求分析
3. 流程按照本協議自動推進，直到需要人類介入或交付完成

### 啟動命令範例

```
人類：「我需要一個用戶註冊功能，包含 email 驗證和密碼強度檢查。」

→ PM Agent 啟動 Phase 0
→ PM 可能回覆釐清問題，或直接產出需求摘要
→ 人類確認後，自治流程自動推進至完成
```

### 監控流程

人類可隨時查看「流程狀態追蹤表」了解當前進度，但**不需要主動干預**，除非收到 `ESCALATE` 信號。

---

## 附錄 B：與現有 Agent Prompt 的對照表

本協議不修改任何現有 Agent System Prompt，而是作為頂層協調層。以下是對照關係：

| 本協議章節 | 對應 Agent Prompt 段落 |
|-----------|----------------------|
| Phase 1 PM 工作流程 | PM Prompt > 需求分析流程 Step 1-5 |
| Phase 1 文件模板 | PM Prompt > 前端/後端/測試需求文件模板 |
| Phase 2 API Contract | PM Prompt > 技術需求評估 + BE Prompt > API 設計規範 |
| Phase 3 FE 開發 | FE Prompt > TDD 開發流程 4.1-4.5 |
| Phase 3 BE 開發 | BE Prompt > TDD 開發流程 Step 1-3 + 接收 PM 需求文件流程 |
| Phase 3 自我驗證 | FE Prompt > PR 規範 + BE Prompt > 自我驗證 Step 4 |
| Phase 3 DEV_COMPLETE 前檢 | FE Prompt > 13.5 DEV_COMPLETE 前置檢查 + BE Prompt > DEV_COMPLETE 前置檢查 |
| Phase 3→4 跨領域審查 | Code Reviewer Prompt > 審查項目 1-3（需求覆蓋、可達性、Contract 合規） |
| Phase 4 QA 驗證 | QA Prompt > 測試工作流程 Step 1-6 |
| Phase 4 Bug 管理 | QA Prompt > Bug 管理（回報格式、嚴重程度、生命週期） |
| Phase 5 自治修復 | PM Prompt > 驗收不通過自治流程（完整引用） |
| 品質閘門 | 各 Agent Prompt 中的驗證標準彙整 |
| 信號協議 | 各 Agent Prompt > 與團隊互動協議（標準化） |

---

## 附錄 C：配套文件索引

本協議與以下配套文件共同構成完整的自治並行開發體系：

| 文件 | 用途 | 使用者 |
|------|------|--------|
| `autonomous-workflow.md` | 頂層流程協議（本文件） | Orchestrator + 所有 Agent |
| `orchestrator-agent-system-prompt.md` | 流程指揮者 Agent 定義 | Orchestrator Agent |
| `product-owner-agent-system-prompt.md` | PM Agent 定義 | PM Agent |
| `frontend-engineer-system-prompt.md` | 前端工程師 Agent 定義 | FE Agent |
| `backend-engineer-agent-system-prompt.md` | 後端工程師 Agent 定義 | BE Agent |
| `qa-engineer-agent-system-prompt.md` | 測試工程師 Agent 定義 | QA Agent |
| `signal-protocol.md` | 信號格式規格與範例 | Orchestrator + 所有 Agent |
| `api-contract-template.md` | API Contract 標準化模板 | PM（建立）、FE/BE（審閱） |
| `flow-state-tracker-template.md` | 流程狀態追蹤表模板 | Orchestrator（維護） |
| `quality-gate-checklist.md` | 品質閘門可執行檢查清單 | Orchestrator（判定） |
