import Foundation

// MARK: - 相對日期：「明天」「下週二」由程式依今天換算（AI 自己算曾把週三說的「下週二」算成週六）
// AI 只抄日期的原文出處，這裡判斷是不是相對說法、換算成哪一天。純函式，不碰 AI

public enum DateReading: Equatable {
    case date(Date)          // 換算出的日期（當天 0 點）
    case unresolvedRelative  // 是相對說法但換算不了（下下週、週末、下個月…）：不要填，改問顧客
    case notRelative         // 不是相對說法（10/21、10月21日…）：照絕對日期處理
}

private func matches(_ text: String, _ pattern: String) -> Bool {
    text.range(of: pattern, options: .regularExpression) != nil
}

// 換算不了的相對說法（先比對，免得「下下週二」被當成「下週二」、「大後天」被當成「後天」）
private let unresolvedPatterns = [
    #"下\s*下\s*(週|周|禮拜|礼拜|星期|個|个)"#, #"大\s*(後|后)\s*天"#, #"再\s*(來|来)\s*週"#,
    #"(週|周)\s*末"#, "weekend", "土日",
    #"(下|這|这|本|上)\s*(個|个)?\s*月"#, #"next\s+month"#, "来月", "來月", "月底", "月末",
    #"(過|过)\s*(幾|几)\s*天"#, #"再\s*(過|过)"#, #"(天|日|週|周|星期|禮拜|礼拜|月)\s*(以)?(後|后)"#,
    #"\bin\s+\S+\s+(days?|weeks?)"#, #"(days?|weeks?)\s+(later|from\s+now)"#, #"week\s+after\s+next"#,
    "改天", "某天", "哪天", #"\blater\b"#, #"\bsomeday\b"#,
]

// 今天＋N 天的說法（「後天」要在「明天」之前比對；「day after tomorrow」要在「tomorrow」之前）
private let dayWords: [(pattern: String, offset: Int)] = [
    (#"day\s+after\s+tomorrow"#, 2), ("明後日", 2), ("あさって", 2), ("後天", 2), ("后天", 2),
    ("明天", 1), ("明日", 1), ("明晚", 1), ("tomorrow", 1), ("あした", 1),
    ("今天", 0), ("今日", 0), ("本日", 0), ("今晚", 0), ("today", 0), ("tonight", 0),
]

// 星期幾（週一＝0 … 週日＝6）
private let chineseWeekday = #"(週|周|星期|禮拜|礼拜)\s*([一二三四五六日天])"#
private let japaneseWeekday = #"([月火水木金土日])曜"#
private let englishWeekday = #"\b(monday|tuesday|wednesday|thursday|friday|saturday|sunday|mon|tues|tue|wed|thurs|thur|thu|fri|sat|sun)\b"#

private func weekday(in text: String) -> Int? {
    for pattern in [chineseWeekday, japaneseWeekday, englishWeekday] {
        guard let r = text.range(of: pattern, options: .regularExpression) else { continue }
        let m = String(text[r])
        for (i, names) in [["一", "月", "mon"], ["二", "火", "tue"], ["三", "水", "wed"], ["四", "木", "thu"],
                           ["五", "金", "fri"], ["六", "土", "sat"], ["日", "天", "sun"]].enumerated() {
            // 中文的「週」「星期」本身不含星期幾，只看最後一個字；日文看「曜」前面那個字；英文看開頭
            let key = pattern == englishWeekday ? String(m.prefix(3)) : pattern == japaneseWeekday ? String(m.prefix(1)) : String(m.suffix(1))
            if names.contains(key) { return i }
        }
    }
    return nil
}

private let nextWeek = [#"下\s*(個|个)?\s*(週|周|禮拜|礼拜|星期)"#, #"next\s+week"#, #"\bnext\s+(mon|tue|wed|thu|fri|sat|sun)"#, "来週", "來週", "らいしゅう"]
private let thisWeek = [#"(這|这|今)\s*(個|个)?\s*(週|周|禮拜|礼拜|星期)"#, #"本\s*(週|周)"#, #"this\s+week"#, #"\bthis\s+(mon|tue|wed|thu|fri|sat|sun)"#, "今週"]

/// 讀日期的原文出處：相對說法依 today 換算。「下週」是下一個星期一開始的那一週；只說「週X」是今天起最近的那一天（含今天）
public func readDate(_ quote: String, today: Date, calendar: Calendar = Calendar(identifier: .gregorian)) -> DateReading {
    let q = quote.lowercased()
    // 寫了月日（10/21（三）、10月21日星期三）就是絕對日期，星期幾只是附註
    if matches(q, #"\d+\s*[/\-.月]\s*\d+|\d+\s*(日|號|号)"#) { return .notRelative }
    if unresolvedPatterns.contains(where: { matches(q, $0) }) { return .unresolvedRelative }

    let start = calendar.startOfDay(for: today)
    func plus(_ days: Int) -> DateReading { .date(calendar.date(byAdding: .day, value: days, to: start)!) }
    if let w = dayWords.first(where: { matches(q, $0.pattern) }) { return plus(w.offset) }

    let todayIndex = (calendar.component(.weekday, from: start) + 5) % 7   // weekday：1=週日、2=週一
    let isNext = nextWeek.contains { matches(q, $0) }, isThis = thisWeek.contains { matches(q, $0) }
    guard let target = weekday(in: q) else { return isNext || isThis ? .unresolvedRelative : .notRelative }
    if isNext { return plus(7 - todayIndex + target) }   // 下一個星期一＋星期幾
    if isThis { return plus(target - todayIndex) }       // 本週一＋星期幾（可能已經過了，交給「已經過了」的檢查）
    return plus((target - todayIndex + 7) % 7)
}

/// 出處是不是原句裡某個換算不了的說法的一部分：AI 曾把「下下週三」只抄成「下週三」，被當成下週三。
/// 只看和出處重疊的說法，句子別處的「晚點再過去」不影響
public func partOfUnresolved(_ quote: String, in text: String) -> Bool {
    let t = text.lowercased()
    func ranges(_ pattern: String) -> [NSRange] {
        (try? NSRegularExpression(pattern: pattern))?.matches(in: t, range: NSRange(t.startIndex..., in: t)).map(\.range) ?? []
    }
    let quotes = ranges(NSRegularExpression.escapedPattern(for: quote.lowercased()))
    return unresolvedPatterns.flatMap(ranges).contains { u in quotes.contains { NSIntersectionRange($0, u).length > 0 } }
}

// 明確寫出的月日（10/19、10-19、10.19、10月19日）：前後不能接著數字（電話、年份），後面不能是人數或時間（4-6人、2.5點）
private let monthDayPattern = try! NSRegularExpression(
    pattern: #"(?<!\d)(\d{1,2})\s*(?:[/\-.]|月)\s*(\d{1,2})(?!\d)(?!\s*(個|个|位|人|點|点|時|时|分))"#)

/// 讀出訊息裡明確寫的月日（不靠 AI 填的整數；AI 曾把「10/19」的月日填成 0）。只留合法範圍，重複的去掉
public func explicitMonthDays(in text: String) -> [(month: Int, day: Int)] {
    var found: [(month: Int, day: Int)] = []
    for m in monthDayPattern.matches(in: text, range: NSRange(text.startIndex..., in: text)) {
        guard let mr = Range(m.range(at: 1), in: text), let dr = Range(m.range(at: 2), in: text),
              let month = Int(text[mr]), let day = Int(text[dr]),
              (1...12).contains(month), (1...31).contains(day), !found.contains(where: { $0 == (month, day) }) else { continue }
        found.append((month, day))
    }
    return found
}

/// 只有月、日的絕對日期補上年份：今年；這天在今天往前 7 天之前就算明年（7 天內維持今年，回「這天已經過了」）。
/// 不存在的日期（2/30）照樣組成字串，交給日期檢查回「看不懂」
public func absoluteDate(month: Int, day: Int, today: Date, calendar: Calendar = Calendar(identifier: .gregorian)) -> String {
    var year = calendar.component(.year, from: today)
    let comps = DateComponents(year: year, month: month, day: day)
    if comps.isValidDate(in: calendar), let d = calendar.date(from: comps),
       let limit = calendar.date(byAdding: .day, value: -7, to: calendar.startOfDay(for: today)), d < limit {
        year += 1
    }
    return String(format: "%04d-%02d-%02d", year, month, day)
}
