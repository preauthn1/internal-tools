// 客服機器人的執行檔：終端機對話，以及 BOT_INTENT_EVAL、BOT_TESTSET 兩種量測模式（核心在 Sources/BotCore）
// 編譯與執行：swift build -c release && .build/release/bot
import Foundation
import FoundationModels
import BotCore
import Storage

// 啟動時先讀設定檔並驗證，有錯就印出來後結束（函式庫的全域變數第一次用到才初始化，這裡照原本的順序先讀）
_ = (store, replies, faqs, examples)

// MARK: - 紀錄與流程測試追蹤

let logDir = env["CHAT_LOG_DIR"].map { URL(fileURLWithPath: $0) } ?? dir.appendingPathComponent("logs")
try? FileManager.default.createDirectory(at: logDir, withIntermediateDirectories: true)
let logURL = logDir.appendingPathComponent("bot_log_\(Int(Date().timeIntervalSince1970)).md")
var log = "# bot.swift 對話紀錄（\(Date())）\n\n"
func save() { try? log.write(to: logURL, atomically: true, encoding: .utf8) }

let traceURL = env["FLOW_TRACE"].map { URL(fileURLWithPath: $0) }
func trace(_ obj: [String: Any]) {
    guard let traceURL, let data = try? JSONSerialization.data(withJSONObject: obj),
          let line = String(data: data, encoding: .utf8) else { return }
    if let h = try? FileHandle(forWritingTo: traceURL) {
        h.seekToEndOfFile(); h.write((line + "\n").data(using: .utf8)!); try? h.close()
    } else {
        try? (line + "\n").write(to: traceURL, atomically: true, encoding: .utf8)
    }
}

// 回覆顧客的那個動作：(回覆種類, 文字)；沒有回覆（已轉人工）就是 ("silent", "")
func replyOf(_ actions: [BotAction]) -> (String, String) {
    for case let .reply(key, text, _) in actions.reversed() { return (key, text) }
    return ("silent", "")
}

// 所有 AI 呼叫經過同一個排隊（AIQueue.shared，跨程序用檔案鎖），量測結果附上等了多久（loop/score.py 會讀這一行）
func printLockWait() {
    print(String(format: "等鎖總時間：%.1f 秒（AI 呼叫 %d 次）", AIQueue.shared.lockWaitSeconds, AIQueue.shared.callCount))
}

// 所有模式共用同一個模型：BOT_AI_CACHE=1 時先查 AI 回答快取（見 AICache.swift），命中就不排隊
let model = CachedModel.fromEnvironment(base: QueuedModel(base: AppleModel()))
// 有開快取才印；對話模式的 stdout 被 eval/run_flows.py 丟掉，改印到 stderr
func printCacheStats(toStderr: Bool = false) {
    guard let cache = model as? CachedModel else { return }
    if toStderr { FileHandle.standardError.write((cache.summary + "\n").data(using: .utf8)!) } else { print(cache.summary) }
}

// MARK: - 主迴圈

guard case .available = SystemLanguageModel.default.availability else {
    print("❌ 本機模型無法使用：\(SystemLanguageModel.default.availability)"); exit(1)
}

// MARK: - 只測分類（BOT_INTENT_EVAL=題組路徑；ROUNDS=輪數，預設 1：AI 用 greedy，多跑幾輪結果都一樣）
// 直接呼叫正式流程同一個 classify()（每題給一位新顧客，依題目設定狀態），避免測試和實際行為不一致

if let evalPath = env["BOT_INTENT_EVAL"] {
    struct Item: Decodable { let text: String; let accept: [String]; let state: String? }
    struct ItemSet: Decodable { let items: [Item] }
    // 題組已依 Go 版的 11 類重新標記（#26）：Swift 版的打招呼、道謝、道別、閒聊、問店家合併成 general 再比對，
    // 舊標記的題組照樣能比（10 類換算的分數）
    let merged: Set<String> = ["greeting", "thanks", "goodbye", "small_talk", "store_question"]
    func label(_ s: String) -> String { merged.contains(s) ? "general" : s }
    let items = try! JSONDecoder().decode(ItemSet.self, from: loadTestFile(evalPath)).items
        .map { Item(text: $0.text, accept: $0.accept.map(label), state: $0.state) }
    let rounds = Int(env["ROUNDS"] ?? "") ?? 1
    var results: [[String]] = Array(repeating: [], count: items.count)   // 每題每輪的分類結果
    let bot = Bot(model: model)
    for _ in 1...rounds {
        for (i, item) in items.enumerated() {
            var customer = Customer()
            if item.state == "ticket" { customer.openTicket = true }
            if let s = item.state, s.hasPrefix("booking:") { customer.askingField = String(s.dropFirst(8)) }
            results[i].append(await bot.classify(item.text, as: customer).intent)
        }
    }
    results = results.map { $0.map(label) }
    // 依「最理想的類別」分組統計
    var byLabel: [String: (ok: Int, total: Int)] = [:]
    for (i, item) in items.enumerated() {
        let ok = results[i].filter(item.accept.contains).count
        byLabel[item.accept[0], default: (0, 0)].ok += ok
        byLabel[item.accept[0], default: (0, 0)].total += rounds
    }
    print("========== 各類別（\(rounds) 輪）==========")
    for (label, r) in byLabel.sorted(by: { $0.key < $1.key }) {
        print("\(label.padding(toLength: 16, withPad: " ", startingAt: 0)) \(r.ok)/\(r.total)")
    }
    print("\n========== 答錯的（各輪結果）==========")
    for (i, item) in items.enumerated() where results[i].contains(where: { !item.accept.contains($0) }) {
        let st = item.state.map { "［\($0)］" } ?? ""
        print("\(st)\(item.text)　應為 \(item.accept.joined(separator: "/"))　→ \(results[i].joined(separator: "、"))")
    }
    let ok = items.indices.map { i in results[i].filter(items[i].accept.contains).count }.reduce(0, +)
    let unstable = items.indices.filter { Set(results[$0]).count > 1 }.count
    print(String(format: "\n總計：%d/%d（%.0f%%）　不穩定：%d 題", ok, items.count * rounds, Double(ok) / Double(items.count * rounds) * 100, unstable))
    printLockWait(); printCacheStats()
    exit(0)
}

// MARK: - 單題題庫（BOT_TESTSET=題庫路徑，可用逗號接多個；ROUNDS=輪數，預設 1：AI 用 greedy，多跑幾輪結果都一樣）
// 每題換新顧客（新的 Bot）、走完整的 receive()，評分規則和 run_testset.swift 相同：
//   答錯 = 回答了但不該回答，或回答內容不對（最嚴重）；該答沒答 = 應該回答卻轉店員／請顧客說清楚
// 題目和 faq.txt 的問法完全相同時另外標出（測不出挑 FAQ 的能力），統計同時列「排除重複後」

if let setPaths = env["BOT_TESTSET"] {
    struct Case: Decodable { let id, q, expect: String; let cat: String?; let all, any, none: [String]?; let lang: String?; let holdout: Bool? }
    struct CaseSet: Decodable { let cases: [Case] }
    let showHoldout = env["HOLDOUT"] == "1"
    let cases = setPaths.split(separator: ",").flatMap { p in
        try! JSONDecoder().decode(CaseSet.self, from: loadTestFile(String(p))).cases
    }.filter { ($0.holdout ?? false) == showHoldout }
     .filter { c in env["IDS"].map { $0.split(separator: ",").contains(Substring(c.id)) } ?? true }   // 快速題組：IDS=id1,id2…
    let rounds = Int(env["ROUNDS"] ?? "") ?? 1
    let faqPhrases = Set(faqs.flatMap { $0.questions.split(separator: "/").map { normalized(String($0)) } })
    func leaked(_ c: Case) -> Bool { faqPhrases.contains(normalized(c.q)) }

    func judge(_ c: Case, key: String, reply: String) -> String {   // "ok" / "wrong" / "missed"
        let answered = key == "answer"
        func contentOK() -> Bool {
            if let all = c.all, all.contains(where: { !reply.contains($0) }) { return false }
            if let any = c.any, !any.isEmpty, !any.contains(where: reply.contains) { return false }
            if let none = c.none, none.contains(where: reply.contains) { return false }
            if let lang = c.lang, detectLanguage(reply) != lang { return false }
            return true
        }
        switch c.expect {
        case "answer": return answered ? (contentOK() ? "ok" : "wrong") : "missed"
        case "no_answer": return answered ? "wrong" : "ok"
        default: return answered && !contentOK() ? "wrong" : "ok"
        }
    }

    var rows: [[String: Any]] = []
    var outcome: [String: [String]] = [:], picks: [String: [String]] = [:], keys: [String: [String]] = [:]
    for round in 1...rounds {
        for c in cases {
            lastPick = nil
            let (key, reply) = replyOf(await Bot(model: model).receive(IncomingMessage(platform: "eval", userID: c.id, text: c.q, time: now())))
            let r = judge(c, key: key, reply: reply)
            outcome[c.id, default: []].append(r)
            picks[c.id, default: []].append(lastPick ?? "-")
            keys[c.id, default: []].append(key)
            rows.append(["id": c.id, "q": c.q, "round": round, "key": key, "pick": lastPick ?? "-", "reply": reply, "result": r, "leaked": leaked(c)])
        }
        print("第 \(round) 輪完成"); fflush(stdout)
    }

    func summary(_ subset: [Case]) -> String {
        let all = subset.flatMap { outcome[$0.id]! }
        let ok = all.filter { $0 == "ok" }.count, wrong = all.filter { $0 == "wrong" }.count, missed = all.filter { $0 == "missed" }.count
        return String(format: "%d/%d（%.0f%%）　答錯 %d　該答沒答 %d", ok, all.count, Double(ok) / Double(max(all.count, 1)) * 100, wrong, missed)
    }
    print("\n========== 各類別（\(rounds) 輪）==========")
    for cat in Array(Set(cases.compactMap(\.cat))).sorted() {
        print("\(cat.padding(toLength: 10, withPad: "　", startingAt: 0)) \(summary(cases.filter { $0.cat == cat }))")
    }
    print("\n========== 沒通過的（各輪：結果／挑 FAQ／回覆種類）==========")
    for c in cases where outcome[c.id]!.contains(where: { $0 != "ok" }) {
        let detail = zip(zip(outcome[c.id]!, picks[c.id]!), keys[c.id]!).map { "\($0.0 == "ok" ? "✓" : $0.0 == "wrong" ? "✗答錯" : "△沒答")/\($0.1)/\($1)" }
        print("\(leaked(c) ? "［重複］" : "")[\(c.id)] \(c.q)　應為 \(c.expect)　→ \(detail.joined(separator: "、"))")
    }
    let clean = cases.filter { !leaked($0) }
    print("\n全部：\(summary(cases))")
    print("排除與 FAQ 問法重複的 \(cases.count - clean.count) 題：\(summary(clean))")
    let out = dir.appendingPathComponent("results/testset-bot-\(Int(Date().timeIntervalSince1970)).json")
    if let data = try? JSONSerialization.data(withJSONObject: rows, options: [.prettyPrinted]) { try? data.write(to: out) }
    printLockWait(); printCacheStats()
    print("結果：\(out.path)")
    exit(0)
}

print("""
🍵 森日茶屋客服（新版 bot.swift）
   /reset 換新顧客  /as 使用者ID 切換顧客  /reload 重讀設定  /me 名字  /pending 待處理清單  /quit 離開
   店員：/reply 編號 文字  /confirm 編號  /reject 編號 [原因]  /endhuman [使用者ID] 結束轉人工
   紀錄：\(logURL.path)

""")

let icons = ["answer": "🤖", "fallback": "📋", "booking_done": "📅", "booking_cancel": "📅", "unsupported": "🌐", "silent": "🙋"]
// 記憶體資料庫：和正式系統同一套資料表，但不寫入正式資料（離開就消失）
let memoryDB = try! BotDatabase()
let bot = Bot(customers: SQLiteCustomerStore(memoryDB), log: SQLiteConversationLog(memoryDB), model: model)
var userID = "user1"   // 目前在說話的顧客（/as 切換）
var me: CustomerKey { CustomerKey(platform: "terminal", userID: userID) }

// 店員操作產生的動作：轉給顧客的回覆
func show(_ actions: [BotAction]) {
    for case let .reply(key, text, chatID) in actions {
        print("📨 → 顧客 \(chatID)：\(text)  \u{1B}[90m[\(key)]\u{1B}[0m")
        log += "**店員 → \(chatID)**（\(key)）：\(text)\n\n"
    }
    print(""); save()
}
func staff(_ op: () async throws -> [BotAction]) async {
    do { show(try await op()) } catch { print("⚠️ \(error)\n") }
}
// 「/confirm 3 原因…」→ (3, "原因…")
func ticketArgs(_ line: String) -> (Int, String)? {
    let parts = line.split(separator: " ", maxSplits: 2).map(String.init)
    guard parts.count >= 2, let id = Int(parts[1].trimmingCharacters(in: CharacterSet(charactersIn: "#"))) else { return nil }
    return (id, parts.count > 2 ? parts[2].trimmingCharacters(in: .whitespaces) : "")
}

while true {
    print("👤 \(userID)：", terminator: ""); fflush(stdout)
    guard let line = readLine()?.trimmingCharacters(in: .whitespaces) else { break }
    if line.isEmpty { continue }
    switch line {
    case "/quit", "/exit":
        save(); printLockWait(); printCacheStats(toStderr: true); print("👋 紀錄已存：\(logURL.path)"); exit(0)
    case "/reset":
        trace(["key": "reset"])
        bot.reset(me)
        log += "---\n*（換一位新顧客：\(userID)）*\n\n"; print("🔄 已換一位新顧客\n"); continue
    case "/reload":
        // 設定有錯：印出錯誤，繼續用舊的設定
        do { applyConfig(try loadConfig()) } catch { print("⚠️ 設定有錯，沿用舊的設定：\n\(error)\n"); continue }
        print("📄 已重新讀取 store.json、replies.json、faq.txt、examples.jsonl（\(examples.count) 個案例）\n"); continue
    case "/pending":
        // trace 的 pending 是字串陣列，內容含原本的通知文字（流程測試用子字串比對）
        let items = bot.pendingTickets.map { "#\($0.id) \($0.text)" }
        trace(["key": "pending", "pending": items])
        print(items.isEmpty ? "📭 沒有待處理項目\n"
                            : "📬 待店員處理：\n" + zip(items, bot.pendingTickets).map { "   \($0)　\u{1B}[90m（\($1.customer.userID)）\u{1B}[0m" }.joined(separator: "\n") + "\n")
        continue
    case _ where line.hasPrefix("/as "):
        userID = String(line.dropFirst(4)).trimmingCharacters(in: .whitespaces)
        let c = bot.customer(me)
        print("👥 切換到顧客 \(userID)（\(c.name ?? "沒有名字")·\(c.lang)\(c.human ? "·轉人工中" : "")）\n")
        log += "---\n*（切換到顧客 \(userID)）*\n\n"; continue
    case _ where line.hasPrefix("/me "):
        let name = String(line.dropFirst(4)).trimmingCharacters(in: .whitespaces)
        bot.update(me) { $0.name = name }
        print("🪪 已記住顧客名字：\(name)\n"); continue
    case _ where line.hasPrefix("/reply "):
        guard case let (id, text)? = ticketArgs(line), !text.isEmpty else { print("用法：/reply 編號 文字\n"); continue }
        await staff { try await bot.staffReply(ticketID: id, text: text) }; continue
    case _ where line.hasPrefix("/confirm "):
        guard case let (id, _)? = ticketArgs(line) else { print("用法：/confirm 編號\n"); continue }
        await staff { try await bot.confirmBooking(id) }; continue
    case _ where line.hasPrefix("/reject "):
        guard case let (id, reason)? = ticketArgs(line) else { print("用法：/reject 編號 [原因]\n"); continue }
        await staff { try await bot.rejectBooking(id, reason: reason) }; continue
    case _ where line == "/endhuman" || line.hasPrefix("/endhuman "):
        let id = String(line.dropFirst(9)).trimmingCharacters(in: .whitespaces)
        let key = id.isEmpty ? me : CustomerKey(platform: "terminal", userID: id)
        show(await bot.endHandoff(key))
        print("✅ 顧客 \(key.userID) 已結束轉人工，恢復 AI 回覆\n"); continue
    default: break
    }

    let t = Date()
    let actions = await bot.receive(IncomingMessage(platform: "terminal", userID: userID, text: line, time: now()))
    let secs = String(format: "%.2f", Date().timeIntervalSince(t))
    let (key, reply) = replyOf(actions)
    let lang = bot.customer(me).lang
    trace(["key": key, "reply": reply, "lang": lang])
    for case let .openTicket(ticket) in actions { print("📬 已開待處理 #\(ticket.id)：\(ticket.text)") }
    let icon = icons[key] ?? (key.hasPrefix("ask_") || key.hasPrefix("err_") ? "📅" : "📋")
    print(key == "silent" ? "🙋 （轉人工中，AI 不回覆；/endhuman 結束）\n"
                          : "\(icon) 小茶：\(reply)  \u{1B}[90m[\(key)·\(lang)] \(secs)s\u{1B}[0m\n")
    log += "**顧客 \(userID)**：\(line)\n\n**小茶**（\(key)）：\(reply)\n\n"
    save()
}
save()
printLockWait(); printCacheStats(toStderr: true)
