# Back Deploy Dev/Prod Split Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 將 `tools/` 依專案分層（`tools/back/`），並新增 `back/deploy-dev` 與 `back/deploy-prod` 兩個 Makefile target，讓 dev / prod 部署可以獨立執行。

**Architecture:** 在各 merge / push shell script 中改用 `${BRANCHES:-預設值}` 讓外部可透過環境變數覆寫目標分支，Makefile 利用 inline env var 傳入 dev 或 prod 的分支清單，達到不新增 script 的零重複方案。

**Tech Stack:** GNU Make, Bash, git

---

## 分支對照表

| 環境 | Source | Dev 分支 | Prod 分支 |
|---|---|---|---|
| tpower | `tpower-master` | `tpower-develop` | `fon666 vplus33 yes2win tpower-new alfawin77-new` |
| vworld | `vw-master` | `vw-develop` | `vwb2c` |
| wintoto | `wintoto`（即 prod source）| `wtt-develop` | 無 merge target，只 push `wintoto` |

---

## 異動檔案清單

| 動作 | 路徑 | 說明 |
|---|---|---|
| git mv | `tools/*.sh` → `tools/back/*.sh` | 依專案分層 |
| Modify | `tools/back/merge_tpower.sh` | BRANCHES 改為可覆寫 |
| Modify | `tools/back/merge_vworld.sh` | BRANCHES 改為可覆寫 |
| Modify | `tools/back/push_tpower.sh` | BRANCHES 改為可覆寫 |
| Modify | `tools/back/push_vworld.sh` | BRANCHES 改為可覆寫 |
| Modify | `tools/back/push_wintoto.sh` | BRANCHES 改為可覆寫 |
| Modify | `tools/back/merge_push.sh` | 更新內部路徑 |
| Modify | `Makefile` | 更新 TOOLS_DIR、新增 target、更新 .PHONY 與 help |
| Modify | `README.md` | 補充新 target 說明 |

---

## Task 1：移動 tools/ 至 tools/back/

**Files:**
- `tools/*.sh` → `tools/back/*.sh`（git mv）
- Modify: `Makefile`（更新 TOOLS_DIR）

- [ ] **Step 1：用 git mv 搬移所有 sh 檔**

```bash
mkdir -p tools/back
git mv tools/clear_feature_branch.sh tools/back/
git mv tools/merge_push.sh          tools/back/
git mv tools/merge_tpower.sh        tools/back/
git mv tools/merge_vworld.sh        tools/back/
git mv tools/merge_wintoto.sh       tools/back/
git mv tools/pull_branch.sh         tools/back/
git mv tools/push_branch.sh         tools/back/
git mv tools/push_tpower.sh         tools/back/
git mv tools/push_vworld.sh         tools/back/
git mv tools/push_wintoto.sh        tools/back/
git mv tools/tag_vworld.sh          tools/back/
```

- [ ] **Step 2：更新 Makefile 的 TOOLS_DIR**

在 `Makefile` 第 4 行，將：
```makefile
TOOLS_DIR := $(CURDIR)/tools
```
改為：
```makefile
TOOLS_DIR := $(CURDIR)/tools/back
```

- [ ] **Step 3：驗證 Makefile 路徑語法正確**

```bash
make --dry-run back/pull 2>&1 | head -5
```
預期輸出包含：`tools/back/pull_branch.sh`

- [ ] **Step 4：Commit**

```bash
git add Makefile
git commit -m "refactor: 將 tools/ 移至 tools/back/ 以支援多專案擴充"
```

---

## Task 2：更新 merge_push.sh 內部路徑

**Files:**
- Modify: `tools/back/merge_push.sh`

- [ ] **Step 1：確認目前內容**

```bash
cat tools/back/merge_push.sh
```

目前內容：
```bash
#!/bin/bash
sh tools/merge_tpower.sh && sh tools/merge_vworld.sh && sh tools/merge_wintoto.sh && sh tools/push_branch.sh
```

- [ ] **Step 2：更新路徑**

將 `tools/back/merge_push.sh` 改為：
```bash
#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
sh "$SCRIPT_DIR/merge_tpower.sh" && \
sh "$SCRIPT_DIR/merge_vworld.sh" && \
sh "$SCRIPT_DIR/merge_wintoto.sh" && \
sh "$SCRIPT_DIR/push_branch.sh"
```

- [ ] **Step 3：語法檢查**

```bash
bash -n tools/back/merge_push.sh
```
預期：無輸出（語法正確）

- [ ] **Step 4：Commit**

```bash
git add tools/back/merge_push.sh
git commit -m "fix: 更新 merge_push.sh 使用相對於腳本的路徑"
```

---

## Task 3：修改 merge scripts 支援 BRANCHES 環境變數覆寫

**Files:**
- Modify: `tools/back/merge_tpower.sh`
- Modify: `tools/back/merge_vworld.sh`
- `tools/back/merge_wintoto.sh`：不需修改（只有 `wtt-develop` 一個 dev 目標，prod 無 merge）

- [ ] **Step 1：修改 merge_tpower.sh**

將第 4 行：
```bash
BRANCHES="fon666 vplus33 yes2win tpower-new alfawin77-new tpower-develop"
```
改為：
```bash
BRANCHES="${BRANCHES:-fon666 vplus33 yes2win tpower-new alfawin77-new tpower-develop}"
```

- [ ] **Step 2：修改 merge_vworld.sh**

將第 4 行：
```bash
BRANCHES="vw-develop vwb2c"
```
改為：
```bash
BRANCHES="${BRANCHES:-vw-develop vwb2c}"
```

- [ ] **Step 3：語法檢查**

```bash
bash -n tools/back/merge_tpower.sh && bash -n tools/back/merge_vworld.sh
```
預期：無輸出

- [ ] **Step 4：Commit**

```bash
git add tools/back/merge_tpower.sh tools/back/merge_vworld.sh
git commit -m "feat: merge scripts 支援 BRANCHES 環境變數覆寫"
```

---

## Task 4：修改 push scripts 支援 BRANCHES 環境變數覆寫

**Files:**
- Modify: `tools/back/push_tpower.sh`
- Modify: `tools/back/push_vworld.sh`
- Modify: `tools/back/push_wintoto.sh`

- [ ] **Step 1：修改 push_tpower.sh**

將第 3 行：
```bash
BRANCHES="tpower-master fon666 vplus33 yes2win tpower-new alfawin77-new tpower-develop"
```
改為：
```bash
BRANCHES="${BRANCHES:-tpower-master fon666 vplus33 yes2win tpower-new alfawin77-new tpower-develop}"
```

- [ ] **Step 2：修改 push_vworld.sh**

將第 3 行：
```bash
BRANCHES="vw-master vw-develop vwb2c"
```
改為：
```bash
BRANCHES="${BRANCHES:-vw-master vw-develop vwb2c}"
```

- [ ] **Step 3：修改 push_wintoto.sh**

將第 3 行：
```bash
BRANCHES="wintoto wtt-develop"
```
改為：
```bash
BRANCHES="${BRANCHES:-wintoto wtt-develop}"
```

- [ ] **Step 4：語法檢查**

```bash
bash -n tools/back/push_tpower.sh && \
bash -n tools/back/push_vworld.sh && \
bash -n tools/back/push_wintoto.sh
```
預期：無輸出

- [ ] **Step 5：Commit**

```bash
git add tools/back/push_tpower.sh tools/back/push_vworld.sh tools/back/push_wintoto.sh
git commit -m "feat: push scripts 支援 BRANCHES 環境變數覆寫"
```

---

## Task 5：新增 back/deploy-dev 與 back/deploy-prod

**Files:**
- Modify: `Makefile`

dev 流程：
- tpower：merge `tpower-master → tpower-develop`，push `tpower-develop`
- vworld：merge `vw-master → vw-develop`，push `vw-develop`
- wintoto：merge `wintoto → wtt-develop`，push `wtt-develop`

prod 流程：
- tpower：merge `tpower-master → fon666 vplus33 yes2win tpower-new alfawin77-new`，push 含 source
- vworld：merge `vw-master → vwb2c`，push 含 source，建立 tag
- wintoto：無 merge，只 push `wintoto`

- [ ] **Step 1：在 .PHONY 加入新 target**

在 `Makefile` .PHONY 的最後一行加入：
```makefile
        back/deploy-dev back/deploy-prod
```

- [ ] **Step 2：在 help 區塊補充說明**

在 `back/deploy-wintoto` 那行 printf 之後插入：
```makefile
	@printf "  make back/deploy-dev             部署全環境 dev 分支（*-develop）\n"
	@printf "  make back/deploy-prod            部署全環境 prod 分支\n"
```

- [ ] **Step 3：在 Makefile 的 back/deploy-wintoto 之後新增兩個 target**

```makefile
back/deploy-dev:
	@printf "\033[1;32m>>> [BACK] DEPLOY DEV：tpower-develop / vw-develop / wtt-develop <<<\033[0m\n"
	@printf "\033[1;33m>>> [BACK] tpower dev：merge <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="tpower-develop" sh $(TOOLS_DIR)/merge_tpower.sh
	@printf "\033[1;33m>>> [BACK] tpower dev：push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="tpower-develop" sh $(TOOLS_DIR)/push_tpower.sh
	@printf "\033[1;33m>>> [BACK] vworld dev：merge <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="vw-develop" sh $(TOOLS_DIR)/merge_vworld.sh
	@printf "\033[1;33m>>> [BACK] vworld dev：push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="vw-develop" sh $(TOOLS_DIR)/push_vworld.sh
	@printf "\033[1;33m>>> [BACK] wintoto dev：merge <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/merge_wintoto.sh
	@printf "\033[1;33m>>> [BACK] wintoto dev：push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="wtt-develop" sh $(TOOLS_DIR)/push_wintoto.sh

back/deploy-prod:
	@printf "\033[1;32m>>> [BACK] DEPLOY PROD：fon666 vplus33 yes2win tpower-new alfawin77-new / vwb2c / wintoto <<<\033[0m\n"
	@printf "\033[1;33m>>> [BACK] tpower prod：merge <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="fon666 vplus33 yes2win tpower-new alfawin77-new" sh $(TOOLS_DIR)/merge_tpower.sh
	@printf "\033[1;33m>>> [BACK] tpower prod：push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="tpower-master fon666 vplus33 yes2win tpower-new alfawin77-new" sh $(TOOLS_DIR)/push_tpower.sh
	@printf "\033[1;33m>>> [BACK] vworld prod：merge <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="vwb2c" sh $(TOOLS_DIR)/merge_vworld.sh
	@printf "\033[1;33m>>> [BACK] vworld prod：push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="vw-master vwb2c" sh $(TOOLS_DIR)/push_vworld.sh
	$(MAKE) back/tag-vworld
	@printf "\033[1;33m>>> [BACK] wintoto prod：push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="wintoto" sh $(TOOLS_DIR)/push_wintoto.sh
```

- [ ] **Step 4：驗證 Makefile 語法**

```bash
make --dry-run back/deploy-dev 2>&1 | head -30
make --dry-run back/deploy-prod 2>&1 | head -30
```
預期：印出各步驟 cd + sh 指令，無 `Error` 字樣

- [ ] **Step 5：Commit**

```bash
git add Makefile
git commit -m "feat: 新增 back/deploy-dev 與 back/deploy-prod target"
```

---

## Task 6：更新 README

**Files:**
- Modify: `README.md`

- [ ] **Step 1：在「單一環境部署」表格下方新增段落**

在 README.md 的「### 單一環境部署」段落之後加入：

```markdown
### Dev / Prod 部署

| 指令 | 說明 | 分支 |
|---|---|---|
| `make back/deploy-dev` | 部署全環境 dev 分支 | tpower-develop, vw-develop, wtt-develop |
| `make back/deploy-prod` | 部署全環境 prod 分支 | fon666, vplus33, yes2win, tpower-new, alfawin77-new, vwb2c, wintoto |
```

- [ ] **Step 2：Commit**

```bash
git add README.md
git commit -m "docs: 補充 back/deploy-dev 與 back/deploy-prod 說明"
```

---

## Self-Review

**Spec coverage：**
- [x] tools 依專案分層（Task 1）
- [x] merge scripts 支援 BRANCHES 覆寫（Task 3）
- [x] push scripts 支援 BRANCHES 覆寫（Task 4）
- [x] back/deploy-dev（Task 5）
- [x] back/deploy-prod（Task 5）
- [x] wintoto prod 無 merge，只 push wintoto source branch（Task 5 Step 3）
- [x] vworld prod 包含 tag（Task 5 Step 3）

**Placeholder scan：** 無 TBD / TODO

**Type consistency：** BRANCHES 變數名稱在所有 script 與 Makefile inline 傳值中一致
