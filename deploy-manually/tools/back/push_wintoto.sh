#!/bin/bash

BRANCHES="${BRANCHES:-wintoto wtt-develop}"

for branch in $BRANCHES; do
    echo "開始處理分支: $branch"

    git checkout $branch
    if [ $? -ne 0 ]; then
        echo "切換到 $branch 分支失敗"
        continue
    fi

    git push origin $branch
    if [ $? -ne 0 ]; then
        echo "推送 $branch 失敗"
        continue
    fi

    echo "成功推送 $branch"
done

echo "wintoto 所有分支處理完成"
