# CHANGELOG

## [Unreleased] - 2026-04-23

### Changed
- `rd/` 指令格式從 `rd/操作-站台` 改為 `rd/站台/操作`（例如 `rd/alfawin77/deploy`）
- 每個 `rd/` 指令執行前加入確認提示，顯示操作與站台名稱，輸入 `y/Y` 才繼續執行

### Docs
- README 新增 `make -n <target>` 說明，提示使用者可先預覽指令再執行

---

## [1.2.0] - 2026-04-22

### Added
- 新增 `api/` 系列指令，對應 API_Agent 專案（pull-all、push-all、merge-all、deploy-dev、deploy-prod、各站台 deploy-prod/dev、clear-feature）

### Changed
- 統一全站術語：「環境」→「站台」

---

## [1.1.0] - 2026-04-22

### Added
- 新增 `back/deploy-dev`、`back/deploy-prod` 批次部署全站台指令
- merge / push scripts 支援 `BRANCHES` 環境變數覆寫預設分支清單

### Changed
- `back/` 指令格式改為 `back/站台/動作`（例如 `back/tpower/deploy-prod`）
- 拆分 `back/deploy-dev`、`back/deploy-prod` 取代舊的批次 target
- 重命名 `pull-all` / `push-all`，拆分 `merge-*-dev/prod`
- `tools/` 目錄移至 `tools/back/`，預留多專案擴充空間

### Fixed
- `merge_push.sh` 改用相對於腳本自身的路徑，修正執行位置不同導致的路徑錯誤

---

## [1.0.0] - 2026-04-20

### Added
- 初始化部署腳本：`rd/`（gcloud RD Studio 部署）與 `back/`（WEB_AgentBackend 分支管理）基礎指令
