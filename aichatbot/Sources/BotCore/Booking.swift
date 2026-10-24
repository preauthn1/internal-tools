import Foundation

// MARK: - 名字（由分類那次呼叫一起抽出，這裡只做驗證）

func verifiedName(_ raw: String, in text: String) -> String? {
    let name = raw.trimmingCharacters(in: .whitespaces)
    var core = name
    for title in ["先生", "小姐", "女士", "太太", "様", "さん", "Mr. ", "Ms. "] { core = core.replacingOccurrences(of: title, with: "") }
    // 名字必須真的出現在訊息裡（防止 AI 編名字）
    guard !core.isEmpty, text.localizedCaseInsensitiveContains(core.trimmingCharacters(in: .whitespaces)) else { return nil }
    return name.count == 1 && name.unicodeScalars.allSatisfy(\.properties.isIdeographic) ? "\(name)先生／小姐" : name
}

// 問姓名時的短答整句就是名字（去掉頭尾標點和「我是」「です」這類前後綴）；有數字或是完整句子就不算
public func nameFromShortAnswer(_ text: String) -> String? {
    var name = text.trimmingCharacters(in: .whitespacesAndNewlines.union(.punctuationCharacters).union(.symbols))
    for prefix in ["我是", "我叫", "叫我", "My name is ", "I'm ", "I am "] where name.lowercased().hasPrefix(prefix.lowercased()) {
        name = String(name.dropFirst(prefix.count)); break
    }
    for suffix in ["と申します", "と言います", "です"] where name.hasSuffix(suffix) { name = String(name.dropLast(suffix.count)); break }
    name = name.trimmingCharacters(in: .whitespaces)
    guard !name.isEmpty, isShortAnswer(text), name.contains(where: \.isLetter), !name.contains(where: \.isNumber) else { return nil }
    return name
}

// MARK: - 預約（程式控制流程，AI 只抽欄位）

let bookingOrder = ["date", "time", "people", "name", "phone"]
let fieldNames = ["date": "日期", "time": "時間", "people": "人數", "name": "預約人姓名", "phone": "電話"]
let dayFormatter: DateFormatter = {
    let f = DateFormatter(); f.dateFormat = "yyyy-MM-dd"; f.locale = Locale(identifier: "en_US_POSIX"); return f
}()

func numbers(in text: String) -> Set<Int> {
    let regex = try! NSRegularExpression(pattern: #"\d+"#)
    return Set(regex.matches(in: text, range: NSRange(text.startIndex..., in: text)).compactMap {
        Range($0.range, in: text).flatMap { Int(text[$0]) }
    })
}

// 人數必須在訊息的數字裡找得到依據（日期、時間、電話的數字不能拿來當人數）
func numbersSupport(_ value: String, in text: String) -> Bool {
    let stripped = text.replacingOccurrences(of: #"\d+\s*[/:\-.月]\s*\d+|\d{6,}"#, with: " ", options: .regularExpression)
    return Int(value.filter(\.isNumber)).map(numbers(in: stripped).contains) ?? false
}

@MainActor func validate(_ key: String, _ value: String, today: Date) -> (String?, String) {
    switch key {
    case "date":
        guard let d = dayFormatter.date(from: value) else { return ("err_date", value) }
        if d < Calendar.current.startOfDay(for: today) { return ("err_past", value) }
        if store.closedWeekdays.contains(Calendar.current.component(.weekday, from: d)) { return ("err_monday", value) }
        return (nil, value)
    case "time":
        let p = value.split(separator: ":").compactMap { Int($0) }
        guard var h = p.first else { return ("err_time", value) }
        if h < store.openHour, h + 12 < store.closeHour { h += 12 }   // 「兩點」多半是下午
        if h < store.openHour || h >= store.closeHour { return ("err_hours", value) }
        return (nil, String(format: "%02d:%02d", h, p.count > 1 ? p[1] : 0))
    case "people":
        guard let n = Int(value.filter(\.isNumber)), n > 0 else { return ("err_people", value) }
        return n > store.maxPeople ? ("err_too_many", value) : (nil, String(n))
    case "phone":
        let digits = value.filter(\.isNumber)
        return (9...10).contains(digits.count) ? (nil, digits) : ("err_phone", value)
    default:
        return (nil, value)
    }
}

// 日期、時間也用整數：月、日、時、分由程式組成，並和原文出處的數字比對；相對日期（明天、下週二）由程式換算
let bookingFields: [Field] = [
    .text("date_quote", "訊息中提到日期的原文，逐字照抄；未提及則為空字串"),
    .integer("date_month", "預約日期的月；有寫月日就照填；只說相對日期（明天、下週二…）或未提及則為 0"),
    .integer("date_day", "預約日期的日；有寫月日就照填；只說相對日期（明天、下週二…）或未提及則為 0"),
    .text("time_quote", "訊息中提到時間的原文，逐字照抄；未提及則為空字串"),
    .integer("time_hour", "時，24 小時制（下午兩點=14）；未提及則為 0"),
    .integer("time_minute", "分（兩點半=30）；未提及則為 0"),
    // 人數用整數型別：用文字欄位時 AI 偶爾填成「{10}」「)」（guided generation 會強制只能輸出數字）
    .integer("people", "人數；未提及則為 0"),
    .text("people_quote", "訊息中提到人數的原文，逐字照抄；未提及則為空字串"),
    .text("name", "預約人姓名；未提及則為空字串"),
    .text("phone", "電話號碼；未提及則為空字串"),
]

extension Turn {
    // probe=true：顧客在預約途中說了別的話，沒抽到新欄位就回傳 nil，交回一般流程
    // intent：這句的分類，決定問姓名時能不能把短答整句當名字
    func bookingStep(_ text: String, probe: Bool, intent: String) async -> (String, String)? {
        let context = customer.askingField.map { "客服剛剛問了顧客的\(fieldNames[$0]!)，這句話很可能就是在回答它。" } ?? ""
        let c = await ask("Booking", text, instructions: "從顧客訊息抽取預約資訊。今天是 \(dayFormatter.string(from: now))。\(context)只說相對日期（明天、下週二…）時照抄原文，月、日填 0，程式會換算；有寫月日就照填。顧客可能使用任何語言。", fields: bookingFields)
        func quote(_ k: String) -> String { (c?.string(k + "_quote") ?? "").trimmingCharacters(in: .whitespaces) }
        // 出處要逐字出現在訊息裡，而且不能和其他欄位的出處重疊（日期和時間常寫在一起「10/19 下午兩點」，彼此重疊不算）
        func quoted(_ k: String) -> Bool {
            let q = quote(k)
            let others = (k == "people" ? ["date", "time"] : ["people"]).map(quote).filter { !$0.isEmpty }
            return !q.isEmpty && text.contains(q) && !others.contains { $0.contains(q) || q.contains($0) }
        }
        // 月日時分要和原文一致（每組至少一個數字出現）：有出處看出處，出處沒有阿拉伯數字（「下午兩點」）不檢查；
        // 沒有出處就看整句的數字，剛問的欄位整句沒數字也採用，其他欄位一定要找到（防止 AI 自己編）
        func supported(_ k: String, _ needs: [Set<Int>]) -> Bool {
            var source = quoted(k) ? quote(k) : text
            // 時間不看日期、電話的數字（「10/19 下午兩點」的 10、19 不是時間）
            if k == "time" { source = source.replacingOccurrences(of: #"\d+\s*[/\-.月]\s*\d+\s*(日|號|号)?|\d{6,}"#, with: " ", options: .regularExpression) }
            let n = numbers(in: source)
            if n.isEmpty { return quoted(k) || k == customer.askingField }
            return needs.allSatisfy { !$0.isDisjoint(with: n) }
        }
        func date() -> String {
            // 出處不在原句裡（AI 改寫或沒抄）就看整句
            let q = quote("date"), inText = !q.isEmpty && text.contains(q)
            let written = explicitMonthDays(in: text)
            // 原句是換算不了的說法（下下週三）又沒寫月日：AI 漏抄成「下週三」也不填
            if written.isEmpty, inText, partOfUnresolved(q, in: text) { return "" }
            // 相對說法由程式依今天換算，AI 填的月日不用（AI 曾把週三說的「下週二」填成週六）；換算不了就不填，改問顧客
            switch readDate(inText ? q : text, today: now) {
            case .date(let d): return dayFormatter.string(from: d)
            case .unresolvedRelative: return ""
            case .notRelative: break
            }
            let m = c?.int("date_month") ?? 0, d = c?.int("date_day") ?? 0
            // 寫了月日就由程式讀出來，AI 填的月日不用（AI 曾把「10/19」的月日填成 0，明確的日期被丟掉）；
            // 寫了好幾個日期時不自己猜，AI 填的要是其中之一
            let fromQuote = inText ? explicitMonthDays(in: q) : []
            let candidates = fromQuote.isEmpty ? written : fromQuote
            if candidates.count == 1 { return absoluteDate(month: candidates[0].month, day: candidates[0].day, today: now) }
            if !candidates.isEmpty { return candidates.contains { $0 == (m, d) } ? absoluteDate(month: m, day: d, today: now) : "" }
            // 出處沒寫日期（AI 抄錯段落），整句是相對說法就用整句換算
            if inText {
                switch readDate(text, today: now) {
                case .date(let day): return dayFormatter.string(from: day)
                case .unresolvedRelative: return ""
                case .notRelative: break
                }
            }
            // 讀不到月日（「十月十九日」「19號」）才用 AI 填的整數，並和出處的數字比對
            guard (1...12).contains(m), (1...31).contains(d), supported("date", [[m], [d]]) else { return "" }
            return absoluteDate(month: m, day: d, today: now)
        }
        func time() -> String {
            let h = c?.int("time_hour") ?? 0, m = c?.int("time_minute") ?? 0
            guard (1...23).contains(h), (0...59).contains(m) else { return "" }
            // 「2點半」的分沒有數字
            let half = m == 30 && (quoted("time") ? quote("time") : text).contains("半")
            guard supported("time", [[h, h % 12]] + (m == 0 || half ? [] : [[m]])) else { return "" }
            return String(format: "%d:%02d", h, m)
        }
        func field(_ k: String) -> String {
            switch k {
            case "date": return date()
            case "time": return time()
            case "people": return (c?.int(k)).map { $0 > 0 ? String($0) : "" } ?? ""
            case "name":
                let raw = (c?.string(k) ?? "").trimmingCharacters(in: .whitespaces)
                // 剛問姓名、AI 抽的名字不能用（空的或不在原文裡）時，短答整句就是名字
                // （AI 曾把「たなか」改寫掉，沒通過原文檢查，預約一直卡在問姓名）
                let usable = !raw.isEmpty && text.localizedCaseInsensitiveContains(raw) && raw.contains(where: \.isLetter)
                // 只在分類成預約或自我介紹時才這樣做，而且整句都是漢字時不用（AI 會照抄中文名字）：
                // 「不知道」「好」曾被收成名字（#17 第二版驗收）
                if !usable, k == customer.askingField, ["booking", "self_intro"].contains(intent),
                   let name = nameFromShortAnswer(text), !name.unicodeScalars.allSatisfy(\.properties.isIdeographic) { return name }
                return raw
            default: return (c?.string(k) ?? "").trimmingCharacters(in: .whitespaces)
            }
        }

        var problems: [String] = [], found = false
        for key in bookingOrder {
            let raw = field(key)
            if raw.isEmpty || raw == "0" || customer.booking[key] == raw { continue }
            switch key {
            case "people":
                let quote = quote(key)
                // 值要和原文出處一致：出處抄對「25 個人」、值卻填 10（取自日期 10/21）曾被收下。
                // 出處有阿拉伯數字時，人數必須是其中之一；不一致就當作沒抓到，改成再問一次（「四個人」這種沒有數字的不檢查）
                if text.contains(quote), case let n = numbers(in: quote), !n.isEmpty, !n.contains(Int(raw) ?? -1) { continue }
                // 剛問的欄位直接採用；其他欄位要有原文出處或數字依據（防止 AI 自己編）
                if key != customer.askingField && !quoted(key) && !numbersSupport(raw, in: text) { continue }
            case "phone":
                let digits = raw.filter(\.isNumber)
                if digits.isEmpty || !text.filter(\.isNumber).contains(digits) { continue }
            case "name":
                // 姓名要出現在原文裡，而且要有文字（電話號碼曾被收成姓名）
                if !text.localizedCaseInsensitiveContains(raw) || !raw.contains(where: \.isLetter) { continue }
            default: break
            }
            found = true
            let (problem, cleaned) = validate(key, raw, today: now)
            if let problem { problems.append(problem) } else { customer.booking[key] = cleaned }
        }
        if probe && !found { return nil }
        return bookingNext(problems)
    }

    // 依目前收到的欄位決定下一步：回報問題、問下一個欄位，或收齊後送出預約申請
    func bookingNext(_ problems: [String] = []) -> (String, String) {
        if customer.booking["name"] == nil, let name = customer.name { customer.booking["name"] = name }
        if let p = problems.first {
            customer.askingField = customer.askingField ?? bookingOrder.first { customer.booking[$0] == nil }
            return (p, R(p))
        }
        if let next = bookingOrder.first(where: { customer.booking[$0] == nil }) {
            customer.askingField = next
            return ("ask_" + next, R("ask_" + next))
        }
        let b = customer.booking
        let summary = "\(b["date"]!) \(b["time"]!), \(b["people"]!)\(R("people_unit")), \(b["name"]!), \(b["phone"]!)"
        // 預約申請留著等店員確認或拒絕（編號和待處理事項相同）
        let t = openTicket(.booking, "📅 預約申請：\(summary)", waiting: false)
        tickets.add(BookingRequest(id: t.id, customer: self.key, chatID: chatID, summary: summary, fields: b))
        customer.booking = [:]; customer.askingField = nil
        return ("booking_done", R("booking_done", ["summary": summary]))
    }
}
