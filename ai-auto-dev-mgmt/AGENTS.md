# NewBestGate — 自治並行開發管理流程

本專案定義了一套 **AI Agent 自治並行開發工作流程**，用於端對端自主軟體交付。

## 核心概念

6 個專業 AI Agent 透過結構化信號協定協作，經過 6 個階段的品質門控，從需求到交付全流程自主運作，僅在關鍵決策點（需求確認、升級、最終驗收）才需要人類介入。

## Agent 角色

| Agent | 職責 |
|-------|------|
| Orchestrator | 中央協調者，管理階段轉換、信號路由、品質門控 |
| Product Owner (PM) | 需求分析、User Story 拆解、驗收標準定義 |
| Frontend Engineer | Vue 3 + TypeScript + TDD 前端開發 |
| Backend Engineer | Laravel + TDD 後端開發 |
| Code Reviewer | 獨立第三方審查：需求覆蓋、路由可達性、API 合約合規 |
| QA Engineer | 測試計劃、測試執行、驗收建議 |

## 工作流程階段

1. **Phase 0** — 需求輸入與人類確認 (G0)
2. **Phase 1** — 需求分析與文件產出 (G1)
3. **Phase 2** — API 合約制定與鎖定 (G2)
4. **Phase 2.5** — 專案初始化與環境建置 (G2.5)
5. **Phase 3** — 前後端並行開發 + Code Review (G3)
6. **Phase 4** — QA 驗證 (G4)
7. **Phase 5** — 自我修復循環（最多 3 輪，失敗則升級人類）(G5)
8. **Phase 6** — 交付 (G6)

## 文件結構

所有工作流程定義文件位於 `autonomous-dev-workflow/` 目錄：

```
autonomous-dev-workflow/
├── orchestrator-agent-system-prompt.md    # Orchestrator Agent 系統提示
├── product-owner-agent-system-prompt.md   # PM Agent 系統提示
├── frontend-engineer-system-prompt.md     # 前端工程師 Agent 系統提示
├── backend-engineer-agent-system-prompt.md # 後端工程師 Agent 系統提示
├── code-reviewer-agent-system-prompt.md   # Code Reviewer Agent 系統提示
├── qa-engineer-agent-system-prompt.md     # QA 工程師 Agent 系統提示
├── autonomous-workflow.md                 # 完整工作流程規範
├── signal-protocol.md                     # Agent 間信號通訊協定
├── quality-gate-checklist.md              # 品質門檢查清單 (G0-G6)
├── api-contract-template.md               # API 合約模板
└── flow-state-tracker-template.md         # 流程狀態追蹤模板
```

## 關鍵機制

- **品質門控 (Quality Gates)** — G0 到 G6，每個階段必須通過才能推進
- **信號協定 (Signal Protocol)** — 25+ 種結構化信號，分 P0-P3 優先級
- **自我修復 (Self-Healing)** — QA 發現問題後自動修復，最多 3 輪
- **TDD 優先** — 所有工程師必須先寫測試再實作
- **並行開發** — 前後端在 Phase 3 同時開發，透過 API 合約解耦
