#!/usr/bin/env python3
"""一次性工具：把題庫寫死的 2026 年 10 月日期換成相對日期標記 {{date:±N:格式}}（規則見 testdates.py）。

寫死日期時的「今天」是 2026-10-07，對應的基準日是 2026-10-19，所以 10/20 → {{date:+1:M/D}}。
驗證題組不能打開來看：這支程式只輸出每個檔案替換了幾處、還剩幾個像日期的字串，不印任何內容。

執行：python3 eval/relativize_dates.py eval/flow_tests.json eval/holdout/flow_holdout.json …
"""
import datetime, re, sys

BASE = datetime.date(2026, 10, 19)
RULES = [   # 先比對完整日期，避免 2026-10-20 被當成 10-20
    (re.compile(r"2026-10-(\d{2})(?!\d)"), "YYYY-MM-DD"),
    (re.compile(r"(?<![\d/])10/(\d{1,2})(?![\d/])"), "M/D"),
    (re.compile(r"(?<![\d-])10-(\d{2})(?![\d-])"), "MM-DD"),
    (re.compile(r"(?<!\d)10月(\d{1,2})日"), "M月D日"),
    (re.compile(r"(?<!\d)10月(\d{1,2})號"), "M月D號"),
]
LEFTOVER = re.compile(r"(?<![\d/{])\d{1,2}/\d{1,2}(?![\d/])|\d{1,2}月\d{1,2}[日號]")


def convert(text):
    n = 0
    for pattern, fmt in RULES:
        def sub(m):
            nonlocal n
            day = int(m.group(1))
            if not 1 <= day <= 31: return m.group(0)
            n += 1
            return f"{{{{date:{(datetime.date(2026, 10, day) - BASE).days:+d}:{fmt}}}}}"
        text = pattern.sub(sub, text)
    return text, n


for path in sys.argv[1:]:
    text = open(path, encoding="utf-8").read()
    new, n = convert(text)
    if n: open(path, "w", encoding="utf-8").write(new)
    print(f"{path}：替換 {n} 處，還剩 {len(LEFTOVER.findall(new))} 個像日期的字串沒換")
