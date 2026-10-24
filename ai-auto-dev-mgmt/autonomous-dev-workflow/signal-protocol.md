# 信號協議執行規格（Signal Protocol Specification）

> **版本**：1.0
> **配合文件**：`autonomous-workflow.md`、`orchestrator-agent-system-prompt.md`

---

## 1. 信號格式標準

每個信號必須以結構化 Markdown 格式傳遞，包含以下欄位：

```markdown
## SIGNAL: [信號名稱]

- **發送者**：[Agent 角色]
- **接收者**：[Agent 角色(s)]
- **時間**：[timestamp]
- **Phase**：Phase [N]
- **關聯 Feature**：[Feature Name]

### 攜帶資訊
[信號的具體內容]

### 需要回應
[是 / 否] — 若是，預期回應信號為 [信號名稱]
```

---

## 2. 信號分類與優先級

### 2.1 優先級定義

| 優先級 | 標記 | 說明 | 處理時限 |
|--------|------|------|---------|
| **P0 — 緊急** | ⚡ | 阻塞信號、上報信號 | 立即處理 |
| **P1 — 高** | 🔴 | 閘門判定、Bug 通報 | 當回合內處理 |
| **P2 — 正常** | 🟡 | 審閱結果、完成通知 | 正常流程處理 |
| **P3 — 資訊** | 🟢 | 狀態更新、記錄用 | 記錄即可 |

### 2.2 信號分類表

| 信號名稱 | 優先級 | 類別 | 需回應 |
|---------|--------|------|--------|
| `ESCALATE` | ⚡ P0 | 上報 | 是（人類決策） |
| `DEV_BLOCKED` | ⚡ P0 | 阻塞 | 是（PM 處理） |
| `BUG_FOUND` | 🔴 P1 | 品質 | 是（FE/BE 確認） |
| `INIT_START` | 🟡 P2 | 啟動 | 是（INIT_COMPLETE） |
| `INIT_COMPLETE` | 🟡 P2 | 完成 | 否 |
| `G0_PASSED` ~ `G5_PASSED`（含 `G2.5_PASSED`） | 🔴 P1 | 閘門 | 否（推進流程） |

> **備註**：所有 `Gx_PASSED` 信號（`G0_PASSED` ~ `G5_PASSED`，含 `G2.5_PASSED`）的發送者均為 **Orchestrator**。
| `REQ_QUESTION` | 🔴 P1 | 釐清 | 是（PM 回覆） |
| `CONTRACT_CHANGE` | 🔴 P1 | 協商 | 是（PM 仲裁） |
| `REQ_ISSUE` | 🔴 P1 | 品質 | 是（PM 判斷） |
| `REQ_READY` | 🟡 P2 | 交付 | 是（FE/BE/QA 審閱） |
| `CONTRACT_DRAFT` | 🟡 P2 | 交付 | 是（FE/BE 審閱） |
| `CONTRACT_LOCKED` | 🟡 P2 | 決策 | 否 |
| `DEV_COMPLETE` | 🟡 P2 | 完成 | 否（Orchestrator 記錄） |
| `FIX_ASSIGN` | 🟡 P2 | 指派 | 是（FE/BE 確認） |
| `FIX_COMPLETE` | 🟡 P2 | 完成 | 否 |
| `QA_REPORT` | 🟡 P2 | 交付 | 是（PM 確認） |
| `RETEST_REPORT` | 🟡 P2 | 交付 | 是（PM 確認） |
| `DELIVERY_READY` | 🟡 P2 | 交付 | 是（人類驗收） |
| `REQ_REVIEW_OK` | 🟢 P3 | 確認 | 否 |
| `CONTRACT_OK` | 🟢 P3 | 確認 | 否 |
| `DEV_START` | 🟢 P3 | 啟動 | 否 |
| `QA_START` | 🟢 P3 | 啟動 | 否 |
| `TESTCASE_READY` | 🟢 P3 | 完成 | 否 |
| `FE_BE_DISCUSS` | 🟢 P3 | 協作 | 否（記錄用） |
| `REQ_CHANGE` | 🟡 P2 | 變更 | 是（全員確認） |
| `RELEASE_APPROVED` | 🟢 P3 | 結束 | 否 |
| `REQ_CLARIFY` | 🔴 P1 | 釐清 | 是（人類回覆） |
| `REQ_INPUT` | 🟡 P2 | 輸入 | 否 |

---

## 3. 各信號詳細規格

### 3.1 需求相關信號

#### REQ_CLARIFY
```markdown
## SIGNAL: REQ_CLARIFY

- **發送者**：PM
- **接收者**：人類
- **Phase**：Phase 0
- **需要回應**：是 → REQ_INPUT

### 釐清問題清單
1. [問題 1 — 上下文與為何需要釐清]
2. [問題 2]
3. ...

### 目前已知資訊
[PM 目前已掌握的需求資訊摘要]
```

#### REQ_READY
```markdown
## SIGNAL: REQ_READY

- **發送者**：PM
- **接收者**：FE, BE, QA
- **Phase**：Phase 1
- **需要回應**：是 → REQ_REVIEW_OK 或 REQ_QUESTION

### 已產出文件
1. 前端需求文件：[內容/連結]
2. 後端需求文件：[內容/連結]
3. 測試驗收文件：[內容/連結]

### User Story 清單摘要
| Story ID | 標題 | 優先級 | Story Point |
|----------|------|--------|-------------|
| US-XXX-001 | | Must | 5 |

### 請各角色審閱對應文件，回覆 REQ_REVIEW_OK 或 REQ_QUESTION
```

#### REQ_QUESTION
```markdown
## SIGNAL: REQ_QUESTION

- **發送者**：[FE / BE / QA]
- **接收者**：PM
- **Phase**：Phase 1
- **需要回應**：是 → PM 直接回覆

### 問題清單
| # | 問題 | 相關文件段落 | 阻塞程度 |
|---|------|------------|---------|
| 1 | [問題描述] | [引用段落] | 阻塞 / 非阻塞 |
| 2 | | | |

### 建議方案（若有）
[提問者的建議]
```

#### REQ_CHANGE
```markdown
## SIGNAL: REQ_CHANGE

- **發送者**：PM
- **接收者**：FE, BE, QA
- **Phase**：Phase [N]
- **需要回應**：是 → 各方確認影響

### 變更內容
| 項目 | 變更前 | 變更後 | 原因 |
|------|--------|--------|------|
| | | | |

### 影響評估
- 影響的 User Story：[列表]
- 影響的 AC：[列表]
- 影響的 API Contract：[是/否]
- Story Point 變更：[+N / -N / 不變]

### 需要各角色確認影響範圍並調整工作計畫
```

#### REQ_ISSUE
```markdown
## SIGNAL: REQ_ISSUE

- **發送者**：QA
- **接收者**：PM
- **Phase**：Phase 4 / 5
- **需要回應**：是 → PM 判斷處理方式

### 問題描述
[需求本身的矛盾或遺漏描述]

### 發現情境
- 測試案例：[TC-XXX-XXX]
- 對應 AC：[AC-XX-XXX]
- 矛盾點：[具體描述]

### QA 觀察
[QA 的分析與建議]
```

### 3.2 Contract 相關信號

#### CONTRACT_DRAFT
```markdown
## SIGNAL: CONTRACT_DRAFT

- **發送者**：PM
- **接收者**：FE, BE
- **Phase**：Phase 2
- **需要回應**：是 → CONTRACT_OK 或 CONTRACT_CHANGE

### API Contract 草案
[完整 API Contract 文件內容，使用 api-contract-template.md 格式]

### 請 FE 和 BE 各自審閱，回覆確認或修改建議
```

#### CONTRACT_CHANGE
```markdown
## SIGNAL: CONTRACT_CHANGE

- **發送者**：[FE / BE]
- **接收者**：PM
- **Phase**：Phase 2
- **需要回應**：是 → PM 仲裁結果

### 修改建議
| API 端點 | 現有定義 | 建議修改 | 修改原因 |
|---------|---------|---------|---------|
| | | | |

### 影響說明
[修改對前端/後端的影響分析]
```

#### CONTRACT_LOCKED
```markdown
## SIGNAL: CONTRACT_LOCKED

- **發送者**：PM
- **接收者**：FE, BE, QA
- **Phase**：Phase 2
- **需要回應**：否

### 最終 API Contract
[完整鎖定版 API Contract]

### Contract 版本：v[X.Y]
### 狀態：🔓 已鎖定

### ⚠️ 此後任何 Contract 變更須走變更流程
```

### 3.2.5 專案初始化信號

#### INIT_START
```markdown
## SIGNAL: INIT_START

- **發送者**：Orchestrator
- **接收者**：FE, BE
- **Phase**：Phase 2.5
- **需要回應**：是 → INIT_COMPLETE

### 初始化指令
- **BE**：初始化 Laravel 專案、基礎 Migration
- **FE**：初始化 Vue 3 + TS 專案、安裝依賴、設定 MSW Mock
- **專案路徑**：[人類指定路徑]
- **API Contract 參照**：[Contract 文件連結]
```

#### INIT_COMPLETE
```markdown
## SIGNAL: INIT_COMPLETE

- **發送者**：[FE / BE]
- **接收者**：Orchestrator
- **Phase**：Phase 2.5
- **需要回應**：否

### 初始化結果
| 檢查項目 | 結果 |
|---------|------|
| 專案可啟動 | ✅ / ❌ |
| 基礎測試可執行 | ✅ / ❌ |
| Git 初始 commit 完成 | ✅ / ❌ |

### 專案路徑
[實際路徑]

### 備註
[任何初始化過程中的注意事項]
```

### 3.3 開發相關信號

#### DEV_BLOCKED
```markdown
## SIGNAL: DEV_BLOCKED

- **發送者**：[FE / BE]
- **接收者**：PM
- **Phase**：Phase 3
- **優先級**：⚡ P0
- **需要回應**：是 → PM 處理方案

### 阻塞描述
- **阻塞類型**：需求不清 / 技術障礙 / Contract 問題 / 環境問題 / 依賴其他 Agent
- **阻塞的 AC**：[AC-XX-XXX]
- **具體描述**：[詳細說明]

### 已嘗試的解決方式
[列出已嘗試的方案]

### 需要的支援
[具體需要什麼幫助]
```

#### DEV_COMPLETE
```markdown
## SIGNAL: DEV_COMPLETE

- **發送者**：[FE / BE]
- **接收者**：PM, QA
- **Phase**：Phase 3
- **需要回應**：否

### 開發報告
[依各自 prompt 定義的開發報告格式]

### 提測通知
[依各自 prompt 定義的提測通知格式]

### 自我驗證結果
| 驗證項目 | 結果 |
|---------|------|
| [項目] | ✅ / ❌ |

### 環境說明
[環境相關注意事項]
```

### 3.4 測試相關信號

#### BUG_FOUND
```markdown
## SIGNAL: BUG_FOUND

- **發送者**：QA
- **接收者**：[FE / BE]（依 Bug 歸屬）, PM
- **Phase**：Phase 4 / 5
- **優先級**：🔴 P1
- **需要回應**：是 → FE/BE 確認收到

### Bug Report
[使用 QA prompt 定義的 Bug Report 格式]

### Bug ID：BUG-[模組]-[序號]
### 歸屬：[FE] / [BE]
### 嚴重程度：Critical / Major / Minor / Trivial
```

#### QA_REPORT
```markdown
## SIGNAL: QA_REPORT

- **發送者**：QA
- **接收者**：PM
- **Phase**：Phase 4
- **需要回應**：是 → PM 確認並決定下一步

### 測試報告
[使用 QA prompt 定義的測試報告格式]

### 驗收建議：✅ 建議通過 / ⚠️ 有條件通過 / ❌ 不建議通過
### 結論說明：[理由]
```

### 3.5 修復相關信號

#### FIX_ASSIGN
```markdown
## SIGNAL: FIX_ASSIGN

- **發送者**：PM
- **接收者**：[FE / BE]
- **Phase**：Phase 5
- **回合**：第 [N] / 3 回合
- **需要回應**：是 → FIX_COMPLETE

### 修復任務清單
| Bug ID / AC 編號 | 問題描述 | 修復方案 | 優先級 |
|-----------------|---------|---------|--------|
| | | | |

### 修復要求
- 必須使用 TDD 方式修復（先寫重現 Bug 的失敗測試）
- 修復後必須通過自我驗證
- 完成後發送 FIX_COMPLETE
```

#### FIX_COMPLETE
```markdown
## SIGNAL: FIX_COMPLETE

- **發送者**：[FE / BE]
- **接收者**：QA, PM
- **Phase**：Phase 5
- **需要回應**：否（QA 自動開始重新驗證）

### 修復報告
| Bug ID / AC 編號 | 修復方式 | 新增測試 | 自我驗證 |
|-----------------|---------|---------|---------|
| | | [測試名稱] | ✅ / ❌ |

### 影響範圍
[修復可能影響的其他功能]
```

#### RETEST_REPORT
```markdown
## SIGNAL: RETEST_REPORT

- **發送者**：QA
- **接收者**：PM
- **Phase**：Phase 5
- **回合**：第 [N] / 3 回合
- **需要回應**：是 → PM 判定是否通過

### 重測結果
| Bug ID / AC 編號 | 重測結果 | 備註 |
|-----------------|---------|------|
| | ✅ 修復 / ❌ 未修復 / ⚠️ 部分修復 | |

### 回歸測試結果
| 測試案例 | 結果 | 備註 |
|---------|------|------|
| | ✅ / ❌ | |

### 新發現問題（若有）
| Bug ID | 描述 | 嚴重程度 |
|--------|------|---------|
| | | |

### 驗收建議：✅ / ⚠️ / ❌
```

### 3.6 交付與上報信號

#### ESCALATE
```markdown
## SIGNAL: ESCALATE

- **發送者**：PM
- **接收者**：人類
- **Phase**：Phase 5
- **優先級**：⚡ P0
- **需要回應**：是 → 人類決策

### ⚠️ 自治流程已暫停

### 上報原因
[3 回合修復未通過 / 重大需求衝突 / 安全疑慮 / ...]

### 自治流程摘要報告
[使用 PM prompt 定義的「自治流程摘要報告」格式]

### 需要人類決策
1. [待決定事項]
2. [待決定事項]

### PM 建議方案
[附上 PM 的建議]
```

#### DELIVERY_READY
```markdown
## SIGNAL: DELIVERY_READY

- **發送者**：PM
- **接收者**：人類
- **Phase**：Phase 6
- **需要回應**：是 → ACCEPTED 或 REJECTED

### 交付清單
[使用 autonomous-workflow.md Phase 6 定義的交付清單格式]

### 流程摘要
- 總耗時：
- 經歷的 Phase：0 → 1 → 2 → 3 → 4 → [5 →] 6
- 自治修復回合：[N] 回合（若有）
- 總 Bug 數：[N]（已全部 Closed）

### 請人類驗收
```

---

## 4. 信號流範例

### 4.1 Happy Path（一次通過）

```
人類 ──[需求]──▶ Orchestrator ──▶ PM
PM ──[REQ_CLARIFY]──▶ Orchestrator ──▶ 人類
人類 ──[REQ_INPUT]──▶ Orchestrator ──▶ PM
PM ──[需求摘要]──▶ Orchestrator → G0 ✅
Orchestrator ──[G0_PASSED]──▶ PM + FE + BE + QA

PM ──[REQ_READY]──▶ Orchestrator ──▶ FE + BE + QA
FE ──[REQ_REVIEW_OK]──▶ Orchestrator
BE ──[REQ_REVIEW_OK]──▶ Orchestrator
QA ──[REQ_REVIEW_OK]──▶ Orchestrator → G1 ✅

PM ──[CONTRACT_DRAFT]──▶ Orchestrator ──▶ FE + BE
FE ──[CONTRACT_OK]──▶ Orchestrator
BE ──[CONTRACT_OK]──▶ Orchestrator
PM ──[CONTRACT_LOCKED]──▶ Orchestrator ──▶ ALL → G2 ✅

Orchestrator ──[INIT_START]──▶ FE + BE （並行，worktree 隔離）
FE ──[INIT_COMPLETE]──▶ Orchestrator
BE ──[INIT_COMPLETE]──▶ Orchestrator → G2.5 ✅

Orchestrator ──[DEV_START]──▶ FE + BE + QA （並行）
FE ──[DEV_COMPLETE]──▶ Orchestrator
BE ──[DEV_COMPLETE]──▶ Orchestrator
QA ──[TESTCASE_READY]──▶ Orchestrator → G3 ✅

Orchestrator ──[INTEGRATION_START]──▶ ALL
QA ──[QA_START]──▶ Orchestrator
QA ──[QA_REPORT: ✅]──▶ Orchestrator → G4 ✅ → 跳至 Phase 6

PM ──[DELIVERY_READY]──▶ Orchestrator ──▶ 人類
人類 ──[ACCEPTED]──▶ Orchestrator
Orchestrator ──[RELEASE_APPROVED]──▶ ALL → 流程結束 ✅
```

### 4.2 修復路徑（需 1 回合修復）

```
... Phase 4 ...
QA ──[BUG_FOUND: BUG-AUTH-001 [FE]]──▶ Orchestrator ──▶ FE + PM
QA ──[BUG_FOUND: BUG-AUTH-002 [BE]]──▶ Orchestrator ──▶ BE + PM
QA ──[QA_REPORT: ⚠️]──▶ Orchestrator → G4 → Phase 5

PM ──[FIX_ASSIGN: BUG-AUTH-001]──▶ Orchestrator ──▶ FE
PM ──[FIX_ASSIGN: BUG-AUTH-002]──▶ Orchestrator ──▶ BE （並行）
FE ──[FIX_COMPLETE]──▶ Orchestrator ──▶ QA + PM
BE ──[FIX_COMPLETE]──▶ Orchestrator ──▶ QA + PM
QA ──[RETEST_REPORT: ✅]──▶ Orchestrator → G5 ✅ → Phase 6
```
