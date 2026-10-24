#!/bin/sh

# 設定分支名稱
BRANCH="vw-master"

# 切換到目標分支
git checkout $BRANCH

# 檢查切換是否成功
if [ $? -ne 0 ]; then
    echo "切換到 $BRANCH 分支失敗"
    exit 1
fi

# 確保本地分支為最新
git pull origin $BRANCH

# 取得現在時間（格式：_onlinev2_YYYYMMDDHHMM）
TAG_NAME="_online_$(date +%Y%m%d%H%M)"

# 建立 tag
git tag "$TAG_NAME"

# push tag 到遠端
git push origin "$TAG_NAME"

echo "Tag $TAG_NAME 已建立並推送成功。"
