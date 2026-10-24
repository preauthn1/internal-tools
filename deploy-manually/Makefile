SHELL := /bin/bash
-include .env
AGENT_BACKEND_REPO_DIR  ?= $(HOME)/projects/www/WEB_AgentBackend
API_AGENT_REPO_DIR      ?= $(HOME)/projects/www/API_Agent
TOOLS_DIR := $(CURDIR)/tools/back

.PHONY: help \
        rd/alfawin77/deploy rd/alfawin77/ssh rd/alfawin77/log \
        rd/yes2win/deploy   rd/yes2win/ssh   rd/yes2win/log \
        rd/tpower/deploy    rd/tpower/ssh    rd/tpower/log \
        rd/vworld/deploy    rd/vworld/ssh    rd/vworld/log \
        rd/me77/deploy      rd/me77/ssh      rd/me77/log \
        rd/dev/deploy       rd/dev/ssh       rd/dev/log \
        back/pull-all back/push-all back/merge-all \
        back/deploy-dev back/deploy-prod \
        back/tpower/deploy-prod back/tpower/deploy-dev \
        back/vworld/deploy-prod back/vworld/deploy-dev back/vworld/tag \
        back/wintoto/deploy-prod back/wintoto/deploy-dev \
        back/clear-feature \
        api/pull-all api/push-all api/merge-all \
        api/deploy-dev api/deploy-prod \
        api/tpower/deploy-prod api/tpower/deploy-dev \
        api/vworld/deploy-prod api/vworld/deploy-dev api/vworld/tag \
        api/wintoto/deploy-prod api/wintoto/deploy-dev \
        api/clear-feature

# ===================================================
# 預設 target：顯示說明
# ===================================================
help:
	@printf "\033[1;36m======  rd-studio  ======\033[0m\n"
	@printf "  make rd/alfawin77/deploy    部署 alfawin77\n"
	@printf "  make rd/yes2win/deploy      部署 yes2win\n"
	@printf "  make rd/tpower/deploy       部署 tpower\n"
	@printf "  make rd/vworld/deploy       部署 vworld\n"
	@printf "  make rd/me77/deploy         部署 me77\n"
	@printf "  make rd/dev/deploy          部署 dev\n"
	@printf "  make rd/<站台>/ssh          SSH 進入站台\n"
	@printf "  make rd/<站台>/log          查看站台 git log\n"
	@printf "\n"
	@printf "\033[1;36m====  WEB_AgentBackend  ====\033[0m\n"
	@printf "  make back/pull-all                   拉取所有分支\n"
	@printf "  make back/push-all                   推送所有分支\n"
	@printf "  make back/deploy-dev                 部署全站台 dev 分支（*-develop）\n"
	@printf "  make back/deploy-prod                部署全站台 prod 分支\n"
	@printf "  make back/merge-all                  合併全部站台所有分支（不推送）\n"
	@printf "  make back/tpower/deploy-prod         合併並推送 tpower prod 環境\n"
	@printf "  make back/tpower/deploy-dev          合併並推送 tpower-develop\n"
	@printf "  make back/vworld/deploy-prod         合併並推送 vworld prod 環境\n"
	@printf "  make back/vworld/deploy-dev          合併並推送 vw-develop\n"
	@printf "  make back/vworld/tag                 為 vw-master 建立時間戳記 tag\n"
	@printf "  make back/wintoto/deploy-prod        推送 wintoto prod 環境\n"
	@printf "  make back/wintoto/deploy-dev         合併並推送 wtt-develop\n"
	@printf "  make back/clear-feature              刪除本地 feature/* 分支\n"
	@printf "  make back/clear-feature DRYRUN=1     預覽要刪除的 feature/* 分支\n"
	@printf "\n"
	@printf "\033[1;36m====  API_Agent  ====\033[0m\n"
	@printf "  make api/pull-all                    拉取所有分支\n"
	@printf "  make api/push-all                    推送所有分支\n"
	@printf "  make api/deploy-dev                  部署全站台 dev 分支（*-develop）\n"
	@printf "  make api/deploy-prod                 部署全站台 prod 分支\n"
	@printf "  make api/merge-all                   合併全部站台所有分支（不推送）\n"
	@printf "  make api/tpower/deploy-prod          合併並推送 tpower prod 環境\n"
	@printf "  make api/tpower/deploy-dev           合併並推送 tpower-develop\n"
	@printf "  make api/vworld/deploy-prod          合併並推送 vworld prod 環境\n"
	@printf "  make api/vworld/deploy-dev           合併並推送 vw-develop\n"
	@printf "  make api/vworld/tag                  為 vw-master 建立時間戳記 tag\n"
	@printf "  make api/wintoto/deploy-prod         推送 wintoto prod 環境\n"
	@printf "  make api/wintoto/deploy-dev          合併並推送 wtt-develop\n"
	@printf "  make api/clear-feature               刪除本地 feature/* 分支\n"
	@printf "  make api/clear-feature DRYRUN=1      預覽要刪除的 feature/* 分支\n"

# ===================================================
# rd-studio
# ===================================================

define CONFIRM
	@printf "\033[1;33m即將執行：$(1)\033[0m\n"
	@read -p "確定要繼續？ [y/N] " _c && { [ "$$_c" = "y" ] || [ "$$_c" = "Y" ]; } || { echo "已取消"; exit 1; }
endef

# ---- alfawin77 ----
rd/alfawin77/deploy:
	$(call CONFIRM,rd/alfawin77/deploy)
	@echo "\033[1;33m>>> [RD] DEPLOY alfawin77 <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-alfawin77" --project "gff-alfawin77" --command "sudo -i sh /vwsource/update_live.sh"

rd/alfawin77/ssh:
	$(call CONFIRM,rd/alfawin77/ssh)
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-alfawin77" --project "gff-alfawin77"

rd/alfawin77/log:
	$(call CONFIRM,rd/alfawin77/log)
	@echo "\033[1;36m>>> [RD] GIT LOG alfawin77 <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-alfawin77" --project "gff-alfawin77" --command "cd /vwsource/RD-Studio-01 && sudo git log -1"

# ---- yes2win ----
rd/yes2win/deploy:
	$(call CONFIRM,rd/yes2win/deploy)
	@echo "\033[1;33m>>> [RD] DEPLOY yes2win <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-yes2win" --project "gff-yes2win" --command "sudo -i sh /vwsource/update_live.sh"

rd/yes2win/ssh:
	$(call CONFIRM,rd/yes2win/ssh)
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-yes2win" --project "gff-yes2win"

rd/yes2win/log:
	$(call CONFIRM,rd/yes2win/log)
	@echo "\033[1;36m>>> [RD] GIT LOG yes2win <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-yes2win" --project "gff-yes2win" --command "cd /vwsource/RD-Studio-01 && sudo git log -1"

# ---- tpower ----
rd/tpower/deploy:
	$(call CONFIRM,rd/tpower/deploy)
	@echo "\033[1;33m>>> [RD] DEPLOY tpower <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-a" "rd-sudio-01-tpower" --project "tp-wz3-tpower" --command "sudo -i sh /vwsource/update_live.sh"

rd/tpower/ssh:
	$(call CONFIRM,rd/tpower/ssh)
	gcloud compute ssh --zone "asia-southeast1-a" "rd-sudio-01-tpower" --project "tp-wz3-tpower"

rd/tpower/log:
	$(call CONFIRM,rd/tpower/log)
	@echo "\033[1;36m>>> [RD] GIT LOG tpower <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-a" "rd-sudio-01-tpower" --project "tp-wz3-tpower" --command "cd /vwsource/RD-Studio-01 && sudo git log -1"

# ---- vworld ----
rd/vworld/deploy:
	$(call CONFIRM,rd/vworld/deploy)
	@echo "\033[1;33m>>> [RD] DEPLOY vworld <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-a" "rd-studio-01" --project "tp-wz3-online" --command "sudo -i sh ~vwsource/update_live.sh"

rd/vworld/ssh:
	$(call CONFIRM,rd/vworld/ssh)
	gcloud compute ssh --zone "asia-southeast1-a" "rd-studio-01" --project "tp-wz3-online"

rd/vworld/log:
	$(call CONFIRM,rd/vworld/log)
	@echo "\033[1;36m>>> [RD] GIT LOG vworld <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-a" "rd-studio-01" --project "tp-wz3-online" --command "sudo bash -c 'cd /home/vwsource/RD-Studio-01 && git log -1'"

# ---- me77 ----
rd/me77/deploy:
	$(call CONFIRM,rd/me77/deploy)
	@echo "\033[1;33m>>> [RD] DEPLOY me77 <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-vwb2c" --project "gff-vwb2c" --command "sudo -i sh /vwsource/update_live.sh"

rd/me77/ssh:
	$(call CONFIRM,rd/me77/ssh)
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-vwb2c" --project "gff-vwb2c"

rd/me77/log:
	$(call CONFIRM,rd/me77/log)
	@echo "\033[1;36m>>> [RD] GIT LOG me77 <<<\033[0m"
	gcloud compute ssh --zone "asia-southeast1-c" "rd-studio-vwb2c" --project "gff-vwb2c" --command "cd /vwsource/RD-Studio-01 && sudo git log -1"

# ---- dev ----
rd/dev/deploy:
	$(call CONFIRM,rd/dev/deploy)
	@echo "\033[1;33m>>> [RD] DEPLOY dev <<<\033[0m"
	gcloud compute ssh --zone "asia-east1-a" "rd-studio-01" --project "tp-wz3-rd" --command "sudo -i sh ~vwsource/update_rd.sh"

rd/dev/ssh:
	$(call CONFIRM,rd/dev/ssh)
	gcloud compute ssh --zone "asia-east1-a" "rd-studio-01" --project "tp-wz3-rd"

rd/dev/log:
	$(call CONFIRM,rd/dev/log)
	@echo "\033[1;36m>>> [RD] GIT LOG dev <<<\033[0m"
	gcloud compute ssh --zone "asia-east1-a" "rd-studio-01" --project "tp-wz3-rd" --command "sudo bash -c 'cd /home/vwsource/RD-Studio-01 && git log -10'"

# ===================================================
# WEB_AgentBackend — 拉取 / 推送
# ===================================================
back/pull-all:
	$(call CONFIRM,back/pull-all)
	@printf "\033[1;33m>>> [BACK] PULL 所有分支 <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/pull_branch.sh

back/push-all:
	$(call CONFIRM,back/push-all)
	@printf "\033[1;33m>>> [BACK] PUSH 所有分支 <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/push_branch.sh

# ===================================================
# WEB_AgentBackend — 合併（全部環境）
# ===================================================
back/merge-all:
	$(call CONFIRM,back/merge-all)
	@printf "\033[1;33m>>> [BACK] MERGE ALL <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/merge_tpower.sh
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/merge_vworld.sh
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/merge_wintoto.sh

# ===================================================
# WEB_AgentBackend — tpower
# ===================================================
back/tpower/deploy-prod:
	$(call CONFIRM,back/tpower/deploy-prod)
	@printf "\033[1;32m>>> [BACK] DEPLOY tpower prod：merge + push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="fon666 vplus33 yes2win tpower-new alfawin77-new" sh $(TOOLS_DIR)/merge_tpower.sh
	@printf "\033[1;33m>>> [BACK] PUSH tpower prod 分支 <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="tpower-master fon666 vplus33 yes2win tpower-new alfawin77-new" sh $(TOOLS_DIR)/push_tpower.sh

back/tpower/deploy-dev:
	$(call CONFIRM,back/tpower/deploy-dev)
	@printf "\033[1;32m>>> [BACK] DEPLOY tpower dev：merge + push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="tpower-develop" sh $(TOOLS_DIR)/merge_tpower.sh
	@printf "\033[1;33m>>> [BACK] PUSH tpower-develop <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="tpower-develop" sh $(TOOLS_DIR)/push_tpower.sh

# ===================================================
# WEB_AgentBackend — vworld
# ===================================================
back/vworld/deploy-prod:
	$(call CONFIRM,back/vworld/deploy-prod)
	@printf "\033[1;32m>>> [BACK] DEPLOY vworld prod：merge + push + tag <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="vwb2c" sh $(TOOLS_DIR)/merge_vworld.sh
	@printf "\033[1;33m>>> [BACK] PUSH vworld prod 分支 <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="vw-master vwb2c" sh $(TOOLS_DIR)/push_vworld.sh
	$(MAKE) back/vworld/tag

back/vworld/deploy-dev:
	$(call CONFIRM,back/vworld/deploy-dev)
	@printf "\033[1;32m>>> [BACK] DEPLOY vworld dev：merge + push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="vw-develop" sh $(TOOLS_DIR)/merge_vworld.sh
	@printf "\033[1;33m>>> [BACK] PUSH vw-develop <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="vw-develop" sh $(TOOLS_DIR)/push_vworld.sh

back/vworld/tag:
	$(call CONFIRM,back/vworld/tag)
	@printf "\033[1;33m>>> [BACK] TAG vw-master 建立時間戳記 tag <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/tag_vworld.sh

# ===================================================
# WEB_AgentBackend — wintoto
# ===================================================
back/wintoto/deploy-prod:
	$(call CONFIRM,back/wintoto/deploy-prod)
	@printf "\033[1;32m>>> [BACK] DEPLOY wintoto prod：push <<<\033[0m\n"
	@printf "\033[1;33m>>> [BACK] PUSH wintoto prod 分支 <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="wintoto" sh $(TOOLS_DIR)/push_wintoto.sh

back/wintoto/deploy-dev:
	$(call CONFIRM,back/wintoto/deploy-dev)
	@printf "\033[1;32m>>> [BACK] DEPLOY wintoto dev：merge + push <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/merge_wintoto.sh
	@printf "\033[1;33m>>> [BACK] PUSH wtt-develop <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && BRANCHES="wtt-develop" sh $(TOOLS_DIR)/push_wintoto.sh

# ===================================================
# WEB_AgentBackend — Dev / Prod deploy（跨環境）
# ===================================================
back/deploy-dev:
	$(call CONFIRM,back/deploy-dev)
	@printf "\033[1;32m>>> [BACK] DEPLOY DEV：tpower-develop / vw-develop / wtt-develop <<<\033[0m\n"
	$(MAKE) back/tpower/deploy-dev
	$(MAKE) back/vworld/deploy-dev
	$(MAKE) back/wintoto/deploy-dev

back/deploy-prod:
	$(call CONFIRM,back/deploy-prod)
	@printf "\033[1;32m>>> [BACK] DEPLOY PROD：fon666 vplus33 yes2win tpower-new alfawin77-new / vwb2c / wintoto <<<\033[0m\n"
	$(MAKE) back/tpower/deploy-prod
	$(MAKE) back/vworld/deploy-prod
	$(MAKE) back/wintoto/deploy-prod

# ===================================================
# WEB_AgentBackend — 清理 feature 分支
# 用法：
#   make back/clear-feature           → 直接刪除
#   make back/clear-feature DRYRUN=1  → 僅預覽，不刪除
# ===================================================
back/clear-feature:
	$(call CONFIRM,back/clear-feature)
	@printf "\033[1;33m>>> [BACK] CLEAR feature/* 分支$(if $(DRYRUN), (dryrun 模式),) <<<\033[0m\n"
	cd $(AGENT_BACKEND_REPO_DIR) && sh $(TOOLS_DIR)/clear_feature_branch.sh $(if $(DRYRUN),--dryrun,)

# ===================================================
# API_Agent — 拉取 / 推送
# ===================================================
api/pull-all:
	$(call CONFIRM,api/pull-all)
	@printf "\033[1;33m>>> [API] PULL 所有分支 <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && sh $(TOOLS_DIR)/pull_branch.sh

api/push-all:
	$(call CONFIRM,api/push-all)
	@printf "\033[1;33m>>> [API] PUSH 所有分支 <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && sh $(TOOLS_DIR)/push_branch.sh

# ===================================================
# API_Agent — 合併（全部環境）
# ===================================================
api/merge-all:
	$(call CONFIRM,api/merge-all)
	@printf "\033[1;33m>>> [API] MERGE ALL <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && sh $(TOOLS_DIR)/merge_tpower.sh
	cd $(API_AGENT_REPO_DIR) && sh $(TOOLS_DIR)/merge_vworld.sh
	cd $(API_AGENT_REPO_DIR) && sh $(TOOLS_DIR)/merge_wintoto.sh

# ===================================================
# API_Agent — tpower
# ===================================================
api/tpower/deploy-prod:
	$(call CONFIRM,api/tpower/deploy-prod)
	@printf "\033[1;32m>>> [API] DEPLOY tpower prod：merge + push <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="fon666 vplus33 yes2win tpower-new alfawin77-new" sh $(TOOLS_DIR)/merge_tpower.sh
	@printf "\033[1;33m>>> [API] PUSH tpower prod 分支 <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="tpower-master fon666 vplus33 yes2win tpower-new alfawin77-new" sh $(TOOLS_DIR)/push_tpower.sh

api/tpower/deploy-dev:
	$(call CONFIRM,api/tpower/deploy-dev)
	@printf "\033[1;32m>>> [API] DEPLOY tpower dev：merge + push <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="tpower-develop" sh $(TOOLS_DIR)/merge_tpower.sh
	@printf "\033[1;33m>>> [API] PUSH tpower-develop <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="tpower-develop" sh $(TOOLS_DIR)/push_tpower.sh

# ===================================================
# API_Agent — vworld
# ===================================================
api/vworld/deploy-prod:
	$(call CONFIRM,api/vworld/deploy-prod)
	@printf "\033[1;32m>>> [API] DEPLOY vworld prod：merge + push + tag <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="vwb2c" sh $(TOOLS_DIR)/merge_vworld.sh
	@printf "\033[1;33m>>> [API] PUSH vworld prod 分支 <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="vw-master vwb2c" sh $(TOOLS_DIR)/push_vworld.sh
	$(MAKE) api/vworld/tag

api/vworld/deploy-dev:
	$(call CONFIRM,api/vworld/deploy-dev)
	@printf "\033[1;32m>>> [API] DEPLOY vworld dev：merge + push <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="vw-develop" sh $(TOOLS_DIR)/merge_vworld.sh
	@printf "\033[1;33m>>> [API] PUSH vw-develop <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="vw-develop" sh $(TOOLS_DIR)/push_vworld.sh

api/vworld/tag:
	$(call CONFIRM,api/vworld/tag)
	@printf "\033[1;33m>>> [API] TAG vw-master 建立時間戳記 tag <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && sh $(TOOLS_DIR)/tag_vworld.sh

# ===================================================
# API_Agent — wintoto
# ===================================================
api/wintoto/deploy-prod:
	$(call CONFIRM,api/wintoto/deploy-prod)
	@printf "\033[1;32m>>> [API] DEPLOY wintoto prod：push <<<\033[0m\n"
	@printf "\033[1;33m>>> [API] PUSH wintoto prod 分支 <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="wintoto" sh $(TOOLS_DIR)/push_wintoto.sh

api/wintoto/deploy-dev:
	$(call CONFIRM,api/wintoto/deploy-dev)
	@printf "\033[1;32m>>> [API] DEPLOY wintoto dev：merge + push <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && sh $(TOOLS_DIR)/merge_wintoto.sh
	@printf "\033[1;33m>>> [API] PUSH wtt-develop <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && BRANCHES="wtt-develop" sh $(TOOLS_DIR)/push_wintoto.sh

# ===================================================
# API_Agent — Dev / Prod deploy（跨環境）
# ===================================================
api/deploy-dev:
	$(call CONFIRM,api/deploy-dev)
	@printf "\033[1;32m>>> [API] DEPLOY DEV：tpower-develop / vw-develop / wtt-develop <<<\033[0m\n"
	$(MAKE) api/tpower/deploy-dev
	$(MAKE) api/vworld/deploy-dev
	$(MAKE) api/wintoto/deploy-dev

api/deploy-prod:
	$(call CONFIRM,api/deploy-prod)
	@printf "\033[1;32m>>> [API] DEPLOY PROD：fon666 vplus33 yes2win tpower-new alfawin77-new / vwb2c / wintoto <<<\033[0m\n"
	$(MAKE) api/tpower/deploy-prod
	$(MAKE) api/vworld/deploy-prod
	$(MAKE) api/wintoto/deploy-prod

# ===================================================
# API_Agent — 清理 feature 分支
# 用法：
#   make api/clear-feature           → 直接刪除
#   make api/clear-feature DRYRUN=1  → 僅預覽，不刪除
# ===================================================
api/clear-feature:
	$(call CONFIRM,api/clear-feature)
	@printf "\033[1;33m>>> [API] CLEAR feature/* 分支$(if $(DRYRUN), (dryrun 模式),) <<<\033[0m\n"
	cd $(API_AGENT_REPO_DIR) && sh $(TOOLS_DIR)/clear_feature_branch.sh $(if $(DRYRUN),--dryrun,)
