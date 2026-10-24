# deploy-manually

> 本目录原始 README 见文末，本节为快照补充信息。

原始路径：`deploy-manually`　|　文件 **15**　|　主语言：Shell 11

## 目录构成

```
tools/back/  — 11
docs/superpowers/  — 1
```

## 关键文件

- CI/构建：`Makefile`

---

## 原始说明

# Deploy Manually

透過 `make` 指令對專案 `WEB_AgentBackend`、`API_Agent` 執行部署、分支管理等操作。  
在執行操作之前，請先確認你的 feature 是否已經放在本地的 master 分支了（`tpower-master` / `vw-master` / `wintoto`）。

## 環境需求

使用前請確認已安裝以下工具：

| 工具 | 用途 | 安裝指令 |
|---|---|---|
| `make` | 執行 Makefile 指令 | `sudo apt install make` |
| `git` | 分支合併與推送 | `sudo apt install git` |
| `gcloud` | SSH 登入與部署 rd-studio | 參考 [Google Cloud SDK 安裝文件](https://example.invalid/sdk/docs/install) |

---

## 環境設定

各專案目錄路徑可透過 `.env` 檔案自訂：

```bash
cp .env.example .env
```

編輯 `.env`：

```env
# WEB_AgentBackend 的路徑
AGENT_BACKEND_REPO_DIR = /your/path/to/WEB_AgentBackend

# API_Agent 的路徑
API_AGENT_REPO_DIR     = /your/path/to/API_Agent
```

也可以在執行時直接覆蓋：

```bash
make back/tpower/deploy-prod AGENT_BACKEND_REPO_DIR=/your/path/to/WEB_AgentBackend
make api/tpower/deploy-prod  API_AGENT_REPO_DIR=/your/path/to/API_Agent
```

> `.env` 已加入 `.gitignore`，不會被提交。

---

> **不確定指令會做什麼？** 可以先用 `make -n <target>` 預覽將執行的指令，不會實際執行任何操作。
> 例如：`make -n rd/alfawin77/deploy`

指令格式以 `專案/動作` 或 `專案/站台/動作` 區分：

- `rd/` — RD Studio（gcloud 部署）
- `back/` — WEB_AgentBackend（git merge / push）
- `back/<站台>/` — 指定 WEB_AgentBackend `tpower` / `vworld` / `wintoto` 等站台的操作
- `api/` — API_Agent（git merge / push）
- `api/<站台>/` — 指定 API_Agent `tpower` / `vworld` / `wintoto` 等站台的操作

---

## rd-studio

對各站台執行上版、SSH 登入、查看 git log。

| 站台 | deploy | ssh | log |
|---|---|---|---|
| alfawin77 | `make rd/alfawin77/deploy` | `make rd/alfawin77/ssh` | `make rd/alfawin77/log` |
| yes2win | `make rd/yes2win/deploy` | `make rd/yes2win/ssh` | `make rd/yes2win/log` |
| tpower | `make rd/tpower/deploy` | `make rd/tpower/ssh` | `make rd/tpower/log` |
| vworld | `make rd/vworld/deploy` | `make rd/vworld/ssh` | `make rd/vworld/log` |
| me77 | `make rd/me77/deploy` | `make rd/me77/ssh` | `make rd/me77/log` |
| dev | `make rd/dev/deploy` | `make rd/dev/ssh` | `make rd/dev/log` |

---

## WEB_AgentBackend

對 WEB_AgentBackend 執行分支合併與推送。

### 工作流程

1. 將 feature branch 合併到各站台的 master 分支（`tpower-master` / `vw-master` / `wintoto`）
2. 執行 `make back/deploy-dev`：將各站台 master → develop 分支並推送
3. 請 QA 驗證 dev 環境
4. 驗證通過後執行 `make back/deploy-prod`：將各站台 master → prod 分支並推送
5. 請 QA 驗證 prod 環境

### Dev / Prod 批次部署

| 指令 | 說明 | 分支 |
|---|---|---|
| `make back/deploy-dev` | 部署全站台 dev 分支 | tpower-develop, vw-develop, wtt-develop |
| `make back/deploy-prod` | 部署全站台 prod 分支 | fon666, vplus33, yes2win, tpower-new, alfawin77-new, vwb2c, wintoto |

### 單一站台部署（merge + push）

| 站台 | 指令 | 說明 |
|---|---|---|
| tpower | `make back/tpower/deploy-prod` | tpower-master → fon666, vplus33, yes2win, tpower-new, alfawin77-new |
| tpower | `make back/tpower/deploy-dev` | tpower-master → tpower-develop |
| vworld | `make back/vworld/deploy-prod` | vw-master → vwb2c（完成後自動建立時間戳記 tag） |
| vworld | `make back/vworld/deploy-dev` | vw-master → vw-develop |
| wintoto | `make back/wintoto/deploy-prod` | push wintoto |
| wintoto | `make back/wintoto/deploy-dev` | wintoto → wtt-develop |

### 批次操作

| 指令 | 說明 |
|---|---|
| `make back/merge-all` | 合併全部站台所有分支（不推送） |
| `make back/pull-all` | 拉取所有分支最新代碼 |
| `make back/push-all` | 推送所有分支到遠端 |

### 其他

| 指令 | 說明 |
|---|---|
| `make back/vworld/tag` | 為 vw-master 建立時間戳記 tag |
| `make back/clear-feature` | 刪除本地所有 feature/* 分支 |
| `make back/clear-feature DRYRUN=1` | 預覽要刪除的 feature/* 分支（不實際刪除） |

---

## API_Agent

對 API_Agent 執行分支合併與推送，指令結構與 WEB_AgentBackend 完全相同，前綴改為 `api/`。

### Dev / Prod 批次部署

| 指令 | 說明 | 分支 |
|---|---|---|
| `make api/deploy-dev` | 部署全站台 dev 分支 | tpower-develop, vw-develop, wtt-develop |
| `make api/deploy-prod` | 部署全站台 prod 分支 | fon666, vplus33, yes2win, tpower-new, alfawin77-new, vwb2c, wintoto |

### 單一站台部署（merge + push）

| 站台 | 指令 | 說明 |
|---|---|---|
| tpower | `make api/tpower/deploy-prod` | tpower-master → fon666, vplus33, yes2win, tpower-new, alfawin77-new |
| tpower | `make api/tpower/deploy-dev` | tpower-master → tpower-develop |
| vworld | `make api/vworld/deploy-prod` | vw-master → vwb2c（完成後自動建立時間戳記 tag） |
| vworld | `make api/vworld/deploy-dev` | vw-master → vw-develop |
| wintoto | `make api/wintoto/deploy-prod` | push wintoto |
| wintoto | `make api/wintoto/deploy-dev` | wintoto → wtt-develop |

### 批次操作

| 指令 | 說明 |
|---|---|
| `make api/merge-all` | 合併全部站台所有分支（不推送） |
| `make api/pull-all` | 拉取所有分支最新代碼 |
| `make api/push-all` | 推送所有分支到遠端 |

### 其他

| 指令 | 說明 |
|---|---|
| `make api/vworld/tag` | 為 vw-master 建立時間戳記 tag |
| `make api/clear-feature` | 刪除本地所有 feature/* 分支 |
| `make api/clear-feature DRYRUN=1` | 預覽要刪除的 feature/* 分支（不實際刪除） |

