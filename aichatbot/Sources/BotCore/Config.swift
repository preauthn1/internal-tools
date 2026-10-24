// 新版客服機器人（重新設計）：明確的對話狀態 ＋ 一次 AI 分類 ＋ 程式依表分派
// 執行：swift build -c release && .build/release/bot（進入點和量測模式在 Sources/bot/main.swift）
// 指令：/reset 換一位新顧客  /as 切換顧客  /reload 重讀設定  /me 名字  /pending 待處理清單  /quit 離開
//       店員：/reply 編號 文字  /confirm 編號  /reject 編號 [原因]  /endhuman 結束轉人工
//
// 設計原則（來自 eval/ 的測試結果）：
//   - AI 只做選擇題：意圖分類、挑 FAQ、讀訂單數量、抽預約欄位、抽名字；不讓 AI 自由寫回覆
//   - 規則與計算交給程式：價格、折扣、預約檢查都依 store.json
//   - 對話狀態存在 Customer 裡（預約問到哪、有沒有在等店員、追問幾次），不靠比對上一句回覆；依「平台＋使用者 ID」分開保存
//   - 訊息進、動作出：Bot.receive 收一則統一格式的訊息，回傳回覆、開待處理事項或不回覆（見 Bot.swift）
//   - AI 呼叫一次只跑一個（同時跑品質會下降），一般訊息只呼叫 1～2 次
//   - 所有回覆都來自 replies.json、faq.txt 或程式計算的結果
//   - 設定檔在 BOT_CONFIG_DIR（沒設就是 eval/），啟動時驗證，有錯就印出哪個檔案、哪裡錯後結束；/reload 有錯時沿用舊設定
import Foundation
import NaturalLanguage

// 設定檔資料夾：BOT_CONFIG_DIR，沒設就用專案根目錄的 eval/（這個檔在 Sources/BotCore/ 底下，往上三層是專案根目錄；
// 只在開發機有效，部署到別台機器時一定要設 BOT_CONFIG_DIR）
public let env = ProcessInfo.processInfo.environment
public let dir = env["BOT_CONFIG_DIR"].map { URL(fileURLWithPath: $0, isDirectory: true) }
    ?? URL(fileURLWithPath: #filePath).deletingLastPathComponent().deletingLastPathComponent()
        .deletingLastPathComponent().appendingPathComponent("eval", isDirectory: true)

// MARK: - 設定檔

struct MenuItem: Decodable { let id: String; let price: Int; let aliases: [String]; let names: [String: String] }
public struct Store: Decodable {
    struct Discount: Decodable { let minTotal: Int; let rate: Double }
    let name: String
    let openHour, closeHour, maxPeople, cardMinimum: Int
    let closedWeekdays: [Int]   // Calendar 的星期：1=週日、2=週一…
    let pickupDiscount: Discount
    let menu: [MenuItem]
}
public struct FAQ { let id: String; public let questions: String; let answers: [String: String] }

// 設定檔的錯誤：哪個檔案、哪裡錯（JSON 的欄位路徑，或 txt/jsonl 的行號）
public struct ConfigError: Error, CustomStringConvertible {
    public let file: String
    public let detail: String
    public var description: String { "\(file)：\(detail)" }
}
public struct ConfigErrors: Error, CustomStringConvertible {
    public let errors: [ConfigError]
    public var description: String { errors.map { "  - \($0)" }.joined(separator: "\n") }
}

func readConfigFile(_ name: String, in d: URL) throws -> Data {
    let url = d.appendingPathComponent(name)
    guard FileManager.default.fileExists(atPath: url.path) else { throw ConfigError(file: name, detail: "找不到檔案（\(url.path)）") }
    do { return try Data(contentsOf: url) } catch { throw ConfigError(file: name, detail: "讀不到檔案：\(error.localizedDescription)") }
}
func readConfigText(_ name: String, in d: URL) throws -> String {
    guard let text = String(data: try readConfigFile(name, in: d), encoding: .utf8) else {
        throw ConfigError(file: name, detail: "不是 UTF-8 文字檔")
    }
    return text
}
// JSON 語法錯誤的說明（含錯誤位置）
func jsonSyntaxDetail(_ error: Error) -> String {
    let e = error as NSError
    return "JSON 格式錯誤：" + ((e.userInfo[NSDebugDescriptionErrorKey] as? String) ?? e.localizedDescription)
}

public func loadStore(from d: URL = dir) throws -> Store {
    let data = try readConfigFile("store.json", in: d)
    // 欄位路徑寫成 pickupDiscount.minTotal、menu[2].price
    func path(_ keys: [CodingKey]) -> String {
        keys.map { $0.intValue.map { "[\($0)]" } ?? ".\($0.stringValue)" }.joined().trimmingCharacters(in: CharacterSet(charactersIn: "."))
    }
    do {
        return try JSONDecoder().decode(Store.self, from: data)
    } catch let DecodingError.keyNotFound(key, c) {
        throw ConfigError(file: "store.json", detail: "缺少欄位 \(path(c.codingPath + [key]))")
    } catch let DecodingError.typeMismatch(type, c) {
        throw ConfigError(file: "store.json", detail: "欄位 \(path(c.codingPath)) 的型別不對（應為 \(type)）")
    } catch let DecodingError.valueNotFound(type, c) {
        throw ConfigError(file: "store.json", detail: "欄位 \(path(c.codingPath)) 沒有值（應為 \(type)）")
    } catch let DecodingError.dataCorrupted(c) {
        let detail = c.codingPath.isEmpty ? jsonSyntaxDetail(c.underlyingError ?? DecodingError.dataCorrupted(c))
                                          : "欄位 \(path(c.codingPath))：\(c.debugDescription)"
        throw ConfigError(file: "store.json", detail: detail)
    } catch {
        throw ConfigError(file: "store.json", detail: "\(error)")
    }
}

// 每個模板是 {"zh": …, "en": …, "ja": …} 或 {"*": …}：一定要有 zh 或 *，值是字串（en、ja 可缺，缺的時候用 zh）
// 值是字串的項目（例如「說明」）是給人看的註解，略過
public func loadReplies(from d: URL = dir) throws -> [String: [String: String]] {
    let data = try readConfigFile("replies.json", in: d)
    let obj: Any
    do { obj = try JSONSerialization.jsonObject(with: data) } catch {
        throw ConfigError(file: "replies.json", detail: jsonSyntaxDetail(error))
    }
    guard let top = obj as? [String: Any] else { throw ConfigError(file: "replies.json", detail: "最外層要是 JSON 物件 {…}") }
    var result: [String: [String: String]] = [:], problems: [String] = []
    for (key, value) in top.sorted(by: { $0.key < $1.key }) {
        if value is String { continue }
        guard let table = value as? [String: Any] else { problems.append("模板 \(key) 要是物件，例如 {\"zh\": \"…\"}"); continue }
        let bad = table.filter { !($0.value is String) }.keys.sorted()
        if !bad.isEmpty { problems.append("模板 \(key) 的 \(bad.joined(separator: "、")) 不是字串"); continue }
        guard table["zh"] != nil || table["*"] != nil else { problems.append("模板 \(key) 缺少 zh（或 *）"); continue }
        result[key] = table.compactMapValues { $0 as? String }
    }
    if !problems.isEmpty { throw ConfigError(file: "replies.json", detail: problems.joined(separator: "；")) }
    return result
}

// 一行一題：id | 問法（用 / 分隔）| 中文答案 [| 英文答案 [| 日文答案]]；# 開頭和空行略過
public func loadFAQ(from d: URL = dir) throws -> [FAQ] {
    let text = try readConfigText("faq.txt", in: d)
    var result: [FAQ] = [], problems: [String] = []
    for (i, raw) in text.components(separatedBy: "\n").enumerated() {
        let line = raw.trimmingCharacters(in: .whitespaces)
        if line.isEmpty || line.hasPrefix("#") { continue }
        let cols = line.components(separatedBy: " | ").map { $0.trimmingCharacters(in: .whitespaces) }
        guard cols.count >= 3 else { problems.append("第 \(i + 1) 行：欄位不足（要用「 | 」分隔 id | 問法 | 答案）"); continue }
        var answers = ["zh": cols[2]]
        if cols.count > 3 { answers["en"] = cols[3] }
        if cols.count > 4 { answers["ja"] = cols[4] }
        result.append(FAQ(id: cols[0], questions: cols[1], answers: answers))
    }
    if !problems.isEmpty { throw ConfigError(file: "faq.txt", detail: problems.joined(separator: "；")) }
    return result
}

// 程式用到的固定回覆模板（啟動驗證和 bot-tests 共用；新增 R("…")、replyText("…") 時要加進來，bot-tests 會掃原始碼檢查）
public let usedReplyKeys: Set<String> = Set([
    // Handle.swift：R(intent) 會用到的意圖，和其他寫死的模板
    "claim", "off_topic", "greeting", "thanks", "goodbye", "small_talk", "complaint", "unclear",
    "unsupported", "booking_cancel", "nudge", "waiting", "handoff", "you_are", "bot_identity", "ask_name_soft", "self_intro",
    // StoreQuestion.swift
    "calc_total", "calc_pickup_discount", "calc_pickup_none", "calc_delivery", "calc_card_ok", "calc_card_no", "fallback",
    // Booking.swift：預約檢查的問題代碼（validate 回傳）、送出預約
    "err_date", "err_past", "err_monday", "err_time", "err_hours", "err_people", "err_too_many", "err_phone",
    "people_unit", "booking_done",
    // Bot.swift：店員確認或拒絕預約
    "booking_confirmed", "booking_rejected", "booking_rejected_reason",
    // 店員群組（#11）：店員回覆的前綴、顧客傳非文字訊息
    "staff_reply_prefix", "text_only",
]).union(bookingOrder.map { "ask_" + $0 })

public func missingReplyTemplates(_ replies: [String: [String: String]]) -> [String] {
    usedReplyKeys.filter { replies[$0] == nil }.sorted()
}

// 全部設定：啟動和 /reload 時一次讀完、驗證完才換上
public struct ConfigBundle {
    public let store: Store
    public let replies: [String: [String: String]]
    public let faqs: [FAQ]
    public let examples: [Example]
}

// 四個檔都讀一遍，收集所有錯誤（不在第一個錯就停），再檢查程式用到的回覆模板都在
public func loadConfig(from d: URL = dir) throws -> ConfigBundle {
    var errors: [ConfigError] = []
    func attempt<T>(_ load: () throws -> T) -> T? {
        do { return try load() } catch let e as ConfigError { errors.append(e) } catch { errors.append(ConfigError(file: "?", detail: "\(error)")) }
        return nil
    }
    let store = attempt { try loadStore(from: d) }
    let replies = attempt { try loadReplies(from: d) }
    let faqs = attempt { try loadFAQ(from: d) }
    let examples = attempt { try loadExamples(from: d) }
    if let replies, case let missing = missingReplyTemplates(replies), !missing.isEmpty {
        errors.append(ConfigError(file: "replies.json", detail: "缺少程式用到的回覆模板：\(missing.joined(separator: "、"))"))
    }
    guard errors.isEmpty, let store, let replies, let faqs, let examples else { throw ConfigErrors(errors: errors) }
    return ConfigBundle(store: store, replies: replies, faqs: faqs, examples: examples)
}

public func applyConfig(_ c: ConfigBundle) {
    store = c.store; replies = c.replies; faqs = c.faqs
    examples = c.examples; exampleVectors = [:]
}

// 啟動時用：設定有錯就印出每個錯誤後結束（exit code 1），不當掉
public func loadConfigOrExit(from d: URL = dir) -> ConfigBundle {
    do { return try loadConfig(from: d) } catch {
        let msg = "❌ 設定檔有錯（\(d.path)；可用 BOT_CONFIG_DIR 指定設定檔資料夾）：\n\(error)\n"
        FileHandle.standardError.write(msg.data(using: .utf8)!)
        exit(1)
    }
}

// 第一次用到設定（store、replies…）時讀取並驗證
public let initialConfig = loadConfigOrExit()
public var store = initialConfig.store
public var replies = initialConfig.replies
public var faqs = initialConfig.faqs

// MARK: - 回覆文字

let supportedLanguages: Set<String> = ["zh", "en", "ja"]

public func replyText(_ key: String, _ vars: [String: String] = [:], lang: String) -> String {
    guard let table = replies[key] else { return key }
    var s = (table[lang] ?? table["*"] ?? table["zh"] ?? "").components(separatedBy: "|").randomElement()!
    let all = vars.merging(["open": String(store.openHour), "close": String(store.closeHour), "max": String(store.maxPeople),
                            "min": String(store.pickupDiscount.minTotal)]) { a, _ in a }
    for (k, v) in all { s = s.replacingOccurrences(of: "{\(k)}", with: v) }
    return s
}

extension Turn {
    // 用這位顧客的語言
    func R(_ key: String, _ vars: [String: String] = [:]) -> String { replyText(key, vars, lang: customer.lang) }
}

// MARK: - 語言（程式判斷，不用 AI）

public func detectLanguage(_ text: String) -> String? {
    // 先看文字系統：假名→日文、韓文字母→韓文、漢字→中文（「可以用LINE Pay嗎」不會被判成英文）
    let scalars = text.unicodeScalars
    if scalars.contains(where: { (0x3040...0x30FF).contains($0.value) }) { return "ja" }
    if scalars.contains(where: { (0xAC00...0xD7AF).contains($0.value) || (0x1100...0x11FF).contains($0.value) }) { return "ko" }
    if scalars.contains(where: { $0.properties.isIdeographic }) { return "zh" }
    guard text.filter(\.isLetter).count >= 2 else { return nil }   // 表情、標點：沿用上一句
    let r = NLLanguageRecognizer()
    r.languageConstraints = [.english, .vietnamese, .spanish, .french, .german, .indonesian, .malay, .portuguese, .italian, .russian, .thai]
    r.processString(text)
    let hypotheses = r.languageHypotheses(withMaximum: 3).sorted { $0.value > $1.value }
    guard let best = hypotheses.first else { return nil }
    let code = { (l: NLLanguage) in String(l.rawValue.prefix(2)) }
    // 短句容易誤判（「Hi, I am Kevin」曾被判成 80% 印尼文）：候選有英文就用英文
    if code(best.key) != "en", text.count < 40, hypotheses.contains(where: { code($0.key) == "en" && $0.value >= 0.1 }) { return "en" }
    if best.value < 0.5 { return text.unicodeScalars.filter(\.properties.isAlphabetic).allSatisfy(\.isASCII) ? "en" : nil }
    return code(best.key)
}

// 很短的回答（名字、電話、人數）：預約問欄位時不依它切換語言（#17：中文預約問到姓名回「Bill」，之後都變成英文）
// 用空白分詞的文字（英文…）不超過 2 個詞；中日韓（不用空白）去掉標點後不超過 6 個字；只有數字、符號也算短
public func isShortAnswer(_ text: String) -> Bool {
    let cjk = text.unicodeScalars.contains {
        (0x3040...0x30FF).contains($0.value) || (0xAC00...0xD7AF).contains($0.value) || (0x1100...0x11FF).contains($0.value) || $0.properties.isIdeographic
    }
    if cjk { return text.filter { $0.isLetter || $0.isNumber }.count <= 6 }
    let words = text.split { !$0.isLetter && !$0.isNumber }.filter { $0.contains(where: \.isLetter) }
    return words.count <= 2
}

// 除錯：BOT_DEBUG=1 時印出每個 AI 判斷的結果（含 AI 失敗時用的預設值）
func debug(_ s: String) { if env["BOT_DEBUG"] == "1" { print("   🔍 \(s)") } }
