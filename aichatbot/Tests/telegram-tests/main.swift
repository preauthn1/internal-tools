// Webhook 層測試：swift build && .build/debug/telegram-tests（全部通過 exit 0，有失敗 exit 1）
// 在本機啟動伺服器、送真的 HTTP 請求；Telegram API 換成本機的假端點，記錄每個請求的 method 和內容，可以設定前幾次回 500。
// 核心用照腳本回答的假模型（ScriptedModel），不需要本機模型；重試間隔和「正在輸入」的間隔調短，避免拖慢速度
import Foundation
import Hummingbird
import BotCore
import Telegram
import Storage

_ = (store, replies, faqs, examples)

var failures = 0, current = ""
@MainActor func test(_ name: String, _ body: @MainActor () async throws -> Void) async {
    current = name
    let before = failures
    do { try await body() } catch { check(false, "丟出錯誤：\(error)") }
    print(failures == before ? "✔ \(name)" : "✘ \(name)")
}
func check(_ ok: Bool, _ what: @autoclosure () -> String = "", line: Int = #line) {
    if !ok { failures += 1; print("   ✘ 第 \(line) 行：\(what())") }
}

// 等到條件成立（逾時回傳 false）
@MainActor func waitUntil(timeout: Double = 5, _ cond: () -> Bool) async -> Bool {
    let end = Date().addingTimeInterval(timeout)
    while Date() < end {
        if cond() { return true }
        try? await Task.sleep(for: .milliseconds(20))
    }
    return cond()
}

final class Box<T>: @unchecked Sendable {
    private let lock = NSLock()
    private var _value: T?
    var value: T? {
        get { lock.withLock { _value } }
        set { lock.withLock { _value = newValue } }
    }
}

// MARK: - 假 Telegram API：POST /bot<token>/<method>

final class FakeTelegram: @unchecked Sendable {
    // messageID：sendMessage 成功時回給伺服器的訊息編號（店員群組的通知要用它來回覆）
    struct Call { let path: String; let method: String; let body: [String: Any]; let messageID: Int64? }
    private let lock = NSLock()
    private var _calls: [Call] = []
    private var failing: [String: Int] = [:]   // method → 還要回幾次 500
    private var nextMessageID: Int64 = 5000
    let port = Box<Int>()

    var calls: [Call] { lock.withLock { _calls } }
    func calls(_ method: String, chat: String? = nil) -> [Call] {
        calls.filter { $0.method == method && (chat == nil || $0.body["chat_id"] as? String == chat) }
    }
    func fail(_ method: String, times: Int) { lock.withLock { failing[method] = times } }

    func start() {
        let router = Router()
        router.post("/:bot/:method") { request, context -> Response in
            let method = context.parameters.get("method") ?? ""
            let buffer = try await request.body.collect(upTo: 1 << 20)
            let body = (try? JSONSerialization.jsonObject(with: Data(buffer.readableBytesView))) as? [String: Any] ?? [:]
            let (fail, id): (Bool, Int64?) = self.lock.withLock {
                if let n = self.failing[method], n > 0 {
                    self.failing[method] = n - 1
                    self._calls.append(Call(path: request.uri.path, method: method, body: body, messageID: nil))
                    return (true, nil)
                }
                var newID: Int64?
                if method == "sendMessage" { self.nextMessageID += 1; newID = self.nextMessageID }
                self._calls.append(Call(path: request.uri.path, method: method, body: body, messageID: newID))
                return (false, newID)
            }
            let result = id.map { #"{"message_id":\#($0)}"# } ?? "true"
            return fail
                ? Response(status: .internalServerError, body: .init(byteBuffer: ByteBuffer(string: #"{"ok":false,"description":"boom"}"#)))
                : Response(status: .ok, body: .init(byteBuffer: ByteBuffer(string: #"{"ok":true,"result":\#(result)}"#)))
        }
        let app = Application(router: router, configuration: .init(address: .hostname("0.0.0.0", port: 0)),
                              onServerRunning: { self.port.value = $0.localAddress?.port })
        Task { try? await app.run() }
    }
}

// MARK: - 準備：假端點＋被測伺服器

let token = [REDACTED]"
let secret = [REDACTED]"
let staffChat: Int64 = -1001
let logLines = Box<[String]>()
logLines.value = []
// 測試用的日誌不遮罩：這樣才驗得出程式本身沒有把 token、secret 寫進日誌（正式執行另外會遮罩）
let log = Log(write: { line in logLines.value = (logLines.value ?? []) + [line] })
func logs() -> [String] { logLines.value ?? [] }

let fake = FakeTelegram()
fake.start()
guard await waitUntil({ fake.port.value != nil }), let fakePort = fake.port.value else { print("❌ 假端點沒有啟動"); exit(1) }

var config = TelegramConfig(token: [REDACTED], secret: [REDACTED], publicURL: "https://bot.example.test/", staffChatID: staffChat)
config.apiBase = "http://0.0.0.0:\(fakePort)"
config.port = 0
config.retryDelays = [0.01, 0.02, 0.04]
config.typingInterval = 0.2
let modelDelay = 0.7   // 每次 AI 呼叫的時間：比「正在輸入」間隔的兩倍長
let model = ScriptedModel(delay: modelDelay)   // 沒有腳本：分類失敗，回「請說清楚」
let server = TelegramServer(config: config, bot: Bot(model: model), log: log)
let serverPort = Box<Int>()
Task { try? await server.run { serverPort.value = $0 } }
guard await waitUntil({ serverPort.value != nil }), let port = serverPort.value else { print("❌ 伺服器沒有啟動"); exit(1) }

var nextUpdateID = 1000
func update(_ fields: [String: Any]) -> [String: Any] {
    nextUpdateID += 1
    return fields.merging(["update_id": nextUpdateID]) { a, _ in a }
}
func message(chat: Int64, type: String, from: Int64? = nil, text: String? = "你好") -> [String: Any] {
    var m: [String: Any] = ["message_id": 1, "date": 0, "chat": ["id": chat, "type": type]]
    if let from { m["from"] = ["id": from, "is_bot": false, "first_name": "小明"] }
    if let text { m["text"] = text }
    return m
}
func privateText(_ user: Int64, _ text: String) -> [String: Any] {
    update(["message": message(chat: user, type: "private", from: user, text: text)])
}

func post(_ body: [String: Any], secret: [REDACTED] = secret, to target: Int? = nil) async throws -> Int {
    var r = URLRequest(url: URL(string: "http://0.0.0.0:\(target ?? port)/telegram/webhook")!)
    r.httpMethod = "POST"
    r.setValue("application/json", forHTTPHeaderField: "Content-Type")
    r.httpBody = try JSONSerialization.data(withJSONObject: body)
    if let secret { r.setValue(secret, forHTTPHeaderField: "X-Telegram-Bot-Api-Secret-Token") }
    let (_, response) = try await URLSession.shared.data(for: r)
    return (response as! HTTPURLResponse).statusCode
}
// 送出並等背景處理完
@MainActor func postAndWait(_ body: [String: Any]) async throws {
    let before = server.processedCount
    check(try await post(body) == 200)
    check(await waitUntil { server.processedCount > before }, "背景處理逾時")
}

let greetings = Set((replies["greeting"]?["zh"] ?? "").components(separatedBy: "|"))

// MARK: - 啟動

await test("啟動時設定 webhook：帶 secret，只收一般訊息和按鈕回呼") {
    check(await waitUntil { !fake.calls("setWebhook").isEmpty }, "沒有呼叫 setWebhook")
    guard let call = fake.calls("setWebhook").first else { return }
    check(call.path.removingPercentEncoding == "/bot\(token)/setWebhook", "網址的 token 不對")
    check(call.body["url"] as? String == "https://bot.example.test/telegram/webhook", "\(call.body)")
    check(call.body["secret_token"] as? String == secret)
    check(call.body["allowed_updates"] as? [String] == ["message", "callback_query"], "\(call.body)")
}

// MARK: - 收請求

await test("secret 錯或沒帶：回 401，不處理") {
    let before = fake.calls.count, processed = server.processedCount
    check(try await post(privateText(201, "你好"), secret: [REDACTED]") == 401)
    check(try await post(privateText(201, "你好"), secret: [REDACTED] == 401)
    check(try await post(privateText(201, "你好"), secret: [REDACTED] + "x") == 401)
    try? await Task.sleep(for: .milliseconds(200))
    check(fake.calls.count == before, "不該呼叫 Telegram API：\(fake.calls.dropFirst(before).map(\.method))")
    check(server.processedCount == processed)
}

await test("同一個 update 送兩次：都回 200，只回覆一次") {
    let u = privateText(202, "/start")
    try await postAndWait(u)
    check(try await post(u) == 200)
    try? await Task.sleep(for: .milliseconds(200))
    check(fake.calls("sendMessage", chat: "202").count == 1, "\(fake.calls("sendMessage", chat: "202").count)")
    check(logs().contains { $0.contains("重複") })
}

await test("讀不懂的內容：回 200（Telegram 才不會一直重送），不處理") {
    var r = URLRequest(url: URL(string: "http://0.0.0.0:\(port)/telegram/webhook")!)
    r.httpMethod = "POST"; r.httpBody = Data("not json".utf8)
    r.setValue(secret, forHTTPHeaderField: "X-Telegram-Bot-Api-Secret-Token")
    let (_, response) = try await URLSession.shared.data(for: r)
    check((response as! HTTPURLResponse).statusCode == 200)
}

// MARK: - 私訊

await test("私訊：先回 200 再處理；先送正在輸入，處理超過間隔會再送，最後送出回覆") {
    let calls = model.requests.count
    let start = Date()
    check(try await post(privateText(203, "你們營業到幾點呀")) == 200)
    let elapsed = Date().timeIntervalSince(start)
    check(elapsed < modelDelay, "回 200 花了 \(elapsed) 秒，應該在處理完之前就回")
    check(fake.calls("sendMessage", chat: "203").isEmpty, "回 200 時還不該送出回覆")

    check(await waitUntil { !fake.calls("sendMessage", chat: "203").isEmpty }, "沒有送出回覆")
    let sequence = fake.calls.filter { $0.body["chat_id"] as? String == "203" }.map(\.method)
    check(sequence.first == "sendChatAction", "\(sequence)")
    let typingBeforeReply = sequence.prefix { $0 != "sendMessage" }.count
    check(typingBeforeReply >= 2, "正在輸入只送了 \(typingBeforeReply) 次：\(sequence)")
    check(fake.calls("sendChatAction", chat: "203").allSatisfy { $0.body["action"] as? String == "typing" })
    let reply = fake.calls("sendMessage", chat: "203").first?.body["text"] as? String ?? ""
    check(!reply.isEmpty)
    check(model.requests.count == calls + 1, "應該呼叫 AI 一次（分類）")

    // 回覆送出後停止送正在輸入
    try? await Task.sleep(for: .milliseconds(500))
    let typingAfter = fake.calls("sendChatAction", chat: "203").count
    try? await Task.sleep(for: .milliseconds(500))
    check(fake.calls("sendChatAction", chat: "203").count == typingAfter, "回覆後還在送正在輸入")
}

await test("/start：回打招呼，不呼叫 AI，不送正在輸入") {
    let calls = model.requests.count
    try await postAndWait(privateText(204, "/start"))
    try await postAndWait(privateText(205, "/start promo-1"))
    for chat in ["204", "205"] {
        let sent = fake.calls("sendMessage", chat: chat)
        check(sent.count == 1 && greetings.contains(sent[0].body["text"] as? String ?? ""), "\(sent.map(\.body))")
        check(fake.calls("sendChatAction", chat: chat).isEmpty)
    }
    check(model.requests.count == calls, "/start 不該呼叫 AI")
}

await test("私訊是不認得的種類（位置）：不回覆") {
    var m = message(chat: 206, type: "private", from: 206, text: nil)
    m["location"] = ["latitude": 25.0, "longitude": 121.5]
    try await postAndWait(update(["message": m]))
    check(fake.calls.filter { $0.body["chat_id"] as? String == "206" }.isEmpty)
}

// MARK: - 店員群組

let staffChatID = String(staffChat)
// 店員群組收到的通知（不含回在通知上的結果）
func notices() -> [FakeTelegram.Call] {
    fake.calls("sendMessage", chat: staffChatID).filter { $0.body["reply_parameters"] == nil }
}
func replyTarget(_ c: FakeTelegram.Call?) -> Int64? {
    ((c?.body["reply_parameters"] as? [String: Any])?["message_id"] as? Int).map(Int64.init)
}
func text(_ c: FakeTelegram.Call?) -> String { c?.body["text"] as? String ?? "" }
// 店員在群組（或其他群組）發言；replyTo 是回覆哪一則
func groupMessage(replyTo: Int64?, text: String?, chat: Int64 = staffChat, extra: [String: Any] = [:]) -> [String: Any] {
    var m = message(chat: chat, type: "supergroup", from: 301, text: text)
    if let replyTo { m["reply_to_message"] = ["message_id": replyTo, "date": 0, "chat": ["id": chat, "type": "supergroup"], "text": "通知"] }
    return update(["message": m.merging(extra) { _, b in b }])
}
// 顧客私訊（可以帶照片、貼圖這類欄位）
func privateMessage(_ user: Int64, text: String?, extra: [String: Any] = [:], from: [String: Any]? = nil) -> [String: Any] {
    var m = message(chat: user, type: "private", from: user, text: text)
    if let from { m["from"] = from }
    return update(["message": m.merging(extra) { _, b in b }])
}

await test("開待處理事項：通知店員群組，附上類型、Telegram 名稱、@username、使用者 ID、原話") {
    let before = notices().count
    try await postAndWait(privateMessage(601, text: "안녕하세요 담배",
                                         from: ["id": 601, "is_bot": false, "first_name": "Min", "last_name": "Kim", "username": "minkim"]))
    let sent = Array(notices().dropFirst(before))
    check(sent.count == 1, "\(sent.map(\.body))")
    let n = text(sent.first)
    check(n.contains("不支援的語言") && n.contains("Min Kim") && n.contains("@minkim") && n.contains("ID 601") && n.contains("안녕하세요 담배"), n)
    check(sent.first?.messageID != nil)
    check(fake.calls("sendMessage", chat: "601").count == 1, "顧客也要收到回覆")
}

await test("店員回覆通知：顧客收到加前綴的原文，群組在通知上回 ✅") {
    guard let notice = notices().last(where: { text($0).contains("안녕하세요 담배") })?.messageID else { check(false, "沒有通知"); return }
    let calls = model.requests.count
    try await postAndWait(groupMessage(replyTo: notice, text: "請用中文或英文喔"))
    let toCustomer = fake.calls("sendMessage", chat: "601")
    check(toCustomer.count == 2 && text(toCustomer.last) == "💬 店員回覆：請用中文或英文喔", "\(toCustomer.map(\.body))")
    let feedback = fake.calls("sendMessage", chat: staffChatID).last
    check(text(feedback) == "✅ 已送給顧客" && replyTarget(feedback) == notice, "\(feedback?.body ?? [:])")
    check(model.requests.count == calls, "店員回覆不經過 AI")

    // 已回覆過的通知照樣轉送
    try await postAndWait(groupMessage(replyTo: notice, text: "補充：營業到九點"))
    check(text(fake.calls("sendMessage", chat: "601").last) == "💬 店員回覆：補充：營業到九點")
}

await test("店員回覆照片、貼圖：不轉送，在通知上回「目前只能轉送文字」") {
    guard let notice = notices().last(where: { text($0).contains("안녕하세요 담배") })?.messageID else { check(false, "沒有通知"); return }
    let before = fake.calls("sendMessage", chat: "601").count
    try await postAndWait(groupMessage(replyTo: notice, text: nil, extra: ["photo": [["file_id": "p", "width": 1, "height": 1]]]))
    try await postAndWait(groupMessage(replyTo: notice, text: nil, extra: ["sticker": ["file_id": "s", "emoji": "👍"]]))
    check(fake.calls("sendMessage", chat: "601").count == before, "不該轉給顧客")
    let feedback = fake.calls("sendMessage", chat: staffChatID).suffix(2)
    check(feedback.allSatisfy { text($0) == "目前只能轉送文字" && replyTarget($0) == notice }, "\(feedback.map(\.body))")
}

await test("店員群組沒有回覆通知的發言、回覆的不是通知：忽略，不在群組回話") {
    let before = fake.calls.filter { $0.method == "sendMessage" }.count
    try await postAndWait(groupMessage(replyTo: nil, text: "大家午安"))
    try await postAndWait(groupMessage(replyTo: 424242, text: "這是回覆別的訊息"))
    check(fake.calls.filter { $0.method == "sendMessage" }.count == before, "不該送出任何訊息")
    check(logs().contains { $0.contains("不是回覆通知") } && logs().contains { $0.contains("不是待處理事項的通知") })
}

await test("非店員群組回覆通知：不轉送") {
    guard let notice = notices().last(where: { text($0).contains("안녕하세요 담배") })?.messageID else { check(false, "沒有通知"); return }
    let before = fake.calls("sendMessage", chat: "601").count
    try await postAndWait(groupMessage(replyTo: notice, text: "我不是店員", chat: -2002))
    check(fake.calls("sendMessage", chat: "601").count == before)
    check(fake.calls.filter { $0.body["chat_id"] as? String == "-2002" }.isEmpty)
}

await test("轉給顧客失敗：群組在通知上回 ❌") {
    let before = notices().count
    try await postAndWait(privateText(602, "안녕하세요"))
    guard let notice = notices().dropFirst(before).first?.messageID else { check(false, "沒有通知"); return }
    let customerCalls = fake.calls("sendMessage", chat: "602").count
    fake.fail("sendMessage", times: 4)   // 顧客那則（1+3 次重試）全部失敗，群組那則成功
    try await postAndWait(groupMessage(replyTo: notice, text: "您好"))
    check(fake.calls("sendMessage", chat: "602").count == customerCalls + 4)
    let feedback = fake.calls("sendMessage", chat: staffChatID).last
    check(text(feedback) == "❌ 傳送失敗" && replyTarget(feedback) == notice, "\(feedback?.body ?? [:])")
}

// MARK: - 非文字訊息

await test("顧客傳照片（附說明文字）、語音、檔案：回「目前只看得懂文字」並通知店員群組，不呼叫 AI") {
    let calls = model.requests.count
    let cases: [(Int64, [String: Any], String)] = [
        (611, ["photo": [["file_id": "p", "width": 1, "height": 1]], "caption": "這個還有貨嗎"], "圖片"),
        (612, ["voice": ["file_id": "v", "duration": 3]], "語音"),
        (613, ["document": ["file_id": "d", "file_name": "menu.pdf"]], "檔案"),
    ]
    for (user, extra, kind) in cases {
        let before = notices().count
        try await postAndWait(privateMessage(user, text: nil, extra: extra))
        let sent = fake.calls("sendMessage", chat: String(user))
        check(sent.count == 1 && text(sent.first) == replyText("text_only", lang: "zh"), "\(sent.map(\.body))")
        let n = text(notices().dropFirst(before).first)
        check(n.contains("📎") && n.contains(kind), n)
        if user == 611 { check(n.contains("這個還有貨嗎"), n) }
    }
    check(model.requests.count == calls)
}

await test("顧客傳貼圖：當閒聊回覆，不開待處理事項、不呼叫 AI") {
    let calls = model.requests.count, before = notices().count
    try await postAndWait(privateMessage(614, text: nil, extra: ["sticker": ["file_id": "s", "emoji": "😂"]]))
    let small = Set((replies["small_talk"]?["zh"] ?? "").components(separatedBy: "|"))
    let sent = fake.calls("sendMessage", chat: "614")
    check(sent.count == 1 && small.contains(text(sent.first)), "\(sent.map(\.body))")
    check(notices().count == before)
    check(model.requests.count == calls)
}

// MARK: - 分流

await test("其他群組、超級群組、頻道：忽略") {
    let calls = model.requests.count, apiCalls = fake.calls.filter { $0.method != "sendChatAction" }.count
    try await postAndWait(update(["message": message(chat: -2002, type: "group", from: 302, text: "你好")]))
    try await postAndWait(update(["message": message(chat: -1003, type: "supergroup", from: 303, text: "/start")]))
    try await postAndWait(update(["channel_post": message(chat: -1004, type: "channel", text: "公告")]))
    for chat in ["-2002", "-1003", "-1004"] {
        check(fake.calls.filter { $0.body["chat_id"] as? String == chat }.isEmpty, "不該送到 \(chat)")
    }
    check(fake.calls.filter { $0.method != "sendChatAction" }.count == apiCalls)
    check(model.requests.count == calls)
}

// MARK: - API 失敗重試

await test("API 前 2 次失敗：重試後成功，共送 3 次") {
    let logCount = logs().count
    fake.fail("sendMessage", times: 2)
    try await postAndWait(privateText(401, "/start"))
    check(fake.calls("sendMessage", chat: "401").count == 3, "\(fake.calls("sendMessage", chat: "401").count)")
    check(!logs().dropFirst(logCount).contains { $0.contains("sendMessage 失敗 4 次") })
}

await test("API 一直失敗：送 4 次後記錄錯誤，伺服器繼續處理下一個 update") {
    fake.fail("sendMessage", times: 100)
    try await postAndWait(privateText(402, "/start"))
    check(fake.calls("sendMessage", chat: "402").count == 4, "\(fake.calls("sendMessage", chat: "402").count)")
    check(logs().contains { $0.contains("❌") && $0.contains("sendMessage 失敗 4 次") }, "沒有錯誤紀錄")
    fake.fail("sendMessage", times: 0)
    try await postAndWait(privateText(403, "/start"))
    check(fake.calls("sendMessage", chat: "403").count == 1, "伺服器沒有繼續處理")
}

await test("網路錯誤（連不上）：重試後記錄錯誤，錯誤訊息不帶網址") {
    let api = TelegramAPI(token: [REDACTED], baseURL: "http://0.0.0.0:1", retryDelays: [0.01], log: log)
    check(await api.sendMessage("1", "hi") == false)
    check(logs().contains { $0.contains("sendMessage 失敗 2 次") })
}

// MARK: - 資料庫：重開伺服器

// 用資料庫檔案啟動一台伺服器（和 server 執行檔相同的接法），回傳伺服器和它的 port
@MainActor func startServer(_ db: BotDatabase, model: ChatModel = ScriptedModel()) async -> (TelegramServer, Int)? {
    await startServer(Bot(customers: SQLiteCustomerStore(db), tickets: SQLiteTicketStore(db), log: SQLiteConversationLog(db), model: model),
                      dedup: SQLiteUpdateDedup(db))
}
@MainActor func startServer(_ bot: Bot, dedup: UpdateDedup) async -> (TelegramServer, Int)? {
    let s = TelegramServer(config: config, bot: bot, dedup: dedup, log: log)
    let p = Box<Int>()
    Task { try? await s.run { p.value = $0 } }
    guard await waitUntil({ p.value != nil }), let port = p.value else { return nil }
    return (s, port)
}

await test("重開伺服器（同一個資料庫）：同一個 update 不再處理，顧客狀態和對話紀錄接續") {
    let dir = FileManager.default.temporaryDirectory.appendingPathComponent("telegram-tests-\(getpid())")
    defer { try? FileManager.default.removeItem(at: dir) }
    let path = dir.appendingPathComponent("bot.sqlite").path
    let u = privateText(501, "안녕하세요")

    // 第一台：處理完就不再用（伺服器在測試程序結束前不會真的停下，之後的請求都送到第二台）
    guard let (first, firstPort) = await startServer(try BotDatabase(path: path)) else { check(false, "第一台沒有啟動"); return }
    check(try await post(u, to: firstPort) == 200)
    check(await waitUntil { first.processedCount == 1 }, "第一台處理逾時")
    check(fake.calls("sendMessage", chat: "501").count == 1)

    let db = try BotDatabase(path: path)
    guard let (second, secondPort) = await startServer(db) else { check(false, "第二台沒有啟動"); return }
    check(try await post(u, to: secondPort) == 200)
    try? await Task.sleep(for: .milliseconds(200))
    check(second.processedCount == 0, "同一個 update 不該再處理")
    check(fake.calls("sendMessage", chat: "501").count == 1, "\(fake.calls("sendMessage", chat: "501").count)")

    check(try await post(privateText(501, "😀"), to: secondPort) == 200)
    check(await waitUntil { second.processedCount == 1 }, "第二台處理逾時")
    check(fake.calls("sendMessage", chat: "501").count == 2)
    let key = CustomerKey(platform: "telegram", userID: "501")
    let c = SQLiteCustomerStore(db).load(key)
    check(c.lang == "ko" && c.openTicket, "第一台存的顧客狀態沒有接續：\(c)")
    check(SQLiteConversationLog(db).entries(for: key).map(\.direction) == [.customer, .bot, .customer, .bot])
}

await test("追問 3 次：緊急通知引用原本那張；店員回覆緊急通知時兩張都改成已回覆") {
    let followUp: [String: FieldValue] = ["intent": .text("follow_up"), "self_name": .text("")]
    let db = try BotDatabase()
    let bot = Bot(customers: SQLiteCustomerStore(db), tickets: SQLiteTicketStore(db), log: SQLiteConversationLog(db),
                  model: ScriptedModel(["Intent": [followUp, followUp, followUp]]))
    guard let (s, p) = await startServer(bot, dedup: SQLiteUpdateDedup(db)) else { check(false, "伺服器沒有啟動"); return }
    func send(_ body: [String: Any]) async throws {
        let before = s.processedCount
        check(try await post(body, to: p) == 200)
        check(await waitUntil { s.processedCount > before }, "背景處理逾時")
    }

    let before = notices().count
    try await send(privateText(701, "안녕하세요"))
    for _ in 0..<3 { try await send(privateText(701, "好了沒")) }
    let sent = Array(notices().dropFirst(before))
    check(sent.count == 2, "\(sent.map(\.body))")
    guard sent.count == 2, let urgentNotice = sent[1].messageID else { return }
    check(text(sent[1]).contains("緊急") && text(sent[1]).contains("안녕하세요"), text(sent[1]))

    try await send(groupMessage(replyTo: urgentNotice, text: "抱歉久等了"))
    check(text(fake.calls("sendMessage", chat: "701").last) == "💬 店員回覆：抱歉久等了")
    check(bot.tickets.tickets.count == 2 && bot.tickets.tickets.allSatisfy(\.done), "\(bot.tickets.tickets)")
    check(!bot.customer(CustomerKey(platform: "telegram", userID: "701")).openTicket)
}

// MARK: - 預約按鈕（另開一台資料庫版的伺服器；顧客 801～803 先停在預約的最後一步，只差電話）

let bookingIntent: [String: FieldValue] = ["intent": .text("booking"), "self_name": .text("")]
let bookingDB = try BotDatabase()
let bookingBot = Bot(customers: SQLiteCustomerStore(bookingDB), tickets: SQLiteTicketStore(bookingDB), log: SQLiteConversationLog(bookingDB),
                     model: ScriptedModel(["Intent": Array(repeating: bookingIntent, count: 3),
                                           "Booking": Array(repeating: ["phone": .text("0912345678")], count: 3)]))
guard let (bookingServer, bookingPort) = await startServer(bookingBot, dedup: SQLiteUpdateDedup(bookingDB)) else { print("❌ 預約伺服器沒有啟動"); exit(1) }
@MainActor func sendToBookingServer(_ body: [String: Any]) async throws {
    let before = bookingServer.processedCount
    check(try await post(body, to: bookingPort) == 200)
    check(await waitUntil { bookingServer.processedCount > before }, "背景處理逾時")
}
// 顧客走完預約，回傳預約編號和群組那則通知
@MainActor func book(_ user: Int64) async throws -> (id: Int, notice: FakeTelegram.Call)? {
    bookingBot.update(CustomerKey(platform: "telegram", userID: String(user))) {
        $0.booking = ["date": "2026-10-21", "time": "14:00", "people": "4", "name": "小康"]; $0.askingField = "phone"
    }
    let before = notices().count
    try await sendToBookingServer(privateText(user, "0912345678"))
    guard let notice = notices().dropFirst(before).first, let id = bookingBot.tickets.tickets.last?.id else { return nil }
    return (id, notice)
}
var nextCallbackID = 0
// 店員按下按鈕；chat 是按鈕所在的聊天室（nil = Telegram 沒帶訊息）
func callback(_ data: String?, messageID: Int64, chat: Int64? = staffChat, text: String = "通知",
              from: [String: Any] = ["id": 301, "is_bot": false, "first_name": "小明", "last_name": "王", "username": "ming"]) -> [String: Any] {
    nextCallbackID += 1
    var q: [String: Any] = ["id": "cb\(nextCallbackID)", "from": from, "chat_instance": "x"]
    if let data { q["data"] = data }
    if let chat { q["message"] = ["message_id": messageID, "date": 0, "chat": ["id": chat, "type": "supergroup"], "text": text] }
    return update(["callback_query": q])
}
func callbackAnswer(_ id: Int) -> FakeTelegram.Call? {
    fake.calls("answerCallbackQuery").last { $0.body["callback_query_id"] as? String == "cb\(id)" }
}

var firstBooking: (id: Int, notice: FakeTelegram.Call)?

await test("預約通知：附「確認／拒絕」按鈕，callback_data 帶預約編號") {
    firstBooking = try await book(801)
    guard let (id, notice) = firstBooking else { check(false, "沒有預約通知"); return }
    check(text(notice).contains("預約") && text(notice).contains("0912345678"), text(notice))
    let keyboard = (notice.body["reply_markup"] as? [String: Any])?["inline_keyboard"] as? [[[String: Any]]]
    let data = keyboard?.flatMap { $0.compactMap { $0["callback_data"] as? String } }
    check(data == ["confirm:\(id)", "reject:\(id)"], "\(notice.body)")
    check(bookingBot.booking(id)?.status == .pending)
}

await test("店員按確認：顧客收到確認，群組通知改成結果和店員並拿掉按鈕，預約改成已確認") {
    guard let (id, notice) = firstBooking, let messageID = notice.messageID else { check(false, "沒有預約通知"); return }
    let before = fake.calls("sendMessage", chat: "801").count
    try await sendToBookingServer(callback("confirm:\(id)", messageID: messageID, text: text(notice)))

    let toCustomer = fake.calls("sendMessage", chat: "801")
    check(toCustomer.count == before + 1, "\(toCustomer.map(\.body))")
    let summary = bookingBot.booking(id)?.summary ?? ""
    check(text(toCustomer.last) == replyText("booking_confirmed", ["summary": summary], lang: "zh"), text(toCustomer.last))
    let edit = fake.calls("editMessageText", chat: staffChatID).last
    check((edit?.body["message_id"] as? Int).map(Int64.init) == messageID, "\(edit?.body ?? [:])")
    check(text(edit).hasPrefix(text(notice)) && text(edit).contains("✅ 已確認") && text(edit).contains("小明 王 @ming"), text(edit))
    check(edit?.body["reply_markup"] == nil, "按鈕應該拿掉")
    check(callbackAnswer(nextCallbackID)?.body["text"] as? String == "已確認")
    let b = bookingBot.booking(id)
    check(b?.status == .confirmed && b?.handledBy == "小明 王 @ming", "\(String(describing: b))")

    // 已處理過的再按（確認或拒絕）：只回「已處理」，顧客不會再收到
    let edits = fake.calls("editMessageText").count
    for data in ["confirm:\(id)", "reject:\(id)"] {
        try await sendToBookingServer(callback(data, messageID: messageID))
        check(callbackAnswer(nextCallbackID)?.body["text"] as? String == "已處理", "\(data)")
    }
    check(fake.calls("sendMessage", chat: "801").count == before + 1)
    check(fake.calls("editMessageText").count == edits)
    check(bookingBot.booking(id)?.status == .confirmed)
}

await test("店員按拒絕：顧客收到「無法安排」，群組通知顯示已拒絕和店員，預約改成已拒絕") {
    guard let (id, notice) = try await book(802), let messageID = notice.messageID else { check(false, "沒有預約通知"); return }
    try await sendToBookingServer(callback("reject:\(id)", messageID: messageID, text: text(notice), from: ["id": 302, "is_bot": false, "first_name": "阿華"]))

    let summary = bookingBot.booking(id)?.summary ?? ""
    check(text(fake.calls("sendMessage", chat: "802").last) == replyText("booking_rejected", ["summary": summary], lang: "zh"),
          text(fake.calls("sendMessage", chat: "802").last))
    let edit = fake.calls("editMessageText", chat: staffChatID).last
    check(text(edit).contains("❌ 已拒絕（操作：阿華）"), text(edit))
    check(callbackAnswer(nextCallbackID)?.body["text"] as? String == "已拒絕")
    check(bookingBot.booking(id)?.status == .rejected && bookingBot.booking(id)?.handledBy == "阿華")
}

await test("不是店員群組的按鈕回呼、沒帶訊息的回呼：只回應回呼，不動預約、不送給顧客") {
    guard let (id, notice) = try await book(803), let messageID = notice.messageID else { check(false, "沒有預約通知"); return }
    let before = fake.calls("sendMessage", chat: "803").count, edits = fake.calls("editMessageText").count
    try await sendToBookingServer(callback("confirm:\(id)", messageID: messageID, chat: -2002))
    check(callbackAnswer(nextCallbackID) != nil, "還是要回應回呼")
    try await sendToBookingServer(callback("confirm:\(id)", messageID: messageID, chat: nil))
    check(callbackAnswer(nextCallbackID) != nil, "還是要回應回呼")
    check(fake.calls("sendMessage", chat: "803").count == before)
    check(fake.calls("editMessageText").count == edits)
    check(fake.calls.filter { $0.body["chat_id"] as? String == "-2002" }.isEmpty)
    check(bookingBot.booking(id)?.status == .pending)
}

await test("看不懂的 callback_data、編號不存在：回「無法處理」「找不到這筆預約」，不動任何預約") {
    let edits = fake.calls("editMessageText").count
    let pending = bookingBot.tickets.bookings.filter { $0.status == .pending }.map(\.id)
    for data in ["confirm", "approve:1", "confirm:abc", nil] {
        try await sendToBookingServer(callback(data, messageID: 1))
        check(callbackAnswer(nextCallbackID)?.body["text"] as? String == "無法處理", "\(String(describing: data))")
    }
    try await sendToBookingServer(callback("confirm:9999", messageID: 1))
    check(callbackAnswer(nextCallbackID)?.body["text"] as? String == "找不到這筆預約")
    check(fake.calls("editMessageText").count == edits)
    check(bookingBot.tickets.bookings.filter { $0.status == .pending }.map(\.id) == pending)
}


// 按下店員群組通知上的按鈕（chat 是按鈕所在的聊天室）
func callback(_ data: String, notice: Int64, chat: Int64 = staffChat, id: String = "cb-\(nextUpdateID)") -> [String: Any] {
    update(["callback_query": ["id": id, "from": ["id": 301, "is_bot": false, "first_name": "阿店"], "data": data,
                               "message": ["message_id": notice, "date": 0, "chat": ["id": chat, "type": "supergroup"]]]])
}

await test("轉人工：通知附原因、最近 5 則對話和按鈕；店員和顧客來回對話；按「結束轉人工」AI 恢復，重複按不重複處理") {
    let offTopic: [String: FieldValue] = ["intent": .text("off_topic"), "self_name": .text("")]
    let db = try BotDatabase()
    let bot = Bot(customers: SQLiteCustomerStore(db), tickets: SQLiteTicketStore(db), log: SQLiteConversationLog(db),
                  model: ScriptedModel(["Intent": [offTopic, offTopic, offTopic]]))
    guard let (s, p) = await startServer(bot, dedup: SQLiteUpdateDedup(db)) else { check(false, "伺服器沒有啟動"); return }
    func send(_ body: [String: Any]) async throws {
        let before = s.processedCount
        check(try await post(body, to: p) == 200)
        check(await waitUntil { s.processedCount > before }, "背景處理逾時")
    }
    let customer = CustomerKey(platform: "telegram", userID: "801")

    let before = notices().count
    for t in ["今天天氣如何", "你喜歡貓嗎", "講個笑話"] { try await send(privateText(801, t)) }
    let sent = Array(notices().dropFirst(before))
    check(sent.count == 1, "\(sent.map(\.body))")
    guard let notice = sent.first, let noticeID = notice.messageID else { return }
    let n = text(notice)
    check(n.contains("轉人工") && n.contains("連續 3 次"), n)
    // 最近 5 則：第 1 則顧客訊息之後的 5 則（含觸發轉人工的那句和轉人工回覆）
    check(!n.contains("今天天氣如何") && n.contains("👤 顧客：你喜歡貓嗎") && n.contains("👤 顧客：講個笑話"), n)
    // 舊的在前、同一輪顧客在前機器人在後：天氣的回覆、貓、貓的回覆、笑話、轉人工回覆
    let history = n.components(separatedBy: "\n").filter { $0.hasPrefix("👤 顧客：") || $0.hasPrefix("🤖 機器人：") }
    check(history.map { $0.hasPrefix("👤") ? "顧客" : "機器人" } == ["機器人", "顧客", "機器人", "顧客", "機器人"], n)
    check(history.filter { $0.hasPrefix("👤") } == ["👤 顧客：你喜歡貓嗎", "👤 顧客：講個笑話"], n)
    let keyboard = (notice.body["reply_markup"] as? [String: Any])?["inline_keyboard"] as? [[[String: Any]]]
    let ticketID = bot.tickets.tickets.last { $0.kind == .handoff }?.id ?? -1
    check(keyboard?.first?.first?["text"] as? String == "結束轉人工", "\(notice.body)")
    check(keyboard?.first?.first?["callback_data"] as? String == "endhuman:\(ticketID)", "\(notice.body)")

    // 轉人工中：顧客私訊不回、仍寫進紀錄；店員回覆通知可以來回對話
    let toCustomer = fake.calls("sendMessage", chat: "801").count
    try await send(privateText(801, "有人在嗎"))
    check(fake.calls("sendMessage", chat: "801").count == toCustomer, "轉人工中不該回覆")
    check(SQLiteConversationLog(db).entries(for: customer).contains { $0.direction == .customer && $0.text == "有人在嗎" })
    try await send(groupMessage(replyTo: noticeID, text: "我是店員，請說"))
    try await send(privateText(801, "想問包場"))
    try await send(groupMessage(replyTo: noticeID, text: "可以喔"))
    check(fake.calls("sendMessage", chat: "801").suffix(2).map(text) == ["💬 店員回覆：我是店員，請說", "💬 店員回覆：可以喔"],
          "\(fake.calls("sendMessage", chat: "801").map(\.body))")

    // 按下按鈕
    let group = fake.calls("sendMessage", chat: staffChatID).count
    try await send(callback("endhuman:\(ticketID)", notice: noticeID, id: "cb-end-1"))
    let answer = fake.calls("answerCallbackQuery").last
    check(answer?.body["callback_query_id"] as? String == "cb-end-1" && answer?.body["text"] as? String == "已恢復 AI 回覆", "\(answer?.body ?? [:])")
    let done = fake.calls("sendMessage", chat: staffChatID).dropFirst(group)
    check(done.count == 1 && text(done.first).contains("✅ 已結束轉人工") && text(done.first).contains("阿店") && replyTarget(done.first) == noticeID,
          "\(done.map(\.body))")
    check(!bot.customer(customer).human)
    try await send(privateText(801, "👍"))
    check(fake.calls("sendMessage", chat: "801").count == toCustomer + 3, "AI 沒有恢復回覆")

    // 重複按：只回應回呼，群組不再發訊息、狀態不變
    bot.update(customer) { $0.offTopicStreak = 1 }
    let groupAfter = fake.calls("sendMessage", chat: staffChatID).count
    try await send(callback("endhuman:\(ticketID)", notice: noticeID, id: "cb-end-2"))
    check(fake.calls("answerCallbackQuery").last?.body["text"] as? String == "已經結束轉人工了")
    check(fake.calls("sendMessage", chat: staffChatID).count == groupAfter)
    check(bot.customer(customer).offTopicStreak == 1)

    // 不存在的編號、看不懂的 data：只回應回呼
    try await send(callback("endhuman:9999", notice: noticeID, id: "cb-end-3"))
    check(fake.calls("answerCallbackQuery").last?.body["text"] as? String == "找不到這件事項")
    try await send(callback("confirm:1", notice: noticeID, id: "cb-end-4"))
    check(fake.calls("answerCallbackQuery").last?.body["callback_query_id"] as? String == "cb-end-4")
    check(fake.calls("sendMessage", chat: staffChatID).count == groupAfter)

    // 其他群組的按鈕回呼：不改狀態，只回應回呼讓按鈕停止轉圈（和預約按鈕 #13 的處理一致）
    bot.update(customer) { $0.human = true }
    let answers = fake.calls("answerCallbackQuery").count
    try await send(callback("endhuman:\(ticketID)", notice: noticeID, chat: -2002))
    check(fake.calls("answerCallbackQuery").count == answers + 1 && bot.customer(customer).human)
}

// MARK: - 日誌

await test("日誌裡沒有 token 和 secret") {
    let all = logs()
    check(!all.isEmpty)
    check(!all.contains { $0.contains(token) || $0.contains("SECRET-TOKEN") }, "日誌含 token")
    check(!all.contains { $0.contains(secret) }, "日誌含 secret")
    check(!all.contains { $0.contains("bot123456") || $0.contains("bot123456".addingPercentEncoding(withAllowedCharacters: .urlPathAllowed)!) },
          "日誌含 API 網址")
}

await test("正式執行的日誌會遮罩 token 和 secret（最後一道防線）") {
    let lines = Box<[String]>()
    let masked = Log(secrets: [token, secret], write: { lines.value = (lines.value ?? []) + [$0] })
    masked.error("url=https://example.invalid/bot\(token)/sendMessage secret=[REDACTED]")
    let line = lines.value?.first ?? ""
    check(!line.contains(token) && !line.contains(secret) && line.contains("bot***/sendMessage"), line)
}

await test("設定：缺環境變數時只說缺哪個，不印值") {
    do {
        _ = try TelegramConfig.fromEnvironment(["TELEGRAM_BOT_TOKEN": token])
        check(false, "應該丟出錯誤")
    } catch {
        let s = "\(error)"
        check(s.contains("TELEGRAM_WEBHOOK_SECRET") && s.contains("TELEGRAM_PUBLIC_URL") && !s.contains(token), s)
    }
    let c = try TelegramConfig.fromEnvironment(["TELEGRAM_BOT_TOKEN": token, "TELEGRAM_WEBHOOK_SECRET": secret,
                                                "TELEGRAM_PUBLIC_URL": "https://a.example/", "TELEGRAM_STAFF_CHAT_ID": "-100", "PORT": "9000"])
    check(c.webhookURL == "https://a.example/telegram/webhook" && c.staffChatID == -100 && c.port == 9000)
    check((try? TelegramConfig.fromEnvironment(["TELEGRAM_BOT_TOKEN": token, "TELEGRAM_WEBHOOK_SECRET": "有空白 的",
                                                "TELEGRAM_PUBLIC_URL": "https://a.example"])) == nil, "secret 只能用英數字、_、-")
}

print(failures == 0 ? "\n全部通過" : "\n失敗 \(failures) 項")
exit(failures == 0 ? 0 : 1)
