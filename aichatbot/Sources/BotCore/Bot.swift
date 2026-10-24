import Foundation

// MARK: - 對外介面：訊息進、動作出（LINE、終端機都走這裡）
// 顧客狀態依「平台＋使用者 ID」從 CustomerStore 讀出來、處理完存回去；待處理事項和預約申請存在 TicketStore；
// 每則訊息（顧客、機器人、店員）寫進 ConversationLogStore，機器人那筆附上這次的 AI 判斷。
// 同一位顧客的處理一定依序進行：receive 讀出狀態、等 AI 好幾秒、再整份存回，
// 店員操作（回覆、確認／拒絕預約、結束轉人工）要等這段做完才改，不然改的狀態會被蓋掉

@MainActor public final class Bot {
    public let customers: CustomerStore
    public let tickets: TicketStore
    public let log: ConversationLogStore
    public let model: ChatModel

    // model 預設是經過排隊的 Apple 本機模型；測試可以換成照腳本回答的 ScriptedModel
    public init(customers: CustomerStore? = nil, tickets: TicketStore? = nil, log: ConversationLogStore? = nil, model: ChatModel? = nil) {
        self.customers = customers ?? InMemoryCustomerStore()
        self.tickets = tickets ?? InMemoryTicketStore()
        self.log = log ?? InMemoryConversationLog()
        self.model = model ?? QueuedModel(base: AppleModel())
    }

    // MARK: 依顧客加鎖

    private var busy: Set<CustomerKey> = []
    private var waiters: [CustomerKey: [CheckedContinuation<Void, Never>]] = [:]

    // 同一位顧客一次只跑一段，後來的照順序排隊（不同顧客互不影響）。都在 MainActor 上，判斷和登記之間不會被插隊
    func serialized<T>(_ key: CustomerKey, _ body: () async throws -> T) async rethrows -> T {
        if busy.contains(key) {
            await withCheckedContinuation { waiters[key, default: []].append($0) }
        } else {
            busy.insert(key)
        }
        defer {
            // 有人在等就直接交給下一個（busy 不清，免得被新來的插隊）
            if var queue = waiters[key], !queue.isEmpty {
                let next = queue.removeFirst()
                waiters[key] = queue.isEmpty ? nil : queue
                next.resume()
            } else {
                busy.remove(key)
            }
        }
        return try await body()
    }

    // 回傳的順序：開待處理事項在前、回覆顧客在後；已轉人工時只回 [.silent]
    public func receive(_ m: IncomingMessage) async -> [BotAction] {
        await serialized(m.customer) {
            let tracing = TracingModel(base: model)
            let turn = Turn(customers.load(m.customer), key: m.customer, chatID: m.chatID, now: m.time, tickets: tickets, model: tracing)
            let (key, text) = if let media = m.media { turn.handleMedia(media, caption: m.text) } else { await turn.handle(m.text) }
            customers.save(m.customer, turn.customer)
            log.append(LogEntry(time: m.time, customer: m.customer, chatID: m.chatID, direction: .customer, text: turn.original))
            log.append(LogEntry(time: m.time, customer: m.customer, chatID: m.chatID, direction: .bot,
                                text: key == "silent" ? "" : text, replyKey: key, traces: tracing.traces))
            return turn.actions + [key == "silent" ? .silent : .reply(key: key, text: text, chatID: m.chatID)]
        }
    }

    // 只做分類（BOT_INTENT_EVAL 用）：和正式流程同一個 classify，狀態由呼叫的人給
    public func classify(_ text: String, as customer: Customer) async -> (intent: String, name: String?) {
        await Turn(customer, key: CustomerKey(platform: "eval", userID: "-"), chatID: "-", tickets: tickets, model: model).classify(text)
    }

    // MARK: 查詢（終端機、測試用）

    public var pendingTickets: [Ticket] { tickets.tickets.filter { !$0.done } }
    public func customer(_ key: CustomerKey) -> Customer { customers.load(key) }
    public func reset(_ key: CustomerKey) { customers.save(key, Customer()) }
    public func update(_ key: CustomerKey, _ change: (inout Customer) -> Void) {
        var c = customers.load(key); change(&c); customers.save(key, c)
    }

    // MARK: 店員群組：通知送出後記下群組訊息，店員「回覆」通知時用它找回待處理事項

    public func recordNotification(ticketID: Int, messageID: Int64) { tickets.setGroupMessage(ticketID: ticketID, messageID: messageID) }
    public func ticket(forGroupMessage messageID: Int64) -> Ticket? { tickets.ticket(groupMessageID: messageID) }

    // MARK: 店員操作（都進那位顧客的鎖，和顧客訊息依序處理）

    // 店員回覆：店員的原文（不經過 AI）加上顧客語言的前綴轉給顧客，這件事不再算「等店員」（已轉人工的狀態不變）。
    // 回覆緊急通知時對應回原本那張，兩張（和同一張的其他緊急通知）都改成已回覆；已回覆過的照樣轉送（店員可以補充好幾句）
    public func staffReply(ticketID: Int, text: String) async throws -> [BotAction] {
        guard let t = tickets.ticket(ticketID) else { throw StaffError.notFound(ticketID) }
        return await serialized(t.customer) {
            let root = t.relatedTo ?? t.id
            for var related in tickets.tickets where !related.done && ([t.id, root].contains(related.id) || related.relatedTo == root) {
                related.done = true; tickets.update(related)
            }
            update(t.customer) { $0.openTicket = false; $0.followUps = 0 }
            let reply = replyText("staff_reply_prefix", ["text": text], lang: customers.load(t.customer).lang)
            logStaff(t.customer, chatID: t.chatID, key: "staff_reply", text: reply)
            return [.reply(key: "staff_reply", text: reply, chatID: t.chatID)]
        }
    }

    func logStaff(_ customer: CustomerKey, chatID: String, key: String, text: String) {
        log.append(LogEntry(time: Date(), customer: customer, chatID: chatID, direction: .staff, text: text, replyKey: key))
    }

    // by：操作的店員（顯示名稱），存進預約；終端機不帶
    public func confirmBooking(_ id: Int, by staff: String? = nil) async throws -> [BotAction] {
        try await decide(id, .confirmed, reason: nil, by: staff)
    }
    public func rejectBooking(_ id: Int, reason: String? = nil, by staff: String? = nil) async throws -> [BotAction] {
        try await decide(id, .rejected, reason: reason, by: staff)
    }
    public func booking(_ id: Int) -> BookingRequest? { tickets.booking(id) }

    // 只有待確認的預約能處理，已確認或已拒絕的不重複處理；回覆用顧客自己的語言
    func decide(_ id: Int, _ status: BookingStatus, reason: String?, by staff: String?) async throws -> [BotAction] {
        guard let found = tickets.booking(id) else { throw StaffError.notFound(id) }
        return try await serialized(found.customer) {
            // 排隊時可能已經被處理了，進鎖後再讀一次
            guard var b = tickets.booking(id), b.status == .pending else { throw StaffError.alreadyHandled(id) }
            b.status = status; b.handledBy = staff; b.handledAt = Date(); tickets.update(b)
            if var t = tickets.ticket(id) { t.done = true; tickets.update(t) }
            let lang = customers.load(b.customer).lang
            let key = status == .confirmed ? "booking_confirmed" : "booking_rejected"
            var text = replyText(key, ["summary": b.summary], lang: lang)
            if let reason = reason?.trimmingCharacters(in: .whitespaces), !reason.isEmpty {
                text += replyText("booking_rejected_reason", ["reason": reason], lang: lang)
            }
            logStaff(b.customer, chatID: b.chatID, key: key, text: text)
            return [.reply(key: key, text: text, chatID: b.chatID)]
        }
    }

    // 結束轉人工：之後恢復 AI 回覆（不另外通知顧客）
    public func endHandoff(_ key: CustomerKey) async -> [BotAction] {
        await serialized(key) {
            finishHandoff(key)
            return []
        }
    }

    public enum EndHandoffResult: Equatable { case ended, alreadyEnded }

    // 店員群組按轉人工通知上的「結束轉人工」：顧客已經不在轉人工中（重複按、按到舊通知）就不改任何狀態
    public func endHandoff(ticketID: Int) async throws -> EndHandoffResult {
        guard let t = tickets.ticket(ticketID) else { throw StaffError.notFound(ticketID) }
        return await serialized(t.customer) {
            guard customers.load(t.customer).human else { return .alreadyEnded }
            finishHandoff(t.customer)
            logStaff(t.customer, chatID: t.chatID, key: "end_handoff", text: "（結束轉人工）")
            return .ended
        }
    }

    // 狀態歸零，這位顧客還沒完成的轉人工事項改成已完成
    func finishHandoff(_ key: CustomerKey) {
        update(key) { $0.human = false; $0.offTopicStreak = 0; $0.followUps = 0; $0.openTicket = false }
        for var t in tickets.tickets where t.customer == key && t.kind == .handoff && !t.done {
            t.done = true; tickets.update(t)
        }
    }

    // 這位顧客最近的對話（舊的在前，不含轉人工中 AI 沒回的空白紀錄）：轉人工通知附給店員看
    public func recentConversation(_ key: CustomerKey, limit: Int = 5) -> [LogEntry] { log.recent(key, limit: limit) }
}
