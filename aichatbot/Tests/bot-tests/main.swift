// 對話層測試：swift build && .build/debug/bot-tests（全部通過 exit 0，有失敗 exit 1）
// Command Line Tools 沒有 XCTest，Swift Testing 的巨集也載不到，所以寫成執行檔、用簡單的斷言。
// 不會用到本機模型：會呼叫 AI 的路徑用照腳本回答的假模型（ScriptedModel），沒有 Apple Intelligence 也能跑
import Foundation
import FoundationModels
import BotCore
import Storage

// 設定檔驗證的子程序：bot-tests --config-child，讀 BOT_CONFIG_DIR 的設定，印出設定檔資料夾和 FAQ 題數就結束（設定有錯會 exit 1）
if CommandLine.arguments.count == 2, CommandLine.arguments[1] == "--config-child" {
    _ = (store, replies, faqs, examples)
    print("\(dir.path) \(faqs.count)"); exit(0)
}

_ = (store, replies, faqs, examples)

// 跨程序排隊測試的子程序：bot-tests --queue-child <鎖檔> <日誌>，用同一把鎖做一次慢的 AI 呼叫就結束
if CommandLine.arguments.count == 4, CommandLine.arguments[1] == "--queue-child" {
    let q = AIQueue(lockPath: CommandLine.arguments[2], logPath: CommandLine.arguments[3])
    _ = await QueuedModel(base: ScriptedModel(delay: 0.5), queue: q).respond(ModelRequest(name: "Slow", text: "", instructions: "", fields: []))
    exit(0)
}

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
func throwsStaffError(_ op: () async throws -> Any) async -> Bool {
    do { _ = try await op(); return false } catch { return error is StaffError }
}

let alice = CustomerKey(platform: "terminal", userID: "alice")
let bob = CustomerKey(platform: "terminal", userID: "bob")

func message(_ key: CustomerKey, _ text: String) -> IncomingMessage {
    IncomingMessage(platform: key.platform, userID: key.userID, text: text)
}
func replyOf(_ actions: [BotAction]) -> (key: String, text: String, chatID: String)? {
    for case let .reply(key, text, chatID) in actions.reversed() { return (key, text, chatID) }
    return nil
}

// 準備一筆待確認的預約（和 bookingStep 完成時存的一樣：待處理事項和預約申請同一個編號）
@MainActor func pendingBooking(_ bot: Bot, for key: CustomerKey, summary: String = "2026-10-20 14:00, 4 位, 小康, 0912345678") -> Int {
    let t = bot.tickets.add(kind: .booking, customer: key, chatID: key.userID, text: "📅 預約申請：\(summary)")
    bot.tickets.add(BookingRequest(id: t.id, customer: key, chatID: key.userID, summary: summary, fields: [:]))
    return t.id
}

// MARK: - 店員操作

await test("店員回覆：原文加上前綴轉給顧客，清掉等店員的狀態，轉人工不變") {
    let bot = Bot()
    let t = bot.tickets.add(kind: .question, customer: alice, chatID: "chat-a", text: "📝 待補問題（知識庫沒有）：有吸菸區嗎")
    bot.update(alice) { $0.openTicket = true; $0.followUps = 2; $0.human = true }

    let actions = try await bot.staffReply(ticketID: t.id, text: "沒有喔")

    check(actions == [.reply(key: "staff_reply", text: "💬 店員回覆：沒有喔", chatID: "chat-a")], "\(actions)")
    let c = bot.customer(alice)
    check(!c.openTicket)
    check(c.followUps == 0)
    check(c.human)
    check(bot.pendingTickets.isEmpty)
}

await test("店員回覆：前綴用顧客的語言（中英日），店員的原文不動") {
    let bot = Bot()
    for (lang, prefix) in [("zh", "💬 店員回覆："), ("en", "💬 Reply from our staff: "), ("ja", "💬 スタッフより：")] {
        let key = CustomerKey(platform: "terminal", userID: lang)
        bot.update(key) { $0.lang = lang }
        let t = bot.tickets.add(kind: .question, customer: key, chatID: lang, text: "📝")
        let r = replyOf(try await bot.staffReply(ticketID: t.id, text: "OK，明天見 {name}"))
        check(r?.text == prefix + "OK，明天見 {name}", r?.text ?? "")
    }
}

await test("店員回覆：已回覆過的照樣轉送") {
    let bot = Bot()
    let t = bot.tickets.add(kind: .question, customer: alice, chatID: "alice", text: "📝")
    _ = try await bot.staffReply(ticketID: t.id, text: "第一句")
    check(replyOf(try await bot.staffReply(ticketID: t.id, text: "補充一句"))?.text.hasSuffix("補充一句") == true)
}

await test("店員回覆：編號不存在") {
    let bot = Bot()
    check(await throwsStaffError { try await bot.staffReply(ticketID: 99, text: "hi") })
}

await test("確認預約：狀態改成已確認，回覆含預約內容") {
    let bot = Bot()
    let id = pendingBooking(bot, for: alice)

    let actions = try await bot.confirmBooking(id)

    check(actions.count == 1, "\(actions)")
    let r = replyOf(actions)
    check(r?.key == "booking_confirmed")
    check(r?.chatID == "alice")
    check(r?.text.contains("2026-10-20 14:00") == true, r?.text ?? "")
    check(bot.tickets.bookings.first?.status == .confirmed)
    check(bot.pendingTickets.isEmpty)
}

await test("確認預約：用顧客自己的語言") {
    let bot = Bot()
    bot.update(bob) { $0.lang = "en" }
    let id = pendingBooking(bot, for: bob, summary: "2026-10-20 14:00, 4 people, Kevin, 0912345678")

    let r = replyOf(try await bot.confirmBooking(id))
    check(r?.text.contains("Kevin") == true, r?.text ?? "")
    check(r.map { detectLanguage($0.text) } == "en", r?.text ?? "")
}

await test("拒絕預約（附原因）：狀態改成已拒絕，回覆含原因") {
    let bot = Bot()
    let id = pendingBooking(bot, for: alice)

    let r = replyOf(try await bot.rejectBooking(id, reason: "當天包場"))
    check(r?.key == "booking_rejected")
    check(r?.text.contains("2026-10-20 14:00") == true, r?.text ?? "")
    check(r?.text.contains("當天包場") == true, r?.text ?? "")
    check(bot.tickets.bookings.first?.status == .rejected)
}

await test("預約只能處理一次") {
    let bot = Bot()
    let id = pendingBooking(bot, for: alice)
    _ = try await bot.rejectBooking(id)

    check(await throwsStaffError { try await bot.confirmBooking(id) })
    check(await throwsStaffError { try await bot.rejectBooking(id, reason: "再拒絕一次") })
    check(bot.tickets.bookings.first?.status == .rejected)
}

await test("確認或拒絕：編號不存在，或不是預約") {
    let bot = Bot()
    let t = bot.tickets.add(kind: .claim, customer: alice, chatID: "alice", text: "❓ 顧客聲稱／更正：上次說打五折")
    check(await throwsStaffError { try await bot.confirmBooking(t.id) })
    check(await throwsStaffError { try await bot.rejectBooking(42) })
}

await test("結束轉人工：恢復 AI 回覆") {
    let bot = Bot()
    bot.update(alice) { $0.human = true; $0.offTopicStreak = 3; $0.openTicket = true; $0.followUps = 1 }
    check(await bot.receive(message(alice, "哈囉")) == [.silent])

    check(await bot.endHandoff(alice).isEmpty)

    let c = bot.customer(alice)
    check(!c.human)
    check(c.offTopicStreak == 0)
    check(!c.openTicket)
    // 只有符號的訊息不呼叫 AI：恢復後會回覆，不再是 silent
    check(replyOf(await bot.receive(message(alice, "👍")))?.key == "small_talk")
}

let offTopic: [String: FieldValue] = ["intent": .text("off_topic"), "self_name": .text("")]

await test("假模型：轉人工 → AI 不回（訊息仍寫進紀錄）→ 結束轉人工 → AI 恢復") {
    let log = InMemoryConversationLog()
    let bot = Bot(log: log, model: ScriptedModel(["Intent": [offTopic, offTopic, offTopic]]))
    check(replyOf(await bot.receive(message(alice, "今天天氣如何")))?.key == "off_topic")
    check(replyOf(await bot.receive(message(alice, "你喜歡貓嗎")))?.key == "off_topic")
    let third = await bot.receive(message(alice, "講個笑話"))
    check(replyOf(third)?.key == "handoff", "\(third)")
    guard case let .openTicket(t)? = third.first else { check(false, "\(third)"); return }
    check(t.kind == .handoff && t.text.contains("轉人工") && t.text.contains("3 次"), t.text)
    check(bot.customer(alice).human)

    check(await bot.receive(message(alice, "有人在嗎")) == [.silent])
    check(log.entries.contains { $0.customer == alice && $0.direction == .customer && $0.text == "有人在嗎" }, "轉人工中的訊息沒寫進紀錄")

    check(try await bot.endHandoff(ticketID: t.id) == .ended)
    check(!bot.customer(alice).human && bot.customer(alice).offTopicStreak == 0)
    check(bot.tickets.ticket(t.id)?.done == true)
    check(log.entries.last?.direction == .staff && log.entries.last?.replyKey == "end_handoff")
    check(replyOf(await bot.receive(message(alice, "👍")))?.key == "small_talk")

    // 重複按、按到舊通知：不改狀態
    bot.update(alice) { $0.offTopicStreak = 2 }
    check(try await bot.endHandoff(ticketID: t.id) == .alreadyEnded)
    check(bot.customer(alice).offTopicStreak == 2)
    check(await throwsStaffError { try await bot.endHandoff(ticketID: 999) })
}

await test("最近對話：只取這位顧客最後 5 則，舊的在前，略過轉人工中不回覆的空白紀錄") {
    let db = try BotDatabase()
    for store in [InMemoryConversationLog(), SQLiteConversationLog(db)] as [ConversationLogStore] {
        let bot = Bot(log: store, model: ScriptedModel())
        for text in ["!", "!!", "!!!", "!!!!"] { _ = await bot.receive(message(alice, text)) }   // 只有符號：不呼叫 AI
        _ = await bot.receive(message(bob, "?"))
        bot.update(alice) { $0.human = true }
        _ = await bot.receive(message(alice, "?!"))
        let recent = bot.recentConversation(alice)
        check(recent.count == 5 && recent.allSatisfy { $0.customer == alice && !$0.text.isEmpty }, "\(recent)")
        // 每輪先寫顧客、再寫機器人；最後一輪轉人工中，只有顧客那句
        check(recent.map(\.direction) == [.customer, .bot, .customer, .bot, .customer], "\(recent.map(\.direction))")
        check(recent[0].text == "!!!" && recent[2].text == "!!!!" && recent.last?.text == "?!", "\(recent.map(\.text))")
        check(bot.recentConversation(bob, limit: 10).map(\.text).first == "?")
    }
}

await test("最近對話：時間完全相同時依寫入順序（同一輪顧客在前、機器人在後）") {
    let db = try BotDatabase()
    for store in [InMemoryConversationLog(), SQLiteConversationLog(db)] as [ConversationLogStore] {
        let now = Date()
        for i in 1...4 {
            store.append(LogEntry(time: now, customer: alice, chatID: "alice", direction: .customer, text: "問\(i)"))
            store.append(LogEntry(time: now, customer: alice, chatID: "alice", direction: .bot, text: "答\(i)", replyKey: "small_talk"))
        }
        let recent = store.recent(alice, limit: 5)
        check(recent.map(\.text) == ["答2", "問3", "答3", "問4", "答4"], "\(recent.map(\.text))")
    }
}

// MARK: - 店員群組：轉店員、回覆、追問、非文字訊息、依序處理

let followUp: [String: FieldValue] = ["intent": .text("follow_up"), "self_name": .text("")]

await test("假模型：轉店員 → 店員回覆 → 顧客收到原文，待處理事項改成已回覆") {
    let bot = Bot(model: ScriptedModel(["Intent": [["intent": .text("store_question"), "self_name": .text("")]]]))
    let actions = await bot.receive(message(alice, "有吸菸區嗎"))
    guard case let .openTicket(t)? = actions.first else { check(false, "\(actions)"); return }
    check(t.kind == .question && t.original == "有吸菸區嗎", "\(t)")
    check(bot.customer(alice).openTicket)

    let reply = try await bot.staffReply(ticketID: t.id, text: "沒有吸菸區喔")

    check(reply == [.reply(key: "staff_reply", text: "💬 店員回覆：沒有吸菸區喔", chatID: "alice")], "\(reply)")
    check(bot.tickets.ticket(t.id)?.done == true)
    check(!bot.customer(alice).openTicket && bot.customer(alice).followUps == 0)
}

await test("假模型：追問 3 次送緊急通知，引用原本那張；回覆緊急通知時兩張都改成已回覆") {
    let bot = Bot(model: ScriptedModel(["Intent": [followUp, followUp, followUp]]))
    let original = bot.tickets.add(kind: .question, customer: alice, chatID: "alice", original: "有吸菸區嗎", text: "📝 待補問題（知識庫沒有）：有吸菸區嗎")
    _ = bot.tickets.add(kind: .question, customer: bob, chatID: "bob", original: "bob 的問題", text: "📝 bob")
    bot.update(alice) { $0.openTicket = true }

    check(replyOf(await bot.receive(message(alice, "好了沒")))?.key == "waiting")
    check(replyOf(await bot.receive(message(alice, "好了沒")))?.key == "waiting")
    let third = await bot.receive(message(alice, "到底好了沒"))
    check(replyOf(third)?.key == "nudge")
    guard case let .openTicket(urgent)? = third.first else { check(false, "\(third)"); return }
    check(urgent.kind == .urgent && urgent.relatedTo == original.id, "\(urgent)")
    check(urgent.text.contains("有吸菸區嗎") && urgent.text.contains("#\(original.id)"), urgent.text)

    _ = try await bot.staffReply(ticketID: urgent.id, text: "沒有喔")
    check(bot.tickets.ticket(original.id)?.done == true && bot.tickets.ticket(urgent.id)?.done == true)
    check(bot.pendingTickets.map(\.customer) == [bob], "\(bot.pendingTickets)")
    check(!bot.customer(alice).openTicket && bot.customer(alice).followUps == 0)
}

await test("非文字訊息：圖片回「目前只看得懂文字」並開待處理事項，說明文字放進通知；不呼叫 AI") {
    let model = ScriptedModel()
    let bot = Bot(model: model)
    let actions = await bot.receive(IncomingMessage(platform: "telegram", userID: "U1", text: "這個還有貨嗎", media: "圖片"))
    check(actions.count == 2, "\(actions)")
    guard case let .openTicket(t)? = actions.first else { check(false, "\(actions)"); return }
    check(t.kind == .question && t.text.contains("📎") && t.text.contains("圖片") && t.text.contains("這個還有貨嗎"), t.text)
    check(replyOf(actions)?.key == "text_only")
    check(model.requests.isEmpty)
    check(bot.customer(CustomerKey(platform: "telegram", userID: "U1")).openTicket)

    bot.update(bob) { $0.lang = "en" }
    let en = await bot.receive(IncomingMessage(platform: bob.platform, userID: bob.userID, text: "", media: "語音"))
    check(replyOf(en)?.text == replyText("text_only", lang: "en"))
}

await test("非文字訊息：轉人工中不回覆，但仍開待處理事項") {
    let bot = Bot()
    bot.update(alice) { $0.human = true }
    let actions = await bot.receive(IncomingMessage(platform: alice.platform, userID: alice.userID, text: "", media: "檔案"))
    check(actions.count == 2 && actions.last == .silent, "\(actions)")
    guard case let .openTicket(t)? = actions.first else { check(false, "\(actions)"); return }
    check(t.text.contains("📎") && t.text.contains("檔案"), t.text)
}

await test("依序處理：顧客訊息正在等 AI 時店員回覆，兩邊的狀態修改都保留") {
    let model = ScriptedModel(["Intent": [["intent": .text("small_talk"), "self_name": .text("")]]], delay: 0.5)
    let bot = Bot(model: model)
    let t = bot.tickets.add(kind: .question, customer: alice, chatID: "alice", text: "📝")
    bot.update(alice) { $0.openTicket = true; $0.followUps = 2 }

    async let customerTurn = bot.receive(message(alice, "Sounds good, thank you so much"))
    // 等顧客那則讀出狀態、開始等 AI，店員才回覆
    while model.requests.isEmpty { try? await Task.sleep(for: .milliseconds(10)) }
    let staff = try await bot.staffReply(ticketID: t.id, text: "好的")
    let turn = await customerTurn

    check(replyOf(turn)?.key == "small_talk", "\(turn)")
    check(replyOf(staff)?.text.hasPrefix("💬 Reply from our staff") == true, "店員回覆應該等顧客那則處理完，用切換後的語言：\(staff)")
    let c = bot.customer(alice)
    check(c.lang == "en", "顧客那則的修改被蓋掉：\(c)")
    check(!c.openTicket && c.followUps == 0, "店員回覆的修改被蓋掉：\(c)")
}

await test("依序處理：不同顧客不用互相等") {
    let model = ScriptedModel(["Intent": [["intent": .text("small_talk"), "self_name": .text("")]]], delay: 0.5)
    let bot = Bot(model: model)
    async let slow = bot.receive(message(alice, "你好嗎今天"))
    while model.requests.isEmpty { try? await Task.sleep(for: .milliseconds(10)) }
    let start = Date()
    _ = await bot.endHandoff(bob)
    check(Date().timeIntervalSince(start) < 0.3, "Bob 等了 Alice")
    _ = await slow
}

// MARK: - 訊息進、動作出；顧客狀態分開

await test("不支援的語言：先開待處理事項，再回覆") {
    let bot = Bot()
    let actions = await bot.receive(IncomingMessage(platform: "line", userID: "U1", chatID: "C1", text: "안녕하세요"))

    check(actions.count == 2, "\(actions)")
    guard case let .openTicket(t)? = actions.first, let r = replyOf(actions) else { check(false, "\(actions)"); return }
    check(t.kind == .unsupported)
    check(t.customer == CustomerKey(platform: "line", userID: "U1"))
    check(t.chatID == "C1")
    check(t.text.contains("안녕하세요"))
    check(r.key == "unsupported")
    check(r.chatID == "C1")
    check(bot.pendingTickets == [t])
    check(bot.customer(t.customer).openTicket)
}

await test("兩位顧客交錯對話，狀態互不干擾") {
    let bot = Bot()
    _ = await bot.receive(message(alice, "안녕하세요"))
    let b1 = await bot.receive(message(bob, "!!!"))
    _ = await bot.receive(message(alice, "😀"))
    let b2 = await bot.receive(message(bob, "🙂"))

    check(bot.customer(alice).lang == "ko")
    check(bot.customer(alice).openTicket)
    check(bot.customer(bob).lang == "zh")
    check(!bot.customer(bob).openTicket)
    // Alice 的語言不影響 Bob：Bob 收到的都是一般回覆，沒有開待處理事項
    for actions in [b1, b2] {
        check(actions.count == 1, "\(actions)")
        check(replyOf(actions)?.key == "small_talk")
        check(replyOf(actions)?.chatID == "bob")
    }
    check(bot.pendingTickets.allSatisfy { $0.customer == alice })
}

await test("轉人工、店員回覆、結束轉人工只影響那位顧客") {
    let bot = Bot()
    bot.update(alice) { $0.human = true }
    bot.update(bob) { $0.openTicket = true; $0.followUps = 2 }
    let t = bot.tickets.add(kind: .handoff, customer: alice, chatID: "alice", text: "🙋 轉人工：連續 3 次無關或聲稱的訊息")

    check(await bot.receive(message(alice, "!!")) == [.silent])
    check(await bot.receive(message(bob, "!!")) != [.silent])

    _ = try await bot.staffReply(ticketID: t.id, text: "您好，我是店員")
    check(bot.customer(bob).openTicket)
    check(bot.customer(bob).followUps == 2)

    _ = await bot.endHandoff(alice)
    check(!bot.customer(alice).human)
    check(bot.customer(bob).openTicket)
}

await test("同一個使用者 ID 在不同平台是不同顧客") {
    let bot = Bot()
    _ = await bot.receive(IncomingMessage(platform: "line", userID: "same", text: "안녕하세요"))
    check(bot.customer(CustomerKey(platform: "line", userID: "same")).lang == "ko")
    check(bot.customer(CustomerKey(platform: "terminal", userID: "same")).lang == "zh")
}

// MARK: - 假模型：會呼叫 AI 的路徑

let monday = Calendar.current.date(from: DateComponents(year: 2026, month: 10, day: 19, hour: 12))!

await test("假模型：預約對話走到開出預約申請") {
    let model = ScriptedModel([
        "Intent": [["intent": .text("booking"), "self_name": .text("")],
                   ["intent": .text("booking"), "self_name": .text("")],
                   ["intent": .text("booking"), "self_name": .text("")]],
        "Booking": [[:],
                    ["date_quote": .text("10/21"), "date_month": .integer(10), "date_day": .integer(21), "time_quote": .text("下午2點"), "time_hour": .integer(14),
                     "people": .integer(4), "people_quote": .text("4位")],
                    ["name": .text("小康"), "phone": .text("0912345678")]],
    ])
    let bot = Bot(model: model)
    func say(_ text: String) async -> [BotAction] {
        await bot.receive(IncomingMessage(platform: "line", userID: "U1", text: text, time: monday))
    }

    check(replyOf(await say("我想訂位喔"))?.key == "ask_date")
    check(replyOf(await say("10/21 下午2點 4位"))?.key == "ask_name")
    let actions = await say("小康 0912345678")

    check(replyOf(actions)?.key == "booking_done", "\(actions)")
    guard case let .openTicket(t)? = actions.first else { check(false, "\(actions)"); return }
    check(t.kind == .booking)
    check(bot.tickets.bookings.first?.fields == ["date": "2026-10-21", "time": "14:00", "people": "4", "name": "小康", "phone": "0912345678"],
          "\(bot.tickets.bookings)")
    check(bot.customer(CustomerKey(platform: "line", userID: "U1")).askingField == nil)
    // 第一句沒抽到任何欄位：先挑 FAQ 確認不是在問預約規則（腳本沒給 Pick，挑不到就照常問日期）
    check(model.requests.map(\.name) == ["Intent", "Booking", "Pick", "Intent", "Booking", "Intent", "Booking"], "\(model.requests.map(\.name))")
    check(model.requests[1].text == "我想訂位喔")
}

// 用假模型走完預約對話（顧客 line:U1），回傳開出來的預約編號
@MainActor func completeBooking(_ bot: Bot) async -> Int? {
    for text in ["我想訂位喔", "10/21 下午2點 4位"] {
        _ = await bot.receive(IncomingMessage(platform: "line", userID: "U1", text: text, time: monday))
    }
    let actions = await bot.receive(IncomingMessage(platform: "line", userID: "U1", text: "小康 0912345678", time: monday))
    guard case let .openTicket(t)? = actions.first, t.kind == .booking else { check(false, "\(actions)"); return nil }
    return t.id
}
func bookingScript() -> ScriptedModel {
    ScriptedModel([
        "Intent": [["intent": .text("booking"), "self_name": .text("")],
                   ["intent": .text("booking"), "self_name": .text("")],
                   ["intent": .text("booking"), "self_name": .text("")]],
        "Booking": [[:],
                    ["date_quote": .text("10/21"), "date_month": .integer(10), "date_day": .integer(21), "time_quote": .text("下午2點"), "time_hour": .integer(14),
                     "people": .integer(4), "people_quote": .text("4位")],
                    ["name": .text("小康"), "phone": .text("0912345678")]],
    ])
}

await test("假模型：預約完成 → 店員確認：顧客收到確認，預約改成已確認、記下店員") {
    let bot = Bot(model: bookingScript())
    guard let id = await completeBooking(bot) else { return }
    check(bot.booking(id)?.status == .pending)

    let r = replyOf(try await bot.confirmBooking(id, by: "小明 @ming"))
    check(r?.key == "booking_confirmed" && r?.chatID == "U1", "\(String(describing: r))")
    check(r?.text.contains("小康") == true, r?.text ?? "")
    let b = bot.booking(id)
    check(b?.status == .confirmed && b?.handledBy == "小明 @ming" && b?.handledAt != nil, "\(String(describing: b))")
    check(bot.pendingTickets.isEmpty)
}

await test("假模型：預約完成 → 店員拒絕（不帶原因）：顧客收到「無法安排」，預約改成已拒絕") {
    let bot = Bot(model: bookingScript())
    guard let id = await completeBooking(bot) else { return }

    let r = replyOf(try await bot.rejectBooking(id, by: "小明"))
    check(r?.key == "booking_rejected" && r?.chatID == "U1", "\(String(describing: r))")
    check(r?.text == replyText("booking_rejected", ["summary": bot.booking(id)?.summary ?? ""], lang: "zh"), r?.text ?? "")
    check(bot.booking(id)?.status == .rejected && bot.booking(id)?.handledBy == "小明")
    check(await throwsStaffError { try await bot.confirmBooking(id, by: "小華") })
    check(bot.booking(id)?.handledBy == "小明", "重複按不該改掉處理的店員")
}

await test("短答判斷：名字、電話是短答，完整句子不是") {
    for s in ["Bill", "Bill Chen", "たなか", "김민수", "田中です", "0912-345-678", "4", "Bill 0912345678"] { check(isShortAnswer(s), s) }
    for s in ["What time do you close today?", "Can I change it to 3 people?", "営業時間を教えてください"] { check(!isShortAnswer(s), s) }
}

// #17：預約問欄位時回外文名字、電話，不切換語言；改說完整的外文句子才切換
@MainActor func midBooking(lang: String = "zh", asking: String, intent: String, booking: [String: FieldValue] = [:]) -> Bot {
    let bot = Bot(model: ScriptedModel(["Intent": [["intent": .text(intent), "self_name": .text("")]], "Booking": [booking]]))
    bot.update(alice) {
        $0.lang = lang; $0.askingField = asking
        $0.booking = ["date": "2026-10-21", "time": "14:00", "people": "4"]
        if asking == "phone" { $0.booking["name"] = "小康" }
    }
    return bot
}

await test("預約中回英文、日文、韓文名字：語言維持中文，不開待處理事項") {
    for name in ["Bill", "たなか", "김민수"] {
        let bot = midBooking(asking: "name", intent: "booking", booking: ["name": .text(name)])
        let r = replyOf(await bot.receive(message(alice, name)))
        check(bot.customer(alice).lang == "zh", "\(name)：\(bot.customer(alice).lang)")
        check(r?.key == "ask_phone", "\(name)：\(r?.key ?? "")")
        check(r.map { detectLanguage($0.text) } == "zh", r?.text ?? "")
        check(bot.pendingTickets.isEmpty, "\(name)：\(bot.pendingTickets)")
    }
}

await test("預約中回電話：語言維持中文") {
    let bot = midBooking(asking: "phone", intent: "booking", booking: ["phone": .text("0912-345-678")])
    let r = replyOf(await bot.receive(message(alice, "0912-345-678")))
    check(bot.customer(alice).lang == "zh")
    check(r?.key == "booking_done", r?.key ?? "")
    check(r.map { detectLanguage($0.text) } == "zh", r?.text ?? "")
}

await test("預約中改說完整英文句子：切換成英文") {
    let bot = midBooking(asking: "name", intent: "store_question")
    _ = await bot.receive(message(alice, "What time do you close today?"))
    check(bot.customer(alice).lang == "en", bot.customer(alice).lang)
}

await test("日文顧客預約中回漢字名字：維持日文") {
    let bot = midBooking(lang: "ja", asking: "name", intent: "booking", booking: ["name": .text("田中")])
    _ = await bot.receive(message(alice, "田中"))
    check(bot.customer(alice).lang == "ja", bot.customer(alice).lang)
}

await test("短答當名字：去掉標點和前後綴，有數字或是完整句子不算") {
    for (s, name) in [("たなか", "たなか"), ("Bill。", "Bill"), ("我是小康", "小康"), ("田中です", "田中"), ("Bill Chen", "Bill Chen")] {
        check(nameFromShortAnswer(s) == name, "\(s)：\(nameFromShortAnswer(s) ?? "nil")")
    }
    for s in ["0923456781", "Bill 0912345678", "What time do you close today?", "！！"] { check(nameFromShortAnswer(s) == nil, s) }
}

await test("預約中回假名名字，AI 改寫或沒抽到：收下整句當名字，接著問電話") {
    for aiName in ["田中", "Tanaka", ""] {
        let bot = midBooking(asking: "name", intent: "booking", booking: ["name": .text(aiName)])
        let r = replyOf(await bot.receive(message(alice, "たなか")))
        check(bot.customer(alice).booking["name"] == "たなか", "\(aiName)：\(bot.customer(alice).booking)")
        check(r?.key == "ask_phone", "\(aiName)：\(r?.key ?? "")")
        check(bot.customer(alice).lang == "zh", "\(aiName)：\(bot.customer(alice).lang)")
    }
}

await test("預約中問姓名時回電話或完整句子：不會被收成名字") {
    for text in ["0923456781", "What time do you close today?"] {
        let bot = midBooking(asking: "name", intent: "booking")
        _ = await bot.receive(message(alice, text))
        check(bot.customer(alice).booking["name"] == nil, "\(text)：\(bot.customer(alice).booking)")
    }
}

await test("預約中回假名名字，分成自我介紹、AI 沒抽到：一樣收下，接著問電話") {
    let bot = midBooking(asking: "name", intent: "self_intro", booking: ["name": .text("")])
    let r = replyOf(await bot.receive(message(alice, "たなか")))
    check(bot.customer(alice).booking["name"] == "たなか", "\(bot.customer(alice).booking)")
    check(r?.key == "ask_phone", r?.key ?? "")
}

// #17 第二版驗收：保底只在分類成預約、自我介紹時用，整句漢字也不用（「不知道」「好」曾被收成名字）
await test("預約中問姓名時回「不知道」「好」：不會被收成名字") {
    for (text, intent, want) in [("不知道", "unclear", "unclear"), ("不知道", "booking", "ask_name"), ("好", "small_talk", "small_talk"),
                                 ("Okay", "small_talk", "small_talk")] {
        let bot = midBooking(asking: "name", intent: intent, booking: ["name": .text("")])
        let r = replyOf(await bot.receive(message(alice, text)))
        check(bot.customer(alice).booking["name"] == nil, "\(text)/\(intent)：\(bot.customer(alice).booking)")
        check(r?.key == want, "\(text)/\(intent)：\(r?.key ?? "")")
        check(bot.customer(alice).askingField == "name", "\(text)/\(intent)：\(bot.customer(alice).askingField ?? "nil")")
    }
}

await test("預約中回英文名字，AI 有抽到：不管分類都收下") {
    let bot = midBooking(asking: "name", intent: "small_talk", booking: ["name": .text("Bill")])
    let r = replyOf(await bot.receive(message(alice, "Bill")))
    check(bot.customer(alice).booking["name"] == "Bill", "\(bot.customer(alice).booking)")
    check(r?.key == "ask_phone", r?.key ?? "")
}

await test("沒有在預約時，短的外文照舊切換語言") {
    let bot = Bot(model: ScriptedModel(["Intent": [["intent": .text("small_talk"), "self_name": .text("")]]]))
    _ = await bot.receive(message(alice, "Thanks"))
    check(bot.customer(alice).lang == "en", bot.customer(alice).lang)
}

await test("假模型：AI 失敗時用預設值（分類當成看不懂）") {
    let model = ScriptedModel()
    let bot = Bot(model: model)
    check(replyOf(await bot.receive(message(alice, "你們營業到幾點呀")))?.key == "unclear")
    check(model.requests.map(\.name) == ["Intent"])
    check(bot.pendingTickets.isEmpty)
}

await test("假模型：挑 FAQ 和判斷是否發問都失敗時，轉店員") {
    let model = ScriptedModel(["Intent": [["intent": .text("store_question"), "self_name": .text("")]]])
    let bot = Bot(model: model)
    let actions = await bot.receive(message(alice, "你們營業到幾點呀"))

    check(replyOf(actions)?.key == "fallback", "\(actions)")
    check(model.requests.map(\.name) == ["Intent", "Pick", "Kind"], "\(model.requests.map(\.name))")
    guard case let .openTicket(t)? = actions.first else { check(false, "\(actions)"); return }
    check(t.kind == .question)
}

// MARK: - 相對日期：程式依今天換算

let gregorian = Calendar(identifier: .gregorian)
func day(_ m: Int, _ d: Int) -> Date { gregorian.date(from: DateComponents(year: 2026, month: m, day: d, hour: 12))! }
func read(_ quote: String, _ today: Date) -> String {
    switch readDate(quote, today: today) {
    case .date(let d): let f = DateFormatter(); f.dateFormat = "M/d"; return f.string(from: d)
    case .unresolvedRelative: return "換算不了"
    case .notRelative: return "不是相對"
    }
}

await test("相對日期：今天是週一（10/5）") {
    for (q, want) in [("今天", "10/5"), ("明天", "10/6"), ("後天", "10/7"), ("tomorrow", "10/6"), ("明日", "10/6"), ("明後日", "10/7"),
                      ("the day after tomorrow", "10/7"), ("這週五", "10/9"), ("這禮拜日", "10/11"), ("this Friday", "10/9"),
                      ("今週の金曜日", "10/9"), ("下週二", "10/13"), ("下禮拜一", "10/12"), ("next Tuesday", "10/13"),
                      ("来週の火曜日", "10/13"), ("週二", "10/6"), ("星期三", "10/7"), ("禮拜天", "10/11"), ("Wednesday", "10/7"),
                      ("水曜日", "10/7"), ("週一", "10/5")] {
        check(read(q, day(10, 5)) == want, "\(q)：\(read(q, day(10, 5)))，應為 \(want)")
    }
}

await test("相對日期：今天是週三（10/7）") {
    for (q, want) in [("今天", "10/7"), ("明天", "10/8"), ("後天", "10/9"), ("Tomorrow", "10/8"), ("這週五", "10/9"), ("這週一", "10/5"),
                      ("this Friday", "10/9"), ("下週二", "10/13"), ("下週二下午兩點", "10/13"), ("下禮拜一", "10/12"), ("下星期三", "10/14"),
                      ("next Tuesday", "10/13"), ("来週の火曜日", "10/13"), ("星期三", "10/7"), ("週二", "10/13"), ("禮拜天", "10/11"),
                      ("Wednesday", "10/7"), ("水曜日", "10/7"), ("tue", "10/13")] {
        check(read(q, day(10, 7)) == want, "\(q)：\(read(q, day(10, 7)))，應為 \(want)")
    }
}

await test("相對日期：今天是週日（10/11）") {
    for (q, want) in [("明天", "10/12"), ("後天", "10/13"), ("tomorrow", "10/12"), ("這週日", "10/11"), ("這週一", "10/5"),
                      ("下週一", "10/12"), ("下週二", "10/13"), ("next Sunday", "10/18"), ("来週の日曜日", "10/18"),
                      ("週一", "10/12"), ("星期天", "10/11"), ("日曜日", "10/11"), ("Saturday", "10/17")] {
        check(read(q, day(10, 11)) == want, "\(q)：\(read(q, day(10, 11)))，應為 \(want)")
    }
}

await test("相對日期：換算不了的說法、絕對日期") {
    for q in ["下下週二", "大後天", "週末", "這週末", "this weekend", "next week", "下週", "下個月", "再過三天", "3天後", "in 3 days", "来月"] {
        check(read(q, day(10, 7)) == "換算不了", "\(q)：\(read(q, day(10, 7)))")
    }
    for q in ["10/21", "10月21日", "10/21（三）", "10月21日星期三", "晚上"] {
        check(read(q, day(10, 7)) == "不是相對", "\(q)：\(read(q, day(10, 7)))")
    }
    for q in ["下下禮拜三", "下下星期三", "下下周三", "下 下週三", "我想訂下下週三下午兩點"] {
        check(read(q, day(10, 7)) == "換算不了", "\(q)：\(read(q, day(10, 7)))")
    }
}

await test("出處是換算不了的說法的一部分") {
    check(partOfUnresolved("下週三", in: "我想訂下下週三下午兩點"))
    check(partOfUnresolved("週三", in: "我想訂下下週三下午兩點"))
    check(partOfUnresolved("next Wednesday", in: "the week after next Wednesday"))
    check(!partOfUnresolved("明天", in: "明天晚上7點，我們晚點再過去"))
    check(!partOfUnresolved("下週三", in: "我想訂下週三下午兩點"))
}

await test("明確寫出的月日由程式讀") {
    func md(_ s: String) -> String { explicitMonthDays(in: s).map { "\($0.month)/\($0.day)" }.joined(separator: ",") }
    for (s, want) in [("我想訂 10/19 下午兩點", "10/19"), ("10月19日", "10/19"), ("10月19號晚上", "10/19"), ("10/21（三）", "10/21"),
                      ("10-19", "10/19"), ("10.19", "10/19"), ("10/19 14:00", "10/19"), ("2026/10/19", "10/19"),
                      ("10/19 或 10/20，10/19 優先", "10/19,10/20"), ("電話 0912-345-678", ""), ("14:00", ""), ("4-6人", ""),
                      ("下午兩點", ""), ("13/40", "")] {
        check(md(s) == want, "\(s)：\(md(s))，應為 \(want)")
    }
}

await test("絕對日期補年份：今年；今天往前 7 天之前算明年") {
    check(absoluteDate(month: 10, day: 21, today: day(10, 7)) == "2026-10-21")
    check(absoluteDate(month: 10, day: 1, today: day(10, 7)) == "2026-10-01")   // 7 天內：維持今年，回「已經過了」
    check(absoluteDate(month: 9, day: 1, today: day(10, 7)) == "2027-09-01")
    check(absoluteDate(month: 1, day: 5, today: day(10, 7)) == "2027-01-05")
}

// 今天 2026-10-07（週三）一句話完成日期、時間、人數，回傳回覆種類和存下的預約欄位
@MainActor func bookOnce(_ text: String, _ booking: [String: FieldValue]) async -> (String?, [String: String]) {
    let bot = Bot(model: ScriptedModel(["Intent": [["intent": .text("booking"), "self_name": .text("")]], "Booking": [booking]]))
    let actions = await bot.receive(IncomingMessage(platform: "line", userID: "U1", text: text, time: day(10, 7)))
    return (replyOf(actions)?.key, bot.customer(CustomerKey(platform: "line", userID: "U1")).booking)
}

await test("假模型：「下週二」用程式換算的日期，不用 AI 填的月日") {
    let (key, b) = await bookOnce("我想預約下週二下午2點 4位", ["date_quote": .text("下週二"), "date_month": .integer(10), "date_day": .integer(10),
                                                         "time_quote": .text("下午2點"), "time_hour": .integer(14), "people": .integer(4), "people_quote": .text("4位")])
    check(key == "ask_name", key ?? "")
    check(b == ["date": "2026-10-13", "time": "14:00", "people": "4"], "\(b)")
}

await test("假模型：換算不了的相對說法不填日期，改問日期") {
    let (key, b) = await bookOnce("下下週三下午2點", ["date_quote": .text("下下週三"), "date_month": .integer(10), "date_day": .integer(21),
                                                 "time_quote": .text("下午2點"), "time_hour": .integer(14)])
    check(key == "ask_date", key ?? "")
    check(b["date"] == nil && b["time"] == "14:00", "\(b)")
}

await test("假模型：AI 只抄到「下週三」，原句是「下下週三」，改問日期") {
    for q in ["下週三", "週三"] {
        let (key, b) = await bookOnce("我想訂下下週三下午兩點", ["date_quote": .text(q), "date_month": .integer(10), "date_day": .integer(14),
                                                         "time_quote": .text("下午兩點"), "time_hour": .integer(14)])
        check(key == "ask_date", "\(q)：\(key ?? "")")
        check(b == ["time": "14:00"], "\(q)：\(b)")
    }
}

await test("假模型：寫了月日就由程式讀，AI 填的月日不用") {
    // 10/19 是週一公休：AI 把月日填成 0、或日期和時間出處重疊，日期都要收到，回公休
    for (dq, tq) in [("10/19", "下午兩點"), ("10/19 下午兩點", "下午兩點"), ("10/19 下午兩點", "10/19 下午兩點"), ("", "")] {
        let (key, b) = await bookOnce("我想訂 10/19 下午兩點", ["date_quote": .text(dq), "date_month": .integer(0), "date_day": .integer(0),
                                                          "time_quote": .text(tq), "time_hour": .integer(14)])
        check(key == "err_monday", "\(dq)｜\(tq)：\(key ?? "")")
        check(b == ["time": "14:00"] || (tq.isEmpty && b.isEmpty), "\(dq)｜\(tq)：\(b)")
    }
    var (key, b) = await bookOnce("我想預約 10/20 下午兩點 4位", ["date_quote": .text("10/20 下午兩點"), "date_month": .integer(10), "date_day": .integer(20),
                                                          "time_quote": .text("下午兩點"), "time_hour": .integer(14), "people": .integer(4), "people_quote": .text("4位")])
    check(key == "ask_name", key ?? "")
    check(b == ["date": "2026-10-20", "time": "14:00", "people": "4"], "\(b)")

    // AI 填錯的日也不用
    (key, b) = await bookOnce("10/21 下午2點 4位", ["date_quote": .text("10/21"), "date_month": .integer(10), "date_day": .integer(12),
                                                  "time_quote": .text("下午2點"), "time_hour": .integer(14), "people": .integer(4), "people_quote": .text("4位")])
    check(key == "ask_name", key ?? "")
    check(b == ["date": "2026-10-21", "time": "14:00", "people": "4"], "\(b)")

    // 寫了兩個日期：AI 填的要是其中之一
    (key, b) = await bookOnce("10/21 或 10/22 下午2點", ["date_quote": .text("10/21 或 10/22"), "date_month": .integer(10), "date_day": .integer(12),
                                                     "time_quote": .text("下午2點"), "time_hour": .integer(14)])
    check(key == "ask_date" && b["date"] == nil, "\(key ?? "")：\(b)")
    (key, b) = await bookOnce("10/21 或 10/22 下午2點", ["date_quote": .text("10/21 或 10/22"), "date_month": .integer(10), "date_day": .integer(22),
                                                     "time_quote": .text("下午2點"), "time_hour": .integer(14)])
    check(b["date"] == "2026-10-22", "\(key ?? "")：\(b)")
}

await test("假模型：絕對日期、時間要和原文出處的數字一致") {
    var (key, b) = await bookOnce("21號 下午2點 4位", ["date_quote": .text("21號"), "date_month": .integer(10), "date_day": .integer(12),
                                                    "time_quote": .text("下午2點"), "time_hour": .integer(14), "people": .integer(4), "people_quote": .text("4位")])
    check(key == "ask_date", key ?? "")
    check(b == ["time": "14:00", "people": "4"], "\(b)")

    (key, b) = await bookOnce("10/21 下午2點 4位", ["date_quote": .text("10/21"), "date_month": .integer(10), "date_day": .integer(21),
                                                  "time_quote": .text("下午2點"), "time_hour": .integer(18), "people": .integer(4), "people_quote": .text("4位")])
    check(key == "ask_time", key ?? "")
    check(b == ["date": "2026-10-21", "people": "4"], "\(b)")

    (key, b) = await bookOnce("明天兩點半 4位", ["date_quote": .text("明天"), "time_quote": .text("兩點半"), "time_hour": .integer(14), "time_minute": .integer(30),
                                             "people": .integer(4), "people_quote": .text("4位")])
    check(key == "ask_name", key ?? "")
    check(b == ["date": "2026-10-08", "time": "14:30", "people": "4"], "\(b)")
}

// MARK: - AI 排隊

let scratch = FileManager.default.temporaryDirectory.appendingPathComponent("bot-tests-\(getpid())")
try? FileManager.default.createDirectory(at: scratch, withIntermediateDirectories: true)
let slow = ModelRequest(name: "Slow", text: "", instructions: "", fields: [])

// 日誌每行：程序編號、開始、結束；依開始時間排序後，每一筆都要在前一筆結束後才開始
func readAILog(_ url: URL) -> [(pid: String, start: Double, end: Double)] {
    ((try? String(contentsOf: url, encoding: .utf8)) ?? "").split(separator: "\n").map { line in
        let p = line.split(separator: "\t")
        return (String(p[0]), Double(p[1])!, Double(p[2])!)
    }.sorted { $0.start < $1.start }
}
func overlaps(_ calls: [(pid: String, start: Double, end: Double)]) -> Int {
    zip(calls, calls.dropFirst()).filter { $0.end > $1.start }.count
}

await test("排隊：同一個程序同時兩個呼叫，時間不重疊、記下等鎖時間") {
    let log = scratch.appendingPathComponent("in-process.log")
    let q = AIQueue(lockPath: scratch.appendingPathComponent("in-process.lock").path, logPath: log.path)
    let model = QueuedModel(base: ScriptedModel(delay: 0.2), queue: q)
    async let a = model.respond(slow)
    async let b = model.respond(slow)
    _ = await (a, b)

    let calls = readAILog(log)
    check(calls.count == 2, "\(calls)")
    check(overlaps(calls) == 0, "\(calls)")
    check(q.callCount == 2)
    check(q.lockWaitSeconds >= 0.15, "\(q.lockWaitSeconds)")
}

await test("排隊：兩個程序搶同一把鎖，時間不重疊") {
    let lock = scratch.appendingPathComponent("cross.lock").path, log = scratch.appendingPathComponent("cross.log")
    let child = Process()
    child.executableURL = URL(fileURLWithPath: CommandLine.arguments[0])
    child.arguments = ["--queue-child", lock, log.path]
    try child.run()
    let q = AIQueue(lockPath: lock, logPath: log.path)
    let model = QueuedModel(base: ScriptedModel(delay: 0.3), queue: q)
    for _ in 1...2 { _ = await model.respond(slow) }
    child.waitUntilExit()

    let calls = readAILog(log)
    check(child.terminationStatus == 0)
    check(calls.count == 3, "\(calls)")
    check(Set(calls.map(\.pid)).count == 2, "\(calls)")
    check(overlaps(calls) == 0, "\(calls)")
}

// MARK: - AI 回答快取

let asked = ModelRequest(name: "Intent", text: "請問營業到幾點秘密句", instructions: "你是茶屋客服秘密指示",
                         fields: [.choice("intent", "分類", ["a", "b"]), .text("self_name", "名字")])
let answer: [String: FieldValue] = ["intent": .text("a"), "self_name": .text("")]
var cacheDirs = 0
func cacheDir() -> URL { cacheDirs += 1; return scratch.appendingPathComponent("ai-cache-\(cacheDirs)") }
func cacheFiles(_ dir: URL) -> [URL] {
    (FileManager.default.enumerator(at: dir, includingPropertiesForKeys: nil)?.allObjects as? [URL] ?? [])
        .filter { $0.pathExtension == "json" }
}

await test("快取：同一個請求第二次命中，不再呼叫模型") {
    let base = ScriptedModel(["Intent": [answer, answer]])
    let cache = CachedModel(base: base, dir: cacheDir(), modelVersion: "m1", systemVersion: "s1")
    let a = await cache.respond(asked), b = await cache.respond(asked)
    check(base.requests.count == 1, "\(base.requests.count)")
    check(cache.hits == 1 && cache.misses == 1, cache.summary)
    check(a?.values == answer && b?.values == answer, "\(String(describing: a?.values)) \(String(describing: b?.values))")
    check(cache.summary == "AI 快取：命中 1 次、實際呼叫 1 次", cache.summary)
}

await test("快取：指示、輸入、欄位不同就不命中") {
    let base = ScriptedModel(["Intent": Array(repeating: answer, count: 4)])
    let cache = CachedModel(base: base, dir: cacheDir(), modelVersion: "m1", systemVersion: "s1")
    _ = await cache.respond(asked)
    _ = await cache.respond(ModelRequest(name: asked.name, text: asked.text, instructions: asked.instructions + "改", fields: asked.fields))
    _ = await cache.respond(ModelRequest(name: asked.name, text: asked.text + "改", instructions: asked.instructions, fields: asked.fields))
    _ = await cache.respond(ModelRequest(name: asked.name, text: asked.text, instructions: asked.instructions,
                                         fields: [.choice("intent", "分類", ["a", "c"]), .text("self_name", "名字")]))
    check(base.requests.count == 4 && cache.hits == 0, cache.summary)
}

await test("快取：模型版本或系統版本不同就不命中") {
    let dir = cacheDir()
    let base = ScriptedModel(["Intent": Array(repeating: answer, count: 3)])
    for (m, s) in [("m1", "s1"), ("m2", "s1"), ("m1", "s2")] {
        _ = await CachedModel(base: base, dir: dir, modelVersion: m, systemVersion: s).respond(asked)
    }
    check(base.requests.count == 3, "\(base.requests.count)")
    check(CachedModel.key(asked, modelVersion: "m1", systemVersion: "s1") != CachedModel.key(asked, modelVersion: "m1", systemVersion: "s2"))
    let again = CachedModel(base: base, dir: dir, modelVersion: "m2", systemVersion: "s1")
    _ = await again.respond(asked)
    check(again.hits == 1 && base.requests.count == 3, again.summary)
}

await test("快取：AI 失敗或逾時（nil）不存，下次重新呼叫") {
    let dir = cacheDir()
    let base = ScriptedModel(["Intent": [answer]])   // 第一次之後腳本用完，回傳 nil
    let failing = CachedModel(base: ScriptedModel(), dir: dir, modelVersion: "m1", systemVersion: "s1")
    check(await failing.respond(asked) == nil)
    check(cacheFiles(dir).isEmpty, "\(cacheFiles(dir))")
    let cache = CachedModel(base: base, dir: dir, modelVersion: "m1", systemVersion: "s1")
    check(await cache.respond(asked)?.values == answer)
    check(cache.misses == 1 && base.requests.count == 1, cache.summary)
}

await test("快取：檔名和內容都不含明文輸入") {
    let dir = cacheDir()
    _ = await CachedModel(base: ScriptedModel(["Intent": [answer]]), dir: dir, modelVersion: "m1", systemVersion: "s1").respond(asked)
    let files = cacheFiles(dir)
    check(files.count == 1, "\(files)")
    for f in files {
        let text = (try? String(contentsOf: f, encoding: .utf8)) ?? ""
        for secret in ["秘密句", "秘密指示", "營業", "茶屋"] {
            check(!text.contains(secret) && !f.path.contains(secret), "\(f.path)：\(text)")
        }
    }
}

await test("快取：dev、holdout 分開存，互不命中") {
    let root = cacheDir()
    let base = ScriptedModel(["Intent": [answer, answer]])
    for ns in ["dev", "holdout"] {
        let m = CachedModel.fromEnvironment(base: base, vars: ["BOT_AI_CACHE": "1", "BOT_AI_CACHE_NS": ns, "BOT_AI_CACHE_DIR": root.path])
        _ = await m.respond(asked)
        check(cacheFiles(root.appendingPathComponent(ns)).count == 1, ns)
    }
    check(base.requests.count == 2, "\(base.requests.count)")
}

await test("快取：沒設 BOT_AI_CACHE=1 就不包快取") {
    let base = ScriptedModel()
    check(!(CachedModel.fromEnvironment(base: base, vars: [:]) is CachedModel))
    check(!(CachedModel.fromEnvironment(base: base, vars: ["BOT_AI_CACHE": "0"]) is CachedModel))
    check(CachedModel.fromEnvironment(base: base, vars: ["BOT_AI_CACHE": "1", "BOT_AI_CACHE_DIR": cacheDir().path]) is CachedModel)
}

// MARK: - 模型服務：Go 送來的 JSON schema 轉回欄位定義

await test("模型服務：JSON schema 轉回欄位，順序照 required（properties 是物件，解析後順序會亂）") {
    let json = #"""
    {"type":"object","properties":{"people":{"type":"integer","description":"人數"},
     "intent":{"type":"string","description":"意圖","enum":["booking","other"]},"name":{"type":"string","description":"名字"}},
     "required":["intent","name","people"],"additionalProperties":false}
    """#
    let fields = fieldsFromJSONSchema(try JSONSerialization.jsonObject(with: Data(json.utf8)))
    check(fields?.count == 3, "\(String(describing: fields))")
    if let f = fields, f.count == 3 {
        if case let .choice(n, d, o) = f[0] { check(n == "intent" && d == "意圖" && o == ["booking", "other"]) } else { check(false, "\(f[0])") }
        if case let .text(n, d) = f[1] { check(n == "name" && d == "名字") } else { check(false, "\(f[1])") }
        if case let .integer(n, d) = f[2] { check(n == "people" && d == "人數") } else { check(false, "\(f[2])") }
    }
    check(fieldsFromJSONSchema(["type": "object", "properties": [:], "required": []]).map(\.isEmpty) == true)
    // 不支援的型別、required 列了沒定義的欄位：整個不收
    check(fieldsFromJSONSchema(["properties": ["x": ["type": "boolean"]], "required": ["x"]]) == nil)
    check(fieldsFromJSONSchema(["properties": [:], "required": ["x"]]) == nil)
}

// MARK: - 有逾時的呼叫：模型服務失敗時要在 log 分得出是逾時還是本機模型丟出錯誤（#59）

await test("逾時呼叫：超過秒數算逾時") {
    let result = await runWithTimeout(0.05) { try await Task.sleep(for: .seconds(1)); return 1 }
    if case .timedOut = result {} else { check(false, "\(result)") }
}

struct FakeModelError: Error, CustomStringConvertible { var description: String { "內容被擋下" } }
await test("逾時呼叫：丟出錯誤算錯誤，留下型別和訊息") {
    let result = await runWithTimeout(5) { () async throws -> Int in throw FakeModelError() }
    if case let .failed(type, message) = result {
        check(type == "FakeModelError" && message == "內容被擋下", "\(type) \(message)")
    } else { check(false, "\(result)") }
}

await test("逾時呼叫：超過本機模型上下文的錯誤標成 contextExceededType（模型服務回 502 帶 reason，#65）") {
    let result = await runWithTimeout(5) { () async throws -> Int in
        throw LanguageModelSession.GenerationError.exceededContextWindowSize(.init(debugDescription: "8192"))
    }
    if case let .failed(type, _) = result {
        check(type == contextExceededType, type)
    } else { check(false, "\(result)") }
}

await test("逾時呼叫：安全過濾擋下的錯誤標成 guardrailType（模型服務回 502 帶 reason guardrail，#68）") {
    let result = await runWithTimeout(5) { () async throws -> Int in
        throw LanguageModelSession.GenerationError.guardrailViolation(.init(debugDescription: "May contain unsafe content"))
    }
    if case let .failed(type, _) = result {
        check(type == guardrailType, type)
    } else { check(false, "\(result)") }
    if #available(macOS 27, *) {
        let newer = await runWithTimeout(5) { () async throws -> Int in
            throw LanguageModelError.guardrailViolation(.init(debugDescription: "May contain unsafe content"))
        }
        if case let .failed(type, _) = newer {
            check(type == guardrailType, type)
        } else { check(false, "\(newer)") }
    }
}

await test("本機模型：只改寫既有內容的翻譯用寬鬆的安全過濾，其他請求照舊（#68）") {
    check(AppleModel.usesPermissiveGuardrails("Translate"))
    check(AppleModel.usesPermissiveGuardrails("TranslateCustomer"))
    for name in ["Intent", "FAQ", "Ask", "Pick", "Answer", "Booking"] {
        check(!AppleModel.usesPermissiveGuardrails(name), name)
    }
}

await test("逾時呼叫：時間內回答就拿到結果") {
    let result = await runWithTimeout(5) { 42 }
    if case .value(42) = result {} else { check(false, "\(result)") }
}

// MARK: - 資料庫

await test("資料庫：同一個檔案開兩次核心，預約問到一半能接著完成") {
    let path = scratch.appendingPathComponent("resume.sqlite").path
    let u1 = CustomerKey(platform: "telegram", userID: "U1")
    func say(_ bot: Bot, _ text: String) async -> [BotAction] {
        await bot.receive(IncomingMessage(platform: u1.platform, userID: u1.userID, text: text, time: monday))
    }

    // 第一個核心：問到名字就「關掉」
    var halfway = Customer()
    do {
        let db = try BotDatabase(path: path)
        let bot = Bot(customers: SQLiteCustomerStore(db), log: SQLiteConversationLog(db), model: ScriptedModel([
            "Intent": [["intent": .text("booking"), "self_name": .text("")], ["intent": .text("booking"), "self_name": .text("")]],
            "Booking": [[:], ["date_quote": .text("10/21"), "date_month": .integer(10), "date_day": .integer(21), "time_quote": .text("下午2點"), "time_hour": .integer(14),
                              "people": .integer(4), "people_quote": .text("4位")]],
        ]))
        check(replyOf(await say(bot, "我想訂位喔"))?.key == "ask_date")
        check(replyOf(await say(bot, "10/21 下午2點 4位"))?.key == "ask_name")
        halfway = bot.customer(u1)
        check(halfway.askingField != nil, "\(halfway)")
    }

    // 第二個核心：同一個檔案重開，接著回答名字和電話
    let db = try BotDatabase(path: path)
    let bot = Bot(customers: SQLiteCustomerStore(db), log: SQLiteConversationLog(db), model: ScriptedModel([
        "Intent": [["intent": .text("booking"), "self_name": .text("")]],
        "Booking": [["name": .text("小康"), "phone": .text("0912345678")]],
    ]))
    check(bot.customer(u1) == halfway, "\(bot.customer(u1))")
    let actions = await say(bot, "小康 0912345678")
    check(replyOf(actions)?.key == "booking_done", "\(actions)")
    check(bot.tickets.bookings.first?.fields == ["date": "2026-10-21", "time": "14:00", "people": "4", "name": "小康", "phone": "0912345678"],
          "\(bot.tickets.bookings)")
    check(bot.customer(u1).askingField == nil)
}

await test("資料庫：對話紀錄有顧客、機器人、店員三個方向，機器人那筆附上 AI 判斷") {
    let db = try BotDatabase()
    let log = SQLiteConversationLog(db)
    let bot = Bot(customers: SQLiteCustomerStore(db), log: log,
                  model: ScriptedModel(["Intent": [["intent": .text("store_question"), "self_name": .text("")]]]))
    let actions = await bot.receive(IncomingMessage(platform: "telegram", userID: "U2", chatID: "C2", text: "你們營業到幾點呀"))
    guard case let .openTicket(t)? = actions.first else { check(false, "\(actions)"); return }
    _ = try await bot.staffReply(ticketID: t.id, text: "營業到晚上九點")

    let entries = log.entries(for: CustomerKey(platform: "telegram", userID: "U2"))
    check(entries.map(\.direction) == [.customer, .bot, .staff], "\(entries.map(\.direction))")
    guard entries.count == 3 else { return }
    check(entries[0].text == "你們營業到幾點呀" && entries[0].chatID == "C2" && entries[0].traces.isEmpty)
    check(entries[1].replyKey == "fallback" && entries[1].text == replyOf(actions)?.text)
    let traces = entries[1].traces
    check(traces.map(\.name) == ["Intent", "Pick", "Kind"], "\(traces.map(\.name))")
    check(traces.first?.text == "你們營業到幾點呀")
    check(traces.first?.values["intent"] == .text("store_question"), "\(traces.first?.values ?? [:])")
    check(traces.first?.failed == false)
    // 腳本沒寫到的題目 = AI 失敗
    check(traces.dropFirst().allSatisfy { $0.failed && $0.values.isEmpty }, "\(traces)")
    check(traces.allSatisfy { $0.ms >= 0 })
    check(entries[2].text == "💬 店員回覆：營業到晚上九點" && entries[2].replyKey == "staff_reply" && entries[2].chatID == "C2")
}

await test("資料庫：轉人工中不回覆也記一筆（silent），確認預約記成店員") {
    let db = try BotDatabase()
    let log = SQLiteConversationLog(db)
    let bot = Bot(customers: SQLiteCustomerStore(db), log: log)
    bot.update(alice) { $0.human = true }
    check(await bot.receive(message(alice, "哈囉")) == [.silent])
    _ = try await bot.confirmBooking(pendingBooking(bot, for: alice))

    let entries = log.entries(for: alice)
    check(entries.map(\.direction) == [.customer, .bot, .staff], "\(entries.map(\.direction))")
    guard entries.count == 3 else { return }
    check(entries.dropFirst().map { $0.replyKey ?? "" } == ["silent", "booking_confirmed"], "\(entries)")
    check(entries[1].text.isEmpty && entries[1].traces.isEmpty)
}

await test("資料庫：清掉超過 90 天的對話紀錄、超過 7 天的 update 紀錄") {
    let db = try BotDatabase()
    let log = SQLiteConversationLog(db), dedup = SQLiteUpdateDedup(db)
    let now = Date(), day = 86400.0
    for days in [91.0, 89] {
        log.append(LogEntry(time: now.addingTimeInterval(-days * day), customer: alice, chatID: "alice", direction: .customer, text: "\(Int(days)) 天前"))
    }
    check(dedup.firstTime(8, now: now.addingTimeInterval(-8 * day)))
    check(dedup.firstTime(6, now: now.addingTimeInterval(-6 * day)))

    let n = try db.purge(now: now)

    check(n.logs == 1 && n.updates == 1, "\(n)")
    check(log.entries().map(\.text) == ["89 天前"], "\(log.entries().map(\.text))")
    check(!dedup.firstTime(6), "6 天前的 update 應該還記得")
    check(dedup.firstTime(8), "8 天前的 update 應該已經清掉")
}

await test("資料庫：記憶體資料庫（終端機版）可以存顧客狀態，不會建檔") {
    let db = try BotDatabase()
    let bot = Bot(customers: SQLiteCustomerStore(db), log: SQLiteConversationLog(db))
    _ = await bot.receive(message(bob, "안녕하세요"))
    check(bot.customer(bob).lang == "ko")
    check(SQLiteCustomerStore(db).load(bob).openTicket)
    check(SQLiteCustomerStore(db).load(alice) == Customer())
}

await test("資料庫：待處理事項存進資料庫，用群組訊息找得回來，重開後讀得回來") {
    let path = scratch.appendingPathComponent("tickets.sqlite").path
    do {
        let db = try BotDatabase(path: path)
        let bot = Bot(customers: SQLiteCustomerStore(db), tickets: SQLiteTicketStore(db), log: SQLiteConversationLog(db),
                      model: ScriptedModel(["Intent": [["intent": .text("store_question"), "self_name": .text("")]]]))
        let actions = await bot.receive(message(alice, "有吸菸區嗎"))
        guard case let .openTicket(t)? = actions.first else { check(false, "\(actions)"); return }
        let urgent = bot.tickets.add(kind: .urgent, customer: alice, chatID: "alice", original: "好了沒", text: "🔥", relatedTo: t.id)
        bot.recordNotification(ticketID: t.id, messageID: 9001)
        bot.recordNotification(ticketID: urgent.id, messageID: 9002)
        check(bot.ticket(forGroupMessage: 9001) == t)
        check(bot.ticket(forGroupMessage: 1234) == nil)
        check(bot.pendingTickets.count == 2)
    }

    let db = try BotDatabase(path: path)
    let bot = Bot(customers: SQLiteCustomerStore(db), tickets: SQLiteTicketStore(db), log: SQLiteConversationLog(db))
    guard let urgent = bot.ticket(forGroupMessage: 9002) else { check(false, "重開後找不到"); return }
    check(urgent.kind == .urgent && urgent.original == "好了沒" && urgent.relatedTo != nil, "\(urgent)")
    check(bot.ticket(forGroupMessage: 9001)?.original == "有吸菸區嗎")
    _ = try await bot.staffReply(ticketID: urgent.id, text: "沒有喔")
    check(SQLiteTicketStore(db).tickets.allSatisfy(\.done), "\(SQLiteTicketStore(db).tickets)")
    check(!bot.customer(alice).openTicket)
}

await test("資料庫：預約申請存進資料庫，確認後狀態和店員重開後讀得回來") {
    let path = scratch.appendingPathComponent("bookings.sqlite").path
    var id = 0
    do {
        let db = try BotDatabase(path: path)
        let bot = Bot(customers: SQLiteCustomerStore(db), tickets: SQLiteTicketStore(db), log: SQLiteConversationLog(db), model: bookingScript())
        guard let booked = await completeBooking(bot) else { return }
        id = booked
    }

    let db = try BotDatabase(path: path)
    let bot = Bot(customers: SQLiteCustomerStore(db), tickets: SQLiteTicketStore(db), log: SQLiteConversationLog(db))
    guard let b = bot.booking(id) else { check(false, "重開後找不到預約 #\(id)"); return }
    check(b.fields == ["date": "2026-10-21", "time": "14:00", "people": "4", "name": "小康", "phone": "0912345678"], "\(b.fields)")
    check(b.status == .pending && b.customer == CustomerKey(platform: "line", userID: "U1") && b.chatID == "U1" && b.summary.contains("小康"), "\(b)")
    check(bot.tickets.bookings.map(\.id) == [id])

    check(replyOf(try await bot.confirmBooking(id, by: "小明"))?.key == "booking_confirmed")
    let saved = SQLiteTicketStore(try BotDatabase(path: path)).booking(id)
    check(saved?.status == .confirmed && saved?.handledBy == "小明" && saved?.handledAt != nil, "\(String(describing: saved))")
    check(await throwsStaffError { try await bot.rejectBooking(id) })
    check(await throwsStaffError { try await bot.confirmBooking(id + 100) })
}

await test("資料庫：顧客狀態少了欄位（舊資料）用預設值") {
    let c = try JSONDecoder().decode(Customer.self, from: Data(#"{"lang":"en","askingField":"name"}"#.utf8))
    check(c.lang == "en" && c.askingField == "name" && !c.openTicket && c.booking.isEmpty && c.followUps == 0)
}

// MARK: - 設定檔：位置（BOT_CONFIG_DIR）、啟動驗證、回覆模板

await test("回覆模板：程式用到的模板在 replies.json 都有") {
    check(missingReplyTemplates(replies).isEmpty, "缺少：\(missingReplyTemplates(replies))")
}

await test("回覆模板：原始碼裡寫死的模板名稱都在 usedReplyKeys 清單裡") {
    let sources = URL(fileURLWithPath: #filePath).deletingLastPathComponent().deletingLastPathComponent()
        .deletingLastPathComponent().appendingPathComponent("Sources")
    let files = (FileManager.default.enumerator(at: sources, includingPropertiesForKeys: nil)?.allObjects as? [URL] ?? [])
        .filter { $0.pathExtension == "swift" }
    // R("…")、replyText("…")（組出來的 R("ask_" + 欄位) 不算），以及條件式或回傳值裡的模板名稱（"calc_card_ok"、"err_date"、"booking_confirmed"）
    let patterns = [#"\b(?:R|replyText)\("([a-z_]+)"\s*[,)]"#, #""((?:err|calc|booking)_[a-z_]+)""#].map { try! NSRegularExpression(pattern: $0) }
    var found: Set<String> = []
    for file in files {
        let text = try String(contentsOf: file, encoding: .utf8)
        for p in patterns {
            for m in p.matches(in: text, range: NSRange(text.startIndex..., in: text)) {
                found.insert(String(text[Range(m.range(at: 1), in: text)!]))
            }
        }
    }
    check(found.count >= 20, "只找到 \(found.count) 個，掃描的路徑可能不對：\(sources.path)")
    check(found.subtracting(usedReplyKeys).isEmpty, "沒列在 usedReplyKeys：\(found.subtracting(usedReplyKeys).sorted())")
}

// 把目前的四個設定檔複製到暫存資料夾，再依測試改壞其中一個
func configCopy(_ name: String, edit: [String: (String) -> String] = [:]) throws -> URL {
    let d = scratch.appendingPathComponent("config-\(name)")
    try? FileManager.default.removeItem(at: d)
    try FileManager.default.createDirectory(at: d, withIntermediateDirectories: true)
    for f in ["store.json", "replies.json", "faq.txt", "examples.jsonl"] {
        let text = try String(contentsOf: dir.appendingPathComponent(f), encoding: .utf8)
        try (edit[f]?(text) ?? text).write(to: d.appendingPathComponent(f), atomically: true, encoding: .utf8)
    }
    return d
}
func configErrors(_ d: URL) -> String? {
    do { _ = try loadConfig(from: d); return nil } catch { return "\(error)" }
}

await test("設定檔：目前的設定沒有錯") {
    let d = try configCopy("ok")
    check(configErrors(d) == nil, configErrors(d) ?? "")
    let c = try loadConfig(from: d)
    check(c.faqs.count == faqs.count && c.examples.count == examples.count && c.replies == replies)
}

await test("設定檔：store.json 少一個逗號，指出檔案和 JSON 格式錯誤") {
    let d = try configCopy("comma", edit: ["store.json": { $0.replacingOccurrences(of: #""openHour""#, with: #""x": 1 "openHour""#) }])
    let e = configErrors(d) ?? ""
    check(e.contains("store.json") && e.contains("JSON 格式錯誤"), e)
}

await test("設定檔：store.json 缺欄位或型別不對，指出欄位") {
    let missing = try configCopy("missing", edit: ["store.json": { $0.replacingOccurrences(of: #""maxPeople""#, with: #""maxPeopleX""#) }])
    check(configErrors(missing)?.contains("store.json：缺少欄位 maxPeople") == true, configErrors(missing) ?? "沒有錯誤")
    let type = try configCopy("type", edit: ["store.json": { $0.replacingOccurrences(of: #""price": 120,"#, with: #""price": "120","#) }])
    let e = configErrors(type) ?? ""
    check(e.contains("store.json：欄位 menu[0].price 的型別不對"), e)
}

await test("設定檔：replies.json 模板沒有 zh、值不是字串，或少了程式用到的模板") {
    let d = try configCopy("replies", edit: ["replies.json": {
        $0.replacingOccurrences(of: #""fallback": {"zh": "#, with: #""fallback": {"zz": "#)
          .replacingOccurrences(of: #""nudge": {"zh": "#, with: #""nudge": {"zh": 1, "x": "#)
    }])
    let e = configErrors(d) ?? ""
    check(e.contains("replies.json") && e.contains("模板 fallback 缺少 zh"), e)
    check(e.contains("模板 nudge 的 zh 不是字串"), e)
    let missing = try configCopy("replies-missing", edit: ["replies.json": { $0.replacingOccurrences(of: #""greeting": {"#, with: #""greetingX": {"#) }])
    let m = configErrors(missing) ?? ""
    check(m.contains("replies.json：缺少程式用到的回覆模板：greeting"), m)
}

await test("設定檔：faq.txt、examples.jsonl 格式不對的行，指出檔名和行號（註解、空行略過）") {
    let d = try configCopy("lines", edit: [
        "faq.txt": { "# 註解\n\n只有 | 兩欄\n" + $0 },
        "examples.jsonl": { "\n不是 JSON\n{\"text\": \"你好\"}\n" + $0 },
    ])
    let e = configErrors(d) ?? ""
    check(e.contains("faq.txt：第 3 行"), e)
    check(e.contains("examples.jsonl：第 2 行") && e.contains("第 3 行：缺少字串欄位"), e)
}

await test("設定檔：BOT_CONFIG_DIR 指定位置；設定有錯時印出錯誤、exit 1，不當掉") {
    func run(_ d: URL) throws -> (Int32, String) {
        let child = Process(), out = Pipe()
        child.executableURL = URL(fileURLWithPath: CommandLine.arguments[0])
        child.arguments = ["--config-child"]
        child.environment = ProcessInfo.processInfo.environment.merging(["BOT_CONFIG_DIR": d.path]) { $1 }
        child.standardOutput = out; child.standardError = out
        try child.run(); child.waitUntilExit()
        return (child.terminationStatus, String(data: out.fileHandleForReading.readDataToEndOfFile(), encoding: .utf8) ?? "")
    }
    let ok = try configCopy("env", edit: ["faq.txt": { $0 + "\nextra | 額外的問法 | 額外的答案\n" }])
    let (status, output) = try run(ok)
    check(status == 0 && output.contains(ok.path) && output.contains(" \(faqs.count + 1)"), "\(status) \(output)")

    let bad = try configCopy("env-bad", edit: ["store.json": { $0.replacingOccurrences(of: #""maxPeople""#, with: #""maxPeopleX""#) }])
    let (badStatus, badOutput) = try run(bad)
    check(badStatus == 1 && badOutput.contains("設定檔有錯") && badOutput.contains("store.json：缺少欄位 maxPeople"), "\(badStatus) \(badOutput)")
}
// MARK: - 相似案例改用預先算好的向量（BOT_EMBED_FILE，Go 的 e5 向量，技術驗證用）

await test("BOT_EMBED_FILE：相似案例從向量檔找，不限同一種文字系統；檔案裡沒有的訊息不附案例") {
    let d = scratch.appendingPathComponent("embed")
    try FileManager.default.createDirectory(at: d, withIntermediateDirectories: true)
    try "{\"text\": \"你好\", \"intent\": \"greeting\"}\n{\"text\": \"再見\", \"intent\": \"goodbye\"}\n{\"text\": \"Bye\", \"intent\": \"goodbye\"}\n"
        .write(to: d.appendingPathComponent("examples.jsonl"), atomically: true, encoding: .utf8)
    let file = d.appendingPathComponent("vectors.json")
    try #"{"こんにちは": [1, 0.1], "你好": [0.9, 0.2], "再見": [0.1, 1], "Bye": [0, 1]}"#.write(to: file, atomically: true, encoding: .utf8)
    examples = try loadExamples(from: d)
    fileVectors = try loadEmbedFile(file)
    check(similarExamples("こんにちは", k: 2).map(\.text) == ["你好", "再見"], "\(similarExamples("こんにちは", k: 2).map(\.text))")
    check(similarExamples("沒算過的句子").isEmpty)
    fileVectors = nil; examples = initialConfig.examples
}

await test("BOT_EMBED_FILE：檔案不是「句子 → 數字陣列」時丟出錯誤") {
    let file = scratch.appendingPathComponent("bad-vectors.json")
    try #"{"你好": "x"}"#.write(to: file, atomically: true, encoding: .utf8)
    check((try? loadEmbedFile(file)) == nil)
}
try? FileManager.default.removeItem(at: scratch)

print(failures == 0 ? "\n全部通過" : "\n失敗 \(failures) 項")
exit(failures == 0 ? 0 : 1)
