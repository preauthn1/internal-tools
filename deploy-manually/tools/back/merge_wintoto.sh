#!/bin/bash

# 定義要合併的分支列表（使用空格分隔）
BRANCHES="wtt-develop"
SOURCE_BRANCH="wintoto"

# 遍歷分支列表並進行合併
for branch in $BRANCHES; do
    echo "開始處理分支: $branch"
    
    # 切換到目標分支
    git checkout $branch
    if [ $? -ne 0 ]; then
        echo "切換到 $branch 分支失敗"
        continue
    fi
    
    # 拉取目標分支最新代碼
    git pull origin $branch
    if [ $? -ne 0 ]; then
        echo "拉取 $branch 最新代碼失敗"
        continue
    fi
    
    # 合併 tpower-master 到當前分支，並提供合併訊息
    git merge $SOURCE_BRANCH -m "Merge branch '$SOURCE_BRANCH' into $branch"
    if [ $? -ne 0 ]; then
        echo "合併 $SOURCE_BRANCH 到 $branch 失敗"
        continue
    fi
    
    echo "成功將 $SOURCE_BRANCH 合併到 $branch"
done

echo "所有分支處理完成"