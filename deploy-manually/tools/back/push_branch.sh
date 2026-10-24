#!/bin/bash

# 定義要push的分支列表（使用空格分隔）
BRANCHES="wtt-develop vwb2c fon666 vplus33 yes2win tpower-develop tpower-master tpower-new alfawin77-new vw-develop vw-master wintoto"

# 遍歷分支列表並進行push
for branch in $BRANCHES; do
    echo "開始處理分支: $branch"
    
    # 切換到目標分支
    git checkout $branch
    if [ $? -ne 0 ]; then
        echo "切換到 $branch 分支失敗"
        continue
    fi
    
    # 推送到目標分支
    git push origin $branch
    if [ $? -ne 0 ]; then
        echo "推送 $branch 失敗"
        continue
    fi

    echo "成功推送 $branch"
done

echo "所有分支處理完成"