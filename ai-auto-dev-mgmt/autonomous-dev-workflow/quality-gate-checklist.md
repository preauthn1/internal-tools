# 品質閘門檢查清單（Quality Gate Checklist）

> **使用時機**：每次 Phase 轉換時，由 Orchestrator Agent 執行閘門判定
> **判定規則**：所有項目必須為 ✅ 才能通過閘門。任何 ❌ 項目必須先解決。

---

## 角色責任矩陣

| 閘門 | PM | FE | BE | QA | Code Reviewer | Orchestrator |
|------|----|----|----|----|---------------|-------------|
| **G0** | 🔵 判定者 | — | — | — | — | ✅ 確認 |
| **G1** | 🔵 產出者 | 👀 審閱 | 👀 審閱 | 👀 審閱 | — | ✅ 確認 |
| **G2** | 🔵 仲裁者 | 👀 確認 | 👀 確認 | — | — | ✅ 確認 |
| **G2.5** | — | 🔵 執行者 | 🔵 執行者 | 🔵 執行者 | — | ✅ 確認 |
| **G3-FE** | — | 🔵 自評者 | — | — | — | ✅ 確認 |
| **G3-BE** | — | — | 🔵 自評者 | — | — | ✅ 確認 |
| **G3-Review** | — | — | — | — | 🔵 審查者 | ✅ 確認 |
| **G4** | 👀 評估 | — | — | 🔵 驗證者 | — | ✅ 確認 |
| **G5** | 🔵 驗收者 | 🔧 修復 | 🔧 修復 | 🔵 重測者 | — | ✅ 確認 |
| **G6** | 🔵 交付者 | — | — | — | — | ✅ 確認 + 人類審閱 |

> 🔵 主要負責 | 👀 參與審閱 | 🔧 修復角色 | ✅ 閘門確認 | — 不參與

---

## G0：需求輸入 → 需求分析

**判定者**：PM + Orchestrator
**觸發條件**：PM 完成需求收集

| # | 檢查項目 | 判定方式 | 結果 |
|---|---------|---------|------|
| 1 | 需求描述完整（5W1H 已涵蓋） | PM 確認已覆蓋 Who/What/When/Where/Why/How | ⬜ |
| 2 | 使用者角色（Persona）已識別 | PM 已定義至少一個主要 Persona | ⬜ |
| 3 | 業務目標明確 | PM 能以一句話說明此功能的業務價值 | ⬜ |
| 4 | 人類已確認需求摘要 | 收到人類的明確確認回覆 | ⬜ |

**通過** → 發送 `G0_PASSED`，進入 Phase 1
**未通過** → PM 繼續向人類釐清（`REQ_CLARIFY`）

---

## G1：需求分析 → API Contract 協商

**判定者**：PM + Orchestrator
**觸發條件**：PM 產出三份需求文件，FE/BE/QA 完成審閱

| # | 檢查項目 | 判定方式 | 結果 |
|---|---------|---------|------|
| 1 | 前端需求文件已產出 | 文件包含：頁面描述、UI 需求、User Flow、表單驗證、API 對接、前端 AC | ⬜ |
| 2 | 後端需求文件已產出 | 文件包含：商業邏輯、API 設計、資料模型、權限控制、效能需求、後端 AC | ⬜ |
| 3 | 測試驗收文件已產出 | 文件包含：功能概述、User Story 清單、測試情境（正向/負向/邊界）、測試資料 | ⬜ |
| 4 | 所有 User Story 通過品質檢查清單 | 逐項確認 PM Prompt 中的 10 項品質檢查 | ⬜ |
| 5 | FE 審閱通過 | 收到 FE 的 `REQ_REVIEW_OK` 信號 | ⬜ |
| 6 | BE 審閱通過 | 收到 BE 的 `REQ_REVIEW_OK` 信號 | ⬜ |
| 7 | QA 審閱通過 | 收到 QA 的 `REQ_REVIEW_OK` 信號 | ⬜ |

**User Story 品質檢查子清單**（對應 PM Prompt）：
| # | 子項目 | 結果 |
|---|--------|------|
| 4a | 使用者角色明確且具體 | ⬜ |
| 4b | 功能描述清晰無歧義 | ⬜ |
| 4c | 業務價值明確說明 | ⬜ |
| 4d | 驗收標準涵蓋正向與異常情境 | ⬜ |
| 4e | 驗收標準使用 Given-When-Then 格式 | ⬜ |
| 4f | Story Point 已評估 | ⬜ |
| 4g | 依賴關係已標注 | ⬜ |
| 4h | 前端需求已拆分並標注 | ⬜ |
| 4i | 後端需求已拆分並標注 | ⬜ |
| 4j | 可測試性已確認 | ⬜ |

**通過** → 發送 `G1_PASSED`，進入 Phase 2
**未通過** → PM 修改文件，重新發送審閱

---

## G2：API Contract 協商 → 並行開發

**判定者**：PM + Orchestrator
**觸發條件**：API Contract 協商完成

| # | 檢查項目 | 判定方式 | 結果 |
|---|---------|---------|------|
| 1 | API Contract 文件完整 | 所有端點含：方法、URL、請求/回應格式、錯誤碼、驗證規則 | ⬜ |
| 2 | FE 確認 Contract 可支撐前端需求 | 收到 FE 的 `CONTRACT_OK` 信號 | ⬜ |
| 3 | BE 確認 Contract 技術可行 | 收到 BE 的 `CONTRACT_OK` 信號 | ⬜ |
| 4 | QA 確認 Contract 可測試 | QA 確認每個端點都有可驗證的預期行為 | ⬜ |
| 5 | Contract 狀態為 🔓 已鎖定 | Contract 文件標記為已鎖定，版本號已定義 | ⬜ |
| 6 | QA 測試計畫初稿完成 | QA 已產出測試策略、測試範圍、風險評估 | ⬜ |

**Contract 完整性子清單**：
| # | 子項目 | 結果 |
|---|--------|------|
| 1a | 共用約定已定義（回應格式、分頁、排序、時間格式） | ⬜ |
| 1b | 認證方式已定義 | ⬜ |
| 1c | HTTP 狀態碼對照已定義 | ⬜ |
| 1d | 每個端點有完整的請求規格 | ⬜ |
| 1e | 每個端點有完整的成功回應規格 | ⬜ |
| 1f | 每個端點有完整的錯誤回應規格 | ⬜ |
| 1g | FE 型別定義建議已提供 | ⬜ |

**通過** → 發送 `G2_PASSED` + `CONTRACT_LOCKED`，進入 Phase 2.5
**未通過** → PM 仲裁衝突，重新協商

---

## G2.5：專案初始化 → 並行開發

**判定者**：Orchestrator
**觸發條件**：FE 和 BE 都發送 `INIT_COMPLETE`

| # | 檢查項目 | 判定方式 | 結果 |
|---|---------|---------|------|
| 1 | BE 專案可正常啟動 | `php artisan serve` 成功 | ⬜ |
| 2 | FE 專案可正常啟動 | `npm run dev` 成功 | ⬜ |
| 3 | BE 基礎測試可執行 | `php artisan test` 通過（即使只有預設測試） | ⬜ |
| 4 | FE 基礎測試可執行 | `npm run test` 通過（即使只有預設測試） | ⬜ |
| 5 | MSW Mock 已設定 | 依據 API Contract 建立 handlers | ⬜ |
| 6 | 兩個專案皆已 Git 初始 commit | `git log` 有初始 commit | ⬜ |

**通過** → 發送 `G2.5_PASSED`，進入 Phase 3
**未通過** → 對應 Agent 修復問題，重新初始化

---

## G3-FE：前端開發自我驗證

**判定者**：FE（自評）+ PM + Orchestrator（確認）
**觸發條件**：FE 發送 `DEV_COMPLETE`

| # | 檢查項目 | 判定方式 | 結果 |
|---|---------|---------|------|
| 1 | Vitest 單元測試全部通過 | `vitest run` 輸出 0 failures | ⬜ |
| 2 | vue-tsc 型別檢查無錯誤 | `vue-tsc --noEmit` 輸出 0 errors | ⬜ |
| 3 | ESLint 無警告或錯誤 | `eslint .` 輸出 0 warnings, 0 errors | ⬜ |
| 4 | Playwright E2E 測試通過 | `npx playwright test` 在 Mock API 環境下全部通過 | ⬜ |
| 5 | 所有前端 AC 有對應測試 | 每個 AC-FE-xxx 至少有一個測試案例 | ⬜ |
| 6 | 無 console.log / debugger 殘留 | 程式碼搜索無 console.log 或 debugger | ⬜ |
| 7 | 前端開發報告已產出 | 報告包含：實作概述、元件清單、測試報告、AC 對應 | ⬜ |
| 8 | 提測通知已產出 | 通知包含：功能入口、測試帳號、AC 覆蓋 | ⬜ |

**通過** → G3-FE ✅（等待 G3-BE 一起判定 G3）
**未通過** → FE 修復問題，重新自我驗證

---

## G3-BE：後端開發自我驗證

**判定者**：BE（自評）+ PM + Orchestrator（確認）
**觸發條件**：BE 發送 `DEV_COMPLETE`

| # | 檢查項目 | 判定方式 | 結果 |
|---|---------|---------|------|
| 1 | php artisan test 全部通過 | `php artisan test` 輸出 0 failures | ⬜ |
| 2 | 單元測試覆蓋率 ≥ 90% | `--coverage` 輸出 Unit coverage ≥ 90% | ⬜ |
| 3 | 功能測試覆蓋率 ≥ 85% | `--coverage` 輸出 Feature coverage ≥ 85% | ⬜ |
| 4 | 整合測試覆蓋率 ≥ 80% | `--coverage` 輸出 Integration coverage ≥ 80% | ⬜ |
| 5 | Laravel Pint 無風格問題 | `./vendor/bin/pint --test` 輸出 0 dirty files | ⬜ |
| 6 | PHPStan / Larastan 靜態分析通過 | `./vendor/bin/phpstan analyse` 輸出 0 errors | ⬜ |
| 7 | 所有後端 AC 有對應測試 | 每個 AC-BE-xxx 至少有一個測試案例 | ⬜ |
| 8 | 無 dd() / dump() 殘留 | 程式碼搜索無 dd() 或 dump() | ⬜ |
| 9 | API 端點文件已產出 | Scribe / Swagger / 手動撰寫的 API 文件 | ⬜ |
| 10 | 後端開發報告已產出 | 報告包含：實作概述、DB 變更、API 清單、測試報告、AC 對應 | ⬜ |
| 11 | 提測通知已產出 | 通知包含：API 端點、環境設定、Migration 狀態 | ⬜ |

**通過** → G3-BE ✅（若 G3-FE ✅ **且** G3-BE ✅ → 啟動 Code Reviewer Agent 進行 G3-Review 審查）
**未通過** → BE 修復問題，重新自我驗證

---

## G3-Review：跨領域審查（Code Reviewer Agent）

> G3-FE ✅ 且 G3-BE ✅ 後，啟動 Code Reviewer Agent 進行跨領域審查。

| # | 檢查項目 | 判定標準 | 結果 |
|---|---------|---------|------|
| 1 | 需求覆蓋矩陣 | 所有需求元件/頁面有對應原始碼（覆蓋率 100%） | ⬜ |
| 2 | 頁面可達性 | 所有路由可從 UI 導航到達（無孤島頁面） | ⬜ |
| 3 | API Contract 合規 | 所有端點在 BE Route + FE Service + MSW Handler 中有對應 | ⬜ |
| 4 | GAP_REPORT 判定 | PASS | ⬜ |

**通過** → G3-Review ✅
**未通過（BLOCKED）** → 依 Gap 清單路由回 FE/BE 修復，修復後重新審查

---

## G3 綜合判定

| 子閘門 | 結果 |
|--------|------|
| G3-FE | ⬜/✅/❌ |
| G3-BE | ⬜/✅/❌ |
| G3-Review | ⬜/✅/❌ |
| **G3 綜合** | **G3-FE ✅ 且 G3-BE ✅ 且 G3-Review ✅ → 通過** |

**通過** → 發送 `G3_PASSED`，進入 Phase 4
**部分通過** → 等待未完成方或修復 Gap，不推進

---

## G4：QA 驗證 → 自治修復 / 交付

**判定者**：QA + PM + Orchestrator
**觸發條件**：QA 發送 `QA_REPORT`

| # | 檢查項目 | 判定方式 | 結果 |
|---|---------|---------|------|
| 1 | QA 測試報告已完成 | 報告包含：測試概況、AC 覆蓋、Bug 清單、風險、驗收建議 | ⬜ |
| 2 | 驗收建議已明確標注 | QA 明確標記 ✅ / ⚠️ / ❌ | ⬜ |
| 3 | 所有 Bug 已分類 | 每個 Bug 有嚴重程度（Critical/Major/Minor/Trivial） | ⬜ |
| 4 | P0 測試案例已全部執行 | 無跳過的 P0 案例 | ⬜ |
| 5 | P1 測試案例已全部執行 | 無跳過的 P1 案例 | ⬜ |

**路由判定**：

| QA 驗收建議 | 判定條件（依 QA Prompt） | 路由 |
|------------|------------------------|------|
| ✅ 建議通過 | 所有 AC 通過、無 Critical/Major Bug、通過率 ≥ 95% | → Phase 6 |
| ⚠️ 有條件通過 | 核心 AC 通過、無 Critical Bug、Minor Bug ≤ 3、通過率 ≥ 85% | → Phase 5 |
| ❌ 不建議通過 | 有 Critical Bug 未修復、核心 AC 未通過、通過率 < 85% | → Phase 5 |

**通過（✅）** → 發送 `G4_PASSED`，跳至 Phase 6
**未通過（⚠️/❌）** → 發送 `G4_PASSED`，進入 Phase 5

---

## G5：自治修復 → 交付 / 上報

**判定者**：PM + Orchestrator
**觸發條件**：QA 發送 `RETEST_REPORT`

| 條件 | 結果 | 動作 |
|------|------|------|
| 回合 ≤ 3 且 QA 驗收 ✅ | G5 ✅ 通過 | → Phase 6 |
| 回合 ≤ 3 且 QA 驗收 ⚠️/❌ | G5 ⏳ 繼續 | → 回合 +1，回到 Phase 5 Step A |
| 回合 > 3 | G5 ❌ 上報 | → ESCALATE → 暫停流程 |

**上報檢查清單**（回合 > 3 時）：
| # | 檢查項目 | 結果 |
|---|---------|------|
| 1 | PM 已產出「自治流程摘要報告」 | ⬜ |
| 2 | 報告包含三回合歷程摘要 | ⬜ |
| 3 | 報告包含未解決問題清單 | ⬜ |
| 4 | 報告包含 PM 建議方案（至少兩個） | ⬜ |
| 5 | 報告包含需要人類決策的事項 | ⬜ |

**通過（✅）** → 發送 `G5_PASSED`，進入 Phase 6
**上報** → 發送 `ESCALATE`，暫停流程等待人類決策

---

## 快速參考卡

```
G0: 需求完整 + 人類確認                    → Phase 1
G1: 三份文件 + 三方審閱OK                   → Phase 2
G2: Contract 鎖定 + 全員確認 + QA 計畫      → Phase 2.5
G2.5: 專案可啟動 + 基礎測試 + MSW 設定          → Phase 3
G3: G3-FE ✅ + G3-BE ✅                     → Phase 4
G4: QA 報告 + 驗收建議                      → Phase 5 or 6
G5: QA ✅ (≤3回合) or ESCALATE (>3回合)     → Phase 6 or 人類
```
