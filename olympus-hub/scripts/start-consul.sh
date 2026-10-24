#!/bin/bash

# 檢查環境變數
if [ -z "$INITIAL_MANAGEMENT_TOKEN" ]; then
    echo "錯誤: 請先設定 INITIAL_MANAGEMENT_TOKEN 環境變數。"
    echo "範例: export INITIAL_MANAGEMENT_TOKEN=[REDACTED]"
    exit 1
fi

echo "正在啟動 Consul 容器..."

docker run -d \
  --name consul \
  -p 8500:8500 \
  -p 8600:8600/udp \
  -v $(pwd)/realms:/consul/realms \
  hashicorp/consul:1.19 \
  agent -server -bootstrap-expect=1 \
  -ui -client=0.0.0.0 \
  -data-dir=/consul/data \
  -hcl="acl { enabled = true, default_policy = \"deny\", enable_token_persistence = true, tokens { initial_management = \"${INITIAL_MANAGEMENT_TOKEN}\" } }"

if [ $? -eq 0 ]; then
    echo "------------------------------------------------"
    echo "✨ Consul 服務已成功啟動！"
    echo "🔗 介面位址: http://localhost:8500/ui"
    echo "🔑 認證 Token: [REDACTED]"
    echo "------------------------------------------------"
else
    echo "❌ 啟動失敗，請檢查 Docker 狀態或容器名稱是否重複。"
fi
