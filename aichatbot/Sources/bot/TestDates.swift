// 題庫的相對日期：把 {{date:+1:M/D}} 這類標記換成以執行當天換算的實際日期，題目才不會過期。
// 規則和 eval/testdates.py 一致：
// {{date:±N:格式}}：基準日＝「今天 + 12 天」當天或之後的第一個星期一，值＝基準日 + 位移天數。
// {{today:±N:格式}}：今天 + N 天。
// {{weekday:next|this|near:1-7:格式}}：1＝週一，next＝下週X、this＝這週X、near＝只說週X（和核心的 readDate 相同規則）。
// 今天可用 BOT_TODAY=YYYY-MM-DD 指定（同時也是對話的「今天」，見 now()）。
import Foundation

private let calendar = Calendar(identifier: .gregorian)

/// 量測和對話用的「現在」：設了 BOT_TODAY 就用那天、時間照現在
func now() -> Date {
    let f = DateFormatter(); f.dateFormat = "yyyy-MM-dd"; f.locale = Locale(identifier: "en_US_POSIX")
    guard let s = ProcessInfo.processInfo.environment["BOT_TODAY"], let day = f.date(from: s) else { return Date() }
    let t = calendar.dateComponents([.hour, .minute, .second], from: Date())
    return calendar.date(bySettingHour: t.hour!, minute: t.minute!, second: t.second!, of: day)!
}

private func format(_ d: Date, _ f: String) -> String {
    let (y, mo, da) = (calendar.component(.year, from: d), calendar.component(.month, from: d), calendar.component(.day, from: d))
    switch f {
    case "M/D": return "\(mo)/\(da)"
    case "MM-DD": return String(format: "%02d-%02d", mo, da)
    case "M月D日": return "\(mo)月\(da)日"
    case "M月D號": return "\(mo)月\(da)號"
    default: return String(format: "%04d-%02d-%02d", y, mo, da)
    }
}

func resolveTestDates(_ text: String) -> String {
    let today = calendar.startOfDay(for: now())
    let start = calendar.date(byAdding: .day, value: 12, to: today)!
    let toMonday = (9 - calendar.component(.weekday, from: start)) % 7   // weekday：1=週日、2=週一
    let base = calendar.date(byAdding: .day, value: toMonday, to: start)!
    let weekdayIndex = (calendar.component(.weekday, from: today) + 5) % 7   // 週一＝0
    let formats = "(M/D|MM-DD|M月D日|M月D號|YYYY-MM-DD)"
    // 每種標記：pattern 的最後一組是格式，前面幾組交給 value 算出日期
    func replace(_ text: String, _ pattern: String, _ value: ([String]) -> Date) -> String {
        let regex = try! NSRegularExpression(pattern: pattern + formats + #"\}\}"#)
        var out = text
        for m in regex.matches(in: text, range: NSRange(text.startIndex..., in: text)).reversed() {
            let groups = (1..<m.numberOfRanges).map { (text as NSString).substring(with: m.range(at: $0)) }
            out = (out as NSString).replacingCharacters(in: m.range, with: format(value(groups), groups.last!))
        }
        return out
    }
    var out = replace(text, #"\{\{date:([+-]\d+):"#) { calendar.date(byAdding: .day, value: Int($0[0])!, to: base)! }
    out = replace(out, #"\{\{today:([+-]\d+):"#) { calendar.date(byAdding: .day, value: Int($0[0])!, to: today)! }
    return replace(out, #"\{\{weekday:(next|this|near):([1-7]):"#) {
        let target = Int($0[1])! - 1
        let offset = $0[0] == "next" ? 7 - weekdayIndex + target : $0[0] == "this" ? target - weekdayIndex : (target - weekdayIndex + 7) % 7
        return calendar.date(byAdding: .day, value: offset, to: today)!
    }
}

/// 讀題庫檔並換掉相對日期
func loadTestFile(_ path: String) -> Data {
    let text = try! String(contentsOf: URL(fileURLWithPath: path), encoding: .utf8)
    return Data(resolveTestDates(text).utf8)
}
