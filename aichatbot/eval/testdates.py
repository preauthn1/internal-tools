"""題庫的相對日期：把 {{date:+1:M/D}} 這類標記換成以執行當天換算的實際日期，題目才不會過期。

{{date:±N:格式}}：基準日＝「今天 + 12 天」當天或之後的第一個星期一；標記的值＝基準日 + 位移天數。
2026-10-07 執行時基準日是 10/19，所以 {{date:+0:M/D}} → 10/19（週一，公休）、{{date:+1:M/D}} → 10/20。
位移是負數就是基準日之前（例如測「過去的日期」）。

{{today:±N:格式}}：今天 + N 天（測「明天」「後天」這類說法）。
{{weekday:next|this|near:1-7:格式}}：1＝週一 … 7＝週日，規則和核心的相對日期換算（Sources/BotCore/RelativeDate.swift）相同：
next＝下週X（下一個星期一開始的那一週）、this＝這週X（本週一開始，可能已經過了）、near＝只說週X（今天起最近的那一天，含今天）。
2026-10-07（週三）執行時 {{weekday:next:2:M/D}} → 10/13、{{weekday:this:1:M/D}} → 10/5、{{weekday:near:3:M/D}} → 10/7。

格式：M/D、MM-DD、M月D日、M月D號、YYYY-MM-DD
今天可用環境變數 BOT_TODAY=YYYY-MM-DD 指定（測試「一個月後」題目還有效）。
規則要和 Sources/bot/TestDates.swift 一致。

>>> d = datetime.date(2026, 10, 7)
>>> resolve("{{date:+1:M/D}} {{today:+1:MM-DD}} {{weekday:next:2:M/D}} {{weekday:this:1:M月D日}} {{weekday:near:3:YYYY-MM-DD}}", d)
'10/20 10-08 10/13 10月5日 2026-10-07'
>>> resolve("{{weekday:next:1:M/D}} {{weekday:near:1:M/D}}", datetime.date(2026, 10, 11))
'10/12 10/12'
"""
import datetime, os, re

FORMATS = r"(M/D|MM-DD|M月D日|M月D號|YYYY-MM-DD)"
PATTERN = re.compile(r"\{\{date:([+-]\d+):" + FORMATS + r"\}\}")
TODAY = re.compile(r"\{\{today:([+-]\d+):" + FORMATS + r"\}\}")
WEEKDAY = re.compile(r"\{\{weekday:(next|this|near):([1-7]):" + FORMATS + r"\}\}")


def today():
    t = os.environ.get("BOT_TODAY")
    return datetime.date.fromisoformat(t) if t else datetime.date.today()


def base_monday(day):
    d = day + datetime.timedelta(days=12)
    return d + datetime.timedelta(days=(7 - d.weekday()) % 7)


def weekday_date(day, kind, n):
    target, now = n - 1, day.weekday()   # 週一＝0
    offset = {"next": 7 - now + target, "this": target - now, "near": (target - now) % 7}[kind]
    return day + datetime.timedelta(days=offset)


def fmt(d, f):
    return {"M/D": f"{d.month}/{d.day}", "MM-DD": f"{d.month:02d}-{d.day:02d}",
            "M月D日": f"{d.month}月{d.day}日", "M月D號": f"{d.month}月{d.day}號",
            "YYYY-MM-DD": d.isoformat()}[f]


def resolve(text, day=None):
    day = day or today()
    base = base_monday(day)
    text = PATTERN.sub(lambda m: fmt(base + datetime.timedelta(days=int(m.group(1))), m.group(2)), text)
    text = TODAY.sub(lambda m: fmt(day + datetime.timedelta(days=int(m.group(1))), m.group(2)), text)
    return WEEKDAY.sub(lambda m: fmt(weekday_date(day, m.group(1), int(m.group(2))), m.group(3)), text)


if __name__ == "__main__":
    import doctest
    doctest.testmod()
